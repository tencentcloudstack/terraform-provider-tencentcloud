## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/cloudaudit/resource_tc_audit_track.go` 的 `storage` block schema 中新增 `compress` 字段（`schema.TypeInt`, `Optional: true`，Description 注明取值 `1:压缩 2:不压缩`）

## 2. CRUD 函数实现

- [x] 2.1 在 `resourceTencentCloudAuditTrackCreate` 的 `storage` 处理分支（`helper.InterfacesHeadMap(d, "storage")`）中新增 `compress` 到 `request.Storage.Compress` 的映射（`helper.IntUint64`），保持 nil 时不传值
- [x] 2.2 在 `resourceTencentCloudAuditTrackRead` 的 `track.Storage != nil` 分支的 `storageMap` 中新增 `compress` 回填逻辑（先判断 `track.Storage.Compress != nil`）
- [x] 2.3 在 `resourceTencentCloudAuditTrackUpdate` 的 `d.HasChange("storage")` 分支中新增 `compress` 到 `request.Storage.Compress` 的映射（`helper.IntUint64`）

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/cloudaudit/resource_tc_audit_track_test.go` 中使用 gomonkey 对云 API 进行 mock，补充 `compress` 字段的 Create/Read/Update 单元测试用例
- [x] 3.2 使用 `go test -gcflags=all=-l` 运行涉及的单元测试文件，确保通过

## 4. 文档

- [x] 4.1 更新 `tencentcloud/services/cloudaudit/resource_tc_audit_track.md`，在 Example Usage 的 `storage` block 中补充 `compress` 字段示例（由 make doc 最终生成 website/docs）
