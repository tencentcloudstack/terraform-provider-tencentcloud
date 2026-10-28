## ADDED Requirements

### Requirement: Support time_zone parameter in CLS COS shipper resource

The `tencentcloud_cls_cos_shipper` resource SHALL support an optional `time_zone` parameter that maps to the `TimeZone` field in the CLS cloud API. This parameter is used to generate time variables in the file path when shipping logs to COS.

The parameter SHALL be:
- Type: string
- Optional: true
- Computed: true
- Description: "The timezone used to generate time variables in the COS file path. Supports GMT/UTC timezone formats (e.g., GMT+08:00, UTC+08:00)."

#### Scenario: Create shipper with time_zone specified

- **WHEN** user creates a `tencentcloud_cls_cos_shipper` resource with `time_zone = "GMT+08:00"`
- **THEN** the `CreateShipper` API request SHALL include `TimeZone` set to `"GMT+08:00"`
- **AND** the resource state SHALL store the `time_zone` value after successful creation

#### Scenario: Create shipper without time_zone specified

- **WHEN** user creates a `tencentcloud_cls_cos_shipper` resource without specifying `time_zone`
- **THEN** the `CreateShipper` API request SHALL NOT include `TimeZone` parameter
- **AND** the resource state SHALL read back whatever value the API returns for `TimeZone` in the Read method

#### Scenario: Read shipper with time_zone

- **WHEN** the resource Read method is called
- **AND** the `ShipperInfo.TimeZone` field in the `DescribeShippers` response is not nil
- **THEN** the resource state SHALL set `time_zone` to the value of `ShipperInfo.TimeZone`

#### Scenario: Read shipper without time_zone

- **WHEN** the resource Read method is called
- **AND** the `ShipperInfo.TimeZone` field in the `DescribeShippers` response is nil
- **THEN** the resource state SHALL NOT set `time_zone` (skip the `d.Set` call)

#### Scenario: Update shipper time_zone

- **WHEN** user updates the `time_zone` field of an existing `tencentcloud_cls_cos_shipper` resource
- **THEN** the `ModifyShipper` API request SHALL include `TimeZone` set to the new value
- **AND** the resource state SHALL reflect the updated value after successful modification

#### Scenario: Import existing shipper with time_zone

- **WHEN** user imports an existing CLS COS shipper resource
- **THEN** the Read method SHALL correctly populate the `time_zone` field from the API response
