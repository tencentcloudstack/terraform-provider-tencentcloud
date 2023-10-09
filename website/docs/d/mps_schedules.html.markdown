---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_schedules"
sidebar_current: "docs-tencentcloud-datasource-mps_schedules"
description: |-
  Use this data source to query detailed information of mps schedules
---

# tencentcloud_mps_schedules

Use this data source to query detailed information of mps schedules

## Example Usage

### Query the enabled schedules.

```hcl
data "tencentcloud_mps_schedules" "schedules" {
  status = "Enabled"
}
```

### Query the specified one.

```hcl
data "tencentcloud_mps_schedules" "schedules" {
  schedule_ids = [% d]
  trigger_type = "CosFileUpload"
  status       = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `schedule_ids` - (Optional, Set: [`Int`]) The IDs of the schemes to query. Array length limit: 100.
* `status` - (Optional, String) The scheme status. Valid values:`Enabled`, `Disabled`. If you do not specify this parameter, all schemes will be returned regardless of the status.
* `trigger_type` - (Optional, String) The trigger type. Valid values:`CosFileUpload`: The scheme is triggered when a file is uploaded to Tencent Cloud Object Storage (COS).`AwsS3FileUpload`: The scheme is triggered when a file is uploaded to AWS S3.If you do not specify this parameter or leave it empty, all schemes will be returned regardless of the trigger type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `schedule_info_set` - The information of the schemes.
  * `activities` - The subtasks of the scheme.Note: This field may return null, indicating that no valid values can be obtained.
    * `activity_para` - The parameters of a subtask.Note: This field may return null, indicating that no valid values can be obtained.
      * `adaptive_dynamic_streaming_task` - An adaptive bitrate streaming task.
        * `add_on_subtitles` - The subtitle file to add.Note: This field may return null, indicating that no valid values can be obtained.
          * `subtitle` - The subtitle file.Note: This field may return null, indicating that no valid values can be obtained.
            * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
              * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
              * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
              * `region` - The region of the COS bucket, such as `ap-chongqing`.
            * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
              * `s3_bucket` - The AWS S3 bucket.
              * `s3_object` - The path of the AWS S3 object.
              * `s3_region` - The region of the AWS S3 bucket.
              * `s3_secret_id` - The key ID required to access the AWS S3 object.
              * `s3_secret_key` - The key required to access the AWS S3 object.
            * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
            * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
              * `url` - URL of a video.
          * `type` - The mode. Valid values:`subtitle-stream`: Add a subtitle track.`close-caption-708`: Embed EA-708 subtitles in SEI frames.`close-caption-608`: Embed CEA-608 subtitles in SEI frames.Note: This field may return null, indicating that no valid values can be obtained.
        * `definition` - Adaptive bitrate streaming template ID.
        * `output_object_path` - The relative or absolute output path of the manifest file after being transcoded to adaptive bitrate streaming. If this parameter is left empty, a relative path in the following format will be used by default: `{inputName}_adaptiveDynamicStreaming_{definition}.{format}`.
        * `output_storage` - Target bucket of an output file after being transcoded to adaptive bitrate streaming. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: this field may return null, indicating that no valid values can be obtained.
          * `cos_output_storage` - The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
            * `bucket` - The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
            * `region` - The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.
          * `s3_output_storage` - The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
            * `s3_bucket` - The AWS S3 bucket.
            * `s3_region` - The region of the AWS S3 bucket.
            * `s3_secret_id` - The key ID required to upload files to the AWS S3 object.
            * `s3_secret_key` - The key required to upload files to the AWS S3 object.
          * `type` - The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
        * `segment_object_name` - The relative output path of the segment file after being transcoded to adaptive bitrate streaming (in HLS format only). If this parameter is left empty, a relative path in the following format will be used by default: `{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{format}`.
        * `sub_stream_object_name` - The relative output path of the substream file after being transcoded to adaptive bitrate streaming. If this parameter is left empty, a relative path in the following format will be used by default: `{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{format}`.
        * `watermark_set` - List of up to 10 image or text watermarks.
          * `definition` - ID of a watermarking template.
          * `end_time_offset` - End time offset of a watermark in seconds.If this parameter is left empty or 0 is entered, the watermark will exist till the last video frame;If this value is greater than 0 (e.g., n), the watermark will exist till second n;If this value is smaller than 0 (e.g., -n), the watermark will exist till second n before the last video frame.
          * `raw_parameter` - Custom watermark parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the watermark parameter preferably.Custom watermark parameter is not available for screenshot.
            * `coordinate_origin` - Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text.Default value: TopLeft.
            * `image_template` - Image watermark template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.
              * `height` - Watermark height. % and px formats are supported:If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px.Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.
              * `image_content` - Input content of watermark image. JPEG and PNG images are supported.
                * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
                  * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
                  * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
                  * `region` - The region of the COS bucket, such as `ap-chongqing`.
                * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
                  * `s3_bucket` - The AWS S3 bucket.
                  * `s3_object` - The path of the AWS S3 object.
                  * `s3_region` - The region of the AWS S3 bucket.
                  * `s3_secret_id` - The key ID required to access the AWS S3 object.
                  * `s3_secret_key` - The key required to access the AWS S3 object.
                * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
                * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
                  * `url` - URL of a video.
              * `repeat_type` - Repeat type of an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends.
              * `width` - Watermark width. % and px formats are supported:If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.
            * `type` - Watermark type. Valid values:image: image watermark.
            * `x_pos` - The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.
            * `y_pos` - The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.
          * `start_time_offset` - Start time offset of a watermark in seconds. If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame.If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame;If this value is greater than 0 (e.g., n), the watermark will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the watermark will appear at second n before the last video frame.
          * `svg_content` - SVG content of up to 2,000,000 characters. This field is required only when the watermark type is `SVG`.SVG watermark is not available for screenshot.
          * `text_content` - Text content of up to 100 characters. This field is required only when the watermark type is text.Text watermark is not available for screenshot.
      * `ai_analysis_task` - A content analysis task.
        * `definition` - Video content analysis template ID.
        * `extended_parameter` - An extended parameter, whose value is a stringfied JSON.Note: This parameter is for customers with special requirements. It needs to be customized offline.Note: This field may return null, indicating that no valid values can be obtained.
      * `ai_content_review_task` - A content moderation task.
        * `definition` - Video content audit template ID.
      * `ai_recognition_task` - A content recognition task.
        * `definition` - Intelligent video recognition template ID.
      * `animated_graphic_task` - An animated screenshot generation task.
        * `definition` - Animated image generating template ID.
        * `end_time_offset` - End time of an animated image in a video in seconds.
        * `output_object_path` - Output path to a generated animated image file, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_animatedGraphic_{definition}.{format}`.
        * `output_storage` - Target bucket of a generated animated image file. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.
          * `cos_output_storage` - The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
            * `bucket` - The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
            * `region` - The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.
          * `s3_output_storage` - The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
            * `s3_bucket` - The AWS S3 bucket.
            * `s3_region` - The region of the AWS S3 bucket.
            * `s3_secret_id` - The key ID required to upload files to the AWS S3 object.
            * `s3_secret_key` - The key required to upload files to the AWS S3 object.
          * `type` - The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
        * `start_time_offset` - Start time of an animated image in a video in seconds.
      * `image_sprite_task` - An image sprite generation task.
        * `definition` - ID of an image sprite generating template.
        * `object_number_format` - Rule of the `{number}` variable in the image sprite output path.Note: This field may return null, indicating that no valid values can be obtained.
          * `increment` - Increment of the `{number}` variable. Default value: 1.
          * `initial_value` - Start value of the `{number}` variable. Default value: 0.
          * `min_length` - Minimum length of the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. Default value: 1.
          * `place_holder` - Placeholder used when the `{number}` variable length is below the minimum requirement. Default value: 0.
        * `output_object_path` - Output path to a generated image sprite file, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_imageSprite_{definition}_{number}.{format}`.
        * `output_storage` - Target bucket of a generated image sprite. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.
          * `cos_output_storage` - The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
            * `bucket` - The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
            * `region` - The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.
          * `s3_output_storage` - The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
            * `s3_bucket` - The AWS S3 bucket.
            * `s3_region` - The region of the AWS S3 bucket.
            * `s3_secret_id` - The key ID required to upload files to the AWS S3 object.
            * `s3_secret_key` - The key required to upload files to the AWS S3 object.
          * `type` - The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
        * `web_vtt_object_name` - Output path to the WebVTT file after an image sprite is generated, which can only be a relative path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_imageSprite_{definition}.{format}`.
      * `sample_snapshot_task` - A sampled screencapturing task.
        * `definition` - Sampled screencapturing template ID.
        * `object_number_format` - Rule of the `{number}` variable in the sampled screenshot output path.Note: This field may return null, indicating that no valid values can be obtained.
          * `increment` - Increment of the `{number}` variable. Default value: 1.
          * `initial_value` - Start value of the `{number}` variable. Default value: 0.
          * `min_length` - Minimum length of the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. Default value: 1.
          * `place_holder` - Placeholder used when the `{number}` variable length is below the minimum requirement. Default value: 0.
        * `output_object_path` - Output path to a generated sampled screenshot, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_sampleSnapshot_{definition}_{number}.{format}`.
        * `output_storage` - Target bucket of a sampled screenshot. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.
          * `cos_output_storage` - The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
            * `bucket` - The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
            * `region` - The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.
          * `s3_output_storage` - The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
            * `s3_bucket` - The AWS S3 bucket.
            * `s3_region` - The region of the AWS S3 bucket.
            * `s3_secret_id` - The key ID required to upload files to the AWS S3 object.
            * `s3_secret_key` - The key required to upload files to the AWS S3 object.
          * `type` - The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
        * `watermark_set` - List of up to 10 image or text watermarks.Note: This field may return null, indicating that no valid values can be obtained.
          * `definition` - ID of a watermarking template.
          * `end_time_offset` - End time offset of a watermark in seconds.If this parameter is left empty or 0 is entered, the watermark will exist till the last video frame;If this value is greater than 0 (e.g., n), the watermark will exist till second n;If this value is smaller than 0 (e.g., -n), the watermark will exist till second n before the last video frame.
          * `raw_parameter` - Custom watermark parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the watermark parameter preferably.Custom watermark parameter is not available for screenshot.
            * `coordinate_origin` - Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text.Default value: TopLeft.
            * `image_template` - Image watermark template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.
              * `height` - Watermark height. % and px formats are supported:If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px.Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.
              * `image_content` - Input content of watermark image. JPEG and PNG images are supported.
                * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
                  * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
                  * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
                  * `region` - The region of the COS bucket, such as `ap-chongqing`.
                * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
                  * `s3_bucket` - The AWS S3 bucket.
                  * `s3_object` - The path of the AWS S3 object.
                  * `s3_region` - The region of the AWS S3 bucket.
                  * `s3_secret_id` - The key ID required to access the AWS S3 object.
                  * `s3_secret_key` - The key required to access the AWS S3 object.
                * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
                * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
                  * `url` - URL of a video.
              * `repeat_type` - Repeat type of an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends.
              * `width` - Watermark width. % and px formats are supported:If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.
            * `type` - Watermark type. Valid values:image: image watermark.
            * `x_pos` - The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.
            * `y_pos` - The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.
          * `start_time_offset` - Start time offset of a watermark in seconds. If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame.If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame;If this value is greater than 0 (e.g., n), the watermark will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the watermark will appear at second n before the last video frame.
          * `svg_content` - SVG content of up to 2,000,000 characters. This field is required only when the watermark type is `SVG`.SVG watermark is not available for screenshot.
          * `text_content` - Text content of up to 100 characters. This field is required only when the watermark type is text.Text watermark is not available for screenshot.
      * `snapshot_by_time_offset_task` - A time point screencapturing task.
        * `definition` - ID of a time point screencapturing template.
        * `ext_time_offset_set` - List of screenshot time points in the format of `s` or `%`:If the string ends in `s`, it means that the time point is in seconds; for example, `3.5s` means that the time point is the 3.5th second;If the string ends in `%`, it means that the time point is the specified percentage of the video duration; for example, `10%` means that the time point is 10% of the video duration.
        * `object_number_format` - Rule of the `{number}` variable in the time point screenshot output path.Note: This field may return null, indicating that no valid values can be obtained.
          * `increment` - Increment of the `{number}` variable. Default value: 1.
          * `initial_value` - Start value of the `{number}` variable. Default value: 0.
          * `min_length` - Minimum length of the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. Default value: 1.
          * `place_holder` - Placeholder used when the `{number}` variable length is below the minimum requirement. Default value: 0.
        * `output_object_path` - Output path to a generated time point screenshot, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_snapshotByTimeOffset_{definition}_{number}.{format}`.
        * `output_storage` - Target bucket of a generated time point screenshot file. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.
          * `cos_output_storage` - The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
            * `bucket` - The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
            * `region` - The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.
          * `s3_output_storage` - The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
            * `s3_bucket` - The AWS S3 bucket.
            * `s3_region` - The region of the AWS S3 bucket.
            * `s3_secret_id` - The key ID required to upload files to the AWS S3 object.
            * `s3_secret_key` - The key required to upload files to the AWS S3 object.
          * `type` - The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
        * `watermark_set` - List of up to 10 image or text watermarks.Note: This field may return null, indicating that no valid values can be obtained.
          * `definition` - ID of a watermarking template.
          * `end_time_offset` - End time offset of a watermark in seconds.If this parameter is left empty or 0 is entered, the watermark will exist till the last video frame;If this value is greater than 0 (e.g., n), the watermark will exist till second n;If this value is smaller than 0 (e.g., -n), the watermark will exist till second n before the last video frame.
          * `raw_parameter` - Custom watermark parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the watermark parameter preferably.Custom watermark parameter is not available for screenshot.
            * `coordinate_origin` - Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text.Default value: TopLeft.
            * `image_template` - Image watermark template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.
              * `height` - Watermark height. % and px formats are supported:If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px.Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.
              * `image_content` - Input content of watermark image. JPEG and PNG images are supported.
                * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
                  * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
                  * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
                  * `region` - The region of the COS bucket, such as `ap-chongqing`.
                * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
                  * `s3_bucket` - The AWS S3 bucket.
                  * `s3_object` - The path of the AWS S3 object.
                  * `s3_region` - The region of the AWS S3 bucket.
                  * `s3_secret_id` - The key ID required to access the AWS S3 object.
                  * `s3_secret_key` - The key required to access the AWS S3 object.
                * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
                * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
                  * `url` - URL of a video.
              * `repeat_type` - Repeat type of an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends.
              * `width` - Watermark width. % and px formats are supported:If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.
            * `type` - Watermark type. Valid values:image: image watermark.
            * `x_pos` - The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.
            * `y_pos` - The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.
          * `start_time_offset` - Start time offset of a watermark in seconds. If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame.If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame;If this value is greater than 0 (e.g., n), the watermark will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the watermark will appear at second n before the last video frame.
          * `svg_content` - SVG content of up to 2,000,000 characters. This field is required only when the watermark type is `SVG`.SVG watermark is not available for screenshot.
          * `text_content` - Text content of up to 100 characters. This field is required only when the watermark type is text.Text watermark is not available for screenshot.
      * `transcode_task` - A transcoding task.
        * `definition` - ID of a video transcoding template.
        * `end_time_offset` - End time offset of a transcoded video, in seconds.If this parameter is left empty or set to 0, the transcoded video will end at the same time as the original video.If this parameter is set to a positive number (n for example), the transcoded video will end at the nth second of the original video.If this parameter is set to a negative number (-n for example), the transcoded video will end at the nth second before the end of the original video.
        * `head_tail_parameter` - Opening and closing credits parametersNote: this field may return `null`, indicating that no valid value was found.
          * `head_set` - Opening credits list.
            * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
              * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
              * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
              * `region` - The region of the COS bucket, such as `ap-chongqing`.
            * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
              * `s3_bucket` - The AWS S3 bucket.
              * `s3_object` - The path of the AWS S3 object.
              * `s3_region` - The region of the AWS S3 bucket.
              * `s3_secret_id` - The key ID required to access the AWS S3 object.
              * `s3_secret_key` - The key required to access the AWS S3 object.
            * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
            * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
              * `url` - URL of a video.
          * `tail_set` - Closing credits list.
            * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
              * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
              * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
              * `region` - The region of the COS bucket, such as `ap-chongqing`.
            * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
              * `s3_bucket` - The AWS S3 bucket.
              * `s3_object` - The path of the AWS S3 object.
              * `s3_region` - The region of the AWS S3 bucket.
              * `s3_secret_id` - The key ID required to access the AWS S3 object.
              * `s3_secret_key` - The key required to access the AWS S3 object.
            * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
            * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
              * `url` - URL of a video.
        * `mosaic_set` - List of blurs. Up to 10 ones can be supported.
          * `coordinate_origin` - Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the blur is in the top-left corner of the image or text.Default value: TopLeft.
          * `end_time_offset` - End time offset of blur in seconds.If this parameter is left empty or 0 is entered, the blur will exist till the last video frame;If this value is greater than 0 (e.g., n), the blur will exist till second n;If this value is smaller than 0 (e.g., -n), the blur will exist till second n before the last video frame.
          * `height` - Blur height. % and px formats are supported:If the string ends in %, the `Height` of the blur will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the blur will be in px; for example, `100px` means that `Height` is 100 px.Default value: 10%.
          * `start_time_offset` - Start time offset of blur in seconds. If this parameter is left empty or 0 is entered, the blur will appear upon the first video frame.If this parameter is left empty or 0 is entered, the blur will appear upon the first video frame;If this value is greater than 0 (e.g., n), the blur will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the blur will appear at second n before the last video frame.
          * `width` - Blur width. % and px formats are supported:If the string ends in %, the `Width` of the blur will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the blur will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.
          * `x_pos` - The horizontal position of the origin of the blur relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the blur will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the blur will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.
          * `y_pos` - Vertical position of the origin of blur relative to the origin of coordinates of video. % and px formats are supported:If the string ends in %, the `YPos` of the blur will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the blur will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.
        * `object_number_format` - Rule of the `{number}` variable in the output path after transcoding.Note: This field may return null, indicating that no valid values can be obtained.
          * `increment` - Increment of the `{number}` variable. Default value: 1.
          * `initial_value` - Start value of the `{number}` variable. Default value: 0.
          * `min_length` - Minimum length of the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. Default value: 1.
          * `place_holder` - Placeholder used when the `{number}` variable length is below the minimum requirement. Default value: 0.
        * `output_object_path` - Path to a primary output file, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_transcode_{definition}.{format}`.
        * `output_storage` - Target bucket of an output file. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.
          * `cos_output_storage` - The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
            * `bucket` - The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
            * `region` - The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.
          * `s3_output_storage` - The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
            * `s3_bucket` - The AWS S3 bucket.
            * `s3_region` - The region of the AWS S3 bucket.
            * `s3_secret_id` - The key ID required to upload files to the AWS S3 object.
            * `s3_secret_key` - The key required to upload files to the AWS S3 object.
          * `type` - The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
        * `override_parameter` - Video transcoding custom parameter, which is valid when `Definition` is not 0.When any parameters in this structure are entered, they will be used to override corresponding parameters in templates.This parameter is used in highly customized scenarios. We recommend you only use `Definition` to specify the transcoding parameter.Note: this field may return `null`, indicating that no valid value was found.
          * `add_on_subtitles` - The subtitle file to add.Note: This field may return null, indicating that no valid values can be obtained.
            * `subtitle` - The subtitle file.Note: This field may return null, indicating that no valid values can be obtained.
              * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
                * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
                * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
                * `region` - The region of the COS bucket, such as `ap-chongqing`.
              * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
                * `s3_bucket` - The AWS S3 bucket.
                * `s3_object` - The path of the AWS S3 object.
                * `s3_region` - The region of the AWS S3 bucket.
                * `s3_secret_id` - The key ID required to access the AWS S3 object.
                * `s3_secret_key` - The key required to access the AWS S3 object.
              * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
              * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
                * `url` - URL of a video.
            * `type` - The mode. Valid values:`subtitle-stream`: Add a subtitle track.`close-caption-708`: Embed EA-708 subtitles in SEI frames.`close-caption-608`: Embed CEA-608 subtitles in SEI frames.Note: This field may return null, indicating that no valid values can be obtained.
          * `addon_audio_stream` - The information of the external audio track to add.Note: This field may return null, indicating that no valid values can be obtained.
            * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
              * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
              * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
              * `region` - The region of the COS bucket, such as `ap-chongqing`.
            * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
              * `s3_bucket` - The AWS S3 bucket.
              * `s3_object` - The path of the AWS S3 object.
              * `s3_region` - The region of the AWS S3 bucket.
              * `s3_secret_id` - The key ID required to access the AWS S3 object.
              * `s3_secret_key` - The key required to access the AWS S3 object.
            * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
            * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
              * `url` - URL of a video.
          * `audio_template` - Audio stream configuration parameter.
            * `audio_channel` - Audio channel system. Valid values:1: Mono2: Dual6: StereoWhen the media is packaged in audio format (FLAC, OGG, MP3, M4A), the sound channel cannot be set to stereo.
            * `bitrate` - Audio stream bitrate in Kbps. Value range: 0 and [26, 256]. If the value is 0, the bitrate of the audio stream will be the same as that of the original audio.
            * `codec` - Audio stream codec.When the outer `Container` parameter is `mp3`, the valid value is:libmp3lame.When the outer `Container` parameter is `ogg` or `flac`, the valid value is:flac.When the outer `Container` parameter is `m4a`, the valid values include:libfdk_aac;libmp3lame;ac3.When the outer `Container` parameter is `mp4` or `flv`, the valid values include:libfdk_aac: More suitable for mp4;libmp3lame: More suitable for flv;mp2.When the outer `Container` parameter is `hls`, the valid values include:libfdk_aac;libmp3lame.
            * `sample_rate` - Audio stream sample rate. Valid values:32,00044,10048,000In Hz.
            * `stream_selects` - The audio tracks to retain. All audio tracks are retained by default.
          * `container` - Container format. Valid values: mp4, flv, hls, mp3, flac, ogg, and m4a; mp3, flac, ogg, and m4a are formats of audio files.
          * `remove_audio` - Whether to remove audio data. Valid values:0: retain1: remove.
          * `remove_video` - Whether to remove video data. Valid values:0: retain1: remove.
          * `std_ext_info` - An extended field for transcoding.Note: This field may return null, indicating that no valid values can be obtained.
          * `subtitle_template` - The subtitle settings.Note: This field may return null, indicating that no valid values can be obtained.
            * `font_alpha` - The text transparency. Value range: 0-1.`0`: Fully transparent.`1`: Fully opaque.Default value: 1.Note: This field may return null, indicating that no valid values can be obtained.
            * `font_color` - The font color in 0xRRGGBB format. Default value: 0xFFFFFF (white).Note: This field may return null, indicating that no valid values can be obtained.
            * `font_size` - The font size (pixels). If this is not specified, the font size in the subtitle file will be used.Note: This field may return null, indicating that no valid values can be obtained.
            * `font_type` - The font. Valid values:`hei.ttf`: Heiti.`song.ttf`: Songti.`simkai.ttf`: Kaiti.`arial.ttf`: Arial.The default is `hei.ttf`.Note: This field may return null, indicating that no valid values can be obtained.
            * `path` - The URL of the subtitles to add to the video.Note: This field may return null, indicating that no valid values can be obtained.
            * `stream_index` - The subtitle track to add to the video. If both `Path` and `StreamIndex` are specified, `Path` will be used. You need to specify at least one of the two parameters.Note: This field may return null, indicating that no valid values can be obtained.
          * `tehd_config` - The TSC transcoding parameters.Note: This field may return null, indicating that no valid values can be obtained.
            * `max_video_bitrate` - The maximum video bitrate. If this parameter is not specified, no modifications will be made.Note: This field may return null, indicating that no valid values can be obtained.
            * `type` - The TSC type. Valid values:`TEHD-100`: TSC-100 (video TSC). `TEHD-200`: TSC-200 (audio TSC). If this parameter is left blank, no modification will be made.Note: This field may return null, indicating that no valid values can be obtained.
          * `video_template` - Video stream configuration parameter.
            * `bitrate` - Bitrate of a video stream in Kbps. Value range: 0 and [128, 35,000].If the value is 0, the bitrate of the video will be the same as that of the source video.
            * `codec` - The video codec. Valid values:libx264: H.264libx265: H.265av1: AOMedia Video 1Note: You must specify a resolution (not higher than 640 x 480) if the H.265 codec is used.Note: You can only use the AOMedia Video 1 codec for MP4 files.
            * `content_adapt_stream` - Whether to enable adaptive encoding. Valid values:0: Disable1: EnableDefault value: 0. If this parameter is set to `1`, multiple streams with different resolutions and bitrates will be generated automatically. The highest resolution, bitrate, and quality of the streams are determined by the values of `width` and `height`, `Bitrate`, and `Vcrf` in `VideoTemplate` respectively. If these parameters are not set in `VideoTemplate`, the highest resolution generated will be the same as that of the source video, and the highest video quality will be close to VMAF 95. To use this parameter or learn about the billing details of adaptive encoding, please contact your sales rep.
            * `fill_type` - Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: stretch: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer;black: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks.white: fill with white. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with white color blocks.gauss: fill with Gaussian blur. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with Gaussian blur.
            * `fps` - Video frame rate in Hz. Value range: [0, 100].If the value is 0, the frame rate will be the same as that of the source video.
            * `gop` - Frame interval between I keyframes. Value range: 0 and [1,100000]. If this parameter is 0, the system will automatically set the GOP length.
            * `height` - Maximum value of the height (or short side) of a video stream in px. Value range: 0 and [128, 4,096].
            * `resolution_adaptive` - Resolution adaption. Valid values:open: Enabled. When resolution adaption is enabled, `Width` indicates the long side of a video, while `Height` indicates the short side.close: Disabled. When resolution adaption is disabled, `Width` indicates the width of a video, while `Height` indicates the height.Note: When resolution adaption is enabled, `Width` cannot be smaller than `Height`.
            * `vcrf` - The control factor of video constant bitrate. Value range: [0, 51]. This parameter will be disabled if you enter `0`.It is not recommended to specify this parameter if there are no special requirements.
            * `width` - Maximum value of the width (or long side) of a video stream in px. Value range: 0 and [128, 4,096].If both `Width` and `Height` are 0, the resolution will be the same as that of the source video;If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled;If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled;If both `Width` and `Height` are not 0, the custom resolution will be used.
        * `raw_parameter` - Custom video transcoding parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the transcoding parameter preferably.
          * `audio_template` - Audio stream configuration parameter. This field is required when `RemoveAudio` is 0.
            * `audio_channel` - Audio channel system. Valid values:1: Mono2: Dual6: StereoWhen the media is packaged in audio format (FLAC, OGG, MP3, M4A), the sound channel cannot be set to stereo.Default value: 2.
            * `bitrate` - Audio stream bitrate in Kbps. Value range: 0 and [26, 256].If the value is 0, the bitrate of the audio stream will be the same as that of the original audio.
            * `codec` - Audio stream codec.When the outer `Container` parameter is `mp3`, the valid value is:libmp3lame.When the outer `Container` parameter is `ogg` or `flac`, the valid value is:flac.When the outer `Container` parameter is `m4a`, the valid values include:libfdk_aac;libmp3lame;ac3.When the outer `Container` parameter is `mp4` or `flv`, the valid values include:libfdk_aac: more suitable for mp4;libmp3lame: more suitable for flv.When the outer `Container` parameter is `hls`, the valid values include:libfdk_aac;libmp3lame.
            * `sample_rate` - Audio stream sample rate. Valid values:32,00044,10048,000In Hz.
          * `container` - Container. Valid values: mp4; flv; hls; mp3; flac; ogg; m4a. Among them, mp3, flac, ogg, and m4a are for audio files.
          * `remove_audio` - Whether to remove audio data. Valid values:0: retain;1: remove.Default value: 0.
          * `remove_video` - Whether to remove video data. Valid values:0: retain;1: remove.Default value: 0.
          * `tehd_config` - TESHD transcoding parameter.
            * `max_video_bitrate` - Maximum bitrate, which is valid when `Type` is `TESHD`. If this parameter is left empty or 0 is entered, there will be no upper limit for bitrate.
            * `type` - TESHD type. Valid values:`TEHD-100`: TESHD-100. If this parameter is left empty, TESHD will not be enabled.
          * `video_template` - Video stream configuration parameter. This field is required when `RemoveVideo` is 0.
            * `bitrate` - The video bitrate (Kbps). Value range: 0 and [128, 35000].If the value is 0, the bitrate of the video will be the same as that of the source video.
            * `codec` - The video codec. Valid values:`libx264`: H.264`libx265`: H.265`av1`: AOMedia Video 1Note: You must specify a resolution (not higher than 640 x 480) if the H.265 codec is used.Note: You can only use the AOMedia Video 1 codec for MP4 files.
            * `fill_type` - The fill mode, which indicates how a video is resized when the video&#39;s original aspect ratio is different from the target aspect ratio. Valid values:stretch: Stretch the image frame by frame to fill the entire screen. The video image may become squashed or stretched after transcoding.black: Keep the image&#39;s original aspect ratio and fill the blank space with black bars.white: Keep the image&#39;s original aspect ratio and fill the blank space with white bars.gauss: Keep the image&#39;s original aspect ratio and apply Gaussian blur to the blank space.Default value: black.Note: Only `stretch` and `black` are supported for adaptive bitrate streaming.
            * `fps` - The video frame rate (Hz). Value range: [0, 100].If the value is 0, the frame rate will be the same as that of the source video.Note: For adaptive bitrate streaming, the value range of this parameter is [0, 60].
            * `gop` - Frame interval between I keyframes. Value range: 0 and [1,100000].If this parameter is 0 or left empty, the system will automatically set the GOP length.
            * `height` - Maximum value of the height (or short side) of a video stream in px. Value range: 0 and [128, 4,096].If both `Width` and `Height` are 0, the resolution will be the same as that of the source video;If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled;If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled;If both `Width` and `Height` are not 0, the custom resolution will be used.Default value: 0.
            * `resolution_adaptive` - Resolution adaption. Valid values:open: Enabled. When resolution adaption is enabled, `Width` indicates the long side of a video, while `Height` indicates the short side.close: Disabled. When resolution adaption is disabled, `Width` indicates the width of a video, while `Height` indicates the height.Default value: open.Note: When resolution adaption is enabled, `Width` cannot be smaller than `Height`.
            * `vcrf` - The control factor of video constant bitrate. Value range: [1, 51]If this parameter is specified, CRF (a bitrate control method) will be used for transcoding. (Video bitrate will no longer take effect.)It is not recommended to specify this parameter if there are no special requirements.
            * `width` - Maximum value of the width (or long side) of a video stream in px. Value range: 0 and [128, 4,096].If both `Width` and `Height` are 0, the resolution will be the same as that of the source video;If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled;If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled;If both `Width` and `Height` are not 0, the custom resolution will be used.Default value: 0.
        * `segment_object_name` - Path to an output file part (the path to ts during transcoding to HLS), which can only be a relative path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_transcode_{definition}_{number}.{format}`.
        * `start_time_offset` - Start time offset of a transcoded video, in seconds.If this parameter is left empty or set to 0, the transcoded video will start at the same time as the original video.If this parameter is set to a positive number (n for example), the transcoded video will start at the nth second of the original video.If this parameter is set to a negative number (-n for example), the transcoded video will start at the nth second before the end of the original video.
        * `watermark_set` - List of up to 10 image or text watermarks.Note: This field may return null, indicating that no valid values can be obtained.
          * `definition` - ID of a watermarking template.
          * `end_time_offset` - End time offset of a watermark in seconds.If this parameter is left empty or 0 is entered, the watermark will exist till the last video frame;If this value is greater than 0 (e.g., n), the watermark will exist till second n;If this value is smaller than 0 (e.g., -n), the watermark will exist till second n before the last video frame.
          * `raw_parameter` - Custom watermark parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the watermark parameter preferably.Custom watermark parameter is not available for screenshot.
            * `coordinate_origin` - Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text.Default value: TopLeft.
            * `image_template` - Image watermark template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.
              * `height` - Watermark height. % and px formats are supported:If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px.Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.
              * `image_content` - Input content of watermark image. JPEG and PNG images are supported.
                * `cos_input_info` - The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.
                  * `bucket` - The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.
                  * `object` - The path of the object to process, such as `/movie/201907/WildAnimal.mov`.
                  * `region` - The region of the COS bucket, such as `ap-chongqing`.
                * `s3_input_info` - The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
                  * `s3_bucket` - The AWS S3 bucket.
                  * `s3_object` - The path of the AWS S3 object.
                  * `s3_region` - The region of the AWS S3 bucket.
                  * `s3_secret_id` - The key ID required to access the AWS S3 object.
                  * `s3_secret_key` - The key required to access the AWS S3 object.
                * `type` - The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.
                * `url_input_info` - The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.
                  * `url` - URL of a video.
              * `repeat_type` - Repeat type of an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends.
              * `width` - Watermark width. % and px formats are supported:If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.
            * `type` - Watermark type. Valid values:image: image watermark.
            * `x_pos` - The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.
            * `y_pos` - The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.
          * `start_time_offset` - Start time offset of a watermark in seconds. If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame.If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame;If this value is greater than 0 (e.g., n), the watermark will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the watermark will appear at second n before the last video frame.
          * `svg_content` - SVG content of up to 2,000,000 characters. This field is required only when the watermark type is `SVG`.SVG watermark is not available for screenshot.
          * `text_content` - Text content of up to 100 characters. This field is required only when the watermark type is text.Text watermark is not available for screenshot.
    * `activity_type` - The subtask type.`input`: The start.`output`: The end.`action-trans`: Transcoding.`action-samplesnapshot`: Sampled screencapturing.`action-AIAnalysis`: Content analysis.`action-AIRecognition`: Content recognition.`action-aiReview`: Content moderation.`action-animated-graphics`: Animated screenshot generation.`action-image-sprite`: Image sprite generation.`action-snapshotByTimeOffset`: Time point screencapturing.`action-adaptive-substream`: Adaptive bitrate streaming.Note: This field may return null, indicating that no valid values can be obtained.
    * `reardrive_index` - The indexes of the subsequent actions.Note: This field may return null, indicating that no valid values can be obtained.
  * `create_time` - The creation time in [ISO date format](https://intl.cloud.tencent.com/document/product/862/37710?from_cn_redirect=1#52).Note: This field may return null, indicating that no valid values can be obtained.
  * `output_dir` - The directory to save the output file.Note: This field may return null, indicating that no valid values can be obtained.
  * `output_storage` - The bucket to save the output file.Note: This field may return null, indicating that no valid values can be obtained.
    * `cos_output_storage` - The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.
      * `bucket` - The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.
      * `region` - The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.
    * `s3_output_storage` - The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.
      * `s3_bucket` - The AWS S3 bucket.
      * `s3_region` - The region of the AWS S3 bucket.
      * `s3_secret_id` - The key ID required to upload files to the AWS S3 object.
      * `s3_secret_key` - The key required to upload files to the AWS S3 object.
    * `type` - The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.
  * `schedule_id` - The scheme ID.
  * `schedule_name` - The scheme name.Note: This field may return null, indicating that no valid values can be obtained.
  * `status` - The scheme status. Valid values:`Enabled``Disabled`Note: This field may return null, indicating that no valid values can be obtained.
  * `task_notify_config` - The notification configuration.Note: This field may return null, indicating that no valid values can be obtained.
    * `aws_sqs` - The AWS SQS queue. This parameter is required if `NotifyType` is `AWS-SQS`.Note: This field may return null, indicating that no valid values can be obtained.
      * `s3_secret_id` - The key ID required to read from/write to the SQS queue.
      * `s3_secret_key` - The key required to read from/write to the SQS queue.
      * `sqs_queue_name` - The name of the SQS queue.
      * `sqs_region` - The region of the SQS queue.
    * `cmq_model` - The CMQ or TDMQ-CMQ model. Valid values: Queue, Topic.
    * `cmq_region` - The CMQ or TDMQ-CMQ region, such as `sh` (Shanghai) or `bj` (Beijing).
    * `notify_mode` - Workflow notification method. Valid values: Finish, Change. If this parameter is left empty, `Finish` will be used.
    * `notify_type` - The notification type. Valid values:`CMQ`: This value is no longer used. Please use `TDMQ-CMQ` instead.`TDMQ-CMQ`: Message queue`URL`: If `NotifyType` is set to `URL`, HTTP callbacks are sent to the URL specified by `NotifyUrl`. HTTP and JSON are used for the callbacks. The packet contains the response parameters of the `ParseNotification` API.`SCF`: This notification type is not recommended. You need to configure it in the SCF console.`AWS-SQS`: AWS queue. This type is only supported for AWS tasks, and the queue must be in the same region as the AWS bucket.Note: If you do not pass this parameter or pass in an empty string, `CMQ` will be used. To use a different notification type, specify this parameter accordingly.
    * `notify_url` - HTTP callback URL, required if `NotifyType` is set to `URL`.
    * `queue_name` - The CMQ or TDMQ-CMQ queue to receive notifications. This parameter is valid when `CmqModel` is `Queue`.
    * `topic_name` - The CMQ or TDMQ-CMQ topic to receive notifications. This parameter is valid when `CmqModel` is `Topic`.
  * `trigger` - The trigger of the scheme.Note: This field may return null, indicating that no valid values can be obtained.
    * `aws_s3_file_upload_trigger` - The AWS S3 trigger. This parameter is valid and required if `Type` is `AwsS3FileUpload`.Note: Currently, the key for the AWS S3 bucket, the trigger SQS queue, and the callback SQS queue must be the same.Note: This field may return null, indicating that no valid values can be obtained.
      * `aws_sqs` - The SQS queue of the AWS S3 bucket.Note: The queue must be in the same region as the bucket.Note: This field may return null, indicating that no valid values can be obtained.
        * `s3_secret_id` - The key ID required to read from/write to the SQS queue.
        * `s3_secret_key` - The key required to read from/write to the SQS queue.
        * `sqs_queue_name` - The name of the SQS queue.
        * `sqs_region` - The region of the SQS queue.
      * `dir` - The bucket directory bound. It must be an absolute path that starts and ends with `/`, such as `/movie/201907/`. If you do not specify this, the root directory will be bound.	.
      * `formats` - The file formats that will trigger the scheme, such as [mp4, flv, mov]. If you do not specify this, the upload of files in any format will trigger the scheme.	.
      * `s3_bucket` - The AWS S3 bucket bound to the scheme.
      * `s3_region` - The region of the AWS S3 bucket.
      * `s3_secret_id` - The key ID of the AWS S3 bucket.Note: This field may return null, indicating that no valid values can be obtained.
      * `s3_secret_key` - The key of the AWS S3 bucket.Note: This field may return null, indicating that no valid values can be obtained.
    * `cos_file_upload_trigger` - This parameter is required and valid when `Type` is `CosFileUpload`, indicating the COS trigger rule.Note: This field may return null, indicating that no valid values can be obtained.
      * `bucket` - Name of the COS bucket bound to a workflow, such as `TopRankVideo-125xxx88`.
      * `dir` - Input path directory bound to a workflow, such as `/movie/201907/`. If this parameter is left empty, the `/` root directory will be used.
      * `formats` - Format list of files that can trigger a workflow, such as [mp4, flv, mov]. If this parameter is left empty, files in all formats can trigger the workflow.
      * `region` - Region of the COS bucket bound to a workflow, such as `ap-chongiqng`.
    * `type` - The trigger type. Valid values:`CosFileUpload`: Tencent Cloud COS trigger.`AwsS3FileUpload`: AWS S3 trigger. Currently, this type is only supported for transcoding tasks and schemes (not supported for workflows).
  * `update_time` - The last updated time in [ISO date format](https://intl.cloud.tencent.com/document/product/862/37710?from_cn_redirect=1#52).Note: This field may return null, indicating that no valid values can be obtained.


