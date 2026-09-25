package rootly

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Configure shapes the generated kinds. Groups follow the resource family and
// kinds carry the full name -- upjet's default would split
// rootly_escalation_policy into group "escalation", kind "Policy", and a kind
// named Policy tells you nothing at a glance in a cluster full of policies.
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("rootly_team", func(r *ujconfig.Resource) {
		r.ShortGroup = "team"
		r.Kind = "Team"
	})
	p.AddResourceConfigurator("rootly_escalation_policy", func(r *ujconfig.Resource) {
		r.ShortGroup = "escalation"
		r.Kind = "EscalationPolicy"
	})
	p.AddResourceConfigurator("rootly_escalation_level", func(r *ujconfig.Resource) {
		r.ShortGroup = "escalation"
		r.Kind = "EscalationLevel"
		// A level is meaningless outside its policy, and the id is a Rootly
		// UUID nobody should be copy-pasting: reference the policy by name.
		r.References["escalation_policy_id"] = ujconfig.Reference{
			TerraformName: "rootly_escalation_policy",
		}
	})
	p.AddResourceConfigurator("rootly_schedule", func(r *ujconfig.Resource) {
		r.ShortGroup = "schedule"
		r.Kind = "Schedule"
	})
	p.AddResourceConfigurator("rootly_heartbeat", func(r *ujconfig.Resource) {
		r.ShortGroup = "heartbeat"
		r.Kind = "Heartbeat"
	})
	p.AddResourceConfigurator("rootly_alerts_source", func(r *ujconfig.Resource) {
		r.ShortGroup = "alerts"
		r.Kind = "AlertsSource"
	})
}
