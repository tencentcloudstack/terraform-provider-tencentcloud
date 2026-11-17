## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/mongodb/resource_tc_mongodb_instance.go` 的 `mongodbInstanceInfo` schema map 中新增 `cpu` 字段：`Type: schema.TypeInt`、`Optional: true`、`Computed: true`，描述说明单位为 C 且支持的规格可通过 `DescribeSpecInfo` 接口获取

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` 的 `UpgradeInstance` 方法中新增 `params["cpu"]` 处理：当���在时使用 `helper.Int64(v.(int))` 构造 `*int64` 并赋值给 `request.Cpu`（匹配 `ModifyDBInstanceSpecRequest.Cpu` 的 `*int64` 类型）

## 3. CRUD 函数实现

- [x] 3.1 在 `resourceTencentCloudMongodbInstanceUpdate` 中将变配触发条件从 `d.HasChange("memory") || d.HasChange("volume") || d.HasChange("node_num")` 扩展为包含 `d.HasChange("cpu")`
- [x] 3.2 在 Update 变配分支中读取 `cpu` 值并放入 `params["cpu"]`（使用 `d.GetOk`/`d.Get`，与现有 `node_num`/`in_maintenance` 传参模式一致）
- [x] 3.3 在 `resourceTencentCloudMongodbInstanceRead` 中，当 `instance.CpuNum != nil` 时执行 `_ = d.Set("cpu", int(*instance.CpuNum))`；不将 `CpuNum` 加入 `CheckNil` 强校验列表

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/mongodb/resource_tc_mongodb_instance_test.go` 中新增基于 gomonkey mock 的单元测试：mock `ModifyDBInstanceSpec` 与 `DescribeDBInstanceDeal`，验证 `cpu` 变更触发变配且 `request.Cpu` 被正确设置
- [x] 4.2 新增 Read 回填单元测试：mock `DescribeDBInstances`（`DescribeInstanceById`）返回 `CpuNum` 非空，验证 `d.Set("cpu", ...)` 被正确调用
- [x] 4.3 新增 Read 处理 `CpuNum` 为 nil 的单元测试，验证 Read 不报错且不设置 `cpu`
- [x] 4.4 使用 `go test -gcflags=all=-l` 运行所涉测试文件确保通过

## 5. 文档

- [x] 5.1 修改 `tencentcloud/services/mongodb/resource_tc_mongodb_instance.md`，补充 `cpu` 参数说明（一句话描述带上 MongoDB 产品名），在示例用法中补充 `cpu` 变更示例；不手动添加 `Argument Reference`/`Attribute Reference` 部分

## 6. 代码正确性检查

- [x] 6.1 核对 `cpu` 参数仅在 `ModifyDBInstanceSpec`（Update 路径）中作为入参，确认 `ModifyDBInstanceSpecRequest.Cpu` 字段存在且类型为 `*int64`
- [x] 6.2 核对 Read 回填来源为 `DescribeDBInstances` 返回的 `InstanceDetail.CpuNum`（`*uint64`），确认字段存在
- [x] 6.3 核对 `cpu` 未被误加入 Create（`CreateDBInstanceHour`/`CreateDBInstance` 使用的是 `CpuCore` 而非 `Cpu`，本次不扩展）

## 7. 收尾阶段（由 tfpacer-finalize 统一执行）

- [ ] 7.1 执行 `gofmt` 格式化本次变更的 Go 文件
- [ ] 7.2 执行 `make doc` 生成 `website/docs/r/mongodb_instance.html.markdown` 等文档
- [ ] 7.3 在 `.changelog/` 目录下创建 changelog 文件（enhancement 类型）
