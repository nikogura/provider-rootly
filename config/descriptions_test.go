// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/registry"
)

// Every kind this provider generates must say what it is: the description is
// what `kubectl explain` shows. Adding a resource to ExternalNameConfigs
// without describing it fails here.
func TestEveryGeneratedKindIsDescribed(t *testing.T) {
	t.Parallel()

	for name := range ExternalNameConfigs {
		description, ok := KindDescriptions[name]
		if !ok || strings.TrimSpace(description) == "" {
			t.Errorf("resource %q has no kind description", name)
		}
	}

	for name := range KindDescriptions {
		if _, ok := ExternalNameConfigs[name]; !ok {
			t.Errorf("kind description for %q, which is not a generated resource", name)
		}
	}
}

func TestDescribeKinds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		resource *ujconfig.Resource
		want     string
	}{
		{
			name:     "ResourceWithoutScrapedMetadata",
			resource: &ujconfig.Resource{Name: "rootly_team"},
			want:     KindDescriptions["rootly_team"],
		},
		{
			name: "ScrapedDescriptionIsReplaced",
			resource: &ujconfig.Resource{
				Name:         "rootly_schedule_rotation",
				MetaResource: &registry.Resource{Description: "Manages a schedule rotation."},
			},
			want: KindDescriptions["rootly_schedule_rotation"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			DescribeKinds()(tt.resource)

			if tt.resource.MetaResource == nil {
				t.Fatal("MetaResource is nil, so the kind has no description")
			}
			if tt.resource.MetaResource.Description != tt.want {
				t.Errorf("Description = %q, want %q", tt.resource.MetaResource.Description, tt.want)
			}
		})
	}

	unknown := &ujconfig.Resource{Name: "rootly_no_such_resource"}
	DescribeKinds()(unknown)

	if unknown.MetaResource != nil {
		t.Error("a resource this provider does not generate was given metadata")
	}
}

// Both providers must carry the descriptions through the full pipeline.
func TestProvidersDescribeEveryKind(t *testing.T) {
	t.Parallel()

	providers := map[string]func() *ujconfig.Provider{
		"Cluster":    GetProvider,
		"Namespaced": GetProviderNamespaced,
	}

	for scope, get := range providers {
		t.Run(scope, func(t *testing.T) {
			t.Parallel()

			for name, r := range get().Resources {
				if r.MetaResource == nil || r.MetaResource.Description != KindDescriptions[name] {
					t.Errorf("%s: kind description not applied", name)
				}
			}
		})
	}
}

// The generated CRDs are what users read. A kind whose description was never
// set renders the template's missing value into the CRD verbatim.
func TestGeneratedCRDsHaveNoMissingValues(t *testing.T) {
	t.Parallel()

	crds, err := filepath.Glob("../package/crds/*.yaml")
	if err != nil {
		t.Fatalf("cannot list generated CRDs: %v", err)
	}
	if len(crds) == 0 {
		t.Fatal("no generated CRDs found")
	}

	for _, path := range crds {
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatalf("cannot read %s: %v", path, readErr)
		}
		if strings.Contains(string(content), "<no value>") {
			t.Errorf("%s contains an unrendered template value", filepath.Base(path))
		}
	}
}
