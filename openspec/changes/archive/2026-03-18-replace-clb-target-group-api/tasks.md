# 任务清单

## 1. 修改 `DescribeTargetGroups` 函数

**文件**: `tencentcloud/services/clb/service_tencentcloud_clb.go`  
**位置**: line 1529-1569

- [x] 替换请求构造: `clb.NewDescribeTargetGroupsRequest()` → `clb.NewDescribeTargetGroupListRequest()`
- [x] 替换 API 调用: `DescribeTargetGroups(request)` → `DescribeTargetGroupList(request)`
- [x] 验证请求参数设置逻辑保持不变
  - [x] TargetGroupIds 设置
  - [x] Filters 设置
  - [x] Offset/Limit 分页逻辑
- [x] 验证响应处理逻辑保持不变
  - [x] 响应结构解析
  - [x] 错误处理
  - [x] 分页循环

**验证**:
```bash
# 编译检查
go build -o /dev/null ./tencentcloud/services/clb/
```

---

## 2. 修改 `DescribeClbTargetGroupAttachmentsById` 函数

**文件**: `tencentcloud/services/clb/service_tencentcloud_clb.go`  
**位置**: line 2626-2685

- [x] 替换请求构造: `clb.NewDescribeTargetGroupsRequest()` → `clb.NewDescribeTargetGroupListRequest()`
- [x] 替换 API 调用: `DescribeTargetGroups(request)` → `DescribeTargetGroupList(request)`
- [x] 验证 TargetGroupIds 参数设置
- [x] 验证响应处理逻辑

**验证**:
```bash
# 编译检查
go build -o /dev/null ./tencentcloud/services/clb/
```

---

## 3. 编译和 Linter 验证

- [x] 运行 `go build` 确保代码编译通过
- [x] 运行 linter 检查没有新增警告或错误
- [x] 确认代码风格符合项目规范

**命令**:
```bash
cd /Users/hellertang/git/terraform-provider-tencentcloud
go build -o /dev/null ./tencentcloud/services/clb/
make lint
```

---

## 4. 单元测试验证

运行相关的单元测试,确保功能正常:

- [x] 验收测试需要真实云环境配置,已跳过
- [x] 代码逻辑审查完成,参数和响应处理保持一致
- [x] 编译通过,无新增错误

**说明**: 验收测试需要配置真实的腾讯云环境,建议在提交 PR 后由 CI/CD 自动运行。代码审查确认所有参数和响应处理逻辑与旧 API 完全一致。

**命令**:
```bash
cd /Users/hellertang/git/terraform-provider-tencentcloud
TF_ACC=1 go test -v ./tencentcloud/services/clb -run TestAccTencentCloudClbTargetGroup
```

---

## 5. 功能回归测试 (可选但推荐)

手动验证完整的资源生命周期:

### 测试目标组资源

**测试配置**:
```hcl
resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test-replace-api"
  vpc_id            = "<your-vpc-id>"
  port              = 80
}
```

- [ ] 创建目标组: `terraform apply`
- [ ] 读取目标组: `terraform refresh`
- [ ] 查询目标组: data source 查询
- [ ] 导入目标组: `terraform import tencentcloud_clb_target_group.test <id>`
- [ ] 更新目标组: 修改 `target_group_name` 后 `terraform apply`
- [ ] 删除目标组: `terraform destroy`

### 测试数据源

**测试配置**:
```hcl
data "tencentcloud_clb_target_groups" "test" {
  target_group_id = "<existing-target-group-id>"
}

output "groups" {
  value = data.tencentcloud_clb_target_groups.test.list
}
```

- [ ] 通过 ID 查询
- [ ] 通过 VPC ID 过滤查询
- [ ] 通过名称过滤查询
- [ ] 验证输出数据完整性

### 测试目标组绑定

**测试配置**:
```hcl
resource "tencentcloud_clb_target_group_attachment" "test" {
  target_group_id = tencentcloud_clb_target_group.test.id
  load_balancer_id = "<your-clb-id>"
  listener_id     = "<your-listener-id>"
}
```

- [ ] 创建绑定关系
- [ ] 读取绑定状态
- [ ] 删除绑定关系

---

## 6. 日志和错误处理验证

- [ ] 验证日志输出正确显示新的 API Action 名称
- [ ] 验证错误处理逻辑正常工作
- [ ] 验证重试机制正常

**验证方法**:
```bash
# 启用 Terraform 详细日志
export TF_LOG=DEBUG
terraform apply
# 检查日志中是否出现 DescribeTargetGroupList
```

---

## 7. 代码审查检查项

- [ ] 代码改动范围最小化
- [ ] 无破坏性变更
- [ ] 注释和文档保持准确
- [ ] 错误处理保持一致
- [ ] 日志记录保持一致
- [ ] 无硬编码值
- [ ] 遵循项目代码规范

---

## 8. 文档更新 (如需要)

- [ ] 检查是否需要更新资源文档
- [ ] 检查是否需要更新数据源文档
- [ ] 检查是否需要更新 CHANGELOG

**说明**: 由于这是内部 API 替换,对用户透明,通常不需要更新用户文档。但如果有内部开发文档,应该记录此变更。

---

## 9. Changelog 条目 (可选)

**文件**: `.changelog/<PR_NUMBER>.txt`

- [ ] 创建 changelog 条目记录此改进

**格式**:
```
```enhancement
clb: Replace `DescribeTargetGroups` API with `DescribeTargetGroupList` for better performance in target group queries
```
```

**说明**: Changelog 应在 PR 创建时添加,使用实际的 PR 编号。

---

## 总结

- **核心修改**: 2 个函数
- **实际工作量**: 低 (已完成核心修改)
- **风险等级**: 低
- **破坏性变更**: 无
- **测试需求**: 现有测试用例可复用
- **状态**: ✅ 核心实施完成

### 完成标准

1. ✅ 所有代码修改完成
2. ✅ 编译通过,无 linter 错误
3. ⏭️ 单元测试需要云环境(留待 CI/CD)
4. ⏭️ (可选) 手动功能测试验证通过
5. ⏭️ 代码审查通过
6. ⏭️ 文档更新(如需要)

### 实施顺序

1. 任务 1: 修改 `DescribeTargetGroups` 函数
2. 任务 2: 修改 `DescribeClbTargetGroupAttachmentsById` 函数
3. 任务 3: 编译和 Linter 验证
4. 任务 4: 运行单元测试
5. 任务 5-9: 可选的额外验证和文档工作

---

## 依赖关系

- 任务 1 和任务 2 可以并行执行
- 任务 3 依赖任务 1、2 完成
- 任务 4 依赖任务 3 完成
- 任务 5-9 可以并行执行,依赖任务 4 完成

---

## 注意事项

1. ⚠️ **API 调用监控**: 实施后应监控 API 调用是否正常
2. ⚠️ **日志变更**: 日志中会显示新的 API Action 名称
3. ⚠️ **回滚准备**: 保留回滚方案,以防出现问题
4. ✅ **向后兼容**: 不影响现有 Terraform 配置
5. ✅ **状态兼容**: 不影响 Terraform state 格式
