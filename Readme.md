# What is dialoginsight? #
[dialoginsight](https://github.com/rrb3942/dialoginsight) allows you to export dynamic and detailed dialog profiles from [OpenSIPs](https://opensips.org) as [Prometheus](https://prometheus.io/) style metrics.

### OpenSIPs Statistics ###
This application does not export standard OpenSIPs statistics. If are looking to export OpenSIPs statistics it has been supported directly by the [prometheus](https://opensips.org/docs/modules/3.2.x/prometheus.html) module since version 3.2. For older versions you can use [opensips_exporter](https://github.com/VoIPGRID/opensips_exporter).
## Features ##
* Supports OpenSIPs 3.2+ via the [mi_http](https://opensips.org/docs/modules/3.2.x/mi_http.html) module and JSON-RPC 2.0
* Export all dialog profiles or only a specific list
* Ability to use dynamic metric labels on exported metrics from the OpenSIPs script.

# Configuration #
Configuration is done via a json config file (typically /etc/dialoginsight/config.json) or command line flags.

## Configuration Settings ##
These may be used as command line flags or as fields in the json configuration.
* `listen` - Local IP and port for the prometheus exporter to listen on. Metrics are exported under http://listen/metrics (default "127.0.0.1:10337")
* `opensips_mi` - URL to the mi_http instance for OpenSIPs. (default "http://127.0.0.1:8888/mi")
* `export_all` - Whether or not to export all dialog profiles from the instance. (default "true")
* `export_profiles` - List of dialog profiles to export. Used if export_all is set to false.
* `replication_hints` - Provide a mapping from OpenSIPs reported name to shared/replicated tagged names. (Example: { "sharedprofile": [ "sharedprofile/s", "sharedprofile/b" ] })
* `insight_label` - Dialog value starting prefix to indicate it is an insight value (contains labels to process). (default "insight")
* `timeout` - Timeout duration for OpenSIPs API requests. (default "2s")
* `idle_remove` - If a metric is idle for this long it will be removed from memory. (default "1m")
* `enable_profiling` - Enables access to profiling via http://listen/debug/pprof/ (default "false")

## Additional Flags ##
* `config` - Allows specifying the configuration file to use.

# Exporting Dialog Profiles #
## Export Namespaces ##
Depending on the type of profile they will be exported under either `dialoginsight_exported_profile_` (standard OpenSIPs profiles) or `dialoginsight_profile_` (extended insight profiles) metric prefixes.

Each profile is exported as its own metric of type `dialogs`

For example, lets says you have a profile named `global` that all calls are a part of.

	set_dlg_profile("global")
This will be exported as:

	dialoginsight_exported_profile_global_dialogs 1

## Profiles with Values ##
Dialog profiles with values will be exported with a label of name 'value' containing the profile value. For example:

	set_dlg_profile("customer", "1234")
Will be exported as:

	dialoginsight_exported_profile_customer_dialogs{value="1234"} 1
## Shared or Replicated Profiles ##
Currently OpenSIPs reports the name of shared and/or replicated profiles without their tags, but the tags must be included in API calls for them to work correctly. To work around this the `replication_hints` configuration option is provided to allow remapping the reported name to names with the correct tags.

When a profile indicates that it it shared or replicated we will automatically include the `shared` or `replicated` labels in the metric with the value of `yes`. If these labels overlap with `insight` labels, the `insight` labels will overwrite them.

## Insight ##
Adding insight to the exported profiles allows you to set dynamic metric labels and increase visibility into active calls.

### Insight Values ###
Insight values are profile values that start with `insight_label` followed by a `:` and a series of `label=value` pairs separated by `;`

`label` must follow Prometheus conventions and match the pattern of `'^[a-zA-Z_][a-zA-Z0-9_]*'`
* If a `label` is not valid the `label=value` pair will be silently dropped

`value` may contain any valid UTF-8 character, except `';'`

Example format:

	insight: cust=1234;carrier=4567;some_stat=5

The format is very similar to a SIP header whose body contains only a list of parameters with values.

### Adding Insight ###
You can can add insight to existing dialog profiles, or use separate profiles.

Lets take our customer example from above

	set_dlg_profile("customer", "1234")

Exports as:

	dialoginsight_exported_profile_customer_dialogs{value="1234"} 1

Lets add insight to track which source IP the customer is sending from:

	set_dlg_profile("customer", "1234")
	set_dlg_profile("customer", "insight: cust=1234;scr_ip=$si")
These will be exported as:

	dialoginsight_exported_profile_customer_dialogs{value="1234"} 1
	dialoginsight_profile_customer_dialogs{cust="1234",src_ip="192.168.168.1"} 1

### Notes on Cardinality ###
`dialoginsight` makes it very easy to create detailed views of your active call dialogs. However one pitfall that you should be aware of is high [cardinality](https://grafana.com/blog/2022/02/15/what-are-cardinality-spikes-and-why-do-they-matter/). High cardinality can cause performance issues, both with exporting metrics and querying them. Therefore it is important to be mindful of how many potential metrics you may be exporting.

Lets take the example of:

	set_dlg_profile("tracking", "insight: cust=1234;carrier=4567")
If you have 100 customers, and 100 carriers, this can potentially result in 10,000 (100\*100) unique combinations. While that is not bad by itself, lets say we decide to also track the destination state for U.S. Domestic calls:

	set_dlg_profile("tracking", "insight: cust=1234;carrier=4567;called_state=NY")
We now have a potential for 500,000 (100*100\*50) unique combinations. This shows how quickly cardinality can grow.

#### Ways to reduce cardinality ####
One way to limit cardinality is to use less combinations of labels in a single metric. If we take the example from above, and instead track the `called_state` by `cust` and `carrier` individually:

	set_dlg_profile("tracking", "insight: cust=1234;carrier=4567")
	set_dlg_profile("cust_called_state", "insight: cust=1234;called_state=NY")
	set_dlg_profile("carrier_called_state", "insight: carrier=4567;called_state=NY")
Our potential unique combinations is greatly reduced to 20,000 ((100\*100) + (100\*50) + (100\*50)). You lose granularity in exchange for reduced cardinality.

## Issues or Pull Requests ##
https://github.com/rrb3942/dialoginsight/issues

# License #
MIT

***Copyright (c) 2023 Ryan Bullock***
