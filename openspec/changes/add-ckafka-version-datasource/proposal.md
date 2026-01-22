# Change: Add CKafka Version Data Source

## Why

用户需要通过 Terraform 查询 CKafka 实例的版本信息。当前 Provider 已经支持 CKafka 实例的创建和管理（`tencentcloud_ckafka_instance`），以及多个 CKafka 相关的数据源，但缺少查询实例版本信息的数据源。

这会导致用户无法：
1. 通过 Terraform 查询 CKafka 实例的当前 Kafka 版本和 Broker 版本
2. 获取平台支持的最新 Broker 版本列表，用于版本升级决策
3. 在 Terraform 配置中基于版本信息进行条件判断和资源配置

## What Changes

- 新增 Data Source: `tencentcloud_ckafka_version`
- 实现对 CKafka API `DescribeCkafkaVersion` 接口的调用
- 支持通过实例 ID 查询版本信息：
  - `instance_id`: CKafka 实例 ID（必填）
  - `result_output_file`: 输出结果到文件（可选）
- 返回版本信息：
  - `kafka_version`: 当前 Kafka 大版本号
  - `cur_broker_version`: 当前运行的 Broker 版本号
  - `latest_broker_versions`: 平台支持的最新 Broker 版本列表

## Impact

- **新增能力**: CKafka 实例版本信息查询
- **受影响的服务**: CKafka (tencentcloud/services/ckafka)
- **新增文件**:
  - `tencentcloud/services/ckafka/data_source_tc_ckafka_version.go`
  - `tencentcloud/services/ckafka/data_source_tc_ckafka_version.md`
  - `tencentcloud/services/ckafka/data_source_tc_ckafka_version_test.go`
  - Provider 注册代码需要添加此 data source
- **API 依赖**: 
  - CKafka API v20190819: `DescribeCkafkaVersion`
  - 文档: https://cloud.tencent.com/document/api/597/125479
- **兼容性**: 无破坏性变更，纯新增功能