// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Both providers must expose exactly the six kinds this provider exists for,
// under the names the custom configurators assign. This exercises the full
// pipeline: schema parse, include list, external-name defaults and the
// per-resource configurators.
func TestGetProvider(t *testing.T) {
	t.Parallel()

	wantKinds := map[string]struct {
		shortGroup string
		kind       string
	}{
		"rootly_team":              {shortGroup: "team", kind: "Team"},
		"rootly_escalation_policy": {shortGroup: "escalation", kind: "EscalationPolicy"},
		"rootly_escalation_level":  {shortGroup: "escalation", kind: "EscalationLevel"},
		"rootly_schedule":          {shortGroup: "schedule", kind: "Schedule"},
		"rootly_heartbeat":         {shortGroup: "heartbeat", kind: "Heartbeat"},
		"rootly_alerts_source":     {shortGroup: "alerts", kind: "AlertsSource"},
	}

	tests := []struct {
		name string
		get  func() *ujconfig.Provider
	}{
		{
			name: "Cluster",
			get:  GetProvider,
		},
		{
			name: "Namespaced",
			get:  GetProviderNamespaced,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := tt.get()
			if p == nil {
				t.Fatal("provider is nil")
			}

			if len(p.Resources) != len(wantKinds) {
				got := make([]string, 0, len(p.Resources))
				for name := range p.Resources {
					got = append(got, name)
				}
				t.Fatalf("expected exactly %d resources, got %d: %v", len(wantKinds), len(p.Resources), got)
			}

			for tfName, want := range wantKinds {
				r, ok := p.Resources[tfName]
				if !ok {
					t.Errorf("resource %q not generated", tfName)
					continue
				}
				if r.ShortGroup != want.shortGroup {
					t.Errorf("%s: ShortGroup = %q, want %q", tfName, r.ShortGroup, want.shortGroup)
				}
				if r.Kind != want.kind {
					t.Errorf("%s: Kind = %q, want %q", tfName, r.Kind, want.kind)
				}
				// Rootly assigns UUIDs; every resource must take its external
				// name from the provider.
				if !r.ExternalName.DisableNameInitializer {
					t.Errorf("%s: external name is not IdentifierFromProvider", tfName)
				}
			}

			level, ok := p.Resources["rootly_escalation_level"]
			if !ok {
				t.Fatal("rootly_escalation_level not generated")
			}
			ref, ok := level.References["escalation_policy_id"]
			if !ok {
				t.Fatal("escalation_policy_id has no cross-resource reference")
			}
			if ref.TerraformName != "rootly_escalation_policy" {
				t.Errorf("escalation_policy_id references %q, want rootly_escalation_policy", ref.TerraformName)
			}
		})
	}
}
