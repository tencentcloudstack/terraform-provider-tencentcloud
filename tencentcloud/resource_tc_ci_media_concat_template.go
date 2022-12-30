/*
Provides a resource to create a ci media_concat_template

Example Usage

```hcl
resource "tencentcloud_ci_media_concat_template" "media_concat_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "concat_templates"
  concat_template {
		concat_fragment {
			url = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"
			mode = "Start"
		}
    concat_fragment {
			url = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"
			mode = "End"
		}
		audio {
			codec = "mp3"
			samplerate = ""
			bitrate = ""
			channels = ""
		}
		video {
			codec = "H.264"
			width = "1280"
			height = ""
      bitrate = "1000"
			fps = "25"
			crf = ""
			remove = ""
			rotate = ""
		}
		container {
			format = "mp4"
		}
		audio_mix {
			audio_source = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
			mix_mode = "Once"
			replace = "true"
			effect_config {
				enable_start_fadein = "true"
				start_fadein_time = "3"
				enable_end_fadeout = "false"
				end_fadeout_time = "0"
				enable_bgm_fade = "true"
				bgm_fade_time = "1.7"
			}
		}
  }
}
```

Import

ci media_concat_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_concat_template.media_concat_template id=terraform-ci-xxxxxx#t1cb115dfa1fcc414284f83b7c69bcedcf
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

func resourceTencentCloudCiMediaConcatTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaConcatTemplateCreate,
		Read:   resourceTencentCloudCiMediaConcatTemplateRead,
		Update: resourceTencentCloudCiMediaConcatTemplateUpdate,
		Delete: resourceTencentCloudCiMediaConcatTemplateDelete,
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

			"concat_template": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "stitching template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"concat_fragment": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Package format.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Splicing object address.",
									},
									"mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "node type, `start`, `end`.",
									},
								},
							},
						},
						"audio": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
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
								},
							},
						},
						"video": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
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
										Description: "Original audio bit rate, unit: Kbps, Value range: [8, 1000].",
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
									"rotate": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Rotation angle, Value range: [0, 360), Unit: degree.",
									},
								},
							},
						},
						"container": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Only splicing without transcoding.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"format": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Container format: mp4, flv, hls, ts, mp3, aac.",
									},
								},
							},
						},
						"audio_mix": {
							Type:        schema.TypeList,
							Optional:    true,
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
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaConcatTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_concat_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaConcatTemplateOptions{
			Tag: "Concat",
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

	if dMap, ok := helper.InterfacesHeadMap(d, "concat_template"); ok {
		concatTemplate := cos.ConcatTemplate{}
		if v, ok := dMap["concat_fragment"]; ok {
			for _, item := range v.([]interface{}) {
				concatFragmentMap := item.(map[string]interface{})
				concatFragment := cos.ConcatFragment{}
				if v, ok := concatFragmentMap["url"]; ok {
					concatFragment.Url = v.(string)
				}
				if v, ok := concatFragmentMap["mode"]; ok {
					concatFragment.Mode = v.(string)
				}
				concatTemplate.ConcatFragment = append(concatTemplate.ConcatFragment, concatFragment)
			}
		}
		if audioMap, ok := helper.InterfaceToMap(dMap, "audio"); ok {
			audio := cos.Audio{}
			if v, ok := audioMap["codec"]; ok {
				audio.Codec = v.(string)
			}
			if v, ok := audioMap["samplerate"]; ok {
				audio.Samplerate = v.(string)
			}
			if v, ok := audioMap["bitrate"]; ok {
				audio.Bitrate = v.(string)
			}
			if v, ok := audioMap["channels"]; ok {
				audio.Channels = v.(string)
			}
			concatTemplate.Audio = &audio
		}
		if videoMap, ok := helper.InterfaceToMap(dMap, "video"); ok {
			video := cos.Video{}
			if v, ok := videoMap["codec"]; ok {
				video.Codec = v.(string)
			}
			if v, ok := videoMap["width"]; ok {
				video.Width = v.(string)
			}
			if v, ok := videoMap["height"]; ok {
				video.Height = v.(string)
			}
			if v, ok := videoMap["bitrate"]; ok {
				video.Bitrate = v.(string)
			}
			if v, ok := videoMap["fps"]; ok {
				video.Fps = v.(string)
			}
			if v, ok := videoMap["crf"]; ok {
				video.Crf = v.(string)
			}
			if v, ok := videoMap["remove"]; ok {
				video.Remove = v.(string)
			}
			if v, ok := videoMap["rotate"]; ok {
				video.Rotate = v.(string)
			}
			concatTemplate.Video = &video
		}
		if containerMap, ok := helper.InterfaceToMap(dMap, "container"); ok {
			container := cos.Container{}
			if v, ok := containerMap["format"]; ok {
				container.Format = v.(string)
			}
			concatTemplate.Container = &container
		}
		if v, ok := dMap["audio_mix"]; ok {
			for _, item := range v.([]interface{}) {
				audioMixMap := item.(map[string]interface{})
				audioMix := cos.AudioMix{}
				if v, ok := audioMixMap["audio_source"]; ok {
					audioMix.AudioSource = v.(string)
				}
				if v, ok := audioMixMap["mix_mode"]; ok {
					audioMix.MixMode = v.(string)
				}
				if v, ok := audioMixMap["replace"]; ok {
					audioMix.Replace = v.(string)
				}
				if effectConfigMap, ok := helper.InterfaceToMap(audioMixMap, "effect_config"); ok {
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
				concatTemplate.AudioMix = append(concatTemplate.AudioMix, audioMix)
			}
		}
		request.ConcatTemplate = &concatTemplate
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaConcatTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaConcatTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaConcatTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaConcatTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaConcatTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_concat_template.read")()
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

	mediaConcatTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaConcatTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if mediaConcatTemplate.Name != "" {
		_ = d.Set("name", mediaConcatTemplate.Name)
	}

	if mediaConcatTemplate.ConcatTemplate != nil {
		concatTemplateMap := map[string]interface{}{}

		if mediaConcatTemplate.ConcatTemplate.ConcatFragment != nil {
			concatFragmentList := []interface{}{}
			for _, concatFragment := range mediaConcatTemplate.ConcatTemplate.ConcatFragment {
				concatFragmentMap := map[string]interface{}{}
				if concatFragment.Url != "" {
					concatFragmentMap["url"] = concatFragment.Url
				}
				if concatFragment.Mode != "" {
					concatFragmentMap["mode"] = concatFragment.Mode
				}
				concatFragmentList = append(concatFragmentList, concatFragmentMap)
			}
			concatTemplateMap["concat_fragment"] = []interface{}{concatFragmentList}
		}

		if mediaConcatTemplate.ConcatTemplate.Audio != nil {
			audioMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Audio.Codec != "" {
				audioMap["codec"] = mediaConcatTemplate.ConcatTemplate.Audio.Codec
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Samplerate != "" {
				audioMap["samplerate"] = mediaConcatTemplate.ConcatTemplate.Audio.Samplerate
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Bitrate != "" {
				audioMap["bitrate"] = mediaConcatTemplate.ConcatTemplate.Audio.Bitrate
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Channels != "" {
				audioMap["channels"] = mediaConcatTemplate.ConcatTemplate.Audio.Channels
			}

			concatTemplateMap["audio"] = []interface{}{audioMap}
		}

		if mediaConcatTemplate.ConcatTemplate.Video != nil {
			videoMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Video.Codec != "" {
				videoMap["codec"] = mediaConcatTemplate.ConcatTemplate.Video.Codec
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Width != "" {
				videoMap["width"] = mediaConcatTemplate.ConcatTemplate.Video.Width
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Height != "" {
				videoMap["height"] = mediaConcatTemplate.ConcatTemplate.Video.Height
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Bitrate != "" {
				videoMap["bitrate"] = mediaConcatTemplate.ConcatTemplate.Video.Bitrate
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Fps != "" {
				videoMap["fps"] = mediaConcatTemplate.ConcatTemplate.Video.Fps
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Crf != "" {
				videoMap["crf"] = mediaConcatTemplate.ConcatTemplate.Video.Crf
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Remove != "" {
				videoMap["remove"] = mediaConcatTemplate.ConcatTemplate.Video.Remove
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Rotate != "" {
				videoMap["rotate"] = mediaConcatTemplate.ConcatTemplate.Video.Rotate
			}

			concatTemplateMap["video"] = []interface{}{videoMap}
		}

		if mediaConcatTemplate.ConcatTemplate.Container != nil {
			containerMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Container.Format != "" {
				containerMap["format"] = mediaConcatTemplate.ConcatTemplate.Container.Format
			}

			concatTemplateMap["container"] = []interface{}{containerMap}
		}

		if mediaConcatTemplate.ConcatTemplate.AudioMix != nil {
			audioMixList := []interface{}{}
			for _, audioMix := range mediaConcatTemplate.ConcatTemplate.AudioMix {
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

			concatTemplateMap["audio_mix"] = []interface{}{audioMixList}
		}

		_ = d.Set("concat_template", []interface{}{concatTemplateMap})
	}

	return nil
}

func resourceTencentCloudCiMediaConcatTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_concat_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaConcatTemplateOptions{
		Tag: "Concat",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("concat_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "concat_template"); ok {
			concatTemplate := cos.ConcatTemplate{}
			if v, ok := dMap["concat_fragment"]; ok {
				for _, item := range v.([]interface{}) {
					concatFragmentMap := item.(map[string]interface{})
					concatFragment := cos.ConcatFragment{}
					if v, ok := concatFragmentMap["url"]; ok {
						concatFragment.Url = v.(string)
					}
					if v, ok := concatFragmentMap["mode"]; ok {
						concatFragment.Mode = v.(string)
					}
					concatTemplate.ConcatFragment = append(concatTemplate.ConcatFragment, concatFragment)
				}
			}
			if audioMap, ok := helper.InterfaceToMap(dMap, "audio"); ok {
				audio := cos.Audio{}
				if v, ok := audioMap["codec"]; ok {
					audio.Codec = v.(string)
				}
				if v, ok := audioMap["samplerate"]; ok {
					audio.Samplerate = v.(string)
				}
				if v, ok := audioMap["bitrate"]; ok {
					audio.Bitrate = v.(string)
				}
				if v, ok := audioMap["channels"]; ok {
					audio.Channels = v.(string)
				}
				concatTemplate.Audio = &audio
			}
			if videoMap, ok := helper.InterfaceToMap(dMap, "video"); ok {
				video := cos.Video{}
				if v, ok := videoMap["codec"]; ok {
					video.Codec = v.(string)
				}
				if v, ok := videoMap["width"]; ok {
					video.Width = v.(string)
				}
				if v, ok := videoMap["height"]; ok {
					video.Height = v.(string)
				}
				if v, ok := videoMap["bitrate"]; ok {
					video.Bitrate = v.(string)
				}
				if v, ok := videoMap["fps"]; ok {
					video.Fps = v.(string)
				}
				if v, ok := videoMap["crf"]; ok {
					video.Crf = v.(string)
				}
				if v, ok := videoMap["remove"]; ok {
					video.Remove = v.(string)
				}
				if v, ok := videoMap["rotate"]; ok {
					video.Rotate = v.(string)
				}
				concatTemplate.Video = &video
			}
			if containerMap, ok := helper.InterfaceToMap(dMap, "container"); ok {
				container := cos.Container{}
				if v, ok := containerMap["format"]; ok {
					container.Format = v.(string)
				}
				concatTemplate.Container = &container
			}
			if v, ok := dMap["audio_mix"]; ok {
				for _, item := range v.([]interface{}) {
					audioMixMap := item.(map[string]interface{})
					audioMix := cos.AudioMix{}
					if v, ok := audioMixMap["audio_source"]; ok {
						audioMix.AudioSource = v.(string)
					}
					if v, ok := audioMixMap["mix_mode"]; ok {
						audioMix.MixMode = v.(string)
					}
					if v, ok := audioMixMap["replace"]; ok {
						audioMix.Replace = v.(string)
					}
					if effectConfigMap, ok := helper.InterfaceToMap(audioMixMap, "effect_config"); ok {
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
					concatTemplate.AudioMix = append(concatTemplate.AudioMix, audioMix)
				}
			}
			request.ConcatTemplate = &concatTemplate
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaConcatTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaConcatTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaConcatTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaConcatTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaConcatTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_concat_template.delete")()
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
