## 1. Schema Definition

- [x] 1.1 Add `time_zone` field to the resource schema in `tencentcloud/services/cls/resource_tc_cls_cos_shipper.go` as an Optional, Computed string field with appropriate description

## 2. CRUD Implementation

- [x] 2.1 Add `time_zone` handling in the Create method (`resourceTencentCloudClsCosShipperCreate`): read from ResourceData and set `request.TimeZone`
- [x] 2.2 Add `time_zone` handling in the Read method (`resourceTencentCloudClsCosShipperRead`): read from `ShipperInfo.TimeZone` with nil check and set to state
- [x] 2.3 Add `time_zone` handling in the Update method (`resourceTencentCloudClsCosShipperUpdate`): handle `d.HasChange("time_zone")` and set `request.TimeZone`

## 3. Unit Tests

- [x] 3.1 Add unit test cases in `tencentcloud/services/cls/resource_tc_cls_cos_shipper_test.go` using gomonkey to mock cloud API calls, verifying the `time_zone` parameter is correctly handled in Create, Read, and Update operations

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/cls/resource_tc_cls_cos_shipper.md` to include `time_zone` parameter in the example usage
