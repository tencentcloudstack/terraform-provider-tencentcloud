/*
Provides a resource to create a ci media_transcode_template

Example Usage

```hcl
resource "tencentcloud_ci_media_transcode_template" "media_transcode_template" {
  bucket = "terraform-ci-1308919341"
  name = "transcode_template"
  container {
		format = "mp4"
		# clip_config {
		# 	duration = ""
		# }
  }
  video {
		codec = "H.264"
		width = "1280"
		# height = ""
		fps = "30"
		remove = "false"
		profile = "high"
		bitrate = "1000"
		# crf = ""
		# gop = ""
		preset = "medium"
		# bufsize = ""
		# maxrate = ""
		# pixfmt = ""
		long_short_mode = "false"
		# rotate = ""
  }
  time_interval {
		start = "0"
		duration = "60"
  }
  audio {
		codec = "aac"
		samplerate = "44100"
		bitrate = "128"
		channels = "4"
		remove = "false"
		keep_two_tracks = "false"
		switch_track = "false"
		sample_format = ""
  }
  trans_config {
		adj_dar_method = "scale"
		is_check_reso = "false"
		reso_adj_method = "1"
		is_check_video_bitrate = "false"
		video_bitrate_adj_method = "0"
		is_check_audio_bitrate = "false"
		audio_bitrate_adj_method = "0"
		delete_metadata = "false"
		is_hdr2_sdr = "false"
		hls_encrypt {
			is_hls_encrypt = "false"
			uri_key = ""
		}
  }
  audio_mix {
		audio_source = "https://terraform-ci-1308919341.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
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
```

Import

ci media_transcode_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_transcode_template.media_transcode_template media_transcode_template_id
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

func resourceTencentCloudCiMediaTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaTranscodeTemplateCreate,
		Read:   resourceTencentCloudCiMediaTranscodeTemplateRead,
		Update: resourceTencentCloudCiMediaTranscodeTemplateUpdate,
		Delete: resourceTencentCloudCiMediaTranscodeTemplateDelete,
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
						"fps": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frame rate, value range: (0, 60], Unit: fps.",
						},
						"remove": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to delete the video stream, true, false.",
						},
						"profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "encoding level, Support baseline, main, high, auto- When Pixfmt is auto, this parameter can only be set to auto, when it is set to other options, the parameter value will be set to auto- baseline: suitable for mobile devices- main: suitable for standard resolution devices- high: suitable for high-resolution devices- Only H.264 supports this parameter.",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bit rate of video output file, value range: [10, 50000], unit: Kbps, auto means adaptive bit rate.",
						},
						"crf": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bit rate-quality control factor, value range: (0, 51], If Crf is set, the setting of Bitrate will be invalid, When Bitrate is empty, the default is 25.",
						},
						"gop": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The maximum number of frames between key frames, value range: [1, 100000].",
						},
						"preset": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video Algorithm Presets- H.264 supports this parameter, the values are veryfast, fast, medium, slow, slower- VP8 supports this parameter, the value is good, realtime- AV1 supports this parameter, the value is 5 (recommended value), 4- H.265 and VP9 do not support this parameter.",
						},
						"bufsize": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "buffer size, Value range: [1000, 128000], Unit: Kb, This parameter is not supported when Codec is VP8/VP9.",
						},
						"maxrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Peak video bit rate, Value range: [10, 50000], Unit: Kbps, This parameter is not supported when Codec is VP8/VP9.",
						},
						"pixfmt": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "video color format, H.264 support: yuv420p, yuv422p, yuv444p, yuvj420p, yuvj422p, yuvj444p, auto, H.265 support: yuv420p, yuv420p10le, auto, This parameter is not supported when Codec is VP8/VP9/AV1.",
						},
						"long_short_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Adaptive length,true, false, This parameter is not supported when Codec is VP8/VP9/AV1.",
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
							Optional:    true,
							Description: "Codec format, value aac, mp3, flac, amr, Vorbis, opus, pcm_s16le.",
						},
						"samplerate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sampling Rate- Unit: Hz- Optional 8000, 11025, 12000, 16000, 22050, 24000, 32000, 44100, 48000, 88200, 96000- Different packages, mp3 supports different sampling rates, as shown in the table below- When Codec is set to amr, only 8000 is supported- When Codec is set to opus, it supports 8000, 16000, 24000, 48000.",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Original audio bit rate, unit: Kbps, Value range: [8, 1000].",
						},
						"channels": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "number of channels- When Codec is set to aac/flac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3/opus, support 1, 2- When Codec is set to Vorbis, only 2 is supported- When Codec is set to amr, only 1 is supported- When Codec is set to pcm_s16le, only 1 and 2 are supported- When the encapsulation format is dash, 8 is not supported.",
						},
						"remove": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to delete the source audio stream, the value is true, false.",
						},
						"keep_two_tracks": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Keep dual audio tracks, the value is true, false. This parameter is invalid when Video.Codec is H.265.",
						},
						"switch_track": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Convert track, the value is true, false. This parameter is invalid when Video.Codec is H.265.",
						},
						"sample_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sampling bit width- When Codec is set to aac, support fltp- When Codec is set to mp3, fltp, s16p, s32p are supported- When Codec is set to flac, s16, s32, s16p, s32p are supported- When Codec is set to amr, support s16, s16p- When Codec is set to opus, support s16- When Codec is set to pcm_s16le, support s16- When Codec is set to Vorbis, support fltp- This parameter is invalid when Video.Codec is H.265.",
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
						"hls_encrypt": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "hls encryption configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_hls_encrypt": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to enable HLS encryption, support encryption when Container.Format is hls.",
									},
									"uri_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "HLS encrypted key, this parameter is only meaningful when IsHlsEncrypt is true.",
									},
								},
							},
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

func resourceTencentCloudCiMediaTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaTranscodeTemplateOptions{
			Tag: "Transcode",
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
		if v, ok := dMap["fps"]; ok {
			video.Fps = v.(string)
		}
		if v, ok := dMap["remove"]; ok {
			video.Remove = v.(string)
		}
		if v, ok := dMap["profile"]; ok {
			video.Profile = v.(string)
		}
		if v, ok := dMap["bitrate"]; ok {
			video.Bitrate = v.(string)
		}
		if v, ok := dMap["crf"]; ok {
			video.Crf = v.(string)
		}
		if v, ok := dMap["gop"]; ok {
			video.Gop = v.(string)
		}
		if v, ok := dMap["preset"]; ok {
			video.Preset = v.(string)
		}
		if v, ok := dMap["bufsize"]; ok {
			video.Bufsize = v.(string)
		}
		if v, ok := dMap["maxrate"]; ok {
			video.Maxrate = v.(string)
		}
		if v, ok := dMap["pixfmt"]; ok {
			video.Pixfmt = v.(string)
		}
		if v, ok := dMap["long_short_mode"]; ok {
			video.LongShortMode = v.(string)
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
		if v, ok := dMap["keep_two_tracks"]; ok {
			audio.KeepTwoTracks = v.(string)
		}
		if v, ok := dMap["switch_track"]; ok {
			audio.SwitchTrack = v.(string)
		}
		if v, ok := dMap["sample_format"]; ok {
			audio.SampleFormat = v.(string)
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
		if hlsEncryptMap, ok := helper.InterfaceToMap(dMap, "hls_encrypt"); ok {
			hlsEncrypt := cos.HlsEncrypt{}
			if v, ok := hlsEncryptMap["is_hls_encrypt"]; ok {
				if v.(string) == "true" {
					hlsEncrypt.IsHlsEncrypt = true
				} else {
					hlsEncrypt.IsHlsEncrypt = false
				}
			}
			if v, ok := hlsEncryptMap["uri_key"]; ok {
				hlsEncrypt.UriKey = v.(string)
			}
			transConfig.HlsEncrypt = &hlsEncrypt
		}
		request.TransConfig = &transConfig
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
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaTranscodeTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaTranscodeTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_template.read")()
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
	mediaTranscodeTemplate := template.TransTpl
	if mediaTranscodeTemplate.Container != nil {
		containerMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Container.Format != "" {
			containerMap["format"] = mediaTranscodeTemplate.Container.Format
		}

		if mediaTranscodeTemplate.Container.ClipConfig != nil {
			clipConfigMap := map[string]interface{}{}

			if mediaTranscodeTemplate.Container.ClipConfig.Duration != "" {
				clipConfigMap["duration"] = mediaTranscodeTemplate.Container.ClipConfig.Duration
			}

			containerMap["clip_config"] = []interface{}{clipConfigMap}
		}

		_ = d.Set("container", []interface{}{containerMap})
	}

	if mediaTranscodeTemplate.Video != nil {
		videoMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Video.Codec != "" {
			videoMap["codec"] = mediaTranscodeTemplate.Video.Codec
		}

		if mediaTranscodeTemplate.Video.Width != "" {
			videoMap["width"] = mediaTranscodeTemplate.Video.Width
		}

		if mediaTranscodeTemplate.Video.Height != "" {
			videoMap["height"] = mediaTranscodeTemplate.Video.Height
		}

		if mediaTranscodeTemplate.Video.Fps != "" {
			videoMap["fps"] = mediaTranscodeTemplate.Video.Fps
		}

		if mediaTranscodeTemplate.Video.Remove != "" {
			videoMap["remove"] = mediaTranscodeTemplate.Video.Remove
		}

		if mediaTranscodeTemplate.Video.Profile != "" {
			videoMap["profile"] = mediaTranscodeTemplate.Video.Profile
		}

		if mediaTranscodeTemplate.Video.Bitrate != "" {
			videoMap["bitrate"] = mediaTranscodeTemplate.Video.Bitrate
		}

		if mediaTranscodeTemplate.Video.Crf != "" {
			videoMap["crf"] = mediaTranscodeTemplate.Video.Crf
		}

		if mediaTranscodeTemplate.Video.Gop != "" {
			videoMap["gop"] = mediaTranscodeTemplate.Video.Gop
		}

		if mediaTranscodeTemplate.Video.Preset != "" {
			videoMap["preset"] = mediaTranscodeTemplate.Video.Preset
		}

		if mediaTranscodeTemplate.Video.Bufsize != "" {
			videoMap["bufsize"] = mediaTranscodeTemplate.Video.Bufsize
		}

		if mediaTranscodeTemplate.Video.Maxrate != "" {
			videoMap["maxrate"] = mediaTranscodeTemplate.Video.Maxrate
		}

		if mediaTranscodeTemplate.Video.Pixfmt != "" {
			videoMap["pixfmt"] = mediaTranscodeTemplate.Video.Pixfmt
		}

		if mediaTranscodeTemplate.Video.LongShortMode != "" {
			videoMap["long_short_mode"] = mediaTranscodeTemplate.Video.LongShortMode
		}

		if mediaTranscodeTemplate.Video.Rotate != "" {
			videoMap["rotate"] = mediaTranscodeTemplate.Video.Rotate
		}

		_ = d.Set("video", []interface{}{videoMap})
	}

	if mediaTranscodeTemplate.TimeInterval != nil {
		timeIntervalMap := map[string]interface{}{}

		if mediaTranscodeTemplate.TimeInterval.Start != "" {
			timeIntervalMap["start"] = mediaTranscodeTemplate.TimeInterval.Start
		}

		if mediaTranscodeTemplate.TimeInterval.Duration != "" {
			timeIntervalMap["duration"] = mediaTranscodeTemplate.TimeInterval.Duration
		}

		_ = d.Set("time_interval", []interface{}{timeIntervalMap})
	}

	if mediaTranscodeTemplate.Audio != nil {
		audioMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Audio.Codec != "" {
			audioMap["codec"] = mediaTranscodeTemplate.Audio.Codec
		}

		if mediaTranscodeTemplate.Audio.Samplerate != "" {
			audioMap["samplerate"] = mediaTranscodeTemplate.Audio.Samplerate
		}

		if mediaTranscodeTemplate.Audio.Bitrate != "" {
			audioMap["bitrate"] = mediaTranscodeTemplate.Audio.Bitrate
		}

		if mediaTranscodeTemplate.Audio.Channels != "" {
			audioMap["channels"] = mediaTranscodeTemplate.Audio.Channels
		}

		if mediaTranscodeTemplate.Audio.Remove != "" {
			audioMap["remove"] = mediaTranscodeTemplate.Audio.Remove
		}

		if mediaTranscodeTemplate.Audio.KeepTwoTracks != "" {
			audioMap["keep_two_tracks"] = mediaTranscodeTemplate.Audio.KeepTwoTracks
		}

		if mediaTranscodeTemplate.Audio.SwitchTrack != "" {
			audioMap["switch_track"] = mediaTranscodeTemplate.Audio.SwitchTrack
		}

		if mediaTranscodeTemplate.Audio.SampleFormat != "" {
			audioMap["sample_format"] = mediaTranscodeTemplate.Audio.SampleFormat
		}

		_ = d.Set("audio", []interface{}{audioMap})
	}

	if mediaTranscodeTemplate.TransConfig != nil {
		transConfigMap := map[string]interface{}{}

		if mediaTranscodeTemplate.TransConfig.AdjDarMethod != "" {
			transConfigMap["adj_dar_method"] = mediaTranscodeTemplate.TransConfig.AdjDarMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckReso != "" {
			transConfigMap["is_check_reso"] = mediaTranscodeTemplate.TransConfig.IsCheckReso
		}

		if mediaTranscodeTemplate.TransConfig.ResoAdjMethod != "" {
			transConfigMap["reso_adj_method"] = mediaTranscodeTemplate.TransConfig.ResoAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckVideoBitrate != "" {
			transConfigMap["is_check_video_bitrate"] = mediaTranscodeTemplate.TransConfig.IsCheckVideoBitrate
		}

		if mediaTranscodeTemplate.TransConfig.VideoBitrateAdjMethod != "" {
			transConfigMap["video_bitrate_adj_method"] = mediaTranscodeTemplate.TransConfig.VideoBitrateAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckAudioBitrate != "" {
			transConfigMap["is_check_audio_bitrate"] = mediaTranscodeTemplate.TransConfig.IsCheckAudioBitrate
		}

		if mediaTranscodeTemplate.TransConfig.AudioBitrateAdjMethod != "" {
			transConfigMap["audio_bitrate_adj_method"] = mediaTranscodeTemplate.TransConfig.AudioBitrateAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.DeleteMetadata != "" {
			transConfigMap["delete_metadata"] = mediaTranscodeTemplate.TransConfig.DeleteMetadata
		}

		if mediaTranscodeTemplate.TransConfig.IsHdr2Sdr != "" {
			transConfigMap["is_hdr2_sdr"] = mediaTranscodeTemplate.TransConfig.IsHdr2Sdr
		}

		if mediaTranscodeTemplate.TransConfig.HlsEncrypt != nil {
			hlsEncryptMap := map[string]interface{}{}

			hlsEncryptMap["is_hls_encrypt"] = fmt.Sprintf("%v", mediaTranscodeTemplate.TransConfig.HlsEncrypt.IsHlsEncrypt)

			if mediaTranscodeTemplate.TransConfig.HlsEncrypt.UriKey != "" {
				hlsEncryptMap["uri_key"] = mediaTranscodeTemplate.TransConfig.HlsEncrypt.UriKey
			}

			transConfigMap["hls_encrypt"] = []interface{}{hlsEncryptMap}
		}

		_ = d.Set("trans_config", []interface{}{transConfigMap})
	}

	if mediaTranscodeTemplate.AudioMix != nil {
		audioMixArrayList := []interface{}{}
		for _, audioMix := range mediaTranscodeTemplate.AudioMix {
			audioMixArrayMap := map[string]interface{}{}

			if audioMix.AudioSource != "" {
				audioMixArrayMap["audio_source"] = audioMix.AudioSource
			}

			if audioMix.MixMode != "" {
				audioMixArrayMap["mix_mode"] = audioMix.MixMode
			}

			if audioMix.Replace != "" {
				audioMixArrayMap["replace"] = audioMix.Replace
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

				audioMixArrayMap["effect_config"] = []interface{}{effectConfigMap}
			}

			audioMixArrayList = append(audioMixArrayList, audioMixArrayMap)
		}

		_ = d.Set("audio_mix_array", audioMixArrayList)

	}

	return nil
}

func resourceTencentCloudCiMediaTranscodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaTranscodeTemplateOptions{
		Tag: "Transcode",
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
			if v, ok := dMap["fps"]; ok {
				video.Fps = v.(string)
			}
			if v, ok := dMap["remove"]; ok {
				video.Remove = v.(string)
			}
			if v, ok := dMap["profile"]; ok {
				video.Profile = v.(string)
			}
			if v, ok := dMap["bitrate"]; ok {
				video.Bitrate = v.(string)
			}
			if v, ok := dMap["crf"]; ok {
				video.Crf = v.(string)
			}
			if v, ok := dMap["gop"]; ok {
				video.Gop = v.(string)
			}
			if v, ok := dMap["preset"]; ok {
				video.Preset = v.(string)
			}
			if v, ok := dMap["bufsize"]; ok {
				video.Bufsize = v.(string)
			}
			if v, ok := dMap["maxrate"]; ok {
				video.Maxrate = v.(string)
			}
			if v, ok := dMap["pixfmt"]; ok {
				video.Pixfmt = v.(string)
			}
			if v, ok := dMap["long_short_mode"]; ok {
				video.LongShortMode = v.(string)
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
			if v, ok := dMap["keep_two_tracks"]; ok {
				audio.KeepTwoTracks = v.(string)
			}
			if v, ok := dMap["switch_track"]; ok {
				audio.SwitchTrack = v.(string)
			}
			if v, ok := dMap["sample_format"]; ok {
				audio.SampleFormat = v.(string)
			}
			request.Audio = &audio
		}
	}

	if d.HasChange("trans_config") {
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
			if hlsEncryptMap, ok := helper.InterfaceToMap(dMap, "hls_encrypt"); ok {
				hlsEncrypt := cos.HlsEncrypt{}
				if v, ok := hlsEncryptMap["is_hls_encrypt"]; ok {
					if v.(string) == "true" {
						hlsEncrypt.IsHlsEncrypt = true
					} else {
						hlsEncrypt.IsHlsEncrypt = false
					}
				}
				if v, ok := hlsEncryptMap["uri_key"]; ok {
					hlsEncrypt.UriKey = v.(string)
				}
				transConfig.HlsEncrypt = &hlsEncrypt
			}
			request.TransConfig = &transConfig
		}
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

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaTranscodeTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaTranscodeTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_template.delete")()
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
