## ADDED Requirements

### Requirement: proxy_id computed attribute on tencentcloud_teo_l4_proxy resource
The `tencentcloud_teo_l4_proxy` resource SHALL expose a computed attribute `proxy_id` of type `schema.TypeString` with `Computed: true`. This attribute SHALL NOT be user-configurable. The value SHALL be the L4 proxy instance ID returned by the cloud API.

#### Scenario: proxy_id is set after resource creation
- **WHEN** a `tencentcloud_teo_l4_proxy` resource is created via CreateL4Proxy API
- **THEN** the `proxy_id` attribute SHALL be set to the `ProxyId` value returned in the CreateL4Proxy response, as read back by the DescribeL4Proxy API

#### Scenario: proxy_id is populated on resource read
- **WHEN** the read function is called for an existing `tencentcloud_teo_l4_proxy` resource
- **THEN** the `proxy_id` attribute SHALL be set from `respData.ProxyId` returned by the DescribeL4Proxy API

#### Scenario: proxy_id is available after import
- **WHEN** a `tencentcloud_teo_l4_proxy` resource is imported using its composite ID
- **THEN** the `proxy_id` attribute SHALL be populated from the DescribeL4Proxy API response

#### Scenario: proxy_id is not user-configurable
- **WHEN** a user attempts to set `proxy_id` in their Terraform configuration
- **THEN** Terraform SHALL reject the configuration since `proxy_id` is a computed-only attribute
