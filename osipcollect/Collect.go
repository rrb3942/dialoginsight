package osipcollect

import (
	"dialoginsight/metrics"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	exportNamespace  = "dialoginsight_exported_profile_"
	insightNamespace = "dialoginsight_profile_"
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
		if osip.exportAll || osip.profilesToExport[profile.Name] {
			// Dialog with values handling
			if profile.HasValue {
				// Fetch all values associated with the the profile
				vprofile, err := osip.GetProfileWithValues(profile.Name)
				if err != nil {
					// ProfileNotFound is safe to continue on
					// May be cause by a profile going away between us querying the list and checking
					if profileNotFound(err) {
						continue
					}

					return err
				}

				for _, value := range vprofile {
					// Is this an insight value?
					if strings.HasPrefix(value.Value, osip.insightPrefix) {
						osip.insightProfiles.SetWithStrLabels(profile.Name, strings.TrimSpace(value.Value[len(osip.insightPrefix):]), float64(value.Count))
						continue
					}

					osip.exportValueProfiles.SetWithLabels(profile.Name, metrics.Labels{"value": value.Value}, float64(value.Count))
				}
			} else {
				// Non-Value dialogs
				size, err := osip.GetProfileSize(profile.Name)

				if err != nil {
					if profileNotFound(err) {
						continue
					}

					return err
				}

				osip.exportProfiles.Set(profile.Name, float64(size))
			}
		}
	}

	return nil
}

// Implement Promeheus Collector Interface
// Provides metric descriptions.
func (osip *Client) Describe(ch chan<- *prometheus.Desc) {
}
