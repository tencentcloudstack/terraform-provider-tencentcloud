## 1. 准备和验证

- [x] 1.1 检查 tencentcloud_teo_origin_group 资源文件位置
- [x] 1.2 确认 tencentcloud-sdk-go 是否包含 HostHeader 参数支持
- [x] 1.3 查看现有的 CreateOriginGroup API 调用实现，了解请求结构

## 2. Schema 修改

- [x] 2.1 在 tencentcloud/services/teo/resource_tc_teo_origin_group.go 中添加 host_header 字段定义到 Schema
- [x] 2.2 验证 host_header 字段属性：Type: String, Optional: true, Computed: false, ForceNew: false
- [x] 2.3 确保字段命名符合 Terraform Provider 规范（snake_case）

## 3. CRUD 函数实现

- [x] 3.1 在 Create 函数中添加逻辑：如果 host_header 非空，将其添加到 CreateOriginGroup API 请求参数
- [x] 3.2 在 Read 函数中添加逻辑：从 API 响应中读取 HostHeader 值并设置到 state
- [x] 3.3 在 Update 函数中添加逻辑：传递更新后的 host_header 值到 CreateOriginGroup API
- [x] 3.4 确认 Delete 函数无需修改（host_header 不影响删除）

## 4. 测试实现

- [x] 4.1 在 tencentcloud/services/teo/resource_tc_teo_origin_group_test.go 中添加 host_header 字段的单元测试用例
- [x] 4.2 添加测试场景：创建资源时设置 host_header
- [x] 4.3 添加测试场景：创建资源时不设置 host_header
- [x] 4.4 添加测试场景：更新资源的 host_header 字段
- [x] 4.5 添加测试场景：读取资源的 host_header 字段
- [x] 4.6 运行单元测试，确保所有新测试通过

## 5. 文档更新

- [x] 5.1 更新 tencentcloud/services/teo/resource_tc_teo_origin_group.md 示例文件，添加 host_header 使用示例
- [x] 5.2 运行 make doc 命令生成 website/docs/r/teo_origin_group.html.md 文档
- [x] 5.3 验证生成的文档正确显示 host_header 字段及其说明

## 6. 构建和验证

- [x] 6.1 运行 go build 验证代码编译无错误
- [x] 6.2 运行 go vet 检查代码规范
- [ ] 6.3 运行 TF_ACC=1 对应的验收测试（需要真实凭证 - 由用户手动执行）
- [ ] 6.4 手动测试：创建带 host_header 的资源，验证 API 调用和状态（由用户手动执行）
- [ ] 6.5 手动测试：更新 host_header 字段，验证状态更新（由用户手动执行）
- [ ] 6.6 手动测试：删除带 host_header 的资源，验证删除成功（由用户手动执行）

**说明：** 任务 6.3-6.6 需要真实的 Tencent Cloud 凭证，需要用户在具备真实环境的情况下手动执行验收测试。