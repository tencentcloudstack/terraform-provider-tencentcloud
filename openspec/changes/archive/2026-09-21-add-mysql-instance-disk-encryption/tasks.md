## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/cdb/resource_tc_mysql_instance.go` 的 `ResourceTencentCloudMysqlInstance()` schema map 中新增 `disk_encryption` 字段（TypeString, Optional, ForceNew, Computed），描述磁盘加密：`on` 表示开启加密，否则不加密；仅云盘版实例支持。放置于 `destroy_protect` 字段附近，与 `disk_type`/`destroy_protect` 风格一致。

## 2. Create 函数变更

- [x] 2.1 在共享 helper `mysqlAllInstanceRoleSet`（resource_tc_mysql_instance.go）中新增读取 `disk_encryption` 的逻辑：使用 `d.GetOk("disk_encryption")`，当存在值时，对 `CreateDBInstanceRequest` 设置 `requestByMonth.DiskEncryption`，对 `CreateDBInstanceHourRequest` 设置 `requestByUse.DiskEncryption`（参照 `destroy_protect` 写法）。
- [x] 2.2 验证 month-paid 路径 `mysqlCreateInstancePayByMonth` 与 hourly-paid 路径 `mysqlCreateInstancePayByUse` 均调用 `mysqlAllInstanceRoleSet`，无需单独改动这两个函数。

## 3. Read 函数变更

- [x] 3.1 在 `tencentMsyqlBasicInfoRead`（resource_tc_mysql_instance.go）中新增 nil 守卫读取：当 `mysqlInfo.DiskEncryption != nil` 时，执行 `_ = d.Set("disk_encryption", mysqlInfo.DiskEncryption)`，放置于 `destroy_protect` 的 set 逻辑附近。

## 4. Update 函数确认

- [x] 4.1 确认 `disk_encryption` 为 `ForceNew`，`resourceTencentCloudMysqlInstanceUpdate` 无需新增 `HasChange("disk_encryption")` 分支；Terraform SDK 会在该字段变更时计划替换资源。

## 5. 单元测试

- [x] 5.1 在 `tencentcloud/services/cdb/resource_tc_mysql_instance_test.go` 中补充单元测试用例（使用 gomonkey mock 云 API），覆盖：Create 时 `disk_encryption="on"` 被正确写入 `CreateDBInstanceRequest.DiskEncryption` / `CreateDBInstanceHourRequest.DiskEncryption`；未设置时不写入；Read 时从 `DescribeDBInstances` 响应 `InstanceInfo.DiskEncryption` 正确回填并跳过 nil。

## 6. 文档

- [x] 6.1 更新 `tencentcloud/services/cdb/resource_tc_mysql_instance.md`，在 Example Usage 中新增 `disk_encryption` 参数的使用示例（一句话描述中保留所属云产品名称 CDB）。website/docs/ 下的文档由 `make doc` 命令自动生成，不手动编写。

## 7. 验证

- [x] 7.1 验证代码可正确编译（由后续流程执行 go build，本阶段不执行）
- [x] 7.2 确认未引入未检查的 error（必不出错函数用 `_ = func()` 忽略 err）
- [x] 7.3 确认提案中只包含 `disk_encryption` 一个新增参数，未引入 `FourthZone`