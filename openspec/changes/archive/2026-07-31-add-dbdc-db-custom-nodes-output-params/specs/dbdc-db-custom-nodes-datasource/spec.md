## MODIFIED Requirements

### Requirement: Datasource schema definition
The system SHALL define a Terraform datasource schema for `tencentcloud_dbdc_db_custom_nodes` with the following input parameters:
- `node_ids`: Optional, TypeList of string - filter by one or more node IDs (max 100 per request)
- `filters`: Optional, TypeList of schema.Resource with sub-fields `name` (Required, string) and `values` (Required, TypeSet of string) - filter by cluster-id, node-name, status, zone
- `tags`: Optional, TypeList of schema.Resource with sub-fields `key` (Required, string) and `value` (Required, string) - filter by tag key-value pairs
- `result_output_file`: Optional, TypeString - used to save results to file

The system SHALL define computed output parameters:
- `node_set`: Computed, TypeList of schema.Resource containing flattened node attributes
- Each node element SHALL contain: `node_id`, `node_name`, `ssh_endpoint`, `lan_ip`, `cluster_id`, `zone`, `node_type`, `cpu`, `memory`, `system_disk` (TypeList of schema.Resource with `disk_type`, `disk_size`), `data_disks` (TypeList of schema.Resource with `disk_type`, `disk_size`, `disk_name`), `os_name`, `image_id`, `vpc_id`, `subnet_id`, `status`, `charge_type`, `expire_time`, `created_time`, `isolated_time`, `tags` (TypeList of schema.Resource with `key`, `value`), `auto_renew`, `switch_id`, `rack_id`, `host_ip`, `network_mode`, `eni_ip`

#### Scenario: Datasource with filters input
- **WHEN** a user provides `filters` with `name` = "cluster-id" and `values` = ["cluster-123"]
- **THEN** the system SHALL call DescribeDBCustomNodes API with the corresponding Filter parameter and return matching nodes

#### Scenario: Datasource with node_ids input
- **WHEN** a user provides `node_ids` = ["node-1", "node-2"]
- **THEN** the system SHALL call DescribeDBCustomNodes API with NodeIds parameter containing those IDs

#### Scenario: Datasource with tags input
- **WHEN** a user provides `tags` with `key` = "env" and `value` = "prod"
- **THEN** the system SHALL call DescribeDBCustomNodes API with Tags parameter containing that key-value pair

#### Scenario: Datasource with no filter inputs
- **WHEN** a user provides no filter inputs (node_ids, filters, tags all empty)
- **THEN** the system SHALL call DescribeDBCustomNodes API without filters and return all nodes the account has access to

### Requirement: Nil field handling in response
The system SHALL check each response field for nil before calling `d.Set()` or adding it to the output map. Fields that may be nil according to API documentation (`SystemDisk`, `DataDisks`, `Tags`, `NetworkMode`, `EniIP` in DBCustomNode) SHALL be skipped when nil rather than set to empty values.

#### Scenario: Node with nil SystemDisk
- **WHEN** a DBCustomNode has nil `SystemDisk` field
- **THEN** the system SHALL skip setting `system_disk` for that node element

#### Scenario: Node with nil DataDisks
- **WHEN** a DBCustomNode has nil `DataDisks` field
- **THEN** the system SHALL skip setting `data_disks` for that node element

#### Scenario: Node with nil NetworkMode
- **WHEN** a DBCustomNode has nil `NetworkMode` field
- **THEN** the system SHALL skip setting `network_mode` for that node element

#### Scenario: Node with nil EniIP
- **WHEN** a DBCustomNode has nil `EniIP` field
- **THEN** the system SHALL skip setting `eni_ip` for that node element

## ADDED Requirements

### Requirement: Network mode and ENI IP output fields
The system SHALL expose two additional computed string fields inside each `node_set` element:
- `network_mode`: mapped from `DBCustomNode.NetworkMode`, representing the network mode enum (`NetworkModePrivateLink` or `NetworkModeCrossTenantENI`)
- `eni_ip`: mapped from `DBCustomNode.EniIP`, representing the node access IP address when the network mode is `NetworkModeCrossTenantENI`

#### Scenario: Node with NetworkModeCrossTenantENI
- **WHEN** the DescribeDBCustomNodes API returns a DBCustomNode with `NetworkMode` = "NetworkModeCrossTenantENI" and `EniIP` = "10.0.0.5"
- **THEN** the system SHALL set `network_mode` = "NetworkModeCrossTenantENI" and `eni_ip` = "10.0.0.5" for that node element in `node_set`

#### Scenario: Node with NetworkModePrivateLink and nil EniIP
- **WHEN** the DescribeDBCustomNodes API returns a DBCustomNode with `NetworkMode` = "NetworkModePrivateLink" and `EniIP` = nil
- **THEN** the system SHALL set `network_mode` = "NetworkModePrivateLink" and SHALL skip setting `eni_ip` for that node element
