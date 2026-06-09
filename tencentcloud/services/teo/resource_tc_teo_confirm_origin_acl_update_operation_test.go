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

// mockMetaConfirmOriginAclUpdate implements tccommon.ProviderMeta
type mockMetaConfirmOriginAclUpdate struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaConfirmOriginAclUpdate) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaConfirmOriginAclUpdate{}

func newMockMetaConfirmOriginAclUpdate() *mockMetaConfirmOriginAclUpdate {
	return &mockMetaConfirmOriginAclUpdate{client: &connectivity.TencentCloudClient{}}
}

func ptrStringConfirmOriginAclUpdate(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestConfirmOriginAclUpdateOperation" -v -count=1 -gcflags="all=-l"

// TestConfirmOriginAclUpdateOperation_Success tests successful confirmation
func TestConfirmOriginAclUpdateOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaConfirmOriginAclUpdate().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ConfirmOriginACLUpdateWithContext", func(_ context.Context, request *teov20220901.ConfirmOriginACLUpdateRequest) (*teov20220901.ConfirmOriginACLUpdateResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)

		resp := teov20220901.NewConfirmOriginACLUpdateResponse()
		resp.Response = &teov20220901.ConfirmOriginACLUpdateResponseParams{
			RequestId: ptrStringConfirmOriginAclUpdate("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaConfirmOriginAclUpdate()
	res := teo.ResourceTencentCloudTeoConfirmOriginAclUpdateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestConfirmOriginAclUpdateOperation_APIError tests API error handling
func TestConfirmOriginAclUpdateOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaConfirmOriginAclUpdate().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ConfirmOriginACLUpdateWithContext", func(_ context.Context, request *teov20220901.ConfirmOriginACLUpdateRequest) (*teov20220901.ConfirmOriginACLUpdateResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaConfirmOriginAclUpdate()
	res := teo.ResourceTencentCloudTeoConfirmOriginAclUpdateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestConfirmOriginAclUpdateOperation_Read tests Read is no-op
func TestConfirmOriginAclUpdateOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoConfirmOriginAclUpdateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})
	d.SetId("test-id")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestConfirmOriginAclUpdateOperation_Delete tests Delete is no-op
func TestConfirmOriginAclUpdateOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoConfirmOriginAclUpdateOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})
	d.SetId("test-id")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestConfirmOriginAclUpdateOperation_Schema validates schema definition
func TestConfirmOriginAclUpdateOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoConfirmOriginAclUpdateOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)
}
