# Tasks: CLS Alarm Classifications 实施清单

## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/cls/resource_tc_cls_alarm.go` 的 Schema 中添加 `classifications` 字段
- [x] 1.2 设置字段类型为 `schema.TypeList`,元素类型为 `schema.TypeString`
- [x] 1.3 设置字段属性: `Optional: true`, `Computed: true`
- [x] 1.4 添加字段描述信息

## 2. Create 操作实现

- [x] 2.1 在 `resourceTencentCloudClsAlarmCreate` 函数中添加 `classifications` 参数读取逻辑
- [x] 2.2 使用 `d.GetOk("classifications")` 检查用户是否配置该字段
- [x] 2.3 使用 `helper.InterfacesStringsPoint` 转换为 API 所需的 `[]*string` 类型
- [x] 2.4 将转换后的值设置到 `request.Classifications` 参数

## 3. Read 操作实现

- [x] 3.1 在 `resourceTencentCloudClsAlarmRead` 函数中添加 `Classifications` 字段读取逻辑
- [x] 3.2 从 API 响应中获取 `alarm.Classifications` 字段
- [x] 3.3 检查字段非 nil 后使用 `d.Set("classifications", value)` 写入 state
- [x] 3.4 处理空数组情况,避免不必要的 state 更新

## 4. Update 操作实现

- [x] 4.1 在 `resourceTencentCloudClsAlarmUpdate` 函数中添加 `classifications` 变更检测
- [x] 4.2 使用 `d.HasChange("classifications")` 检测字段是否变更
- [x] 4.3 如果变更,使用 `d.GetOk("classifications")` 读取新值
- [x] 4.4 使用 `helper.InterfacesStringsPoint` 转换并设置到 `request.Classifications` 参数

## 5. 验收测试

- [x] 5.1 在 `tencentcloud/services/cls/resource_tc_cls_alarm_test.go` 中添加测试用例
- [x] 5.2 创建测试函数 `TestAccTencentCloudClsAlarmResource_classifications`
- [x] 5.3 测试场景1: 创建带 `classifications` 的告警
- [x] 5.4 测试场景2: 更新告警的 `classifications` 字段
- [x] 5.5 测试场景3: 创建不带 `classifications` 的告警(验证向后兼容)
- [x] 5.6 添加检查步骤验证 state 中 `classifications` 值正确

## 6. 文档更新

- [x] 6.1 更新 `website/docs/r/cls_alarm.html.markdown` 文档
- [x] 6.2 在 Argument Reference 部分添加 `classifications` 字段说明
- [x] 6.3 说明字段类型(list of strings)、是否必填(Optional)、作用
- [x] 6.4 在示例部分添加使用 `classifications` 的配置示例

## 7. 代码质量检查

- [x] 7.1 运行 `make fmt` 格式化代码
- [x] 7.2 运行 `make lint` 检查代码规范
- [x] 7.3 运行 `go build` 确保编译通过
- [x] 7.4 修复所有 linter 警告和错误

## 8. 测试验证

- [x] 8.1 配置测试环境变量 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY`
- [x] 8.2 运行新增的验收测试: `TF_ACC=1 go test -v ./tencentcloud/services/cls -run TestAccTencentCloudClsAlarmResource_classifications`
- [x] 8.3 验证所有测试用例通过
- [x] 8.4 检查测试覆盖率和场景完整性

**注意**: 验收测试需要真实的腾讯云环境和凭证,建议在 CI/CD 流程中运行或手动验证。
