## Why

The `tencentcloud_mariadb_hour_db_instance` resource currently hardcodes database initialization parameters (character_set_server, lower_case_table_names, sync_mode, innodb_page_size) when calling the `InitDBInstances` API. Users cannot customize these parameters, which limits flexibility for different database configurations.

## What Changes

- Add a new Optional + ForceNew parameter `init_params` to the `tencentcloud_mariadb_hour_db_instance` resource schema, allowing users to specify custom initialization parameters when creating an instance.
- The parameter maps to `request.InitParams` in the `CreateHourDBInstance` API, which accepts a list of `DBParamValue` objects (each with `param` and `value` string fields).
- When `init_params` is configured, the user-provided values are used for instance initialization; when not configured, the existing hardcoded defaults are preserved.

## Capabilities

### New Capabilities
- `init-params-config`: Expose the `init_params` parameter in the `tencentcloud_mariadb_hour_db_instance` resource schema to allow user-configurable database initialization parameters during instance creation.

### Modified Capabilities

## Impact

- Resource file: `tencentcloud/services/mariadb/resource_tc_mariadb_hour_db_instance.go` — add schema field and modify create logic
- Test file: `tencentcloud/services/mariadb/resource_tc_mariadb_hour_db_instance_test.go` — add test cases
- Documentation: `tencentcloud/services/mariadb/resource_tc_mariadb_hour_db_instance.md` — update example usage
