## ADDED Requirements

### Requirement: AIGC quota resource CRUD lifecycle
The system SHALL provide a Terraform resource `tencentcloud_vod_aigc_quota` that supports full Create, Read, Update, and Delete operations for managing VOD AIGC quotas.

#### Scenario: Create an Image quota
- **WHEN** the user defines a `tencentcloud_vod_aigc_quota` resource with `sub_app_id`, `quota_type = "Image"`, and `quota_limit`
- **THEN** the provider SHALL call `CreateAigcQuota` API and return the resource state with the composite ID `{sub_app_id}#Image#`

#### Scenario: Create a Text quota with API token
- **WHEN** the user defines a `tencentcloud_vod_aigc_quota` resource with `sub_app_id`, `quota_type = "Text"`, `quota_limit`, and `api_token`
- **THEN** the provider SHALL call `CreateAigcQuota` API with all four parameters and return the resource state with the composite ID `{sub_app_id}#Text#{api_token}`

#### Scenario: Read an existing quota
- **WHEN** the provider reads a `tencentcloud_vod_aigc_quota` resource
- **THEN** the provider SHALL call `DescribeAigcQuotas` API with the parsed `sub_app_id`, `quota_type`, and `api_token` from the composite ID, and populate `quota_limit` and `usage` from the response

#### Scenario: Update quota limit
- **WHEN** the user changes `quota_limit` for an existing `tencentcloud_vod_aigc_quota` resource
- **THEN** the provider SHALL call `ModifyAigcQuota` API with the new `quota_limit` value

#### Scenario: Delete a quota
- **WHEN** the user destroys a `tencentcloud_vod_aigc_quota` resource
- **THEN** the provider SHALL call `DeleteAigcQuota` API and remove the resource from state

### Requirement: Composite ID with import support
The system SHALL use a composite ID format `{sub_app_id}#{quota_type}#{api_token}` for uniquely identifying each AIGC quota resource, and SHALL support Terraform import using this format.

#### Scenario: Import an existing Image quota
- **WHEN** the user imports a resource with ID `12345#Image#`
- **THEN** the provider SHALL parse the ID and call `DescribeAigcQuotas` to populate the resource state

#### Scenario: Import an existing Text quota with API token
- **WHEN** the user imports a resource with ID `12345#Text#my_token_value`
- **THEN** the provider SHALL parse the ID and call `DescribeAigcQuotas` to populate the resource state

### Requirement: Schema field definitions
The resource schema SHALL expose the following fields matching the cloud API parameter semantics:

| Field | Type | Mode | Description |
|---|---|---|---|
| `sub_app_id` | TypeInt | Required, ForceNew | VOD sub-application ID |
| `quota_type` | TypeString | Required, ForceNew | Quota type: Image, Video, or Text |
| `quota_limit` | TypeInt | Required | Quota limit value |
| `api_token` | TypeString | Optional, ForceNew | API token (only meaningful for Text quotas) |
| `usage` | TypeInt | Computed | Current usage amount (read-only from cloud) |

#### Scenario: All required fields provided
- **WHEN** the user provides `sub_app_id`, `quota_type`, and `quota_limit`
- **THEN** the provider SHALL accept the configuration as valid

#### Scenario: Text quota with API token provided
- **WHEN** the user provides `quota_type = "Text"` with an `api_token`
- **THEN** the provider SHALL include the API token in all API calls and the composite ID

### Requirement: Async consistency handling
The system SHALL handle eventual consistency after Create operations by polling `DescribeAigcQuotas` until the newly created quota becomes visible.

#### Scenario: Quota not immediately visible after Create
- **WHEN** the Create API returns success but Describe does not yet show the quota
- **THEN** the provider SHALL retry Describe for up to `ReadRetryTimeout` duration until the quota appears

#### Scenario: Quota not visible after Delete
- **WHEN** the Delete API returns success but Describe still shows the quota
- **THEN** the provider SHALL retry Describe for up to `ReadRetryTimeout` duration until the quota disappears from the list
