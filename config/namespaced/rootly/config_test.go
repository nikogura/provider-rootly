// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package rootly

import (
	"os"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/stretchr/testify/assert"
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
			"rootly_schedule$", "rootly_schedule_rotation$", "rootly_schedule_rotation_user$",
			"rootly_heartbeat$", "rootly_alerts_source$",
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
		{name: "ScheduleRotation", resource: "rootly_schedule_rotation", shortGroup: "schedule", kind: "ScheduleRotation"},
		{name: "ScheduleRotationUser", resource: "rootly_schedule_rotation_user", shortGroup: "schedule", kind: "ScheduleRotationUser"},
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

	references := []struct {
		resource string
		field    string
		target   string
	}{
		{resource: "rootly_escalation_level", field: "escalation_policy_id", target: "rootly_escalation_policy"},
		{resource: "rootly_schedule_rotation", field: "schedule_id", target: "rootly_schedule"},
		{resource: "rootly_schedule_rotation_user", field: "schedule_rotation_id", target: "rootly_schedule_rotation"},
	}

	ignored := p.Resources["rootly_alerts_source"].LateInitializer.IgnoredFields
	assert.ElementsMatch(t, []string{"alert_template_attributes", "alert_source_fields_attributes"}, ignored,
		"the conflicting AlertsSource blocks must not be late-initialized")

	for _, want := range references {
		ref, ok := p.Resources[want.resource].References[want.field]
		if !ok {
			t.Fatalf("%s.%s has no cross-resource reference", want.resource, want.field)
		}
		if ref.TerraformName != want.target {
			t.Errorf("%s.%s references %q, want %s", want.resource, want.field, ref.TerraformName, want.target)
		}
	}
}
