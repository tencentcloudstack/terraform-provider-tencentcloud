## ADDED Requirements

### Requirement: Add cluster_type parameter

The `tencentcloud_tcaplus_cluster` resource SHALL accept an optional `cluster_type` top-level schema parameter (`TypeInt`, Optional, Computed, immutable after creation) mapped to the TcaplusDB `CreateCluster` API `ClusterType` request field (`*int64`, 1=shared, 2=dedicated). On Read, the resource SHALL refresh `cluster_type` from the `DescribeClusters` response `Clusters.ClusterType` field when it is not nil.

#### Scenario: Create cluster with cluster_type unspecified
- **WHEN** `cluster_type` is not set in the Terraform configuration
- **THEN** the resource SHALL create a cluster without passing `ClusterType` to the `CreateCluster` API (preserving backward-compatible shared-cluster behavior) and SHALL populate `cluster_type` from the `DescribeClusters` response on Read

#### Scenario: Create dedicated cluster with cluster_type set to 2
- **WHEN** `cluster_type` is set to `2` in the Terraform configuration
- **THEN** the resource SHALL pass `ClusterType=2` to the `CreateCluster` API to create a dedicated cluster and SHALL refresh `cluster_type` from the `DescribeClusters` response on Read

#### Scenario: cluster_type is immutable after creation
- **WHEN** `cluster_type` is changed in the Terraform configuration after the cluster is created
- **THEN** the Update function SHALL detect the change via the `immutableArgs` check and return a clear error rejecting the modification (the cluster SHALL NOT be destroyed and recreated)

### Requirement: Add resource_tags nested block

The `tencentcloud_tcaplus_cluster` resource SHALL accept an optional `resource_tags` nested block (`TypeList`, Optional, immutable after creation) mapped to the TcaplusDB `CreateCluster` API `ResourceTags` request field (`[]*TagInfoUnit`). Each element SHALL expose `tag_key` (string) and `tag_value` (string) sub-fields. The resource SHALL pass the configured tags to the `CreateCluster` API on creation. Because the `DescribeClusters` response does not return cluster-level tags, the `resource_tags` block SHALL NOT be refreshed on Read and the configured value SHALL remain authoritative in Terraform state.

#### Scenario: Create cluster with resource tags
- **WHEN** the `resource_tags` block contains one or more `tag_key`/`tag_value` entries
- **THEN** the resource SHALL pass each entry as a `TagInfoUnit` to the `CreateCluster` API `ResourceTags` field

#### Scenario: Create cluster without resource tags
- **WHEN** the `resource_tags` block is absent from the configuration
- **THEN** the resource SHALL not set `ResourceTags` on the `CreateCluster` API request (preserving backward compatibility)

#### Scenario: resource_tags is immutable after creation
- **WHEN** `resource_tags` is changed in the Terraform configuration after the cluster is created
- **THEN** the Update function SHALL detect the change via the `immutableArgs` check and return a clear error rejecting the modification

### Requirement: Add server_list nested block

The `tencentcloud_tcaplus_cluster` resource SHALL accept an optional `server_list` nested block (`TypeList`, Optional, Computed, immutable after creation) mapped to the TcaplusDB `CreateCluster` API `ServerList` request field (`[]*MachineInfo`). For creation, each element SHALL expose `machine_type` (string) and `machine_num` (int) sub-fields, which SHALL be passed to the `CreateCluster` API. On Read, the `server_list` block SHALL be refreshed from the `DescribeClusters` response `Clusters.ServerList` (`[]*ServerDetailInfo`), populating `server_uid` (string, Computed), `machine_type` (string, Computed), `memory_rate` (int, Computed), `disk_rate` (int, Computed), `read_num` (int, Computed), `write_num` (int, Computed), and `version` (string, Computed) with nil checks before each set.

#### Scenario: Create dedicated cluster with server_list
- **WHEN** the `server_list` block contains one or more entries with `machine_type` and `machine_num`
- **THEN** the resource SHALL pass each entry as a `MachineInfo` to the `CreateCluster` API `ServerList` field

#### Scenario: Read server_list from DescribeClusters
- **WHEN** the `DescribeClusters` response returns a non-empty `ServerList` for the cluster
- **THEN** the resource SHALL populate the `server_list` block from `ServerDetailInfo` elements, setting each available field (`server_uid`, `machine_type`, `memory_rate`, `disk_rate`, `read_num`, `write_num`, `version`) only when the corresponding response field is not nil

#### Scenario: server_list is immutable after creation
- **WHEN** `server_list` is changed in the Terraform configuration after the cluster is created
- **THEN** the Update function SHALL detect the change via the `immutableArgs` check and return a clear error rejecting the modification

### Requirement: Add proxy_list nested block

The `tencentcloud_tcaplus_cluster` resource SHALL accept an optional `proxy_list` nested block (`TypeList`, Optional, Computed, immutable after creation) mapped to the TcaplusDB `CreateCluster` API `ProxyList` request field (`[]*MachineInfo`). For creation, each element SHALL expose `machine_type` (string) and `machine_num` (int) sub-fields, which SHALL be passed to the `CreateCluster` API. On Read, the `proxy_list` block SHALL be refreshed from the `DescribeClusters` response `Clusters.ProxyList` (`[]*ProxyDetailInfo`), populating `proxy_uid` (string, Computed), `machine_type` (string, Computed), `process_speed` (int, Computed), `average_process_delay` (int, Computed), `slow_process_speed` (int, Computed), and `version` (string, Computed) with nil checks before each set.

#### Scenario: Create dedicated cluster with proxy_list
- **WHEN** the `proxy_list` block contains one or more entries with `machine_type` and `machine_num`
- **THEN** the resource SHALL pass each entry as a `MachineInfo` to the `CreateCluster` API `ProxyList` field

#### Scenario: Read proxy_list from DescribeClusters
- **WHEN** the `DescribeClusters` response returns a non-empty `ProxyList` for the cluster
- **THEN** the resource SHALL populate the `proxy_list` block from `ProxyDetailInfo` elements, setting each available field (`proxy_uid`, `machine_type`, `process_speed`, `average_process_delay`, `slow_process_speed`, `version`) only when the corresponding response field is not nil

#### Scenario: proxy_list is immutable after creation
- **WHEN** `proxy_list` is changed in the Terraform configuration after the cluster is created
- **THEN** the Update function SHALL detect the change via the `immutableArgs` check and return a clear error rejecting the modification

### Requirement: Extend CreateCluster service-layer function

The TcaplusDB service-layer `CreateCluster` function SHALL accept the new parameters (`resourceTags []*tcaplusdb.TagInfoUnit`, `serverList []*tcaplusdb.MachineInfo`, `proxyList []*tcaplusdb.MachineInfo`, `clusterType int64`) in addition to the existing parameters, and SHALL pass them into the `CreateClusterRequest` only when non-empty/non-zero, preserving backward compatibility for existing callers. The function return type SHALL remain `(string, error)` (the cluster id).

#### Scenario: CreateCluster called with all new parameters empty
- **WHEN** `resourceTags`, `serverList`, `proxyList` are empty and `clusterType` is 0
- **THEN** the service-layer function SHALL not set `ResourceTags`, `ServerList`, `ProxyList`, or `ClusterType` on the `CreateClusterRequest`, producing the same request as before this change

#### Scenario: CreateCluster called with dedicated cluster parameters
- **WHEN** `clusterType` is 2 and `serverList`/`proxyList` contain entries and `resourceTags` contains entries
- **THEN** the service-layer function SHALL set `ClusterType`, `ServerList`, `ProxyList`, and `ResourceTags` on the `CreateClusterRequest` before invoking the API
