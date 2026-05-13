## MODIFIED Requirements

### Requirement: Remove zone_setting computed attribute
The `tencentcloud_teo_l7_acc_setting` resource SHALL NOT have a `zone_setting` computed attribute. The `zone_name` and `zone_config` top-level attributes provide the same data.

#### Scenario: zone_setting is removed from schema
- **WHEN** a user inspects the `tencentcloud_teo_l7_acc_setting` resource schema
- **THEN** there SHALL NOT be a `zone_setting` attribute

## ADDED Requirements

### Requirement: network_error_logging configuration in zone_config
The `tencentcloud_teo_l7_acc_setting` resource's `zone_config` attribute SHALL include a `network_error_logging` sub-field of type `TypeList` with `MaxItems: 1` and `Optional: true`, mapping to `ZoneConfig.NetworkErrorLogging` (`NetworkErrorLoggingParameters`) in the cloud API.

The `network_error_logging` block SHALL contain:
- `switch` (TypeString, Optional): Whether to enable network error logging. Valid values: `on`, `off`.

#### Scenario: network_error_logging is set via Update
- **WHEN** a user configures `zone_config.network_error_logging` with `switch = "on"` or `switch = "off"`
- **THEN** the `ModifyL7AccSetting` API SHALL be called with `ZoneConfig.NetworkErrorLogging.Switch` set to the specified value

#### Scenario: network_error_logging is read from API
- **WHEN** the Read function is called for an existing `tencentcloud_teo_l7_acc_setting` resource
- **THEN** the `zone_config.network_error_logging.switch` attribute SHALL be populated from `respData.ZoneConfig.NetworkErrorLogging.Switch` when not nil
