---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_edit_media_operation"
sidebar_current: "docs-tencentcloud-resource-mps_edit_media_operation"
description: |-
  Provides a resource to create a mps edit_media_operation
---

# tencentcloud_mps_edit_media_operation

Provides a resource to create a mps edit_media_operation

## Example Usage

### Operation through COS

```hcl
resource "tencentcloud_cos_bucket" "output" {
  bucket = "tf-bucket-mps-output-${local.app_id}"
}

data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

resource "tencentcloud_mps_edit_media_operation" "operation" {
  file_infos {
    input_info {
      type = "COS"
      cos_input_info {
        bucket = data.tencentcloud_cos_bucket_object.object.bucket
        region = "%s"
        object = data.tencentcloud_cos_bucket_object.object.key
      }
    }
    start_time_offset = 60
    end_time_offset   = 120
  }
  output_storage {
    type = "COS"
    cos_output_storage {
      bucket = tencentcloud_cos_bucket.output.bucket
      region = "%s"
    }
  }
  output_object_path = "/output"
}
```

## Argument Reference

The following arguments are supported:

* `file_infos` - (Required, List, ForceNew) Information of input video file.
* `output_object_path` - (Required, String, ForceNew) The path to save the media processing output file.
* `output_storage` - (Required, List, ForceNew) The storage location of the media processing output file.
* `output_config` - (Optional, List, ForceNew) Configuration for output files of video editing.
* `session_context` - (Optional, String, ForceNew) The source context which is used to pass through the user request information. The task flow status change callback will return the value of this field. It can contain up to 1,000 characters.
* `session_id` - (Optional, String, ForceNew) The ID used for deduplication. If there was a request with the same ID in the last three days, the current request will return an error. The ID can contain up to 50 characters. If this parameter is left empty or an empty string is entered, no deduplication will be performed.
* `task_notify_config` - (Optional, List, ForceNew) Event notification information of task. If this parameter is left empty, no event notifications will be obtained.
* `tasks_priority` - (Optional, Int, ForceNew) Task priority. The higher the value, the higher the priority. Value range: [-10,10]. If this parameter is left empty, 0 will be used.

The `aws_sqs` object of `task_notify_config` supports the following:

* `sqs_queue_name` - (Required, String) The name of the SQS queue.
* `sqs_region` - (Required, String) The region of the SQS queue.
* `s3_secret_id` - (Optional, String) The key ID required to read from/write to the SQS queue.
* `s3_secret_key` - (Optional, String) The key required to read from/write to the SQS queue.

The `cos_input_info` object of `input_info` supports the following:

* `bucket` - (Required, String) The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
* `object` - (Required, String) The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
* `region` - (Required, String) The region of the COS bucket, such as `ap-chongqing`.

The `cos_output_storage` object of `output_storage` supports the following:

* `bucket` - (Optional, String) The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
* `region` - (Optional, String) The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.

The `file_infos` object supports the following:

* `input_info` - (Required, List) Video input information.
* `end_time_offset` - (Optional, Float64) End time offset of video clipping in seconds.
* `start_time_offset` - (Optional, Float64) Start time offset of video clipping in seconds.

The `input_info` object of `file_infos` supports the following:

* `type` - (Required, String) The input type. Valid values: `COS`: A COS bucket address.  `URL`: A URL.  `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
* `cos_input_info` - (Optional, List) The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
* `s3_input_info` - (Optional, List) The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
* `url_input_info` - (Optional, List) The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.

The `output_config` object supports the following:

* `container` - (Optional, String) Format. Valid values: `mp4` (default), `hls`, `mov`, `flv`, `avi`.
* `type` - (Optional, String) The editing mode. Valid values are `normal` and `fast`. The default is `normal`, which indicates precise editing.

The `output_storage` object supports the following:

* `type` - (Required, String) The storage type for a media processing output file. Valid values: `COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
* `cos_output_storage` - (Optional, List) The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
* `s3_output_storage` - (Optional, List) The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.

The `s3_input_info` object of `input_info` supports the following:

* `s3_bucket` - (Required, String) The AWS S3 bucket.
* `s3_object` - (Required, String) The path of the AWS S3 object.
* `s3_region` - (Required, String) The region of the AWS S3 bucket.
* `s3_secret_id` - (Optional, String) The key ID required to access the AWS S3 object.
* `s3_secret_key` - (Optional, String) The key required to access the AWS S3 object.

The `s3_output_storage` object of `output_storage` supports the following:

* `s3_bucket` - (Required, String) The AWS S3 bucket.
* `s3_region` - (Required, String) The region of the AWS S3 bucket.
* `s3_secret_id` - (Optional, String) The key ID required to upload files to the AWS S3 object.
* `s3_secret_key` - (Optional, String) The key required to upload files to the AWS S3 object.

The `task_notify_config` object supports the following:

* `aws_sqs` - (Optional, List) The AWS SQS queue. This parameter is required if `NotifyType` is `AWS-SQS`.Note: This field may return null, indicating that no valid values can be obtained.
* `cmq_model` - (Optional, String) The CMQ or TDMQ-CMQ model. Valid values: Queue, Topic.
* `cmq_region` - (Optional, String) The CMQ or TDMQ-CMQ region, such as `sh` (Shanghai) or `bj` (Beijing).
* `notify_mode` - (Optional, String) Workflow notification method. Valid values: Finish, Change. If this parameter is left empty, `Finish` will be used.
* `notify_type` - (Optional, String) The notification type. Valid values: `CMQ`: This value is no longer used. Please use `TDMQ-CMQ` instead. `TDMQ-CMQ`: Message queue. `URL`: If `NotifyType` is set to `URL`, HTTP callbacks are sent to the URL specified by `NotifyUrl`. HTTP and JSON are used for the callbacks. The packet contains the response parameters of the `ParseNotification` API. `SCF`: This notification type is not recommended. You need to configure it in the SCF console. `AWS-SQS`: AWS queue. This type is only supported for AWS tasks, and the queue must be in the same region as the AWS bucket. If you do not pass this parameter or pass in an empty string, `CMQ` will be used. To use a different notification type, specify this parameter accordingly.
* `notify_url` - (Optional, String) HTTP callback URL, required if `NotifyType` is set to `URL`.
* `queue_name` - (Optional, String) The CMQ or TDMQ-CMQ queue to receive notifications. This parameter is valid when `CmqModel` is `Queue`.
* `topic_name` - (Optional, String) The CMQ or TDMQ-CMQ topic to receive notifications. This parameter is valid when `CmqModel` is `Topic`.

The `url_input_info` object of `input_info` supports the following:

* `url` - (Required, String) URL of a video.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



