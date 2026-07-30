## Context

TKE exposes cluster-level roll-out sequence tags through two `tke/v20180525` SDK APIs that already exist in the vendored SDK:

- `ModifyClusterRollOutSequenceTags(ClusterID, Tags []*Tag)` — sets/replaces the tag list on a cluster; an empty `Tags` removes all tags.
- `DescribeClusterRollOutSequenceTags(Offset, Limit, Filters []*Filter)` — returns `ClusterTags []*ClusterRollOutSequenceTag` and `TotalCount`.

SDK types (verified in `vendor/.../tke/v20180525/models.go`):
- `Tag{ Key *string, Value *string }`
- `ClusterRollOutSequenceTag{ ClusterID *string, ClusterName *string, Tags []*Tag, Region *string }`
- `Filter{ Name *string, Values []*string }`

This is a **config-type** resource: it manages a property collection (roll-out sequence tags) of an existing cluster rather than creating/destroying an independent object. The cluster ID is the natural unique identifier.

The closest existing reference in the same service is `tencentcloud_kubernetes_roll_out_sequence` (`resource_tc_kubernetes_roll_out_sequence.go`), which demonstrates the TKE client usage, paginated Describe loop, retry pattern, and tag list flattening. Code style additionally follows `tencentcloud_igtm_monitor`.

## Goals / Non-Goals

**Goals:**
- Provide full CRUD for cluster roll-out sequence tags keyed by `ClusterID`.
- Keep the schema arguments aligned exactly with `ModifyClusterRollOutSequenceTags` input (`ClusterID` + `Tags`).
- Use retry on every API call; nil-safe access on all response fields.
- Ship resource example `.md`, website docs, unit test, and provider registration.

**Non-Goals:**
- No management of the cluster itself (assumes the cluster already exists).
- No exposure of `Offset`/`Limit` to users — pagination is handled internally with the documented maximum page size.
- No async waiting/Timeouts block: `ModifyClusterRollOutSequenceTags` is a synchronous metadata write with no async status to poll.

## Decisions

### Resource name and file layout
- Resource: `tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config`
- Files under `tencentcloud/services/tke/`:
  - `resource_tc_kubernetes_cluster_roll_out_sequence_tag_config.go`
  - `resource_tc_kubernetes_cluster_roll_out_sequence_tag_config.md`
  - `resource_tc_kubernetes_cluster_roll_out_sequence_tag_config_test.go`
- Website doc: `website/docs/r/kubernetes_cluster_roll_out_sequence_tag_config.html.markdown`
- Constructor `ResourceTencentCloudKubernetesClusterRollOutSequenceTagConfig()`, registered in `provider.go`.

### Schema (must match ModifyClusterRollOutSequenceTags input)
- `cluster_id` (TypeString, Required, ForceNew) → `ClusterID`. ForceNew because it is the resource identity.
- `tags` (TypeList, Required) → `Tags`, element is a nested resource:
  - `key` (TypeString, Required) → `Tag.Key`
  - `value` (TypeString, Required) → `Tag.Value`

  Rationale: the SDK `Tag` is a flat `Key`/`Value` string pair (distinct from `SequenceTag` used by `tencentcloud_kubernetes_roll_out_sequence`, whose `Value` is a list). The schema mirrors the SDK `Tag` exactly to satisfy the "schema == Modify input" rule.

### CRUD mapping (config-type, following `tencentcloud_waf_owasp_rule_status_config`)
- **Create**: set `d.SetId(cluster_id)` and delegate to Update — no direct API call in Create.
- **Read**: delegate to service-layer `TkeService.DescribeKubernetesClusterRollOutSequenceTagConfigById(ctx, clusterId)` which pages through `DescribeClusterRollOutSequenceTagsWithContext` (`Limit = 100`, retry + ratelimit, accumulate until `Offset+Limit >= TotalCount`) and returns the entry whose `ClusterID == clusterId`; if none found or its tag list is empty, `d.SetId("")`. Flatten `Tags` into the `tags` list with nil checks.
- **Update**: rebuild full `Tags` from `tags` and call `ModifyClusterRollOutSequenceTagsWithContext` inside `resource.Retry(WriteRetryTimeout, ...)`, then Read.
- **Delete**: call `ModifyClusterRollOutSequenceTags` with `ClusterID = d.Id()` and `Tags = []*Tag{}` (empty list) inside retry — this removes all tags (per business rule).

### Filter name for Describe
- Use the cluster-id filter to scope the describe query. The exact filter `Name` value will be confirmed against the API doc / examples during implementation (`ClusterId`); pagination still guards correctness if the filter is ignored by the backend.

## Risks / Trade-offs

- [The `DescribeClusterRollOutSequenceTags` filter key name (`ClusterId` vs `ClusterID`) is not 100% certain from the type definition] → Mitigation: rely on full pagination + client-side match on `ClusterID`, so Read is correct even if the server-side filter is permissive; confirm the exact filter name during apply.
- [`Limit` documented default is 20; max may differ] → Mitigation: use 100 as the page size and loop with `TotalCount`; if the backend caps lower, pagination still collects all pages.
- [Config-type resource with ForceNew `cluster_id`] → changing the cluster recreates the resource, which is the intended semantics (tags belong to a specific cluster).

## Migration Plan

Purely additive; no migration. New resource registered alongside existing TKE resources. Rollback = revert the additive files and provider registration.
