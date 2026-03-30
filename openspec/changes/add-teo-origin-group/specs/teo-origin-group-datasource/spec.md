## ADDED Requirements

### Requirement: Query origin group by ID
The data source SHALL support querying a TEO origin group by its unique identifier. The data source MUST retrieve all configuration details of the specified origin group, including name, type, origin records, host header, references, and timestamps.

#### Scenario: Successful query by origin group ID
- **WHEN** user provides a valid origin group ID
- **THEN** the data source returns the origin group details with all configuration fields
- **AND** the data source includes the origin group name
- **AND** the data source includes the origin group type (GENERAL or HTTP)
- **AND** the data source includes the list of origin records with their configurations
- **AND** the data source includes the host header if configured
- **AND** the data source includes the list of reference instances
- **AND** the data source includes creation and update timestamps

#### Scenario: Query non-existent origin group
- **WHEN** user provides an origin group ID that does not exist
- **THEN** the data source returns an error indicating the origin group was not found
- **AND** no data is returned to the user

### Requirement: Support zone ID parameter
The data source SHALL require a zone ID parameter to specify the TEO zone (site) context. The zone ID MUST be used to scope the origin group query to the correct site context.

#### Scenario: Query with valid zone ID
- **WHEN** user provides a valid zone ID and origin group ID
- **THEN** the data source successfully queries the origin group within the specified zone
- **AND** the returned origin group belongs to the specified zone

### Requirement: Return origin records configuration
The data source SHALL return the complete configuration of origin records, including record value, type, record ID, weight, private authentication settings, and private parameters for COS/S3 origins.

#### Scenario: Origin records with different types
- **WHEN** origin group contains multiple origin records
- **THEN** the data source returns all records with their types (IP_DOMAIN, COS, or AWS_S3)
- **AND** the data source includes record IDs for each origin
- **AND** the data source includes weight values for weighted load balancing
- **AND** the data source includes private authentication settings for COS/S3 origins

#### Scenario: Origin records with private parameters
- **WHEN** origin group contains COS or AWS_S3 type records with private authentication
- **THEN** the data source includes private_parameters list
- **AND** each parameter contains name and value fields
- **AND** parameter names include AccessKeyId, SecretAccessKey, SignatureVersion, and Region as applicable

### Requirement: Return reference instances
The data source SHALL return a list of instances that reference the origin group, including instance type, instance ID, and instance name for each reference.

#### Scenario: Origin group referenced by multiple services
- **WHEN** origin group is referenced by acceleration domains, rule engines, or load balancers
- **THEN** the data source returns all reference instances
- **AND** each reference includes the instance type (AccelerationDomain, RuleEngine, Loadbalance, ApplicationProxy)
- **AND** each reference includes the instance ID
- **AND** each reference includes the instance name

### Requirement: Computed output fields
The data source SHALL provide all origin group fields as computed output fields. The data source MUST NOT support creating, updating, or deleting origin groups.

#### Scenario: All fields are read-only
- **WHEN** user queries the data source
- **THEN** all returned fields are computed
- **AND** no fields can be set by the user
- **AND** no CRUD operations are supported

### Requirement: ID format
The data source SHALL use the origin group ID as the primary identifier. The ID MUST be a string value that uniquely identifies the origin group within the TEO service.

#### Scenario: Data source ID matches origin group ID
- **WHEN** user queries an origin group
- **THEN** the data source ID matches the origin group ID returned by the API
- **AND** the ID can be used to uniquely identify the origin group
