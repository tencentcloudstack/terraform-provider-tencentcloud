package tag_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

type mockMetaTagAttachmentV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTagAttachmentV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTagAttachmentV2{}

func newMockMetaTagAttachmentV2() *mockMetaTagAttachmentV2 {
	return &mockMetaTagAttachmentV2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTRT(s string) *string {
	return &s
}

const (
	testTagKeyTRT      = "env"
	testTagValueTRT    = "prod"
	testTagValueNewTRT = "staging"
	testResourceSixTRT = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
	testCompositeIDTRT = "env#qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
)

// go test ./tencentcloud/services/tag/ -run "TestTagAttachmentV2" -v -count=1 -gcflags="all=-l"

// TestTagAttachmentV2_Create_Success tests Create with required parameters.
func TestTagAttachmentV2_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	// mock AddResourceTagWithContext
	patches.ApplyMethodFunc(tagClient, "AddResourceTagWithContext", func(_ context.Context, request *tag.AddResourceTagRequest) (*tag.AddResourceTagResponse, error) {
		assert.NotNil(t, request.TagKey)
		assert.Equal(t, testTagKeyTRT, *request.TagKey)
		assert.NotNil(t, request.TagValue)
		assert.Equal(t, testTagValueTRT, *request.TagValue)
		assert.NotNil(t, request.Resource)
		assert.Equal(t, testResourceSixTRT, *request.Resource)

		resp := tag.NewAddResourceTagResponse()
		resp.Response = &tag.AddResourceTagResponseParams{
			RequestId: ptrStringTRT("fake-request-id-create"),
		}
		return resp, nil
	})

	// mock GetResources (called by Read after Create)
	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		assert.NotNil(t, request.ResourceList)
		assert.Equal(t, testResourceSixTRT, *request.ResourceList[0])

		resp := tag.NewGetResourcesResponse()
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{
				{
					Resource: ptrStringTRT(testResourceSixTRT),
					Tags: []*tag.Tag{
						{
							TagKey:   ptrStringTRT(testTagKeyTRT),
							TagValue: ptrStringTRT(testTagValueTRT),
							Category: ptrStringTRT("Custom"),
						},
					},
				},
			},
			RequestId: ptrStringTRT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueTRT,
		"resource":  testResourceSixTRT,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, testCompositeIDTRT, d.Id())
	assert.Equal(t, testTagKeyTRT, d.Get("tag_key"))
	assert.Equal(t, testTagValueTRT, d.Get("tag_value"))
	assert.Equal(t, testResourceSixTRT, d.Get("resource"))
}

// TestTagAttachmentV2_Create_EmptyResponse tests Create when API returns empty response.
func TestTagAttachmentV2_Create_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	// mock AddResourceTagWithContext returns nil Response
	patches.ApplyMethodFunc(tagClient, "AddResourceTagWithContext", func(_ context.Context, request *tag.AddResourceTagRequest) (*tag.AddResourceTagResponse, error) {
		resp := tag.NewAddResourceTagResponse()
		// Response is nil
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueTRT,
		"resource":  testResourceSixTRT,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Response is nil")
	assert.Empty(t, d.Id())
}

// TestTagAttachmentV2_Create_APIError tests Create when API returns an error.
func TestTagAttachmentV2_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "AddResourceTagWithContext", func(_ context.Context, request *tag.AddResourceTagRequest) (*tag.AddResourceTagResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceInUse.TagDuplicate, Message=tag already bound")
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueTRT,
		"resource":  testResourceSixTRT,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TagDuplicate")
}

// TestTagAttachmentV2_Read_Found tests Read when the binding exists.
func TestTagAttachmentV2_Read_Found(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		resp := tag.NewGetResourcesResponse()
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{
				{
					Resource: ptrStringTRT(testResourceSixTRT),
					Tags: []*tag.Tag{
						{
							TagKey:   ptrStringTRT(testTagKeyTRT),
							TagValue: ptrStringTRT(testTagValueTRT),
							Category: ptrStringTRT("Custom"),
						},
					},
				},
			},
			RequestId: ptrStringTRT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueTRT,
		"resource":  testResourceSixTRT,
	})
	d.SetId(testCompositeIDTRT)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, testCompositeIDTRT, d.Id())
	assert.Equal(t, testTagKeyTRT, d.Get("tag_key"))
	assert.Equal(t, testTagValueTRT, d.Get("tag_value"))
	assert.Equal(t, testResourceSixTRT, d.Get("resource"))
}

// TestTagAttachmentV2_Read_NotFound tests Read when the binding does not exist.
func TestTagAttachmentV2_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		resp := tag.NewGetResourcesResponse()
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{},
			RequestId:              ptrStringTRT("fake-request-id-read-empty"),
		}
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueTRT,
		"resource":  testResourceSixTRT,
	})
	d.SetId(testCompositeIDTRT)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

// TestTagAttachmentV2_Read_NilTagValue tests Read when the matched row has nil TagValue.
func TestTagAttachmentV2_Read_NilTagValue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		resp := tag.NewGetResourcesResponse()
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{
				{
					Resource: ptrStringTRT(testResourceSixTRT),
					Tags: []*tag.Tag{
						{
							TagKey:   ptrStringTRT(testTagKeyTRT),
							TagValue: nil,
							Category: ptrStringTRT("Custom"),
						},
					},
				},
			},
			RequestId: ptrStringTRT("fake-request-id-read-nil"),
		}
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueTRT,
		"resource":  testResourceSixTRT,
	})
	d.SetId(testCompositeIDTRT)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, testCompositeIDTRT, d.Id())
	assert.Equal(t, testTagKeyTRT, d.Get("tag_key"))
	// tag_value should not be overwritten (nil check skips Set)
	assert.Equal(t, testTagValueTRT, d.Get("tag_value"))
}

// TestTagAttachmentV2_Update_TagValueChange tests Update when tag_value changes.
func TestTagAttachmentV2_Update_TagValueChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	// mock UpdateResourceTagValueWithContext
	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValueWithContext", func(_ context.Context, request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		assert.NotNil(t, request.TagKey)
		assert.Equal(t, testTagKeyTRT, *request.TagKey)
		assert.NotNil(t, request.TagValue)
		assert.Equal(t, testTagValueNewTRT, *request.TagValue)
		assert.NotNil(t, request.Resource)
		assert.Equal(t, testResourceSixTRT, *request.Resource)

		resp := tag.NewUpdateResourceTagValueResponse()
		resp.Response = &tag.UpdateResourceTagValueResponseParams{
			RequestId: ptrStringTRT("fake-request-id-update"),
		}
		return resp, nil
	})

	// mock GetResources (called by Read after Update)
	patches.ApplyMethodFunc(tagClient, "GetResources", func(request *tag.GetResourcesRequest) (*tag.GetResourcesResponse, error) {
		resp := tag.NewGetResourcesResponse()
		resp.Response = &tag.GetResourcesResponseParams{
			ResourceTagMappingList: []*tag.ResourceTagMapping{
				{
					Resource: ptrStringTRT(testResourceSixTRT),
					Tags: []*tag.Tag{
						{
							TagKey:   ptrStringTRT(testTagKeyTRT),
							TagValue: ptrStringTRT(testTagValueNewTRT),
							Category: ptrStringTRT("Custom"),
						},
					},
				},
			},
			RequestId: ptrStringTRT("fake-request-id-read-after-update"),
		}
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueNewTRT,
		"resource":  testResourceSixTRT,
	})
	d.SetId(testCompositeIDTRT)

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, testCompositeIDTRT, d.Id())
	assert.Equal(t, testTagValueNewTRT, d.Get("tag_value"))
}

// TestTagAttachmentV2_Update_EmptyResponse tests Update when API returns empty response.
func TestTagAttachmentV2_Update_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "UpdateResourceTagValueWithContext", func(_ context.Context, request *tag.UpdateResourceTagValueRequest) (*tag.UpdateResourceTagValueResponse, error) {
		resp := tag.NewUpdateResourceTagValueResponse()
		// Response is nil
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueNewTRT,
		"resource":  testResourceSixTRT,
	})
	d.SetId(testCompositeIDTRT)

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Response is nil")
}

// TestTagAttachmentV2_Delete_Success tests Delete success.
func TestTagAttachmentV2_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "DeleteResourceTagWithContext", func(_ context.Context, request *tag.DeleteResourceTagRequest) (*tag.DeleteResourceTagResponse, error) {
		assert.NotNil(t, request.TagKey)
		assert.Equal(t, testTagKeyTRT, *request.TagKey)
		assert.NotNil(t, request.Resource)
		assert.Equal(t, testResourceSixTRT, *request.Resource)

		resp := tag.NewDeleteResourceTagResponse()
		resp.Response = &tag.DeleteResourceTagResponseParams{
			RequestId: ptrStringTRT("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueTRT,
		"resource":  testResourceSixTRT,
	})
	d.SetId(testCompositeIDTRT)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTagAttachmentV2_Delete_EmptyResponse tests Delete when API returns empty response.
func TestTagAttachmentV2_Delete_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaTagAttachmentV2().client, "UseTagClient", tagClient)

	patches.ApplyMethodFunc(tagClient, "DeleteResourceTagWithContext", func(_ context.Context, request *tag.DeleteResourceTagRequest) (*tag.DeleteResourceTagResponse, error) {
		resp := tag.NewDeleteResourceTagResponse()
		// Response is nil
		return resp, nil
	})

	meta := newMockMetaTagAttachmentV2()
	res := svctag.ResourceTencentCloudTagAttachmentV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   testTagKeyTRT,
		"tag_value": testTagValueTRT,
		"resource":  testResourceSixTRT,
	})
	d.SetId(testCompositeIDTRT)

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Response is nil")
}

// TestTagAttachmentV2_Schema validates the schema definition.
func TestTagAttachmentV2_Schema(t *testing.T) {
	res := svctag.ResourceTencentCloudTagAttachmentV2()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "tag_key")
	assert.Contains(t, res.Schema, "tag_value")
	assert.Contains(t, res.Schema, "resource")

	tagKey := res.Schema["tag_key"]
	assert.Equal(t, schema.TypeString, tagKey.Type)
	assert.True(t, tagKey.Required)
	assert.True(t, tagKey.ForceNew)

	tagValue := res.Schema["tag_value"]
	assert.Equal(t, schema.TypeString, tagValue.Type)
	assert.True(t, tagValue.Required)
	assert.False(t, tagValue.ForceNew)

	resource := res.Schema["resource"]
	assert.Equal(t, schema.TypeString, resource.Type)
	assert.True(t, resource.Required)
	assert.True(t, resource.ForceNew)
}
