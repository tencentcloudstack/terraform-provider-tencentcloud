## ADDED Requirements

### Requirement: proxy_id computed attribute on tencentcloud_teo_l4_proxy resource
The `tencentcloud_teo_l4_proxy` resource SHALL expose a computed attribute `proxy_id` of type `schema.TypeString` with `Computed: true`. This attribute SHALL NOT be user-configurable. The value SHALL be the L4 proxy instance ID returned by the cloud API.

#### Scenario: proxy_id is set after resource creation
- **WHEN** a `tencentcloud_teo_l4_proxy` resource is created via CreateL4Proxy API
- **THEN** the `proxy_id` attribute SHALL be set directly from `response.Response.ProxyId` in the create function, and subsequently refreshed by the read function via the DescribeL4Proxy API

#### Scenario: proxy_id is populated on resource read
- **WHEN** the read function is called for an existing `tencentcloud_teo_l4_proxy` resource
- **THEN** the `proxy_id` attribute SHALL be set from `respData.ProxyId` returned by the DescribeL4Proxy API

#### Scenario: proxy_id is available after import
- **WHEN** a `tencentcloud_teo_l4_proxy` resource is imported using its composite ID
- **THEN** the `proxy_id` attribute SHALL be populated from the DescribeL4Proxy API response

#### Scenario: proxy_id is not user-configurable
- **WHEN** a user attempts to set `proxy_id` in their Terraform configuration
- **THEN** Terraform SHALL reject the configuration since `proxy_id` is a computed-only attribute

### Requirement: ddos_protection_config deprecated
The `ddos_protection_config` attribute and all its sub-fields (`level_mainland`, `max_bandwidth_mainland`, `level_overseas`) on the `tencentcloud_teo_l4_proxy` resource SHALL be marked as deprecated from version 1.82.90.

#### Scenario: ddos_protection_config shows deprecation warning
- **WHEN** a user configures `ddos_protection_config` in their Terraform configuration
- **THEN** Terraform SHALL display a deprecation warning indicating the field has been deprecated from version 1.82.90

#### Scenario: ddos_protection_config remains functional
- **WHEN** an existing configuration uses `ddos_protection_config`
- **THEN** the field SHALL continue to function for backward compatibility, with the value read from the API response in the read function
