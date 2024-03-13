package vod

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVodTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodTranscodeTemplateCreate,
		Read:   resourceTencentCloudVodTranscodeTemplateRead,
		Update: resourceTencentCloudVodTranscodeTemplateUpdate,
		Delete: resourceTencentCloudVodTranscodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"container": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The container format. Valid values: `mp4`, `flv`, `hls`, `mp3`, `flac`, `ogg`, `m4a`, `wav` ( `mp3`, `flac`, `ogg`, `m4a`, and `wav` are audio file formats).",
			},

			"sub_app_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Transcoding template name. Length limit: 64 characters.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description. Length limit: 256 characters.",
			},

			"remove_video": {
				Optional: true,
				Type:     schema.TypeInt,
				Description: "Whether to remove video data. Valid values:\n" +
					"- 0: retain\n" +
					"- 1: remove\n" +
					"Default value: 0.",
			},

			"remove_audio": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to remove audio data. Valid values:0: retain 1: remove Default value: 0.",
			},

			"video_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video stream configuration parameter. This field is required when `RemoveVideo` is 0.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The video codec. Valid values:libx264: H.264; libx265: H.265; av1: AOMedia Video 1; H.266: H.266. The AOMedia Video 1 and H.266 codecs can only be used for MP4 files. Only CRF is supported for H.266 currently.",
						},
						"fps": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video frame rate in Hz. Value range: [0,100].If the value is 0, the frame rate will be the same as that of the source video.",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Bitrate of video stream in Kbps. Value range: 0 and [128, 35,000].If the value is 0, the bitrate of the video will be the same as that of the source video.",
						},
						"resolution_adaptive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resolution adaption. Valid values:open: enabled. In this case, `Width` represents the long side of a video, while `Height` the short side;close: disabled. In this case, `Width` represents the width of a video, while `Height` the height.Default value: open.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"width": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum video width (or long side) in pixels. Value range: 0 and [128, 8192].If both `Width` and `Height` are 0, the output resolution will be the same as that of the source video.If `Width` is 0 and `Height` is not, the video width will be proportionally scaled.If `Width` is not 0 and `Height` is, the video height will be proportionally scaled.If neither `Width` nor `Height` is 0, the specified width and height will be used.Default value: 0.",
						},
						"height": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum video height (or short side) in pixels. Value range: 0 and [128, 8192].If both `Width` and `Height` are 0, the output resolution will be the same as that of the source video.If `Width` is 0 and `Height` is not, the video width will be proportionally scaled.If `Width` is not 0 and `Height` is, the video height will be proportionally scaled.If neither `Width` nor `Height` is 0, the specified width and height will be used.Default value: 0.",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Fill type, the way of processing a screenshot when the configured aspect ratio is different from that of the source video. Valid values:stretch: stretches the video image frame by frame to fill the screen. The video image may become squashed or stretched after transcoding.black: fills the uncovered area with black color, without changing the image&#39;s aspect ratio.white: fills the uncovered area with white color, without changing the image&#39;s aspect ratio.gauss: applies Gaussian blur to the uncovered area, without changing the image&#39;s aspect ratio.Default value: black.",
						},
						"vcrf": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The video constant rate factor (CRF). Value range: 1-51.If this parameter is specified, CRF encoding will be used and the bitrate parameter will be ignored.If `Codec` is `H.266`, this parameter is required (`28` is recommended).We don't recommend using this parameter unless you have special requirements.",
						},
						"gop": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "I-frame interval in frames. Valid values: 0 and 1-100000.When this parameter is set to 0 or left empty, `Gop` will be automatically set.",
						},
						"preserve_hdr_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to output an HDR (high dynamic range) video if the source video is HDR. Valid values:ON: If the source video is HDR, output an HDR video; if not, output an SDR (standard dynamic range) video.OFF: Output an SDR video regardless of whether the source video is HDR.Default value: OFF.",
						},
						"codec_tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The codec tag. This parameter is valid only if the H.265 codec is used. Valid values:hvc1hev1Default value: hvc1.",
						},
					},
				},
			},

			"audio_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Audio stream configuration parameter. This field is required when `RemoveAudio` is 0.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The audio codec.If `Container` is `mp3`, the valid value is:`libmp3lame`If `Container` is `ogg` or `flac`, the valid value is:`flac`If `Container` is `m4a`, the valid values are:`libfdk_aac``libmp3lame``ac3`If `Container` is `mp4` or `flv`, the valid values are:`libfdk_aac` (Recommended for MP4)`libmp3lame` (Recommended for FLV)`mp2`If `Container` is `hls`, the valid value is:`libfdk_aac`If `Format` is `HLS` or `MPEG-DASH`, the valid value is:`libfdk_aac`If `Container` is `wav`, the valid value is:`pcm16`.",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Audio stream bitrate in Kbps. Value range: 0 and [26, 256].If the value is 0, the bitrate of the audio stream will be the same as that of the original audio.",
						},
						"sample_rate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The audio sample rate. Valid values:`16000` (valid only if `Codec` is `pcm16`)`32000``44100``48000`Unit: Hz.",
						},
						"audio_channel": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Audio channel system. Valid values:1: mono-channel2: dual-channel6: stereoYou cannot set the sound channel as stereo for media files in container formats for audios (FLAC, OGG, MP3, M4A).Default value: 2.",
						},
					},
				},
			},

			"tehd_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "TESHD transcoding parameter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "TESHD transcoding type. Valid values: TEHD-100, OFF (default).",
						},
						"max_video_bitrate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum bitrate, which is valid when `Type` is `TESHD`.If this parameter is left blank or 0 is entered, there will be no upper limit for bitrate.",
						},
					},
				},
			},

			"segment_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The segment type. This parameter is valid only if `Container` is `hls`. Valid values: `ts`: TS segment; `fmp4`: fMP4 segment Default: `ts`.",
			},
		},
	}
}

func resourceTencentCloudVodTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_transcode_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vod.NewCreateTranscodeTemplateRequest()
		response = vod.NewCreateTranscodeTemplateResponse()
		subAppId string
	)
	if v, ok := d.GetOk("container"); ok {
		request.Container = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sub_app_id"); ok {
		subAppId = helper.IntToStr(v.(int))
		request.SubAppId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("remove_video"); ok {
		request.RemoveVideo = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("remove_audio"); ok {
		request.RemoveAudio = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video_template"); ok {
		videoTemplateInfo := vod.VideoTemplateInfo{}
		if v, ok := dMap["codec"]; ok {
			videoTemplateInfo.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["fps"]; ok {
			videoTemplateInfo.Fps = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["bitrate"]; ok {
			videoTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
		if v, ok := dMap["fill_type"]; ok {
			videoTemplateInfo.FillType = helper.String(v.(string))
		}
		if v, ok := dMap["vcrf"]; ok && v.(int) != 0 {
			videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["gop"]; ok {
			videoTemplateInfo.Gop = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["preserve_hdr_switch"]; ok {
			videoTemplateInfo.PreserveHDRSwitch = helper.String(v.(string))
		}
		if v, ok := dMap["codec_tag"]; ok {
			videoTemplateInfo.CodecTag = helper.String(v.(string))
		}
		request.VideoTemplate = &videoTemplateInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio_template"); ok {
		audioTemplateInfo := vod.AudioTemplateInfo{}
		if v, ok := dMap["codec"]; ok {
			audioTemplateInfo.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["bitrate"]; ok {
			audioTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
		tEHDConfig := vod.TEHDConfig{}
		if v, ok := dMap["type"]; ok {
			tEHDConfig.Type = helper.String(v.(string))
		}
		if v, ok := dMap["max_video_bitrate"]; ok {
			tEHDConfig.MaxVideoBitrate = helper.IntUint64(v.(int))
		}
		request.TEHDConfig = &tEHDConfig
	}

	if v, ok := d.GetOk("segment_type"); ok {
		request.SegmentType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateTranscodeTemplate(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.NonRetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vod transcodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition := *response.Response.Definition
	d.SetId(subAppId + tccommon.FILED_SP + helper.Int64ToStr(definition))

	return resourceTencentCloudVodTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudVodTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_transcode_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("transcode template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	transcodeTemplate, err := service.DescribeVodTranscodeTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToInt64(definition))
	if err != nil {
		return err
	}

	if transcodeTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VodTranscodeTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if transcodeTemplate.Container != nil {
		_ = d.Set("container", transcodeTemplate.Container)
	}

	_ = d.Set("sub_app_id", helper.StrToInt(subAppId))

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

		if transcodeTemplate.VideoTemplate.FillType != nil {
			videoTemplateMap["fill_type"] = transcodeTemplate.VideoTemplate.FillType
		}

		if transcodeTemplate.VideoTemplate.Vcrf != nil && *transcodeTemplate.VideoTemplate.Vcrf != 0 {
			videoTemplateMap["vcrf"] = transcodeTemplate.VideoTemplate.Vcrf
		}

		if transcodeTemplate.VideoTemplate.Gop != nil {
			videoTemplateMap["gop"] = transcodeTemplate.VideoTemplate.Gop
		}

		if transcodeTemplate.VideoTemplate.PreserveHDRSwitch != nil {
			videoTemplateMap["preserve_hdr_switch"] = transcodeTemplate.VideoTemplate.PreserveHDRSwitch
		}

		if transcodeTemplate.VideoTemplate.CodecTag != nil {
			videoTemplateMap["codec_tag"] = transcodeTemplate.VideoTemplate.CodecTag
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

	if transcodeTemplate.SegmentType != nil {
		_ = d.Set("segment_type", transcodeTemplate.SegmentType)
	}

	return nil
}

func resourceTencentCloudVodTranscodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_transcode_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewModifyTranscodeTemplateRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("transcode template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	request.SubAppId = helper.StrToUint64Point(subAppId)
	request.Definition = helper.StrToInt64Point(definition)

	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("container") || d.HasChange("name") || d.HasChange("comment") || d.HasChange("remove_video") || d.HasChange("remove_audio") || d.HasChange("video_template") || d.HasChange("audio_template") || d.HasChange("tehd_config") || d.HasChange("segment_type") {
		if v, ok := d.GetOk("container"); ok {
			request.Container = helper.String(v.(string))
		}
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
		if v, ok := d.GetOkExists("remove_video"); ok {
			request.RemoveVideo = helper.IntInt64(v.(int))
		}
		if v, ok := d.GetOkExists("remove_audio"); ok {
			request.RemoveAudio = helper.IntInt64(v.(int))
		}
		if dMap, ok := helper.InterfacesHeadMap(d, "video_template"); ok {
			videoTemplateInfo := vod.VideoTemplateInfoForUpdate{}
			if v, ok := dMap["codec"]; ok {
				videoTemplateInfo.Codec = helper.String(v.(string))
			}
			if v, ok := dMap["fps"]; ok {
				videoTemplateInfo.Fps = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["bitrate"]; ok {
				videoTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
			if v, ok := dMap["fill_type"]; ok {
				videoTemplateInfo.FillType = helper.String(v.(string))
			}
			if v, ok := dMap["vcrf"]; ok {
				videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["gop"]; ok {
				videoTemplateInfo.Gop = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["preserve_hdr_switch"]; ok {
				videoTemplateInfo.PreserveHDRSwitch = helper.String(v.(string))
			}
			if v, ok := dMap["codec_tag"]; ok {
				videoTemplateInfo.CodecTag = helper.String(v.(string))
			}
			request.VideoTemplate = &videoTemplateInfo
		}
		if dMap, ok := helper.InterfacesHeadMap(d, "audio_template"); ok {
			audioTemplateInfo := vod.AudioTemplateInfoForUpdate{}
			if v, ok := dMap["codec"]; ok {
				audioTemplateInfo.Codec = helper.String(v.(string))
			}
			if v, ok := dMap["bitrate"]; ok {
				audioTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
			tEHDConfig := vod.TEHDConfigForUpdate{}
			if v, ok := dMap["type"]; ok {
				tEHDConfig.Type = helper.String(v.(string))
			}
			if v, ok := dMap["max_video_bitrate"]; ok {
				tEHDConfig.MaxVideoBitrate = helper.IntUint64(v.(int))
			}
			request.TEHDConfig = &tEHDConfig
		}
		if v, ok := d.GetOk("segment_type"); ok {
			request.SegmentType = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyTranscodeTemplate(request)
		if e != nil {
			return resource.RetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vod transcodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVodTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudVodTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_transcode_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("transcode template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	if err := service.DeleteVodTranscodeTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToInt64(definition)); err != nil {
		return err
	}

	return nil
}
