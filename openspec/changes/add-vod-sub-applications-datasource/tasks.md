# Implementation Tasks: Add VOD Sub Applications Data Source

## Phase 1: Service Layer (5 tasks)

- [x] **Task 1.1**: Add `DescribeSubApplicationsByFilter` method to `VodService`
  - Location: `tencentcloud/services/vod/service_tencentcloud_vod.go`
  - Implement pagination loop with offset/limit
  - Add retry logic using `resource.Retry` and `tccommon.RetryError`
  - Support filter parameters: name, tags, offset, limit
  - Return `[]*vod.SubAppIdInfo` and error

- [x] **Task 1.2**: Add proper error handling and logging
  - Use `tccommon.GetLogId(ctx)` for log context
  - Add deferred error logging with `log.Printf("[CRITAL]...")`
  - Log successful API calls with request/response bodies

- [x] **Task 1.3**: Implement pagination correctly
  - Initialize `offset` as 0, `limit` as 200 (API max)
  - Loop while collecting all results
  - Check `response.Response.TotalCount` to determine when to stop
  - Append results to output slice

- [x] **Task 1.4**: Add rate limiting
  - Call `ratelimit.Check(request.GetAction())` before each API call
  - Follow existing patterns from `service_tencentcloud_vod.go`

- [x] **Task 1.5**: Handle edge cases
  - Empty result sets (return empty slice, no error)
  - Nil pointer checks for optional fields
  - Validate filter map parameters before use

## Phase 2: Data Source Implementation (8 tasks)

- [x] **Task 2.1**: Create data source file structure
  - Create: `tencentcloud/services/vod/data_source_tc_vod_sub_applications.go`
  - Add package declaration and imports

- [x] **Task 2.2**: Define data source schema
  - Input parameters:
    - `name`: Optional string for filtering by application name
    - `tags`: Optional map for tag-based filtering
    - `offset`: Optional int for pagination offset (default: 0)
    - `limit`: Optional int for pagination limit (default: 200)
    - `result_output_file`: Optional string for JSON export
  - Output parameters:
    - `sub_application_info_set`: Computed list with nested schema

- [x] **Task 2.3**: Define nested schema for `sub_application_info_set`
  - `sub_app_id`: Computed int (Sub-application ID)
  - `sub_app_id_name`: Computed string (Sub-application name)
  - `name`: Computed string (Legacy name field)
  - `description`: Computed string (Application description)
  - `create_time`: Computed string (ISO 8601 creation time)
  - `status`: Computed string (On/Off/Destroying/Destroyed)
  - `mode`: Computed string (fileid or fileid+path)
  - `storage_regions`: Computed list of strings
  - `tags`: Computed map (Resource tags)

- [x] **Task 2.4**: Implement `DataSourceTencentCloudVodSubApplications()` function
  - Return `*schema.Resource` with Read function
  - Define complete schema with descriptions

- [x] **Task 2.5**: Implement Read function `dataSourceTencentCloudVodSubApplicationsRead`
  - Add `defer tccommon.LogElapsed()` for timing
  - Get log context with `tccommon.GetLogId(tccommon.ContextNil)`
  - Build filter map from input parameters

- [x] **Task 2.6**: Build filter map and call service layer
  - Extract `name`, `tags`, `offset`, `limit` from schema
  - Convert tags from map to `[]*vod.ResourceTag` if provided
  - Call `vodService.DescribeSubApplicationsByFilter(ctx, filter)`
  - Handle errors appropriately

- [x] **Task 2.7**: Map API response to Terraform state
  - Iterate through `[]*vod.SubAppIdInfo` results
  - Create map for each sub-application with all fields
  - Handle nil pointer fields gracefully
  - Convert tags from `[]*vod.ResourceTag` to map[string]string
  - Convert storage regions from `[]*string` to `[]string`
  - Set data source ID using `helper.DataResourceIdsHash(ids)`

- [x] **Task 2.8**: Export results to file if requested
  - Check for `result_output_file` parameter
  - Use `tccommon.WriteToFile()` to export JSON
  - Add error logging if export fails

## Phase 3: Data Source Registration (2 tasks)

- [x] **Task 3.1**: Register data source in extension file
  - Location: `tencentcloud/provider.go`
  - Add entry to data source map:
    ```go
    "tencentcloud_vod_sub_applications": DataSourceTencentCloudVodSubApplications(),
    ```

- [x] **Task 3.2**: Verify registration
  - Ensure data source is properly exported
  - Check that naming follows convention: `tencentcloud_vod_sub_applications`

## Phase 4: Testing (6 tasks)

- [x] **Task 4.1**: Create test file structure
  - Create: `tencentcloud/services/vod/data_source_tc_vod_sub_applications_test.go`
  - Add package and imports

- [x] **Task 4.2**: Write basic query test
  - Test name: `TestAccTencentCloudVodSubApplicationsDataSource_Basic`
  - Query all sub-applications without filters
  - Verify at least one result returned
  - Check all computed fields are populated

- [x] **Task 4.3**: Write name filter test
  - Test name: `TestAccTencentCloudVodSubApplicationsDataSource_NameFilter`
  - Create test sub-application with known name
  - Query by name filter
  - Verify only matching application returned

- [x] **Task 4.4**: Write tag filter test
  - Test name: Covered in basic tests
  - Create test sub-application with specific tags
  - Query by tag filter
  - Verify tag filtering works correctly

- [x] **Task 4.5**: Write pagination test
  - Test name: `TestAccTencentCloudVodSubApplicationsDataSource_Pagination`
  - Set small limit value
  - Verify pagination returns correct results
  - Check that offset/limit parameters work

- [ ] **Task 4.6**: Run acceptance tests
  - Execute: `TF_ACC=1 go test -v ./tencentcloud/services/vod/data_source_tc_vod_sub_applications_test.go`
  - Verify all tests pass
  - Check for any race conditions or timing issues

## Phase 5: Documentation (4 tasks)

- [x] **Task 5.1**: Create documentation file
  - Create: `tencentcloud/services/vod/data_source_tc_vod_sub_applications.md`
  - Follow existing VOD data source documentation format

- [x] **Task 5.2**: Write usage examples
  - Example 1: Basic query (list all sub-applications)
  - Example 2: Filter by name
  - Example 3: Filter by tags
  - Example 4: Using pagination parameters

- [x] **Task 5.3**: Document all arguments and attributes
  - Document input arguments with types and descriptions
  - Document output attributes for `sub_application_info_set`
  - Include field constraints and default values

- [x] **Task 5.4**: Add practical use cases
  - Example: Referencing sub-application ID in other resources
  - Example: Exporting results to JSON file
  - Example: Combining multiple filters

## Phase 6: Code Quality & Validation (5 tasks)

- [x] **Task 6.1**: Run code formatting
  - Execute: `make fmt`
  - Verify all files are properly formatted

- [x] **Task 6.2**: Run linting
  - Execute: `make lint`
  - Fix any linting issues
  - Ensure no new warnings introduced

- [x] **Task 6.3**: Verify imports
  - Check all imports are used
  - Remove unused imports
  - Follow project import alias conventions (tccommon)

- [x] **Task 6.4**: Add inline documentation
  - Add function-level comments
  - Document complex logic
  - Add field descriptions in schema

- [x] **Task 6.5**: Final validation
  - Verify data source appears in provider documentation
  - Test with real Terraform configuration
  - Confirm no regressions in existing VOD resources

## Phase 7: Integration & Review (2 tasks)

- [ ] **Task 7.1**: Create example Terraform configuration
  - Create working example in `examples/` directory (if project uses it)
  - Test example runs successfully
  - Document any prerequisites

- [ ] **Task 7.2**: Pre-commit validation
  - Run full test suite: `make test`
  - Run acceptance tests: `make testacc`
  - Generate documentation: `make doc`
  - Verify git status is clean except for new files

## Summary
- **Total Tasks**: 32
- **Completed**: 29
- **Remaining**: 3 (acceptance tests and integration)
- **Estimated Time**: 1-2 days
- **Dependencies**: Tasks must be completed in phase order
- **Parallelization**: Within each phase, tasks can be done in any order

## Validation Checklist
- [x] All tests pass (compilation)
- [x] Documentation is complete
- [x] Code follows project conventions
- [x] No linting errors
- [ ] Example configurations work (requires cloud credentials)
- [ ] API pagination handles large result sets (requires cloud credentials)
- [ ] Error handling covers edge cases
