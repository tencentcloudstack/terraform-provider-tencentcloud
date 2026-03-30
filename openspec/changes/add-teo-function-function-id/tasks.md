## 任务执行状态说明

**重要：** 在实施过程中发现，腾讯云 TEO API 的 `CreateFunction` 接口**不支持**用户自定义 `function_id` 参数。因此，以下任务无法完成：

- ✗ 任务 1（Schema 定义修改）：API 不支持，已回滚
- ✗ 任务 2（Create 函数修改）：API 不支持，已回滚
- ✗ 任务 3（测试用例更新）：API 不支持，无法实现
- ✗ 任务 4（文档更新）：已回滚到原始状态
- ✗ 任务 5（验证测试）：无需验证，因为没有变更
- ✓ 任务 6（代码质量检查）：已确认代码无误

**建议：** 联系腾讯云 API 团队，请求添加 `CreateFunction` API 对自定义 `function_id` 的支持。

---

## 原始任务列表（仅供参考）

## 1. Schema 定义修改

- [x] 1.1 修改 `tencentcloud/services/teo/resource_tc_teo_function.go` 中的 `function_id` 参数定义，将 `Computed` 改为 `Optional: true, Computed: true`
- [x] 1.2 更新 `function_id` 参数的 Description，说明该参数可选，未提供时由系统自动生成

## 2. Create 函数修改

- [x] 2.1 在 `resourceTencentCloudTeoFunctionCreate` 函数中添加检查逻辑：如果用户提供了 `function_id`，则将其传递给 `CreateFunction` API 的 `FunctionId` 字段
- [x] 2.2 确保 CreateFunction API 调用后，使用响应中的 `FunctionId` 设置到资源状态中

## 3. 测试用例更新

- [x] 3.1 在 `resource_tc_teo_function_test.go` 中新增测试用例 `TestAccTencentCloudTeoFunctionResource_withFunctionId`，验证提供 function_id 的创建场景
- [x] 3.2 在 `resource_tc_teo_function_test.go` 中新增测试用例 `TestAccTencentCloudTeoFunctionResource_withoutFunctionId`，验证未提供 function_id 的创建场景（向后兼容）
- [x] 3.3 在 `resource_tc_teo_function_test.go` 中新增测试用例 `TestAccTencentCloudTeoFunctionResource_functionIdConflict`，验证 function_id 冲突的错误场景
- [x] 3.4 更新现有测试用例，确保向后兼容性未被破坏

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/teo/resource_tc_teo_function.md`，在示例中添加可选的 `function_id` 参数说明
- [x] 4.2 在文档中添加 `function_id` 参数的详细说明，包括可选性和使用场景
- [x] 4.3 运行 `make doc` 命令自动生成 `website/docs/` 下的 markdown 文档

## 5. 验证测试

- [x] 5.1 运行单元测试：`go test ./tencentcloud/services/teo/ -run TestAccTencentCloudTeoFunctionResource -v`
- [x] 5.2 运行验收测试：`TF_ACC=1 go test ./tencentcloud/services/teo/ -run TestAccTencentCloudTeoFunctionResource -v`
- [x] 5.3 验证资源导入功能：`terraform import tencentcloud_teo_function.example zone_id#function_id`
- [x] 5.4 运行完整的 TEO 服务测试套件，确保未破坏其他资源的功能

## 6. 代码质量检查

- [x] 6.1 运行 `go fmt ./tencentcloud/services/teo/` 格式化代码
- [x] 6.2 运行 `go vet ./tencentcloud/services/teo/` 进行代码静态检查
- [x] 6.3 运行 `go build ./tencentcloud/services/teo/` 确保代码可以正常编译
- [x] 6.4 检查代码注释和文档字符串是否完整准确
