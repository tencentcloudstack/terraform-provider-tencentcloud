## ADDED Requirements

### Requirement: Activate CLS for the account

The resource `tencentcloud_cls_open_service_operation` SHALL activate CLS for the current account via the `OpenClsService` API. As a one-time operation resource, it MUST auto-generate its unique ID. The `OpenClsService` API takes no input parameters.

#### Scenario: Open the CLS service

- **WHEN** the user applies the resource
- **THEN** the resource calls `OpenClsService` wrapped in a retry mechanism
- **AND** sets an auto-generated unique resource ID
- **AND** then reads the current service status

### Requirement: Expose CLS service status

The resource SHALL read the current CLS activation status via the `GetClsService` API and expose it as a computed attribute.

#### Scenario: Read service status

- **WHEN** the resource is read
- **THEN** the resource calls `GetClsService` wrapped in a retry mechanism
- **AND** populates the computed `status` attribute from the response in a nil-safe manner (0: opened, 1: not opened)

### Requirement: Delete is a no-op

Because account-level CLS activation cannot be reverted through the API, the delete lifecycle SHALL only remove the resource from Terraform state without calling any API.

#### Scenario: Destroy the operation resource

- **WHEN** the resource is destroyed
- **THEN** no API is called and the resource is removed from state
