# teo-just-in-time-transcode-template Specification

## Purpose

定义 TEO 服务即时转码模板资源的行为规范。

## ADDED Requirements

### Requirement: Resource Schema Definition

Resource 必须支持以下输入参数：

**Required Parameters (ForceNew):**
- `zone_id` (String, Required): 站点 ID，用于标识 TE0 站点
- `template_name` (String, Required): 即时转码模板名称，长度限制：64 个字符

**Optional Parameters (ForceNew):**
- `comment` (String, Optional): 模板描述信息，长度限制：256 个字符
- `video_stream_switch` (String, Optional, Default: "on"): 启用视频流开关，取值：on/off
- `audio_stream_switch` (String, Optional, Default: "on"): 启用音频流开关，取值：on/off
- `video_template` (Block, Optional): 视频流配置参数，当 video_stream_switch 为 on 时该字段必填
  - `codec` (String, Optional): 视频流编码格式，可选值：H.264/H.265
  - `fps` (Float, Optional): 视频帧率，取值范围：[0, 30]，单位：Hz，默认值：0
  - `bitrate` (Int, Optional): 视频流码率，取值范围：0 和 [128, 10000]，单位：kbps，默认值：0
  - `resolution_adaptive` (String, Optional, Default: "open"): 分辨率自适应，可选值：open/close
  - `width` (Int, Optional): 视频流宽度或长边的最大值，取值范围：0 和 [128, 1920]，单位：px，默认值：0
  - `height` (Int, Optional): 视频流高度或短边的最大值，取值范围：0 和 [128, 1080]，单位：px，默认值：0
  - `fill_type` (String, Optional, Default: "black"): 填充方式，可选值：stretch/black/white/gauss
- `audio_template` (Block, Optional): 音频流配置参数，当 audio_stream_switch 为 on 时该字段必填
  - `codec` (String, Optional): 音频流编码格式，可选值：libfdk_aac
  - `audio_channel` (Int, Optional, Default: 2): 音频通道数，可选值：2

**Computed Parameters:**
- `template_id` (String, Computed): 即时转码模板唯一标识
- `create_time` (String, Computed): 模板创建时间，使用 ISO 日期格式
- `update_time` (String, Computed): 模板最后修改时间，使用 ISO 日期格式

**Timeouts:**
- `create`: 创建操作的默认超时时间，建议 30 分钟
- `delete`: 删除操作的默认超时时间，建议 30 分钟

#### Scenario: Successful resource creation with minimal required parameters

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-12345678"
  template_name = "my-transcode-template"
}
```

- **WHEN** user provides only required parameters (`zone_id` and `template_name`)
- **THEN** the resource creates a template with default values for optional parameters
- **AND** `video_stream_switch` defaults to "on"
- **AND** `audio_stream_switch` defaults to "on"
- **AND** sets the resource ID to `zone_id#template_id`
- **AND** calls `DescribeJustInTimeTranscodeTemplates` to verify template creation

#### Scenario: Successful resource creation with full parameters

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-12345678"
  template_name = "my-transcode-template"
  comment       = "My custom transcode template"
  video_stream_switch = "on"
  audio_stream_switch = "on"
  video_template {
    codec                 = "H.264"
    fps                   = 30
    bitrate               = 2000
    resolution_adaptive   = "open"
    width                 = 1280
    height                = 720
    fill_type             = "black"
  }
  audio_template {
    codec        = "libfdk_aac"
    audio_channel = 2
  }
}
```

- **WHEN** user provides all parameters including video_template and audio_template
- **THEN** the resource creates a template with all specified configurations
- **AND** validates that video_template is provided when video_stream_switch is "on"
- **AND** validates that audio_template is provided when audio_stream_switch is "on"
- **AND** sets the resource ID to `zone_id#template_id`

#### Scenario: Handle missing required parameters

- **WHEN** user omits `zone_id` parameter
- **THEN** the resource returns validation error indicating `zone_id` is required

- **WHEN** user omits `template_name` parameter
- **THEN** the resource returns validation error indicating `template_name` is required

#### Scenario: Handle parameter validation

- **WHEN** user provides invalid `video_stream_switch` value (not "on" or "off")
- **THEN** the resource returns validation error indicating valid values are "on" or "off"

- **WHEN** user provides `video_template` but sets `video_stream_switch` to "off"
- **THEN** the resource returns validation error indicating video_template should not be provided when video_stream_switch is "off"

- **WHEN** user provides `audio_template` but sets `audio_stream_switch` to "off"
- **THEN** the resource returns validation error indicating audio_template should not be provided when audio_stream_switch is "off"

#### Scenario: Handle empty zone_id

- **WHEN** user provides empty `zone_id` parameter
- **THEN** the resource returns validation error indicating `zone_id` cannot be empty

#### Scenario: Handle empty template_name

- **WHEN** user provides empty `template_name` parameter
- **THEN** the resource returns validation error indicating `template_name` cannot be empty

#### Scenario: Handle template_name length exceeding limit

- **WHEN** user provides `template_name` with length exceeding 64 characters
- **THEN** the resource returns validation error indicating template_name length cannot exceed 64 characters

#### Scenario: Handle comment length exceeding limit

- **WHEN** user provides `comment` with length exceeding 256 characters
- **THEN** the resource returns validation error indicating comment length cannot exceed 256 characters

### Requirement: Create Operation

Resource 必须调用 `CreateJustInTimeTranscodeTemplate` API 接口创建即时转码模板。

#### Scenario: Successful API call

- **WHEN** the resource is created with valid parameters
- **THEN** it calls `CreateJustInTimeTranscodeTemplateWithContext` with the provided parameters
- **AND** maps Terraform parameters to API request parameters correctly:
  - `zone_id` → `request.ZoneId`
  - `template_name` → `request.TemplateName`
  - `comment` → `request.Comment`
  - `video_stream_switch` → `request.VideoStreamSwitch`
  - `audio_stream_switch` → `request.AudioStreamSwitch`
  - `video_template` → `request.VideoTemplate`
  - `audio_template` → `request.AudioTemplate`
- **AND** uses `tccommon.WriteRetryTimeout` for retry logic on transient failures
- **AND** sets the resource ID to `zone_id#template_id` after receiving template_id from API response
- **AND** calls Read operation to update state with full template details
- **AND** logs API calls using standard logging patterns with `tccommon.LogElapsed`
- **AND** performs consistency check with `tccommon.InconsistentCheck`

#### Scenario: API response mapping

- **WHEN** the API returns successful response
- **THEN** the resource maps API response parameters to Terraform schema:
  - `response.Response.TemplateId` → `template_id`
  - Other parameters are filled by Read operation
- **AND** stores `template_id` in the resource state

#### Scenario: API rate limiting handling

- **WHEN** API returns rate limiting error (429)
- **THEN** the resource retries the request according to `tccommon.WriteRetryTimeout`
- **AND** logs retry attempts for debugging
- **AND** eventually succeeds if rate limiting is temporary

#### Scenario: API authentication failure

- **WHEN** API returns authentication error due to invalid credentials
- **THEN** the resource returns the original API error to help with troubleshooting
- **AND** does not retry (authentication errors are not transient)

#### Scenario: Invalid zone_id

- **WHEN** API returns error indicating invalid `zone_id`
- **THEN** the resource returns the original API error message
- **AND** helps user identify which parameter is invalid

#### Scenario: Network timeout

- **WHEN** API call times out due to network issues
- **THEN** the resource retries according to retry policy
- **AND** returns timeout error after exhausting retries
- **AND** respects the user-configured `create` timeout

#### Scenario: Async operation polling

- **WHEN** API returns success but template creation is asynchronous
- **THEN** the resource polls `DescribeJustInTimeTranscodeTemplates` to verify template exists
- **AND** waits up to the user-configured `create` timeout
- **AND** uses 5-second polling interval
- **AND** completes successfully when template is found in the list

### Requirement: Read Operation

Resource 必须调用 `DescribeJustInTimeTranscodeTemplates` API 接口读取即时转码模板信息。

#### Scenario: Successful read of existing template

- **WHEN** the resource is read with valid `zone_id` and `template_id`
- **THEN** it calls `DescribeJustInTimeTranscodeTemplatesWithContext` with:
  - `request.ZoneId` set to `zone_id`
  - `request.Filters` containing a filter for `template-id` with value `template_id`
- **AND** maps API response parameters to Terraform schema:
  - `TemplateId` → `template_id`
  - `TemplateName` → `template_name`
  - `Comment` → `comment`
  - `Type` → (computed field, not exposed to user)
  - `VideoStreamSwitch` → `video_stream_switch`
  - `AudioStreamSwitch` → `audio_stream_switch`
  - `VideoTemplate` → `video_template` block with all nested fields
  - `AudioTemplate` → `audio_template` block with all nested fields
  - `CreateTime` → `create_time`
  - `UpdateTime` → `update_time`
- **AND** logs API calls using standard logging patterns with `tccommon.LogElapsed`
- **AND** performs consistency check with `tccommon.InconsistentCheck`

#### Scenario: Read returns template not found

- **WHEN** the resource is read but the template does not exist in the API response
- **THEN** the resource returns nil (resource not found)
- **AND** Terraform will recreate the resource if it still exists in the state

#### Scenario: API rate limiting handling during read

- **WHEN** API returns rate limiting error (429) during read operation
- **THEN** the resource retries the request according to `tccommon.Retry()`
- **AND** logs retry attempts for debugging
- **AND** eventually succeeds if rate limiting is temporary

#### Scenario: Multiple templates with same ID (should not happen)

- **WHEN** the API returns multiple templates with the same template_id
- **THEN** the resource uses the first template in the list
- **AND** logs a warning about potential duplicate template IDs

#### Scenario: Network timeout during read

- **WHEN** API call times out due to network issues during read
- **THEN** the resource retries according to retry policy
- **AND** returns timeout error after exhausting retries

### Requirement: Delete Operation

Resource 必须调用 `DeleteJustInTimeTranscodeTemplates` API 接口删除即时转码模板。

#### Scenario: Successful delete

- **WHEN** the resource is deleted with valid `zone_id` and `template_id`
- **THEN** it calls `DeleteJustInTimeTranscodeTemplatesWithContext` with:
  - `request.ZoneId` set to `zone_id`
  - `request.TemplateIds` containing an array with the `template_id`
- **AND** uses `tccommon.WriteRetryTimeout` for retry logic on transient failures
- **AND** logs API calls using standard logging patterns with `tccommon.LogElapsed`
- **AND** performs consistency check with `tccommon.InconsistentCheck`

#### Scenario: Async delete verification

- **WHEN** API returns success but template deletion is asynchronous
- **THEN** the resource polls `DescribeJustInTimeTranscodeTemplates` to verify template is deleted
- **AND** waits up to the user-configured `delete` timeout
- **AND** uses 5-second polling interval
- **AND** completes successfully when template is no longer found in the list

#### Scenario: Delete template that does not exist

- **WHEN** user attempts to delete a template that does not exist
- **THEN** the resource calls the delete API
- **AND** if API returns success, the delete operation completes
- **AND** if API returns error indicating template not found, the resource removes the resource from state

#### Scenario: API rate limiting handling during delete

- **WHEN** API returns rate limiting error (429) during delete operation
- **THEN** the resource retries the request according to `tccommon.WriteRetryTimeout`
- **AND** logs retry attempts for debugging
- **AND** eventually succeeds if rate limiting is temporary

#### Scenario: Network timeout during delete

- **WHEN** API call times out due to network issues during delete
- **THEN** the resource retries according to retry policy
- **AND** returns timeout error after exhausting retries
- **AND** respects the user-configured `delete` timeout

### Requirement: Update Operation

Resource 必须不支持 Update 操作，所有参数变更触发重新创建（ForceNew）。

#### Scenario: Parameter change triggers recreation

- **WHEN** user changes any parameter (e.g., `template_name`, `comment`, `video_template`, etc.)
- **THEN** Terraform destroys the existing resource instance
- **AND** creates a new resource instance with the updated parameters
- **AND** calls `DeleteJustInTimeTranscodeTemplates` API to delete the old template
- **AND** calls `CreateJustInTimeTranscodeTemplate` API to create a new template

#### Scenario: No direct update method

- **WHEN** user attempts to use Terraform's `terraform apply` to update parameters
- **THEN** the resource performs destroy+create cycle instead of update
- **AND** the new template gets a new `template_id` from the API

#### Scenario: Partial parameter update

- **WHEN** user changes only one parameter (e.g., only `comment`)
- **THEN** the entire resource is recreated with all parameters
- **AND** the old template is deleted and a new template is created with the updated comment

### Requirement: Error Handling and Logging

Resource 必须正确处理错误并提供充分的日志信息。

#### Scenario: API error logging

- **WHEN** an API call fails
- **THEN** the resource logs the error with request and response details
- **AND** includes the API action name for easier debugging
- **AND** returns the error to user

#### Scenario: Operation duration logging

- **WHEN** any CRUD operation is executed
- **THEN** the resource logs the operation duration using `tccommon.LogElapsed`
- **AND** includes the resource name and operation type in log message

#### Scenario: Consistency check

- **WHEN** any CRUD operation is executed
- **THEN** the resource performs consistency check using `tccommon.InconsistentCheck`
- **AND** detects and reports state inconsistencies

#### Scenario: Timeout error handling

- **WHEN** an operation times out (create or delete)
- **THEN** the resource returns a timeout error
- **AND** includes the timeout value and operation type in the error message
- **AND** suggests user to increase the timeout value if needed

### Requirement: Unit Testing

Resource 必须提供完整的单元测试覆盖，使用 mock（gomonkey）方法。

#### Scenario: Test successful creation

- **WHEN** unit test runs with valid mock API response
- **THEN** the test verifies `CreateJustInTimeTranscodeTemplate` is called with correct parameters
- **AND** the test verifies resource ID is set correctly to `zone_id#template_id`
- **AND** the test verifies retry logic on transient errors
- **AND** the test verifies `DescribeJustInTimeTranscodeTemplates` is called for verification

#### Scenario: Test API error handling

- **WHEN** unit test simulates API error during creation
- **THEN** the test verifies error is returned correctly
- **AND** the test verifies error message is preserved

#### Scenario: Test missing required parameters

- **WHEN** unit test runs with missing required parameters
- **THEN** the test verifies validation error is returned
- **AND** the test verifies error message indicates which parameter is missing

#### Scenario: Test parameter validation

- **WHEN** unit test runs with invalid parameter values
- **THEN** the test verifies validation error is returned
- **AND** the test verifies error message indicates which parameter is invalid

#### Scenario: Test read operation

- **WHEN** unit test mocks successful read response
- **THEN** the test verifies `DescribeJustInTimeTranscodeTemplates` is called with correct filters
- **AND** the test verifies response is correctly mapped to schema fields

#### Scenario: Test read not found

- **WHEN** unit test mocks read response with no matching template
- **THEN** the test verifies nil is returned

#### Scenario: Test delete operation

- **WHEN** unit test mocks successful delete response
- **THEN** the test verifies `DeleteJustInTimeTranscodeTemplates` is called with correct parameters
- **AND** the test verifies delete verification polling works

#### Scenario: Test delete with async verification

- **WHEN** unit test simulates async deletion
- **THEN** the test verifies polling logic is called
- **AND** the test verifies deletion is considered successful after template is not found

#### Scenario: Test timeout handling

- **WHEN** unit test simulates timeout during create or delete
- **THEN** the test verifies timeout error is returned
- **AND** the test verifies timeout value is respected

### Requirement: Documentation

Resource 必须提供清晰的使用文档。

#### Scenario: Documentation includes parameters

- **WHEN** users read the resource documentation
- **THEN** they see clear descriptions of all parameters
- **AND** they see default values for optional parameters
- **AND** they see parameter validation rules and constraints
- **AND** they see usage examples showing how to create the resource

#### Scenario: Documentation includes video_template details

- **WHEN** users read the resource documentation
- **THEN** they see detailed descriptions of video_template fields
- **AND** they understand the relationship between video_stream_switch and video_template
- **AND** they see examples of video_template configuration

#### Scenario: Documentation includes audio_template details

- **WHEN** users read the resource documentation
- **THEN** they see detailed descriptions of audio_template fields
- **AND** they understand the relationship between audio_stream_switch and audio_template
- **AND** they see examples of audio_template configuration

#### Scenario: Documentation includes notes

- **WHEN** users read the resource documentation
- **THEN** they see notes about the resource characteristics (ForceNew)
- **AND** they understand that all parameter changes trigger resource recreation
- **AND** they see information about timeout configuration
- **AND** they understand the resource ID format

#### Scenario: Documentation includes examples

- **WHEN** users read the resource documentation
- **THEN** they see complete examples for common use cases:
  - Minimal configuration (only required parameters)
  - Full configuration with all parameters
  - Video-only template
  - Audio-only template
- **AND** examples include explanations of what each parameter does

### Requirement: Resource ID Handling

Resource 必须正确处理复合 ID 格式。

#### Scenario: Parse resource ID

- **WHEN** resource ID is `zone-abc123#tpl-def456`
- **THEN** the resource correctly parses `zone_id` as `zone-abc123`
- **AND** the resource correctly parses `template_id` as `tpl-def456`

#### Scenario: Handle malformed resource ID

- **WHEN** resource ID is malformed (e.g., missing `#` separator)
- **THEN** the resource returns a clear error indicating the ID format is incorrect
- **AND** the error message explains the expected format: `zone_id#template_id`

#### Scenario: Handle zone_id or template_id containing `#`

- **WHEN** `zone_id` or `template_id` contains `#` character
- **THEN** the resource splits on the first `#` only
- **AND** the rest of the string is treated as template_id

### Requirement: Timeouts Configuration

Resource 必须支持 configurable timeouts.

#### Scenario: Default timeouts

- **WHEN** user does not specify custom timeouts
- **THEN** the resource uses default timeout values:
  - `create`: 30 minutes
  - `delete`: 30 minutes

#### Scenario: Custom timeouts

- **WHEN** user specifies custom timeouts in the resource configuration
- **THEN** the resource uses the user-specified timeout values
- **AND** the custom values override the defaults

#### Scenario: Timeout during create

- **WHEN** template creation takes longer than the configured `create` timeout
- **THEN** the resource returns a timeout error
- **AND** the error message includes the configured timeout value
- **AND** suggests increasing the timeout if necessary

#### Scenario: Timeout during delete

- **WHEN** template deletion takes longer than the configured `delete` timeout
- **THEN** the resource returns a timeout error
- **AND** the error message includes the configured timeout value
- **AND** suggests increasing the timeout if necessary
