package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider. Every Rootly resource is identified by an API-assigned ID (a
// UUID), so nothing here can be named in advance: the external name is
// always learned from the provider after creation. This map is also the
// provider's scope -- only resources listed here are generated.
var ExternalNameConfigs = map[string]config.ExternalName{
	"rootly_team":              config.IdentifierFromProvider,
	"rootly_escalation_policy": config.IdentifierFromProvider,
	"rootly_escalation_level":  config.IdentifierFromProvider,
	"rootly_schedule":          config.IdentifierFromProvider,
	"rootly_heartbeat":         config.IdentifierFromProvider,
	"rootly_alerts_source":     config.IdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to
// v1beta1 assuming they will be tested.
func ExternalNameConfigurations() (opt config.ResourceOption) {
	opt = func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
	return opt
}

// ExternalNameConfigured returns the list of all resources whose external
// name is configured manually.
func ExternalNameConfigured() (list []string) {
	list = make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		list[i] = name + "$"
		i++
	}
	return list
}
