## ADDED Requirements

### Requirement: SCF Function Qualifier Parameter

The `tencentcloud_scf_function` resource SHALL support an optional `qualifier` parameter that allows users to specify a function version number or alias for API operations that support it.

#### Scenario: User specifies qualifier in configuration

- **WHEN** a user sets `qualifier = "my-alias"` in their `tencentcloud_scf_function` resource configuration
- **THEN** the provider SHALL pass this qualifier value to the `GetFunction` API when reading the function state
- **AND** the provider SHALL pass this qualifier value to the `DeleteFunction` API when deleting the function
- **AND** the provider SHALL pass this qualifier value to the `CreateTrigger` and `DeleteTrigger` APIs when managing triggers

#### Scenario: User does not specify qualifier

- **WHEN** a user does not set the `qualifier` parameter in their configuration
- **THEN** the provider SHALL NOT pass a qualifier to the `GetFunction`, `DeleteFunction`, `CreateTrigger`, or `DeleteTrigger` APIs (allowing the API to use its default behavior, typically `$LATEST`)
- **AND** the provider SHALL read back the qualifier value returned by the `GetFunction` API response and store it in Terraform state

#### Scenario: Provider reads qualifier from API response

- **WHEN** the provider calls the `GetFunction` API
- **THEN** the provider SHALL read the `response.Qualifier` field (top-level function qualifier) from the response and store it in the `qualifier` attribute of the Terraform state
- **AND** the provider SHALL read the `response.Triggers[].Qualifier` field from each trigger in the response

#### Scenario: Backward compatibility with existing configurations

- **WHEN** an existing Terraform configuration that does not include the `qualifier` parameter is applied
- **THEN** the resource SHALL continue to function identically to before the change
- **AND** no state migration or manual intervention SHALL be required