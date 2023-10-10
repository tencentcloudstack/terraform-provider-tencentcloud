---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_process_live_stream_operation"
sidebar_current: "docs-tencentcloud-resource-mps_process_live_stream_operation"
description: |-
  Provides a resource to create a mps process_live_stream_operation
---

# tencentcloud_mps_process_live_stream_operation

Provides a resource to create a mps process_live_stream_operation

## Example Usage

### Process mps live stream through CMQ

```hcl
resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-process-live-stream-output-${local.app_id}"
  force_clean = true
  acl         = "public-read"
}

resource "tencentcloud_mps_process_live_stream_operation" "operation" {
  url = "http://www.abc.com/abc.m3u8"
  task_notify_config {
    cmq_model   = "Queue"
    cmq_region  = "gz"
    queue_name  = "test"
    topic_name  = "test"
    notify_type = "CMQ"
  }

  output_storage {
    type = "COS"
    cos_output_storage {
      bucket = tencentcloud_cos_bucket.output.bucket
      region = "%s"
    }
  }

  output_dir = "/output/"

  ai_content_review_task {
    definition = 10
  }

  ai_recognition_task {
    definition = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `task_notify_config` - (Required, List, ForceNew) Event notification information of a task, which is used to specify the live stream processing result.
* `url` - (Required, String, ForceNew) Live stream URL, which must be a live stream file address. RTMP, HLS, and FLV are supported.
* `ai_analysis_task` - (Optional, List, ForceNew) AI video intelligent analysis input parameter types.
* `ai_content_review_task` - (Optional, List, ForceNew) Type parameter of a video content audit task.
* `ai_quality_control_task` - (Optional, List, ForceNew) The parameters for a video quality control task.
* `ai_recognition_task` - (Optional, List, ForceNew) Type parameter of video content recognition task.
* `output_dir` - (Optional, String, ForceNew) Target directory of a live stream processing output file, such as `/movie/201909/`. If this parameter is left empty, the `/` directory will be used.
* `output_storage` - (Optional, List, ForceNew) Target bucket of a live stream processing output file. This parameter is required if a file will be output.
* `schedule_id` - (Optional, Int, ForceNew) The scheme ID.Note 1: About `OutputStorage` and `OutputDir`:If an output storage and directory are specified for a subtask of the scheme, those output settings will be applied.If an output storage and directory are not specified for the subtasks of a scheme, the output parameters passed in the `ProcessMedia` API will be applied.Note 2: If `TaskNotifyConfig` is specified, the specified settings will be used instead of the default callback settings of the scheme.
* `session_context` - (Optional, String, ForceNew) The source context which is used to pass through the user request information. The task flow status change callback will return the value of this field. It can contain up to 1,000 characters.
* `session_id` - (Optional, String, ForceNew) The ID used for deduplication. If there was a request with the same ID in the last seven days, the current request will return an error. The ID can contain up to 50 characters. If this parameter is left empty or an empty string is entered, no deduplication will be performed.

The `ai_analysis_task` object supports the following:

* `definition` - (Required, Int) Video content analysis template ID.
* `extended_parameter` - (Optional, String) An extended parameter, whose value is a stringfied JSON.Note: This parameter is for customers with special requirements. It needs to be customized offline.Note: This field may return null, indicating that no valid values can be obtained.

The `ai_content_review_task` object supports the following:

* `definition` - (Required, Int) Video content audit template ID.

The `ai_quality_control_task` object supports the following:

* `channel_ext_para` - (Optional, String) The channel extension parameter, which is a serialized JSON string.Note: This field may return null, indicating that no valid values can be obtained.
* `definition` - (Optional, Int) The ID of the quality control template.Note: This field may return null, indicating that no valid values can be obtained.

The `ai_recognition_task` object supports the following:

* `definition` - (Required, Int) Intelligent video recognition template ID.

The `cos_output_storage` object supports the following:

* `bucket` - (Optional, String) The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
* `region` - (Optional, String) The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.

The `output_storage` object supports the following:

* `type` - (Required, String) The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS.`AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
* `cos_output_storage` - (Optional, List) The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
* `s3_output_storage` - (Optional, List) The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.

The `s3_output_storage` object supports the following:

* `s3_bucket` - (Required, String) The AWS S3 bucket.
* `s3_region` - (Required, String) The region of the AWS S3 bucket.
* `s3_secret_id` - (Optional, String) The key ID required to upload files to the AWS S3 object.
* `s3_secret_key` - (Optional, String) The key required to upload files to the AWS S3 object.

The `task_notify_config` object supports the following:

* `cmq_model` - (Optional, String) CMQ model. There are two types: `Queue` and `Topic`. Currently, only `Queue` is supported.
* `cmq_region` - (Optional, String) CMQ region, such as `sh` and `bj`.
* `notify_type` - (Optional, String) The notification type, `CMQ` by default. If this parameter is set to `URL`, HTTP callbacks are sent to the URL specified by `NotifyUrl`.Note: If you do not pass this parameter or pass in an empty string, `CMQ` will be used. To use a different notification type, specify this parameter accordingly.
* `notify_url` - (Optional, String) HTTP callback URL, required if `NotifyType` is set to `URL`.
* `queue_name` - (Optional, String) This parameter is valid when the model is `Queue`, indicating the name of the CMQ queue for receiving event notifications.
* `topic_name` - (Optional, String) This parameter is valid when the model is `Topic`, indicating the name of the CMQ topic for receiving event notifications.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



