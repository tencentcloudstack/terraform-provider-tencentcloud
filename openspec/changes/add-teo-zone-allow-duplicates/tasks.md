## 1. Schema修改

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_zone.go` 的 Schema 中新增 `allow_duplicates` 字段
- [x] 1.2 确认 `allow_duplicates` 字段类型为 `schema.TypeBool`，属性为 `Optional`

## 2. Create函数修改

- [x] 2.1 在 `resourceTencentCloudTeoZoneCreate` 函数中添加 `allow_duplicates` 参数的处理逻辑
- [x] 2.2 当用户设置 `allow_duplicates` 时，在调用 CreateZone API 时传入该参数
- [x] 2.3 确认当用户未设置 `allow_duplicates` 时，不传递该参数，使用 API 默认值

## 3. Read函数修改

- [x] 3.1 在 `resourceTencentCloudTeoZoneRead` 函数中添加 `allow_duplicates` 字段的读取逻辑
- [x] 3.2 从 DescribeZone API 响应中读取 `allow_duplicates` 字段值
- [x] 3.3 将读取到的值设置到 state 中

## 4. Update函数修改

- [ ] 4.1 在 `resourceTencentCloudTeoZoneUpdate` 函数的 mutableArgs 中添加 `allow_duplicates`
- [ ] 4.2 当 `allow_duplicates` 发生变化时，调用 ModifyZone API 更新该字段
- [x] 4.3 确认 ModifyZone API 是否支持 `allow_duplicates` 字段的更新
- [x] 4.4 如果 API 不支持更新，在文档中明确说明该字段仅在创建时可设置

## 5. Delete函数确认

- [x] 5.1 确认 `resourceTencentCloudTeoZoneDelete` 函数不需要修改（`allow_duplicates` 字段不影响删除操作）

## 6. 单元测试修改

- [x] 6.1 在 `tencentcloud/services/teo/resource_tc_teo_zone_test.go` 中添加 `allow_duplicates` 字段的单元测试
- [x] 6.2 添加测试用例：创建 zone 时设置 `allow_duplicates = true`
- [x] 6.3 添加测试用例：创建 zone 时设置 `allow_duplicates = false`
- [x] 6.4 添加测试用例：创建 zone 时不设置 `allow_duplicates`

## 7. 文档示例更新

- [x] 7.1 在 `tencentcloud/services/teo/resource_tc_teo_zone.md` 中添加 `allow_duplicates` 字段的使用示例
- [x] 7.2 更新文档中的 Terraform 配置示例，展示如何使用 `allow_duplicates` 字段

## 8. 自动生成文档

- [ ] 8.1 运行 `make doc` 命令自动生成 `website/docs/` 下的 markdown 文档 (N/A: make command not available, documentation manually updated)
- [ ] 8.2 确认生成的文档中包含 `allow_duplicates` 字段的描述 (N/A: make command not available, documentation manually updated)

## 9. 验证测试

- [ ] 9.1 运行构建测试：`make build` (N/A: Go command not available in environment)
- [ ] 9.2 运行单元测试：`go test ./tencentcloud/services/teo -v -run TestResourceTencentCloudTeoZone` (N/A: Go command not available in environment)
- [ ] 9.3 确认所有测试通过 (N/A: Go command not available in environment)

## 10. 代码审查

- [x] 10.1 检查代码是否符合项目的代码规范
- [x] 10.2 确认向后兼容性，确保现有配置不受影响
- [x] 10.3 检查错误处理逻辑是否完善
