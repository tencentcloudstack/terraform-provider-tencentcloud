## 1. Schema Definition

- [x] 1.1 在 `tencentcloud/services/cls/resource_tc_cls_alarm_notice.go` 的 `Schema` map 中新增 `secure_detail_status` 字段：类型 `schema.TypeInt`，`Optional: true`，`Computed: true`，`Description` 说明枚举值 1（关闭，默认）和 2（开启）

## 2. Create Logic Implementation

- [x] 2.1 在 `resourceTencentCloudClsAlarmNoticeCreate` 函数中，使用 `d.GetOkExists("secure_detail_status")` 读取参数值
- [x] 2.2 通过 `helper.IntUint64(v.(int))` 转换后赋值给 `request.SecureDetailStatus`

## 3. Read Logic Implementation

- [x] 3.1 在 `resourceTencentCloudClsAlarmNoticeRead` 函数中，判断 `alarmNotice.SecureDetailStatus != nil`，非 nil 时调用 `d.Set("secure_detail_status", alarmNotice.SecureDetailStatus)` 回写 state

## 4. Update Logic Implementation

- [x] 4.1 在 `resourceTencentCloudClsAlarmNoticeUpdate` 函数中，将 `secure_detail_status` 添加到 `mutableArgs` 数组
- [x] 4.2 在 `needChange` 块内，使用 `d.GetOkExists("secure_detail_status")` 读取并通过 `helper.IntUint64()` 赋值给 `request.SecureDetailStatus`

## 5. Documentation

- [x] 5.1 更新 `tencentcloud/services/cls/resource_tc_cls_alarm_notice.md`，在示例中添加 `secure_detail_status` 参数使用示例

## 6. Testing

- [x] 6.1 在 `tencentcloud/services/cls/resource_tc_cls_alarm_notice_test.go` 中补充 `secure_detail_status` 参数的单元测试用例（使用 gomonkey mock 云 API）
