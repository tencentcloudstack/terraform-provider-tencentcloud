## 1. 资源代码实现

- [x] 1.1 创建 `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_operation.go`，实现 schema 定义（blueprint_id: Required/ForceNew/TypeString, account_ids: Required/ForceNew/TypeList(of TypeString)）和 CRUD 函数（Create 调用 ShareBlueprintAcrossAccounts API，Read/Delete 为空实现，无 Update）
- [x] 1.2 在 `tencentcloud/provider.go` 中注册资源 `tencentcloud_lighthouse_share_blueprint_across_account`

## 2. 资源文档

- [x] 2.1 创建 `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_operation.md` 资源样例文档，包含一句话描述和 Example Usage

## 3. 单元测试

- [x] 3.1 创建 `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_operation_test.go`，使用 gomonkey mock `ShareBlueprintAcrossAccounts` API，测试 Create 函数的成功和失败场景

## 4. 验证

- [x] 4.1 使用 `go test -gcflags=all=-l` 运行单元测试，确保测试通过
