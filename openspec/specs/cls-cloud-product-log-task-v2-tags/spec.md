# cls-cloud-product-log-task-v2-tags Specification

## Purpose
TBD - created by archiving change add-cls-cloud-product-log-task-v2-tags. Update Purpose after archive.
## Requirements
### Requirement: Support tags parameter on cls cloud product log task v2 resource

The `tencentcloud_cls_cloud_product_log_task_v2` resource SHALL support an optional `tags` parameter that binds tags to the logset and topic associated with the cloud product log collection task. Tags SHALL be passed to the `CreateCloudProductLogCollection` API on creation and to the `ModifyCloudProductLogCollection` API on update. Tags MUST NOT force resource recreation; they MUST be modifiable in-place.

**Rationale**: The cloud API already accepts a `Tags` field on both create and modify operations to bind tags to the associated logset and topic. Exposing this through Terraform enables full lifecycle management of tags without manual console operations.

#### Scenario: Create a cloud product log task with tags

- **WHEN** a user creates a `tencentcloud_cls_cloud_product_log_task_v2` resource specifying:
```hcl
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-0an6hpv3"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"

  tags = {
    Environment = "production"
    Team        = "backend"
  }
}
```
- **THEN** the resource is created and the `Tags` field is passed to the `CreateCloudProductLogCollection` API request
- **AND** the tags are bound to the associated logset and topic
- **AND** terraform plan shows no changes after apply

#### Scenario: Create a cloud product log task without tags

- **WHEN** a user creates a `tencentcloud_cls_cloud_product_log_task_v2` resource without specifying the `tags` parameter
- **THEN** the resource is created and the `Tags` field is not set on the `CreateCloudProductLogCollection` API request
- **AND** existing behavior is preserved (backward compatible)

#### Scenario: Update tags on an existing cloud product log task

- **WHEN** a user updates the `tags` parameter on an existing `tencentcloud_cls_cloud_product_log_task_v2` resource (adding, removing, or modifying tag key-value pairs)
- **THEN** the `ModifyCloudProductLogCollection` API is called with the new `Tags` value
- **AND** the resource is NOT recreated
- **AND** the updated tags are applied to the associated logset and topic

#### Scenario: Read tags from an existing cloud product log task

- **WHEN** the resource Read operation executes for a `tencentcloud_cls_cloud_product_log_task_v2` resource that has tags bound
- **THEN** the tags are read back from the API response and populated into the Terraform state
- **AND** the `tags` field in state reflects the tags bound to the associated logset and topic

#### Scenario: Remove all tags from an existing cloud product log task

- **WHEN** a user removes the `tags` parameter (or sets it to an empty map) on an existing `tencentcloud_cls_cloud_product_log_task_v2` resource
- **THEN** the `ModifyCloudProductLogCollection` API is called with an empty `Tags` value
- **AND** the tags are removed from the associated logset and topic
- **AND** the resource is NOT recreated

