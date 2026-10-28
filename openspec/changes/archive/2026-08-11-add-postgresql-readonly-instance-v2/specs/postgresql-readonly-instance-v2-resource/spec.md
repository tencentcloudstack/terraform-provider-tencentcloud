## ADDED Requirements

### Requirement: Create PostgreSQL readonly instance

The system SHALL allow users to create a PostgreSQL readonly instance by calling the `CreateReadOnlyDBInstance` cloud API. The resource SHALL accept all input parameters of `CreateReadOnlyDBInstance` as schema fields, with `zone`, `master_db_instance_id`, `spec_code`, `storage`, `instance_count`, and `period` as required fields.

The resource ID SHALL be set to the first element of the `DBInstanceIdSet` returned by the API. After creation, the system SHALL poll `DescribeDBInstanceAttribute` until the instance status is `running` to confirm the async creation has taken effect.

#### Scenario: Successful creation of a readonly instance
- **WHEN** the user applies a `tencentcloud_postgresql_readonly_instance_v2` resource configuration with valid required fields (`zone`, `master_db_instance_id`, `spec_code`, `storage`, `instance_count`, `period`)
- **THEN** the system SHALL call `CreateReadOnlyDBInstance` with the provided parameters, set the resource ID from `DBInstanceIdSet[0]`, and poll until the instance is `running`

#### Scenario: Creation returns empty instance ID set
- **WHEN** `CreateReadOnlyDBInstance` returns successfully but `DBInstanceIdSet` is empty
- **THEN** the system SHALL return a `NonRetryableError` indicating the instance ID is empty, without setting the resource ID

#### Scenario: Creation API call fails
- **WHEN** `CreateReadOnlyDBInstance` returns an error
- **THEN** the system SHALL wrap the error with `tccommon.RetryError()` and retry up to `WriteRetryTimeout`, returning the wrapped error if retries are exhausted

### Requirement: Read PostgreSQL readonly instance

The system SHALL allow users to read a PostgreSQL readonly instance by calling the `DescribeDBInstanceAttribute` cloud API with the `DBInstanceId` obtained from the resource ID. The system SHALL backfill all available schema fields from the returned `DBInstance` structure.

When the API returns an empty response (`response == nil` or `DBInstance == nil`), the system SHALL print a现场 log with the resource ID before clearing the ID with `d.SetId("")`.

#### Scenario: Successful read of an existing instance
- **WHEN** the user runs `terraform plan` or `terraform refresh` on an existing `tencentcloud_postgresql_readonly_instance_v2` resource
- **THEN** the system SHALL call `DescribeDBInstanceAttribute` with the resource ID, and set all schema fields from the returned `DBInstance` structure after nil-checking each field

#### Scenario: Read returns empty response
- **WHEN** `DescribeDBInstanceAttribute` returns an empty `DBInstance` (instance not found or deleted)
- **THEN** the system SHALL print `log.Printf("[CRUD] postgresql_readonly_instance_v2 id=%s", d.Id())` to preserve the现场, then call `d.SetId("")` to remove the resource from state

### Requirement: Update PostgreSQL readonly instance (immutable mode)

The system SHALL treat the resource as CRD-only (Create-Read-Delete). All top-level fields except the ID SHALL be immutable. The Update method SHALL check changed fields against an `immutableArgs` array; if any immutable field has changed, the system SHALL return an error.

#### Scenario: Attempt to update an immutable field
- **WHEN** the user changes any top-level field (e.g., `storage`, `spec_code`, `name`) on an existing `tencentcloud_postgresql_readonly_instance_v2` resource
- **THEN** the system SHALL detect the change via `immutableArgs` check and return an error indicating the field cannot be updated; the resource SHALL be recreated if ForceNew fields change

### Requirement: Delete (isolate) PostgreSQL readonly instance

The system SHALL allow users to delete a PostgreSQL readonly instance by calling the `IsolateDBInstances` cloud API with the `DBInstanceIdSet` containing the resource ID. The system SHALL wrap API errors with `tccommon.RetryError()` and retry up to `WriteRetryTimeout`.

After isolation, the system SHALL poll `DescribeDBInstanceAttribute` until the instance status becomes `isolated` to confirm the async isolation has taken effect.

#### Scenario: Successful isolation of an instance
- **WHEN** the user runs `terraform destroy` on a `tencentcloud_postgresql_readonly_instance_v2` resource
- **THEN** the system SHALL call `IsolateDBInstances` with the resource ID, and poll `DescribeDBInstanceAttribute` until the status is `isolated`

#### Scenario: Isolation API call fails
- **WHEN** `IsolateDBInstances` returns an error
- **THEN** the system SHALL wrap the error with `tccommon.RetryError()` and retry up to `WriteRetryTimeout`, returning the wrapped error if retries are exhausted

### Requirement: Import PostgreSQL readonly instance

The system SHALL support importing an existing PostgreSQL readonly instance via `terraform import` using the instance ID. The import SHALL use `schema.ImportStatePassthrough`.

#### Scenario: Successful import by instance ID
- **WHEN** the user runs `terraform import tencentcloud_postgresql_readonly_instance_v2.foo <instance_id>`
- **THEN** the system SHALL set the resource ID to the provided instance ID and run the Read operation to populate state

### Requirement: Provider registration and documentation

The system SHALL register the `tencentcloud_postgresql_readonly_instance_v2` resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`. The system SHALL create a resource documentation file `resource_tc_postgresql_readonly_instance_v2.md` in the postgresql service directory.

#### Scenario: Resource is registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_postgresql_readonly_instance_v2` SHALL be listed in the provider's resources map and documented in `provider.md`

### Requirement: Unit tests with gomonkey mock

The system SHALL provide unit tests in `resource_tc_postgresql_readonly_instance_v2_test.go` using gomonkey to mock cloud API calls. Tests SHALL cover Create, Read, and Delete business logic paths including success and error scenarios.

#### Scenario: Unit tests verify business logic
- **WHEN** the unit tests are executed
- **THEN** the tests SHALL mock `CreateReadOnlyDBInstance`, `DescribeDBInstanceAttribute`, and `IsolateDBInstances` responses, and assert correct state transitions without calling real cloud APIs
