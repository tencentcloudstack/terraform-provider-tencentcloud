package mps

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsTranscodeTemplateCreate,
		Read:   resourceTencentCloudMpsTranscodeTemplateRead,
		Update: resourceTencentCloudMpsTranscodeTemplateUpdate,
		Delete: resourceTencentCloudMpsTranscodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"container": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Transcoding template name, length limit: 64 characters.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description information, length limit: 256 characters.",
			},

			"remove_video": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to remove video data, value:0: reserved.1: remove.Default: 0.",
			},

			"remove_audio": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to remove audio data, value:0: reserved.1: remove.Default: 0.",
			},

			"video_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
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
							Description: "Filling method, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling method:stretch: Stretch, stretch each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the aspect ratio of the video unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video unchanged, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and fill the rest of the edge with Gaussian blur.Default: black.Note: Adaptive stream only supports stretch, black.",
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
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
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
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
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

			"enhance_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Audio and video enhancement configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"video_enhance": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Video Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frame_rate": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Interpolation frame rate configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"fps": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Frame rate, value range: [0, 100], unit: Hz.Default value: 0.Note: For transcoding, this parameter will override the Fps inside the VideoTemplate.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"super_resolution": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Super resolution configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type, optional value:lq: super-resolution for low-definition video with more noise.hq: super resolution for high-definition video.Default value: lq.Note: This field may return null, indicating that no valid value can be obtained.",
												},
												"size": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Super resolution multiple, optional value:2: currently only supports 2x super resolution.Default value: 2.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"hdr": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "HDR configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type, optional value: HDR10/HLG.Default value: HDR10.Note: The encoding method of video needs to be libx265.Note: Video encoding bit depth is 10.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"denoise": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Video Noise Reduction Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type, optional value: weak/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"image_quality_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Comprehensive Enhanced Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type, optional value: weak/normal/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"color_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Color Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type, optional value: weak/normal/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"sharp_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Detail Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"intensity": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"face_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Face Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"intensity": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"low_light_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Low Light Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type, optional value: normal.Default value: normal.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"scratch_repair": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "De-scratch configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"intensity": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"artifact_repair": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "De-artifact (glitch) configuration.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability configuration switch, optional value: ON/OFF.Default value: ON.",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
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
	}
}

func resourceTencentCloudMpsTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_transcode_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateTranscodeTemplateRequest()
		response   = mps.NewCreateTranscodeTemplateResponse()
		definition int64
	)
	if v, ok := d.GetOk("container"); ok {
		request.Container = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, _ := d.GetOk("remove_video"); v != nil {
		request.RemoveVideo = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("remove_audio"); v != nil {
		request.RemoveAudio = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video_template"); ok {
		videoTemplateInfo := mps.VideoTemplateInfo{}
		if v, ok := dMap["codec"]; ok {
			videoTemplateInfo.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["fps"]; ok {
			videoTemplateInfo.Fps = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["bitrate"]; ok {
			videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["resolution_adaptive"]; ok {
			videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
		}
		if v, ok := dMap["width"]; ok {
			videoTemplateInfo.Width = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["height"]; ok {
			videoTemplateInfo.Height = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["gop"]; ok {
			videoTemplateInfo.Gop = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["fill_type"]; ok {
			videoTemplateInfo.FillType = helper.String(v.(string))
		}
		if v, ok := dMap["vcrf"]; ok {
			videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
		}
		request.VideoTemplate = &videoTemplateInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio_template"); ok {
		audioTemplateInfo := mps.AudioTemplateInfo{}
		if v, ok := dMap["codec"]; ok {
			audioTemplateInfo.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["bitrate"]; ok {
			audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sample_rate"]; ok {
			audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["audio_channel"]; ok {
			audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
		}
		request.AudioTemplate = &audioTemplateInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "tehd_config"); ok {
		tEHDConfig := mps.TEHDConfig{}
		if v, ok := dMap["type"]; ok {
			tEHDConfig.Type = helper.String(v.(string))
		}
		if v, ok := dMap["max_video_bitrate"]; ok {
			tEHDConfig.MaxVideoBitrate = helper.IntInt64(v.(int))
		}
		request.TEHDConfig = &tEHDConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "enhance_config"); ok {
		enhanceConfig := mps.EnhanceConfig{}
		if videoEnhanceMap, ok := helper.InterfaceToMap(dMap, "video_enhance"); ok {
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
		request.EnhanceConfig = &enhanceConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateTranscodeTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps transcodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.Int64ToStr(definition))

	return resourceTencentCloudMpsTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudMpsTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_transcode_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	transcodeTemplate, err := service.DescribeMpsTranscodeTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if transcodeTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsTranscodeTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if transcodeTemplate.Container != nil {
		_ = d.Set("container", transcodeTemplate.Container)
	}

	if transcodeTemplate.Name != nil {
		_ = d.Set("name", transcodeTemplate.Name)
	}

	if transcodeTemplate.Comment != nil {
		_ = d.Set("comment", transcodeTemplate.Comment)
	}

	if transcodeTemplate.RemoveVideo != nil {
		_ = d.Set("remove_video", transcodeTemplate.RemoveVideo)
	}

	if transcodeTemplate.RemoveAudio != nil {
		_ = d.Set("remove_audio", transcodeTemplate.RemoveAudio)
	}

	if transcodeTemplate.VideoTemplate != nil {
		videoTemplateMap := map[string]interface{}{}

		if transcodeTemplate.VideoTemplate.Codec != nil {
			videoTemplateMap["codec"] = transcodeTemplate.VideoTemplate.Codec
		}

		if transcodeTemplate.VideoTemplate.Fps != nil {
			videoTemplateMap["fps"] = transcodeTemplate.VideoTemplate.Fps
		}

		if transcodeTemplate.VideoTemplate.Bitrate != nil {
			videoTemplateMap["bitrate"] = transcodeTemplate.VideoTemplate.Bitrate
		}

		if transcodeTemplate.VideoTemplate.ResolutionAdaptive != nil {
			videoTemplateMap["resolution_adaptive"] = transcodeTemplate.VideoTemplate.ResolutionAdaptive
		}

		if transcodeTemplate.VideoTemplate.Width != nil {
			videoTemplateMap["width"] = transcodeTemplate.VideoTemplate.Width
		}

		if transcodeTemplate.VideoTemplate.Height != nil {
			videoTemplateMap["height"] = transcodeTemplate.VideoTemplate.Height
		}

		if transcodeTemplate.VideoTemplate.Gop != nil {
			videoTemplateMap["gop"] = transcodeTemplate.VideoTemplate.Gop
		}

		if transcodeTemplate.VideoTemplate.FillType != nil {
			videoTemplateMap["fill_type"] = transcodeTemplate.VideoTemplate.FillType
		}

		if transcodeTemplate.VideoTemplate.Vcrf != nil {
			videoTemplateMap["vcrf"] = transcodeTemplate.VideoTemplate.Vcrf
		}

		_ = d.Set("video_template", []interface{}{videoTemplateMap})
	}

	if transcodeTemplate.AudioTemplate != nil {
		audioTemplateMap := map[string]interface{}{}

		if transcodeTemplate.AudioTemplate.Codec != nil {
			audioTemplateMap["codec"] = transcodeTemplate.AudioTemplate.Codec
		}

		if transcodeTemplate.AudioTemplate.Bitrate != nil {
			audioTemplateMap["bitrate"] = transcodeTemplate.AudioTemplate.Bitrate
		}

		if transcodeTemplate.AudioTemplate.SampleRate != nil {
			audioTemplateMap["sample_rate"] = transcodeTemplate.AudioTemplate.SampleRate
		}

		if transcodeTemplate.AudioTemplate.AudioChannel != nil {
			audioTemplateMap["audio_channel"] = transcodeTemplate.AudioTemplate.AudioChannel
		}

		_ = d.Set("audio_template", []interface{}{audioTemplateMap})
	}

	if transcodeTemplate.TEHDConfig != nil {
		tEHDConfigMap := map[string]interface{}{}

		if transcodeTemplate.TEHDConfig.Type != nil {
			tEHDConfigMap["type"] = transcodeTemplate.TEHDConfig.Type
		}

		if transcodeTemplate.TEHDConfig.MaxVideoBitrate != nil {
			tEHDConfigMap["max_video_bitrate"] = transcodeTemplate.TEHDConfig.MaxVideoBitrate
		}

		_ = d.Set("tehd_config", []interface{}{tEHDConfigMap})
	}

	if transcodeTemplate.EnhanceConfig != nil {
		enhanceConfigMap := map[string]interface{}{}

		if transcodeTemplate.EnhanceConfig.VideoEnhance != nil {
			videoEnhanceMap := map[string]interface{}{}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate != nil {
				frameRateMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate.Switch != nil {
					frameRateMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate.Fps != nil {
					frameRateMap["fps"] = transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate.Fps
				}

				videoEnhanceMap["frame_rate"] = []interface{}{frameRateMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution != nil {
				superResolutionMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Switch != nil {
					superResolutionMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Type != nil {
					superResolutionMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Type
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Size != nil {
					superResolutionMap["size"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Size
				}

				videoEnhanceMap["super_resolution"] = []interface{}{superResolutionMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr != nil {
				hdrMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr.Switch != nil {
					hdrMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr.Type != nil {
					hdrMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr.Type
				}

				videoEnhanceMap["hdr"] = []interface{}{hdrMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise != nil {
				denoiseMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise.Switch != nil {
					denoiseMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise.Type != nil {
					denoiseMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise.Type
				}

				videoEnhanceMap["denoise"] = []interface{}{denoiseMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance != nil {
				imageQualityEnhanceMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance.Switch != nil {
					imageQualityEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance.Type != nil {
					imageQualityEnhanceMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance.Type
				}

				videoEnhanceMap["image_quality_enhance"] = []interface{}{imageQualityEnhanceMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance != nil {
				colorEnhanceMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance.Switch != nil {
					colorEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance.Type != nil {
					colorEnhanceMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance.Type
				}

				videoEnhanceMap["color_enhance"] = []interface{}{colorEnhanceMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance != nil {
				sharpEnhanceMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance.Switch != nil {
					sharpEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance.Intensity != nil {
					sharpEnhanceMap["intensity"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance.Intensity
				}

				videoEnhanceMap["sharp_enhance"] = []interface{}{sharpEnhanceMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance != nil {
				faceEnhanceMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance.Switch != nil {
					faceEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance.Intensity != nil {
					faceEnhanceMap["intensity"] = transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance.Intensity
				}

				videoEnhanceMap["face_enhance"] = []interface{}{faceEnhanceMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance != nil {
				lowLightEnhanceMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance.Switch != nil {
					lowLightEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance.Type != nil {
					lowLightEnhanceMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance.Type
				}

				videoEnhanceMap["low_light_enhance"] = []interface{}{lowLightEnhanceMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair != nil {
				scratchRepairMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair.Switch != nil {
					scratchRepairMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair.Intensity != nil {
					scratchRepairMap["intensity"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair.Intensity
				}

				videoEnhanceMap["scratch_repair"] = []interface{}{scratchRepairMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair != nil {
				artifactRepairMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair.Switch != nil {
					artifactRepairMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair.Type != nil {
					artifactRepairMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair.Type
				}

				videoEnhanceMap["artifact_repair"] = []interface{}{artifactRepairMap}
			}

			enhanceConfigMap["video_enhance"] = []interface{}{videoEnhanceMap}
		}

		_ = d.Set("enhance_config", []interface{}{enhanceConfigMap})
	}

	return nil
}

func resourceTencentCloudMpsTranscodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_transcode_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyTranscodeTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToInt64Point(definition)

	if d.HasChange("container") {
		if v, ok := d.GetOk("container"); ok {
			request.Container = helper.String(v.(string))
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
	}

	if d.HasChange("remove_video") {
		if v, _ := d.GetOk("remove_video"); v != nil {
			request.RemoveVideo = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("remove_audio") {
		if v, _ := d.GetOk("remove_audio"); v != nil {
			request.RemoveAudio = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("video_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "video_template"); ok {
			videoTemplateInfo := mps.VideoTemplateInfoForUpdate{}
			if v, ok := dMap["codec"]; ok {
				videoTemplateInfo.Codec = helper.String(v.(string))
			}
			if v, ok := dMap["fps"]; ok {
				videoTemplateInfo.Fps = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["bitrate"]; ok {
				videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["resolution_adaptive"]; ok {
				videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
			}
			if v, ok := dMap["width"]; ok {
				videoTemplateInfo.Width = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["height"]; ok {
				videoTemplateInfo.Height = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["gop"]; ok {
				videoTemplateInfo.Gop = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["fill_type"]; ok {
				videoTemplateInfo.FillType = helper.String(v.(string))
			}
			if v, ok := dMap["vcrf"]; ok {
				videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
			}
			request.VideoTemplate = &videoTemplateInfo
		}
	}

	if d.HasChange("audio_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "audio_template"); ok {
			audioTemplateInfo := mps.AudioTemplateInfoForUpdate{}
			if v, ok := dMap["codec"]; ok {
				audioTemplateInfo.Codec = helper.String(v.(string))
			}
			if v, ok := dMap["bitrate"]; ok {
				audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["sample_rate"]; ok {
				audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["audio_channel"]; ok {
				audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
			}
			request.AudioTemplate = &audioTemplateInfo
		}
	}

	if d.HasChange("tehd_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "tehd_config"); ok {
			tEHDConfig := mps.TEHDConfigForUpdate{}
			if v, ok := dMap["type"]; ok {
				tEHDConfig.Type = helper.String(v.(string))
			}
			if v, ok := dMap["max_video_bitrate"]; ok {
				tEHDConfig.MaxVideoBitrate = helper.IntInt64(v.(int))
			}
			request.TEHDConfig = &tEHDConfig
		}
	}

	if d.HasChange("enhance_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "enhance_config"); ok {
			enhanceConfig := mps.EnhanceConfig{}
			if videoEnhanceMap, ok := helper.InterfaceToMap(dMap, "video_enhance"); ok {
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
			request.EnhanceConfig = &enhanceConfig
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyTranscodeTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps transcodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudMpsTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_transcode_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsTranscodeTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
