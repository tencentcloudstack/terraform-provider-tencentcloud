/*
Use this data source to query detailed information of mps workflows

Example Usage

```hcl
data "tencentcloud_mps_workflows" "workflows" {
  workflow_ids = &lt;nil&gt;
  status = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  total_count = &lt;nil&gt;
  workflow_info_set {
		workflow_id = &lt;nil&gt;
		workflow_name = &lt;nil&gt;
		status = &lt;nil&gt;
		trigger {
			type = "CosFileUpload"
			cos_file_upload_trigger {
				bucket = "TopRankVideo-125xxx88"
				region = "ap-chongqing"
				dir = "/movie/201907/"
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
		media_process_task {
			transcode_task_set {
				definition = &lt;nil&gt;
				raw_parameter {
					container = &lt;nil&gt;
					remove_video = &lt;nil&gt;
					remove_audio = 0
					video_template {
						codec = &lt;nil&gt;
						fps = &lt;nil&gt;
						bitrate = &lt;nil&gt;
						resolution_adaptive = "open"
						width = 0
						height = 0
						gop = &lt;nil&gt;
						fill_type = "black"
						vcrf = &lt;nil&gt;
					}
					audio_template {
						codec = &lt;nil&gt;
						bitrate = &lt;nil&gt;
						sample_rate = &lt;nil&gt;
						audio_channel = 2
					}
					t_e_h_d_config {
						type = &lt;nil&gt;
						max_video_bitrate = &lt;nil&gt;
					}
				}
				override_parameter {
					container = &lt;nil&gt;
					remove_video = &lt;nil&gt;
					remove_audio = &lt;nil&gt;
					video_template {
						codec = &lt;nil&gt;
						fps = &lt;nil&gt;
						bitrate = &lt;nil&gt;
						resolution_adaptive = &lt;nil&gt;
						width = &lt;nil&gt;
						height = &lt;nil&gt;
						gop = &lt;nil&gt;
						fill_type = &lt;nil&gt;
						vcrf = &lt;nil&gt;
						content_adapt_stream = 0
					}
					audio_template {
						codec = &lt;nil&gt;
						bitrate = &lt;nil&gt;
						sample_rate = &lt;nil&gt;
						audio_channel = &lt;nil&gt;
						stream_selects = &lt;nil&gt;
					}
					t_e_h_d_config {
						type = &lt;nil&gt;
						max_video_bitrate = &lt;nil&gt;
					}
					subtitle_template {
						path = &lt;nil&gt;
						stream_index = &lt;nil&gt;
						font_type = "hei.ttf"
						font_size = &lt;nil&gt;
						font_color = "0xFFFFFF"
						font_alpha =
					}
				}
				watermark_set {
					definition = &lt;nil&gt;
					raw_parameter {
						type = &lt;nil&gt;
						coordinate_origin = "TopLeft"
						x_pos = "0px"
						y_pos = "0px"
						image_template {
							image_content {
								type = "COS"
								cos_input_info {
									bucket = "TopRankVideo-125xxx88"
									region = "ap-chongqing"
									object = "/movie/201907/WildAnimal.mov"
								}
								url_input_info {
									url = &lt;nil&gt;
								}
							}
							width = "10%"
							height = "0px"
							repeat_type = &lt;nil&gt;
						}
					}
					text_content = &lt;nil&gt;
					svg_content = &lt;nil&gt;
					start_time_offset = &lt;nil&gt;
					end_time_offset = &lt;nil&gt;
				}
				mosaic_set {
					coordinate_origin = "TopLeft"
					x_pos = "0px"
					y_pos = "0px"
					width = "10%"
					height = "10%"
					start_time_offset = &lt;nil&gt;
					end_time_offset = &lt;nil&gt;
				}
				start_time_offset = &lt;nil&gt;
				end_time_offset = &lt;nil&gt;
				output_storage {
					type = "COS"
					cos_output_storage {
						bucket = "TopRankVideo-125xxx88"
						region = "ap-chongqinq"
					}
				}
				output_object_path = &lt;nil&gt;
				segment_object_name = &lt;nil&gt;
				object_number_format {
					initial_value = 0
					increment = 1
					min_length = 1
					place_holder = "0"
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
							url = &lt;nil&gt;
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
							url = &lt;nil&gt;
						}
					}
				}
			}
			animated_graphic_task_set {
				definition = &lt;nil&gt;
				start_time_offset = &lt;nil&gt;
				end_time_offset = &lt;nil&gt;
				output_storage {
					type = "COS"
					cos_output_storage {
						bucket = "TopRankVideo-125xxx88"
						region = "ap-chongqinq"
					}
				}
				output_object_path = &lt;nil&gt;
			}
			snapshot_by_time_offset_task_set {
				definition = &lt;nil&gt;
				ext_time_offset_set = &lt;nil&gt;
				time_offset_set = &lt;nil&gt;
				watermark_set {
					definition = &lt;nil&gt;
					raw_parameter {
						type = &lt;nil&gt;
						coordinate_origin = &lt;nil&gt;
						x_pos = "0px"
						y_pos = "0px"
						image_template {
							image_content {
								type = "COS"
								cos_input_info {
									bucket = "TopRankVideo-125xxx88"
									region = "ap-chongqing"
									object = "/movie/201907/WildAnimal.mov"
								}
								url_input_info {
									url = &lt;nil&gt;
								}
							}
							width = "10%"
							height = "0px"
							repeat_type = &lt;nil&gt;
						}
					}
					text_content = &lt;nil&gt;
					svg_content = &lt;nil&gt;
					start_time_offset = &lt;nil&gt;
					end_time_offset = &lt;nil&gt;
				}
				output_storage {
					type = "COS"
					cos_output_storage {
						bucket = "TopRankVideo-125xxx88"
						region = "ap-chongqinq"
					}
				}
				output_object_path = &lt;nil&gt;
				object_number_format {
					initial_value = 0
					increment = 1
					min_length = 1
					place_holder = "0"
				}
			}
			sample_snapshot_task_set {
				definition = &lt;nil&gt;
				watermark_set {
					definition = &lt;nil&gt;
					raw_parameter {
						type = &lt;nil&gt;
						coordinate_origin = "TopLeft"
						x_pos = "0px"
						y_pos = "0px"
						image_template {
							image_content {
								type = "COS"
								cos_input_info {
									bucket = "TopRankVideo-125xxx88"
									region = "ap-chongqing"
									object = "/movie/201907/WildAnimal.mov"
								}
								url_input_info {
									url = &lt;nil&gt;
								}
							}
							width = "10%"
							height = "0px"
							repeat_type = "repeat"
						}
					}
					text_content = &lt;nil&gt;
					svg_content = &lt;nil&gt;
					start_time_offset = &lt;nil&gt;
					end_time_offset = &lt;nil&gt;
				}
				output_storage {
					type = "COS"
					cos_output_storage {
						bucket = "TopRankVideo-125xxx88"
						region = "ap-chongqinq"
					}
				}
				output_object_path = &lt;nil&gt;
				object_number_format {
					initial_value = 0
					increment = 1
					min_length = 1
					place_holder = "0"
				}
			}
			image_sprite_task_set {
				definition = &lt;nil&gt;
				output_storage {
					type = "COS"
					cos_output_storage {
						bucket = "TopRankVideo-125xxx88"
						region = "ap-chongqinq"
					}
				}
				output_object_path = &lt;nil&gt;
				web_vtt_object_name = &lt;nil&gt;
				object_number_format {
					initial_value = 0
					increment = 1
					min_length = 1
					place_holder = "0"
				}
			}
			adaptive_dynamic_streaming_task_set {
				definition = &lt;nil&gt;
				watermark_set {
					definition = &lt;nil&gt;
					raw_parameter {
						type = &lt;nil&gt;
						coordinate_origin = "TopLeft"
						x_pos = "0px"
						y_pos = "0px"
						image_template {
							image_content {
								type = "COS"
								cos_input_info {
									bucket = "TopRankVideo-125xxx88"
									region = "ap-chongqing"
									object = "/movie/201907/WildAnimal.mov"
								}
								url_input_info {
									url = &lt;nil&gt;
								}
							}
							width = "10%"
							height = "0px"
							repeat_type = "repeat"
						}
					}
					text_content = &lt;nil&gt;
					svg_content = &lt;nil&gt;
					start_time_offset = &lt;nil&gt;
					end_time_offset = &lt;nil&gt;
				}
				output_storage {
					type = "COS"
					cos_output_storage {
						bucket = "TopRankVideo-125xxx88"
						region = "ap-chongqinq"
					}
				}
				output_object_path = &lt;nil&gt;
				sub_stream_object_name = &lt;nil&gt;
				segment_object_name = &lt;nil&gt;
			}
		}
		ai_content_review_task {
			definition = &lt;nil&gt;
		}
		ai_analysis_task {
			definition = &lt;nil&gt;
			extended_parameter = &lt;nil&gt;
		}
		ai_recognition_task {
			definition = &lt;nil&gt;
		}
		task_notify_config {
			cmq_model = &lt;nil&gt;
			cmq_region = &lt;nil&gt;
			topic_name = &lt;nil&gt;
			queue_name = &lt;nil&gt;
			notify_mode = &lt;nil&gt;
			notify_type = &lt;nil&gt;
			notify_url = &lt;nil&gt;
		}
		task_priority = &lt;nil&gt;
		output_dir = "/movie/201907/"
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;

  }
  request_id = &lt;nil&gt;
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMpsWorkflows() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsWorkflowsRead,
		Schema: map[string]*schema.Schema{
			"workflow_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Workflow ID filter condition, array length limit: 100.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workflow status, value range:Enabled, Disabled.If this parameter is not filled in, the workflow status will not be distinguished.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset, default: 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return the number of records, default value: 10, maximum value: 100.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"workflow_info_set": {
				Type:        schema.TypeList,
				Description: "Workflow information array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workflow_id": {
							Type:        schema.TypeInt,
							Description: "Workflow ID.",
						},
						"workflow_name": {
							Type:        schema.TypeString,
							Description: "Workflow name.",
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Workflow status, value range:Enabled, Disabled.",
						},
						"trigger": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "The input rule bound to the workflow, when the uploaded video hits the rule to this object, the workflow will be triggered.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "The type of trigger, currently only supports CosFileUpload.",
									},
									"cos_file_upload_trigger": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Mandatory and valid when Type is CosFileUpload, the rule is triggered for COS.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Description: "The name of the COS Bucket bound to the workflow.",
												},
												"region": {
													Type:        schema.TypeString,
													Description: "The park to which the COS Bucket bound to the workflow belongs.",
												},
												"dir": {
													Type:        schema.TypeString,
													Description: "The input path directory of the workflow binding must be an absolute path, that is, start and end with `/`.",
												},
												"formats": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "A list of file formats that are allowed to be triggered by the workflow, if not filled in, it means that files of all formats can trigger the workflow.",
												},
											},
										},
									},
								},
							},
						},
						"output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "File output storage location for media processing.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Description: "The type of media processing output object storage location, now only supports COS.",
									},
									"cos_output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
												},
												"region": {
													Type:        schema.TypeString,
													Description: "The park of the target Bucket for the output of the file generated by media processing. If not filled, it means inheriting from the upper layer.",
												},
											},
										},
									},
								},
							},
						},
						"media_process_task": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Media Processing Type Task Parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"transcode_task_set": {
										Type:        schema.TypeList,
										Description: "Video transcoding task list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Description: "Video Transcoding Template ID.",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Video transcoding custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios. It is recommended that you use Definition to specify transcoding parameters first.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"container": {
																Type:        schema.TypeString,
																Description: "Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.",
															},
															"remove_video": {
																Type:        schema.TypeInt,
																Description: "Whether to remove video data, value:0: reserved.1: remove.Default: 0.",
															},
															"remove_audio": {
																Type:        schema.TypeInt,
																Description: "Whether to remove audio data, value:0: reserved.1: remove.Default: 0.",
															},
															"video_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Video stream configuration parameters, when RemoveVideo is 0, this field is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"codec": {
																			Type:        schema.TypeString,
																			Description: "Encoding format of the video stream, optional value:libx264: H.264 encoding.libx265: H.265 encoding.av1: AOMedia Video 1 encoding.Note: Currently H.265 encoding must specify a resolution, and it needs to be within 640*480.Note: av1 encoded containers currently only support mp4.",
																		},
																		"fps": {
																			Type:        schema.TypeInt,
																			Description: "Video frame rate, value range: [0, 100], unit: Hz.When the value is 0, it means that the frame rate is consistent with the original video.Note: The value range for adaptive code rate is [0, 60].",
																		},
																		"bitrate": {
																			Type:        schema.TypeInt,
																			Description: "Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.When the value is 0, it means that the video bit rate is consistent with the original video.",
																		},
																		"resolution_adaptive": {
																			Type:        schema.TypeString,
																			Description: "Adaptive resolution, optional values:```open: open, at this time, Width represents the long side of the video, Height represents the short side of the video.close: close, at this time, Width represents the width of the video, and Height represents the height of the video.Default: open.Note: In adaptive mode, Width cannot be smaller than Height.",
																		},
																		"width": {
																			Type:        schema.TypeInt,
																			Description: "The maximum value of video stream width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default: 0.",
																		},
																		"height": {
																			Type:        schema.TypeInt,
																			Description: "The maximum value of video stream height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default: 0.",
																		},
																		"gop": {
																			Type:        schema.TypeInt,
																			Description: "The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.When filling 0 or not filling, the system will automatically set the gop length.",
																		},
																		"fill_type": {
																			Type:        schema.TypeString,
																			Description: "Filling method, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling method:stretch: Stretch, stretch each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the aspect ratio of the video unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video unchanged, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and fill the rest of the edge with Gaussian blur.Default: black.Note: Adaptive stream only supports stretch, black.",
																		},
																		"vcrf": {
																			Type:        schema.TypeInt,
																			Description: "Video constant bit rate control factor, the value range is [1, 51].If this parameter is specified, the code rate control method of CRF will be used for transcoding (the video code rate will no longer take effect).If there is no special requirement, it is not recommended to specify this parameter.",
																		},
																	},
																},
															},
															"audio_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Audio stream configuration parameters, when RemoveAudio is 0, this field is required.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"codec": {
																			Type:        schema.TypeString,
																			Description: "Encoding format of frequency stream.When the outer parameter Container is mp3, the optional value is:libmp3lame.When the outer parameter Container is ogg or flac, the optional value is:flac.When the outer parameter Container is m4a, the optional value is:libfdk_aac.libmp3lame.ac3.When the outer parameter Container is mp4 or flv, the optional value is:libfdk_aac: more suitable for mp4.libmp3lame: more suitable for flv.When the outer parameter Container is hls, the optional value is:libfdk_aac.libmp3lame.",
																		},
																		"bitrate": {
																			Type:        schema.TypeInt,
																			Description: "Bit rate of the audio stream, value range: 0 and [26, 256], unit: kbps.When the value is 0, it means that the audio bit rate is consistent with the original audio.",
																		},
																		"sample_rate": {
																			Type:        schema.TypeInt,
																			Description: "Sampling rate of audio stream, optional value.32000.44100.48000.Unit: Hz.",
																		},
																		"audio_channel": {
																			Type:        schema.TypeInt,
																			Description: "Audio channel mode, optional values:`1: single channel.2: Dual channel.6: Stereo.When the package format of the media is an audio format (flac, ogg, mp3, m4a), the number of channels is not allowed to be set to stereo.Default: 2.",
																		},
																	},
																},
															},
															"t_e_h_d_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Ultra-fast HD transcoding parameters.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Description: "Extremely high-definition type, optional value:TEHD-100: Extreme HD-100.Not filling means that the ultra-fast high-definition is not enabled.",
																		},
																		"max_video_bitrate": {
																			Type:        schema.TypeInt,
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
													Description: "Video transcoding custom parameters, valid when Definition is not filled with 0.When some transcoding parameters in this structure are filled in, the parameters in the transcoding template will be overwritten with the filled parameters.This parameter is used in highly customized scenarios, it is recommended that you only use Definition to specify transcoding parameters.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"container": {
																Type:        schema.TypeString,
																Description: "Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.",
															},
															"remove_video": {
																Type:        schema.TypeInt,
																Description: "Whether to remove video data, value:0: reserved.1: remove.",
															},
															"remove_audio": {
																Type:        schema.TypeInt,
																Description: "Whether to remove audio data, value:0: reserved.1: remove.",
															},
															"video_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Video streaming configuration parameters.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"codec": {
																			Type:        schema.TypeString,
																			Description: "Encoding format of the video stream, optional value:libx264: H.264 encoding.libx265: H.265 encoding.av1: AOMedia Video 1 encoding.Note: Currently H.265 encoding must specify a resolution, and it needs to be within 640*480.Note: av1 encoded containers currently only support mp4.",
																		},
																		"fps": {
																			Type:        schema.TypeInt,
																			Description: "Video frame rate, value range: [0, 100], unit: Hz.When the value is 0, it means that the frame rate is consistent with the original video.",
																		},
																		"bitrate": {
																			Type:        schema.TypeInt,
																			Description: "Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.When the value is 0, it means that the video bit rate is consistent with the original video.",
																		},
																		"resolution_adaptive": {
																			Type:        schema.TypeString,
																			Description: "Adaptive resolution, optional values:```open: open, at this time, Width represents the long side of the video, Height represents the short side of the video.close: close, at this time, Width represents the width of the video, and Height represents the height of the video.Note: In adaptive mode, Width cannot be smaller than Height.",
																		},
																		"width": {
																			Type:        schema.TypeInt,
																			Description: "The maximum value of video stream width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.",
																		},
																		"height": {
																			Type:        schema.TypeInt,
																			Description: "The maximum value of video stream height (or short side), value range: 0 and [128, 4096], unit: px.",
																		},
																		"gop": {
																			Type:        schema.TypeInt,
																			Description: "The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.When filling 0 or not filling, the system will automatically set the gop length.",
																		},
																		"fill_type": {
																			Type:        schema.TypeString,
																			Description: "Filling method, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling method:stretch: Stretch, stretch each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the aspect ratio of the video unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video unchanged, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and fill the rest of the edge with Gaussian blur.",
																		},
																		"vcrf": {
																			Type:        schema.TypeInt,
																			Description: "Video constant bit rate control factor, the value range is [1, 51], Fill in 0 to disable this parameter.If there is no special requirement, it is not recommended to specify this parameter.",
																		},
																		"content_adapt_stream": {
																			Type:        schema.TypeInt,
																			Description: "Content Adaptive Encoding. optional value:0: not open.1: open.Default: 0.When this parameter is turned on, multiple code streams with different resolutions and different bit rates will be adaptively generated. The width and height of the VideoTemplate are the maximum resolutions among the multiple code streams, and the bit rates in the VideoTemplate are multiple code rates. The highest bit rate in the stream, the vcrf in VideoTemplate is the highest quality among multiple bit streams. When the resolution, bit rate and vcrf are not set, the highest resolution generated by the ContentAdaptStream parameter is the resolution of the video source, and the video quality is close to vmaf95. To enable this parameter or learn about billing details, please contact your Tencent Cloud Business.",
																		},
																	},
																},
															},
															"audio_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Audio stream configuration parameters.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"codec": {
																			Type:        schema.TypeString,
																			Description: "Encoding format of frequency stream.When the outer parameter Container is mp3, the optional value is:libmp3lame.When the outer parameter Container is ogg or flac, the optional value is:flac.When the outer parameter Container is m4a, the optional value is:libfdk_aac.libmp3lame.ac3.When the outer parameter Container is mp4 or flv, the optional value is:libfdk_aac: more suitable for mp4.libmp3lame: more suitable for flv.When the outer parameter Container is hls, the optional value is:libfdk_aac.libmp3lame.",
																		},
																		"bitrate": {
																			Type:        schema.TypeInt,
																			Description: "Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.When the value is 0, it means that the video bit rate is consistent with the original video.",
																		},
																		"sample_rate": {
																			Type:        schema.TypeInt,
																			Description: "Sampling rate of audio stream, optional value.32000.44100.48000.Unit: Hz.",
																		},
																		"audio_channel": {
																			Type:        schema.TypeInt,
																			Description: "Audio channel mode, optional values:`1: single channel.2: Dual channel.6: Stereo.When the package format of the media is an audio format (flac, ogg, mp3, m4a), the number of channels is not allowed to be set to stereo.",
																		},
																		"stream_selects": {
																			Type: schema.TypeSet,
																			Elem: &schema.Schema{
																				Type: schema.TypeInt,
																			},
																			Description: "Specifies the audio track to preserve for the output. The default is to keep all sources.",
																		},
																	},
																},
															},
															"t_e_h_d_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Ultra-fast HD transcoding parameters.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Description: "Extremely high-definition type, optional value:TEHD-100: Extreme HD-100.Not filling means that the ultra-fast high-definition is not enabled.",
																		},
																		"max_video_bitrate": {
																			Type:        schema.TypeInt,
																			Description: "The upper limit of the video bit rate, No filling means no modification.",
																		},
																	},
																},
															},
															"subtitle_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Subtitle Stream Configuration Parameters.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": {
																			Type:        schema.TypeString,
																			Description: "The address of the subtitle file to be compressed into the video.",
																		},
																		"stream_index": {
																			Type:        schema.TypeInt,
																			Description: "Specifies the subtitle track to be compressed into the video. If there is a specified Path, the Path has a higher priority. Path and StreamIndex specify at least one.",
																		},
																		"font_type": {
																			Type:        schema.TypeString,
																			Description: "Font type.hei.ttf, song.ttf, simkai.ttf, arial.ttf.Default: hei.ttf.",
																		},
																		"font_size": {
																			Type:        schema.TypeString,
																			Description: "Font size, format: Npx, N is a value, if not specified, the subtitle file shall prevail.",
																		},
																		"font_color": {
																			Type:        schema.TypeString,
																			Description: "Font color, format: 0xRRGGBB, default value: 0xFFFFFF (white).",
																		},
																		"font_alpha": {
																			Type:        schema.TypeFloat,
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
													Description: "Watermark list, support multiple pictures or text watermarks, up to 10.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Description: "Watermark Template ID.",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Watermark custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.Watermark custom parameters do not support screenshot watermarking.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Description: "Watermark type, optional value:image: image watermark.",
																		},
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Description: "Origin position, currently only supports:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.Default: TopLeft.",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
																		},
																		"image_template": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Description: "Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_content": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Description: "The input content of the watermark image. Support jpeg, png image format.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Description: "Enter the type of source object, which supports COS and URL.",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Description: "The name of the COS Bucket where the media processing object file is located.",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Description: "Input path for media processing object files.",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
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
																						Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.Default: 10%.",
																					},
																					"height": {
																						Type:        schema.TypeString,
																						Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.",
																					},
																					"repeat_type": {
																						Type:        schema.TypeString,
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
																Description: "Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.Text watermark does not support screenshot watermarking.",
															},
															"svg_content": {
																Type:        schema.TypeString,
																Description: "SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.SVG watermark does not support screenshot watermarking.",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Description: "The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Description: "End time offset of watermark, unit: second.Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.",
															},
														},
													},
												},
												"mosaic_set": {
													Type:        schema.TypeList,
													Description: "Mosaic list, up to 10 sheets can be supported.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"coordinate_origin": {
																Type:        schema.TypeString,
																Description: "Origin position, currently only supports:TopLeft: Indicates that the coordinate origin is located in the upper left corner of the video image, and the origin of the mosaic is the upper left corner of the picture or textDefault: TopLeft.",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
															},
															"width": {
																Type:        schema.TypeString,
																Description: "The width of the mosaic. Support %, px two formats:When the string ends with %, it means that the mosaic Width is the percentage size of the video width, such as 10% means that the Width is 10% of the video width.The string ends with px, indicating that the mosaic Width unit is pixels, such as 100px indicates that the Width is 100 pixels.Default: 10%.",
															},
															"height": {
																Type:        schema.TypeString,
																Description: "The height of the mosaic. Support %, px two formats.When the string ends with %, it means that the mosaic Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the mosaic Height unit is pixel, such as 100px means that the Height is 100 pixels.Default: 10%.",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Description: "The start time offset of the mosaic, unit: second. Do not fill or fill in 0, which means that the mosaic will start to appear when the screen appears.Fill in or fill in 0, which means that the mosaic will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the mosaic appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the mosaic starts to appear n seconds before the end of the screen.",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Description: "The end time offset of the mosaic, unit: second.Fill in or fill in 0, indicating that the mosaic continues until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the mosaic lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the mosaic lasts until it disappears n seconds before the end of the screen.",
															},
														},
													},
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Description: "The start time offset of the transcoded video, unit: second.Do not fill in or fill in 0, indicating that the transcoded video starts from the beginning of the original video.When the value is greater than 0 (assumed to be n), it means that the transcoded video starts from the nth second position of the original video.When the value is less than 0 (assumed to be -n), it means that the transcoded video starts from the position n seconds before the end of the original video.",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Description: "End time offset of video after transcoding, unit: second.Do not fill in or fill in 0, indicating that the transcoded video continues until the end of the original video..When the value is greater than 0 (assumed to be n), it means that the transcoded video lasts until the nth second of the original video and terminates.When the value is less than 0 (assumed to be -n), it means that the transcoded video lasts until n seconds before the end of the original video..",
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "The target storage of the transcoded file, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Description: "The type of media processing output object storage location, now only supports COS.",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
																		},
																		"region": {
																			Type:        schema.TypeString,
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
													Description: "The output path of the main file after transcoding can be a relative path or an absolute path. If not filled, the default is a relative path: {inputName}_transcode_{definition}.{format}.",
												},
												"segment_object_name": {
													Type:        schema.TypeString,
													Description: "The output path of the transcoded fragment file (the path of ts when transcoding HLS), can only be a relative path. If not filled, the default is: `{inputName}_transcode_{definition}_{number}.{format}.",
												},
												"object_number_format": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Rules for the `{number}` variable in the output path after transcoding.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"initial_value": {
																Type:        schema.TypeInt,
																Description: "The starting value of `{number}` variable, the default is 0.",
															},
															"increment": {
																Type:        schema.TypeInt,
																Description: "The growth step of the `{number}` variable, the default is 1.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Description: "The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.",
															},
															"place_holder": {
																Type:        schema.TypeString,
																Description: "When the length of the `{number}` variable is insufficient, a placeholder is added. Default is 0.",
															},
														},
													},
												},
												"head_tail_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Opening and ending parameters.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"head_set": {
																Type:        schema.TypeList,
																Description: "Title list.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Description: "Enter the type of source object, which supports COS and URL.",
																		},
																		"cos_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Description: "The name of the COS Bucket where the media processing object file is located.",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																					},
																					"object": {
																						Type:        schema.TypeString,
																						Description: "Input path for media processing object files.",
																					},
																				},
																			},
																		},
																		"url_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
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
																Description: "Ending List.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Description: "Enter the type of source object, which supports COS and URL.",
																		},
																		"cos_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Description: "The name of the COS Bucket where the media processing object file is located.",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																					},
																					"object": {
																						Type:        schema.TypeString,
																						Description: "Input path for media processing object files.",
																					},
																				},
																			},
																		},
																		"url_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
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
										Description: "Video Rotation Map Task List.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Description: "Video turntable template id.",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Description: "The start time of the animation in the video, in seconds.",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Description: "The end time of the animation in the video, in seconds.",
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "The target storage of the transcoded file, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Description: "The type of media processing output object storage location, now only supports COS.",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
																		},
																		"region": {
																			Type:        schema.TypeString,
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
													Description: "The output path of the file after rotating the image, which can be a relative path or an absolute path. If not filled, the default is a relative path: {inputName}_animatedGraphic_{definition}.{format}.",
												},
											},
										},
									},
									"snapshot_by_time_offset_task_set": {
										Type:        schema.TypeList,
										Description: "Screenshot the task list of the video according to the time point.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Description: "Specified time point screenshot template ID.",
												},
												"ext_time_offset_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "Screenshot time point list, the time point supports two formats: s and %:When the string ends with s, it means that the time point is in seconds, such as 3.5s means that the time point is the 3.5th second.When the string ends with %, it means that the time point is the percentage of the video duration, such as 10% means that the time point is the first 10% of the time in the video.",
												},
												"time_offset_set": {
													Type:        schema.TypeList,
													Description: "Screenshot time point list, the unit is &amp;lt;font color=red&amp;gt;seconds&amp;lt;/font&amp;gt;. This parameter is no longer recommended, it is recommended that you use the ExtTimeOffsetSet parameter.",
												},
												"watermark_set": {
													Type:        schema.TypeList,
													Description: "Watermark list, support multiple pictures or text watermarks, up to 10.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Description: "Watermark Template ID.",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Watermark custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.Watermark custom parameters do not support screenshot watermarking.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Description: "Watermark type, optional value:image: image watermark.",
																		},
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Description: "Origin position, currently only supports:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.Default: TopLeft.",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
																		},
																		"image_template": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Description: "Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_content": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Description: "The input content of the watermark image. Support jpeg, png image format.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Description: "Enter the type of source object, which supports COS and URL.",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Description: "The name of the COS Bucket where the media processing object file is located.",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Description: "Input path for media processing object files.",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
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
																						Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.Default: 10%.",
																					},
																					"height": {
																						Type:        schema.TypeString,
																						Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.",
																					},
																					"repeat_type": {
																						Type:        schema.TypeString,
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
																Description: "Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.Text watermark does not support screenshot watermarking.",
															},
															"svg_content": {
																Type:        schema.TypeString,
																Description: "SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.SVG watermark does not support screenshot watermarking.",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Description: "The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Description: "End time offset of watermark, unit: second.Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.",
															},
														},
													},
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "The target storage of the file after the screenshot at the time point, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Description: "The type of media processing output object storage location, now only supports COS.",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
																		},
																		"region": {
																			Type:        schema.TypeString,
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
													Description: "The output path of the picture file after the snapshot at the time point can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_snapshotByTimeOffset_{definition}_{number}.{format}`.",
												},
												"object_number_format": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Rules for the `{number}` variable in the output path after the screenshot at the time point.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"initial_value": {
																Type:        schema.TypeInt,
																Description: "The starting value of `{number}` variable, the default is 0.",
															},
															"increment": {
																Type:        schema.TypeInt,
																Description: "The growth step of `{number}` variable, default is 1.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Description: "The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.",
															},
															"place_holder": {
																Type:        schema.TypeString,
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
										Description: "Screenshot task list for video sampling.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Description: "Sample screenshot template ID.",
												},
												"watermark_set": {
													Type:        schema.TypeList,
													Description: "Watermark list, support multiple pictures or text watermarks, up to 10.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Description: "Watermark Template ID.",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Watermark custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.Watermark custom parameters do not support screenshot watermarking.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Description: "Watermark type, optional value:image: image watermark.",
																		},
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Description: "Origin position, currently only supports:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.Default: TopLeft.",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
																		},
																		"image_template": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Description: "Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_content": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Description: "The input content of the watermark image. Support jpeg, png image format.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Description: "Enter the type of source object, which supports COS and URL.",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Description: "The name of the COS Bucket where the media processing object file is located.",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Description: "Input path for media processing object files.",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
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
																						Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.Default: 10%.",
																					},
																					"height": {
																						Type:        schema.TypeString,
																						Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.",
																					},
																					"repeat_type": {
																						Type:        schema.TypeString,
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
																Description: "Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.Text watermark does not support screenshot watermarking.",
															},
															"svg_content": {
																Type:        schema.TypeString,
																Description: "SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.SVG watermark does not support screenshot watermarking.",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Description: "The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Description: "End time offset of watermark, unit: second.Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.",
															},
														},
													},
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "The target storage of the file after the screenshot at the time point, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Description: "The type of media processing output object storage location, now only supports COS.",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
																		},
																		"region": {
																			Type:        schema.TypeString,
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
													Description: "The output path of the image file after sampling the screenshot, which can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_sampleSnapshot_{definition}_{number}.{format}`.",
												},
												"object_number_format": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Rules for the `{number}` variable in the output path after sampling the screenshot.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"initial_value": {
																Type:        schema.TypeInt,
																Description: "The starting value of `{number}` variable, the default is 0.",
															},
															"increment": {
																Type:        schema.TypeInt,
																Description: "The growth step of the `{number}` variable, the default is 1.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Description: "The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.",
															},
															"place_holder": {
																Type:        schema.TypeString,
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
										Description: "Sprite image capture task list for video.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Description: "Sprite Illustration Template ID.",
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "The target storage of the file after the sprite image is intercepted, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Description: "The type of media processing output object storage location, now only supports COS.",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
																		},
																		"region": {
																			Type:        schema.TypeString,
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
													Description: "After capturing the sprite image, the output path of the sprite image file can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_imageSprite_{definition}_{number}.{format}`.",
												},
												"web_vtt_object_name": {
													Type:        schema.TypeString,
													Description: "After capturing the sprite image, the output path of the Web VTT file can only be a relative path. If not filled, the default is a relative path: `{inputName}_imageSprite_{definition}.{format}`.",
												},
												"object_number_format": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Rules for the `{number}` variable in the output path after intercepting the Sprite image.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"initial_value": {
																Type:        schema.TypeInt,
																Description: "The starting value of `{number}` variable, the default is 0.",
															},
															"increment": {
																Type:        schema.TypeInt,
																Description: "The growth step of the `{number}` variable, the default is 1.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Description: "The minimum length of the `{number}` variable, if insufficient, placeholders will be filled. Default is 1.",
															},
															"place_holder": {
																Type:        schema.TypeString,
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
										Description: "Transfer Adaptive Code Stream Task List.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Description: "Transfer Adaptive Code Stream Template ID.",
												},
												"watermark_set": {
													Type:        schema.TypeList,
													Description: "Watermark list, support multiple pictures or text watermarks, up to 10.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Description: "Watermark Template ID.",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Watermark custom parameters, valid when Definition is filled with 0.This parameter is used in highly customized scenarios, it is recommended that you use Definition to specify watermark parameters first.Watermark custom parameters do not support screenshot watermarking.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Description: "Watermark type, optional value:image: image watermark.",
																		},
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Description: "Origin position, currently only supports:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.Default: TopLeft.",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default: 0px.",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats:When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default: 0px.",
																		},
																		"image_template": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Description: "Image watermark template, when Type is image, this field is required. When Type is text, this field is invalid.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_content": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Description: "The input content of the watermark image. Support jpeg, png image format.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Description: "Enter the type of source object, which supports COS and URL.",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Description: "The name of the COS Bucket where the media processing object file is located.",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Description: "Input path for media processing object files.",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
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
																						Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels.Default: 10%.",
																					},
																					"height": {
																						Type:        schema.TypeString,
																						Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.Default value: 0px, indicating that Height is scaled according to the aspect ratio of the original watermark image.",
																					},
																					"repeat_type": {
																						Type:        schema.TypeString,
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
																Description: "Text content, the length does not exceed 100 characters. Fill in only when the watermark type is text watermark.Text watermark does not support screenshot watermarking.",
															},
															"svg_content": {
																Type:        schema.TypeString,
																Description: "SVG content. The length cannot exceed 2000000 characters. Fill in only if the watermark type is SVG watermark.SVG watermark does not support screenshot watermarking.",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Description: "The start time offset of the watermark, unit: second. Do not fill in or fill in 0, which means that the watermark will start to appear when the screen appears.Do not fill in or fill in 0, which means the watermark will appear from the beginning of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark appears from the nth second of the screen.When the value is less than 0 (assumed to be -n), it means that the watermark starts to appear n seconds before the end of the screen.",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Description: "End time offset of watermark, unit: second.Do not fill in or fill in 0, indicating that the watermark lasts until the end of the screen.When the value is greater than 0 (assumed to be n), it means that the watermark lasts until the nth second and disappears.When the value is less than 0 (assumed to be -n), it means that the watermark lasts until it disappears n seconds before the end of the screen.",
															},
														},
													},
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "The target storage of the file after converting to the adaptive code stream, if not filled, it will inherit the OutputStorage value of the upper layer.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Description: "The type of media processing output object storage location, now only supports COS.",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Description: "Valid when Type is COS, this item is required, indicating the media processing COS output location.Note: This field may return null, indicating that no valid value can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Description: "The target Bucket name of the file output generated by media processing, if not filled, it means the upper layer.",
																		},
																		"region": {
																			Type:        schema.TypeString,
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
													Description: "After converting to an adaptive stream, the output path of the manifest file can be a relative path or an absolute path. If not filled, the default is a relative path: `{inputName}_adaptiveDynamicStreaming_{definition}.{format}`.",
												},
												"sub_stream_object_name": {
													Type:        schema.TypeString,
													Description: "After converting to an adaptive stream, the output path of the sub-stream file can only be a relative path. If not filled, the default is a relative path: {inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{format}`.",
												},
												"segment_object_name": {
													Type:        schema.TypeString,
													Description: "After converting to an adaptive stream (only HLS), the output path of the fragmented file can only be a relative path. If not filled, the default is a relative path: `{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{format}`.",
												},
											},
										},
									},
								},
							},
						},
						"ai_content_review_task": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Video Content Moderation Type Task Parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Description: "Video Content Review Template ID.",
									},
								},
							},
						},
						"ai_analysis_task": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Video Content Analysis Type Task Parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Description: "Video Content Analysis Template ID.",
									},
									"extended_parameter": {
										Type:        schema.TypeString,
										Description: "Extension parameter whose value is a serialized json string.Note: This parameter is a customized demand parameter, which requires offline docking.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"ai_recognition_task": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Video content recognition type task parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Description: "Video Intelligent Recognition Template ID.",
									},
								},
							},
						},
						"task_notify_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "The event notification configuration of the task, if it is not filled, it means that the event notification will not be obtained.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cmq_model": {
										Type:        schema.TypeString,
										Description: "CMQ or TDMQ-CMQ model, there are two kinds of Queue and Topic.",
									},
									"cmq_region": {
										Type:        schema.TypeString,
										Description: "Region of CMQ or TDMQ-CMQ, such as sh, bj, etc.",
									},
									"topic_name": {
										Type:        schema.TypeString,
										Description: "Valid when the model is a Topic, indicating the topic name of the CMQ or TDMQ-CMQ that receives event notifications.",
									},
									"queue_name": {
										Type:        schema.TypeString,
										Description: "Valid when the model is Queue, indicating the queue name of the CMQ or TDMQ-CMQ that receives the event notification.",
									},
									"notify_mode": {
										Type:        schema.TypeString,
										Description: "The mode of the workflow notification, the possible values are Finish and Change, leaving blank means Finish.",
									},
									"notify_type": {
										Type:        schema.TypeString,
										Description: "Notification type, optional value:CMQ: offline, it is recommended to switch to TDMQ-CMQ.TDMQ-CMQ: message queue.URL: When the URL is specified, the HTTP callback is pushed to the address specified by NotifyUrl, the callback protocol is http+json, and the package body content is the same as the output parameters of the parsing event notification interface.SCF: not recommended, additional configuration of SCF in the console is required.Note: CMQ is the default when not filled or empty, if you need to use other types, you need to fill in the corresponding type value.",
									},
									"notify_url": {
										Type:        schema.TypeString,
										Description: "HTTP callback address, required when NotifyType is URL.",
									},
								},
							},
						},
						"task_priority": {
							Type:        schema.TypeInt,
							Description: "The priority of the workflow, the larger the value, the higher the priority, the value range is -10 to 10, and blank means 0.",
						},
						"output_dir": {
							Type:        schema.TypeString,
							Description: "The target directory of the output file generated by media processing.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Workflow creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Workflow last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
					},
				},
			},

			"request_id": {
				Type:        schema.TypeString,
				Description: "Unique request ID, returned for every request. The RequestId of the request needs to be provided when locating the problem.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMpsWorkflowsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_workflows.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("workflow_ids"); ok {
		workflowIdsSet := v.(*schema.Set).List()
		for i := range workflowIdsSet {
			workflowIds := workflowIdsSet[i].(int)
			paramMap["WorkflowIds"] = append(paramMap["WorkflowIds"], helper.IntInt64(workflowIds))
		}
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("workflow_info_set"); ok {
		workflowInfoSetSet := v.([]interface{})
		tmpSet := make([]*mps.WorkflowInfo, 0, len(workflowInfoSetSet))

		for _, item := range workflowInfoSetSet {
			workflowInfo := mps.WorkflowInfo{}
			workflowInfoMap := item.(map[string]interface{})

			if v, ok := workflowInfoMap["workflow_id"]; ok {
				workflowInfo.WorkflowId = helper.IntInt64(v.(int))
			}
			if v, ok := workflowInfoMap["workflow_name"]; ok {
				workflowInfo.WorkflowName = helper.String(v.(string))
			}
			if v, ok := workflowInfoMap["status"]; ok {
				workflowInfo.Status = helper.String(v.(string))
			}
			if triggerMap, ok := helper.InterfaceToMap(workflowInfoMap, "trigger"); ok {
				workflowTrigger := mps.WorkflowTrigger{}
				if v, ok := triggerMap["type"]; ok {
					workflowTrigger.Type = helper.String(v.(string))
				}
				if cosFileUploadTriggerMap, ok := helper.InterfaceToMap(triggerMap, "cos_file_upload_trigger"); ok {
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
						cosFileUploadTrigger.Formats = helper.InterfacesStringsPoint(formatsSet)
					}
					workflowTrigger.CosFileUploadTrigger = &cosFileUploadTrigger
				}
				workflowInfo.Trigger = &workflowTrigger
			}
			if outputStorageMap, ok := helper.InterfaceToMap(workflowInfoMap, "output_storage"); ok {
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
				workflowInfo.OutputStorage = &taskOutputStorage
			}
			if mediaProcessTaskMap, ok := helper.InterfaceToMap(workflowInfoMap, "media_process_task"); ok {
				mediaProcessTaskInput := mps.MediaProcessTaskInput{}
				if v, ok := mediaProcessTaskMap["transcode_task_set"]; ok {
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
									videoTemplateInfo.Fps = helper.IntUint64(v.(int))
								}
								if v, ok := videoTemplateMap["bitrate"]; ok {
									videoTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
									audioTemplateInfo.Bitrate = helper.IntUint64(v.(int))
								}
								if v, ok := audioTemplateMap["sample_rate"]; ok {
									audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
								}
								if v, ok := audioTemplateMap["audio_channel"]; ok {
									audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
								}
								rawTranscodeParameter.AudioTemplate = &audioTemplateInfo
							}
							if tEHDConfigMap, ok := helper.InterfaceToMap(rawParameterMap, "t_e_h_d_config"); ok {
								tEHDConfig := mps.TEHDConfig{}
								if v, ok := tEHDConfigMap["type"]; ok {
									tEHDConfig.Type = helper.String(v.(string))
								}
								if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
									tEHDConfig.MaxVideoBitrate = helper.IntUint64(v.(int))
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
									videoTemplateInfoForUpdate.Fps = helper.IntUint64(v.(int))
								}
								if v, ok := videoTemplateMap["bitrate"]; ok {
									videoTemplateInfoForUpdate.Bitrate = helper.IntUint64(v.(int))
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
									audioTemplateInfoForUpdate.Bitrate = helper.IntUint64(v.(int))
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
							if tEHDConfigMap, ok := helper.InterfaceToMap(overrideParameterMap, "t_e_h_d_config"); ok {
								tEHDConfigForUpdate := mps.TEHDConfigForUpdate{}
								if v, ok := tEHDConfigMap["type"]; ok {
									tEHDConfigForUpdate.Type = helper.String(v.(string))
								}
								if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
									tEHDConfigForUpdate.MaxVideoBitrate = helper.IntUint64(v.(int))
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
				if v, ok := mediaProcessTaskMap["animated_graphic_task_set"]; ok {
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
				if v, ok := mediaProcessTaskMap["snapshot_by_time_offset_task_set"]; ok {
					for _, item := range v.([]interface{}) {
						snapshotByTimeOffsetTaskSetMap := item.(map[string]interface{})
						snapshotByTimeOffsetTaskInput := mps.SnapshotByTimeOffsetTaskInput{}
						if v, ok := snapshotByTimeOffsetTaskSetMap["definition"]; ok {
							snapshotByTimeOffsetTaskInput.Definition = helper.IntUint64(v.(int))
						}
						if v, ok := snapshotByTimeOffsetTaskSetMap["ext_time_offset_set"]; ok {
							extTimeOffsetSetSet := v.(*schema.Set).List()
							snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet = helper.InterfacesStringsPoint(extTimeOffsetSetSet)
						}
						if v, ok := snapshotByTimeOffsetTaskSetMap["time_offset_set"]; ok {
							for _, item := range v.([]interface{}) {
								timeOffsetSetMap := item.(map[string]interface{})
								float := mps.float{}
								snapshotByTimeOffsetTaskInput.TimeOffsetSet = append(snapshotByTimeOffsetTaskInput.TimeOffsetSet, &float)
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
				if v, ok := mediaProcessTaskMap["sample_snapshot_task_set"]; ok {
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
				if v, ok := mediaProcessTaskMap["image_sprite_task_set"]; ok {
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
				if v, ok := mediaProcessTaskMap["adaptive_dynamic_streaming_task_set"]; ok {
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
				workflowInfo.MediaProcessTask = &mediaProcessTaskInput
			}
			if aiContentReviewTaskMap, ok := helper.InterfaceToMap(workflowInfoMap, "ai_content_review_task"); ok {
				aiContentReviewTaskInput := mps.AiContentReviewTaskInput{}
				if v, ok := aiContentReviewTaskMap["definition"]; ok {
					aiContentReviewTaskInput.Definition = helper.IntUint64(v.(int))
				}
				workflowInfo.AiContentReviewTask = &aiContentReviewTaskInput
			}
			if aiAnalysisTaskMap, ok := helper.InterfaceToMap(workflowInfoMap, "ai_analysis_task"); ok {
				aiAnalysisTaskInput := mps.AiAnalysisTaskInput{}
				if v, ok := aiAnalysisTaskMap["definition"]; ok {
					aiAnalysisTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := aiAnalysisTaskMap["extended_parameter"]; ok {
					aiAnalysisTaskInput.ExtendedParameter = helper.String(v.(string))
				}
				workflowInfo.AiAnalysisTask = &aiAnalysisTaskInput
			}
			if aiRecognitionTaskMap, ok := helper.InterfaceToMap(workflowInfoMap, "ai_recognition_task"); ok {
				aiRecognitionTaskInput := mps.AiRecognitionTaskInput{}
				if v, ok := aiRecognitionTaskMap["definition"]; ok {
					aiRecognitionTaskInput.Definition = helper.IntUint64(v.(int))
				}
				workflowInfo.AiRecognitionTask = &aiRecognitionTaskInput
			}
			if taskNotifyConfigMap, ok := helper.InterfaceToMap(workflowInfoMap, "task_notify_config"); ok {
				taskNotifyConfig := mps.TaskNotifyConfig{}
				if v, ok := taskNotifyConfigMap["cmq_model"]; ok {
					taskNotifyConfig.CmqModel = helper.String(v.(string))
				}
				if v, ok := taskNotifyConfigMap["cmq_region"]; ok {
					taskNotifyConfig.CmqRegion = helper.String(v.(string))
				}
				if v, ok := taskNotifyConfigMap["topic_name"]; ok {
					taskNotifyConfig.TopicName = helper.String(v.(string))
				}
				if v, ok := taskNotifyConfigMap["queue_name"]; ok {
					taskNotifyConfig.QueueName = helper.String(v.(string))
				}
				if v, ok := taskNotifyConfigMap["notify_mode"]; ok {
					taskNotifyConfig.NotifyMode = helper.String(v.(string))
				}
				if v, ok := taskNotifyConfigMap["notify_type"]; ok {
					taskNotifyConfig.NotifyType = helper.String(v.(string))
				}
				if v, ok := taskNotifyConfigMap["notify_url"]; ok {
					taskNotifyConfig.NotifyUrl = helper.String(v.(string))
				}
				workflowInfo.TaskNotifyConfig = &taskNotifyConfig
			}
			if v, ok := workflowInfoMap["task_priority"]; ok {
				workflowInfo.TaskPriority = helper.IntInt64(v.(int))
			}
			if v, ok := workflowInfoMap["output_dir"]; ok {
				workflowInfo.OutputDir = helper.String(v.(string))
			}
			if v, ok := workflowInfoMap["create_time"]; ok {
				workflowInfo.CreateTime = helper.String(v.(string))
			}
			if v, ok := workflowInfoMap["update_time"]; ok {
				workflowInfo.UpdateTime = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &workflowInfo)
		}
		paramMap["workflow_info_set"] = tmpSet
	}

	if v, ok := d.GetOk("request_id"); ok {
		paramMap["RequestId"] = helper.String(v.(string))
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var workflowInfoSet []*mps.WorkflowInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsWorkflowsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		workflowInfoSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(workflowInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(workflowInfoSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
