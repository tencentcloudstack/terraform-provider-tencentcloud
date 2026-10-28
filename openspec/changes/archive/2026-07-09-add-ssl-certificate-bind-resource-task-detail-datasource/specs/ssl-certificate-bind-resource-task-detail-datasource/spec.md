## ADDED Requirements

### Requirement: Data Source Schema Definition

Data Source `tencentcloud_ssl_certificate_bind_resource_task_detail` MUST support the following input parameters and output attributes.

**Input Parameters:**
- `task_id` (String, Required): 任务 ID，根据 CreateCertificateBindResourceSyncTask 得到的任务 ID 查询绑定云资源结果
- `resource_types` (List of String, Optional): 查询资源类型的结果详情，不传则查询所有。取值支持：clb、cdn、ddos、live、vod、waf、apigateway、teo、tke、cos、tse、tcb、tdmq、mqtt、gaap、scf
- `regions` (List of String, Optional): 查询地域列表的数据，clb、tke、waf、apigateway、tcb、cos、tse 支持地域查询
- `result_output_file` (String, Optional): 输出结果到文件

**Output Attributes:**
- `status` (Int): 关联云资源异步查询结果状态，0 表示查询中，1 表示查询成功，2 表示查询异常
- `cache_time` (String): 当前结果缓存时间
- `clb` (List): 关联 CLB 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 CLB 实例详情）
- `cdn` (List): 关联 CDN 资源详情列表，每个元素包含：`total_count`、`error`、`instance_list`（嵌套 CDN 实例详情）
- `waf` (List): 关联 WAF 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 WAF 实例详情）
- `ddos` (List): 关联 DDoS 资源详情列表，每个元素包含：`total_count`、`error`、`instance_list`（嵌套 DDoS 实例详情）
- `live` (List): 关联 Live 资源详情列表，每个元素包含：`total_count`、`error`、`instance_list`（嵌套 Live 实例详情）
- `vod` (List): 关联 VOD 资源详情列表，每个元素包含：`total_count`、`error`、`instance_list`（嵌套 VOD 实例详情）
- `tke` (List): 关联 TKE 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 TKE 实例详情）
- `apigateway` (List): 关联 APIGateway 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 APIGateway 实例详情）
- `tcb` (List): 关联 TCB 资源详情列表，每个元素包含：`region`、`error`、`environments`（嵌套 TCB 环境详情）
- `teo` (List): 关联 TEO 资源详情列表，每个元素包含：`total_count`、`error`、`instance_list`（嵌套 TEO 实例详情）
- `tse` (List): 关联 TSE 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 TSE 实例详情）
- `cos` (List): 关联 COS 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 COS 实例详情）
- `tdmq` (List): 关联 TDMQ 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 TDMQ 实例详情）
- `mqtt` (List): 关联 MQTT 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 MQTT 实例详情）
- `gaap` (List): 关联 GAAP 资源详情列表，每个元素包含：`total_count`、`error`、`instance_list`（嵌套 GAAP 实例详情）
- `scf` (List): 关联 SCF 资源详情列表，每个元素包含：`region`、`total_count`、`error`、`instance_list`（嵌套 SCF 实例详情）

#### Scenario: Query task detail with task_id only
- **WHEN** user provides only `task_id`
- **THEN** the data source queries all resource types' bind details for the given task
- **AND** returns the status, cache_time and all 16 cloud resource detail lists

#### Scenario: Query task detail with resource_types filter
- **WHEN** user provides `task_id` and `resource_types` (e.g., `["clb", "cdn"]`)
- **THEN** the data source queries only the specified resource types' bind details

#### Scenario: Query task detail with regions filter
- **WHEN** user provides `task_id` and `regions`
- **THEN** the data source queries the bind details for the specified regions (supported by clb, tke, waf, apigateway, tcb, cos, tse)

#### Scenario: Output to file
- **WHEN** user specifies `result_output_file` parameter
- **THEN** the bind resource detail information is written to the specified file in JSON format

### Requirement: Async Task Polling

Data Source MUST correctly handle async task status. The `DescribeCertificateBindResourceTaskDetail` API returns a `Status` field, where 0 means querying in progress, 1 means query success, and 2 means query exception.

#### Scenario: Task is still in progress
- **WHEN** the API returns `Status` = 0 (querying in progress)
- **THEN** the data source retries the query until `Status` is not 0, within `tccommon.ReadRetryTimeout`

#### Scenario: Task completed successfully
- **WHEN** the API returns `Status` = 1 (query success)
- **THEN** the data source returns the full bind resource detail results

#### Scenario: Task failed with error
- **WHEN** the API returns `Status` = 2 (query exception)
- **THEN** the data source returns the result without further retry, allowing the user to inspect the `error` field of each resource type

### Requirement: API Integration

Data Source MUST integrate with the SSL API v20191205 `DescribeCertificateBindResourceTaskDetail` interface.

#### Scenario: API call with pagination
- **WHEN** the data source is read
- **THEN** it calls `DescribeCertificateBindResourceTaskDetail` with `Limit` set to 100 (the maximum allowed by the API) and iterates `Offset` internally until all results are retrieved
- **AND** does not expose `limit`/`offset` parameters to the user

#### Scenario: API call wrapped in retry
- **WHEN** the data source is read
- **THEN** the API call is wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`
- **AND** uses `tccommon.RetryError()` to wrap errors returned by the API

#### Scenario: Empty response handling
- **WHEN** the API returns an empty response (`response == nil` or `response.Response == nil`)
- **THEN** the data source returns `NonRetryableError` instead of clearing the id, to avoid data loss from transient API fluctuations
- **AND** logs `[DATASOURCE] read empty, skip SetId` for troubleshooting

### Requirement: Error Handling

Data Source MUST correctly handle various error conditions.

#### Scenario: Invalid task ID
- **WHEN** user provides an invalid or non-existent task ID
- **THEN** the data source returns the API error message (e.g., `FailedOperation.CertificateSyncTaskIdInvalid`)

#### Scenario: API permission error
- **WHEN** the API call fails due to insufficient permissions (e.g., `FailedOperation.RoleNotFoundAuthorization`)
- **THEN** the data source returns the original API error message to help with troubleshooting

#### Scenario: Network or service error
- **WHEN** the API call fails due to network issues or service unavailability
- **THEN** the data source retries according to the configured retry policy and returns an appropriate error if all retries fail

### Requirement: Nested Resource Detail Schema

Data Source MUST completely map the nested structure of each cloud resource detail, following the struct definitions in the vendor SDK.

#### Scenario: CLB instance detail structure
- **WHEN** the API returns CLB instance details
- **THEN** each CLB instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each CLB instance detail contains: `load_balancer_id`, `load_balancer_name`, `forward`, `listeners` (nested list)
- **AND** each CLB listener contains: `listener_id`, `listener_name`, `sni_switch`, `protocol`, `certificate` (nested struct), `rules` (nested list), `no_match_domains`
- **AND** each CLB listener rule contains: `location_id`, `domain`, `is_match`, `certificate` (nested struct), `no_match_domains`, `url`
- **AND** the certificate struct contains: `cert_id`, `dns_names`, `cert_ca_id`, `s_s_l_mode`

#### Scenario: CDN instance detail structure
- **WHEN** the API returns CDN instance details
- **THEN** each CDN instance list element contains: `total_count`, `error`, `instance_list` (nested list)
- **AND** each CDN instance detail contains: `domain`, `cert_id`, `status`, `https_billing_switch`

#### Scenario: WAF instance detail structure
- **WHEN** the API returns WAF instance details
- **THEN** each WAF instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each WAF instance detail contains: `domain`, `cert_id`, `keepalive`

#### Scenario: DDOS instance detail structure
- **WHEN** the API returns DDOS instance details
- **THEN** each DDOS instance list element contains: `total_count`, `error`, `instance_list` (nested list)
- **AND** each DDOS instance detail contains: `domain`, `instance_id`, `protocol`, `cert_id`, `virtual_port`

#### Scenario: LIVE instance detail structure
- **WHEN** the API returns LIVE instance details
- **THEN** each LIVE instance list element contains: `total_count`, `error`, `instance_list` (nested list)
- **AND** each LIVE instance detail contains: `domain`, `cert_id`, `status`

#### Scenario: VOD instance detail structure
- **WHEN** the API returns VOD instance details
- **THEN** each VOD instance list element contains: `total_count`, `error`, `instance_list` (nested list)
- **AND** each VOD instance detail contains: `domain`, `cert_id`

#### Scenario: TKE instance detail structure
- **WHEN** the API returns TKE instance details
- **THEN** each TKE instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each TKE instance detail contains: `cluster_id`, `cluster_name`, `cluster_type`, `cluster_version`, `namespace_list` (nested list)
- **AND** each TKE namespace detail contains: `name`, `secret_list` (nested list)
- **AND** each TKE secret detail contains: `name`, `cert_id`, `ingress_list` (nested list), `no_match_domains`
- **AND** each TKE ingress detail contains: `ingress_name`, `tls_domains`, `domains`

#### Scenario: APIGateway instance detail structure
- **WHEN** the API returns APIGateway instance details
- **THEN** each APIGateway instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each APIGateway instance detail contains: `service_id`, `service_name`, `domain`, `cert_id`, `protocol`

#### Scenario: TCB environment structure
- **WHEN** the API returns TCB details
- **THEN** each TCB instance list element contains: `region`, `error`, `environments` (nested list)
- **AND** each TCB environment wrapper contains: `environment` (nested struct), `access_service` (nested struct), `host_service` (nested struct)
- **AND** the TCB environment struct contains: `id`, `source`, `name`, `status`
- **AND** the TCB access service contains: `instance_list` (nested list with `domain`, `status`, `union_status`, `is_preempted`, `i_c_p_status`, `old_certificate_id`), `total_count`
- **AND** the TCB host service contains: `instance_list` (nested list with `domain`, `status`, `d_n_s_status`, `old_certificate_id`), `total_count`

#### Scenario: TEO instance detail structure
- **WHEN** the API returns TEO instance details
- **THEN** each TEO instance list element contains: `total_count`, `error`, `instance_list` (nested list)
- **AND** each TEO instance detail contains: `host`, `cert_id`, `zone_id`, `status`, `algorithm`

#### Scenario: TSE instance detail structure
- **WHEN** the API returns TSE instance details
- **THEN** each TSE instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each TSE instance detail contains: `gateway_id`, `gateway_name`, `certificate_list` (nested list)
- **AND** each gateway certificate contains: `id`, `name`, `bind_domains`, `cert_source`, `cert_id`

#### Scenario: COS instance detail structure
- **WHEN** the API returns COS instance details
- **THEN** each COS instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each COS instance detail contains: `domain`, `cert_id`, `status`, `bucket`, `region`

#### Scenario: TDMQ instance detail structure
- **WHEN** the API returns TDMQ instance details
- **THEN** each TDMQ instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each TDMQ instance detail contains: `instance_id`, `instance_name`, `instance_status`, `cert_id`, `ca_cert_id`, `no_match_domains`

#### Scenario: MQTT instance detail structure
- **WHEN** the API returns MQTT instance details
- **THEN** each MQTT instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each MQTT instance detail contains: `instance_id`, `instance_name`, `instance_status`, `no_match_domains`, `server_cert_id_list`, `ca_cert_id_list`

#### Scenario: GAAP instance detail structure
- **WHEN** the API returns GAAP instance details
- **THEN** each GAAP instance list element contains: `total_count`, `error`, `instance_list` (nested list)
- **AND** each GAAP instance detail contains: `instance_id`, `instance_name`, `listener_list` (nested list)
- **AND** each GAAP listener detail contains: `listener_status`, `listener_id`, `listener_name`, `no_match_domains`, `cert_id_list`, `protocol`

#### Scenario: SCF instance detail structure
- **WHEN** the API returns SCF instance details
- **THEN** each SCF instance list element contains: `region`, `total_count`, `error`, `instance_list` (nested list)
- **AND** each SCF instance detail contains: `certificate_id`, `protocol`, `domain`, `region`
