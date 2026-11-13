# mongodb-scale-up-db-instance-cpu Specification

## Purpose
TBD - created by archiving change add-mongodb-scale-up-db-instance-cpu. Update Purpose after archive.
## Requirements
### Requirement: MongoDB Scale Up DB Instance CPU Resource Schema
The `tencentcloud_mongodb_scale_up_db_instance_cpu` resource SHALL define a schema to manage MongoDB instance elastic CPU scale-up lifecycle, including input fields `instance_id` (Optional, ForceNew) and `extra_cpu` (Optional, not ForceNew), plus computed output fields `flow_id`, `scalable`, `is_scaled`, `is_locked`, `reason`, `max_extra_cpu`, `trigger_type`, `scale_up_time`.

#### Scenario: Schema defines instance_id field
- **GIVEN** the `tencentcloud_mongodb_scale_up_db_instance_cpu` resource schema definition
- **THEN** `instance_id` SHALL be of type `schema.TypeString`
- **AND** `instance_id` SHALL have `Optional: true` and `ForceNew: true`
- **AND** the field description SHALL mention it is the MongoDB instance ID

#### Scenario: Schema defines extra_cpu field
- **GIVEN** the `tencentcloud_mongodb_scale_up_db_instance_cpu` resource schema definition
- **THEN** `extra_cpu` SHALL be of type `schema.TypeInt`
- **AND** `extra_cpu` SHALL have `Optional: true` and SHALL NOT set `ForceNew`
- **AND** the field description SHALL document it as the incremental CPU cores to scale up per node

#### Scenario: Schema defines computed output fields
- **GIVEN** the resource schema definition
- **THEN** `flow_id`, `scalable`, `is_scaled`, `is_locked`, `reason`, `max_extra_cpu`, `trigger_type`, `scale_up_time` SHALL each have `Computed: true`
- **AND** `flow_id` SHALL be `schema.TypeInt`, `scalable`/`is_scaled`/`is_locked` SHALL be `schema.TypeBool`, `reason`/`scale_up_time` SHALL be `schema.TypeString`, `max_extra_cpu`/`trigger_type` SHALL be `schema.TypeInt`

#### Scenario: Resource supports import
- **GIVEN** the resource definition
- **THEN** an `Importer` SHALL be present with `State: schema.ImportStatePassthrough`
- **AND** import SHALL accept the instance id as the resource id

### Requirement: MongoDB Scale Up DB Instance CPU Create Operation
The resource Create operation SHALL call the `ScaleUpDBInstanceCpu` API to manually enable elastic CPU scale-up, handle it as an asynchronous operation, and set the resource id before polling.

#### Scenario: Create calls ScaleUpDBInstanceCpu with request fields
- **GIVEN** a user creates `tencentcloud_mongodb_scale_up_db_instance_cpu` with `instance_id` and `extra_cpu`
- **WHEN** the Create operation runs
- **THEN** the system SHALL build a `ScaleUpDBInstanceCpuRequest` with `InstanceId` and `ExtraCpu` set from the schema
- **AND** the system SHALL call `UseMongodbClient().ScaleUpDBInstanceCpu` within a `resource.Retry(tccommon.WriteRetryTimeout, ...)`, wrapping errors with `tccommon.RetryError`

#### Scenario: Create checks response not nil
- **GIVEN** the ScaleUpDBInstanceCpu API returns
- **WHEN** the result is `nil` or `result.Response` is `nil`
- **THEN** the system SHALL return a `NonRetryableError`

#### Scenario: Create sets id before polling for async completion
- **GIVEN** the ScaleUpDBInstanceCpu API returns a non-nil `Response.FlowId`
- **WHEN** the create flow proceeds
- **THEN** the system SHALL print logId and d.Id() for troubleshooting
- **AND** the system SHALL call `d.SetId(instanceId)` BEFORE starting async polling
- **AND** the system SHALL set the `flow_id` field from `Response.FlowId` (after nil check)

#### Scenario: Create polls CheckDBInstanceElasticCpuScalable until scaled
- **GIVEN** the ScaleUpDBInstanceCpu call succeeded
- **WHEN** polling for async completion
- **THEN** the system SHALL call `CheckDBInstanceElasticCpuScalable` within `resource.Retry(3*tccommon.ReadRetryTimeout, ...)`
- **AND** the system SHALL continue retrying while `IsScaled` is false or `IsLocked` is true
- **AND** API errors SHALL be wrapped with `tccommon.RetryError`

### Requirement: MongoDB Scale Up DB Instance CPU Read Operation
The resource Read operation SHALL call `CheckDBInstanceElasticCpuScalable` to query current elastic CPU scale-up status and populate state.

#### Scenario: Read populates all computed fields when present
- **GIVEN** an existing `tencentcloud_mongodb_scale_up_db_instance_cpu` resource
- **WHEN** the Read operation calls `CheckDBInstanceElasticCpuScalable` with `InstanceId` from the resource id
- **AND** a response field is not nil
- **THEN** the system SHALL set the corresponding schema field (`scalable`, `is_scaled`, `is_locked`, `reason`, `max_extra_cpu`, `extra_cpu`, `trigger_type`, `scale_up_time`)
- **AND** the system SHALL set `instance_id` to the resource id

#### Scenario: Read handles empty response preserving context
- **GIVEN** the CheckDBInstanceElasticCpuScalable response is nil or `response.Response` is nil
- **WHEN** Read runs
- **THEN** the system SHALL first print `log.Printf("[CRUD] tencentcloud_mongodb_scale_up_db_instance_cpu id=%s", d.Id())`
- **AND** then call `d.SetId("")`
- **AND** SHALL NOT call `d.SetId("")` before logging

#### Scenario: Read skips nil fields
- **GIVEN** the response returns a nil field
- **WHEN** Read sets fields
- **THEN** the system SHALL NOT call `d.Set(...)` for that nil field

### Requirement: MongoDB Scale Up DB Instance CPU Update Operation
The resource Update operation SHALL re-trigger `ScaleUpDBInstanceCpu` when `extra_cpu` changes, since `instance_id` is immutable.

#### Scenario: Update re-scales when extra_cpu changes
- **GIVEN** an existing resource
- **WHEN** the user changes `extra_cpu` (and not `instance_id`)
- **THEN** the system SHALL detect `d.HasChange("extra_cpu")`
- **AND** the system SHALL call `ScaleUpDBInstanceCpu` with the new `ExtraCpu` within `resource.Retry(tccommon.WriteRetryTimeout, ...)`
- **AND** the system SHALL set `flow_id` from the new response FlowId after nil check
- **AND** the system SHALL poll `CheckDBInstanceElasticCpuScalable` until `IsScaled` is true

#### Scenario: Update ignores instance_id change
- **GIVEN** the resource schema
- **THEN** `instance_id` SHALL have `ForceNew: true`, so changes trigger recreation rather than update

### Requirement: MongoDB Scale Up DB Instance CPU Delete Operation
The resource Delete operation SHALL call `ScaleDownDBInstanceCpu` to disable elastic CPU scale-up and shrink back, treating it as asynchronous.

#### Scenario: Delete calls ScaleDownDBInstanceCpu
- **GIVEN** a `tencentcloud_mongodb_scale_up_db_instance_cpu` resource is destroyed
- **WHEN** the Delete operation runs
- **THEN** the system SHALL build a `ScaleDownDBInstanceCpuRequest` with `InstanceId` from the resource id
- **AND** the system SHALL call `ScaleDownDBInstanceCpu` within `resource.Retry(tccommon.WriteRetryTimeout, ...)`, wrapping errors with `tccommon.RetryError`

#### Scenario: Delete polls CheckDBInstanceElasticCpuScalable until not scaled
- **GIVEN** the ScaleDownDBInstanceCpu call succeeded
- **WHEN** polling for async completion
- **THEN** the system SHALL call `CheckDBInstanceElasticCpuScalable` within `resource.Retry(3*tccommon.ReadRetryTimeout, ...)`
- **AND** the system SHALL continue retrying while `IsScaled` is true or `IsLocked` is true

#### Scenario: Delete tolerates already-not-scaled state
- **GIVEN** the ScaleDown API returns an error indicating the instance is not in scaled state
- **WHEN** Delete runs
- **THEN** the system SHALL log the error and return nil (treat as already gone), avoiding destroy blocking

### Requirement: MongoDB Scale Up DB Instance CPU Service Layer
The service layer SHALL provide methods wrapping the three cloud APIs with retry handling and async polling, following the `MongodbService` pattern.

#### Scenario: Service method wraps ScaleUpDBInstanceCpu
- **GIVEN** `MongodbService`
- **THEN** a method SHALL call `UseMongodbClient().ScaleUpDBInstanceCpu` within `resource.Retry(tccommon.WriteRetryTimeout, ...)`, returning `FlowId` and error

#### Scenario: Service method wraps ScaleDownDBInstanceCpu
- **GIVEN** `MongodbService`
- **THEN** a method SHALL call `UseMongodbClient().ScaleDownDBInstanceCpu` within `resource.Retry(tccommon.WriteRetryTimeout, ...)`, returning `FlowId` and error

#### Scenario: Service method wraps CheckDBInstanceElasticCpuScalable
- **GIVEN** `MongodbService`
- **THEN** a method SHALL call `UseMongodbClient().CheckDBInstanceElasticCpuScalable` within `resource.Retry(tccommon.ReadRetryTimeout, ...)`, returning the response and error

#### Scenario: Service method polls elastic cpu scaled status
- **GIVEN** `MongodbService`
- **WHEN** polling for scale-up completion
- **THEN** a method SHALL retry calling `CheckDBInstanceElasticCpuScalable` until `IsScaled == true` (scale-up) or `IsScaled == false` (scale-down), with timeout `3*tccommon.ReadRetryTimeout`
- **AND** `IsLocked == true` or not-yet-expected state SHALL produce `resource.RetryableError`
- **AND** API errors SHALL be wrapped with `tccommon.RetryError`

### Requirement: MongoDB Scale Up DB Instance CPU Provider Registration
The provider SHALL register the new resource so it is usable in Terraform configurations.

#### Scenario: Resource registered in provider
- **GIVEN** `tencentcloud/provider.go`
- **THEN** the ResourcesMap SHALL contain `"tencentcloud_mongodb_scale_up_db_instance_cpu": mongodb.ResourceTencentCloudMongodbScaleUpDbInstanceCpu()`

### Requirement: MongoDB Scale Up DB Instance CPU Unit Tests
The resource SHALL have unit tests using gomonkey mocks (not Terraform acceptance suites) covering Create/Read/Update/Delete business logic.

#### Scenario: Unit test mocks cloud APIs via gomonkey
- **GIVEN** the resource test file
- **THEN** tests SHALL mock `UseMongodbClient().ScaleUpDBInstanceCpu`, `ScaleDownDBInstanceCpu`, and `CheckDBInstanceElasticCpuScalable` using gomonkey
- **AND** tests SHALL verify Create sets id and polls until scaled
- **AND** tests SHALL verify Delete polls until not scaled and tolerates already-not-scaled
- **AND** tests SHALL verify Update re-scales on extra_cpu change

