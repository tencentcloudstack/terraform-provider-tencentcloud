package tse_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	tseService "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tse"
)

type mockMetaIpRestriction struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaIpRestriction) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaIpRestriction{}

func newMockMetaIpRestriction() *mockMetaIpRestriction {
	return &mockMetaIpRestriction{client: &connectivity.TencentCloudClient{}}
}

func ptrStringIpR(s string) *string {
	return &s
}

func ptrBoolIpR(b bool) *bool {
	return &b
}

func TestTseCngwIPRestriction_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tseClient := &tse.Client{}
	patches.ApplyMethodReturn(newMockMetaIpRestriction().client, "UseTseClient", tseClient)

	patches.ApplyMethodFunc(tseClient, "CreateOrModifyCloudNativeAPIGatewayIPRestrictionWithContext", func(_ context.Context, request *tse.CreateOrModifyCloudNativeAPIGatewayIPRestrictionRequest) (*tse.CreateOrModifyCloudNativeAPIGatewayIPRestrictionResponse, error) {
		assert.NotNil(t, request.GatewayId)
		assert.Equal(t, "gw-1", *request.GatewayId)
		assert.NotNil(t, request.SourceType)
		assert.Equal(t, "route", *request.SourceType)
		assert.NotNil(t, request.SourceId)
		assert.Equal(t, "r-1", *request.SourceId)
		assert.NotNil(t, request.Enabled)
		assert.Equal(t, true, *request.Enabled)
		assert.NotNil(t, request.RestrictionType)
		assert.Equal(t, "whiteList", *request.RestrictionType)
		assert.NotNil(t, request.AddressList)
		assert.Equal(t, 2, len(request.AddressList))

		resp := tse.NewCreateOrModifyCloudNativeAPIGatewayIPRestrictionResponse()
		resp.Response = &tse.CreateOrModifyCloudNativeAPIGatewayIPRestrictionResponseParams{
			RequestId: ptrStringIpR("fake-request-id-create"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tseClient, "DescribeCloudNativeAPIGatewayIPRestrictionWithContext", func(_ context.Context, request *tse.DescribeCloudNativeAPIGatewayIPRestrictionRequest) (*tse.DescribeCloudNativeAPIGatewayIPRestrictionResponse, error) {
		assert.NotNil(t, request.GatewayId)
		assert.Equal(t, "gw-1", *request.GatewayId)
		assert.NotNil(t, request.SourceType)
		assert.Equal(t, "route", *request.SourceType)
		assert.NotNil(t, request.SourceId)
		assert.Equal(t, "r-1", *request.SourceId)

		resp := tse.NewDescribeCloudNativeAPIGatewayIPRestrictionResponse()
		resp.Response = &tse.DescribeCloudNativeAPIGatewayIPRestrictionResponseParams{
			Result: &tse.DescribeKongIpRestrictionResult{
				SourceType:      ptrStringIpR("route"),
				SourceId:        ptrStringIpR("r-1"),
				Enabled:         ptrBoolIpR(true),
				RestrictionType: ptrStringIpR("whiteList"),
				AddressList:     []*string{ptrStringIpR("10.0.0.0/8"), ptrStringIpR("192.168.1.1")},
			},
			RequestId: ptrStringIpR("fake-request-id-describe"),
		}
		return resp, nil
	})

	meta := newMockMetaIpRestriction()
	res := tseService.ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"gateway_id":       "gw-1",
		"source_type":      "route",
		"source_id":        "r-1",
		"enabled":          true,
		"restriction_type": "whiteList",
		"address_list":     []interface{}{"10.0.0.0/8", "192.168.1.1"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "gw-1#route#r-1", d.Id())
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, "whiteList", d.Get("restriction_type"))
	assert.Equal(t, "gw-1", d.Get("gateway_id"))
	assert.Equal(t, "route", d.Get("source_type"))
	assert.Equal(t, "r-1", d.Get("source_id"))
}

func TestTseCngwIPRestriction_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tseClient := &tse.Client{}
	patches.ApplyMethodReturn(newMockMetaIpRestriction().client, "UseTseClient", tseClient)

	patches.ApplyMethodFunc(tseClient, "DescribeCloudNativeAPIGatewayIPRestrictionWithContext", func(_ context.Context, request *tse.DescribeCloudNativeAPIGatewayIPRestrictionRequest) (*tse.DescribeCloudNativeAPIGatewayIPRestrictionResponse, error) {
		resp := tse.NewDescribeCloudNativeAPIGatewayIPRestrictionResponse()
		resp.Response = &tse.DescribeCloudNativeAPIGatewayIPRestrictionResponseParams{
			Result: &tse.DescribeKongIpRestrictionResult{
				SourceType:      ptrStringIpR("route"),
				SourceId:        ptrStringIpR("r-1"),
				Enabled:         ptrBoolIpR(true),
				RestrictionType: ptrStringIpR("whiteList"),
				AddressList:     []*string{ptrStringIpR("10.0.0.0/8")},
			},
			RequestId: ptrStringIpR("fake-request-id-describe"),
		}
		return resp, nil
	})

	meta := newMockMetaIpRestriction()
	res := tseService.ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"gateway_id":       "gw-1",
		"source_type":      "route",
		"source_id":        "r-1",
		"enabled":          true,
		"restriction_type": "whiteList",
		"address_list":     []interface{}{"10.0.0.0/8"},
	})
	d.SetId("gw-1#route#r-1")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "gw-1#route#r-1", d.Id())
	assert.Equal(t, true, d.Get("enabled"))
	assert.Equal(t, "whiteList", d.Get("restriction_type"))
	assert.Equal(t, "gw-1", d.Get("gateway_id"))
	assert.Equal(t, "route", d.Get("source_type"))
	assert.Equal(t, "r-1", d.Get("source_id"))
}

func TestTseCngwIPRestriction_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tseClient := &tse.Client{}
	patches.ApplyMethodReturn(newMockMetaIpRestriction().client, "UseTseClient", tseClient)

	patches.ApplyMethodFunc(tseClient, "DescribeCloudNativeAPIGatewayIPRestrictionWithContext", func(_ context.Context, request *tse.DescribeCloudNativeAPIGatewayIPRestrictionRequest) (*tse.DescribeCloudNativeAPIGatewayIPRestrictionResponse, error) {
		resp := tse.NewDescribeCloudNativeAPIGatewayIPRestrictionResponse()
		resp.Response = &tse.DescribeCloudNativeAPIGatewayIPRestrictionResponseParams{
			Result:    nil,
			RequestId: ptrStringIpR("fake-request-id-not-found"),
		}
		return resp, nil
	})

	meta := newMockMetaIpRestriction()
	res := tseService.ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"gateway_id":       "gw-1",
		"source_type":      "route",
		"source_id":        "r-1",
		"enabled":          true,
		"restriction_type": "whiteList",
		"address_list":     []interface{}{"10.0.0.0/8"},
	})
	d.SetId("gw-1#route#r-1")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTseCngwIPRestriction_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tseClient := &tse.Client{}
	patches.ApplyMethodReturn(newMockMetaIpRestriction().client, "UseTseClient", tseClient)

	patches.ApplyMethodFunc(tseClient, "CreateOrModifyCloudNativeAPIGatewayIPRestrictionWithContext", func(_ context.Context, request *tse.CreateOrModifyCloudNativeAPIGatewayIPRestrictionRequest) (*tse.CreateOrModifyCloudNativeAPIGatewayIPRestrictionResponse, error) {
		assert.NotNil(t, request.GatewayId)
		assert.Equal(t, "gw-1", *request.GatewayId)
		assert.NotNil(t, request.SourceType)
		assert.Equal(t, "route", *request.SourceType)
		assert.NotNil(t, request.SourceId)
		assert.Equal(t, "r-1", *request.SourceId)
		assert.NotNil(t, request.Enabled)
		assert.Equal(t, true, *request.Enabled)
		assert.NotNil(t, request.RestrictionType)
		assert.Equal(t, "blackList", *request.RestrictionType)
		assert.NotNil(t, request.AddressList)
		assert.Equal(t, 1, len(request.AddressList))

		resp := tse.NewCreateOrModifyCloudNativeAPIGatewayIPRestrictionResponse()
		resp.Response = &tse.CreateOrModifyCloudNativeAPIGatewayIPRestrictionResponseParams{
			RequestId: ptrStringIpR("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(tseClient, "DescribeCloudNativeAPIGatewayIPRestrictionWithContext", func(_ context.Context, request *tse.DescribeCloudNativeAPIGatewayIPRestrictionRequest) (*tse.DescribeCloudNativeAPIGatewayIPRestrictionResponse, error) {
		resp := tse.NewDescribeCloudNativeAPIGatewayIPRestrictionResponse()
		resp.Response = &tse.DescribeCloudNativeAPIGatewayIPRestrictionResponseParams{
			Result: &tse.DescribeKongIpRestrictionResult{
				SourceType:      ptrStringIpR("route"),
				SourceId:        ptrStringIpR("r-1"),
				Enabled:         ptrBoolIpR(true),
				RestrictionType: ptrStringIpR("blackList"),
				AddressList:     []*string{ptrStringIpR("2.2.2.2")},
			},
			RequestId: ptrStringIpR("fake-request-id-describe"),
		}
		return resp, nil
	})

	meta := newMockMetaIpRestriction()
	res := tseService.ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"gateway_id":       "gw-1",
		"source_type":      "route",
		"source_id":        "r-1",
		"enabled":          true,
		"restriction_type": "blackList",
		"address_list":     []interface{}{"2.2.2.2"},
	})
	d.SetId("gw-1#route#r-1")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestTseCngwIPRestriction_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	tseClient := &tse.Client{}
	patches.ApplyMethodReturn(newMockMetaIpRestriction().client, "UseTseClient", tseClient)

	patches.ApplyMethodFunc(tseClient, "DeleteCloudNativeAPIGatewayIPRestrictionWithContext", func(_ context.Context, request *tse.DeleteCloudNativeAPIGatewayIPRestrictionRequest) (*tse.DeleteCloudNativeAPIGatewayIPRestrictionResponse, error) {
		assert.NotNil(t, request.GatewayId)
		assert.Equal(t, "gw-1", *request.GatewayId)
		assert.NotNil(t, request.SourceType)
		assert.Equal(t, "route", *request.SourceType)
		assert.NotNil(t, request.SourceId)
		assert.Equal(t, "r-1", *request.SourceId)

		resp := tse.NewDeleteCloudNativeAPIGatewayIPRestrictionResponse()
		resp.Response = &tse.DeleteCloudNativeAPIGatewayIPRestrictionResponseParams{
			RequestId: ptrStringIpR("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaIpRestriction()
	res := tseService.ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"gateway_id":       "gw-1",
		"source_type":      "route",
		"source_id":        "r-1",
		"enabled":          true,
		"restriction_type": "whiteList",
		"address_list":     []interface{}{"10.0.0.0/8"},
	})
	d.SetId("gw-1#route#r-1")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
