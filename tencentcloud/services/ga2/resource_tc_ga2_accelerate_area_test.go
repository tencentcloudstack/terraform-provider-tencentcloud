package ga2_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcga2 "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

type mockMetaForGa2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForGa2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForGa2{}

func newMockMetaForGa2() *mockMetaForGa2 {
	return &mockMetaForGa2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringGa2(s string) *string {
	return &s
}

func ptrUint64Ga2(i uint64) *uint64 {
	return &i
}

// go test ./tencentcloud/services/ga2/ -run "TestGa2AccelerateArea" -v -count=1 -gcflags="all=-l"

// TestGa2AccelerateArea_Create_Success tests the Create function
func TestGa2AccelerateArea_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Mock CreateAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "CreateAccelerateAreas", func(request *ga2v20250115.CreateAccelerateAreasRequest) (*ga2v20250115.CreateAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewCreateAccelerateAreasResponse()
		resp.Response = &ga2v20250115.CreateAccelerateAreasResponseParams{
			TaskId:    ptrStringGa2("task-123"),
			RequestId: ptrStringGa2("req-123"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResult - return SUCCESS
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResult", func(request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("req-456"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerateAreas for the Read call after Create
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AcceleratorAreaId: ptrStringGa2("area-001"),
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(10),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					IpAddress:         []*string{ptrStringGa2("1.1.1.1")},
					IpAddressInfoSet: []*ga2v20250115.IpAddressInfoSet{
						{
							IpAddress: ptrStringGa2("1.1.1.1"),
							IspType:   ptrStringGa2("BGP"),
						},
					},
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("req-789"),
		}
		return resp, nil
	})

	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
		"accelerator_areas": []interface{}{
			map[string]interface{}{
				"accelerate_region": "ap-guangzhou",
				"bandwidth":         10,
				"isp_type":          "BGP",
				"ip_version":        "IPv4",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-12345678", d.Id())
}

// TestGa2AccelerateArea_Read_Success tests the Read function
func TestGa2AccelerateArea_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AcceleratorAreaId: ptrStringGa2("area-001"),
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(10),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					IpAddress:         []*string{ptrStringGa2("1.1.1.1")},
					IpAddressInfoSet: []*ga2v20250115.IpAddressInfoSet{
						{
							IpAddress: ptrStringGa2("1.1.1.1"),
							IspType:   ptrStringGa2("BGP"),
						},
					},
				},
				{
					AcceleratorAreaId: ptrStringGa2("area-002"),
					AccelerateRegion:  ptrStringGa2("ap-shanghai"),
					Bandwidth:         ptrUint64Ga2(20),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					IpAddress:         []*string{ptrStringGa2("2.2.2.2")},
					IpAddressInfoSet:  nil,
				},
			},
			TotalCount: ptrUint64Ga2(2),
			RequestId:  ptrStringGa2("req-123"),
		}
		return resp, nil
	})

	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
	})
	d.SetId("ga-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-12345678", d.Get("global_accelerator_id"))

	areaSet := d.Get("accelerate_area_set").([]interface{})
	assert.Equal(t, 2, len(areaSet))

	area1 := areaSet[0].(map[string]interface{})
	assert.Equal(t, "area-001", area1["accelerator_area_id"])
	assert.Equal(t, "ap-guangzhou", area1["accelerate_region"])
	assert.Equal(t, 10, area1["bandwidth"])
	assert.Equal(t, "BGP", area1["isp_type"])
	assert.Equal(t, "IPv4", area1["ip_version"])

	area2 := areaSet[1].(map[string]interface{})
	assert.Equal(t, "area-002", area2["accelerator_area_id"])
	assert.Equal(t, "ap-shanghai", area2["accelerate_region"])
	assert.Equal(t, 20, area2["bandwidth"])
}

// TestGa2AccelerateArea_Read_Empty tests Read when no areas exist
func TestGa2AccelerateArea_Read_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{},
			TotalCount:        ptrUint64Ga2(0),
			RequestId:         ptrStringGa2("req-123"),
		}
		return resp, nil
	})

	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
	})
	d.SetId("ga-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestGa2AccelerateArea_Update_Success tests the Update function
func TestGa2AccelerateArea_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Mock ModifyAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "ModifyAccelerateAreas", func(request *ga2v20250115.ModifyAccelerateAreasRequest) (*ga2v20250115.ModifyAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewModifyAccelerateAreasResponse()
		resp.Response = &ga2v20250115.ModifyAccelerateAreasResponseParams{
			TaskId:    ptrStringGa2("task-456"),
			RequestId: ptrStringGa2("req-456"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResult
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResult", func(request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("req-789"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerateAreas for the Read call after Update
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AcceleratorAreaId: ptrStringGa2("area-001"),
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(50),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					IpAddress:         []*string{ptrStringGa2("1.1.1.1")},
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("req-101"),
		}
		return resp, nil
	})

	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := res.TestResourceData()
	d.SetId("ga-12345678")
	_ = d.Set("global_accelerator_id", "ga-12345678")
	_ = d.Set("accelerator_areas", []interface{}{
		map[string]interface{}{
			"accelerate_region": "ap-guangzhou",
			"bandwidth":         50,
			"isp_type":          "BGP",
			"ip_version":        "IPv4",
		},
	})
	_ = d.Set("accelerate_area_set", []interface{}{
		map[string]interface{}{
			"accelerator_area_id": "area-001",
			"accelerate_region":   "ap-guangzhou",
			"bandwidth":           10,
			"isp_type":            "BGP",
			"ip_version":          "IPv4",
			"ip_address":          []interface{}{"1.1.1.1"},
			"ip_address_info_set": []interface{}{},
		},
	})

	// Simulate HasChange by marking the field as changed
	d.MarkNewResource()

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestGa2AccelerateArea_Delete_Success tests the Delete function
func TestGa2AccelerateArea_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas to return existing areas
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AcceleratorAreaId: ptrStringGa2("area-001"),
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(10),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
				},
				{
					AcceleratorAreaId: ptrStringGa2("area-002"),
					AccelerateRegion:  ptrStringGa2("ap-shanghai"),
					Bandwidth:         ptrUint64Ga2(20),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
				},
			},
			TotalCount: ptrUint64Ga2(2),
			RequestId:  ptrStringGa2("req-123"),
		}
		return resp, nil
	})

	// Mock DeleteAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "DeleteAccelerateAreas", func(request *ga2v20250115.DeleteAccelerateAreasRequest) (*ga2v20250115.DeleteAccelerateAreasResponse, error) {
		assert.Equal(t, "ga-12345678", *request.GlobalAcceleratorId)
		assert.Equal(t, 2, len(request.AcceleratorAreaIds))
		resp := ga2v20250115.NewDeleteAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DeleteAccelerateAreasResponseParams{
			TaskId:    ptrStringGa2("task-789"),
			RequestId: ptrStringGa2("req-456"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResult
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResult", func(request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("req-789"),
		}
		return resp, nil
	})

	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
	})
	d.SetId("ga-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestGa2AccelerateArea_Delete_NoAreas tests Delete when no areas exist
func TestGa2AccelerateArea_Delete_NoAreas(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas to return empty
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{},
			TotalCount:        ptrUint64Ga2(0),
			RequestId:         ptrStringGa2("req-123"),
		}
		return resp, nil
	})

	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
	})
	d.SetId("ga-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestGa2AccelerateArea_Schema tests the schema definition
func TestGa2AccelerateArea_Schema(t *testing.T) {
	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	assert.NotNil(t, res)

	// Check global_accelerator_id
	assert.Contains(t, res.Schema, "global_accelerator_id")
	gaIdSchema := res.Schema["global_accelerator_id"]
	assert.Equal(t, schema.TypeString, gaIdSchema.Type)
	assert.True(t, gaIdSchema.Required)
	assert.True(t, gaIdSchema.ForceNew)

	// Check accelerator_areas
	assert.Contains(t, res.Schema, "accelerator_areas")
	areasSchema := res.Schema["accelerator_areas"]
	assert.Equal(t, schema.TypeList, areasSchema.Type)
	assert.True(t, areasSchema.Required)

	// Check accelerate_area_set
	assert.Contains(t, res.Schema, "accelerate_area_set")
	areaSetSchema := res.Schema["accelerate_area_set"]
	assert.Equal(t, schema.TypeList, areaSetSchema.Type)
	assert.True(t, areaSetSchema.Computed)

	// Check Importer
	assert.NotNil(t, res.Importer)

	// Check Timeouts
	assert.NotNil(t, res.Timeouts)
	assert.NotNil(t, res.Timeouts.Create)
	assert.NotNil(t, res.Timeouts.Update)
	assert.NotNil(t, res.Timeouts.Delete)
}
