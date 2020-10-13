/*
Provide a resource to create a VOD adaptive dynamic streaming template.

Example Usage

```hcl
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec   = "libx264"
      fps     = 3
      bitrate = 128
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 128
      sample_rate = 32000
    }
    remove_audio = true
  }
  stream_info {
    video {
      codec   = "libx264"
      fps     = 4
      bitrate = 256
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 256
      sample_rate = 44100
    }
    remove_audio = true
  }
}
```

Import

Vod adaptive dynamic streaming template can be imported using the id, e.g.

```
$ terraform import tencentcloud_vod_adaptive_dynamic_streaming_template.foo 169141
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplate() *schema.Resource {
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
				ValidateFunc: validateStringLengthInRange(1, 64),
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
				Description: "Whether to prohibit transcoding video from low bitrate to high bitrate. Valid values: `false`: no, `true`: yes. Default value: `false`.",
			},
			"disable_higher_video_resolution": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to prohibit transcoding from low resolution to high resolution. Valid values: `false`: no, `true`: yes. Default value: `false`.",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 256),
				Description:  "Template description. Length limit: 256 characters.",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.",
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
										Description: "Video stream encoder. Valid values: `libx264`: H.264, `libx265`: H.265, `av1`: AOMedia Video 1. Currently, a resolution within 640x480 must be specified for `H.265`. and the `av1` container only supports mp4.",
									},
									"fps": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validateIntegerInRange(0, 60),
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
										Description: "Resolution adaption. Valid values: `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height. Default value: `true`. Note: this field may return null, indicating that no valid values can be obtained.",
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
										ValidateFunc: validateAllowedStringValue([]string{"stretch", "black"}),
										Description:  "Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. Default value: black. Note: this field may return null, indicating that no valid values can be obtained.",
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
										Description: "Audio stream encoder. Valid value are: `libfdk_aac` and `libmp3lame`, while `libfdk_aac` is recommended.",
									},
									"bitrate": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Audio stream bitrate in Kbps. Value range: `0` and `[26, 256]`. If the value is `0`, the bitrate of the audio stream will be the same as that of the original audio.",
									},
									"sample_rate": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Audio stream sample rate. Valid values: `32000`, `44100`, `48000`, in Hz.",
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
	defer logElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.create")()

	var (
		logId   = getLogId(contextNil)
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
	if v, ok := d.GetOk("sub_app_id"); ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}
	streamInfos := d.Get("stream_info").([]interface{})
	request.StreamInfos = make([]*vod.AdaptiveStreamTemplate, 0, len(streamInfos))
	for _, item := range streamInfos {
		v := item.(map[string]interface{})
		video := v["video"].([]interface{})[0].(map[string]interface{})
		audio := v["audio"].([]interface{})[0].(map[string]interface{})
		rAudio := REMOVE_AUDIO_TO_UNINT[v["remove_audio"].(bool)]
		request.StreamInfos = append(request.StreamInfos, &vod.AdaptiveStreamTemplate{
			Video: &vod.VideoTemplateInfo{
				Codec:              helper.String(video["codec"].(string)),
				Fps:                helper.IntUint64(video["fps"].(int)),
				Bitrate:            helper.IntUint64(video["bitrate"].(int)),
				ResolutionAdaptive: helper.String(RESOLUTION_ADAPTIVE_TO_STRING[video["resolution_adaptive"].(bool)]),
				Width:              helper.IntUint64(video["width"].(int)),
				Height:             helper.IntUint64(video["height"].(int)),
				FillType:           helper.String(video["fill_type"].(string)),
			},
			Audio: &vod.AudioTemplateInfo{
				Codec:        helper.String(audio["codec"].(string)),
				Bitrate:      helper.IntUint64(audio["bitrate"].(int)),
				SampleRate:   helper.IntUint64(audio["sample_rate"].(int)),
				AudioChannel: helper.Int64(VOD_AUDIO_CHANNEL_TYPE_TO_INT[audio["audio_channel"].(string)]),
			},
			RemoveAudio: &rAudio,
		})
	}

	var response *vod.CreateAdaptiveDynamicStreamingTemplateResponse
	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().CreateAdaptiveDynamicStreamingTemplate(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		return fmt.Errorf("for vod adaptive dynamic streaming template creation, response is nil")
	}
	d.SetId(strconv.FormatUint(*response.Response.Definition, 10))

	return resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead(d, meta)
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		id         = d.Id()
		client     = meta.(*TencentCloudClient).apiV3Conn
		vodService = VodService{client: client}
	)
	// waiting for refreshing cache
	time.Sleep(30 * time.Second)
	template, has, err := vodService.DescribeAdaptiveDynamicStreamingTemplatesById(ctx, id)
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
		})
	}
	_ = d.Set("stream_info", streamInfos)

	return nil
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.update")()

	var (
		logId      = getLogId(contextNil)
		request    = vod.NewModifyAdaptiveDynamicStreamingTemplateRequest()
		id         = d.Id()
		changeFlag = false
	)

	idUint, _ := strconv.ParseUint(id, 0, 64)
	request.Definition = &idUint
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
	if d.HasChange("sub_app_id") {
		changeFlag = true
		request.SubAppId = helper.IntUint64(d.Get("sub_app_id").(int))
	}
	if d.HasChange("stream_info") {
		changeFlag = true
		streamInfos := d.Get("stream_info").([]interface{})
		request.StreamInfos = make([]*vod.AdaptiveStreamTemplate, 0, len(streamInfos))
		for _, item := range streamInfos {
			v := item.(map[string]interface{})
			video := v["video"].([]interface{})[0].(map[string]interface{})
			audio := v["audio"].([]interface{})[0].(map[string]interface{})
			rAudio := REMOVE_AUDIO_TO_UNINT[v["remove_audio"].(bool)]
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
				},
				Audio: &vod.AudioTemplateInfo{
					Codec:        helper.String(audio["codec"].(string)),
					Bitrate:      helper.IntUint64(audio["bitrate"].(int)),
					SampleRate:   helper.IntUint64(audio["sample_rate"].(int)),
					AudioChannel: helper.Int64(VOD_AUDIO_CHANNEL_TYPE_TO_INT[audio["audio_channel"].(string)]),
				},
				RemoveAudio: &rAudio,
			})
		}
	}

	if changeFlag {
		var err error
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ModifyAdaptiveDynamicStreamingTemplate(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return retryError(err)
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
	defer logElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	vodService := VodService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if err := vodService.DeleteAdaptiveDynamicStreamingTemplate(ctx, id, uint64(d.Get("sub_app_id").(int))); err != nil {
		return err
	}

	return nil
}
