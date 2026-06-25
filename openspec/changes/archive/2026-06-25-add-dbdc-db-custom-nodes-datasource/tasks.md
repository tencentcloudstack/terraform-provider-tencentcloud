## 1. Service Layer

- [x] 1.1 Add `DescribeDbdcDbCustomNodesByFilter` method in `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` with internal pagination (Limit=100, increment Offset until all results collected), request construction from paramMap, and response parsing returning NodeSet list

## 2. Datasource Schema and Read Function

- [x] 2.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.go` with `DataSourceTencentCloudDbdcDbCustomNodes()` schema definition including input parameters: `node_ids` (Optional, TypeList of string), `filters` (Optional, TypeList with `name` and `values` sub-fields), `tags` (Optional, TypeList with `key` and `value` sub-fields), `result_output_file` (Optional, TypeString), and computed output `node_set` (TypeList of schema.Resource with all flattened DBCustomNode fields)
- [x] 2.2 Implement `dataSourceTencentCloudDbdcDbCustomNodesRead` function with: defer LogElapsed/InconsistentCheck, paramMap construction from schema inputs, resource.Retry with ReadRetryTimeout calling service method, nil checks before d.Set(), NonRetryableError on empty response, and d.SetId(helper.BuildToken()) on success

## 3. Provider Registration

- [x] 3.1 Add `tencentcloud_dbdc_db_custom_nodes` datasource entry in `tencentcloud/provider.go` mapping to `dbdc.DataSourceTencentCloudDbdcDbCustomNodes()`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.md` with: one-sentence description ("Use this data source to query DB Custom node list in dbdc product"), Example Usage section showing filters/node_ids/tags inputs and node_set output, no Import section, no Argument/Attribute Reference sections

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes_test.go` with gomonkey mock-based unit tests covering: Read function request construction, response parsing with nil field handling, and pagination logic; runnable with `go test -gcflags=all=-l`
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` to verify all test cases pass
