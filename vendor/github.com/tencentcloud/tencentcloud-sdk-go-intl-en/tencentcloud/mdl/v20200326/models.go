// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
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
    "encoding/json"
    tcerr "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/errors"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/http"
)

type AVTemplate struct {
	// Name of an audio/video transcoding template, which can contain 1-20 case-sensitive letters and digits
	Name *string `json:"Name,omitempty" name:"Name"`

	// Whether video is needed. `0`: not needed; `1`: needed
	NeedVideo *uint64 `json:"NeedVideo,omitempty" name:"NeedVideo"`

	// Video codec. Valid values: `H264`, `H265`. If this parameter is left empty, the original video codec will be used.
	Vcodec *string `json:"Vcodec,omitempty" name:"Vcodec"`

	// Video width. Value range: (0, 3000]. The value must be an integer multiple of 4. If this parameter is left empty, the original video width will be used.
	Width *uint64 `json:"Width,omitempty" name:"Width"`

	// Video height. Value range: (0, 3000]. The value must be an integer multiple of 4. If this parameter is left empty, the original video height will be used.
	Height *uint64 `json:"Height,omitempty" name:"Height"`

	// Video frame rate. Value range: [1, 240]. If this parameter is left empty, the original frame rate will be used.
	Fps *uint64 `json:"Fps,omitempty" name:"Fps"`

	// Whether to enable top speed codec transcoding. Valid values: `CLOSE` (disable), `OPEN` (enable). Default value: `CLOSE`
	TopSpeed *string `json:"TopSpeed,omitempty" name:"TopSpeed"`

	// Compression ratio for top speed codec transcoding. Value range: [0, 50]. The lower the compression ratio, the higher the image quality.
	BitrateCompressionRatio *uint64 `json:"BitrateCompressionRatio,omitempty" name:"BitrateCompressionRatio"`

	// Whether audio is needed. `0`: not needed; `1`: needed
	NeedAudio *int64 `json:"NeedAudio,omitempty" name:"NeedAudio"`

	// Audio codec. Valid value: `AAC` (default)
	Acodec *string `json:"Acodec,omitempty" name:"Acodec"`

	// Audio bitrate. If this parameter is left empty, the original bitrate will be used.
	// Valid values: `6000`, `7000`, `8000`, `10000`, `12000`, `14000`, `16000`, `20000`, `24000`, `28000`, `32000`, `40000`, `48000`, `56000`, `64000`, `80000`, `96000`, `112000`, `128000`, `160000`, `192000`, `224000`, `256000`, `288000`, `320000`, `384000`, `448000`, `512000`, `576000`, `640000`, `768000`, `896000`, `1024000`
	AudioBitrate *uint64 `json:"AudioBitrate,omitempty" name:"AudioBitrate"`

	// Video bitrate. Value range: [50000, 40000000]. The value must be an integer multiple of 1000. If this parameter is left empty, the original bitrate will be used.
	VideoBitrate *uint64 `json:"VideoBitrate,omitempty" name:"VideoBitrate"`

	// Bitrate control mode. Valid values: `CBR`, `ABR` (default)
	RateControlMode *string `json:"RateControlMode,omitempty" name:"RateControlMode"`

	// Watermark ID
	WatermarkId *string `json:"WatermarkId,omitempty" name:"WatermarkId"`

	// Whether to convert audio to text. `0` (default): No; `1`: Yes.
	SmartSubtitles *uint64 `json:"SmartSubtitles,omitempty" name:"SmartSubtitles"`

	// The subtitle settings. Currently, the following subtitles are supported:
	// `eng2eng`: English speech to English text.
	// `eng2chs`: English speech to Chinese text. 
	// `eng2chseng`: English speech to English and Chinese text. 
	// `chs2chs`: Chinese speech to Chinese text.   
	// `chs2eng`: Chinese speech to English text. 
	// `chs2chseng`: Chinese speech to Chinese and English text.
	SubtitleConfiguration *string `json:"SubtitleConfiguration,omitempty" name:"SubtitleConfiguration"`
}

type AttachedInput struct {
	// Input ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// Audio selector for the input. There can be 0 to 20 audio selectors.
	// Note: this field may return `null`, indicating that no valid value was found.
	AudioSelectors []*AudioSelectorInfo `json:"AudioSelectors,omitempty" name:"AudioSelectors"`

	// Pull mode. If the input type is `HLS_PULL` or `MP4_PULL`, you can set this parameter to `LOOP` or `ONCE`. `LOOP` is the default value.
	// Note: this field may return `null`, indicating that no valid value was found.
	PullBehavior *string `json:"PullBehavior,omitempty" name:"PullBehavior"`

	// Input failover configuration
	// Note: this field may return `null`, indicating that no valid value was found.
	FailOverSettings *FailOverSettings `json:"FailOverSettings,omitempty" name:"FailOverSettings"`
}

type AudioPidSelectionInfo struct {
	// Audio `Pid`. Default value: 0.
	Pid *uint64 `json:"Pid,omitempty" name:"Pid"`
}

type AudioPipelineInputStatistics struct {
	// Audio FPS.
	Fps *uint64 `json:"Fps,omitempty" name:"Fps"`

	// Audio bitrate in bps.
	Rate *uint64 `json:"Rate,omitempty" name:"Rate"`

	// Audio `Pid`, which is available only if the input is `rtp/udp`.
	Pid *int64 `json:"Pid,omitempty" name:"Pid"`
}

type AudioSelectorInfo struct {
	// Audio name, which can contain 1-32 letters, digits, and underscores.
	Name *string `json:"Name,omitempty" name:"Name"`

	// Audio `Pid` selection.
	AudioPidSelection *AudioPidSelectionInfo `json:"AudioPidSelection,omitempty" name:"AudioPidSelection"`
}

type AudioTemplateInfo struct {
	// Only `AttachedInputs.AudioSelectors.Name` can be selected. This parameter is required for RTP_PUSH and UDP_PUSH.
	AudioSelectorName *string `json:"AudioSelectorName,omitempty" name:"AudioSelectorName"`

	// Audio transcoding template name, which can contain 1-20 letters and digits.
	Name *string `json:"Name,omitempty" name:"Name"`

	// Audio codec. Valid value: AAC. Default value: AAC.
	Acodec *string `json:"Acodec,omitempty" name:"Acodec"`

	// Audio bitrate. If this parameter is left empty, the original value will be used.
	// Valid values: 6000, 7000, 8000, 10000, 12000, 14000, 16000, 20000, 24000, 28000, 32000, 40000, 48000, 56000, 64000, 80000, 96000, 112000, 128000, 160000, 192000, 224000, 256000, 288000, 320000, 384000, 448000, 512000, 576000, 640000, 768000, 896000, 1024000
	AudioBitrate *uint64 `json:"AudioBitrate,omitempty" name:"AudioBitrate"`

	// Audio language code, whose length is always 3 characters.
	LanguageCode *string `json:"LanguageCode,omitempty" name:"LanguageCode"`
}

type ChannelAlertInfos struct {
	// Alarm details of pipeline 0 under this channel.
	Pipeline0 []*ChannelPipelineAlerts `json:"Pipeline0,omitempty" name:"Pipeline0"`

	// Alarm details of pipeline 1 under this channel.
	Pipeline1 []*ChannelPipelineAlerts `json:"Pipeline1,omitempty" name:"Pipeline1"`
}

type ChannelInputStatistics struct {
	// Input ID.
	InputId *string `json:"InputId,omitempty" name:"InputId"`

	// Input statistics.
	Statistics *InputStatistics `json:"Statistics,omitempty" name:"Statistics"`
}

type ChannelOutputsStatistics struct {
	// Output group name.
	OutputGroupName *string `json:"OutputGroupName,omitempty" name:"OutputGroupName"`

	// Output group statistics.
	Statistics *OutputsStatistics `json:"Statistics,omitempty" name:"Statistics"`
}

type ChannelPipelineAlerts struct {
	// Alarm start time in UTC time.
	SetTime *string `json:"SetTime,omitempty" name:"SetTime"`

	// Alarm end time in UTC time.
	// This time is available only after the alarm ends.
	ClearTime *string `json:"ClearTime,omitempty" name:"ClearTime"`

	// Alarm type.
	Type *string `json:"Type,omitempty" name:"Type"`

	// Alarm details.
	Message *string `json:"Message,omitempty" name:"Message"`
}

type CreateImageSettings struct {
	// Image file format. Valid values: png, jpg.
	ImageType *string `json:"ImageType,omitempty" name:"ImageType"`

	// Base64 encoded image content
	ImageContent *string `json:"ImageContent,omitempty" name:"ImageContent"`

	// Origin. Valid values: TOP_LEFT, BOTTOM_LEFT, TOP_RIGHT, BOTTOM_RIGHT.
	Location *string `json:"Location,omitempty" name:"Location"`

	// The watermark’s horizontal distance from the origin as a percentage of the video width. Value range: 0-100. Default: 10.
	XPos *int64 `json:"XPos,omitempty" name:"XPos"`

	// The watermark’s vertical distance from the origin as a percentage of the video height. Value range: 0-100. Default: 10.
	YPos *int64 `json:"YPos,omitempty" name:"YPos"`

	// The watermark image’s width as a percentage of the video width. Value range: 0-100. Default: 10.
	// `0` means to scale the width proportionally to the height.
	// You cannot set both `Width` and `Height` to `0`.
	Width *int64 `json:"Width,omitempty" name:"Width"`

	// The watermark image’s height as a percentage of the video height. Value range: 0-100. Default: 10.
	// `0` means to scale the height proportionally to the width.
	// You cannot set both `Width` and `Height` to `0`.
	Height *int64 `json:"Height,omitempty" name:"Height"`
}

// Predefined struct for user
type CreateStreamLiveChannelRequestParams struct {
	// Channel name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Inputs to attach. You can attach 1 to 5 inputs.
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitempty" name:"AttachedInputs"`

	// Configuration information of the channel’s output groups. Quantity: [1, 10]
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitempty" name:"OutputGroups"`

	// Audio transcoding templates. Quantity: [1, 20]
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitempty" name:"AudioTemplates"`

	// Video transcoding templates. Quantity: [1, 10]
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates. Quantity: [1, 10]
	AVTemplates []*AVTemplate `json:"AVTemplates,omitempty" name:"AVTemplates"`

	// Event settings
	PlanSettings *PlanSettings `json:"PlanSettings,omitempty" name:"PlanSettings"`

	// The callback settings.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitempty" name:"EventNotifySettings"`
}

type CreateStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Inputs to attach. You can attach 1 to 5 inputs.
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitempty" name:"AttachedInputs"`

	// Configuration information of the channel’s output groups. Quantity: [1, 10]
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitempty" name:"OutputGroups"`

	// Audio transcoding templates. Quantity: [1, 20]
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitempty" name:"AudioTemplates"`

	// Video transcoding templates. Quantity: [1, 10]
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates. Quantity: [1, 10]
	AVTemplates []*AVTemplate `json:"AVTemplates,omitempty" name:"AVTemplates"`

	// Event settings
	PlanSettings *PlanSettings `json:"PlanSettings,omitempty" name:"PlanSettings"`

	// The callback settings.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitempty" name:"EventNotifySettings"`
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
	delete(f, "PlanSettings")
	delete(f, "EventNotifySettings")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateStreamLiveChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveChannelResponseParams struct {
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Name *string `json:"Name,omitempty" name:"Name"`

	// Input type
	// Valid values: `RTMP_PUSH`, `RTP_PUSH`, `UDP_PUSH`, `RTMP_PULL`, `HLS_PULL`, `MP4_PULL`
	Type *string `json:"Type,omitempty" name:"Type"`

	// ID of the input security group to attach
	// You can attach only one security group to an input.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitempty" name:"SecurityGroupIds"`

	// Input settings. For the type `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitempty" name:"InputSettings"`
}

type CreateStreamLiveInputRequest struct {
	*tchttp.BaseRequest
	
	// Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Input type
	// Valid values: `RTMP_PUSH`, `RTP_PUSH`, `UDP_PUSH`, `RTMP_PULL`, `HLS_PULL`, `MP4_PULL`
	Type *string `json:"Type,omitempty" name:"Type"`

	// ID of the input security group to attach
	// You can attach only one security group to an input.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitempty" name:"SecurityGroupIds"`

	// Input settings. For the type `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitempty" name:"InputSettings"`
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
	Id *string `json:"Id,omitempty" name:"Id"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Name *string `json:"Name,omitempty" name:"Name"`

	// Allowlist entries. Quantity: [1, 10]
	Whitelist []*string `json:"Whitelist,omitempty" name:"Whitelist"`
}

type CreateStreamLiveInputSecurityGroupRequest struct {
	*tchttp.BaseRequest
	
	// Input security group name, which can contain case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Allowlist entries. Quantity: [1, 10]
	Whitelist []*string `json:"Whitelist,omitempty" name:"Whitelist"`
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
	Id *string `json:"Id,omitempty" name:"Id"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Event configuration
	Plan *PlanReq `json:"Plan,omitempty" name:"Plan"`
}

type CreateStreamLivePlanRequest struct {
	*tchttp.BaseRequest
	
	// ID of the channel for which you want to configure an event
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Event configuration
	Plan *PlanReq `json:"Plan,omitempty" name:"Plan"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Name *string `json:"Name,omitempty" name:"Name"`

	// Watermark type. Valid values: STATIC_IMAGE, TEXT.
	Type *string `json:"Type,omitempty" name:"Type"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	ImageSettings *CreateImageSettings `json:"ImageSettings,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	TextSettings *CreateTextSettings `json:"TextSettings,omitempty" name:"TextSettings"`
}

type CreateStreamLiveWatermarkRequest struct {
	*tchttp.BaseRequest
	
	// Watermark name
	Name *string `json:"Name,omitempty" name:"Name"`

	// Watermark type. Valid values: STATIC_IMAGE, TEXT.
	Type *string `json:"Type,omitempty" name:"Type"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	ImageSettings *CreateImageSettings `json:"ImageSettings,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	TextSettings *CreateTextSettings `json:"TextSettings,omitempty" name:"TextSettings"`
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
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "CreateStreamLiveWatermarkRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type CreateStreamLiveWatermarkResponseParams struct {
	// Watermark ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Text *string `json:"Text,omitempty" name:"Text"`

	// Origin. Valid values: TOP_LEFT, BOTTOM_LEFT, TOP_RIGHT, BOTTOM_RIGHT.
	Location *string `json:"Location,omitempty" name:"Location"`

	// The watermark’s horizontal distance from the origin as a percentage of the video width. Value range: 0-100. Default: 10.
	XPos *int64 `json:"XPos,omitempty" name:"XPos"`

	// The watermark’s vertical distance from the origin as a percentage of the video height. Value range: 0-100. Default: 10.
	YPos *int64 `json:"YPos,omitempty" name:"YPos"`

	// Font size. Value range: 25-50.
	FontSize *int64 `json:"FontSize,omitempty" name:"FontSize"`

	// Font color, which is an RGB color value. Default value: 0x000000.
	FontColor *string `json:"FontColor,omitempty" name:"FontColor"`
}

type DashRemuxSettingsInfo struct {
	// Segment duration in ms. Value range: [1000,30000]. Default value: 4000. The value can only be a multiple of 1,000.
	SegmentDuration *uint64 `json:"SegmentDuration,omitempty" name:"SegmentDuration"`

	// Number of segments. Value range: [1,30]. Default value: 5.
	SegmentNumber *uint64 `json:"SegmentNumber,omitempty" name:"SegmentNumber"`

	// Whether to enable multi-period. Valid values: CLOSE/OPEN. Default value: CLOSE.
	PeriodTriggers *string `json:"PeriodTriggers,omitempty" name:"PeriodTriggers"`

	// The HLS package type when the H.265 codec is used. Valid values: `hvc1`, `hev1` (default).
	H265PackageType *string `json:"H265PackageType,omitempty" name:"H265PackageType"`
}

// Predefined struct for user
type DeleteStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`
}

type DeleteStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`
}

type DeleteStreamLiveInputRequest struct {
	*tchttp.BaseRequest
	
	// Input ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`
}

type DeleteStreamLiveInputSecurityGroupRequest struct {
	*tchttp.BaseRequest
	
	// Input security group ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Name of the event to delete
	EventName *string `json:"EventName,omitempty" name:"EventName"`
}

type DeleteStreamLivePlanRequest struct {
	*tchttp.BaseRequest
	
	// ID of the channel whose event is to be deleted
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Name of the event to delete
	EventName *string `json:"EventName,omitempty" name:"EventName"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`
}

type DeleteStreamLiveWatermarkRequest struct {
	*tchttp.BaseRequest
	
	// Watermark ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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

type DescribeImageSettings struct {
	// Origin
	Location *string `json:"Location,omitempty" name:"Location"`

	// The watermark image’s horizontal distance from the origin as a percentage of the video width
	XPos *int64 `json:"XPos,omitempty" name:"XPos"`

	// The watermark image’s vertical distance from the origin as a percentage of the video height
	YPos *int64 `json:"YPos,omitempty" name:"YPos"`

	// The watermark image’s width as a percentage of the video width
	Width *int64 `json:"Width,omitempty" name:"Width"`

	// The watermark image’s height as a percentage of the video height
	Height *int64 `json:"Height,omitempty" name:"Height"`
}

// Predefined struct for user
type DescribeStreamLiveChannelAlertsRequestParams struct {
	// Channel ID
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`
}

type DescribeStreamLiveChannelAlertsRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`
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
	// Alarm information of the channel’s two pipelines
	Infos *ChannelAlertInfos `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query statistics in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`

	// Data collection interval. Valid values: `5s`, `1min` (default), `5min`, `15min`
	Period *string `json:"Period,omitempty" name:"Period"`
}

type DescribeStreamLiveChannelInputStatisticsRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query statistics in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`

	// Data collection interval. Valid values: `5s`, `1min` (default), `5min`, `15min`
	Period *string `json:"Period,omitempty" name:"Period"`
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
	Infos []*ChannelInputStatistics `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query logs in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`
}

type DescribeStreamLiveChannelLogsRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query logs in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`
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
	Infos *PipelineLogInfo `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query statistics in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`

	// Data collection interval. Valid values: `5s`, `1min` (default), `5min`, `15min`
	Period *string `json:"Period,omitempty" name:"Period"`
}

type DescribeStreamLiveChannelOutputStatisticsRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// Start time for query, which is 1 hour ago by default. You can query statistics in the last 7 days.
	// UTC time, such as `2020-01-01T12:00:00Z`
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// End time for query, which is 1 hour after `StartTime` by default
	// UTC time, such as `2020-01-01T12:00:00Z`
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`

	// Data collection interval. Valid values: `5s`, `1min` (default), `5min`, `15min`
	Period *string `json:"Period,omitempty" name:"Period"`
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
	Infos []*ChannelOutputsStatistics `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`
}

type DescribeStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	Info *StreamLiveChannelInfo `json:"Info,omitempty" name:"Info"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Infos []*StreamLiveChannelInfo `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`
}

type DescribeStreamLiveInputRequest struct {
	*tchttp.BaseRequest
	
	// Input ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	Info *InputInfo `json:"Info,omitempty" name:"Info"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`
}

type DescribeStreamLiveInputSecurityGroupRequest struct {
	*tchttp.BaseRequest
	
	// Input security group ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	Info *InputSecurityGroupInfo `json:"Info,omitempty" name:"Info"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Infos []*InputSecurityGroupInfo `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Infos []*InputInfo `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`
}

type DescribeStreamLivePlansRequest struct {
	*tchttp.BaseRequest
	
	// ID of the channel whose events you want to query
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`
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
	Infos []*PlanResp `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Info *StreamLiveRegionInfo `json:"Info,omitempty" name:"Info"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	StartDayTime *string `json:"StartDayTime,omitempty" name:"StartDayTime"`

	// The query end time (UTC+8) in the format of yyyy-MM-dd.
	// You can only query data in the last month (not including the current day).
	EndDayTime *string `json:"EndDayTime,omitempty" name:"EndDayTime"`

	// The channel ID (optional).
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// The number of pages. Default value: 1.
	// The value cannot exceed 100.
	PageNum *int64 `json:"PageNum,omitempty" name:"PageNum"`

	// The number of records per page. Default value: 10.
	// Value range: 1-1000.
	PageSize *int64 `json:"PageSize,omitempty" name:"PageSize"`
}

type DescribeStreamLiveTranscodeDetailRequest struct {
	*tchttp.BaseRequest
	
	// The query start time (UTC+8) in the format of yyyy-MM-dd.
	// You can only query data in the last month (not including the current day).
	StartDayTime *string `json:"StartDayTime,omitempty" name:"StartDayTime"`

	// The query end time (UTC+8) in the format of yyyy-MM-dd.
	// You can only query data in the last month (not including the current day).
	EndDayTime *string `json:"EndDayTime,omitempty" name:"EndDayTime"`

	// The channel ID (optional).
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// The number of pages. Default value: 1.
	// The value cannot exceed 100.
	PageNum *int64 `json:"PageNum,omitempty" name:"PageNum"`

	// The number of records per page. Default value: 10.
	// Value range: 1-1000.
	PageSize *int64 `json:"PageSize,omitempty" name:"PageSize"`
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
	Infos []*DescribeTranscodeDetailInfo `json:"Infos,omitempty" name:"Infos"`

	// The number of the current page.
	PageNum *int64 `json:"PageNum,omitempty" name:"PageNum"`

	// The number of records per page.
	PageSize *int64 `json:"PageSize,omitempty" name:"PageSize"`

	// The total number of records.
	TotalNum *int64 `json:"TotalNum,omitempty" name:"TotalNum"`

	// The total number of pages.
	TotalPage *int64 `json:"TotalPage,omitempty" name:"TotalPage"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`
}

type DescribeStreamLiveWatermarkRequest struct {
	*tchttp.BaseRequest
	
	// Watermark ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	Info *DescribeWatermarkInfo `json:"Info,omitempty" name:"Info"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Infos []*DescribeWatermarkInfo `json:"Infos,omitempty" name:"Infos"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Text *string `json:"Text,omitempty" name:"Text"`

	// Origin
	Location *string `json:"Location,omitempty" name:"Location"`

	// The watermark image’s horizontal distance from the origin as a percentage of the video width
	XPos *int64 `json:"XPos,omitempty" name:"XPos"`

	// The watermark image’s vertical distance from the origin as a percentage of the video height
	YPos *int64 `json:"YPos,omitempty" name:"YPos"`

	// Font size
	FontSize *int64 `json:"FontSize,omitempty" name:"FontSize"`

	// Font color
	FontColor *string `json:"FontColor,omitempty" name:"FontColor"`
}

type DescribeTranscodeDetailInfo struct {
	// The channel ID.
	ChannelId *string `json:"ChannelId,omitempty" name:"ChannelId"`

	// The start time (UTC+8) of transcoding in the format of yyyy-MM-dd HH:mm:ss.
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// The end time (UTC+8) of transcoding in the format of yyyy-MM-dd HH:mm:ss.
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`

	// The duration (s) of transcoding.
	Duration *int64 `json:"Duration,omitempty" name:"Duration"`

	// The encoding method.
	// Examples:
	// `liveprocessor_H264`: Live transcoding-H264
	// `liveprocessor_H265`: Live transcoding-H265
	// `topspeed_H264`: Top speed codec-H264
	// `topspeed_H265`: Top speed codec-H265
	ModuleCodec *string `json:"ModuleCodec,omitempty" name:"ModuleCodec"`

	// The target bitrate (Kbps).
	Bitrate *int64 `json:"Bitrate,omitempty" name:"Bitrate"`

	// The transcoding type.
	Type *string `json:"Type,omitempty" name:"Type"`

	// The push domain name.
	PushDomain *string `json:"PushDomain,omitempty" name:"PushDomain"`

	// The target resolution.
	Resolution *string `json:"Resolution,omitempty" name:"Resolution"`
}

type DescribeWatermarkInfo struct {
	// Watermark ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// Watermark name
	Name *string `json:"Name,omitempty" name:"Name"`

	// Watermark type. Valid values: STATIC_IMAGE, TEXT.
	Type *string `json:"Type,omitempty" name:"Type"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	// Note: This field may return `null`, indicating that no valid value was found.
	ImageSettings *DescribeImageSettings `json:"ImageSettings,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	// Note: This field may return `null`, indicating that no valid value was found.
	TextSettings *DescribeTextSettings `json:"TextSettings,omitempty" name:"TextSettings"`

	// Last modified time (UTC+0) of the watermark, in the format of `2020-01-01T12:00:00Z`
	// Note: This field may return `null`, indicating that no valid value was found.
	UpdateTime *string `json:"UpdateTime,omitempty" name:"UpdateTime"`

	// List of channel IDs the watermark is bound to
	// Note: This field may return `null`, indicating that no valid value was found.
	AttachedChannels []*string `json:"AttachedChannels,omitempty" name:"AttachedChannels"`
}

type DestinationInfo struct {
	// Relay destination address. Length limit: [1,512].
	OutputUrl *string `json:"OutputUrl,omitempty" name:"OutputUrl"`

	// Authentication key. Length limit: [1,128].
	// Note: this field may return null, indicating that no valid values can be obtained.
	AuthKey *string `json:"AuthKey,omitempty" name:"AuthKey"`

	// Authentication username. Length limit: [1,128].
	// Note: this field may return null, indicating that no valid values can be obtained.
	Username *string `json:"Username,omitempty" name:"Username"`

	// Authentication password. Length limit: [1,128].
	// Note: this field may return null, indicating that no valid values can be obtained.
	Password *string `json:"Password,omitempty" name:"Password"`
}

type DrmKey struct {
	// DRM key, which is a 32-bit hexadecimal string.
	// Note: uppercase letters in the string will be automatically converted to lowercase ones.
	Key *string `json:"Key,omitempty" name:"Key"`

	// Required for Widevine encryption. Valid values: SD, HD, UHD1, UHD2, AUDIO, ALL.
	// ALL refers to all tracks. If this parameter is set to ALL, no other tracks can be added.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Track *string `json:"Track,omitempty" name:"Track"`

	// Required for Widevine encryption. It is a 32-bit hexadecimal string.
	// Note: uppercase letters in the string will be automatically converted to lowercase ones.
	// Note: this field may return null, indicating that no valid values can be obtained.
	KeyId *string `json:"KeyId,omitempty" name:"KeyId"`

	// Required when FairPlay uses the AES encryption method. It is a 32-bit hexadecimal string.
	// For more information about this parameter, please see: 
	// https://tools.ietf.org/html/rfc3826
	// Note: uppercase letters in the string will be automatically converted to lowercase ones.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Iv *string `json:"Iv,omitempty" name:"Iv"`

	// The URI of the license server when AES-128 is used. This parameter may be empty.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	KeyUri *string `json:"KeyUri,omitempty" name:"KeyUri"`
}

type DrmSettingsInfo struct {
	// Whether to enable DRM encryption. Valid values: `CLOSE` (disable), `OPEN` (enable). Default value: `CLOSE`
	// DRM encryption is supported only for HLS, DASH, HLS_ARCHIVE, DASH_ARCHIVE, HLS_MEDIAPACKAGE, and DASH_MEDIAPACKAGE outputs.
	State *string `json:"State,omitempty" name:"State"`

	// Valid values: `CustomDRMKeys` (default value), `SDMCDRM`
	// `CustomDRMKeys` means encryption keys customized by users.
	// `SDMCDRM` means the DRM key management system of SDMC.
	Scheme *string `json:"Scheme,omitempty" name:"Scheme"`

	// If `Scheme` is set to `CustomDRMKeys`, this parameter is required.
	// If `Scheme` is set to `SDMCDRM`, this parameter is optional. It supports digits, letters, hyphens, and underscores and must contain 1 to 36 characters. If it is not specified, the value of `ChannelId` will be used.
	ContentId *string `json:"ContentId,omitempty" name:"ContentId"`

	// The key customized by the content user, which is required when `Scheme` is set to CustomDRMKeys.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Keys []*DrmKey `json:"Keys,omitempty" name:"Keys"`

	// SDMC key configuration. This parameter is used when `Scheme` is set to `SDMCDRM`.
	// Note: This field may return `null`, indicating that no valid value was found.
	SDMCSettings *SDMCSettingsInfo `json:"SDMCSettings,omitempty" name:"SDMCSettings"`

	// The DRM type. Valid values: `FAIRPLAY`, `WIDEVINE`, `AES128`. For HLS, this can be `FAIRPLAY` or `AES128`. For DASH, this can only be `WIDEVINE`.
	// Note: This field may return `null`, indicating that no valid values can be obtained.
	DrmType *string `json:"DrmType,omitempty" name:"DrmType"`
}

type EventNotifySetting struct {
	// The callback configuration for push events.
	PushEventSettings *PushEventSetting `json:"PushEventSettings,omitempty" name:"PushEventSettings"`
}

type EventSettingsDestinationReq struct {
	// URL of the COS bucket to save recording files
	Url *string `json:"Url,omitempty" name:"Url"`
}

type EventSettingsDestinationResp struct {
	// URL of the COS bucket where recording files are saved
	Url *string `json:"Url,omitempty" name:"Url"`
}

type EventSettingsReq struct {
	// Valid values: `INPUT_SWITCH`, `TIMED_RECORD`. If it is not specified, `INPUT_SWITCH` will be used.
	EventType *string `json:"EventType,omitempty" name:"EventType"`

	// ID of the input to attach, which is required if `EventType` is `INPUT_SWITCH`
	InputAttachment *string `json:"InputAttachment,omitempty" name:"InputAttachment"`

	// Name of the output group to attach. This parameter is required if `EventType` is `TIMED_RECORD`.
	OutputGroupName *string `json:"OutputGroupName,omitempty" name:"OutputGroupName"`

	// Name of the manifest file for timed recording, which must end with `.m3u8` for HLS and `.mpd` for DASH. This parameter is required if `EventType` is `TIMED_RECORD`.
	ManifestName *string `json:"ManifestName,omitempty" name:"ManifestName"`

	// URL of the COS bucket to save recording files. This parameter is required if `EventType` is `TIMED_RECORD`. It may contain 1 or 2 URLs. The first URL corresponds to pipeline 0 and the second pipeline 1.
	Destinations []*EventSettingsDestinationReq `json:"Destinations,omitempty" name:"Destinations"`
}

type EventSettingsResp struct {
	// Only `INPUT_SWITCH` is supported currently.
	EventType *string `json:"EventType,omitempty" name:"EventType"`

	// ID of the input attached, which is not empty if `EventType` is `INPUT_SWITCH`
	InputAttachment *string `json:"InputAttachment,omitempty" name:"InputAttachment"`

	// Name of the output group attached. This parameter is not empty if `EventType` is `TIMED_RECORD`.
	OutputGroupName *string `json:"OutputGroupName,omitempty" name:"OutputGroupName"`

	// Name of the manifest file for timed recording, which ends with `.m3u8` for HLS and `.mpd` for DASH. This parameter is not empty if `EventType` is `TIMED_RECORD`.
	ManifestName *string `json:"ManifestName,omitempty" name:"ManifestName"`

	// URL of the COS bucket where recording files are saved. This parameter is not empty if `EventType` is `TIMED_RECORD`. It may contain 1 or 2 URLs. The first URL corresponds to pipeline 0 and the second pipeline 1.
	Destinations []*EventSettingsDestinationResp `json:"Destinations,omitempty" name:"Destinations"`
}

type FailOverSettings struct {
	// ID of the backup input
	// Note: this field may return `null`, indicating that no valid value was found.
	SecondaryInputId *string `json:"SecondaryInputId,omitempty" name:"SecondaryInputId"`

	// The wait time (ms) for triggering failover after the primary input becomes unavailable. Value range: [1000, 86400000]. Default value: `3000`
	LossThreshold *int64 `json:"LossThreshold,omitempty" name:"LossThreshold"`

	// Failover policy. Valid values: `CURRENT_PREFERRED` (default), `PRIMARY_PREFERRED`
	RecoverBehavior *string `json:"RecoverBehavior,omitempty" name:"RecoverBehavior"`
}

type HlsRemuxSettingsInfo struct {
	// Segment duration in ms. Value range: [1000,30000]. Default value: 4000. The value can only be a multiple of 1,000.
	SegmentDuration *uint64 `json:"SegmentDuration,omitempty" name:"SegmentDuration"`

	// Number of segments. Value range: [1,30]. Default value: 5.
	SegmentNumber *uint64 `json:"SegmentNumber,omitempty" name:"SegmentNumber"`

	// Whether to enable PDT insertion. Valid values: CLOSE/OPEN. Default value: CLOSE.
	PdtInsertion *string `json:"PdtInsertion,omitempty" name:"PdtInsertion"`

	// PDT duration in seconds. Value range: (0,3000]. Default value: 600.
	PdtDuration *uint64 `json:"PdtDuration,omitempty" name:"PdtDuration"`

	// Audio/Video packaging scheme. Valid values: `SEPARATE`, `MERGE`
	Scheme *string `json:"Scheme,omitempty" name:"Scheme"`

	// The segment type. Valid values: `ts` (default), `fmp4`.
	// Currently, fMP4 segments do not support DRM or time shifting.
	SegmentType *string `json:"SegmentType,omitempty" name:"SegmentType"`

	// The HLS package type when the H.265 codec is used. Valid values: `hvc1`, `hev1` (default).
	H265PackageType *string `json:"H265PackageType,omitempty" name:"H265PackageType"`
}

type InputInfo struct {
	// Input region.
	Region *string `json:"Region,omitempty" name:"Region"`

	// Input ID.
	Id *string `json:"Id,omitempty" name:"Id"`

	// Input name.
	Name *string `json:"Name,omitempty" name:"Name"`

	// Input type.
	Type *string `json:"Type,omitempty" name:"Type"`

	// Array of security groups associated with input.
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitempty" name:"SecurityGroupIds"`

	// Array of channels associated with input.
	// Note: this field may return null, indicating that no valid values can be obtained.
	AttachedChannels []*string `json:"AttachedChannels,omitempty" name:"AttachedChannels"`

	// Input configuration array.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitempty" name:"InputSettings"`
}

type InputSecurityGroupInfo struct {
	// Input security group ID.
	Id *string `json:"Id,omitempty" name:"Id"`

	// Input security group name.
	Name *string `json:"Name,omitempty" name:"Name"`

	// List of allowlist entries.
	Whitelist []*string `json:"Whitelist,omitempty" name:"Whitelist"`

	// List of bound input streams.
	// Note: this field may return null, indicating that no valid values can be obtained.
	OccupiedInputs []*string `json:"OccupiedInputs,omitempty" name:"OccupiedInputs"`

	// Input security group address.
	Region *string `json:"Region,omitempty" name:"Region"`
}

type InputSettingInfo struct {
	// Application name, which is valid if `Type` is `RTMP_PUSH` and can contain 1-32 letters and digits
	// Note: This field may return `null`, indicating that no valid value was found.
	AppName *string `json:"AppName,omitempty" name:"AppName"`

	// Stream name, which is valid if `Type` is `RTMP_PUSH` and can contain 1-32 letters and digits
	// Note: This field may return `null`, indicating that no valid value was found.
	StreamName *string `json:"StreamName,omitempty" name:"StreamName"`

	// Source URL, which is valid if `Type` is `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL` and can contain 1-512 characters
	// Note: This field may return `null`, indicating that no valid value was found.
	SourceUrl *string `json:"SourceUrl,omitempty" name:"SourceUrl"`

	// RTP/UDP input address, which does not need to be entered for the input parameter.
	// Note: this field may return null, indicating that no valid values can be obtained.
	InputAddress *string `json:"InputAddress,omitempty" name:"InputAddress"`

	// Source type for stream pulling and relaying. To pull content from private-read COS buckets under the current account, set this parameter to `TencentCOS`; otherwise, leave it empty.
	// Note: this field may return `null`, indicating that no valid value was found.
	SourceType *string `json:"SourceType,omitempty" name:"SourceType"`

	// Delayed time (ms) for playback, which is valid if `Type` is `RTMP_PUSH`
	// Value range: 0 (default) or 10000-600000
	// The value must be a multiple of 1,000.
	// Note: This field may return `null`, indicating that no valid value was found.
	DelayTime *int64 `json:"DelayTime,omitempty" name:"DelayTime"`

	// The domain of an SRT_PUSH address. If this is a request parameter, you don’t need to specify it.
	// Note: This field may return `null`, indicating that no valid value was found.
	InputDomain *string `json:"InputDomain,omitempty" name:"InputDomain"`

	// The username, which is used for authentication.
	// Note: This field may return `null`, indicating that no valid value was found.
	UserName *string `json:"UserName,omitempty" name:"UserName"`

	// The password, which is used for authentication.
	// Note: This field may return `null`, indicating that no valid value was found.
	Password *string `json:"Password,omitempty" name:"Password"`
}

type InputStatistics struct {
	// Input statistics of pipeline 0.
	Pipeline0 []*PipelineInputStatistics `json:"Pipeline0,omitempty" name:"Pipeline0"`

	// Input statistics of pipeline 1.
	Pipeline1 []*PipelineInputStatistics `json:"Pipeline1,omitempty" name:"Pipeline1"`
}

type InputStreamInfo struct {
	// The input stream address.
	InputAddress *string `json:"InputAddress,omitempty" name:"InputAddress"`

	// The input stream path.
	AppName *string `json:"AppName,omitempty" name:"AppName"`

	// The input stream name.
	StreamName *string `json:"StreamName,omitempty" name:"StreamName"`

	// The input stream status. `1` indicates the stream is active.
	Status *int64 `json:"Status,omitempty" name:"Status"`
}

type LogInfo struct {
	// Log type.
	// It contains the value of `StreamStart` which refers to the push information.
	Type *string `json:"Type,omitempty" name:"Type"`

	// Time when the log is printed.
	Time *string `json:"Time,omitempty" name:"Time"`

	// Log details.
	Message *LogMessageInfo `json:"Message,omitempty" name:"Message"`
}

type LogMessageInfo struct {
	// Push information.
	// Note: this field may return null, indicating that no valid values can be obtained.
	StreamInfo *StreamInfo `json:"StreamInfo,omitempty" name:"StreamInfo"`
}

// Predefined struct for user
type ModifyStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// Channel name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Inputs to attach. You can attach 1 to 5 inputs.
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitempty" name:"AttachedInputs"`

	// Configuration information of the channel’s output groups. Quantity: [1, 10]
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitempty" name:"OutputGroups"`

	// Audio transcoding templates. Quantity: [1, 20]
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitempty" name:"AudioTemplates"`

	// Video transcoding templates. Quantity: [1, 10]
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates. Quantity: [1, 10]
	AVTemplates []*AVTemplate `json:"AVTemplates,omitempty" name:"AVTemplates"`

	// Event settings
	PlanSettings *PlanSettings `json:"PlanSettings,omitempty" name:"PlanSettings"`

	// The callback settings.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitempty" name:"EventNotifySettings"`
}

type ModifyStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// Channel name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Inputs to attach. You can attach 1 to 5 inputs.
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitempty" name:"AttachedInputs"`

	// Configuration information of the channel’s output groups. Quantity: [1, 10]
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitempty" name:"OutputGroups"`

	// Audio transcoding templates. Quantity: [1, 20]
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitempty" name:"AudioTemplates"`

	// Video transcoding templates. Quantity: [1, 10]
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates. Quantity: [1, 10]
	AVTemplates []*AVTemplate `json:"AVTemplates,omitempty" name:"AVTemplates"`

	// Event settings
	PlanSettings *PlanSettings `json:"PlanSettings,omitempty" name:"PlanSettings"`

	// The callback settings.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitempty" name:"EventNotifySettings"`
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
	delete(f, "PlanSettings")
	delete(f, "EventNotifySettings")
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyStreamLiveChannelRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveChannelResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`

	// Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// List of the IDs of the security groups to attach
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitempty" name:"SecurityGroupIds"`

	// Input settings
	// For the type `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
	// This parameter can be left empty for RTP_PUSH and UDP_PUSH inputs.
	// Note: If this parameter is not specified or empty, the original input settings will be used.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitempty" name:"InputSettings"`
}

type ModifyStreamLiveInputRequest struct {
	*tchttp.BaseRequest
	
	// Input ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// List of the IDs of the security groups to attach
	SecurityGroupIds []*string `json:"SecurityGroupIds,omitempty" name:"SecurityGroupIds"`

	// Input settings
	// For the type `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.
	// This parameter can be left empty for RTP_PUSH and UDP_PUSH inputs.
	// Note: If this parameter is not specified or empty, the original input settings will be used.
	InputSettings []*InputSettingInfo `json:"InputSettings,omitempty" name:"InputSettings"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`

	// Input security group name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Allowlist entries (max: 10)
	Whitelist []*string `json:"Whitelist,omitempty" name:"Whitelist"`
}

type ModifyStreamLiveInputSecurityGroupRequest struct {
	*tchttp.BaseRequest
	
	// Input security group ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// Input security group name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Allowlist entries (max: 10)
	Whitelist []*string `json:"Whitelist,omitempty" name:"Whitelist"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Id *string `json:"Id,omitempty" name:"Id"`

	// Watermark name
	Name *string `json:"Name,omitempty" name:"Name"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	ImageSettings *CreateImageSettings `json:"ImageSettings,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	TextSettings *CreateTextSettings `json:"TextSettings,omitempty" name:"TextSettings"`
}

type ModifyStreamLiveWatermarkRequest struct {
	*tchttp.BaseRequest
	
	// Watermark ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// Watermark name
	Name *string `json:"Name,omitempty" name:"Name"`

	// Watermark image settings. This parameter is valid if `Type` is `STATIC_IMAGE`.
	ImageSettings *CreateImageSettings `json:"ImageSettings,omitempty" name:"ImageSettings"`

	// Watermark text settings. This parameter is valid if `Type` is `TEXT`.
	TextSettings *CreateTextSettings `json:"TextSettings,omitempty" name:"TextSettings"`
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
	if len(f) > 0 {
		return tcerr.NewTencentCloudSDKError("ClientError.BuildRequestError", "ModifyStreamLiveWatermarkRequest has unknown keys!", "")
	}
	return json.Unmarshal([]byte(s), &r)
}

// Predefined struct for user
type ModifyStreamLiveWatermarkResponseParams struct {
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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

type OutputInfo struct {
	// Output name.
	Name *string `json:"Name,omitempty" name:"Name"`

	// Audio transcoding template name array.
	// Quantity limit: [0,1] for RTMP; [0,20] for others.
	// Note: this field may return null, indicating that no valid values can be obtained.
	AudioTemplateNames []*string `json:"AudioTemplateNames,omitempty" name:"AudioTemplateNames"`

	// Video transcoding template name array. Quantity limit: [0,1].
	// Note: this field may return null, indicating that no valid values can be obtained.
	VideoTemplateNames []*string `json:"VideoTemplateNames,omitempty" name:"VideoTemplateNames"`

	// SCTE-35 information configuration.
	Scte35Settings *Scte35SettingsInfo `json:"Scte35Settings,omitempty" name:"Scte35Settings"`

	// Audio/Video transcoding template name. If `HlsRemuxSettings.Scheme` is `MERGE`, there is 1 audio/video transcoding template. Otherwise, this parameter is empty.
	// Note: this field may return `null`, indicating that no valid value was found.
	AVTemplateNames []*string `json:"AVTemplateNames,omitempty" name:"AVTemplateNames"`
}

type OutputsStatistics struct {
	// Output information of pipeline 0.
	Pipeline0 []*PipelineOutputStatistics `json:"Pipeline0,omitempty" name:"Pipeline0"`

	// Output information of pipeline 1.
	Pipeline1 []*PipelineOutputStatistics `json:"Pipeline1,omitempty" name:"Pipeline1"`
}

type PipelineInputStatistics struct {
	// Data timestamp in seconds.
	Timestamp *uint64 `json:"Timestamp,omitempty" name:"Timestamp"`

	// Input bandwidth in bps.
	NetworkIn *uint64 `json:"NetworkIn,omitempty" name:"NetworkIn"`

	// Video information array.
	// For `rtp/udp` input, the quantity is the number of `Pid` of the input video.
	// For other inputs, the quantity is 1.
	Video []*VideoPipelineInputStatistics `json:"Video,omitempty" name:"Video"`

	// Audio information array.
	// For `rtp/udp` input, the quantity is the number of `Pid` of the input audio.
	// For other inputs, the quantity is 1.
	Audio []*AudioPipelineInputStatistics `json:"Audio,omitempty" name:"Audio"`
}

type PipelineLogInfo struct {
	// Log information of pipeline 0.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pipeline0 []*LogInfo `json:"Pipeline0,omitempty" name:"Pipeline0"`

	// Log information of pipeline 1.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pipeline1 []*LogInfo `json:"Pipeline1,omitempty" name:"Pipeline1"`
}

type PipelineOutputStatistics struct {
	// Timestamp.
	// In seconds, indicating data time.
	Timestamp *uint64 `json:"Timestamp,omitempty" name:"Timestamp"`

	// Output bandwidth in bps.
	NetworkOut *uint64 `json:"NetworkOut,omitempty" name:"NetworkOut"`
}

type PlanReq struct {
	// Event name
	EventName *string `json:"EventName,omitempty" name:"EventName"`

	// Event trigger time settings
	TimingSettings *TimingSettingsReq `json:"TimingSettings,omitempty" name:"TimingSettings"`

	// Event configuration
	EventSettings *EventSettingsReq `json:"EventSettings,omitempty" name:"EventSettings"`
}

type PlanResp struct {
	// Event name
	EventName *string `json:"EventName,omitempty" name:"EventName"`

	// Event trigger time settings
	TimingSettings *TimingSettingsResp `json:"TimingSettings,omitempty" name:"TimingSettings"`

	// Event configuration
	EventSettings *EventSettingsResp `json:"EventSettings,omitempty" name:"EventSettings"`
}

type PlanSettings struct {
	// Timed recording settings
	// Note: This field may return `null`, indicating that no valid value was found.
	TimedRecordSettings *TimedRecordSettings `json:"TimedRecordSettings,omitempty" name:"TimedRecordSettings"`
}

type PushEventSetting struct {
	// The callback URL (required).
	NotifyUrl *string `json:"NotifyUrl,omitempty" name:"NotifyUrl"`

	// The callback key (optional).
	NotifyKey *string `json:"NotifyKey,omitempty" name:"NotifyKey"`
}

type QueryDispatchInputInfo struct {
	// The input ID.
	InputID *string `json:"InputID,omitempty" name:"InputID"`

	// The input name.
	InputName *string `json:"InputName,omitempty" name:"InputName"`

	// The input protocol.
	Protocol *string `json:"Protocol,omitempty" name:"Protocol"`

	// The stream status of the input.
	InputStreamInfoList []*InputStreamInfo `json:"InputStreamInfoList,omitempty" name:"InputStreamInfoList"`
}

// Predefined struct for user
type QueryInputStreamStateRequestParams struct {
	// The StreamLive input ID.
	Id *string `json:"Id,omitempty" name:"Id"`
}

type QueryInputStreamStateRequest struct {
	*tchttp.BaseRequest
	
	// The StreamLive input ID.
	Id *string `json:"Id,omitempty" name:"Id"`
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
	Info *QueryDispatchInputInfo `json:"Info,omitempty" name:"Info"`

	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Name *string `json:"Name,omitempty" name:"Name"`
}

type SDMCSettingsInfo struct {
	// User ID in the SDMC DRM system
	Uid *string `json:"Uid,omitempty" name:"Uid"`

	// Tracks of the SDMC DRM system. This parameter is valid for DASH output groups.
	// `1`: audio
	// `2`: SD
	// `4`: HD
	// `8`: UHD1
	// `16`: UHD2
	// 
	// Default value: `31` (audio + SD + HD + UHD1 + UHD2)
	Tracks *int64 `json:"Tracks,omitempty" name:"Tracks"`

	// Key ID in the SDMC DRM system; required
	SecretId *string `json:"SecretId,omitempty" name:"SecretId"`

	// Key in the SDMC DRM system; required
	SecretKey *string `json:"SecretKey,omitempty" name:"SecretKey"`

	// Key request URL of the SDMC DRM system, which is `https://uat.multidrm.tv/cpix/2.2/getcontentkey` by default
	Url *string `json:"Url,omitempty" name:"Url"`

	// Token name in an SDMC key request URL, which is `token` by default
	TokenName *string `json:"TokenName,omitempty" name:"TokenName"`
}

type Scte35SettingsInfo struct {
	// Whether to pass through SCTE-35 information. Valid values: NO_PASSTHROUGH/PASSTHROUGH. Default value: NO_PASSTHROUGH.
	Behavior *string `json:"Behavior,omitempty" name:"Behavior"`
}

// Predefined struct for user
type StartStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`
}

type StartStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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

// Predefined struct for user
type StopStreamLiveChannelRequestParams struct {
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`
}

type StopStreamLiveChannelRequest struct {
	*tchttp.BaseRequest
	
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`
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
	// The unique request ID, which is returned for each request. RequestId is required for locating a problem.
	RequestId *string `json:"RequestId,omitempty" name:"RequestId"`
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
	Pid *int64 `json:"Pid,omitempty" name:"Pid"`

	// Audio codec.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Codec *string `json:"Codec,omitempty" name:"Codec"`

	// Audio frame rate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Fps *int64 `json:"Fps,omitempty" name:"Fps"`

	// Audio bitrate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Rate *int64 `json:"Rate,omitempty" name:"Rate"`

	// Audio sample rate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	SampleRate *int64 `json:"SampleRate,omitempty" name:"SampleRate"`
}

type StreamInfo struct {
	// Client IP.
	ClientIp *string `json:"ClientIp,omitempty" name:"ClientIp"`

	// Video information of pushed streams.
	Video []*StreamVideoInfo `json:"Video,omitempty" name:"Video"`

	// Audio information of pushed streams.
	Audio []*StreamAudioInfo `json:"Audio,omitempty" name:"Audio"`

	// SCTE-35 information of pushed streams.
	Scte35 []*StreamScte35Info `json:"Scte35,omitempty" name:"Scte35"`
}

type StreamLiveChannelInfo struct {
	// Channel ID
	Id *string `json:"Id,omitempty" name:"Id"`

	// Channel status
	State *string `json:"State,omitempty" name:"State"`

	// Information of attached inputs
	AttachedInputs []*AttachedInput `json:"AttachedInputs,omitempty" name:"AttachedInputs"`

	// Information of output groups
	OutputGroups []*StreamLiveOutputGroupsInfo `json:"OutputGroups,omitempty" name:"OutputGroups"`

	// Channel name
	Name *string `json:"Name,omitempty" name:"Name"`

	// Audio transcoding templates
	// Note: this field may return `null`, indicating that no valid value was found.
	AudioTemplates []*AudioTemplateInfo `json:"AudioTemplates,omitempty" name:"AudioTemplates"`

	// Video transcoding templates
	// Note: this field may return `null`, indicating that no valid value was found.
	VideoTemplates []*VideoTemplateInfo `json:"VideoTemplates,omitempty" name:"VideoTemplates"`

	// Audio/Video transcoding templates
	// Note: this field may return `null`, indicating that no valid value was found.
	AVTemplates []*AVTemplate `json:"AVTemplates,omitempty" name:"AVTemplates"`

	// Event settings
	// Note: This field may return `null`, indicating that no valid value was found.
	PlanSettings *PlanSettings `json:"PlanSettings,omitempty" name:"PlanSettings"`

	// The callback settings.
	// Note: This field may return `null`, indicating that no valid value was found.
	EventNotifySettings *EventNotifySetting `json:"EventNotifySettings,omitempty" name:"EventNotifySettings"`
}

type StreamLiveOutputGroupsInfo struct {
	// Output group name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the channel level
	Name *string `json:"Name,omitempty" name:"Name"`

	// Output protocol
	// Valid values: `HLS`, `DASH`, `HLS_ARCHIVE`, `HLS_STREAM_PACKAGE`, `DASH_STREAM_PACKAGE`
	Type *string `json:"Type,omitempty" name:"Type"`

	// Output information
	// If the type is RTMP or RTP, only one output is allowed; if it is HLS or DASH, 1-10 outputs are allowed.
	Outputs []*OutputInfo `json:"Outputs,omitempty" name:"Outputs"`

	// Relay destinations. Quantity: [1, 2]
	Destinations []*DestinationInfo `json:"Destinations,omitempty" name:"Destinations"`

	// HLS protocol configuration information, which takes effect only for HLS/HLS_ARCHIVE outputs
	// Note: this field may return `null`, indicating that no valid value was found.
	HlsRemuxSettings *HlsRemuxSettingsInfo `json:"HlsRemuxSettings,omitempty" name:"HlsRemuxSettings"`

	// DRM configuration information
	// Note: this field may return `null`, indicating that no valid value was found.
	DrmSettings *DrmSettingsInfo `json:"DrmSettings,omitempty" name:"DrmSettings"`

	// DASH protocol configuration information, which takes effect only for DASH/DASH_ARCHIVE outputs
	// Note: this field may return `null`, indicating that no valid value was found.
	DashRemuxSettings *DashRemuxSettingsInfo `json:"DashRemuxSettings,omitempty" name:"DashRemuxSettings"`

	// StreamPackage configuration information, which is required if the output type is StreamPackage
	// Note: this field may return `null`, indicating that no valid value was found.
	StreamPackageSettings *StreamPackageSettingsInfo `json:"StreamPackageSettings,omitempty" name:"StreamPackageSettings"`

	// Time-shift configuration information
	// Note: This field may return `null`, indicating that no valid value was found.
	TimeShiftSettings *TimeShiftSettingsInfo `json:"TimeShiftSettings,omitempty" name:"TimeShiftSettings"`
}

type StreamLiveRegionInfo struct {
	// List of StreamLive regions
	Regions []*RegionInfo `json:"Regions,omitempty" name:"Regions"`
}

type StreamPackageSettingsInfo struct {
	// Channel ID in StreamPackage
	Id *string `json:"Id,omitempty" name:"Id"`
}

type StreamScte35Info struct {
	// SCTE-35 `Pid`.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pid *int64 `json:"Pid,omitempty" name:"Pid"`
}

type StreamVideoInfo struct {
	// Video `Pid`.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Pid *int64 `json:"Pid,omitempty" name:"Pid"`

	// Video codec.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Codec *string `json:"Codec,omitempty" name:"Codec"`

	// Video frame rate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Fps *int64 `json:"Fps,omitempty" name:"Fps"`

	// Video bitrate.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Rate *int64 `json:"Rate,omitempty" name:"Rate"`

	// Video width.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Width *int64 `json:"Width,omitempty" name:"Width"`

	// Video height.
	// Note: this field may return null, indicating that no valid values can be obtained.
	Height *int64 `json:"Height,omitempty" name:"Height"`
}

type TimeShiftSettingsInfo struct {
	// Whether to enable time shifting. Valid values: `OPEN`; `CLOSE`
	// Note: This field may return `null`, indicating that no valid value was found.
	State *string `json:"State,omitempty" name:"State"`

	// Domain name bound for time shifting
	// Note: This field may return `null`, indicating that no valid value was found.
	PlayDomain *string `json:"PlayDomain,omitempty" name:"PlayDomain"`

	// Allowable time-shift period (s). Value range: [600, 1209600]. Default value: 300
	// Note: This field may return `null`, indicating that no valid value was found.
	StartoverWindow *int64 `json:"StartoverWindow,omitempty" name:"StartoverWindow"`
}

type TimedRecordSettings struct {
	// Whether to automatically delete finished recording events. Valid values: `CLOSE`, `OPEN`. If this parameter is left empty, `CLOSE` will be used.
	// If it is set to `OPEN`, a recording event will be deleted 7 days after it is finished.
	// Note: This field may return `null`, indicating that no valid value was found.
	AutoClear *string `json:"AutoClear,omitempty" name:"AutoClear"`
}

type TimingSettingsReq struct {
	// Event trigger type. Valid values: `FIXED_TIME`, `IMMEDIATE`. This parameter is required if `EventType` is `INPUT_SWITCH`.
	StartType *string `json:"StartType,omitempty" name:"StartType"`

	// This parameter is required if `EventType` is `INPUT_SWITCH` and `StartType` is `FIXED_TIME`.
	// It must be in UTC format, e.g., `2020-01-01T12:00:00Z`.
	Time *string `json:"Time,omitempty" name:"Time"`

	// This parameter is required if `EventType` is `TIMED_RECORD`.
	// It specifies the recording start time in UTC format (e.g., `2020-01-01T12:00:00Z`) and must be at least 1 minute later than the current time.
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// This parameter is required if `EventType` is `TIMED_RECORD`.
	// It specifies the recording end time in UTC format (e.g., `2020-01-01T12:00:00Z`) and must be at least 1 minute later than the recording start time.
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`
}

type TimingSettingsResp struct {
	// Event trigger type
	StartType *string `json:"StartType,omitempty" name:"StartType"`

	// Not empty if `StartType` is `FIXED_TIME`
	// UTC time, such as `2020-01-01T12:00:00Z`
	Time *string `json:"Time,omitempty" name:"Time"`

	// This parameter cannot be empty if `EventType` is `TIMED_RECORD`.
	// It indicates the start time for recording in UTC format (e.g., `2020-01-01T12:00:00Z`) and must be at least 1 minute later than the current time.
	StartTime *string `json:"StartTime,omitempty" name:"StartTime"`

	// This parameter cannot be empty if `EventType` is `TIMED_RECORD`.
	// It indicates the end time for recording in UTC format (e.g., `2020-01-01T12:00:00Z`) and must be at least 1 minute later than the start time for recording.
	EndTime *string `json:"EndTime,omitempty" name:"EndTime"`
}

type VideoPipelineInputStatistics struct {
	// Video FPS.
	Fps *uint64 `json:"Fps,omitempty" name:"Fps"`

	// Video bitrate in bps.
	Rate *uint64 `json:"Rate,omitempty" name:"Rate"`

	// Video `Pid`, which is available only if the input is `rtp/udp`.
	Pid *int64 `json:"Pid,omitempty" name:"Pid"`
}

type VideoTemplateInfo struct {
	// Video transcoding template name, which can contain 1-20 letters and digits.
	Name *string `json:"Name,omitempty" name:"Name"`

	// Video codec. Valid values: H264/H265. If this parameter is left empty, the original value will be used.
	Vcodec *string `json:"Vcodec,omitempty" name:"Vcodec"`

	// Video bitrate. Value range: [50000,40000000]. The value can only be a multiple of 1,000. If this parameter is left empty, the original value will be used.
	VideoBitrate *uint64 `json:"VideoBitrate,omitempty" name:"VideoBitrate"`

	// Video width. Value range: (0,3000]. The value can only be a multiple of 4. If this parameter is left empty, the original value will be used.
	Width *uint64 `json:"Width,omitempty" name:"Width"`

	// Video height. Value range: (0,3000]. The value can only be a multiple of 4. If this parameter is left empty, the original value will be used.
	Height *uint64 `json:"Height,omitempty" name:"Height"`

	// Video frame rate. Value range: [1,240]. If this parameter is left empty, the original value will be used.
	Fps *uint64 `json:"Fps,omitempty" name:"Fps"`

	// Whether to enable top speed codec. Valid value: CLOSE/OPEN. Default value: CLOSE.
	TopSpeed *string `json:"TopSpeed,omitempty" name:"TopSpeed"`

	// Top speed codec compression ratio. Value range: [0,50]. The lower the compression ratio, the higher the image quality.
	BitrateCompressionRatio *uint64 `json:"BitrateCompressionRatio,omitempty" name:"BitrateCompressionRatio"`

	// Bitrate control mode. Valid values: `CBR`, `ABR` (default)
	RateControlMode *string `json:"RateControlMode,omitempty" name:"RateControlMode"`

	// Watermark ID
	// Note: This field may return `null`, indicating that no valid value was found.
	WatermarkId *string `json:"WatermarkId,omitempty" name:"WatermarkId"`
}