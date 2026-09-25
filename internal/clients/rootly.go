// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
// SPDX-FileCopyrightText: 2026 Nik Ogura <nik.ogura@gmail.com>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/v2/pkg/terraform"

	clusterv1beta1 "github.com/nikogura/provider-rootly/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/nikogura/provider-rootly/apis/namespaced/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal rootly credentials as JSON"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) (fn terraform.SetupFn) {
	fn = func(ctx context.Context, client client.Client, mg resource.Managed) (setup terraform.Setup, err error) {
		setup = terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		var pcSpec *namespacedv1beta1.ProviderConfigSpec
		pcSpec, err = resolveProviderConfig(ctx, client, mg)
		if err != nil {
			err = errors.Wrap(err, "cannot resolve provider config")
			return setup, err
		}

		var data []byte
		data, err = resource.CommonCredentialExtractor(ctx, pcSpec.Credentials.Source, client, pcSpec.Credentials.CommonCredentialSelectors)
		if err != nil {
			err = errors.Wrap(err, errExtractCredentials)
			return setup, err
		}
		creds := map[string]string{}
		err = json.Unmarshal(data, &creds)
		if err != nil {
			err = errors.Wrap(err, errUnmarshalCredentials)
			return setup, err
		}

		// Set credentials in Terraform provider configuration. Rootly's auth is
		// a single bearer token; api_host is optional and only present for
		// self-hosted or regional endpoints.
		setup.Configuration = map[string]any{}
		if v, ok := creds["api_token"]; ok {
			setup.Configuration["api_token"] = v
		}
		if v, ok := creds["api_host"]; ok {
			setup.Configuration["api_host"] = v
		}
		return setup, err
	}
	return fn
}

func toSharedPCSpec(pc *clusterv1beta1.ProviderConfig) (spec *namespacedv1beta1.ProviderConfigSpec, err error) {
	if pc == nil {
		return spec, err
	}
	var data []byte
	data, err = json.Marshal(pc.Spec)
	if err != nil {
		return spec, err
	}

	var mSpec namespacedv1beta1.ProviderConfigSpec
	err = json.Unmarshal(data, &mSpec)
	spec = &mSpec
	return spec, err
}

func resolveProviderConfig(ctx context.Context, crClient client.Client, mg resource.Managed) (spec *namespacedv1beta1.ProviderConfigSpec, err error) {
	switch managed := mg.(type) {
	case resource.LegacyManaged: //nolint:staticcheck // still handling cluster-scoped behavior
		spec, err = resolveLegacy(ctx, crClient, managed)
		return spec, err
	case resource.ModernManaged:
		spec, err = resolveModern(ctx, crClient, managed)
		return spec, err
	default:
		err = errors.New("resource is not a managed resource")
		return spec, err
	}
}

func resolveLegacy(ctx context.Context, client client.Client, mg resource.LegacyManaged) (spec *namespacedv1beta1.ProviderConfigSpec, err error) { //nolint:staticcheck // still handling cluster-scoped behavior
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		err = errors.New(errNoProviderConfig)
		return spec, err
	}
	pc := &clusterv1beta1.ProviderConfig{}
	err = client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc)
	if err != nil {
		err = errors.Wrap(err, errGetProviderConfig)
		return spec, err
	}

	t := resource.NewLegacyProviderConfigUsageTracker(client, &clusterv1beta1.ProviderConfigUsage{})
	err = t.Track(ctx, mg)
	if err != nil {
		err = errors.Wrap(err, errTrackUsage)
		return spec, err
	}

	spec, err = toSharedPCSpec(pc)
	return spec, err
}

func resolveModern(ctx context.Context, crClient client.Client, mg resource.ModernManaged) (spec *namespacedv1beta1.ProviderConfigSpec, err error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		err = errors.New(errNoProviderConfig)
		return spec, err
	}

	var pcRuntimeObj runtime.Object
	pcRuntimeObj, err = crClient.Scheme().New(namespacedv1beta1.SchemeGroupVersion.WithKind(configRef.Kind))
	if err != nil {
		err = errors.Wrap(err, "unknown GVK for ProviderConfig")
		return spec, err
	}
	pcObj, ok := pcRuntimeObj.(client.Object)
	if !ok {
		// This indicates a programming error, types are not properly generated
		err = errors.New(" is not an Object")
		return spec, err
	}

	// Namespace will be ignored if the PC is a cluster-scoped type
	err = crClient.Get(ctx, types.NamespacedName{Name: configRef.Name, Namespace: mg.GetNamespace()}, pcObj)
	if err != nil {
		err = errors.Wrap(err, errGetProviderConfig)
		return spec, err
	}

	var pcSpec namespacedv1beta1.ProviderConfigSpec
	pcu := &namespacedv1beta1.ProviderConfigUsage{}
	switch pc := pcObj.(type) {
	case *namespacedv1beta1.ProviderConfig:
		pcSpec = pc.Spec
		if pcSpec.Credentials.SecretRef != nil {
			pcSpec.Credentials.SecretRef.Namespace = mg.GetNamespace()
		}
	case *namespacedv1beta1.ClusterProviderConfig:
		pcSpec = pc.Spec
	default:
		err = errors.New("unknown provider config type")
		return spec, err
	}
	t := resource.NewProviderConfigUsageTracker(crClient, pcu)
	err = t.Track(ctx, mg)
	if err != nil {
		err = errors.Wrap(err, errTrackUsage)
		return spec, err
	}
	spec = &pcSpec
	return spec, err
}
