## Why

用户需要在 Tencent Cloud EdgeOne (TEO) 产品中使用即时转码模板功能来配置和管理视频/音频流的转码规则。目前 Terraform Provider 中缺少对该资源的支持，用户无法通过 Terraform 自动化创建和管理 TEO 即时转码模板。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_just_in_time_transcode_template`，支持即时转码模板的完整生命周期管理
- 实现资源的创建、读取、更新、删除（CRUD）操作
- 支持视频流和音频流的独立配置和控制
- 支持通过 zone_id 和 template_id 进行资源标识和查询

## Capabilities

### New Capabilities
- `teo-just-in-time-transcode-template`: 支持 TEO 即时转码模板的管理，包括模板的创建、查询和删除操作，支持视频流和音频流的独立配置

### Modified Capabilities
（无现有能力的需求变更）

## Impact

- **新增文件**：
  - `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.go` - 资源主文件
  - `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template_test.go` - 资源测试文件
  - `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.md` - 资源文档

- **依赖更新**：使用现有的 `teo/v20220901` SDK 包，无需新增依赖

- **API 接口**：
  - `CreateJustInTimeTranscodeTemplate` - 创建即时转码模板
  - `DescribeJustInTimeTranscodeTemplates` - 查询即时转码模板列表（用于 Read 操作）
  - `DeleteJustInTimeTranscodeTemplates` - 删除即时转码模板

- **注意事项**：
  - 资源 ID 格式：`zone_id#template_id`（复合 ID）
  - 需要处理异步操作，调用创建/删除接口后需轮询查询接口确认生效
  - 必须保持向后兼容，不影响现有 TEO 资源
