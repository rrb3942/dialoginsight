package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/maps"
)

// DynamicLabelGauges handles managing prometheus gauges with dynamic labels.
type DynamicLabelGauges struct {
	now         time.Time
	gauges      map[string]*prometheus.GaugeVec
	lastActive  map[string]time.Time
	namespace   string
	metricName  string
	help        string
	idleCleanup time.Duration
}

func NewDynamicLabelGauges(namespace, metricName, help string, idleCleanup time.Duration) *DynamicLabelGauges {
	dyng := new(DynamicLabelGauges)
	dyng.namespace = namespace
	dyng.metricName = metricName
	dyng.help = help
	dyng.gauges = make(map[string]*prometheus.GaugeVec)
	dyng.lastActive = make(map[string]time.Time)
	dyng.idleCleanup = idleCleanup

	return dyng
}

// All gauges should be reset before every collections
// Try to reuse gauges where possible.
func (dyng *DynamicLabelGauges) Reset() {
	// Reset current polling time
	dyng.now = time.Now()

	// Remove idle and inactive namespaces
	for g, t := range dyng.lastActive {
		if dyng.now.Sub(t) > dyng.idleCleanup {
			delete(dyng.gauges, g)
			delete(dyng.lastActive, g)
		}
	}

	// Reset any active namespaces
	for _, g := range dyng.gauges {
		g.Reset()
	}
}

func (dyng *DynamicLabelGauges) Collect(ch chan<- prometheus.Metric) {
	for k, t := range dyng.lastActive {
		// Only return namespaces that were collected during this period
		if t.Equal(dyng.now) {
			if g, found := dyng.gauges[k]; found {
				g.Collect(ch)
			}
		}
	}
}

// Sets the value for a given label, creating the gauge if one does not yet exist for the label keys.
func (dyng *DynamicLabelGauges) Set(labels Labels, value float64) {
	mapKey := labels.MapKey()
	g, found := dyng.gauges[mapKey]

	if !found {
		g = NewLabeledGauge(dyng.namespace, dyng.metricName, maps.Keys(labels), dyng.help)
		dyng.gauges[mapKey] = g
	}

	g.With(prometheus.Labels(labels)).Set(value)

	dyng.lastActive[mapKey] = dyng.now
}

// Sets the value for a given label, creating the gauge if one does not yet exist for the label keys.
func (dyng *DynamicLabelGauges) SetWithStrLabels(labelStr string, labels Labels, value float64) {
	maps.Copy(labels, NewLabelsFromString(labelStr))
	dyng.Set(labels, value)
}
