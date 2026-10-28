## ADDED Requirements

### Requirement: Create KMS instance via Terraform framework action

The provider SHALL provide a terraform-plugin-framework action `tencentcloud_kms_instance_operation` that wraps the KMS cloud API `CreateKmsInstance` to purchase and activate a prepaid KMS instance. The action SHALL only implement the Invoke operation; no Read, Update, or Delete is applicable.

#### Scenario: Invoke with required and optional parameters

- **WHEN** a user invokes `tencentcloud_kms_instance_operation` with `instance_charge_type` set to `PREPAID`, `product_version` set to `professional`, a `client_token`, an `instance_charge_prepaid` block containing `period` and `renew_flag`, and a `tags` list with `tag_key`/`tag_value` entries
- **THEN** the provider SHALL call `CreateKmsInstance` with the corresponding request fields and SHALL log the returned `instance_id`, `deal_name`, and `bill_id` for observability

#### Scenario: Empty response from CreateKmsInstance

- **WHEN** `CreateKmsInstance` returns a response with a nil `Response` or an empty `InstanceId`
- **THEN** the provider SHALL return an error diagnostic and SHALL NOT report success

#### Scenario: API call failure with retry

- **WHEN** `CreateKmsInstance` returns an error
- **THEN** the provider SHALL wrap the error with `tccommon.RetryError` and retry up to `WriteRetryTimeout`

### Requirement: Action registration and documentation

The provider SHALL register `tencentcloud_kms_instance_operation` in `tencentcloud/framework/registry.go` under the `actionFactories` slice. It SHALL NOT be registered in the SDKv2 `provider.go` ResourcesMap. The provider SHALL list the action in `tencentcloud/provider.md` under a new `Action` sub-section of the Key Management Service(KMS) product node. An action documentation file `action_tc_kms_instance_operation.md` SHALL be created with a one-line description mentioning KMS, a NOTE requiring Terraform 1.14+, an Example Usage section, and no Import section.

#### Scenario: Framework registration

- **WHEN** the provider is built
- **THEN** `tencentcloud_kms_instance_operation` SHALL map to `kms.NewKmsInstanceOperation()` in the framework `actionFactories`

#### Scenario: Documentation availability

- **WHEN** documentation is generated via `make doc`
- **THEN** the `tencentcloud_kms_instance_operation` action page SHALL exist under `website/docs/actions/` and SHALL NOT contain a manually-written `Argument Reference` or `Attribute Reference` section

### Requirement: Unit tests with mocked cloud API

A unit test file `action_tc_kms_instance_operation_test.go` SHALL be created using gomonkey to mock the KMS cloud API client and tftypes raw values to build InvokeRequest, testing business logic only without invoking real cloud APIs.

#### Scenario: Mocked invoke success

- **WHEN** the unit test runs Invoke with mocked `UseKmsClient` and `CreateKmsInstanceWithContext` returning a valid response
- **THEN** the test SHALL verify no error diagnostics are produced and the request fields are filled correctly