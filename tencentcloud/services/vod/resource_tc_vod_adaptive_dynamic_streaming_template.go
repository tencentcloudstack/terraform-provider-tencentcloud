package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudVodAdaptiveDynamicStreamingTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodAdaptiveDynamicStreamingTemplateCreate,
		Read:   resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead,
		Update: resourceTencentCloudVodAdaptiveDynamicStreamingTemplateUpdate,
		Delete: resourceTencentCloudVodAdaptiveDynamicStreamingTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"format": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Adaptive bitstream format. Valid values: `HLS`.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 64),
				Description:  "Template name. Length limit: 64 characters.",
			},
			"drm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "DRM scheme type. Valid values: `SimpleAES`. If this field is an empty string, DRM will not be performed on the video.",
			},
			"disable_higher_video_bitrate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to prohibit transcoding video from low bitrate to high bitrate. Valid values: `false`,`true`. `false`: no, `true`: yes. Default value: `false`.",
			},
			"disable_higher_video_resolution": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to prohibit transcoding from low resolution to high resolution. Valid values: `false`,`true`. `false`: no, `true`: yes. Default value: `false`.",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 256),
				Description:  "Template description. Length limit: 256 characters.",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.",
			},
			"stream_info": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of AdaptiveStreamTemplate parameter information of output substream for adaptive bitrate streaming. Up to 10 substreams can be output. Note: the frame rate of all substreams must be the same; otherwise, the frame rate of the first substream will be used as the output frame rate.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"video": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							MinItems:    1,
							Description: "Video parameter information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Video stream encoder. Valid values: `libx264`,`libx265`,`av1`. `libx264`: H.264, `libx265`: H.265, `av1`: AOMedia Video 1. Currently, a resolution within 640x480 must be specified for `H.265`. and the `av1` container only supports mp4.",
									},
									"fps": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: tccommon.ValidateIntegerInRange(0, 60),
										Description:  "Video frame rate in Hz. Value range: `[0, 60]`. If the value is `0`, the frame rate will be the same as that of the source video.",
									},
									"bitrate": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Bitrate of video stream in Kbps. Value range: `0` and `[128, 35000]`. If the value is `0`, the bitrate of the video will be the same as that of the source video.",
									},
									"resolution_adaptive": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Resolution adaption. Valid values: `true`,`false`. `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height. Default value: `true`. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"width": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: "Maximum value of the width (or long side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"height": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: "Maximum value of the height (or short side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"fill_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "black",
										ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"stretch", "black"}),
										Description:  "Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. Default value: black. Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"vcrf": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										Description: "Video constant bit rate control factor, value range is [1,51].\n" +
											"Note:\n" +
											"- If this parameter is specified, the bitrate control method of CRF will be used for transcoding (the video bitrate will no longer take effect);\n" +
											"- This field is required when the video stream encoding format is H.266. The recommended value is 28;\n" +
											"- If there are no special requirements, it is not recommended to specify this parameter.",
									},
									"gop": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Interval between Keyframe I frames, value range: 0 and [1, 100000], unit: number of frames. When you fill in 0 or leave it empty, the gop length is automatically set.",
									},
									"preserve_hdr_switch": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										Description: "Whether the transcoding output still maintains HDR when the original video is HDR (High Dynamic Range). Value range:\n" +
											"- ON: if the original file is HDR, the transcoding output remains HDR;, otherwise the transcoding output is SDR (Standard Dynamic Range);\n" +
											"- OFF: regardless of whether the original file is HDR or SDR, the transcoding output is SDR;\n" +
											"Default value: OFF.",
									},
									"codec_tag": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										Description: "Encoding label, valid only if the encoding format of the video stream is H.265 encoding. Available values:\n" +
											"- hvc1: stands for hvc1 tag;\n" +
											"- hev1: stands for the hev1 tag;\n" +
											"Default value: hvc1.",
									},
								},
							},
						},
						"audio": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							MinItems:    1,
							Description: "Audio parameter information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Audio stream encoder. Valid value are: `libfdk_aac` and `libmp3lame`. while `libfdk_aac` is recommended.",
									},
									"bitrate": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Audio stream bitrate in Kbps. Value range: `0` and `[26, 256]`. If the value is `0`, the bitrate of the audio stream will be the same as that of the original audio.",
									},
									"sample_rate": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Audio stream sample rate. Valid values: `32000`, `44100`, `48000`Hz.",
									},
									"audio_channel": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     VOD_AUDIO_CHANNEL_DUAL,
										Description: fmt.Sprintf("Audio channel system. Valid values: %s, %s, %s. Default value: %s.", VOD_AUDIO_CHANNEL_MONO, VOD_AUDIO_CHANNEL_DUAL, VOD_AUDIO_CHANNEL_STEREO, VOD_AUDIO_CHANNEL_DUAL),
									},
								},
							},
						},
						"remove_audio": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to remove audio stream. Valid values: `false`: no, `true`: yes. `false` by default.",
						},
						"remove_video": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to remove video stream. Valid values: `false`: no, `true`: yes. `false` by default.",
						},
						"tehd_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							MinItems:    1,
							Description: "Extremely fast HD transcoding parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Extreme high-speed HD type, available values:\n" +
											"- TEHD-100: super high definition-100th;\n" +
											"- OFF: turn off Ultra High definition.",
									},
									"max_video_bitrate": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Video bitrate limit, which is valid when Type specifies extreme speed HD type. If you leave it empty or enter 0, there is no video bitrate limit.",
									},
								},
							},
						},
					},
				},
			},
			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of template in ISO date format.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of template in ISO date format.",
			},
		},
	}
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.create")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vod.NewCreateAdaptiveDynamicStreamingTemplateRequest()
	)

	request.Format = helper.String(d.Get("format").(string))
	request.Name = helper.String(d.Get("name").(string))
	if v, ok := d.GetOk("drm_type"); ok {
		request.DrmType = helper.String(v.(string))
	}
	request.DisableHigherVideoBitrate = helper.Uint64(DISABLE_HIGHER_VIDEO_BITRATE_TO_UNINT[d.Get("disable_higher_video_bitrate").(bool)])
	request.DisableHigherVideoResolution = helper.Uint64(DISABLE_HIGHER_VIDEO_RESOLUTION_TO_UNINT[d.Get("disable_higher_video_resolution").(bool)])
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}
	var resourceId string
	if v, ok := d.GetOk("sub_app_id"); ok {
		subAppId := v.(int)
		resourceId += helper.IntToStr(subAppId)
		resourceId += tccommon.FILED_SP
		request.SubAppId = helper.IntUint64(subAppId)
	}
	streamInfos := d.Get("stream_info").([]interface{})
	request.StreamInfos = make([]*vod.AdaptiveStreamTemplate, 0, len(streamInfos))
	for _, item := range streamInfos {
		v := item.(map[string]interface{})
		video := v["video"].([]interface{})[0].(map[string]interface{})
		audio := v["audio"].([]interface{})[0].(map[string]interface{})
		rAudio := REMOVE_AUDIO_TO_UNINT[v["remove_audio"].(bool)]
		videoTemplateInfo := &vod.VideoTemplateInfo{
			Codec:              helper.String(video["codec"].(string)),
			Fps:                helper.IntUint64(video["fps"].(int)),
			Bitrate:            helper.IntUint64(video["bitrate"].(int)),
			ResolutionAdaptive: helper.String(RESOLUTION_ADAPTIVE_TO_STRING[video["resolution_adaptive"].(bool)]),
			Width:              helper.IntUint64(video["width"].(int)),
			Height:             helper.IntUint64(video["height"].(int)),
			FillType:           helper.String(video["fill_type"].(string)),
		}
		var rVideo uint64
		if v, ok := video["remove_video"]; ok && v.(bool) {
			rVideo = REMOVE_AUDIO_TO_UNINT[v.(bool)]
		}
		if v, ok := video["vcrf"]; ok && v.(int) != 0 {
			videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
		}
		if v, ok := video["gop"]; ok {
			videoTemplateInfo.Gop = helper.IntUint64(v.(int))
		}
		if v, ok := video["preserve_hdr_switch"]; ok && v.(string) != "" {
			videoTemplateInfo.PreserveHDRSwitch = helper.String(v.(string))
		}
		if v, ok := video["codec_tag"]; ok && v.(string) != "" {
			videoTemplateInfo.CodecTag = helper.String(v.(string))
		}

		var tehdConfig map[string]interface{}
		if len(v["tehd_config"].([]interface{})) > 0 {
			tehdConfig = v["tehd_config"].([]interface{})[0].(map[string]interface{})
		}
		request.StreamInfos = append(request.StreamInfos, &vod.AdaptiveStreamTemplate{

			Video: videoTemplateInfo,
			Audio: &vod.AudioTemplateInfo{
				Codec:        helper.String(audio["codec"].(string)),
				Bitrate:      helper.IntUint64(audio["bitrate"].(int)),
				SampleRate:   helper.IntUint64(audio["sample_rate"].(int)),
				AudioChannel: helper.Int64(VOD_AUDIO_CHANNEL_TYPE_TO_INT[audio["audio_channel"].(string)]),
			},
			RemoveAudio: &rAudio,
			RemoveVideo: &rVideo,
			TEHDConfig: func() *vod.TEHDConfig {
				if tehdConfig == nil {
					return nil
				}
				tehd := &vod.TEHDConfig{
					Type: helper.String(tehdConfig["type"].(string)),
				}
				if v, ok := tehdConfig["max_video_bitrate"]; ok {
					tehd.MaxVideoBitrate = helper.IntUint64(v.(int))
				}
				return tehd
			}(),
		})
	}

	var response *vod.CreateAdaptiveDynamicStreamingTemplateResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateAdaptiveDynamicStreamingTemplate(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(err)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), strconv.ErrRange.Error())
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		return fmt.Errorf("for vod adaptive dynamic streaming template creation, response is nil")
	}
	resourceId += strconv.FormatUint(*response.Response.Definition, 10)
	d.SetId(resourceId)

	return resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead(d, meta)
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		subAppId   int
		definition string
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		vodService = VodService{client: client}
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
	} else {
		definition = d.Id()
	}
	// waiting for refreshing cache
	time.Sleep(30 * time.Second)
	template, has, err := vodService.DescribeAdaptiveDynamicStreamingTemplatesById(ctx, definition, subAppId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("format", template.Format)
	_ = d.Set("name", template.Name)
	_ = d.Set("drm_type", template.DrmType)
	_ = d.Set("disable_higher_video_bitrate", *template.DisableHigherVideoBitrate == 1)
	_ = d.Set("disable_higher_video_resolution", *template.DisableHigherVideoResolution == 1)
	_ = d.Set("comment", template.Comment)
	_ = d.Set("create_time", template.CreateTime)
	_ = d.Set("update_time", template.UpdateTime)

	var streamInfos = make([]interface{}, 0, len(template.StreamInfos))
	for _, v := range template.StreamInfos {
		streamInfos = append(streamInfos, map[string]interface{}{
			"video": []map[string]interface{}{
				{
					"codec":               v.Video.Codec,
					"fps":                 v.Video.Fps,
					"bitrate":             v.Video.Bitrate,
					"resolution_adaptive": *v.Video.ResolutionAdaptive == "open",
					"width":               v.Video.Width,
					"height":              v.Video.Height,
					"fill_type":           v.Video.FillType,
					"vcrf":                v.Video.Vcrf,
					"gop":                 v.Video.Gop,
					"preserve_hdr_switch": v.Video.PreserveHDRSwitch,
					"codec_tag":           v.Video.CodecTag,
				},
			},
			"audio": []map[string]interface{}{
				{
					"codec":         v.Audio.Codec,
					"bitrate":       v.Audio.Bitrate,
					"sample_rate":   v.Audio.SampleRate,
					"audio_channel": VOD_AUDIO_CHANNEL_TYPE_TO_STRING[*v.Audio.AudioChannel],
				},
			},
			"remove_audio": *v.RemoveAudio == 1,
			"remove_video": *v.RemoveVideo == 1,
			"tehd_config": func() []map[string]interface{} {
				if v.TEHDConfig == nil {
					return nil
				}
				return []map[string]interface{}{
					{
						"type":              v.TEHDConfig.Type,
						"max_video_bitrate": v.TEHDConfig.MaxVideoBitrate,
					},
				}
			}(),
		})
	}
	_ = d.Set("stream_info", streamInfos)
	if subAppId != 0 {
		_ = d.Set("sub_app_id", subAppId)
	}

	return nil
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = vod.NewModifyAdaptiveDynamicStreamingTemplateRequest()
		changeFlag = false
		subAppId   int
		definition string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
		request.SubAppId = helper.IntUint64(subAppId)
	} else {
		definition = d.Id()
		if v, ok := d.GetOk("sub_app_id"); ok {
			request.SubAppId = helper.IntUint64(v.(int))
		}
	}

	request.Definition = helper.StrToUint64Point(definition)

	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("format") {
		changeFlag = true
		request.Format = helper.String(d.Get("format").(string))
	}
	if d.HasChange("name") {
		changeFlag = true
		request.Name = helper.String(d.Get("name").(string))
	}
	if d.HasChange("disable_higher_video_bitrate") {
		changeFlag = true
		request.DisableHigherVideoBitrate = helper.Uint64(DISABLE_HIGHER_VIDEO_BITRATE_TO_UNINT[d.Get("disable_higher_video_bitrate").(bool)])
	}
	if d.HasChange("disable_higher_video_resolution") {
		changeFlag = true
		request.DisableHigherVideoResolution = helper.Uint64(DISABLE_HIGHER_VIDEO_RESOLUTION_TO_UNINT[d.Get("disable_higher_video_resolution").(bool)])
	}
	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("stream_info") {
		changeFlag = true
		streamInfos := d.Get("stream_info").([]interface{})
		request.StreamInfos = make([]*vod.AdaptiveStreamTemplate, 0, len(streamInfos))
		for _, item := range streamInfos {
			v := item.(map[string]interface{})
			video := v["video"].([]interface{})[0].(map[string]interface{})
			audio := v["audio"].([]interface{})[0].(map[string]interface{})
			var tehdConfig map[string]interface{}
			if len(v["tehd_config"].([]interface{})) > 0 {
				tehdConfig = v["tehd_config"].([]interface{})[0].(map[string]interface{})
			}
			rAudio := REMOVE_AUDIO_TO_UNINT[v["remove_audio"].(bool)]
			var rVideo uint64
			if v, ok := video["remove_video"]; ok && v.(bool) {
				rVideo = REMOVE_AUDIO_TO_UNINT[v.(bool)]
			}
			request.StreamInfos = append(request.StreamInfos, &vod.AdaptiveStreamTemplate{
				Video: &vod.VideoTemplateInfo{
					Codec:              helper.String(video["codec"].(string)),
					Fps:                helper.IntUint64(video["fps"].(int)),
					Bitrate:            helper.IntUint64(video["bitrate"].(int)),
					ResolutionAdaptive: helper.String(RESOLUTION_ADAPTIVE_TO_STRING[video["resolution_adaptive"].(bool)]),
					Width: func(width int) *uint64 {
						if width == 0 {
							return nil
						}
						return helper.IntUint64(width)
					}(video["width"].(int)),
					Height: func(height int) *uint64 {
						if height == 0 {
							return nil
						}
						return helper.IntUint64(height)
					}(video["height"].(int)),
					FillType: helper.String(video["fill_type"].(string)),
					Vcrf: func() *uint64 {
						if v, ok := video["vcrf"]; !ok || v.(int) == 0 {
							return nil
						}
						return helper.IntUint64(video["vcrf"].(int))
					}(),
					Gop: func() *uint64 {
						if _, ok := video["gop"]; !ok {
							return nil
						}
						return helper.IntUint64(video["gop"].(int))
					}(),
					PreserveHDRSwitch: func() *string {
						if v, ok := video["preserve_hdr_switch"]; !ok || v.(string) == "" {
							return nil
						}
						return helper.String(video["preserve_hdr_switch"].(string))
					}(),
					CodecTag: func() *string {
						if v, ok := video["codec_tag"]; !ok || v.(string) == "" {
							return nil
						}
						return helper.String(video["codec_tag"].(string))
					}(),
				},
				Audio: &vod.AudioTemplateInfo{
					Codec:        helper.String(audio["codec"].(string)),
					Bitrate:      helper.IntUint64(audio["bitrate"].(int)),
					SampleRate:   helper.IntUint64(audio["sample_rate"].(int)),
					AudioChannel: helper.Int64(VOD_AUDIO_CHANNEL_TYPE_TO_INT[audio["audio_channel"].(string)]),
				},
				RemoveAudio: &rAudio,
				RemoveVideo: &rVideo,
				TEHDConfig: func() *vod.TEHDConfig {
					if tehdConfig == nil {
						return nil
					}
					tehd := &vod.TEHDConfig{
						Type: helper.String(tehdConfig["type"].(string)),
					}
					if v, ok := tehdConfig["max_video_bitrate"]; ok {
						tehd.MaxVideoBitrate = helper.IntUint64(v.(int))
					}
					return tehd
				}(),
			})
		}
	}

	if changeFlag {
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyAdaptiveDynamicStreamingTemplate(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		subAppId   int
		definition string
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
	} else {
		definition = d.Id()
		if v, ok := d.GetOk("sub_app_id"); ok {
			subAppId = v.(int)
		}
	}
	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := vodService.DeleteAdaptiveDynamicStreamingTemplate(ctx, definition, uint64(subAppId)); err != nil {
		return err
	}

	return nil
}
