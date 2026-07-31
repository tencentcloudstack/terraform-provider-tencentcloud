# dbdc-db-custom-node-create-params Specification

## Purpose
TBD - created by archiving change add-dbdc-db-custom-node-params. Update Purpose after archive.
## Requirements
### Requirement: Charge type on db custom node creation
The `tencentcloud_dbdc_db_custom_node` resource SHALL support an optional `charge_type` parameter (TypeString, Optional, Computed, ForceNew) that is passed to the `CreateDBCustomNodes` API as `ChargeType`. Valid values are `PREPAID` and `POSTPAID`. When omitted, the provider SHALL NOT set `ChargeType` in the request (the API defaults to `PREPAID`). The value SHALL be refreshed from `DescribeDBCustomNodes` (`DBCustomNode.ChargeType`) on Read.

#### Scenario: Create node with explicit charge_type
- **WHEN** a user specifies `charge_type = "PREPAID"` in the resource configuration
- **THEN** the provider SHALL pass `ChargeType=PREPAID` in the `CreateDBCustomNodes` request

#### Scenario: Create node without charge_type
- **WHEN** a user does NOT specify `charge_type`
- **THEN** the provider SHALL NOT set `ChargeType` in the request, and SHALL read the API-applied value back into state

#### Scenario: Changing charge_type forces recreation
- **WHEN** a user changes `charge_type` on an existing node
- **THEN** the provider SHALL recreate the resource (ForceNew) because no Modify API exists

### Requirement: Network mode on db custom node creation
The `tencentcloud_dbdc_db_custom_node` resource SHALL support an optional `network_mode` parameter (TypeString, Optional, ForceNew) passed to the `CreateDBCustomNodes` API as `NetworkMode`. Valid values are `privatelink` and `cross_tenant_eni`. When omitted, the provider SHALL NOT set `NetworkMode` (API defaults to `privatelink`). The value SHALL be refreshed from `DescribeDBCustomNodes` (`DBCustomNode.NetworkMode`) on Read.

#### Scenario: Create node with cross_tenant_eni network mode
- **WHEN** a user specifies `network_mode = "cross_tenant_eni"`
- **THEN** the provider SHALL pass `NetworkMode=cross_tenant_eni` in the request

#### Scenario: Create node without network_mode
- **WHEN** a user does NOT specify `network_mode`
- **THEN** the provider SHALL NOT set `NetworkMode`, and SHALL read the API-applied value back into state

#### Scenario: Changing network_mode forces recreation
- **WHEN** a user changes `network_mode` on an existing node
- **THEN** the provider SHALL recreate the resource (ForceNew)

### Requirement: System disk configuration on db custom node creation
The `tencentcloud_dbdc_db_custom_node` resource SHALL support an optional `system_disk` block (TypeList, MaxItems 1, Optional, Computed, ForceNew) whose `disk_type` and `disk_size` are passed to the `CreateDBCustomNodes` API as `SystemDisk.DiskType` and `SystemDisk.DiskSize`. When omitted, the provider SHALL NOT set `SystemDisk` (API applies its default). The value SHALL be refreshed from `DescribeDBCustomNodes` (`DBCustomNode.SystemDisk`) on Read.

#### Scenario: Create node with explicit system disk
- **WHEN** a user specifies a `system_disk` block with `disk_type = "CLOUD_HSSD"` and `disk_size = 100`
- **THEN** the provider SHALL set `SystemDisk={DiskType:"CLOUD_HSSD", DiskSize:100}` in the request

#### Scenario: Create node without system_disk
- **WHEN** a user does NOT specify `system_disk`
- **THEN** the provider SHALL NOT set `SystemDisk`, and SHALL read the API-applied value back into state

#### Scenario: Changing system_disk forces recreation
- **WHEN** a user changes `system_disk` on an existing node
- **THEN** the provider SHALL recreate the resource (ForceNew)

### Requirement: Data disks configuration on db custom node creation
The `tencentcloud_dbdc_db_custom_node` resource SHALL support an optional `data_disks` block (TypeList, Optional, Computed, ForceNew) whose entries' `disk_type`, `disk_size`, and `disk_name` are passed to the `CreateDBCustomNodes` API as `DataDisks[]`. When omitted, the provider SHALL NOT set `DataDisks` (API applies its default). The value SHALL be refreshed from `DescribeDBCustomNodes` (`DBCustomNode.DataDisks`) on Read.

#### Scenario: Create node with explicit data disks
- **WHEN** a user specifies a `data_disks` block with `disk_type`, `disk_size`, and `disk_name`
- **THEN** the provider SHALL build the `DataDisks` array and pass it in the request

#### Scenario: Create node without data_disks
- **WHEN** a user does NOT specify `data_disks`
- **THEN** the provider SHALL NOT set `DataDisks`, and SHALL read the API-applied value back into state

#### Scenario: Changing data_disks forces recreation
- **WHEN** a user changes `data_disks` on an existing node
- **THEN** the provider SHALL recreate the resource (ForceNew)

### Requirement: Host name on db custom node creation
The `tencentcloud_dbdc_db_custom_node` resource SHALL support an optional `host_name` parameter (TypeString, Optional, ForceNew) passed to the `CreateDBCustomNodes` API as `HostName`. Because `DescribeDBCustomNodes` does NOT return `HostName`, `host_name` SHALL NOT be Computed; the configured value SHALL be retained in state on Read (write-only). Import verification SHALL ignore `host_name`.

#### Scenario: Create node with host_name
- **WHEN** a user specifies `host_name = "my-node-1"`
- **THEN** the provider SHALL pass `HostName=my-node-1` in the request

#### Scenario: Read does not clear host_name
- **WHEN** the provider reads an existing node that was created with `host_name`
- **THEN** the provider SHALL leave the configured `host_name` value in state (it is not overwritten, because the API does not return it)

#### Scenario: Changing host_name forces recreation
- **WHEN** a user changes `host_name` on an existing node
- **THEN** the provider SHALL recreate the resource (ForceNew)

### Requirement: Security group ids on db custom node
The `tencentcloud_dbdc_db_custom_node` resource SHALL support an optional `security_group_ids` parameter (TypeList of String, Optional, mutable — NOT ForceNew) passed to the `CreateDBCustomNodes` API as `SecurityGroupIds` on Create. On Update, when `security_group_ids` changes, the provider SHALL call `ModifyDBCustomNodeSecurityGroups` with the full new list (set/overwrite semantic). On Read, the provider SHALL refresh `security_group_ids` from the `DescribeDBCustomNodeSecurityGroups` API (`Groups[].SecurityGroupId`), because `DBCustomNode` does not return security groups.

#### Scenario: Create node with security group ids
- **WHEN** a user specifies `security_group_ids = ["sg-xxx", "sg-yyy"]`
- **THEN** the provider SHALL set `SecurityGroupIds=["sg-xxx","sg-yyy"]` in the `CreateDBCustomNodes` request

#### Scenario: Create node without security_group_ids
- **WHEN** a user does NOT specify `security_group_ids`
- **THEN** the provider SHALL NOT set `SecurityGroupIds` in the request

#### Scenario: Update security group ids in place
- **WHEN** a user changes `security_group_ids` on an existing node
- **THEN** the provider SHALL call `ModifyDBCustomNodeSecurityGroups` with the full new list and SHALL NOT recreate the node

#### Scenario: Read refreshes security group ids
- **WHEN** the provider reads an existing node
- **THEN** the provider SHALL call `DescribeDBCustomNodeSecurityGroups` and populate `security_group_ids` from `Groups[].SecurityGroupId`

