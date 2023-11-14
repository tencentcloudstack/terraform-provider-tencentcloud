/*
Provides a resource to create a ci media_transcode_template

Example Usage

```hcl
resource "tencentcloud_ci_media_transcode_template" "media_transcode_template" {
  name = &lt;nil&gt;
  container {
		format = &lt;nil&gt;
		clip_config {
			duration = &lt;nil&gt;
		}

  }
  video {
		codec = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		fps = &lt;nil&gt;
		remove = &lt;nil&gt;
		profile = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		crf = &lt;nil&gt;
		gop = &lt;nil&gt;
		preset = &lt;nil&gt;
		bufsize = &lt;nil&gt;
		maxrate = &lt;nil&gt;
		pixfmt = &lt;nil&gt;
		long_short_mode = &lt;nil&gt;
		rotate = &lt;nil&gt;

  }
  time_interval {
		start = &lt;nil&gt;
		duration = &lt;nil&gt;

  }
  audio {
		codec = &lt;nil&gt;
		samplerate = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		channels = &lt;nil&gt;
		remove = &lt;nil&gt;
		keep_two_tracks = &lt;nil&gt;
		switch_track = &lt;nil&gt;
		sample_format = &lt;nil&gt;

  }
  trans_config {
		adj_dar_method = &lt;nil&gt;
		is_check_reso = &lt;nil&gt;
		reso_adj_method = &lt;nil&gt;
		is_check_video_bitrate = &lt;nil&gt;
		video_bitrate_adj_method = &lt;nil&gt;
		is_check_audio_bitrate = &lt;nil&gt;
		audio_bitrate_adj_method = &lt;nil&gt;
		is_check_video_fps = &lt;nil&gt;
		video_fps_adj_method = &lt;nil&gt;
		delete_metadata = &lt;nil&gt;
		is_hdr2_sdr = &lt;nil&gt;
		transcode_index = &lt;nil&gt;
		hls_encrypt {
			is_hls_encrypt = &lt;nil&gt;
			uri_key = &lt;nil&gt;
		}
		dash_encrypt {
			is_encrypt = &lt;nil&gt;
			uri_key = &lt;nil&gt;
		}

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
```

Import

ci media_transcode_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_transcode_template.media_transcode_template media_transcode_template_id
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
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"container": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Container format.",
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
				Description: "Video information, do not upload Video, which is equivalent to deleting video information.",
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
							Description: "Width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.",
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
							Description: "Encoding level, Support baseline, main, high, auto- When Pixfmt is auto, this parameter can only be set to auto, when it is set to other options, the parameter value will be set to auto- baseline: suitable for mobile devices- main: suitable for standard resolution devices- high: suitable for high-resolution devices- Only H.264 supports this parameter.",
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
							Description: "Buffer size, Value range: [1000, 128000], Unit: Kb, This parameter is not supported when Codec is VP8/VP9.",
						},
						"maxrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Peak video bit rate, Value range: [10, 50000], Unit: Kbps, This parameter is not supported when Codec is VP8/VP9.",
						},
						"pixfmt": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video color format, H.264 support: yuv420p, yuv422p, yuv444p, yuvj420p, yuvj422p, yuvj444p, auto, H.265 support: yuv420p, yuv420p10le, auto, This parameter is not supported when Codec is VP8/VP9/AV1.",
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
				Description: "Time interval.",
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
							Description: "Duration, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.",
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
							Description: "Number of channels- When Codec is set to aac/flac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3/opus, support 1, 2- When Codec is set to Vorbis, only 2 is supported- When Codec is set to amr, only 1 is supported- When Codec is set to pcm_s16le, only 1 and 2 are supported- When the encapsulation format is dash, 8 is not supported.",
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
				Description: "Transcoding configuration.",
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
						"is_check_video_fps": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to check the video frame rate, true, false, When false, transcode according to configuration parameters.",
						},
						"video_fps_adj_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video frame rate adjustment method, the value is 0, 1; when the output video frame rate is greater than the original video frame rate, 0 means to use the original video frame rate; 1 means return to transcoding failure, Take effect when IsCheckVideoFps is true.",
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
						"transcode_index": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the stream number to be processed, corresponding to Response.MediaInfo.Stream.Video.Index in the media information and Response.MediaInfo.Stream.Audio.Index.",
						},
						"hls_encrypt": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Hls encryption configuration.",
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
						"dash_encrypt": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dash encryption configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_encrypt": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to enable DASH encryption, support encryption when Container.Format is hls.",
									},
									"uri_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "DASH encrypted key, this parameter is only meaningful when IsEncrypt is true.",
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
	}
}

func resourceTencentCloudCiMediaTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ci.NewCreateMediaTranscodeTemplateRequest()
		response   = ci.NewCreateMediaTranscodeTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := ci.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = helper.String(v.(string))
		}
		if clipConfigMap, ok := helper.InterfaceToMap(dMap, "clip_config"); ok {
			snapshot := ci.Snapshot{}
			if v, ok := clipConfigMap["duration"]; ok {
				snapshot.Duration = helper.String(v.(string))
			}
			container.ClipConfig = &snapshot
		}
		request.Container = &container
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
		video := ci.Video{}
		if v, ok := dMap["codec"]; ok {
			video.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["width"]; ok {
			video.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			video.Height = helper.String(v.(string))
		}
		if v, ok := dMap["fps"]; ok {
			video.Fps = helper.String(v.(string))
		}
		if v, ok := dMap["remove"]; ok {
			video.Remove = helper.String(v.(string))
		}
		if v, ok := dMap["profile"]; ok {
			video.Profile = helper.String(v.(string))
		}
		if v, ok := dMap["bitrate"]; ok {
			video.Bitrate = helper.String(v.(string))
		}
		if v, ok := dMap["crf"]; ok {
			video.Crf = helper.String(v.(string))
		}
		if v, ok := dMap["gop"]; ok {
			video.Gop = helper.String(v.(string))
		}
		if v, ok := dMap["preset"]; ok {
			video.Preset = helper.String(v.(string))
		}
		if v, ok := dMap["bufsize"]; ok {
			video.Bufsize = helper.String(v.(string))
		}
		if v, ok := dMap["maxrate"]; ok {
			video.Maxrate = helper.String(v.(string))
		}
		if v, ok := dMap["pixfmt"]; ok {
			video.Pixfmt = helper.String(v.(string))
		}
		if v, ok := dMap["long_short_mode"]; ok {
			video.LongShortMode = helper.String(v.(string))
		}
		if v, ok := dMap["rotate"]; ok {
			video.Rotate = helper.String(v.(string))
		}
		request.Video = &video
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "time_interval"); ok {
		timeInterval := ci.TimeInterval{}
		if v, ok := dMap["start"]; ok {
			timeInterval.Start = helper.String(v.(string))
		}
		if v, ok := dMap["duration"]; ok {
			timeInterval.Duration = helper.String(v.(string))
		}
		request.TimeInterval = &timeInterval
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio"); ok {
		audio := ci.Audio{}
		if v, ok := dMap["codec"]; ok {
			audio.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["samplerate"]; ok {
			audio.Samplerate = helper.String(v.(string))
		}
		if v, ok := dMap["bitrate"]; ok {
			audio.Bitrate = helper.String(v.(string))
		}
		if v, ok := dMap["channels"]; ok {
			audio.Channels = helper.String(v.(string))
		}
		if v, ok := dMap["remove"]; ok {
			audio.Remove = helper.String(v.(string))
		}
		if v, ok := dMap["keep_two_tracks"]; ok {
			audio.KeepTwoTracks = helper.String(v.(string))
		}
		if v, ok := dMap["switch_track"]; ok {
			audio.SwitchTrack = helper.String(v.(string))
		}
		if v, ok := dMap["sample_format"]; ok {
			audio.SampleFormat = helper.String(v.(string))
		}
		request.Audio = &audio
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "trans_config"); ok {
		transConfig := ci.TransConfig{}
		if v, ok := dMap["adj_dar_method"]; ok {
			transConfig.AdjDarMethod = helper.String(v.(string))
		}
		if v, ok := dMap["is_check_reso"]; ok {
			transConfig.IsCheckReso = helper.String(v.(string))
		}
		if v, ok := dMap["reso_adj_method"]; ok {
			transConfig.ResoAdjMethod = helper.String(v.(string))
		}
		if v, ok := dMap["is_check_video_bitrate"]; ok {
			transConfig.IsCheckVideoBitrate = helper.String(v.(string))
		}
		if v, ok := dMap["video_bitrate_adj_method"]; ok {
			transConfig.VideoBitrateAdjMethod = helper.String(v.(string))
		}
		if v, ok := dMap["is_check_audio_bitrate"]; ok {
			transConfig.IsCheckAudioBitrate = helper.String(v.(string))
		}
		if v, ok := dMap["audio_bitrate_adj_method"]; ok {
			transConfig.AudioBitrateAdjMethod = helper.String(v.(string))
		}
		if v, ok := dMap["is_check_video_fps"]; ok {
			transConfig.IsCheckVideoFps = helper.String(v.(string))
		}
		if v, ok := dMap["video_fps_adj_method"]; ok {
			transConfig.VideoFpsAdjMethod = helper.String(v.(string))
		}
		if v, ok := dMap["delete_metadata"]; ok {
			transConfig.DeleteMetadata = helper.String(v.(string))
		}
		if v, ok := dMap["is_hdr2_sdr"]; ok {
			transConfig.IsHdr2Sdr = helper.String(v.(string))
		}
		if v, ok := dMap["transcode_index"]; ok {
			transConfig.TranscodeIndex = helper.String(v.(string))
		}
		if hlsEncryptMap, ok := helper.InterfaceToMap(dMap, "hls_encrypt"); ok {
			hlsEncrypt := ci.HlsEncrypt{}
			if v, ok := hlsEncryptMap["is_hls_encrypt"]; ok {
				hlsEncrypt.IsHlsEncrypt = helper.String(v.(string))
			}
			if v, ok := hlsEncryptMap["uri_key"]; ok {
				hlsEncrypt.UriKey = helper.String(v.(string))
			}
			transConfig.HlsEncrypt = &hlsEncrypt
		}
		if v, ok := dMap["dash_encrypt"]; ok {
			transConfig.DashEncrypt = helper.String(v.(string))
		}
		request.TransConfig = &transConfig
	}

	if v, ok := d.GetOk("audio_mix"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			audioMix := ci.AudioMix{}
			if v, ok := dMap["audio_source"]; ok {
				audioMix.AudioSource = helper.String(v.(string))
			}
			if v, ok := dMap["mix_mode"]; ok {
				audioMix.MixMode = helper.String(v.(string))
			}
			if v, ok := dMap["replace"]; ok {
				audioMix.Replace = helper.String(v.(string))
			}
			if effectConfigMap, ok := helper.InterfaceToMap(dMap, "effect_config"); ok {
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
			request.AudioMix = append(request.AudioMix, &audioMix)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaTranscodeTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaTranscodeTemplateId := d.Id()

	mediaTranscodeTemplate, err := service.DescribeCiMediaTranscodeTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaTranscodeTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaTranscodeTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaTranscodeTemplate.Name != nil {
		_ = d.Set("name", mediaTranscodeTemplate.Name)
	}

	if mediaTranscodeTemplate.Container != nil {
		containerMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Container.Format != nil {
			containerMap["format"] = mediaTranscodeTemplate.Container.Format
		}

		if mediaTranscodeTemplate.Container.ClipConfig != nil {
			clipConfigMap := map[string]interface{}{}

			if mediaTranscodeTemplate.Container.ClipConfig.Duration != nil {
				clipConfigMap["duration"] = mediaTranscodeTemplate.Container.ClipConfig.Duration
			}

			containerMap["clip_config"] = []interface{}{clipConfigMap}
		}

		_ = d.Set("container", []interface{}{containerMap})
	}

	if mediaTranscodeTemplate.Video != nil {
		videoMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Video.Codec != nil {
			videoMap["codec"] = mediaTranscodeTemplate.Video.Codec
		}

		if mediaTranscodeTemplate.Video.Width != nil {
			videoMap["width"] = mediaTranscodeTemplate.Video.Width
		}

		if mediaTranscodeTemplate.Video.Height != nil {
			videoMap["height"] = mediaTranscodeTemplate.Video.Height
		}

		if mediaTranscodeTemplate.Video.Fps != nil {
			videoMap["fps"] = mediaTranscodeTemplate.Video.Fps
		}

		if mediaTranscodeTemplate.Video.Remove != nil {
			videoMap["remove"] = mediaTranscodeTemplate.Video.Remove
		}

		if mediaTranscodeTemplate.Video.Profile != nil {
			videoMap["profile"] = mediaTranscodeTemplate.Video.Profile
		}

		if mediaTranscodeTemplate.Video.Bitrate != nil {
			videoMap["bitrate"] = mediaTranscodeTemplate.Video.Bitrate
		}

		if mediaTranscodeTemplate.Video.Crf != nil {
			videoMap["crf"] = mediaTranscodeTemplate.Video.Crf
		}

		if mediaTranscodeTemplate.Video.Gop != nil {
			videoMap["gop"] = mediaTranscodeTemplate.Video.Gop
		}

		if mediaTranscodeTemplate.Video.Preset != nil {
			videoMap["preset"] = mediaTranscodeTemplate.Video.Preset
		}

		if mediaTranscodeTemplate.Video.Bufsize != nil {
			videoMap["bufsize"] = mediaTranscodeTemplate.Video.Bufsize
		}

		if mediaTranscodeTemplate.Video.Maxrate != nil {
			videoMap["maxrate"] = mediaTranscodeTemplate.Video.Maxrate
		}

		if mediaTranscodeTemplate.Video.Pixfmt != nil {
			videoMap["pixfmt"] = mediaTranscodeTemplate.Video.Pixfmt
		}

		if mediaTranscodeTemplate.Video.LongShortMode != nil {
			videoMap["long_short_mode"] = mediaTranscodeTemplate.Video.LongShortMode
		}

		if mediaTranscodeTemplate.Video.Rotate != nil {
			videoMap["rotate"] = mediaTranscodeTemplate.Video.Rotate
		}

		_ = d.Set("video", []interface{}{videoMap})
	}

	if mediaTranscodeTemplate.TimeInterval != nil {
		timeIntervalMap := map[string]interface{}{}

		if mediaTranscodeTemplate.TimeInterval.Start != nil {
			timeIntervalMap["start"] = mediaTranscodeTemplate.TimeInterval.Start
		}

		if mediaTranscodeTemplate.TimeInterval.Duration != nil {
			timeIntervalMap["duration"] = mediaTranscodeTemplate.TimeInterval.Duration
		}

		_ = d.Set("time_interval", []interface{}{timeIntervalMap})
	}

	if mediaTranscodeTemplate.Audio != nil {
		audioMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Audio.Codec != nil {
			audioMap["codec"] = mediaTranscodeTemplate.Audio.Codec
		}

		if mediaTranscodeTemplate.Audio.Samplerate != nil {
			audioMap["samplerate"] = mediaTranscodeTemplate.Audio.Samplerate
		}

		if mediaTranscodeTemplate.Audio.Bitrate != nil {
			audioMap["bitrate"] = mediaTranscodeTemplate.Audio.Bitrate
		}

		if mediaTranscodeTemplate.Audio.Channels != nil {
			audioMap["channels"] = mediaTranscodeTemplate.Audio.Channels
		}

		if mediaTranscodeTemplate.Audio.Remove != nil {
			audioMap["remove"] = mediaTranscodeTemplate.Audio.Remove
		}

		if mediaTranscodeTemplate.Audio.KeepTwoTracks != nil {
			audioMap["keep_two_tracks"] = mediaTranscodeTemplate.Audio.KeepTwoTracks
		}

		if mediaTranscodeTemplate.Audio.SwitchTrack != nil {
			audioMap["switch_track"] = mediaTranscodeTemplate.Audio.SwitchTrack
		}

		if mediaTranscodeTemplate.Audio.SampleFormat != nil {
			audioMap["sample_format"] = mediaTranscodeTemplate.Audio.SampleFormat
		}

		_ = d.Set("audio", []interface{}{audioMap})
	}

	if mediaTranscodeTemplate.TransConfig != nil {
		transConfigMap := map[string]interface{}{}

		if mediaTranscodeTemplate.TransConfig.AdjDarMethod != nil {
			transConfigMap["adj_dar_method"] = mediaTranscodeTemplate.TransConfig.AdjDarMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckReso != nil {
			transConfigMap["is_check_reso"] = mediaTranscodeTemplate.TransConfig.IsCheckReso
		}

		if mediaTranscodeTemplate.TransConfig.ResoAdjMethod != nil {
			transConfigMap["reso_adj_method"] = mediaTranscodeTemplate.TransConfig.ResoAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckVideoBitrate != nil {
			transConfigMap["is_check_video_bitrate"] = mediaTranscodeTemplate.TransConfig.IsCheckVideoBitrate
		}

		if mediaTranscodeTemplate.TransConfig.VideoBitrateAdjMethod != nil {
			transConfigMap["video_bitrate_adj_method"] = mediaTranscodeTemplate.TransConfig.VideoBitrateAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckAudioBitrate != nil {
			transConfigMap["is_check_audio_bitrate"] = mediaTranscodeTemplate.TransConfig.IsCheckAudioBitrate
		}

		if mediaTranscodeTemplate.TransConfig.AudioBitrateAdjMethod != nil {
			transConfigMap["audio_bitrate_adj_method"] = mediaTranscodeTemplate.TransConfig.AudioBitrateAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckVideoFps != nil {
			transConfigMap["is_check_video_fps"] = mediaTranscodeTemplate.TransConfig.IsCheckVideoFps
		}

		if mediaTranscodeTemplate.TransConfig.VideoFpsAdjMethod != nil {
			transConfigMap["video_fps_adj_method"] = mediaTranscodeTemplate.TransConfig.VideoFpsAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.DeleteMetadata != nil {
			transConfigMap["delete_metadata"] = mediaTranscodeTemplate.TransConfig.DeleteMetadata
		}

		if mediaTranscodeTemplate.TransConfig.IsHdr2Sdr != nil {
			transConfigMap["is_hdr2_sdr"] = mediaTranscodeTemplate.TransConfig.IsHdr2Sdr
		}

		if mediaTranscodeTemplate.TransConfig.TranscodeIndex != nil {
			transConfigMap["transcode_index"] = mediaTranscodeTemplate.TransConfig.TranscodeIndex
		}

		if mediaTranscodeTemplate.TransConfig.HlsEncrypt != nil {
			hlsEncryptMap := map[string]interface{}{}

			if mediaTranscodeTemplate.TransConfig.HlsEncrypt.IsHlsEncrypt != nil {
				hlsEncryptMap["is_hls_encrypt"] = mediaTranscodeTemplate.TransConfig.HlsEncrypt.IsHlsEncrypt
			}

			if mediaTranscodeTemplate.TransConfig.HlsEncrypt.UriKey != nil {
				hlsEncryptMap["uri_key"] = mediaTranscodeTemplate.TransConfig.HlsEncrypt.UriKey
			}

			transConfigMap["hls_encrypt"] = []interface{}{hlsEncryptMap}
		}

		if mediaTranscodeTemplate.TransConfig.DashEncrypt != nil {
		}

		_ = d.Set("trans_config", []interface{}{transConfigMap})
	}

	if mediaTranscodeTemplate.AudioMix != nil {
		audioMixList := []interface{}{}
		for _, audioMix := range mediaTranscodeTemplate.AudioMix {
			audioMixMap := map[string]interface{}{}

			if mediaTranscodeTemplate.AudioMix.AudioSource != nil {
				audioMixMap["audio_source"] = mediaTranscodeTemplate.AudioMix.AudioSource
			}

			if mediaTranscodeTemplate.AudioMix.MixMode != nil {
				audioMixMap["mix_mode"] = mediaTranscodeTemplate.AudioMix.MixMode
			}

			if mediaTranscodeTemplate.AudioMix.Replace != nil {
				audioMixMap["replace"] = mediaTranscodeTemplate.AudioMix.Replace
			}

			if mediaTranscodeTemplate.AudioMix.EffectConfig != nil {
				effectConfigMap := map[string]interface{}{}

				if mediaTranscodeTemplate.AudioMix.EffectConfig.EnableStartFadein != nil {
					effectConfigMap["enable_start_fadein"] = mediaTranscodeTemplate.AudioMix.EffectConfig.EnableStartFadein
				}

				if mediaTranscodeTemplate.AudioMix.EffectConfig.StartFadeinTime != nil {
					effectConfigMap["start_fadein_time"] = mediaTranscodeTemplate.AudioMix.EffectConfig.StartFadeinTime
				}

				if mediaTranscodeTemplate.AudioMix.EffectConfig.EnableEndFadeout != nil {
					effectConfigMap["enable_end_fadeout"] = mediaTranscodeTemplate.AudioMix.EffectConfig.EnableEndFadeout
				}

				if mediaTranscodeTemplate.AudioMix.EffectConfig.EndFadeoutTime != nil {
					effectConfigMap["end_fadeout_time"] = mediaTranscodeTemplate.AudioMix.EffectConfig.EndFadeoutTime
				}

				if mediaTranscodeTemplate.AudioMix.EffectConfig.EnableBgmFade != nil {
					effectConfigMap["enable_bgm_fade"] = mediaTranscodeTemplate.AudioMix.EffectConfig.EnableBgmFade
				}

				if mediaTranscodeTemplate.AudioMix.EffectConfig.BgmFadeTime != nil {
					effectConfigMap["bgm_fade_time"] = mediaTranscodeTemplate.AudioMix.EffectConfig.BgmFadeTime
				}

				audioMixMap["effect_config"] = []interface{}{effectConfigMap}
			}

			audioMixList = append(audioMixList, audioMixMap)
		}

		_ = d.Set("audio_mix", audioMixList)

	}

	return nil
}

func resourceTencentCloudCiMediaTranscodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_transcode_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaTranscodeTemplateRequest()

	mediaTranscodeTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "container", "video", "time_interval", "audio", "trans_config", "audio_mix"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaTranscodeTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaTranscodeTemplate failed, reason:%+v", logId, err)
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
	mediaTranscodeTemplateId := d.Id()

	if err := service.DeleteCiMediaTranscodeTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
