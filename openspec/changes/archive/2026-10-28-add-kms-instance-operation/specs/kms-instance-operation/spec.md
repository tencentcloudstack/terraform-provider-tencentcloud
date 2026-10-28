## ADDED Requirements

### Requirement: Create KMS instance via Terraform operation resource

The provider SHALL provide a resource `tencentcloud_kms_instance_operation` of type RESOURCE_KIND_OPERATION that wraps the KMS cloud API `CreateKmsInstance` to purchase and activate a prepaid KMS instance. The resource SHALL only implement the Create operation; Read, Update, and Delete SHALL be empty no-ops.

#### Scenario: Create a prepaid KMS instance with required parameters

- **WHEN** a user applies a `tencentcloud_kms_instance_operation` resource with `instance_charge_type` set to `PREPAID`, `product_version` set to `professional`, and an `instance_charge_prepaid` block containing `period` and `renew_flag`
- **THEN** the provider SHALL call `CreateKmsInstance` with the corresponding request fields and SHALL set `instance_id`, `deal_name`, and `bill_id` as computed attributes on the resource state

#### Scenario: Create with optional client token and tags

- **WHEN** a user provides `client_token` and a `tags` list with `tag_key`/`tag_value` entries
- **THEN** the provider SHALL pass `ClientToken` and `Tags` array to `CreateKmsInstance`

#### Scenario: Empty response from CreateKmsInstance

- **WHEN** `CreateKmsInstance` returns a nil `Response`, a nil `InstanceId`, or an empty `InstanceId` string
- **THEN** the provider SHALL return a `NonRetryableError` and SHALL NOT set an empty id into the state

#### Scenario: API call failure with retry

- **WHEN** `CreateKmsInstance` returns an error
- **THEN** the provider SHALL wrap the error with `tccommon.RetryError` and retry up to `WriteRetryTimeout`, and SHALL NOT perform any SetId inside the retry block

#### Scenario: Read and delete are no-ops

- **WHEN** terraform invokes Read, Update, or Delete on the resource
- **THEN** the provider SHALL return nil without performing any cloud API call, because this is a one-time operation resource that records no state

### Requirement: Resource registration and documentation

The provider SHALL register `tencentcloud_kms_instance_operation` in `tencentcloud/provider.go` and list it in `tencentcloud/provider.md` under the Key Management Service(KMS) Resource section. A resource documentation file `resource_tc_kms_kms_instance_operation.md` SHALL be created with a one-line description mentioning KMS, an Example Usage section, and no Import section.

#### Scenario: Provider registration

- **WHEN** the provider is built
- **THEN** `tencentcloud_kms_instance_operation` SHALL map to `kms.ResourceTencentCloudKmsInstanceOperation()` in the provider resource map

#### Scenario: Documentation availability

- **WHEN** documentation is generated via `make doc`
- **THEN** the `tencentcloud_kms_instance_operation` resource page SHALL exist and SHALL NOT contain a manually-written `Argument Reference` or `Attribute Reference` section

### Requirement: Unit tests with mocked cloud API

A unit test file `resource_tc_kms_kms_instance_operation_test.go` SHALL be created using gomonkey to mock the KMS cloud API client, testing business logic only without invoking real cloud APIs or the terraform test suite.

#### Scenario: Mocked create success

- **WHEN** the unit test runs the create function with mocked `UseKmsClient` and `CreateKmsInstanceWithContext` returning a valid `InstanceId`
- **THEN** the test SHALL verify the resource id is set to the returned `InstanceId` without making real API calls