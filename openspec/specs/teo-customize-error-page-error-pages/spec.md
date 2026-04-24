# teo-customize-error-page-error-pages Specification

## Purpose
TBD - created by archiving change add-teo-customize-error-page-error-pages. Update Purpose after archive.
## Requirements
### Requirement: error_pages computed attribute
The `tencentcloud_teo_customize_error_page` resource SHALL expose an `error_pages` computed attribute of type list, where each element contains the following fields from the `DescribeCustomErrorPages` API response:
- `page_id` (string): Custom error page ID
- `name` (string): Custom error page name
- `content_type` (string): Custom error page type
- `description` (string): Custom error page description
- `content` (string): Custom error page content
- `references` (list of string): List of business IDs that reference this error page

The `error_pages` attribute SHALL be computed-only and SHALL NOT be user-configurable.

#### Scenario: Read resource with error_pages populated
- **WHEN** a `tencentcloud_teo_customize_error_page` resource is read and the `DescribeCustomErrorPages` API returns error pages data with `References` field populated
- **THEN** the `error_pages` attribute SHALL be set with the flattened list of error page objects, each including the `references` sub-field containing the list of `BusinessId` values from the `References` array

#### Scenario: Read resource with empty references
- **WHEN** a `tencentcloud_teo_customize_error_page` resource is read and the `DescribeCustomErrorPages` API returns an error page with no references
- **THEN** the `references` sub-field within the `error_pages` element SHALL be an empty list

#### Scenario: Backward compatibility with existing state
- **WHEN** an existing `tencentcloud_teo_customize_error_page` resource is read without the `error_pages` attribute in its state
- **THEN** the resource SHALL continue to function correctly and the `error_pages` attribute SHALL be populated on the next read operation
