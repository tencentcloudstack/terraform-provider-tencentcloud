# Design: tencentcloud_cls_kafka_consumer Resource

## Architecture

新增资源文件 `tencentcloud/services/cls/resource_tc_cls_kafka_consumer.go`，严格参考 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go` 的代码风格。

### 文件结构

```
tencentcloud/services/cls/resource_tc_cls_kafka_consumer.go  # 资源 CRUD 逻辑（新增）
tencentcloud/provider.go                                      # 注册资源（修改）
```

## Detailed Design

### 1. Resource Schema 定义

```go
func ResourceTencentCloudClsKafkaConsumer() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudClsKafkaConsumerCreate,
        Read:   resourceTencentCloudClsKafkaConsumerRead,
        Update: resourceTencentCloudClsKafkaConsumerUpdate,
        Delete: resourceTencentCloudClsKafkaConsumerDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: map[string]*schema.Schema{...},
    }
}
```

**Schema 字段**:

| 字段名 | 类型 | 必选/可选 | ForceNew | Computed | 描述 |
|--------|------|-----------|----------|----------|------|
| `from_topic_id` | TypeString | Required | Yes | No | 日志主题 ID |
| `compression` | TypeInt | Optional | No | Yes | 压缩方式：0-NONE，2-SNAPPY，3-LZ4 |
| `consumer_content` | TypeList(MaxItems:1) | Optional | No | Yes | Kafka 协议消费数据格式 |
| `consumer_content.format` | TypeInt | Optional | No | Yes | 消费数据格式：0-原始内容，1-JSON |
| `consumer_content.enable_tag` | TypeBool | Optional | No | Yes | 是否投递 TAG 信息 |
| `consumer_content.meta_fields` | TypeList(String) | Optional | No | Yes | 元数据信息列表 |
| `consumer_content.tag_transaction` | TypeInt | Optional | No | Yes | tag 数据处理方式：1-不平铺，2-平铺 |
| `consumer_content.json_type` | TypeInt | Optional | No | Yes | 消费数据 JSON 格式：1-不转义，2-转义 |
| `topic_id` | TypeString | - | No | Yes | KafkaConsumer 消费时使用的 Topic 参数（只读） |

### 2. Create 函数

```go
func resourceTencentCloudClsKafkaConsumerCreate(d *schema.ResourceData, meta interface{}) error
```

**流程**:
1. 构建 `OpenKafkaConsumerRequest`
2. 设置 `FromTopicId`（必选）、`Compression`（可选）、`ConsumerContent`（可选）
3. 使用 `resource.Retry` + `tccommon.WriteRetryTimeout` 调用 API
4. 通过 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().OpenKafkaConsumerWithContext(ctx, request)` 调用
5. 检查 response 非 nil
6. 使用 `d.GetOk("from_topic_id")` 获取的值作为资源 ID：`d.SetId(fromTopicId)`
7. 调用 Read 函数刷新状态

### 3. Read 函数

```go
func resourceTencentCloudClsKafkaConsumerRead(d *schema.ResourceData, meta interface{}) error
```

**流程**:
1. 获取资源 ID 即 `fromTopicId = d.Id()`
2. 构建 `DescribeKafkaConsumerRequest`，设置 `FromTopicId`
3. 使用 `resource.Retry` + `tccommon.ReadRetryTimeout` 调用 API
4. 通过 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DescribeKafkaConsumerWithContext(ctx, request)` 调用
5. 检查 `Status` 字段，若为 `false` 表示资源已被关闭，设置 `d.SetId("")` 返回
6. 将返回的 `Compression`、`ConsumerContent`、`TopicID` 设置到 state

### 4. Update 函数

```go
func resourceTencentCloudClsKafkaConsumerUpdate(d *schema.ResourceData, meta interface{}) error
```

**流程**:
1. 获取资源 ID
2. 检测 `compression`、`consumer_content` 是否有变更（`d.HasChange`）
3. 如有变更，构建 `ModifyKafkaConsumerRequest`
4. 设置 `FromTopicId` 以及变更的字段
5. 使用 `resource.Retry` + `tccommon.WriteRetryTimeout` 调用 API
6. 调用 Read 函数刷新状态

### 5. Delete 函数

```go
func resourceTencentCloudClsKafkaConsumerDelete(d *schema.ResourceData, meta interface{}) error
```

**流程**:
1. 获取资源 ID
2. 构建 `CloseKafkaConsumerRequest`，设置 `FromTopicId`
3. 使用 `resource.Retry` + `tccommon.WriteRetryTimeout` 调用 API
4. 返回 nil

### 6. Provider 注册

在 `tencentcloud/provider.go` 的 `ResourcesMap` 中添加：
```go
"tencentcloud_cls_kafka_consumer": cls.ResourceTencentCloudClsKafkaConsumer(),
```

## SDK 类型映射

| SDK 类型 | 字段 | Go 类型 |
|----------|------|---------|
| OpenKafkaConsumerRequest.FromTopicId | *string | |
| OpenKafkaConsumerRequest.Compression | *int64 | |
| OpenKafkaConsumerRequest.ConsumerContent | *KafkaConsumerContent | |
| KafkaConsumerContent.Format | *int64 | |
| KafkaConsumerContent.EnableTag | *bool | |
| KafkaConsumerContent.MetaFields | []*string | |
| KafkaConsumerContent.TagTransaction | *int64 | |
| KafkaConsumerContent.JsonType | *int64 | |
| DescribeKafkaConsumerResponse.Status | *bool | |
| DescribeKafkaConsumerResponse.TopicID | *string | |
| OpenKafkaConsumerResponse.TopicID | *string | |

## 代码风格约束

严格参考 `resource_tc_igtm_strategy.go` 的代码风格：
- 使用 `tccommon.NewResourceLifeCycleHandleFuncContext` 创建 context
- 使用 `helper.String()` / `helper.Int64()` 等辅助函数
- 使用 `resource.Retry` + `tccommon.WriteRetryTimeout` / `tccommon.ReadRetryTimeout` 包装 API 调用
- 日志格式：`log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, ...)`
- 错误日志：`log.Printf("[CRITAL]%s create cls kafka consumer failed, reason:%+v", logId, reqErr)`
- 使用 `GetAPIV3Conn().UseClsClient()` 获取客户端

## HCL 使用示例

```hcl
resource "tencentcloud_cls_kafka_consumer" "example" {
  from_topic_id = "57f5808c-4a55-11eb-b378-0242ac130002"
  compression   = 2

  consumer_content {
    format          = 1
    enable_tag      = true
    meta_fields     = ["__SOURCE__", "__FILENAME__", "__TIMESTAMP__", "__HOSTNAME__", "__PKGID__"]
    tag_transaction = 1
    json_type       = 1
  }
}
```
