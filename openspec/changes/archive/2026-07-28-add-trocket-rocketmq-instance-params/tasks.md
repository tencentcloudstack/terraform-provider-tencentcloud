## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/trocket/resource_tc_trocket_rocketmq_instance.go` 的 `ResourceTencentCloudTrocketRocketmqInstance()` Schema map 中新增五个 Optional 字段：`pay_mode` (TypeInt)、`renew_flag` (TypeInt)、`time_span` (TypeInt)、`max_topic_num` (TypeInt)、`zone_ids` (TypeList, Elem TypeInt)，均带 Description
- [x] 1.2 确认新字段均为 Optional、不设 ForceNew、不修改任何既有字段，保证向后兼容

## 2. Create 操作

- [x] 2.1 在 `resourceTencentCloudTrocketRocketmqInstanceCreate` 中，当 `pay_mode` 设置时填充 `request.PayMode` (helper.IntInt64)
- [x] 2.2 当 `renew_flag` 设置时填充 `request.RenewFlag`
- [x] 2.3 当 `time_span` 设置时填充 `request.TimeSpan`
- [x] 2.4 当 `max_topic_num` 设置时填充 `request.MaxTopicNum`
- [x] 2.5 当 `zone_ids` 设置时遍历列表，将 int 转为 `*int64` 填充 `request.ZoneIds` ([]*int64)
- [x] 2.6 确认 Create 调用后检查 response/InstanceId 非空，复用现有等待逻辑

## 3. Read 操作

- [x] 3.1 在 `resourceTencentCloudTrocketRocketmqInstanceRead` 中，当 `rocketmqInstance.PayMode` 非 nil 时做枚举映射（POSTPAID→0, PREPAID→1）并 `d.Set("pay_mode", ...)`
- [x] 3.2 当 `rocketmqInstance.RenewFlag` 非 nil 时 `d.Set("renew_flag", ...)`
- [x] 3.3 当 `rocketmqInstance.ZoneIds` 非 nil 时将 []*int64 转为 []interface{} int 并 `d.Set("zone_ids", ...)`
- [x] 3.4 不回填 `time_span`、`max_topic_num`（无对应响应字段）

## 4. Update 操作

- [x] 4.1 将 `pay_mode`、`renew_flag`、`time_span` 加入 Update 方法的 `immutableArgs` 数组，变更时返回 `fmt.Errorf("argument %s cannot be changed", v)`
- [x] 4.2 当 `max_topic_num` 变更时填充 `request1.MaxTopicNum` (helper.IntInt64) 并置 `needModifyInstance = true`
- [x] 4.3 当 `zone_ids` 变更时遍历列表，将 int 转为 `*string` (helper.IntToStr) 填充 `request1.ZoneIds` ([]*string) 并置 `needModifyInstance = true`
- [x] 4.4 确认 `needModifyInstance` 触发 `ModifyInstance` 调用并复用现有等待逻辑

## 5. 单元测试

- [x] 5.1 在 `tencentcloud/services/trocket/resource_tc_trocket_rocketmq_instance_test.go` 中补充测试用例，覆盖新增 `pay_mode`、`renew_flag`、`time_span`、`max_topic_num`、`zone_ids` 的 create 与 update（max_topic_num/zone_ids 原地更新）场景，使用 Terraform 测试套件（TF_ACC）

## 6. 文档

- [x] 6.1 更新 `tencentcloud/services/trocket/resource_tc_trocket_rocketmq_instance.md`，在示例 HCL 中补充新字段演示（不添加 Argument Reference / Attribute Reference 段落）

## 7. 验证

- [x] 7.1 代码正确性检查：确认 Create 字段在 `CreateInstanceRequest` 中存在、Update 字段在 `ModifyInstanceRequest` 中存在、Read 字段在 `DescribeInstanceResponseParams` 中存在
- [x] 7.2 最终由收尾阶段执行 `gofmt` 格式化与 `make doc` 文档生成（移交 tfpacer-finalize skill 执行）
