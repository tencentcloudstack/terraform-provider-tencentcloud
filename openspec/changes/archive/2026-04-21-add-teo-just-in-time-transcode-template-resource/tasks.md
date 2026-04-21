## 1. Resource Schema Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.go`
- [x] 1.2 Define resource schema with all required and optional parameters
- [x] 1.3 Set ForceNew for zone_id, template_name, comment; remove ForceNew from video_stream_switch, audio_stream_switch, video_template, audio_template
- [x] 1.4 Add Timeout configuration for create, delete, and read operations (default: 10 minutes)
- [x] 1.5 Define composite ID format as `zone_id#template_id`
- [x] 1.6 Implement resource schema validation rules (template_name max 64 chars, comment max 256 chars)
- [x] 1.7 Implement video_stream_switch and audio_stream_switch validation (only "on" or "off")

## 2. Video Template Schema Implementation

- [x] 2.1 Define video_template as TypeList with MaxItems: 1
- [x] 2.2 Define video_codec field (TypeString, Optional)
- [x] 2.3 Define fps field (TypeFloat, Optional, default 0, range [0, 30])
- [x] 2.4 Define bitrate field (TypeInt, Optional, default 0, range [128, 10000])
- [x] 2.5 Define resolution_adaptive field (TypeString, Optional, default "open")
- [x] 2.6 Define width field (TypeInt, Optional, default 0, range [128, 1920])
- [x] 2.7 Define height field (TypeInt, Optional, default 0, range [128, 1080])
- [x] 2.8 Define fill_type field (TypeString, Optional, default "black")

## 3. Audio Template Schema Implementation

- [x] 3.1 Define audio_template as TypeList with MaxItems: 1
- [x] 3.2 Define codec field (TypeString, Optional)
- [x] 3.3 Define audio_channel field (TypeInt, Optional, default 2)

## 4. Create Function Implementation

- [x] 4.1 Implement resourceTencentCloudTeoJustInTimeTranscodeTemplateCreate function
- [x] 4.2 Parse and validate input parameters from schema
- [x] 4.3 Build CreateJustInTimeTranscodeTemplateRequest from schema data
- [x] 4.4 Call TEO API CreateJustInTimeTranscodeTemplate with retry logic
- [x] 4.5 Handle API errors and return appropriate error messages
- [x] 4.6 Set composite resource ID using tccommon.FILED_SP as separator (`zone_id#template_id`)

## 4a. Update Function Implementation

- [x] 4a.1 Implement resourceTencentCloudTeoJustInTimeTranscodeTemplateUpdate (no-op, logs warning, calls Read)
- [x] 4.7 Set template_id in state
- [x] 4.8 Implement polling mechanism to wait for template creation completion
- [x] 4.9 Call Read function to verify template is available
- [x] 4.10 Add error handling for creation timeout

## 5. Read Function Implementation

- [x] 5.1 Implement resourceTencentCloudTeoJustInTimeTranscodeTemplateRead function
- [x] 5.2 Read zone_id and template_id directly from d.Get() instead of parsing from Id()
- [x] 5.3 Build DescribeJustInTimeTranscodeTemplatesRequest with zone_id
- [x] 5.4 Set filter for template_id to find the specific template
- [x] 5.5 Call TEO API DescribeJustInTimeTranscodeTemplates with retry logic
- [x] 5.6 Handle template not found scenario (return nil, no error)
- [x] 5.7 Parse API response and populate schema fields
- [x] 5.8 Map VideoTemplateInfo from API to video_template schema
- [x] 5.9 Map AudioTemplateInfo from API to audio_template schema
- [x] 5.10 Handle API errors with clear error messages

## 6. Delete Function Implementation

- [x] 6.1 Implement resourceTencentCloudTeoJustInTimeTranscodeTemplateDelete function
- [x] 6.2 Read zone_id and template_id directly from d.Get() instead of parsing from Id()
- [x] 6.3 Build DeleteJustInTimeTranscodeTemplatesRequest with zone_id and template_id
- [x] 6.4 Call TEO API DeleteJustInTimeTranscodeTemplates with retry logic
- [x] 6.5 Handle template not found scenario (ignore error, remove from state)
- [x] 6.6 Implement polling mechanism to wait for template deletion completion
- [x] 6.7 Add error handling for deletion timeout
- [x] 6.8 Clear resource ID from state upon successful deletion

## 7. Resource Registration

- [x] 7.1 Register resource in tencentcloud/provider.go ResourcesMap
- [x] 7.2 Add resource entry in tencentcloud/provider.md under TencentCloud EdgeOne(TEO) Resource section
- [x] 7.3 Ensure resource follows TEO service naming convention

## 8. Helper Functions Implementation

- [x] 8.1 Implement resourceTencentCloudTeoJustInTimeTranscodeTemplateParseId function
- [x] 8.2 Implement function to build VideoTemplateInfo from schema data
- [x] 8.3 Implement function to build AudioTemplateInfo from schema data
- [x] 8.4 Implement function to map API response to video_template schema
- [x] 8.5 Implement function to map API response to audio_template schema
- [x] 8.6 Add defer tccommon.LogElapsed() calls for performance tracking
- [x] 8.7 Add defer tccommon.InconsistentCheck() calls for consistency

## 9. Unit Tests Implementation

- [x] 9.1 Create test file `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template_test.go` (package teo_test)
- [x] 9.2 Add test for Create success scenario with gomonkey mock
- [x] 9.3 Add test for Create API error handling
- [x] 9.4 Add test for Read success scenario
- [x] 9.5 Add test for Read not-found scenario
- [x] 9.6 Add test for Delete success scenario
- [x] 9.7 Add test for Delete API error handling
- [x] 9.8 Add test for Update (no-op that calls Read)
- [x] 9.9 Add test for schema validation (fields, types, Required/Optional/ForceNew/Computed, Update not nil)
- [x] 9.10 Use mockMeta pattern from teo_test package (consistent with identify_zone_operation and create_cls_index_operation tests)

## 10. Resource Documentation

- [x] 10.1 Create documentation file `tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.md`
- [x] 10.2 Follow gendoc convention: one-line description with TEO product name
- [x] 10.3 Add Example Usage section with full HCL configuration
- [x] 10.4 Add Import section with import command example
- [x] 10.5 Argument/Attributes Reference auto-generated from Go schema Description fields

## 11. Verification

- [x] 11.1 Run go build to ensure code compiles
- [x] 11.2 Run unit tests with `go test ./tencentcloud/services/teo/...`
- [x] 11.3 Ensure all tests pass
- [x] 11.4 Verify resource registration in service file
- [x] 11.5 Check for any linting warnings
- [x] 11.6 Review code for consistency with existing TEO resources
- [x] 11.7 Verify all error messages are clear and helpful
- [x] 11.8 Validate documentation examples are accurate
