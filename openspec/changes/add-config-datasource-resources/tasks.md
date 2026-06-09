## 1. Service Layer

- [x] 1.1 Append `DescribeConfigDiscoveredResourcesByFilter()` — wraps `ListDiscoveredResources` with NextToken loop + retry (MaxResults=200)
- [x] 1.2 Append `DescribeConfigResourceTypes()` — single call with retry, no pagination

## 2. Data Source: tencentcloud_config_discovered_resources

- [x] 2.1 Create `data_source_tc_config_discovered_resources.go` with schema and read handler
- [x] 2.2 Optional filters: `filters` (list), `tags` (list), `order_type`; computed `resource_list`
- [x] 2.3 Read handler: build request, call service, flatten into `resource_list`

## 3. Data Source: tencentcloud_config_resource_types

- [x] 3.1 Create `data_source_tc_config_resource_types.go` with schema and read handler
- [x] 3.2 Computed `resource_type_list`; no input filters
- [x] 3.3 Read handler: call service, flatten into `resource_type_list`

## 4. Provider Registration

- [x] 4.1 Register `tencentcloud_config_discovered_resources`
- [x] 4.2 Register `tencentcloud_config_resource_types`

## 5. Documentation & Tests

- [x] 5.1 Create `data_source_tc_config_discovered_resources.md`
- [x] 5.2 Create `data_source_tc_config_resource_types.md`
- [x] 5.3 Create `data_source_tc_config_discovered_resources_test.go`
- [x] 5.4 Create `data_source_tc_config_resource_types_test.go`

