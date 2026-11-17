## 1. Schema Definition

- [x] 1.1 Add `period` schema field (TypeInt, Optional, Computed, not ForceNew, description "Log storage lifecycle in days. Standard storage supports 1-3600; 3640 means permanent retention. Defaults to 30 when unset.") to `ResourceTencentCloudClbLogTopic()` in `tencentcloud/services/clb/resource_tc_clb_log_topic.go`

## 2. Service Layer Changes

- [x] 2.1 In `ClbService.CreateTopic` (`tencentcloud/services/clb/service_tencentcloud_clb.go`), read `params["period"]` and, when present, set `request.Period = common.Uint64Ptr(uint64(period.(int)))` on the `clb.CreateTopicRequest`

## 3. Create Function Changes

- [x] 3.1 In `resourceTencentCloudClbInstanceTopicCreate`, read `period` from schema data (`d.GetOk("period")`) and add it to the `params` map passed to `ClbService.CreateTopic`

## 4. Read Function Changes

- [x] 4.1 In `resourceTencentCloudClbInstanceTopicRead`, after the existing `_ = d.Set(...)` calls, check that `res.Period` is non-nil before setting `period` into state via `_ = d.Set("period", res.Period)`

## 5. Update Function Changes

- [x] 5.1 In `resourceTencentCloudClbInstanceTopicUpdate`, build a single `cls.NewModifyTopicRequest()` when any of `status`, `tags`, or `period` changes, populate `Period` (cast to `*int64` via `helper.Int64`) only when `period` changes, and issue one retried `ModifyTopic` call that updates all changed fields together

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/clb/resource_tc_clb_log_topic.md` with a `period` usage example in the Example Usage block

## 7. Testing

- [x] 7.1 Add mock-based unit tests (gomonkey) in `tencentcloud/services/clb/resource_tc_clb_log_topic_test.go` covering Create with `period`, Read populating `period` from `TopicInfo.Period`, and Update changing `period` via `ModifyTopic`

## 8. Validation

- [x] 8.1 Verify the code compiles (no `go build`/`go vet`; ensure generated code is syntactically correct and all returned errors are handled)
- [x] 8.2 Confirm no lint errors (unused variables handled via `_ =` where a function cannot fail)
