package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// DynamicGauges handles managing prometheus gauges with dynamic namespaces and labels
type DynamicGauges struct {
	// Fields used for creating new gauges
	namespacePrefix string
	metricName      string
	help            string
	// Mapping of namespaces -> dyanmiclabelgauges
	namespaces map[string]*DynamicLabelGauges
	// Track which namespaces are active
	lastActive  map[string]time.Time
	now         time.Time
	idleCleanup time.Duration
}

func NewDynamicGauges(namespacePrefix, metricName, help string, idleCleanup time.Duration) *DynamicGauges {
	dyng := new(DynamicGauges)
	dyng.namespacePrefix = namespacePrefix
	dyng.metricName = metricName
	dyng.help = help
	dyng.namespaces = make(map[string]*DynamicLabelGauges)
	dyng.lastActive = make(map[string]time.Time)
	dyng.idleCleanup = idleCleanup
	return dyng
}

// All namespaces should be reset before every collections
// Try to reuse namespaces where possible
func (dyng *DynamicGauges) Reset() {
	// Reset current polling time
	dyng.now = time.Now()

	// Remove idle and inactive namespaces
	for g, t := range dyng.lastActive {
		if dyng.now.Sub(t) > dyng.idleCleanup {
			delete(dyng.namespaces, g)
			delete(dyng.lastActive, g)
		}
	}

	// Reset any active namespaces
	for _, g := range dyng.namespaces {
		g.Reset()
	}
}

func (dyng *DynamicGauges) Collect(ch chan<- prometheus.Metric) {
	for k, t := range dyng.lastActive {
		// Only return namespaces that were collected during this period
		if t.Equal(dyng.now) {
			if g, found := dyng.namespaces[k]; found {
				g.Collect(ch)
			}
		}
	}
}

// Sets the value for a given namespace and label, creating the namespace and gauge if one does not yet exist for the label keys
func (dyng *DynamicGauges) SetWithLabels(namespace string, labels Labels, value float64) {
	g, found := dyng.namespaces[namespace]

	if !found {
		g = NewDynamicLabelGauges(dyng.namespacePrefix+namespace, dyng.metricName, dyng.help, dyng.idleCleanup)
		dyng.namespaces[namespace] = g
	}

	g.Set(labels, value)
	dyng.lastActive[namespace] = dyng.now
}

func (dyng *DynamicGauges) Set(namespace string, value float64) {
	dyng.SetWithLabels(namespace, nil, value)
}

func (dyng *DynamicGauges) SetWithStrLabels(namespace, strLabels string, value float64) {
	dyng.SetWithLabels(namespace, NewLabelsFromString(strLabels), value)
}
