package teo_test

import (
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

// mockMeta implements tccommon.ProviderMeta
type mockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMeta{}

func newMockMeta() *mockMeta {
	return &mockMeta{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/teo/ -run "TestCreateCLSIndexOperation" -v -count=1 -gcflags="all=-l"
// TestCreateCLSIndexOperation_Success tests Create calls API and sets ID
func TestCreateCLSIndexOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch CreateCLSIndex to return success
	patches.ApplyMethodFunc(teoClient, "CreateCLSIndex", func(request *teov20220901.CreateCLSIndexRequest) (*teov20220901.CreateCLSIndexResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "task-87654321", *request.TaskId)
		resp := teov20220901.NewCreateCLSIndexResponse()
		resp.Response = &teov20220901.CreateCLSIndexResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"task_id": "task-87654321",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())
}

// TestCreateCLSIndexOperation_APIError tests Create handles API error
func TestCreateCLSIndexOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateCLSIndex", func(request *teov20220901.CreateCLSIndexRequest) (*teov20220901.CreateCLSIndexResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"task_id": "task-87654321",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestCreateCLSIndexOperation_Read tests Read is no-op
func TestCreateCLSIndexOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"task_id": "task-87654321",
	})
	d.SetId("zone-12345678")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestCreateCLSIndexOperation_Delete tests Delete is no-op
func TestCreateCLSIndexOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"task_id": "task-87654321",
	})
	d.SetId("zone-12345678")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestCreateCLSIndexOperation_Schema validates schema definition
func TestCreateCLSIndexOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "task_id")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	taskId := res.Schema["task_id"]
	assert.Equal(t, schema.TypeString, taskId.Type)
	assert.True(t, taskId.Required)
	assert.True(t, taskId.ForceNew)
}

func ptrString(s string) *string {
	return &s
}
