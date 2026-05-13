package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaForMultiPathGatewayLine implements tccommon.ProviderMeta
type mockMetaForMultiPathGatewayLine struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForMultiPathGatewayLine) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForMultiPathGatewayLine{}

func newMockMetaForMultiPathGatewayLine() *mockMetaForMultiPathGatewayLine {
	return &mockMetaForMultiPathGatewayLine{client: &connectivity.TencentCloudClient{}}
}

func ptrStrMPGL(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestMultiPathGatewayLine" -v -count=1 -gcflags="all=-l"

// TestMultiPathGatewayLine_Create_Success tests Create calls API and sets composite ID
func TestMultiPathGatewayLine_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGatewayLine().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "custom", *request.LineType)
		assert.Equal(t, "1.2.3.4:80", *request.LineAddress)
		resp := teov20220901.NewCreateMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayLineResponseParams{
			LineId:    ptrStrMPGL("line-1"),
			RequestId: ptrStrMPGL("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      ptrStrMPGL("line-1"),
				LineType:    ptrStrMPGL("custom"),
				LineAddress: ptrStrMPGL("1.2.3.4:80"),
			},
			RequestId: ptrStrMPGL("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGatewayLine()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:80",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#gw-abcdefgh#line-1", d.Id())
}

// TestMultiPathGatewayLine_Create_APIError tests Create handles API error
func TestMultiPathGatewayLine_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGatewayLine().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForMultiPathGatewayLine()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-invalid",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:80",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestMultiPathGatewayLine_Create_WithProxy tests Create with proxy type
func TestMultiPathGatewayLine_Create_WithProxy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGatewayLine().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "proxy", *request.LineType)
		assert.Equal(t, "sid-38hbn26osico", *request.ProxyId)
		assert.Equal(t, "rule-abcdef", *request.RuleId)
		resp := teov20220901.NewCreateMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayLineResponseParams{
			LineId:    ptrStrMPGL("line-2"),
			RequestId: ptrStrMPGL("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      ptrStrMPGL("line-2"),
				LineType:    ptrStrMPGL("proxy"),
				LineAddress: ptrStrMPGL("5.6.7.8:443"),
				ProxyId:     ptrStrMPGL("sid-38hbn26osico"),
				RuleId:      ptrStrMPGL("rule-abcdef"),
			},
			RequestId: ptrStrMPGL("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGatewayLine()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "proxy",
		"line_address": "5.6.7.8:443",
		"proxy_id":     "sid-38hbn26osico",
		"rule_id":      "rule-abcdef",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#gw-abcdefgh#line-2", d.Id())
}

// TestMultiPathGatewayLine_Read_Success tests Read populates state from API
func TestMultiPathGatewayLine_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGatewayLine().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "line-1", *request.LineId)
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      ptrStrMPGL("line-1"),
				LineType:    ptrStrMPGL("custom"),
				LineAddress: ptrStrMPGL("1.2.3.4:80"),
			},
			RequestId: ptrStrMPGL("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGatewayLine()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:80",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-1")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "custom", d.Get("line_type"))
	assert.Equal(t, "1.2.3.4:80", d.Get("line_address"))
	assert.Equal(t, "line-1", d.Get("line_id"))
}

// TestMultiPathGatewayLine_Read_NotFound tests Read handles resource not found
func TestMultiPathGatewayLine_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGatewayLine().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Line not found")
	})

	meta := newMockMetaForMultiPathGatewayLine()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:80",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-1")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestMultiPathGatewayLine_Update_Success tests Update calls Modify API
func TestMultiPathGatewayLine_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGatewayLine().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.ModifyMultiPathGatewayLineRequest) (*teov20220901.ModifyMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "line-1", *request.LineId)
		assert.Equal(t, "custom", *request.LineType)
		assert.Equal(t, "5.6.7.8:443", *request.LineAddress)
		resp := teov20220901.NewModifyMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.ModifyMultiPathGatewayLineResponseParams{
			RequestId: ptrStrMPGL("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      ptrStrMPGL("line-1"),
				LineType:    ptrStrMPGL("custom"),
				LineAddress: ptrStrMPGL("5.6.7.8:443"),
			},
			RequestId: ptrStrMPGL("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGatewayLine()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "5.6.7.8:443",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-1")

	// Simulate a change by setting old state differently
	_ = d.Set("line_address", "5.6.7.8:443")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestMultiPathGatewayLine_Delete_Success tests Delete calls API successfully
func TestMultiPathGatewayLine_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGatewayLine().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.DeleteMultiPathGatewayLineRequest) (*teov20220901.DeleteMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "line-1", *request.LineId)
		resp := teov20220901.NewDeleteMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DeleteMultiPathGatewayLineResponseParams{
			RequestId: ptrStrMPGL("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGatewayLine()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:80",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-1")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestMultiPathGatewayLine_Delete_APIError tests Delete handles API error
func TestMultiPathGatewayLine_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGatewayLine().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.DeleteMultiPathGatewayLineRequest) (*teov20220901.DeleteMultiPathGatewayLineResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Line not found")
	})

	meta := newMockMetaForMultiPathGatewayLine()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:80",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-1")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestMultiPathGatewayLine_Schema validates schema definition
func TestMultiPathGatewayLine_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "gateway_id")
	assert.Contains(t, res.Schema, "line_type")
	assert.Contains(t, res.Schema, "line_address")
	assert.Contains(t, res.Schema, "proxy_id")
	assert.Contains(t, res.Schema, "rule_id")
	assert.Contains(t, res.Schema, "line_id")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	gatewayId := res.Schema["gateway_id"]
	assert.Equal(t, schema.TypeString, gatewayId.Type)
	assert.True(t, gatewayId.Required)
	assert.True(t, gatewayId.ForceNew)

	lineType := res.Schema["line_type"]
	assert.Equal(t, schema.TypeString, lineType.Type)
	assert.True(t, lineType.Required)
	assert.False(t, lineType.ForceNew)

	lineAddress := res.Schema["line_address"]
	assert.Equal(t, schema.TypeString, lineAddress.Type)
	assert.True(t, lineAddress.Required)
	assert.False(t, lineAddress.ForceNew)

	proxyId := res.Schema["proxy_id"]
	assert.Equal(t, schema.TypeString, proxyId.Type)
	assert.True(t, proxyId.Optional)
	assert.False(t, proxyId.Required)

	ruleId := res.Schema["rule_id"]
	assert.Equal(t, schema.TypeString, ruleId.Type)
	assert.True(t, ruleId.Optional)
	assert.False(t, ruleId.Required)

	lineId := res.Schema["line_id"]
	assert.Equal(t, schema.TypeString, lineId.Type)
	assert.True(t, lineId.Computed)
	assert.False(t, lineId.Required)
	assert.False(t, lineId.Optional)
}
