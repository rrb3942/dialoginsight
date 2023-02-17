# What is dialoginsight? #
dialoginsight allows you to export dialog profiles from [OpenSIPs](https://opensips.org) as [Prometheus](https://prometheus.io/) style metrics.

## Features ##
* Supports OpenSIPs 3.2+ via the [mi_http module](https://opensips.org/docs/modules/3.2.x/mi_http.html) and JSON-RPC 2.0
* Export all dialog profiles or only a specific list
* Ability to use dynamic metric labels on the exported metrics from the OpenSIPs script.

# Configuration #
Configuration is done via a json config file (typically /etc/dialogsight/config.json) or command line flags.

## Configuration Settings ##
These may be used as command line flags or as fields in the json configuration.
* `listen` - Local IP and port for the prometheus exporter to listen on. (default "127.0.0.1:10337")
* `opensips_mi` - URL to the mi_http instance for opensips. (default "http://127.0.0.1:8888/mi")
* `export_all` - Whether or not to export all dialog profiles from the instance. (default "true")
* `export_profiles` - List of dialog profiles to export. Used if export_all is set to false.
* `insight_label` - Dialog value starting prefix to indicate it is an insight value (contains labels to process). (default "insight")
* `timeout` - Timeout duration for opensips API requests. (default "2s")
* `idle_remove` - If a metric is idle for this long it will be removed from memory. (default "1m")

## Additional Flags ##
* `-config` - Allows specifying the configuration file to use.

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

## Insight ##
Adding insight to the exported metrics allows you to increase visibility into active calls.

### Insight Values ###
Insight values are profile values that start with `insight_label` followed by a `:` and a series of `label=value` pairs separated by `;`

`label` must follow Prometheus conventions and match the pattern of `'^[a-zA-Z_][a-zA-Z0-9_]*'`
* If a `label` is not valid the `label=value` pair will be silently dropped

`value` may be any valid UTF-8 character, except `';'`

Example format:

	insight: cust=1234;carrier=4567;some_stat=5

The format is very similar to a SIP header whos body contains only a list of parameters with values.

### Adding Insight ###
You can can add insight to existing dialog profiles, or use separate profiles.

Lets take our customer example from above

	set_dlg_profile("customer", "1234")

Exports as:

	dialoginsight_exported_profile_customer_dialogs{value="1234"} 1

Lets add insight to track which source IPs the customer is sending from:

	set_dlg_profile("customer", "1234")
	set_dlg_profile("customer", "insight: cust=1234;scr_ip=$si")
These will be exported as:

	dialoginsight_exported_profile_customer_dialogs{value="1234"} 1
	dialoginsight_profile_customer_dialogs{cust="1234",src_ip="192.168.168.1"} 1

### Notes on Cardinality ###
`dialogsight` makes it very easy to create detailed views of your active call dialogs. However one pitfall that you should be aware of is high [cardinality](https://grafana.com/blog/2022/02/15/what-are-cardinality-spikes-and-why-do-they-matter/). High cardinality can cause performance issues, both with exporting metrics and querying them. Therefore it is important to be mindful of how many potential metrics you may be exporting.

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
Our potential unique combinations is 20,000 ((100\*100) + (100\*50) + (100\*50)). You lose granularity in exchange for reduced cardinality.