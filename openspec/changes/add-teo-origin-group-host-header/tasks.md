## 1. 调研和准备

- [x] 1.1 查阅腾讯云 CreateOriginGroup API 文档，确认 HostHeader 字段的准确类型（string、*string 或其他）
- [x] 1.2 确认 ModifyOriginGroup API 是否支持更新 HostHeader 字段
- [x] 1.3 确认 DescribeOriginGroup 或 ReadOriginGroup API 响应中 HostHeader 字段的返回格式
- [x] 1.4 检查现有 tencentcloud_teo_origin_group 资源代码结构（resource_tencentcloud_teo_origin_group.go）
- [x] 1.5 确认 host_header 参数是否需要 Computed（取决于 API 是否返回默认值）

## 2. Schema 定义更新

- [x] 2.1 在 tencentcloud_teo_origin_group 资源 Schema 中添加 `host_header` 参数（TypeString，Optional，Computed）
- [x] 2.2 为 host_header 参数添加合适的描述信息（"Custom Host header to use when making requests to the origin"）
- [x] 2.3 确保参数位于 Schema 根级别（与资源结构一致）

## 3. Create 函数实现

- [x] 3.1 在 Create 函数中获取 host_header 参数值（`d.Get("host_header")`）
- [x] 3.2 将 host_header 值映射到 CreateOriginGroup API 请求的 HostHeader 字段
- [x] 3.3 处理空字符串情况（根据 API 要求决定是否传递或省略该字段）
- [x] 3.4 在创建成功后设置 state 中的 host_header 值（`d.Set("host_header", params.HostHeader)`）

## 4. Read 函数实现

- [x] 4.1 在 Read 函数中从 API 响应读取 HostHeader 字段
- [x] 4.2 将 API 响应中的 HostHeader 值设置到 state（`d.Set("host_header", resp.HostHeader)`）
- [x] 4.3 处理 API 未返回 HostHeader 的情况（不设置 state 或设置为 nil）

## 5. Update 函数实现

- [x] 5.1 在 Update 函数中检测 host_header 参数变化（对比 old 和 new 值）
- [x] 5.2 当 host_header 发生变化时，调用 ModifyOriginGroup API 更新该参数
- [x] 5.3 处理新增 host_header（从无到有）的场景
- [x] 5.4 处理删除 host_header（从有到无）的场景
- [x] 5.5 处理修改 host_header（值变化）的场景
- [ ] 5.6 更新 state 中的 host_header 值

## 6. 文档和示例更新

- [ ] 6.1 更新 resource_tc_teo_origin_group.md 示例文件，添加 host_header 参数使用示例
- [ ] 6.2 在示例中演示典型用例（如设置自定义 Host 头）
- [ ] 6.3 确保示例中的参数类型和值正确

## 7. 测试用例实现

- [ ] 7.1 新增测试用例：创建包含 host_header 的 origin_group 资源
- [ ] 7.2 新增测试用例：更新 origin_group 的 host_header 值
- [ ] 7.3 新增测试用例：为现有 origin_group 添加 host_header
- [ ] 7.4 新增测试用例：从 origin_group 删除 host_header
- [ ] 7.5 新增测试用例：创建不包含 host_header 的 origin_group（向后兼容性）
- [ ] 7.6 确保所有现有测试用例仍通过

## 8. 代码验证

- [ ] 8.1 运行 go build 确保代码编译无误
- [ ] 8.2 运行 go vet 或 golangci-lint 检查代码质量
- [ ] 8.3 运行相关单元测试（resource_tencentcloud_teo_origin_group_test.go）
- [ ] 8.4 运行完整测试套件确保现有功能不受影响
- [ ] 8.5 运行文档生成命令（make doc）验证文档更新

## 9. 最终检查

- [ ] 9.1 确认所有代码变更符合 Terraform Provider 最佳实践
- [ ] 9.2 确认 host_header 参数在所有 CRUD 场景下正确工作
- [ ] 9.3 确认向后兼容性，现有资源配置无需修改
- [ ] 9.4 确认文档和示例完整且准确
- [ ] 9.5 检查代码注释和错误处理是否完善
