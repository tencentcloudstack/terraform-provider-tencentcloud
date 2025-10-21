---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_media_meta_data"
sidebar_current: "docs-tencentcloud-datasource-mps_media_meta_data"
description: |-
  Use this data source to query detailed information of mps media_meta_data
---

# tencentcloud_mps_media_meta_data

Use this data source to query detailed information of mps media_meta_data

## Example Usage

### Query the mps media meta data through COS

```hcl
data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

data "tencentcloud_mps_media_meta_data" "metadata" {
  input_info {
    type = "COS"
    cos_input_info {
      bucket = data.tencentcloud_cos_bucket_object.object.bucket
      region = "%s"
      object = data.tencentcloud_cos_bucket_object.object.key
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `input_info` - (Required, List) Input information of file for metadata getting.
* `result_output_file` - (Optional, String) Used to save results.

The `cos_input_info` object of `input_info` supports the following:

* `bucket` - (Required, String) The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
* `object` - (Required, String) The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
* `region` - (Required, String) The region of the COS bucket, such as `ap-chongqing`.

The `input_info` object supports the following:

* `type` - (Required, String) The input type. Valid values:`COS`: A COS bucket address.`URL`: A URL.`AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
* `cos_input_info` - (Optional, List) The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
* `s3_input_info` - (Optional, List) The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
* `url_input_info` - (Optional, List) The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.

The `s3_input_info` object of `input_info` supports the following:

* `s3_bucket` - (Required, String) The AWS S3 bucket.
* `s3_object` - (Required, String) The path of the AWS S3 object.
* `s3_region` - (Required, String) The region of the AWS S3 bucket.
* `s3_secret_id` - (Optional, String) The key ID required to access the AWS S3 object.
* `s3_secret_key` - (Optional, String) The key required to access the AWS S3 object.

The `url_input_info` object of `input_info` supports the following:

* `url` - (Required, String) URL of a video.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `meta_data` - Media metadata.
  * `audio_duration` - Audio duration in seconds.Note: This field may return null, indicating that no valid values can be obtained.
  * `audio_stream_set` - Audio stream information.Note: This field may return null, indicating that no valid values can be obtained.
    * `bitrate` - Bitrate of an audio stream in bps.Note: This field may return null, indicating that no valid values can be obtained.
    * `channel` - Number of sound channels, e.g., 2Note: this field may return `null`, indicating that no valid value was found.
    * `codec` - Audio stream codec, such as aac.Note: This field may return null, indicating that no valid values can be obtained.
    * `sampling_rate` - Sample rate of an audio stream in Hz.Note: This field may return null, indicating that no valid values can be obtained.
  * `bitrate` - Sum of the average bitrate of a video stream and that of an audio stream in bps.Note: This field may return null, indicating that no valid values can be obtained.
  * `container` - Container, such as m4a and mp4.Note: This field may return null, indicating that no valid values can be obtained.
  * `duration` - Video duration in seconds.Note: This field may return null, indicating that no valid values can be obtained.
  * `height` - Maximum value of the height of a video stream in px.Note: This field may return null, indicating that no valid values can be obtained.
  * `rotate` - Selected angle during video recording in degrees.Note: This field may return null, indicating that no valid values can be obtained.
  * `size` - Size of an uploaded media file in bytes (which is the sum of size of m3u8 and ts files if the video is in HLS format).Note: This field may return null, indicating that no valid values can be obtained.
  * `video_duration` - Video duration in seconds.Note: This field may return null, indicating that no valid values can be obtained.
  * `video_stream_set` - Video stream information.Note: This field may return null, indicating that no valid values can be obtained.
    * `bitrate` - Bitrate of a video stream in bps.Note: This field may return null, indicating that no valid values can be obtained.
    * `codec` - Video stream codec, such as h264.Note: This field may return null, indicating that no valid values can be obtained.
    * `color_primaries` - Color primariesNote: this field may return `null`, indicating that no valid value was found.
    * `color_space` - Color spaceNote: this field may return `null`, indicating that no valid value was found.
    * `color_transfer` - Color transferNote: this field may return `null`, indicating that no valid value was found.
    * `fps` - Frame rate in Hz.Note: This field may return null, indicating that no valid values can be obtained.
    * `hdr_type` - HDR typeNote: This field may return `null`, indicating that no valid value was found.
    * `height` - Height of a video stream in px.Note: This field may return null, indicating that no valid values can be obtained.
    * `width` - Width of a video stream in px.Note: This field may return null, indicating that no valid values can be obtained.
  * `width` - Maximum value of the width of a video stream in px.Note: This field may return null, indicating that no valid values can be obtained.


