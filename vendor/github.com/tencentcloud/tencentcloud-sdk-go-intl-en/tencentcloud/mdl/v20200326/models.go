// Copyright (c) 2017-2025 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20200326

import (
    tcerr "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/json"
)

type AVTemplate struct {
	// Name of an audio/video transcoding template, which can contain 1-20 case-sensitive letters and digits
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Whether video is needed. `0`: not needed; `1`: needed
	NeedVideo *uint64 `json:"NeedVideo,omitnil,omitempty" name:"NeedVideo"`

	// Video codec. Valid values: `H264`, `H265`. If this parameter is left empty, the original video codec will be used.
	Vcodec *string `json:"Vcodec,omitnil,omitempty" name:"Vcodec"`

	// Video width. Value range: (0, 4096]. The value must be an integer multiple of 2. If this parameter is left empty, the original video width will be used.
	Width *uint64 `json:"Width,omitnil,omitempty" name:"Width"`

	// Video height. Value range: (0, 4096]. The value must be an integer multiple of 2. If this parameter is left empty, the original video height will be used.
	Height *uint64 `json:"Height,omitnil,omitempty" name:"Height"`

	// Video frame rate. Value range: [1, 240]. If this parameter is left empty, the original frame rate will be used.
	Fps *uint64 `json:"Fps,omitnil,omitempty" name:"Fps"`

	// Whether to enable top speed codec transcoding. Valid values: `CLOSE` (disable), `OPEN` (enable). Default value: `CLOSE`
	TopSpeed *string `json:"TopSpeed,omitnil,omitempty" name:"TopSpeed"`

	// Compression ratio for top speed codec transcoding. Value range: [0, 50]. The lower the compression ratio, the higher the image quality.
	BitrateCompressionRatio *uint64 `json:"BitrateCompressionRatio,omitnil,omitempty" name:"BitrateCompressionRatio"`

	// Whether audio is needed. `0`: not needed; `1`: needed
	NeedAudio *int64 `json:"NeedAudio,omitnil,omitempty" name:"NeedAudio"`

	// Audio encoding format, only `AAC` and `PASSTHROUGH` are available, with `AAC` as the default.
	Acodec *string `json:"Acodec,omitnil,omitempty" name:"Acodec"`

	// Audio bitrate. If this parameter is left empty, the original bitrate will be used.
	// Valid values: `6000`, `7000`, `8000`, `10000`, `12000`, `14000`, `16000`, `20000`, `24000`, `28000`, `32000`, `40000`, `48000`, `56000`, `64000`, `80000`, `96000`, `112000`, `128000`, `160000`, `192000`, `224000`, `256000`, `288000`, `320000`, `384000`, `448000`, `512000`, `576000`, `640000`, `768000`, `896000`, `1024000`
	AudioBitrate *uint64 `json:"AudioBitrate,omitnil,omitempty" name:"AudioBitrate"`

	// Video bitrate. Value range: [50000, 40000000]. The value must be an integer multiple of 1000. If this parameter is left empty, the original bitrate will be used.
	VideoBitrate *uint64 `json:"VideoBitrate,omitnil,omitempty" name:"VideoBitrate"`

	// Bitrate control mode. Valid values: `CBR`, `ABR` (default), `VBR`.
	RateControlMode *string `json:"RateControlMode,omitnil,omitempty" name:"RateControlMode"`

	// Watermark ID
	WatermarkId *string `json:"WatermarkId,omitnil,omitempty" name:"WatermarkId"`

	// Whether to convert audio to text. `0` (default): No; `1`: Yes.
	SmartSubtitles *uint64 `json:"SmartSubtitles,omitnil,omitempty" name:"SmartSubtitles"`

	// The subtitle settings. Currently, the following subtitles are supported:
	// `eng2eng`: English speech to English text.
	// `eng2chs`: English speech to Chinese text. 
	// `eng2chseng`: English speech to English and Chinese text. 
	// `chs2chs`: Chinese speech to Chinese text.   
	// `chs2eng`: Chinese speech to English text. 
	// `chs2chseng`: Chinese speech to Chinese and English text.
	SubtitleConfiguration *string `json:"SubtitleConfiguration,omitnil,omitempty" name:"SubtitleConfiguration"`

	// Whether to enable the face blur function, 1 is on, 0 is off, and the default is 0.
	FaceBlurringEnabled *uint64 `json:"FaceBlurringEnabled,omitnil,omitempty" name:"FaceBlurringEnabled"`

	// Only AttachedInputs.AudioSelectors.Name can be selected. The following types need to be filled in: 'RTP_PUSH', 'SRT_PUSH', 'UDP_PUSH', 'RTP-FEC_PUSH'.
	AudioSelectorName *string `json:"AudioSelectorName,omitnil,omitempty" name:"AudioSelectorName"`

	// Audio transcoding special configuration information.
	AudioNormalization *AudioNormalizationSettings `json:"AudioNormalization,omitnil,omitempty" name:"AudioNormalization"`

	// Audio sampling rate, unit HZ.
	AudioSampleRate *uint64 `json:"AudioSampleRate,omitnil,omitempty" name:"AudioSampleRate"`

	// This field indicates how to specify the output video frame rate. If FOLLOW_SOURCE is selected, the output video frame rate will be set equal to the input video frame rate of the first input. If SPECIFIED_FRACTION is selected, the output video frame rate is determined by the fraction (frame rate numerator and frame rate denominator). If SPECIFIED_HZ is selected, the frame rate of the output video is determined by the HZ you enter.
	FrameRateType *string `json:"FrameRateType,omitnil,omitempty" name:"FrameRateType"`

	// Valid when the FrameRateType type you select is SPECIFIED_FRACTION, the output frame rate numerator setting.
	FrameRateNumerator *uint64 `json:"FrameRateNumerator,omitnil,omitempty" name:"FrameRateNumerator"`

	// Valid when the FrameRateType type you select is SPECIFIED_FRACTION, the output frame rate denominator setting.
	FrameRateDenominator *uint64 `json:"FrameRateDenominator,omitnil,omitempty" name:"FrameRateDenominator"`

	// The number of B frames can be selected from 1 to 3.
	BFramesNum *uint64 `json:"BFramesNum,omitnil,omitempty" name:"BFramesNum"`

	// The number of reference frames can be selected from 1 to 16.
	RefFramesNum *uint64 `json:"RefFramesNum,omitnil,omitempty" name:"RefFramesNum"`

	// Additional video bitrate configuration.
	AdditionalRateSettings *AdditionalRateSetting `json:"AdditionalRateSettings,omitnil,omitempty" name:"AdditionalRateSettings"`

	// Video encoding configuration.
	VideoCodecDetails *VideoCodecDetail `json:"VideoCodecDetails,omitnil,omitempty" name:"VideoCodecDetails"`

	// Audio encoding configuration.
	AudioCodecDetails *AudioCodecDetail `json:"AudioCodecDetails,omitnil,omitempty" name:"AudioCodecDetails"`

	// Whether to enable multiple audio tracks 0: Not required 1: Required Default value 0.
	MultiAudioTrackEnabled *uint64 `json:"MultiAudioTrackEnabled,omitnil,omitempty" name:"MultiAudioTrackEnabled"`

	// Quantity limit 0-20 Valid when MultiAudioTrackEnabled is turned on.
	AudioTracks []*AudioTrackInfo `json:"AudioTracks,omitnil,omitempty" name:"AudioTracks"`

	// Do you want to enable video enhancement? 1: Enable 0: Do not enable.
	VideoEnhanceEnabled *uint64 `json:"VideoEnhanceEnabled,omitnil,omitempty" name:"VideoEnhanceEnabled"`

	// Video enhancement configuration array.
	VideoEnhanceSettings []*VideoEnhanceSetting `json:"VideoEnhanceSettings,omitnil,omitempty" name:"VideoEnhanceSettings"`

	// Key frame interval, 300-10000, optional.
	GopSize *int64 `json:"GopSize,omitnil,omitempty" name:"GopSize"`

	// Keyframe units, only support MILLISECONDS (milliseconds).
	GopSizeUnits *string `json:"GopSizeUnits,omitnil,omitempty" name:"GopSizeUnits"`

	// Color space setting.
	ColorSpaceSettings *ColorSpaceSetting `json:"ColorSpaceSettings,omitnil,omitempty" name:"ColorSpaceSettings"`

	// Traceability watermark.
	ForensicWatermarkIds []*string `json:"ForensicWatermarkIds,omitnil,omitempty" name:"ForensicWatermarkIds"`
}

type AbWatermarkDetectionInfo struct {
	// Task ID
	TaskId *string `json:"TaskId,omitnil,omitempty" name:"TaskId"`

	// Types of testing
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// State
	State *string `json:"State,omitnil,omitempty" name:"State"`

	// Result
	Result *string `json:"Result,omitnil,omitempty" name:"Result"`

	// Error code
	ErrorCode *int64 `json:"ErrorCode,omitnil,omitempty" name:"ErrorCode"`

	// Error message
	ErrorMsg *string `json:"ErrorMsg,omitnil,omitempty" name:"ErrorMsg"`

	// Input information
	InputInfo *AbWatermarkInputInfo `json:"InputInfo,omitnil,omitempty" name:"InputInfo"`

	// Task notification configuration
	TaskNotifyConfig *TaskNotifyConfig `json:"TaskNotifyConfig,omitnil,omitempty" name:"TaskNotifyConfig"`

	// Create time
	CreateTime *int64 `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`

	// Update time
	UpdateTime *int64 `json:"UpdateTime,omitnil,omitempty" name:"UpdateTime"`

	// Input file information
	InputFileInfo *InputFileInfo `json:"InputFileInfo,omitnil,omitempty" name:"InputFileInfo"`
}

type AbWatermarkInputInfo struct {
	// Input type, optional URL/COS, currently only supports URL
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// URL input information
	UrlInputInfo *UrlInputInfo `json:"UrlInputInfo,omitnil,omitempty" name:"UrlInputInfo"`
}

type AbWatermarkSettingsReq struct {
	// Optional values: A/B.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`
}

type AbWatermarkSettingsResp struct {
	// AB watermark type.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Watermark payload.
	Content *string `json:"Content,omitnil,omitempty" name:"Content"`
}

type AdBreakSetting struct {
	// Advertising type, currently supports L-SQUEEZE
	Format *string `json:"Format,omitnil,omitempty" name:"Format"`

	// Duration, in milliseconds, requires 1000<duration<=600000. The current accuracy is seconds, which is a multiple of 1000
	Duration *uint64 `json:"Duration,omitnil,omitempty" name:"Duration"`

	// L-type compression recovery configuration
	LSqueezeSetting *LSqueezeSetting `json:"LSqueezeSetting,omitnil,omitempty" name:"LSqueezeSetting"`

	// AdSource type, supports UPLOAD_CREATIVES
	AdSource *string `json:"AdSource,omitnil,omitempty" name:"AdSource"`
}

type AdditionalRateSetting struct {
	// The maximum bit rate in a VBR scenario must be a multiple of 1000 and between 50000 - 40000000.
	VideoMaxBitrate *uint64 `json:"VideoMaxBitrate,omitnil,omitempty" name:"VideoMaxBitrate"`

	// Cache configuration supports configuring a Max Bitrate value of 1-4 times.
	BufferSize *uint64 `json:"BufferSize,omitnil,omitempty" name:"BufferSize"`

	// VBR scene is valid, video quality level, only supports user input numbers between 1-51.
	QualityLevel *uint64 `json:"QualityLevel,omitnil,omitempty" name:"QualityLevel"`
}

type AmazonS3Settings struct {
	// Access key ID of the S3 sub-account.
	AccessKeyID *string `json:"AccessKeyID,omitnil,omitempty" name:"AccessKeyID"`

	// Secret access key of the S3 sub-account.
	SecretAccessKey *string `json:"SecretAccessKey,omitnil,omitempty" name:"SecretAccessKey"`

	// Region of S3.
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// Bucket name of S3.
	Bucket *string `json:"Bucket,omitnil,omitempty" name:"Bucket"`

	// File output path, which can be empty. If it is not empty, it starts with / and ends with /.
	FilePath *string `json:"FilePath,omitnil,omitempty" name:"FilePath"`

	// User-defined name, supports alphanumeric characters, underscores, and hyphens, with a length between 1 and 32 characters.
	FileName *string `json:"FileName,omitnil,omitempty" name:"FileName"`

	// File suffix, only supports `jpg`.
	FileExt *string `json:"FileExt,omitnil,omitempty" name:"FileExt"`

	// Support `unix` or `utc0`, default unix.
	TimeFormat *string `json:"TimeFormat,omitnil,omitempty" name:"TimeFormat"`
}

type AttachedInput struct {
	// Input ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Audio selector for the input. There can be 0 to 20 audio selectors.
	// Note: this field may return `null`, indicating that no valid value was found.
	AudioSelectors []*AudioSelectorInfo `json:"AudioSelectors,omitnil,omitempty" name:"AudioSelectors"`

	// Pull mode. If the input type is `HLS_PULL` or `MP4_PULL`, you can set this parameter to `LOOP` or `ONCE`. `LOOP` is the default value.
	// Note: this field may return `null`, indicating that no valid value was found.
	PullBehavior *string `json:"PullBehavior,omitnil,omitempty" name:"PullBehavior"`

	// Input failover configuration
	// Note: this field may return `null`, indicating that no valid value was found.
	FailOverSettings *FailOverSettings `json:"FailOverSettings,omitnil,omitempty" name:"FailOverSettings"`

	// Caption selector for the input. There can be 0 to 1 audio selectors.
	CaptionSelectors []*CaptionSelector `json:"CaptionSelectors,omitnil,omitempty" name:"CaptionSelectors"`
}

type AudioCodecDetail struct {
	// Channel configuration, optional values: MONO (mono), STEREO (two-channel), 5.1 (surround).
	ChannelMode *string `json:"ChannelMode,omitnil,omitempty" name:"ChannelMode"`

	// Level in aac case, optional values: "LC" "HE-AAC" "HE-AACV2".
	Profile *string `json:"Profile,omitnil,omitempty" name:"Profile"`
}

type AudioNormalizationSettings struct {
	// Whether to enable special configuration for audio transcoding: 1: Enable 0: Disable, the default value is 0.
	AudioNormalizationEnabled *uint64 `json:"AudioNormalizationEnabled,omitnil,omitempty" name:"AudioNormalizationEnabled"`

	// Loudness value, floating-point number, rounded to one decimal place, range -5 to -70.
	TargetLUFS *float64 `json:"TargetLUFS,omitnil,omitempty" name:"TargetLUFS"`
}

type AudioPidSelectionInfo struct {
	// Audio `Pid`. Default value: 0.
	Pid *uint64 `json:"Pid,omitnil,omitempty" name:"Pid"`
}

type AudioPipelineInputStatistics struct {
	// Audio FPS.
	Fps *uint64 `json:"Fps,omitnil,omitempty" name:"Fps"`

	// Audio bitrate in bps.
	Rate *uint64 `json:"Rate,omitnil,omitempty" name:"Rate"`

	// Audio `Pid`, which is available only if the input is `rtp/udp`.
	Pid *int64 `json:"Pid,omitnil,omitempty" name:"Pid"`
}

type AudioSelectorInfo struct {
	// Audio name, which can contain 1-32 letters, digits, and underscores.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Audio `Pid` selection.
	AudioPidSelection *AudioPidSelectionInfo `json:"AudioPidSelection,omitnil,omitempty" name:"AudioPidSelection"`

	// Audio input type, optional values: 'PID_SELECTOR' 'TRACK_SELECTOR', default value PID_SELECTOR.
	AudioSelectorType *string `json:"AudioSelectorType,omitnil,omitempty" name:"AudioSelectorType"`

	// AudioTrack configuration.
	AudioTrackSelection *InputTracks `json:"AudioTrackSelection,omitnil,omitempty" name:"AudioTrackSelection"`
}

type AudioTemplateInfo struct {
	// Only `AttachedInputs.AudioSelectors.Name` can be selected. This parameter is required for RTP_PUSH and UDP_PUSH.
	AudioSelectorName *string `json:"AudioSelectorName,omitnil,omitempty" name:"AudioSelectorName"`

	// Audio transcoding template name, which can contain 1-20 letters and digits.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Audio encoding format, only `AAC` and `PASSTHROUGH` are available, with `AAC` as the default.
	Acodec *string `json:"Acodec,omitnil,omitempty" name:"Acodec"`

	// Audio bitrate. If this parameter is left empty, the original value will be used.
	// Valid values: 6000, 7000, 8000, 10000, 12000, 14000, 16000, 20000, 24000, 28000, 32000, 40000, 48000, 56000, 64000, 80000, 96000, 112000, 128000, 160000, 192000, 224000, 256000, 288000, 320000, 384000, 448000, 512000, 576000, 640000, 768000, 896000, 1024000
	AudioBitrate *uint64 `json:"AudioBitrate,omitnil,omitempty" name:"AudioBitrate"`

	// Audio language code, which length is between 2 and 20.
	LanguageCode *string `json:"LanguageCode,omitnil,omitempty" name:"LanguageCode"`

	// Audio transcoding special configuration information.
	AudioNormalization *AudioNormalizationSettings `json:"AudioNormalization,omitnil,omitempty" name:"AudioNormalization"`

	// Audio sampling rate, unit HZ.
	AudioSampleRate *uint64 `json:"AudioSampleRate,omitnil,omitempty" name:"AudioSampleRate"`

	// Audio encoding parameters.
	AudioCodecDetails *AudioCodecDetail `json:"AudioCodecDetails,omitnil,omitempty" name:"AudioCodecDetails"`

	// Audio language description, which maximum length is 100.
	LanguageDescription *string `json:"LanguageDescription,omitnil,omitempty" name:"LanguageDescription"`
}

type AudioTrackInfo struct {
	// User input is limited to letters and numbers, the length should not exceed 20, and should not be repeated in the same channel.
	TrackName *string `json:"TrackName,omitnil,omitempty" name:"TrackName"`

	// Audio encoding format, only `AAC` and `PASSTHROUGH` are available, with `AAC` as the default.
	AudioCodec *string `json:"AudioCodec,omitnil,omitempty" name:"AudioCodec"`

	// Audio bitrate.
	AudioBitrate *uint64 `json:"AudioBitrate,omitnil,omitempty" name:"AudioBitrate"`

	// Audio sample rate.
	AudioSampleRate *uint64 `json:"AudioSampleRate,omitnil,omitempty" name:"AudioSampleRate"`

	// Only values defined by AttachedInputs.$.AudioSelectors.$.audioPidSelection.pid can be entered.
	AudioSelectorName *string `json:"AudioSelectorName,omitnil,omitempty" name:"AudioSelectorName"`

	// Audio loudness configuration.
	AudioNormalization *AudioNormalizationSettings `json:"AudioNormalization,omitnil,omitempty" name:"AudioNormalization"`

	// Audio encoding configuration.
	AudioCodecDetails *AudioCodecDetail `json:"AudioCodecDetails,omitnil,omitempty" name:"AudioCodecDetails"`
}

type CaptionSelector struct {
	// Caption selector name, which can contain 1-32 letters, digits, and underscores.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Caption source type, only support `SCTE-128`.
	CaptionSourceType *string `json:"CaptionSourceType,omitnil,omitempty" name:"CaptionSourceType"`
}

type ChannelAlertInfos struct {
	// Alarm details of pipeline 0 under this channel.
	Pipeline0 []*ChannelPipelineAlerts `json:"Pipeline0,omitnil,omitempty" name:"Pipeline0"`

	// Alarm details of pipeline 1 under this channel.
	Pipeline1 []*ChannelPipelineAlerts `json:"Pipeline1,omitnil,omitempty" name:"Pipeline1"`

	// Pipeline 0 total active alarm count
	PipelineAActiveAlerts *int64 `json:"PipelineAActiveAlerts,omitnil,omitempty" name:"PipelineAActiveAlerts"`

	// Pipeline 1 total active alarm count
	PipelineBActiveAlerts *int64 `json:"PipelineBActiveAlerts,omitnil,omitempty" name:"PipelineBActiveAlerts"`
}

type ChannelInputStatistics struct {
	// Input ID.
	InputId *string `json:"InputId,omitnil,omitempty" name:"InputId"`

	// Input statistics.
	Statistics *InputStatistics `json:"Statistics,omitnil,omitempty" name:"Statistics"`
}

type ChannelOutputsStatistics struct {
	// Output group name.
	OutputGroupName *string `json:"OutputGroupName,omitnil,omitempty" name:"OutputGroupName"`

	// Output group statistics.
	Statistics *OutputsStatistics `json:"Statistics,omitnil,omitempty" name:"Statistics"`
}

type ChannelPipelineAlerts struct {
	// Alarm start time in UTC time.
	SetTime *string `json:"SetTime,omitnil,omitempty" name:"SetTime"`

	// Alarm end time in UTC time.
	// This time is available only after the alarm ends.
	ClearTime *string `json:"ClearTime,omitnil,omitempty" name:"ClearTime"`

	// Alarm type.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Alarm details.
	Message *string `json:"Message,omitnil,omitempty" name:"Message"`
}

type ColorSpaceSetting struct {
	// Color space, supports `PASSTHROUGH` (transparent transmission, only supports H265); optional.
	ColorSpace *string `json:"ColorSpace,omitnil,omitempty" name:"ColorSpace"`
}

type CosSettings struct {
	// Region of COS.
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// Bucket name of COS.
	Bucket *string `json:"Bucket,omitnil,omitempty" name:"Bucket"`

	// File output path, which can be empty. If it is not empty, it  ends with /.
	FilePath *string `json:"FilePath,omitnil,omitempty" name:"FilePath"`

	// User-defined name, supports alphanumeric characters, underscores, and hyphens, with a length between 1 and 32 characters.
	FileName *string `json:"FileName,omitnil,omitempty" name:"FileName"`

	// File suffix, only supports `jpg`.
	FileExt *string `json:"FileExt,omitnil,omitempty" name:"FileExt"`

	// Support `unix` or `utc0`, default unix.
	TimeFormat *string `json:"TimeFormat,omitnil,omitempty" name:"TimeFormat"`
}

type CreateImageSettings struct {
	// Image file format. Valid values: png, jpg.
	ImageType *string `json:"ImageType,omitnil,omitempty" name:"ImageType"`

	// Base64 encoded image content
	ImageContent *string `json:"ImageContent,omitnil,omitempty" name:"ImageContent"`

	// Origin. Valid values: TOP_LEFT, BOTTOM_LEFT, TOP_RIGHT, BOTTOM_RIGHT.
	Location *string `json:"Location,omitnil,omitempty" name:"Location"`

	// The watermark's horizontal distance from the origin as a percentage of the video width. Value range: 0-100. Default: 10.
	XPos *int64 `json:"XPos,omitnil,omitempty" name:"XPos"`

	// The watermark's vertical distance from the origin as a percentage of the video height. Value range: 0-100. Default: 10.
	YPos *int64 `json:"YPos,omitnil,omitempty" name:"YPos"`

	// The watermark image's width as a percentage of the video width. Value range: 0-100. Default: 10.
	// `0` means to scale the width proportionally to the height.
	// You cannot set both `Width` and `Height` to `0`.
	Width *int64 `json:"Width,omitnil,omitempty" name:"Width"`

	// The watermark image's height as a percentage of the video height. Value range: 0-100. Default: 10.
	// `0` means to scale the height proportionally to the width.
	// You cannot set both `Width` and `Height` to `0`.
	Height *int64 `json:"Height,omitnil,omitempty" name:"Height"`
}

// Predefined struct for user
type CreateStreamLiveChannelRequestParams struct {
	// Channel name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Inputs to attach. You can attach 1 to 5 inputs.
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitnil,omitempty" name:"AttachedInputs"`

	// Configuration information of the channel's output groups. Quantity: [1, 10]
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitnil,omitempty" name:"OutputGroups"`

	// Audio transcoding templates. Quantity: [1, 20]
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitnil,omitempty" name:"AudioTemplates"`

	// Video transcoding templates. Quantity: [1, 10]
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitnil,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates. Quantity: [1, 10]
	AVTemplates []*AVTemplate `json:"AVTemplates,omitnil,omitempty" name:"AVTemplates"`

	// Subtitle template configuration.
	CaptionTemplates []*SubtitleConf `json:"CaptionTemplates,omitnil,omitempty" name:"CaptionTemplates"`

	// Event settings
	PlanSettings *PlanSettings `json:"PlanSettings,omitnil,omitempty" name:"PlanSettings"`

	// The callback settings.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitnil,omitempty" name:"EventNotifySettings"`

	// Complement the last video frame settings.
	InputLossBehavior *InputLossBehaviorInfo `json:"InputLossBehavior,omitnil,omitempty" name:"InputLossBehavior"`

	// Pipeline configuration.
	PipelineInputSettings *PipelineInputSettingsInfo `json:"PipelineInputSettings,omitnil,omitempty" name:"PipelineInputSettings"`

	// Recognition configuration for input content.
	InputAnalysisSettings *InputAnalysisInfo `json:"InputAnalysisSettings,omitnil,omitempty" name:"InputAnalysisSettings"`

	// Console tag list.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`

	// Frame capture templates.
	FrameCaptureTemplates []*FrameCaptureTemplate `json:"FrameCaptureTemplates,omitnil,omitempty" name:"FrameCaptureTemplates"`

	// General settings.
	GeneralSettings *GeneralSetting `json:"GeneralSettings,omitnil,omitempty" name:"GeneralSettings"`
}

type CreateStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Inputs to attach. You can attach 1 to 5 inputs.
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitnil,omitempty" name:"AttachedInputs"`

	// Configuration information of the channel's output groups. Quantity: [1, 10]
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitnil,omitempty" name:"OutputGroups"`

	// Audio transcoding templates. Quantity: [1, 20]
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitnil,omitempty" name:"AudioTemplates"`

	// Video transcoding templates. Quantity: [1, 10]
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitnil,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates. Quantity: [1, 10]
	AVTemplates []*AVTemplate `json:"AVTemplates,omitnil,omitempty" name:"AVTemplates"`

	// Subtitle template configuration.
	CaptionTemplates []*SubtitleConf `json:"CaptionTemplates,omitnil,omitempty" name:"CaptionTemplates"`

	// Event settings
	PlanSettings *PlanSettings `json:"PlanSettings,omitnil,omitempty" name:"PlanSettings"`

	// The callback settings.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitnil,omitempty" name:"EventNotifySettings"`

	// Complement the last video frame settings.
	InputLossBehavior *InputLossBehaviorInfo `json:"InputLossBehavior,omitnil,omitempty" name:"InputLossBehavior"`

	// Pipeline configuration.
	PipelineInputSettings *PipelineInputSettingsInfo `json:"PipelineInputSettings,omitnil,omitempty" name:"PipelineInputSettings"`

	// Recognition configuration for input content.
	InputAnalysisSettings *InputAnalysisInfo `json:"InputAnalysisSettings,omitnil,omitempty" name:"InputAnalysisSettings"`

	// Console tag list.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`

	// Frame capture templates.
	FrameCaptureTemplates []*FrameCaptureTemplate `json:"FrameCaptureTemplates,omitnil,omitempty" name:"FrameCaptureTemplates"`

	// General settings.
	GeneralSettings *GeneralSetting `json:"GeneralSettings,omitnil,omitempty" name:"GeneralSettings"`
}

func (r *CreateStreamLiveChannelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLiveChannelRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "AttachedInputs")
	delete(f, "OutputGroups")
	delete(f, "AudioTemplates")
	delete(f, "VideoTemplates")
	delete(f, "AVTemplates")
	delete(f, "CaptionTemplates")
	delete(f, "PlanSettings")
	delete(f, "EventNotifySettings")
	delete(f, "InputLossBehavior")
	delete(f, "PipelineInputSettings")
	delete(f, "InputAnalysisSettings")
	delete(f, "Tags")
	delete(f, "FrameCaptureTemplates")
	delete(f, "GeneralSettings")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateStreamLiveChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveChannelResponseParams struct {
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Tag prompt information, this information will be attached when the tag operation fails.
	TagMsg *string `json:"TagMsg,omitnil,omitempty" name:"TagMsg"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateStreamLiveChannelResponse struct {
	*tchttp.BaseResponse
	Response *CreateStreamLiveChannelResponseParams `json:"Response"`
}

func (r *CreateStreamLiveChannelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLiveChannelResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveInputRequestParams struct {
	// Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Input type
	// Valid values: `RTMP_PUSH`, `RTP_PUSH`, `UDP_PUSH`, `RTMP_PULL`, `HLS_PULL`, `MP4_PULL`,`RTP-FEC_PUSH`,`RTSP_PULL`,`SRT_PUSH `,`SRT_PULL `
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// ID of the input security group to attach
	// You can attach only one security group to an input.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Input settings. 
	// For the type:
	// `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`,`RTSP_PULL`,`SRT_PULL` or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
	// For the type:
	// `SRT_PUSH`, 0 or 2 inputs of the corresponding type can be configured.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitnil,omitempty" name:"InputSettings"`
}

type CreateStreamLiveInputRequest struct {
	*tchttp.BaseRequest
	
	// Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Input type
	// Valid values: `RTMP_PUSH`, `RTP_PUSH`, `UDP_PUSH`, `RTMP_PULL`, `HLS_PULL`, `MP4_PULL`,`RTP-FEC_PUSH`,`RTSP_PULL`,`SRT_PUSH `,`SRT_PULL `
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// ID of the input security group to attach
	// You can attach only one security group to an input.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Input settings. 
	// For the type:
	// `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`,`RTSP_PULL`,`SRT_PULL` or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
	// For the type:
	// `SRT_PUSH`, 0 or 2 inputs of the corresponding type can be configured.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitnil,omitempty" name:"InputSettings"`
}

func (r *CreateStreamLiveInputRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLiveInputRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "Type")
	delete(f, "SecurityGroupIds")
	delete(f, "InputSettings")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateStreamLiveInputRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveInputResponseParams struct {
	// Input ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateStreamLiveInputResponse struct {
	*tchttp.BaseResponse
	Response *CreateStreamLiveInputResponseParams `json:"Response"`
}

func (r *CreateStreamLiveInputResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLiveInputResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveInputSecurityGroupRequestParams struct {
	// Input security group name, which can contain case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Allowlist entries. Quantity: [1, 10]
	Whitelist []*string `json:"Whitelist,omitnil,omitempty" name:"Whitelist"`
}

type CreateStreamLiveInputSecurityGroupRequest struct {
	*tchttp.BaseRequest
	
	// Input security group name, which can contain case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Allowlist entries. Quantity: [1, 10]
	Whitelist []*string `json:"Whitelist,omitnil,omitempty" name:"Whitelist"`
}

func (r *CreateStreamLiveInputSecurityGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLiveInputSecurityGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "Whitelist")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateStreamLiveInputSecurityGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveInputSecurityGroupResponseParams struct {
	// Security group ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateStreamLiveInputSecurityGroupResponse struct {
	*tchttp.BaseResponse
	Response *CreateStreamLiveInputSecurityGroupResponseParams `json:"Response"`
}

func (r *CreateStreamLiveInputSecurityGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLiveInputSecurityGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLivePlanRequestParams struct {
	// ID of the channel for which you want to configure an event
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Event configuration
	Plan *PlanReq `json:"Plan,omitnil,omitempty" name:"Plan"`
}

type CreateStreamLivePlanRequest struct {
	*tchttp.BaseRequest
	
	// ID of the channel for which you want to configure an event
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Event configuration
	Plan *PlanReq `json:"Plan,omitnil,omitempty" name:"Plan"`
}

func (r *CreateStreamLivePlanRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLivePlanRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChannelId")
	delete(f, "Plan")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateStreamLivePlanRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLivePlanResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateStreamLivePlanResponse struct {
	*tchttp.BaseResponse
	Response *CreateStreamLivePlanResponseParams `json:"Response"`
}

func (r *CreateStreamLivePlanResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLivePlanResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveWatermarkRequestParams struct {
	// Watermark name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Watermark type. Valid values: STATIC_IMAGE, TEXT.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	ImageSettings *CreateImageSettings `json:"ImageSettings,omitnil,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	TextSettings *CreateTextSettings `json:"TextSettings,omitnil,omitempty" name:"TextSettings"`

	// AB watermark configuration
	AbWatermarkSettings *AbWatermarkSettingsReq `json:"AbWatermarkSettings,omitnil,omitempty" name:"AbWatermarkSettings"`
}

type CreateStreamLiveWatermarkRequest struct {
	*tchttp.BaseRequest
	
	// Watermark name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Watermark type. Valid values: STATIC_IMAGE, TEXT.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	ImageSettings *CreateImageSettings `json:"ImageSettings,omitnil,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	TextSettings *CreateTextSettings `json:"TextSettings,omitnil,omitempty" name:"TextSettings"`

	// AB watermark configuration
	AbWatermarkSettings *AbWatermarkSettingsReq `json:"AbWatermarkSettings,omitnil,omitempty" name:"AbWatermarkSettings"`
}

func (r *CreateStreamLiveWatermarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLiveWatermarkRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Name")
	delete(f, "Type")
	delete(f, "ImageSettings")
	delete(f, "TextSettings")
	delete(f, "AbWatermarkSettings")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateStreamLiveWatermarkRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveWatermarkResponseParams struct {
	// Watermark ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateStreamLiveWatermarkResponse struct {
	*tchttp.BaseResponse
	Response *CreateStreamLiveWatermarkResponseParams `json:"Response"`
}

func (r *CreateStreamLiveWatermarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateStreamLiveWatermarkResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type CreateTextSettings struct {
	// Text
	Text *string `json:"Text,omitnil,omitempty" name:"Text"`

	// Origin. Valid values: TOP_LEFT, BOTTOM_LEFT, TOP_RIGHT, BOTTOM_RIGHT.
	Location *string `json:"Location,omitnil,omitempty" name:"Location"`

	// The watermark's horizontal distance from the origin as a percentage of the video width. Value range: 0-100. Default: 10.
	XPos *int64 `json:"XPos,omitnil,omitempty" name:"XPos"`

	// The watermark's vertical distance from the origin as a percentage of the video height. Value range: 0-100. Default: 10.
	YPos *int64 `json:"YPos,omitnil,omitempty" name:"YPos"`

	// Font size. Value range: 25-50.
	FontSize *int64 `json:"FontSize,omitnil,omitempty" name:"FontSize"`

	// Font color, which is an RGB color value. Default value: 0x000000.
	FontColor *string `json:"FontColor,omitnil,omitempty" name:"FontColor"`
}

// Predefined struct for user
type CreateWatermarkDetectionRequestParams struct {
	// Task type, currently supports ExtractVideoABWatermarkId
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Input information
	InputInfo *AbWatermarkInputInfo `json:"InputInfo,omitnil,omitempty" name:"InputInfo"`

	// Input file information
	InputFileInfo *InputFileInfo `json:"InputFileInfo,omitnil,omitempty" name:"InputFileInfo"`

	// Input notification configuration
	TaskNotifyConfig *TaskNotifyConfig `json:"TaskNotifyConfig,omitnil,omitempty" name:"TaskNotifyConfig"`
}

type CreateWatermarkDetectionRequest struct {
	*tchttp.BaseRequest
	
	// Task type, currently supports ExtractVideoABWatermarkId
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Input information
	InputInfo *AbWatermarkInputInfo `json:"InputInfo,omitnil,omitempty" name:"InputInfo"`

	// Input file information
	InputFileInfo *InputFileInfo `json:"InputFileInfo,omitnil,omitempty" name:"InputFileInfo"`

	// Input notification configuration
	TaskNotifyConfig *TaskNotifyConfig `json:"TaskNotifyConfig,omitnil,omitempty" name:"TaskNotifyConfig"`
}

func (r *CreateWatermarkDetectionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateWatermarkDetectionRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Type")
	delete(f, "InputInfo")
	delete(f, "InputFileInfo")
	delete(f, "TaskNotifyConfig")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateWatermarkDetectionRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateWatermarkDetectionResponseParams struct {
	// Task ID
	TaskId *string `json:"TaskId,omitnil,omitempty" name:"TaskId"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateWatermarkDetectionResponse struct {
	*tchttp.BaseResponse
	Response *CreateWatermarkDetectionResponseParams `json:"Response"`
}

func (r *CreateWatermarkDetectionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *CreateWatermarkDetectionResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DashRemuxSettingsInfo struct {
	// Segment duration in ms. Value range: [1000,30000]. Default value: 4000. The value can only be a multiple of 1,000.
	SegmentDuration *uint64 `json:"SegmentDuration,omitnil,omitempty" name:"SegmentDuration"`

	// Number of segments. Value range: [1,30]. Default value: 5.
	SegmentNumber *uint64 `json:"SegmentNumber,omitnil,omitempty" name:"SegmentNumber"`

	// Whether to enable multi-period. Valid values: CLOSE/OPEN. Default value: CLOSE.
	PeriodTriggers *string `json:"PeriodTriggers,omitnil,omitempty" name:"PeriodTriggers"`

	// The HLS package type when the H.265 codec is used. Valid values: `hvc1`, `hev1` (default).
	H265PackageType *string `json:"H265PackageType,omitnil,omitempty" name:"H265PackageType"`
}

// Predefined struct for user
type DeleteStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type DeleteStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *DeleteStreamLiveChannelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLiveChannelRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteStreamLiveChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLiveChannelResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteStreamLiveChannelResponse struct {
	*tchttp.BaseResponse
	Response *DeleteStreamLiveChannelResponseParams `json:"Response"`
}

func (r *DeleteStreamLiveChannelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLiveChannelResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLiveInputRequestParams struct {
	// Input ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type DeleteStreamLiveInputRequest struct {
	*tchttp.BaseRequest
	
	// Input ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *DeleteStreamLiveInputRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLiveInputRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteStreamLiveInputRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLiveInputResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteStreamLiveInputResponse struct {
	*tchttp.BaseResponse
	Response *DeleteStreamLiveInputResponseParams `json:"Response"`
}

func (r *DeleteStreamLiveInputResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLiveInputResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLiveInputSecurityGroupRequestParams struct {
	// Input security group ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type DeleteStreamLiveInputSecurityGroupRequest struct {
	*tchttp.BaseRequest
	
	// Input security group ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *DeleteStreamLiveInputSecurityGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLiveInputSecurityGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteStreamLiveInputSecurityGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLiveInputSecurityGroupResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteStreamLiveInputSecurityGroupResponse struct {
	*tchttp.BaseResponse
	Response *DeleteStreamLiveInputSecurityGroupResponseParams `json:"Response"`
}

func (r *DeleteStreamLiveInputSecurityGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLiveInputSecurityGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLivePlanRequestParams struct {
	// ID of the channel whose event is to be deleted
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Name of the event to delete
	EventName *string `json:"EventName,omitnil,omitempty" name:"EventName"`
}

type DeleteStreamLivePlanRequest struct {
	*tchttp.BaseRequest
	
	// ID of the channel whose event is to be deleted
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Name of the event to delete
	EventName *string `json:"EventName,omitnil,omitempty" name:"EventName"`
}

func (r *DeleteStreamLivePlanRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLivePlanRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChannelId")
	delete(f, "EventName")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteStreamLivePlanRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLivePlanResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteStreamLivePlanResponse struct {
	*tchttp.BaseResponse
	Response *DeleteStreamLivePlanResponseParams `json:"Response"`
}

func (r *DeleteStreamLivePlanResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLivePlanResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLiveWatermarkRequestParams struct {
	// Watermark ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type DeleteStreamLiveWatermarkRequest struct {
	*tchttp.BaseRequest
	
	// Watermark ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *DeleteStreamLiveWatermarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLiveWatermarkRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DeleteStreamLiveWatermarkRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DeleteStreamLiveWatermarkResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DeleteStreamLiveWatermarkResponse struct {
	*tchttp.BaseResponse
	Response *DeleteStreamLiveWatermarkResponseParams `json:"Response"`
}

func (r *DeleteStreamLiveWatermarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DeleteStreamLiveWatermarkResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DeliveryRestrictionsInfo struct {
	// Corresponds to SCTE-35 web_delivery_allowed_flag parameter.
	WebDeliveryAllowed *string `json:"WebDeliveryAllowed,omitnil,omitempty" name:"WebDeliveryAllowed"`

	// Corresponds to SCTE-35 no_regional_blackout_flag parameter.
	NoRegionalBlackout *string `json:"NoRegionalBlackout,omitnil,omitempty" name:"NoRegionalBlackout"`

	// Corresponds to SCTE-35 archive_allowed_flag.
	ArchiveAllowed *string `json:"ArchiveAllowed,omitnil,omitempty" name:"ArchiveAllowed"`

	// Corresponds to SCTE-35 device_restrictions parameter.
	DeviceRestrictions *string `json:"DeviceRestrictions,omitnil,omitempty" name:"DeviceRestrictions"`
}

type DescribeImageSettings struct {
	// Origin
	Location *string `json:"Location,omitnil,omitempty" name:"Location"`

	// The watermark image’s horizontal distance from the origin as a percentage of the video width
	XPos *int64 `json:"XPos,omitnil,omitempty" name:"XPos"`

	// The watermark image’s vertical distance from the origin as a percentage of the video height
	YPos *int64 `json:"YPos,omitnil,omitempty" name:"YPos"`

	// The watermark image’s width as a percentage of the video width
	Width *int64 `json:"Width,omitnil,omitempty" name:"Width"`

	// The watermark image’s height as a percentage of the video height
	Height *int64 `json:"Height,omitnil,omitempty" name:"Height"`
}

// Predefined struct for user
type DescribeStreamLiveChannelAlertsRequestParams struct {
	// Channel ID
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`
}

type DescribeStreamLiveChannelAlertsRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`
}

func (r *DescribeStreamLiveChannelAlertsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelAlertsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChannelId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveChannelAlertsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelAlertsResponseParams struct {
	// Alarm information of the channel's two pipelines
	Infos *ChannelAlertInfos `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveChannelAlertsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveChannelAlertsResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveChannelAlertsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelAlertsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelInputStatisticsRequestParams struct {
	// Channel ID
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query statistics in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Data collection interval. Valid values: `5s`, `1min` (default), `5min`, `15min`
	Period *string `json:"Period,omitnil,omitempty" name:"Period"`
}

type DescribeStreamLiveChannelInputStatisticsRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query statistics in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Data collection interval. Valid values: `5s`, `1min` (default), `5min`, `15min`
	Period *string `json:"Period,omitnil,omitempty" name:"Period"`
}

func (r *DescribeStreamLiveChannelInputStatisticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelInputStatisticsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChannelId")
	delete(f, "StartTime")
	delete(f, "EndTime")
	delete(f, "Period")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveChannelInputStatisticsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelInputStatisticsResponseParams struct {
	// Channel input statistics
	Infos []*ChannelInputStatistics `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveChannelInputStatisticsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveChannelInputStatisticsResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveChannelInputStatisticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelInputStatisticsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelLogsRequestParams struct {
	// Channel ID
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query logs in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`
}

type DescribeStreamLiveChannelLogsRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query logs in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`
}

func (r *DescribeStreamLiveChannelLogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelLogsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChannelId")
	delete(f, "StartTime")
	delete(f, "EndTime")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveChannelLogsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelLogsResponseParams struct {
	// Pipeline push information
	Infos *PipelineLogInfo `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveChannelLogsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveChannelLogsResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveChannelLogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelLogsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelOutputStatisticsRequestParams struct {
	// Channel ID
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query statistics in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Data collection interval. Valid values: `5s`, `1min` (default), `5min`, `15min`
	Period *string `json:"Period,omitnil,omitempty" name:"Period"`
}

type DescribeStreamLiveChannelOutputStatisticsRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query statistics in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Data collection interval. Valid values: `5s`, `1min` (default), `5min`, `15min`
	Period *string `json:"Period,omitnil,omitempty" name:"Period"`
}

func (r *DescribeStreamLiveChannelOutputStatisticsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelOutputStatisticsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChannelId")
	delete(f, "StartTime")
	delete(f, "EndTime")
	delete(f, "Period")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveChannelOutputStatisticsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelOutputStatisticsResponseParams struct {
	// Channel output information
	Infos []*ChannelOutputsStatistics `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveChannelOutputStatisticsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveChannelOutputStatisticsResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveChannelOutputStatisticsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelOutputStatisticsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type DescribeStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *DescribeStreamLiveChannelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelResponseParams struct {
	// Channel information
	Info *StreamLiveChannelInfo `json:"Info,omitnil,omitempty" name:"Info"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveChannelResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveChannelResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveChannelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelsRequestParams struct {

}

type DescribeStreamLiveChannelsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeStreamLiveChannelsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveChannelsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveChannelsResponseParams struct {
	// List of channel information
	// Note: this field may return `null`, indicating that no valid value was found.
	Infos []*StreamLiveChannelInfo `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveChannelsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveChannelsResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveChannelsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveChannelsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveInputRequestParams struct {
	// Input ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type DescribeStreamLiveInputRequest struct {
	*tchttp.BaseRequest
	
	// Input ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *DescribeStreamLiveInputRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveInputRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveInputRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveInputResponseParams struct {
	// Input information
	Info *InputInfo `json:"Info,omitnil,omitempty" name:"Info"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveInputResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveInputResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveInputResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveInputResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveInputSecurityGroupRequestParams struct {
	// Input security group ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type DescribeStreamLiveInputSecurityGroupRequest struct {
	*tchttp.BaseRequest
	
	// Input security group ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *DescribeStreamLiveInputSecurityGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveInputSecurityGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveInputSecurityGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveInputSecurityGroupResponseParams struct {
	// Input security group information
	Info *InputSecurityGroupInfo `json:"Info,omitnil,omitempty" name:"Info"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveInputSecurityGroupResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveInputSecurityGroupResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveInputSecurityGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveInputSecurityGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveInputSecurityGroupsRequestParams struct {

}

type DescribeStreamLiveInputSecurityGroupsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeStreamLiveInputSecurityGroupsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveInputSecurityGroupsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveInputSecurityGroupsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveInputSecurityGroupsResponseParams struct {
	// List of input security group information
	Infos []*InputSecurityGroupInfo `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveInputSecurityGroupsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveInputSecurityGroupsResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveInputSecurityGroupsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveInputSecurityGroupsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveInputsRequestParams struct {

}

type DescribeStreamLiveInputsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeStreamLiveInputsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveInputsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveInputsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveInputsResponseParams struct {
	// List of input information
	// Note: this field may return `null`, indicating that no valid value was found.
	Infos []*InputInfo `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveInputsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveInputsResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveInputsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveInputsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLivePlansRequestParams struct {
	// ID of the channel whose events you want to query
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`
}

type DescribeStreamLivePlansRequest struct {
	*tchttp.BaseRequest
	
	// ID of the channel whose events you want to query
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`
}

func (r *DescribeStreamLivePlansRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLivePlansRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "ChannelId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLivePlansRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLivePlansResponseParams struct {
	// List of event information
	// Note: this field may return `null`, indicating that no valid value was found.
	Infos []*PlanResp `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLivePlansResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLivePlansResponseParams `json:"Response"`
}

func (r *DescribeStreamLivePlansResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLivePlansResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveRegionsRequestParams struct {

}

type DescribeStreamLiveRegionsRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeStreamLiveRegionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveRegionsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveRegionsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveRegionsResponseParams struct {
	// StreamLive region information
	Info *StreamLiveRegionInfo `json:"Info,omitnil,omitempty" name:"Info"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveRegionsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveRegionsResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveRegionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveRegionsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveTranscodeDetailRequestParams struct {
	// The query start time (UTC+8) in the format of yyyy-MM-dd.
	// You can only query data in the last month (not including the current day).
	StartDayTime *string `json:"StartDayTime,omitnil,omitempty" name:"StartDayTime"`

	// The query end time (UTC+8) in the format of yyyy-MM-dd.
	// You can only query data in the last month (not including the current day).
	EndDayTime *string `json:"EndDayTime,omitnil,omitempty" name:"EndDayTime"`

	// The channel ID (optional).
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// The number of pages. Default value: 1.
	// The value cannot exceed 100.
	PageNum *int64 `json:"PageNum,omitnil,omitempty" name:"PageNum"`

	// The number of records per page. Default value: 10.
	// Value range: 1-1000.
	PageSize *int64 `json:"PageSize,omitnil,omitempty" name:"PageSize"`
}

type DescribeStreamLiveTranscodeDetailRequest struct {
	*tchttp.BaseRequest
	
	// The query start time (UTC+8) in the format of yyyy-MM-dd.
	// You can only query data in the last month (not including the current day).
	StartDayTime *string `json:"StartDayTime,omitnil,omitempty" name:"StartDayTime"`

	// The query end time (UTC+8) in the format of yyyy-MM-dd.
	// You can only query data in the last month (not including the current day).
	EndDayTime *string `json:"EndDayTime,omitnil,omitempty" name:"EndDayTime"`

	// The channel ID (optional).
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// The number of pages. Default value: 1.
	// The value cannot exceed 100.
	PageNum *int64 `json:"PageNum,omitnil,omitempty" name:"PageNum"`

	// The number of records per page. Default value: 10.
	// Value range: 1-1000.
	PageSize *int64 `json:"PageSize,omitnil,omitempty" name:"PageSize"`
}

func (r *DescribeStreamLiveTranscodeDetailRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveTranscodeDetailRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "StartDayTime")
	delete(f, "EndDayTime")
	delete(f, "ChannelId")
	delete(f, "PageNum")
	delete(f, "PageSize")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveTranscodeDetailRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveTranscodeDetailResponseParams struct {
	// A list of the transcoding information.
	Infos []*DescribeTranscodeDetailInfo `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The number of the current page.
	PageNum *int64 `json:"PageNum,omitnil,omitempty" name:"PageNum"`

	// The number of records per page.
	PageSize *int64 `json:"PageSize,omitnil,omitempty" name:"PageSize"`

	// The total number of records.
	TotalNum *int64 `json:"TotalNum,omitnil,omitempty" name:"TotalNum"`

	// The total number of pages.
	TotalPage *int64 `json:"TotalPage,omitnil,omitempty" name:"TotalPage"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveTranscodeDetailResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveTranscodeDetailResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveTranscodeDetailResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveTranscodeDetailResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveWatermarkRequestParams struct {
	// Watermark ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type DescribeStreamLiveWatermarkRequest struct {
	*tchttp.BaseRequest
	
	// Watermark ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *DescribeStreamLiveWatermarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveWatermarkRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveWatermarkRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveWatermarkResponseParams struct {
	// Watermark information
	Info *DescribeWatermarkInfo `json:"Info,omitnil,omitempty" name:"Info"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveWatermarkResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveWatermarkResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveWatermarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveWatermarkResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveWatermarksRequestParams struct {

}

type DescribeStreamLiveWatermarksRequest struct {
	*tchttp.BaseRequest
	
}

func (r *DescribeStreamLiveWatermarksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveWatermarksRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeStreamLiveWatermarksRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeStreamLiveWatermarksResponseParams struct {
	// List of watermark information
	Infos []*DescribeWatermarkInfo `json:"Infos,omitnil,omitempty" name:"Infos"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeStreamLiveWatermarksResponse struct {
	*tchttp.BaseResponse
	Response *DescribeStreamLiveWatermarksResponseParams `json:"Response"`
}

func (r *DescribeStreamLiveWatermarksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeStreamLiveWatermarksResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeTextSettings struct {
	// Text
	Text *string `json:"Text,omitnil,omitempty" name:"Text"`

	// Origin
	Location *string `json:"Location,omitnil,omitempty" name:"Location"`

	// The watermark image's horizontal distance from the origin as a percentage of the video width
	XPos *int64 `json:"XPos,omitnil,omitempty" name:"XPos"`

	// The watermark image's vertical distance from the origin as a percentage of the video height
	YPos *int64 `json:"YPos,omitnil,omitempty" name:"YPos"`

	// Font size
	FontSize *int64 `json:"FontSize,omitnil,omitempty" name:"FontSize"`

	// Font color
	FontColor *string `json:"FontColor,omitnil,omitempty" name:"FontColor"`
}

type DescribeTranscodeDetailInfo struct {
	// The channel ID.
	ChannelId *string `json:"ChannelId,omitnil,omitempty" name:"ChannelId"`

	// The start time (UTC+8) of transcoding in the format of yyyy-MM-dd HH:mm:ss.
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// The end time (UTC+8) of transcoding in the format of yyyy-MM-dd HH:mm:ss.
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// The duration (s) of transcoding.
	Duration *int64 `json:"Duration,omitnil,omitempty" name:"Duration"`

	// The encoding method.
	// Examples:
	// `liveprocessor_H264`: Live transcoding-H264
	// `liveprocessor_H265`: Live transcoding-H265
	// `topspeed_H264`: Top speed codec-H264
	// `topspeed_H265`: Top speed codec-H265
	ModuleCodec *string `json:"ModuleCodec,omitnil,omitempty" name:"ModuleCodec"`

	// The target bitrate (Kbps).
	Bitrate *int64 `json:"Bitrate,omitnil,omitempty" name:"Bitrate"`

	// The transcoding type.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// The push domain name.
	PushDomain *string `json:"PushDomain,omitnil,omitempty" name:"PushDomain"`

	// The target resolution.
	Resolution *string `json:"Resolution,omitnil,omitempty" name:"Resolution"`
}

// Predefined struct for user
type DescribeWatermarkDetectionRequestParams struct {
	// Task Id
	TaskId *string `json:"TaskId,omitnil,omitempty" name:"TaskId"`
}

type DescribeWatermarkDetectionRequest struct {
	*tchttp.BaseRequest
	
	// Task Id
	TaskId *string `json:"TaskId,omitnil,omitempty" name:"TaskId"`
}

func (r *DescribeWatermarkDetectionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeWatermarkDetectionRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "TaskId")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeWatermarkDetectionRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeWatermarkDetectionResponseParams struct {
	// Detecting task related information
	TaskInfo *AbWatermarkDetectionInfo `json:"TaskInfo,omitnil,omitempty" name:"TaskInfo"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeWatermarkDetectionResponse struct {
	*tchttp.BaseResponse
	Response *DescribeWatermarkDetectionResponseParams `json:"Response"`
}

func (r *DescribeWatermarkDetectionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeWatermarkDetectionResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeWatermarkDetectionsRequestParams struct {
	// Start time, 2022-12-04T16:50:00+08:00
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// End time, 2022-12-04T17:50:00+08:00, maximum supported query range of 7 days
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Number of pages queried
	PageNum *int64 `json:"PageNum,omitnil,omitempty" name:"PageNum"`

	// Single page quantity, 1-100
	PageSize *int64 `json:"PageSize,omitnil,omitempty" name:"PageSize"`
}

type DescribeWatermarkDetectionsRequest struct {
	*tchttp.BaseRequest
	
	// Start time, 2022-12-04T16:50:00+08:00
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// End time, 2022-12-04T17:50:00+08:00, maximum supported query range of 7 days
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Number of pages queried
	PageNum *int64 `json:"PageNum,omitnil,omitempty" name:"PageNum"`

	// Single page quantity, 1-100
	PageSize *int64 `json:"PageSize,omitnil,omitempty" name:"PageSize"`
}

func (r *DescribeWatermarkDetectionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeWatermarkDetectionsRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "StartTime")
	delete(f, "EndTime")
	delete(f, "PageNum")
	delete(f, "PageSize")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "DescribeWatermarkDetectionsRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type DescribeWatermarkDetectionsResponseParams struct {
	// Watermark detection information
	TaskInfos []*AbWatermarkDetectionInfo `json:"TaskInfos,omitnil,omitempty" name:"TaskInfos"`

	// number of tasks
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type DescribeWatermarkDetectionsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeWatermarkDetectionsResponseParams `json:"Response"`
}

func (r *DescribeWatermarkDetectionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *DescribeWatermarkDetectionsResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type DescribeWatermarkInfo struct {
	// Watermark ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Watermark name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Watermark type. Valid values: STATIC_IMAGE, TEXT.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	// Note: This field may return `null`, indicating that no valid value was found.
	ImageSettings *DescribeImageSettings `json:"ImageSettings,omitnil,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	// Note: This field may return `null`, indicating that no valid value was found.
	TextSettings *DescribeTextSettings `json:"TextSettings,omitnil,omitempty" name:"TextSettings"`

	// Last modified time (UTC+0) of the watermark, in the format of `2020-01-01T12:00:00Z`
	// Note: This field may return `null`, indicating that no valid value was found.
	UpdateTime *string `json:"UpdateTime,omitnil,omitempty" name:"UpdateTime"`

	// List of channel IDs the watermark is bound to
	// Note: This field may return `null`, indicating that no valid value was found.
	AttachedChannels []*string `json:"AttachedChannels,omitnil,omitempty" name:"AttachedChannels"`

	// AB watermark configuration.
	AbWatermarkSettings *AbWatermarkSettingsResp `json:"AbWatermarkSettings,omitnil,omitempty" name:"AbWatermarkSettings"`
}

type DestinationInfo struct {
	// Relay destination address. Length limit: [1,512].
	OutputUrl *string `json:"OutputUrl,omitnil,omitempty" name:"OutputUrl"`

	// Authentication key. Length limit: [1,128].
	// Note: this field may return null, indicating that no valid values can be obtained.
	AuthKey *string `json:"AuthKey,omitnil,omitempty" name:"AuthKey"`

	// Authentication username. Length limit: [1,128].
	// Note: this field may return null, indicating that no valid values can be obtained.
	Username *string `json:"Username,omitnil,omitempty" name:"Username"`

	// Authentication password. Length limit: [1,128].
	// Note: this field may return null, indicating that no valid values can be obtained.
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// The destination type of the retweet. Currently available values are: Standard, AWS_MediaPackageV1, AWS_MediaPackageV2. The default is: Standard. When the output group type is FRAME_CAPTURE, valid values are: AWS_AmazonS3, COS.
	DestinationType *string `json:"DestinationType,omitnil,omitempty" name:"DestinationType"`

	// Aws S3 destination setting.
	AmazonS3Settings *AmazonS3Settings `json:"AmazonS3Settings,omitnil,omitempty" name:"AmazonS3Settings"`

	// Cos destination setting.
	CosSettings *CosSettings `json:"CosSettings,omitnil,omitempty" name:"CosSettings"`
}

type DrmKey struct {
	// DRM key, which is a 32-bit hexadecimal string.
	// Note: uppercase letters in the string will be automatically converted to lowercase ones.
	Key *string `json:"Key,omitnil,omitempty" name:"Key"`

	// Required for Widevine encryption. Valid values: SD, HD, UHD1, UHD2, AUDIO, ALL.
	// ALL refers to all tracks. If this parameter is set to ALL, no other tracks can be added.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Track *string `json:"Track,omitnil,omitempty" name:"Track"`

	// Required for Widevine encryption. It is a 32-bit hexadecimal string.
	// Note: uppercase letters in the string will be automatically converted to lowercase ones.
	// Note: this field may return null, indicating that no valid values can be obtained.
	KeyId *string `json:"KeyId,omitnil,omitempty" name:"KeyId"`

	// Required when FairPlay uses the AES encryption method. It is a 32-bit hexadecimal string.
	// For more information about this parameter, please see: 
	// https://tools.ietf.org/html/rfc3826
	// Note: uppercase letters in the string will be automatically converted to lowercase ones.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Iv *string `json:"Iv,omitnil,omitempty" name:"Iv"`

	// The URI of the license server when AES-128 is used. This parameter may be empty.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	KeyUri *string `json:"KeyUri,omitnil,omitempty" name:"KeyUri"`
}

type DrmSettingsInfo struct {
	// Whether to enable DRM encryption. Valid values: `CLOSE` (disable), `OPEN` (enable). Default value: `CLOSE`
	// DRM encryption is supported only for HLS, DASH, HLS_ARCHIVE, DASH_ARCHIVE, HLS_MEDIAPACKAGE, and DASH_MEDIAPACKAGE outputs.
	State *string `json:"State,omitnil,omitempty" name:"State"`

	// Valid values: `CustomDRMKeys` (default value), `SDMCDRM`
	// `CustomDRMKeys` means encryption keys customized by users.
	// `SDMCDRM` means the DRM key management system of SDMC.
	Scheme *string `json:"Scheme,omitnil,omitempty" name:"Scheme"`

	// If `Scheme` is set to `CustomDRMKeys`, this parameter is required.
	// If `Scheme` is set to `SDMCDRM`, this parameter is optional. It supports digits, letters, hyphens, and underscores and must contain 1 to 36 characters. If it is not specified, the value of `ChannelId` will be used.
	ContentId *string `json:"ContentId,omitnil,omitempty" name:"ContentId"`

	// The key customized by the content user, which is required when `Scheme` is set to CustomDRMKeys.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Keys []*DrmKey `json:"Keys,omitnil,omitempty" name:"Keys"`

	// SDMC key configuration. This parameter is used when `Scheme` is set to `SDMCDRM`.
	// Note: This field may return `null`, indicating that no valid value was found.
	SDMCSettings *SDMCSettingsInfo `json:"SDMCSettings,omitnil,omitempty" name:"SDMCSettings"`

	// Optional Types:
	// `FAIRPLAY`, `WIDEVINE`, `PLAYREADY`, `AES128`
	// 
	// HLS-TS supports `FAIRPLAY` and `AES128`.
	// 
	// HLS-FMP4 supports `FAIRPLAY`, `WIDEVINE`, `PLAYREADY`, `AES128`, and combinations of two or three from `FAIRPLAY`, `WIDEVINE`, and `PLAYREADY` (concatenated with commas, e.g., "FAIRPLAY,WIDEVINE,PLAYREADY").
	// 
	// DASH supports `WIDEVINE`, `PLAYREADY`, and combinations of `PLAYREADY` and `WIDEVINE` (concatenated with commas, e.g., "PLAYREADY,WIDEVINE").
	DrmType *string `json:"DrmType,omitnil,omitempty" name:"DrmType"`
}

type EventNotifySetting struct {
	// The callback configuration for push events.
	PushEventSettings *PushEventSetting `json:"PushEventSettings,omitnil,omitempty" name:"PushEventSettings"`
}

type EventSettingsDestinationReq struct {
	// URL of the COS bucket to save recording files
	Url *string `json:"Url,omitnil,omitempty" name:"Url"`
}

type EventSettingsDestinationResp struct {
	// URL of the COS bucket where recording files are saved
	Url *string `json:"Url,omitnil,omitempty" name:"Url"`
}

type EventSettingsReq struct {
	// Valid values: `INPUT_SWITCH`, `TIMED_RECORD`, `SCTE35_TIME_SIGNAL`, `SCTE35_SPLICE_INSERT`, `SCTE35_RETURN_TO_NETWORK`,`TIMED_METADATA `,`STATIC_IMAGE_ACTIVATE `,`STATIC_IMAGE_DEACTIVATE `. If it is not specified, `INPUT_SWITCH` will be used.
	EventType *string `json:"EventType,omitnil,omitempty" name:"EventType"`

	// ID of the input to attach, which is required if `EventType` is `INPUT_SWITCH`
	InputAttachment *string `json:"InputAttachment,omitnil,omitempty" name:"InputAttachment"`

	// When the type is FIXED_PTS, it is mandatory and defaults to 0
	PipelineId *int64 `json:"PipelineId,omitnil,omitempty" name:"PipelineId"`

	// Name of the output group to attach. This parameter is required if `EventType` is `TIMED_RECORD`.
	OutputGroupName *string `json:"OutputGroupName,omitnil,omitempty" name:"OutputGroupName"`

	// Name of the manifest file for timed recording, which must end with `.m3u8` for HLS and `.mpd` for DASH. This parameter is required if `EventType` is `TIMED_RECORD`.
	ManifestName *string `json:"ManifestName,omitnil,omitempty" name:"ManifestName"`

	// URL of the COS bucket to save recording files. This parameter is required if `EventType` is `TIMED_RECORD`. It may contain 1 or 2 URLs. The first URL corresponds to pipeline 0 and the second pipeline 1.
	Destinations []*EventSettingsDestinationReq `json:"Destinations,omitnil,omitempty" name:"Destinations"`

	// SCTE-35 configuration information.
	SCTE35SegmentationDescriptor []*SegmentationDescriptorInfo `json:"SCTE35SegmentationDescriptor,omitnil,omitempty" name:"SCTE35SegmentationDescriptor"`

	// A 32-bit unique segmentation event identifier.Only one occurrence of a given segmentation_event_id value shall be active at any one time.
	SpliceEventID *uint64 `json:"SpliceEventID,omitnil,omitempty" name:"SpliceEventID"`

	// The duration of the segment in 90kHz ticks.It used to  give the splicer an indication of when the break will be over and when the network In Point will occur. If not specifyed,the splice_insert will continue when enter a return_to_network to end the splice_insert at the appropriate time.
	SpliceDuration *uint64 `json:"SpliceDuration,omitnil,omitempty" name:"SpliceDuration"`

	// Meta information plan configuration.
	TimedMetadataSetting *TimedMetadataInfo `json:"TimedMetadataSetting,omitnil,omitempty" name:"TimedMetadataSetting"`

	// Static image activate setting.
	StaticImageActivateSetting *StaticImageActivateSetting `json:"StaticImageActivateSetting,omitnil,omitempty" name:"StaticImageActivateSetting"`

	// Static image deactivate setting.
	StaticImageDeactivateSetting *StaticImageDeactivateSetting `json:"StaticImageDeactivateSetting,omitnil,omitempty" name:"StaticImageDeactivateSetting"`

	// Dynamic graphic overlay activate configuration
	MotionGraphicsActivateSetting *MotionGraphicsActivateSetting `json:"MotionGraphicsActivateSetting,omitnil,omitempty" name:"MotionGraphicsActivateSetting"`

	// Ad Settings
	AdBreakSetting *AdBreakSetting `json:"AdBreakSetting,omitnil,omitempty" name:"AdBreakSetting"`
}

type EventSettingsResp struct {
	// Valid values: `INPUT_SWITCH`, `TIMED_RECORD`, `SCTE35_TIME_SIGNAL`, `SCTE35_SPLICE_INSERT`, `SCTE35_RETURN_TO_NETWORK`, `STATIC_IMAGE_ACTIVATE`, `STATIC_IMAGE_DEACTIVATE`.
	EventType *string `json:"EventType,omitnil,omitempty" name:"EventType"`

	// ID of the input attached, which is not empty if `EventType` is `INPUT_SWITCH`
	InputAttachment *string `json:"InputAttachment,omitnil,omitempty" name:"InputAttachment"`

	// When the type is FIXED_PTS, it is mandatory and defaults to 0
	PipelineId *int64 `json:"PipelineId,omitnil,omitempty" name:"PipelineId"`

	// Name of the output group attached. This parameter is not empty if `EventType` is `TIMED_RECORD`.
	OutputGroupName *string `json:"OutputGroupName,omitnil,omitempty" name:"OutputGroupName"`

	// Name of the manifest file for timed recording, which ends with `.m3u8` for HLS and `.mpd` for DASH. This parameter is not empty if `EventType` is `TIMED_RECORD`.
	ManifestName *string `json:"ManifestName,omitnil,omitempty" name:"ManifestName"`

	// URL of the COS bucket where recording files are saved. This parameter is not empty if `EventType` is `TIMED_RECORD`. It may contain 1 or 2 URLs. The first URL corresponds to pipeline 0 and the second pipeline 1.
	Destinations []*EventSettingsDestinationResp `json:"Destinations,omitnil,omitempty" name:"Destinations"`

	// SCTE-35 configuration information.
	SCTE35SegmentationDescriptor []*SegmentationDescriptorRespInfo `json:"SCTE35SegmentationDescriptor,omitnil,omitempty" name:"SCTE35SegmentationDescriptor"`

	// A 32-bit unique segmentation event identifier.Only one occurrence of a given segmentation_event_id value shall be active at any one time.
	SpliceEventID *uint64 `json:"SpliceEventID,omitnil,omitempty" name:"SpliceEventID"`

	// The duration of the segment in 90kHz ticks.It used to  give the splicer an indication of when the break will be over and when the network In Point will occur. If not specifyed,the splice_insert will continue when enter a return_to_network to end the splice_insert at the appropriate time.
	SpliceDuration *string `json:"SpliceDuration,omitnil,omitempty" name:"SpliceDuration"`

	// Meta information plan configuration.
	TimedMetadataSetting *TimedMetadataInfo `json:"TimedMetadataSetting,omitnil,omitempty" name:"TimedMetadataSetting"`

	// Static image activate setting.
	StaticImageActivateSetting *StaticImageActivateSetting `json:"StaticImageActivateSetting,omitnil,omitempty" name:"StaticImageActivateSetting"`

	// Static image deactivate setting.
	StaticImageDeactivateSetting *StaticImageDeactivateSetting `json:"StaticImageDeactivateSetting,omitnil,omitempty" name:"StaticImageDeactivateSetting"`

	// Dynamic graphic overlay activate configuration.
	MotionGraphicsActivateSetting *MotionGraphicsActivateSetting `json:"MotionGraphicsActivateSetting,omitnil,omitempty" name:"MotionGraphicsActivateSetting"`

	// Ad Settings
	AdBreakSetting *AdBreakSetting `json:"AdBreakSetting,omitnil,omitempty" name:"AdBreakSetting"`
}

type FailOverSettings struct {
	// ID of the backup input
	// Note: this field may return `null`, indicating that no valid value was found.
	SecondaryInputId *string `json:"SecondaryInputId,omitnil,omitempty" name:"SecondaryInputId"`

	// The wait time (ms) for triggering failover after the primary input becomes unavailable. Value range: [1000, 86400000]. Default value: `3000`
	LossThreshold *int64 `json:"LossThreshold,omitnil,omitempty" name:"LossThreshold"`

	// Failover policy. Valid values: `CURRENT_PREFERRED` (default), `PRIMARY_PREFERRED`
	RecoverBehavior *string `json:"RecoverBehavior,omitnil,omitempty" name:"RecoverBehavior"`
}

type FrameCaptureTemplate struct {
	// Name of frame capture template, limited to uppercase and lowercase letters and numbers, with a length between 1 and 20 characters.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Width of frame capture, optional, input range is from 0 to 3000, must be a multiple of 2.
	Width *uint64 `json:"Width,omitnil,omitempty" name:"Width"`

	// Height of frame capture, optional, input range is from 0 to 3000, must be a multiple of 2.
	Height *uint64 `json:"Height,omitnil,omitempty" name:"Height"`

	// Interval of frame capture, an integer between 1 and 3600.
	CaptureInterval *uint64 `json:"CaptureInterval,omitnil,omitempty" name:"CaptureInterval"`

	// Interval units of frame capture, only supports SECONDS.
	CaptureIntervalUnits *string `json:"CaptureIntervalUnits,omitnil,omitempty" name:"CaptureIntervalUnits"`

	// Scaling behavior of frame capture, supports DEFAULT or STRETCH_TO_OUTPUT, with DEFAULT being the default option.
	ScalingBehavior *string `json:"ScalingBehavior,omitnil,omitempty" name:"ScalingBehavior"`

	// Sharpness, an integer between 0 and 100.
	Sharpness *uint64 `json:"Sharpness,omitnil,omitempty" name:"Sharpness"`
}

type GeneralSetting struct {
	// Static graphic overlay configuration.
	StaticImageSettings *StaticImageSettings `json:"StaticImageSettings,omitnil,omitempty" name:"StaticImageSettings"`

	// Dynamic graphic overlay configuration.
	MotionGraphicsSettings *MotionGraphicsSetting `json:"MotionGraphicsSettings,omitnil,omitempty" name:"MotionGraphicsSettings"`

	// Thumbnail Configuration.
	ThumbnailSettings *ThumbnailSettings `json:"ThumbnailSettings,omitnil,omitempty" name:"ThumbnailSettings"`
}

// Predefined struct for user
type GetAbWatermarkPlayUrlRequestParams struct {
	// Client UUID, 32-bit unsigned integer, [0, 4294967295].
	Uuid *uint64 `json:"Uuid,omitnil,omitempty" name:"Uuid"`

	// Channel ID of Stream Package.
	StreamPackageChannelId *string `json:"StreamPackageChannelId,omitnil,omitempty" name:"StreamPackageChannelId"`

	// Original play URL.
	Url *string `json:"Url,omitnil,omitempty" name:"Url"`
}

type GetAbWatermarkPlayUrlRequest struct {
	*tchttp.BaseRequest
	
	// Client UUID, 32-bit unsigned integer, [0, 4294967295].
	Uuid *uint64 `json:"Uuid,omitnil,omitempty" name:"Uuid"`

	// Channel ID of Stream Package.
	StreamPackageChannelId *string `json:"StreamPackageChannelId,omitnil,omitempty" name:"StreamPackageChannelId"`

	// Original play URL.
	Url *string `json:"Url,omitnil,omitempty" name:"Url"`
}

func (r *GetAbWatermarkPlayUrlRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *GetAbWatermarkPlayUrlRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Uuid")
	delete(f, "StreamPackageChannelId")
	delete(f, "Url")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "GetAbWatermarkPlayUrlRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type GetAbWatermarkPlayUrlResponseParams struct {
	// The play URL after adding token.
	Url *string `json:"Url,omitnil,omitempty" name:"Url"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type GetAbWatermarkPlayUrlResponse struct {
	*tchttp.BaseResponse
	Response *GetAbWatermarkPlayUrlResponseParams `json:"Response"`
}

func (r *GetAbWatermarkPlayUrlResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *GetAbWatermarkPlayUrlResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type HighlightInfo struct {
	// Whether to enable input recognition 0: Disable 1 Enable Default value 0 Disable.
	HighlightEnabled *uint64 `json:"HighlightEnabled,omitnil,omitempty" name:"HighlightEnabled"`

	// The product where the results are saved, optional: COS. Currently, only Tencent Cloud COS is supported. In the future, it will be connected to AWS S3 and COS will be used by default.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Valid when Type is COS, the region where COS is stored.
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// Valid when Type is COS, the bucket name stored in COS.
	Bucket *string `json:"Bucket,omitnil,omitempty" name:"Bucket"`

	// Valid when Type is COS, the path where cos is stored.
	Path *string `json:"Path,omitnil,omitempty" name:"Path"`

	// Valid when Type is COS, the file name stored in cos.
	Filename *string `json:"Filename,omitnil,omitempty" name:"Filename"`

	// Valid when Type is COS, the file name suffix stored in COS is automatically generated in the time format, optional values: unix, utc. Unix is the second-level timestamp and UTC is the year, month and day represented by the zero time zone.
	TimestampFormat *string `json:"TimestampFormat,omitnil,omitempty" name:"TimestampFormat"`

	// Audio selector list is optional and can be empty. If not filled in, an audio will be used as the output of the recognition result by default.
	AudioSelectorNames []*string `json:"AudioSelectorNames,omitnil,omitempty" name:"AudioSelectorNames"`
}

type HlsRemuxSettingsInfo struct {
	// Segment duration in ms. Value range: [1000,30000]. Default value: 4000. The value can only be a multiple of 1,000.
	SegmentDuration *uint64 `json:"SegmentDuration,omitnil,omitempty" name:"SegmentDuration"`

	// Number of segments. Value range: [3,30]. Default value: 5.
	SegmentNumber *uint64 `json:"SegmentNumber,omitnil,omitempty" name:"SegmentNumber"`

	// Whether to enable PDT insertion. Valid values: CLOSE/OPEN. Default value: CLOSE.
	PdtInsertion *string `json:"PdtInsertion,omitnil,omitempty" name:"PdtInsertion"`

	// PDT duration in seconds. Value range: (0,3000]. Default value: 600.
	PdtDuration *uint64 `json:"PdtDuration,omitnil,omitempty" name:"PdtDuration"`

	// Audio/Video packaging scheme. Valid values: `SEPARATE`, `MERGE`. Default value is: SEPARATE.
	Scheme *string `json:"Scheme,omitnil,omitempty" name:"Scheme"`

	// The segment type. Valid values: `ts` (default), `fmp4`.
	// Currently, fMP4 segments do not support DRM or time shifting.
	SegmentType *string `json:"SegmentType,omitnil,omitempty" name:"SegmentType"`

	// The HLS package type when the H.265 codec is used. Valid values: `hvc1`, `hev1` (default).
	H265PackageType *string `json:"H265PackageType,omitnil,omitempty" name:"H265PackageType"`

	// Whether to enable low latency 0:CLOSE, 1:OPEN, default value: 0.
	LowLatency *uint64 `json:"LowLatency,omitnil,omitempty" name:"LowLatency"`

	// Low latency slice size, unit ms. Value range: integer [200-HlsRemuxSettings.SegmentDuration] Default value: 500ms.
	PartialSegmentDuration *uint64 `json:"PartialSegmentDuration,omitnil,omitempty" name:"PartialSegmentDuration"`

	// Low latency slice playback position, unit ms. Value range: integer [3*HlsRemuxSettings.PartiSegmentDuration - 3*HlsRemuxSettings.SegmentDuration], Default value: 3*HlsRemuxSettings.PartiSegmentDuration.
	PartialSegmentPlaySite *uint64 `json:"PartialSegmentPlaySite,omitnil,omitempty" name:"PartialSegmentPlaySite"`

	// Hls main m3u8 file sorting rules by bitrate, optional values: 1: video bitrate ascending order; 2: video bitrate descending order. Default value: 1.
	StreamOrder *uint64 `json:"StreamOrder,omitnil,omitempty" name:"StreamOrder"`

	// Whether the Hls main m3u8 file contains resolution information, optional values: 1: INCLUDE includes video resolution; 2: EXCLUDE does not include video resolution. Default value: 1.
	VideoResolution *uint64 `json:"VideoResolution,omitnil,omitempty" name:"VideoResolution"`

	// Whether to include the `EXT-X-ENDLIST` tag, 1 includes  `EXT-X-ENDLIST` tag, 2 does not include  `EXT-X-ENDLIST` tag; the default value is 1.
	EndListTag *int64 `json:"EndListTag,omitnil,omitempty" name:"EndListTag"`

	// Optional: `ENHANCED_SCTE35`, `DATERANGE`; default value: `ENHANCED_SCTE35`.
	AdMarkupType *string `json:"AdMarkupType,omitnil,omitempty" name:"AdMarkupType"`
}

type InputAnalysisInfo struct {
	// Highlight configuration.
	HighlightSetting *HighlightInfo `json:"HighlightSetting,omitnil,omitempty" name:"HighlightSetting"`
}

type InputFileInfo struct {
	// Segment duration, in milliseconds, ranging from 1000-10000, must be a multiple of 1000. The input video duration should be between SegmentDuration * 90 and SegmentDuration * 180
	SegmentDuration *int64 `json:"SegmentDuration,omitnil,omitempty" name:"SegmentDuration"`
}

type InputInfo struct {
	// Input region.
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`

	// Input ID.
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Input name.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Input type.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Array of security groups associated with input.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Array of channels associated with input.
	// Note: this field may return null, indicating that no valid values can be obtained.
	AttachedChannels []*string `json:"AttachedChannels,omitnil,omitempty" name:"AttachedChannels"`

	// Input configuration array.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitnil,omitempty" name:"InputSettings"`
}

type InputLossBehaviorInfo struct {
	// The time to fill in the last video frame, unit ms, range 0-1000000, 1000000 means always inserting, default 0 means filling in black screen frame.
	RepeatLastFrameMs *uint64 `json:"RepeatLastFrameMs,omitnil,omitempty" name:"RepeatLastFrameMs"`

	// Fill frame type, COLOR means solid color filling, IMAGE means picture filling, the default is COLOR.
	InputLossImageType *string `json:"InputLossImageType,omitnil,omitempty" name:"InputLossImageType"`

	// When the type is COLOR, the corresponding rgb value
	ColorRGB *string `json:"ColorRGB,omitnil,omitempty" name:"ColorRGB"`

	// When the type is IMAGE, the corresponding image url value
	ImageUrl *string `json:"ImageUrl,omitnil,omitempty" name:"ImageUrl"`
}

type InputSecurityGroupInfo struct {
	// Input security group ID.
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Input security group name.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// List of allowlist entries.
	Whitelist []*string `json:"Whitelist,omitnil,omitempty" name:"Whitelist"`

	// List of bound input streams.
	// Note: this field may return null, indicating that no valid values can be obtained.
	OccupiedInputs []*string `json:"OccupiedInputs,omitnil,omitempty" name:"OccupiedInputs"`

	// Input security group address.
	Region *string `json:"Region,omitnil,omitempty" name:"Region"`
}

type InputSettingInfo struct {
	// Application name, which is valid if `Type` is `RTMP_PUSH` or `RTMPS_PUSH`, and can contain 1-32 letters and digits
	// Note: This field may return `null`, indicating that no valid value was found.
	AppName *string `json:"AppName,omitnil,omitempty" name:"AppName"`

	// Stream name, which is valid if `Type` is `RTMP_PUSH` or `RTMPS_PUSH`, and can contain 1-32 letters and digits
	// Note: This field may return `null`, indicating that no valid value was found.
	StreamName *string `json:"StreamName,omitnil,omitempty" name:"StreamName"`

	// Source URL, which is valid if `Type` is `RTMP_PULL`, `HLS_PULL`,  `MP4_PULL`, `RTSP_PULL` or `SRT_PULL`, and can contain 1-512 characters
	// Note: This field may return `null`, indicating that no valid value was found.
	SourceUrl *string `json:"SourceUrl,omitnil,omitempty" name:"SourceUrl"`

	// RTP/UDP input address, which does not need to be entered for the input parameter.
	// Note: this field may return null, indicating that no valid values can be obtained.
	InputAddress *string `json:"InputAddress,omitnil,omitempty" name:"InputAddress"`

	// Source type for stream pulling and relaying. To pull content from private-read COS buckets under the current account, set this parameter to `TencentCOS`; otherwise, leave it empty.
	// Note: this field may return `null`, indicating that no valid value was found.
	SourceType *string `json:"SourceType,omitnil,omitempty" name:"SourceType"`

	// Delayed time (ms) for playback, which is valid if `Type` is `RTMP_PUSH` or `RTMPS_PUSH`.
	// Value range: 0 (default) or 10000-600000.
	// The value must be a multiple of 1,000.
	// Note: This field may return `null`, indicating that no valid value was found.
	DelayTime *int64 `json:"DelayTime,omitnil,omitempty" name:"DelayTime"`

	// The domain name of the SRT_PUSH push address. No need to fill in the input parameter.
	InputDomain *string `json:"InputDomain,omitnil,omitempty" name:"InputDomain"`

	// The username, which is used for authentication.
	// Note: This field may return `null`, indicating that no valid value was found.
	UserName *string `json:"UserName,omitnil,omitempty" name:"UserName"`

	// The password, which is used for authentication.
	// Note: This field may return `null`, indicating that no valid value was found.
	Password *string `json:"Password,omitnil,omitempty" name:"Password"`

	// This parameter is valid when the input source is HLS_PULL and MP4_PULL. It indicates the type of file the source is. The optional values are: LIVE, VOD. Please note that if you do not enter this parameter, the system will take the default input value VOD.
	ContentType *string `json:"ContentType,omitnil,omitempty" name:"ContentType"`
}

type InputStatistics struct {
	// Input statistics of pipeline 0.
	Pipeline0 []*PipelineInputStatistics `json:"Pipeline0,omitnil,omitempty" name:"Pipeline0"`

	// Input statistics of pipeline 1.
	Pipeline1 []*PipelineInputStatistics `json:"Pipeline1,omitnil,omitempty" name:"Pipeline1"`
}

type InputStreamInfo struct {
	// The input stream address.
	InputAddress *string `json:"InputAddress,omitnil,omitempty" name:"InputAddress"`

	// The input stream path.
	AppName *string `json:"AppName,omitnil,omitempty" name:"AppName"`

	// The input stream name.
	StreamName *string `json:"StreamName,omitnil,omitempty" name:"StreamName"`

	// The input stream status. `1` indicates the stream is active.
	Status *int64 `json:"Status,omitnil,omitempty" name:"Status"`
}

type InputTrack struct {
	// Audio track index 1-based index mapping to the specified audio track integer starting from 1.
	TrackIndex *uint64 `json:"TrackIndex,omitnil,omitempty" name:"TrackIndex"`
}

type InputTracks struct {
	// Audio track configuration information.
	Tracks []*InputTrack `json:"Tracks,omitnil,omitempty" name:"Tracks"`
}

type LSqueezeSetting struct {
	// Advertising benchmark position, 0 top left, 1 top right, 2 bottom right, 3 bottom left, default value 0, corresponding TOP_LEFT,TOP_RIGHT,BOTTOM_RIGHT,BOTTOM_LEFT
	Location *uint64 `json:"Location,omitnil,omitempty" name:"Location"`

	// The default value for the percentage in the X-axis direction is 20, with a range of 0-50
	OffsetX *uint64 `json:"OffsetX,omitnil,omitempty" name:"OffsetX"`

	// The default value for the percentage in the Y-axis direction is 20, with a range of 0-50
	OffsetY *uint64 `json:"OffsetY,omitnil,omitempty" name:"OffsetY"`

	// Background image URL, starting with http/https and ending in jpg/jpeg/png
	BackgroundImgUrl *string `json:"BackgroundImgUrl,omitnil,omitempty" name:"BackgroundImgUrl"`

	// Compress time. Unit ms, default value 2000, range: 500-10000, SqueezeInPeriod+SqueezeOutPeriod cannot be greater than duration, included in duration
	SqueezeInPeriod *uint64 `json:"SqueezeInPeriod,omitnil,omitempty" name:"SqueezeInPeriod"`

	// Restore to full screen time. Unit ms, default value 2000, range 500-10000, SqueezeInPeriod+SqueezeOutPeriod cannot be greater than duration, included in duration
	SqueezeOutPeriod *uint64 `json:"SqueezeOutPeriod,omitnil,omitempty" name:"SqueezeOutPeriod"`
}

type LogInfo struct {
	// Log type.
	// It contains the value of `StreamStart` which refers to the push information.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Time when the log is printed.
	Time *string `json:"Time,omitnil,omitempty" name:"Time"`

	// Log details.
	Message *LogMessageInfo `json:"Message,omitnil,omitempty" name:"Message"`
}

type LogMessageInfo struct {
	// Push information.
	// Note: this field may return null, indicating that no valid values can be obtained.
	StreamInfo *StreamInfo `json:"StreamInfo,omitnil,omitempty" name:"StreamInfo"`
}

// Predefined struct for user
type ModifyStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Channel name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Inputs to attach. You can attach 1 to 5 inputs.
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitnil,omitempty" name:"AttachedInputs"`

	// Configuration information of the channel's output groups. Quantity: [1, 10]
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitnil,omitempty" name:"OutputGroups"`

	// Audio transcoding templates. Quantity: [1, 20]
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitnil,omitempty" name:"AudioTemplates"`

	// Video transcoding templates. Quantity: [1, 10]
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitnil,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates. Quantity: [1, 10]
	AVTemplates []*AVTemplate `json:"AVTemplates,omitnil,omitempty" name:"AVTemplates"`

	// Subtitle template configuration.
	CaptionTemplates []*SubtitleConf `json:"CaptionTemplates,omitnil,omitempty" name:"CaptionTemplates"`

	// Event settings
	PlanSettings *PlanSettings `json:"PlanSettings,omitnil,omitempty" name:"PlanSettings"`

	// The callback settings.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitnil,omitempty" name:"EventNotifySettings"`

	// Complement the last video frame settings.
	InputLossBehavior *InputLossBehaviorInfo `json:"InputLossBehavior,omitnil,omitempty" name:"InputLossBehavior"`

	// Pipeline configuration.
	PipelineInputSettings *PipelineInputSettingsInfo `json:"PipelineInputSettings,omitnil,omitempty" name:"PipelineInputSettings"`

	// Recognition configuration for input content.
	InputAnalysisSettings *InputAnalysisInfo `json:"InputAnalysisSettings,omitnil,omitempty" name:"InputAnalysisSettings"`

	// Console tag list.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`

	// Frame capture templates.
	FrameCaptureTemplates []*FrameCaptureTemplate `json:"FrameCaptureTemplates,omitnil,omitempty" name:"FrameCaptureTemplates"`

	// General settings.
	GeneralSettings *GeneralSetting `json:"GeneralSettings,omitnil,omitempty" name:"GeneralSettings"`
}

type ModifyStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Channel name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Inputs to attach. You can attach 1 to 5 inputs.
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitnil,omitempty" name:"AttachedInputs"`

	// Configuration information of the channel's output groups. Quantity: [1, 10]
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitnil,omitempty" name:"OutputGroups"`

	// Audio transcoding templates. Quantity: [1, 20]
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitnil,omitempty" name:"AudioTemplates"`

	// Video transcoding templates. Quantity: [1, 10]
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitnil,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates. Quantity: [1, 10]
	AVTemplates []*AVTemplate `json:"AVTemplates,omitnil,omitempty" name:"AVTemplates"`

	// Subtitle template configuration.
	CaptionTemplates []*SubtitleConf `json:"CaptionTemplates,omitnil,omitempty" name:"CaptionTemplates"`

	// Event settings
	PlanSettings *PlanSettings `json:"PlanSettings,omitnil,omitempty" name:"PlanSettings"`

	// The callback settings.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitnil,omitempty" name:"EventNotifySettings"`

	// Complement the last video frame settings.
	InputLossBehavior *InputLossBehaviorInfo `json:"InputLossBehavior,omitnil,omitempty" name:"InputLossBehavior"`

	// Pipeline configuration.
	PipelineInputSettings *PipelineInputSettingsInfo `json:"PipelineInputSettings,omitnil,omitempty" name:"PipelineInputSettings"`

	// Recognition configuration for input content.
	InputAnalysisSettings *InputAnalysisInfo `json:"InputAnalysisSettings,omitnil,omitempty" name:"InputAnalysisSettings"`

	// Console tag list.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`

	// Frame capture templates.
	FrameCaptureTemplates []*FrameCaptureTemplate `json:"FrameCaptureTemplates,omitnil,omitempty" name:"FrameCaptureTemplates"`

	// General settings.
	GeneralSettings *GeneralSetting `json:"GeneralSettings,omitnil,omitempty" name:"GeneralSettings"`
}

func (r *ModifyStreamLiveChannelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyStreamLiveChannelRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	delete(f, "Name")
	delete(f, "AttachedInputs")
	delete(f, "OutputGroups")
	delete(f, "AudioTemplates")
	delete(f, "VideoTemplates")
	delete(f, "AVTemplates")
	delete(f, "CaptionTemplates")
	delete(f, "PlanSettings")
	delete(f, "EventNotifySettings")
	delete(f, "InputLossBehavior")
	delete(f, "PipelineInputSettings")
	delete(f, "InputAnalysisSettings")
	delete(f, "Tags")
	delete(f, "FrameCaptureTemplates")
	delete(f, "GeneralSettings")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyStreamLiveChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveChannelResponseParams struct {
	// Tag prompt information, this information will be attached when the tag operation fails.
	TagMsg *string `json:"TagMsg,omitnil,omitempty" name:"TagMsg"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyStreamLiveChannelResponse struct {
	*tchttp.BaseResponse
	Response *ModifyStreamLiveChannelResponseParams `json:"Response"`
}

func (r *ModifyStreamLiveChannelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyStreamLiveChannelResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveInputRequestParams struct {
	// Input ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// List of the IDs of the security groups to attach
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Input settings. 
	// For the type:
	// `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`,`RTSP_PULL`,`SRT_PULL` or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
	// For the type:
	// `SRT_PUSH`, 0 or 2 inputs of the corresponding type can be configured.
	// This parameter can be left empty for RTP_PUSH and UDP_PUSH inputs.
	// 
	// Note: If this parameter is not specified or empty, the original input settings will be used.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitnil,omitempty" name:"InputSettings"`
}

type ModifyStreamLiveInputRequest struct {
	*tchttp.BaseRequest
	
	// Input ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// List of the IDs of the security groups to attach
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitnil,omitempty" name:"SecurityGroupIds"`

	// Input settings. 
	// For the type:
	// `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`,`RTSP_PULL`,`SRT_PULL` or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
	// For the type:
	// `SRT_PUSH`, 0 or 2 inputs of the corresponding type can be configured.
	// This parameter can be left empty for RTP_PUSH and UDP_PUSH inputs.
	// 
	// Note: If this parameter is not specified or empty, the original input settings will be used.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitnil,omitempty" name:"InputSettings"`
}

func (r *ModifyStreamLiveInputRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyStreamLiveInputRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	delete(f, "Name")
	delete(f, "SecurityGroupIds")
	delete(f, "InputSettings")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyStreamLiveInputRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveInputResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyStreamLiveInputResponse struct {
	*tchttp.BaseResponse
	Response *ModifyStreamLiveInputResponseParams `json:"Response"`
}

func (r *ModifyStreamLiveInputResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyStreamLiveInputResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveInputSecurityGroupRequestParams struct {
	// Input security group ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Input security group name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Allowlist entries (max: 10)
	Whitelist []*string `json:"Whitelist,omitnil,omitempty" name:"Whitelist"`
}

type ModifyStreamLiveInputSecurityGroupRequest struct {
	*tchttp.BaseRequest
	
	// Input security group ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Input security group name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Allowlist entries (max: 10)
	Whitelist []*string `json:"Whitelist,omitnil,omitempty" name:"Whitelist"`
}

func (r *ModifyStreamLiveInputSecurityGroupRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyStreamLiveInputSecurityGroupRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	delete(f, "Name")
	delete(f, "Whitelist")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyStreamLiveInputSecurityGroupRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveInputSecurityGroupResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyStreamLiveInputSecurityGroupResponse struct {
	*tchttp.BaseResponse
	Response *ModifyStreamLiveInputSecurityGroupResponseParams `json:"Response"`
}

func (r *ModifyStreamLiveInputSecurityGroupResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyStreamLiveInputSecurityGroupResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveWatermarkRequestParams struct {
	// Watermark ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Watermark name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	ImageSettings *CreateImageSettings `json:"ImageSettings,omitnil,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	TextSettings *CreateTextSettings `json:"TextSettings,omitnil,omitempty" name:"TextSettings"`


	AbWatermarkSettings *AbWatermarkSettingsReq `json:"AbWatermarkSettings,omitnil,omitempty" name:"AbWatermarkSettings"`
}

type ModifyStreamLiveWatermarkRequest struct {
	*tchttp.BaseRequest
	
	// Watermark ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Watermark name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	ImageSettings *CreateImageSettings `json:"ImageSettings,omitnil,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	TextSettings *CreateTextSettings `json:"TextSettings,omitnil,omitempty" name:"TextSettings"`

	AbWatermarkSettings *AbWatermarkSettingsReq `json:"AbWatermarkSettings,omitnil,omitempty" name:"AbWatermarkSettings"`
}

func (r *ModifyStreamLiveWatermarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyStreamLiveWatermarkRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	delete(f, "Name")
	delete(f, "ImageSettings")
	delete(f, "TextSettings")
	delete(f, "AbWatermarkSettings")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyStreamLiveWatermarkRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveWatermarkResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type ModifyStreamLiveWatermarkResponse struct {
	*tchttp.BaseResponse
	Response *ModifyStreamLiveWatermarkResponseParams `json:"Response"`
}

func (r *ModifyStreamLiveWatermarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *ModifyStreamLiveWatermarkResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type MotionGraphicsActivateSetting struct {
	// Duration in ms, valid when MOTION_Graphics_ACTIVATE, required; An integer in the range of 0-86400000, where 0 represents the duration until the end of the live stream.
	Duration *int64 `json:"Duration,omitnil,omitempty" name:"Duration"`

	// The address of HTML5 needs to comply with the format specification of http/https.
	Url *string `json:"Url,omitnil,omitempty" name:"Url"`
}

type MotionGraphicsSetting struct {
	// Whether to enable dynamic graphic overlay, '0' not enabled, '1' enabled; Default 0.
	MotionGraphicsOverlayEnabled *int64 `json:"MotionGraphicsOverlayEnabled,omitnil,omitempty" name:"MotionGraphicsOverlayEnabled"`
}

type OutputInfo struct {
	// Output name.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Audio transcoding template name array.
	// Quantity limit: [0,1] for RTMP; [0,20] for others.
	// Note: this field may return null, indicating that no valid values can be obtained.
	AudioTemplateNames []*string `json:"AudioTemplateNames,omitnil,omitempty" name:"AudioTemplateNames"`

	// Video transcoding template name array. Quantity limit: [0,1].
	// Note: this field may return null, indicating that no valid values can be obtained.
	VideoTemplateNames []*string `json:"VideoTemplateNames,omitnil,omitempty" name:"VideoTemplateNames"`

	// SCTE-35 information configuration.
	Scte35Settings *Scte35SettingsInfo `json:"Scte35Settings,omitnil,omitempty" name:"Scte35Settings"`

	// Audio/Video transcoding template name. If `HlsRemuxSettings.Scheme` is `MERGE`, there is 1 audio/video transcoding template. Otherwise, this parameter is empty.
	// Note: this field may return `null`, indicating that no valid value was found.
	AVTemplateNames []*string `json:"AVTemplateNames,omitnil,omitempty" name:"AVTemplateNames"`

	// For the subtitle template used, only the AVTemplateNames is valid.
	CaptionTemplateNames []*string `json:"CaptionTemplateNames,omitnil,omitempty" name:"CaptionTemplateNames"`

	// Meta information controls configuration.
	TimedMetadataSettings *TimedMetadataSettingInfo `json:"TimedMetadataSettings,omitnil,omitempty" name:"TimedMetadataSettings"`

	// Frame capture template name array. Quantity limit: [0,1].
	FrameCaptureTemplateNames []*string `json:"FrameCaptureTemplateNames,omitnil,omitempty" name:"FrameCaptureTemplateNames"`

	// Name modification for sub m3u8.
	NameModifier *string `json:"NameModifier,omitnil,omitempty" name:"NameModifier"`
}

type OutputsStatistics struct {
	// Output information of pipeline 0.
	Pipeline0 []*PipelineOutputStatistics `json:"Pipeline0,omitnil,omitempty" name:"Pipeline0"`

	// Output information of pipeline 1.
	Pipeline1 []*PipelineOutputStatistics `json:"Pipeline1,omitnil,omitempty" name:"Pipeline1"`
}

type PipelineInputSettingsInfo struct {
	// Pipeline failover configuration, the valid value is: 1.PIPELINE_FAILOVER (channels are mutually failover); 2.PIPELINE_FILLING (channels fill in themselves). Default value: PIPELINE_FILLING. The specific content is specified by FaultBehavior.
	FaultBehavior *string `json:"FaultBehavior,omitnil,omitempty" name:"FaultBehavior"`
}

type PipelineInputStatistics struct {
	// Data timestamp in seconds.
	Timestamp *uint64 `json:"Timestamp,omitnil,omitempty" name:"Timestamp"`

	// Input bandwidth in bps.
	NetworkIn *uint64 `json:"NetworkIn,omitnil,omitempty" name:"NetworkIn"`

	// Video information array.
	// For `rtp/udp` input, the quantity is the number of `Pid` of the input video.
	// For other inputs, the quantity is 1.
	Video []*VideoPipelineInputStatistics `json:"Video,omitnil,omitempty" name:"Video"`

	// Audio information array.
	// For `rtp/udp` input, the quantity is the number of `Pid` of the input audio.
	// For other inputs, the quantity is 1.
	Audio []*AudioPipelineInputStatistics `json:"Audio,omitnil,omitempty" name:"Audio"`

	// Session ID
	SessionId *string `json:"SessionId,omitnil,omitempty" name:"SessionId"`

	// Rtt time, in milliseconds
	RTT *int64 `json:"RTT,omitnil,omitempty" name:"RTT"`

	// Is the Network parameter valid? 0 indicates invalid, 1 indicates valid
	NetworkValid *int64 `json:"NetworkValid,omitnil,omitempty" name:"NetworkValid"`
}

type PipelineLogInfo struct {
	// Log information of pipeline 0.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pipeline0 []*LogInfo `json:"Pipeline0,omitnil,omitempty" name:"Pipeline0"`

	// Log information of pipeline 1.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pipeline1 []*LogInfo `json:"Pipeline1,omitnil,omitempty" name:"Pipeline1"`
}

type PipelineOutputStatistics struct {
	// Timestamp.
	// In seconds, indicating data time.
	Timestamp *uint64 `json:"Timestamp,omitnil,omitempty" name:"Timestamp"`

	// Output bandwidth in bps.
	NetworkOut *uint64 `json:"NetworkOut,omitnil,omitempty" name:"NetworkOut"`

	// Is the Network parameter valid? 0 indicates invalid, 1 indicates valid
	NetworkValid *int64 `json:"NetworkValid,omitnil,omitempty" name:"NetworkValid"`
}

type PlanReq struct {
	// Event name
	EventName *string `json:"EventName,omitnil,omitempty" name:"EventName"`

	// Event trigger time settings
	TimingSettings *TimingSettingsReq `json:"TimingSettings,omitnil,omitempty" name:"TimingSettings"`

	// Event configuration
	EventSettings *EventSettingsReq `json:"EventSettings,omitnil,omitempty" name:"EventSettings"`
}

type PlanResp struct {
	// Event name
	EventName *string `json:"EventName,omitnil,omitempty" name:"EventName"`

	// Event trigger time settings
	TimingSettings *TimingSettingsResp `json:"TimingSettings,omitnil,omitempty" name:"TimingSettings"`

	// Event configuration
	EventSettings *EventSettingsResp `json:"EventSettings,omitnil,omitempty" name:"EventSettings"`
}

type PlanSettings struct {
	// Timed recording settings
	// Note: This field may return `null`, indicating that no valid value was found.
	TimedRecordSettings *TimedRecordSettings `json:"TimedRecordSettings,omitnil,omitempty" name:"TimedRecordSettings"`
}

type PushEventSetting struct {
	// The callback URL (required).
	NotifyUrl *string `json:"NotifyUrl,omitnil,omitempty" name:"NotifyUrl"`

	// The callback key (optional).
	NotifyKey *string `json:"NotifyKey,omitnil,omitempty" name:"NotifyKey"`
}

type QueryDispatchInputInfo struct {
	// The input ID.
	InputID *string `json:"InputID,omitnil,omitempty" name:"InputID"`

	// The input name.
	InputName *string `json:"InputName,omitnil,omitempty" name:"InputName"`

	// The input protocol.
	Protocol *string `json:"Protocol,omitnil,omitempty" name:"Protocol"`

	// The stream status of the input.
	InputStreamInfoList []*InputStreamInfo `json:"InputStreamInfoList,omitnil,omitempty" name:"InputStreamInfoList"`
}

// Predefined struct for user
type QueryInputStreamStateRequestParams struct {
	// The StreamLive input ID.Currently, only RTMP_PUSH and RTMPS_PUSH are supported
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type QueryInputStreamStateRequest struct {
	*tchttp.BaseRequest
	
	// The StreamLive input ID.Currently, only RTMP_PUSH and RTMPS_PUSH are supported
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *QueryInputStreamStateRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *QueryInputStreamStateRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "QueryInputStreamStateRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type QueryInputStreamStateResponseParams struct {
	// The information of the StreamLive input queried.
	Info *QueryDispatchInputInfo `json:"Info,omitnil,omitempty" name:"Info"`

	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type QueryInputStreamStateResponse struct {
	*tchttp.BaseResponse
	Response *QueryInputStreamStateResponseParams `json:"Response"`
}

func (r *QueryInputStreamStateResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *QueryInputStreamStateResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type RegionInfo struct {
	// Region name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
}

type SDMCSettingsInfo struct {
	// User ID in the SDMC DRM system
	Uid *string `json:"Uid,omitnil,omitempty" name:"Uid"`

	// Tracks of the SDMC DRM system. This parameter is valid for DASH output groups.
	// `1`: audio
	// `2`: SD
	// `4`: HD
	// `8`: UHD1
	// `16`: UHD2
	// 
	// Default value: `31` (audio + SD + HD + UHD1 + UHD2)
	Tracks *int64 `json:"Tracks,omitnil,omitempty" name:"Tracks"`

	// Key ID in the SDMC DRM system; required
	SecretId *string `json:"SecretId,omitnil,omitempty" name:"SecretId"`

	// Key in the SDMC DRM system; required
	SecretKey *string `json:"SecretKey,omitnil,omitempty" name:"SecretKey"`

	// Key request URL of the SDMC DRM system, which is `https://uat.multidrm.tv/cpix/2.2/getcontentkey` by default
	Url *string `json:"Url,omitnil,omitempty" name:"Url"`

	// Token name in an SDMC key request URL, which is `token` by default
	TokenName *string `json:"TokenName,omitnil,omitempty" name:"TokenName"`
}

type Scte35SettingsInfo struct {
	// Whether to pass through SCTE-35 information. Valid values: NO_PASSTHROUGH/PASSTHROUGH. Default value: NO_PASSTHROUGH.
	Behavior *string `json:"Behavior,omitnil,omitempty" name:"Behavior"`
}

type SegmentationDescriptorInfo struct {
	// A 32-bit unique segmentation event identifier. Only one occurrence of a given segmentation_event_id value shall be active at any one time.
	EventID *uint64 `json:"EventID,omitnil,omitempty" name:"EventID"`

	// Indicates that a previously sent segmentation event, identified by segmentation_event_id, has been cancelled.
	EventCancelIndicator *uint64 `json:"EventCancelIndicator,omitnil,omitempty" name:"EventCancelIndicator"`

	// Distribution configuration.
	DeliveryRestrictions *DeliveryRestrictionsInfo `json:"DeliveryRestrictions,omitnil,omitempty" name:"DeliveryRestrictions"`

	// The duration of the segment in 90kHz ticks. indicat when the segment will be over and when the next segmentation message will occur.Shall be 0 for end messages.the time signal will continue until insert a cancellation message when not specify the duration.
	Duration *uint64 `json:"Duration,omitnil,omitempty" name:"Duration"`

	// Corresponds to SCTE-35 segmentation_upid_type parameter.
	UPIDType *uint64 `json:"UPIDType,omitnil,omitempty" name:"UPIDType"`

	// Corresponds to SCTE-35 segmentation_upid. 
	UPID *string `json:"UPID,omitnil,omitempty" name:"UPID"`

	// Corresponds to SCTE-35 segmentation_type_id.
	TypeID *uint64 `json:"TypeID,omitnil,omitempty" name:"TypeID"`

	// Corresponds to SCTE-35 segment_num. This field provides support for numbering segments within a given collection of segments.
	Num *uint64 `json:"Num,omitnil,omitempty" name:"Num"`

	// Corresponds to SCTE-35 segment_expected.This field provides a count of the expected number of individual segments within a collection of segments.
	Expected *uint64 `json:"Expected,omitnil,omitempty" name:"Expected"`

	// Corresponds to SCTE-35 sub_segment_num.This field provides identification for a specific sub-segment within a collection of sub-segments.
	SubSegmentNum *uint64 `json:"SubSegmentNum,omitnil,omitempty" name:"SubSegmentNum"`

	// Corresponds to SCTE-35 sub_segments_expected.This field provides a count of the expected number of individual sub-segments within the collection of sub-segments.
	SubSegmentsExpected *uint64 `json:"SubSegmentsExpected,omitnil,omitempty" name:"SubSegmentsExpected"`
}

type SegmentationDescriptorRespInfo struct {
	// A 32-bit unique segmentation event identifier. Only one occurrence of a given segmentation_event_id value shall be active at any one time.
	EventID *uint64 `json:"EventID,omitnil,omitempty" name:"EventID"`

	// Indicates that a previously sent segmentation event, identified by segmentation_event_id, has been cancelled.
	EventCancelIndicator *uint64 `json:"EventCancelIndicator,omitnil,omitempty" name:"EventCancelIndicator"`

	// Distribution configuration.
	DeliveryRestrictions *DeliveryRestrictionsInfo `json:"DeliveryRestrictions,omitnil,omitempty" name:"DeliveryRestrictions"`

	// The duration of the segment in 90kHz ticks. indicat when the segment will be over and when the next segmentation message will occur.Shall be 0 for end messages.the time signal will continue until insert a cancellation message when not specify the duration.
	Duration *string `json:"Duration,omitnil,omitempty" name:"Duration"`

	// Corresponds to SCTE-35 segmentation_upid_type parameter.
	UPIDType *uint64 `json:"UPIDType,omitnil,omitempty" name:"UPIDType"`

	// Corresponds to SCTE-35 segmentation_upid. 
	UPID *string `json:"UPID,omitnil,omitempty" name:"UPID"`

	// Corresponds to SCTE-35 segmentation_type_id.
	TypeID *uint64 `json:"TypeID,omitnil,omitempty" name:"TypeID"`

	// Corresponds to SCTE-35 segment_num. This field provides support for numbering segments within a given collection of segments.
	Num *uint64 `json:"Num,omitnil,omitempty" name:"Num"`

	// Corresponds to SCTE-35 segment_expected.This field provides a count of the expected number of individual segments within a collection of segments.
	Expected *uint64 `json:"Expected,omitnil,omitempty" name:"Expected"`

	// Corresponds to SCTE-35 sub_segment_num.This field provides identification for a specific sub-segment within a collection of sub-segments.
	SubSegmentNum *uint64 `json:"SubSegmentNum,omitnil,omitempty" name:"SubSegmentNum"`

	// Corresponds to SCTE-35 sub_segments_expected.This field provides a count of the expected number of individual sub-segments within the collection of sub-segments.
	SubSegmentsExpected *uint64 `json:"SubSegmentsExpected,omitnil,omitempty" name:"SubSegmentsExpected"`
}

// Predefined struct for user
type StartStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type StartStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *StartStreamLiveChannelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StartStreamLiveChannelRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "StartStreamLiveChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type StartStreamLiveChannelResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type StartStreamLiveChannelResponse struct {
	*tchttp.BaseResponse
	Response *StartStreamLiveChannelResponseParams `json:"Response"`
}

func (r *StartStreamLiveChannelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StartStreamLiveChannelResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type StaticImageActivateSetting struct {
	// The address of the image to be inserted, starting with http or https and ending with .png .PNG .bmp .BMP .tga .TGA.
	ImageUrl *string `json:"ImageUrl,omitnil,omitempty" name:"ImageUrl"`

	// The layer of the superimposed image, 0-7; the default value is 0, and a higher layer means it is on the top.
	Layer *int64 `json:"Layer,omitnil,omitempty" name:"Layer"`

	// Opacity, range 0-100; the default value is 100, which means completely opaque.
	Opacity *int64 `json:"Opacity,omitnil,omitempty" name:"Opacity"`

	// The distance from the left edge in pixels; the default value is 0 and the maximum value is 4096.
	ImageX *int64 `json:"ImageX,omitnil,omitempty" name:"ImageX"`

	// The distance from the top edge in pixels; the default value is 0 and the maximum value is 2160.
	ImageY *int64 `json:"ImageY,omitnil,omitempty" name:"ImageY"`

	// The width of the image superimposed on the video frame, in pixels. The default value is empty (not set), which means using the original image size. The minimum value is 1 and the maximum value is 4096.
	Width *int64 `json:"Width,omitnil,omitempty" name:"Width"`

	// The height of the image superimposed on the video frame, in pixels. The default value is empty (not set), which means the original image size is used. The minimum value is 1 and the maximum value is 2160.
	Height *int64 `json:"Height,omitnil,omitempty" name:"Height"`

	// Overlay duration, in milliseconds, range 0-86400000; default value 0, 0 means continuous.
	Duration *int64 `json:"Duration,omitnil,omitempty" name:"Duration"`

	// Fade-in duration, in milliseconds, range 0-5000; default value 0, 0 means no fade-in effect.
	FadeIn *int64 `json:"FadeIn,omitnil,omitempty" name:"FadeIn"`

	// Fade-out duration, in milliseconds, range 0-5000; default value 0, 0 means no fade-out effect.
	FadeOut *int64 `json:"FadeOut,omitnil,omitempty" name:"FadeOut"`
}

type StaticImageDeactivateSetting struct {
	// The overlay level to be canceled, range 0-7, default value 0.
	Layer *int64 `json:"Layer,omitnil,omitempty" name:"Layer"`

	// Fade-out duration, in milliseconds, range 0-5000; default value 0, 0 means no fade-out effect.
	FadeOut *int64 `json:"FadeOut,omitnil,omitempty" name:"FadeOut"`
}

type StaticImageSettings struct {
	// Whether to enable global static image overlay, 0: Disable, 1: Enable; Default value: 0.
	GlobalImageOverlayEnabled *int64 `json:"GlobalImageOverlayEnabled,omitnil,omitempty" name:"GlobalImageOverlayEnabled"`
}

// Predefined struct for user
type StopStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type StopStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

func (r *StopStreamLiveChannelRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StopStreamLiveChannelRequest) FromJsonString(s string) error {
	f := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		return err
	}
	delete(f, "Id")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "StopStreamLiveChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type StopStreamLiveChannelResponseParams struct {
	// The unique request ID, generated by the server, will be returned for every request (if the request fails to reach the server for other reasons, the request will not obtain a RequestId). RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type StopStreamLiveChannelResponse struct {
	*tchttp.BaseResponse
	Response *StopStreamLiveChannelResponseParams `json:"Response"`
}

func (r *StopStreamLiveChannelResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

// FromJsonString It is highly **NOT** recommended to use this function
// because it has no param check, nor strict type check
func (r *StopStreamLiveChannelResponse) FromJsonString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

type StreamAudioInfo struct {
	// Audio `Pid`.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pid *int64 `json:"Pid,omitnil,omitempty" name:"Pid"`

	// Audio codec.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Codec *string `json:"Codec,omitnil,omitempty" name:"Codec"`

	// Audio frame rate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Fps *int64 `json:"Fps,omitnil,omitempty" name:"Fps"`

	// Audio bitrate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Rate *int64 `json:"Rate,omitnil,omitempty" name:"Rate"`

	// Audio sample rate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	SampleRate *int64 `json:"SampleRate,omitnil,omitempty" name:"SampleRate"`
}

type StreamInfo struct {
	// Client IP.
	ClientIp *string `json:"ClientIp,omitnil,omitempty" name:"ClientIp"`

	// Video information of pushed streams.
	Video []*StreamVideoInfo `json:"Video,omitnil,omitempty" name:"Video"`

	// Audio information of pushed streams.
	Audio []*StreamAudioInfo `json:"Audio,omitnil,omitempty" name:"Audio"`

	// SCTE-35 information of pushed streams.
	Scte35 []*StreamScte35Info `json:"Scte35,omitnil,omitempty" name:"Scte35"`
}

type StreamLiveChannelInfo struct {
	// Channel ID
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`

	// Channel status
	State *string `json:"State,omitnil,omitempty" name:"State"`

	// Information of attached inputs
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitnil,omitempty" name:"AttachedInputs"`

	// Information of output groups
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitnil,omitempty" name:"OutputGroups"`

	// Channel name
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Audio transcoding templates
	// Note: this field may return `null`, indicating that no valid value was found.
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitnil,omitempty" name:"AudioTemplates"`

	// Video transcoding templates
	// Note: this field may return `null`, indicating that no valid value was found.
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitnil,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates
	// Note: this field may return `null`, indicating that no valid value was found.
	AVTemplates []*AVTemplate `json:"AVTemplates,omitnil,omitempty" name:"AVTemplates"`

	// Caption templates.
	CaptionTemplates []*SubtitleConf `json:"CaptionTemplates,omitnil,omitempty" name:"CaptionTemplates"`

	// Event settings
	// Note: This field may return `null`, indicating that no valid value was found.
	PlanSettings *PlanSettings `json:"PlanSettings,omitnil,omitempty" name:"PlanSettings"`

	// The callback settings.
	// Note: This field may return `null`, indicating that no valid value was found.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitnil,omitempty" name:"EventNotifySettings"`

	// Supplement the last video frame configuration settings.
	InputLossBehavior *InputLossBehaviorInfo `json:"InputLossBehavior,omitnil,omitempty" name:"InputLossBehavior"`

	// Pipeline configuration.
	PipelineInputSettings *PipelineInputSettingsInfo `json:"PipelineInputSettings,omitnil,omitempty" name:"PipelineInputSettings"`

	// Recognition configuration for input content.
	InputAnalysisSettings *InputAnalysisInfo `json:"InputAnalysisSettings,omitnil,omitempty" name:"InputAnalysisSettings"`

	// Console tag list.
	Tags []*Tag `json:"Tags,omitnil,omitempty" name:"Tags"`

	// Frame capture templates.
	FrameCaptureTemplates []*FrameCaptureTemplate `json:"FrameCaptureTemplates,omitnil,omitempty" name:"FrameCaptureTemplates"`

	// General settings.
	GeneralSettings *GeneralSetting `json:"GeneralSettings,omitnil,omitempty" name:"GeneralSettings"`
}

type StreamLiveOutputGroupsInfo struct {
	// Output group name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the channel level
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Output protocol
	// Valid values: `HLS`, `DASH`, `HLS_ARCHIVE`, 
	//  `DASH_ARCHIVE`, `HLS_STREAM_PACKAGE`, 
	//  `DASH_STREAM_PACKAGE`, 
	//  `FRAME_CAPTURE`, `RTP`, `RTMP`, `M2TS`.
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Output information
	// If the type is RTMP, RTP or FRAME_CAPTURE, only one output is allowed; if it is HLS or DASH, 1-10 outputs are allowed.
	Outputs []*OutputInfo `json:"Outputs,omitnil,omitempty" name:"Outputs"`

	// Relay destinations. Quantity: [1, 2]
	Destinations []*DestinationInfo `json:"Destinations,omitnil,omitempty" name:"Destinations"`

	// HLS protocol configuration information, which takes effect only for HLS/HLS_ARCHIVE/HLS_STREAM_PACKAGE outputs.
	// Note: this field may return `null`, indicating that no valid value was found.
	HlsRemuxSettings *HlsRemuxSettingsInfo `json:"HlsRemuxSettings,omitnil,omitempty" name:"HlsRemuxSettings"`

	// DRM configuration information
	// Note: this field may return `null`, indicating that no valid value was found.
	DrmSettings *DrmSettingsInfo `json:"DrmSettings,omitnil,omitempty" name:"DrmSettings"`

	// DASH protocol configuration information, which takes effect only for DASH/DASH_ARCHIVE outputs
	// Note: this field may return `null`, indicating that no valid value was found.
	DashRemuxSettings *DashRemuxSettingsInfo `json:"DashRemuxSettings,omitnil,omitempty" name:"DashRemuxSettings"`

	// StreamPackage configuration information, which is required if the output type is StreamPackage
	// Note: this field may return `null`, indicating that no valid value was found.
	StreamPackageSettings *StreamPackageSettingsInfo `json:"StreamPackageSettings,omitnil,omitempty" name:"StreamPackageSettings"`

	// Time-shift configuration information
	// Note: This field may return `null`, indicating that no valid value was found.
	TimeShiftSettings *TimeShiftSettingsInfo `json:"TimeShiftSettings,omitnil,omitempty" name:"TimeShiftSettings"`
}

type StreamLiveRegionInfo struct {
	// List of StreamLive regions
	Regions []*RegionInfo `json:"Regions,omitnil,omitempty" name:"Regions"`
}

type StreamPackageSettingsInfo struct {
	// Channel ID in StreamPackage
	Id *string `json:"Id,omitnil,omitempty" name:"Id"`
}

type StreamScte35Info struct {
	// SCTE-35 `Pid`.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pid *int64 `json:"Pid,omitnil,omitempty" name:"Pid"`
}

type StreamVideoInfo struct {
	// Video `Pid`.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pid *int64 `json:"Pid,omitnil,omitempty" name:"Pid"`

	// Video codec.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Codec *string `json:"Codec,omitnil,omitempty" name:"Codec"`

	// Video frame rate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Fps *int64 `json:"Fps,omitnil,omitempty" name:"Fps"`

	// Video bitrate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Rate *int64 `json:"Rate,omitnil,omitempty" name:"Rate"`

	// Video width.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Width *int64 `json:"Width,omitnil,omitempty" name:"Width"`

	// Video height.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Height *int64 `json:"Height,omitnil,omitempty" name:"Height"`
}

type SubtitleConf struct {
	// Template name.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Name of caption selector. Required when CaptionSource selects `INPUT`.
	CaptionSelectorName *string `json:"CaptionSelectorName,omitnil,omitempty" name:"CaptionSelectorName"`

	// Optional values: INPUT (source subtitle information), ANALYSIS (intelligent speech recognition to subtitles).
	CaptionSource *string `json:"CaptionSource,omitnil,omitempty" name:"CaptionSource"`

	// Optional values: 1 Source, 2 Source+Target, 3 Target (original language only, original language + translation language, translation language). Required when CaptionSource selects `ANALYSIS `. When outputting as WebVTT, a single template can only output one language.
	ContentType *uint64 `json:"ContentType,omitnil,omitempty" name:"ContentType"`

	// Output mode: 1 Burn in, 2 Embedded, 3 WebVTT. Support `2` when CaptionSource selects `INPUT`. Support `1` and `3` when CaptionSource selects `ANALYSIS `.
	TargetType *uint64 `json:"TargetType,omitnil,omitempty" name:"TargetType"`

	// Original phonetic language.
	// Optional values: Chinese, English, Japanese, Korean. Required when CaptionSource selects `ANALYSIS `.
	SourceLanguage *string `json:"SourceLanguage,omitnil,omitempty" name:"SourceLanguage"`

	// Target language.
	// Optional values: Chinese, English, Japanese, Korean. Required when CaptionSource selects `ANALYSIS `.
	TargetLanguage *string `json:"TargetLanguage,omitnil,omitempty" name:"TargetLanguage"`

	// Font style configuration. Required when CaptionSource selects `ANALYSIS `.
	FontStyle *SubtitleFontConf `json:"FontStyle,omitnil,omitempty" name:"FontStyle"`

	// There are two modes: STEADY and DYNAMIC, corresponding to steady state and unstable state respectively; the default is STEADY. Required when CaptionSource selects `ANALYSIS `. When the output is WebVTT, only STEADY can be selected.
	StateEffectMode *string `json:"StateEffectMode,omitnil,omitempty" name:"StateEffectMode"`

	// Steady-state delay time, unit seconds; optional values: 10, 20, default 10. Required when CaptionSource selects `ANALYSIS `.
	SteadyStateDelayedTime *uint64 `json:"SteadyStateDelayedTime,omitnil,omitempty" name:"SteadyStateDelayedTime"`

	// Audio selector name, required for generating WebVTT subtitles using speech recognition, can be empty.
	AudioSelectorName *string `json:"AudioSelectorName,omitnil,omitempty" name:"AudioSelectorName"`

	// Format configuration for speech recognition output on WebVTT.
	WebVTTFontStyle *WebVTTFontStyle `json:"WebVTTFontStyle,omitnil,omitempty" name:"WebVTTFontStyle"`

	// Language code, length 2-20. ISO 639-2 three-digit code is recommend.
	LanguageCode *string `json:"LanguageCode,omitnil,omitempty" name:"LanguageCode"`

	// Language description, less than 100 characters in length.
	LanguageDescription *string `json:"LanguageDescription,omitnil,omitempty" name:"LanguageDescription"`
}

type SubtitleFontConf struct {
	// Line spacing.
	LineSpacing *uint64 `json:"LineSpacing,omitnil,omitempty" name:"LineSpacing"`

	// Margins.
	Margins *uint64 `json:"Margins,omitnil,omitempty" name:"Margins"`

	// Rows.
	Lines *uint64 `json:"Lines,omitnil,omitempty" name:"Lines"`

	// Number of characters per line.
	CharactersPerLine *uint64 `json:"CharactersPerLine,omitnil,omitempty" name:"CharactersPerLine"`

	// Original font Helvetica: simhei.ttf Song Dynasty: simsun.ttc Dynacw Diamond Black: hkjgh.ttf Helvetica font: helvetica.ttf; Need to be set in Source or Source+Target mode
	SourceTextFont *string `json:"SourceTextFont,omitnil,omitempty" name:"SourceTextFont"`

	// Font color is represented by 6 RGB hexadecimal characters.
	TextColor *string `json:"TextColor,omitnil,omitempty" name:"TextColor"`

	// The background color is represented by 6 RGB hexadecimal characters.
	BackgroundColor *string `json:"BackgroundColor,omitnil,omitempty" name:"BackgroundColor"`

	// Background transparency, a number from 0-100.
	BackgroundAlpha *uint64 `json:"BackgroundAlpha,omitnil,omitempty" name:"BackgroundAlpha"`

	// Preview copy.
	PreviewContent *string `json:"PreviewContent,omitnil,omitempty" name:"PreviewContent"`

	// Preview window height.
	PreviewWindowHeight *uint64 `json:"PreviewWindowHeight,omitnil,omitempty" name:"PreviewWindowHeight"`

	// Preview window width.
	PreviewWindowWidth *uint64 `json:"PreviewWindowWidth,omitnil,omitempty" name:"PreviewWindowWidth"`

	// Translation language font, the enumeration value is the same as Font, the fonts supported by the language need to be distinguished; TextColor needs to be set in Target or Source+Target mode
	TranslatedTextFont *string `json:"TranslatedTextFont,omitnil,omitempty" name:"TranslatedTextFont"`
}

type Tag struct {
	// Tag key, for restrictions please refer to the tag documentation: https://www.tencentcloud.com/document/product/651/13354.
	TagKey *string `json:"TagKey,omitnil,omitempty" name:"TagKey"`

	// Tag value, for restrictions please refer to the tag documentation: https://www.tencentcloud.com/document/product/651/13354.
	TagValue *string `json:"TagValue,omitnil,omitempty" name:"TagValue"`

	// Tag type, optional; for documentation please refer to: https://www.tencentcloud.com/document/product/651/33023#tag.
	Category *string `json:"Category,omitnil,omitempty" name:"Category"`
}

type TaskNotifyConfig struct {
	// Notification type. Currently only supports URLs
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Callback URL. Starting with http/https, supporting a maximum of 512 characters
	NotifyUrl *string `json:"NotifyUrl,omitnil,omitempty" name:"NotifyUrl"`
}

type ThumbnailSettings struct {
	// Generate thumbnail ,0: Disabled ,1: Enabled , Default: 0
	ThumbnailEnabled *int64 `json:"ThumbnailEnabled,omitnil,omitempty" name:"ThumbnailEnabled"`
}

type TimeShiftSettingsInfo struct {
	// Whether to enable time shifting. Valid values: `OPEN`; `CLOSE`
	// Note: This field may return `null`, indicating that no valid value was found.
	State *string `json:"State,omitnil,omitempty" name:"State"`

	// Domain name bound for time shifting
	// Note: This field may return `null`, indicating that no valid value was found.
	PlayDomain *string `json:"PlayDomain,omitnil,omitempty" name:"PlayDomain"`

	// Allowable time-shift period (s). Value range: [300, 2592000]. Default value: 300Note: This field may return `null`, indicating that no valid value was found.
	StartoverWindow *int64 `json:"StartoverWindow,omitnil,omitempty" name:"StartoverWindow"`
}

type TimedMetadataInfo struct {
	// Base64-encoded id3 metadata information, with a maximum limit of 1024 characters.
	ID3 *string `json:"ID3,omitnil,omitempty" name:"ID3"`
}

type TimedMetadataSettingInfo struct {
	// Whether to transparently transmit ID3 information, optional values: 0:NO_PASSTHROUGH, 1:PASSTHROUGH, default 0.
	Behavior *uint64 `json:"Behavior,omitnil,omitempty" name:"Behavior"`
}

type TimedRecordSettings struct {
	// Whether to automatically delete finished recording events. Valid values: `CLOSE`, `OPEN`. If this parameter is left empty, `CLOSE` will be used.
	// If it is set to `OPEN`, a recording event will be deleted 7 days after it is finished.
	// Note: This field may return `null`, indicating that no valid value was found.
	AutoClear *string `json:"AutoClear,omitnil,omitempty" name:"AutoClear"`
}

type TimingSettingsReq struct {
	// Event trigger type. Valid values: `FIXED_TIME`, `IMMEDIATE`,`FIXED_PTS `. This parameter is required if `EventType` is `INPUT_SWITCH`.
	StartType *string `json:"StartType,omitnil,omitempty" name:"StartType"`

	// This parameter is required if `EventType` is `INPUT_SWITCH` and `StartType` is `FIXED_TIME`.
	// It must be in UTC format, e.g., `2020-01-01T12:00:00Z`.
	Time *string `json:"Time,omitnil,omitempty" name:"Time"`

	// This parameter is required if `EventType` is `TIMED_RECORD`.
	// It specifies the recording start time in UTC format (e.g., `2020-01-01T12:00:00Z`) and must be at least 1 minute later than the current time.
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// This parameter is required if `EventType` is `TIMED_RECORD`.
	// It specifies the recording end time in UTC format (e.g., `2020-01-01T12:00:00Z`) and must be at least 1 minute later than the recording start time.
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Effective only when StartType is FIXED_PTS, with a range of 1-8589934592
	PTS *uint64 `json:"PTS,omitnil,omitempty" name:"PTS"`
}

type TimingSettingsResp struct {
	// Event trigger type
	StartType *string `json:"StartType,omitnil,omitempty" name:"StartType"`

	// Not empty if `StartType` is `FIXED_TIME`
	// UTC time, such as `2020-01-01T12:00:00Z`
	Time *string `json:"Time,omitnil,omitempty" name:"Time"`

	// This parameter cannot be empty if `EventType` is `TIMED_RECORD`.
	// It indicates the start time for recording in UTC format (e.g., `2020-01-01T12:00:00Z`) and must be at least 1 minute later than the current time.
	StartTime *string `json:"StartTime,omitnil,omitempty" name:"StartTime"`

	// This parameter cannot be empty if `EventType` is `TIMED_RECORD`.
	// It indicates the end time for recording in UTC format (e.g., `2020-01-01T12:00:00Z`) and must be at least 1 minute later than the start time for recording.
	EndTime *string `json:"EndTime,omitnil,omitempty" name:"EndTime"`

	// Effective only when StartType is FIXED_PTS, with a range of 1-8589934592
	PTS *uint64 `json:"PTS,omitnil,omitempty" name:"PTS"`
}

type UrlInputInfo struct {
	// Video URL, starting with http/https, supports up to 512 characters, currently only supports complete single file videos, does not support streaming formats based on playlists and segments (such as HLS or DASH)
	Url *string `json:"Url,omitnil,omitempty" name:"Url"`
}

type VideoCodecDetail struct {
	// The three image quality levels of h264 include: BASELINE, HIGH, and MAIN. The default option is MAIN.
	Profile *string `json:"Profile,omitnil,omitempty" name:"Profile"`

	// Profile corresponding codec performance, options include: 1, 1.1, 1.2, 1.3, 2, 2.1, 2.2, 2.3, 3, 3.1, 3.2, 4, 4.1, 4.2, 5, 5.1, AUTO. The default option is AUTO.
	Level *string `json:"Level,omitnil,omitempty" name:"Level"`

	// Codecs include entropy coding and lossless coding, and options include: CABAC and CAVLC. The default option is CABAC. .
	EntropyEncoding *string `json:"EntropyEncoding,omitnil,omitempty" name:"EntropyEncoding"`

	// Mode, options include: AUTO, HIGH, HIGHER, LOW, MAX, MEDIUM, OFF. The default option is: AUTO. .
	AdaptiveQuantization *string `json:"AdaptiveQuantization,omitnil,omitempty" name:"AdaptiveQuantization"`

	// Analyze subsequent encoded frames in advance, options include: HIGH, LOW, MEDIUM. The default option is: MEDIUM. .
	LookAheadRateControl *string `json:"LookAheadRateControl,omitnil,omitempty" name:"LookAheadRateControl"`
}

type VideoEnhanceSetting struct {
	// Video enhancement types, optional: "GameEnhance", "ColorEnhance", "Debur", "Comprehensive", "Denoising", "SR", "OutdoorSportsCompetitions", "IndoorSportsCompetitions", "ShowEnhance"
	Type *string `json:"Type,omitnil,omitempty" name:"Type"`

	// Video enhancement intensity, 0-1.0, granularity 0.1
	Strength *float64 `json:"Strength,omitnil,omitempty" name:"Strength"`
}

type VideoPipelineInputStatistics struct {
	// Video FPS.
	Fps *uint64 `json:"Fps,omitnil,omitempty" name:"Fps"`

	// Video bitrate in bps.
	Rate *uint64 `json:"Rate,omitnil,omitempty" name:"Rate"`

	// Video `Pid`, which is available only if the input is `rtp/udp`.
	Pid *int64 `json:"Pid,omitnil,omitempty" name:"Pid"`
}

type VideoTemplateInfo struct {
	// Video transcoding template name, which can contain 1-20 letters and digits.
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`

	// Video codec. Valid values: H264/H265. If this parameter is left empty, the original value will be used.
	Vcodec *string `json:"Vcodec,omitnil,omitempty" name:"Vcodec"`

	// Video bitrate. Value range: [50000,40000000]. The value can only be a multiple of 1,000. If this parameter is left empty, the original value will be used.
	VideoBitrate *uint64 `json:"VideoBitrate,omitnil,omitempty" name:"VideoBitrate"`

	// Video width. Value range: (0,4096]. The value can only be a multiple of 2. If this parameter is left empty, the original value will be used.
	Width *uint64 `json:"Width,omitnil,omitempty" name:"Width"`

	// Video height. Value range: (0,4096]. The value can only be a multiple of 2. If this parameter is left empty, the original value will be used.
	Height *uint64 `json:"Height,omitnil,omitempty" name:"Height"`

	// Video frame rate. Value range: [1,240]. If this parameter is left empty, the original value will be used.
	Fps *uint64 `json:"Fps,omitnil,omitempty" name:"Fps"`

	// Whether to enable top speed codec. Valid value: CLOSE/OPEN. Default value: CLOSE.
	TopSpeed *string `json:"TopSpeed,omitnil,omitempty" name:"TopSpeed"`

	// Top speed codec compression ratio. Value range: [0,50]. The lower the compression ratio, the higher the image quality.
	BitrateCompressionRatio *uint64 `json:"BitrateCompressionRatio,omitnil,omitempty" name:"BitrateCompressionRatio"`

	// Bitrate control mode. Valid values: `CBR`, `ABR` (default), `VBR`.
	RateControlMode *string `json:"RateControlMode,omitnil,omitempty" name:"RateControlMode"`

	// Watermark ID
	// Note: This field may return `null`, indicating that no valid value was found.
	WatermarkId *string `json:"WatermarkId,omitnil,omitempty" name:"WatermarkId"`

	// Whether to enable the face blur function, 1 is on, 0 is off, and the default is 0.
	FaceBlurringEnabled *uint64 `json:"FaceBlurringEnabled,omitnil,omitempty" name:"FaceBlurringEnabled"`

	// This field indicates how to specify the output video frame rate. If FOLLOW_SOURCE is selected, the output video frame rate will be set equal to the input video frame rate of the first input. If SPECIFIED_FRACTION is selected, the output video frame rate is determined by the fraction (frame rate numerator and frame rate denominator). If SPECIFIED_HZ is selected, the frame rate of the output video is determined by the HZ you enter.
	FrameRateType *string `json:"FrameRateType,omitnil,omitempty" name:"FrameRateType"`

	// Valid when the FrameRateType type you select is SPECIFIED_FRACTION, the output frame rate numerator setting.
	FrameRateNumerator *uint64 `json:"FrameRateNumerator,omitnil,omitempty" name:"FrameRateNumerator"`

	// Valid when the FrameRateType type you select is SPECIFIED_FRACTION, the output frame rate denominator setting.
	FrameRateDenominator *uint64 `json:"FrameRateDenominator,omitnil,omitempty" name:"FrameRateDenominator"`

	// The number of B frames can be selected from 1 to 3.
	BFramesNum *uint64 `json:"BFramesNum,omitnil,omitempty" name:"BFramesNum"`

	// The number of reference frames can be selected from 1 to 16.
	RefFramesNum *uint64 `json:"RefFramesNum,omitnil,omitempty" name:"RefFramesNum"`

	// Additional video bitrate configuration.
	AdditionalRateSettings *AdditionalRateSetting `json:"AdditionalRateSettings,omitnil,omitempty" name:"AdditionalRateSettings"`

	// Video encoding configuration.
	VideoCodecDetails *VideoCodecDetail `json:"VideoCodecDetails,omitnil,omitempty" name:"VideoCodecDetails"`

	// Video enhancement switch, 1: on 0: off.
	VideoEnhanceEnabled *uint64 `json:"VideoEnhanceEnabled,omitnil,omitempty" name:"VideoEnhanceEnabled"`

	// Video enhancement parameter array.
	VideoEnhanceSettings []*VideoEnhanceSetting `json:"VideoEnhanceSettings,omitnil,omitempty" name:"VideoEnhanceSettings"`

	// Color space setting.
	ColorSpaceSettings *ColorSpaceSetting `json:"ColorSpaceSettings,omitnil,omitempty" name:"ColorSpaceSettings"`

	// Traceability watermark.
	ForensicWatermarkIds []*string `json:"ForensicWatermarkIds,omitnil,omitempty" name:"ForensicWatermarkIds"`
}

type WebVTTFontStyle struct {
	// Text color, RGB hexadecimal representation, 6 hexadecimal characters (no # needed).
	TextColor *string `json:"TextColor,omitnil,omitempty" name:"TextColor"`

	// Background color, RGB hexadecimal representation, 6 hexadecimal characters (no # needed).
	BackgroundColor *string `json:"BackgroundColor,omitnil,omitempty" name:"BackgroundColor"`

	// Background opacity parameter, a number from 0 to 100, with 0 being the default for full transparency.
	BackgroundAlpha *int64 `json:"BackgroundAlpha,omitnil,omitempty" name:"BackgroundAlpha"`

	// Font size, in units of vh (1% of height), default value 0 means automatic.
	FontSize *int64 `json:"FontSize,omitnil,omitempty" name:"FontSize"`

	// The position of the text box, default value AUTO, can be empty; represents the percentage of video height, supports integers from 0 to 100.
	Line *string `json:"Line,omitnil,omitempty" name:"Line"`

	// The alignment of the text box on the Line. Optional values: START, CENTER, END. Which can be empty.
	LineAlignment *string `json:"LineAlignment,omitnil,omitempty" name:"LineAlignment"`

	// The text box is positioned in another direction as a percentage of the video's width. It defaults to AUTO and can be empty.
	Position *string `json:"Position,omitnil,omitempty" name:"Position"`

	// The alignment of the text box on the Position. Optional values are LINE_LEFT, LINE_RIGHT, CENTER, and AUTO. The default value is AUTO, and it can be empty.
	PositionAlignment *string `json:"PositionAlignment,omitnil,omitempty" name:"PositionAlignment"`

	// Text box size, a percentage of video width/height, with values (0, 100), default AUTO, can be empty.
	CueSize *string `json:"CueSize,omitnil,omitempty" name:"CueSize"`

	// Text alignment, with possible values  START, CENTER, END, LEFT, and RIGHT; the default value is CENTER, which can be empty.
	TextAlignment *string `json:"TextAlignment,omitnil,omitempty" name:"TextAlignment"`
}