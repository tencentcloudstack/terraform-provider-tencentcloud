package vod_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcvod "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vod"
)

type mockMetaVodAigcQuota struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaVodAigcQuota) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaVodAigcQuota{}

func newMockMetaVodAigcQuota() *mockMetaVodAigcQuota {
	return &mockMetaVodAigcQuota{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringVod(s string) *string { return &s }
func ptrUint64Vod(u uint64) *uint64 { return &u }

// go test ./tencentcloud/services/vod/ -run "TestUnitVodAigcQuota" -v -count=1 -gcflags="all=-l"

func TestUnitVodAigcQuotaCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaVodAigcQuota()
	vodClient := &vod.Client{}
	patches.ApplyMethodReturn(meta.client, "UseVodClient", vodClient)

	createCalled := false
	patches.ApplyMethodFunc(vodClient, "CreateAigcQuotaWithContext", func(ctx context.Context, request *vod.CreateAigcQuotaRequest) (*vod.CreateAigcQuotaResponse, error) {
		createCalled = true
		assert.Equal(t, uint64(251006666), *request.SubAppId)
		assert.Equal(t, "Image", *request.QuotaType)
		assert.Equal(t, uint64(100), *request.QuotaLimit)
		resp := &vod.CreateAigcQuotaResponse{}
		resp.Response = &vod.CreateAigcQuotaResponseParams{
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	describeCalled := false
	patches.ApplyMethodFunc(vodClient, "DescribeAigcQuotasWithContext", func(ctx context.Context, request *vod.DescribeAigcQuotasRequest) (*vod.DescribeAigcQuotasResponse, error) {
		describeCalled = true
		assert.Equal(t, uint64(251006666), *request.SubAppId)
		assert.Equal(t, "Image", *request.QuotaType)
		resp := &vod.DescribeAigcQuotasResponse{}
		resp.Response = &vod.DescribeAigcQuotasResponseParams{
			QuotaSet: []*vod.AigcQuotaItem{
				{
					QuotaType:  ptrStringVod("Image"),
					QuotaLimit: ptrUint64Vod(100),
					Usage:      ptrUint64Vod(0),
				},
			},
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	res := svcvod.ResourceTencentCloudVodAigcQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"sub_app_id":  251006666,
		"quota_type":  "Image",
		"quota_limit": 100,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.True(t, createCalled)
	assert.True(t, describeCalled)
	assert.Equal(t, "251006666#Image", d.Id())
}

func TestUnitVodAigcQuotaCreateText(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaVodAigcQuota()
	vodClient := &vod.Client{}
	patches.ApplyMethodReturn(meta.client, "UseVodClient", vodClient)

	createCalled := false
	patches.ApplyMethodFunc(vodClient, "CreateAigcQuotaWithContext", func(ctx context.Context, request *vod.CreateAigcQuotaRequest) (*vod.CreateAigcQuotaResponse, error) {
		createCalled = true
		assert.Equal(t, uint64(251006666), *request.SubAppId)
		assert.Equal(t, "Text", *request.QuotaType)
		assert.Equal(t, uint64(5000), *request.QuotaLimit)
		assert.Equal(t, "my-token", *request.ApiToken)
		resp := &vod.CreateAigcQuotaResponse{}
		resp.Response = &vod.CreateAigcQuotaResponseParams{
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vodClient, "DescribeAigcQuotasWithContext", func(ctx context.Context, request *vod.DescribeAigcQuotasRequest) (*vod.DescribeAigcQuotasResponse, error) {
		assert.Equal(t, uint64(251006666), *request.SubAppId)
		assert.Equal(t, "Text", *request.QuotaType)
		assert.Equal(t, "my-token", *request.ApiToken)
		resp := &vod.DescribeAigcQuotasResponse{}
		resp.Response = &vod.DescribeAigcQuotasResponseParams{
			QuotaSet: []*vod.AigcQuotaItem{
				{
					QuotaType:  ptrStringVod("Text"),
					ApiToken:   ptrStringVod("my-token"),
					QuotaLimit: ptrUint64Vod(5000),
					Usage:      ptrUint64Vod(100),
				},
			},
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	res := svcvod.ResourceTencentCloudVodAigcQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"sub_app_id":  251006666,
		"quota_type":  "Text",
		"quota_limit": 5000,
		"api_token":   "my-token",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.True(t, createCalled)
	assert.Equal(t, "251006666#Text#my-token", d.Id())
}

func TestUnitVodAigcQuotaRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaVodAigcQuota()
	vodClient := &vod.Client{}
	patches.ApplyMethodReturn(meta.client, "UseVodClient", vodClient)

	patches.ApplyMethodFunc(vodClient, "DescribeAigcQuotasWithContext", func(ctx context.Context, request *vod.DescribeAigcQuotasRequest) (*vod.DescribeAigcQuotasResponse, error) {
		assert.Equal(t, uint64(251006666), *request.SubAppId)
		assert.Equal(t, "Image", *request.QuotaType)
		resp := &vod.DescribeAigcQuotasResponse{}
		resp.Response = &vod.DescribeAigcQuotasResponseParams{
			QuotaSet: []*vod.AigcQuotaItem{
				{
					QuotaType:  ptrStringVod("Image"),
					QuotaLimit: ptrUint64Vod(200),
					Usage:      ptrUint64Vod(50),
				},
			},
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	res := svcvod.ResourceTencentCloudVodAigcQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"sub_app_id":  251006666,
		"quota_type":  "Image",
		"quota_limit": 200,
	})
	d.SetId("251006666#Image")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 251006666, d.Get("sub_app_id"))
	assert.Equal(t, "Image", d.Get("quota_type"))
	assert.Equal(t, 200, d.Get("quota_limit"))
	assert.Equal(t, 50, d.Get("usage"))
}

func TestUnitVodAigcQuotaReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaVodAigcQuota()
	vodClient := &vod.Client{}
	patches.ApplyMethodReturn(meta.client, "UseVodClient", vodClient)

	patches.ApplyMethodFunc(vodClient, "DescribeAigcQuotasWithContext", func(ctx context.Context, request *vod.DescribeAigcQuotasRequest) (*vod.DescribeAigcQuotasResponse, error) {
		resp := &vod.DescribeAigcQuotasResponse{}
		resp.Response = &vod.DescribeAigcQuotasResponseParams{
			QuotaSet:  []*vod.AigcQuotaItem{},
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	res := svcvod.ResourceTencentCloudVodAigcQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"sub_app_id":  251006666,
		"quota_type":  "Image",
		"quota_limit": 200,
	})
	d.SetId("251006666#Image")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestUnitVodAigcQuotaUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaVodAigcQuota()
	vodClient := &vod.Client{}
	patches.ApplyMethodReturn(meta.client, "UseVodClient", vodClient)

	modifyCalled := false
	patches.ApplyMethodFunc(vodClient, "ModifyAigcQuotaWithContext", func(ctx context.Context, request *vod.ModifyAigcQuotaRequest) (*vod.ModifyAigcQuotaResponse, error) {
		modifyCalled = true
		assert.Equal(t, uint64(251006666), *request.SubAppId)
		assert.Equal(t, "Image", *request.QuotaType)
		assert.Equal(t, uint64(300), *request.QuotaLimit)
		resp := &vod.ModifyAigcQuotaResponse{}
		resp.Response = &vod.ModifyAigcQuotaResponseParams{
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vodClient, "DescribeAigcQuotasWithContext", func(ctx context.Context, request *vod.DescribeAigcQuotasRequest) (*vod.DescribeAigcQuotasResponse, error) {
		resp := &vod.DescribeAigcQuotasResponse{}
		resp.Response = &vod.DescribeAigcQuotasResponseParams{
			QuotaSet: []*vod.AigcQuotaItem{
				{
					QuotaType:  ptrStringVod("Image"),
					QuotaLimit: ptrUint64Vod(300),
					Usage:      ptrUint64Vod(50),
				},
			},
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	res := svcvod.ResourceTencentCloudVodAigcQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"sub_app_id":  251006666,
		"quota_type":  "Image",
		"quota_limit": 300,
	})
	d.SetId("251006666#Image")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, 300, d.Get("quota_limit"))
}

func TestUnitVodAigcQuotaDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaVodAigcQuota()
	vodClient := &vod.Client{}
	patches.ApplyMethodReturn(meta.client, "UseVodClient", vodClient)

	deleteCalled := false
	patches.ApplyMethodFunc(vodClient, "DeleteAigcQuotaWithContext", func(ctx context.Context, request *vod.DeleteAigcQuotaRequest) (*vod.DeleteAigcQuotaResponse, error) {
		deleteCalled = true
		assert.Equal(t, uint64(251006666), *request.SubAppId)
		assert.Equal(t, "Image", *request.QuotaType)
		resp := &vod.DeleteAigcQuotaResponse{}
		resp.Response = &vod.DeleteAigcQuotaResponseParams{
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	// After delete, describe returns empty
	patches.ApplyMethodFunc(vodClient, "DescribeAigcQuotasWithContext", func(ctx context.Context, request *vod.DescribeAigcQuotasRequest) (*vod.DescribeAigcQuotasResponse, error) {
		resp := &vod.DescribeAigcQuotasResponse{}
		resp.Response = &vod.DescribeAigcQuotasResponseParams{
			QuotaSet:  []*vod.AigcQuotaItem{},
			RequestId: ptrStringVod("fake-request-id"),
		}
		return resp, nil
	})

	res := svcvod.ResourceTencentCloudVodAigcQuota()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"sub_app_id":  251006666,
		"quota_type":  "Image",
		"quota_limit": 100,
	})
	d.SetId("251006666#Image")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}

func TestUnitVodAigcQuotaParseId(t *testing.T) {
	// two-part id (Image/Video, no api_token)
	subAppId, quotaType, apiToken, err := svcvod.ParseVodAigcQuotaId("251006666#Image")
	assert.NoError(t, err)
	assert.Equal(t, uint64(251006666), subAppId)
	assert.Equal(t, "Image", quotaType)
	assert.Equal(t, "", apiToken)

	// two-part id (Video)
	subAppId, quotaType, apiToken, err = svcvod.ParseVodAigcQuotaId("251006666#Video")
	assert.NoError(t, err)
	assert.Equal(t, uint64(251006666), subAppId)
	assert.Equal(t, "Video", quotaType)
	assert.Equal(t, "", apiToken)

	// three-part id with empty api_token (Image)
	subAppId, quotaType, apiToken, err = svcvod.ParseVodAigcQuotaId("251006666#Image#")
	assert.NoError(t, err)
	assert.Equal(t, uint64(251006666), subAppId)
	assert.Equal(t, "Image", quotaType)
	assert.Equal(t, "", apiToken)

	// three-part id with api_token (Text)
	subAppId, quotaType, apiToken, err = svcvod.ParseVodAigcQuotaId("251006666#Text#my-token")
	assert.NoError(t, err)
	assert.Equal(t, uint64(251006666), subAppId)
	assert.Equal(t, "Text", quotaType)
	assert.Equal(t, "my-token", apiToken)

	// Invalid: single part
	_, _, _, err = svcvod.ParseVodAigcQuotaId("251006666")
	assert.Error(t, err)

	// Invalid: too many parts
	_, _, _, err = svcvod.ParseVodAigcQuotaId("251006666#Image#token#extra")
	assert.Error(t, err)

	// Invalid: non-numeric sub_app_id (2-part)
	_, _, _, err = svcvod.ParseVodAigcQuotaId("abc#Image")
	assert.Error(t, err)

	// Invalid: non-numeric sub_app_id (3-part)
	_, _, _, err = svcvod.ParseVodAigcQuotaId("abc#Image#")
	assert.Error(t, err)

	// Invalid: empty sub_app_id
	_, _, _, err = svcvod.ParseVodAigcQuotaId("#Image")
	assert.Error(t, err)
}
