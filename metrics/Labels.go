package metrics

import (
	"sort"
	"strconv"
	"strings"

	"github.com/carlmjohnson/bytemap"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/maps"
)

// Prometheus labels can only contain these values.
const promLabelStart = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const promLabelSet = promLabelStart + "_0123456789"

// bytemap lets us easily check if a label only contains valid characters.
var labelSet = bytemap.Make(promLabelSet)

// Extened prometheus.Labels.
type Labels prometheus.Labels

// Generate a key from the label keys that we can use for a map lookup.
func (l Labels) MapKey() string {
	// Hangle empty maps with an empty string
	if len(l) == 0 {
		return ""
	}

	keys := maps.Keys(l)
	sort.Strings(keys)

	return strings.Join(keys, ";")
}

// Parses a string in the format of 'label1=value;label2=value....' into prometheus label.
func NewLabelsFromString(lstr string) Labels {
	labels := make(Labels)

	for _, tupleStr := range strings.Split(lstr, ";") {
		if len(tupleStr) > 0 {
			if label, value, found := strings.Cut(strings.TrimSpace(tupleStr), "="); found {
				// Skip labels that are not valid
				// Labels follow the pattern [a-zA-Z_][a-zA-Z0-9_]*
				if strings.Contains(promLabelStart, label[0:1]) && labelSet.Contains(label) {
					if u, err := strconv.Unquote(value); err == nil {
						labels[label] = u
					} else {
						labels[label] = value
					}
				}
			}
		}
	}

	return labels
}
