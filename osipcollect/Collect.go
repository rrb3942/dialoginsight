package osipcollect

import (
	"dialoginsight/metrics"
	"log"
	"maps"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	exportNamespace  = "dialoginsight_exported_profile_"
	insightNamespace = "dialoginsight_profile_"
	yes              = "yes"
)

// Implement Promeheus Collector Interface
// Collects and provides metrics.
func (osip *Client) Collect(ch chan<- prometheus.Metric) {
	// Only allow one collect to run at a time
	osip.mu.Lock()
	defer osip.mu.Unlock()

	osip.exportProfiles.Reset()
	osip.exportValueProfiles.Reset()
	osip.insightProfiles.Reset()

	if err := osip.setProfiles(); err != nil {
		log.Println(err)
		return
	}

	osip.exportProfiles.Collect(ch)
	osip.exportValueProfiles.Collect(ch)
	osip.insightProfiles.Collect(ch)
}

func (osip *Client) setProfiles() error {
	// API Call to opensips to get list of profiles
	profiles, err := osip.GetProfileList()

	if err != nil {
		return err
	}

	for _, profile := range profiles {
		if !osip.exportAll && !osip.profilesToExport[profile.Name] {
			continue
		}

		lookupNames, found := osip.replicationHints[profile.Name]

		if !found {
			lookupNames = []string{profile.Name}
		}

		for _, lookupName := range lookupNames {
			// Fetch size for every profile, as this contains the name after tags as well as replication and shared info
			size, err := osip.GetProfileSize(lookupName)

			if err != nil {
				if profileNotFound(err) {
					continue
				}

				return err
			}

			replLabels := make(metrics.Labels)

			if size.Shared == yes {
				replLabels["shared"] = yes
			}

			if size.Replicated == yes {
				replLabels["replicated"] = yes
			}

			// Dialog with values handling
			if profile.HasValue {
				// Fetch all values associated with the the profile
				vprofile, err := osip.GetProfileWithValues(lookupName)
				if err != nil {
					// ProfileNotFound is safe to continue on
					// May be cause by a profile going away between us querying the list and checking
					if profileNotFound(err) {
						continue
					}

					return err
				}

				for _, value := range vprofile {
					labels := maps.Clone(replLabels)

					// Is this an insight value?
					if strings.HasPrefix(value.Value, osip.insightPrefix) {
						osip.insightProfiles.SetWithStrLabels(size.Name, strings.TrimSpace(value.Value[len(osip.insightPrefix):]), labels, float64(value.Count))
						continue
					}

					labels["value"] = value.Value

					osip.exportValueProfiles.SetWithLabels(size.Name, labels, float64(value.Count))
				}
			} else {
				// Non-Value dialogs
				if len(replLabels) > 0 {
					osip.exportValueProfiles.SetWithLabels(size.Name, replLabels, float64(size.Count))
				} else {
					osip.exportProfiles.Set(size.Name, float64(size.Count))
				}
			}
		}
	}

	return nil
}

// Implement Promeheus Collector Interface
// Provides metric descriptions.
func (osip *Client) Describe(_ chan<- *prometheus.Desc) {
}
