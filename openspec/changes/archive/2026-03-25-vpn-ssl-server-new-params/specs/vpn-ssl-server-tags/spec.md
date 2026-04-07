## ADDED Requirements

### Requirement: Manage tags for VPN SSL Server resources
The `tencentcloud_vpn_ssl_server` resource SHALL support resource tagging by providing the `tags` parameter as a key-value map.

#### Scenario: Create server with tags
- **WHEN** user provides `tags = { Environment = "production", Owner = "team-a" }` in configuration
- **THEN** the VPN SSL Server MUST be created with these tags attached
- **AND** the tags MUST be visible in TencentCloud console and API responses

#### Scenario: Create server without tags (default)
- **WHEN** user does not specify `tags` in configuration
- **THEN** the VPN SSL Server MUST be created without any tags
- **AND** the state MUST not contain tags attribute or show it as empty

#### Scenario: Add tags to existing server
- **WHEN** user adds `tags` parameter to an existing VPN SSL Server configuration
- **THEN** the server MUST be updated with the new tags
- **AND** no new resource MUST be created (in-place update)

#### Scenario: Modify tags on existing server
- **WHEN** user changes tag values or adds/removes tag keys
- **THEN** the server tags MUST be updated to match the new configuration
- **AND** only the changed tags MUST be reflected in the Terraform diff

#### Scenario: Remove all tags from server
- **WHEN** user removes the `tags` parameter from configuration
- **THEN** all tags MUST be removed from the VPN SSL Server
- **AND** the state MUST reflect empty tags

### Requirement: Support standard tag key-value constraints
The resource SHALL accept tag keys and values that conform to TencentCloud tagging standards.

#### Scenario: Valid tag format
- **WHEN** user provides tags with valid keys and values (alphanumeric, hyphens, underscores)
- **THEN** the tags MUST be accepted and applied to the resource

#### Scenario: Invalid tag format
- **WHEN** user provides tags with invalid characters or exceeding length limits
- **THEN** the cloud API MUST reject the request with validation error
- **AND** the error MUST be propagated to Terraform output

### Requirement: Read tags state from cloud API
The resource SHALL read and store tags when refreshing resource state.

#### Scenario: Read server with tags
- **WHEN** Terraform reads an existing VPN SSL Server with tags
- **THEN** the state MUST reflect all current tags as key-value pairs
- **AND** the tags MUST match the cloud resource state

#### Scenario: Handle missing tags in API response
- **WHEN** the DescribeVpnGatewaySslServers API does not return tags field (or returns null)
- **THEN** the state SHOULD preserve the configured tags (Computed behavior)
- **OR** treat it as empty tags if the API explicitly returns empty list

#### Scenario: Detect tag drift
- **WHEN** tags are modified outside Terraform (via console or other tools)
- **THEN** Terraform refresh MUST detect the drift
- **AND** the plan MUST show the differences between desired and actual state

### Requirement: Use Provider's tag helper functions
The resource implementation SHALL use the Provider's common tag helper functions (`helper.GetTags`) for consistency with other resources.

#### Scenario: Consistent tag handling across resources
- **WHEN** implementing tag support
- **THEN** the code MUST use `helper.GetTags(d, "tags")` for reading configuration
- **AND** the code MUST follow the same tag-to-SDK mapping pattern as other TencentCloud resources (e.g., CLS, CVM)
