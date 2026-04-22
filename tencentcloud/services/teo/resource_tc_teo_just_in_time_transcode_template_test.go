package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestJustInTimeTranscodeTemplate" -v -count=1 -gcflags="all=-l"

// TestJustInTimeTranscodeTemplate_Create_Success tests Create calls API and sets ID
func TestJustInTimeTranscodeTemplate_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateJustInTimeTranscodeTemplate", func(request *teov20220901.CreateJustInTimeTranscodeTemplateRequest) (*teov20220901.CreateJustInTimeTranscodeTemplateResponse, error) {
		resp := teov20220901.NewCreateJustInTimeTranscodeTemplateResponse()
		resp.Response = &teov20220901.CreateJustInTimeTranscodeTemplateResponseParams{
			TemplateId: ptrString("tpl-abcdefghij"),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeJustInTimeTranscodeTemplates", func(request *teov20220901.DescribeJustInTimeTranscodeTemplatesRequest) (*teov20220901.DescribeJustInTimeTranscodeTemplatesResponse, error) {
		resp := teov20220901.NewDescribeJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount: ptrUint64(1),
			TemplateSet: []*teov20220901.JustInTimeTranscodeTemplate{
				{
					TemplateId:        ptrString("tpl-abcdefghij"),
					TemplateName:      ptrString("test-template"),
					Comment:           ptrString("test comment"),
					VideoStreamSwitch: ptrString("on"),
					AudioStreamSwitch: ptrString("on"),
					Type:              ptrString("custom"),
					VideoTemplate: &teov20220901.VideoTemplateInfo{
						Codec:              ptrString("H.264"),
						Fps:                ptrFloat64(30),
						Bitrate:            ptrUint64(2000),
						ResolutionAdaptive: ptrString("open"),
						Width:              ptrUint64(1280),
						Height:             ptrUint64(720),
						FillType:           ptrString("black"),
					},
					AudioTemplate: &teov20220901.AudioTemplateInfo{
						Codec:        ptrString("libfdk_aac"),
						AudioChannel: ptrUint64(2),
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-1234567890",
		"template_name":       "test-template",
		"comment":             "test comment",
		"video_stream_switch": "on",
		"audio_stream_switch": "on",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#tpl-abcdefghij", d.Id())
}

// TestJustInTimeTranscodeTemplate_Create_APIError tests Create handles API error
func TestJustInTimeTranscodeTemplate_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateJustInTimeTranscodeTemplate", func(request *teov20220901.CreateJustInTimeTranscodeTemplateRequest) (*teov20220901.CreateJustInTimeTranscodeTemplateResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-invalid",
		"template_name": "test-template",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestJustInTimeTranscodeTemplate_Read_Success tests Read retrieves template data
func TestJustInTimeTranscodeTemplate_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeJustInTimeTranscodeTemplates", func(request *teov20220901.DescribeJustInTimeTranscodeTemplatesRequest) (*teov20220901.DescribeJustInTimeTranscodeTemplatesResponse, error) {
		resp := teov20220901.NewDescribeJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount: ptrUint64(1),
			TemplateSet: []*teov20220901.JustInTimeTranscodeTemplate{
				{
					TemplateId:        ptrString("tpl-abcdefghij"),
					TemplateName:      ptrString("test-template"),
					Comment:           ptrString("test comment"),
					VideoStreamSwitch: ptrString("on"),
					AudioStreamSwitch: ptrString("on"),
					Type:              ptrString("custom"),
					VideoTemplate: &teov20220901.VideoTemplateInfo{
						Codec:              ptrString("H.264"),
						Fps:                ptrFloat64(30),
						Bitrate:            ptrUint64(2000),
						ResolutionAdaptive: ptrString("open"),
						Width:              ptrUint64(1280),
						Height:             ptrUint64(720),
						FillType:           ptrString("black"),
					},
					AudioTemplate: &teov20220901.AudioTemplateInfo{
						Codec:        ptrString("libfdk_aac"),
						AudioChannel: ptrUint64(2),
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-1234567890",
		"template_name": "test-template",
		"template_id":   "tpl-abcdefghij",
	})
	d.SetId("zone-1234567890#tpl-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-template", d.Get("template_name"))
	assert.Equal(t, "test comment", d.Get("comment"))
	assert.Equal(t, "on", d.Get("video_stream_switch"))
	assert.Equal(t, "on", d.Get("audio_stream_switch"))
}

// TestJustInTimeTranscodeTemplate_Read_NotFound tests Read handles template not found
func TestJustInTimeTranscodeTemplate_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeJustInTimeTranscodeTemplates", func(request *teov20220901.DescribeJustInTimeTranscodeTemplatesRequest) (*teov20220901.DescribeJustInTimeTranscodeTemplatesResponse, error) {
		resp := teov20220901.NewDescribeJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount:  ptrUint64(0),
			TemplateSet: []*teov20220901.JustInTimeTranscodeTemplate{},
			RequestId:   ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-1234567890",
		"template_name": "test-template",
		"template_id":   "tpl-abcdefghij",
	})
	d.SetId("zone-1234567890#tpl-abcdefghij")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestJustInTimeTranscodeTemplate_Update tests Update is a no-op that calls Read
func TestJustInTimeTranscodeTemplate_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeJustInTimeTranscodeTemplates", func(request *teov20220901.DescribeJustInTimeTranscodeTemplatesRequest) (*teov20220901.DescribeJustInTimeTranscodeTemplatesResponse, error) {
		resp := teov20220901.NewDescribeJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount: ptrUint64(1),
			TemplateSet: []*teov20220901.JustInTimeTranscodeTemplate{
				{
					TemplateId:        ptrString("tpl-abcdefghij"),
					TemplateName:      ptrString("test-template"),
					Comment:           ptrString("test comment"),
					VideoStreamSwitch: ptrString("on"),
					AudioStreamSwitch: ptrString("on"),
					Type:              ptrString("custom"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-1234567890",
		"template_name": "test-template",
		"template_id":   "tpl-abcdefghij",
	})
	d.SetId("zone-1234567890#tpl-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestJustInTimeTranscodeTemplate_Delete_Success tests Delete removes template
func TestJustInTimeTranscodeTemplate_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteJustInTimeTranscodeTemplates", func(request *teov20220901.DeleteJustInTimeTranscodeTemplatesRequest) (*teov20220901.DeleteJustInTimeTranscodeTemplatesResponse, error) {
		resp := teov20220901.NewDeleteJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DeleteJustInTimeTranscodeTemplatesResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-1234567890",
		"template_name": "test-template",
		"template_id":   "tpl-abcdefghij",
	})
	d.SetId("zone-1234567890#tpl-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestJustInTimeTranscodeTemplate_Delete_APIError tests Delete handles API error
func TestJustInTimeTranscodeTemplate_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteJustInTimeTranscodeTemplates", func(request *teov20220901.DeleteJustInTimeTranscodeTemplatesRequest) (*teov20220901.DeleteJustInTimeTranscodeTemplatesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Template not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-1234567890",
		"template_name": "test-template",
		"template_id":   "tpl-abcdefghij",
	})
	d.SetId("zone-1234567890#tpl-abcdefghij")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestJustInTimeTranscodeTemplate_Schema validates schema definition
func TestJustInTimeTranscodeTemplate_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields with ForceNew
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "template_name")
	templateName := res.Schema["template_name"]
	assert.Equal(t, schema.TypeString, templateName.Type)
	assert.True(t, templateName.Required)
	assert.True(t, templateName.ForceNew)

	// Check optional fields WITHOUT ForceNew
	assert.Contains(t, res.Schema, "comment")
	comment := res.Schema["comment"]
	assert.Equal(t, schema.TypeString, comment.Type)
	assert.True(t, comment.Optional)
	assert.False(t, comment.ForceNew)

	assert.Contains(t, res.Schema, "video_stream_switch")
	videoSwitch := res.Schema["video_stream_switch"]
	assert.Equal(t, schema.TypeString, videoSwitch.Type)
	assert.True(t, videoSwitch.Optional)
	assert.False(t, videoSwitch.ForceNew)

	assert.Contains(t, res.Schema, "audio_stream_switch")
	audioSwitch := res.Schema["audio_stream_switch"]
	assert.Equal(t, schema.TypeString, audioSwitch.Type)
	assert.True(t, audioSwitch.Optional)
	assert.False(t, audioSwitch.ForceNew)

	assert.Contains(t, res.Schema, "video_template")
	videoTemplate := res.Schema["video_template"]
	assert.Equal(t, schema.TypeList, videoTemplate.Type)
	assert.True(t, videoTemplate.Optional)
	assert.False(t, videoTemplate.ForceNew)

	assert.Contains(t, res.Schema, "audio_template")
	audioTemplate := res.Schema["audio_template"]
	assert.Equal(t, schema.TypeList, audioTemplate.Type)
	assert.True(t, audioTemplate.Optional)
	assert.False(t, audioTemplate.ForceNew)

	// Check computed fields
	assert.Contains(t, res.Schema, "template_id")
	templateId := res.Schema["template_id"]
	assert.Equal(t, schema.TypeString, templateId.Type)
	assert.True(t, templateId.Computed)

	assert.Contains(t, res.Schema, "type")
	typeField := res.Schema["type"]
	assert.Equal(t, schema.TypeString, typeField.Type)
	assert.True(t, typeField.Computed)

	assert.Contains(t, res.Schema, "create_time")
	assert.Contains(t, res.Schema, "update_time")
}

func ptrFloat64(f float64) *float64 {
	return &f
}

func ptrUint64(u uint64) *uint64 {
	return &u
}
