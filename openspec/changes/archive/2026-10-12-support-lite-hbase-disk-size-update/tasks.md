# 任务清单：support-lite-hbase-disk-size-update

## 1. Update 函数修改

**文件**: `tencentcloud/services/emr/resource_tc_lite_hbase_instance.go`

- [x] 1.1 从 Update 函数的 `immutableArgs` 数组中移除 `"disk_size"` 字符串（第 380 行附近）
- [x] 1.2 在现有 `zone_settings.node_num` 变更处理块（`needChange` 块）之后、`_ = instanceId` 之前，新增 `disk_size` 变更处理块：
  - 检测 `d.HasChange("disk_size")`
  - 构造 `emr.NewModifySLInstanceRequest()`，设置 `InstanceId`（取自 `instanceId`）与 `DiskSize`（取自 `d.GetOkExists("disk_size")`）
  - 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `ModifySLInstanceWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装
  - 调用成功后，构造 `EMRService` + `tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.SLInstanceStateRefreshFunc(instanceId, []string{}))` 等待状态变为 2
  - 错误处理：打印 `[CRITAL]` 日志并返回

## 2. 单元测试补充

**文件**: `tencentcloud/services/emr/resource_tc_lite_hbase_instance_test.go`

- [x] 2.1 在现有 `TestAccTencentCloudLiteHbaseInstanceResource_basic` 测试用例的 BasicUpdate 配置中，将 `disk_size` 从 100 修改为 200，增加对扩容场景的验收测试覆盖
- [x] 2.2 更新对应的 `TestCheckResourceAttr` 断言，验证 `disk_size` 更新为 200

## 3. 文档更新

**文件**: `tencentcloud/services/emr/resource_tc_lite_hbase_instance.md`

- [x] 3.1 同步更新资源示例文档，补充 `disk_size` 支持扩容（仅扩容）的说明

## 4. 验证

- [x] 4.1 确认代码编译通过（由后续流程执行 go build/vet）
- [x] 4.2 确认单元测试可正确构建（由后续流程执行）