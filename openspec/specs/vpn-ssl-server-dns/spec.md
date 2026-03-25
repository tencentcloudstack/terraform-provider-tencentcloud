## ADDED Requirements

### Requirement: Configure custom DNS servers for VPN SSL Server
The `tencentcloud_vpn_ssl_server` resource SHALL support configuring custom DNS servers by providing the `dns_servers` nested object with `primary_dns` and `secondary_dns` fields.

#### Scenario: Configure primary DNS only
- **WHEN** user provides `dns_servers { primary_dns = "8.8.8.8" }` in Terraform configuration
- **THEN** the VPN SSL Server MUST be created with primary DNS set to "8.8.8.8"
- **AND** secondary DNS MUST use system default or remain unset

#### Scenario: Configure both primary and secondary DNS
- **WHEN** user provides `dns_servers { primary_dns = "8.8.8.8"; secondary_dns = "8.8.4.4" }` in configuration
- **THEN** the VPN SSL Server MUST be created with primary DNS "8.8.8.8" and secondary DNS "8.8.4.4"
- **AND** VPN clients MUST receive these DNS servers upon connection

#### Scenario: Create server without custom DNS (default)
- **WHEN** user does not specify `dns_servers` in configuration
- **THEN** the VPN SSL Server MUST use cloud platform default DNS servers
- **AND** no custom DNS configuration MUST be applied

#### Scenario: Update DNS configuration on existing server
- **WHEN** user changes `dns_servers` values on an existing VPN SSL Server
- **THEN** the server MUST be updated with new DNS configuration
- **AND** no new resource MUST be created (in-place update)

#### Scenario: Remove custom DNS configuration
- **WHEN** user removes the `dns_servers` block from configuration
- **THEN** the server MUST revert to default DNS configuration
- **AND** VPN clients MUST receive default DNS servers upon reconnection

### Requirement: Validate DNS server IP addresses
The resource SHALL accept valid IP addresses for primary and secondary DNS servers.

#### Scenario: Invalid DNS IP address format
- **WHEN** user provides an invalid IP address (e.g., "999.999.999.999" or "invalid")
- **THEN** the cloud API MUST reject the request with validation error
- **AND** the error MUST be propagated to Terraform output

#### Scenario: Empty DNS server value
- **WHEN** user provides `dns_servers { primary_dns = "" }` (empty string)
- **THEN** the provider SHOULD treat it as unset
- **AND** no DNS configuration MUST be sent to the API

### Requirement: Read DNS configuration state
The resource SHALL read and store DNS configuration when refreshing resource.

#### Scenario: Read server with custom DNS
- **WHEN** Terraform reads an existing VPN SSL Server with custom DNS configured
- **THEN** the state MUST reflect the current `primary_dns` value
- **AND** the state MUST reflect the current `secondary_dns` value (if set)

#### Scenario: Read server with default DNS
- **WHEN** Terraform reads an existing VPN SSL Server without custom DNS
- **THEN** the `dns_servers` MUST be empty or reflect computed default values
- **AND** no unexpected diff MUST be shown in plan

### Requirement: Handle API response variations for DNS
The resource SHALL handle cases where the API may or may not return DNS configuration in describe responses.

#### Scenario: API does not return DNS in describe
- **WHEN** the DescribeVpnGatewaySslServers API omits `DnsServers` field
- **THEN** the state MUST retain the configured values (Computed attribute behavior)
- **AND** Terraform MUST NOT show unnecessary diffs
