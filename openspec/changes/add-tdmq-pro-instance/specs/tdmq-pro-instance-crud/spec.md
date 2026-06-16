## ADDED Requirements

### Requirement: Create TDMQ professional cluster instance
The system SHALL create a TDMQ professional cluster instance by calling the `CreateProCluster` API with the following parameters: `zone_ids` (required), `product_name` (required), `storage_size` (required), `vpc` (required, containing `vpc_id` and `subnet_id`), `cluster_name` (optional), `time_span` (optional, default 1), `auto_renew_flag` (optional), and `auto_voucher` (optional). Upon successful creation, the system SHALL store the returned `ClusterId` as the resource ID.

#### Scenario: Successful creation with all required parameters
- **WHEN** user provides `zone_ids`, `product_name`, `storage_size`, and `vpc` configuration in the Terraform resource block
- **THEN** the system calls `CreateProCluster` API and sets the resource ID to the returned `ClusterId`

#### Scenario: Creation returns empty response
- **WHEN** the `CreateProCluster` API returns a nil response or empty `ClusterId`
- **THEN** the system SHALL return a `NonRetryableError` indicating the creation failed

#### Scenario: Creation with optional parameters
- **WHEN** user provides `cluster_name`, `time_span`, `auto_renew_flag`, and `auto_voucher` in addition to required parameters
- **THEN** the system passes all parameters to the `CreateProCluster` API

### Requirement: Read TDMQ professional cluster instance
The system SHALL read a TDMQ professional cluster instance by calling the `DescribeClusters` API with the `ClusterIdList` filter set to the resource's `ClusterId`. The system SHALL populate the resource state with the cluster's attributes from the response.

#### Scenario: Successful read of existing cluster
- **WHEN** the `DescribeClusters` API returns a `ClusterSet` containing the target cluster
- **THEN** the system SHALL set the resource attributes (`cluster_name`, `remark`, `public_access_enabled`, `status`, `cluster_id`, etc.) from the response

#### Scenario: Cluster not found
- **WHEN** the `DescribeClusters` API returns an empty `ClusterSet` or nil response
- **THEN** the system SHALL log the cluster ID for debugging, then call `d.SetId("")` to remove the resource from state

### Requirement: Update TDMQ professional cluster instance
The system SHALL update a TDMQ professional cluster instance by calling the `ModifyCluster` API. Only the following fields are updatable: `cluster_name`, `remark`, and `public_access_enabled`.

#### Scenario: Update cluster name
- **WHEN** user changes the `cluster_name` attribute
- **THEN** the system calls `ModifyCluster` API with the new `ClusterName` value

#### Scenario: Update remark
- **WHEN** user changes the `remark` attribute
- **THEN** the system calls `ModifyCluster` API with the new `Remark` value

#### Scenario: Enable public access
- **WHEN** user sets `public_access_enabled` to true
- **THEN** the system calls `ModifyCluster` API with `PublicAccessEnabled` set to true

#### Scenario: Update triggers read
- **WHEN** the `ModifyCluster` API call succeeds
- **THEN** the system SHALL call the Read function to refresh the resource state

### Requirement: Delete TDMQ professional cluster instance
The system SHALL delete a TDMQ professional cluster instance by calling the `DeleteCluster` API with the resource's `ClusterId`.

#### Scenario: Successful deletion
- **WHEN** the `DeleteCluster` API call succeeds
- **THEN** the system SHALL remove the resource from Terraform state

#### Scenario: Cluster already deleted
- **WHEN** the `DeleteCluster` API returns a not-found error
- **THEN** the system SHALL treat it as a successful deletion

### Requirement: Import TDMQ professional cluster instance
The system SHALL support importing an existing TDMQ professional cluster by its `ClusterId`.

#### Scenario: Import by cluster ID
- **WHEN** user runs `terraform import tencentcloud_tdmq_pro_instance.example <cluster_id>`
- **THEN** the system SHALL call the Read function to populate the resource state from the existing cluster

### Requirement: Create-only fields trigger recreation
The system SHALL mark the following fields as `ForceNew`: `zone_ids`, `product_name`, `auto_renew_flag`, `time_span`, `auto_voucher`, `storage_size`, `vpc`. Changing any of these fields SHALL trigger resource destruction and recreation.

#### Scenario: Change zone_ids triggers recreation
- **WHEN** user modifies the `zone_ids` attribute
- **THEN** Terraform plan SHALL show the resource will be destroyed and recreated

#### Scenario: Change storage_size triggers recreation
- **WHEN** user modifies the `storage_size` attribute
- **THEN** Terraform plan SHALL show the resource will be destroyed and recreated
