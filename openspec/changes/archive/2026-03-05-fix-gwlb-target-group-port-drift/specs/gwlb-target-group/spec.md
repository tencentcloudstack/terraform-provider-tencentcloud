# GWLB Target Group Specification Changes

## MODIFIED Requirements

### Requirement: Port Field Configuration
The `port` field SHALL support both user-specified and API-provided default values without causing configuration drift.

**Changes**:
- Field SHALL have `Computed: true` attribute to accept API defaults
- Field SHALL remain `Optional: true` to allow user specification
- When user does not specify `port`, API-provided default SHALL be accepted
- When user specifies `port`, the specified value SHALL be used

#### Scenario: Port not specified by user
- **GIVEN** a user creates a target group resource
- **WHEN** the user does not specify the `port` field
- **THEN** the API returns a default port value
- **AND** Terraform SHALL accept the API default without detecting drift
- **AND** subsequent `terraform plan` SHALL show no changes

#### Scenario: Port explicitly specified by user
- **GIVEN** a user creates a target group resource
- **WHEN** the user explicitly specifies `port = 6081`
- **THEN** Terraform SHALL use the specified value
- **AND** the API SHALL store and return the specified value
- **AND** subsequent `terraform plan` SHALL show no changes

#### Scenario: Existing resources without drift
- **GIVEN** an existing target group resource created without `port` field
- **WHEN** the user runs `terraform plan` after the fix
- **THEN** no drift SHALL be detected
- **AND** the resource SHALL show as up-to-date
