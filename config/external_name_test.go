// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"strings"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

func TestExternalNameConfigured(t *testing.T) {
	t.Parallel()

	list := ExternalNameConfigured()

	if len(list) != len(ExternalNameConfigs) {
		t.Fatalf("expected %d entries, got %d", len(ExternalNameConfigs), len(list))
	}

	seen := map[string]bool{}
	for _, entry := range list {
		if !strings.HasSuffix(entry, "$") {
			t.Errorf("entry %q is not anchored with $", entry)
		}
		seen[strings.TrimSuffix(entry, "$")] = true
	}

	for name := range ExternalNameConfigs {
		if !seen[name] {
			t.Errorf("resource %q missing from configured list", name)
		}
	}
}

func TestExternalNameConfigurations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		resource     string
		wantDisabled bool
	}{
		{
			name:         "ConfiguredResourceGetsIdentifierFromProvider",
			resource:     "rootly_team",
			wantDisabled: true,
		},
		{
			name:         "UnconfiguredResourceIsLeftAlone",
			resource:     "rootly_no_such_resource",
			wantDisabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := &ujconfig.Resource{Name: tt.resource}
			ExternalNameConfigurations()(r)

			// IdentifierFromProvider is the only external-name config this
			// provider uses, and DisableNameInitializer is its signature.
			if r.ExternalName.DisableNameInitializer != tt.wantDisabled {
				t.Errorf("DisableNameInitializer = %v, want %v",
					r.ExternalName.DisableNameInitializer, tt.wantDisabled)
			}
		})
	}
}
