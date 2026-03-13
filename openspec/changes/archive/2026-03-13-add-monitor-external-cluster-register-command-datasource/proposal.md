# Change: 新增 Monitor External Cluster Register Command 数据源

## Why
目前 terraform-provider-tencentcloud 缺少获取 Prometheus 外部集群注册命令的数据源。用户在创建外部集群后,需要能够通过 Terraform 获取该集群的注册命令,以便在外部 Kubernetes 集群中执行相应的安装脚本。

## What Changes
- 新增数据源: `tencentcloud_monitor_external_cluster_register_command`
- 支持查询指定 TMP 实例和外部集群的注册命令
- 返回包含安装脚本、配置等完整注册信息

## Impact
- 新增文件: `tencentcloud/services/monitor/data_source_tc_monitor_external_cluster_register_command.go`
- 修改文件: `tencentcloud/provider.go` (注册新数据源)
- 新增测试文件: `tencentcloud/services/monitor/data_source_tc_monitor_external_cluster_register_command_test.go`
- 新增文档: `website/docs/d/monitor_external_cluster_register_command.html.markdown`
- 受影响的包: `tencentcloud/services/monitor`
- 依赖服务: Monitor SDK v20180724 (现有依赖,无需新增)

## API Details

### 查询接口 (DescribeExternalClusterRegisterCommand)
- **接口名称**: DescribeExternalClusterRegisterCommand
- **接口文档**: https://cloud.tencent.com/document/api/248/118965
- **主要参数**:
  - `InstanceId` (必选): TMP 实例 ID
  - `ClusterId` (必选): 集群 ID
- **返回字段**:
  - 注册命令相关信息(需根据接口文档确认具体字段)

## DataSource Design

### 资源 ID 格式
- 格式: `{instanceId}#{clusterId}`
- 示例: `prom-abcd#ecls-1234`
- ID 由 `InstanceId` 和 `ClusterId` 组成

### Schema 设计
```hcl
data "tencentcloud_monitor_external_cluster_register_command" "example" {
  instance_id = "prom-abcd"    # Required
  cluster_id  = "ecls-1234"     # Required
}

output "register_command" {
  value = data.tencentcloud_monitor_external_cluster_register_command.example
}
```

## Implementation Notes

### 重要注意事项
1. **必传参数**: `InstanceId` 和 `ClusterId` 都是必传参数
   
2. **资源唯一 ID**: 数据源的唯一标识为 `instanceId#clusterId`

3. **参考实现**: 代码格式和结构严格参考 `tencentcloud_igtm_instance_list` 数据源
   - 文件位置: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/data_source_tc_igtm_instance_list.go`
   - 遵循相同的代码组织模式、错误处理、重试逻辑和日志记录

4. **服务层方法**: 需要在 `service_tencentcloud_monitor.go` 中添加辅助方法:
   - `DescribeExternalClusterRegisterCommand(ctx, instanceId, clusterId)` 方法

5. **数据源特点**:
   - 数据源是只读操作,不需要 Create/Update/Delete 方法
   - 只需实现 Read 方法
   - 使用 `helper.BuildToken()` 或复合 ID 作为数据源 ID
   - 支持 `result_output_file` 参数用于导出结果

## References
- 参考数据源实现: `data_source_tc_igtm_instance_list.go`
- Monitor 服务层: `service_tencentcloud_monitor.go`
- 相关资源: `resource_tc_monitor_external_cluster.go` (外部集群资源管理)
