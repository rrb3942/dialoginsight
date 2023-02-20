package osipcollect

import (
	"dialoginsight/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
)

const (
	exportNamespace  = "dialoginsight_exported_profile_"
	insightNamespace = "dialoginsight_profile_"
)

// Implement Promeheus Collector Interface
// Collects and provides metrics
func (osip *Client) Collect(ch chan<- prometheus.Metric) {
	// Only allow one collect to run at a time
	osip.Lock()
	defer osip.Unlock()

	osip.exportProfiles.Reset()
	osip.exportValueProfiles.Reset()
	osip.insightProfiles.Reset()

	// API Call to opensips to get list of profiles
	profiles, err := osip.GetProfileList()

	if err != nil {
		log.Println(err)
		return
	}

	for _, profile := range profiles {
		if osip.exportAll || osip.profilesToExport[profile.Name] {
			// Dialog with values handling
			if profile.HasValue {
				// Fetch all values associated with the the profile
				if vprofile, err := osip.GetProfileWithValues(profile.Name); err != nil {
					// ProfileNotFound is safe to continue on
					//May be cause by a profile going away between us querying the list and checking
					if !profileNotFound(err) {
						log.Println(err)
						return
					}
				} else {
					for _, value := range vprofile {
						// Is this an insight value?
						if strings.HasPrefix(value.Value, osip.insightPrefix) {
							osip.insightProfiles.SetWithStrLabels(profile.Name, strings.TrimSpace(value.Value[len(osip.insightPrefix):]), float64(value.Count))
						} else {
							osip.exportValueProfiles.SetWithLabels(profile.Name, metrics.Labels{"value": value.Value}, float64(value.Count))
						}
					}
				}
				// Non-Value dialogs
			} else {
				if size, err := osip.GetProfileSize(profile.Name); err != nil {
					if !profileNotFound(err) {
						log.Println(err)
						return
					}
				} else {
					osip.exportProfiles.Set(profile.Name, float64(size))
				}
			}
		}
	}

	osip.exportProfiles.Collect(ch)
	osip.exportValueProfiles.Collect(ch)
	osip.insightProfiles.Collect(ch)
}

// Implement Promeheus Collector Interface
// Provides metric descriptions
func (osip *Client) Describe(ch chan<- *prometheus.Desc) {
}
