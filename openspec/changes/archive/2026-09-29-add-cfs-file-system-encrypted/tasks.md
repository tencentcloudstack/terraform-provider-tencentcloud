## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/cfs/resource_tc_cfs_file_system.go` 的 Schema map 中新增 `encrypted` 字段：`TypeBool, Optional, ForceNew: true`，描述"是否加密文件系统，true 表示加密，false 表示不加密。修改此参数会触发资源重建。"
- [x] 1.2 在 Update 函数的 `immutableArgs` 数组中添加 `"encrypted"`（与 ForceNew 双重保险，确保修改时不可原地更新）

## 2. Create 函数修改

- [x] 2.1 在 `resourceTencentCloudCfsFileSystemCreate` 中读取 `encrypted` 参数并设置到 `CreateCfsFileSystemRequest.Encrypted`：`request.Encrypted = helper.Bool(d.Get("encrypted").(bool))`

## 3. Read 函数修改

- [x] 3.1 在 `resourceTencentCloudCfsFileSystemRead` 中，当 `fileSystem.Encrypted != nil` 时调用 `_ = d.Set("encrypted", *fileSystem.Encrypted)` 同步加密状态到 state

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/cfs/resource_tc_cfs_file_system_test.go` 中使用 mock（gomonkey）方式补充单元测试用例，覆盖 Create 设置 Encrypted 入参、Read 同步 Encrypted 出参的业务逻辑

## 5. 文档更新

- [x] 5.1 更新 `tencentcloud/services/cfs/resource_tc_cfs_file_system.md` 源文档，添加 `encrypted` 参数说明
- [x] 5.2 通过收尾阶段的 `make doc` 命令生成 `website/docs/r/cfs_file_system.html.markdown`（禁止手动编写）

## 依赖关系
- 任务 2、3 依赖于任务 1 完成
- 任务 4 依赖于任务 1、2、3 完成
- 任务 5 依赖于所有代码修改完成