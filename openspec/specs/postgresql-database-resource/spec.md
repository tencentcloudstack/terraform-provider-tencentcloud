# postgresql-database-resource Specification

## Purpose
Defines the behavior of the `tencentcloud_postgresql_database` resource, which manages the lifecycle of a PostgreSQL database within a DB instance.
## Requirements
### Requirement: Resource schema for tencentcloud_postgresql_database
The system SHALL provide a Terraform resource named `tencentcloud_postgresql_database` that manages a PostgreSQL database within a DB instance. The resource SHALL support the following schema fields:
- `db_instance_id` (string, Required, ForceNew) — the DB instance ID.
- `database_name` (string, Required, ForceNew) — the database name.
- `database_owner` (string, Required) — the database owner account.
- `encoding` (string, Optional, Computed, ForceNew) — the database character encoding (default UTF8).
- `collate` (string, Optional, Computed, ForceNew) — the database collation rule.
- `ctype` (string, Optional, Computed, ForceNew) — the database character classification.

The resource SHALL support import via a composite ID of format `db_instance_id#database_name`.

#### Scenario: Resource with all required fields
- **WHEN** a user defines a `tencentcloud_postgresql_database` resource with `db_instance_id`, `database_name`, and `database_owner`
- **THEN** the resource SHALL be accepted as valid configuration

#### Scenario: Resource with optional encoding fields
- **WHEN** a user defines a `tencentcloud_postgresql_database` resource with `encoding`, `collate`, and `ctype` set
- **THEN** the resource SHALL accept these optional fields and mark them as ForceNew and Computed so that server-defaulted values (e.g. `collate`/`ctype` defaulting to `C`) do not trigger a replacement diff when omitted from configuration

#### Scenario: Import existing database
- **WHEN** a user imports `tencentcloud_postgresql_database` using a composite ID `db_instance_id#database_name`
- **THEN** the system SHALL split the ID and populate the state with the database details

### Requirement: Service-layer encapsulation
The system SHALL encapsulate all PostgreSQL database SDK calls as reusable methods on the `PostgresqlService` (in `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go`): `CreatePostgresqlDatabase`, `DescribePostgresqlDatabaseById`, `ModifyPostgresqlDatabaseOwner`, and `DeletePostgresqlDatabaseById`. Each write method SHALL apply `ratelimit.Check(request.GetAction())` and `resource.Retry(tccommon.WriteRetryTimeout, ...)`, returning `tccommon.RetryError` on SDK failure and a `NonRetryableError` on a nil response. The read method SHALL apply `ratelimit.Check(request.GetAction())` and return `nil` (no database) when the database is not found. The resource CRUD functions SHALL invoke these service methods rather than calling the SDK client directly.

#### Scenario: Write operations go through the service layer
- **WHEN** the resource creates, updates, or deletes a database
- **THEN** the resource SHALL call the corresponding `PostgresqlService` method instead of constructing an SDK request directly

#### Scenario: Read operation goes through the service layer
- **WHEN** the resource reads a database
- **THEN** the resource SHALL call `DescribePostgresqlDatabaseById` and treat a `nil` return value as "not found"

### Requirement: Create database
The system SHALL create a PostgreSQL database by calling the `CreatePostgresqlDatabase` service method, which invokes the `CreateDatabase` API with `DBInstanceId`, `DatabaseName`, `DatabaseOwner`, and optionally `Encoding`, `Collate`, `Ctype`. After successful creation, the system SHALL set the resource ID to `db_instance_id#database_name` (using `FILED_SP` separator) and call the Read function to populate state.

The service method SHALL verify that the API response is not nil before proceeding. If the API call fails, the service method SHALL wrap the error using `tccommon.RetryError`.

#### Scenario: Successful database creation
- **WHEN** the user creates a `tencentcloud_postgresql_database` with valid parameters
- **THEN** the system SHALL call `CreateDatabase` and set the composite ID `db_instance_id#database_name`

#### Scenario: Create with optional encoding
- **WHEN** the user provides `encoding`, `collate`, or `ctype` during creation
- **THEN** the system SHALL pass these values to the `CreateDatabase` API request

#### Scenario: Create API returns nil response
- **WHEN** the `CreateDatabase` API returns a nil response
- **THEN** the system SHALL return a `NonRetryableError`

### Requirement: Read database
The system SHALL read a PostgreSQL database by calling the `DescribePostgresqlDatabaseById` service method, which invokes the `DescribeDatabases` API with `DBInstanceId` and a `database-name` filter. The service method SHALL iterate the returned `Databases` array and match the exact `DatabaseName`, returning `nil` if no exact match is found. If the database is found, the system SHALL set `database_name`, `database_owner`, `encoding`, `collate`, and `ctype` from the `Database` struct fields. If the database is not found, the system SHALL set the resource ID to empty string after printing a log with the current ID.

#### Scenario: Database found
- **WHEN** the system reads a database that exists
- **THEN** the system SHALL populate state fields from the `Database` struct

#### Scenario: Database not found
- **WHEN** the system reads a database that does not exist
- **THEN** the system SHALL print a log with the resource ID and set the resource ID to empty string

#### Scenario: Setting fields from nil response values
- **WHEN** the `Database` struct has nil field values
- **THEN** the system SHALL skip calling `d.Set()` for those nil fields

### Requirement: Update database owner
The system SHALL update the database owner by calling the `ModifyPostgresqlDatabaseOwner` service method, which invokes the `ModifyDatabaseOwner` API with `DBInstanceId`, `DatabaseName`, and `DatabaseOwner`, when the `database_owner` field changes. If `encoding`, `collate`, or `ctype` change, the Terraform SDK SHALL trigger recreation (ForceNew) rather than an in-place update.

#### Scenario: Update database owner
- **WHEN** the user changes the `database_owner` field
- **THEN** the system SHALL call `ModifyDatabaseOwner` with the new owner value

#### Scenario: Attempting to update immutable fields
- **WHEN** the user changes `encoding`, `collate`, or `ctype`
- **THEN** the Terraform SDK SHALL trigger resource recreation (ForceNew)

### Requirement: Delete database
The system SHALL delete a PostgreSQL database by calling the `DeletePostgresqlDatabaseById` service method, which invokes the `DeleteDatabase` API with `DBInstanceId` and `DatabaseName`. The system SHALL parse the composite ID to extract these values.

#### Scenario: Successful database deletion
- **WHEN** the user deletes a `tencentcloud_postgresql_database` resource
- **THEN** the system SHALL call `DeletePostgresqlDatabaseById` with `DBInstanceId` and `DatabaseName` parsed from the composite ID

### Requirement: Provider registration
The system SHALL register the `tencentcloud_postgresql_database` resource in `tencentcloud/provider.go` and add the resource name to `tencentcloud/provider.md` in the PostgreSQL section.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** the `tencentcloud_postgresql_database` resource SHALL be registered and available for use

### Requirement: Documentation
The system SHALL provide a documentation file `resource_tc_postgresql_database.md` that includes a one-line description referencing TencentDB for PostgreSQL, example usage with HCL, and import instructions using the composite ID.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a `resource_tc_postgresql_database.md` file SHALL exist with description, example usage, and import section
