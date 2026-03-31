# Proposal: Add tencentcloud_cls_kafka_consumer Resource

## Summary

新增 CLS（日志服务）Kafka 协议消费资源 `tencentcloud_cls_kafka_consumer`，用于管理日志主题的 Kafka 协议消费功能。通过该资源可以打开、查询、修改和关闭 Kafka 协议消费。

## Problem Statement

当前 Terraform Provider 缺少对 CLS Kafka 协议消费功能的原生资源支持。用户无法通过 Terraform 管理日志主题的 Kafka 协议消费配置（包括压缩方式、消费数据格式等），影响了基础设施即代码的完整性。

## Proposed Solution

新增一个 Terraform Resource `tencentcloud_cls_kafka_consumer`，实现完整的 CRUD 操作：

- **Create**: 调用 `OpenKafkaConsumer` 接口打开 Kafka 协议消费
- **Read**: 调用 `DescribeKafkaConsumer` 接口查询 Kafka 协议消费信息
- **Update**: 调用 `ModifyKafkaConsumer` 接口修改 Kafka 协议消费配置
- **Delete**: 调用 `CloseKafkaConsumer` 接口关闭 Kafka 协议消费

### 资源 Schema

| 属性 | 类型 | 必选 | ForceNew | 描述 |
|------|------|------|----------|------|
| `from_topic_id` | String | 是 | 是 | 日志主题 ID |
| `compression` | Integer | 否 | 否 | 压缩方式：0-NONE，2-SNAPPY，3-LZ4 |
| `consumer_content` | List(MaxItems:1) | 否 | 否 | Kafka 协议消费数据格式 |
| `consumer_content.format` | Integer | 否 | 否 | 消费数据格式：0-原始内容，1-JSON |
| `consumer_content.enable_tag` | Boolean | 否 | 否 | 是否投递 TAG 信息 |
| `consumer_content.meta_fields` | List(String) | 否 | 否 | 元数据信息列表 |
| `consumer_content.tag_transaction` | Integer | 否 | 否 | tag 数据处理方式：1-不平铺，2-平铺 |
| `consumer_content.json_type` | Integer | 否 | 否 | 消费数据 JSON 格式：1-不转义，2-转义 |
| `topic_id` | String | Computed | - | KafkaConsumer 消费时使用的 Topic ID（由 OpenKafkaConsumer 返回） |

### 资源唯一标识

`OpenKafkaConsumer` 接口返回 `TopicID`，该资源以 `from_topic_id` 作为唯一标识（资源 ID）。

## Impact

### Affected Resources

- 新增 `tencentcloud_cls_kafka_consumer` 资源

### Affected Files

- 新增 `tencentcloud/services/cls/resource_tc_cls_kafka_consumer.go`（资源 CRUD 逻辑）
- 修改 `tencentcloud/provider.go`（注册资源）

### API Dependencies

| 操作 | 接口 | 文档 |
|------|------|------|
| Create | OpenKafkaConsumer | https://cloud.tencent.com/document/api/614/72339 |
| Read | DescribeKafkaConsumer | https://cloud.tencent.com/document/api/614/95719 |
| Update | ModifyKafkaConsumer | https://cloud.tencent.com/document/api/614/95720 |
| Delete | CloseKafkaConsumer | https://cloud.tencent.com/document/api/614/72340 |

## Reference

- 代码风格参考：`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`
- SDK 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`
- SDK Client 调用方式：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient()`
