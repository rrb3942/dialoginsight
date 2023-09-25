package osipcollect

import (
	"strings"
)

const (
	hasTag = 2
)

func parseSharedTags(profiles []string) (export map[string]bool, hints map[string][]string) {
	export = make(map[string]bool)
	hints = make(map[string][]string)

	for _, profile := range profiles {
		split := strings.SplitN(profile, "/", hasTag)
		if len(split) == hasTag {
			hints[split[0]] = append(hints[split[0]], profile)
		}

		export[split[0]] = true
	}

	return export, hints
}
