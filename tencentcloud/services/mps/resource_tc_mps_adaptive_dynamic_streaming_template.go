package mps

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsAdaptiveDynamicStreamingTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateCreate,
		Read:   resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateRead,
		Update: resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateUpdate,
		Delete: resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"format": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Adaptive transcoding format, value range:HLS, MPEG-DASH.",
			},

			"stream_infos": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Convert adaptive code stream to output sub-stream parameter information, and output up to 10 sub-streams.Note: The frame rate of each sub-stream must be consistent; if not, the frame rate of the first sub-stream is used as the output frame rate.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"video": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Video parameter information.",
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
										Description: "Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.Note: In adaptive mode, Width cannot be smaller than Height.",
									},
									"width": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum value of the width (or long side) of the video streaming, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
									},
									"height": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum value of the height (or short side) of the video streaming, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
									},
									"gop": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.When filling 0 or not filling, the system will automatically set the gop length.",
									},
									"fill_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and use Gaussian blur for the rest of the edge.Default value: black.Note: Adaptive stream only supports stretch, black.",
									},
									"vcrf": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Video constant bit rate control factor, the value range is [1, 51].If this parameter is specified, the code rate control method of CRF will be used for transcoding (the video code rate will no longer take effect).If there is no special requirement, it is not recommended to specify this parameter.",
									},
								},
							},
						},
						"audio": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Audio parameter information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Encoding format of audio stream.When the outer parameter Container is mp3, the optional value is:libmp3lame.When the outer parameter Container is ogg or flac, the optional value is:flac.When the outer parameter Container is m4a, the optional value is:libfdk_aac.libmp3lame.ac3.When the outer parameter Container is mp4 or flv, the optional value is:libfdk_aac: more suitable for mp4.libmp3lame: more suitable for flv.When the outer parameter Container is hls, the optional value is:libfdk_aac.libmp3lame.",
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
						"remove_audio": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to remove audio stream, value:0: reserved.1: remove.",
						},
						"remove_video": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Whether to remove video stream, value:0: reserved.1: remove.",
						},
					},
				},
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template name, length limit: 64 characters.",
			},

			"disable_higher_video_bitrate": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to prohibit video from low bit rate to high bit rate, value range:0: no.1: yes.Default value: 0.",
			},

			"disable_higher_video_resolution": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to prohibit the conversion of video resolution to high resolution, value range:0: no.1: yes.Default value: 0.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description information, length limit: 256 characters.",
			},

			"pure_audio": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Indicates whether it is audio-only. 0 means video template, 1 means audio-only template.\nWhen the value is 1.\n1. StreamInfos.N.RemoveVideo=1\n2. StreamInfos.N.RemoveAudio=0\n3. StreamInfos.N.Video.Codec=copy\nWhen the value is 0.\n1. StreamInfos.N.Video.Codec cannot be copy.\n2. StreamInfos.N.Video.Fps cannot be null.\nNote: This value only distinguishes template types. The task uses the values of RemoveAudio and RemoveVideo.",
			},

			"segment_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Segment type. Valid values: \nts-segment: HLS+TS segment\nts-byterange: HLS+TS byte range\nmp4-segment: HLS+MP4 segment\nmp4-byterange: HLS/DASH+MP4 byte range\nts-packed-audio: HLS+TS+Packed Audio segment\nmp4-packed-audio: HLS+MP4+Packed Audio segment\nts-ts-segment: HLS+TS+TS segment\nts-ts-byterange: HLS+TS+TS byte range\nmp4-mp4-segment: HLS+MP4+MP4 segment\nmp4-mp4-byterange: HLS/DASH+MP4+MP4 byte range\nts-packed-audio-byterange: HLS+TS+Packed Audio byte range\nmp4-packed-audio-byterange: HLS+MP4+Packed Audio byte range.\n Default value: ts-segment. Note: The segment format for adaptive bitrate streaming is determined by this field. For DASH format, SegmentType can only be mp4-byterange or mp4-mp4-byterange.",
			},
		},
	}
}

func resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_adaptive_dynamic_streaming_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mps.NewCreateAdaptiveDynamicStreamingTemplateRequest()
		response   = mps.NewCreateAdaptiveDynamicStreamingTemplateResponse()
		definition uint64
	)

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_infos"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			adaptiveStreamTemplate := mps.AdaptiveStreamTemplate{}
			if videoMap, ok := helper.InterfaceToMap(dMap, "video"); ok {
				videoTemplateInfo := mps.VideoTemplateInfo{}
				if v, ok := videoMap["codec"]; ok {
					videoTemplateInfo.Codec = helper.String(v.(string))
				}

				if v, ok := videoMap["fps"]; ok {
					videoTemplateInfo.Fps = helper.IntInt64(v.(int))
				}

				if v, ok := videoMap["bitrate"]; ok {
					videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
				}

				if v, ok := videoMap["resolution_adaptive"]; ok {
					videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
				}

				if v, ok := videoMap["width"]; ok {
					videoTemplateInfo.Width = helper.IntUint64(v.(int))
				}

				if v, ok := videoMap["height"]; ok {
					videoTemplateInfo.Height = helper.IntUint64(v.(int))
				}

				if v, ok := videoMap["gop"]; ok {
					videoTemplateInfo.Gop = helper.IntUint64(v.(int))
				}

				if v, ok := videoMap["fill_type"]; ok {
					videoTemplateInfo.FillType = helper.String(v.(string))
				}

				if v, ok := videoMap["vcrf"]; ok {
					videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
				}

				adaptiveStreamTemplate.Video = &videoTemplateInfo
			}

			if audioMap, ok := helper.InterfaceToMap(dMap, "audio"); ok {
				audioTemplateInfo := mps.AudioTemplateInfo{}
				if v, ok := audioMap["codec"]; ok {
					audioTemplateInfo.Codec = helper.String(v.(string))
				}

				if v, ok := audioMap["bitrate"]; ok {
					audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
				}

				if v, ok := audioMap["sample_rate"]; ok {
					audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
				}

				if v, ok := audioMap["audio_channel"]; ok {
					audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
				}

				adaptiveStreamTemplate.Audio = &audioTemplateInfo
			}

			if v, ok := dMap["remove_audio"]; ok {
				adaptiveStreamTemplate.RemoveAudio = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["remove_video"]; ok {
				adaptiveStreamTemplate.RemoveVideo = helper.IntUint64(v.(int))
			}

			request.StreamInfos = append(request.StreamInfos, &adaptiveStreamTemplate)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disable_higher_video_bitrate"); ok {
		request.DisableHigherVideoBitrate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("disable_higher_video_resolution"); ok {
		request.DisableHigherVideoResolution = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("pure_audio"); ok {
		request.PureAudio = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("segment_type"); ok {
		request.SegmentType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateAdaptiveDynamicStreamingTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mps adaptive dynamic streaming template failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mps adaptiveDynamicStreamingTemplate failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Definition == nil {
		return fmt.Errorf("Definition is nil.")
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateRead(d, meta)
}

func resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_adaptive_dynamic_streaming_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		definition = d.Id()
	)

	adaptiveDynamicStreamingTemplate, err := service.DescribeMpsAdaptiveDynamicStreamingTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if adaptiveDynamicStreamingTemplate == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mps_adaptive_dynamic_streaming_template` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if adaptiveDynamicStreamingTemplate.Format != nil {
		_ = d.Set("format", adaptiveDynamicStreamingTemplate.Format)
	}

	if adaptiveDynamicStreamingTemplate.StreamInfos != nil {
		streamInfosList := []interface{}{}
		for _, streamInfos := range adaptiveDynamicStreamingTemplate.StreamInfos {
			streamInfosMap := map[string]interface{}{}
			if streamInfos.Video != nil {
				videoMap := map[string]interface{}{}
				if streamInfos.Video.Codec != nil {
					videoMap["codec"] = streamInfos.Video.Codec
				}

				if streamInfos.Video.Fps != nil {
					videoMap["fps"] = streamInfos.Video.Fps
				}

				if streamInfos.Video.Bitrate != nil {
					videoMap["bitrate"] = streamInfos.Video.Bitrate
				}

				if streamInfos.Video.ResolutionAdaptive != nil {
					videoMap["resolution_adaptive"] = streamInfos.Video.ResolutionAdaptive
				}

				if streamInfos.Video.Width != nil {
					videoMap["width"] = streamInfos.Video.Width
				}

				if streamInfos.Video.Height != nil {
					videoMap["height"] = streamInfos.Video.Height
				}

				if streamInfos.Video.Gop != nil {
					videoMap["gop"] = streamInfos.Video.Gop
				}

				if streamInfos.Video.FillType != nil {
					videoMap["fill_type"] = streamInfos.Video.FillType
				}

				if streamInfos.Video.Vcrf != nil {
					videoMap["vcrf"] = streamInfos.Video.Vcrf
				}

				streamInfosMap["video"] = []interface{}{videoMap}
			}

			if streamInfos.Audio != nil {
				audioMap := map[string]interface{}{}
				if streamInfos.Audio.Codec != nil {
					audioMap["codec"] = streamInfos.Audio.Codec
				}

				if streamInfos.Audio.Bitrate != nil {
					audioMap["bitrate"] = streamInfos.Audio.Bitrate
				}

				if streamInfos.Audio.SampleRate != nil {
					audioMap["sample_rate"] = streamInfos.Audio.SampleRate
				}

				if streamInfos.Audio.AudioChannel != nil {
					audioMap["audio_channel"] = streamInfos.Audio.AudioChannel
				}

				streamInfosMap["audio"] = []interface{}{audioMap}
			}

			if streamInfos.RemoveAudio != nil {
				streamInfosMap["remove_audio"] = streamInfos.RemoveAudio
			}

			if streamInfos.RemoveVideo != nil {
				streamInfosMap["remove_video"] = streamInfos.RemoveVideo
			}

			streamInfosList = append(streamInfosList, streamInfosMap)
		}

		_ = d.Set("stream_infos", streamInfosList)

	}

	if adaptiveDynamicStreamingTemplate.Name != nil {
		_ = d.Set("name", adaptiveDynamicStreamingTemplate.Name)
	}

	if adaptiveDynamicStreamingTemplate.DisableHigherVideoBitrate != nil {
		_ = d.Set("disable_higher_video_bitrate", adaptiveDynamicStreamingTemplate.DisableHigherVideoBitrate)
	}

	if adaptiveDynamicStreamingTemplate.DisableHigherVideoResolution != nil {
		_ = d.Set("disable_higher_video_resolution", adaptiveDynamicStreamingTemplate.DisableHigherVideoResolution)
	}

	if adaptiveDynamicStreamingTemplate.Comment != nil {
		_ = d.Set("comment", adaptiveDynamicStreamingTemplate.Comment)
	}

	if adaptiveDynamicStreamingTemplate.PureAudio != nil {
		_ = d.Set("pure_audio", adaptiveDynamicStreamingTemplate.PureAudio)
	}

	if adaptiveDynamicStreamingTemplate.SegmentType != nil {
		_ = d.Set("segment_type", adaptiveDynamicStreamingTemplate.SegmentType)
	}

	return nil
}

func resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_adaptive_dynamic_streaming_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mps.NewModifyAdaptiveDynamicStreamingTemplateRequest()
		definition = d.Id()
	)

	needChange := false
	mutableArgs := []string{"format", "stream_infos", "name", "disable_higher_video_bitrate", "disable_higher_video_resolution", "comment", "pure_audio", "segment_type"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("format"); ok {
			request.Format = helper.String(v.(string))
		}

		if v, ok := d.GetOk("stream_infos"); ok {
			for _, item := range v.([]interface{}) {
				adaptiveStreamTemplateMap := item.(map[string]interface{})
				adaptiveStreamTemplate := mps.AdaptiveStreamTemplate{}
				if videoMap, ok := helper.InterfaceToMap(adaptiveStreamTemplateMap, "video"); ok {
					videoTemplateInfo := mps.VideoTemplateInfo{}
					if v, ok := videoMap["codec"]; ok {
						videoTemplateInfo.Codec = helper.String(v.(string))
					}

					if v, ok := videoMap["fps"]; ok {
						videoTemplateInfo.Fps = helper.IntInt64(v.(int))
					}

					if v, ok := videoMap["bitrate"]; ok {
						videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
					}

					if v, ok := videoMap["resolution_adaptive"]; ok {
						videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
					}

					if v, ok := videoMap["width"]; ok {
						videoTemplateInfo.Width = helper.IntUint64(v.(int))
					}

					if v, ok := videoMap["height"]; ok {
						videoTemplateInfo.Height = helper.IntUint64(v.(int))
					}

					if v, ok := videoMap["gop"]; ok {
						videoTemplateInfo.Gop = helper.IntUint64(v.(int))
					}

					if v, ok := videoMap["fill_type"]; ok {
						videoTemplateInfo.FillType = helper.String(v.(string))
					}

					if v, ok := videoMap["vcrf"]; ok {
						videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
					}

					adaptiveStreamTemplate.Video = &videoTemplateInfo
				}

				if audioMap, ok := helper.InterfaceToMap(adaptiveStreamTemplateMap, "audio"); ok {
					audioTemplateInfo := mps.AudioTemplateInfo{}
					if v, ok := audioMap["codec"]; ok {
						audioTemplateInfo.Codec = helper.String(v.(string))
					}

					if v, ok := audioMap["bitrate"]; ok {
						audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
					}

					if v, ok := audioMap["sample_rate"]; ok {
						audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
					}

					if v, ok := audioMap["audio_channel"]; ok {
						audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
					}

					adaptiveStreamTemplate.Audio = &audioTemplateInfo
				}

				if v, ok := adaptiveStreamTemplateMap["remove_audio"]; ok {
					adaptiveStreamTemplate.RemoveAudio = helper.IntUint64(v.(int))
				}

				if v, ok := adaptiveStreamTemplateMap["remove_video"]; ok {
					adaptiveStreamTemplate.RemoveVideo = helper.IntUint64(v.(int))
				}

				request.StreamInfos = append(request.StreamInfos, &adaptiveStreamTemplate)
			}
		}

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("disable_higher_video_bitrate"); ok {
			request.DisableHigherVideoBitrate = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("disable_higher_video_resolution"); ok {
			request.DisableHigherVideoResolution = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("pure_audio"); ok {
			request.PureAudio = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("segment_type"); ok {
			request.SegmentType = helper.String(v.(string))
		}

		request.Definition = helper.StrToUint64Point(definition)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyAdaptiveDynamicStreamingTemplate(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mps adaptiveDynamicStreamingTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateRead(d, meta)
}

func resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_adaptive_dynamic_streaming_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		definition = d.Id()
	)

	if err := service.DeleteMpsAdaptiveDynamicStreamingTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
