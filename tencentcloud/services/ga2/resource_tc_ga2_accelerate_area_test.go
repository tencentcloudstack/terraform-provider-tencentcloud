package ga2_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

// ---- Mock helpers ----

type ga2MockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *ga2MockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &ga2MockMeta{}

func newGa2MockMeta() *ga2MockMeta {
	return &ga2MockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string {
	return &s
}

func ptrUint64(i uint64) *uint64 {
	return &i
}

// ---- Unit Tests ----

// go test ./tencentcloud/services/ga2/ -run "TestGa2AccelerateArea_" -v -count=1 -gcflags="all=-l"

func TestGa2AccelerateArea_Schema(t *testing.T) {
	res := ga2.ResourceTencentCloudGa2AccelerateArea()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "global_accelerator_id")
	assert.Contains(t, res.Schema, "accelerator_areas")
	assert.Contains(t, res.Schema, "accelerate_area_set")
	assert.Contains(t, res.Schema, "task_id")

	// Check global_accelerator_id
	globalAcceleratorId := res.Schema["global_accelerator_id"]
	assert.Equal(t, schema.TypeString, globalAcceleratorId.Type)
	assert.True(t, globalAcceleratorId.Required)
	assert.True(t, globalAcceleratorId.ForceNew)

	// Check accelerator_areas
	acceleratorAreas := res.Schema["accelerator_areas"]
	assert.Equal(t, schema.TypeList, acceleratorAreas.Type)
	assert.True(t, acceleratorAreas.Required)

	elem := acceleratorAreas.Elem.(*schema.Resource)
	assert.Contains(t, elem.Schema, "accelerate_region")
	assert.Contains(t, elem.Schema, "bandwidth")
	assert.Contains(t, elem.Schema, "isp_type")
	assert.Contains(t, elem.Schema, "ip_version")
	assert.Contains(t, elem.Schema, "accelerator_area_id")
	assert.Contains(t, elem.Schema, "ip_address")
	assert.Contains(t, elem.Schema, "ip_address_info_set")

	// Check accelerate_area_set is Computed
	accelerateAreaSet := res.Schema["accelerate_area_set"]
	assert.Equal(t, schema.TypeList, accelerateAreaSet.Type)
	assert.True(t, accelerateAreaSet.Computed)

	// Check task_id is Computed
	taskId := res.Schema["task_id"]
	assert.Equal(t, schema.TypeString, taskId.Type)
	assert.True(t, taskId.Computed)

	// Check Timeouts
	assert.NotNil(t, res.Timeouts)
	assert.NotNil(t, res.Timeouts.Create)
	assert.NotNil(t, res.Timeouts.Update)
	assert.NotNil(t, res.Timeouts.Delete)
}

func TestGa2AccelerateArea_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newGa2MockMeta().client, "UseGa2V20250115Client", ga2Client)

	// Mock CreateAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "CreateAccelerateAreas",
		func(request *ga2v20250115.CreateAccelerateAreasRequest) (*ga2v20250115.CreateAccelerateAreasResponse, error) {
			assert.NotNil(t, request.GlobalAcceleratorId)
			assert.Equal(t, "ga-12345678", *request.GlobalAcceleratorId)
			assert.Len(t, request.AcceleratorAreas, 1)
			assert.Equal(t, "ap-guangzhou", *request.AcceleratorAreas[0].AccelerateRegion)
			assert.Equal(t, uint64(10), *request.AcceleratorAreas[0].Bandwidth)

			resp := ga2v20250115.NewCreateAccelerateAreasResponse()
			resp.Response = &ga2v20250115.CreateAccelerateAreasResponseParams{
				TaskId:    ptrString("task-001"),
				RequestId: ptrString("req-001"),
			}
			return resp, nil
		},
	)

	// Mock DescribeAccelerateAreas (for polling and read)
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas",
		func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
			resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
			resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
				TotalCount: ptrUint64(1),
				AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
					{
						AccelerateRegion:  ptrString("ap-guangzhou"),
						Bandwidth:         ptrUint64(10),
						IspType:           ptrString("BGP"),
						IpVersion:         ptrString("IPv4"),
						AcceleratorAreaId: ptrString("area-001"),
						IpAddress:         []*string{ptrString("1.2.3.4")},
						IpAddressInfoSet: []*ga2v20250115.IpAddressInfoSet{
							{
								IpAddress: ptrString("1.2.3.4"),
								IspType:   ptrString("BGP"),
							},
						},
					},
				},
				RequestId: ptrString("req-002"),
			}
			return resp, nil
		},
	)

	meta := newGa2MockMeta()
	res := ga2.ResourceTencentCloudGa2AccelerateArea()
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

	// Verify accelerate_area_set is populated
	areaSet := d.Get("accelerate_area_set").([]interface{})
	assert.Len(t, areaSet, 1)
	areaMap := areaSet[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou", areaMap["accelerate_region"])
	assert.Equal(t, 10, areaMap["bandwidth"])
	assert.Equal(t, "BGP", areaMap["isp_type"])
	assert.Equal(t, "IPv4", areaMap["ip_version"])
	assert.Equal(t, "area-001", areaMap["accelerator_area_id"])
}

func TestGa2AccelerateArea_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newGa2MockMeta().client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas",
		func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
			assert.Equal(t, "ga-12345678", *request.GlobalAcceleratorId)

			resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
			resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
				TotalCount: ptrUint64(1),
				AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
					{
						AccelerateRegion:  ptrString("ap-shanghai"),
						Bandwidth:         ptrUint64(20),
						IspType:           ptrString("BGP"),
						IpVersion:         ptrString("IPv4"),
						AcceleratorAreaId: ptrString("area-002"),
						IpAddress:         []*string{ptrString("5.6.7.8")},
						IpAddressInfoSet: []*ga2v20250115.IpAddressInfoSet{
							{
								IpAddress: ptrString("5.6.7.8"),
								IspType:   ptrString("BGP"),
							},
						},
					},
				},
				RequestId: ptrString("req-003"),
			}
			return resp, nil
		},
	)

	meta := newGa2MockMeta()
	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
		"accelerator_areas":     []interface{}{},
	})
	d.SetId("ga-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-12345678", d.Id())

	// Verify accelerate_area_set
	areaSet := d.Get("accelerate_area_set").([]interface{})
	assert.Len(t, areaSet, 1)
	areaMap := areaSet[0].(map[string]interface{})
	assert.Equal(t, "ap-shanghai", areaMap["accelerate_region"])
	assert.Equal(t, 20, areaMap["bandwidth"])
	assert.Equal(t, "area-002", areaMap["accelerator_area_id"])
}

func TestGa2AccelerateArea_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newGa2MockMeta().client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas returns empty
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas",
		func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
			resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
			resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
				TotalCount:        ptrUint64(0),
				AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{},
				RequestId:         ptrString("req-004"),
			}
			return resp, nil
		},
	)

	meta := newGa2MockMeta()
	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
		"accelerator_areas":     []interface{}{},
	})
	d.SetId("ga-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestGa2AccelerateArea_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newGa2MockMeta().client, "UseGa2V20250115Client", ga2Client)

	// Mock ModifyAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "ModifyAccelerateAreas",
		func(request *ga2v20250115.ModifyAccelerateAreasRequest) (*ga2v20250115.ModifyAccelerateAreasResponse, error) {
			assert.Equal(t, "ga-12345678", *request.GlobalAcceleratorId)
			assert.Len(t, request.AcceleratorAreas, 1)
			assert.Equal(t, "ap-guangzhou", *request.AcceleratorAreas[0].AccelerateRegion)
			assert.Equal(t, uint64(50), *request.AcceleratorAreas[0].Bandwidth)
			assert.Equal(t, "area-001", *request.AcceleratorAreas[0].AcceleratorAreaId)

			resp := ga2v20250115.NewModifyAccelerateAreasResponse()
			resp.Response = &ga2v20250115.ModifyAccelerateAreasResponseParams{
				TaskId:    ptrString("task-002"),
				RequestId: ptrString("req-005"),
			}
			return resp, nil
		},
	)

	// Mock DescribeAccelerateAreas (for polling and read)
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas",
		func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
			resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
			resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
				TotalCount: ptrUint64(1),
				AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
					{
						AccelerateRegion:  ptrString("ap-guangzhou"),
						Bandwidth:         ptrUint64(50),
						IspType:           ptrString("BGP"),
						IpVersion:         ptrString("IPv4"),
						AcceleratorAreaId: ptrString("area-001"),
						IpAddress:         []*string{ptrString("1.2.3.4")},
						IpAddressInfoSet: []*ga2v20250115.IpAddressInfoSet{
							{
								IpAddress: ptrString("1.2.3.4"),
								IspType:   ptrString("BGP"),
							},
						},
					},
				},
				RequestId: ptrString("req-006"),
			}
			return resp, nil
		},
	)

	meta := newGa2MockMeta()
	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
		"accelerator_areas": []interface{}{
			map[string]interface{}{
				"accelerate_region":   "ap-guangzhou",
				"bandwidth":           50,
				"isp_type":            "BGP",
				"ip_version":          "IPv4",
				"accelerator_area_id": "area-001",
			},
		},
	})
	d.SetId("ga-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify accelerate_area_set is updated
	areaSet := d.Get("accelerate_area_set").([]interface{})
	assert.Len(t, areaSet, 1)
	areaMap := areaSet[0].(map[string]interface{})
	assert.Equal(t, 50, areaMap["bandwidth"])
}

func TestGa2AccelerateArea_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newGa2MockMeta().client, "UseGa2V20250115Client", ga2Client)

	describeCallCount := 0

	// Mock DescribeAccelerateAreas (first call returns areas, second returns empty)
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas",
		func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
			describeCallCount++
			resp := ga2v20250115.NewDescribeAccelerateAreasResponse()

			if describeCallCount <= 1 {
				// First call: return existing areas for deletion
				resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
					TotalCount: ptrUint64(1),
					AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
						{
							AccelerateRegion:  ptrString("ap-guangzhou"),
							Bandwidth:         ptrUint64(10),
							IspType:           ptrString("BGP"),
							IpVersion:         ptrString("IPv4"),
							AcceleratorAreaId: ptrString("area-001"),
						},
					},
					RequestId: ptrString("req-007"),
				}
			} else {
				// Subsequent calls: areas deleted
				resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
					TotalCount:        ptrUint64(0),
					AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{},
					RequestId:         ptrString("req-009"),
				}
			}
			return resp, nil
		},
	)

	// Mock DeleteAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "DeleteAccelerateAreas",
		func(request *ga2v20250115.DeleteAccelerateAreasRequest) (*ga2v20250115.DeleteAccelerateAreasResponse, error) {
			assert.Equal(t, "ga-12345678", *request.GlobalAcceleratorId)
			assert.Len(t, request.AcceleratorAreaIds, 1)
			assert.Equal(t, "area-001", *request.AcceleratorAreaIds[0])

			resp := ga2v20250115.NewDeleteAccelerateAreasResponse()
			resp.Response = &ga2v20250115.DeleteAccelerateAreasResponseParams{
				TaskId:    ptrString("task-003"),
				RequestId: ptrString("req-008"),
			}
			return resp, nil
		},
	)

	meta := newGa2MockMeta()
	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
		"accelerator_areas":     []interface{}{},
	})
	d.SetId("ga-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
