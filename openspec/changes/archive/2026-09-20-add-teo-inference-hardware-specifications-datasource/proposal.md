## Why

TEO（EdgeOne）推理服务在创建时需要指定硬件规格（`HardwareSpecId`），而该规格 ID 需要通过 `DescribeInferenceHardwareSpecifications` 接口按站点查询获取。当前 Terraform Provider 缺少查询推理硬件规格列表的数据源，用户无法在 Terraform 配置中直接获取可用硬件规格，只能手动查询后填写，影响自动化部署体验。

## What Changes

- 新增 Data Source: `tencentcloud_teo_inference_hardware_specifications`
- 实现对 TEO API `DescribeInferenceHardwareSpecifications` 接口的调用，按站点查询推理硬件规格列表
- 支持通过站点 ID 查询：
  - `zone_id`: 站点 ID（必填）
- 返回硬件规格列表信息，每个规格包含以下字段：
  - `spec`: 规格标识（已废弃，参考使用 `hardware_spec_id`）
  - `hardware_spec_id`: 规格唯一标识 ID
  - `name`: 规格名称
  - `gpu_num`: 规格默认分配的 GPU 卡数
  - `cpu_num`: 规格默认分配的 CPU 核数
  - `mem_size`: 规格默认分配的内存大小（MB）
  - `gpu_mem_size`: 规格默认分配的显存大小（MB）
  - `disk_size`: 规格默认分配的磁盘大小（MB）
  - `allowed_gpu_nums`: 规格当前支持的 GPU 卡数列表

## Capabilities

### New Capabilities
- `teo-inference-hardware-specifications-datasource`: 查询 TEO 推理硬件规格列表的数据源，支持按站点 ID 获取可用硬件规格及其资源配置信息。

### Modified Capabilities
<!-- 无现有 capability 的需求变更 -->

## Impact

- **新增能力**: TEO 推理硬件规格列表查询
- **受影响的服务**: TEO (tencentcloud/services/teo)
- **新增文件**:
  - `tencentcloud/services/teo/data_source_tc_teo_inference_hardware_specifications.go`
  - `tencentcloud/services/teo/data_source_tc_teo_inference_hardware_specifications.md`
  - `tencentcloud/services/teo/data_source_tc_teo_inference_hardware_specifications_test.go`
- **修改文件**:
  - `tencentcloud/provider.go`: 添加数据源注册代码
  - `tencentcloud/provider.md`: 添加数据源文档条目
- **API 依赖**: 
  - TEO API v20220901: `DescribeInferenceHardwareSpecifications`
- **兼容性**: 无破坏性变更，纯新增功能
