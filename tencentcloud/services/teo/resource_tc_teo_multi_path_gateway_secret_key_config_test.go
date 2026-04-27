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

// mockMetaForSecretKey implements tccommon.ProviderMeta
type mockMetaForSecretKey struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForSecretKey) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForSecretKey{}

func newMockMetaForSecretKey() *mockMetaForSecretKey {
	return &mockMetaForSecretKey{client: &connectivity.TencentCloudClient{}}
}

func ptrStringForSecretKey(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestMultiPathGatewaySecretKeyConfig" -v -count=1 -gcflags="all=-l"

// TestMultiPathGatewaySecretKeyConfig_CreateSuccess tests Create sets zone_id as ID and calls CreateMultiPathGatewaySecretKey API
func TestMultiPathGatewaySecretKeyConfig_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch IsNewResource to return true during Create
	patches.ApplyMethodReturn(&schema.ResourceData{}, "IsNewResource", true)

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecretKey().client, "UseTeoV20220901Client", teoClient)

	// Patch CreateMultiPathGatewaySecretKeyWithContext (called by Update during Create)
	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewaySecretKeyWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewaySecretKeyRequest) (*teov20220901.CreateMultiPathGatewaySecretKeyResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "dGVzdC1zZWNyZXQta2V5", *request.SecretKey)
		resp := teov20220901.NewCreateMultiPathGatewaySecretKeyResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewaySecretKeyResponseParams{
			RequestId: ptrStringForSecretKey("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeMultiPathGatewaySecretKey (called by Read after Update)
	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewaySecretKey", func(request *teov20220901.DescribeMultiPathGatewaySecretKeyRequest) (*teov20220901.DescribeMultiPathGatewaySecretKeyResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		resp := teov20220901.NewDescribeMultiPathGatewaySecretKeyResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaySecretKeyResponseParams{
			SecretKey: ptrStringForSecretKey("dGVzdC1zZWNyZXQta2V5"),
			RequestId: ptrStringForSecretKey("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecretKey()
	res := teo.ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"secret_key": "dGVzdC1zZWNyZXQta2V5",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())
	assert.Equal(t, "dGVzdC1zZWNyZXQta2V5", d.Get("secret_key"))
}

// TestMultiPathGatewaySecretKeyConfig_CreateAPIError tests Create handles API error
func TestMultiPathGatewaySecretKeyConfig_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch IsNewResource to return true during Create
	patches.ApplyMethodReturn(&schema.ResourceData{}, "IsNewResource", true)

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecretKey().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewaySecretKeyWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewaySecretKeyRequest) (*teov20220901.CreateMultiPathGatewaySecretKeyResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForSecretKey()
	res := teo.ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-invalid",
		"secret_key": "dGVzdC1zZWNyZXQta2V5",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestMultiPathGatewaySecretKeyConfig_ReadSuccess tests Read populates state from API
func TestMultiPathGatewaySecretKeyConfig_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecretKey().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewaySecretKey", func(request *teov20220901.DescribeMultiPathGatewaySecretKeyRequest) (*teov20220901.DescribeMultiPathGatewaySecretKeyResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewaySecretKeyResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaySecretKeyResponseParams{
			SecretKey: ptrStringForSecretKey("cmVhZC1zZWNyZXQta2V5"),
			RequestId: ptrStringForSecretKey("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecretKey()
	res := teo.ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"secret_key": "cmVhZC1zZWNyZXQta2V5",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())
	assert.Equal(t, "cmVhZC1zZWNyZXQta2V5", d.Get("secret_key"))
	assert.Equal(t, "zone-12345678", d.Get("zone_id"))
}

// TestMultiPathGatewaySecretKeyConfig_ReadNotFound tests Read removes resource when not found
func TestMultiPathGatewaySecretKeyConfig_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecretKey().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewaySecretKey", func(request *teov20220901.DescribeMultiPathGatewaySecretKeyRequest) (*teov20220901.DescribeMultiPathGatewaySecretKeyResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaForSecretKey()
	res := teo.ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-notfound",
		"secret_key": "dGVzdC1zZWNyZXQta2V5",
	})
	d.SetId("zone-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestMultiPathGatewaySecretKeyConfig_UpdateSuccess tests Update calls Modify API
func TestMultiPathGatewaySecretKeyConfig_UpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch IsNewResource to return false during Update
	patches.ApplyMethodReturn(&schema.ResourceData{}, "IsNewResource", false)

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecretKey().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewaySecretKeyWithContext", func(_ context.Context, request *teov20220901.ModifyMultiPathGatewaySecretKeyRequest) (*teov20220901.ModifyMultiPathGatewaySecretKeyResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "dXBkYXRlZC1zZWNyZXQta2V5", *request.SecretKey)
		resp := teov20220901.NewModifyMultiPathGatewaySecretKeyResponse()
		resp.Response = &teov20220901.ModifyMultiPathGatewaySecretKeyResponseParams{
			RequestId: ptrStringForSecretKey("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewaySecretKey", func(request *teov20220901.DescribeMultiPathGatewaySecretKeyRequest) (*teov20220901.DescribeMultiPathGatewaySecretKeyResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewaySecretKeyResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaySecretKeyResponseParams{
			SecretKey: ptrStringForSecretKey("dXBkYXRlZC1zZWNyZXQta2V5"),
			RequestId: ptrStringForSecretKey("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecretKey()
	res := teo.ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"secret_key": "dXBkYXRlZC1zZWNyZXQta2V5",
	})
	d.SetId("zone-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "dXBkYXRlZC1zZWNyZXQta2V5", d.Get("secret_key"))
}

// TestMultiPathGatewaySecretKeyConfig_UpdateAPIError tests Update handles API error
func TestMultiPathGatewaySecretKeyConfig_UpdateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch IsNewResource to return false during Update
	patches.ApplyMethodReturn(&schema.ResourceData{}, "IsNewResource", false)

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecretKey().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewaySecretKeyWithContext", func(_ context.Context, request *teov20220901.ModifyMultiPathGatewaySecretKeyRequest) (*teov20220901.ModifyMultiPathGatewaySecretKeyResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid secret_key")
	})

	meta := newMockMetaForSecretKey()
	res := teo.ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"secret_key": "invalid-key",
	})
	d.SetId("zone-12345678")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestMultiPathGatewaySecretKeyConfig_Delete tests Delete is no-op
func TestMultiPathGatewaySecretKeyConfig_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"secret_key": "dGVzdC1zZWNyZXQta2V5",
	})
	d.SetId("zone-12345678")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestMultiPathGatewaySecretKeyConfig_Schema validates schema definition
func TestMultiPathGatewaySecretKeyConfig_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoMultiPathGatewaySecretKeyConfig()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "secret_key")

	// Check zone_id
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check secret_key
	secretKey := res.Schema["secret_key"]
	assert.Equal(t, schema.TypeString, secretKey.Type)
	assert.True(t, secretKey.Required)
	assert.True(t, secretKey.Sensitive)
}
