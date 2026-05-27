package ga2_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

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

func ptrString(s string) *string    { return &s }
func ptrFloat64(f float64) *float64 { return &f }

// go test ./tencentcloud/services/ga2/ -run "TestGa2CrossBorderSettlementDataSource" -v -count=1 -gcflags="all=-l"

// TestGa2CrossBorderSettlementDataSource_ReadSuccess tests successful read with traffic data
func TestGa2CrossBorderSettlementDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeCrossBorderSettlementWithContext", func(ctx context.Context, request *ga2v20250115.DescribeCrossBorderSettlementRequest) (*ga2v20250115.DescribeCrossBorderSettlementResponse, error) {
		assert.Equal(t, "ga2-test123", *request.GlobalAcceleratorId)
		assert.Equal(t, "ap-guangzhou", *request.AccelerateRegion)
		assert.Equal(t, "ap-singapore", *request.EndpointGroupRegion)
		assert.Equal(t, uint64(202501), *request.SettlementMonth)

		resp := ga2v20250115.NewDescribeCrossBorderSettlementResponse()
		resp.Response = &ga2v20250115.DescribeCrossBorderSettlementResponseParams{
			Traffic:   ptrFloat64(123.456789),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := ga2.DataSourceTencentCloudGa2CrossBorderSettlement()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-test123",
		"accelerate_region":     "ap-guangzhou",
		"endpoint_group_region": "ap-singapore",
		"settlement_month":      202501,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "ga2-test123#ap-guangzhou#ap-singapore#202501", d.Id())

	traffic := d.Get("traffic").(float64)
	assert.Equal(t, 123.456789, traffic)
}

// TestGa2CrossBorderSettlementDataSource_ReadNilTraffic tests read when API returns nil traffic
func TestGa2CrossBorderSettlementDataSource_ReadNilTraffic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeCrossBorderSettlementWithContext", func(ctx context.Context, request *ga2v20250115.DescribeCrossBorderSettlementRequest) (*ga2v20250115.DescribeCrossBorderSettlementResponse, error) {
		resp := ga2v20250115.NewDescribeCrossBorderSettlementResponse()
		resp.Response = &ga2v20250115.DescribeCrossBorderSettlementResponseParams{
			Traffic:   nil,
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := ga2.DataSourceTencentCloudGa2CrossBorderSettlement()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga2-nil-traffic",
		"accelerate_region":     "ap-beijing",
		"endpoint_group_region": "ap-tokyo",
		"settlement_month":      202502,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	traffic := d.Get("traffic").(float64)
	assert.Equal(t, float64(0), traffic)
}

// TestGa2CrossBorderSettlementDataSource_Schema validates schema definition
func TestGa2CrossBorderSettlementDataSource_Schema(t *testing.T) {
	res := ga2.DataSourceTencentCloudGa2CrossBorderSettlement()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "global_accelerator_id")
	assert.Contains(t, res.Schema, "accelerate_region")
	assert.Contains(t, res.Schema, "endpoint_group_region")
	assert.Contains(t, res.Schema, "settlement_month")
	assert.Contains(t, res.Schema, "traffic")
	assert.Contains(t, res.Schema, "result_output_file")

	// global_accelerator_id is Required TypeString
	gaId := res.Schema["global_accelerator_id"]
	assert.Equal(t, schema.TypeString, gaId.Type)
	assert.True(t, gaId.Required)

	// accelerate_region is Required TypeString
	accelRegion := res.Schema["accelerate_region"]
	assert.Equal(t, schema.TypeString, accelRegion.Type)
	assert.True(t, accelRegion.Required)

	// endpoint_group_region is Required TypeString
	egRegion := res.Schema["endpoint_group_region"]
	assert.Equal(t, schema.TypeString, egRegion.Type)
	assert.True(t, egRegion.Required)

	// settlement_month is Required TypeInt
	month := res.Schema["settlement_month"]
	assert.Equal(t, schema.TypeInt, month.Type)
	assert.True(t, month.Required)

	// traffic is Computed TypeFloat
	traffic := res.Schema["traffic"]
	assert.Equal(t, schema.TypeFloat, traffic.Type)
	assert.True(t, traffic.Computed)

	// result_output_file is Optional TypeString
	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}
