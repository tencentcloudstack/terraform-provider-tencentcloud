## Why

为 TEO（TencentCloud EdgeOne）产品新增即时转码模板资源，以满足用户通过 Terraform 管理 TEO 即时转码功能的需求。该资源允许用户定义和管理实时视频转码模板，支持视频和音频流的转码配置，为边缘计算场景提供灵活的媒体处理能力。

## What Changes

- **新增资源**: `tencentcloud_teo_just_in_time_transcode_template`
  - 支持创建即时转码模板，包含视频流和音频流配置
  - 支持查询模板详情和模板列表
  - 支持删除模板
  - 支持更新模板配置（通过删除+重建实现）
  - 资源 ID 格式：`zone_id#template_id`
  - 支持异步操作，包含完整的超时和重试机制

## Capabilities

### New Capabilities

- `teo-just-in-time-transcode-template`: 管理 TEO 即时转码模板资源，包括创建模板、查询模板、更新配置和删除模板的操作。

### Modified Capabilities

- 无

## Impact

**影响范围**：
- 代码文件：
  - 新增 `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.go`
  - 新增 `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template_test.go`
  - 新增 `website/docs/r/teo_just_in_time_transcode_template.md`
  - 修改 `tencentcloud/services/teo/tencentcloud_teo_schema.go`（注册资源）

- API 接口：
  - 使用 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包
  - 调用以下接口：
    - `CreateJustInTimeTranscodeTemplate`
    - `DescribeJustInTimeTranscodeTemplates`
    - `DeleteJustInTimeTranscodeTemplates`

- 依赖：
  - 新资源使用现有的 Terraform Plugin SDK v2 基础设施
  - 使用现有的通用工具函数：`tccommon.Retry()`、`tccommon.LogElapsed()`、`tccommon.InconsistentCheck()`

- 兼容性：
  - 新增资源，不影响现有资源和数据源
  - 向后兼容，无破坏性变更
