## 1. Schema Definition

- [x] 1.1 修改 `alarm_notice_ids` 字段:将 `Required: true` 改为 `Optional: true`,添加 `ExactlyOneOf: []string{"monitor_notice"}`
- [x] 1.2 添加 `monitor_notice` schema 块,类型为 `schema.TypeList`,`MaxItems: 1`,`Optional: true`,`ExactlyOneOf: []string{"alarm_notice_ids"}`
- [x] 1.3 在 `monitor_notice` 块内定义 `notices` 字段:类型 `schema.TypeList`,包含 `notice_id`(Required, TypeString)、`content_tmpl_id`(Optional, TypeString)、`alarm_levels`(Required, TypeList of TypeInt)

## 2. Create Logic Implementation

- [x] 2.1 在 `resourceTencentCloudClsAlarmCreate` 函数中添加 `monitor_notice` 参数处理逻辑
- [x] 2.2 使用 `d.GetOk("monitor_notice")` 检查是否配置了 `monitor_notice`
- [x] 2.3 遍历 `monitor_notice.notices` 列表,构建 `cls.MonitorNoticeRule` 对象数组
- [x] 2.4 将 `MonitorNoticeRule` 数组赋值给 `cls.MonitorNotice.Notices`,并将 `MonitorNotice` 赋值给 `request.MonitorNotice`

## 3. Read Logic Implementation

- [x] 3.1 在 `resourceTencentCloudClsAlarmRead` 函数中添加 `MonitorNotice` 读取逻辑
- [x] 3.2 检查 API 响应 `alarm.MonitorNotice` 是否为 nil,非 nil 时才进行字段读取
- [x] 3.3 将 `alarm.MonitorNotice.Notices` 转换为 Terraform state 格式(map 数组)
- [x] 3.4 使用 `d.Set("monitor_notice", ...)` 将数据写入 state

## 4. Update Logic Implementation

- [x] 4.1 在 `resourceTencentCloudClsAlarmUpdate` 函数中添加 `monitor_notice` 参数处理逻辑(与 Create 逻辑类似)
- [x] 4.2 检测 `monitor_notice` 字段的变更(`d.HasChange("monitor_notice")`)
- [x] 4.3 构建 `ModifyAlarmRequest.MonitorNotice` 并调用更新接口

## 5. Testing

- [x] 5.1 在 `resource_tc_cls_alarm_test.go` 中添加新测试函数 `TestAccTencentCloudClsAlarmResource_monitorNotice`
- [x] 5.2 编写测试配置 `testAccClsAlarmWithMonitorNotice`,使用 `monitor_notice` 块创建告警
- [x] 5.3 编写测试配置 `testAccClsAlarmUpdateMonitorNotice`,测试 `monitor_notice` 的更新场景
- [x] 5.4 添加测试步骤验证 `monitor_notice.notices.#`、`monitor_notice.notices.0.notice_id` 等字段值
- [x] 5.5 添加测试场景:从 `alarm_notice_ids` 切换到 `monitor_notice`(测试互斥约束和更新逻辑)
- [ ] 5.6 本地运行验收测试 `TF_ACC=1 go test -v -run TestAccTencentCloudClsAlarmResource_monitorNotice`

## 6. Documentation

- [ ] 6.1 更新 `examples/tencentcloud-cls-alarm/main.tf` 示例文件,添加 `monitor_notice` 使用示例
- [ ] 6.2 在示例中注释说明 `alarm_notice_ids` 和 `monitor_notice` 的互斥关系
- [ ] 6.3 运行 `make doc` 生成 `website/docs/r/cls_alarm.html.markdown` 文档
- [ ] 6.4 检查生成的文档是否正确包含 `monitor_notice` 块的描述和示例

## 7. Validation and Cleanup

- [x] 7.1 运行 `make build` 确保代码编译通过
- [ ] 7.2 运行 `make lint` 检查代码规范
- [ ] 7.3 运行完整的验收测试套件 `TF_ACC=1 go test ./tencentcloud/services/cls/... -v`
- [ ] 7.4 验证向后兼容性:确认现有使用 `alarm_notice_ids` 的测试用例仍然通过
- [ ] 7.5 手动测试互斥约束:尝试同时配置两个字段,验证 Terraform plan 是否报错
