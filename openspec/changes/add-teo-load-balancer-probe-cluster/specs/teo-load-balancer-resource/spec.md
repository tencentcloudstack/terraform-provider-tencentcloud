## MODIFIED Requirements

### Requirement: HealthChecker nested block
The system SHALL support `health_checker` as a TypeList (MaxItems: 1) nested block with the following fields: `type` (required), `port` (optional), `interval` (optional), `timeout` (optional), `health_threshold` (optional), `critical_threshold` (optional), `path` (optional), `method` (optional), `expected_codes` (optional, list of strings), `headers` (optional, list of nested blocks with `key` and `value`), `follow_redirect` (optional), `send_context` (optional), `recv_context` (optional), `probe_cluster` (optional, TypeString).

The `probe_cluster` field specifies the probe cluster region scope for health checks. Valid values: `global`, `mainland_china`, `overseas`. Default: `global`. Only supported when the load balancer instance type is `HTTP_V2`.

The system SHALL map `probe_cluster` to the `HealthChecker.ProbeCluster` API field in Create and Update operations via `buildTeoHealthChecker`, and populate it from the API response in Read operations via `flattenTeoHealthChecker`.

#### Scenario: Configure HTTP-health-checker with probe_cluster
- **WHEN** user specifies `health_checker` with type="HTTP", port=80, interval=30, path="/health", probe_cluster="mainland_china"
- **THEN** the system SHALL create a `HealthChecker` API struct with the corresponding fields including ProbeCluster="mainland_china"

#### Scenario: probe_cluster not specified
- **WHEN** user does not specify `probe_cluster` in the health_checker block
- **THEN** the system SHALL not send ProbeCluster in the Create/Modify request (nil), and the cloud API will default to `global`

#### Scenario: Read populates probe_cluster from API response
- **WHEN** the API response's HealthChecker.ProbeCluster is non-nil
- **THEN** the system SHALL populate the `probe_cluster` schema field from the response value

#### Scenario: Read skips probe_cluster when API returns nil
- **WHEN** the API response's HealthChecker.ProbeCluster is nil
- **THEN** the system SHALL not set the `probe_cluster` schema field

#### Scenario: No health checker specified
- **WHEN** user does not specify `health_checker`
- **THEN** the system SHALL not send HealthChecker in the Create/Modify request, and the cloud API will default to NoCheck (no health check)

### Requirement: Schema field definitions
The system SHALL define the following schema fields for `tencentcloud_teo_load_balancer`:

#### Scenario: Required fields
- **WHEN** user creates a load balancer resource
- **THEN** the following fields SHALL be required: `zone_id` (ForceNew), `name`, `type` (ForceNew), `origin_groups`

#### Scenario: Optional fields
- **WHEN** user creates a load balancer resource
- **THEN** the following fields SHALL be optional: `health_checker`, `steering_policy`, `failover_policy`

#### Scenario: Computed fields
- **WHEN** Terraform reads the resource state
- **THEN** the following fields SHALL be computed (read-only): `instance_id`, `status`, `origin_group_health_status`, `l4_used_list`, `l7_used_list`, `references`

#### Scenario: type field description includes HTTP_V2
- **WHEN** the resource schema is defined
- **THEN** the `type` field Description SHALL include three valid values: `HTTP_V2` (HTTP-specific V2, supports selecting probe region, lower probe volume and lower sensitivity), `HTTP` (HTTP-specific V1, does not support selecting probe region, higher probe volume and higher sensitivity), `GENERAL` (general type, only supports general origin groups, referenced by site acceleration and Layer-4 proxy, does not support selecting probe region)

### Requirement: Resource Update Operation supports probe_cluster
When `health_checker` is specified in the Terraform configuration, the system SHALL include the `probe_cluster` value (if set) in the `HealthChecker` struct passed to the `ModifyLoadBalancer` API call, alongside other health checker fields.

#### Scenario: Update health_checker with probe_cluster
- **WHEN** user updates the `health_checker` block including `probe_cluster`
- **THEN** the system SHALL call `ModifyLoadBalancer` with a `HealthChecker` struct containing the `ProbeCluster` field set to the configured value

#### Scenario: Update health_checker without probe_cluster
- **WHEN** user updates the `health_checker` block without specifying `probe_cluster`
- **THEN** the system SHALL call `ModifyLoadBalancer` with a `HealthChecker` struct where `ProbeCluster` is nil