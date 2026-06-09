## Why

TKE supports assigning roll-out sequence tags (e.g. `Env`, `Protection-Level`) to a cluster so that platform-level batch operations can roll out across clusters in a controlled order. There is currently no Terraform resource to manage these cluster-level roll-out sequence tags, forcing users to set them manually in the console or via raw API calls.

## What Changes

- Add a new **config-type** resource `tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config` that manages the roll-out sequence tags bound to a single TKE cluster.
- Create lifecycle backed by `ModifyClusterRollOutSequenceTags` (set the desired tag list on a cluster).
- Read lifecycle backed by `DescribeClusterRollOutSequenceTags` (paginated lookup filtered by cluster ID).
- Update lifecycle reuses `ModifyClusterRollOutSequenceTags` to converge to the desired tag list.
- Delete lifecycle calls `ModifyClusterRollOutSequenceTags` with an **empty `Tags` list**, which removes all roll-out sequence tags from the cluster.
- The resource unique ID is the cluster ID (`ClusterID`).
- Provider registration, resource example doc (`.md`), website documentation, and unit test are added.

## Capabilities

### New Capabilities
- `kubernetes-cluster-roll-out-sequence-tag-config`: Manage the roll-out sequence tags (key/value pairs) bound to a TKE cluster, including create/read/update/delete via the Modify/Describe ClusterRollOutSequenceTags APIs.

### Modified Capabilities
<!-- None: this is a brand-new resource and does not change requirements of existing capabilities. -->

## Impact

- New file: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_roll_out_sequence_tag_config.go`
- New file: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_roll_out_sequence_tag_config.md`
- New file: `tencentcloud/services/tke/resource_tc_kubernetes_cluster_roll_out_sequence_tag_config_test.go`
- New file: `website/docs/r/kubernetes_cluster_roll_out_sequence_tag_config.html.markdown`
- Modified: `tencentcloud/provider.go` (register resource), `tencentcloud/provider.md`, `website/tencentcloud.erb`
- SDK: uses existing `tke/v20180525` APIs `ModifyClusterRollOutSequenceTags` and `DescribeClusterRollOutSequenceTags` (already present in vendored SDK — no SDK changes required).
- No breaking changes; purely additive.
