/*
Provides a resource to create a mps workflow

Example Usage

```hcl
resource "tencentcloud_mps_workflow" "workflow" {
  output_dir    = "/"
  task_priority = 0
  workflow_name = "tf-workflow"

  media_process_task {
    adaptive_dynamic_streaming_task_set {
      definition             = 12
      output_object_path     = "/out"
      segment_object_name    = "/out"
      sub_stream_object_name = "/out/out/"

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }

    snapshot_by_time_offset_task_set {
      definition          = 10
      ext_time_offset_set = [
        "1s",
      ]
      output_object_path  = "/snapshot/"
      time_offset_set     = []

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }

    animated_graphic_task_set {
      definition         = 20000
      end_time_offset    = 0
      output_object_path = "/test/"
      start_time_offset  = 0

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }
  }

  ai_analysis_task {
    definition = 20
  }

  ai_content_review_task {
    definition = 20
  }

  ai_recognition_task {
    definition = 20
  }

  output_storage {
    type = "COS"

    cos_output_storage {
      bucket = "cos-lock-1308919341"
      region = "ap-guangzhou"
    }
  }

  trigger {
    type = "CosFileUpload"

    cos_file_upload_trigger {
      bucket = "cos-lock-1308919341"
      dir    = "/"
      region = "ap-guangzhou"
    }
  }
}

```

Import

mps workflow can be imported using the id, e.g.

```
terraform import tencentcloud_mps_workflow.workflow workflow_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsWorkflowCreate,
		Read:   resourceTencentCloudMpsWorkflowRead,
		Update: resourceTencentCloudMpsWorkflowUpdate,
		Delete: resourceTencentCloudMpsWorkflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"workflow_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Workflow name, up to 128 characters. The name is unique for the same user.",
			},

			"trigger": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The trigger rule bound to the workflow, when the uploaded video hits the rule to this object, the workflow will be triggered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of trigger, currently only supports CosFileUpload.",
						},
						"cos_file_upload_trigger": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Mandatory and valid when Type is CosFileUpload, the rule is triggered for COS.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the COS Bucket bound to the workflow.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The park to which the COS Bucket bound to the workflow belongs.",
									},
									"dir": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The input path directory of the workflow binding must be an absolute path, that is, start and end with `/`.",
									},
									"formats": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "A list of file formats that are allowed to be triggered by the workflow, if not filled in, it means that files of all formats can trigger the workflow.",
									},
								},
							},
						},
					},
				},
			},

			"output_storage": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "File output storage location for media processing. If left blank, the storage location in Trigger will be inherited.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of media processing output object storage location, now only supports COS.",
						},
						"cos_output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.",
									},
								},
							},
						},
					},
				},
			},

			"output_dir": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The target directory of the output file generated by media processing, if not filled, it means that it is consistent with the directory where the trigger file is located.",
			},

			"media_process_task": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Media Processing Type Task Parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transcode_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Video Transcoding Task List.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Video Transcoding Template ID.",
									},
									"raw_parameter": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Video transcoding custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios. It is recommended that you use Definition to specify transcoding parameters first.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.",
												},
												"remove_video": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Whether to remove video data, value:0: reserved.1: remove.Default: 0.",
												},
												"remove_audio": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Whether to remove audio data, value:0: reserved.1: remove.Default: 0.",
												},
												"video_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Video stream configuration parameters, when RemoveVideo is 0, this field is required.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Encoding format of the video stream, optional value:libx264: H.264 encoding.libx265: H.265 encoding.av1: AOMedia Video 1 encoding.Note: Currently H.265 encoding must specify a resolution, and it needs to be within 640*480.Note: av1 encoded containers currently only support mp4.",
															},
															"fps": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Video frame rate, value range: [0, 100], unit: Hz.When the value is 0, it means that the frame rate is consistent with the original video.Note: The value range for adaptive code rate is [0, 60].",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.When the value is 0, it means that the video bit rate is consistent with the original video.",
															},
															"resolution_adaptive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Adaptive resolution, optional values:```open: open, at this time, Width represents the long side of the video, Height represents the short side of the video.close: close, at this time, Width represents the width of the video, and Height represents the height of the video.Default: open.Note: In adaptive mode, Width cannot be smaller than Height.",
															},
															"width": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The maximum value of video stream width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default: 0.",
															},
															"height": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The maximum value of video stream height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default: 0.",
															},
															"gop": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.When filling 0 or not filling, the system will automatically set the gop length.",
															},
															"fill_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Filling method, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling method:stretch: Stretch, stretch each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched;.black: Leave black, keep the aspect ratio of the video unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video unchanged, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and fill the rest of the edge with Gaussian blur.Default: black.Note: Adaptive stream only supports stretch, black.",
															},
															"vcrf": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Video constant bit rate control factor, the value range is [1, 51].If this parameter is specified, the code rate control method of CRF will be used for transcoding (the video code rate will no longer take effect).If there is no special requirement, it is not recommended to specify this parameter.",
															},
														},
													},
												},
												"audio_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Audio stream configuration parameters, when RemoveAudio is 0, this field is required.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Encoding format of frequency stream.When the outer parameter Container is mp3, the optional value is:libmp3lame.When the outer parameter Container is ogg or flac, the optional value is:flac.When the outer parameter Container is m4a, the optional value is:libfdk_aac.libmp3lame.ac3.When the outer parameter Container is mp4 or flv, the optional value is:libfdk_aac: more suitable for mp4.libmp3lame: more suitable for flv.When the outer parameter Container is hls, the optional value is:libfdk_aac.libmp3lame.",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Bit rate of the audio stream, value range: 0 and [26, 256], unit: kbps.When the value is 0, it means that the audio bit rate is consistent with the original audio.",
															},
															"sample_rate": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Sampling rate of audio stream, optional value.32000.44100.48000.Unit: Hz.",
															},
															"audio_channel": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio channel mode, optional values:`1: single channel.2: Dual channel.6: Stereo.When the package format of the media is an audio format (flac, ogg, mp3, m4a), the number of channels is not allowed to be set to stereo.Default: 2.",
															},
														},
													},
												},
												"tehd_config": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Ultra-fast HD transcoding parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Extremely high-definition type, optional value:TEHD-100: Extreme HD-100.Not filling means that the ultra-fast high-definition is not enabled.",
															},
															"max_video_bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The upper limit of the video bit rate, which is valid when the Type specifies the ultra-fast HD type.Do not fill in or fill in 0 means that there is no upper limit on the video bit rate.",
															},
														},
													},
												},
											},
										},
									},
									"override_parameter": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Video transcoding custom parameters, valid when Definition is not filled with 0.When some transcoding parameters in this structure are filled in, the parameters in the transcoding template will be overwritten with the filled parameters.This parameter is used in highly customized scenarios, it is recommended that you only use Definition to specify transcoding parameters.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.",
												},
												"remove_video": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Whether to remove video data, value:0: reserved.1: remove.",
												},
												"remove_audio": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Whether to remove audio data, value:0: reserved.1: remove.",
												},
												"video_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Video streaming configuration parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Encoding format of the video stream, optional value:libx264: H.264 encoding.libx265: H.265 encoding.av1: AOMedia Video 1 encoding.Note: Currently H.265 encoding must specify a resolution, and it needs to be within 640*480.Note: av1 encoded containers currently only support mp4.",
															},
															"fps": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Video frame rate, value range: [0, 100], unit: Hz.When the value is 0, it means that the frame rate is consistent with the original video.",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.When the value is 0, it means that the video bit rate is consistent with the original video.",
															},
															"resolution_adaptive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Adaptive resolution, optional values:```open: open, at this time, Width represents the long side of the video, Height represents the short side of the video.close: close, at this time, Width represents the width of the video, and Height represents the height of the video.Note: In adaptive mode, Width cannot be smaller than Height.",
															},
															"width": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The maximum value of video stream width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.",
															},
															"height": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The maximum value of video stream height (or short side), value range: 0 and [128, 4096], unit: px.",
															},
															"gop": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.When filling 0 or not filling, the system will automatically set the gop length.",
															},
															"fill_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Filling method, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling;. Optional filling method:stretch: Stretch, stretch each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched; black: Leave black, keep the aspect ratio of the video unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video unchanged, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and fill the rest of the edge with Gaussian blur.",
															},
															"vcrf": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Video constant bit rate control factor, the value range is [1, 51], Fill in 0 to disable this parameter.If there is no special requirement, it is not recommended to specify this parameter.",
															},
															"content_adapt_stream": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Content Adaptive Encoding. optional value:0: not open.1: open.Default: 0.When this parameter is turned on, multiple code streams with different resolutions and different bit rates will be adaptively generated. The width and height of the VideoTemplate are the maximum resolutions among the multiple code streams, and the bit rates in the VideoTemplate are multiple code rates. The highest bit rate in the stream, the vcrf in VideoTemplate is the highest quality among multiple bit streams. When the resolution, bit rate and vcrf are not set, the highest resolution generated by the ContentAdaptStream parameter is the resolution of the video source, and the video quality is close to vmaf95. To enable this parameter or learn about billing details, please contact your Tencent Cloud Business.",
															},
														},
													},
												},
												"audio_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Audio stream configuration parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Encoding format of frequency stream.When the outer parameter Container is mp3, the optional value is:libmp3lame.When the outer parameter Container is ogg or flac, the optional value is:flac.When the outer parameter Container is m4a, the optional value is:libfdk_aac.libmp3lame.ac3.When the outer parameter Container is mp4 or flv, the optional value is:libfdk_aac: more suitable for mp4.libmp3lame: more suitable for flv.When the outer parameter Container is hls, the optional value is:libfdk_aac.libmp3lame.",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.When the value is 0, it means that the video bit rate is consistent with the original video.",
															},
															"sample_rate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Sampling rate of audio stream, optional value.32000.44100.48000.Unit: Hz.",
															},
															"audio_channel": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio channel mode, optional values:`1: single channel.2: Dual channel.6: Stereo.When the package format of the media is an audio format (flac, ogg, mp3, m4a), the number of channels is not allowed to be set to stereo.",
															},
															"stream_selects": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeInt,
																},
																Optional:    true,
																Description: "Specifies the audio track to preserve for the output. The default is to keep all sources.",
															},
														},
													},
												},
												"tehd_config": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Ultra-fast HD transcoding parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Extremely high-definition type, optional value:TEHD-100: Extreme HD-100.Not filling means that the ultra-fast high-definition is not enabled.",
															},
															"max_video_bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The upper limit of the video bit rate, No filling means no modification.",
															},
														},
													},
												},
												"subtitle_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Subtitle Stream Configuration Parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The address of the subtitle file to be compressed into the video.",
															},
															"stream_index": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the subtitle track to be compressed into the video. If there is a specified Path, the Path has a higher priority. Path and StreamIndex specify at least one.",
															},
															"font_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Font type.hei.ttf, song.ttf, simkai.ttf, arial.ttf.Default: hei.ttf.",
															},
															"font_size": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Font size, format: Npx, N is a value, if not specified, the subtitle file shall prevail.",
															},
															"font_color": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Font color, format: 0xRRGGBB, default value: 0xFFFFFF (white).",
															},
															"font_alpha": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "Text transparency, value range: (0, 1].0: fully transparent.1: fully opaque.Default: 1.",
															},
														},
													},
												},
											},
										},
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Watermark list, support multiple pictures or text watermarks, up to 10.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Watermark Template ID.",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Watermark custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.Watermark custom parameters do not support screenshot watermarking.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark type, optional value:image: image watermark.",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position, currently only supports:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.Default: TopLeft.",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "The input content of the watermark image. Support jpeg, png image format.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Enter the type of source object, which supports COS and URL.",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The name of the COS Bucket where the media processing object file is located.",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Input path for media processing object files.",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Video URL.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.Default: 10%.",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark repeat type. Usage scenario: The watermark is a dynamic image. Ranges.once: After the dynamic watermark is played, it will no longer appear.repeat_last_frame: After the watermark is played, stay on the last frame.repeat: the watermark loops until the end of the video (default).",
																		},
																	},
																},
															},
														},
													},
												},
												"text_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.Text watermark does not support screenshot watermarking.",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.SVG watermark does not support screenshot watermarking.",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "End time offset of watermark, unit: second.Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.",
												},
											},
										},
									},
									"mosaic_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Mosaic list, up to 10 sheets can be supported.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"coordinate_origin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Origin position, currently only supports:TopLeft: Indicates that the coordinate origin is located in the upper left corner of the video image, and the origin of the mosaic is the upper left corner of the picture or textDefault: TopLeft.",
												},
												"x_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
												},
												"y_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
												},
												"width": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The width of the mosaic. Support %, px two formats:When the string ends with %, it means that the mosaic Width is the percentage size of the video width, such as 10% means that the Width is 10% of the video width.The string ends with px, indicating that the mosaic Width unit is pixels, such as 100px indicates that the Width is 100 pixels.Default: 10%.",
												},
												"height": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The height of the mosaic. Support %, px two formats.When the string ends with %, it means that the mosaic Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the mosaic Height unit is pixel, such as 100px means that the Height is 100 pixels.Default: 10%.",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "The start time offset of the mosaic, unit: second. Do not fill or fill in 0, which means that the mosaic will start to appear when the screen appears.Fill in or fill in 0, which means that the mosaic will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the mosaic appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the mosaic starts to appear n seconds before the end of the screen.",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "The end time offset of the mosaic, unit: second.Fill in or fill in 0, indicating that the mosaic continues until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the mosaic lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the mosaic lasts until it disappears n seconds before the end of the screen.",
												},
											},
										},
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "The start time offset of the transcoded video, unit: second.Do not fill in or fill in 0, indicating that the transcoded video starts from the beginning of the original video.When the value is greater than 0 (assumed to be n), it means that the transcoded video starts from the nth second position of the original video.When the value is less than 0 (assumed to be -n), it means that the transcoded video starts from the position n seconds before the end of the original video.",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "End time offset of video after transcoding, unit: second.Do not fill in or fill in 0, indicating that the transcoded video continues until the end of the original video.When the value is greater than 0 (assumed to be n), it means that the transcoded video lasts until the nth second of the original video and terminates.When the value is less than 0 (assumed to be -n), it means that the transcoded video lasts until n seconds before the end of the original video.",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The target storage of the transcoded file, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The type of media processing output object storage location, now only supports COS.",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The output path of the main file after transcoding can be a relative path or an absolute path. If not filled, the default is a relative path: {inputName}_transcode_{definition}.{format}.",
									},
									"segment_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The output path of the transcoded fragment file (the path of ts when transcoding HLS), can only be a relative path. If not filled, the default is: `{inputName}_transcode_{definition}_{number}.{format}.",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rules for the `{number}` variable in the output path after transcoding.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The starting value of `{number}` variable, the default is 0.",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The growth step of the `{number}` variable, the default is 1.",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "When the length of the `{number}` variable is insufficient, a placeholder is added. Default is 0.",
												},
											},
										},
									},
									"head_tail_parameter": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Opening and ending parameters.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"head_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Title list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Enter the type of source object, which supports COS and URL.",
															},
															"cos_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The name of the COS Bucket where the media processing object file is located.",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Input path for media processing object files.",
																		},
																	},
																},
															},
															"url_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Video URL.",
																		},
																	},
																},
															},
														},
													},
												},
												"tail_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Ending List.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Enter the type of source object, which supports COS and URL.",
															},
															"cos_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The name of the COS Bucket where the media processing object file is located.",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Input path for media processing object files.",
																		},
																	},
																},
															},
															"url_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Video URL.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"animated_graphic_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Video Rotation Map Task List.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Video turntable template id.",
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "The start time of the animation in the video, in seconds.",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "The end time of the animation in the video, in seconds.",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The target storage of the transcoded file, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The type of media processing output object storage location, now only supports COS.",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The output path of the file after rotating the image, which can be a relative path or an absolute path. If not filled, the default is a relative path: {inputName}_animatedGraphic_{definition}.{format}.",
									},
								},
							},
						},
						"snapshot_by_time_offset_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Screenshot the task list of the video according to the time point.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specified time point screenshot template ID.",
									},
									"ext_time_offset_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Screenshot time point list, the time point supports two formats: s and %:;When the string ends with s, it means that the time point is in seconds, such as 3.5s means that the time point is the 3.5th second.When the string ends with %, it means that the time point is the percentage of the video duration, such as 10% means that the time point is the first 10% of the time in the video.",
									},
									"time_offset_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
										Optional:    true,
										Description: "Screenshot time point list, the unit is &lt;font color=red&gt;seconds&lt;/font&gt;. This parameter is no longer recommended, it is recommended that you use the ExtTimeOffsetSet parameter.",
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Watermark list, support multiple pictures or text watermarks, up to 10.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Watermark Template ID.",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Watermark custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.Watermark custom parameters do not support screenshot watermarking.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark type, optional value:image: image watermark.",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position, currently only supports:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.Default: TopLeft.",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "The input content of the watermark image. Support jpeg, png image format.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Enter the type of source object, which supports COS and URL.",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The name of the COS Bucket where the media processing object file is located.",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Input path for media processing object files.",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Video URL.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.Default: 10%.",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark repeat type. Usage scenario: The watermark is a dynamic image. Ranges.once: After the dynamic watermark is played, it will no longer appear.repeat_last_frame: After the watermark is played, stay on the last frame.repeat: the watermark loops until the end of the video (default).",
																		},
																	},
																},
															},
														},
													},
												},
												"text_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.Text watermark does not support screenshot watermarking.",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.SVG watermark does not support screenshot watermarking.",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "End time offset of watermark, unit: second.Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The target storage of the file after the screenshot at the time point, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The type of media processing output object storage location, now only supports COS.",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The output path of the picture file after the snapshot at the time point can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_snapshotByTimeOffset_{definition}_{number}.{format}`.",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rules for the `{number}` variable in the output path after the screenshot at the time point.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The starting value of `{number}` variable, the default is 0.",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The growth step of `{number}` variable, default is 1.",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "When the length of the `{number}` variable is insufficient, a placeholder is added. Default is 0.",
												},
											},
										},
									},
								},
							},
						},
						"sample_snapshot_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Screenshot task list for video sampling.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Sample screenshot template ID.",
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Watermark list, support multiple pictures or text watermarks, up to 10.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Watermark Template ID.",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Watermark custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.Watermark custom parameters do not support screenshot watermarking.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark type, optional value:image: image watermark.",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position, currently only supports:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.Default: TopLeft.",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "The input content of the watermark image. Support jpeg, png image format.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Enter the type of source object, which supports COS and URL.",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The name of the COS Bucket where the media processing object file is located.",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Input path for media processing object files.",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Video URL.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.Default: 10%.",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark repeat type. Usage scenario: The watermark is a dynamic image. Ranges.once: After the dynamic watermark is played, it will no longer appear.repeat_last_frame: After the watermark is played, stay on the last frame.repeat: the watermark loops until the end of the video (default).",
																		},
																	},
																},
															},
														},
													},
												},
												"text_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.Text watermark does not support screenshot watermarking.",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.SVG watermark does not support screenshot watermarking.",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "End time offset of watermark, unit: second.Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The target storage of the file after the screenshot at the time point, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The type of media processing output object storage location, now only supports COS.",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The output path of the image file after sampling the screenshot, which can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_sampleSnapshot_{definition}_{number}.{format}`.",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rules for the `{number}` variable in the output path after sampling the screenshot.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The starting value of `{number}` variable, the default is 0.",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The growth step of the `{number}` variable, the default is 1.",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "When the length of the `{number}` variable is insufficient, a placeholder is added. Default is 0.",
												},
											},
										},
									},
								},
							},
						},
						"image_sprite_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Sprite image capture task list for video.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Sprite Illustration Template ID.",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The target storage of the file after the sprite image is intercepted, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The type of media processing output object storage location, now only supports COS.",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After capturing the sprite image, the output path of the sprite image file can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_imageSprite_{definition}_{number}.{format}`.",
									},
									"web_vtt_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After capturing the sprite image, the output path of the Web VTT file can only be a relative path. If not filled, the default is a relative path: `{inputName}_imageSprite_{definition}.{format}`.",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rules for the `{number}` variable in the output path after intercepting the Sprite image.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The starting value of `{number}` variable, the default is 0.",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The growth step of the `{number}` variable, the default is 1.",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "When the length of the `{number}` variable is insufficient, a placeholder is added. Default is 0.",
												},
											},
										},
									},
								},
							},
						},
						"adaptive_dynamic_streaming_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Transfer Adaptive Code Stream Task List.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Transfer Adaptive Code Stream Template ID.",
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Watermark list, support multiple pictures or text watermarks, up to 10.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Watermark Template ID.",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Watermark custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.Watermark custom parameters do not support screenshot watermarking.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark type, optional value:image: image watermark.",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position, currently only supports:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.Default: TopLeft.",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "The input content of the watermark image. Support jpeg, png image format.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Enter the type of source object, which supports COS and URL.",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The name of the COS Bucket where the media processing object file is located.",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Input path for media processing object files.",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Video URL.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.Default: 10%.",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark repeat type. Usage scenario: The watermark is a dynamic image. Ranges.once: After the dynamic watermark is played, it will no longer appear.repeat_last_frame: After the watermark is played, stay on the last frame.repeat: the watermark loops until the end of the video (default).",
																		},
																	},
																},
															},
														},
													},
												},
												"text_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.Text watermark does not support screenshot watermarking.",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.SVG watermark does not support screenshot watermarking.",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "End time offset of watermark, unit: second.Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The target storage of the file after converting to the adaptive code stream, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The type of media processing output object storage location, now only supports COS.",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After converting to an adaptive stream, the output path of the manifest file can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_adaptiveDynamicStreaming_{definition}.{format}`.",
									},
									"sub_stream_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After converting to an adaptive stream, the output path of the sub-stream file can only be a relative path. If not filled, the default is a relative path: {inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{format}`.",
									},
									"segment_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After converting to an adaptive stream (only HLS), the output path of the fragmented file can only be a relative path. If not filled, the default is a relative path: `{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{format}`.",
									},
								},
							},
						},
					},
				},
			},

			"ai_content_review_task": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video Content Moderation Type Task Parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video Content Review Template ID.",
						},
					},
				},
			},

			"ai_analysis_task": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video Content Analysis Type Task Parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video Content Analysis Template ID.",
						},
						"extended_parameter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Extension parameter whose value is a serialized json string.Note: This parameter is a customized demand parameter, which requires offline docking.Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"ai_recognition_task": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video content recognition type task parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video Intelligent Recognition Template ID.",
						},
					},
				},
			},

			"task_notify_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The event notification configuration of the task, if it is not filled, it means that the event notification will not be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmq_model": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ or TDMQ-CMQ model, there are two kinds of Queue and Topic.",
						},
						"cmq_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region of CMQ or TDMQ-CMQ, such as sh, bj, etc.",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Valid when the model is a Topic, indicating the topic name of the CMQ or TDMQ-CMQ that receives event notifications.",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Valid when the model is Queue, indicating the queue name of the CMQ or TDMQ-CMQ that receives the event notification.",
						},
						"notify_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The mode of the workflow notification, the possible values are Finish and Change, leaving blank means Finish.",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification type, optional value:CMQ: offline, it is recommended to switch to TDMQ-CMQ.TDMQ-CMQ: message queue.URL: When the URL is specified, the HTTP callback is pushed to the address specified by NotifyUrl, the callback protocol is http+json, and the package body content is the same as the output parameters of the parsing event notification interface.SCF: not recommended, additional configuration of SCF in the console is required.Note: CMQ is the default when not filled or empty, if you need to use other types, you need to fill in the corresponding type value.",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback address, required when NotifyType is URL.",
						},
					},
				},
			},

			"task_priority": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     0,
				Description: "The priority of the workflow, the larger the value, the higher the priority, the value range is -10 to 10, and blank means 0.",
			},
		},
	}
}

func resourceTencentCloudMpsWorkflowCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_workflow.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mps.NewCreateWorkflowRequest()
		response   = mps.NewCreateWorkflowResponse()
		workflowId int64
	)
	if v, ok := d.GetOk("workflow_name"); ok {
		request.WorkflowName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "trigger"); ok {
		workflowTrigger := mps.WorkflowTrigger{}
		if v, ok := dMap["type"]; ok {
			workflowTrigger.Type = helper.String(v.(string))
		}
		if cosFileUploadTriggerMap, ok := helper.InterfaceToMap(dMap, "cos_file_upload_trigger"); ok {
			cosFileUploadTrigger := mps.CosFileUploadTrigger{}
			if v, ok := cosFileUploadTriggerMap["bucket"]; ok {
				cosFileUploadTrigger.Bucket = helper.String(v.(string))
			}
			if v, ok := cosFileUploadTriggerMap["region"]; ok {
				cosFileUploadTrigger.Region = helper.String(v.(string))
			}
			if v, ok := cosFileUploadTriggerMap["dir"]; ok {
				cosFileUploadTrigger.Dir = helper.String(v.(string))
			}
			if v, ok := cosFileUploadTriggerMap["formats"]; ok {
				formatsSet := v.(*schema.Set).List()
				for i := range formatsSet {
					formats := formatsSet[i].(string)
					cosFileUploadTrigger.Formats = append(cosFileUploadTrigger.Formats, &formats)
				}
			}
			workflowTrigger.CosFileUploadTrigger = &cosFileUploadTrigger
		}
		request.Trigger = &workflowTrigger
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "output_storage"); ok {
		taskOutputStorage := mps.TaskOutputStorage{}
		if v, ok := dMap["type"]; ok {
			taskOutputStorage.Type = helper.String(v.(string))
		}
		if cosOutputStorageMap, ok := helper.InterfaceToMap(dMap, "cos_output_storage"); ok {
			cosOutputStorage := mps.CosOutputStorage{}
			if v, ok := cosOutputStorageMap["bucket"]; ok {
				cosOutputStorage.Bucket = helper.String(v.(string))
			}
			if v, ok := cosOutputStorageMap["region"]; ok {
				cosOutputStorage.Region = helper.String(v.(string))
			}
			taskOutputStorage.CosOutputStorage = &cosOutputStorage
		}
		request.OutputStorage = &taskOutputStorage
	}

	if v, ok := d.GetOk("output_dir"); ok {
		request.OutputDir = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "media_process_task"); ok {
		mediaProcessTaskInput := mps.MediaProcessTaskInput{}
		if v, ok := dMap["transcode_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				transcodeTaskSetMap := item.(map[string]interface{})
				transcodeTaskInput := mps.TranscodeTaskInput{}
				if v, ok := transcodeTaskSetMap["definition"]; ok {
					transcodeTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if rawParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "raw_parameter"); ok {
					rawTranscodeParameter := mps.RawTranscodeParameter{}
					if v, ok := rawParameterMap["container"]; ok {
						rawTranscodeParameter.Container = helper.String(v.(string))
					}
					if v, ok := rawParameterMap["remove_video"]; ok {
						rawTranscodeParameter.RemoveVideo = helper.IntInt64(v.(int))
					}
					if v, ok := rawParameterMap["remove_audio"]; ok {
						rawTranscodeParameter.RemoveAudio = helper.IntInt64(v.(int))
					}
					if videoTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "video_template"); ok {
						videoTemplateInfo := mps.VideoTemplateInfo{}
						if v, ok := videoTemplateMap["codec"]; ok {
							videoTemplateInfo.Codec = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["fps"]; ok {
							videoTemplateInfo.Fps = helper.IntInt64(v.(int))
						}
						if v, ok := videoTemplateMap["bitrate"]; ok {
							videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
						}
						if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
							videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["width"]; ok {
							videoTemplateInfo.Width = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["height"]; ok {
							videoTemplateInfo.Height = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["gop"]; ok {
							videoTemplateInfo.Gop = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["fill_type"]; ok {
							videoTemplateInfo.FillType = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["vcrf"]; ok {
							videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
						}
						rawTranscodeParameter.VideoTemplate = &videoTemplateInfo
					}
					if audioTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "audio_template"); ok {
						audioTemplateInfo := mps.AudioTemplateInfo{}
						if v, ok := audioTemplateMap["codec"]; ok {
							audioTemplateInfo.Codec = helper.String(v.(string))
						}
						if v, ok := audioTemplateMap["bitrate"]; ok {
							audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
						}
						if v, ok := audioTemplateMap["sample_rate"]; ok {
							audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
						}
						if v, ok := audioTemplateMap["audio_channel"]; ok {
							audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
						}
						rawTranscodeParameter.AudioTemplate = &audioTemplateInfo
					}
					if tEHDConfigMap, ok := helper.InterfaceToMap(rawParameterMap, "tehd_config"); ok {
						tEHDConfig := mps.TEHDConfig{}
						if v, ok := tEHDConfigMap["type"]; ok {
							tEHDConfig.Type = helper.String(v.(string))
						}
						if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
							tEHDConfig.MaxVideoBitrate = helper.IntInt64(v.(int))
						}
						rawTranscodeParameter.TEHDConfig = &tEHDConfig
					}
					transcodeTaskInput.RawParameter = &rawTranscodeParameter
				}
				if overrideParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "override_parameter"); ok {
					overrideTranscodeParameter := mps.OverrideTranscodeParameter{}
					if v, ok := overrideParameterMap["container"]; ok {
						overrideTranscodeParameter.Container = helper.String(v.(string))
					}
					if v, ok := overrideParameterMap["remove_video"]; ok {
						overrideTranscodeParameter.RemoveVideo = helper.IntUint64(v.(int))
					}
					if v, ok := overrideParameterMap["remove_audio"]; ok {
						overrideTranscodeParameter.RemoveAudio = helper.IntUint64(v.(int))
					}
					if videoTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "video_template"); ok {
						videoTemplateInfoForUpdate := mps.VideoTemplateInfoForUpdate{}
						if v, ok := videoTemplateMap["codec"]; ok {
							videoTemplateInfoForUpdate.Codec = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["fps"]; ok {
							videoTemplateInfoForUpdate.Fps = helper.IntInt64(v.(int))
						}
						if v, ok := videoTemplateMap["bitrate"]; ok {
							videoTemplateInfoForUpdate.Bitrate = helper.IntInt64(v.(int))
						}
						if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
							videoTemplateInfoForUpdate.ResolutionAdaptive = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["width"]; ok {
							videoTemplateInfoForUpdate.Width = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["height"]; ok {
							videoTemplateInfoForUpdate.Height = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["gop"]; ok {
							videoTemplateInfoForUpdate.Gop = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["fill_type"]; ok {
							videoTemplateInfoForUpdate.FillType = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["vcrf"]; ok {
							videoTemplateInfoForUpdate.Vcrf = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["content_adapt_stream"]; ok {
							videoTemplateInfoForUpdate.ContentAdaptStream = helper.IntUint64(v.(int))
						}
						overrideTranscodeParameter.VideoTemplate = &videoTemplateInfoForUpdate
					}
					if audioTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "audio_template"); ok {
						audioTemplateInfoForUpdate := mps.AudioTemplateInfoForUpdate{}
						if v, ok := audioTemplateMap["codec"]; ok {
							audioTemplateInfoForUpdate.Codec = helper.String(v.(string))
						}
						if v, ok := audioTemplateMap["bitrate"]; ok {
							audioTemplateInfoForUpdate.Bitrate = helper.IntInt64(v.(int))
						}
						if v, ok := audioTemplateMap["sample_rate"]; ok {
							audioTemplateInfoForUpdate.SampleRate = helper.IntUint64(v.(int))
						}
						if v, ok := audioTemplateMap["audio_channel"]; ok {
							audioTemplateInfoForUpdate.AudioChannel = helper.IntInt64(v.(int))
						}
						if v, ok := audioTemplateMap["stream_selects"]; ok {
							streamSelectsSet := v.(*schema.Set).List()
							for i := range streamSelectsSet {
								streamSelects := streamSelectsSet[i].(int)
								audioTemplateInfoForUpdate.StreamSelects = append(audioTemplateInfoForUpdate.StreamSelects, helper.IntInt64(streamSelects))
							}
						}
						overrideTranscodeParameter.AudioTemplate = &audioTemplateInfoForUpdate
					}
					if tEHDConfigMap, ok := helper.InterfaceToMap(overrideParameterMap, "tehd_config"); ok {
						tEHDConfigForUpdate := mps.TEHDConfigForUpdate{}
						if v, ok := tEHDConfigMap["type"]; ok {
							tEHDConfigForUpdate.Type = helper.String(v.(string))
						}
						if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
							tEHDConfigForUpdate.MaxVideoBitrate = helper.IntInt64(v.(int))
						}
						overrideTranscodeParameter.TEHDConfig = &tEHDConfigForUpdate
					}
					if subtitleTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "subtitle_template"); ok {
						subtitleTemplate := mps.SubtitleTemplate{}
						if v, ok := subtitleTemplateMap["path"]; ok {
							subtitleTemplate.Path = helper.String(v.(string))
						}
						if v, ok := subtitleTemplateMap["stream_index"]; ok {
							subtitleTemplate.StreamIndex = helper.IntInt64(v.(int))
						}
						if v, ok := subtitleTemplateMap["font_type"]; ok {
							subtitleTemplate.FontType = helper.String(v.(string))
						}
						if v, ok := subtitleTemplateMap["font_size"]; ok {
							subtitleTemplate.FontSize = helper.String(v.(string))
						}
						if v, ok := subtitleTemplateMap["font_color"]; ok {
							subtitleTemplate.FontColor = helper.String(v.(string))
						}
						if v, ok := subtitleTemplateMap["font_alpha"]; ok {
							subtitleTemplate.FontAlpha = helper.Float64(v.(float64))
						}
						overrideTranscodeParameter.SubtitleTemplate = &subtitleTemplate
					}
					transcodeTaskInput.OverrideParameter = &overrideTranscodeParameter
				}
				if v, ok := transcodeTaskSetMap["watermark_set"]; ok {
					for _, item := range v.([]interface{}) {
						watermarkSetMap := item.(map[string]interface{})
						watermarkInput := mps.WatermarkInput{}
						if v, ok := watermarkSetMap["definition"]; ok {
							watermarkInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
							rawWatermarkParameter := mps.RawWatermarkParameter{}
							if v, ok := rawParameterMap["type"]; ok {
								rawWatermarkParameter.Type = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["coordinate_origin"]; ok {
								rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["x_pos"]; ok {
								rawWatermarkParameter.XPos = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["y_pos"]; ok {
								rawWatermarkParameter.YPos = helper.String(v.(string))
							}
							if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
								rawImageWatermarkInput := mps.RawImageWatermarkInput{}
								if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := imageContentMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
										cosInputInfo := mps.CosInputInfo{}
										if v, ok := cosInputInfoMap["bucket"]; ok {
											cosInputInfo.Bucket = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["region"]; ok {
											cosInputInfo.Region = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["object"]; ok {
											cosInputInfo.Object = helper.String(v.(string))
										}
										mediaInputInfo.CosInputInfo = &cosInputInfo
									}
									if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									rawImageWatermarkInput.ImageContent = &mediaInputInfo
								}
								if v, ok := imageTemplateMap["width"]; ok {
									rawImageWatermarkInput.Width = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["height"]; ok {
									rawImageWatermarkInput.Height = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["repeat_type"]; ok {
									rawImageWatermarkInput.RepeatType = helper.String(v.(string))
								}
								rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
							}
							watermarkInput.RawParameter = &rawWatermarkParameter
						}
						if v, ok := watermarkSetMap["text_content"]; ok {
							watermarkInput.TextContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["svg_content"]; ok {
							watermarkInput.SvgContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["start_time_offset"]; ok {
							watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := watermarkSetMap["end_time_offset"]; ok {
							watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						transcodeTaskInput.WatermarkSet = append(transcodeTaskInput.WatermarkSet, &watermarkInput)
					}
				}
				if v, ok := transcodeTaskSetMap["mosaic_set"]; ok {
					for _, item := range v.([]interface{}) {
						mosaicSetMap := item.(map[string]interface{})
						mosaicInput := mps.MosaicInput{}
						if v, ok := mosaicSetMap["coordinate_origin"]; ok {
							mosaicInput.CoordinateOrigin = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["x_pos"]; ok {
							mosaicInput.XPos = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["y_pos"]; ok {
							mosaicInput.YPos = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["width"]; ok {
							mosaicInput.Width = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["height"]; ok {
							mosaicInput.Height = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["start_time_offset"]; ok {
							mosaicInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := mosaicSetMap["end_time_offset"]; ok {
							mosaicInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						transcodeTaskInput.MosaicSet = append(transcodeTaskInput.MosaicSet, &mosaicInput)
					}
				}
				if v, ok := transcodeTaskSetMap["start_time_offset"]; ok {
					transcodeTaskInput.StartTimeOffset = helper.Float64(v.(float64))
				}
				if v, ok := transcodeTaskSetMap["end_time_offset"]; ok {
					transcodeTaskInput.EndTimeOffset = helper.Float64(v.(float64))
				}
				if outputStorageMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					transcodeTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := transcodeTaskSetMap["output_object_path"]; ok {
					transcodeTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if v, ok := transcodeTaskSetMap["segment_object_name"]; ok {
					transcodeTaskInput.SegmentObjectName = helper.String(v.(string))
				}
				if objectNumberFormatMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "object_number_format"); ok {
					numberFormat := mps.NumberFormat{}
					if v, ok := objectNumberFormatMap["initial_value"]; ok {
						numberFormat.InitialValue = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["increment"]; ok {
						numberFormat.Increment = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["min_length"]; ok {
						numberFormat.MinLength = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["place_holder"]; ok {
						numberFormat.PlaceHolder = helper.String(v.(string))
					}
					transcodeTaskInput.ObjectNumberFormat = &numberFormat
				}
				if headTailParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "head_tail_parameter"); ok {
					headTailParameter := mps.HeadTailParameter{}
					if v, ok := headTailParameterMap["head_set"]; ok {
						for _, item := range v.([]interface{}) {
							headSetMap := item.(map[string]interface{})
							mediaInputInfo := mps.MediaInputInfo{}
							if v, ok := headSetMap["type"]; ok {
								mediaInputInfo.Type = helper.String(v.(string))
							}
							if cosInputInfoMap, ok := helper.InterfaceToMap(headSetMap, "cos_input_info"); ok {
								cosInputInfo := mps.CosInputInfo{}
								if v, ok := cosInputInfoMap["bucket"]; ok {
									cosInputInfo.Bucket = helper.String(v.(string))
								}
								if v, ok := cosInputInfoMap["region"]; ok {
									cosInputInfo.Region = helper.String(v.(string))
								}
								if v, ok := cosInputInfoMap["object"]; ok {
									cosInputInfo.Object = helper.String(v.(string))
								}
								mediaInputInfo.CosInputInfo = &cosInputInfo
							}
							if urlInputInfoMap, ok := helper.InterfaceToMap(headSetMap, "url_input_info"); ok {
								urlInputInfo := mps.UrlInputInfo{}
								if v, ok := urlInputInfoMap["url"]; ok {
									urlInputInfo.Url = helper.String(v.(string))
								}
								mediaInputInfo.UrlInputInfo = &urlInputInfo
							}
							headTailParameter.HeadSet = append(headTailParameter.HeadSet, &mediaInputInfo)
						}
					}
					if v, ok := headTailParameterMap["tail_set"]; ok {
						for _, item := range v.([]interface{}) {
							tailSetMap := item.(map[string]interface{})
							mediaInputInfo := mps.MediaInputInfo{}
							if v, ok := tailSetMap["type"]; ok {
								mediaInputInfo.Type = helper.String(v.(string))
							}
							if cosInputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "cos_input_info"); ok {
								cosInputInfo := mps.CosInputInfo{}
								if v, ok := cosInputInfoMap["bucket"]; ok {
									cosInputInfo.Bucket = helper.String(v.(string))
								}
								if v, ok := cosInputInfoMap["region"]; ok {
									cosInputInfo.Region = helper.String(v.(string))
								}
								if v, ok := cosInputInfoMap["object"]; ok {
									cosInputInfo.Object = helper.String(v.(string))
								}
								mediaInputInfo.CosInputInfo = &cosInputInfo
							}
							if urlInputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "url_input_info"); ok {
								urlInputInfo := mps.UrlInputInfo{}
								if v, ok := urlInputInfoMap["url"]; ok {
									urlInputInfo.Url = helper.String(v.(string))
								}
								mediaInputInfo.UrlInputInfo = &urlInputInfo
							}
							headTailParameter.TailSet = append(headTailParameter.TailSet, &mediaInputInfo)
						}
					}
					transcodeTaskInput.HeadTailParameter = &headTailParameter
				}
				mediaProcessTaskInput.TranscodeTaskSet = append(mediaProcessTaskInput.TranscodeTaskSet, &transcodeTaskInput)
			}
		}
		if v, ok := dMap["animated_graphic_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				animatedGraphicTaskSetMap := item.(map[string]interface{})
				animatedGraphicTaskInput := mps.AnimatedGraphicTaskInput{}
				if v, ok := animatedGraphicTaskSetMap["definition"]; ok {
					animatedGraphicTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := animatedGraphicTaskSetMap["start_time_offset"]; ok {
					animatedGraphicTaskInput.StartTimeOffset = helper.Float64(v.(float64))
				}
				if v, ok := animatedGraphicTaskSetMap["end_time_offset"]; ok {
					animatedGraphicTaskInput.EndTimeOffset = helper.Float64(v.(float64))
				}
				if outputStorageMap, ok := helper.InterfaceToMap(animatedGraphicTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					animatedGraphicTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := animatedGraphicTaskSetMap["output_object_path"]; ok {
					animatedGraphicTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				mediaProcessTaskInput.AnimatedGraphicTaskSet = append(mediaProcessTaskInput.AnimatedGraphicTaskSet, &animatedGraphicTaskInput)
			}
		}
		if v, ok := dMap["snapshot_by_time_offset_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				snapshotByTimeOffsetTaskSetMap := item.(map[string]interface{})
				snapshotByTimeOffsetTaskInput := mps.SnapshotByTimeOffsetTaskInput{}
				if v, ok := snapshotByTimeOffsetTaskSetMap["definition"]; ok {
					snapshotByTimeOffsetTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := snapshotByTimeOffsetTaskSetMap["ext_time_offset_set"]; ok {
					extTimeOffsetSetSet := v.(*schema.Set).List()
					for i := range extTimeOffsetSetSet {
						extTimeOffsetSet := extTimeOffsetSetSet[i].(string)
						snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet = append(snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet, &extTimeOffsetSet)
					}
				}
				if v, _ := d.GetOk("time_offset_set"); v != nil {
					timeOffsetSetSet := v.(*schema.Set).List()
					for i := range timeOffsetSetSet {
						timeOffsetSet := timeOffsetSetSet[i].(float64)
						snapshotByTimeOffsetTaskInput.TimeOffsetSet = append(snapshotByTimeOffsetTaskInput.TimeOffsetSet, &timeOffsetSet)
					}
				}

				if v, ok := snapshotByTimeOffsetTaskSetMap["watermark_set"]; ok {
					for _, item := range v.([]interface{}) {
						watermarkSetMap := item.(map[string]interface{})
						watermarkInput := mps.WatermarkInput{}
						if v, ok := watermarkSetMap["definition"]; ok {
							watermarkInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
							rawWatermarkParameter := mps.RawWatermarkParameter{}
							if v, ok := rawParameterMap["type"]; ok {
								rawWatermarkParameter.Type = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["coordinate_origin"]; ok {
								rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["x_pos"]; ok {
								rawWatermarkParameter.XPos = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["y_pos"]; ok {
								rawWatermarkParameter.YPos = helper.String(v.(string))
							}
							if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
								rawImageWatermarkInput := mps.RawImageWatermarkInput{}
								if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := imageContentMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
										cosInputInfo := mps.CosInputInfo{}
										if v, ok := cosInputInfoMap["bucket"]; ok {
											cosInputInfo.Bucket = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["region"]; ok {
											cosInputInfo.Region = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["object"]; ok {
											cosInputInfo.Object = helper.String(v.(string))
										}
										mediaInputInfo.CosInputInfo = &cosInputInfo
									}
									if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									rawImageWatermarkInput.ImageContent = &mediaInputInfo
								}
								if v, ok := imageTemplateMap["width"]; ok {
									rawImageWatermarkInput.Width = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["height"]; ok {
									rawImageWatermarkInput.Height = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["repeat_type"]; ok {
									rawImageWatermarkInput.RepeatType = helper.String(v.(string))
								}
								rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
							}
							watermarkInput.RawParameter = &rawWatermarkParameter
						}
						if v, ok := watermarkSetMap["text_content"]; ok {
							watermarkInput.TextContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["svg_content"]; ok {
							watermarkInput.SvgContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["start_time_offset"]; ok {
							watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := watermarkSetMap["end_time_offset"]; ok {
							watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						snapshotByTimeOffsetTaskInput.WatermarkSet = append(snapshotByTimeOffsetTaskInput.WatermarkSet, &watermarkInput)
					}
				}
				if outputStorageMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					snapshotByTimeOffsetTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := snapshotByTimeOffsetTaskSetMap["output_object_path"]; ok {
					snapshotByTimeOffsetTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if objectNumberFormatMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskSetMap, "object_number_format"); ok {
					numberFormat := mps.NumberFormat{}
					if v, ok := objectNumberFormatMap["initial_value"]; ok {
						numberFormat.InitialValue = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["increment"]; ok {
						numberFormat.Increment = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["min_length"]; ok {
						numberFormat.MinLength = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["place_holder"]; ok {
						numberFormat.PlaceHolder = helper.String(v.(string))
					}
					snapshotByTimeOffsetTaskInput.ObjectNumberFormat = &numberFormat
				}
				mediaProcessTaskInput.SnapshotByTimeOffsetTaskSet = append(mediaProcessTaskInput.SnapshotByTimeOffsetTaskSet, &snapshotByTimeOffsetTaskInput)
			}
		}
		if v, ok := dMap["sample_snapshot_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				sampleSnapshotTaskSetMap := item.(map[string]interface{})
				sampleSnapshotTaskInput := mps.SampleSnapshotTaskInput{}
				if v, ok := sampleSnapshotTaskSetMap["definition"]; ok {
					sampleSnapshotTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := sampleSnapshotTaskSetMap["watermark_set"]; ok {
					for _, item := range v.([]interface{}) {
						watermarkSetMap := item.(map[string]interface{})
						watermarkInput := mps.WatermarkInput{}
						if v, ok := watermarkSetMap["definition"]; ok {
							watermarkInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
							rawWatermarkParameter := mps.RawWatermarkParameter{}
							if v, ok := rawParameterMap["type"]; ok {
								rawWatermarkParameter.Type = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["coordinate_origin"]; ok {
								rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["x_pos"]; ok {
								rawWatermarkParameter.XPos = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["y_pos"]; ok {
								rawWatermarkParameter.YPos = helper.String(v.(string))
							}
							if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
								rawImageWatermarkInput := mps.RawImageWatermarkInput{}
								if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := imageContentMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
										cosInputInfo := mps.CosInputInfo{}
										if v, ok := cosInputInfoMap["bucket"]; ok {
											cosInputInfo.Bucket = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["region"]; ok {
											cosInputInfo.Region = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["object"]; ok {
											cosInputInfo.Object = helper.String(v.(string))
										}
										mediaInputInfo.CosInputInfo = &cosInputInfo
									}
									if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									rawImageWatermarkInput.ImageContent = &mediaInputInfo
								}
								if v, ok := imageTemplateMap["width"]; ok {
									rawImageWatermarkInput.Width = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["height"]; ok {
									rawImageWatermarkInput.Height = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["repeat_type"]; ok {
									rawImageWatermarkInput.RepeatType = helper.String(v.(string))
								}
								rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
							}
							watermarkInput.RawParameter = &rawWatermarkParameter
						}
						if v, ok := watermarkSetMap["text_content"]; ok {
							watermarkInput.TextContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["svg_content"]; ok {
							watermarkInput.SvgContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["start_time_offset"]; ok {
							watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := watermarkSetMap["end_time_offset"]; ok {
							watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						sampleSnapshotTaskInput.WatermarkSet = append(sampleSnapshotTaskInput.WatermarkSet, &watermarkInput)
					}
				}
				if outputStorageMap, ok := helper.InterfaceToMap(sampleSnapshotTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					sampleSnapshotTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := sampleSnapshotTaskSetMap["output_object_path"]; ok {
					sampleSnapshotTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if objectNumberFormatMap, ok := helper.InterfaceToMap(sampleSnapshotTaskSetMap, "object_number_format"); ok {
					numberFormat := mps.NumberFormat{}
					if v, ok := objectNumberFormatMap["initial_value"]; ok {
						numberFormat.InitialValue = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["increment"]; ok {
						numberFormat.Increment = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["min_length"]; ok {
						numberFormat.MinLength = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["place_holder"]; ok {
						numberFormat.PlaceHolder = helper.String(v.(string))
					}
					sampleSnapshotTaskInput.ObjectNumberFormat = &numberFormat
				}
				mediaProcessTaskInput.SampleSnapshotTaskSet = append(mediaProcessTaskInput.SampleSnapshotTaskSet, &sampleSnapshotTaskInput)
			}
		}
		if v, ok := dMap["image_sprite_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				imageSpriteTaskSetMap := item.(map[string]interface{})
				imageSpriteTaskInput := mps.ImageSpriteTaskInput{}
				if v, ok := imageSpriteTaskSetMap["definition"]; ok {
					imageSpriteTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if outputStorageMap, ok := helper.InterfaceToMap(imageSpriteTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					imageSpriteTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := imageSpriteTaskSetMap["output_object_path"]; ok {
					imageSpriteTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if v, ok := imageSpriteTaskSetMap["web_vtt_object_name"]; ok {
					imageSpriteTaskInput.WebVttObjectName = helper.String(v.(string))
				}
				if objectNumberFormatMap, ok := helper.InterfaceToMap(imageSpriteTaskSetMap, "object_number_format"); ok {
					numberFormat := mps.NumberFormat{}
					if v, ok := objectNumberFormatMap["initial_value"]; ok {
						numberFormat.InitialValue = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["increment"]; ok {
						numberFormat.Increment = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["min_length"]; ok {
						numberFormat.MinLength = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["place_holder"]; ok {
						numberFormat.PlaceHolder = helper.String(v.(string))
					}
					imageSpriteTaskInput.ObjectNumberFormat = &numberFormat
				}
				mediaProcessTaskInput.ImageSpriteTaskSet = append(mediaProcessTaskInput.ImageSpriteTaskSet, &imageSpriteTaskInput)
			}
		}
		if v, ok := dMap["adaptive_dynamic_streaming_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				adaptiveDynamicStreamingTaskSetMap := item.(map[string]interface{})
				adaptiveDynamicStreamingTaskInput := mps.AdaptiveDynamicStreamingTaskInput{}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["definition"]; ok {
					adaptiveDynamicStreamingTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["watermark_set"]; ok {
					for _, item := range v.([]interface{}) {
						watermarkSetMap := item.(map[string]interface{})
						watermarkInput := mps.WatermarkInput{}
						if v, ok := watermarkSetMap["definition"]; ok {
							watermarkInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
							rawWatermarkParameter := mps.RawWatermarkParameter{}
							if v, ok := rawParameterMap["type"]; ok {
								rawWatermarkParameter.Type = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["coordinate_origin"]; ok {
								rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["x_pos"]; ok {
								rawWatermarkParameter.XPos = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["y_pos"]; ok {
								rawWatermarkParameter.YPos = helper.String(v.(string))
							}
							if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
								rawImageWatermarkInput := mps.RawImageWatermarkInput{}
								if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := imageContentMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
										cosInputInfo := mps.CosInputInfo{}
										if v, ok := cosInputInfoMap["bucket"]; ok {
											cosInputInfo.Bucket = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["region"]; ok {
											cosInputInfo.Region = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["object"]; ok {
											cosInputInfo.Object = helper.String(v.(string))
										}
										mediaInputInfo.CosInputInfo = &cosInputInfo
									}
									if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									rawImageWatermarkInput.ImageContent = &mediaInputInfo
								}
								if v, ok := imageTemplateMap["width"]; ok {
									rawImageWatermarkInput.Width = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["height"]; ok {
									rawImageWatermarkInput.Height = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["repeat_type"]; ok {
									rawImageWatermarkInput.RepeatType = helper.String(v.(string))
								}
								rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
							}
							watermarkInput.RawParameter = &rawWatermarkParameter
						}
						if v, ok := watermarkSetMap["text_content"]; ok {
							watermarkInput.TextContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["svg_content"]; ok {
							watermarkInput.SvgContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["start_time_offset"]; ok {
							watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := watermarkSetMap["end_time_offset"]; ok {
							watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						adaptiveDynamicStreamingTaskInput.WatermarkSet = append(adaptiveDynamicStreamingTaskInput.WatermarkSet, &watermarkInput)
					}
				}
				if outputStorageMap, ok := helper.InterfaceToMap(adaptiveDynamicStreamingTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					adaptiveDynamicStreamingTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["output_object_path"]; ok {
					adaptiveDynamicStreamingTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["sub_stream_object_name"]; ok {
					adaptiveDynamicStreamingTaskInput.SubStreamObjectName = helper.String(v.(string))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["segment_object_name"]; ok {
					adaptiveDynamicStreamingTaskInput.SegmentObjectName = helper.String(v.(string))
				}
				mediaProcessTaskInput.AdaptiveDynamicStreamingTaskSet = append(mediaProcessTaskInput.AdaptiveDynamicStreamingTaskSet, &adaptiveDynamicStreamingTaskInput)
			}
		}
		request.MediaProcessTask = &mediaProcessTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_content_review_task"); ok {
		aiContentReviewTaskInput := mps.AiContentReviewTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiContentReviewTaskInput.Definition = helper.IntUint64(v.(int))
		}
		request.AiContentReviewTask = &aiContentReviewTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_analysis_task"); ok {
		aiAnalysisTaskInput := mps.AiAnalysisTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiAnalysisTaskInput.Definition = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["extended_parameter"]; ok {
			aiAnalysisTaskInput.ExtendedParameter = helper.String(v.(string))
		}
		request.AiAnalysisTask = &aiAnalysisTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_recognition_task"); ok {
		aiRecognitionTaskInput := mps.AiRecognitionTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiRecognitionTaskInput.Definition = helper.IntUint64(v.(int))
		}
		request.AiRecognitionTask = &aiRecognitionTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "task_notify_config"); ok {
		taskNotifyConfig := mps.TaskNotifyConfig{}
		if v, ok := dMap["cmq_model"]; ok {
			taskNotifyConfig.CmqModel = helper.String(v.(string))
		}
		if v, ok := dMap["cmq_region"]; ok {
			taskNotifyConfig.CmqRegion = helper.String(v.(string))
		}
		if v, ok := dMap["topic_name"]; ok {
			taskNotifyConfig.TopicName = helper.String(v.(string))
		}
		if v, ok := dMap["queue_name"]; ok {
			taskNotifyConfig.QueueName = helper.String(v.(string))
		}
		if v, ok := dMap["notify_mode"]; ok {
			taskNotifyConfig.NotifyMode = helper.String(v.(string))
		}
		if v, ok := dMap["notify_type"]; ok {
			taskNotifyConfig.NotifyType = helper.String(v.(string))
		}
		if v, ok := dMap["notify_url"]; ok {
			taskNotifyConfig.NotifyUrl = helper.String(v.(string))
		}
		request.TaskNotifyConfig = &taskNotifyConfig
	}

	if v, _ := d.GetOk("task_priority"); v != nil {
		request.TaskPriority = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateWorkflow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps workflow failed, reason:%+v", logId, err)
		return err
	}

	workflowId = *response.Response.WorkflowId
	d.SetId(helper.Int64ToStr(workflowId))

	return resourceTencentCloudMpsWorkflowRead(d, meta)
}

func resourceTencentCloudMpsWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_workflow.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	workflowId := d.Id()

	workflow, err := service.DescribeMpsWorkflowById(ctx, workflowId)
	if err != nil {
		return err
	}

	if workflow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsWorkflow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if workflow.WorkflowName != nil {
		_ = d.Set("workflow_name", workflow.WorkflowName)
	}

	if workflow.Trigger != nil {
		triggerMap := map[string]interface{}{}

		if workflow.Trigger.Type != nil {
			triggerMap["type"] = workflow.Trigger.Type
		}

		if workflow.Trigger.CosFileUploadTrigger != nil {
			cosFileUploadTriggerMap := map[string]interface{}{}

			if workflow.Trigger.CosFileUploadTrigger.Bucket != nil {
				cosFileUploadTriggerMap["bucket"] = workflow.Trigger.CosFileUploadTrigger.Bucket
			}

			if workflow.Trigger.CosFileUploadTrigger.Region != nil {
				cosFileUploadTriggerMap["region"] = workflow.Trigger.CosFileUploadTrigger.Region
			}

			if workflow.Trigger.CosFileUploadTrigger.Dir != nil {
				cosFileUploadTriggerMap["dir"] = workflow.Trigger.CosFileUploadTrigger.Dir
			}

			if workflow.Trigger.CosFileUploadTrigger.Formats != nil {
				cosFileUploadTriggerMap["formats"] = workflow.Trigger.CosFileUploadTrigger.Formats
			}

			triggerMap["cos_file_upload_trigger"] = []interface{}{cosFileUploadTriggerMap}
		}

		_ = d.Set("trigger", []interface{}{triggerMap})
	}

	if workflow.OutputStorage != nil {
		outputStorageMap := map[string]interface{}{}

		if workflow.OutputStorage.Type != nil {
			outputStorageMap["type"] = workflow.OutputStorage.Type
		}

		if workflow.OutputStorage.CosOutputStorage != nil {
			cosOutputStorageMap := map[string]interface{}{}

			if workflow.OutputStorage.CosOutputStorage.Bucket != nil {
				cosOutputStorageMap["bucket"] = workflow.OutputStorage.CosOutputStorage.Bucket
			}

			if workflow.OutputStorage.CosOutputStorage.Region != nil {
				cosOutputStorageMap["region"] = workflow.OutputStorage.CosOutputStorage.Region
			}

			outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
		}

		_ = d.Set("output_storage", []interface{}{outputStorageMap})
	}

	if workflow.OutputDir != nil {
		_ = d.Set("output_dir", workflow.OutputDir)
	}

	if workflow.MediaProcessTask != nil {
		mediaProcessTaskMap := map[string]interface{}{}

		if workflow.MediaProcessTask.TranscodeTaskSet != nil {
			transcodeTaskSetList := []interface{}{}
			for _, transcodeTaskSet := range workflow.MediaProcessTask.TranscodeTaskSet {
				transcodeTaskSetMap := map[string]interface{}{}

				if transcodeTaskSet.Definition != nil {
					transcodeTaskSetMap["definition"] = transcodeTaskSet.Definition
				}

				if transcodeTaskSet.RawParameter != nil {
					rawParameterMap := map[string]interface{}{}

					if transcodeTaskSet.RawParameter.Container != nil {
						rawParameterMap["container"] = transcodeTaskSet.RawParameter.Container
					}

					if transcodeTaskSet.RawParameter.RemoveVideo != nil {
						rawParameterMap["remove_video"] = transcodeTaskSet.RawParameter.RemoveVideo
					}

					if transcodeTaskSet.RawParameter.RemoveAudio != nil {
						rawParameterMap["remove_audio"] = transcodeTaskSet.RawParameter.RemoveAudio
					}

					if transcodeTaskSet.RawParameter.VideoTemplate != nil {
						videoTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.RawParameter.VideoTemplate.Codec != nil {
							videoTemplateMap["codec"] = transcodeTaskSet.RawParameter.VideoTemplate.Codec
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Fps != nil {
							videoTemplateMap["fps"] = transcodeTaskSet.RawParameter.VideoTemplate.Fps
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Bitrate != nil {
							videoTemplateMap["bitrate"] = transcodeTaskSet.RawParameter.VideoTemplate.Bitrate
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.ResolutionAdaptive != nil {
							videoTemplateMap["resolution_adaptive"] = transcodeTaskSet.RawParameter.VideoTemplate.ResolutionAdaptive
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Width != nil {
							videoTemplateMap["width"] = transcodeTaskSet.RawParameter.VideoTemplate.Width
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Height != nil {
							videoTemplateMap["height"] = transcodeTaskSet.RawParameter.VideoTemplate.Height
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Gop != nil {
							videoTemplateMap["gop"] = transcodeTaskSet.RawParameter.VideoTemplate.Gop
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.FillType != nil {
							videoTemplateMap["fill_type"] = transcodeTaskSet.RawParameter.VideoTemplate.FillType
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Vcrf != nil {
							videoTemplateMap["vcrf"] = transcodeTaskSet.RawParameter.VideoTemplate.Vcrf
						}

						rawParameterMap["video_template"] = []interface{}{videoTemplateMap}
					}

					if transcodeTaskSet.RawParameter.AudioTemplate != nil {
						audioTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.RawParameter.AudioTemplate.Codec != nil {
							audioTemplateMap["codec"] = transcodeTaskSet.RawParameter.AudioTemplate.Codec
						}

						if transcodeTaskSet.RawParameter.AudioTemplate.Bitrate != nil {
							audioTemplateMap["bitrate"] = transcodeTaskSet.RawParameter.AudioTemplate.Bitrate
						}

						if transcodeTaskSet.RawParameter.AudioTemplate.SampleRate != nil {
							audioTemplateMap["sample_rate"] = transcodeTaskSet.RawParameter.AudioTemplate.SampleRate
						}

						if transcodeTaskSet.RawParameter.AudioTemplate.AudioChannel != nil {
							audioTemplateMap["audio_channel"] = transcodeTaskSet.RawParameter.AudioTemplate.AudioChannel
						}

						rawParameterMap["audio_template"] = []interface{}{audioTemplateMap}
					}

					if transcodeTaskSet.RawParameter.TEHDConfig != nil {
						tEHDConfigMap := map[string]interface{}{}

						if transcodeTaskSet.RawParameter.TEHDConfig.Type != nil {
							tEHDConfigMap["type"] = transcodeTaskSet.RawParameter.TEHDConfig.Type
						}

						if transcodeTaskSet.RawParameter.TEHDConfig.MaxVideoBitrate != nil {
							tEHDConfigMap["max_video_bitrate"] = transcodeTaskSet.RawParameter.TEHDConfig.MaxVideoBitrate
						}

						rawParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
					}

					transcodeTaskSetMap["raw_parameter"] = []interface{}{rawParameterMap}
				}

				if transcodeTaskSet.OverrideParameter != nil {
					overrideParameterMap := map[string]interface{}{}

					if transcodeTaskSet.OverrideParameter.Container != nil {
						overrideParameterMap["container"] = transcodeTaskSet.OverrideParameter.Container
					}

					if transcodeTaskSet.OverrideParameter.RemoveVideo != nil {
						overrideParameterMap["remove_video"] = transcodeTaskSet.OverrideParameter.RemoveVideo
					}

					if transcodeTaskSet.OverrideParameter.RemoveAudio != nil {
						overrideParameterMap["remove_audio"] = transcodeTaskSet.OverrideParameter.RemoveAudio
					}

					if transcodeTaskSet.OverrideParameter.VideoTemplate != nil {
						videoTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Codec != nil {
							videoTemplateMap["codec"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Codec
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Fps != nil {
							videoTemplateMap["fps"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Fps
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Bitrate != nil {
							videoTemplateMap["bitrate"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Bitrate
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.ResolutionAdaptive != nil {
							videoTemplateMap["resolution_adaptive"] = transcodeTaskSet.OverrideParameter.VideoTemplate.ResolutionAdaptive
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Width != nil {
							videoTemplateMap["width"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Width
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Height != nil {
							videoTemplateMap["height"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Height
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Gop != nil {
							videoTemplateMap["gop"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Gop
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.FillType != nil {
							videoTemplateMap["fill_type"] = transcodeTaskSet.OverrideParameter.VideoTemplate.FillType
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Vcrf != nil {
							videoTemplateMap["vcrf"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Vcrf
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.ContentAdaptStream != nil {
							videoTemplateMap["content_adapt_stream"] = transcodeTaskSet.OverrideParameter.VideoTemplate.ContentAdaptStream
						}

						overrideParameterMap["video_template"] = []interface{}{videoTemplateMap}
					}

					if transcodeTaskSet.OverrideParameter.AudioTemplate != nil {
						audioTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.Codec != nil {
							audioTemplateMap["codec"] = transcodeTaskSet.OverrideParameter.AudioTemplate.Codec
						}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.Bitrate != nil {
							audioTemplateMap["bitrate"] = transcodeTaskSet.OverrideParameter.AudioTemplate.Bitrate
						}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.SampleRate != nil {
							audioTemplateMap["sample_rate"] = transcodeTaskSet.OverrideParameter.AudioTemplate.SampleRate
						}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.AudioChannel != nil {
							audioTemplateMap["audio_channel"] = transcodeTaskSet.OverrideParameter.AudioTemplate.AudioChannel
						}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.StreamSelects != nil {
							audioTemplateMap["stream_selects"] = transcodeTaskSet.OverrideParameter.AudioTemplate.StreamSelects
						}

						overrideParameterMap["audio_template"] = []interface{}{audioTemplateMap}
					}

					if transcodeTaskSet.OverrideParameter.TEHDConfig != nil {
						tEHDConfigMap := map[string]interface{}{}

						if transcodeTaskSet.OverrideParameter.TEHDConfig.Type != nil {
							tEHDConfigMap["type"] = transcodeTaskSet.OverrideParameter.TEHDConfig.Type
						}

						if transcodeTaskSet.OverrideParameter.TEHDConfig.MaxVideoBitrate != nil {
							tEHDConfigMap["max_video_bitrate"] = transcodeTaskSet.OverrideParameter.TEHDConfig.MaxVideoBitrate
						}

						overrideParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
					}

					if transcodeTaskSet.OverrideParameter.SubtitleTemplate != nil {
						subtitleTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.Path != nil {
							subtitleTemplateMap["path"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.Path
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.StreamIndex != nil {
							subtitleTemplateMap["stream_index"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.StreamIndex
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontType != nil {
							subtitleTemplateMap["font_type"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontType
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontSize != nil {
							subtitleTemplateMap["font_size"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontSize
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontColor != nil {
							subtitleTemplateMap["font_color"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontColor
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontAlpha != nil {
							subtitleTemplateMap["font_alpha"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontAlpha
						}

						overrideParameterMap["subtitle_template"] = []interface{}{subtitleTemplateMap}
					}

					transcodeTaskSetMap["override_parameter"] = []interface{}{overrideParameterMap}
				}

				if transcodeTaskSet.WatermarkSet != nil {
					watermarkSetList := []interface{}{}
					for _, watermarkSet := range transcodeTaskSet.WatermarkSet {
						watermarkSetMap := map[string]interface{}{}

						if watermarkSet.Definition != nil {
							watermarkSetMap["definition"] = watermarkSet.Definition
						}

						if watermarkSet.RawParameter != nil {
							rawParameterMap := map[string]interface{}{}

							if watermarkSet.RawParameter.Type != nil {
								rawParameterMap["type"] = watermarkSet.RawParameter.Type
							}

							if watermarkSet.RawParameter.CoordinateOrigin != nil {
								rawParameterMap["coordinate_origin"] = watermarkSet.RawParameter.CoordinateOrigin
							}

							if watermarkSet.RawParameter.XPos != nil {
								rawParameterMap["x_pos"] = watermarkSet.RawParameter.XPos
							}

							if watermarkSet.RawParameter.YPos != nil {
								rawParameterMap["y_pos"] = watermarkSet.RawParameter.YPos
							}

							if watermarkSet.RawParameter.ImageTemplate != nil {
								imageTemplateMap := map[string]interface{}{}

								if watermarkSet.RawParameter.ImageTemplate.ImageContent != nil {
									imageContentMap := map[string]interface{}{}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.Type != nil {
										imageContentMap["type"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.Type
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo != nil {
										cosInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket != nil {
											cosInputInfoMap["bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region != nil {
											cosInputInfoMap["region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object != nil {
											cosInputInfoMap["object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object
										}

										imageContentMap["cos_input_info"] = []interface{}{cosInputInfoMap}
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo != nil {
										urlInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url != nil {
											urlInputInfoMap["url"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url
										}

										imageContentMap["url_input_info"] = []interface{}{urlInputInfoMap}
									}

									imageTemplateMap["image_content"] = []interface{}{imageContentMap}
								}

								if watermarkSet.RawParameter.ImageTemplate.Width != nil {
									imageTemplateMap["width"] = watermarkSet.RawParameter.ImageTemplate.Width
								}

								if watermarkSet.RawParameter.ImageTemplate.Height != nil {
									imageTemplateMap["height"] = watermarkSet.RawParameter.ImageTemplate.Height
								}

								if watermarkSet.RawParameter.ImageTemplate.RepeatType != nil {
									imageTemplateMap["repeat_type"] = watermarkSet.RawParameter.ImageTemplate.RepeatType
								}

								rawParameterMap["image_template"] = []interface{}{imageTemplateMap}
							}

							watermarkSetMap["raw_parameter"] = []interface{}{rawParameterMap}
						}

						if watermarkSet.TextContent != nil {
							watermarkSetMap["text_content"] = watermarkSet.TextContent
						}

						if watermarkSet.SvgContent != nil {
							watermarkSetMap["svg_content"] = watermarkSet.SvgContent
						}

						if watermarkSet.StartTimeOffset != nil {
							watermarkSetMap["start_time_offset"] = watermarkSet.StartTimeOffset
						}

						if watermarkSet.EndTimeOffset != nil {
							watermarkSetMap["end_time_offset"] = watermarkSet.EndTimeOffset
						}

						watermarkSetList = append(watermarkSetList, watermarkSetMap)
					}

					transcodeTaskSetMap["watermark_set"] = watermarkSetList
				}

				if transcodeTaskSet.MosaicSet != nil {
					mosaicSetList := []interface{}{}
					for _, mosaicSet := range transcodeTaskSet.MosaicSet {
						mosaicSetMap := map[string]interface{}{}

						if mosaicSet.CoordinateOrigin != nil {
							mosaicSetMap["coordinate_origin"] = mosaicSet.CoordinateOrigin
						}

						if mosaicSet.XPos != nil {
							mosaicSetMap["x_pos"] = mosaicSet.XPos
						}

						if mosaicSet.YPos != nil {
							mosaicSetMap["y_pos"] = mosaicSet.YPos
						}

						if mosaicSet.Width != nil {
							mosaicSetMap["width"] = mosaicSet.Width
						}

						if mosaicSet.Height != nil {
							mosaicSetMap["height"] = mosaicSet.Height
						}

						if mosaicSet.StartTimeOffset != nil {
							mosaicSetMap["start_time_offset"] = mosaicSet.StartTimeOffset
						}

						if mosaicSet.EndTimeOffset != nil {
							mosaicSetMap["end_time_offset"] = mosaicSet.EndTimeOffset
						}

						mosaicSetList = append(mosaicSetList, mosaicSetMap)
					}

					transcodeTaskSetMap["mosaic_set"] = mosaicSetList
				}

				if transcodeTaskSet.StartTimeOffset != nil {
					transcodeTaskSetMap["start_time_offset"] = transcodeTaskSet.StartTimeOffset
				}

				if transcodeTaskSet.EndTimeOffset != nil {
					transcodeTaskSetMap["end_time_offset"] = transcodeTaskSet.EndTimeOffset
				}

				if transcodeTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if transcodeTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = transcodeTaskSet.OutputStorage.Type
					}

					if transcodeTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if transcodeTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = transcodeTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if transcodeTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = transcodeTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					transcodeTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if transcodeTaskSet.OutputObjectPath != nil {
					transcodeTaskSetMap["output_object_path"] = transcodeTaskSet.OutputObjectPath
				}

				if transcodeTaskSet.SegmentObjectName != nil {
					transcodeTaskSetMap["segment_object_name"] = transcodeTaskSet.SegmentObjectName
				}

				if transcodeTaskSet.ObjectNumberFormat != nil {
					objectNumberFormatMap := map[string]interface{}{}

					if transcodeTaskSet.ObjectNumberFormat.InitialValue != nil {
						objectNumberFormatMap["initial_value"] = transcodeTaskSet.ObjectNumberFormat.InitialValue
					}

					if transcodeTaskSet.ObjectNumberFormat.Increment != nil {
						objectNumberFormatMap["increment"] = transcodeTaskSet.ObjectNumberFormat.Increment
					}

					if transcodeTaskSet.ObjectNumberFormat.MinLength != nil {
						objectNumberFormatMap["min_length"] = transcodeTaskSet.ObjectNumberFormat.MinLength
					}

					if transcodeTaskSet.ObjectNumberFormat.PlaceHolder != nil {
						objectNumberFormatMap["place_holder"] = transcodeTaskSet.ObjectNumberFormat.PlaceHolder
					}

					transcodeTaskSetMap["object_number_format"] = []interface{}{objectNumberFormatMap}
				}

				if transcodeTaskSet.HeadTailParameter != nil {
					headTailParameterMap := map[string]interface{}{}

					if transcodeTaskSet.HeadTailParameter.HeadSet != nil {
						headSetList := []interface{}{}
						for _, headSet := range transcodeTaskSet.HeadTailParameter.HeadSet {
							headSetMap := map[string]interface{}{}

							if headSet.Type != nil {
								headSetMap["type"] = headSet.Type
							}

							if headSet.CosInputInfo != nil {
								cosInputInfoMap := map[string]interface{}{}

								if headSet.CosInputInfo.Bucket != nil {
									cosInputInfoMap["bucket"] = headSet.CosInputInfo.Bucket
								}

								if headSet.CosInputInfo.Region != nil {
									cosInputInfoMap["region"] = headSet.CosInputInfo.Region
								}

								if headSet.CosInputInfo.Object != nil {
									cosInputInfoMap["object"] = headSet.CosInputInfo.Object
								}

								headSetMap["cos_input_info"] = []interface{}{cosInputInfoMap}
							}

							if headSet.UrlInputInfo != nil {
								urlInputInfoMap := map[string]interface{}{}

								if headSet.UrlInputInfo.Url != nil {
									urlInputInfoMap["url"] = headSet.UrlInputInfo.Url
								}

								headSetMap["url_input_info"] = []interface{}{urlInputInfoMap}
							}

							headSetList = append(headSetList, headSetMap)
						}

						headTailParameterMap["head_set"] = headSetList
					}

					if transcodeTaskSet.HeadTailParameter.TailSet != nil {
						tailSetList := []interface{}{}
						for _, tailSet := range transcodeTaskSet.HeadTailParameter.TailSet {
							tailSetMap := map[string]interface{}{}

							if tailSet.Type != nil {
								tailSetMap["type"] = tailSet.Type
							}

							if tailSet.CosInputInfo != nil {
								cosInputInfoMap := map[string]interface{}{}

								if tailSet.CosInputInfo.Bucket != nil {
									cosInputInfoMap["bucket"] = tailSet.CosInputInfo.Bucket
								}

								if tailSet.CosInputInfo.Region != nil {
									cosInputInfoMap["region"] = tailSet.CosInputInfo.Region
								}

								if tailSet.CosInputInfo.Object != nil {
									cosInputInfoMap["object"] = tailSet.CosInputInfo.Object
								}

								tailSetMap["cos_input_info"] = []interface{}{cosInputInfoMap}
							}

							if tailSet.UrlInputInfo != nil {
								urlInputInfoMap := map[string]interface{}{}

								if tailSet.UrlInputInfo.Url != nil {
									urlInputInfoMap["url"] = tailSet.UrlInputInfo.Url
								}

								tailSetMap["url_input_info"] = []interface{}{urlInputInfoMap}
							}

							tailSetList = append(tailSetList, tailSetMap)
						}

						headTailParameterMap["tail_set"] = tailSetList
					}

					transcodeTaskSetMap["head_tail_parameter"] = []interface{}{headTailParameterMap}
				}

				transcodeTaskSetList = append(transcodeTaskSetList, transcodeTaskSetMap)
			}

			mediaProcessTaskMap["transcode_task_set"] = transcodeTaskSetList
		}

		if workflow.MediaProcessTask.AnimatedGraphicTaskSet != nil {
			animatedGraphicTaskSetList := []interface{}{}
			for _, animatedGraphicTaskSet := range workflow.MediaProcessTask.AnimatedGraphicTaskSet {
				animatedGraphicTaskSetMap := map[string]interface{}{}

				if animatedGraphicTaskSet.Definition != nil {
					animatedGraphicTaskSetMap["definition"] = animatedGraphicTaskSet.Definition
				}

				if animatedGraphicTaskSet.StartTimeOffset != nil {
					animatedGraphicTaskSetMap["start_time_offset"] = animatedGraphicTaskSet.StartTimeOffset
				}

				if animatedGraphicTaskSet.EndTimeOffset != nil {
					animatedGraphicTaskSetMap["end_time_offset"] = animatedGraphicTaskSet.EndTimeOffset
				}

				if animatedGraphicTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if animatedGraphicTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = animatedGraphicTaskSet.OutputStorage.Type
					}

					if animatedGraphicTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if animatedGraphicTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = animatedGraphicTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if animatedGraphicTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = animatedGraphicTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					animatedGraphicTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if animatedGraphicTaskSet.OutputObjectPath != nil {
					animatedGraphicTaskSetMap["output_object_path"] = animatedGraphicTaskSet.OutputObjectPath
				}

				animatedGraphicTaskSetList = append(animatedGraphicTaskSetList, animatedGraphicTaskSetMap)
			}

			mediaProcessTaskMap["animated_graphic_task_set"] = animatedGraphicTaskSetList
		}

		if workflow.MediaProcessTask.SnapshotByTimeOffsetTaskSet != nil {
			snapshotByTimeOffsetTaskSetList := []interface{}{}
			for _, snapshotByTimeOffsetTaskSet := range workflow.MediaProcessTask.SnapshotByTimeOffsetTaskSet {
				snapshotByTimeOffsetTaskSetMap := map[string]interface{}{}

				if snapshotByTimeOffsetTaskSet.Definition != nil {
					snapshotByTimeOffsetTaskSetMap["definition"] = snapshotByTimeOffsetTaskSet.Definition
				}

				if snapshotByTimeOffsetTaskSet.ExtTimeOffsetSet != nil {
					snapshotByTimeOffsetTaskSetMap["ext_time_offset_set"] = snapshotByTimeOffsetTaskSet.ExtTimeOffsetSet
				}

				if snapshotByTimeOffsetTaskSet.TimeOffsetSet != nil {
					snapshotByTimeOffsetTaskSetMap["time_offset_set"] = snapshotByTimeOffsetTaskSet.TimeOffsetSet
				}

				if snapshotByTimeOffsetTaskSet.WatermarkSet != nil {
					watermarkSetList := []interface{}{}
					for _, watermarkSet := range snapshotByTimeOffsetTaskSet.WatermarkSet {
						watermarkSetMap := map[string]interface{}{}

						if watermarkSet.Definition != nil {
							watermarkSetMap["definition"] = watermarkSet.Definition
						}

						if watermarkSet.RawParameter != nil {
							rawParameterMap := map[string]interface{}{}

							if watermarkSet.RawParameter.Type != nil {
								rawParameterMap["type"] = watermarkSet.RawParameter.Type
							}

							if watermarkSet.RawParameter.CoordinateOrigin != nil {
								rawParameterMap["coordinate_origin"] = watermarkSet.RawParameter.CoordinateOrigin
							}

							if watermarkSet.RawParameter.XPos != nil {
								rawParameterMap["x_pos"] = watermarkSet.RawParameter.XPos
							}

							if watermarkSet.RawParameter.YPos != nil {
								rawParameterMap["y_pos"] = watermarkSet.RawParameter.YPos
							}

							if watermarkSet.RawParameter.ImageTemplate != nil {
								imageTemplateMap := map[string]interface{}{}

								if watermarkSet.RawParameter.ImageTemplate.ImageContent != nil {
									imageContentMap := map[string]interface{}{}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.Type != nil {
										imageContentMap["type"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.Type
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo != nil {
										cosInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket != nil {
											cosInputInfoMap["bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region != nil {
											cosInputInfoMap["region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object != nil {
											cosInputInfoMap["object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object
										}

										imageContentMap["cos_input_info"] = []interface{}{cosInputInfoMap}
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo != nil {
										urlInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url != nil {
											urlInputInfoMap["url"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url
										}

										imageContentMap["url_input_info"] = []interface{}{urlInputInfoMap}
									}

									imageTemplateMap["image_content"] = []interface{}{imageContentMap}
								}

								if watermarkSet.RawParameter.ImageTemplate.Width != nil {
									imageTemplateMap["width"] = watermarkSet.RawParameter.ImageTemplate.Width
								}

								if watermarkSet.RawParameter.ImageTemplate.Height != nil {
									imageTemplateMap["height"] = watermarkSet.RawParameter.ImageTemplate.Height
								}

								if watermarkSet.RawParameter.ImageTemplate.RepeatType != nil {
									imageTemplateMap["repeat_type"] = watermarkSet.RawParameter.ImageTemplate.RepeatType
								}

								rawParameterMap["image_template"] = []interface{}{imageTemplateMap}
							}

							watermarkSetMap["raw_parameter"] = []interface{}{rawParameterMap}
						}

						if watermarkSet.TextContent != nil {
							watermarkSetMap["text_content"] = watermarkSet.TextContent
						}

						if watermarkSet.SvgContent != nil {
							watermarkSetMap["svg_content"] = watermarkSet.SvgContent
						}

						if watermarkSet.StartTimeOffset != nil {
							watermarkSetMap["start_time_offset"] = watermarkSet.StartTimeOffset
						}

						if watermarkSet.EndTimeOffset != nil {
							watermarkSetMap["end_time_offset"] = watermarkSet.EndTimeOffset
						}

						watermarkSetList = append(watermarkSetList, watermarkSetMap)
					}

					snapshotByTimeOffsetTaskSetMap["watermark_set"] = watermarkSetList
				}

				if snapshotByTimeOffsetTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if snapshotByTimeOffsetTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = snapshotByTimeOffsetTaskSet.OutputStorage.Type
					}

					if snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					snapshotByTimeOffsetTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if snapshotByTimeOffsetTaskSet.OutputObjectPath != nil {
					snapshotByTimeOffsetTaskSetMap["output_object_path"] = snapshotByTimeOffsetTaskSet.OutputObjectPath
				}

				if snapshotByTimeOffsetTaskSet.ObjectNumberFormat != nil {
					objectNumberFormatMap := map[string]interface{}{}

					if snapshotByTimeOffsetTaskSet.ObjectNumberFormat.InitialValue != nil {
						objectNumberFormatMap["initial_value"] = snapshotByTimeOffsetTaskSet.ObjectNumberFormat.InitialValue
					}

					if snapshotByTimeOffsetTaskSet.ObjectNumberFormat.Increment != nil {
						objectNumberFormatMap["increment"] = snapshotByTimeOffsetTaskSet.ObjectNumberFormat.Increment
					}

					if snapshotByTimeOffsetTaskSet.ObjectNumberFormat.MinLength != nil {
						objectNumberFormatMap["min_length"] = snapshotByTimeOffsetTaskSet.ObjectNumberFormat.MinLength
					}

					if snapshotByTimeOffsetTaskSet.ObjectNumberFormat.PlaceHolder != nil {
						objectNumberFormatMap["place_holder"] = snapshotByTimeOffsetTaskSet.ObjectNumberFormat.PlaceHolder
					}

					snapshotByTimeOffsetTaskSetMap["object_number_format"] = []interface{}{objectNumberFormatMap}
				}

				snapshotByTimeOffsetTaskSetList = append(snapshotByTimeOffsetTaskSetList, snapshotByTimeOffsetTaskSetMap)
			}

			mediaProcessTaskMap["snapshot_by_time_offset_task_set"] = snapshotByTimeOffsetTaskSetList
		}

		if workflow.MediaProcessTask.SampleSnapshotTaskSet != nil {
			sampleSnapshotTaskSetList := []interface{}{}
			for _, sampleSnapshotTaskSet := range workflow.MediaProcessTask.SampleSnapshotTaskSet {
				sampleSnapshotTaskSetMap := map[string]interface{}{}

				if sampleSnapshotTaskSet.Definition != nil {
					sampleSnapshotTaskSetMap["definition"] = sampleSnapshotTaskSet.Definition
				}

				if sampleSnapshotTaskSet.WatermarkSet != nil {
					watermarkSetList := []interface{}{}
					for _, watermarkSet := range sampleSnapshotTaskSet.WatermarkSet {
						watermarkSetMap := map[string]interface{}{}

						if watermarkSet.Definition != nil {
							watermarkSetMap["definition"] = watermarkSet.Definition
						}

						if watermarkSet.RawParameter != nil {
							rawParameterMap := map[string]interface{}{}

							if watermarkSet.RawParameter.Type != nil {
								rawParameterMap["type"] = watermarkSet.RawParameter.Type
							}

							if watermarkSet.RawParameter.CoordinateOrigin != nil {
								rawParameterMap["coordinate_origin"] = watermarkSet.RawParameter.CoordinateOrigin
							}

							if watermarkSet.RawParameter.XPos != nil {
								rawParameterMap["x_pos"] = watermarkSet.RawParameter.XPos
							}

							if watermarkSet.RawParameter.YPos != nil {
								rawParameterMap["y_pos"] = watermarkSet.RawParameter.YPos
							}

							if watermarkSet.RawParameter.ImageTemplate != nil {
								imageTemplateMap := map[string]interface{}{}

								if watermarkSet.RawParameter.ImageTemplate.ImageContent != nil {
									imageContentMap := map[string]interface{}{}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.Type != nil {
										imageContentMap["type"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.Type
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo != nil {
										cosInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket != nil {
											cosInputInfoMap["bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region != nil {
											cosInputInfoMap["region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object != nil {
											cosInputInfoMap["object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object
										}

										imageContentMap["cos_input_info"] = []interface{}{cosInputInfoMap}
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo != nil {
										urlInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url != nil {
											urlInputInfoMap["url"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url
										}

										imageContentMap["url_input_info"] = []interface{}{urlInputInfoMap}
									}

									imageTemplateMap["image_content"] = []interface{}{imageContentMap}
								}

								if watermarkSet.RawParameter.ImageTemplate.Width != nil {
									imageTemplateMap["width"] = watermarkSet.RawParameter.ImageTemplate.Width
								}

								if watermarkSet.RawParameter.ImageTemplate.Height != nil {
									imageTemplateMap["height"] = watermarkSet.RawParameter.ImageTemplate.Height
								}

								if watermarkSet.RawParameter.ImageTemplate.RepeatType != nil {
									imageTemplateMap["repeat_type"] = watermarkSet.RawParameter.ImageTemplate.RepeatType
								}

								rawParameterMap["image_template"] = []interface{}{imageTemplateMap}
							}

							watermarkSetMap["raw_parameter"] = []interface{}{rawParameterMap}
						}

						if watermarkSet.TextContent != nil {
							watermarkSetMap["text_content"] = watermarkSet.TextContent
						}

						if watermarkSet.SvgContent != nil {
							watermarkSetMap["svg_content"] = watermarkSet.SvgContent
						}

						if watermarkSet.StartTimeOffset != nil {
							watermarkSetMap["start_time_offset"] = watermarkSet.StartTimeOffset
						}

						if watermarkSet.EndTimeOffset != nil {
							watermarkSetMap["end_time_offset"] = watermarkSet.EndTimeOffset
						}

						watermarkSetList = append(watermarkSetList, watermarkSetMap)
					}

					sampleSnapshotTaskSetMap["watermark_set"] = watermarkSetList
				}

				if sampleSnapshotTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if sampleSnapshotTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = sampleSnapshotTaskSet.OutputStorage.Type
					}

					if sampleSnapshotTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if sampleSnapshotTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = sampleSnapshotTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if sampleSnapshotTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = sampleSnapshotTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					sampleSnapshotTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if sampleSnapshotTaskSet.OutputObjectPath != nil {
					sampleSnapshotTaskSetMap["output_object_path"] = sampleSnapshotTaskSet.OutputObjectPath
				}

				if sampleSnapshotTaskSet.ObjectNumberFormat != nil {
					objectNumberFormatMap := map[string]interface{}{}

					if sampleSnapshotTaskSet.ObjectNumberFormat.InitialValue != nil {
						objectNumberFormatMap["initial_value"] = sampleSnapshotTaskSet.ObjectNumberFormat.InitialValue
					}

					if sampleSnapshotTaskSet.ObjectNumberFormat.Increment != nil {
						objectNumberFormatMap["increment"] = sampleSnapshotTaskSet.ObjectNumberFormat.Increment
					}

					if sampleSnapshotTaskSet.ObjectNumberFormat.MinLength != nil {
						objectNumberFormatMap["min_length"] = sampleSnapshotTaskSet.ObjectNumberFormat.MinLength
					}

					if sampleSnapshotTaskSet.ObjectNumberFormat.PlaceHolder != nil {
						objectNumberFormatMap["place_holder"] = sampleSnapshotTaskSet.ObjectNumberFormat.PlaceHolder
					}

					sampleSnapshotTaskSetMap["object_number_format"] = []interface{}{objectNumberFormatMap}
				}

				sampleSnapshotTaskSetList = append(sampleSnapshotTaskSetList, sampleSnapshotTaskSetMap)
			}

			mediaProcessTaskMap["sample_snapshot_task_set"] = sampleSnapshotTaskSetList
		}

		if workflow.MediaProcessTask.ImageSpriteTaskSet != nil {
			imageSpriteTaskSetList := []interface{}{}
			for _, imageSpriteTaskSet := range workflow.MediaProcessTask.ImageSpriteTaskSet {
				imageSpriteTaskSetMap := map[string]interface{}{}

				if imageSpriteTaskSet.Definition != nil {
					imageSpriteTaskSetMap["definition"] = imageSpriteTaskSet.Definition
				}

				if imageSpriteTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if imageSpriteTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = imageSpriteTaskSet.OutputStorage.Type
					}

					if imageSpriteTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if imageSpriteTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = imageSpriteTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if imageSpriteTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = imageSpriteTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					imageSpriteTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if imageSpriteTaskSet.OutputObjectPath != nil {
					imageSpriteTaskSetMap["output_object_path"] = imageSpriteTaskSet.OutputObjectPath
				}

				if imageSpriteTaskSet.WebVttObjectName != nil {
					imageSpriteTaskSetMap["web_vtt_object_name"] = imageSpriteTaskSet.WebVttObjectName
				}

				if imageSpriteTaskSet.ObjectNumberFormat != nil {
					objectNumberFormatMap := map[string]interface{}{}

					if imageSpriteTaskSet.ObjectNumberFormat.InitialValue != nil {
						objectNumberFormatMap["initial_value"] = imageSpriteTaskSet.ObjectNumberFormat.InitialValue
					}

					if imageSpriteTaskSet.ObjectNumberFormat.Increment != nil {
						objectNumberFormatMap["increment"] = imageSpriteTaskSet.ObjectNumberFormat.Increment
					}

					if imageSpriteTaskSet.ObjectNumberFormat.MinLength != nil {
						objectNumberFormatMap["min_length"] = imageSpriteTaskSet.ObjectNumberFormat.MinLength
					}

					if imageSpriteTaskSet.ObjectNumberFormat.PlaceHolder != nil {
						objectNumberFormatMap["place_holder"] = imageSpriteTaskSet.ObjectNumberFormat.PlaceHolder
					}

					imageSpriteTaskSetMap["object_number_format"] = []interface{}{objectNumberFormatMap}
				}

				imageSpriteTaskSetList = append(imageSpriteTaskSetList, imageSpriteTaskSetMap)
			}

			mediaProcessTaskMap["image_sprite_task_set"] = imageSpriteTaskSetList
		}

		if workflow.MediaProcessTask.AdaptiveDynamicStreamingTaskSet != nil {
			adaptiveDynamicStreamingTaskSetList := []interface{}{}
			for _, adaptiveDynamicStreamingTaskSet := range workflow.MediaProcessTask.AdaptiveDynamicStreamingTaskSet {
				adaptiveDynamicStreamingTaskSetMap := map[string]interface{}{}

				if adaptiveDynamicStreamingTaskSet.Definition != nil {
					adaptiveDynamicStreamingTaskSetMap["definition"] = adaptiveDynamicStreamingTaskSet.Definition
				}

				if adaptiveDynamicStreamingTaskSet.WatermarkSet != nil {
					watermarkSetList := []interface{}{}
					for _, watermarkSet := range adaptiveDynamicStreamingTaskSet.WatermarkSet {
						watermarkSetMap := map[string]interface{}{}

						if watermarkSet.Definition != nil {
							watermarkSetMap["definition"] = watermarkSet.Definition
						}

						if watermarkSet.RawParameter != nil {
							rawParameterMap := map[string]interface{}{}

							if watermarkSet.RawParameter.Type != nil {
								rawParameterMap["type"] = watermarkSet.RawParameter.Type
							}

							if watermarkSet.RawParameter.CoordinateOrigin != nil {
								rawParameterMap["coordinate_origin"] = watermarkSet.RawParameter.CoordinateOrigin
							}

							if watermarkSet.RawParameter.XPos != nil {
								rawParameterMap["x_pos"] = watermarkSet.RawParameter.XPos
							}

							if watermarkSet.RawParameter.YPos != nil {
								rawParameterMap["y_pos"] = watermarkSet.RawParameter.YPos
							}

							if watermarkSet.RawParameter.ImageTemplate != nil {
								imageTemplateMap := map[string]interface{}{}

								if watermarkSet.RawParameter.ImageTemplate.ImageContent != nil {
									imageContentMap := map[string]interface{}{}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.Type != nil {
										imageContentMap["type"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.Type
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo != nil {
										cosInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket != nil {
											cosInputInfoMap["bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region != nil {
											cosInputInfoMap["region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object != nil {
											cosInputInfoMap["object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object
										}

										imageContentMap["cos_input_info"] = []interface{}{cosInputInfoMap}
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo != nil {
										urlInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url != nil {
											urlInputInfoMap["url"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url
										}

										imageContentMap["url_input_info"] = []interface{}{urlInputInfoMap}
									}

									imageTemplateMap["image_content"] = []interface{}{imageContentMap}
								}

								if watermarkSet.RawParameter.ImageTemplate.Width != nil {
									imageTemplateMap["width"] = watermarkSet.RawParameter.ImageTemplate.Width
								}

								if watermarkSet.RawParameter.ImageTemplate.Height != nil {
									imageTemplateMap["height"] = watermarkSet.RawParameter.ImageTemplate.Height
								}

								if watermarkSet.RawParameter.ImageTemplate.RepeatType != nil {
									imageTemplateMap["repeat_type"] = watermarkSet.RawParameter.ImageTemplate.RepeatType
								}

								rawParameterMap["image_template"] = []interface{}{imageTemplateMap}
							}

							watermarkSetMap["raw_parameter"] = []interface{}{rawParameterMap}
						}

						if watermarkSet.TextContent != nil {
							watermarkSetMap["text_content"] = watermarkSet.TextContent
						}

						if watermarkSet.SvgContent != nil {
							watermarkSetMap["svg_content"] = watermarkSet.SvgContent
						}

						if watermarkSet.StartTimeOffset != nil {
							watermarkSetMap["start_time_offset"] = watermarkSet.StartTimeOffset
						}

						if watermarkSet.EndTimeOffset != nil {
							watermarkSetMap["end_time_offset"] = watermarkSet.EndTimeOffset
						}

						watermarkSetList = append(watermarkSetList, watermarkSetMap)
					}

					adaptiveDynamicStreamingTaskSetMap["watermark_set"] = watermarkSetList
				}

				if adaptiveDynamicStreamingTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if adaptiveDynamicStreamingTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = adaptiveDynamicStreamingTaskSet.OutputStorage.Type
					}

					if adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					adaptiveDynamicStreamingTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if adaptiveDynamicStreamingTaskSet.OutputObjectPath != nil {
					adaptiveDynamicStreamingTaskSetMap["output_object_path"] = adaptiveDynamicStreamingTaskSet.OutputObjectPath
				}

				if adaptiveDynamicStreamingTaskSet.SubStreamObjectName != nil {
					adaptiveDynamicStreamingTaskSetMap["sub_stream_object_name"] = adaptiveDynamicStreamingTaskSet.SubStreamObjectName
				}

				if adaptiveDynamicStreamingTaskSet.SegmentObjectName != nil {
					adaptiveDynamicStreamingTaskSetMap["segment_object_name"] = adaptiveDynamicStreamingTaskSet.SegmentObjectName
				}

				adaptiveDynamicStreamingTaskSetList = append(adaptiveDynamicStreamingTaskSetList, adaptiveDynamicStreamingTaskSetMap)
			}

			mediaProcessTaskMap["adaptive_dynamic_streaming_task_set"] = adaptiveDynamicStreamingTaskSetList
		}

		_ = d.Set("media_process_task", []interface{}{mediaProcessTaskMap})
	}

	if workflow.AiContentReviewTask != nil {
		aiContentReviewTaskMap := map[string]interface{}{}

		if workflow.AiContentReviewTask.Definition != nil {
			aiContentReviewTaskMap["definition"] = workflow.AiContentReviewTask.Definition
		}

		_ = d.Set("ai_content_review_task", []interface{}{aiContentReviewTaskMap})
	}

	if workflow.AiAnalysisTask != nil {
		aiAnalysisTaskMap := map[string]interface{}{}

		if workflow.AiAnalysisTask.Definition != nil {
			aiAnalysisTaskMap["definition"] = workflow.AiAnalysisTask.Definition
		}

		if workflow.AiAnalysisTask.ExtendedParameter != nil {
			aiAnalysisTaskMap["extended_parameter"] = workflow.AiAnalysisTask.ExtendedParameter
		}

		_ = d.Set("ai_analysis_task", []interface{}{aiAnalysisTaskMap})
	}

	if workflow.AiRecognitionTask != nil {
		aiRecognitionTaskMap := map[string]interface{}{}

		if workflow.AiRecognitionTask.Definition != nil {
			aiRecognitionTaskMap["definition"] = workflow.AiRecognitionTask.Definition
		}

		_ = d.Set("ai_recognition_task", []interface{}{aiRecognitionTaskMap})
	}

	if workflow.TaskNotifyConfig != nil {
		taskNotifyConfigMap := map[string]interface{}{}

		if workflow.TaskNotifyConfig.CmqModel != nil {
			taskNotifyConfigMap["cmq_model"] = workflow.TaskNotifyConfig.CmqModel
		}

		if workflow.TaskNotifyConfig.CmqRegion != nil {
			taskNotifyConfigMap["cmq_region"] = workflow.TaskNotifyConfig.CmqRegion
		}

		if workflow.TaskNotifyConfig.TopicName != nil {
			taskNotifyConfigMap["topic_name"] = workflow.TaskNotifyConfig.TopicName
		}

		if workflow.TaskNotifyConfig.QueueName != nil {
			taskNotifyConfigMap["queue_name"] = workflow.TaskNotifyConfig.QueueName
		}

		if workflow.TaskNotifyConfig.NotifyMode != nil {
			taskNotifyConfigMap["notify_mode"] = workflow.TaskNotifyConfig.NotifyMode
		}

		if workflow.TaskNotifyConfig.NotifyType != nil {
			taskNotifyConfigMap["notify_type"] = workflow.TaskNotifyConfig.NotifyType
		}

		if workflow.TaskNotifyConfig.NotifyUrl != nil {
			taskNotifyConfigMap["notify_url"] = workflow.TaskNotifyConfig.NotifyUrl
		}

		_ = d.Set("task_notify_config", []interface{}{taskNotifyConfigMap})
	}

	if workflow.TaskPriority != nil {
		_ = d.Set("task_priority", workflow.TaskPriority)
	}

	return nil
}

func resourceTencentCloudMpsWorkflowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_workflow.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewResetWorkflowRequest()

	workflowId := d.Id()

	request.WorkflowId = helper.StrToInt64Point(workflowId)

	isChanged := false

	mutableArgs := []string{
		"workflow_name", "trigger", "output_storage",
		"output_dir", "media_process_task", "ai_content_review_task",
		"ai_analysis_task", "ai_recognition_task", "task_notify_config", "task_priority",
	}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			isChanged = true
			break
		}
	}

	if isChanged {
		if v, ok := d.GetOk("workflow_name"); ok {
			request.WorkflowName = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "trigger"); ok {
			workflowTrigger := mps.WorkflowTrigger{}
			if v, ok := dMap["type"]; ok {
				workflowTrigger.Type = helper.String(v.(string))
			}
			if cosFileUploadTriggerMap, ok := helper.InterfaceToMap(dMap, "cos_file_upload_trigger"); ok {
				cosFileUploadTrigger := mps.CosFileUploadTrigger{}
				if v, ok := cosFileUploadTriggerMap["bucket"]; ok {
					cosFileUploadTrigger.Bucket = helper.String(v.(string))
				}
				if v, ok := cosFileUploadTriggerMap["region"]; ok {
					cosFileUploadTrigger.Region = helper.String(v.(string))
				}
				if v, ok := cosFileUploadTriggerMap["dir"]; ok {
					cosFileUploadTrigger.Dir = helper.String(v.(string))
				}
				if v, ok := cosFileUploadTriggerMap["formats"]; ok {
					formatsSet := v.(*schema.Set).List()
					for i := range formatsSet {
						formats := formatsSet[i].(string)
						cosFileUploadTrigger.Formats = append(cosFileUploadTrigger.Formats, &formats)
					}
				}
				workflowTrigger.CosFileUploadTrigger = &cosFileUploadTrigger
			}
			request.Trigger = &workflowTrigger
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "output_storage"); ok {
			taskOutputStorage := mps.TaskOutputStorage{}
			if v, ok := dMap["type"]; ok {
				taskOutputStorage.Type = helper.String(v.(string))
			}
			if cosOutputStorageMap, ok := helper.InterfaceToMap(dMap, "cos_output_storage"); ok {
				cosOutputStorage := mps.CosOutputStorage{}
				if v, ok := cosOutputStorageMap["bucket"]; ok {
					cosOutputStorage.Bucket = helper.String(v.(string))
				}
				if v, ok := cosOutputStorageMap["region"]; ok {
					cosOutputStorage.Region = helper.String(v.(string))
				}
				taskOutputStorage.CosOutputStorage = &cosOutputStorage
			}
			request.OutputStorage = &taskOutputStorage
		}

		if v, ok := d.GetOk("output_dir"); ok {
			request.OutputDir = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "media_process_task"); ok {
			mediaProcessTaskInput := mps.MediaProcessTaskInput{}
			if v, ok := dMap["transcode_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					transcodeTaskSetMap := item.(map[string]interface{})
					transcodeTaskInput := mps.TranscodeTaskInput{}
					if v, ok := transcodeTaskSetMap["definition"]; ok {
						transcodeTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if rawParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "raw_parameter"); ok {
						rawTranscodeParameter := mps.RawTranscodeParameter{}
						if v, ok := rawParameterMap["container"]; ok {
							rawTranscodeParameter.Container = helper.String(v.(string))
						}
						if v, ok := rawParameterMap["remove_video"]; ok {
							rawTranscodeParameter.RemoveVideo = helper.IntInt64(v.(int))
						}
						if v, ok := rawParameterMap["remove_audio"]; ok {
							rawTranscodeParameter.RemoveAudio = helper.IntInt64(v.(int))
						}
						if videoTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "video_template"); ok {
							videoTemplateInfo := mps.VideoTemplateInfo{}
							if v, ok := videoTemplateMap["codec"]; ok {
								videoTemplateInfo.Codec = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["fps"]; ok {
								videoTemplateInfo.Fps = helper.IntInt64(v.(int))
							}
							if v, ok := videoTemplateMap["bitrate"]; ok {
								videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
							}
							if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
								videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["width"]; ok {
								videoTemplateInfo.Width = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["height"]; ok {
								videoTemplateInfo.Height = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["gop"]; ok {
								videoTemplateInfo.Gop = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["fill_type"]; ok {
								videoTemplateInfo.FillType = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["vcrf"]; ok {
								videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
							}
							rawTranscodeParameter.VideoTemplate = &videoTemplateInfo
						}
						if audioTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "audio_template"); ok {
							audioTemplateInfo := mps.AudioTemplateInfo{}
							if v, ok := audioTemplateMap["codec"]; ok {
								audioTemplateInfo.Codec = helper.String(v.(string))
							}
							if v, ok := audioTemplateMap["bitrate"]; ok {
								audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
							}
							if v, ok := audioTemplateMap["sample_rate"]; ok {
								audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
							}
							if v, ok := audioTemplateMap["audio_channel"]; ok {
								audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
							}
							rawTranscodeParameter.AudioTemplate = &audioTemplateInfo
						}
						if tEHDConfigMap, ok := helper.InterfaceToMap(rawParameterMap, "tehd_config"); ok {
							tEHDConfig := mps.TEHDConfig{}
							if v, ok := tEHDConfigMap["type"]; ok {
								tEHDConfig.Type = helper.String(v.(string))
							}
							if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
								tEHDConfig.MaxVideoBitrate = helper.IntInt64(v.(int))
							}
							rawTranscodeParameter.TEHDConfig = &tEHDConfig
						}
						transcodeTaskInput.RawParameter = &rawTranscodeParameter
					}
					if overrideParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "override_parameter"); ok {
						overrideTranscodeParameter := mps.OverrideTranscodeParameter{}
						if v, ok := overrideParameterMap["container"]; ok {
							overrideTranscodeParameter.Container = helper.String(v.(string))
						}
						if v, ok := overrideParameterMap["remove_video"]; ok {
							overrideTranscodeParameter.RemoveVideo = helper.IntUint64(v.(int))
						}
						if v, ok := overrideParameterMap["remove_audio"]; ok {
							overrideTranscodeParameter.RemoveAudio = helper.IntUint64(v.(int))
						}
						if videoTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "video_template"); ok {
							videoTemplateInfoForUpdate := mps.VideoTemplateInfoForUpdate{}
							if v, ok := videoTemplateMap["codec"]; ok {
								videoTemplateInfoForUpdate.Codec = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["fps"]; ok {
								videoTemplateInfoForUpdate.Fps = helper.IntInt64(v.(int))
							}
							if v, ok := videoTemplateMap["bitrate"]; ok {
								videoTemplateInfoForUpdate.Bitrate = helper.IntInt64(v.(int))
							}
							if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
								videoTemplateInfoForUpdate.ResolutionAdaptive = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["width"]; ok {
								videoTemplateInfoForUpdate.Width = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["height"]; ok {
								videoTemplateInfoForUpdate.Height = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["gop"]; ok {
								videoTemplateInfoForUpdate.Gop = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["fill_type"]; ok {
								videoTemplateInfoForUpdate.FillType = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["vcrf"]; ok {
								videoTemplateInfoForUpdate.Vcrf = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["content_adapt_stream"]; ok {
								videoTemplateInfoForUpdate.ContentAdaptStream = helper.IntUint64(v.(int))
							}
							overrideTranscodeParameter.VideoTemplate = &videoTemplateInfoForUpdate
						}
						if audioTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "audio_template"); ok {
							audioTemplateInfoForUpdate := mps.AudioTemplateInfoForUpdate{}
							if v, ok := audioTemplateMap["codec"]; ok {
								audioTemplateInfoForUpdate.Codec = helper.String(v.(string))
							}
							if v, ok := audioTemplateMap["bitrate"]; ok {
								audioTemplateInfoForUpdate.Bitrate = helper.IntInt64(v.(int))
							}
							if v, ok := audioTemplateMap["sample_rate"]; ok {
								audioTemplateInfoForUpdate.SampleRate = helper.IntUint64(v.(int))
							}
							if v, ok := audioTemplateMap["audio_channel"]; ok {
								audioTemplateInfoForUpdate.AudioChannel = helper.IntInt64(v.(int))
							}
							if v, ok := audioTemplateMap["stream_selects"]; ok {
								streamSelectsSet := v.(*schema.Set).List()
								for i := range streamSelectsSet {
									streamSelects := streamSelectsSet[i].(int)
									audioTemplateInfoForUpdate.StreamSelects = append(audioTemplateInfoForUpdate.StreamSelects, helper.IntInt64(streamSelects))
								}
							}
							overrideTranscodeParameter.AudioTemplate = &audioTemplateInfoForUpdate
						}
						if tEHDConfigMap, ok := helper.InterfaceToMap(overrideParameterMap, "tehd_config"); ok {
							tEHDConfigForUpdate := mps.TEHDConfigForUpdate{}
							if v, ok := tEHDConfigMap["type"]; ok {
								tEHDConfigForUpdate.Type = helper.String(v.(string))
							}
							if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
								tEHDConfigForUpdate.MaxVideoBitrate = helper.IntInt64(v.(int))
							}
							overrideTranscodeParameter.TEHDConfig = &tEHDConfigForUpdate
						}
						if subtitleTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "subtitle_template"); ok {
							subtitleTemplate := mps.SubtitleTemplate{}
							if v, ok := subtitleTemplateMap["path"]; ok {
								subtitleTemplate.Path = helper.String(v.(string))
							}
							if v, ok := subtitleTemplateMap["stream_index"]; ok {
								subtitleTemplate.StreamIndex = helper.IntInt64(v.(int))
							}
							if v, ok := subtitleTemplateMap["font_type"]; ok {
								subtitleTemplate.FontType = helper.String(v.(string))
							}
							if v, ok := subtitleTemplateMap["font_size"]; ok {
								subtitleTemplate.FontSize = helper.String(v.(string))
							}
							if v, ok := subtitleTemplateMap["font_color"]; ok {
								subtitleTemplate.FontColor = helper.String(v.(string))
							}
							if v, ok := subtitleTemplateMap["font_alpha"]; ok {
								subtitleTemplate.FontAlpha = helper.Float64(v.(float64))
							}
							overrideTranscodeParameter.SubtitleTemplate = &subtitleTemplate
						}
						transcodeTaskInput.OverrideParameter = &overrideTranscodeParameter
					}
					if v, ok := transcodeTaskSetMap["watermark_set"]; ok {
						for _, item := range v.([]interface{}) {
							watermarkSetMap := item.(map[string]interface{})
							watermarkInput := mps.WatermarkInput{}
							if v, ok := watermarkSetMap["definition"]; ok {
								watermarkInput.Definition = helper.IntUint64(v.(int))
							}
							if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
								rawWatermarkParameter := mps.RawWatermarkParameter{}
								if v, ok := rawParameterMap["type"]; ok {
									rawWatermarkParameter.Type = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["coordinate_origin"]; ok {
									rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["x_pos"]; ok {
									rawWatermarkParameter.XPos = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["y_pos"]; ok {
									rawWatermarkParameter.YPos = helper.String(v.(string))
								}
								if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
									rawImageWatermarkInput := mps.RawImageWatermarkInput{}
									if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := imageContentMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
											cosInputInfo := mps.CosInputInfo{}
											if v, ok := cosInputInfoMap["bucket"]; ok {
												cosInputInfo.Bucket = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["region"]; ok {
												cosInputInfo.Region = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["object"]; ok {
												cosInputInfo.Object = helper.String(v.(string))
											}
											mediaInputInfo.CosInputInfo = &cosInputInfo
										}
										if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										rawImageWatermarkInput.ImageContent = &mediaInputInfo
									}
									if v, ok := imageTemplateMap["width"]; ok {
										rawImageWatermarkInput.Width = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["height"]; ok {
										rawImageWatermarkInput.Height = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["repeat_type"]; ok {
										rawImageWatermarkInput.RepeatType = helper.String(v.(string))
									}
									rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
								}
								watermarkInput.RawParameter = &rawWatermarkParameter
							}
							if v, ok := watermarkSetMap["text_content"]; ok {
								watermarkInput.TextContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["svg_content"]; ok {
								watermarkInput.SvgContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["start_time_offset"]; ok {
								watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := watermarkSetMap["end_time_offset"]; ok {
								watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							transcodeTaskInput.WatermarkSet = append(transcodeTaskInput.WatermarkSet, &watermarkInput)
						}
					}
					if v, ok := transcodeTaskSetMap["mosaic_set"]; ok {
						for _, item := range v.([]interface{}) {
							mosaicSetMap := item.(map[string]interface{})
							mosaicInput := mps.MosaicInput{}
							if v, ok := mosaicSetMap["coordinate_origin"]; ok {
								mosaicInput.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["x_pos"]; ok {
								mosaicInput.XPos = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["y_pos"]; ok {
								mosaicInput.YPos = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["width"]; ok {
								mosaicInput.Width = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["height"]; ok {
								mosaicInput.Height = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["start_time_offset"]; ok {
								mosaicInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := mosaicSetMap["end_time_offset"]; ok {
								mosaicInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							transcodeTaskInput.MosaicSet = append(transcodeTaskInput.MosaicSet, &mosaicInput)
						}
					}
					if v, ok := transcodeTaskSetMap["start_time_offset"]; ok {
						transcodeTaskInput.StartTimeOffset = helper.Float64(v.(float64))
					}
					if v, ok := transcodeTaskSetMap["end_time_offset"]; ok {
						transcodeTaskInput.EndTimeOffset = helper.Float64(v.(float64))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						transcodeTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := transcodeTaskSetMap["output_object_path"]; ok {
						transcodeTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := transcodeTaskSetMap["segment_object_name"]; ok {
						transcodeTaskInput.SegmentObjectName = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "object_number_format"); ok {
						numberFormat := mps.NumberFormat{}
						if v, ok := objectNumberFormatMap["initial_value"]; ok {
							numberFormat.InitialValue = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["increment"]; ok {
							numberFormat.Increment = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["min_length"]; ok {
							numberFormat.MinLength = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["place_holder"]; ok {
							numberFormat.PlaceHolder = helper.String(v.(string))
						}
						transcodeTaskInput.ObjectNumberFormat = &numberFormat
					}
					if headTailParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "head_tail_parameter"); ok {
						headTailParameter := mps.HeadTailParameter{}
						if v, ok := headTailParameterMap["head_set"]; ok {
							for _, item := range v.([]interface{}) {
								headSetMap := item.(map[string]interface{})
								mediaInputInfo := mps.MediaInputInfo{}
								if v, ok := headSetMap["type"]; ok {
									mediaInputInfo.Type = helper.String(v.(string))
								}
								if cosInputInfoMap, ok := helper.InterfaceToMap(headSetMap, "cos_input_info"); ok {
									cosInputInfo := mps.CosInputInfo{}
									if v, ok := cosInputInfoMap["bucket"]; ok {
										cosInputInfo.Bucket = helper.String(v.(string))
									}
									if v, ok := cosInputInfoMap["region"]; ok {
										cosInputInfo.Region = helper.String(v.(string))
									}
									if v, ok := cosInputInfoMap["object"]; ok {
										cosInputInfo.Object = helper.String(v.(string))
									}
									mediaInputInfo.CosInputInfo = &cosInputInfo
								}
								if urlInputInfoMap, ok := helper.InterfaceToMap(headSetMap, "url_input_info"); ok {
									urlInputInfo := mps.UrlInputInfo{}
									if v, ok := urlInputInfoMap["url"]; ok {
										urlInputInfo.Url = helper.String(v.(string))
									}
									mediaInputInfo.UrlInputInfo = &urlInputInfo
								}
								headTailParameter.HeadSet = append(headTailParameter.HeadSet, &mediaInputInfo)
							}
						}
						if v, ok := headTailParameterMap["tail_set"]; ok {
							for _, item := range v.([]interface{}) {
								tailSetMap := item.(map[string]interface{})
								mediaInputInfo := mps.MediaInputInfo{}
								if v, ok := tailSetMap["type"]; ok {
									mediaInputInfo.Type = helper.String(v.(string))
								}
								if cosInputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "cos_input_info"); ok {
									cosInputInfo := mps.CosInputInfo{}
									if v, ok := cosInputInfoMap["bucket"]; ok {
										cosInputInfo.Bucket = helper.String(v.(string))
									}
									if v, ok := cosInputInfoMap["region"]; ok {
										cosInputInfo.Region = helper.String(v.(string))
									}
									if v, ok := cosInputInfoMap["object"]; ok {
										cosInputInfo.Object = helper.String(v.(string))
									}
									mediaInputInfo.CosInputInfo = &cosInputInfo
								}
								if urlInputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "url_input_info"); ok {
									urlInputInfo := mps.UrlInputInfo{}
									if v, ok := urlInputInfoMap["url"]; ok {
										urlInputInfo.Url = helper.String(v.(string))
									}
									mediaInputInfo.UrlInputInfo = &urlInputInfo
								}
								headTailParameter.TailSet = append(headTailParameter.TailSet, &mediaInputInfo)
							}
						}
						transcodeTaskInput.HeadTailParameter = &headTailParameter
					}
					mediaProcessTaskInput.TranscodeTaskSet = append(mediaProcessTaskInput.TranscodeTaskSet, &transcodeTaskInput)
				}
			}
			if v, ok := dMap["animated_graphic_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					animatedGraphicTaskSetMap := item.(map[string]interface{})
					animatedGraphicTaskInput := mps.AnimatedGraphicTaskInput{}
					if v, ok := animatedGraphicTaskSetMap["definition"]; ok {
						animatedGraphicTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := animatedGraphicTaskSetMap["start_time_offset"]; ok {
						animatedGraphicTaskInput.StartTimeOffset = helper.Float64(v.(float64))
					}
					if v, ok := animatedGraphicTaskSetMap["end_time_offset"]; ok {
						animatedGraphicTaskInput.EndTimeOffset = helper.Float64(v.(float64))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(animatedGraphicTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						animatedGraphicTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := animatedGraphicTaskSetMap["output_object_path"]; ok {
						animatedGraphicTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					mediaProcessTaskInput.AnimatedGraphicTaskSet = append(mediaProcessTaskInput.AnimatedGraphicTaskSet, &animatedGraphicTaskInput)
				}
			}
			if v, ok := dMap["snapshot_by_time_offset_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					snapshotByTimeOffsetTaskSetMap := item.(map[string]interface{})
					snapshotByTimeOffsetTaskInput := mps.SnapshotByTimeOffsetTaskInput{}
					if v, ok := snapshotByTimeOffsetTaskSetMap["definition"]; ok {
						snapshotByTimeOffsetTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := snapshotByTimeOffsetTaskSetMap["ext_time_offset_set"]; ok {
						extTimeOffsetSetSet := v.(*schema.Set).List()
						for i := range extTimeOffsetSetSet {
							extTimeOffsetSet := extTimeOffsetSetSet[i].(string)
							snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet = append(snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet, &extTimeOffsetSet)
						}
					}
					if v, _ := d.GetOk("time_offset_set"); v != nil {
						timeOffsetSetSet := v.(*schema.Set).List()
						for i := range timeOffsetSetSet {
							timeOffsetSet := timeOffsetSetSet[i].(float64)
							snapshotByTimeOffsetTaskInput.TimeOffsetSet = append(snapshotByTimeOffsetTaskInput.TimeOffsetSet, &timeOffsetSet)
						}
					}

					if v, ok := snapshotByTimeOffsetTaskSetMap["watermark_set"]; ok {
						for _, item := range v.([]interface{}) {
							watermarkSetMap := item.(map[string]interface{})
							watermarkInput := mps.WatermarkInput{}
							if v, ok := watermarkSetMap["definition"]; ok {
								watermarkInput.Definition = helper.IntUint64(v.(int))
							}
							if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
								rawWatermarkParameter := mps.RawWatermarkParameter{}
								if v, ok := rawParameterMap["type"]; ok {
									rawWatermarkParameter.Type = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["coordinate_origin"]; ok {
									rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["x_pos"]; ok {
									rawWatermarkParameter.XPos = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["y_pos"]; ok {
									rawWatermarkParameter.YPos = helper.String(v.(string))
								}
								if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
									rawImageWatermarkInput := mps.RawImageWatermarkInput{}
									if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := imageContentMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
											cosInputInfo := mps.CosInputInfo{}
											if v, ok := cosInputInfoMap["bucket"]; ok {
												cosInputInfo.Bucket = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["region"]; ok {
												cosInputInfo.Region = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["object"]; ok {
												cosInputInfo.Object = helper.String(v.(string))
											}
											mediaInputInfo.CosInputInfo = &cosInputInfo
										}
										if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										rawImageWatermarkInput.ImageContent = &mediaInputInfo
									}
									if v, ok := imageTemplateMap["width"]; ok {
										rawImageWatermarkInput.Width = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["height"]; ok {
										rawImageWatermarkInput.Height = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["repeat_type"]; ok {
										rawImageWatermarkInput.RepeatType = helper.String(v.(string))
									}
									rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
								}
								watermarkInput.RawParameter = &rawWatermarkParameter
							}
							if v, ok := watermarkSetMap["text_content"]; ok {
								watermarkInput.TextContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["svg_content"]; ok {
								watermarkInput.SvgContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["start_time_offset"]; ok {
								watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := watermarkSetMap["end_time_offset"]; ok {
								watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							snapshotByTimeOffsetTaskInput.WatermarkSet = append(snapshotByTimeOffsetTaskInput.WatermarkSet, &watermarkInput)
						}
					}
					if outputStorageMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						snapshotByTimeOffsetTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := snapshotByTimeOffsetTaskSetMap["output_object_path"]; ok {
						snapshotByTimeOffsetTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskSetMap, "object_number_format"); ok {
						numberFormat := mps.NumberFormat{}
						if v, ok := objectNumberFormatMap["initial_value"]; ok {
							numberFormat.InitialValue = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["increment"]; ok {
							numberFormat.Increment = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["min_length"]; ok {
							numberFormat.MinLength = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["place_holder"]; ok {
							numberFormat.PlaceHolder = helper.String(v.(string))
						}
						snapshotByTimeOffsetTaskInput.ObjectNumberFormat = &numberFormat
					}
					mediaProcessTaskInput.SnapshotByTimeOffsetTaskSet = append(mediaProcessTaskInput.SnapshotByTimeOffsetTaskSet, &snapshotByTimeOffsetTaskInput)
				}
			}
			if v, ok := dMap["sample_snapshot_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					sampleSnapshotTaskSetMap := item.(map[string]interface{})
					sampleSnapshotTaskInput := mps.SampleSnapshotTaskInput{}
					if v, ok := sampleSnapshotTaskSetMap["definition"]; ok {
						sampleSnapshotTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := sampleSnapshotTaskSetMap["watermark_set"]; ok {
						for _, item := range v.([]interface{}) {
							watermarkSetMap := item.(map[string]interface{})
							watermarkInput := mps.WatermarkInput{}
							if v, ok := watermarkSetMap["definition"]; ok {
								watermarkInput.Definition = helper.IntUint64(v.(int))
							}
							if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
								rawWatermarkParameter := mps.RawWatermarkParameter{}
								if v, ok := rawParameterMap["type"]; ok {
									rawWatermarkParameter.Type = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["coordinate_origin"]; ok {
									rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["x_pos"]; ok {
									rawWatermarkParameter.XPos = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["y_pos"]; ok {
									rawWatermarkParameter.YPos = helper.String(v.(string))
								}
								if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
									rawImageWatermarkInput := mps.RawImageWatermarkInput{}
									if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := imageContentMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
											cosInputInfo := mps.CosInputInfo{}
											if v, ok := cosInputInfoMap["bucket"]; ok {
												cosInputInfo.Bucket = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["region"]; ok {
												cosInputInfo.Region = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["object"]; ok {
												cosInputInfo.Object = helper.String(v.(string))
											}
											mediaInputInfo.CosInputInfo = &cosInputInfo
										}
										if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										rawImageWatermarkInput.ImageContent = &mediaInputInfo
									}
									if v, ok := imageTemplateMap["width"]; ok {
										rawImageWatermarkInput.Width = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["height"]; ok {
										rawImageWatermarkInput.Height = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["repeat_type"]; ok {
										rawImageWatermarkInput.RepeatType = helper.String(v.(string))
									}
									rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
								}
								watermarkInput.RawParameter = &rawWatermarkParameter
							}
							if v, ok := watermarkSetMap["text_content"]; ok {
								watermarkInput.TextContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["svg_content"]; ok {
								watermarkInput.SvgContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["start_time_offset"]; ok {
								watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := watermarkSetMap["end_time_offset"]; ok {
								watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							sampleSnapshotTaskInput.WatermarkSet = append(sampleSnapshotTaskInput.WatermarkSet, &watermarkInput)
						}
					}
					if outputStorageMap, ok := helper.InterfaceToMap(sampleSnapshotTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						sampleSnapshotTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := sampleSnapshotTaskSetMap["output_object_path"]; ok {
						sampleSnapshotTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(sampleSnapshotTaskSetMap, "object_number_format"); ok {
						numberFormat := mps.NumberFormat{}
						if v, ok := objectNumberFormatMap["initial_value"]; ok {
							numberFormat.InitialValue = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["increment"]; ok {
							numberFormat.Increment = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["min_length"]; ok {
							numberFormat.MinLength = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["place_holder"]; ok {
							numberFormat.PlaceHolder = helper.String(v.(string))
						}
						sampleSnapshotTaskInput.ObjectNumberFormat = &numberFormat
					}
					mediaProcessTaskInput.SampleSnapshotTaskSet = append(mediaProcessTaskInput.SampleSnapshotTaskSet, &sampleSnapshotTaskInput)
				}
			}
			if v, ok := dMap["image_sprite_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					imageSpriteTaskSetMap := item.(map[string]interface{})
					imageSpriteTaskInput := mps.ImageSpriteTaskInput{}
					if v, ok := imageSpriteTaskSetMap["definition"]; ok {
						imageSpriteTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(imageSpriteTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						imageSpriteTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := imageSpriteTaskSetMap["output_object_path"]; ok {
						imageSpriteTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := imageSpriteTaskSetMap["web_vtt_object_name"]; ok {
						imageSpriteTaskInput.WebVttObjectName = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(imageSpriteTaskSetMap, "object_number_format"); ok {
						numberFormat := mps.NumberFormat{}
						if v, ok := objectNumberFormatMap["initial_value"]; ok {
							numberFormat.InitialValue = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["increment"]; ok {
							numberFormat.Increment = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["min_length"]; ok {
							numberFormat.MinLength = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["place_holder"]; ok {
							numberFormat.PlaceHolder = helper.String(v.(string))
						}
						imageSpriteTaskInput.ObjectNumberFormat = &numberFormat
					}
					mediaProcessTaskInput.ImageSpriteTaskSet = append(mediaProcessTaskInput.ImageSpriteTaskSet, &imageSpriteTaskInput)
				}
			}
			if v, ok := dMap["adaptive_dynamic_streaming_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					adaptiveDynamicStreamingTaskSetMap := item.(map[string]interface{})
					adaptiveDynamicStreamingTaskInput := mps.AdaptiveDynamicStreamingTaskInput{}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["definition"]; ok {
						adaptiveDynamicStreamingTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["watermark_set"]; ok {
						for _, item := range v.([]interface{}) {
							watermarkSetMap := item.(map[string]interface{})
							watermarkInput := mps.WatermarkInput{}
							if v, ok := watermarkSetMap["definition"]; ok {
								watermarkInput.Definition = helper.IntUint64(v.(int))
							}
							if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
								rawWatermarkParameter := mps.RawWatermarkParameter{}
								if v, ok := rawParameterMap["type"]; ok {
									rawWatermarkParameter.Type = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["coordinate_origin"]; ok {
									rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["x_pos"]; ok {
									rawWatermarkParameter.XPos = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["y_pos"]; ok {
									rawWatermarkParameter.YPos = helper.String(v.(string))
								}
								if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
									rawImageWatermarkInput := mps.RawImageWatermarkInput{}
									if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := imageContentMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
											cosInputInfo := mps.CosInputInfo{}
											if v, ok := cosInputInfoMap["bucket"]; ok {
												cosInputInfo.Bucket = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["region"]; ok {
												cosInputInfo.Region = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["object"]; ok {
												cosInputInfo.Object = helper.String(v.(string))
											}
											mediaInputInfo.CosInputInfo = &cosInputInfo
										}
										if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										rawImageWatermarkInput.ImageContent = &mediaInputInfo
									}
									if v, ok := imageTemplateMap["width"]; ok {
										rawImageWatermarkInput.Width = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["height"]; ok {
										rawImageWatermarkInput.Height = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["repeat_type"]; ok {
										rawImageWatermarkInput.RepeatType = helper.String(v.(string))
									}
									rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
								}
								watermarkInput.RawParameter = &rawWatermarkParameter
							}
							if v, ok := watermarkSetMap["text_content"]; ok {
								watermarkInput.TextContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["svg_content"]; ok {
								watermarkInput.SvgContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["start_time_offset"]; ok {
								watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := watermarkSetMap["end_time_offset"]; ok {
								watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							adaptiveDynamicStreamingTaskInput.WatermarkSet = append(adaptiveDynamicStreamingTaskInput.WatermarkSet, &watermarkInput)
						}
					}
					if outputStorageMap, ok := helper.InterfaceToMap(adaptiveDynamicStreamingTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						adaptiveDynamicStreamingTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["output_object_path"]; ok {
						adaptiveDynamicStreamingTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["sub_stream_object_name"]; ok {
						adaptiveDynamicStreamingTaskInput.SubStreamObjectName = helper.String(v.(string))
					}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["segment_object_name"]; ok {
						adaptiveDynamicStreamingTaskInput.SegmentObjectName = helper.String(v.(string))
					}
					mediaProcessTaskInput.AdaptiveDynamicStreamingTaskSet = append(mediaProcessTaskInput.AdaptiveDynamicStreamingTaskSet, &adaptiveDynamicStreamingTaskInput)
				}
			}
			request.MediaProcessTask = &mediaProcessTaskInput
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ai_content_review_task"); ok {
			aiContentReviewTaskInput := mps.AiContentReviewTaskInput{}
			if v, ok := dMap["definition"]; ok {
				aiContentReviewTaskInput.Definition = helper.IntUint64(v.(int))
			}
			request.AiContentReviewTask = &aiContentReviewTaskInput
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ai_analysis_task"); ok {
			aiAnalysisTaskInput := mps.AiAnalysisTaskInput{}
			if v, ok := dMap["definition"]; ok {
				aiAnalysisTaskInput.Definition = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["extended_parameter"]; ok {
				aiAnalysisTaskInput.ExtendedParameter = helper.String(v.(string))
			}
			request.AiAnalysisTask = &aiAnalysisTaskInput
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ai_recognition_task"); ok {
			aiRecognitionTaskInput := mps.AiRecognitionTaskInput{}
			if v, ok := dMap["definition"]; ok {
				aiRecognitionTaskInput.Definition = helper.IntUint64(v.(int))
			}
			request.AiRecognitionTask = &aiRecognitionTaskInput
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "task_notify_config"); ok {
			taskNotifyConfig := mps.TaskNotifyConfig{}
			if v, ok := dMap["cmq_model"]; ok {
				taskNotifyConfig.CmqModel = helper.String(v.(string))
			}
			if v, ok := dMap["cmq_region"]; ok {
				taskNotifyConfig.CmqRegion = helper.String(v.(string))
			}
			if v, ok := dMap["topic_name"]; ok {
				taskNotifyConfig.TopicName = helper.String(v.(string))
			}
			if v, ok := dMap["queue_name"]; ok {
				taskNotifyConfig.QueueName = helper.String(v.(string))
			}
			if v, ok := dMap["notify_mode"]; ok {
				taskNotifyConfig.NotifyMode = helper.String(v.(string))
			}
			if v, ok := dMap["notify_type"]; ok {
				taskNotifyConfig.NotifyType = helper.String(v.(string))
			}
			if v, ok := dMap["notify_url"]; ok {
				taskNotifyConfig.NotifyUrl = helper.String(v.(string))
			}
			request.TaskNotifyConfig = &taskNotifyConfig
		}

		if v, _ := d.GetOk("task_priority"); v != nil {
			request.TaskPriority = helper.IntInt64(v.(int))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ResetWorkflow(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mps workflow failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsWorkflowRead(d, meta)
}

func resourceTencentCloudMpsWorkflowDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_workflow.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	workflowId := d.Id()

	if err := service.DeleteMpsWorkflowById(ctx, workflowId); err != nil {
		return err
	}

	return nil
}
