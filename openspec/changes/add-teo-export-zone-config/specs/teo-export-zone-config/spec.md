## ADDED Requirements

### Requirement: Export TEO Zone Config

The system SHALL provide a data source that allows users to export the complete configuration of a TEO (TencentCloud EdgeOne) zone.

#### Scenario: Query zone config by zone ID
- **GIVEN** a TEO zone exists with zone ID "zone-12345678"
- **WHEN** the user queries the data source with `zone_id = "zone-12345678"`
- **THEN** the system returns the complete zone configuration including basic info, acceleration settings, security rules, and origin settings

#### Scenario: Query zone config by zone name
- **GIVEN** a TEO zone exists with zone name "example.com"
- **WHEN** the user queries the data source with `zone_name = "example.com"`
- **THEN** the system returns the complete zone configuration corresponding to that zone name

#### Scenario: Handle missing zone
- **GIVEN** no TEO zone exists with zone ID "zone-nonexistent"
- **WHEN** the user queries the data source with `zone_id = "zone-nonexistent"`
- **THEN** the system returns an appropriate error message indicating the zone was not found

#### Scenario: Export includes basic configuration
- **GIVEN** a TEO zone exists with basic configuration
- **WHEN** the user queries the data source
- **THEN** the exported configuration includes zone ID, zone name, status, and other basic information

#### Scenario: Export includes acceleration settings
- **GIVEN** a TEO zone exists with acceleration configuration
- **WHEN** the user queries the data source
- **THEN** the exported configuration includes acceleration settings such as caching rules, cache key, and optimization settings

#### Scenario: Export includes security rules
- **GIVEN** a TEO zone exists with security configuration
- **WHEN** the user queries the data source
- **THEN** the exported configuration includes security rules such as WAF rules, rate limiting, and access control settings

#### Scenario: Export includes origin settings
- **GIVEN** a TEO zone exists with origin configuration
- **WHEN** the user queries the data source
- **THEN** the exported configuration includes origin server settings, protocol, and load balancing configuration

#### Scenario: Query with both zone ID and zone name
- **GIVEN** a TEO zone exists
- **WHEN** the user queries the data source with both `zone_id` and `zone_name` parameters
- **THEN** the system uses the zone_id parameter and ignores zone_name (zone_id takes precedence)
