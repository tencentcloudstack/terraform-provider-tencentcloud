# kms-instance-operation Specification

## Purpose

定义 KMS 服务 framework action `tencentcloud_kms_instance` 的行为规范，封装 `CreateKmsInstance` 云 API，用于购买并开通预付费 KMS 实例。该 action 为一次性操作（仅实现 Invoke），不涉及 Read / Update / Delete，操作完成后不持久化云侧状态。

## Requirements

### Requirement: Create KMS instance via Terraform framework action

The provider SHALL provide a terraform-plugin-framework action `tencentcloud_kms_instance` that wraps the KMS cloud API `CreateKmsInstance` to purchase and activate a prepaid KMS instance. The action SHALL only implement the Invoke operation; no Read, Update, or Delete is applicable.

#### Scenario: Invoke with required and optional parameters

- **WHEN** a user invokes `tencentcloud_kms_instance` with `instance_charge_type` set to `PREPAID`, `product_version` set to `professional`, a `client_token`, an `instance_charge_prepaid` block containing `period` and `renew_flag`, and a `tags` list with `tag_key`/`tag_value` entries
- **THEN** the provider SHALL call `CreateKmsInstance` with the corresponding request fields and SHALL log the returned `instance_id`, `deal_name`, and `bill_id` for observability

#### Scenario: Empty response from CreateKmsInstance

- **WHEN** `CreateKmsInstance` returns a response with a nil `Response`, or an `InstanceId` that is nil or empty
- **THEN** the provider SHALL log the empty result and return an error diagnostic, and SHALL NOT report success

#### Scenario: API call failure with retry

- **WHEN** `CreateKmsInstance` returns an error
- **THEN** the provider SHALL wrap the error with `tccommon.RetryError` and retry up to `WriteRetryTimeout`

### Requirement: Action registration, service helper and documentation

The provider SHALL register `tencentcloud_kms_instance` in `tencentcloud/framework/registry.go` under the `actionFactories` slice. It SHALL NOT be registered in the SDKv2 `provider.go` ResourcesMap. The cloud API call SHALL live in the shared service file `tencentcloud/services/kms/service_tencentcloud_kms.go` as `KmsService.CreateKmsInstance`, wrapped with `resource.Retry(tccommon.WriteRetryTimeout, ...)`. The provider SHALL list the action in `tencentcloud/provider.md` under a new `Action` sub-section of the Key Management Service(KMS) product node. An action documentation file `action_tc_kms_instance.md` SHALL be created with a one-line description mentioning KMS, a NOTE requiring Terraform 1.14+, an Example Usage section, and no Import section.

#### Scenario: Framework registration

- **WHEN** the provider is built
- **THEN** `tencentcloud_kms_instance` SHALL map to `kms.NewKmsInstance()` in the framework `actionFactories`

#### Scenario: Documentation availability

- **WHEN** documentation is generated via `make doc`
- **THEN** the `tencentcloud_kms_instance` action page SHALL exist under `website/docs/actions/` and SHALL NOT contain a manually-written `Argument Reference` or `Attribute Reference` section

### Requirement: Unit tests with mocked cloud API

A unit test file `action_tc_kms_instance_test.go` SHALL be created using gomonkey to mock the KMS cloud API client and tftypes raw values to build InvokeRequest, testing business logic only without invoking real cloud APIs.

#### Scenario: Mocked invoke success

- **WHEN** the unit test runs Invoke with mocked `UseKmsClient` and `CreateKmsInstanceWithContext` returning a valid response
- **THEN** the test SHALL verify no error diagnostics are produced, the request fields are filled correctly, and `Metadata` reports the type name `tencentcloud_kms_instance`
