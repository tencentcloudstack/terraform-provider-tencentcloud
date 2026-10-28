## ADDED Requirements

### Requirement: User can configure init_params for mariadb hour db instance

The `tencentcloud_mariadb_hour_db_instance` resource SHALL expose an `init_params` parameter that allows users to specify database initialization parameters at instance creation time. The parameter SHALL be Optional and ForceNew. Each element in the list SHALL contain a `param` field (string, Required) representing the parameter name and a `value` field (string, Required) representing the parameter value.

#### Scenario: User provides custom init_params
- **WHEN** user specifies `init_params` with custom values (e.g., character_set_server=utf8, lower_case_table_names=0)
- **THEN** the resource SHALL pass those values via the `CreateHourDBInstance` API's `InitParams` field and skip the hardcoded `InitDBInstances` call

#### Scenario: User does not provide init_params
- **WHEN** user does not specify `init_params` in the resource configuration
- **THEN** the resource SHALL use the existing hardcoded default values via the `InitDBInstances` API call (character_set_server=utf8mb4, lower_case_table_names=1, sync_mode=2, innodb_page_size=16384)

#### Scenario: Changing init_params forces recreation
- **WHEN** user modifies the `init_params` value in an existing resource configuration
- **THEN** Terraform SHALL plan to destroy and recreate the instance because `init_params` is marked ForceNew
