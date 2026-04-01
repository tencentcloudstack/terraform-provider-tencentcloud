## ADDED Requirements

### Requirement: Enable access policy control for VPN SSL Server
The `tencentcloud_vpn_ssl_server` resource SHALL support enabling access policy control by providing the `access_policy_enabled` parameter.

#### Scenario: Enable access policy control
- **WHEN** user sets `access_policy_enabled = true` in Terraform configuration
- **THEN** the VPN SSL Server MUST be created with access policy control enabled
- **AND** the server MUST enforce access policies (configured separately)

#### Scenario: Create server without access policy control (default)
- **WHEN** user does not specify `access_policy_enabled` or sets it to `false`
- **THEN** the VPN SSL Server MUST be created with access policy control disabled (default behavior)
- **AND** no access restrictions MUST be enforced beyond standard VPN authentication

#### Scenario: Update access policy control on existing server
- **WHEN** user changes `access_policy_enabled` from `false` to `true`
- **THEN** the VPN SSL Server MUST be updated to enable access policy control
- **AND** no new resource MUST be created (in-place update)

#### Scenario: Disable access policy control
- **WHEN** user changes `access_policy_enabled` from `true` to `false`
- **THEN** the VPN SSL Server MUST be updated to disable access policy control
- **AND** previously configured policies MUST no longer be enforced

### Requirement: Read access policy control state
The resource SHALL read and store access policy control state when refreshing resource.

#### Scenario: Read policy-enabled server state
- **WHEN** Terraform reads an existing VPN SSL Server with access policy enabled
- **THEN** the state MUST reflect `access_policy_enabled = true`

#### Scenario: Read non-policy server state
- **WHEN** Terraform reads an existing VPN SSL Server without access policy
- **THEN** the state MUST reflect `access_policy_enabled = false` or computed default

### Requirement: Document access policy configuration scope
The resource documentation SHALL clarify that `access_policy_enabled` only controls the feature switch, and detailed access policies must be configured through other means (console or other resources).

#### Scenario: User expects policy details in this resource
- **WHEN** documentation is consulted for access policy configuration
- **THEN** it MUST clearly state that this parameter only enables/disables the feature
- **AND** it MUST provide guidance on where to configure the actual access policies
