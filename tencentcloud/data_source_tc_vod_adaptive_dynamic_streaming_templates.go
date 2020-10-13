/*
Use this data source to query detailed information of VOD adaptive dynamic streaming templates.

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

data "tencentcloud_vod_adaptive_dynamic_streaming_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVodAdaptiveDynamicStreamingTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodAdaptiveDynamicStreamingTemplatesRead,

		Schema: map[string]*schema.Schema{
			"definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique ID filter of adaptive dynamic streaming template.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template type filter. Valid values: `Preset`: preset template; `Custom`: custom template.",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"template_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of adaptive dynamic streaming templates. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of adaptive dynamic streaming template.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template type filter. Valid values: `Preset`: preset template; `Custom`: custom template.",
						},
						"format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Adaptive bitstream format.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template name.",
						},
						"drm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DRM scheme type.",
						},
						"disable_higher_video_bitrate": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to prohibit transcoding video from low bitrate to high bitrate. `false`: no, `true`: yes.",
						},
						"disable_higher_video_resolution": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to prohibit transcoding from low resolution to high resolution. `false`: no, `true`: yes.",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template description.",
						},
						"stream_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of AdaptiveStreamTemplate parameter information of output substream for adaptive bitrate streaming.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"video": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Video parameter information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"codec": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Video stream encoder. Valid values: `libx264`: H.264, `libx265`: H.265, `av1`: AOMedia Video 1. Currently, a resolution within 640x480 must be specified for `H.265`. and the `av1` container only supports mp4.",
												},
												"fps": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Video frame rate in Hz. Value range: `[0, 60]`. If the value is `0`, the frame rate will be the same as that of the source video.",
												},
												"bitrate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Bitrate of video stream in Kbps. Value range: `0` and `[128, 35000]`. If the value is `0`, the bitrate of the video will be the same as that of the source video.",
												},
												"resolution_adaptive": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Resolution adaption. Valid values: `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height. Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"width": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum value of the width (or long side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"height": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum value of the height (or short side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"fill_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. Note: this field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"audio": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Audio parameter information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"codec": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Audio stream encoder. Valid value are: `libfdk_aac` and `libmp3lame`.",
												},
												"bitrate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Audio stream bitrate in Kbps. Value range: `0` and `[26, 256]`. If the value is `0`, the bitrate of the audio stream will be the same as that of the original audio.",
												},
												"sample_rate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Audio stream sample rate. Valid values: `32000`, `44100`, `48000`, in Hz.",
												},
												"audio_channel": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: fmt.Sprintf("Audio channel system. Valid values: %s, %s, %s.", VOD_AUDIO_CHANNEL_MONO, VOD_AUDIO_CHANNEL_DUAL, VOD_AUDIO_CHANNEL_STEREO),
												},
											},
										},
									},
									"remove_audio": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to remove audio stream. `false`: no, `true`: yes.",
									},
								},
							},
						},
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
				},
			},
		},
	}
}

func dataSourceTencentCloudVodAdaptiveDynamicStreamingTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vod_adaptive_dynamic_streaming_templates.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	filter := make(map[string]interface{})
	if v, ok := d.GetOk("definition"); ok {
		filter["definitions"] = []string{v.(string)}
	}
	if v, ok := d.GetOk("type"); ok {
		filter["type"] = v.(string)
	}
	if v, ok := d.GetOk("sub_app_id"); ok {
		filter["sub_appid"] = v.(int)
	}

	vodService := VodService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	templates, err := vodService.DescribeAdaptiveDynamicStreamingTemplatesByFilter(ctx, filter)
	if err != nil {
		return err
	}

	templatesList := make([]map[string]interface{}, 0, len(templates))
	ids := make([]string, 0, len(templates))
	for _, item := range templates {
		templatesList = append(templatesList, func() map[string]interface{} {
			definitionStr := strconv.FormatUint(*item.Definition, 10)
			mapping := map[string]interface{}{
				"definition":                      definitionStr,
				"type":                            item.Type,
				"format":                          item.Format,
				"name":                            item.Name,
				"drm_type":                        item.DrmType,
				"disable_higher_video_bitrate":    *item.DisableHigherVideoBitrate == 1,
				"disable_higher_video_resolution": *item.DisableHigherVideoResolution == 1,
				"comment":                         item.Comment,
				"create_time":                     item.CreateTime,
				"update_time":                     item.UpdateTime,
			}
			var streamInfos = make([]interface{}, 0, len(item.StreamInfos))
			for _, v := range item.StreamInfos {
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
			mapping["stream_info"] = streamInfos
			ids = append(ids, definitionStr)
			return mapping
		}())
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("template_list", templatesList); e != nil {
		log.Printf("[CRITAL]%s provider set vod adaptive dynamic streaming template list fail, reason:%s ", logId, e.Error())
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), templatesList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
