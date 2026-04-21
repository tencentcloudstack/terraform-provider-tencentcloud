package teo

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestResourceTencentCloudTeoJustInTimeTranscodeTemplate_Create(t *testing.T) {
	var (
		zoneId            = "zone-1234567890"
		templateId        = "tpl-abcdefghij"
		templateName      = "test-template"
		comment           = "test comment"
		videoStreamSwitch = "on"
		audioStreamSwitch = "on"
	)

	r := ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := r.Data(nil)
	d.Set("zone_id", zoneId)
	d.Set("template_name", templateName)
	d.Set("comment", comment)
	d.Set("video_stream_switch", videoStreamSwitch)
	d.Set("audio_stream_switch", audioStreamSwitch)
	d.Set("video_template", []interface{}{
		map[string]interface{}{
			"video_codec":         "H.264",
			"fps":                 float64(30),
			"bitrate":             2000,
			"resolution_adaptive": "open",
			"width":               1280,
			"height":              720,
			"fill_type":           "black",
		},
	})
	d.Set("audio_template", []interface{}{
		map[string]interface{}{
			"codec":         "libfdk_aac",
			"audio_channel": 2,
		},
	})

	createResp := &teo.CreateJustInTimeTranscodeTemplateResponse{
		Response: &teo.CreateJustInTimeTranscodeTemplateResponseParams{
			TemplateId: helper.String(templateId),
			RequestId:  helper.String("test-request-id"),
		},
	}

	describeResp := &teo.DescribeJustInTimeTranscodeTemplatesResponse{
		Response: &teo.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount: helper.Uint64(1),
			TemplateSet: []*teo.JustInTimeTranscodeTemplate{
				{
					TemplateId:        helper.String(templateId),
					TemplateName:      helper.String(templateName),
					Comment:           helper.String(comment),
					VideoStreamSwitch: helper.String(videoStreamSwitch),
					AudioStreamSwitch: helper.String(audioStreamSwitch),
					Type:              helper.String("custom"),
					VideoTemplate: &teo.VideoTemplateInfo{
						Codec:              helper.String("H.264"),
						Fps:                helper.Float64(30),
						Bitrate:            helper.Uint64(2000),
						ResolutionAdaptive: helper.String("open"),
						Width:              helper.Uint64(1280),
						Height:             helper.Uint64(720),
						FillType:           helper.String("black"),
					},
					AudioTemplate: &teo.AudioTemplateInfo{
						Codec:        helper.String("libfdk_aac"),
						AudioChannel: helper.Uint64(2),
					},
				},
			},
			RequestId: helper.String("test-request-id"),
		},
	}

	// Mock CreateJustInTimeTranscodeTemplate
	patch1 := gomonkey.ApplyFunc(common.NewClient, func(region string, secretId, secretKey, token string, clientProfile *common.ClientProfile) (client *common.Client, err error) {
		return nil, nil
	})
	defer patch1.Reset()

	// Mock the client method
	patch2 := gomonkey.ApplyMethod((*resource.ResourceData)(nil), "Id", func(_ *resource.ResourceData) string {
		return ""
	})
	defer patch2.Reset()

	// Test Create
	createFunc := func() (interface{}, error) {
		meta := &tccommon.ProviderMeta{
			TencentV2Client: &tencentcloudV2Client{
				teoClient: &teo.Client{},
			},
		}
		return nil, resourceTencentCloudTeoJustInTimeTranscodeTemplateCreate(d, meta)
	}

	// Mock API call
	patch3 := gomonkey.ApplyMethod(&teo.Client{}, "CreateJustInTimeTranscodeTemplate", func(_ *teo.Client, req *teo.CreateJustInTimeTranscodeTemplateRequest) (resp *teo.CreateJustInTimeTranscodeTemplateResponse, err error) {
		return createResp, nil
	})
	defer patch3.Reset()

	patch4 := gomonkey.ApplyMethod(&teo.Client{}, "DescribeJustInTimeTranscodeTemplates", func(_ *teo.Client, req *teo.DescribeJustInTimeTranscodeTemplatesRequest) (resp *teo.DescribeJustInTimeTranscodeTemplatesResponse, err error) {
		return describeResp, nil
	})
	defer patch4.Reset()

	err := createFunc()
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	expectedId := fmt.Sprintf("%s#%s", zoneId, templateId)
	if d.Id() != expectedId {
		t.Errorf("expected id %s, got %s", expectedId, d.Id())
	}
}

func TestResourceTencentCloudTeoJustInTimeTranscodeTemplate_Read(t *testing.T) {
	var (
		zoneId            = "zone-1234567890"
		templateId        = "tpl-abcdefghij"
		templateName      = "test-template"
		comment           = "test comment"
		videoStreamSwitch = "on"
		audioStreamSwitch = "on"
	)

	describeResp := &teo.DescribeJustInTimeTranscodeTemplatesResponse{
		Response: &teo.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount: helper.Uint64(1),
			TemplateSet: []*teo.JustInTimeTranscodeTemplate{
				{
					TemplateId:        helper.String(templateId),
					TemplateName:      helper.String(templateName),
					Comment:           helper.String(comment),
					VideoStreamSwitch: helper.String(videoStreamSwitch),
					AudioStreamSwitch: helper.String(audioStreamSwitch),
					Type:              helper.String("custom"),
					VideoTemplate: &teo.VideoTemplateInfo{
						Codec:              helper.String("H.264"),
						Fps:                helper.Float64(30),
						Bitrate:            helper.Uint64(2000),
						ResolutionAdaptive: helper.String("open"),
						Width:              helper.Uint64(1280),
						Height:             helper.Uint64(720),
						FillType:           helper.String("black"),
					},
					AudioTemplate: &teo.AudioTemplateInfo{
						Codec:        helper.String("libfdk_aac"),
						AudioChannel: helper.Uint64(2),
					},
				},
			},
			RequestId: helper.String("test-request-id"),
		},
	}

	r := ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := r.Data(nil)
	d.SetId(fmt.Sprintf("%s#%s", zoneId, templateId))

	// Mock API call
	patch := gomonkey.ApplyMethod(&teo.Client{}, "DescribeJustInTimeTranscodeTemplates", func(_ *teo.Client, req *teo.DescribeJustInTimeTranscodeTemplatesRequest) (resp *teo.DescribeJustInTimeTranscodeTemplatesResponse, err error) {
		return describeResp, nil
	})
	defer patch.Reset()

	readFunc := func() error {
		meta := &tccommon.ProviderMeta{
			TencentV2Client: &tencentcloudV2Client{
				teoClient: &teo.Client{},
			},
		}
		return resourceTencentCloudTeoJustInTimeTranscodeTemplateRead(context.Background(), d, meta)
	}

	err := readFunc()
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	// Verify schema values
	if d.Get("template_name") != templateName {
		t.Errorf("expected template_name %s, got %s", templateName, d.Get("template_name"))
	}
	if d.Get("comment") != comment {
		t.Errorf("expected comment %s, got %s", comment, d.Get("comment"))
	}
	if d.Get("video_stream_switch") != videoStreamSwitch {
		t.Errorf("expected video_stream_switch %s, got %s", videoStreamSwitch, d.Get("video_stream_switch"))
	}
	if d.Get("audio_stream_switch") != audioStreamSwitch {
		t.Errorf("expected audio_stream_switch %s, got %s", audioStreamSwitch, d.Get("audio_stream_switch"))
	}
}

func TestResourceTencentCloudTeoJustInTimeTranscodeTemplate_Delete(t *testing.T) {
	var (
		zoneId     = "zone-1234567890"
		templateId = "tpl-abcdefghij"
	)

	deleteResp := &teo.DeleteJustInTimeTranscodeTemplatesResponse{
		Response: &teo.DeleteJustInTimeTranscodeTemplatesResponseParams{
			RequestId: helper.String("test-request-id"),
		},
	}

	r := ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := r.Data(nil)
	d.SetId(fmt.Sprintf("%s#%s", zoneId, templateId))

	// Mock API call
	patch := gomonkey.ApplyMethod(&teo.Client{}, "DeleteJustInTimeTranscodeTemplates", func(_ *teo.Client, req *teo.DeleteJustInTimeTranscodeTemplatesRequest) (resp *teo.DeleteJustInTimeTranscodeTemplatesResponse, err error) {
		return deleteResp, nil
	})
	defer patch.Reset()

	deleteFunc := func() error {
		meta := &tccommon.ProviderMeta{
			TencentV2Client: &tencentcloudV2Client{
				teoClient: &teo.Client{},
			},
		}
		return resourceTencentCloudTeoJustInTimeTranscodeTemplateDelete(context.Background(), d, meta)
	}

	err := deleteFunc()
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if d.Id() != "" {
		t.Errorf("expected id to be empty after delete, got %s", d.Id())
	}
}

func TestResourceTencentCloudTeoJustInTimeTranscodeTemplate_ParseId(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		wantZone string
		wantTemp string
		wantErr  bool
	}{
		{
			name:     "valid id",
			id:       "zone-123#tpl-456",
			wantZone: "zone-123",
			wantTemp: "tpl-456",
			wantErr:  false,
		},
		{
			name:     "valid id with special chars",
			id:       "zone_abc-123#tpl_xyz-456",
			wantZone: "zone_abc-123",
			wantTemp: "tpl_xyz-456",
			wantErr:  false,
		},
		{
			name:    "invalid id - no separator",
			id:      "zone-123",
			wantErr: true,
		},
		{
			name:    "invalid id - too many separators",
			id:      "zone-123#tpl-456#extra",
			wantErr: true,
		},
		{
			name:    "invalid id - empty parts",
			id:      "#",
			wantErr: true,
		},
		{
			name:    "invalid id - empty first part",
			id:      "#tpl-456",
			wantErr: true,
		},
		{
			name:    "invalid id - empty second part",
			id:      "zone-123#",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zone, temp, err := resourceTencentCloudTeoJustInTimeTranscodeTemplateParseId(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if zone != tt.wantZone {
					t.Errorf("ParseId() zone = %v, want %v", zone, tt.wantZone)
				}
				if temp != tt.wantTemp {
					t.Errorf("ParseId() template = %v, want %v", temp, tt.wantTemp)
				}
			}
		})
	}
}

func TestBuildVideoTemplateInfo(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
		want *teo.VideoTemplateInfo
	}{
		{
			name: "full configuration",
			data: map[string]interface{}{
				"video_codec":         "H.264",
				"fps":                 float64(30),
				"bitrate":             2000,
				"resolution_adaptive": "open",
				"width":               1280,
				"height":              720,
				"fill_type":           "black",
			},
			want: &teo.VideoTemplateInfo{
				Codec:              helper.String("H.264"),
				Fps:                helper.Float64(30),
				Bitrate:            helper.Uint64(2000),
				ResolutionAdaptive: helper.String("open"),
				Width:              helper.Uint64(1280),
				Height:             helper.Uint64(720),
				FillType:           helper.String("black"),
			},
		},
		{
			name: "minimal configuration",
			data: map[string]interface{}{
				"video_codec": "H.265",
			},
			want: &teo.VideoTemplateInfo{
				Codec: helper.String("H.265"),
			},
		},
		{
			name: "empty configuration",
			data: map[string]interface{}{},
			want: &teo.VideoTemplateInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildVideoTemplateInfo(tt.data)
			if got.Codec != nil && tt.want.Codec != nil && *got.Codec != *tt.want.Codec {
				t.Errorf("Codec mismatch, got %v, want %v", *got.Codec, *tt.want.Codec)
			}
			if got.Fps != nil && tt.want.Fps != nil && *got.Fps != *tt.want.Fps {
				t.Errorf("Fps mismatch, got %v, want %v", *got.Fps, *tt.want.Fps)
			}
			if got.Bitrate != nil && tt.want.Bitrate != nil && *got.Bitrate != *tt.want.Bitrate {
				t.Errorf("Bitrate mismatch, got %v, want %v", *got.Bitrate, *tt.want.Bitrate)
			}
		})
	}
}

func TestBuildAudioTemplateInfo(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
		want *teo.AudioTemplateInfo
	}{
		{
			name: "full configuration",
			data: map[string]interface{}{
				"codec":         "libfdk_aac",
				"audio_channel": 2,
			},
			want: &teo.AudioTemplateInfo{
				Codec:        helper.String("libfdk_aac"),
				AudioChannel: helper.Uint64(2),
			},
		},
		{
			name: "minimal configuration",
			data: map[string]interface{}{
				"codec": "libfdk_aac",
			},
			want: &teo.AudioTemplateInfo{
				Codec: helper.String("libfdk_aac"),
			},
		},
		{
			name: "empty configuration",
			data: map[string]interface{}{},
			want: &teo.AudioTemplateInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildAudioTemplateInfo(tt.data)
			if got.Codec != nil && tt.want.Codec != nil && *got.Codec != *tt.want.Codec {
				t.Errorf("Codec mismatch, got %v, want %v", *got.Codec, *tt.want.Codec)
			}
			if got.AudioChannel != nil && tt.want.AudioChannel != nil && *got.AudioChannel != *tt.want.AudioChannel {
				t.Errorf("AudioChannel mismatch, got %v, want %v", *got.AudioChannel, *tt.want.AudioChannel)
			}
		})
	}
}

func TestMapVideoTemplateInfoToSchema(t *testing.T) {
	tests := []struct {
		name string
		info *teo.VideoTemplateInfo
		want map[string]interface{}
	}{
		{
			name: "full configuration",
			info: &teo.VideoTemplateInfo{
				Codec:              helper.String("H.264"),
				Fps:                helper.Float64(30),
				Bitrate:            helper.Uint64(2000),
				ResolutionAdaptive: helper.String("open"),
				Width:              helper.Uint64(1280),
				Height:             helper.Uint64(720),
				FillType:           helper.String("black"),
			},
			want: map[string]interface{}{
				"video_codec":         "H.264",
				"fps":                 float64(30),
				"bitrate":             2000,
				"resolution_adaptive": "open",
				"width":               1280,
				"height":              720,
				"fill_type":           "black",
			},
		},
		{
			name: "nil values",
			info: &teo.VideoTemplateInfo{},
			want: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapVideoTemplateInfoToSchema(tt.info)
			if got["video_codec"] != tt.want["video_codec"] {
				t.Errorf("video_codec mismatch, got %v, want %v", got["video_codec"], tt.want["video_codec"])
			}
			if got["fps"] != tt.want["fps"] {
				t.Errorf("fps mismatch, got %v, want %v", got["fps"], tt.want["fps"])
			}
			if got["bitrate"] != tt.want["bitrate"] {
				t.Errorf("bitrate mismatch, got %v, want %v", got["bitrate"], tt.want["bitrate"])
			}
		})
	}
}

func TestMapAudioTemplateInfoToSchema(t *testing.T) {
	tests := []struct {
		name string
		info *teo.AudioTemplateInfo
		want map[string]interface{}
	}{
		{
			name: "full configuration",
			info: &teo.AudioTemplateInfo{
				Codec:        helper.String("libfdk_aac"),
				AudioChannel: helper.Uint64(2),
			},
			want: map[string]interface{}{
				"codec":         "libfdk_aac",
				"audio_channel": 2,
			},
		},
		{
			name: "nil values",
			info: &teo.AudioTemplateInfo{},
			want: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapAudioTemplateInfoToSchema(tt.info)
			if got["codec"] != tt.want["codec"] {
				t.Errorf("codec mismatch, got %v, want %v", got["codec"], tt.want["codec"])
			}
			if got["audio_channel"] != tt.want["audio_channel"] {
				t.Errorf("audio_channel mismatch, got %v, want %v", got["audio_channel"], tt.want["audio_channel"])
			}
		})
	}
}

// Mock structs for testing
type tccommon struct {
	TencentV2Client *tencentcloudV2Client
}

type tencentcloudV2Client struct {
	teoClient interface{}
}
