## 1. Schema Definition

- [x] 1.1 Add `zone_id` field (Computed TypeString) to the `tencentcloud_teo_zone` resource schema in `tencentcloud/services/teo/resource_tc_teo_zone.go`

## 2. CRUD Function Updates

- [x] 2.1 Update `resourceTencentCloudTeoZoneRead` to set `zone_id` from `respData.ZoneId` in the API response
- [x] 2.2 Update `resourceTencentCloudTeoZoneDelete` to read `zone_id` from `d.Get("zone_id")` and use it as the `ZoneId` in the `DeleteZone` API request, with fallback to `d.Id()` if `zone_id` is empty

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/resource_tc_teo_zone.md` to add `zone_id` parameter documentation

## 4. Unit Tests

- [x] 4.1 Update `tencentcloud/services/teo/resource_tc_teo_zone_test.go` to add unit test for `zone_id` parameter in read and delete functions

## 5. Verification

- [x] 5.1 Run unit tests with `go test -gcflags=all=-l` to verify the changes compile and tests pass
