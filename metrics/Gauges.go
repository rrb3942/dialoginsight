package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Wrapper func to create a new prometheus gaugevec with labels.
func NewLabeledGauge(namespace, name string, labels []string, help string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      name,
			Help:      help,
		},
		labels,
	)
}
