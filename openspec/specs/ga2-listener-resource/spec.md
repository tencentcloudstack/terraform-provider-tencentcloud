## ADDED Requirements

### Requirement: Create GA2 Listener
The system SHALL provide a Terraform resource `tencentcloud_ga2_listener` that creates a GA2 listener by calling the CreateListener API. The resource SHALL support the following input parameters: `global_accelerator_id` (Required, ForceNew), `name` (Required), `port_ranges` (Required, ForceNew, nested with `from_port` and `to_port`), `description` (Optional), `listener_type` (Optional, ForceNew), `protocol` (Optional, ForceNew), `idle_timeout` (Optional), `get_real_ip_type` (Optional), `client_affinity` (Optional), `request_timeout` (Optional), `x_forwarded_for_real_ip` (Optional), `certification_type` (Optional), `cipher_policy_id` (Optional), `server_certificates` (Optional), `client_ca_certificates` (Optional). After creation, the system SHALL poll DescribeTaskResult until the task succeeds, then set the resource ID as `global_accelerator_id#listener_id`.

#### Scenario: Successful listener creation
- **WHEN** user applies a Terraform configuration with valid `tencentcloud_ga2_listener` resource
- **THEN** the system calls CreateListener API, polls DescribeTaskResult until SUCCESS, and stores the composite ID `global_accelerator_id#listener_id`

#### Scenario: CreateListener returns nil response
- **WHEN** CreateListener API returns nil Response or nil ListenerId
- **THEN** the system returns a NonRetryableError

### Requirement: Read GA2 Listener
The system SHALL read the listener state by calling DescribeListeners API with a filter on `listener-id`. The system SHALL set all readable attributes from the ListenerSet response including: `global_accelerator_id`, `listener_id`, `name`, `description`, `protocol`, `port_ranges`, `x_forwarded_for_real_ip`, `client_affinity`, `client_affinity_time`, `certification_type`, `server_certificates`, `client_ca_certificates`, `cipher_policy_id`, `request_timeout`, `listener_type`, `get_real_ip_type`, `idle_timeout`. Each field SHALL only be set if the corresponding response field is non-nil.

#### Scenario: Listener exists
- **WHEN** DescribeListeners returns a matching listener in ListenerSet
- **THEN** the system sets all non-nil fields from the response into Terraform state

#### Scenario: Listener not found
- **WHEN** DescribeListeners returns no matching listener
- **THEN** the system removes the resource from state by calling `d.SetId("")`

### Requirement: Update GA2 Listener
The system SHALL update the listener by calling ModifyListener API when any of the following fields change: `name`, `description`, `idle_timeout`, `client_affinity`, `client_affinity_time`, `request_timeout`, `x_forwarded_for_real_ip`, `certification_type`, `cipher_policy_id`, `server_certificates`, `client_ca_certificates`, `get_real_ip_type`. After modification, the system SHALL poll DescribeTaskResult until the task succeeds.

#### Scenario: Successful listener update
- **WHEN** user modifies an updatable field and applies
- **THEN** the system calls ModifyListener API with changed fields, polls DescribeTaskResult until SUCCESS, then reads the updated state

#### Scenario: No changes detected
- **WHEN** no updatable fields have changed
- **THEN** the system skips the ModifyListener call and returns the current state

### Requirement: Delete GA2 Listener
The system SHALL delete the listener by calling DeleteListener API with `global_accelerator_id` and `listener_id`. After deletion, the system SHALL poll DescribeTaskResult until the task succeeds.

#### Scenario: Successful listener deletion
- **WHEN** user destroys the `tencentcloud_ga2_listener` resource
- **THEN** the system calls DeleteListener API, polls DescribeTaskResult until SUCCESS

#### Scenario: DeleteListener returns nil TaskId
- **WHEN** DeleteListener API returns nil Response or nil TaskId
- **THEN** the system returns a NonRetryableError

### Requirement: Import GA2 Listener
The system SHALL support importing an existing listener using the composite ID format `global_accelerator_id#listener_id` (separated by tccommon.FILED_SP).

#### Scenario: Import with valid composite ID
- **WHEN** user runs `terraform import tencentcloud_ga2_listener.example ga-xxx#lbl-xxx`
- **THEN** the system parses the composite ID and reads the listener state successfully

#### Scenario: Import with invalid ID format
- **WHEN** user provides an ID that does not contain exactly one separator
- **THEN** the system returns an error indicating the expected format

### Requirement: Register resource in provider
The system SHALL register `tencentcloud_ga2_listener` in `tencentcloud/provider.go` ResourcesMap and add the corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Provider includes ga2_listener
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_ga2_listener` is available as a managed resource

### Requirement: Unit tests with gomonkey mock
The system SHALL provide unit tests in `resource_tc_ga2_listener_test.go` using gomonkey to mock GA2 API calls. Tests SHALL cover Create, Read, Update, and Delete operations and SHALL pass with `go test -gcflags=all=-l`.

#### Scenario: Unit tests pass
- **WHEN** running `go test -gcflags=all=-l` on the test file
- **THEN** all test cases pass successfully

### Requirement: Resource documentation
The system SHALL provide a documentation file at `tencentcloud/services/ga2/resource_tc_ga2_listener.md` with Example Usage and Import sections. The Import section SHALL document the composite ID format `global_accelerator_id#listener_id`.

#### Scenario: Documentation includes example and import
- **WHEN** the documentation is generated
- **THEN** it contains a valid Example Usage with all required fields and an Import section showing the composite ID format
