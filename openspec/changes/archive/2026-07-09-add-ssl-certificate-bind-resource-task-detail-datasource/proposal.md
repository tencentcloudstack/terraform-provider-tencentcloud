## Why

用户需要通过 Terraform 查询证书关联云资源异步任务的结果详情。当前 SSL 服务已支持通过 `CreateCertificateBindResourceSyncTask` 创建查询证书关联云资源的异步任务，但缺少查询该任务结果详情的数据源，导致用户无法在 Terraform 中获取证书关联的各云资源（CLB、CDN、WAF、DDoS、Live、VOD、TKE、APIGateway、TCB、TEO、TSE、COS、TDMQ、MQTT、GAAP、SCF）的详细绑定信息。

## What Changes

- 新增 Data Source: `tencentcloud_ssl_certificate_bind_resource_task_detail`
- 实现对 SSL API `DescribeCertificateBindResourceTaskDetail` 接口的调用（包名: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205`）
- 该接口为异步查询接口，返回 `Status` 字段（0=查询中，1=查询成功，2=查询异常），需轮询直到任务完成（Status != 0）
- 支持以下查询参数：
  - `task_id` (String, Required): 任务 ID
  - `resource_types` (List of String, Optional): 查询资源类型，不传则查询所有
  - `regions` (List of String, Optional): 查询地域列表
  - `result_output_file` (String, Optional): 输出结果到文件
- 返回各类云资源关联详情：
  - `clb`: 关联 CLB 资源详情列表
  - `cdn`: 关联 CDN 资源详情列表
  - `waf`: 关联 WAF 资源详情列表
  - `ddos`: 关联 DDoS 资源详情列表
  - `live`: 关联 Live 资源详情列表
  - `vod`: 关联 VOD 资源详情列表
  - `tke`: 关联 TKE 资源详情列表
  - `apigateway`: 关联 APIGateway 资源详情列表
  - `tcb`: 关联 TCB 资源详情列表
  - `teo`: 关联 TEO 资源详情列表
  - `tse`: 关联 TSE 资源详情列表
  - `cos`: 关联 COS 资源详情列表
  - `tdmq`: 关联 TDMQ 资源详情列表
  - `mqtt`: 关联 MQTT 资源详情列表
  - `gaap`: 关联 GAAP 资源详情列表
  - `scf`: 关联 SCF 资源详情列表
  - `status`: 异步查询结果状态
  - `cache_time`: 当前结果缓存时间

## Capabilities

### New Capabilities
- `ssl-certificate-bind-resource-task-detail-datasource`: 查询证书关联云资源异步任务结果详情的数据源，支持通过任务 ID 查询证书在各云资源上的绑定详情

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **新增能力**: 证书关联云资源任务结果详情查询
- **受影响的服务**: SSL (tencentcloud/services/ssl)
- **新增文件**:
  - `tencentcloud/services/ssl/data_source_tc_ssl_certificate_bind_resource_task_detail.go`
  - `tencentcloud/services/ssl/data_source_tc_ssl_certificate_bind_resource_task_detail.md`
  - `tencentcloud/services/ssl/data_source_tc_ssl_certificate_bind_resource_task_detail_test.go`
- **修改文件**:
  - `tencentcloud/services/ssl/service_tencent_ssl_certificate.go` — 追加 service 层方法
  - `tencentcloud/provider.go` — 注册新数据源
- **API 依赖**:
  - SSL API v20191205: `DescribeCertificateBindResourceTaskDetail`
- **兼容性**: 无破坏性变更，纯新增功能
