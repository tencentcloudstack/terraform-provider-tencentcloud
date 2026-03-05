# 任务清单: 为 Organization Member Email API 添加重试逻辑

**变更 ID**: `add-retry-logic-org-member-email-api`  
**预计总时间**: 3 小时

---

## 实施任务

### 任务 1: 修改服务层方法添加重试逻辑
**文件**: `tencentcloud/services/tco/service_tencentcloud_organization.go`  
**预计时间**: 30 分钟  
**依赖**: 无

**操作步骤**:
1. 定位到 `DescribeOrganizationOrgMemberEmailById` 方法 (第 460-489 行)
2. 将 API 调用包裹在 `resource.Retry()` 中
3. 使用 `tccommon.ReadRetryTimeout` 作为超时配置
4. 使用 `tccommon.RetryError(e)` 处理错误
5. 确保响应验证逻辑在重试回调函数内部

**预期变更**:
```go
// 当前代码 (第 474-478 行)
response, err := me.client.UseOrganizationClient().DescribeOrganizationMemberEmailBind(request)
if err != nil {
    errRet = err
    return
}

// 修改为
err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
    response, e := me.client.UseOrganizationClient().DescribeOrganizationMemberEmailBind(request)
    if e != nil {
        return tccommon.RetryError(e)
    }
    log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

    if response == nil || response.Response == nil {
        return nil
    }
    if *response.Response.BindId != bindId {
        return nil
    }
    orgMemberEmail = response.Response
    return nil
})
if err != nil {
    errRet = err
    return
}
```

**验证标准**:
- [x] 代码编译通过
- [x] 导入语句正确 (确保 `resource` 包已导入)
- [x] 日志输出移到重试函数内部
- [x] 响应验证逻辑在重试函数内部
- [x] 错误处理保持一致

---

### 任务 2: 代码格式化和静态检查
**预计时间**: 15 分钟  
**依赖**: 任务 1

**操作步骤**:
1. 运行 `make fmt` 格式化代码
2. 运行 `make lint` 检查代码质量
3. 修复任何 lint 警告或错误

**命令**:
```bash
cd /Users/hellertang/git/terraform-provider-tencentcloud
make fmt
make lint
```

**验证标准**:
- [x] `make fmt` 无输出或仅有格式化确认
- [x] `make lint` 无错误和警告
- [x] 代码符合项目编码规范

---

### 任务 3: 运行现有测试
**预计时间**: 30 分钟  
**依赖**: 任务 2

**操作步骤**:
1. 设置测试所需的环境变量
2. 运行 Organization Member Email 资源的验收测试
3. 确保所有测试通过

**命令**:
```bash
# 设置环境变量
export TF_ACC=1
export TENCENTCLOUD_SECRET_ID=<your_secret_id>
export TENCENTCLOUD_SECRET_KEY=<your_secret_key>

# 运行测试
cd /Users/hellertang/git/terraform-provider-tencentcloud
go test -v ./tencentcloud/services/tco -run TestAccTencentCloudOrganizationOrgMemberEmailResource -timeout 120m
```

**验证标准**:
- [x] 测试环境变量已正确设置
- [x] 测试成功运行并通过
- [x] 无新增的测试失败
- [x] 资源的 Create, Read, Update, Delete 操作均正常

---

### 任务 4: 手动功能验证 (可选)
**预计时间**: 30 分钟  
**依赖**: 任务 3

**操作步骤**:
1. 创建测试 Terraform 配置
2. 执行 `terraform apply` 创建资源
3. 执行 `terraform refresh` 验证读取操作
4. 执行 `terraform destroy` 清理资源
5. 观察日志,确认重试机制在需要时生效

**测试配置示例**:
```hcl
resource "tencentcloud_organization_org_member_email" "test" {
  member_uin   = 100012345678
  email        = "test@example.com"
  country_code = "+86"
  phone        = "13800138000"
}
```

**验证标准**:
- [ ] 资源成功创建
- [ ] 资源状态正确读取
- [ ] 重试逻辑在需要时生效 (查看日志中的 RETRY 标记)
- [ ] 资源成功销毁

---

### 任务 5: 更新变更日志
**预计时间**: 10 分钟  
**依赖**: 任务 3

**操作步骤**:
1. 在 `.changelog/` 目录下创建新的变更日志文件
2. 文件名格式: `<PR_NUMBER>.txt` (PR 创建后填写)
3. 描述变更内容

**变更日志内容模板**:
```txt
```enhancement
tco: Add retry logic to DescribeOrganizationMemberEmailBind API call in service layer
```
```

**验证标准**:
- [x] 变更日志文件已创建
- [x] 描述清晰准确
- [x] 分类正确 (enhancement)

---

### 任务 6: 代码审查准备
**预计时间**: 15 分钟  
**依赖**: 任务 5

**操作步骤**:
1. 检查所有修改的文件
2. 确保没有遗漏的 TODO 或注释
3. 验证提交信息清晰
4. 准备 PR 描述

**PR 描述模板**:
```markdown
## 描述
为 `DescribeOrganizationOrgMemberEmailById` 服务层方法添加重试逻辑,提高 API 调用的可靠性。

## 变更类型
- [x] 增强 (Enhancement)
- [ ] Bug 修复
- [ ] 新功能
- [ ] 破坏性变更

## 变更内容
- 使用 `resource.Retry()` 包裹 `DescribeOrganizationMemberEmailBind` API 调用
- 使用 `tccommon.RetryError()` 智能处理可重试错误
- 使用标准的 `tccommon.ReadRetryTimeout` 配置

## 测试
- [x] 现有验收测试通过
- [x] 代码通过 lint 检查
- [x] 手动功能验证 (可选)

## 相关 Issue
Closes #<issue_number> (如果有)
```

**验证标准**:
- [ ] 所有文件变更已审查
- [ ] 提交信息符合规范
- [ ] PR 描述完整准确

---

### 任务 7: 创建 Pull Request
**预计时间**: 10 分钟  
**依赖**: 任务 6

**操作步骤**:
1. 创建功能分支 (如果尚未创建)
2. 提交所有变更
3. 推送到远程仓库
4. 创建 Pull Request
5. 链接相关 Issue (如果有)

**Git 命令**:
```bash
git checkout -b add-retry-logic-org-member-email-api
git add tencentcloud/services/tco/service_tencentcloud_organization.go
git add .changelog/<PR_NUMBER>.txt
git commit -m "tco: Add retry logic to DescribeOrganizationMemberEmailBind API call"
git push origin add-retry-logic-org-member-email-api
```

**验证标准**:
- [ ] 分支已创建并推送
- [ ] Pull Request 已创建
- [ ] PR 描述完整
- [ ] CI/CD 检查开始运行

---

### 任务 8: 响应代码审查
**预计时间**: 30 分钟  
**依赖**: 任务 7

**操作步骤**:
1. 等待代码审查反馈
2. 根据反馈进行修改
3. 推送更新
4. 回复审查意见

**验证标准**:
- [ ] 所有审查意见已解决
- [ ] 代码审查批准
- [ ] CI/CD 检查全部通过

---

### 任务 9: 合并与部署
**预计时间**: 10 分钟  
**依赖**: 任务 8

**操作步骤**:
1. 确认所有检查通过
2. 合并 Pull Request
3. 删除功能分支
4. 验证变更已合并到主分支

**验证标准**:
- [ ] PR 已合并
- [ ] 功能分支已删除
- [ ] 主分支包含变更
- [ ] CI/CD 在主分支上运行成功

---

## 任务依赖关系图

```
任务 1 (修改代码)
    ↓
任务 2 (格式化和 lint)
    ↓
任务 3 (运行测试)
    ↓
任务 4 (手动验证 - 可选)
    ↓
任务 5 (更新变更日志)
    ↓
任务 6 (代码审查准备)
    ↓
任务 7 (创建 PR)
    ↓
任务 8 (响应代码审查)
    ↓
任务 9 (合并与部署)
```

---

## 并行任务机会

- 任务 4 (手动验证) 可以与任务 5 (更新变更日志) 并行
- 但建议按顺序执行以确保质量

---

## 回滚计划

如果变更导致问题:
1. 创建 revert PR 恢复变更
2. 或通过 `git revert` 命令恢复提交
3. 重新评估并修复问题后再次提交

**回滚命令**:
```bash
git revert <commit_hash>
git push origin master
```

---

## 完成标准

所有任务完成且满足以下条件:
- ✅ 代码已合并到主分支
- ✅ 所有测试通过
- ✅ 代码审查批准
- ✅ 变更日志已更新
- ✅ 无遗留的 TODO 或问题
