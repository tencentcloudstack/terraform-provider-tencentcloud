## ADDED Requirements

### Requirement: Configure SSO authentication for VPN SSL Server
The `tencentcloud_vpn_ssl_server` resource SHALL support enabling SSO (Single Sign-On) authentication by providing `sso_enabled` and `saml_data` parameters.

#### Scenario: Enable SSO authentication with SAML data
- **WHEN** user sets `sso_enabled = true` and provides valid `saml_data` in Terraform configuration
- **THEN** the VPN SSL Server MUST be created with SSO authentication enabled
- **AND** the SAML data MUST be correctly configured in the server

#### Scenario: Create server without SSO (default behavior)
- **WHEN** user does not specify `sso_enabled` or sets it to `false`
- **THEN** the VPN SSL Server MUST be created with SSO authentication disabled (default behavior)
- **AND** the `saml_data` parameter MUST be ignored if provided

#### Scenario: Update SSO configuration on existing server
- **WHEN** user changes `sso_enabled` from `false` to `true` and adds `saml_data`
- **THEN** the VPN SSL Server MUST be updated to enable SSO authentication
- **AND** no new resource MUST be created (in-place update)

#### Scenario: Disable SSO on existing server
- **WHEN** user changes `sso_enabled` from `true` to `false`
- **THEN** the VPN SSL Server MUST be updated to disable SSO authentication
- **AND** the `saml_data` MUST be cleared from the server configuration

### Requirement: Validate SSO configuration prerequisites
The resource SHALL validate that `saml_data` is provided when `sso_enabled` is `true`.

#### Scenario: SSO enabled without SAML data
- **WHEN** user sets `sso_enabled = true` but does not provide `saml_data`
- **THEN** the Terraform plan SHOULD show a warning (validation handled by cloud API)
- **AND** the cloud API MAY reject the request with an appropriate error message

### Requirement: Handle SSO whitelist requirement
The resource documentation SHALL clearly indicate that SSO feature requires whitelist approval from TencentCloud.

#### Scenario: User attempts SSO without whitelist
- **WHEN** user enables SSO but account is not whitelisted
- **THEN** the cloud API MUST return an error indicating whitelist requirement
- **AND** the error message MUST be propagated to the Terraform output

### Requirement: Read SSO configuration state
The resource SHALL read and store SSO configuration state when refreshing resource.

#### Scenario: Read SSO-enabled server state
- **WHEN** Terraform reads an existing VPN SSL Server with SSO enabled
- **THEN** the state MUST reflect `sso_enabled = true`
- **AND** the state SHOULD reflect the SAML configuration (if API returns it)

#### Scenario: Read non-SSO server state
- **WHEN** Terraform reads an existing VPN SSL Server without SSO
- **THEN** the state MUST reflect `sso_enabled = false` or computed default
- **AND** the `saml_data` MUST be empty or absent from state
