/*
Provides a resource to create a ci media_concat_template

Example Usage

```hcl
resource "tencentcloud_ci_media_concat_template" "media_concat_template" {
  name = &lt;nil&gt;
  concat_template {
		concat_fragment {
			url = &lt;nil&gt;
			mode = &lt;nil&gt;
		}
		audio {
			codec = &lt;nil&gt;
			samplerate = &lt;nil&gt;
			bitrate = &lt;nil&gt;
			channels = &lt;nil&gt;
		}
		video {
			codec = &lt;nil&gt;
			width = &lt;nil&gt;
			height = &lt;nil&gt;
			bitrate = &lt;nil&gt;
			fps = &lt;nil&gt;
			crf = &lt;nil&gt;
			remove = &lt;nil&gt;
			rotate = &lt;nil&gt;
		}
		container {
			format = &lt;nil&gt;
		}
		audio_mix {
			audio_source = &lt;nil&gt;
			mix_mode = &lt;nil&gt;
			replace = &lt;nil&gt;
			effect_config {
				enable_start_fadein = &lt;nil&gt;
				start_fadein_time = &lt;nil&gt;
				enable_end_fadeout = &lt;nil&gt;
				end_fadeout_time = &lt;nil&gt;
				enable_bgm_fade = &lt;nil&gt;
				bgm_fade_time = &lt;nil&gt;
			}
		}

  }
}
```

Import

ci media_concat_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_concat_template.media_concat_template media_concat_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	ci "github.com/tencentyun/cos-go-sdk-v5"
	"log"
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
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"concat_template": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Stitching template.",
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
										Description: "Node type, `start`, `end`.",
									},
								},
							},
						},
						"audio": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Audio parameters, the target file does not require Audio information, need to set Audio.Remove to true.",
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
										Description: "Number of channels- When Codec is set to aac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3, support 1, 2.",
									},
								},
							},
						},
						"video": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Video information, do not upload Video, which is equivalent to deleting video information.",
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
										Description: "Width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.",
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
							Description: "Mixing parameters.",
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
													Description: "Enable fade in.",
												},
												"start_fadein_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Fade in duration, greater than 0, support floating point numbers.",
												},
												"enable_end_fadeout": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Enable fade out.",
												},
												"end_fadeout_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Fade out time, greater than 0, support floating point numbers.",
												},
												"enable_bgm_fade": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Enable bgm conversion fade in.",
												},
												"bgm_fade_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Bgm transition fade-in duration, support floating point numbers.",
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

	var (
		request    = ci.NewCreateMediaConcatTemplateRequest()
		response   = ci.NewCreateMediaConcatTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "concat_template"); ok {
		concatTemplate := ci.ConcatTemplate{}
		if v, ok := dMap["concat_fragment"]; ok {
			for _, item := range v.([]interface{}) {
				concatFragmentMap := item.(map[string]interface{})
				concatTemplate := ci.ConcatTemplate{}
				if v, ok := concatFragmentMap["url"]; ok {
					concatTemplate.Url = helper.String(v.(string))
				}
				if v, ok := concatFragmentMap["mode"]; ok {
					concatTemplate.Mode = helper.String(v.(string))
				}
				concatTemplate.ConcatFragment = append(concatTemplate.ConcatFragment, &concatTemplate)
			}
		}
		if audioMap, ok := helper.InterfaceToMap(dMap, "audio"); ok {
			audio := ci.Audio{}
			if v, ok := audioMap["codec"]; ok {
				audio.Codec = helper.String(v.(string))
			}
			if v, ok := audioMap["samplerate"]; ok {
				audio.Samplerate = helper.String(v.(string))
			}
			if v, ok := audioMap["bitrate"]; ok {
				audio.Bitrate = helper.String(v.(string))
			}
			if v, ok := audioMap["channels"]; ok {
				audio.Channels = helper.String(v.(string))
			}
			concatTemplate.Audio = &audio
		}
		if videoMap, ok := helper.InterfaceToMap(dMap, "video"); ok {
			video := ci.Video{}
			if v, ok := videoMap["codec"]; ok {
				video.Codec = helper.String(v.(string))
			}
			if v, ok := videoMap["width"]; ok {
				video.Width = helper.String(v.(string))
			}
			if v, ok := videoMap["height"]; ok {
				video.Height = helper.String(v.(string))
			}
			if v, ok := videoMap["bitrate"]; ok {
				video.Bitrate = helper.String(v.(string))
			}
			if v, ok := videoMap["fps"]; ok {
				video.Fps = helper.String(v.(string))
			}
			if v, ok := videoMap["crf"]; ok {
				video.Crf = helper.String(v.(string))
			}
			if v, ok := videoMap["remove"]; ok {
				video.Remove = helper.String(v.(string))
			}
			if v, ok := videoMap["rotate"]; ok {
				video.Rotate = helper.String(v.(string))
			}
			concatTemplate.Video = &video
		}
		if containerMap, ok := helper.InterfaceToMap(dMap, "container"); ok {
			container := ci.Container{}
			if v, ok := containerMap["format"]; ok {
				container.Format = helper.String(v.(string))
			}
			concatTemplate.Container = &container
		}
		if v, ok := dMap["audio_mix"]; ok {
			for _, item := range v.([]interface{}) {
				audioMixMap := item.(map[string]interface{})
				audioMix := ci.AudioMix{}
				if v, ok := audioMixMap["audio_source"]; ok {
					audioMix.AudioSource = helper.String(v.(string))
				}
				if v, ok := audioMixMap["mix_mode"]; ok {
					audioMix.MixMode = helper.String(v.(string))
				}
				if v, ok := audioMixMap["replace"]; ok {
					audioMix.Replace = helper.String(v.(string))
				}
				if effectConfigMap, ok := helper.InterfaceToMap(audioMixMap, "effect_config"); ok {
					effectConfig := ci.EffectConfig{}
					if v, ok := effectConfigMap["enable_start_fadein"]; ok {
						effectConfig.EnableStartFadein = helper.String(v.(string))
					}
					if v, ok := effectConfigMap["start_fadein_time"]; ok {
						effectConfig.StartFadeinTime = helper.String(v.(string))
					}
					if v, ok := effectConfigMap["enable_end_fadeout"]; ok {
						effectConfig.EnableEndFadeout = helper.String(v.(string))
					}
					if v, ok := effectConfigMap["end_fadeout_time"]; ok {
						effectConfig.EndFadeoutTime = helper.String(v.(string))
					}
					if v, ok := effectConfigMap["enable_bgm_fade"]; ok {
						effectConfig.EnableBgmFade = helper.String(v.(string))
					}
					if v, ok := effectConfigMap["bgm_fade_time"]; ok {
						effectConfig.BgmFadeTime = helper.String(v.(string))
					}
					audioMix.EffectConfig = &effectConfig
				}
				concatTemplate.AudioMix = append(concatTemplate.AudioMix, &audioMix)
			}
		}
		request.ConcatTemplate = &concatTemplate
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaConcatTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaConcatTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaConcatTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaConcatTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_concat_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaConcatTemplateId := d.Id()

	mediaConcatTemplate, err := service.DescribeCiMediaConcatTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaConcatTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaConcatTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaConcatTemplate.Name != nil {
		_ = d.Set("name", mediaConcatTemplate.Name)
	}

	if mediaConcatTemplate.ConcatTemplate != nil {
		concatTemplateMap := map[string]interface{}{}

		if mediaConcatTemplate.ConcatTemplate.ConcatFragment != nil {
			concatFragmentList := []interface{}{}
			for _, concatFragment := range mediaConcatTemplate.ConcatTemplate.ConcatFragment {
				concatFragmentMap := map[string]interface{}{}

				if concatFragment.Url != nil {
					concatFragmentMap["url"] = concatFragment.Url
				}

				if concatFragment.Mode != nil {
					concatFragmentMap["mode"] = concatFragment.Mode
				}

				concatFragmentList = append(concatFragmentList, concatFragmentMap)
			}

			concatTemplateMap["concat_fragment"] = []interface{}{concatFragmentList}
		}

		if mediaConcatTemplate.ConcatTemplate.Audio != nil {
			audioMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Audio.Codec != nil {
				audioMap["codec"] = mediaConcatTemplate.ConcatTemplate.Audio.Codec
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Samplerate != nil {
				audioMap["samplerate"] = mediaConcatTemplate.ConcatTemplate.Audio.Samplerate
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Bitrate != nil {
				audioMap["bitrate"] = mediaConcatTemplate.ConcatTemplate.Audio.Bitrate
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Channels != nil {
				audioMap["channels"] = mediaConcatTemplate.ConcatTemplate.Audio.Channels
			}

			concatTemplateMap["audio"] = []interface{}{audioMap}
		}

		if mediaConcatTemplate.ConcatTemplate.Video != nil {
			videoMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Video.Codec != nil {
				videoMap["codec"] = mediaConcatTemplate.ConcatTemplate.Video.Codec
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Width != nil {
				videoMap["width"] = mediaConcatTemplate.ConcatTemplate.Video.Width
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Height != nil {
				videoMap["height"] = mediaConcatTemplate.ConcatTemplate.Video.Height
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Bitrate != nil {
				videoMap["bitrate"] = mediaConcatTemplate.ConcatTemplate.Video.Bitrate
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Fps != nil {
				videoMap["fps"] = mediaConcatTemplate.ConcatTemplate.Video.Fps
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Crf != nil {
				videoMap["crf"] = mediaConcatTemplate.ConcatTemplate.Video.Crf
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Remove != nil {
				videoMap["remove"] = mediaConcatTemplate.ConcatTemplate.Video.Remove
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Rotate != nil {
				videoMap["rotate"] = mediaConcatTemplate.ConcatTemplate.Video.Rotate
			}

			concatTemplateMap["video"] = []interface{}{videoMap}
		}

		if mediaConcatTemplate.ConcatTemplate.Container != nil {
			containerMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Container.Format != nil {
				containerMap["format"] = mediaConcatTemplate.ConcatTemplate.Container.Format
			}

			concatTemplateMap["container"] = []interface{}{containerMap}
		}

		if mediaConcatTemplate.ConcatTemplate.AudioMix != nil {
			audioMixList := []interface{}{}
			for _, audioMix := range mediaConcatTemplate.ConcatTemplate.AudioMix {
				audioMixMap := map[string]interface{}{}

				if audioMix.AudioSource != nil {
					audioMixMap["audio_source"] = audioMix.AudioSource
				}

				if audioMix.MixMode != nil {
					audioMixMap["mix_mode"] = audioMix.MixMode
				}

				if audioMix.Replace != nil {
					audioMixMap["replace"] = audioMix.Replace
				}

				if audioMix.EffectConfig != nil {
					effectConfigMap := map[string]interface{}{}

					if audioMix.EffectConfig.EnableStartFadein != nil {
						effectConfigMap["enable_start_fadein"] = audioMix.EffectConfig.EnableStartFadein
					}

					if audioMix.EffectConfig.StartFadeinTime != nil {
						effectConfigMap["start_fadein_time"] = audioMix.EffectConfig.StartFadeinTime
					}

					if audioMix.EffectConfig.EnableEndFadeout != nil {
						effectConfigMap["enable_end_fadeout"] = audioMix.EffectConfig.EnableEndFadeout
					}

					if audioMix.EffectConfig.EndFadeoutTime != nil {
						effectConfigMap["end_fadeout_time"] = audioMix.EffectConfig.EndFadeoutTime
					}

					if audioMix.EffectConfig.EnableBgmFade != nil {
						effectConfigMap["enable_bgm_fade"] = audioMix.EffectConfig.EnableBgmFade
					}

					if audioMix.EffectConfig.BgmFadeTime != nil {
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

	request := ci.NewUpdateMediaConcatTemplateRequest()

	mediaConcatTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "concat_template"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaConcatTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaConcatTemplate failed, reason:%+v", logId, err)
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
	mediaConcatTemplateId := d.Id()

	if err := service.DeleteCiMediaConcatTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
