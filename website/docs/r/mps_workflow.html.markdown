---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_workflow"
sidebar_current: "docs-tencentcloud-resource-mps_workflow"
description: |-
  Provides a resource to create a mps workflow
---

# tencentcloud_mps_workflow

Provides a resource to create a mps workflow

## Example Usage

```hcl
resource "tencentcloud_mps_workflow" "workflow" {
  workflow_name = & lt ; nil & gt ;
  trigger {
    type = "CosFileUpload"
    cos_file_upload_trigger {
      bucket  = "TopRankVideo-125xxx88"
      region  = "ap-chongqing"
      dir     = "/movie/201907/"
      formats =
    }

  }
  output_storage {
    type = "COS"
    cos_output_storage {
      bucket = "TopRankVideo-125xxx88"
      region = "ap-chongqing"
    }

  }
  output_dir = "/movie/201907/"
  media_process_task {
    transcode_task_set {
      definition = & lt ; nil & gt ;
      raw_parameter {
        container    = & lt ; nil & gt ;
        remove_video = 0
        remove_audio = 0
        video_template {
          codec               = & lt ; nil & gt ;
          fps                 = & lt ; nil & gt ;
          bitrate             = & lt ; nil & gt ;
          resolution_adaptive = "open"
          width               = 0
          height              = 0
          gop                 = & lt ; nil & gt ;
          fill_type           = "black"
          vcrf                = & lt ; nil & gt ;
        }
        audio_template {
          codec         = & lt ; nil & gt ;
          bitrate       = & lt ; nil & gt ;
          sample_rate   = & lt ; nil & gt ;
          audio_channel = 2
        }
        t_e_h_d_config {
          type              = & lt ; nil & gt ;
          max_video_bitrate = & lt ; nil & gt ;
        }
      }
      override_parameter {
        container    = & lt ; nil & gt ;
        remove_video = & lt ; nil & gt ;
        remove_audio = & lt ; nil & gt ;
        video_template {
          codec                = & lt ; nil & gt ;
          fps                  = & lt ; nil & gt ;
          bitrate              = & lt ; nil & gt ;
          resolution_adaptive  = & lt ; nil & gt ;
          width                = & lt ; nil & gt ;
          height               = & lt ; nil & gt ;
          gop                  = & lt ; nil & gt ;
          fill_type            = & lt ; nil & gt ;
          vcrf                 = & lt ; nil & gt ;
          content_adapt_stream = 0
        }
        audio_template {
          codec          = & lt ; nil & gt ;
          bitrate        = & lt ; nil & gt ;
          sample_rate    = & lt ; nil & gt ;
          audio_channel  = & lt ; nil & gt ;
          stream_selects = & lt ; nil & gt ;
        }
        t_e_h_d_config {
          type              = & lt ; nil & gt ;
          max_video_bitrate = & lt ; nil & gt ;
        }
        subtitle_template {
          path         = & lt ; nil & gt ;
          stream_index = & lt ; nil & gt ;
          font_type    = "hei.ttf"
          font_size    = & lt ; nil & gt ;
          font_color   = "0xFFFFFF"
          font_alpha   =
        }
      }
      watermark_set {
        definition = & lt ; nil & gt ;
        raw_parameter {
          type              = & lt ; nil & gt ;
          coordinate_origin = "TopLeft"
          x_pos             = "0px"
          y_pos             = "0px"
          image_template {
            image_content {
              type = "COS"
              cos_input_info {
                bucket = "TopRankVideo-125xxx88"
                region = "ap-chongqing"
                object = "/movie/201907/WildAnimal.mov"
              }
              url_input_info {
                url = & lt ; nil & gt ;
              }
            }
            width       = "10%"
            height      = "0px"
            repeat_type = & lt ; nil & gt ;
          }
        }
        text_content      = & lt ; nil & gt ;
        svg_content       = & lt ; nil & gt ;
        start_time_offset = & lt ; nil & gt ;
        end_time_offset   = & lt ; nil & gt ;
      }
      mosaic_set {
        coordinate_origin = "TopLeft"
        x_pos             = "0px"
        y_pos             = "0px"
        width             = "10%"
        height            = "10%"
        start_time_offset = & lt ; nil & gt ;
        end_time_offset   = & lt ; nil & gt ;
      }
      start_time_offset = & lt ; nil & gt ;
      end_time_offset   = & lt ; nil & gt ;
      output_storage {
        type = "COS"
        cos_output_storage {
          bucket = "TopRankVideo-125xxx88"
          region = "ap-chongqinq"
        }
      }
      output_object_path  = & lt ; nil & gt ;
      segment_object_name = & lt ; nil & gt ;
      object_number_format {
        initial_value = 0
        increment     = 1
        min_length    = 1
        place_holder  = "0"
      }
      head_tail_parameter {
        head_set {
          type = "COS"
          cos_input_info {
            bucket = "TopRankVideo-125xxx88"
            region = "ap-chongqing"
            object = "/movie/201907/WildAnimal.mov"
          }
          url_input_info {
            url = & lt ; nil & gt ;
          }
        }
        tail_set {
          type = "COS"
          cos_input_info {
            bucket = "TopRankVideo-125xxx88"
            region = "ap-chongqing"
            object = "/movie/201907/WildAnimal.mov"
          }
          url_input_info {
            url = & lt ; nil & gt ;
          }
        }
      }
    }
    animated_graphic_task_set {
      definition        = & lt ; nil & gt ;
      start_time_offset = & lt ; nil & gt ;
      end_time_offset   = & lt ; nil & gt ;
      output_storage {
        type = "COS"
        cos_output_storage {
          bucket = "TopRankVideo-125xxx88"
          region = "ap-chongqinq"
        }
      }
      output_object_path = & lt ; nil & gt ;
    }
    snapshot_by_time_offset_task_set {
      definition          = & lt ; nil & gt ;
      ext_time_offset_set = & lt ; nil & gt ;
      time_offset_set     = & lt ; nil & gt ;
      watermark_set {
        definition = & lt ; nil & gt ;
        raw_parameter {
          type              = & lt ; nil & gt ;
          coordinate_origin = "TopLeft"
          x_pos             = "0px"
          y_pos             = "0px"
          image_template {
            image_content {
              type = "COS"
              cos_input_info {
                bucket = "TopRankVideo-125xxx88"
                region = "ap-chongqing"
                object = "/movie/201907/WildAnimal.mov"
              }
              url_input_info {
                url = & lt ; nil & gt ;
              }
            }
            width       = "10%"
            height      = "0px"
            repeat_type = & lt ; nil & gt ;
          }
        }
        text_content      = & lt ; nil & gt ;
        svg_content       = & lt ; nil & gt ;
        start_time_offset = & lt ; nil & gt ;
        end_time_offset   = & lt ; nil & gt ;
      }
      output_storage {
        type = "COS"
        cos_output_storage {
          bucket = "TopRankVideo-125xxx88"
          region = "ap-chongqinq"
        }
      }
      output_object_path = & lt ; nil & gt ;
      object_number_format {
        initial_value = 0
        increment     = 1
        min_length    = 1
        place_holder  = "0"
      }
    }
    sample_snapshot_task_set {
      definition = & lt ; nil & gt ;
      watermark_set {
        definition = & lt ; nil & gt ;
        raw_parameter {
          type              = & lt ; nil & gt ;
          coordinate_origin = "TopLeft"
          x_pos             = "0px"
          y_pos             = "0px"
          image_template {
            image_content {
              type = "COS"
              cos_input_info {
                bucket = "TopRankVideo-125xxx88"
                region = "ap-chongqing"
                object = "/movie/201907/WildAnimal.mov"
              }
              url_input_info {
                url = & lt ; nil & gt ;
              }
            }
            width       = "10%"
            height      = "0px"
            repeat_type = "repeat"
          }
        }
        text_content      = & lt ; nil & gt ;
        svg_content       = & lt ; nil & gt ;
        start_time_offset = & lt ; nil & gt ;
        end_time_offset   = & lt ; nil & gt ;
      }
      output_storage {
        type = "COS"
        cos_output_storage {
          bucket = "TopRankVideo-125xxx88"
          region = "ap-chongqinq"
        }
      }
      output_object_path = & lt ; nil & gt ;
      object_number_format {
        initial_value = 0
        increment     = 1
        min_length    = 1
        place_holder  = "0"
      }
    }
    image_sprite_task_set {
      definition = & lt ; nil & gt ;
      output_storage {
        type = "COS"
        cos_output_storage {
          bucket = "TopRankVideo-125xxx88"
          region = "ap-chongqinq"
        }
      }
      output_object_path  = & lt ; nil & gt ;
      web_vtt_object_name = & lt ; nil & gt ;
      object_number_format {
        initial_value = 0
        increment     = 1
        min_length    = 1
        place_holder  = "0"
      }
    }
    adaptive_dynamic_streaming_task_set {
      definition = & lt ; nil & gt ;
      watermark_set {
        definition = & lt ; nil & gt ;
        raw_parameter {
          type              = & lt ; nil & gt ;
          coordinate_origin = "TopLeft"
          x_pos             = "0px"
          y_pos             = "0px"
          image_template {
            image_content {
              type = "COS"
              cos_input_info {
                bucket = "TopRankVideo-125xxx88"
                region = "ap-chongqing"
                object = "/movie/201907/WildAnimal.mov"
              }
              url_input_info {
                url = & lt ; nil & gt ;
              }
            }
            width       = "10%"
            height      = "0px"
            repeat_type = "repeat"
          }
        }
        text_content      = & lt ; nil & gt ;
        svg_content       = & lt ; nil & gt ;
        start_time_offset = & lt ; nil & gt ;
        end_time_offset   = & lt ; nil & gt ;
      }
      output_storage {
        type = "COS"
        cos_output_storage {
          bucket = "TopRankVideo-125xxx88"
          region = "ap-chongqinq"
        }
      }
      output_object_path     = & lt ; nil & gt ;
      sub_stream_object_name = & lt ; nil & gt ;
      segment_object_name    = & lt ; nil & gt ;
    }

  }
  ai_content_review_task {
    definition = & lt ; nil & gt ;

  }
  ai_analysis_task {
    definition         = & lt ; nil & gt ;
    extended_parameter = & lt ; nil & gt ;

  }
  ai_recognition_task {
    definition = & lt ; nil & gt ;

  }
  task_notify_config {
    cmq_model   = & lt ; nil & gt ;
    cmq_region  = & lt ; nil & gt ;
    topic_name  = & lt ; nil & gt ;
    queue_name  = & lt ; nil & gt ;
    notify_mode = & lt ; nil & gt ;
    notify_type = & lt ; nil & gt ;
    notify_url  = & lt ; nil & gt ;

  }
  task_priority = & lt ; nil & gt ;
}
```

## Argument Reference

The following arguments are supported:

* `trigger` - (Required, List) The trigger rule bound to the workflow, when the uploaded video hits the rule to this object, the workflow will be triggered.
* `workflow_name` - (Required, String) Workflow name, up to 128 characters. The name is unique for the same user.
* `ai_analysis_task` - (Optional, List) Video Content Analysis Type Task Parameters.
* `ai_content_review_task` - (Optional, List) Video Content Moderation Type Task Parameters.
* `ai_recognition_task` - (Optional, List) Video content recognition type task parameters.
* `media_process_task` - (Optional, List) Media Processing Type Task Parameters.
* `output_dir` - (Optional, String) The target directory of the output file generated by media processing, if not filled, it means that it is consistent with the directory where the trigger file is located.
* `output_storage` - (Optional, List) File output storage location for media processing. If left blank, the storage location in Trigger will be inherited.
* `task_notify_config` - (Optional, List) The event notification configuration of the task, if it is not filled, it means that the event notification will not be obtained.
* `task_priority` - (Optional, Int) The priority of the workflow, the larger the value, the higher the priority, the value range is -10 to 10, and blank means 0.

The `adaptive_dynamic_streaming_task_set` object supports the following:

* `definition` - (Required, Int) Transfer Adaptive Code Stream Template ID.
* `output_object_path` - (Optional, String) After converting to an adaptive stream, the output path of the manifest file can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_adaptiveDynamicStreaming_{definition}.{format}`.
* `output_storage` - (Optional, List) &quot;The target storage of the file after converting to the adaptive code stream, if not filled, it will inherit the OutputStorage value of the upper layer.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `segment_object_name` - (Optional, String) After converting to an adaptive stream (only HLS), the output path of the fragmented file can only be a relative path. If not filled, the default is a relative path: `{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{format}`.
* `sub_stream_object_name` - (Optional, String) After converting to an adaptive stream, the output path of the sub-stream file can only be a relative path. If not filled, the default is a relative path: {inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{format}`.
* `watermark_set` - (Optional, List) Watermark list, support multiple pictures or text watermarks, up to 10.

The `ai_analysis_task` object supports the following:

* `definition` - (Required, Int) Video Content Analysis Template ID.
* `extended_parameter` - (Optional, String) &quot;Extension parameter whose value is a serialized json string.&quot;&quot;Note: This parameter is a customized demand parameter, which requires offline docking.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.

The `ai_content_review_task` object supports the following:

* `definition` - (Required, Int) Video Content Review Template ID.

The `ai_recognition_task` object supports the following:

* `definition` - (Required, Int) Video Intelligent Recognition Template ID.

The `animated_graphic_task_set` object supports the following:

* `definition` - (Required, Int) Video turntable template id.
* `end_time_offset` - (Required, Float64) The end time of the animation in the video, in seconds.
* `start_time_offset` - (Required, Float64) The start time of the animation in the video, in seconds.
* `output_object_path` - (Optional, String) The output path of the file after rotating the image, which can be a relative path or an absolute path. If not filled, the default is a relative path: {inputName}_animatedGraphic_{definition}.{format}.
* `output_storage` - (Optional, List) &quot;The target storage of the transcoded file, if not filled, it will inherit the OutputStorage value of the upper layer.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.

The `audio_template` object supports the following:

* `audio_channel` - (Optional, Int) &quot;Audio channel mode, optional values:`&quot;&quot;1: single channel.&quot;&quot;2: Dual channel.&quot;&quot;6: Stereo.&quot;&quot;When the package format of the media is an audio format (flac, ogg, mp3, m4a), the number of channels is not allowed to be set to stereo.&quot;.
* `bitrate` - (Optional, Int) &quot;Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.&quot;&quot;When the value is 0, it means that the video bit rate is consistent with the original video.&quot;.
* `codec` - (Optional, String) &quot;Encoding format of frequency stream.&quot;&quot;When the outer parameter Container is mp3, the optional value is:&quot;&quot;libmp3lame.&quot;&quot;When the outer parameter Container is ogg or flac, the optional value is:&quot;&quot;flac.&quot;&quot;When the outer parameter Container is m4a, the optional value is:&quot;&quot;libfdk_aac.&quot;&quot;libmp3lame.&quot;&quot;ac3.&quot;&quot;When the outer parameter Container is mp4 or flv, the optional value is:&quot;&quot;libfdk_aac: more suitable for mp4.&quot;&quot;libmp3lame: more suitable for flv.&quot;&quot;When the outer parameter Container is hls, the optional value is:&quot;&quot;libfdk_aac.&quot;&quot;libmp3lame.&quot;.
* `sample_rate` - (Optional, Int) &quot;Sampling rate of audio stream, optional value.&quot;&quot;32000.&quot;&quot;44100.&quot;&quot;48000.&quot;&quot;Unit: Hz.&quot;.
* `stream_selects` - (Optional, Set) Specifies the audio track to preserve for the output. The default is to keep all sources.

The `audio_template` object supports the following:

* `bitrate` - (Required, Int) &quot;Bit rate of the audio stream, value range: 0 and [26, 256], unit: kbps.&quot;&quot;When the value is 0, it means that the audio bit rate is consistent with the original audio.&quot;.
* `codec` - (Required, String) &quot;Encoding format of frequency stream.&quot;&quot;When the outer parameter Container is mp3, the optional value is:&quot;&quot;libmp3lame.&quot;&quot;When the outer parameter Container is ogg or flac, the optional value is:&quot;&quot;flac.&quot;&quot;When the outer parameter Container is m4a, the optional value is:&quot;&quot;libfdk_aac.&quot;&quot;libmp3lame.&quot;&quot;ac3.&quot;&quot;When the outer parameter Container is mp4 or flv, the optional value is:&quot;&quot;libfdk_aac: more suitable for mp4.&quot;&quot;libmp3lame: more suitable for flv.&quot;&quot;When the outer parameter Container is hls, the optional value is:&quot;&quot;libfdk_aac.&quot;&quot;libmp3lame.&quot;.
* `sample_rate` - (Required, Int) &quot;Sampling rate of audio stream, optional value.&quot;&quot;32000.&quot;&quot;44100.&quot;&quot;48000.&quot;&quot;Unit: Hz.&quot;.
* `audio_channel` - (Optional, Int) &quot;Audio channel mode, optional values:`&quot;&quot;1: single channel.&quot;&quot;2: Dual channel.&quot;&quot;6: Stereo.&quot;&quot;When the package format of the media is an audio format (flac, ogg, mp3, m4a), the number of channels is not allowed to be set to stereo.&quot;&quot;Default: 2.&quot;.

The `cos_file_upload_trigger` object supports the following:

* `bucket` - (Required, String) The name of the COS Bucket bound to the workflow.
* `region` - (Required, String) The park to which the COS Bucket bound to the workflow belongs.
* `dir` - (Optional, String) The input path directory of the workflow binding must be an absolute path, that is, start and end with `/`.
* `formats` - (Optional, Set) A list of file formats that are allowed to be triggered by the workflow, if not filled in, it means that files of all formats can trigger the workflow.

The `cos_input_info` object supports the following:

* `bucket` - (Required, String) The name of the COS Bucket where the media processing object file is located.
* `object` - (Required, String) Input path for media processing object files.
* `region` - (Required, String) The park to which the COS Bucket where the media processing target file resides belongs.

The `cos_output_storage` object supports the following:

* `bucket` - (Optional, String) The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.
* `region` - (Optional, String) The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.

The `head_set` object supports the following:

* `type` - (Required, String) Enter the type of source object, which supports COS and URL.
* `cos_input_info` - (Optional, List) Valid when Type is COS, this item is required, indicating media processing COS object information.
* `url_input_info` - (Optional, List) &quot;Valid when Type is URL, this item is required, indicating media processing URL object information.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.

The `head_tail_parameter` object supports the following:

* `head_set` - (Optional, List) Title list.
* `tail_set` - (Optional, List) Ending List.

The `image_content` object supports the following:

* `type` - (Required, String) Enter the type of source object, which supports COS and URL.
* `cos_input_info` - (Optional, List) Valid when Type is COS, this item is required, indicating media processing COS object information.
* `url_input_info` - (Optional, List) &quot;Valid when Type is URL, this item is required, indicating media processing URL object information.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.

The `image_sprite_task_set` object supports the following:

* `definition` - (Required, Int) Sprite Illustration Template ID.
* `object_number_format` - (Optional, List) &quot;Rules for the `{number}` variable in the output path after intercepting the Sprite image.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `output_object_path` - (Optional, String) After capturing the sprite image, the output path of the sprite image file can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_imageSprite_{definition}_{number}.{format}`.
* `output_storage` - (Optional, List) &quot;The target storage of the file after the sprite image is intercepted, if not filled, it will inherit the OutputStorage value of the upper layer.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `web_vtt_object_name` - (Optional, String) After capturing the sprite image, the output path of the Web VTT file can only be a relative path. If not filled, the default is a relative path: `{inputName}_imageSprite_{definition}.{format}`.

The `image_template` object supports the following:

* `image_content` - (Required, List) The input content of the watermark image. Support jpeg, png image format.
* `height` - (Optional, String) &quot;The height of the watermark. Support %, px two formats:&quot;&quot;When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.&quot;&quot;When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.&quot;&quot;Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.&quot;.
* `repeat_type` - (Optional, String) &quot;Watermark repeat type. Usage scenario: The watermark is a dynamic image. Ranges.&quot;&quot;once: After the dynamic watermark is played, it will no longer appear.&quot;&quot;repeat_last_frame: After the watermark is played, stay on the last frame.&quot;&quot;repeat: the watermark loops until the end of the video (default).&quot;.
* `width` - (Optional, String) &quot;The width of the watermark. Support %, px two formats:&quot;&quot;When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.&quot;&quot;When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.&quot;&quot;Default: 10%.&quot;.

The `media_process_task` object supports the following:

* `adaptive_dynamic_streaming_task_set` - (Optional, List) Transfer Adaptive Code Stream Task List.
* `animated_graphic_task_set` - (Optional, List) Video Rotation Map Task List.
* `image_sprite_task_set` - (Optional, List) Sprite image capture task list for video.
* `sample_snapshot_task_set` - (Optional, List) Screenshot task list for video sampling.
* `snapshot_by_time_offset_task_set` - (Optional, List) Screenshot the task list of the video according to the time point.
* `transcode_task_set` - (Optional, List) Video Transcoding Task List.

The `mosaic_set` object supports the following:

* `coordinate_origin` - (Optional, String) &quot;Origin position, currently only supports:&quot;&quot;TopLeft: Indicates that the coordinate origin is located in the upper left corner of the video image, and the origin of the mosaic is the upper left corner of the picture or text&quot;&quot;Default: TopLeft.&quot;.
* `end_time_offset` - (Optional, Float64) &quot;The end time offset of the mosaic, unit: second.&quot;&quot;Fill in or fill in 0, indicating that the mosaic continues until the end of the screen.&quot;&quot;When the value is greater than 0 (assumed to be n), it means that the mosaic lasts until the nth second and disappears.&quot;&quot;When the value is less than 0 (assumed to be -n), it means that the mosaic lasts until it disappears n seconds before the end of the screen.&quot;.
* `height` - (Optional, String) &quot;The height of the mosaic. Support %, px two formats.&quot;&quot;When the string ends with %, it means that the mosaic Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.&quot;&quot;When the string ends with px, it means that the mosaic Height unit is pixel, such as 100px means that the Height is 100 pixels.&quot;&quot;Default: 10%.&quot;.
* `start_time_offset` - (Optional, Float64) &quot;The start time offset of the mosaic, unit: second. Do not fill or fill in 0, which means that the mosaic will start to appear when the screen appears.&quot;&quot;Fill in or fill in 0, which means that the mosaic will appear from the beginning of the screen.&quot;&quot;When the value is greater than 0 (assumed to be n), it means that the mosaic appears from the nth second of the screen.&quot;&quot;When the value is less than 0 (assumed to be -n), it means that the mosaic starts to appear n seconds before the end of the screen.&quot;.
* `width` - (Optional, String) &quot;The width of the mosaic. Support %, px two formats:&quot;&quot;When the string ends with %, it means that the mosaic Width is the percentage size of the video width, such as 10% means that the Width is 10% of the video width.&quot;&quot;The string ends with px, indicating that the mosaic Width unit is pixels, such as 100px indicates that the Width is 100 pixels.&quot;&quot;Default: 10%.&quot;.
* `x_pos` - (Optional, String) &quot;The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:&quot;&quot;When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.&quot;&quot;When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.&quot;&quot;Default: 0px.&quot;.
* `y_pos` - (Optional, String) &quot;The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:&quot;&quot;When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.&quot;&quot;When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.&quot;&quot;Default: 0px.&quot;.

The `object_number_format` object supports the following:

* `increment` - (Optional, Int) The growth step of `{number}` variable, default is 1.
* `initial_value` - (Optional, Int) The starting value of `{number}` variable, the default is 0.
* `min_length` - (Optional, Int) The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.
* `place_holder` - (Optional, String) When the length of the `{number}` variable is insufficient, a placeholder is added. Default is &quot;0&quot;.

The `object_number_format` object supports the following:

* `increment` - (Optional, Int) The growth step of the `{number}` variable, the default is 1.
* `initial_value` - (Optional, Int) The starting value of `{number}` variable, the default is 0.
* `min_length` - (Optional, Int) The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.
* `place_holder` - (Optional, String) When the length of the `{number}` variable is insufficient, a placeholder is added. Default is &quot;0&quot;.

The `output_storage` object supports the following:

* `type` - (Required, String) The type of media processing output object storage location, now only supports COS.
* `cos_output_storage` - (Optional, List) &quot;Valid when Type is COS, this item is required, indicating the media processing COS output location.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.

The `override_parameter` object supports the following:

* `audio_template` - (Optional, List) Audio stream configuration parameters.
* `container` - (Optional, String) Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.
* `remove_audio` - (Optional, Int) &quot;Whether to remove audio data, value:&quot;&quot;0: reserved.&quot;&quot;1: remove.&quot;.
* `remove_video` - (Optional, Int) &quot;Whether to remove video data, value:&quot;&quot;0: reserved.&quot;&quot;1: remove.&quot;.
* `subtitle_template` - (Optional, List) Subtitle Stream Configuration Parameters.
* `t_e_h_d_config` - (Optional, List) Ultra-fast HD transcoding parameters.
* `video_template` - (Optional, List) Video streaming configuration parameters.

The `raw_parameter` object supports the following:

* `container` - (Required, String) Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.
* `audio_template` - (Optional, List) Audio stream configuration parameters, when RemoveAudio is 0, this field is required.
* `remove_audio` - (Optional, Int) &quot;Whether to remove audio data, value:&quot;&quot;0: reserved.&quot;&quot;1: remove.&quot;&quot;Default: 0.&quot;.
* `remove_video` - (Optional, Int) &quot;Whether to remove video data, value:&quot;&quot;0: reserved.&quot;&quot;1: remove.&quot;&quot;Default: 0.&quot;.
* `t_e_h_d_config` - (Optional, List) Ultra-fast HD transcoding parameters.
* `video_template` - (Optional, List) Video stream configuration parameters, when RemoveVideo is 0, this field is required.

The `raw_parameter` object supports the following:

* `type` - (Required, String) &quot;Watermark type, optional value:&quot;&quot;image: image watermark.&quot;.
* `coordinate_origin` - (Optional, String) &quot;Origin position, currently only supports:&quot;&quot;TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.&quot;&quot;Default: TopLeft.&quot;.
* `image_template` - (Optional, List) Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.
* `x_pos` - (Optional, String) &quot;The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:&quot;&quot;When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.&quot;&quot;When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.&quot;&quot;Default: 0px.&quot;.
* `y_pos` - (Optional, String) &quot;The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:&quot;&quot;When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.&quot;&quot;When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.&quot;&quot;Default: 0px.&quot;.

The `sample_snapshot_task_set` object supports the following:

* `definition` - (Required, Int) Sample screenshot template ID.
* `object_number_format` - (Optional, List) &quot;Rules for the `{number}` variable in the output path after sampling the screenshot.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `output_object_path` - (Optional, String) The output path of the image file after sampling the screenshot, which can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_sampleSnapshot_{definition}_{number}.{format}`.
* `output_storage` - (Optional, List) &quot;The target storage of the file after the screenshot at the time point, if not filled, it will inherit the OutputStorage value of the upper layer.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `watermark_set` - (Optional, List) Watermark list, support multiple pictures or text watermarks, up to 10.

The `snapshot_by_time_offset_task_set` object supports the following:

* `definition` - (Required, Int) Specified time point screenshot template ID.
* `ext_time_offset_set` - (Optional, Set) &quot;Screenshot time point list, the time point supports two formats: s and %:&quot;&quot;When the string ends with s, it means that the time point is in seconds, such as 3.5s means that the time point is the 3.5th second.&quot;&quot;When the string ends with %, it means that the time point is the percentage of the video duration, such as 10% means that the time point is the first 10% of the time in the video&quot;.
* `object_number_format` - (Optional, List) &quot;Rules for the `{number}` variable in the output path after the screenshot at the time point.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `output_object_path` - (Optional, String) The output path of the picture file after the snapshot at the time point can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_snapshotByTimeOffset_{definition}_{number}.{format}`.
* `output_storage` - (Optional, List) &quot;The target storage of the file after the screenshot at the time point, if not filled, it will inherit the OutputStorage value of the upper layer.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `time_offset_set` - (Optional, ) Screenshot time point list, the unit is &lt;font color=red&gt;seconds&lt;/font&gt;. This parameter is no longer recommended, it is recommended that you use the ExtTimeOffsetSet parameter.
* `watermark_set` - (Optional, List) Watermark list, support multiple pictures or text watermarks, up to 10.

The `subtitle_template` object supports the following:

* `font_alpha` - (Optional, Float64) &quot;Text transparency, value range: (0, 1].&quot;&quot;0: fully transparent.&quot;&quot;1: fully opaque.&quot;&quot;Default: 1.&quot;.
* `font_color` - (Optional, String) Font color, format: 0xRRGGBB, default value: 0xFFFFFF (white).
* `font_size` - (Optional, String) Font size, format: Npx, N is a value, if not specified, the subtitle file shall prevail.
* `font_type` - (Optional, String) &quot;Font type.&quot;&quot;hei.ttf, song.ttf, simkai.ttf, arial.ttf.&quot;&quot;Default: hei.ttf&quot;.
* `path` - (Optional, String) The address of the subtitle file to be compressed into the video.
* `stream_index` - (Optional, Int) Specifies the subtitle track to be compressed into the video. If there is a specified Path, the Path has a higher priority. Path and StreamIndex specify at least one.

The `t_e_h_d_config` object supports the following:

* `max_video_bitrate` - (Optional, Int) The upper limit of the video bit rate, No filling means no modification.
* `type` - (Optional, String) &quot;Extremely high-definition type, optional value:&quot;&quot;TEHD-100: Extreme HD-100.&quot;&quot;Not filling means that the ultra-fast high-definition is not enabled.&quot;.

The `t_e_h_d_config` object supports the following:

* `type` - (Required, String) &quot;Extremely high-definition type, optional value:&quot;&quot;TEHD-100: Extreme HD-100.&quot;&quot;Not filling means that the ultra-fast high-definition is not enabled.&quot;.
* `max_video_bitrate` - (Optional, Int) &quot;The upper limit of the video bit rate, which is valid when the Type specifies the ultra-fast HD type.&quot;&quot;Do not fill in or fill in 0 means that there is no upper limit on the video bit rate.&quot;.

The `tail_set` object supports the following:

* `type` - (Required, String) Enter the type of source object, which supports COS and URL.
* `cos_input_info` - (Optional, List) Valid when Type is COS, this item is required, indicating media processing COS object information.
* `url_input_info` - (Optional, List) &quot;Valid when Type is URL, this item is required, indicating media processing URL object information.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.

The `task_notify_config` object supports the following:

* `cmq_model` - (Optional, String) CMQ or TDMQ-CMQ model, there are two kinds of Queue and Topic.
* `cmq_region` - (Optional, String) Region of CMQ or TDMQ-CMQ, such as sh, bj, etc.
* `notify_mode` - (Optional, String) The mode of the workflow notification, the possible values are Finish and Change, leaving blank means Finish.
* `notify_type` - (Optional, String) &quot;Notification type, optional value:&quot;&quot;CMQ: offline, it is recommended to switch to TDMQ-CMQ.&quot;&quot;TDMQ-CMQ: message queue.&quot;&quot;URL: When the URL is specified, the HTTP callback is pushed to the address specified by NotifyUrl, the callback protocol is http+json, and the package body content is the same as the output parameters of the parsing event notification interface.&quot;&quot;SCF: not recommended, additional configuration of SCF in the console is required.&quot;&quot;Note: CMQ is the default when not filled or empty, if you need to use other types, you need to fill in the corresponding type value.&quot;.
* `notify_url` - (Optional, String) HTTP callback address, required when NotifyType is URL.
* `queue_name` - (Optional, String) Valid when the model is Queue, indicating the queue name of the CMQ or TDMQ-CMQ that receives the event notification.
* `topic_name` - (Optional, String) Valid when the model is a Topic, indicating the topic name of the CMQ or TDMQ-CMQ that receives event notifications.

The `transcode_task_set` object supports the following:

* `definition` - (Required, Int) Video Transcoding Template ID.
* `end_time_offset` - (Optional, Float64) &quot;End time offset of video after transcoding, unit: second.&quot;&quot;Do not fill in or fill in 0, indicating that the transcoded video continues until the end of the original video..&quot;&quot;When the value is greater than 0 (assumed to be n), it means that the transcoded video lasts until the nth second of the original video and terminates.&quot;&quot;When the value is less than 0 (assumed to be -n), it means that the transcoded video lasts until n seconds before the end of the original video..&quot;.
* `head_tail_parameter` - (Optional, List) &quot;Opening and ending parameters.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `mosaic_set` - (Optional, List) Mosaic list, up to 10 sheets can be supported.
* `object_number_format` - (Optional, List) &quot;Rules for the `{number}` variable in the output path after transcoding.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `output_object_path` - (Optional, String) The output path of the main file after transcoding can be a relative path or an absolute path. If not filled, the default is a relative path: {inputName}_transcode_{definition}.{format}.
* `output_storage` - (Optional, List) &quot;The target storage of the transcoded file, if not filled, it will inherit the OutputStorage value of the upper layer.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `override_parameter` - (Optional, List) &quot;Video transcoding custom parameters, valid when Definition is not filled with 0.&quot;&quot;When some transcoding parameters in this structure are filled in, the parameters in the transcoding template will be overwritten with the filled parameters.&quot;&quot;This parameter is used in highly customized scenarios, it is recommended that you only use Definition to specify transcoding parameters.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `raw_parameter` - (Optional, List) &quot;Video transcoding custom parameters, valid when Definition is filled with 0.&quot;&quot;This parameter is used in highly customized scenarios. It is recommended that you use Definition to specify transcoding parameters first.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.
* `segment_object_name` - (Optional, String) The output path of the transcoded fragment file (the path of ts when transcoding HLS), can only be a relative path. If not filled, the default is: `{inputName}_transcode_{definition}_{number}.{format}.
* `start_time_offset` - (Optional, Float64) &quot;The start time offset of the transcoded video, unit: second.&quot;&quot;Do not fill in or fill in 0, indicating that the transcoded video starts from the beginning of the original video.&quot;&quot;When the value is greater than 0 (assumed to be n), it means that the transcoded video starts from the nth second position of the original video.&quot;&quot;When the value is less than 0 (assumed to be -n), it means that the transcoded video starts from the position n seconds before the end of the original video.&quot;.
* `watermark_set` - (Optional, List) &quot;Watermark list, support multiple pictures or text watermarks, up to 10.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.

The `trigger` object supports the following:

* `type` - (Required, String) The type of trigger, currently only supports CosFileUpload.
* `cos_file_upload_trigger` - (Optional, List) &quot;Mandatory and valid when Type is CosFileUpload, the rule is triggered for COS.&quot;&quot;Note: This field may return null, indicating that no valid value can be obtained.&quot;.

The `url_input_info` object supports the following:

* `url` - (Required, String) Video URL.

The `video_template` object supports the following:

* `bitrate` - (Optional, Int) &quot;Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.&quot;&quot;When the value is 0, it means that the video bit rate is consistent with the original video.&quot;.
* `codec` - (Optional, String) &quot;Encoding format of the video stream, optional value:&quot;&quot;libx264: H.264 encoding.&quot;&quot;libx265: H.265 encoding.&quot;&quot;av1: AOMedia Video 1 encoding.&quot;&quot;Note: Currently H.265 encoding must specify a resolution, and it needs to be within 640*480.&quot;&quot;Note: av1 encoded containers currently only support mp4.&quot;.
* `content_adapt_stream` - (Optional, Int) &quot;Content Adaptive Encoding. optional value:&quot;&quot;0: not open.&quot;&quot;1: open.&quot;&quot;Default: 0.&quot;&quot;When this parameter is turned on, multiple code streams with different resolutions and different bit rates will be adaptively generated. The width and height of the VideoTemplate are the maximum resolutions among the multiple code streams, and the bit rates in the VideoTemplate are multiple code rates. The highest bit rate in the stream, the vcrf in VideoTemplate is the highest quality among multiple bit streams. When the resolution, bit rate and vcrf are not set, the highest resolution generated by the ContentAdaptStream parameter is the resolution of the video source, and the video quality is close to vmaf95. To enable this parameter or learn about billing details, please contact your Tencent Cloud Business.&quot;.
* `fill_type` - (Optional, String) &quot;Filling method, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling;. Optional filling method:&quot;&quot;stretch: Stretch, stretch each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched; black: Leave black, keep the aspect ratio of the video unchanged, and fill the rest of the edge with black.&quot;&quot;white: Leave blank, keep the aspect ratio of the video unchanged, and fill the rest of the edge with white.&quot;&quot;gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and fill the rest of the edge with Gaussian blur.&quot;.
* `fps` - (Optional, Int) &quot;Video frame rate, value range: [0, 100], unit: Hz.&quot;&quot;When the value is 0, it means that the frame rate is consistent with the original video.&quot;.
* `gop` - (Optional, Int) &quot;The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.&quot;&quot;When filling 0 or not filling, the system will automatically set the gop length.&quot;.
* `height` - (Optional, Int) The maximum value of video stream height (or short side), value range: 0 and [128, 4096], unit: px.
* `resolution_adaptive` - (Optional, String) &quot;Adaptive resolution, optional values:```&quot;&quot;open: open, at this time, Width represents the long side of the video, Height represents the short side of the video.&quot;&quot;close: close, at this time, Width represents the width of the video, and Height represents the height of the video.&quot;&quot;Note: In adaptive mode, Width cannot be smaller than Height.&quot;.
* `vcrf` - (Optional, Int) &quot;Video constant bit rate control factor, the value range is [1, 51], Fill in 0 to disable this parameter.&quot;&quot;If there is no special requirement, it is not recommended to specify this parameter.&quot;.
* `width` - (Optional, Int) &quot;The maximum value of video stream width (or long side), value range: 0 and [128, 4096], unit: px.&quot;&quot;When Width and Height are both 0, the resolution is the same.&quot;&quot;When Width is 0 and Height is not 0, Width is scaled proportionally.&quot;&quot;When Width is not 0 and Height is 0, Height is scaled proportionally.&quot;&quot;When both Width and Height are not 0, the resolution is specified by the user.&quot;.

The `video_template` object supports the following:

* `bitrate` - (Required, Int) &quot;Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.&quot;&quot;When the value is 0, it means that the video bit rate is consistent with the original video.&quot;.
* `codec` - (Required, String) &quot;Encoding format of the video stream, optional value:&quot;&quot;libx264: H.264 encoding.&quot;&quot;libx265: H.265 encoding.&quot;&quot;av1: AOMedia Video 1 encoding.&quot;&quot;Note: Currently H.265 encoding must specify a resolution, and it needs to be within 640*480.&quot;&quot;Note: av1 encoded containers currently only support mp4.&quot;.
* `fps` - (Required, Int) &quot;Video frame rate, value range: [0, 100], unit: Hz.&quot;&quot;When the value is 0, it means that the frame rate is consistent with the original video.&quot;&quot;Note: The value range for adaptive code rate is [0, 60].&quot;.
* `fill_type` - (Optional, String) &quot;Filling method, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling method:&quot;&quot;stretch: Stretch, stretch each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched;.&quot;&quot;black: Leave black, keep the aspect ratio of the video unchanged, and fill the rest of the edge with black.&quot;&quot;white: Leave blank, keep the aspect ratio of the video unchanged, and fill the rest of the edge with white.&quot;&quot;gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and fill the rest of the edge with Gaussian blur.&quot;&quot;Default: black.&quot;&quot;Note: Adaptive stream only supports stretch, black.&quot;.
* `gop` - (Optional, Int) &quot;The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.&quot;&quot;When filling 0 or not filling, the system will automatically set the gop length.&quot;.
* `height` - (Optional, Int) &quot;The maximum value of video stream height (or short side), value range: 0 and [128, 4096], unit: px.&quot;&quot;When Width and Height are both 0, the resolution is the same.&quot;&quot;When Width is 0 and Height is not 0, Width is scaled proportionally.&quot;&quot;When Width is not 0 and Height is 0, Height is scaled proportionally.&quot;&quot;When both Width and Height are not 0, the resolution is specified by the user.&quot;&quot;Default: 0.&quot;.
* `resolution_adaptive` - (Optional, String) &quot;Adaptive resolution, optional values:```&quot;&quot;open: open, at this time, Width represents the long side of the video, Height represents the short side of the video.&quot;&quot;close: close, at this time, Width represents the width of the video, and Height represents the height of the video.&quot;&quot;Default: open.&quot;&quot;Note: In adaptive mode, Width cannot be smaller than Height.&quot;.
* `vcrf` - (Optional, Int) &quot;Video constant bit rate control factor, the value range is [1, 51].&quot;&quot;If this parameter is specified, the code rate control method of CRF will be used for transcoding (the video code rate will no longer take effect).&quot;&quot;If there is no special requirement, it is not recommended to specify this parameter.&quot;.
* `width` - (Optional, Int) &quot;The maximum value of video stream width (or long side), value range: 0 and [128, 4096], unit: px.&quot;&quot;When Width and Height are both 0, the resolution is the same.&quot;&quot;When Width is 0 and Height is not 0, Width is scaled proportionally.&quot;&quot;When Width is not 0 and Height is 0, Height is scaled proportionally.&quot;&quot;When both Width and Height are not 0, the resolution is specified by the user.&quot;&quot;Default: 0&quot;.

The `watermark_set` object supports the following:

* `definition` - (Required, Int) Watermark Template ID.
* `end_time_offset` - (Optional, Float64) &quot;End time offset of watermark, unit: second.&quot;&quot;Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.&quot;&quot;When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.&quot;&quot;When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.&quot;.
* `raw_parameter` - (Optional, List) &quot;Watermark custom parameters, valid when Definition is filled with 0.&quot;&quot;This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.&quot;&quot;Watermark custom parameters do not support screenshot watermarking.&quot;.
* `start_time_offset` - (Optional, Float64) &quot;The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.&quot;&quot;Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.&quot;&quot;When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.&quot;&quot;When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.&quot;.
* `svg_content` - (Optional, String) &quot;SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.&quot;&quot;SVG watermark does not support screenshot watermarking.&quot;.
* `text_content` - (Optional, String) &quot;Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.&quot;&quot;Text watermark does not support screenshot watermarking.&quot;.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps workflow can be imported using the id, e.g.

```
terraform import tencentcloud_mps_workflow.workflow workflow_id
```

