# dbdc-db-custom-cluster-deletion-protection Specification

## Purpose
TBD - created by archiving change add-dbdc-db-custom-cluster-deletion-protection. Update Purpose after archive.
## Requirements
### Requirement: Deletion protection on cluster creation
The `tencentcloud_dbdc_db_custom_cluster` resource SHALL support an optional `deletion_protection` parameter (TypeBool) that is passed to the `CreateDBCustomCluster` API as `DeletionProtection`. When the user does not specify this parameter, the provider SHALL NOT set `DeletionProtection` in the API request, allowing the API to apply its default value of `true`.

#### Scenario: Create cluster with deletion protection enabled
- **WHEN** a user specifies `deletion_protection = true` in the `tencentcloud_dbdc_db_custom_cluster` resource configuration
- **THEN** the provider SHALL pass `DeletionProtection=true` in the `CreateDBCustomCluster` API request

#### Scenario: Create cluster with deletion protection disabled
- **WHEN** a user specifies `deletion_protection = false` in the `tencentcloud_dbdc_db_custom_cluster` resource configuration
- **THEN** the provider SHALL pass `DeletionProtection=false` in the `CreateDBCustomCluster` API request

#### Scenario: Create cluster without explicit deletion protection
- **WHEN** a user does NOT specify `deletion_protection` in the `tencentcloud_dbdc_db_custom_cluster` resource configuration
- **THEN** the provider SHALL NOT set `DeletionProtection` in the `CreateDBCustomCluster` API request (API defaults to `true`)

### Requirement: Read deletion protection from cluster detail
The `tencentcloud_dbdc_db_custom_cluster` resource SHALL read `DeletionProtection` from the `DescribeDBCustomClusterDetail` API response and populate the `deletion_protection` field in Terraform state.

#### Scenario: Read existing cluster with deletion protection
- **WHEN** the provider reads an existing `tencentcloud_dbdc_db_custom_cluster` resource
- **AND** the `DescribeDBCustomClusterDetail` response contains `DeletionProtection` that is not nil
- **THEN** the provider SHALL set `deletion_protection` in state to the value from the response

#### Scenario: Read existing cluster when DeletionProtection is nil
- **WHEN** the provider reads an existing `tencentcloud_dbdc_db_custom_cluster` resource
- **AND** the `DescribeDBCustomClusterDetail` response contains `DeletionProtection` that is nil
- **THEN** the provider SHALL NOT set `deletion_protection` in state

### Requirement: Update deletion protection via ModifyDBCustomClusterAttributes
The `tencentcloud_dbdc_db_custom_cluster` resource SHALL support updating the `deletion_protection` parameter by calling the `ModifyDBCustomClusterAttributes` API with `ClusterId` and `DeletionProtection` fields. The update SHALL be triggered when `d.HasChange("deletion_protection")` is true.

#### Scenario: Update deletion protection from true to false
- **WHEN** a user changes `deletion_protection` from `true` to `false` on an existing `tencentcloud_dbdc_db_custom_cluster` resource
- **THEN** the provider SHALL call `ModifyDBCustomClusterAttributes` with `ClusterId` set to the resource ID and `DeletionProtection=false`

#### Scenario: Update deletion protection from false to true
- **WHEN** a user changes `deletion_protection` from `false` to `true` on an existing `tencentcloud_dbdc_db_custom_cluster` resource
- **THEN** the provider SHALL call `ModifyDBCustomClusterAttributes` with `ClusterId` set to the resource ID and `DeletionProtection=true`

#### Scenario: Update without deletion protection change
- **WHEN** a user updates the `tencentcloud_dbdc_db_custom_cluster` resource without changing `deletion_protection`
- **THEN** the provider SHALL NOT call `ModifyDBCustomClusterAttributes` for deletion protection

