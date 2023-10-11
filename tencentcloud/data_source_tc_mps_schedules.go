/*
Use this data source to query detailed information of mps schedules

Example Usage

Query the enabled schedules.

```hcl
data "tencentcloud_mps_schedules" "schedules" {
  status       = "Enabled"
}
```

Query the specified one.

```hcl
data "tencentcloud_mps_schedules" "schedules" {
  schedule_ids = [%d]
  trigger_type = "CosFileUpload"
  status       = "Enabled"
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

func dataSourceTencentCloudMpsSchedules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsSchedulesRead,
		Schema: map[string]*schema.Schema{
			"schedule_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The IDs of the schemes to query. Array length limit: 100.",
			},

			"trigger_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The trigger type. Valid values:`CosFileUpload`: The scheme is triggered when a file is uploaded to Tencent Cloud Object Storage (COS).`AwsS3FileUpload`: The scheme is triggered when a file is uploaded to AWS S3.If you do not specify this parameter or leave it empty, all schemes will be returned regardless of the trigger type.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The scheme status. Valid values:`Enabled`, `Disabled`. If you do not specify this parameter, all schemes will be returned regardless of the status.",
			},

			"schedule_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The information of the schemes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The scheme ID.",
						},
						"schedule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The scheme name.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The scheme status. Valid values:`Enabled``Disabled`Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"trigger": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The trigger of the scheme.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The trigger type. Valid values:`CosFileUpload`: Tencent Cloud COS trigger.`AwsS3FileUpload`: AWS S3 trigger. Currently, this type is only supported for transcoding tasks and schemes (not supported for workflows).",
									},
									"cos_file_upload_trigger": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "This parameter is required and valid when `Type` is `CosFileUpload`, indicating the COS trigger rule.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the COS bucket bound to a workflow, such as `TopRankVideo-125xxx88`.",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Region of the COS bucket bound to a workflow, such as `ap-chongiqng`.",
												},
												"dir": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Input path directory bound to a workflow, such as `/movie/201907/`. If this parameter is left empty, the `/` root directory will be used.",
												},
												"formats": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Format list of files that can trigger a workflow, such as [mp4, flv, mov]. If this parameter is left empty, files in all formats can trigger the workflow.",
												},
											},
										},
									},
									"aws_s3_file_upload_trigger": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The AWS S3 trigger. This parameter is valid and required if `Type` is `AwsS3FileUpload`.Note: Currently, the key for the AWS S3 bucket, the trigger SQS queue, and the callback SQS queue must be the same.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"s3_bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The AWS S3 bucket bound to the scheme.",
												},
												"s3_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The region of the AWS S3 bucket.",
												},
												"dir": {
													Type:     schema.TypeString,
													Computed: true,
													Description: "The bucket directory bound. It must be an absolute path that starts and ends with `/`, such as `/movie/201907/`. If you do not specify this, the root directory will be bound.	.",
												},
												"formats": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed: true,
													Description: "The file formats that will trigger the scheme, such as [mp4, flv, mov]. If you do not specify this, the upload of files in any format will trigger the scheme.	.",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The key ID of the AWS S3 bucket.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The key of the AWS S3 bucket.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"aws_sqs": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The SQS queue of the AWS S3 bucket.Note: The queue must be in the same region as the bucket.Note: This field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"sqs_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The region of the SQS queue.",
															},
															"sqs_queue_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of the SQS queue.",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The key ID required to read from/write to the SQS queue.",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The key required to read from/write to the SQS queue.",
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
						"activities": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subtasks of the scheme.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"activity_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The subtask type.`input`: The start.`output`: The end.`action-trans`: Transcoding.`action-samplesnapshot`: Sampled screencapturing.`action-AIAnalysis`: Content analysis.`action-AIRecognition`: Content recognition.`action-aiReview`: Content moderation.`action-animated-graphics`: Animated screenshot generation.`action-image-sprite`: Image sprite generation.`action-snapshotByTimeOffset`: Time point screencapturing.`action-adaptive-substream`: Adaptive bitrate streaming.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"reardrive_index": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "The indexes of the subsequent actions.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"activity_para": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The parameters of a subtask.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"transcode_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A transcoding task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID of a video transcoding template.",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Custom video transcoding parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the transcoding parameter preferably.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"container": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Container. Valid values: mp4; flv; hls; mp3; flac; ogg; m4a. Among them, mp3, flac, ogg, and m4a are for audio files.",
																		},
																		"remove_video": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Whether to remove video data. Valid values:0: retain;1: remove.Default value: 0.",
																		},
																		"remove_audio": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Whether to remove audio data. Valid values:0: retain;1: remove.Default value: 0.",
																		},
																		"video_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Video stream configuration parameter. This field is required when `RemoveVideo` is 0.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The video codec. Valid values:`libx264`: H.264`libx265`: H.265`av1`: AOMedia Video 1Note: You must specify a resolution (not higher than 640 x 480) if the H.265 codec is used.Note: You can only use the AOMedia Video 1 codec for MP4 files.",
																					},
																					"fps": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The video frame rate (Hz). Value range: [0, 100].If the value is 0, the frame rate will be the same as that of the source video.Note: For adaptive bitrate streaming, the value range of this parameter is [0, 60].",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The video bitrate (Kbps). Value range: 0 and [128, 35000].If the value is 0, the bitrate of the video will be the same as that of the source video.",
																					},
																					"resolution_adaptive": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resolution adaption. Valid values:open: Enabled. When resolution adaption is enabled, `Width` indicates the long side of a video, while `Height` indicates the short side.close: Disabled. When resolution adaption is disabled, `Width` indicates the width of a video, while `Height` indicates the height.Default value: open.Note: When resolution adaption is enabled, `Width` cannot be smaller than `Height`.",
																					},
																					"width": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum value of the width (or long side) of a video stream in px. Value range: 0 and [128, 4,096].If both `Width` and `Height` are 0, the resolution will be the same as that of the source video;If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled;If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled;If both `Width` and `Height` are not 0, the custom resolution will be used.Default value: 0.",
																					},
																					"height": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum value of the height (or short side) of a video stream in px. Value range: 0 and [128, 4,096].If both `Width` and `Height` are 0, the resolution will be the same as that of the source video;If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled;If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled;If both `Width` and `Height` are not 0, the custom resolution will be used.Default value: 0.",
																					},
																					"gop": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Frame interval between I keyframes. Value range: 0 and [1,100000].If this parameter is 0 or left empty, the system will automatically set the GOP length.",
																					},
																					"fill_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The fill mode, which indicates how a video is resized when the video&#39;s original aspect ratio is different from the target aspect ratio. Valid values:stretch: Stretch the image frame by frame to fill the entire screen. The video image may become squashed or stretched after transcoding.black: Keep the image&#39;s original aspect ratio and fill the blank space with black bars.white: Keep the image&#39;s original aspect ratio and fill the blank space with white bars.gauss: Keep the image&#39;s original aspect ratio and apply Gaussian blur to the blank space.Default value: black.Note: Only `stretch` and `black` are supported for adaptive bitrate streaming.",
																					},
																					"vcrf": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The control factor of video constant bitrate. Value range: [1, 51]If this parameter is specified, CRF (a bitrate control method) will be used for transcoding. (Video bitrate will no longer take effect.)It is not recommended to specify this parameter if there are no special requirements.",
																					},
																				},
																			},
																		},
																		"audio_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Audio stream configuration parameter. This field is required when `RemoveAudio` is 0.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Audio stream codec.When the outer `Container` parameter is `mp3`, the valid value is:libmp3lame.When the outer `Container` parameter is `ogg` or `flac`, the valid value is:flac.When the outer `Container` parameter is `m4a`, the valid values include:libfdk_aac;libmp3lame;ac3.When the outer `Container` parameter is `mp4` or `flv`, the valid values include:libfdk_aac: more suitable for mp4;libmp3lame: more suitable for flv.When the outer `Container` parameter is `hls`, the valid values include:libfdk_aac;libmp3lame.",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio stream bitrate in Kbps. Value range: 0 and [26, 256].If the value is 0, the bitrate of the audio stream will be the same as that of the original audio.",
																					},
																					"sample_rate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio stream sample rate. Valid values:32,00044,10048,000In Hz.",
																					},
																					"audio_channel": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio channel system. Valid values:1: Mono2: Dual6: StereoWhen the media is packaged in audio format (FLAC, OGG, MP3, M4A), the sound channel cannot be set to stereo.Default value: 2.",
																					},
																				},
																			},
																		},
																		"tehd_config": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "TESHD transcoding parameter.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "TESHD type. Valid values:`TEHD-100`: TESHD-100. If this parameter is left empty, TESHD will not be enabled.",
																					},
																					"max_video_bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum bitrate, which is valid when `Type` is `TESHD`. If this parameter is left empty or 0 is entered, there will be no upper limit for bitrate.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"override_parameter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Video transcoding custom parameter, which is valid when `Definition` is not 0.When any parameters in this structure are entered, they will be used to override corresponding parameters in templates.This parameter is used in highly customized scenarios. We recommend you only use `Definition` to specify the transcoding parameter.Note: this field may return `null`, indicating that no valid value was found.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"container": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Container format. Valid values: mp4, flv, hls, mp3, flac, ogg, and m4a; mp3, flac, ogg, and m4a are formats of audio files.",
																		},
																		"remove_video": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Whether to remove video data. Valid values:0: retain1: remove.",
																		},
																		"remove_audio": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Whether to remove audio data. Valid values:0: retain1: remove.",
																		},
																		"video_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Video stream configuration parameter.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The video codec. Valid values:libx264: H.264libx265: H.265av1: AOMedia Video 1Note: You must specify a resolution (not higher than 640 x 480) if the H.265 codec is used.Note: You can only use the AOMedia Video 1 codec for MP4 files.",
																					},
																					"fps": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Video frame rate in Hz. Value range: [0, 100].If the value is 0, the frame rate will be the same as that of the source video.",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Bitrate of a video stream in Kbps. Value range: 0 and [128, 35,000].If the value is 0, the bitrate of the video will be the same as that of the source video.",
																					},
																					"resolution_adaptive": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resolution adaption. Valid values:open: Enabled. When resolution adaption is enabled, `Width` indicates the long side of a video, while `Height` indicates the short side.close: Disabled. When resolution adaption is disabled, `Width` indicates the width of a video, while `Height` indicates the height.Note: When resolution adaption is enabled, `Width` cannot be smaller than `Height`.",
																					},
																					"width": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum value of the width (or long side) of a video stream in px. Value range: 0 and [128, 4,096].If both `Width` and `Height` are 0, the resolution will be the same as that of the source video;If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled;If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled;If both `Width` and `Height` are not 0, the custom resolution will be used.",
																					},
																					"height": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum value of the height (or short side) of a video stream in px. Value range: 0 and [128, 4,096].",
																					},
																					"gop": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Frame interval between I keyframes. Value range: 0 and [1,100000]. If this parameter is 0, the system will automatically set the GOP length.",
																					},
																					"fill_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: stretch: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer;black: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks.white: fill with white. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with white color blocks.gauss: fill with Gaussian blur. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with Gaussian blur.",
																					},
																					"vcrf": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The control factor of video constant bitrate. Value range: [0, 51]. This parameter will be disabled if you enter `0`.It is not recommended to specify this parameter if there are no special requirements.",
																					},
																					"content_adapt_stream": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Whether to enable adaptive encoding. Valid values:0: Disable1: EnableDefault value: 0. If this parameter is set to `1`, multiple streams with different resolutions and bitrates will be generated automatically. The highest resolution, bitrate, and quality of the streams are determined by the values of `width` and `height`, `Bitrate`, and `Vcrf` in `VideoTemplate` respectively. If these parameters are not set in `VideoTemplate`, the highest resolution generated will be the same as that of the source video, and the highest video quality will be close to VMAF 95. To use this parameter or learn about the billing details of adaptive encoding, please contact your sales rep.",
																					},
																				},
																			},
																		},
																		"audio_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Audio stream configuration parameter.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Audio stream codec.When the outer `Container` parameter is `mp3`, the valid value is:libmp3lame.When the outer `Container` parameter is `ogg` or `flac`, the valid value is:flac.When the outer `Container` parameter is `m4a`, the valid values include:libfdk_aac;libmp3lame;ac3.When the outer `Container` parameter is `mp4` or `flv`, the valid values include:libfdk_aac: More suitable for mp4;libmp3lame: More suitable for flv;mp2.When the outer `Container` parameter is `hls`, the valid values include:libfdk_aac;libmp3lame.",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio stream bitrate in Kbps. Value range: 0 and [26, 256]. If the value is 0, the bitrate of the audio stream will be the same as that of the original audio.",
																					},
																					"sample_rate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio stream sample rate. Valid values:32,00044,10048,000In Hz.",
																					},
																					"audio_channel": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio channel system. Valid values:1: Mono2: Dual6: StereoWhen the media is packaged in audio format (FLAC, OGG, MP3, M4A), the sound channel cannot be set to stereo.",
																					},
																					"stream_selects": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeInt,
																						},
																						Computed:    true,
																						Description: "The audio tracks to retain. All audio tracks are retained by default.",
																					},
																				},
																			},
																		},
																		"tehd_config": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The TSC transcoding parameters.Note: This field may return null, indicating that no valid values can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The TSC type. Valid values:`TEHD-100`: TSC-100 (video TSC). `TEHD-200`: TSC-200 (audio TSC). If this parameter is left blank, no modification will be made.Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																					"max_video_bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The maximum video bitrate. If this parameter is not specified, no modifications will be made.Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																				},
																			},
																		},
																		"subtitle_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The subtitle settings.Note: This field may return null, indicating that no valid values can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"path": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The URL of the subtitles to add to the video.Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																					"stream_index": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The subtitle track to add to the video. If both `Path` and `StreamIndex` are specified, `Path` will be used. You need to specify at least one of the two parameters.Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																					"font_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The font. Valid values:`hei.ttf`: Heiti.`song.ttf`: Songti.`simkai.ttf`: Kaiti.`arial.ttf`: Arial.The default is `hei.ttf`.Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																					"font_size": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The font size (pixels). If this is not specified, the font size in the subtitle file will be used.Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																					"font_color": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The font color in 0xRRGGBB format. Default value: 0xFFFFFF (white).Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																					"font_alpha": {
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "The text transparency. Value range: 0-1.`0`: Fully transparent.`1`: Fully opaque.Default value: 1.Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																				},
																			},
																		},
																		"addon_audio_stream": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The information of the external audio track to add.Note: This field may return null, indicating that no valid values can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The region of the COS bucket, such as `ap-chongqing`.",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL of a video.",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The AWS S3 bucket.",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The region of the AWS S3 bucket.",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The path of the AWS S3 object.",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The key ID required to access the AWS S3 object.",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The key required to access the AWS S3 object.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"std_ext_info": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "An extended field for transcoding.Note: This field may return null, indicating that no valid values can be obtained.",
																		},
																		"add_on_subtitles": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The subtitle file to add.Note: This field may return null, indicating that no valid values can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The mode. Valid values:`subtitle-stream`: Add a subtitle track.`close-caption-708`: Embed EA-708 subtitles in SEI frames.`close-caption-608`: Embed CEA-608 subtitles in SEI frames.Note: This field may return null, indicating that no valid values can be obtained.",
																					},
																					"subtitle": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The subtitle file.Note: This field may return null, indicating that no valid values can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The region of the COS bucket, such as `ap-chongqing`.",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "URL of a video.",
																											},
																										},
																									},
																								},
																								"s3_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"s3_bucket": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The AWS S3 bucket.",
																											},
																											"s3_region": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The region of the AWS S3 bucket.",
																											},
																											"s3_object": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The path of the AWS S3 object.",
																											},
																											"s3_secret_id": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The key ID required to access the AWS S3 object.",
																											},
																											"s3_secret_key": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The key required to access the AWS S3 object.",
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
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of up to 10 image or text watermarks.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"definition": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID of a watermarking template.",
																		},
																		"raw_parameter": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Custom watermark parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the watermark parameter preferably.Custom watermark parameter is not available for screenshot.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark type. Valid values:image: image watermark.",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text.Default value: TopLeft.",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image watermark template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input content of watermark image. JPEG and PNG images are supported.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The region of the COS bucket, such as `ap-chongqing`.",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL of a video.",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The AWS S3 bucket.",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The region of the AWS S3 bucket.",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The path of the AWS S3 object.",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The key ID required to access the AWS S3 object.",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The key required to access the AWS S3 object.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"width": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark width. % and px formats are supported:If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark height. % and px formats are supported:If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px.Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat type of an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"text_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Text content of up to 100 characters. This field is required only when the watermark type is text.Text watermark is not available for screenshot.",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG content of up to 2,000,000 characters. This field is required only when the watermark type is `SVG`.SVG watermark is not available for screenshot.",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Start time offset of a watermark in seconds. If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame.If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame;If this value is greater than 0 (e.g., n), the watermark will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the watermark will appear at second n before the last video frame.",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "End time offset of a watermark in seconds.If this parameter is left empty or 0 is entered, the watermark will exist till the last video frame;If this value is greater than 0 (e.g., n), the watermark will exist till second n;If this value is smaller than 0 (e.g., -n), the watermark will exist till second n before the last video frame.",
																		},
																	},
																},
															},
															"mosaic_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of blurs. Up to 10 ones can be supported.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the blur is in the top-left corner of the image or text.Default value: TopLeft.",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The horizontal position of the origin of the blur relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the blur will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the blur will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Vertical position of the origin of blur relative to the origin of coordinates of video. % and px formats are supported:If the string ends in %, the `YPos` of the blur will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the blur will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.",
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Blur width. % and px formats are supported:If the string ends in %, the `Width` of the blur will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the blur will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Blur height. % and px formats are supported:If the string ends in %, the `Height` of the blur will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the blur will be in px; for example, `100px` means that `Height` is 100 px.Default value: 10%.",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Start time offset of blur in seconds. If this parameter is left empty or 0 is entered, the blur will appear upon the first video frame.If this parameter is left empty or 0 is entered, the blur will appear upon the first video frame;If this value is greater than 0 (e.g., n), the blur will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the blur will appear at second n before the last video frame.",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "End time offset of blur in seconds.If this parameter is left empty or 0 is entered, the blur will exist till the last video frame;If this value is greater than 0 (e.g., n), the blur will exist till second n;If this value is smaller than 0 (e.g., -n), the blur will exist till second n before the last video frame.",
																		},
																	},
																},
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Start time offset of a transcoded video, in seconds.If this parameter is left empty or set to 0, the transcoded video will start at the same time as the original video.If this parameter is set to a positive number (n for example), the transcoded video will start at the nth second of the original video.If this parameter is set to a negative number (-n for example), the transcoded video will start at the nth second before the end of the original video.",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "End time offset of a transcoded video, in seconds.If this parameter is left empty or set to 0, the transcoded video will end at the same time as the original video.If this parameter is set to a positive number (n for example), the transcoded video will end at the nth second of the original video.If this parameter is set to a negative number (-n for example), the transcoded video will end at the nth second before the end of the original video.",
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target bucket of an output file. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 bucket.",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the AWS S3 bucket.",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key ID required to upload files to the AWS S3 object.",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key required to upload files to the AWS S3 object.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Path to a primary output file, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_transcode_{definition}.{format}`.",
															},
															"segment_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Path to an output file part (the path to ts during transcoding to HLS), which can only be a relative path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_transcode_{definition}_{number}.{format}`.",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule of the `{number}` variable in the output path after transcoding.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start value of the `{number}` variable. Default value: 0.",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment of the `{number}` variable. Default value: 1.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. Default value: 1.",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. Default value: 0.",
																		},
																	},
																},
															},
															"head_tail_parameter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Opening and closing credits parametersNote: this field may return `null`, indicating that no valid value was found.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"head_set": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Opening credits list.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The region of the COS bucket, such as `ap-chongqing`.",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL of a video.",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The AWS S3 bucket.",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The region of the AWS S3 bucket.",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The path of the AWS S3 object.",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The key ID required to access the AWS S3 object.",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The key required to access the AWS S3 object.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"tail_set": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Closing credits list.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The region of the COS bucket, such as `ap-chongqing`.",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL of a video.",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The AWS S3 bucket.",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The region of the AWS S3 bucket.",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The path of the AWS S3 object.",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The key ID required to access the AWS S3 object.",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The key required to access the AWS S3 object.",
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
												"animated_graphic_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "An animated screenshot generation task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Animated image generating template ID.",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Start time of an animated image in a video in seconds.",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "End time of an animated image in a video in seconds.",
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target bucket of a generated animated image file. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 bucket.",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the AWS S3 bucket.",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key ID required to upload files to the AWS S3 object.",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key required to upload files to the AWS S3 object.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output path to a generated animated image file, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_animatedGraphic_{definition}.{format}`.",
															},
														},
													},
												},
												"snapshot_by_time_offset_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A time point screencapturing task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID of a time point screencapturing template.",
															},
															"ext_time_offset_set": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "List of screenshot time points in the format of `s` or `%`:If the string ends in `s`, it means that the time point is in seconds; for example, `3.5s` means that the time point is the 3.5th second;If the string ends in `%`, it means that the time point is the specified percentage of the video duration; for example, `10%` means that the time point is 10% of the video duration.",
															},
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of up to 10 image or text watermarks.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"definition": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID of a watermarking template.",
																		},
																		"raw_parameter": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Custom watermark parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the watermark parameter preferably.Custom watermark parameter is not available for screenshot.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark type. Valid values:image: image watermark.",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text.Default value: TopLeft.",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image watermark template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input content of watermark image. JPEG and PNG images are supported.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The region of the COS bucket, such as `ap-chongqing`.",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL of a video.",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The AWS S3 bucket.",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The region of the AWS S3 bucket.",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The path of the AWS S3 object.",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The key ID required to access the AWS S3 object.",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The key required to access the AWS S3 object.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"width": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark width. % and px formats are supported:If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark height. % and px formats are supported:If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px.Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat type of an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"text_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Text content of up to 100 characters. This field is required only when the watermark type is text.Text watermark is not available for screenshot.",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG content of up to 2,000,000 characters. This field is required only when the watermark type is `SVG`.SVG watermark is not available for screenshot.",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Start time offset of a watermark in seconds. If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame.If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame;If this value is greater than 0 (e.g., n), the watermark will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the watermark will appear at second n before the last video frame.",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "End time offset of a watermark in seconds.If this parameter is left empty or 0 is entered, the watermark will exist till the last video frame;If this value is greater than 0 (e.g., n), the watermark will exist till second n;If this value is smaller than 0 (e.g., -n), the watermark will exist till second n before the last video frame.",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target bucket of a generated time point screenshot file. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 bucket.",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the AWS S3 bucket.",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key ID required to upload files to the AWS S3 object.",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key required to upload files to the AWS S3 object.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output path to a generated time point screenshot, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_snapshotByTimeOffset_{definition}_{number}.{format}`.",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule of the `{number}` variable in the time point screenshot output path.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start value of the `{number}` variable. Default value: 0.",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment of the `{number}` variable. Default value: 1.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. Default value: 1.",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. Default value: 0.",
																		},
																	},
																},
															},
														},
													},
												},
												"sample_snapshot_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A sampled screencapturing task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Sampled screencapturing template ID.",
															},
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of up to 10 image or text watermarks.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"definition": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID of a watermarking template.",
																		},
																		"raw_parameter": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Custom watermark parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the watermark parameter preferably.Custom watermark parameter is not available for screenshot.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark type. Valid values:image: image watermark.",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text.Default value: TopLeft.",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image watermark template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input content of watermark image. JPEG and PNG images are supported.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The region of the COS bucket, such as `ap-chongqing`.",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL of a video.",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The AWS S3 bucket.",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The region of the AWS S3 bucket.",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The path of the AWS S3 object.",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The key ID required to access the AWS S3 object.",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The key required to access the AWS S3 object.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"width": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark width. % and px formats are supported:If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark height. % and px formats are supported:If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px.Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat type of an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"text_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Text content of up to 100 characters. This field is required only when the watermark type is text.Text watermark is not available for screenshot.",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG content of up to 2,000,000 characters. This field is required only when the watermark type is `SVG`.SVG watermark is not available for screenshot.",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Start time offset of a watermark in seconds. If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame.If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame;If this value is greater than 0 (e.g., n), the watermark will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the watermark will appear at second n before the last video frame.",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "End time offset of a watermark in seconds.If this parameter is left empty or 0 is entered, the watermark will exist till the last video frame;If this value is greater than 0 (e.g., n), the watermark will exist till second n;If this value is smaller than 0 (e.g., -n), the watermark will exist till second n before the last video frame.",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target bucket of a sampled screenshot. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 bucket.",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the AWS S3 bucket.",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key ID required to upload files to the AWS S3 object.",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key required to upload files to the AWS S3 object.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output path to a generated sampled screenshot, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_sampleSnapshot_{definition}_{number}.{format}`.",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule of the `{number}` variable in the sampled screenshot output path.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start value of the `{number}` variable. Default value: 0.",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment of the `{number}` variable. Default value: 1.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. Default value: 1.",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. Default value: 0.",
																		},
																	},
																},
															},
														},
													},
												},
												"image_sprite_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "An image sprite generation task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID of an image sprite generating template.",
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target bucket of a generated image sprite. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 bucket.",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the AWS S3 bucket.",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key ID required to upload files to the AWS S3 object.",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key required to upload files to the AWS S3 object.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output path to a generated image sprite file, which can be a relative path or an absolute path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_imageSprite_{definition}_{number}.{format}`.",
															},
															"web_vtt_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Output path to the WebVTT file after an image sprite is generated, which can only be a relative path. If this parameter is left empty, the following relative path will be used by default: `{inputName}_imageSprite_{definition}.{format}`.",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule of the `{number}` variable in the image sprite output path.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start value of the `{number}` variable. Default value: 0.",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment of the `{number}` variable. Default value: 1.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. Default value: 1.",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. Default value: 0.",
																		},
																	},
																},
															},
														},
													},
												},
												"adaptive_dynamic_streaming_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "An adaptive bitrate streaming task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Adaptive bitrate streaming template ID.",
															},
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of up to 10 image or text watermarks.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"definition": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID of a watermarking template.",
																		},
																		"raw_parameter": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Custom watermark parameter, which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to specify the watermark parameter preferably.Custom watermark parameter is not available for screenshot.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark type. Valid values:image: image watermark.",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position, which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text.Default value: TopLeft.",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width;If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height;If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image watermark template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input content of watermark image. JPEG and PNG images are supported.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The region of the COS bucket, such as `ap-chongqing`.",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL of a video.",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The AWS S3 bucket.",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The region of the AWS S3 bucket.",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The path of the AWS S3 object.",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The key ID required to access the AWS S3 object.",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The key required to access the AWS S3 object.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"width": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark width. % and px formats are supported:If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width; for example, `10%` means that `Width` is 10% of the video width;If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px.Default value: 10%.",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark height. % and px formats are supported:If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px.Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat type of an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"text_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Text content of up to 100 characters. This field is required only when the watermark type is text.Text watermark is not available for screenshot.",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG content of up to 2,000,000 characters. This field is required only when the watermark type is `SVG`.SVG watermark is not available for screenshot.",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Start time offset of a watermark in seconds. If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame.If this parameter is left empty or 0 is entered, the watermark will appear upon the first video frame;If this value is greater than 0 (e.g., n), the watermark will appear at second n after the first video frame;If this value is smaller than 0 (e.g., -n), the watermark will appear at second n before the last video frame.",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "End time offset of a watermark in seconds.If this parameter is left empty or 0 is entered, the watermark will exist till the last video frame;If this value is greater than 0 (e.g., n), the watermark will exist till second n;If this value is smaller than 0 (e.g., -n), the watermark will exist till second n before the last video frame.",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target bucket of an output file after being transcoded to adaptive bitrate streaming. If this parameter is left empty, the `OutputStorage` value of the upper folder will be inherited.Note: this field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 bucket.",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The region of the AWS S3 bucket.",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key ID required to upload files to the AWS S3 object.",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The key required to upload files to the AWS S3 object.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative or absolute output path of the manifest file after being transcoded to adaptive bitrate streaming. If this parameter is left empty, a relative path in the following format will be used by default: `{inputName}_adaptiveDynamicStreaming_{definition}.{format}`.",
															},
															"sub_stream_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative output path of the substream file after being transcoded to adaptive bitrate streaming. If this parameter is left empty, a relative path in the following format will be used by default: `{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{format}`.",
															},
															"segment_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative output path of the segment file after being transcoded to adaptive bitrate streaming (in HLS format only). If this parameter is left empty, a relative path in the following format will be used by default: `{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{format}`.",
															},
															"add_on_subtitles": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The subtitle file to add.Note: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The mode. Valid values:`subtitle-stream`: Add a subtitle track.`close-caption-708`: Embed EA-708 subtitles in SEI frames.`close-caption-608`: Embed CEA-608 subtitles in SEI frames.Note: This field may return null, indicating that no valid values can be obtained.",
																		},
																		"subtitle": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The subtitle file.Note: This field may return null, indicating that no valid values can be obtained.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The input type. Valid values:`COS`: A COS bucket address. `URL`: A URL. `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The region of the COS bucket, such as `ap-chongqing`.",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL of a video.",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The AWS S3 bucket.",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The region of the AWS S3 bucket.",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The path of the AWS S3 object.",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The key ID required to access the AWS S3 object.",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The key required to access the AWS S3 object.",
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
												"ai_content_review_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A content moderation task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Video content audit template ID.",
															},
														},
													},
												},
												"ai_analysis_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A content analysis task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Video content analysis template ID.",
															},
															"extended_parameter": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "An extended parameter, whose value is a stringfied JSON.Note: This parameter is for customers with special requirements. It needs to be customized offline.Note: This field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
												"ai_recognition_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A content recognition task.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Intelligent video recognition template ID.",
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
						"output_storage": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The bucket to save the output file.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The storage type for a media processing output file. Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.",
									},
									"cos_output_storage": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.",
												},
											},
										},
									},
									"s3_output_storage": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"s3_bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The AWS S3 bucket.",
												},
												"s3_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The region of the AWS S3 bucket.",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The key ID required to upload files to the AWS S3 object.",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The key required to upload files to the AWS S3 object.",
												},
											},
										},
									},
								},
							},
						},
						"output_dir": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The directory to save the output file.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"task_notify_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The notification configuration.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cmq_model": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CMQ or TDMQ-CMQ model. Valid values: Queue, Topic.",
									},
									"cmq_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CMQ or TDMQ-CMQ region, such as `sh` (Shanghai) or `bj` (Beijing).",
									},
									"topic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CMQ or TDMQ-CMQ topic to receive notifications. This parameter is valid when `CmqModel` is `Topic`.",
									},
									"queue_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CMQ or TDMQ-CMQ queue to receive notifications. This parameter is valid when `CmqModel` is `Queue`.",
									},
									"notify_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workflow notification method. Valid values: Finish, Change. If this parameter is left empty, `Finish` will be used.",
									},
									"notify_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The notification type. Valid values:`CMQ`: This value is no longer used. Please use `TDMQ-CMQ` instead.`TDMQ-CMQ`: Message queue`URL`: If `NotifyType` is set to `URL`, HTTP callbacks are sent to the URL specified by `NotifyUrl`. HTTP and JSON are used for the callbacks. The packet contains the response parameters of the `ParseNotification` API.`SCF`: This notification type is not recommended. You need to configure it in the SCF console.`AWS-SQS`: AWS queue. This type is only supported for AWS tasks, and the queue must be in the same region as the AWS bucket.Note: If you do not pass this parameter or pass in an empty string, `CMQ` will be used. To use a different notification type, specify this parameter accordingly.",
									},
									"notify_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP callback URL, required if `NotifyType` is set to `URL`.",
									},
									"aws_sqs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The AWS SQS queue. This parameter is required if `NotifyType` is `AWS-SQS`.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sqs_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The region of the SQS queue.",
												},
												"sqs_queue_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the SQS queue.",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The key ID required to read from/write to the SQS queue.",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The key required to read from/write to the SQS queue.",
												},
											},
										},
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time in [ISO date format](https://intl.cloud.tencent.com/document/product/862/37710?from_cn_redirect=1#52).Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last updated time in [ISO date format](https://intl.cloud.tencent.com/document/product/862/37710?from_cn_redirect=1#52).Note: This field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMpsSchedulesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_schedules.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("schedule_ids"); ok {
		scheduleIdsSet := v.(*schema.Set).List()
		scheduleIdList := []interface{}{}
		for i := range scheduleIdsSet {
			scheduleIds := scheduleIdsSet[i].(int)
			scheduleIdList = append(scheduleIdList, scheduleIds)
		}
		paramMap["ScheduleIds"] = scheduleIdList
	}

	if v, ok := d.GetOk("trigger_type"); ok {
		paramMap["TriggerType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var scheduleInfoSet []*mps.SchedulesInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsSchedulesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		scheduleInfoSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(scheduleInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(scheduleInfoSet))

	if scheduleInfoSet != nil {
		for _, schedulesInfo := range scheduleInfoSet {
			schedulesInfoMap := map[string]interface{}{}

			if schedulesInfo.ScheduleId != nil {
				schedulesInfoMap["schedule_id"] = schedulesInfo.ScheduleId
			}

			if schedulesInfo.ScheduleName != nil {
				schedulesInfoMap["schedule_name"] = schedulesInfo.ScheduleName
			}

			if schedulesInfo.Status != nil {
				schedulesInfoMap["status"] = schedulesInfo.Status
			}

			if schedulesInfo.Trigger != nil {
				triggerMap := map[string]interface{}{}

				if schedulesInfo.Trigger.Type != nil {
					triggerMap["type"] = schedulesInfo.Trigger.Type
				}

				if schedulesInfo.Trigger.CosFileUploadTrigger != nil {
					cosFileUploadTriggerMap := map[string]interface{}{}

					if schedulesInfo.Trigger.CosFileUploadTrigger.Bucket != nil {
						cosFileUploadTriggerMap["bucket"] = schedulesInfo.Trigger.CosFileUploadTrigger.Bucket
					}

					if schedulesInfo.Trigger.CosFileUploadTrigger.Region != nil {
						cosFileUploadTriggerMap["region"] = schedulesInfo.Trigger.CosFileUploadTrigger.Region
					}

					if schedulesInfo.Trigger.CosFileUploadTrigger.Dir != nil {
						cosFileUploadTriggerMap["dir"] = schedulesInfo.Trigger.CosFileUploadTrigger.Dir
					}

					if schedulesInfo.Trigger.CosFileUploadTrigger.Formats != nil {
						cosFileUploadTriggerMap["formats"] = schedulesInfo.Trigger.CosFileUploadTrigger.Formats
					}

					triggerMap["cos_file_upload_trigger"] = []interface{}{cosFileUploadTriggerMap}
				}

				if schedulesInfo.Trigger.AwsS3FileUploadTrigger != nil {
					awsS3FileUploadTriggerMap := map[string]interface{}{}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3Bucket != nil {
						awsS3FileUploadTriggerMap["s3_bucket"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3Bucket
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3Region != nil {
						awsS3FileUploadTriggerMap["s3_region"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3Region
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.Dir != nil {
						awsS3FileUploadTriggerMap["dir"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.Dir
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.Formats != nil {
						awsS3FileUploadTriggerMap["formats"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.Formats
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3SecretId != nil {
						awsS3FileUploadTriggerMap["s3_secret_id"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3SecretId
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3SecretKey != nil {
						awsS3FileUploadTriggerMap["s3_secret_key"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3SecretKey
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS != nil {
						awsSQSMap := map[string]interface{}{}

						if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSRegion != nil {
							awsSQSMap["sqs_region"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSRegion
						}

						if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSQueueName != nil {
							awsSQSMap["sqs_queue_name"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSQueueName
						}

						if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretId != nil {
							awsSQSMap["s3_secret_id"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretId
						}

						if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretKey != nil {
							awsSQSMap["s3_secret_key"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretKey
						}

						awsS3FileUploadTriggerMap["aws_sqs"] = []interface{}{awsSQSMap}
					}

					triggerMap["aws_s3_file_upload_trigger"] = []interface{}{awsS3FileUploadTriggerMap}
				}

				schedulesInfoMap["trigger"] = []interface{}{triggerMap}
			}

			if schedulesInfo.Activities != nil {
				activitiesList := []interface{}{}
				for _, activities := range schedulesInfo.Activities {
					activitiesMap := map[string]interface{}{}

					if activities.ActivityType != nil {
						activitiesMap["activity_type"] = activities.ActivityType
					}

					if activities.ReardriveIndex != nil {
						activitiesMap["reardrive_index"] = activities.ReardriveIndex
					}

					if activities.ActivityPara != nil {
						activityParaMap := map[string]interface{}{}

						if activities.ActivityPara.TranscodeTask != nil {
							transcodeTaskMap := map[string]interface{}{}

							if activities.ActivityPara.TranscodeTask.Definition != nil {
								transcodeTaskMap["definition"] = activities.ActivityPara.TranscodeTask.Definition
							}

							if activities.ActivityPara.TranscodeTask.RawParameter != nil {
								rawParameterMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.RawParameter.Container != nil {
									rawParameterMap["container"] = activities.ActivityPara.TranscodeTask.RawParameter.Container
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.RemoveVideo != nil {
									rawParameterMap["remove_video"] = activities.ActivityPara.TranscodeTask.RawParameter.RemoveVideo
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.RemoveAudio != nil {
									rawParameterMap["remove_audio"] = activities.ActivityPara.TranscodeTask.RawParameter.RemoveAudio
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate != nil {
									videoTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Codec != nil {
										videoTemplateMap["codec"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Codec
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Fps != nil {
										videoTemplateMap["fps"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Fps
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Bitrate != nil {
										videoTemplateMap["bitrate"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Bitrate
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.ResolutionAdaptive != nil {
										videoTemplateMap["resolution_adaptive"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.ResolutionAdaptive
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Width != nil {
										videoTemplateMap["width"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Width
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Height != nil {
										videoTemplateMap["height"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Height
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Gop != nil {
										videoTemplateMap["gop"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Gop
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.FillType != nil {
										videoTemplateMap["fill_type"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.FillType
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Vcrf != nil {
										videoTemplateMap["vcrf"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Vcrf
									}

									rawParameterMap["video_template"] = []interface{}{videoTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate != nil {
									audioTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Codec != nil {
										audioTemplateMap["codec"] = activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Codec
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Bitrate != nil {
										audioTemplateMap["bitrate"] = activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Bitrate
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.SampleRate != nil {
										audioTemplateMap["sample_rate"] = activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.SampleRate
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.AudioChannel != nil {
										audioTemplateMap["audio_channel"] = activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.AudioChannel
									}

									rawParameterMap["audio_template"] = []interface{}{audioTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig != nil {
									tEHDConfigMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.Type != nil {
										tEHDConfigMap["type"] = activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.Type
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.MaxVideoBitrate != nil {
										tEHDConfigMap["max_video_bitrate"] = activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.MaxVideoBitrate
									}

									rawParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
								}

								transcodeTaskMap["raw_parameter"] = []interface{}{rawParameterMap}
							}

							if activities.ActivityPara.TranscodeTask.OverrideParameter != nil {
								overrideParameterMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.Container != nil {
									overrideParameterMap["container"] = activities.ActivityPara.TranscodeTask.OverrideParameter.Container
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.RemoveVideo != nil {
									overrideParameterMap["remove_video"] = activities.ActivityPara.TranscodeTask.OverrideParameter.RemoveVideo
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.RemoveAudio != nil {
									overrideParameterMap["remove_audio"] = activities.ActivityPara.TranscodeTask.OverrideParameter.RemoveAudio
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate != nil {
									videoTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Codec != nil {
										videoTemplateMap["codec"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Codec
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Fps != nil {
										videoTemplateMap["fps"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Fps
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Bitrate != nil {
										videoTemplateMap["bitrate"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Bitrate
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ResolutionAdaptive != nil {
										videoTemplateMap["resolution_adaptive"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ResolutionAdaptive
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Width != nil {
										videoTemplateMap["width"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Width
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Height != nil {
										videoTemplateMap["height"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Height
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Gop != nil {
										videoTemplateMap["gop"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Gop
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.FillType != nil {
										videoTemplateMap["fill_type"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.FillType
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Vcrf != nil {
										videoTemplateMap["vcrf"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Vcrf
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ContentAdaptStream != nil {
										videoTemplateMap["content_adapt_stream"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ContentAdaptStream
									}

									overrideParameterMap["video_template"] = []interface{}{videoTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate != nil {
									audioTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Codec != nil {
										audioTemplateMap["codec"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Codec
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Bitrate != nil {
										audioTemplateMap["bitrate"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Bitrate
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.SampleRate != nil {
										audioTemplateMap["sample_rate"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.SampleRate
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.AudioChannel != nil {
										audioTemplateMap["audio_channel"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.AudioChannel
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.StreamSelects != nil {
										audioTemplateMap["stream_selects"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.StreamSelects
									}

									overrideParameterMap["audio_template"] = []interface{}{audioTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig != nil {
									tEHDConfigMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.Type != nil {
										tEHDConfigMap["type"] = activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.Type
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.MaxVideoBitrate != nil {
										tEHDConfigMap["max_video_bitrate"] = activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.MaxVideoBitrate
									}

									overrideParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate != nil {
									subtitleTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.Path != nil {
										subtitleTemplateMap["path"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.Path
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.StreamIndex != nil {
										subtitleTemplateMap["stream_index"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.StreamIndex
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontType != nil {
										subtitleTemplateMap["font_type"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontType
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontSize != nil {
										subtitleTemplateMap["font_size"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontSize
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontColor != nil {
										subtitleTemplateMap["font_color"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontColor
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontAlpha != nil {
										subtitleTemplateMap["font_alpha"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontAlpha
									}

									overrideParameterMap["subtitle_template"] = []interface{}{subtitleTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.AddonAudioStream != nil {
									addonAudioStreamList := []interface{}{}
									for _, addonAudioStream := range activities.ActivityPara.TranscodeTask.OverrideParameter.AddonAudioStream {
										addonAudioStreamMap := map[string]interface{}{}

										if addonAudioStream.Type != nil {
											addonAudioStreamMap["type"] = addonAudioStream.Type
										}

										if addonAudioStream.CosInputInfo != nil {
											cosInputInfoMap := map[string]interface{}{}

											if addonAudioStream.CosInputInfo.Bucket != nil {
												cosInputInfoMap["bucket"] = addonAudioStream.CosInputInfo.Bucket
											}

											if addonAudioStream.CosInputInfo.Region != nil {
												cosInputInfoMap["region"] = addonAudioStream.CosInputInfo.Region
											}

											if addonAudioStream.CosInputInfo.Object != nil {
												cosInputInfoMap["object"] = addonAudioStream.CosInputInfo.Object
											}

											addonAudioStreamMap["cos_input_info"] = []interface{}{cosInputInfoMap}
										}

										if addonAudioStream.UrlInputInfo != nil {
											urlInputInfoMap := map[string]interface{}{}

											if addonAudioStream.UrlInputInfo.Url != nil {
												urlInputInfoMap["url"] = addonAudioStream.UrlInputInfo.Url
											}

											addonAudioStreamMap["url_input_info"] = []interface{}{urlInputInfoMap}
										}

										if addonAudioStream.S3InputInfo != nil {
											s3InputInfoMap := map[string]interface{}{}

											if addonAudioStream.S3InputInfo.S3Bucket != nil {
												s3InputInfoMap["s3_bucket"] = addonAudioStream.S3InputInfo.S3Bucket
											}

											if addonAudioStream.S3InputInfo.S3Region != nil {
												s3InputInfoMap["s3_region"] = addonAudioStream.S3InputInfo.S3Region
											}

											if addonAudioStream.S3InputInfo.S3Object != nil {
												s3InputInfoMap["s3_object"] = addonAudioStream.S3InputInfo.S3Object
											}

											if addonAudioStream.S3InputInfo.S3SecretId != nil {
												s3InputInfoMap["s3_secret_id"] = addonAudioStream.S3InputInfo.S3SecretId
											}

											if addonAudioStream.S3InputInfo.S3SecretKey != nil {
												s3InputInfoMap["s3_secret_key"] = addonAudioStream.S3InputInfo.S3SecretKey
											}

											addonAudioStreamMap["s3_input_info"] = []interface{}{s3InputInfoMap}
										}

										addonAudioStreamList = append(addonAudioStreamList, addonAudioStreamMap)
									}

									overrideParameterMap["addon_audio_stream"] = addonAudioStreamList
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.StdExtInfo != nil {
									overrideParameterMap["std_ext_info"] = activities.ActivityPara.TranscodeTask.OverrideParameter.StdExtInfo
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.AddOnSubtitles != nil {
									addOnSubtitlesList := []interface{}{}
									for _, addOnSubtitles := range activities.ActivityPara.TranscodeTask.OverrideParameter.AddOnSubtitles {
										addOnSubtitlesMap := map[string]interface{}{}

										if addOnSubtitles.Type != nil {
											addOnSubtitlesMap["type"] = addOnSubtitles.Type
										}

										if addOnSubtitles.Subtitle != nil {
											subtitleMap := map[string]interface{}{}

											if addOnSubtitles.Subtitle.Type != nil {
												subtitleMap["type"] = addOnSubtitles.Subtitle.Type
											}

											if addOnSubtitles.Subtitle.CosInputInfo != nil {
												cosInputInfoMap := map[string]interface{}{}

												if addOnSubtitles.Subtitle.CosInputInfo.Bucket != nil {
													cosInputInfoMap["bucket"] = addOnSubtitles.Subtitle.CosInputInfo.Bucket
												}

												if addOnSubtitles.Subtitle.CosInputInfo.Region != nil {
													cosInputInfoMap["region"] = addOnSubtitles.Subtitle.CosInputInfo.Region
												}

												if addOnSubtitles.Subtitle.CosInputInfo.Object != nil {
													cosInputInfoMap["object"] = addOnSubtitles.Subtitle.CosInputInfo.Object
												}

												subtitleMap["cos_input_info"] = []interface{}{cosInputInfoMap}
											}

											if addOnSubtitles.Subtitle.UrlInputInfo != nil {
												urlInputInfoMap := map[string]interface{}{}

												if addOnSubtitles.Subtitle.UrlInputInfo.Url != nil {
													urlInputInfoMap["url"] = addOnSubtitles.Subtitle.UrlInputInfo.Url
												}

												subtitleMap["url_input_info"] = []interface{}{urlInputInfoMap}
											}

											if addOnSubtitles.Subtitle.S3InputInfo != nil {
												s3InputInfoMap := map[string]interface{}{}

												if addOnSubtitles.Subtitle.S3InputInfo.S3Bucket != nil {
													s3InputInfoMap["s3_bucket"] = addOnSubtitles.Subtitle.S3InputInfo.S3Bucket
												}

												if addOnSubtitles.Subtitle.S3InputInfo.S3Region != nil {
													s3InputInfoMap["s3_region"] = addOnSubtitles.Subtitle.S3InputInfo.S3Region
												}

												if addOnSubtitles.Subtitle.S3InputInfo.S3Object != nil {
													s3InputInfoMap["s3_object"] = addOnSubtitles.Subtitle.S3InputInfo.S3Object
												}

												if addOnSubtitles.Subtitle.S3InputInfo.S3SecretId != nil {
													s3InputInfoMap["s3_secret_id"] = addOnSubtitles.Subtitle.S3InputInfo.S3SecretId
												}

												if addOnSubtitles.Subtitle.S3InputInfo.S3SecretKey != nil {
													s3InputInfoMap["s3_secret_key"] = addOnSubtitles.Subtitle.S3InputInfo.S3SecretKey
												}

												subtitleMap["s3_input_info"] = []interface{}{s3InputInfoMap}
											}

											addOnSubtitlesMap["subtitle"] = []interface{}{subtitleMap}
										}

										addOnSubtitlesList = append(addOnSubtitlesList, addOnSubtitlesMap)
									}

									overrideParameterMap["add_on_subtitles"] = addOnSubtitlesList
								}

								transcodeTaskMap["override_parameter"] = []interface{}{overrideParameterMap}
							}

							if activities.ActivityPara.TranscodeTask.WatermarkSet != nil {
								watermarkSetList := []interface{}{}
								for _, watermarkSet := range activities.ActivityPara.TranscodeTask.WatermarkSet {
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

												if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo != nil {
													s3InputInfoMap := map[string]interface{}{}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket != nil {
														s3InputInfoMap["s3_bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region != nil {
														s3InputInfoMap["s3_region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object != nil {
														s3InputInfoMap["s3_object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId != nil {
														s3InputInfoMap["s3_secret_id"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey != nil {
														s3InputInfoMap["s3_secret_key"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey
													}

													imageContentMap["s3_input_info"] = []interface{}{s3InputInfoMap}
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

								transcodeTaskMap["watermark_set"] = watermarkSetList
							}

							if activities.ActivityPara.TranscodeTask.MosaicSet != nil {
								mosaicSetList := []interface{}{}
								for _, mosaicSet := range activities.ActivityPara.TranscodeTask.MosaicSet {
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

								transcodeTaskMap["mosaic_set"] = mosaicSetList
							}

							if activities.ActivityPara.TranscodeTask.StartTimeOffset != nil {
								transcodeTaskMap["start_time_offset"] = activities.ActivityPara.TranscodeTask.StartTimeOffset
							}

							if activities.ActivityPara.TranscodeTask.EndTimeOffset != nil {
								transcodeTaskMap["end_time_offset"] = activities.ActivityPara.TranscodeTask.EndTimeOffset
							}

							if activities.ActivityPara.TranscodeTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.TranscodeTask.OutputStorage.Type
								}

								if activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								transcodeTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.TranscodeTask.OutputObjectPath != nil {
								transcodeTaskMap["output_object_path"] = activities.ActivityPara.TranscodeTask.OutputObjectPath
							}

							if activities.ActivityPara.TranscodeTask.SegmentObjectName != nil {
								transcodeTaskMap["segment_object_name"] = activities.ActivityPara.TranscodeTask.SegmentObjectName
							}

							if activities.ActivityPara.TranscodeTask.ObjectNumberFormat != nil {
								objectNumberFormatMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.ObjectNumberFormat.InitialValue != nil {
									objectNumberFormatMap["initial_value"] = activities.ActivityPara.TranscodeTask.ObjectNumberFormat.InitialValue
								}

								if activities.ActivityPara.TranscodeTask.ObjectNumberFormat.Increment != nil {
									objectNumberFormatMap["increment"] = activities.ActivityPara.TranscodeTask.ObjectNumberFormat.Increment
								}

								if activities.ActivityPara.TranscodeTask.ObjectNumberFormat.MinLength != nil {
									objectNumberFormatMap["min_length"] = activities.ActivityPara.TranscodeTask.ObjectNumberFormat.MinLength
								}

								if activities.ActivityPara.TranscodeTask.ObjectNumberFormat.PlaceHolder != nil {
									objectNumberFormatMap["place_holder"] = activities.ActivityPara.TranscodeTask.ObjectNumberFormat.PlaceHolder
								}

								transcodeTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
							}

							if activities.ActivityPara.TranscodeTask.HeadTailParameter != nil {
								headTailParameterMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.HeadTailParameter.HeadSet != nil {
									headSetList := []interface{}{}
									for _, headSet := range activities.ActivityPara.TranscodeTask.HeadTailParameter.HeadSet {
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

										if headSet.S3InputInfo != nil {
											s3InputInfoMap := map[string]interface{}{}

											if headSet.S3InputInfo.S3Bucket != nil {
												s3InputInfoMap["s3_bucket"] = headSet.S3InputInfo.S3Bucket
											}

											if headSet.S3InputInfo.S3Region != nil {
												s3InputInfoMap["s3_region"] = headSet.S3InputInfo.S3Region
											}

											if headSet.S3InputInfo.S3Object != nil {
												s3InputInfoMap["s3_object"] = headSet.S3InputInfo.S3Object
											}

											if headSet.S3InputInfo.S3SecretId != nil {
												s3InputInfoMap["s3_secret_id"] = headSet.S3InputInfo.S3SecretId
											}

											if headSet.S3InputInfo.S3SecretKey != nil {
												s3InputInfoMap["s3_secret_key"] = headSet.S3InputInfo.S3SecretKey
											}

											headSetMap["s3_input_info"] = []interface{}{s3InputInfoMap}
										}

										headSetList = append(headSetList, headSetMap)
									}

									headTailParameterMap["head_set"] = headSetList
								}

								if activities.ActivityPara.TranscodeTask.HeadTailParameter.TailSet != nil {
									tailSetList := []interface{}{}
									for _, tailSet := range activities.ActivityPara.TranscodeTask.HeadTailParameter.TailSet {
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

										if tailSet.S3InputInfo != nil {
											s3InputInfoMap := map[string]interface{}{}

											if tailSet.S3InputInfo.S3Bucket != nil {
												s3InputInfoMap["s3_bucket"] = tailSet.S3InputInfo.S3Bucket
											}

											if tailSet.S3InputInfo.S3Region != nil {
												s3InputInfoMap["s3_region"] = tailSet.S3InputInfo.S3Region
											}

											if tailSet.S3InputInfo.S3Object != nil {
												s3InputInfoMap["s3_object"] = tailSet.S3InputInfo.S3Object
											}

											if tailSet.S3InputInfo.S3SecretId != nil {
												s3InputInfoMap["s3_secret_id"] = tailSet.S3InputInfo.S3SecretId
											}

											if tailSet.S3InputInfo.S3SecretKey != nil {
												s3InputInfoMap["s3_secret_key"] = tailSet.S3InputInfo.S3SecretKey
											}

											tailSetMap["s3_input_info"] = []interface{}{s3InputInfoMap}
										}

										tailSetList = append(tailSetList, tailSetMap)
									}

									headTailParameterMap["tail_set"] = tailSetList
								}

								transcodeTaskMap["head_tail_parameter"] = []interface{}{headTailParameterMap}
							}

							activityParaMap["transcode_task"] = []interface{}{transcodeTaskMap}
						}

						if activities.ActivityPara.AnimatedGraphicTask != nil {
							animatedGraphicTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AnimatedGraphicTask.Definition != nil {
								animatedGraphicTaskMap["definition"] = activities.ActivityPara.AnimatedGraphicTask.Definition
							}

							if activities.ActivityPara.AnimatedGraphicTask.StartTimeOffset != nil {
								animatedGraphicTaskMap["start_time_offset"] = activities.ActivityPara.AnimatedGraphicTask.StartTimeOffset
							}

							if activities.ActivityPara.AnimatedGraphicTask.EndTimeOffset != nil {
								animatedGraphicTaskMap["end_time_offset"] = activities.ActivityPara.AnimatedGraphicTask.EndTimeOffset
							}

							if activities.ActivityPara.AnimatedGraphicTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.Type
								}

								if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								animatedGraphicTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.AnimatedGraphicTask.OutputObjectPath != nil {
								animatedGraphicTaskMap["output_object_path"] = activities.ActivityPara.AnimatedGraphicTask.OutputObjectPath
							}

							activityParaMap["animated_graphic_task"] = []interface{}{animatedGraphicTaskMap}
						}

						if activities.ActivityPara.SnapshotByTimeOffsetTask != nil {
							snapshotByTimeOffsetTaskMap := map[string]interface{}{}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.Definition != nil {
								snapshotByTimeOffsetTaskMap["definition"] = activities.ActivityPara.SnapshotByTimeOffsetTask.Definition
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.ExtTimeOffsetSet != nil {
								snapshotByTimeOffsetTaskMap["ext_time_offset_set"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ExtTimeOffsetSet
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.WatermarkSet != nil {
								watermarkSetList := []interface{}{}
								for _, watermarkSet := range activities.ActivityPara.SnapshotByTimeOffsetTask.WatermarkSet {
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

												if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo != nil {
													s3InputInfoMap := map[string]interface{}{}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket != nil {
														s3InputInfoMap["s3_bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region != nil {
														s3InputInfoMap["s3_region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object != nil {
														s3InputInfoMap["s3_object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId != nil {
														s3InputInfoMap["s3_secret_id"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey != nil {
														s3InputInfoMap["s3_secret_key"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey
													}

													imageContentMap["s3_input_info"] = []interface{}{s3InputInfoMap}
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

								snapshotByTimeOffsetTaskMap["watermark_set"] = watermarkSetList
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.Type
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								snapshotByTimeOffsetTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputObjectPath != nil {
								snapshotByTimeOffsetTaskMap["output_object_path"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputObjectPath
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat != nil {
								objectNumberFormatMap := map[string]interface{}{}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.InitialValue != nil {
									objectNumberFormatMap["initial_value"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.InitialValue
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.Increment != nil {
									objectNumberFormatMap["increment"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.Increment
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.MinLength != nil {
									objectNumberFormatMap["min_length"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.MinLength
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.PlaceHolder != nil {
									objectNumberFormatMap["place_holder"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.PlaceHolder
								}

								snapshotByTimeOffsetTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
							}

							activityParaMap["snapshot_by_time_offset_task"] = []interface{}{snapshotByTimeOffsetTaskMap}
						}

						if activities.ActivityPara.SampleSnapshotTask != nil {
							sampleSnapshotTaskMap := map[string]interface{}{}

							if activities.ActivityPara.SampleSnapshotTask.Definition != nil {
								sampleSnapshotTaskMap["definition"] = activities.ActivityPara.SampleSnapshotTask.Definition
							}

							if activities.ActivityPara.SampleSnapshotTask.WatermarkSet != nil {
								watermarkSetList := []interface{}{}
								for _, watermarkSet := range activities.ActivityPara.SampleSnapshotTask.WatermarkSet {
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

												if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo != nil {
													s3InputInfoMap := map[string]interface{}{}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket != nil {
														s3InputInfoMap["s3_bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region != nil {
														s3InputInfoMap["s3_region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object != nil {
														s3InputInfoMap["s3_object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId != nil {
														s3InputInfoMap["s3_secret_id"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey != nil {
														s3InputInfoMap["s3_secret_key"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey
													}

													imageContentMap["s3_input_info"] = []interface{}{s3InputInfoMap}
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

								sampleSnapshotTaskMap["watermark_set"] = watermarkSetList
							}

							if activities.ActivityPara.SampleSnapshotTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.SampleSnapshotTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.Type
								}

								if activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								sampleSnapshotTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.SampleSnapshotTask.OutputObjectPath != nil {
								sampleSnapshotTaskMap["output_object_path"] = activities.ActivityPara.SampleSnapshotTask.OutputObjectPath
							}

							if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat != nil {
								objectNumberFormatMap := map[string]interface{}{}

								if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.InitialValue != nil {
									objectNumberFormatMap["initial_value"] = activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.InitialValue
								}

								if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.Increment != nil {
									objectNumberFormatMap["increment"] = activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.Increment
								}

								if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.MinLength != nil {
									objectNumberFormatMap["min_length"] = activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.MinLength
								}

								if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.PlaceHolder != nil {
									objectNumberFormatMap["place_holder"] = activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.PlaceHolder
								}

								sampleSnapshotTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
							}

							activityParaMap["sample_snapshot_task"] = []interface{}{sampleSnapshotTaskMap}
						}

						if activities.ActivityPara.ImageSpriteTask != nil {
							imageSpriteTaskMap := map[string]interface{}{}

							if activities.ActivityPara.ImageSpriteTask.Definition != nil {
								imageSpriteTaskMap["definition"] = activities.ActivityPara.ImageSpriteTask.Definition
							}

							if activities.ActivityPara.ImageSpriteTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.ImageSpriteTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.Type
								}

								if activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								imageSpriteTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.ImageSpriteTask.OutputObjectPath != nil {
								imageSpriteTaskMap["output_object_path"] = activities.ActivityPara.ImageSpriteTask.OutputObjectPath
							}

							if activities.ActivityPara.ImageSpriteTask.WebVttObjectName != nil {
								imageSpriteTaskMap["web_vtt_object_name"] = activities.ActivityPara.ImageSpriteTask.WebVttObjectName
							}

							if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat != nil {
								objectNumberFormatMap := map[string]interface{}{}

								if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.InitialValue != nil {
									objectNumberFormatMap["initial_value"] = activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.InitialValue
								}

								if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.Increment != nil {
									objectNumberFormatMap["increment"] = activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.Increment
								}

								if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.MinLength != nil {
									objectNumberFormatMap["min_length"] = activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.MinLength
								}

								if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.PlaceHolder != nil {
									objectNumberFormatMap["place_holder"] = activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.PlaceHolder
								}

								imageSpriteTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
							}

							activityParaMap["image_sprite_task"] = []interface{}{imageSpriteTaskMap}
						}

						if activities.ActivityPara.AdaptiveDynamicStreamingTask != nil {
							adaptiveDynamicStreamingTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.Definition != nil {
								adaptiveDynamicStreamingTaskMap["definition"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.Definition
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.WatermarkSet != nil {
								watermarkSetList := []interface{}{}
								for _, watermarkSet := range activities.ActivityPara.AdaptiveDynamicStreamingTask.WatermarkSet {
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

												if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo != nil {
													s3InputInfoMap := map[string]interface{}{}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket != nil {
														s3InputInfoMap["s3_bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region != nil {
														s3InputInfoMap["s3_region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object != nil {
														s3InputInfoMap["s3_object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId != nil {
														s3InputInfoMap["s3_secret_id"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey != nil {
														s3InputInfoMap["s3_secret_key"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey
													}

													imageContentMap["s3_input_info"] = []interface{}{s3InputInfoMap}
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

								adaptiveDynamicStreamingTaskMap["watermark_set"] = watermarkSetList
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.Type
								}

								if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								adaptiveDynamicStreamingTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputObjectPath != nil {
								adaptiveDynamicStreamingTaskMap["output_object_path"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputObjectPath
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.SubStreamObjectName != nil {
								adaptiveDynamicStreamingTaskMap["sub_stream_object_name"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.SubStreamObjectName
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.SegmentObjectName != nil {
								adaptiveDynamicStreamingTaskMap["segment_object_name"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.SegmentObjectName
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.AddOnSubtitles != nil {
								addOnSubtitlesList := []interface{}{}
								for _, addOnSubtitles := range activities.ActivityPara.AdaptiveDynamicStreamingTask.AddOnSubtitles {
									addOnSubtitlesMap := map[string]interface{}{}

									if addOnSubtitles.Type != nil {
										addOnSubtitlesMap["type"] = addOnSubtitles.Type
									}

									if addOnSubtitles.Subtitle != nil {
										subtitleMap := map[string]interface{}{}

										if addOnSubtitles.Subtitle.Type != nil {
											subtitleMap["type"] = addOnSubtitles.Subtitle.Type
										}

										if addOnSubtitles.Subtitle.CosInputInfo != nil {
											cosInputInfoMap := map[string]interface{}{}

											if addOnSubtitles.Subtitle.CosInputInfo.Bucket != nil {
												cosInputInfoMap["bucket"] = addOnSubtitles.Subtitle.CosInputInfo.Bucket
											}

											if addOnSubtitles.Subtitle.CosInputInfo.Region != nil {
												cosInputInfoMap["region"] = addOnSubtitles.Subtitle.CosInputInfo.Region
											}

											if addOnSubtitles.Subtitle.CosInputInfo.Object != nil {
												cosInputInfoMap["object"] = addOnSubtitles.Subtitle.CosInputInfo.Object
											}

											subtitleMap["cos_input_info"] = []interface{}{cosInputInfoMap}
										}

										if addOnSubtitles.Subtitle.UrlInputInfo != nil {
											urlInputInfoMap := map[string]interface{}{}

											if addOnSubtitles.Subtitle.UrlInputInfo.Url != nil {
												urlInputInfoMap["url"] = addOnSubtitles.Subtitle.UrlInputInfo.Url
											}

											subtitleMap["url_input_info"] = []interface{}{urlInputInfoMap}
										}

										if addOnSubtitles.Subtitle.S3InputInfo != nil {
											s3InputInfoMap := map[string]interface{}{}

											if addOnSubtitles.Subtitle.S3InputInfo.S3Bucket != nil {
												s3InputInfoMap["s3_bucket"] = addOnSubtitles.Subtitle.S3InputInfo.S3Bucket
											}

											if addOnSubtitles.Subtitle.S3InputInfo.S3Region != nil {
												s3InputInfoMap["s3_region"] = addOnSubtitles.Subtitle.S3InputInfo.S3Region
											}

											if addOnSubtitles.Subtitle.S3InputInfo.S3Object != nil {
												s3InputInfoMap["s3_object"] = addOnSubtitles.Subtitle.S3InputInfo.S3Object
											}

											if addOnSubtitles.Subtitle.S3InputInfo.S3SecretId != nil {
												s3InputInfoMap["s3_secret_id"] = addOnSubtitles.Subtitle.S3InputInfo.S3SecretId
											}

											if addOnSubtitles.Subtitle.S3InputInfo.S3SecretKey != nil {
												s3InputInfoMap["s3_secret_key"] = addOnSubtitles.Subtitle.S3InputInfo.S3SecretKey
											}

											subtitleMap["s3_input_info"] = []interface{}{s3InputInfoMap}
										}

										addOnSubtitlesMap["subtitle"] = []interface{}{subtitleMap}
									}

									addOnSubtitlesList = append(addOnSubtitlesList, addOnSubtitlesMap)
								}

								adaptiveDynamicStreamingTaskMap["add_on_subtitles"] = addOnSubtitlesList
							}

							activityParaMap["adaptive_dynamic_streaming_task"] = []interface{}{adaptiveDynamicStreamingTaskMap}
						}

						if activities.ActivityPara.AiContentReviewTask != nil {
							aiContentReviewTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AiContentReviewTask.Definition != nil {
								aiContentReviewTaskMap["definition"] = activities.ActivityPara.AiContentReviewTask.Definition
							}

							activityParaMap["ai_content_review_task"] = []interface{}{aiContentReviewTaskMap}
						}

						if activities.ActivityPara.AiAnalysisTask != nil {
							aiAnalysisTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AiAnalysisTask.Definition != nil {
								aiAnalysisTaskMap["definition"] = activities.ActivityPara.AiAnalysisTask.Definition
							}

							if activities.ActivityPara.AiAnalysisTask.ExtendedParameter != nil {
								aiAnalysisTaskMap["extended_parameter"] = activities.ActivityPara.AiAnalysisTask.ExtendedParameter
							}

							activityParaMap["ai_analysis_task"] = []interface{}{aiAnalysisTaskMap}
						}

						if activities.ActivityPara.AiRecognitionTask != nil {
							aiRecognitionTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AiRecognitionTask.Definition != nil {
								aiRecognitionTaskMap["definition"] = activities.ActivityPara.AiRecognitionTask.Definition
							}

							activityParaMap["ai_recognition_task"] = []interface{}{aiRecognitionTaskMap}
						}

						activitiesMap["activity_para"] = []interface{}{activityParaMap}
					}

					activitiesList = append(activitiesList, activitiesMap)
				}

				schedulesInfoMap["activities"] = activitiesList
			}

			if schedulesInfo.OutputStorage != nil {
				outputStorageMap := map[string]interface{}{}

				if schedulesInfo.OutputStorage.Type != nil {
					outputStorageMap["type"] = schedulesInfo.OutputStorage.Type
				}

				if schedulesInfo.OutputStorage.CosOutputStorage != nil {
					cosOutputStorageMap := map[string]interface{}{}

					if schedulesInfo.OutputStorage.CosOutputStorage.Bucket != nil {
						cosOutputStorageMap["bucket"] = schedulesInfo.OutputStorage.CosOutputStorage.Bucket
					}

					if schedulesInfo.OutputStorage.CosOutputStorage.Region != nil {
						cosOutputStorageMap["region"] = schedulesInfo.OutputStorage.CosOutputStorage.Region
					}

					outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
				}

				if schedulesInfo.OutputStorage.S3OutputStorage != nil {
					s3OutputStorageMap := map[string]interface{}{}

					if schedulesInfo.OutputStorage.S3OutputStorage.S3Bucket != nil {
						s3OutputStorageMap["s3_bucket"] = schedulesInfo.OutputStorage.S3OutputStorage.S3Bucket
					}

					if schedulesInfo.OutputStorage.S3OutputStorage.S3Region != nil {
						s3OutputStorageMap["s3_region"] = schedulesInfo.OutputStorage.S3OutputStorage.S3Region
					}

					if schedulesInfo.OutputStorage.S3OutputStorage.S3SecretId != nil {
						s3OutputStorageMap["s3_secret_id"] = schedulesInfo.OutputStorage.S3OutputStorage.S3SecretId
					}

					if schedulesInfo.OutputStorage.S3OutputStorage.S3SecretKey != nil {
						s3OutputStorageMap["s3_secret_key"] = schedulesInfo.OutputStorage.S3OutputStorage.S3SecretKey
					}

					outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
				}

				schedulesInfoMap["output_storage"] = []interface{}{outputStorageMap}
			}

			if schedulesInfo.OutputDir != nil {
				schedulesInfoMap["output_dir"] = schedulesInfo.OutputDir
			}

			if schedulesInfo.TaskNotifyConfig != nil {
				taskNotifyConfigMap := map[string]interface{}{}

				if schedulesInfo.TaskNotifyConfig.CmqModel != nil {
					taskNotifyConfigMap["cmq_model"] = schedulesInfo.TaskNotifyConfig.CmqModel
				}

				if schedulesInfo.TaskNotifyConfig.CmqRegion != nil {
					taskNotifyConfigMap["cmq_region"] = schedulesInfo.TaskNotifyConfig.CmqRegion
				}

				if schedulesInfo.TaskNotifyConfig.TopicName != nil {
					taskNotifyConfigMap["topic_name"] = schedulesInfo.TaskNotifyConfig.TopicName
				}

				if schedulesInfo.TaskNotifyConfig.QueueName != nil {
					taskNotifyConfigMap["queue_name"] = schedulesInfo.TaskNotifyConfig.QueueName
				}

				if schedulesInfo.TaskNotifyConfig.NotifyMode != nil {
					taskNotifyConfigMap["notify_mode"] = schedulesInfo.TaskNotifyConfig.NotifyMode
				}

				if schedulesInfo.TaskNotifyConfig.NotifyType != nil {
					taskNotifyConfigMap["notify_type"] = schedulesInfo.TaskNotifyConfig.NotifyType
				}

				if schedulesInfo.TaskNotifyConfig.NotifyUrl != nil {
					taskNotifyConfigMap["notify_url"] = schedulesInfo.TaskNotifyConfig.NotifyUrl
				}

				if schedulesInfo.TaskNotifyConfig.AwsSQS != nil {
					awsSQSMap := map[string]interface{}{}

					if schedulesInfo.TaskNotifyConfig.AwsSQS.SQSRegion != nil {
						awsSQSMap["sqs_region"] = schedulesInfo.TaskNotifyConfig.AwsSQS.SQSRegion
					}

					if schedulesInfo.TaskNotifyConfig.AwsSQS.SQSQueueName != nil {
						awsSQSMap["sqs_queue_name"] = schedulesInfo.TaskNotifyConfig.AwsSQS.SQSQueueName
					}

					if schedulesInfo.TaskNotifyConfig.AwsSQS.S3SecretId != nil {
						awsSQSMap["s3_secret_id"] = schedulesInfo.TaskNotifyConfig.AwsSQS.S3SecretId
					}

					if schedulesInfo.TaskNotifyConfig.AwsSQS.S3SecretKey != nil {
						awsSQSMap["s3_secret_key"] = schedulesInfo.TaskNotifyConfig.AwsSQS.S3SecretKey
					}

					taskNotifyConfigMap["aws_sqs"] = []interface{}{awsSQSMap}
				}

				schedulesInfoMap["task_notify_config"] = []interface{}{taskNotifyConfigMap}
			}

			if schedulesInfo.CreateTime != nil {
				schedulesInfoMap["create_time"] = schedulesInfo.CreateTime
			}

			if schedulesInfo.UpdateTime != nil {
				schedulesInfoMap["update_time"] = schedulesInfo.UpdateTime
			}

			ids = append(ids, helper.Int64ToStr(*schedulesInfo.ScheduleId))
			tmpList = append(tmpList, schedulesInfoMap)
		}

		_ = d.Set("schedule_info_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
