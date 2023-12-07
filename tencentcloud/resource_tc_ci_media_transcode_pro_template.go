package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCiMediaTranscodeProTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaTranscodeProTemplateCreate,
		Read:   resourceTencentCloudCiMediaTranscodeProTemplateRead,
		Update: resourceTencentCloudCiMediaTranscodeProTemplateUpdate,
		Delete: resourceTencentCloudCiMediaTranscodeProTemplateDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "bucket name.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"container": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "container format.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Package format.",
						},
						"clip_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Fragment configuration, valid when format is hls and dash.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"duration": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Fragmentation duration, default 5s.",
									},
								},
							},
						},
					},
				},
			},

			"video": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "video information, do not upload Video, which is equivalent to deleting video information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Codec format, default value: `H.264`, when format is WebM, it is VP8, value range: `H.264`, `H.265`, `VP8`, `VP9`, `AV1`.",
						},
						"profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "encoding level, Support baseline, main, high, auto- When Pixfmt is auto, this parameter can only be set to auto, when it is set to other options, the parameter value will be set to auto- baseline: suitable for mobile devices- main: suitable for standard resolution devices- high: suitable for high-resolution devices- Only H.264 supports this parameter.",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "High, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video, must be even.",
						},
						"interlaced": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "field pattern.",
						},
						"fps": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frame rate, value range: (0, 60], Unit: fps.",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bit rate of video output file, value range: [10, 50000], unit: Kbps, auto means adaptive bit rate.",
						},
						"rotate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rotation angle, Value range: [0, 360), Unit: degree.",
						},
					},
				},
			},

			"time_interval": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "time interval.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Starting time, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
						},
						"duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "duration, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
						},
					},
				},
			},

			"audio": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Audio information, do not transmit Audio, which is equivalent to deleting audio information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Codec format, value aac, mp3, flac, amr, Vorbis, opus, pcm_s16le.",
						},
						"remove": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to delete the source audio stream, the value is true, false.",
						},
					},
				},
			},

			"trans_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "transcoding configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"adj_dar_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resolution adjustment method, value scale, crop, pad, none, When the aspect ratio of the output video is different from the original video, adjust the resolution accordingly according to this parameter.",
						},
						"is_check_reso": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to check the resolution, when it is false, transcode according to the configuration parameters.",
						},
						"reso_adj_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resolution adjustment mode, value 0, 1; 0 means use the original video resolution; 1 means return transcoding failed, Take effect when IsCheckReso is true.",
						},
						"is_check_video_bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to check the video code rate, when it is false, transcode according to the configuration parameters.",
						},
						"video_bitrate_adj_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video bit rate adjustment method, value 0, 1; when the output video bit rate is greater than the original video bit rate, 0 means use the original video bit rate; 1 means return transcoding failed, Take effect when IsCheckVideoBitrate is true.",
						},
						"is_check_audio_bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to check the audio code rate, true, false, When false, transcode according to configuration parameters.",
						},
						"audio_bitrate_adj_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Audio bit rate adjustment mode, value 0, 1; when the output audio bit rate is greater than the original audio bit rate, 0 means use the original audio bit rate; 1 means return transcoding failed, Take effect when IsCheckAudioBitrate is true.",
						},
						"delete_metadata": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to delete the MetaData information in the file, true, false, When false, keep source file information.",
						},
						"is_hdr2_sdr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to enable HDR to SDR true, false.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaTranscodeProTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_pro_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaTranscodeProTemplateOptions{
			Tag: "TranscodePro",
		}
		templateId string
		bucket     string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		if clipConfigMap, ok := helper.InterfaceToMap(dMap, "clip_config"); ok {
			clipConfig := cos.ClipConfig{}
			if v, ok := clipConfigMap["duration"]; ok {
				clipConfig.Duration = v.(string)
			}
			container.ClipConfig = &clipConfig
		}
		request.Container = &container
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
		video := cos.TranscodeProVideo{}
		if v, ok := dMap["codec"]; ok {
			video.Codec = v.(string)
		}
		if v, ok := dMap["profile"]; ok {
			video.Profile = v.(string)
		}
		if v, ok := dMap["width"]; ok {
			video.Width = v.(string)
		}
		if v, ok := dMap["height"]; ok {
			video.Height = v.(string)
		}
		if v, ok := dMap["interlaced"]; ok {
			video.Interlaced = v.(string)
		}
		if v, ok := dMap["fps"]; ok {
			video.Fps = v.(string)
		}
		if v, ok := dMap["bitrate"]; ok {
			video.Bitrate = v.(string)
		}
		if v, ok := dMap["rotate"]; ok {
			video.Rotate = v.(string)
		}
		request.Video = &video
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "time_interval"); ok {
		timeInterval := cos.TimeInterval{}
		if v, ok := dMap["start"]; ok {
			timeInterval.Start = v.(string)
		}
		if v, ok := dMap["duration"]; ok {
			timeInterval.Duration = v.(string)
		}
		request.TimeInterval = &timeInterval
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio"); ok {
		audio := cos.TranscodeProAudio{}
		if v, ok := dMap["codec"]; ok {
			audio.Codec = v.(string)
		}
		if v, ok := dMap["remove"]; ok {
			audio.Remove = v.(string)
		}
		request.Audio = &audio
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "trans_config"); ok {
		transConfig := cos.TransConfig{}
		if v, ok := dMap["adj_dar_method"]; ok {
			transConfig.AdjDarMethod = v.(string)
		}
		if v, ok := dMap["is_check_reso"]; ok {
			transConfig.IsCheckReso = v.(string)
		}
		if v, ok := dMap["reso_adj_method"]; ok {
			transConfig.ResoAdjMethod = v.(string)
		}
		if v, ok := dMap["is_check_video_bitrate"]; ok {
			transConfig.IsCheckVideoBitrate = v.(string)
		}
		if v, ok := dMap["video_bitrate_adj_method"]; ok {
			transConfig.VideoBitrateAdjMethod = v.(string)
		}
		if v, ok := dMap["is_check_audio_bitrate"]; ok {
			transConfig.IsCheckAudioBitrate = v.(string)
		}
		if v, ok := dMap["audio_bitrate_adj_method"]; ok {
			transConfig.AudioBitrateAdjMethod = v.(string)
		}
		if v, ok := dMap["delete_metadata"]; ok {
			transConfig.DeleteMetadata = v.(string)
		}
		if v, ok := dMap["is_hdr2_sdr"]; ok {
			transConfig.IsHdr2Sdr = v.(string)
		}
		request.TransConfig = &transConfig
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaTranscodeProTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaTranscodeProTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeProTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaTranscodeProTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeProTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_pro_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	template, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if template.Name != "" {
		_ = d.Set("name", template.Name)
	}
	mediaTranscodeProTemplate := template.TransProTpl
	if mediaTranscodeProTemplate != nil {
		if mediaTranscodeProTemplate.Container != nil {
			containerMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.Container.Format != "" {
				containerMap["format"] = mediaTranscodeProTemplate.Container.Format
			}

			if mediaTranscodeProTemplate.Container.ClipConfig != nil {
				clipConfigMap := map[string]interface{}{}

				if mediaTranscodeProTemplate.Container.ClipConfig.Duration != "" {
					clipConfigMap["duration"] = mediaTranscodeProTemplate.Container.ClipConfig.Duration
				}

				containerMap["clip_config"] = []interface{}{clipConfigMap}
			}

			_ = d.Set("container", []interface{}{containerMap})
		}

		if mediaTranscodeProTemplate.Video != nil {
			videoMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.Video.Codec != "" {
				videoMap["codec"] = mediaTranscodeProTemplate.Video.Codec
			}

			if mediaTranscodeProTemplate.Video.Profile != "" {
				videoMap["profile"] = mediaTranscodeProTemplate.Video.Profile
			}

			if mediaTranscodeProTemplate.Video.Width != "" {
				videoMap["width"] = mediaTranscodeProTemplate.Video.Width
			}

			if mediaTranscodeProTemplate.Video.Height != "" {
				videoMap["height"] = mediaTranscodeProTemplate.Video.Height
			}

			// if mediaTranscodeProTemplate.Video.Interlaced != "" {
			// 	videoMap["interlaced"] = mediaTranscodeProTemplate.Video.Interlaced
			// }

			if mediaTranscodeProTemplate.Video.Fps != "" {
				videoMap["fps"] = mediaTranscodeProTemplate.Video.Fps
			}

			if mediaTranscodeProTemplate.Video.Bitrate != "" {
				videoMap["bitrate"] = mediaTranscodeProTemplate.Video.Bitrate
			}

			if mediaTranscodeProTemplate.Video.Rotate != "" {
				videoMap["rotate"] = mediaTranscodeProTemplate.Video.Rotate
			}

			_ = d.Set("video", []interface{}{videoMap})
		}

		if mediaTranscodeProTemplate.TimeInterval != nil {
			timeIntervalMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.TimeInterval.Start != "" {
				timeIntervalMap["start"] = mediaTranscodeProTemplate.TimeInterval.Start
			}

			if mediaTranscodeProTemplate.TimeInterval.Duration != "" {
				timeIntervalMap["duration"] = mediaTranscodeProTemplate.TimeInterval.Duration
			}

			_ = d.Set("time_interval", []interface{}{timeIntervalMap})
		}

		if mediaTranscodeProTemplate.Audio != nil {
			audioMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.Audio.Codec != "" {
				audioMap["codec"] = mediaTranscodeProTemplate.Audio.Codec
			}

			if mediaTranscodeProTemplate.Audio.Remove != "" {
				audioMap["remove"] = mediaTranscodeProTemplate.Audio.Remove
			}

			_ = d.Set("audio", []interface{}{audioMap})
		}

		if mediaTranscodeProTemplate.TransConfig != nil {
			transConfigMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.TransConfig.AdjDarMethod != "" {
				transConfigMap["adj_dar_method"] = mediaTranscodeProTemplate.TransConfig.AdjDarMethod
			}

			if mediaTranscodeProTemplate.TransConfig.IsCheckReso != "" {
				transConfigMap["is_check_reso"] = mediaTranscodeProTemplate.TransConfig.IsCheckReso
			}

			if mediaTranscodeProTemplate.TransConfig.ResoAdjMethod != "" {
				transConfigMap["reso_adj_method"] = mediaTranscodeProTemplate.TransConfig.ResoAdjMethod
			}

			if mediaTranscodeProTemplate.TransConfig.IsCheckVideoBitrate != "" {
				transConfigMap["is_check_video_bitrate"] = mediaTranscodeProTemplate.TransConfig.IsCheckVideoBitrate
			}

			if mediaTranscodeProTemplate.TransConfig.VideoBitrateAdjMethod != "" {
				transConfigMap["video_bitrate_adj_method"] = mediaTranscodeProTemplate.TransConfig.VideoBitrateAdjMethod
			}

			if mediaTranscodeProTemplate.TransConfig.IsCheckAudioBitrate != "" {
				transConfigMap["is_check_audio_bitrate"] = mediaTranscodeProTemplate.TransConfig.IsCheckAudioBitrate
			}

			if mediaTranscodeProTemplate.TransConfig.AudioBitrateAdjMethod != "" {
				transConfigMap["audio_bitrate_adj_method"] = mediaTranscodeProTemplate.TransConfig.AudioBitrateAdjMethod
			}

			if mediaTranscodeProTemplate.TransConfig.DeleteMetadata != "" {
				transConfigMap["delete_metadata"] = mediaTranscodeProTemplate.TransConfig.DeleteMetadata
			}

			if mediaTranscodeProTemplate.TransConfig.IsHdr2Sdr != "" {
				transConfigMap["is_hdr2_sdr"] = mediaTranscodeProTemplate.TransConfig.IsHdr2Sdr
			}

			_ = d.Set("trans_config", []interface{}{transConfigMap})
		}
	}

	return nil
}

func resourceTencentCloudCiMediaTranscodeProTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_pro_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaTranscodeProTemplateOptions{
		Tag: "TranscodePro",
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		if clipConfigMap, ok := helper.InterfaceToMap(dMap, "clip_config"); ok {
			clipConfig := cos.ClipConfig{}
			if v, ok := clipConfigMap["duration"]; ok {
				clipConfig.Duration = v.(string)
			}
			container.ClipConfig = &clipConfig
		}
		request.Container = &container
	}

	if d.HasChange("video") {
		if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
			video := cos.TranscodeProVideo{}
			if v, ok := dMap["codec"]; ok {
				video.Codec = v.(string)
			}
			if v, ok := dMap["profile"]; ok {
				video.Profile = v.(string)
			}
			if v, ok := dMap["width"]; ok {
				video.Width = v.(string)
			}
			if v, ok := dMap["height"]; ok {
				video.Height = v.(string)
			}
			if v, ok := dMap["interlaced"]; ok {
				video.Interlaced = v.(string)
			}
			if v, ok := dMap["fps"]; ok {
				video.Fps = v.(string)
			}
			if v, ok := dMap["bitrate"]; ok {
				video.Bitrate = v.(string)
			}
			if v, ok := dMap["rotate"]; ok {
				video.Rotate = v.(string)
			}
			request.Video = &video
		}
	}

	if d.HasChange("time_interval") {
		if dMap, ok := helper.InterfacesHeadMap(d, "time_interval"); ok {
			timeInterval := cos.TimeInterval{}
			if v, ok := dMap["start"]; ok {
				timeInterval.Start = v.(string)
			}
			if v, ok := dMap["duration"]; ok {
				timeInterval.Duration = v.(string)
			}
			request.TimeInterval = &timeInterval
		}
	}

	if d.HasChange("audio") {
		if dMap, ok := helper.InterfacesHeadMap(d, "audio"); ok {
			audio := cos.TranscodeProAudio{}
			if v, ok := dMap["codec"]; ok {
				audio.Codec = v.(string)
			}
			if v, ok := dMap["remove"]; ok {
				audio.Remove = v.(string)
			}
			request.Audio = &audio
		}
	}

	if d.HasChange("vidtrans_configeo") {
		if dMap, ok := helper.InterfacesHeadMap(d, "trans_config"); ok {
			transConfig := cos.TransConfig{}
			if v, ok := dMap["adj_dar_method"]; ok {
				transConfig.AdjDarMethod = v.(string)
			}
			if v, ok := dMap["is_check_reso"]; ok {
				transConfig.IsCheckReso = v.(string)
			}
			if v, ok := dMap["reso_adj_method"]; ok {
				transConfig.ResoAdjMethod = v.(string)
			}
			if v, ok := dMap["is_check_video_bitrate"]; ok {
				transConfig.IsCheckVideoBitrate = v.(string)
			}
			if v, ok := dMap["video_bitrate_adj_method"]; ok {
				transConfig.VideoBitrateAdjMethod = v.(string)
			}
			if v, ok := dMap["is_check_audio_bitrate"]; ok {
				transConfig.IsCheckAudioBitrate = v.(string)
			}
			if v, ok := dMap["audio_bitrate_adj_method"]; ok {
				transConfig.AudioBitrateAdjMethod = v.(string)
			}
			if v, ok := dMap["delete_metadata"]; ok {
				transConfig.DeleteMetadata = v.(string)
			}
			if v, ok := dMap["is_hdr2_sdr"]; ok {
				transConfig.IsHdr2Sdr = v.(string)
			}
			request.TransConfig = &transConfig
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaTranscodeProTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaTranscodeProTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeProTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaTranscodeProTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeProTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_pro_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteCiMediaTemplateById(ctx, bucket, templateId); err != nil {
		return err
	}

	return nil
}
