## Why

The `tencentcloud_cls_cos_shipper` resource currently does not support the `TimeZone` parameter, which is used to generate time variables in the file path when shipping logs to COS. Users need this parameter to control the timezone used in COS file path generation, ensuring logs are organized according to their preferred timezone.

## What Changes

- Add a new optional `time_zone` parameter to the `tencentcloud_cls_cos_shipper` resource
- The parameter maps to `TimeZone` in the CLS cloud API (CreateShipper, ModifyShipper, DescribeShippers)
- The parameter accepts GMT/UTC timezone format strings (e.g., "GMT+08:00", "UTC+08:00")
- The parameter is used in Create, Update, and Read operations

## Capabilities

### New Capabilities
- `cls-cos-shipper-timezone`: Add `time_zone` parameter support to the CLS COS shipper resource for controlling timezone in COS file path generation

### Modified Capabilities

## Impact

- `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go`: Add `time_zone` schema field and wire it into Create, Read, Update methods
- `tencentcloud/services/cls/resource_tc_cls_cos_shipper_test.go`: Add unit test cases for the new parameter
- `tencentcloud/services/cls/resource_tc_cls_cos_shipper.md`: Update example usage documentation
- Cloud API SDK: Uses existing `TimeZone` field in `CreateShipperRequest`, `ModifyShipperRequest`, and `ShipperInfo` structs from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`
