# monitor-tmp-instance-long-term-storage-retention-time Specification

## Purpose
Support configuring the archive storage retention duration (LongTermStorageRetentionTime) for Prometheus instances via the `tencentcloud_monitor_tmp_instance` resource's `InstanceAttributes` parameter.

## Requirements
### Requirement: Resource supports long_term_storage_retention_time parameter
The `tencentcloud_monitor_tmp_instance` resource SHALL support an optional `long_term_storage_retention_time` parameter (int type, range 60-730) that allows users to configure the archive storage retention duration for the Prometheus instance.

#### Scenario: Create instance with long_term_storage_retention_time specified
- **WHEN** user specifies `long_term_storage_retention_time` in the resource configuration
- **THEN** the system SHALL pass the `InstanceAttributes` field with key `LongTermStorageRetentionTime` and the specified value to the `CreatePrometheusMultiTenantInstancePostPayMode` API request

#### Scenario: Create instance without long_term_storage_retention_time
- **WHEN** user does not specify `long_term_storage_retention_time` in the resource configuration
- **THEN** the system SHALL NOT pass the `InstanceAttributes` field to the Create API request, and the resource SHALL be created successfully without archive storage configuration

#### Scenario: Read long_term_storage_retention_time from API
- **WHEN** the resource Read function is invoked
- **THEN** the system SHALL read the `InstanceAttributes` array from the `DescribePrometheusInstances` API response, find the entry with key `LongTermStorageRetentionTime`, parse its value as integer, and set it in state

#### Scenario: Modify long_term_storage_retention_time after creation
- **WHEN** user changes the `long_term_storage_retention_time` value in an existing resource configuration
- **THEN** the system SHALL pass the `InstanceAttributes` field with key `LongTermStorageRetentionTime` and the new value to the `ModifyPrometheusInstanceAttributes` API request

#### Scenario: Import existing instance with long_term_storage_retention_time
- **WHEN** user imports an existing Prometheus instance that has `LongTermStorageRetentionTime` set in `InstanceAttributes`
- **THEN** the system SHALL correctly read and populate the `long_term_storage_retention_time` attribute in state
