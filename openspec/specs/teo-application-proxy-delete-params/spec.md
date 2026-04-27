### Requirement: Delete function uses d.Get() for zone_id and proxy_id
The `resourceTencentCloudTeoApplicationProxyDelete` function SHALL obtain `zone_id` and `proxy_id` from `d.Get("zone_id")` and `d.Get("proxy_id")` respectively, instead of parsing them from `d.Id()` by splitting on `FILED_SP`.

#### Scenario: Successful deletion using d.Get() values
- **WHEN** the delete function is called and `d.Get("zone_id")` returns a valid zone ID and `d.Get("proxy_id")` returns a valid proxy ID
- **THEN** the function SHALL use these values as input parameters for the `DeleteApplicationProxy` API request

#### Scenario: zone_id is empty from d.Get()
- **WHEN** the delete function is called and `d.Get("zone_id")` returns an empty string
- **THEN** the function SHALL return an error indicating zone_id is required for deletion

#### Scenario: proxy_id is empty from d.Get()
- **WHEN** the delete function is called and `d.Get("proxy_id")` returns an empty string
- **THEN** the function SHALL return an error indicating proxy_id is required for deletion

### Requirement: Delete function calls DeleteApplicationProxy API directly
The `resourceTencentCloudTeoApplicationProxyDelete` function SHALL construct a `DeleteApplicationProxyRequest` with `ZoneId` and `ProxyId` fields and call the `DeleteApplicationProxy` API directly with retry logic using `tccommon.WriteRetryTimeout`, instead of delegating to `service.DeleteTeoApplicationProxyById()`.

#### Scenario: Successful API call with retry
- **WHEN** the `DeleteApplicationProxy` API is called with valid `ZoneId` and `ProxyId`
- **THEN** the request SHALL include `ZoneId` and `ProxyId` as parameters and use `resource.Retry` with `tccommon.WriteRetryTimeout`

#### Scenario: API call fails with retryable error
- **WHEN** the `DeleteApplicationProxy` API returns a retryable error
- **THEN** the function SHALL retry the request using `tccommon.RetryError()` to wrap the error

### Requirement: Delete function preserves offline-before-delete flow
The `resourceTencentCloudTeoApplicationProxyDelete` function SHALL preserve the existing two-step deletion flow: first setting the proxy status to `offline` via `ModifyApplicationProxyStatus`, then calling `DeleteApplicationProxy`.

#### Scenario: Proxy is online before deletion
- **WHEN** the proxy status is `online`
- **THEN** the function SHALL first call `ModifyApplicationProxyStatus` to set status to `offline`, wait for the status change, then call `DeleteApplicationProxy`

#### Scenario: Proxy is already offline
- **WHEN** the proxy status is `offline`
- **THEN** the function SHALL skip the offline step and directly call `DeleteApplicationProxy`

### Requirement: Unit tests for delete function
The `resource_tc_teo_application_proxy_test.go` file SHALL include unit tests that verify the delete function correctly uses `d.Get()` for `zone_id` and `proxy_id` and properly constructs the `DeleteApplicationProxyRequest`.

#### Scenario: Unit test verifies d.Get() usage in delete
- **WHEN** unit tests are run for the delete function
- **THEN** tests SHALL verify that `zone_id` and `proxy_id` are obtained from `d.Get()` and passed as `ZoneId` and `ProxyId` in the API request
