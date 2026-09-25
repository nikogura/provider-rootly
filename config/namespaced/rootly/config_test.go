// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package rootly

import (
	"os"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Configure only registers configurator callbacks; they run inside
// ConfigureResources. So the test drives the same pipeline the provider
// does, on the real committed schema.
func TestConfigure(t *testing.T) {
	t.Parallel()

	schema, err := os.ReadFile("../../schema.json")
	if err != nil {
		t.Fatalf("cannot read provider schema: %v", err)
	}
	meta, err := os.ReadFile("../../provider-metadata.yaml")
	if err != nil {
		t.Fatalf("cannot read provider metadata: %v", err)
	}

	p := ujconfig.NewProvider(schema, "rootly", "github.com/nikogura/provider-rootly", meta,
		ujconfig.WithRootGroup("rootly.crossplane.io"),
		ujconfig.WithIncludeList([]string{
			"rootly_team$", "rootly_escalation_policy$", "rootly_escalation_level$",
			"rootly_schedule$", "rootly_heartbeat$", "rootly_alerts_source$",
		}))

	Configure(p)
	p.ConfigureResources()

	tests := []struct {
		name       string
		resource   string
		shortGroup string
		kind       string
	}{
		{name: "Team", resource: "rootly_team", shortGroup: "team", kind: "Team"},
		{name: "EscalationPolicy", resource: "rootly_escalation_policy", shortGroup: "escalation", kind: "EscalationPolicy"},
		{name: "EscalationLevel", resource: "rootly_escalation_level", shortGroup: "escalation", kind: "EscalationLevel"},
		{name: "Schedule", resource: "rootly_schedule", shortGroup: "schedule", kind: "Schedule"},
		{name: "Heartbeat", resource: "rootly_heartbeat", shortGroup: "heartbeat", kind: "Heartbeat"},
		{name: "AlertsSource", resource: "rootly_alerts_source", shortGroup: "alerts", kind: "AlertsSource"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, ok := p.Resources[tt.resource]
			if !ok {
				t.Fatalf("resource %q not present", tt.resource)
			}
			if r.ShortGroup != tt.shortGroup {
				t.Errorf("ShortGroup = %q, want %q", r.ShortGroup, tt.shortGroup)
			}
			if r.Kind != tt.kind {
				t.Errorf("Kind = %q, want %q", r.Kind, tt.kind)
			}
		})
	}

	ref, ok := p.Resources["rootly_escalation_level"].References["escalation_policy_id"]
	if !ok {
		t.Fatal("escalation_policy_id has no cross-resource reference")
	}
	if ref.TerraformName != "rootly_escalation_policy" {
		t.Errorf("escalation_policy_id references %q, want rootly_escalation_policy", ref.TerraformName)
	}
}
