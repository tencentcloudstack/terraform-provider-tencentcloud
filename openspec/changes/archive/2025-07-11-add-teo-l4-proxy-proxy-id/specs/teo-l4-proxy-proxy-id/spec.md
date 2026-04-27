## ADDED Requirements

### Requirement: proxy_id computed attribute
The `tencentcloud_teo_l4_proxy` resource SHALL include a `proxy_id` computed string attribute that represents the L4 proxy instance ID returned by the cloud API.

#### Scenario: proxy_id is set after resource creation
- **WHEN** a `tencentcloud_teo_l4_proxy` resource is created successfully via the `CreateL4Proxy` API
- **THEN** the `proxy_id` attribute SHALL be set to the value of `response.Response.ProxyId` from the API response

#### Scenario: proxy_id is set after resource read
- **WHEN** a `tencentcloud_teo_l4_proxy` resource is read via the `DescribeL4Proxy` API
- **THEN** the `proxy_id` attribute SHALL be set to the value of `respData.ProxyId` from the `L4Proxy` struct in the API response

#### Scenario: proxy_id is a computed attribute
- **WHEN** a user defines a `tencentcloud_teo_l4_proxy` resource in their Terraform configuration
- **THEN** the `proxy_id` attribute SHALL NOT be user-settable and SHALL be automatically computed by the provider

#### Scenario: proxy_id is available for reference
- **WHEN** a `tencentcloud_teo_l4_proxy` resource has been created or read
- **THEN** the `proxy_id` attribute SHALL be available for reference in other resources, data sources, or outputs
