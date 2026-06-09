# Tasks: Add tencentcloud_cls_kafka_consumer Resource

## Task 1: 创建资源文件 resource_tc_cls_kafka_consumer.go

**文件**: `tencentcloud/services/cls/resource_tc_cls_kafka_consumer.go`

**代码风格参考**: `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`

### 1.1 Package 和 Import

```go
package cls

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)
```

### 1.2 Resource Schema 定义

实现 `ResourceTencentCloudClsKafkaConsumer()` 函数，定义以下 Schema：

- `from_topic_id`: TypeString, Required, ForceNew, Description: "日志主题ID"
- `compression`: TypeInt, Optional, Computed, Description: "压缩方式：0-NONE，2-SNAPPY，3-LZ4"
- `consumer_content`: TypeList, Optional, Computed, MaxItems:1, Description: "Kafka协议消费数据格式"
  - `format`: TypeInt, Optional, Computed
  - `enable_tag`: TypeBool, Optional, Computed
  - `meta_fields`: TypeList(TypeString), Optional, Computed
  - `tag_transaction`: TypeInt, Optional, Computed
  - `json_type`: TypeInt, Optional, Computed
- `topic_id`: TypeString, Computed, Description: "KafkaConsumer消费时使用的Topic参数"

### 1.3 Create 函数

实现 `resourceTencentCloudClsKafkaConsumerCreate`：

1. 使用 `tccommon.NewResourceLifeCycleHandleFuncContext` 创建 context
2. 构建 `cls.NewOpenKafkaConsumerRequest()`
3. 设置 `FromTopicId`, 可选设置 `Compression` 和 `ConsumerContent`
4. `ConsumerContent` 为嵌套结构 `cls.KafkaConsumerContent`，包含 `Format`/`EnableTag`/`MetaFields`/`TagTransaction`/`JsonType`
5. 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `UseClsClient().OpenKafkaConsumerWithContext(ctx, request)`
6. 检查 response 非 nil
7. 使用 `fromTopicId` 作为资源 ID：`d.SetId(fromTopicId)`
8. 调用 Read 刷新状态

### 1.4 Read 函数

实现 `resourceTencentCloudClsKafkaConsumerRead`：

1. `fromTopicId = d.Id()`
2. 构建 `cls.NewDescribeKafkaConsumerRequest()`，设置 `FromTopicId`
3. 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 `UseClsClient().DescribeKafkaConsumerWithContext(ctx, request)`
4. 检查 `response.Response.Status`，若为 `false` 则 `d.SetId("")` 并返回
5. 设置 state：
   - `d.Set("from_topic_id", fromTopicId)`
   - `d.Set("compression", response.Response.Compression)`
   - `d.Set("topic_id", response.Response.TopicID)`
   - 将 `ConsumerContent` 转换为 `[]map[string]interface{}` 格式设置

### 1.5 Update 函数

实现 `resourceTencentCloudClsKafkaConsumerUpdate`：

1. 获取 `fromTopicId = d.Id()`
2. 检测 `d.HasChange("compression")` 或 `d.HasChange("consumer_content")`
3. 如有变更，构建 `cls.NewModifyKafkaConsumerRequest()`
4. 设置 `FromTopicId`、`Compression`（从当前 state 获取）、`ConsumerContent`（从当前 state 获取）
5. 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `UseClsClient().ModifyKafkaConsumerWithContext(ctx, request)`
6. 调用 Read 刷新状态

### 1.6 Delete 函数

实现 `resourceTencentCloudClsKafkaConsumerDelete`：

1. 获取 `fromTopicId = d.Id()`
2. 构建 `cls.NewCloseKafkaConsumerRequest()`，设置 `FromTopicId`
3. 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `UseClsClient().CloseKafkaConsumerWithContext(ctx, request)`
4. 返回 nil

### 验收标准

- [x] 文件路径正确：`tencentcloud/services/cls/resource_tc_cls_kafka_consumer.go`
- [x] 包名为 `cls`
- [x] Import 语句完整
- [x] Schema 定义包含所有字段
- [x] Create 函数调用 OpenKafkaConsumer 并以 fromTopicId 作为资源 ID
- [x] Read 函数调用 DescribeKafkaConsumer 并正确设置所有 state
- [x] Read 函数检查 Status 为 false 时清除资源
- [x] Update 函数检测变更并调用 ModifyKafkaConsumer
- [x] Delete 函数调用 CloseKafkaConsumer
- [x] 代码风格与 resource_tc_igtm_strategy.go 一致
- [x] 使用 `UseClsClient()` 获取客户端
- [x] 日志格式正确（DEBUG/CRITAL）

---

## Task 2: 注册资源到 Provider

**文件**: `tencentcloud/provider.go`

### 实施步骤

在 `ResourcesMap` 中 CLS 相关资源区域添加：

```go
"tencentcloud_cls_kafka_consumer": cls.ResourceTencentCloudClsKafkaConsumer(),
```

### 验收标准

- [x] 在 provider.go 的 ResourcesMap 中正确注册
- [x] 资源名称为 `tencentcloud_cls_kafka_consumer`
- [x] 引用函数名为 `cls.ResourceTencentCloudClsKafkaConsumer()`

---

## Task 3: 代码编译验证

### 3.1 格式化

```bash
cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
gofmt -w tencentcloud/services/cls/resource_tc_cls_kafka_consumer.go
```

### 3.2 编译

```bash
go build ./tencentcloud/services/cls/
go build ./tencentcloud/...
```

### 验收标准

- [x] `gofmt` 格式化成功
- [x] `go build` 编译成功，无语法错误
- [x] 无 linter 错误（可接受已有的 deprecated 警告）

---

## 验收清单

### 功能完整性

- [x] 资源文件包含完整的 CRUD 四个函数
- [x] Create 调用 OpenKafkaConsumer
- [x] Read 调用 DescribeKafkaConsumer
- [x] Update 调用 ModifyKafkaConsumer
- [x] Delete 调用 CloseKafkaConsumer
- [x] 支持 Import（使用 ImportStatePassthrough）
- [x] 资源 ID 使用 fromTopicId

### 代码质量

- [x] 严格参考 igtm_strategy 代码风格
- [x] 使用 resource.Retry 包装 API 调用
- [x] 完整的错误处理和日志记录
- [x] 正确的 nil 检查

### Schema 完整性

- [x] from_topic_id: Required, ForceNew
- [x] compression: Optional, Computed
- [x] consumer_content: Optional, Computed, MaxItems:1
- [x] consumer_content 子字段完整（format/enable_tag/meta_fields/tag_transaction/json_type）
- [x] topic_id: Computed

---

## 实施顺序

1. **Task 1**: 创建资源文件（核心任务）
2. **Task 2**: 注册到 Provider
3. **Task 3**: 编译验证

---

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

## 参考信息

### API 文档

| 操作 | 接口 | 文档链接 |
|------|------|----------|
| 创建 | OpenKafkaConsumer | https://cloud.tencent.com/document/api/614/72339 |
| 查询 | DescribeKafkaConsumer | https://cloud.tencent.com/document/api/614/95719 |
| 修改 | ModifyKafkaConsumer | https://cloud.tencent.com/document/api/614/95720 |
| 删除 | CloseKafkaConsumer | https://cloud.tencent.com/document/api/614/72340 |

### SDK 结构体

- `cls.OpenKafkaConsumerRequest`: FromTopicId(*string), Compression(*int64), ConsumerContent(*KafkaConsumerContent)
- `cls.OpenKafkaConsumerResponse`: TopicID(*string)
- `cls.DescribeKafkaConsumerRequest`: FromTopicId(*string)
- `cls.DescribeKafkaConsumerResponse`: Status(*bool), TopicID(*string), Compression(*int64), ConsumerContent(*KafkaConsumerContent)
- `cls.ModifyKafkaConsumerRequest`: FromTopicId(*string), Compression(*int64), ConsumerContent(*KafkaConsumerContent)
- `cls.CloseKafkaConsumerRequest`: FromTopicId(*string)
- `cls.KafkaConsumerContent`: Format(*int64), EnableTag(*bool), MetaFields([]*string), TagTransaction(*int64), JsonType(*int64)

### 代码风格参考

**文件**: `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/resource_tc_igtm_strategy.go`

关键模式：
- context 创建：`tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`
- Client 获取：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient()`
- API 调用包装：`resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})`
- 日志格式：`log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())`
- 错误日志：`log.Printf("[CRITAL]%s create cls kafka consumer failed, reason:%+v", logId, reqErr)`
