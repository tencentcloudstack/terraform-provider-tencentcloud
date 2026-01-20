# Spec: CKafka Version Data Source

## Overview

本规范定义了 `tencentcloud_ckafka_version` Data Source 的功能需求，用于查询 CKafka 实例的版本信息。

## ADDED Requirements

### Requirement: Data Source Schema Definition

Data Source 必须支持以下输入参数和输出属性：

**Input Parameters:**
- `instance_id` (String, Required): CKafka 实例 ID
- `result_output_file` (String, Optional): 输出结果到文件

**Output Attributes:**
- `kafka_version` (String): 当前 Kafka 大版本号
- `cur_broker_version` (String): 当前运行的 Broker 版本号
- `latest_broker_versions` (List): 平台支持的最新 Broker 版本列表，每个版本包含：
  - `kafka_version` (String): Kafka 版本号
  - `broker_version` (String): Broker 版本号

#### Scenario: Query instance version successfully

```hcl
data "tencentcloud_ckafka_version" "example" {
  instance_id = "ckafka-bqwlyrg8"
}

output "current_version" {
  value = {
    kafka_version = data.tencentcloud_ckafka_version.example.kafka_version
    broker_version = data.tencentcloud_ckafka_version.example.cur_broker_version
  }
}

output "available_versions" {
  value = data.tencentcloud_ckafka_version.example.latest_broker_versions
}
```

#### Scenario: Handle invalid instance ID
- **WHEN** user provides an invalid or non-existent CKafka instance ID
- **THEN** the data source returns an appropriate error message from the API

#### Scenario: Output to file
- **WHEN** user specifies `result_output_file` parameter
- **THEN** the version information is written to the specified file in JSON format

### Requirement: API Integration

Data Source 必须集成 CKafka API v20190819 的 `DescribeCkafkaVersion` 接口。

#### Scenario: API call execution
- **WHEN** the data source is read
- **THEN** it calls `DescribeCkafkaVersionWithContext` with the provided instance ID
- **AND** handles API rate limiting and retries appropriately using `tccommon.WriteRetryTimeout`
- **AND** processes the response according to Terraform Plugin SDK patterns
- **AND** logs API calls using `tccommon.LogElapsed` and debug logging

### Requirement: Error Handling

Data Source 必须正确处理各种错误情况。

#### Scenario: Instance not found
- **WHEN** the specified CKafka instance does not exist
- **THEN** the data source returns a clear error message indicating the instance was not found

#### Scenario: API permission error
- **WHEN** the API call fails due to insufficient permissions
- **THEN** the data source returns the original API error message to help with troubleshooting

#### Scenario: Network or service error
- **WHEN** the API call fails due to network issues or service unavailability
- **THEN** the data source retries according to the configured retry policy and returns an appropriate error if all retries fail