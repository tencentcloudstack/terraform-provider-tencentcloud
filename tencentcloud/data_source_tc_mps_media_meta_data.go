/*
Use this data source to query detailed information of mps media_meta_data

Example Usage

```hcl
data "tencentcloud_mps_media_meta_data" "media_meta_data" {
  input_info {
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
  meta_data {
		size = &lt;nil&gt;
		container = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		height = &lt;nil&gt;
		width = &lt;nil&gt;
		duration =
		rotate = &lt;nil&gt;
		video_stream_set {
			bitrate = &lt;nil&gt;
			height = &lt;nil&gt;
			width = &lt;nil&gt;
			codec = &lt;nil&gt;
			fps = &lt;nil&gt;
			color_primaries = &lt;nil&gt;
			color_space = &lt;nil&gt;
			color_transfer = &lt;nil&gt;
			hdr_type = &lt;nil&gt;
		}
		audio_stream_set {
			bitrate = &lt;nil&gt;
			sampling_rate = &lt;nil&gt;
			codec = &lt;nil&gt;
			channel = &lt;nil&gt;
		}
		video_duration = &lt;nil&gt;
		audio_duration = &lt;nil&gt;

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

func dataSourceTencentCloudMpsMediaMetaData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsMediaMetaDataRead,
		Schema: map[string]*schema.Schema{
			"input_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "File input information that needs to get meta information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "The type of source object, which supports COS and URL.",
						},
						"cos_input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Valid when Type is COS, this item is required, indicating media processing COS object information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the COS Bucket where the media processing object file is located.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The park to which the COS Bucket where the media processing target file resides belongs.",
									},
									"object": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Input path for media processing object files.",
									},
								},
							},
						},
						"url_input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Valid when Type is URL, this item is required, indicating media processing URL object information.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Video URL.",
									},
								},
							},
						},
					},
				},
			},

			"meta_data": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Media meta data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Description: "The uploaded media file size (when the video is HLS, the size is the sum of m3u8 and ts file sizes), unit: byte.",
						},
						"container": {
							Type:        schema.TypeString,
							Description: "Container type, such as m4a, mp4, etc.",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Description: "The sum of the average bit rate of the video stream and the average bit rate of the audio stream, unit: bps.",
						},
						"height": {
							Type:        schema.TypeInt,
							Description: "The maximum value of video stream height, unit: px.",
						},
						"width": {
							Type:        schema.TypeInt,
							Description: "The maximum value of video stream width, unit: px.",
						},
						"duration": {
							Type:        schema.TypeFloat,
							Description: "Video duration, unit: second.",
						},
						"rotate": {
							Type:        schema.TypeInt,
							Description: "The selected angle during video shooting, unit: degree.",
						},
						"video_stream_set": {
							Type:        schema.TypeList,
							Description: "Video stream information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bitrate": {
										Type:        schema.TypeInt,
										Description: "The bit rate of the video stream, unit: bps.",
									},
									"height": {
										Type:        schema.TypeInt,
										Description: "The height of the video stream, unit: px.",
									},
									"width": {
										Type:        schema.TypeInt,
										Description: "The width of the video stream, unit: px.",
									},
									"codec": {
										Type:        schema.TypeString,
										Description: "The encoding format of the video stream, such as h264.",
									},
									"fps": {
										Type:        schema.TypeInt,
										Description: "Frame rate, unit: hz.",
									},
									"color_primaries": {
										Type:        schema.TypeString,
										Description: "Color primaries.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"color_space": {
										Type:        schema.TypeString,
										Description: "Color space.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"color_transfer": {
										Type:        schema.TypeString,
										Description: "Color transfer.Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"hdr_type": {
										Type:        schema.TypeString,
										Description: "Hdr type.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"audio_stream_set": {
							Type:        schema.TypeList,
							Description: "Audio stream info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bitrate": {
										Type:        schema.TypeInt,
										Description: "The bit rate of the audio stream, unit: bps.",
									},
									"sampling_rate": {
										Type:        schema.TypeInt,
										Description: "Sampling rate of the audio stream, unit: hz.",
									},
									"codec": {
										Type:        schema.TypeString,
										Description: "The encoding format of the audio stream, such as aac.",
									},
									"channel": {
										Type:        schema.TypeInt,
										Description: "Number of audio channels.Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
						"video_duration": {
							Type:        schema.TypeFloat,
							Description: "Video duration, unit: second.",
						},
						"audio_duration": {
							Type:        schema.TypeFloat,
							Description: "Audio duration, unit: second.",
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

func dataSourceTencentCloudMpsMediaMetaDataRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_media_meta_data.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if dMap, ok := helper.InterfacesHeadMap(d, "input_info"); ok {
		mediaInputInfo := mps.MediaInputInfo{}
		if v, ok := dMap["type"]; ok {
			mediaInputInfo.Type = helper.String(v.(string))
		}
		if cosInputInfoMap, ok := helper.InterfaceToMap(dMap, "cos_input_info"); ok {
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
		if urlInputInfoMap, ok := helper.InterfaceToMap(dMap, "url_input_info"); ok {
			urlInputInfo := mps.UrlInputInfo{}
			if v, ok := urlInputInfoMap["url"]; ok {
				urlInputInfo.Url = helper.String(v.(string))
			}
			mediaInputInfo.UrlInputInfo = &urlInputInfo
		}
		paramMap["input_info"] = &mediaInputInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "meta_data"); ok {
		mediaMetaData := mps.MediaMetaData{}
		if v, ok := dMap["size"]; ok {
			mediaMetaData.Size = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["container"]; ok {
			mediaMetaData.Container = helper.String(v.(string))
		}
		if v, ok := dMap["bitrate"]; ok {
			mediaMetaData.Bitrate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["height"]; ok {
			mediaMetaData.Height = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["width"]; ok {
			mediaMetaData.Width = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["duration"]; ok {
			mediaMetaData.Duration = helper.Float64(v.(float64))
		}
		if v, ok := dMap["rotate"]; ok {
			mediaMetaData.Rotate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["video_stream_set"]; ok {
			for _, item := range v.([]interface{}) {
				videoStreamSetMap := item.(map[string]interface{})
				mediaVideoStreamItem := mps.MediaVideoStreamItem{}
				if v, ok := videoStreamSetMap["bitrate"]; ok {
					mediaVideoStreamItem.Bitrate = helper.IntInt64(v.(int))
				}
				if v, ok := videoStreamSetMap["height"]; ok {
					mediaVideoStreamItem.Height = helper.IntInt64(v.(int))
				}
				if v, ok := videoStreamSetMap["width"]; ok {
					mediaVideoStreamItem.Width = helper.IntInt64(v.(int))
				}
				if v, ok := videoStreamSetMap["codec"]; ok {
					mediaVideoStreamItem.Codec = helper.String(v.(string))
				}
				if v, ok := videoStreamSetMap["fps"]; ok {
					mediaVideoStreamItem.Fps = helper.IntInt64(v.(int))
				}
				if v, ok := videoStreamSetMap["color_primaries"]; ok {
					mediaVideoStreamItem.ColorPrimaries = helper.String(v.(string))
				}
				if v, ok := videoStreamSetMap["color_space"]; ok {
					mediaVideoStreamItem.ColorSpace = helper.String(v.(string))
				}
				if v, ok := videoStreamSetMap["color_transfer"]; ok {
					mediaVideoStreamItem.ColorTransfer = helper.String(v.(string))
				}
				if v, ok := videoStreamSetMap["hdr_type"]; ok {
					mediaVideoStreamItem.HdrType = helper.String(v.(string))
				}
				mediaMetaData.VideoStreamSet = append(mediaMetaData.VideoStreamSet, &mediaVideoStreamItem)
			}
		}
		if v, ok := dMap["audio_stream_set"]; ok {
			for _, item := range v.([]interface{}) {
				audioStreamSetMap := item.(map[string]interface{})
				mediaAudioStreamItem := mps.MediaAudioStreamItem{}
				if v, ok := audioStreamSetMap["bitrate"]; ok {
					mediaAudioStreamItem.Bitrate = helper.IntInt64(v.(int))
				}
				if v, ok := audioStreamSetMap["sampling_rate"]; ok {
					mediaAudioStreamItem.SamplingRate = helper.IntInt64(v.(int))
				}
				if v, ok := audioStreamSetMap["codec"]; ok {
					mediaAudioStreamItem.Codec = helper.String(v.(string))
				}
				if v, ok := audioStreamSetMap["channel"]; ok {
					mediaAudioStreamItem.Channel = helper.IntInt64(v.(int))
				}
				mediaMetaData.AudioStreamSet = append(mediaMetaData.AudioStreamSet, &mediaAudioStreamItem)
			}
		}
		if v, ok := dMap["video_duration"]; ok {
			mediaMetaData.VideoDuration = helper.Float64(v.(float64))
		}
		if v, ok := dMap["audio_duration"]; ok {
			mediaMetaData.AudioDuration = helper.Float64(v.(float64))
		}
		paramMap["meta_data"] = &mediaMetaData
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var metaData []*mps.MediaMetaData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsMediaMetaDataByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		metaData = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(metaData))
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), mediaMetaDataMap); e != nil {
			return e
		}
	}
	return nil
}
