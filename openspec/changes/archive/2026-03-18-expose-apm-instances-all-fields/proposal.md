# Change: Expose All DescribeApmInstances Response Fields in tencentcloud_apm_instances Data Source

## Why
当前 `tencentcloud_apm_instances` data source 仅暴露了 `ApmInstanceDetail` 结构体中 16 个字段，SDK 中还有约 30 个字段未暴露，用户无法获取存储用量、日志配置、安全检测开关、Metric 保存时长等关键信息。需要将 `DescribeApmInstances` 接口返回参数全部暴露出来。

## What Changes
- 在 `data_source_tc_apm_instances.go` 的 `instance_list` schema 中新增所有缺失字段定义
- 在 `dataSourceTencentCloudApmInstancesRead` 函数中补充所有缺失字段的数据映射
- 更新 `data_source_tc_apm_instances.md` 文档，补充新增字段说明

### 新增字段列表（按 SDK 结构体顺序）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `amount_of_used_storage` | Float | 存储使用量(MB) |
| `count_of_report_span_per_day` | Int | 日均上报 Span 数 |
| `billing_instance` | Int | 是否已开通计费(0=未开通,1=已开通) |
| `slow_request_saved_threshold` | Int | 采样慢调用保存阈值(ms) |
| `log_region` | String | CLS 日志所在地域 |
| `log_source` | String | 日志源 |
| `is_related_log` | Int | 日志功能开关(0=关,1=开) |
| `log_topic_id` | String | 日志主题 ID |
| `client_count` | Int | 客户端应用数量 |
| `total_count` | Int | 最近2天活跃应用数量 |
| `log_set` | String | CLS 日志集 |
| `metric_duration` | Int | Metric 数据保存时长(天) |
| `custom_show_tags` | List(String) | 用户自定义展示标签列表 |
| `pay_mode_effective` | Bool | 计费模式是否生效 |
| `response_duration_warning_threshold` | Int | 响应时间警示线(ms) |
| `default_tsf` | Int | 是否 TSF 默认业务系统(0=否,1=是) |
| `is_related_dashboard` | Int | 是否关联 Dashboard(0=关,1=开) |
| `dashboard_topic_id` | String | 关联的 Dashboard ID |
| `is_instrumentation_vulnerability_scan` | Int | 是否开启组件漏洞检测(0=关,1=开) |
| `is_sql_injection_analysis` | Int | 是否开启 SQL 注入分析(0=关,1=开) |
| `stop_reason` | Int | 限流原因 |
| `is_remote_command_execution_analysis` | Int | 是否开启远程命令执行检测(0=关,1=开) |
| `is_memory_hijacking_analysis` | Int | 是否开启内存马执行检测(0=关,1=开) |
| `log_index_type` | Int | CLS索引类型(0=全文索引,1=键值索引) |
| `log_trace_id_key` | String | traceId的索引key |
| `is_delete_any_file_analysis` | Int | 是否开启删除任意文件检测(0=关,1=开) |
| `is_read_any_file_analysis` | Int | 是否开启读取任意文件检测(0=关,1=开) |
| `is_upload_any_file_analysis` | Int | 是否开启上传任意文件检测(0=关,1=开) |
| `is_include_any_file_analysis` | Int | 是否开启包含任意文件检测(0=关,1=开) |
| `is_directory_traversal_analysis` | Int | 是否开启目录遍历检测(0=关,1=开) |
| `is_template_engine_injection_analysis` | Int | 是否开启模板引擎注入检测(0=关,1=开) |
| `is_script_engine_injection_analysis` | Int | 是否开启脚本引擎注入检测(0=关,1=开) |
| `is_expression_injection_analysis` | Int | 是否开启表达式注入检测(0=关,1=开) |
| `is_jndi_injection_analysis` | Int | 是否开启JNDI注入检测(0=关,1=开) |
| `is_jni_injection_analysis` | Int | 是否开启JNI注入检测(0=关,1=开) |
| `is_webshell_backdoor_analysis` | Int | 是否开启Webshell后门检测(0=关,1=开) |
| `is_deserialization_analysis` | Int | 是否开启反序列化检测(0=关,1=开) |
| `token` | String | 业务系统鉴权 token |
| `url_long_segment_threshold` | Int | URL长分段收敛阈值 |
| `url_number_segment_threshold` | Int | URL数字分段收敛阈值 |

## Impact
- Affected specs: `apm-instances-datasource`
- Affected code:
  - `tencentcloud/services/apm/data_source_tc_apm_instances.go` — schema 定义 + Read 函数数据映射
  - `tencentcloud/services/apm/data_source_tc_apm_instances.md` — 文档更新
  - `tencentcloud/services/apm/data_source_tc_apm_instances_test.go` — 测试补充新增字段验证
- 非破坏性变更：所有新增字段均为 Computed，不影响现有配置
