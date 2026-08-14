# postgres-database-resource Specification

## Purpose
TBD - created by archiving change add-postgres-database-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource schema for tencentcloud_postgres_database
The system SHALL provide a Terraform resource named `tencentcloud_postgres_database` that manages a PostgreSQL database within a DB instance. The resource SHALL support the following schema fields:
- `db_instance_id` (string, Required, ForceNew) — the DB instance ID.
- `database_name` (string, Required, ForceNew) — the database name.
- `database_owner` (string, Required) — the database owner account.
- `encoding` (string, Optional, ForceNew) — the database character encoding (default UTF8).
- `collate` (string, Optional, ForceNew) — the database collation rule.
- `ctype` (string, Optional, ForceNew) — the database character classification.

The resource SHALL support import via a composite ID of format `db_instance_id#database_name`.

#### Scenario: Resource with all required fields
- **WHEN** a user defines a `tencentcloud_postgres_database` resource with `db_instance_id`, `database_name`, and `database_owner`
- **THEN** the resource SHALL be accepted as valid configuration

#### Scenario: Resource with optional encoding fields
- **WHEN** a user defines a `tencentcloud_postgres_database` resource with `encoding`, `collate`, and `ctype` set
- **THEN** the resource SHALL accept these optional fields and mark them as ForceNew

#### Scenario: Import existing database
- **WHEN** a user imports `tencentcloud_postgres_database` using a composite ID `db_instance_id#database_name`
- **THEN** the system SHALL split the ID and populate the state with the database details

### Requirement: Create database
The system SHALL create a PostgreSQL database by calling the `CreateDatabase` API with `DBInstanceId`, `DatabaseName`, `DatabaseOwner`, and optionally `Encoding`, `Collate`, `Ctype`. After successful creation, the system SHALL set the resource ID to `db_instance_id#database_name` (using `FILED_SP` separator) and call the Read function to populate state.

The system SHALL verify that the API response is not nil before proceeding. If the API call fails, the system SHALL wrap the error using `tccommon.RetryError`.

#### Scenario: Successful database creation
- **WHEN** the user creates a `tencentcloud_postgres_database` with valid parameters
- **THEN** the system SHALL call `CreateDatabase` and set the composite ID `db_instance_id#database_name`

#### Scenario: Create with optional encoding
- **WHEN** the user provides `encoding`, `collate`, or `ctype` during creation
- **THEN** the system SHALL pass these values to the `CreateDatabase` API request

#### Scenario: Create API returns nil response
- **WHEN** the `CreateDatabase` API returns a nil response
- **THEN** the system SHALL return a `NonRetryableError`

### Requirement: Read database
The system SHALL read a PostgreSQL database by calling the `DescribeDatabases` API with `DBInstanceId` and a `database-name` filter. The system SHALL iterate the returned `Databases` array and match the exact `DatabaseName`. If the database is found, the system SHALL set `database_name`, `database_owner`, `encoding`, `collate`, and `ctype` from the `Database` struct fields. If the database is not found, the system SHALL set the resource ID to empty string after printing a log with the current ID.

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
The system SHALL update the database owner by calling the `ModifyDatabaseOwner` API when the `database_owner` field changes. The system SHALL pass `DBInstanceId`, `DatabaseName`, and `DatabaseOwner` to the API. If `encoding`, `collate`, or `ctype` change, the Terraform SDK SHALL trigger recreation (ForceNew) rather than an in-place update.

#### Scenario: Update database owner
- **WHEN** the user changes the `database_owner` field
- **THEN** the system SHALL call `ModifyDatabaseOwner` with the new owner value

#### Scenario: Attempting to update immutable fields
- **WHEN** the user changes `encoding`, `collate`, or `ctype`
- **THEN** the Terraform SDK SHALL trigger resource recreation (ForceNew)

### Requirement: Delete database
The system SHALL delete a PostgreSQL database by calling the `DeleteDatabase` API with `DBInstanceId` and `DatabaseName`. The system SHALL parse the composite ID to extract these values.

#### Scenario: Successful database deletion
- **WHEN** the user deletes a `tencentcloud_postgres_database` resource
- **THEN** the system SHALL call `DeleteDatabase` with `DBInstanceId` and `DatabaseName` parsed from the composite ID

### Requirement: Provider registration
The system SHALL register the `tencentcloud_postgres_database` resource in `tencentcloud/provider.go` and add the resource name to `tencentcloud/provider.md` in the PostgreSQL section.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** the `tencentcloud_postgres_database` resource SHALL be registered and available for use

### Requirement: Documentation
The system SHALL provide a documentation file `resource_tc_postgres_database.md` that includes a one-line description referencing TencentDB for PostgreSQL, example usage with HCL, and import instructions using the composite ID.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a `resource_tc_postgres_database.md` file SHALL exist with description, example usage, and import section

