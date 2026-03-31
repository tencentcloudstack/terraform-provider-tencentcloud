## 1. 代码修改

### 1.1 资源 Schema 修改

- [x] 1.1.1 在 tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go 文件中添加 TotalCount 字段到 Schema
- [x] 1.1.2 将 TotalCount 字段定义为 schema.TypeInt 类型
- [x] 1.1.3 设置 TotalCount 字段的 Optional 属性为 true
- [x] 1.1.4 设置 TotalCount 字段的 Computed 属性为 true
- [x] 1.1.5 添加 TotalCount 字段的描述信息

### 1.2 资源 Read 函数修改

- [x] 1.2.1 在 Read 函数中从 DescribeL7AccRules API 响应中读取 TotalCount 字段
- [x] 1.2.2 添加 TotalCount 字段的空值检查逻辑
- [x] 1.2.3 当 TotalCount 为 null 时，设置默认值为 0
- [x] 1.2.4 使用 d.Set("total_count", value) 将 TotalCount 值设置到资源状态中

### 1.3 单元测试修改

- [x] 1.3.1 在 tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_test.go 文件中添加 TotalCount 字段的测试用例
- [x] 1.3.2 添加测试验证 TotalCount 字段的类型为整型
- [x] 1.3.3 添加测试验证 TotalCount 字段的值正确
- [x] 1.3.4 添加测试验证 TotalCount 字段为 null 时的默认值处理

## 2. 验证任务

### 2.1 单元测试

- [ ] 2.1.1 运行 tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_test.go 单元测试
- [ ] 2.1.2 验证所有单元测试通过
- [ ] 2.1.3 验证新增的 TotalCount 字段测试通过

### 2.2 代码质量检查

- [ ] 2.2.1 运行 go fmt 格式化代码
- [ ] 2.2.2 运行 golint 检查代码规范
- [ ] 2.2.3 运行 go vet 检查代码问题
- [ ] 2.2.4 修复所有代码检查警告和错误

### 2.3 构建验证

- [ ] 2.3.1 运行 go build 编译 Terraform Provider
- [ ] 2.3.2 验证编译成功，无错误

## 3. 文档更新

### 3.1 示例文件更新

- [ ] 3.1.1 检查是否需要更新 tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.md 示例文件
- [ ] 3.1.2 如需要，在示例文件中添加 TotalCount 字段的说明
- [ ] 3.1.3 在示例文件中明确说明 TotalCount 为只读字段

### 3.2 文档生成

- [ ] 3.2.1 运行 make doc 命令生成 website/docs/ 下的 markdown 文档
- [ ] 3.2.2 验证生成的文档中包含 TotalCount 字段
- [ ] 3.2.3 验证文档中 TotalCount 字段的描述正确
