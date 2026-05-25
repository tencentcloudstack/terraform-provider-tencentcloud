package ga2_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcga2 "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

type mockMetaGa2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaGa2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaGa2{}

func newMockMetaGa2() *mockMetaGa2 {
	return &mockMetaGa2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringGa2(s string) *string {
	return &s
}

func ptrUint64Ga2(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/ga2/ -run "TestGa2AccelerateArea_" -v -count=1 -gcflags="all=-l"

func TestGa2AccelerateArea_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock CreateAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "CreateAccelerateAreas", func(_ *ga2v20250115.CreateAccelerateAreasRequest) (*ga2v20250115.CreateAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewCreateAccelerateAreasResponse()
		resp.Response = &ga2v20250115.CreateAccelerateAreasResponseParams{
			TaskId:    ptrStringGa2("task-create-123"),
			RequestId: ptrStringGa2("req-123"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResult
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResult", func(_ *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("req-456"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerateAreas for Read after Create
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(_ *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(10),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					AcceleratorAreaId: ptrStringGa2("area-001"),
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

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test001",
		"accelerator_areas": []interface{}{
			map[string]interface{}{
				"accelerate_region": "ap-guangzhou",
				"bandwidth":         10,
				"isp_type":          "BGP",
				"ip_version":        "IPv4",
			},
		},
	})

	// Set timeout for create
	d.SetId("")
	_ = d.Set("global_accelerator_id", "ga-test001")

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test001", d.Id())
}

func TestGa2AccelerateArea_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(_ *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(10),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					AcceleratorAreaId: ptrStringGa2("area-001"),
					IpAddress:         []*string{ptrStringGa2("1.1.1.1")},
					IpAddressInfoSet: []*ga2v20250115.IpAddressInfoSet{
						{
							IpAddress: ptrStringGa2("1.1.1.1"),
							IspType:   ptrStringGa2("BGP"),
						},
					},
				},
				{
					AccelerateRegion:  ptrStringGa2("ap-shanghai"),
					Bandwidth:         ptrUint64Ga2(20),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					AcceleratorAreaId: ptrStringGa2("area-002"),
					IpAddress:         []*string{ptrStringGa2("2.2.2.2")},
				},
			},
			TotalCount: ptrUint64Ga2(2),
			RequestId:  ptrStringGa2("req-read-001"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test001",
	})
	d.SetId("ga-test001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test001", d.Id())
	assert.Equal(t, "ga-test001", d.Get("global_accelerator_id").(string))

	areaSet := d.Get("accelerate_area_set").([]interface{})
	assert.Equal(t, 2, len(areaSet))
}

func TestGa2AccelerateArea_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas to return empty
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(_ *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{},
			TotalCount:        ptrUint64Ga2(0),
			RequestId:         ptrStringGa2("req-notfound"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test001",
	})
	d.SetId("ga-test001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestGa2AccelerateArea_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas for reading current areas (in update)
	callCount := 0
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(_ *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		callCount++
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(20),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					AcceleratorAreaId: ptrStringGa2("area-001"),
					IpAddress:         []*string{ptrStringGa2("1.1.1.1")},
				},
			},
			TotalCount: ptrUint64Ga2(1),
			RequestId:  ptrStringGa2("req-update-read"),
		}
		return resp, nil
	})

	// Mock ModifyAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "ModifyAccelerateAreas", func(_ *ga2v20250115.ModifyAccelerateAreasRequest) (*ga2v20250115.ModifyAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewModifyAccelerateAreasResponse()
		resp.Response = &ga2v20250115.ModifyAccelerateAreasResponseParams{
			TaskId:    ptrStringGa2("task-modify-123"),
			RequestId: ptrStringGa2("req-modify"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResult
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResult", func(_ *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("req-task"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test001",
		"accelerator_areas": []interface{}{
			map[string]interface{}{
				"accelerate_region": "ap-guangzhou",
				"bandwidth":         20,
				"isp_type":          "BGP",
				"ip_version":        "IPv4",
			},
		},
	})
	d.SetId("ga-test001")

	// Simulate HasChange by marking the field as changed
	// In unit tests with TestResourceDataRaw, HasChange returns true for all set fields
	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-test001", d.Id())
}

func TestGa2AccelerateArea_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas to return areas with IDs
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(_ *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(10),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					AcceleratorAreaId: ptrStringGa2("area-001"),
				},
				{
					AccelerateRegion:  ptrStringGa2("ap-shanghai"),
					Bandwidth:         ptrUint64Ga2(20),
					IspType:           ptrStringGa2("BGP"),
					IpVersion:         ptrStringGa2("IPv4"),
					AcceleratorAreaId: ptrStringGa2("area-002"),
				},
			},
			TotalCount: ptrUint64Ga2(2),
			RequestId:  ptrStringGa2("req-del-read"),
		}
		return resp, nil
	})

	// Mock DeleteAccelerateAreas
	patches.ApplyMethodFunc(ga2Client, "DeleteAccelerateAreas", func(req *ga2v20250115.DeleteAccelerateAreasRequest) (*ga2v20250115.DeleteAccelerateAreasResponse, error) {
		assert.Equal(t, "ga-test001", *req.GlobalAcceleratorId)
		assert.Equal(t, 2, len(req.AcceleratorAreaIds))
		resp := ga2v20250115.NewDeleteAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DeleteAccelerateAreasResponseParams{
			TaskId:    ptrStringGa2("task-delete-123"),
			RequestId: ptrStringGa2("req-delete"),
		}
		return resp, nil
	})

	// Mock DescribeTaskResult
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResult", func(_ *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := ga2v20250115.NewDescribeTaskResultResponse()
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("req-task-del"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test001",
	})
	d.SetId("ga-test001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestGa2AccelerateArea_Delete_NoAreas(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(newMockMetaGa2().client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas to return empty
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(_ *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{},
			TotalCount:        ptrUint64Ga2(0),
			RequestId:         ptrStringGa2("req-empty"),
		}
		return resp, nil
	})

	meta := newMockMetaGa2()
	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-test001",
	})
	d.SetId("ga-test001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestGa2AccelerateArea_Schema(t *testing.T) {
	res := svcga2.ResourceTencentCloudGa2AccelerateArea()
	assert.NotNil(t, res)

	// Verify global_accelerator_id is Required and ForceNew
	assert.Contains(t, res.Schema, "global_accelerator_id")
	gaIdSchema := res.Schema["global_accelerator_id"]
	assert.Equal(t, schema.TypeString, gaIdSchema.Type)
	assert.True(t, gaIdSchema.Required)
	assert.True(t, gaIdSchema.ForceNew)

	// Verify accelerator_areas is Required
	assert.Contains(t, res.Schema, "accelerator_areas")
	areasSchema := res.Schema["accelerator_areas"]
	assert.Equal(t, schema.TypeList, areasSchema.Type)
	assert.True(t, areasSchema.Required)

	// Verify accelerate_area_set is Computed
	assert.Contains(t, res.Schema, "accelerate_area_set")
	areaSetSchema := res.Schema["accelerate_area_set"]
	assert.Equal(t, schema.TypeList, areaSetSchema.Type)
	assert.True(t, areaSetSchema.Computed)

	// Verify task_id is Computed
	assert.Contains(t, res.Schema, "task_id")
	taskIdSchema := res.Schema["task_id"]
	assert.Equal(t, schema.TypeString, taskIdSchema.Type)
	assert.True(t, taskIdSchema.Computed)

	// Verify Timeouts
	assert.NotNil(t, res.Timeouts)
	assert.Equal(t, 10*time.Minute, *res.Timeouts.Create)
	assert.Equal(t, 10*time.Minute, *res.Timeouts.Update)
	assert.Equal(t, 10*time.Minute, *res.Timeouts.Delete)

	// Verify Importer
	assert.NotNil(t, res.Importer)
}

func TestGa2Service_DescribeAccelerateAreas(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	client := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(client, "UseGa2V20250115Client", ga2Client)

	// Mock DescribeAccelerateAreas with pagination
	callCount := 0
	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(req *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		callCount++
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		if callCount == 1 {
			// First page - return full page
			areas := make([]*ga2v20250115.AcceleratorAreas, 100)
			for i := 0; i < 100; i++ {
				areas[i] = &ga2v20250115.AcceleratorAreas{
					AccelerateRegion:  ptrStringGa2("ap-guangzhou"),
					Bandwidth:         ptrUint64Ga2(10),
					AcceleratorAreaId: ptrStringGa2(fmt.Sprintf("area-%d", i)),
				}
			}
			resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
				AccelerateAreaSet: areas,
				TotalCount:        ptrUint64Ga2(101),
				RequestId:         ptrStringGa2("req-page1"),
			}
		} else {
			// Second page - return partial
			resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
				AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
					{
						AccelerateRegion:  ptrStringGa2("ap-shanghai"),
						Bandwidth:         ptrUint64Ga2(20),
						AcceleratorAreaId: ptrStringGa2("area-last"),
					},
				},
				TotalCount: ptrUint64Ga2(101),
				RequestId:  ptrStringGa2("req-page2"),
			}
		}
		return resp, nil
	})

	service := svcga2.NewGa2Service(client)
	ctx := context.Background()
	areas, err := service.DescribeAccelerateAreas(ctx, "ga-test001")
	assert.NoError(t, err)
	assert.Equal(t, 101, len(areas))
	assert.Equal(t, 2, callCount)
}
