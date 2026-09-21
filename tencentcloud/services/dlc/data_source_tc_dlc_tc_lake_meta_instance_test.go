package dlc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcdlc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
)

type mockMetaForDlcTCLakeMetaInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForDlcTCLakeMetaInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForDlcTCLakeMetaInstance{}

func newMockMetaForDlcTCLakeMetaInstance() *mockMetaForDlcTCLakeMetaInstance {
	return &mockMetaForDlcTCLakeMetaInstance{client: &connectivity.TencentCloudClient{}}
}

func ptrStrDlcTCLake(s string) *string {
	return &s
}

// go test ./tencentcloud/services/dlc/ -run "TestDlcTCLakeMetaInstance" -v -count=1 -gcflags="all=-l"

// TestDlcTCLakeMetaInstance_Read_Success tests successful read with a populated status response
func TestDlcTCLakeMetaInstance_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaForDlcTCLakeMetaInstance().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeTCLakeMetaInstance", func(request *dlc.DescribeTCLakeMetaInstanceRequest) (*dlc.DescribeTCLakeMetaInstanceResponse, error) {
		resp := dlc.NewDescribeTCLakeMetaInstanceResponse()
		resp.Response = &dlc.DescribeTCLakeMetaInstanceResponseParams{
			Status:    ptrStrDlcTCLake("Running"),
			RequestId: ptrStrDlcTCLake("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForDlcTCLakeMetaInstance()
	res := svcdlc.DataSourceTencentCloudDlcTCLakeMetaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tc_lake_meta_instance", d.Id())
	assert.Equal(t, "Running", d.Get("status"))
}

// TestDlcTCLakeMetaInstance_Read_NilStatus tests read when API returns nil status
func TestDlcTCLakeMetaInstance_Read_NilStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaForDlcTCLakeMetaInstance().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeTCLakeMetaInstance", func(request *dlc.DescribeTCLakeMetaInstanceRequest) (*dlc.DescribeTCLakeMetaInstanceResponse, error) {
		resp := dlc.NewDescribeTCLakeMetaInstanceResponse()
		resp.Response = &dlc.DescribeTCLakeMetaInstanceResponseParams{
			Status: nil,
		}
		return resp, nil
	})

	meta := newMockMetaForDlcTCLakeMetaInstance()
	res := svcdlc.DataSourceTencentCloudDlcTCLakeMetaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

// TestDlcTCLakeMetaInstance_Read_NilResponse tests read when API returns nil response
func TestDlcTCLakeMetaInstance_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaForDlcTCLakeMetaInstance().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeTCLakeMetaInstance", func(request *dlc.DescribeTCLakeMetaInstanceRequest) (*dlc.DescribeTCLakeMetaInstanceResponse, error) {
		resp := dlc.NewDescribeTCLakeMetaInstanceResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaForDlcTCLakeMetaInstance()
	res := svcdlc.DataSourceTencentCloudDlcTCLakeMetaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

// TestDlcTCLakeMetaInstance_Read_APIError tests read when API returns error
func TestDlcTCLakeMetaInstance_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dlcClient := &dlc.Client{}
	patches.ApplyMethodReturn(newMockMetaForDlcTCLakeMetaInstance().client, "UseDlcClient", dlcClient)

	patches.ApplyMethodFunc(dlcClient, "DescribeTCLakeMetaInstance", func(request *dlc.DescribeTCLakeMetaInstanceRequest) (*dlc.DescribeTCLakeMetaInstanceResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaForDlcTCLakeMetaInstance()
	res := svcdlc.DataSourceTencentCloudDlcTCLakeMetaInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestDlcTCLakeMetaInstance_Schema tests the schema definition
func TestDlcTCLakeMetaInstance_Schema(t *testing.T) {
	res := svcdlc.DataSourceTencentCloudDlcTCLakeMetaInstance()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	// Check computed output field
	assert.True(t, res.Schema["status"].Computed)
	assert.Equal(t, schema.TypeString, res.Schema["status"].Type)

	// Check optional result_output_file field
	assert.True(t, res.Schema["result_output_file"].Optional)
	assert.Equal(t, schema.TypeString, res.Schema["result_output_file"].Type)
}
