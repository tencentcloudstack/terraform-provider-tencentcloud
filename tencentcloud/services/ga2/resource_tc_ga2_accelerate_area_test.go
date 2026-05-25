package ga2_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
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

func ptrString(s string) *string {
	return &s
}

func ptrUint64(i uint64) *uint64 {
	return &i
}

// go test ./tencentcloud/services/ga2/ -run "TestGa2AccelerateArea" -v -count=1 -gcflags="all=-l"

func TestGa2AccelerateArea_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateAccelerateAreasWithContext", func(_ interface{}, request *ga2v20250115.CreateAccelerateAreasRequest) (*ga2v20250115.CreateAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewCreateAccelerateAreasResponse()
		resp.Response = &ga2v20250115.CreateAccelerateAreasResponseParams{
			TaskId:    ptrString("task-123"),
			RequestId: ptrString("req-123"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			TotalCount: ptrUint64(1),
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AcceleratorAreaId: ptrString("area-123"),
					AccelerateRegion:  ptrString("ap-guangzhou"),
					Bandwidth:         ptrUint64(10),
					IspType:           ptrString("BGP"),
					IpVersion:         ptrString("IPv4"),
					IpAddress:         []*string{ptrString("1.1.1.1")},
				},
			},
			RequestId: ptrString("req-456"),
		}
		return resp, nil
	})

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
}

func TestGa2AccelerateArea_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateAccelerateAreasWithContext", func(_ interface{}, request *ga2v20250115.CreateAccelerateAreasRequest) (*ga2v20250115.CreateAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewCreateAccelerateAreasResponse()
		resp.Response = &ga2v20250115.CreateAccelerateAreasResponseParams{
			TaskId:    nil,
			RequestId: ptrString("req-123"),
		}
		return resp, nil
	})

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
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TaskId is nil")
}

func TestGa2AccelerateArea_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			TotalCount: ptrUint64(2),
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AcceleratorAreaId: ptrString("area-001"),
					AccelerateRegion:  ptrString("ap-guangzhou"),
					Bandwidth:         ptrUint64(10),
					IspType:           ptrString("BGP"),
					IpVersion:         ptrString("IPv4"),
					IpAddress:         []*string{ptrString("1.1.1.1")},
				},
				{
					AcceleratorAreaId: ptrString("area-002"),
					AccelerateRegion:  ptrString("ap-shanghai"),
					Bandwidth:         ptrUint64(20),
					IspType:           ptrString("BGP"),
					IpVersion:         ptrString("IPv4"),
					IpAddress:         []*string{ptrString("2.2.2.2")},
				},
			},
			RequestId: ptrString("req-789"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
	})
	d.SetId("ga-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-12345678", d.Id())
	assert.Equal(t, "ga-12345678", d.Get("global_accelerator_id"))
}

func TestGa2AccelerateArea_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			TotalCount:        ptrUint64(0),
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{},
			RequestId:         ptrString("req-000"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
	})
	d.SetId("ga-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestGa2AccelerateArea_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "ModifyAccelerateAreasWithContext", func(_ interface{}, request *ga2v20250115.ModifyAccelerateAreasRequest) (*ga2v20250115.ModifyAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewModifyAccelerateAreasResponse()
		resp.Response = &ga2v20250115.ModifyAccelerateAreasResponseParams{
			TaskId:    ptrString("task-456"),
			RequestId: ptrString("req-456"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			TotalCount: ptrUint64(1),
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AcceleratorAreaId: ptrString("area-123"),
					AccelerateRegion:  ptrString("ap-guangzhou"),
					Bandwidth:         ptrUint64(50),
					IspType:           ptrString("BGP"),
					IpVersion:         ptrString("IPv4"),
					IpAddress:         []*string{ptrString("1.1.1.1")},
				},
			},
			RequestId: ptrString("req-789"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
		"accelerator_areas": []interface{}{
			map[string]interface{}{
				"accelerate_region": "ap-guangzhou",
				"bandwidth":         50,
				"isp_type":          "BGP",
				"ip_version":        "IPv4",
			},
		},
	})
	d.SetId("ga-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga-12345678", d.Id())
}

func TestGa2AccelerateArea_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			TotalCount: ptrUint64(1),
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{
				{
					AcceleratorAreaId: ptrString("area-123"),
					AccelerateRegion:  ptrString("ap-guangzhou"),
					Bandwidth:         ptrUint64(10),
					IspType:           ptrString("BGP"),
					IpVersion:         ptrString("IPv4"),
					IpAddress:         []*string{ptrString("1.1.1.1")},
				},
			},
			RequestId: ptrString("req-111"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(ga2Client, "DeleteAccelerateAreasWithContext", func(_ interface{}, request *ga2v20250115.DeleteAccelerateAreasRequest) (*ga2v20250115.DeleteAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDeleteAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DeleteAccelerateAreasResponseParams{
			TaskId:    ptrString("task-789"),
			RequestId: ptrString("req-222"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
	})
	d.SetId("ga-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestGa2AccelerateArea_Delete_NoAreas(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateAreas", func(request *ga2v20250115.DescribeAccelerateAreasRequest) (*ga2v20250115.DescribeAccelerateAreasResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateAreasResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateAreasResponseParams{
			TotalCount:        ptrUint64(0),
			AccelerateAreaSet: []*ga2v20250115.AcceleratorAreas{},
			RequestId:         ptrString("req-333"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2AccelerateArea()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"global_accelerator_id": "ga-12345678",
	})
	d.SetId("ga-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestGa2AccelerateArea_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	meta := newMockMetaForGa2()
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "CreateAccelerateAreasWithContext", func(_ interface{}, request *ga2v20250115.CreateAccelerateAreasRequest) (*ga2v20250115.CreateAccelerateAreasResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=UnsupportedOperation.InstanceStateNotAllowedOperate")
	})

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
	assert.Error(t, err)
}
