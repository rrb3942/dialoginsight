package osipcollect

import (
	"context"
)

// Structs for deserializing API calls to opensips.
type ProfilesWrapper struct {
	Profiles []Profiles `json:"Profiles"`
}

type Profiles struct {
	Name     string `json:"name"`
	HasValue bool   `json:"has value"`
}

// API Call - list_all_profiles
// Returns a list of all profile names on a server.
func (osip *Client) GetProfileList() ([]Profiles, error) {
	profiles := ProfilesWrapper{}

	ctx, cancel := context.WithTimeout(context.Background(), osip.timeout)
	defer cancel()

	if err := osip.rpc.CallContext(ctx, &profiles, "list_all_profiles"); err != nil {
		return nil, err
	}

	return profiles.Profiles, nil
}
