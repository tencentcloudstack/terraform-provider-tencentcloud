/*
Provides a resource to create a ci media_video_montage_template

Example Usage

```hcl
resource "tencentcloud_ci_media_video_montage_template" "media_video_montage_template" {
  bucket = "terraform-ci-xxxxx"
  name = "video_montage_template"
  duration = "10.5"
  audio {
		codec = "aac"
		samplerate = "44100"
		bitrate = "128"
		channels = "4"
		remove = "false"

  }
  video {
		codec = "H.264"
		width = "1280"
		height = ""
		bitrate = "1000"
		fps = "25"
		crf = ""
		remove = ""
  }
  container {
		format = "mp4"

  }
  audio_mix {
		audio_source = "https://terraform-ci-xxxxx.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
		mix_mode = "Once"
		replace = "true"
		# effect_config {
		# 	enable_start_fadein = ""
		# 	start_fadein_time = ""
		# 	enable_end_fadeout = ""
		# 	end_fadeout_time = ""
		# 	enable_bgm_fade = ""
		# 	bgm_fade_time = ""
		# }

  }
}
```

Import

ci media_video_montage_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_video_montage_template.media_video_montage_template terraform-ci-xxxxxx#t193e5ecc1b8154e57a8376b4405ad9c63
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCiMediaVideoMontageTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaVideoMontageTemplateCreate,
		Read:   resourceTencentCloudCiMediaVideoMontageTemplateRead,
		Update: resourceTencentCloudCiMediaVideoMontageTemplateUpdate,
		Delete: resourceTencentCloudCiMediaVideoMontageTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

			"duration": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Collection duration 1: Default automatic analysis duration, 2: The unit is seconds, 3: Support float format, execution accuracy is accurate to milliseconds.",
			},

			"audio": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "audio parameters, the target file does not require Audio information, need to set Audio.Remove to true.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Codec format, value aac, mp3.",
						},
						"samplerate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sampling Rate- Unit: Hz- Optional 11025, 22050, 32000, 44100, 48000, 96000- Different packages, mp3 supports different sampling rates, as shown in the table below.",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Original audio bit rate, unit: Kbps, Value range: [8, 1000].",
						},
						"channels": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "number of channels- When Codec is set to aac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3, support 1, 2.",
						},
						"remove": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to delete the source audio stream, the value is true, false.",
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
							Required:    true,
							Description: "Codec format `H.264`.",
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
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bit rate of video output file, value range: [10, 50000], unit: Kbps, auto means adaptive bit rate.",
						},
						"fps": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frame rate, value range: (0, 60], Unit: fps.",
						},
						"crf": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bit rate-quality control factor, value range: (0, 51], If Crf is set, the setting of Bitrate will be invalid, When Bitrate is empty, the default is 25.",
						},
						"remove": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to delete the source audio stream, the value is true, false.",
						},
					},
				},
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
							Description: "Container format: mp4, flv, hls, ts, mkv.",
						},
					},
				},
			},

			"audio_mix": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "mixing parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audio_source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The media address of the audio track that needs to be mixed.",
						},
						"mix_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Mixing mode Repeat: background sound loop, Once: The background sound is played once.",
						},
						"replace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to replace the original audio of the Input media file with the mixed audio track media.",
						},
						"effect_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Mix Fade Configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_start_fadein": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "enable fade in.",
									},
									"start_fadein_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Fade in duration, greater than 0, support floating point numbers.",
									},
									"enable_end_fadeout": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "enable fade out.",
									},
									"end_fadeout_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "fade out time, greater than 0, support floating point numbers.",
									},
									"enable_bgm_fade": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Enable bgm conversion fade in.",
									},
									"bgm_fade_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "bgm transition fade-in duration, support floating point numbers.",
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

func resourceTencentCloudCiMediaVideoMontageTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_video_montage_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaVideoMontageTemplateOptions{
			Tag: "VideoMontage",
		}
		bucket     string
		templateId string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("duration"); ok {
		request.Duration = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio"); ok {
		audio := cos.Audio{}
		if v, ok := dMap["codec"]; ok {
			audio.Codec = v.(string)
		}
		if v, ok := dMap["samplerate"]; ok {
			audio.Samplerate = v.(string)
		}
		if v, ok := dMap["bitrate"]; ok {
			audio.Bitrate = v.(string)
		}
		if v, ok := dMap["channels"]; ok {
			audio.Channels = v.(string)
		}
		if v, ok := dMap["remove"]; ok {
			audio.Remove = v.(string)
		}
		request.Audio = &audio
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
		video := cos.Video{}
		if v, ok := dMap["codec"]; ok {
			video.Codec = v.(string)
		}
		if v, ok := dMap["width"]; ok {
			video.Width = v.(string)
		}
		if v, ok := dMap["height"]; ok {
			video.Height = v.(string)
		}
		if v, ok := dMap["bitrate"]; ok {
			video.Bitrate = v.(string)
		}
		if v, ok := dMap["fps"]; ok {
			video.Fps = v.(string)
		}
		if v, ok := dMap["crf"]; ok {
			video.Crf = v.(string)
		}
		if v, ok := dMap["remove"]; ok {
			video.Remove = v.(string)
		}

		request.Video = &video
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		request.Container = &container
	}

	if v, ok := d.GetOk("audio_mix"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			audioMix := cos.AudioMix{}
			if v, ok := dMap["audio_source"]; ok {
				audioMix.AudioSource = v.(string)
			}
			if v, ok := dMap["mix_mode"]; ok {
				audioMix.MixMode = v.(string)
			}
			if v, ok := dMap["replace"]; ok {
				audioMix.Replace = v.(string)
			}
			if effectConfigMap, ok := helper.InterfaceToMap(dMap, "effect_config"); ok {
				effectConfig := cos.EffectConfig{}
				if v, ok := effectConfigMap["enable_start_fadein"]; ok {
					effectConfig.EnableStartFadein = v.(string)
				}
				if v, ok := effectConfigMap["start_fadein_time"]; ok {
					effectConfig.StartFadeinTime = v.(string)
				}
				if v, ok := effectConfigMap["enable_end_fadeout"]; ok {
					effectConfig.EnableEndFadeout = v.(string)
				}
				if v, ok := effectConfigMap["end_fadeout_time"]; ok {
					effectConfig.EndFadeoutTime = v.(string)
				}
				if v, ok := effectConfigMap["enable_bgm_fade"]; ok {
					effectConfig.EnableBgmFade = v.(string)
				}
				if v, ok := effectConfigMap["bgm_fade_time"]; ok {
					effectConfig.BgmFadeTime = v.(string)
				}
				audioMix.EffectConfig = &effectConfig
			}
			request.AudioMix = append(request.AudioMix, audioMix)
		}
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaVideoMontageTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaVideoMontageTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaVideoMontageTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaVideoMontageTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVideoMontageTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_video_montage_template.read")()
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

	_ = d.Set("bucket", bucket)

	if template.Name != "" {
		_ = d.Set("name", template.Name)
	}

	if template.VideoMontage != nil {
		mediaVideoMontageTemplate := template.VideoMontage
		if mediaVideoMontageTemplate.Duration != "" {
			_ = d.Set("duration", mediaVideoMontageTemplate.Duration)
		}

		if mediaVideoMontageTemplate.Audio != nil {
			audioMap := map[string]interface{}{}

			if mediaVideoMontageTemplate.Audio.Codec != "" {
				audioMap["codec"] = mediaVideoMontageTemplate.Audio.Codec
			}

			if mediaVideoMontageTemplate.Audio.Samplerate != "" {
				audioMap["samplerate"] = mediaVideoMontageTemplate.Audio.Samplerate
			}

			if mediaVideoMontageTemplate.Audio.Bitrate != "" {
				audioMap["bitrate"] = mediaVideoMontageTemplate.Audio.Bitrate
			}

			if mediaVideoMontageTemplate.Audio.Channels != "" {
				audioMap["channels"] = mediaVideoMontageTemplate.Audio.Channels
			}

			if mediaVideoMontageTemplate.Audio.Remove != "" {
				audioMap["remove"] = mediaVideoMontageTemplate.Audio.Remove
			}

			_ = d.Set("audio", []interface{}{audioMap})
		}

		if mediaVideoMontageTemplate.Video != nil {
			videoMap := map[string]interface{}{}

			if mediaVideoMontageTemplate.Video.Codec != "" {
				videoMap["codec"] = mediaVideoMontageTemplate.Video.Codec
			}

			if mediaVideoMontageTemplate.Video.Width != "" {
				videoMap["width"] = mediaVideoMontageTemplate.Video.Width
			}

			if mediaVideoMontageTemplate.Video.Height != "" {
				videoMap["height"] = mediaVideoMontageTemplate.Video.Height
			}

			if mediaVideoMontageTemplate.Video.Bitrate != "" {
				videoMap["bitrate"] = mediaVideoMontageTemplate.Video.Bitrate
			}

			if mediaVideoMontageTemplate.Video.Fps != "" {
				videoMap["fps"] = mediaVideoMontageTemplate.Video.Fps
			}

			if mediaVideoMontageTemplate.Video.Crf != "" {
				videoMap["crf"] = mediaVideoMontageTemplate.Video.Crf
			}

			if mediaVideoMontageTemplate.Video.Remove != "" {
				videoMap["remove"] = mediaVideoMontageTemplate.Video.Remove
			}

			_ = d.Set("video", []interface{}{videoMap})
		}

		if mediaVideoMontageTemplate.Container != nil {
			containerMap := map[string]interface{}{}

			if mediaVideoMontageTemplate.Container.Format != "" {
				containerMap["format"] = mediaVideoMontageTemplate.Container.Format
			}

			_ = d.Set("container", []interface{}{containerMap})
		}

		if mediaVideoMontageTemplate.AudioMix != nil {
			audioMixList := []interface{}{}
			for _, audioMix := range mediaVideoMontageTemplate.AudioMix {
				audioMixMap := map[string]interface{}{}

				if audioMix.AudioSource != "" {
					audioMixMap["audio_source"] = audioMix.AudioSource
				}

				if audioMix.MixMode != "" {
					audioMixMap["mix_mode"] = audioMix.MixMode
				}

				if audioMix.Replace != "" {
					audioMixMap["replace"] = audioMix.Replace
				}

				if audioMix.EffectConfig != nil {
					effectConfigMap := map[string]interface{}{}

					if audioMix.EffectConfig.EnableStartFadein != "" {
						effectConfigMap["enable_start_fadein"] = audioMix.EffectConfig.EnableStartFadein
					}

					if audioMix.EffectConfig.StartFadeinTime != "" {
						effectConfigMap["start_fadein_time"] = audioMix.EffectConfig.StartFadeinTime
					}

					if audioMix.EffectConfig.EnableEndFadeout != "" {
						effectConfigMap["enable_end_fadeout"] = audioMix.EffectConfig.EnableEndFadeout
					}

					if audioMix.EffectConfig.EndFadeoutTime != "" {
						effectConfigMap["end_fadeout_time"] = audioMix.EffectConfig.EndFadeoutTime
					}

					if audioMix.EffectConfig.EnableBgmFade != "" {
						effectConfigMap["enable_bgm_fade"] = audioMix.EffectConfig.EnableBgmFade
					}

					if audioMix.EffectConfig.BgmFadeTime != "" {
						effectConfigMap["bgm_fade_time"] = audioMix.EffectConfig.BgmFadeTime
					}

					audioMixMap["effect_config"] = []interface{}{effectConfigMap}
				}

				audioMixList = append(audioMixList, audioMixMap)
			}

			_ = d.Set("audio_mix", audioMixList)
		}
	}

	return nil
}

func resourceTencentCloudCiMediaVideoMontageTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_video_montage_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaVideoMontageTemplateOptions{
		Tag: "VideoMontage",
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

	if v, ok := d.GetOk("duration"); ok {
		request.Duration = v.(string)
	}

	if d.HasChange("audio") {
		if dMap, ok := helper.InterfacesHeadMap(d, "audio"); ok {
			audio := cos.Audio{}
			if v, ok := dMap["codec"]; ok {
				audio.Codec = v.(string)
			}
			if v, ok := dMap["samplerate"]; ok {
				audio.Samplerate = v.(string)
			}
			if v, ok := dMap["bitrate"]; ok {
				audio.Bitrate = v.(string)
			}
			if v, ok := dMap["channels"]; ok {
				audio.Channels = v.(string)
			}
			if v, ok := dMap["remove"]; ok {
				audio.Remove = v.(string)
			}
			request.Audio = &audio
		}
	}

	if d.HasChange("video") {
		if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
			video := cos.Video{}
			if v, ok := dMap["codec"]; ok {
				video.Codec = v.(string)
			}
			if v, ok := dMap["width"]; ok {
				video.Width = v.(string)
			}
			if v, ok := dMap["height"]; ok {
				video.Height = v.(string)
			}
			if v, ok := dMap["bitrate"]; ok {
				video.Bitrate = v.(string)
			}
			if v, ok := dMap["fps"]; ok {
				video.Fps = v.(string)
			}
			if v, ok := dMap["crf"]; ok {
				video.Crf = v.(string)
			}
			if v, ok := dMap["remove"]; ok {
				video.Remove = v.(string)
			}

			request.Video = &video
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		request.Container = &container
	}

	if d.HasChange("audio_mix") {
		if v, ok := d.GetOk("audio_mix"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				audioMix := cos.AudioMix{}
				if v, ok := dMap["audio_source"]; ok {
					audioMix.AudioSource = v.(string)
				}
				if v, ok := dMap["mix_mode"]; ok {
					audioMix.MixMode = v.(string)
				}
				if v, ok := dMap["replace"]; ok {
					audioMix.Replace = v.(string)
				}
				if effectConfigMap, ok := helper.InterfaceToMap(dMap, "effect_config"); ok {
					effectConfig := cos.EffectConfig{}
					if v, ok := effectConfigMap["enable_start_fadein"]; ok {
						effectConfig.EnableStartFadein = v.(string)
					}
					if v, ok := effectConfigMap["start_fadein_time"]; ok {
						effectConfig.StartFadeinTime = v.(string)
					}
					if v, ok := effectConfigMap["enable_end_fadeout"]; ok {
						effectConfig.EnableEndFadeout = v.(string)
					}
					if v, ok := effectConfigMap["end_fadeout_time"]; ok {
						effectConfig.EndFadeoutTime = v.(string)
					}
					if v, ok := effectConfigMap["enable_bgm_fade"]; ok {
						effectConfig.EnableBgmFade = v.(string)
					}
					if v, ok := effectConfigMap["bgm_fade_time"]; ok {
						effectConfig.BgmFadeTime = v.(string)
					}
					audioMix.EffectConfig = &effectConfig
				}
				request.AudioMix = append(request.AudioMix, audioMix)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaVideoMontageTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaVideoMontageTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaVideoMontageTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaVideoMontageTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVideoMontageTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_video_montage_template.delete")()
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
