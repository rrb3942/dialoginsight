package osipcollect

import (
	"context"
	"strings"
)

const (
	errProfileNotFound = "Profile not found"
)

// Structs for deserializing opensips responses.
type ProfileValue struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type ProfileSizeWrapper struct {
	Profile ProfileSize `json:"Profile"`
}

type ProfileSize struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// API Call - profile_get_values
// Returns slice containing all values and counts for a given profile.
func (osip *Client) GetProfileWithValues(profile string) ([]ProfileValue, error) {
	profiles := []ProfileValue{}

	ctx, cancel := context.WithTimeout(context.Background(), osip.timeout)
	defer cancel()

	if err := osip.rpc.CallContext(ctx, &profiles, "profile_get_values", profile); err != nil {
		return nil, err
	}

	return profiles, nil
}

// API Call - profile_get_size
// Returns number of active dialogs in the given profile.
func (osip *Client) GetProfileSize(profile string) (int, error) {
	size := ProfileSizeWrapper{}

	ctx, cancel := context.WithTimeout(context.Background(), osip.timeout)
	defer cancel()

	if err := osip.rpc.CallContext(ctx, &size, "profile_get_size", profile); err != nil {
		return 0, err
	}

	return size.Profile.Count, nil
}

func profileNotFound(err error) bool {
	return strings.EqualFold(err.Error(), errProfileNotFound)
}
