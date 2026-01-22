# Change: Add MQTT Message Enrichment Rule Resource

## Why

用户需要通过 Terraform 管理 MQTT 消息增强规则。当前 Provider 已经支持 MQTT 实例、主题、用户、认证等资源的创建和管理，但缺少消息增强规则的管理功能。

消息增强规则是 MQTT 服务的重要功能，允许用户对消息进行预处理、格式转换、字段增强等操作，这会导致用户无法：
1. 通过 Terraform 创建和管理 MQTT 消息增强规则
2. 自动化配置消息处理流程和增强策略
3. 在基础设施即代码(IaC)中统一管理 MQTT 相关资源
4. 实现消息增强规则的版本控制和批量部署

## What Changes

- 新增 Resource: `tencentcloud_mqtt_message_enrichment_rule`
- 实现对 MQTT API v20240516 接口的调用：
  - `CreateMessageEnrichmentRule`: 创建消息增强规则
  - `DescribeMessageEnrichmentRules`: 查询消息增强规则
  - `ModifyMessageEnrichmentRule`: 修改消息增强规则
  - `DeleteMessageEnrichmentRule`: 删除消息增强规则
- 支持完整的 CRUD 操作：
  - `instance_id`: MQTT 实例 ID（必填）
  - `rule_name`: 规则名称（必填，3-64个字符）
  - `condition`: 规则匹配条件，JSON格式，需要Base64编码（必填）
  - `actions`: 规则执行的动作，JSON格式，需要Base64编码（必填）
  - `priority`: 规则优先级（可选，默认100）
  - `status`: 策略状态（可选，0:未定义；1:激活；2:不激活，默认1）
  - `remark`: 备注信息（可选）

## Impact

- **新增能力**: MQTT 消息增强规则管理
- **受影响的服务**: MQTT (tencentcloud/services/mqtt)
- **新增文件**:
  - `tencentcloud/services/mqtt/resource_tc_mqtt_message_enrichment_rule.go`
  - `tencentcloud/services/mqtt/resource_tc_mqtt_message_enrichment_rule.md`
  - `tencentcloud/services/mqtt/resource_tc_mqtt_message_enrichment_rule_test.go`
  - Provider 注册代码需要添加此 resource
- **API 依赖**: 
  - MQTT API v20240516: `CreateMessageEnrichmentRule`, `DescribeMessageEnrichmentRules`, `ModifyMessageEnrichmentRule`, `DeleteMessageEnrichmentRule`
  - 文档: 
    - https://cloud.tencent.com/document/api/1778/126116
    - https://cloud.tencent.com/document/api/1778/126114
    - https://cloud.tencent.com/document/api/1778/126113
    - https://cloud.tencent.com/document/api/1778/126115
- **兼容性**: 无破坏性变更，纯新增功能
- **代码格式要求**: 生成的 resource 代码格式需要参考 `/Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud/tencentcloud/services/igtm/resource_tc_igtm_monitor.go` 的代码风格，包括：
  - 变量声明块的格式和顺序
  - 函数参数处理的模式
  - 错误处理和重试逻辑的实现
  - Schema 定义的结构和注释风格
  - API 调用和响应处理的标准模式
- **特殊注意**: API 使用 Base64 编码的 JSON 格式传递条件和动作配置，需要在文档中提供清晰的使用示例