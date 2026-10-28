## 1. Schema Definition

- [x] 1.1 Add `tags` schema field to `ResourceTencentCloudClbLogTopic()` in `tencentcloud/services/clb/resource_tc_clb_log_topic.go` as a `TypeMap` of string values (`Elem: &schema.Schema{Type: schema.TypeString}`); `Optional`, no `ForceNew`.

## 2. Service Layer

- [x] 2.1 Update `ClbService.CreateTopic()` in `tencentcloud/services/clb/service_tencentcloud_clb.go` to read `tags` from the `params` map (a `map[string]interface{}`) and convert it to `[]*clb.TagInfo` (mapping map key→`TagKey`, map value→`TagValue`), then set `request.Tags`.

## 3. CRUD Implementation

- [x] 3.1 Update `resourceTencentCloudClbInstanceTopicCreate()` to extract `tags` from schema and pass it into the `params` map passed to `ClbService.CreateTopic()`.
- [x] 3.2 Update `resourceTencentCloudClbInstanceTopicRead()` to flatten `res.Tags` (`[]*cls.Tag` with `Key`/`Value`) into the `tags` schema map (`map[string]string`), with nil-safety checks before calling `d.Set("tags", ...)`.
- [x] 3.3 Update `resourceTencentCloudClbInstanceTopicUpdate()` to call `cls.ModifyTopic` with `Tags` (`[]*cls.Tag` mapping map key→`Key`, map value→`Value`) when `d.HasChange("tags")`, using the existing retry pattern (`tccommon.WriteRetryTimeout` / `tccommon.RetryError`).

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/clb/resource_tc_clb_log_topic.md` with a `tags` example in the Example Usage section (generated `website/docs/` files are produced by `make doc` in the finalize phase).

## 5. Tests

- [x] 5.1 Add test cases for the `tags` parameter in `tencentcloud/services/clb/resource_tc_clb_log_topic_test.go` covering create with tags, read-back of tags, and update of tags (using the existing test suite pattern).
