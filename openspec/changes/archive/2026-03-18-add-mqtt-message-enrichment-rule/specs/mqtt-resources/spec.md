# Spec: MQTT Message Enrichment Rule Resource

## Overview

本规范定义了 `tencentcloud_mqtt_message_enrichment_rule` Resource 的功能需求，用于管理 MQTT 消息增强规则。

## ADDED Requirements

### Requirement: Resource Schema Definition

Resource 必须支持以下输入参数和输出属性：

**Input Parameters:**
- `instance_id` (String, Required): MQTT 实例 ID
- `rule_name` (String, Required): 规则名称，3-64个字符，支持中文、字母、数字、"-"及"_"
- `condition` (String, Required): 规则匹配条件，JSON格式，需要Base64编码
- `actions` (String, Required): 规则执行的动作，JSON格式，需要Base64编码
- `priority` (Int, Optional): 规则优先级，数字越小，优先级越高，默认为100
- `status` (Int, Optional): 策略状态，0:未定义；1:激活；2:不激活，默认为1
- `remark` (String, Optional): 备注信息

**Output Attributes:**
- `id` (String): 资源 ID，格式为 `{instance_id}#{rule_id}`
- `rule_id` (Int): 规则 ID
- `create_time` (String): 创建时间，毫秒级时间戳
- `update_time` (String): 更新时间，毫秒级时间戳

#### Scenario: Create message enrichment rule successfully

```hcl
resource "tencentcloud_mqtt_message_enrichment_rule" "example" {
  instance_id = "mqtt-12345678"
  rule_name   = "temperature-enrichment"
  condition   = base64encode(jsonencode({
    clientId = "sensor-*"
    topic    = "sensors/+/temperature"
  }))
  actions = base64encode(jsonencode({
    messageExpiryInterval = 360
    responseTopic        = "enriched/temperature/${clientid}"
    correlationData      = "${traceid}"
    userProperty = [
      {
        key   = "enriched-time"
        value = "${timestamp}"
      },
      {
        key   = "data-source"
        value = "rule-engine"
      }
    ]
  }))
  priority = 10
  status   = 1
  remark   = "Temperature sensor data enrichment rule"
}
```

#### Scenario: Update message enrichment rule

```hcl
resource "tencentcloud_mqtt_message_enrichment_rule" "example" {
  instance_id = "mqtt-12345678"
  rule_name   = "temperature-enrichment"
  condition   = base64encode(jsonencode({
    clientId = "sensor-*"
    topic    = "sensors/+/temperature"
    username = "iot-device"
  }))
  actions = base64encode(jsonencode({
    messageExpiryInterval = 600
    responseTopic        = "enriched/temperature/v2/${clientid}"
    correlationData      = "${traceid}"
    userProperty = [
      {
        key   = "enriched-time"
        value = "${timestamp}"
      },
      {
        key   = "data-source"
        value = "rule-engine-v2"
      },
      {
        key   = "rule-version"
        value = "2.0"
      }
    ]
  }))
  priority = 5
  status   = 1
  remark   = "Updated temperature sensor data enrichment rule with version info"
}
```

#### Scenario: Handle invalid parameters
- **WHEN** user provides invalid `status` values (not 0, 1, or 2)
- **THEN** the resource returns validation error with supported values
- **WHEN** user provides invalid `rule_name` (empty, too short, or contains invalid characters)
- **THEN** the resource returns validation error with naming requirements
- **WHEN** user provides invalid Base64 encoded JSON in `condition` or `actions`
- **THEN** the resource returns validation error with format requirements

### Requirement: API Integration

Resource 必须集成 MQTT API v20240516 的消息增强规则相关接口。

#### Scenario: Create operation
- **WHEN** the resource is created
- **THEN** it calls `CreateMessageEnrichmentRuleWithContext` with provided parameters
- **AND** handles API rate limiting and retries appropriately using `tccommon.WriteRetryTimeout`
- **AND** sets the resource ID using instance_id and returned rule_id
- **AND** logs API calls using standard logging patterns

#### Scenario: Read operation
- **WHEN** the resource is read
- **THEN** it calls `DescribeMessageEnrichmentRulesWithContext` to get rule details
- **AND** filters results by instance_id and rule_name to find the specific rule
- **AND** updates the Terraform state with current rule configuration
- **AND** handles cases where the rule no longer exists (sets ID to empty)

#### Scenario: Update operation
- **WHEN** the resource is updated
- **THEN** it calls `ModifyMessageEnrichmentRuleWithContext` with rule ID and changed parameters
- **AND** only updates fields that have actually changed
- **AND** handles partial update scenarios correctly

#### Scenario: Delete operation
- **WHEN** the resource is deleted
- **THEN** it calls `DeleteMessageEnrichmentRuleWithContext` with rule ID
- **AND** handles cases where the rule is already deleted gracefully

### Requirement: Error Handling

Resource 必须正确处理各种错误情况。

#### Scenario: Rule not found
- **WHEN** the specified message enrichment rule does not exist during read/update/delete
- **THEN** the resource handles the error appropriately (remove from state for read, error for update/delete)

#### Scenario: Instance not found
- **WHEN** the specified MQTT instance does not exist
- **THEN** the resource returns a clear error message indicating the instance was not found

#### Scenario: API permission error
- **WHEN** the API call fails due to insufficient permissions
- **THEN** the resource returns the original API error message to help with troubleshooting

#### Scenario: Validation errors
- **WHEN** the API returns validation errors (invalid JSON format, Base64 encoding, etc.)
- **THEN** the resource returns clear error messages with validation details

### Requirement: Import Support

Resource 必须支持 Terraform import 功能。

#### Scenario: Import existing rule
- **WHEN** user runs `terraform import tencentcloud_mqtt_message_enrichment_rule.example mqtt-12345678#123456`
- **THEN** the resource imports the existing rule configuration using instance_id and rule_id
- **AND** populates all readable attributes from the API response
- **AND** sets the correct resource ID format

### Requirement: Base64 Encoding Helper

Resource 必须提供清晰的文档说明如何使用Base64编码的JSON格式。

#### Scenario: Documentation examples
- **WHEN** users read the resource documentation
- **THEN** they see clear examples of how to use `base64encode(jsonencode())` functions
- **AND** understand the expected JSON structure for condition and actions fields
- **AND** have working examples they can copy and modify