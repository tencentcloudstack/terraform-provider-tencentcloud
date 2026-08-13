## 1. 资源代码实现

- [x] 1.1 在 `tencentcloud/services/cfw/` 目录下创建 `resource_tc_cfw_cluster_fw_bypass_config.go`，实现 RESOURCE_KIND_CONFIG 资源的 Schema 定义（仅包含 ModifyClusterFwBypass 接口入参：fw_type、ccn_id、enable、nat_ins_id）
- [x] 1.2 实现资源的 Read 方法：调用 `DescribeClusterNatCcnFwSwitchList` 接口，从返回的 Data 列表中找到匹配实例，将 enable（基于 Bypass 字段）写入 state；若实例不存在则打印日志后调用 d.SetId("")
- [x] 1.3 实现资源的 Update 方法：调用 `ModifyClusterFwBypass` 接口，传入 fw_type、ccn_id、enable（以及 fw_type 为 NAT_FW 时的 nat_ins_id），更新完成后调用 Read 方法刷新 state
- [x] 1.4 实现资源的 Create 方法（CONFIG 类型，仅设置 ID 并调用 Update）和 Delete 方法（CONFIG 类型，仅清空 ID）
- [x] 1.5 资源 ID 使用联合 ID：VPC_FW 类型为 `{fw_type}#{ccn_id}`，NAT_FW 类型为 `{fw_type}#{ccn_id}#{nat_ins_id}`，使用 `tccommon.FILED_SP` 作为分隔符

## 2. 文档文件

- [x] 2.1 创建 `tencentcloud/services/cfw/resource_tc_cfw_cluster_fw_bypass_config.md`，包含资源描述、Example Usage（使用 jsonencode 处理 JSON 字段）和 Import 部分（说明联合 ID 格式）

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_cfw_cluster_fw_bypass_config` 资源，参考 `tencentcloud_igtm_strategy` 资源的注册方式
- [x] 3.2 在 `tencentcloud/provider.md` 中添加 `tencentcloud_cfw_cluster_fw_bypass_config` 资源的文档引用

## 4. 单元测试

- [x] 4.1 创建 `tencentcloud/services/cfw/resource_tc_cfw_cluster_fw_bypass_config_test.go`，使用 gomonkey mock `DescribeClusterNatCcnFwSwitchList` 和 `ModifyClusterFwBypass` 云 API，编写 Read 方法的单元测试（覆盖正常路径和实例不存在路径）
- [x] 4.2 编写 Update 方法的单元测试（覆盖正常路径）
- [x] 4.3 使用 `go test -gcflags=all=-l` 运行单元测试，确保所有测试通过
