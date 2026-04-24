## ADDED Requirements

### Requirement: references computed attribute
The `tencentcloud_teo_customize_error_page` resource SHALL expose a `references` computed attribute at the top level of type list of string, containing the `BusinessId` values from the `ErrorPageReference` objects in the `DescribeCustomErrorPages` API response.

The `references` attribute SHALL be computed-only and SHALL NOT be user-configurable.

#### Scenario: Read resource with references populated
- **WHEN** a `tencentcloud_teo_customize_error_page` resource is read and the `DescribeCustomErrorPages` API returns a `CustomErrorPage` with `References` field populated
- **THEN** the `references` attribute SHALL be set with the list of `BusinessId` values from the `References` array

#### Scenario: Read resource with empty references
- **WHEN** a `tencentcloud_teo_customize_error_page` resource is read and the `DescribeCustomErrorPages` API returns a `CustomErrorPage` with no references
- **THEN** the `references` attribute SHALL be an empty list

#### Scenario: Backward compatibility with existing state
- **WHEN** an existing `tencentcloud_teo_customize_error_page` resource is read without the `references` attribute in its state
- **THEN** the resource SHALL continue to function correctly and the `references` attribute SHALL be populated on the next read operation
