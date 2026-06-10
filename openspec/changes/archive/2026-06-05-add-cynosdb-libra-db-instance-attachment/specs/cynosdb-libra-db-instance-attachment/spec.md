## ADDED Requirements

### Requirement: Create LibraDB instance attachment
The system SHALL create a LibraDB read-only analytics engine instance and attach it to a CynosDB cluster by calling the `AddLibraDBInstances` API. The resource SHALL accept all parameters defined in the API including `cluster_id`, `zone`, `cpu`, `mem`, `storage_size`, `pay_mode`, `objects`, `port`, `goods_num`, `instance_name`, `replicas_num`, `instance_type`, `storage_type`, `auto_voucher`, `order_source`, `deal_mode`, `vpc_id`, `subnet_id`, `security_group_ids`, `libra_db_version`, `time_span`, `time_unit`, and `src_instance_id`. After successful creation, the system SHALL store the composite ID (`cluster_id` + FILED_SP + `instance_id`) from the response's `ResourceIds` field. The system SHALL verify that the response and `ResourceIds` are not empty, returning a NonRetryableError otherwise. After creation, the system SHALL poll `DescribeLibraDBInstanceDetail` until the instance is in a ready state.

#### Scenario: Successful creation of LibraDB instance
- **WHEN** user applies a Terraform configuration with valid `tencentcloud_cynosdb_libra_db_instance` resource parameters including required `cluster_id`, `zone`, `cpu`, `mem`, and `storage_size`
- **THEN** the system calls `AddLibraDBInstances` API, extracts the instance ID from `ResourceIds`, sets the composite resource ID, polls until ready, and stores computed attributes (`big_deal_ids`, `tran_id`, `deal_names`, `resource_ids`)

#### Scenario: Create response has empty ResourceIds
- **WHEN** the `AddLibraDBInstances` API returns a response with nil or empty `ResourceIds`
- **THEN** the system returns a NonRetryableError indicating the instance ID was not returned

### Requirement: Read LibraDB instance detail
The system SHALL read the LibraDB instance detail by calling `DescribeLibraDBInstanceDetail` API with `cluster_id` and `instance_id` extracted from the composite resource ID. The system SHALL set all computed attributes from the response including `uin`, `app_id`, `cluster_name`, `instance_name`, `project_id`, `region`, `zone`, `status`, `status_desc`, `libra_db_version`, `cpu`, `memory`, `storage`, `storage_type`, `instance_type`, `instance_role`, `update_time`, `create_time`, `pay_mode`, `period_start_time`, `period_end_time`, `renew_flag`, `net_type`, `vpc_id`, `subnet_id`, `vip`, `vport`, `instance_net_info`, `resource_tags`, `node_info`, `node_count`, and `analysis_upgrade_version_info`. The system SHALL check each response field for nil before calling the corresponding `d.Set()`.

#### Scenario: Successful read of existing instance
- **WHEN** the resource exists and `DescribeLibraDBInstanceDetail` returns valid data
- **THEN** the system sets all non-nil computed attributes in the Terraform state

#### Scenario: Instance not found during read
- **WHEN** `DescribeLibraDBInstanceDetail` returns an error indicating the instance does not exist
- **THEN** the system removes the resource from state by calling `d.SetId("")` and returns nil

### Requirement: Delete LibraDB instance attachment
The system SHALL delete (isolate) the LibraDB analytics cluster by calling `IsolateLibraDBCluster` API with `cluster_id` extracted from the composite resource ID. The system SHALL accept optional `isolate_reason_types` (list of integers) and `isolate_reason` (string) parameters for the delete operation.

#### Scenario: Successful deletion of LibraDB instance
- **WHEN** user destroys the `tencentcloud_cynosdb_libra_db_instance` resource
- **THEN** the system calls `IsolateLibraDBCluster` API with the `cluster_id` and optional isolation reason parameters

#### Scenario: Instance already isolated
- **WHEN** `IsolateLibraDBCluster` returns an error indicating the cluster is already isolated or not found
- **THEN** the system treats the deletion as successful

### Requirement: Resource ID uses composite format
The system SHALL use a composite ID format of `cluster_id` + `tccommon.FILED_SP` + `instance_id` for the resource identifier. The system SHALL parse this composite ID in Read and Delete operations to extract both `cluster_id` and `instance_id`.

#### Scenario: Import with composite ID
- **WHEN** user imports the resource using `terraform import tencentcloud_cynosdb_libra_db_instance.example cluster_id#instance_id`
- **THEN** the system correctly parses the composite ID and reads the instance detail

### Requirement: Immutable parameters enforcement
The system SHALL implement an update function that checks all top-level parameters (except `cluster_id` which is ForceNew) against an immutableArgs list. If any immutable parameter is changed, the system SHALL return an error indicating the parameter cannot be modified.

#### Scenario: Attempt to modify immutable parameter
- **WHEN** user modifies a parameter like `cpu` or `mem` in the Terraform configuration
- **THEN** the system returns an error stating the parameter is immutable and cannot be changed

### Requirement: Retry logic for API calls
The system SHALL wrap all API calls with `resource.Retry` using `tccommon.ReadRetryTimeout` for read operations and `tccommon.WriteRetryTimeout` for write operations. Failed API calls SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: Transient API failure during create
- **WHEN** `AddLibraDBInstances` API call fails with a transient error
- **THEN** the system retries the call within the configured timeout period

#### Scenario: Transient API failure during read
- **WHEN** `DescribeLibraDBInstanceDetail` API call fails with a transient error
- **THEN** the system retries the call within the configured timeout period
