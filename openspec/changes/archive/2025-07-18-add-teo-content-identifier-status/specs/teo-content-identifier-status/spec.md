## ADDED Requirements

### Requirement: Status computed field in teo_content_identifier resource
The `tencentcloud_teo_content_identifier` resource SHALL include a computed `status` field of type `TypeString` that exposes the content identifier's lifecycle status from the TEO cloud API. The field SHALL be populated in the Read method from the `Status` property of the `ContentIdentifier` response struct returned by `DescribeContentIdentifiers`. The field SHALL NOT be user-configurable (not included in Create or Modify API requests).

#### Scenario: Status field is populated after resource creation
- **WHEN** a `tencentcloud_teo_content_identifier` resource is created and the Read method is called
- **THEN** the `status` field SHALL be set to the value returned by `DescribeContentIdentifiers` API's `ContentIdentifier.Status` property (e.g., "active")

#### Scenario: Status field is populated after resource read
- **WHEN** the Read method for `tencentcloud_teo_content_identifier` is called and the API returns `Status` as "active"
- **THEN** the `status` attribute in Terraform state SHALL be "active"

#### Scenario: Status field handles nil response value
- **WHEN** the Read method for `tencentcloud_teo_content_identifier` is called and the API returns `Status` as nil
- **THEN** the `status` attribute SHALL NOT be set (skipped via nil-check), preserving the previous state value

#### Scenario: Status field is populated after resource update
- **WHEN** a `tencentcloud_teo_content_identifier` resource is updated and the Read method is called
- **THEN** the `status` field SHALL reflect the current value from the API response

#### Scenario: Status field does not affect Create API request
- **WHEN** the Create method for `tencentcloud_teo_content_identifier` is called
- **THEN** the `status` field SHALL NOT be included in the `CreateContentIdentifier` API request

#### Scenario: Status field does not affect Modify API request
- **WHEN** the Update method for `tencentcloud_teo_content_identifier` is called
- **THEN** the `status` field SHALL NOT be included in the `ModifyContentIdentifier` API request
