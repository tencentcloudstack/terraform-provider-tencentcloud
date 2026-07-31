## 1. Service Layer

- [x] 1.1 Add `DescribeDBCustomNodeSecurityGroupsById` method in `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go`
  - Call `DescribeDBCustomNodeSecurityGroups` API with `NodeId` parameter
  - Use `tccommon.ReadRetryTimeout` retry with `ratelimit.Check`
  - Check nil response (`result == nil || result.Response == nil`) and return `NonRetryableError`
  - Log empty response with `[DATASOURCE]` prefix on retry exhaustion

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_security_groups.go`
  - Define schema with `node_id` (Required, TypeString) and `groups` (Computed, TypeList)
  - `groups` element contains: `security_group_id`, `security_group_name`, `security_group_remark`, `project_id`, `create_time`, `inbound`, `outbound`
  - `inbound`/`outbound` are `PolicyRule` lists with: `action`, `cidr_ip`, `port_range`, `ip_protocol`, `service_module`, `address_module`, `id`
  - Implement Read function with `defer LogElapsed` and `defer InconsistentCheck`
  - Parse `node_id` from `d.GetOk`, call service layer method
  - Map API response fields to Terraform state, check nil before each `d.Set`
  - Set `d.SetId(helper.BuildToken())` and handle `result_output_file`

## 3. Data Source Test

- [x] 3.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_security_groups_test.go`
  - Use gomonkey mock for cloud API calls (not terraform test suite)
  - Test normal case: valid `node_id` returns security groups
  - Test empty case: `node_id` returns empty groups
  - Test error case: API returns nil response

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_dbdc_db_custom_node_security_groups` in `tencentcloud/provider.go`
  - Add to data sources map: `"tencentcloud_dbdc_db_custom_node_security_groups": dbdc.DataSourceTencentCloudDbdcDbCustomNodeSecurityGroups()`

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_node_security_groups.md`
  - Format: one-line description with "Use this data source to query..."
  - Example Usage section with HCL code
  - No Argument Reference or Attribute Reference sections (auto-generated)