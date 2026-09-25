// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"strings"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	clusterv1beta1 "github.com/nikogura/provider-rootly/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/nikogura/provider-rootly/apis/namespaced/v1beta1"

	clusterteamv1alpha1 "github.com/nikogura/provider-rootly/apis/cluster/team/v1alpha1"
	namespacedteamv1alpha1 "github.com/nikogura/provider-rootly/apis/namespaced/team/v1alpha1"
)

func newTestScheme(t *testing.T) (s *runtime.Scheme) {
	t.Helper()
	s = runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{
		corev1.AddToScheme,
		clusterv1beta1.SchemeBuilder.AddToScheme,
		namespacedv1beta1.SchemeBuilder.AddToScheme,
		clusterteamv1alpha1.SchemeBuilder.AddToScheme,
		namespacedteamv1alpha1.SchemeBuilder.AddToScheme,
	} {
		err := add(s)
		if err != nil {
			t.Fatalf("cannot build scheme: %v", err)
		}
	}
	return s
}

func newLegacyTeam(pcName string) (mr *clusterteamv1alpha1.Team) {
	mr = &clusterteamv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{
			Name: "a-team",
			UID:  types.UID("11111111-2222-3333-4444-555555555555"),
		},
	}
	mr.SetGroupVersionKind(clusterteamv1alpha1.Team_GroupVersionKind)
	if pcName != "" {
		mr.SetProviderConfigReference(&xpv2.Reference{Name: pcName})
	}
	return mr
}

func newModernTeam(ns, pcName, pcKind string) (mr *namespacedteamv1alpha1.Team) {
	mr = &namespacedteamv1alpha1.Team{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "a-team",
			Namespace: ns,
			UID:       types.UID("66666666-7777-8888-9999-000000000000"),
		},
	}
	mr.SetGroupVersionKind(namespacedteamv1alpha1.Team_GroupVersionKind)
	if pcName != "" {
		mr.SetProviderConfigReference(&xpv2.ProviderConfigReference{Name: pcName, Kind: pcKind})
	}
	return mr
}

func secretCredentials(ns, name string, payload string) (secret *corev1.Secret) {
	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
		Data:       map[string][]byte{"credentials": []byte(payload)},
	}
	return secret
}

func clusterPC(name, secretNS, secretName string) (pc *clusterv1beta1.ProviderConfig) {
	pc = &clusterv1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: clusterv1beta1.ProviderConfigSpec{
			Credentials: clusterv1beta1.ProviderCredentials{
				Source: xpv2.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
					SecretRef: &xpv2.SecretKeySelector{
						SecretReference: xpv2.SecretReference{Name: secretName, Namespace: secretNS},
						Key:             "credentials",
					},
				},
			},
		},
	}
	return pc
}

func namespacedPC(ns, name, secretName string) (pc *namespacedv1beta1.ProviderConfig) {
	pc = &namespacedv1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns},
		Spec: namespacedv1beta1.ProviderConfigSpec{
			Credentials: namespacedv1beta1.ProviderCredentials{
				Source: xpv2.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
					SecretRef: &xpv2.SecretKeySelector{
						// Namespace intentionally wrong: resolveModern must
						// pin it to the managed resource's namespace.
						SecretReference: xpv2.SecretReference{Name: secretName, Namespace: "somewhere-else"},
						Key:             "credentials",
					},
				},
			},
		},
	}
	return pc
}

func TestTerraformSetupBuilder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		mg         resource.Managed
		objects    []client.Object
		wantErr    string
		wantConfig map[string]any
	}{
		{
			name: "LegacyClusterScopedHappyPath",
			mg:   newLegacyTeam("default"),
			objects: []client.Object{
				clusterPC("default", "crossplane-system", "rootly-creds"),
				secretCredentials("crossplane-system", "rootly-creds",
					`{"api_token":"tok-123","api_host":"https://api.rootly.example"}`),
			},
			wantConfig: map[string]any{
				"api_token": "tok-123",
				"api_host":  "https://api.rootly.example",
			},
		},
		{
			name: "LegacyTokenOnlyOmitsHost",
			mg:   newLegacyTeam("default"),
			objects: []client.Object{
				clusterPC("default", "crossplane-system", "rootly-creds"),
				secretCredentials("crossplane-system", "rootly-creds", `{"api_token":"tok-123"}`),
			},
			wantConfig: map[string]any{"api_token": "tok-123"},
		},
		{
			name: "ModernNamespacedHappyPath",
			mg:   newModernTeam("demo", "default", "ProviderConfig"),
			objects: []client.Object{
				namespacedPC("demo", "default", "rootly-creds"),
				secretCredentials("demo", "rootly-creds", `{"api_token":"tok-456"}`),
			},
			wantConfig: map[string]any{"api_token": "tok-456"},
		},
		{
			name:    "MissingProviderConfigRef",
			mg:      newLegacyTeam(""),
			wantErr: errNoProviderConfig,
		},
		{
			name:    "ProviderConfigNotFound",
			mg:      newLegacyTeam("absent"),
			wantErr: errGetProviderConfig,
		},
		{
			name: "CredentialsAreNotJSON",
			mg:   newLegacyTeam("default"),
			objects: []client.Object{
				clusterPC("default", "crossplane-system", "rootly-creds"),
				secretCredentials("crossplane-system", "rootly-creds", "not json"),
			},
			wantErr: errUnmarshalCredentials,
		},
		{
			name: "ModernUnknownProviderConfigKind",
			mg:   newModernTeam("demo", "default", "NoSuchKind"),
			objects: []client.Object{
				namespacedPC("demo", "default", "rootly-creds"),
			},
			wantErr: "unknown GVK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kube := fake.NewClientBuilder().
				WithScheme(newTestScheme(t)).
				WithObjects(tt.objects...).
				Build()

			fn := TerraformSetupBuilder("1.5.7", "rootlyhq/rootly", "5.20.1")
			setup, err := fn(context.Background(), kube, tt.mg)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if setup.Version != "1.5.7" {
				t.Errorf("Version = %q, want 1.5.7", setup.Version)
			}
			if setup.Requirement.Source != "rootlyhq/rootly" || setup.Requirement.Version != "5.20.1" {
				t.Errorf("Requirement = %+v, want rootlyhq/rootly@5.20.1", setup.Requirement)
			}

			if len(setup.Configuration) != len(tt.wantConfig) {
				t.Fatalf("Configuration = %v, want %v", setup.Configuration, tt.wantConfig)
			}
			for k, want := range tt.wantConfig {
				if got := setup.Configuration[k]; got != want {
					t.Errorf("Configuration[%q] = %v, want %v", k, got, want)
				}
			}
		})
	}
}

func TestToSharedPCSpec(t *testing.T) {
	t.Parallel()

	t.Run("NilInput", func(t *testing.T) {
		t.Parallel()

		spec, err := toSharedPCSpec(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if spec != nil {
			t.Fatalf("expected nil spec for nil input, got %+v", spec)
		}
	})

	t.Run("RoundTripsCredentials", func(t *testing.T) {
		t.Parallel()

		pc := clusterPC("default", "crossplane-system", "rootly-creds")
		spec, err := toSharedPCSpec(pc)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if spec == nil {
			t.Fatal("spec is nil")
		}
		if spec.Credentials.Source != xpv2.CredentialsSourceSecret {
			t.Errorf("Source = %q, want Secret", spec.Credentials.Source)
		}
		ref := spec.Credentials.SecretRef
		if ref == nil {
			t.Fatal("SecretRef did not survive conversion")
		}
		if ref.Name != "rootly-creds" || ref.Namespace != "crossplane-system" || ref.Key != "credentials" {
			t.Errorf("SecretRef = %+v, want rootly-creds/crossplane-system/credentials", ref)
		}
	})
}
