# provider-rootly

A [Crossplane](https://crossplane.io) provider for [Rootly](https://rootly.com),
generated with [Upjet](https://github.com/crossplane/upjet) from the official
[rootlyhq/terraform-provider-rootly](https://github.com/rootlyhq/terraform-provider-rootly).

Incident-management substrate as managed resources: drive Rootly from git,
alongside the rest of your infrastructure, with the same reconcile loop.

## Managed resources

The SRE set, deliberately small to start:

| Kind | Group | Terraform resource |
|---|---|---|
| `Team` | `team.rootly.crossplane.io` | `rootly_team` |
| `EscalationPolicy` | `escalation.rootly.crossplane.io` | `rootly_escalation_policy` |
| `EscalationLevel` | `escalation.rootly.crossplane.io` | `rootly_escalation_level` |
| `Schedule` | `schedule.rootly.crossplane.io` | `rootly_schedule` |
| `Heartbeat` | `heartbeat.rootly.crossplane.io` | `rootly_heartbeat` |
| `AlertsSource` | `alerts.rootly.crossplane.io` | `rootly_alerts_source` |

Every kind also exists namespaced under `*.rootly.m.crossplane.io` for
Crossplane v2 namespaced composition. Rootly assigns every resource its ID, so
external names are always learned from the API after creation -- never set one
by hand.

The upstream provider covers ~229 resources; widening the scope is one entry in
`config/external_name.go` and a regenerate.

## Install

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-rootly
spec:
  package: ghcr.io/nikogura/provider-rootly:v0.1.0
```

Pin by digest in anything you operate; tags are for humans.

Every green merge to main is a release: CI computes the next semantic version
from conventional commits, publishes the multi-arch image and package to ghcr
under that tag, and cuts the matching GitHub Release. Nothing is tagged by hand.

## Authenticate

One bearer token. The credentials Secret holds a JSON document:

```bash
kubectl create secret generic rootly-creds -n crossplane-system --from-literal=credentials='{"api_token":"<ROOTLY_API_TOKEN>"}'
```

```yaml
apiVersion: rootly.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: rootly-creds
      namespace: crossplane-system
      key: credentials
```

`api_host` is an optional second key, only for self-hosted or regional
endpoints; it defaults to `https://api.rootly.com`.

## Adopting an existing Rootly organization

Start observe-only: create MRs with
`managementPolicies: ["Observe"]` and the `crossplane.io/external-name`
annotation set to the resource's Rootly ID, watch them all come up Synced and
Ready against reality, and only then widen the policies. A wrong external name
does not fail -- full management would build a duplicate beside the real thing
and report success.

## Developing

```bash
make submodules   # once, after clone
make generate     # regenerate from config/schema.json + upstream docs
make build        # binary + package
```

The generation inputs are pinned in the Makefile
(`TERRAFORM_PROVIDER_VERSION`); bumping the upstream provider version is a
version bump plus `make generate`, then read the diff.

## License

Apache-2.0. At runtime the provider drives the MPL-2.0
`terraform-provider-rootly` binary, which is downloaded unmodified at image
build; see NOTICE.
