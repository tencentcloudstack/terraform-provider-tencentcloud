# Spec: APM Instances Data Source

## Overview

本规范定义了 `tencentcloud_apm_instances` Data Source 的功能需求，用于查询 APM (应用性能监控) 业务系统列表。

## ADDED Requirements

### Requirement: Data Source Schema Definition

Data Source 必须支持以下输入参数和输出属性：

**Input Parameters (Optional):**
- `instance_ids` (List of String): 按业务系统 ID 列表精确过滤
- `instance_id` (String): 按业务系统 ID 模糊搜索
- `instance_name` (String): 按业务系统名称模糊搜索
- `tags` (Map): 按标签过滤
- `demo_instance_flag` (Int): 是否查询官方 Demo 实例 (0=非Demo, 1=Demo，默认0)
- `all_regions_flag` (Int): 是否查询全地域实例 (0=不查询, 1=查询，默认0)
- `result_output_file` (String): 输出结果到文件

**Output Attributes:**
- `instance_list` (List): APM 实例列表，每个实例包含以下字段：
  - `instance_id` (String): 实例 ID
  - `instance_name` (String): 实例名称
  - `description` (String): 实例描述
  - `region` (String): 地域
  - `app_id` (Int): 腾讯云账号 AppId
  - `status` (Int): 实例状态
  - `create_uin` (String): 创建者 UIN
  - `trace_duration` (Int): Trace 数据保存时长
  - `span_daily_counters` (Int): 实例日上报 Span 数
  - `pay_mode` (Int): 计费模式
  - `free` (Int): 是否免费版
  - `tags` (Map): 标签列表
  - `err_rate_threshold` (Int): 错误率告警阈值
  - `sample_rate` (Int): 采样率
  - `error_sample` (Int): 错误采样开关
  - 其他相关字段

#### Scenario: Query all APM instances

```hcl
data "tencentcloud_apm_instances" "all" {
}

output "instances" {
  value = data.tencentcloud_apm_instances.all.instance_list
}
```

**期望行为**: 
- 返回当前地域所有 APM 实例
- 不应用任何过滤条件

#### Scenario: Query APM instances by IDs

```hcl
data "tencentcloud_apm_instances" "by_ids" {
  instance_ids = ["apm-xxxxxxxx", "apm-yyyyyyyy"]
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_ids.instance_list
}
```

**期望行为**: 
- 仅返回指定 ID 列表中的 APM 实例
- 使用精确匹配

#### Scenario: Query APM instances by name (fuzzy search)

```hcl
data "tencentcloud_apm_instances" "by_name" {
  instance_name = "test"
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_name.instance_list
}
```

**期望行为**: 
- 返回名称包含 "test" 的所有 APM 实例
- 使用模糊匹配

#### Scenario: Query APM instances by tags

```hcl
data "tencentcloud_apm_instances" "by_tags" {
  tags = {
    "Environment" = "Production"
    "Team"        = "DevOps"
  }
}

output "instances" {
  value = data.tencentcloud_apm_instances.by_tags.instance_list
}
```

**期望行为**: 
- 返回同时包含所有指定标签的 APM 实例
- 标签的 Key 和 Value 必须完全匹配

#### Scenario: Export query results to file

```hcl
data "tencentcloud_apm_instances" "export" {
  instance_name      = "prod"
  result_output_file = "./apm_instances.json"
}
```

**期望行为**: 
- 将查询结果以 JSON 格式写入指定文件
- 文件内容应包含完整的实例列表信息

### Requirement: Service Layer Implementation

必须在 `service_tencentcloud_apm.go` 中实现或复用查询方法支持以下功能：

1. **方法签名** (如需新增):
```go
func (me *ApmService) DescribeApmInstances(ctx context.Context, params map[string]interface{}) (instances []*apm.ApmInstanceDetail, errRet error)
```

2. **参数处理**:
   - 支持 `instance_ids`: 转换为 `[]*string` 数组
   - 支持 `instance_id`: 单个 ID 模糊查询
   - 支持 `instance_name`: 名称模糊查询
   - 支持 `tags`: 转换为 `[]*apm.ApmTag` 数组
   - 支持 `demo_instance_flag`: 整型标志
   - 支持 `all_regions_flag`: 整型标志

3. **错误处理**:
   - API 调用失败时返回详细错误信息
   - 使用重试机制处理临时性错误
   - 记录请求和响应日志

#### Scenario: Service method handles empty result

当 API 返回空列表时，service 方法应：
- 返回空的 instances 切片（不是 nil）
- 不应返回错误
- 记录调试日志

#### Scenario: Service method handles API errors

当 API 返回错误时，service 方法应：
- 将错误封装为 Terraform 可识别的错误类型
- 在日志中记录完整的请求和错误信息
- 支持重试逻辑（使用 `tccommon.ReadRetryTimeout`）

### Requirement: Provider Registration

必须在 `tencentcloud/provider.go` 中的 `DataSourcesMap` 添加新的 data source 注册：

```go
"tencentcloud_apm_instances": apm.DataSourceTencentCloudApmInstances(),
```

#### Scenario: Data source is accessible via Terraform

用户应能够在 Terraform 配置中使用 `data "tencentcloud_apm_instances"` 而不出现 "provider doesn't support data source" 错误。

### Requirement: Documentation

必须提供完整的文档：

1. **Markdown 文档** (`data_source_tc_apm_instances.md`):
   - 一句话描述
   - 至少 3 个使用示例（基本查询、按 ID 过滤、按标签过滤）

2. **Provider 索引** (`tencentcloud/provider.md`):
   - 在 `Application Performance Management(APM)` 部分添加 `Data Source` 节（如不存在）
   - 添加 `tencentcloud_apm_instances` 条目

3. **生成的 HTML 文档**:
   - 运行 `make doc` 后应生成 `website/docs/d/apm_instances.html.markdown`

#### Scenario: User finds data source in documentation

用户查看生成的文档时应能看到：
- Data source 的完整描述
- 所有可用的参数及其说明
- 所有输出属性及其说明
- 实用的代码示例

### Requirement: Testing

必须提供验收测试：

1. **测试文件** (`data_source_tc_apm_instances_test.go`):
   - 实现 `TestAccTencentCloudApmInstancesDataSource_basic` 测试
   - 测试依赖已存在的 APM 实例或在测试中创建临时实例
   - 验证返回的字段非空

2. **测试覆盖**:
   - 基本查询功能
   - 至少一种过滤条件（如按 instance_ids 过滤）
   - 验证返回数据结构完整性

#### Scenario: Basic acceptance test passes

```go
func TestAccTencentCloudApmInstancesDataSource_basic(t *testing.T) {
    t.Parallel()
    resource.Test(t, resource.TestCase{
        PreCheck:  func() { tcacctest.AccPreCheck(t) },
        Providers: tcacctest.AccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccApmInstancesDataSource,
                Check: resource.ComposeTestCheckFunc(
                    tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_apm_instances.instances"),
                    resource.TestCheckResourceAttrSet("data.tencentcloud_apm_instances.instances", "instance_list.#"),
                ),
            },
        },
    })
}
```

## Cross-References

- Related Resource: `tencentcloud_apm_instance` (已存在)
- Related Service: APM Service (`service_tencentcloud_apm.go`)
- API Documentation: https://cloud.tencent.com/document/api/1463/65103
