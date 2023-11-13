/*
Use this data source to query detailed information of mps transcode_templates

Example Usage

```hcl
data "tencentcloud_mps_transcode_templates" "transcode_templates" {
  definitions = &lt;nil&gt;
  type = &lt;nil&gt;
  container_type = &lt;nil&gt;
  t_e_h_d_type = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  transcode_type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  transcode_template_set {
		definition = &lt;nil&gt;
		container = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		type = &lt;nil&gt;
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
		}
		audio_template {
			codec = &lt;nil&gt;
			bitrate = &lt;nil&gt;
			sample_rate = &lt;nil&gt;
			audio_channel = &lt;nil&gt;
		}
		t_e_h_d_config {
			type = &lt;nil&gt;
			max_video_bitrate = &lt;nil&gt;
		}
		container_type = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		enhance_config {
			video_enhance {
				frame_rate {
					switch = &lt;nil&gt;
					fps = &lt;nil&gt;
				}
				super_resolution {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
					size = &lt;nil&gt;
				}
				hdr {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				denoise {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				image_quality_enhance {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				color_enhance {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				sharp_enhance {
					switch = &lt;nil&gt;
					intensity = &lt;nil&gt;
				}
				face_enhance {
					switch = &lt;nil&gt;
					intensity = &lt;nil&gt;
				}
				low_light_enhance {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				scratch_repair {
					switch = &lt;nil&gt;
					intensity = &lt;nil&gt;
				}
				artifact_repair {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
			}
		}

  }
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

func dataSourceTencentCloudMpsTranscodeTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsTranscodeTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The transcoding template uniquely identifies the filter condition, and the array length limit: 100.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template type filter condition, optional value:Preset: system preset template.Custom: user-defined template.",
			},

			"container_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Encapsulation format filter condition, optional value:Video: video format, the encapsulation format board that can contain both video stream and audio stream.PureAudio: pure audio format, which can only contain audio stream encapsulation format.",
			},

			"t_e_h_d_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Extreme HD filter condition, used to filter common transcoding or extreme high-definition transcoding templates, optional value:Common: common transcoding template.TEHD: Extreme HD template.",
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

			"transcode_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template type (replaces old version TEHDType), optional value:Common: common transcoding template.TEHD: Extreme HD template.Enhance: audio and video enhancement template.Empty by default, unlimited types.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"transcode_template_set": {
				Type:        schema.TypeList,
				Description: "Transcoding template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Description: "The unique identifier of the transcoding template.",
						},
						"container": {
							Type:        schema.TypeString,
							Description: "Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Transcoding template name.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "Template description information.",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Template type, optional value:Preset: system preset template.Custom: user-defined template.",
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
							Description: "Video stream configuration parameters, this field is valid only when RemoveVideo is 0.",
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
							Description: "Audio stream configuration parameters, this field is valid only when RemoveAudio is 0.",
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
							Description: "Ultra-fast HD transcoding parameters.Note: This field may return null, indicating that no valid value can be obtained.",
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
						"container_type": {
							Type:        schema.TypeString,
							Description: "Encapsulation format filter condition, optional value:Video: video format, the encapsulation format board that can contain both video stream and audio stream.PureAudio: pure audio format, which can only contain audio stream encapsulation format.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Template creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Template last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"enhance_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Audio and video enhancement configuration.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"video_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Video Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"frame_rate": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Interpolation frame rate configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"fps": {
																Type:        schema.TypeInt,
																Description: "Frame rate, value range: [0, 100], unit: Hz.Default value: 0.Note: For transcoding, this parameter will override the Fps inside the VideoTemplate.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"super_resolution": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Super resolution configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"type": {
																Type:        schema.TypeString,
																Description: "Type, optional value:lq: super-resolution for low-definition video with more noise.hq: super resolution for high-definition video.Default value: lq.Note: This field may return null, indicating that no valid value can be obtained.",
															},
															"size": {
																Type:        schema.TypeInt,
																Description: "Super resolution multiple, optional value:2: currently only supports 2x super resolution.Default value: 2.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"hdr": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "HDR configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"type": {
																Type:        schema.TypeString,
																Description: "Type, optional value: HDR10/HLG.Default value: HDR10.Note: The encoding method of video needs to be libx265.Note: Video encoding bit depth is 10.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"denoise": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Video Noise Reduction Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"type": {
																Type:        schema.TypeString,
																Description: "Type, optional value: weak/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"image_quality_enhance": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Comprehensive Enhanced Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"type": {
																Type:        schema.TypeString,
																Description: "Type, optional value: weak/normal/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"color_enhance": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Color Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"type": {
																Type:        schema.TypeString,
																Description: "Type, optional value: weak/normal/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"sharp_enhance": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Detail Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"intensity": {
																Type:        schema.TypeFloat,
																Description: "Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"face_enhance": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Face Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"intensity": {
																Type:        schema.TypeFloat,
																Description: "Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"low_light_enhance": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "Low Light Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"type": {
																Type:        schema.TypeString,
																Description: "Type, optional value: normal.Default value: normal.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"scratch_repair": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "De-scratch configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"intensity": {
																Type:        schema.TypeFloat,
																Description: "Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.",
															},
														},
													},
												},
												"artifact_repair": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Description: "De-artifact (glitch) configuration.Note: This field may return null, indicating that no valid value can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"switch": {
																Type:        schema.TypeString,
																Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
															},
															"type": {
																Type:        schema.TypeString,
																Description: "Type, optional value: weak/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.",
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMpsTranscodeTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_transcode_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("definitions"); ok {
		definitionsSet := v.(*schema.Set).List()
		for i := range definitionsSet {
			definitions := definitionsSet[i].(int)
			paramMap["Definitions"] = append(paramMap["Definitions"], helper.IntInt64(definitions))
		}
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("container_type"); ok {
		paramMap["ContainerType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("t_e_h_d_type"); ok {
		paramMap["TEHDType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("transcode_type"); ok {
		paramMap["TranscodeType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("transcode_template_set"); ok {
		transcodeTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.TranscodeTemplate, 0, len(transcodeTemplateSetSet))

		for _, item := range transcodeTemplateSetSet {
			transcodeTemplate := mps.TranscodeTemplate{}
			transcodeTemplateMap := item.(map[string]interface{})

			if v, ok := transcodeTemplateMap["definition"]; ok {
				transcodeTemplate.Definition = helper.String(v.(string))
			}
			if v, ok := transcodeTemplateMap["container"]; ok {
				transcodeTemplate.Container = helper.String(v.(string))
			}
			if v, ok := transcodeTemplateMap["name"]; ok {
				transcodeTemplate.Name = helper.String(v.(string))
			}
			if v, ok := transcodeTemplateMap["comment"]; ok {
				transcodeTemplate.Comment = helper.String(v.(string))
			}
			if v, ok := transcodeTemplateMap["type"]; ok {
				transcodeTemplate.Type = helper.String(v.(string))
			}
			if v, ok := transcodeTemplateMap["remove_video"]; ok {
				transcodeTemplate.RemoveVideo = helper.IntInt64(v.(int))
			}
			if v, ok := transcodeTemplateMap["remove_audio"]; ok {
				transcodeTemplate.RemoveAudio = helper.IntInt64(v.(int))
			}
			if videoTemplateMap, ok := helper.InterfaceToMap(transcodeTemplateMap, "video_template"); ok {
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
				transcodeTemplate.VideoTemplate = &videoTemplateInfo
			}
			if audioTemplateMap, ok := helper.InterfaceToMap(transcodeTemplateMap, "audio_template"); ok {
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
				transcodeTemplate.AudioTemplate = &audioTemplateInfo
			}
			if tEHDConfigMap, ok := helper.InterfaceToMap(transcodeTemplateMap, "t_e_h_d_config"); ok {
				tEHDConfig := mps.TEHDConfig{}
				if v, ok := tEHDConfigMap["type"]; ok {
					tEHDConfig.Type = helper.String(v.(string))
				}
				if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
					tEHDConfig.MaxVideoBitrate = helper.IntUint64(v.(int))
				}
				transcodeTemplate.TEHDConfig = &tEHDConfig
			}
			if v, ok := transcodeTemplateMap["container_type"]; ok {
				transcodeTemplate.ContainerType = helper.String(v.(string))
			}
			if v, ok := transcodeTemplateMap["create_time"]; ok {
				transcodeTemplate.CreateTime = helper.String(v.(string))
			}
			if v, ok := transcodeTemplateMap["update_time"]; ok {
				transcodeTemplate.UpdateTime = helper.String(v.(string))
			}
			if enhanceConfigMap, ok := helper.InterfaceToMap(transcodeTemplateMap, "enhance_config"); ok {
				enhanceConfig := mps.EnhanceConfig{}
				if videoEnhanceMap, ok := helper.InterfaceToMap(enhanceConfigMap, "video_enhance"); ok {
					videoEnhanceConfig := mps.VideoEnhanceConfig{}
					if frameRateMap, ok := helper.InterfaceToMap(videoEnhanceMap, "frame_rate"); ok {
						frameRateConfig := mps.FrameRateConfig{}
						if v, ok := frameRateMap["switch"]; ok {
							frameRateConfig.Switch = helper.String(v.(string))
						}
						if v, ok := frameRateMap["fps"]; ok {
							frameRateConfig.Fps = helper.IntUint64(v.(int))
						}
						videoEnhanceConfig.FrameRate = &frameRateConfig
					}
					if superResolutionMap, ok := helper.InterfaceToMap(videoEnhanceMap, "super_resolution"); ok {
						superResolutionConfig := mps.SuperResolutionConfig{}
						if v, ok := superResolutionMap["switch"]; ok {
							superResolutionConfig.Switch = helper.String(v.(string))
						}
						if v, ok := superResolutionMap["type"]; ok {
							superResolutionConfig.Type = helper.String(v.(string))
						}
						if v, ok := superResolutionMap["size"]; ok {
							superResolutionConfig.Size = helper.IntInt64(v.(int))
						}
						videoEnhanceConfig.SuperResolution = &superResolutionConfig
					}
					if hdrMap, ok := helper.InterfaceToMap(videoEnhanceMap, "hdr"); ok {
						hdrConfig := mps.HdrConfig{}
						if v, ok := hdrMap["switch"]; ok {
							hdrConfig.Switch = helper.String(v.(string))
						}
						if v, ok := hdrMap["type"]; ok {
							hdrConfig.Type = helper.String(v.(string))
						}
						videoEnhanceConfig.Hdr = &hdrConfig
					}
					if denoiseMap, ok := helper.InterfaceToMap(videoEnhanceMap, "denoise"); ok {
						videoDenoiseConfig := mps.VideoDenoiseConfig{}
						if v, ok := denoiseMap["switch"]; ok {
							videoDenoiseConfig.Switch = helper.String(v.(string))
						}
						if v, ok := denoiseMap["type"]; ok {
							videoDenoiseConfig.Type = helper.String(v.(string))
						}
						videoEnhanceConfig.Denoise = &videoDenoiseConfig
					}
					if imageQualityEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "image_quality_enhance"); ok {
						imageQualityEnhanceConfig := mps.ImageQualityEnhanceConfig{}
						if v, ok := imageQualityEnhanceMap["switch"]; ok {
							imageQualityEnhanceConfig.Switch = helper.String(v.(string))
						}
						if v, ok := imageQualityEnhanceMap["type"]; ok {
							imageQualityEnhanceConfig.Type = helper.String(v.(string))
						}
						videoEnhanceConfig.ImageQualityEnhance = &imageQualityEnhanceConfig
					}
					if colorEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "color_enhance"); ok {
						colorEnhanceConfig := mps.ColorEnhanceConfig{}
						if v, ok := colorEnhanceMap["switch"]; ok {
							colorEnhanceConfig.Switch = helper.String(v.(string))
						}
						if v, ok := colorEnhanceMap["type"]; ok {
							colorEnhanceConfig.Type = helper.String(v.(string))
						}
						videoEnhanceConfig.ColorEnhance = &colorEnhanceConfig
					}
					if sharpEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "sharp_enhance"); ok {
						sharpEnhanceConfig := mps.SharpEnhanceConfig{}
						if v, ok := sharpEnhanceMap["switch"]; ok {
							sharpEnhanceConfig.Switch = helper.String(v.(string))
						}
						if v, ok := sharpEnhanceMap["intensity"]; ok {
							sharpEnhanceConfig.Intensity = helper.Float64(v.(float64))
						}
						videoEnhanceConfig.SharpEnhance = &sharpEnhanceConfig
					}
					if faceEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "face_enhance"); ok {
						faceEnhanceConfig := mps.FaceEnhanceConfig{}
						if v, ok := faceEnhanceMap["switch"]; ok {
							faceEnhanceConfig.Switch = helper.String(v.(string))
						}
						if v, ok := faceEnhanceMap["intensity"]; ok {
							faceEnhanceConfig.Intensity = helper.Float64(v.(float64))
						}
						videoEnhanceConfig.FaceEnhance = &faceEnhanceConfig
					}
					if lowLightEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "low_light_enhance"); ok {
						lowLightEnhanceConfig := mps.LowLightEnhanceConfig{}
						if v, ok := lowLightEnhanceMap["switch"]; ok {
							lowLightEnhanceConfig.Switch = helper.String(v.(string))
						}
						if v, ok := lowLightEnhanceMap["type"]; ok {
							lowLightEnhanceConfig.Type = helper.String(v.(string))
						}
						videoEnhanceConfig.LowLightEnhance = &lowLightEnhanceConfig
					}
					if scratchRepairMap, ok := helper.InterfaceToMap(videoEnhanceMap, "scratch_repair"); ok {
						scratchRepairConfig := mps.ScratchRepairConfig{}
						if v, ok := scratchRepairMap["switch"]; ok {
							scratchRepairConfig.Switch = helper.String(v.(string))
						}
						if v, ok := scratchRepairMap["intensity"]; ok {
							scratchRepairConfig.Intensity = helper.Float64(v.(float64))
						}
						videoEnhanceConfig.ScratchRepair = &scratchRepairConfig
					}
					if artifactRepairMap, ok := helper.InterfaceToMap(videoEnhanceMap, "artifact_repair"); ok {
						artifactRepairConfig := mps.ArtifactRepairConfig{}
						if v, ok := artifactRepairMap["switch"]; ok {
							artifactRepairConfig.Switch = helper.String(v.(string))
						}
						if v, ok := artifactRepairMap["type"]; ok {
							artifactRepairConfig.Type = helper.String(v.(string))
						}
						videoEnhanceConfig.ArtifactRepair = &artifactRepairConfig
					}
					enhanceConfig.VideoEnhance = &videoEnhanceConfig
				}
				transcodeTemplate.EnhanceConfig = &enhanceConfig
			}
			tmpSet = append(tmpSet, &transcodeTemplate)
		}
		paramMap["transcode_template_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var transcodeTemplateSet []*mps.TranscodeTemplate

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsTranscodeTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		transcodeTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(transcodeTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(transcodeTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
