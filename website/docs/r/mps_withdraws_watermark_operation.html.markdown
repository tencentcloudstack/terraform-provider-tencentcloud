---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_withdraws_watermark_operation"
sidebar_current: "docs-tencentcloud-resource-mps_withdraws_watermark_operation"
description: |-
  Provides a resource to create a mps withdraws_watermark_operation
---

# tencentcloud_mps_withdraws_watermark_operation

Provides a resource to create a mps withdraws_watermark_operation

## Example Usage

### Withdraw the watermark from COS

```hcl
resource "tencentcloud_cos_bucket" "example" {
  bucket = "tf-test-mps-wm-${local.app_id}"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "example" {
  bucket = tencentcloud_cos_bucket.example.bucket
  key    = "/test-file/test.mov"
  source = "/Users/luoyin/Downloads/file_example_MOV_480_700kB.mov"
}

resource "tencentcloud_mps_withdraws_watermark_operation" "operation" {
  input_info {
    type = "COS"
    cos_input_info {
      bucket = tencentcloud_cos_bucket_object.example.bucket
      region = "%s"
      object = tencentcloud_cos_bucket_object.example.key
    }
  }

  session_context = "this is a example session context"
}
```

## Argument Reference

The following arguments are supported:

* `input_info` - (Required, List, ForceNew) Input information of file for metadata getting.
* `session_context` - (Optional, String, ForceNew) The source context which is used to pass through the user request information. The task flow status change callback will return the value of this field.
* `task_notify_config` - (Optional, List, ForceNew) Event notification information of a task. If this parameter is left empty, no event notifications will be obtained.

The `aws_sqs` object of `task_notify_config` supports the following:

* `sqs_queue_name` - (Required, String) The name of the SQS queue.
* `sqs_region` - (Required, String) The region of the SQS queue.
* `s3_secret_id` - (Optional, String) The key ID required to read from/write to the SQS queue.
* `s3_secret_key` - (Optional, String) The key required to read from/write to the SQS queue.

The `cos_input_info` object of `input_info` supports the following:

* `bucket` - (Required, String) The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
* `object` - (Required, String) The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
* `region` - (Required, String) The region of the COS bucket, such as `ap-chongqing`.

The `input_info` object supports the following:

* `type` - (Required, String) The input type. Valid values: `COS`: A COS bucket address.  `URL`: A URL.  `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks..
* `cos_input_info` - (Optional, List) The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
* `s3_input_info` - (Optional, List) The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
* `url_input_info` - (Optional, List) The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.

The `s3_input_info` object of `input_info` supports the following:

* `s3_bucket` - (Required, String) The AWS S3 bucket.
* `s3_object` - (Required, String) The path of the AWS S3 object.
* `s3_region` - (Required, String) The region of the AWS S3 bucket.
* `s3_secret_id` - (Optional, String) The key ID required to access the AWS S3 object.
* `s3_secret_key` - (Optional, String) The key required to access the AWS S3 object.

The `task_notify_config` object supports the following:

* `aws_sqs` - (Optional, List) The AWS SQS queue. This parameter is required if `NotifyType` is `AWS-SQS`.Note: This field may return null, indicating that no valid values can be obtained.
* `cmq_model` - (Optional, String) The CMQ or TDMQ-CMQ model. Valid values: Queue, Topic.
* `cmq_region` - (Optional, String) The CMQ or TDMQ-CMQ region, such as `sh` (Shanghai) or `bj` (Beijing).
* `notify_mode` - (Optional, String) Workflow notification method. Valid values: Finish, Change. If this parameter is left empty, `Finish` will be used.
* `notify_type` - (Optional, String) The notification type. Valid values:  `CMQ`: This value is no longer used. Please use `TDMQ-CMQ` instead.  `TDMQ-CMQ`: Message queue  `URL`: If `NotifyType` is set to `URL`, HTTP callbacks are sent to the URL specified by `NotifyUrl`. HTTP and JSON are used for the callbacks. The packet contains the response parameters of the `ParseNotification` API.  `SCF`: This notification type is not recommended. You need to configure it in the SCF console.  `AWS-SQS`: AWS queue. This type is only supported for AWS tasks, and the queue must be in the same region as the AWS bucket. Note: If you do not pass this parameter or pass in an empty string, `CMQ` will be used. To use a different notification type, specify this parameter accordingly.
* `notify_url` - (Optional, String) HTTP callback URL, required if `NotifyType` is set to `URL`.
* `queue_name` - (Optional, String) The CMQ or TDMQ-CMQ queue to receive notifications. This parameter is valid when `CmqModel` is `Queue`.
* `topic_name` - (Optional, String) The CMQ or TDMQ-CMQ topic to receive notifications. This parameter is valid when `CmqModel` is `Topic`.

The `url_input_info` object of `input_info` supports the following:

* `url` - (Required, String) URL of a video.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



