package tencentcloud

import (
	"context"
	"encoding/json"

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
				Description: "Input information of file for metadata getting.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The input type. Valid values:`COS`: A COS bucket address.`URL`: A URL.`AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
						},
						"cos_input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The region of the COS bucket, such as `ap-chongqing`.",
									},
									"object": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
									},
								},
							},
						},
						"url_input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "URL of a video.",
									},
								},
							},
						},
						"s3_input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"s3_bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The AWS S3 bucket.",
									},
									"s3_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The region of the AWS S3 bucket.",
									},
									"s3_object": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The path of the AWS S3 object.",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key ID required to access the AWS S3 object.",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key required to access the AWS S3 object.",
									},
								},
							},
						},
					},
				},
			},

			"meta_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Media metadata.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of an uploaded media file in bytes (which is the sum of size of m3u8 and ts files if the video is in HLS format).Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"container": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Container, such as m4a and mp4.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sum of the average bitrate of a video stream and that of an audio stream in bps.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"height": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum value of the height of a video stream in px.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum value of the width of a video stream in px.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"duration": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Video duration in seconds.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"rotate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Selected angle during video recording in degrees.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"video_stream_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Video stream information.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bitrate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Bitrate of a video stream in bps.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"height": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Height of a video stream in px.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"width": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Width of a video stream in px.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"codec": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Video stream codec, such as h264.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"fps": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Frame rate in Hz.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"color_primaries": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Color primariesNote: this field may return `null`, indicating that no valid value was found.",
									},
									"color_space": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Color spaceNote: this field may return `null`, indicating that no valid value was found.",
									},
									"color_transfer": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Color transferNote: this field may return `null`, indicating that no valid value was found.",
									},
									"hdr_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HDR typeNote: This field may return `null`, indicating that no valid value was found.",
									},
								},
							},
						},
						"audio_stream_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Audio stream information.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bitrate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Bitrate of an audio stream in bps.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"sampling_rate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Sample rate of an audio stream in Hz.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"codec": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Audio stream codec, such as aac.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"channel": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of sound channels, e.g., 2Note: this field may return `null`, indicating that no valid value was found.",
									},
								},
							},
						},
						"video_duration": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Video duration in seconds.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"audio_duration": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Audio duration in seconds.Note: This field may return null, indicating that no valid values can be obtained.",
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
	mediaInputInfo := mps.MediaInputInfo{}

	paramMap := make(map[string]interface{})
	if dMap, ok := helper.InterfacesHeadMap(d, "input_info"); ok {
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
		if s3InputInfoMap, ok := helper.InterfaceToMap(dMap, "s3_input_info"); ok {
			s3InputInfo := mps.S3InputInfo{}
			if v, ok := s3InputInfoMap["s3_bucket"]; ok {
				s3InputInfo.S3Bucket = helper.String(v.(string))
			}
			if v, ok := s3InputInfoMap["s3_region"]; ok {
				s3InputInfo.S3Region = helper.String(v.(string))
			}
			if v, ok := s3InputInfoMap["s3_object"]; ok {
				s3InputInfo.S3Object = helper.String(v.(string))
			}
			if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
				s3InputInfo.S3SecretId = helper.String(v.(string))
			}
			if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
				s3InputInfo.S3SecretKey = helper.String(v.(string))
			}
			mediaInputInfo.S3InputInfo = &s3InputInfo
		}
		paramMap["InputInfo"] = &mediaInputInfo
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var metaData *mps.MediaMetaData

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

	mediaMetaDataMap := map[string]interface{}{}
	if metaData != nil {
		if metaData.Size != nil {
			mediaMetaDataMap["size"] = metaData.Size
		}

		if metaData.Container != nil {
			mediaMetaDataMap["container"] = metaData.Container
		}

		if metaData.Bitrate != nil {
			mediaMetaDataMap["bitrate"] = metaData.Bitrate
		}

		if metaData.Height != nil {
			mediaMetaDataMap["height"] = metaData.Height
		}

		if metaData.Width != nil {
			mediaMetaDataMap["width"] = metaData.Width
		}

		if metaData.Duration != nil {
			mediaMetaDataMap["duration"] = metaData.Duration
		}

		if metaData.Rotate != nil {
			mediaMetaDataMap["rotate"] = metaData.Rotate
		}

		if metaData.VideoStreamSet != nil {
			videoStreamSetList := []interface{}{}
			for _, videoStreamSet := range metaData.VideoStreamSet {
				videoStreamSetMap := map[string]interface{}{}

				if videoStreamSet.Bitrate != nil {
					videoStreamSetMap["bitrate"] = videoStreamSet.Bitrate
				}

				if videoStreamSet.Height != nil {
					videoStreamSetMap["height"] = videoStreamSet.Height
				}

				if videoStreamSet.Width != nil {
					videoStreamSetMap["width"] = videoStreamSet.Width
				}

				if videoStreamSet.Codec != nil {
					videoStreamSetMap["codec"] = videoStreamSet.Codec
				}

				if videoStreamSet.Fps != nil {
					videoStreamSetMap["fps"] = videoStreamSet.Fps
				}

				if videoStreamSet.ColorPrimaries != nil {
					videoStreamSetMap["color_primaries"] = videoStreamSet.ColorPrimaries
				}

				if videoStreamSet.ColorSpace != nil {
					videoStreamSetMap["color_space"] = videoStreamSet.ColorSpace
				}

				if videoStreamSet.ColorTransfer != nil {
					videoStreamSetMap["color_transfer"] = videoStreamSet.ColorTransfer
				}

				if videoStreamSet.HdrType != nil {
					videoStreamSetMap["hdr_type"] = videoStreamSet.HdrType
				}

				videoStreamSetList = append(videoStreamSetList, videoStreamSetMap)
			}

			mediaMetaDataMap["video_stream_set"] = videoStreamSetList
		}

		if metaData.AudioStreamSet != nil {
			audioStreamSetList := []interface{}{}
			for _, audioStreamSet := range metaData.AudioStreamSet {
				audioStreamSetMap := map[string]interface{}{}

				if audioStreamSet.Bitrate != nil {
					audioStreamSetMap["bitrate"] = audioStreamSet.Bitrate
				}

				if audioStreamSet.SamplingRate != nil {
					audioStreamSetMap["sampling_rate"] = audioStreamSet.SamplingRate
				}

				if audioStreamSet.Codec != nil {
					audioStreamSetMap["codec"] = audioStreamSet.Codec
				}

				if audioStreamSet.Channel != nil {
					audioStreamSetMap["channel"] = audioStreamSet.Channel
				}

				audioStreamSetList = append(audioStreamSetList, audioStreamSetMap)
			}

			mediaMetaDataMap["audio_stream_set"] = audioStreamSetList
		}

		if metaData.VideoDuration != nil {
			mediaMetaDataMap["video_duration"] = metaData.VideoDuration
		}

		if metaData.AudioDuration != nil {
			mediaMetaDataMap["audio_duration"] = metaData.AudioDuration
		}

		_ = d.Set("meta_data", []interface{}{mediaMetaDataMap})
	}

	id, _ := json.Marshal(mediaInputInfo)
	d.SetId(helper.DataResourceIdHash(string(id)))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), mediaMetaDataMap); e != nil {
			return e
		}
	}
	return nil
}
