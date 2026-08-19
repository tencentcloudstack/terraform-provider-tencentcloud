## ADDED Requirements

### Requirement: Manage cluster roll-out sequence tags

The resource `tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config` SHALL manage the roll-out sequence tags bound to a single TKE cluster. The resource unique ID MUST be the cluster ID (`ClusterID`). The set of configurable arguments MUST match the input parameters of the `ModifyClusterRollOutSequenceTags` API (`ClusterID` and a list of `Tags`, each tag being a `Key`/`Value` pair).

#### Scenario: Create tags for a cluster

- **WHEN** the user applies a configuration with `cluster_id` and one or more `tags` blocks
- **THEN** the resource calls `ModifyClusterRollOutSequenceTags` with the given `ClusterID` and `Tags`
- **AND** sets the resource ID to the `cluster_id`
- **AND** the create call is wrapped with a retry mechanism

#### Scenario: Read tags back from the cluster

- **WHEN** the resource is read
- **THEN** the resource calls `DescribeClusterRollOutSequenceTags` filtered by the cluster ID using the maximum page size and pagination
- **AND** populates `cluster_id` and `tags` from the matching `ClusterTags` entry
- **AND** all returned-value access is nil-safe

#### Scenario: Cluster no longer has tags during read

- **WHEN** the describe result contains no matching cluster tag entry (or an empty tag list) for the cluster ID
- **THEN** the resource removes itself from state by setting an empty ID

### Requirement: Update cluster roll-out sequence tags

The resource SHALL converge the cluster's roll-out sequence tags to the desired `tags` list on update by reusing `ModifyClusterRollOutSequenceTags`.

#### Scenario: Modify the tag list

- **WHEN** the user changes the `tags` argument and applies
- **THEN** the resource calls `ModifyClusterRollOutSequenceTags` with the cluster ID and the new full `Tags` list
- **AND** the update call is wrapped with a retry mechanism

### Requirement: Delete cluster roll-out sequence tags

The resource SHALL remove all roll-out sequence tags from the cluster on delete by calling `ModifyClusterRollOutSequenceTags` with an empty `Tags` list.

#### Scenario: Delete removes all tags

- **WHEN** the resource is destroyed
- **THEN** the resource calls `ModifyClusterRollOutSequenceTags` with the cluster ID and an empty `Tags` list
- **AND** the delete call is wrapped with a retry mechanism
