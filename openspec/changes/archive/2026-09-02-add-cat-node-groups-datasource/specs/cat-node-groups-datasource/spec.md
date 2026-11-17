## ADDED Requirements

### Requirement: Data source `tencentcloud_cat_node_groups` SHALL exist and expose `DescribeNodeGroups` filters
The provider SHALL register a data source named `tencentcloud_cat_node_groups` whose Read operation calls the cloud API `DescribeNodeGroups` (package `cat.v20180409`). The data source SHALL expose the following optional filter arguments, each mapped to the corresponding cloud API request field: `node_type` (`NodeType`, list of int), `task_category` (`TaskCategory`, int), `ip_type` (`IPType`, int), `name` (`Name`, string), `region_id` (`RegionID`, int), `district_id` (`DistrictID`, int), `net_service_id` (`NetServiceID`, int), `node_group_type` (`NodeGroupType`, int), `task_type` (`TaskType`, int), `probe_type` (`ProbeType`, int). Only arguments provided by the user SHALL be set on the request; unset arguments SHALL NOT be sent (relying on API defaults).

#### Scenario: All filter arguments optional
- **WHEN** a user declares `data "tencentcloud_cat_node_groups" "groups" {}` with no filter arguments
- **THEN** the Read function SHALL build a `DescribeNodeGroupsRequest` with no filter fields set and the API SHALL return the default (unfiltered) node-group tree

#### Scenario: Single filter applied
- **WHEN** a user sets `node_type = [1]` and `ip_type = 1`
- **THEN** the Read function SHALL set `request.NodeType = []*int64{int64Ptr(1)}` and `request.IPType = int64Ptr(1)`, and SHALL leave all other request fields unset

#### Scenario: node_type accepts multiple values
- **WHEN** a user sets `node_type = [1, 2]`
- **THEN** the Read function SHALL convert each element to `*int64` and assign the resulting slice to `request.NodeType`

### Requirement: The data source SHALL flatten `DescribeNodeGroups` response into computed attributes
The Read function SHALL flatten the `DescribeNodeGroupsResponse.Response` into three computed `TypeList` attributes:
- `node_list`: one entry per `NodeList` item (`NodeTree`). Each entry SHALL expose `id` (`NodeTree.ID`, string), `content` (`NodeTree.Content`, string), and `children` (a `TypeList`).
  - Each `children` entry (`NodeLeaf`) SHALL expose `id` (`NodeLeaf.ID`), `content` (`NodeLeaf.Content`), and `children` (a `TypeList`).
    - Each inner `children` entry (`NodeInfoBase`) SHALL expose `id` (`NodeInfoBase.ID`) and `content` (`NodeInfoBase.Content`).
- `district_list`: one entry per `DistrictList` item (`DistinctOrNetServiceInfo`), exposing `id` (`ID`) and `name` (`Name`).
- `net_service_list`: one entry per `NetServiceList` item (`DistinctOrNetServiceInfo`), exposing `id` (`ID`) and `name` (`Name`).

Before calling `d.Set` for any field, the Read function SHALL skip fields whose underlying pointer is `nil`. Empty but non-nil lists SHALL be flattened to empty Terraform lists (not an error).

#### Scenario: Fully populated response
- **WHEN** `DescribeNodeGroups` returns a `NodeList` with two `NodeTree` entries, each having `Children`, plus non-empty `DistrictList` and `NetServiceList`
- **THEN** `node_list` SHALL contain two entries with nested `children`, `district_list` and `net_service_list` SHALL each contain their respective entries, and every `id`/`content`/`name` value SHALL be populated from the response

#### Scenario: Nil pointer fields skipped
- **WHEN** a `NodeTree` entry has a non-nil `ID` but a nil `Content`
- **THEN** the flattened map SHALL include `id` but SHALL omit the `content` key

### Requirement: The Read function SHALL treat a nil API response as retryable, not as state-clearing
The Read function SHALL wrap the service call in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. Inside the retry block, if the response is structurally nil (`result == nil`, `result.Response == nil`, or all of `NodeList`, `DistrictList`, and `NetServiceList` are nil), the function SHALL return `resource.NonRetryableError(...)` so the outer retry continues — it SHALL NOT call `d.SetId("")`. On the retry failure path it SHALL log `[DATASOURCE] read empty, skip SetId` for diagnostics.

#### Scenario: Transient nil response
- **WHEN** `DescribeNodeGroups` returns a nil response or `Response == nil`
- **THEN** the retry block SHALL return a `NonRetryableError`, the id SHALL NOT be cleared, and the overall Read SHALL eventually fail with a retry-exhausted error

#### Scenario: Empty-but-valid response
- **WHEN** `DescribeNodeGroups` returns a non-nil `Response` where `NodeList`, `DistrictList`, and `NetServiceList` are all empty (non-nil, length 0)
- **THEN** the Read SHALL succeed, set empty `node_list`/`district_list`/`net_service_list`, and set a valid id — SHALL NOT error

### Requirement: The data source id SHALL be deterministic
After a successful read, the Read function SHALL set the data source id to `helper.DataResourceIdsHash(ids)` where `ids` is the slice of top-level `NodeTree.ID` values from the response. This keeps the id deterministic and consistent with the sibling `tencentcloud_cat_node` data source.

#### Scenario: Id derived from node list
- **WHEN** the response `NodeList` contains entries with IDs `"A"` and `"B"`
- **THEN** `d.Id()` SHALL equal `helper.DataResourceIdsHash([]string{"A", "B"})`

### Requirement: `result_output_file` SHALL write the flattened node list
The data source SHALL accept an optional `result_output_file` argument (string). When set to a non-empty path, after a successful read the Read function SHALL write the flattened `node_list` list (`[]map[string]interface{}`) to that file via `tccommon.WriteToFile`.

#### Scenario: Output file written
- **WHEN** `result_output_file = "/tmp/out.json"` and the read succeeds
- **THEN** the flattened node-list maps SHALL be written to `/tmp/out.json`

### Requirement: The provider SHALL register the data source and docs SHALL be generated
The provider entry `tencentcloud/provider.go` SHALL add `"tencentcloud_cat_node_groups": cat.DataSourceTencentCloudCatNodeGroups()` to `DataSourcesMap`, and the provider.go doc-comment index SHALL list the new data source under the CAT product so that `make doc` regenerates `website/docs/d/cat_node_groups.html.markdown` and updates `website/tencentcloud.erb`. A hand-written example file `tencentcloud/services/cat/data_source_tc_cat_node_groups.md` SHALL be created with a one-line description (mentioning CAT) and an `Example Usage` HCL block.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `DataSourcesMap["tencentcloud_cat_node_groups"]` SHALL be non-nil

#### Scenario: Documentation example file
- **WHEN** `make doc` runs
- **THEN** it SHALL generate `website/docs/d/cat_node_groups.html.markdown` from the example `.md` file and the schema `Description` fields

### Requirement: Unit tests SHALL mock the SDK with gomonkey
The test file `tencentcloud/services/cat/data_source_tc_cat_node_groups_test.go` SHALL use gomonkey to mock the CAT SDK client (NOT a Terraform acceptance test suite). Tests SHALL cover: (a) a successful read that flattens a populated `DescribeNodeGroups` response and asserts the `node_list`/`district_list`/`net_service_list` values and a non-empty id; (b) a nil/empty response path; (c) schema field existence and types. The mocked provider meta SHALL implement `tccommon.ProviderMeta`.

#### Scenario: Successful read mocked
- **WHEN** the gomonkey-patched `DescribeNodeGroups` returns a response with one `NodeTree` (with one `NodeLeaf` child having one `NodeInfoBase` grandchild), one district, and one ISP
- **THEN** calling `res.Read(d, meta)` SHALL succeed, `d.Id()` SHALL be non-empty, and `d.Get("node_list")` SHALL contain one entry whose nested `children` contains one entry whose nested `children` contains one entry

#### Scenario: Schema correctness
- **WHEN** the data source schema is inspected
- **THEN** it SHALL contain `node_list`, `district_list`, `net_service_list`, `result_output_file`, and all filter arguments with the correct types
