/*
Use this data source to query detailed information of mps adaptive_dynamic_streaming_templates

Example Usage

```hcl
data "tencentcloud_mps_adaptive_dynamic_streaming_templates" "adaptive_dynamic_streaming_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  adaptive_dynamic_streaming_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		format = &lt;nil&gt;
		stream_infos {
			video {
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
			audio {
				codec = &lt;nil&gt;
				bitrate = &lt;nil&gt;
				sample_rate = &lt;nil&gt;
				audio_channel = 2
			}
			remove_audio = &lt;nil&gt;
			remove_video = &lt;nil&gt;
		}
		disable_higher_video_bitrate = &lt;nil&gt;
		disable_higher_video_resolution = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;

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

func dataSourceTencentCloudMpsAdaptiveDynamicStreamingTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsAdaptiveDynamicStreamingTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Adaptive dynamic streaming template uniquely identifies filter conditions, array length limit: 100.",
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

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template type filter condition, optional value:Preset: system preset template.Custom: user-defined template.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"adaptive_dynamic_streaming_template_set": {
				Type:        schema.TypeList,
				Description: "Adaptive dynamic streaming template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Description: "The unique identifier of the adaptive dynamic streaming template.",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Template type, optional value:Preset: system preset template.Custom: user-defined template.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Adaptive dynamic streaming template name.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "The description information of adaptive dynamic streaming template.",
						},
						"format": {
							Type:        schema.TypeString,
							Description: "Adaptive transcoding format, value range:HLS, MPEG-DASH.",
						},
						"stream_infos": {
							Type:        schema.TypeList,
							Description: "Adaptive code stream input stream parameter information, up to 10 streams can be input.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"video": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Video parameter information.",
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
													Description: "Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.Note: In adaptive mode, Width cannot be smaller than Height.",
												},
												"width": {
													Type:        schema.TypeInt,
													Description: "The maximum value of the width (or long side) of the video streaming, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
												},
												"height": {
													Type:        schema.TypeInt,
													Description: "The maximum value of the height (or short side) of the video streaming, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
												},
												"gop": {
													Type:        schema.TypeInt,
													Description: "The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.When filling 0 or not filling, the system will automatically set the gop length.",
												},
												"fill_type": {
													Type:        schema.TypeString,
													Description: "Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and use Gaussian blur for the rest of the edge.Default value: black.Note: Adaptive stream only supports stretch, black.",
												},
												"vcrf": {
													Type:        schema.TypeInt,
													Description: "Video constant bit rate control factor, the value range is [1, 51].If this parameter is specified, the code rate control method of CRF will be used for transcoding (the video code rate will no longer take effect).If there is no special requirement, it is not recommended to specify this parameter.",
												},
											},
										},
									},
									"audio": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Audio parameter information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"codec": {
													Type:        schema.TypeString,
													Description: "Encoding format of audio stream.When the outer parameter Container is mp3, the optional value is:libmp3lame.When the outer parameter Container is ogg or flac, the optional value is:flac.When the outer parameter Container is m4a, the optional value is:libfdk_aac.libmp3lame.ac3.When the outer parameter Container is mp4 or flv, the optional value is:libfdk_aac: more suitable for mp4.libmp3lame: more suitable for flv.When the outer parameter Container is hls, the optional value is:libfdk_aac.libmp3lame.",
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
									"remove_audio": {
										Type:        schema.TypeInt,
										Description: "Whether to remove audio stream, value:0: reserved.1: remove.",
									},
									"remove_video": {
										Type:        schema.TypeInt,
										Description: "Whether to remove video stream, value:0: reserved.1: remove.",
									},
								},
							},
						},
						"disable_higher_video_bitrate": {
							Type:        schema.TypeInt,
							Description: "Whether to prohibit video from low bit rate to high bit rate, value range:0: no.1: yes.",
						},
						"disable_higher_video_resolution": {
							Type:        schema.TypeInt,
							Description: "Whether to prohibit the conversion of video resolution to high resolution, value range:0 : no.1: yes.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Template creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Template last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
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

func dataSourceTencentCloudMpsAdaptiveDynamicStreamingTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_adaptive_dynamic_streaming_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("definitions"); ok {
		definitionsSet := v.(*schema.Set).List()
		for i := range definitionsSet {
			definitions := definitionsSet[i].(int)
			paramMap["Definitions"] = append(paramMap["Definitions"], helper.IntUint64(definitions))
		}
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("adaptive_dynamic_streaming_template_set"); ok {
		adaptiveDynamicStreamingTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.AdaptiveDynamicStreamingTemplate, 0, len(adaptiveDynamicStreamingTemplateSetSet))

		for _, item := range adaptiveDynamicStreamingTemplateSetSet {
			adaptiveDynamicStreamingTemplate := mps.AdaptiveDynamicStreamingTemplate{}
			adaptiveDynamicStreamingTemplateMap := item.(map[string]interface{})

			if v, ok := adaptiveDynamicStreamingTemplateMap["definition"]; ok {
				adaptiveDynamicStreamingTemplate.Definition = helper.IntUint64(v.(int))
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["type"]; ok {
				adaptiveDynamicStreamingTemplate.Type = helper.String(v.(string))
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["name"]; ok {
				adaptiveDynamicStreamingTemplate.Name = helper.String(v.(string))
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["comment"]; ok {
				adaptiveDynamicStreamingTemplate.Comment = helper.String(v.(string))
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["format"]; ok {
				adaptiveDynamicStreamingTemplate.Format = helper.String(v.(string))
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["stream_infos"]; ok {
				for _, item := range v.([]interface{}) {
					streamInfosMap := item.(map[string]interface{})
					adaptiveStreamTemplate := mps.AdaptiveStreamTemplate{}
					if videoMap, ok := helper.InterfaceToMap(streamInfosMap, "video"); ok {
						videoTemplateInfo := mps.VideoTemplateInfo{}
						if v, ok := videoMap["codec"]; ok {
							videoTemplateInfo.Codec = helper.String(v.(string))
						}
						if v, ok := videoMap["fps"]; ok {
							videoTemplateInfo.Fps = helper.IntUint64(v.(int))
						}
						if v, ok := videoMap["bitrate"]; ok {
							videoTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
					if audioMap, ok := helper.InterfaceToMap(streamInfosMap, "audio"); ok {
						audioTemplateInfo := mps.AudioTemplateInfo{}
						if v, ok := audioMap["codec"]; ok {
							audioTemplateInfo.Codec = helper.String(v.(string))
						}
						if v, ok := audioMap["bitrate"]; ok {
							audioTemplateInfo.Bitrate = helper.IntUint64(v.(int))
						}
						if v, ok := audioMap["sample_rate"]; ok {
							audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
						}
						if v, ok := audioMap["audio_channel"]; ok {
							audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
						}
						adaptiveStreamTemplate.Audio = &audioTemplateInfo
					}
					if v, ok := streamInfosMap["remove_audio"]; ok {
						adaptiveStreamTemplate.RemoveAudio = helper.IntUint64(v.(int))
					}
					if v, ok := streamInfosMap["remove_video"]; ok {
						adaptiveStreamTemplate.RemoveVideo = helper.IntUint64(v.(int))
					}
					adaptiveDynamicStreamingTemplate.StreamInfos = append(adaptiveDynamicStreamingTemplate.StreamInfos, &adaptiveStreamTemplate)
				}
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["disable_higher_video_bitrate"]; ok {
				adaptiveDynamicStreamingTemplate.DisableHigherVideoBitrate = helper.IntUint64(v.(int))
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["disable_higher_video_resolution"]; ok {
				adaptiveDynamicStreamingTemplate.DisableHigherVideoResolution = helper.IntUint64(v.(int))
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["create_time"]; ok {
				adaptiveDynamicStreamingTemplate.CreateTime = helper.String(v.(string))
			}
			if v, ok := adaptiveDynamicStreamingTemplateMap["update_time"]; ok {
				adaptiveDynamicStreamingTemplate.UpdateTime = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &adaptiveDynamicStreamingTemplate)
		}
		paramMap["adaptive_dynamic_streaming_template_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var adaptiveDynamicStreamingTemplateSet []*mps.AdaptiveDynamicStreamingTemplate

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsAdaptiveDynamicStreamingTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		adaptiveDynamicStreamingTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(adaptiveDynamicStreamingTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(adaptiveDynamicStreamingTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
