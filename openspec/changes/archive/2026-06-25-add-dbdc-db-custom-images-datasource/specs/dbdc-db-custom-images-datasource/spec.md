## ADDED Requirements

### Requirement: Data source schema defines image_set with flattened fields
The data source `tencentcloud_dbdc_db_custom_images` SHALL define a `image_set` computed field of type `TypeList` with `Elem: &schema.Resource{}` containing individual computed fields: `image_id` (TypeString), `os_name` (TypeString), `image_type` (TypeString), `architecture` (TypeString). The schema SHALL NOT create a wrapper nesting layer around the image list.

#### Scenario: Schema structure is correctly flattened
- **WHEN** the data source schema is inspected
- **THEN** `image_set` is TypeList with Elem containing a Resource with `image_id`, `os_name`, `image_type`, `architecture` as individual Computed fields at the same level

### Requirement: Data source has no user-facing filter parameters
The data source SHALL only expose `result_output_file` as an Optional input parameter. Offset and Limit SHALL NOT be exposed to users, as pagination SHALL be handled internally by the service layer.

#### Scenario: No pagination parameters in schema
- **WHEN** the data source schema is inspected
- **THEN** there is no `offset` or `limit` field in the schema, only `result_output_file` and `image_set`

### Requirement: Read function queries DescribeDBCustomImages with retry
The Read function SHALL call the service layer's `DescribeDBCustomImagesByFilter` method wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. Errors from the API SHALL be wrapped with `tccommon.RetryError(e)`.

#### Scenario: Successful API call returns image list
- **WHEN** the Read function is invoked and the DescribeDBCustomImages API succeeds
- **THEN** the `image_set` field is populated with all available images, and `d.SetId(helper.BuildToken())` is called

#### Scenario: API call fails transiently
- **WHEN** the DescribeDBCustomImages API returns a retryable error
- **THEN** the retry mechanism continues attempting until success or timeout

### Requirement: Service layer handles automatic pagination
The service layer method `DescribeDBCustomImagesByFilter` SHALL implement pagination with Limit=100, accumulating all results across pages. Each page request SHALL be wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. The loop SHALL break when the number of returned items is less than the limit.

#### Scenario: Multiple pages of results
- **WHEN** the API returns more than 100 images
- **THEN** the service layer makes multiple paginated requests and accumulates all images into a single result list

#### Scenario: Single page of results
- **WHEN** the API returns fewer than 100 images
- **THEN** the service layer makes one request and returns the complete list

### Requirement: Empty response returns NonRetryableError
In the service layer, if `DescribeDBCustomImages` returns `nil` response, `nil` Response, or empty `ImageSet`, the method SHALL return `resource.NonRetryableError` with a descriptive error message, rather than silently clearing the data source ID. A `log.Printf("[DATASOURCE] read empty, skip SetId")` SHALL be logged on the retry failure path.

#### Scenario: API returns nil response
- **WHEN** the DescribeDBCustomImages API returns nil or empty response
- **THEN** NonRetryableError is returned, and the data source ID is not cleared

### Requirement: Nil-check before setting each field
In the Read function, each field of DBCustomImage SHALL be checked for nil before being set into the Terraform state. Fields with nil values SHALL NOT be set.

#### Scenario: Image with nil fields
- **WHEN** an image in the API response has nil OsName
- **THEN** the `os_name` field is not set for that image entry, but other non-nil fields are set

### Requirement: Data source is registered in provider.go and provider.md
The data source SHALL be registered in `tencentcloud/provider.go` DataSourcesMap as `"tencentcloud_dbdc_db_custom_images": dbdc.DataSourceTencentCloudDbdcDbCustomImages()`, and listed in `tencentcloud/provider.md`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_dbdc_db_custom_images` is available as a data source

### Requirement: Unit tests use gomonkey mock pattern
The test file `data_source_tc_dbdc_db_custom_images_test.go` SHALL use gomonkey to mock the `DescribeDBCustomImages` client method. Tests SHALL verify schema structure, Read function behavior with mock responses, nil-field handling, and empty-response error handling. Tests SHALL be runnable with `go test -gcflags="all=-l"`.

#### Scenario: Basic read test with mock data
- **WHEN** the Read function is called with a mocked DescribeDBCustomImages returning two images
- **THEN** `image_set` contains two entries with correct field values, and `d.Id()` is not empty

#### Scenario: Empty response test
- **WHEN** the Read function is called with a mocked DescribeDBCustomImages returning empty ImageSet
- **THEN** an error is returned (NonRetryableError propagated)

### Requirement: .md documentation file follows gendoc format
The `data_source_tc_dbdc_db_custom_images.md` file SHALL contain a one-line description mentioning the DBDC product name, Example Usage section with HCL examples, and no Argument Reference or Attribute Reference sections (those are auto-generated).

#### Scenario: Documentation format
- **WHEN** the .md file is generated
- **THEN** it starts with "Use this data source to query ..." mentioning DBDC, contains Example Usage, and does not contain Argument Reference or Attribute Reference sections
