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

type mockMetaGa2GlobalAccelerator struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaGa2GlobalAccelerator) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaGa2GlobalAccelerator{}

func newMockMetaGa2GlobalAccelerator() *mockMetaGa2GlobalAccelerator {
	return &mockMetaGa2GlobalAccelerator{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringGa2(s string) *string {
	return &s
}

func ptrBoolGa2(b bool) *bool {
	return &b
}

// go test ./tencentcloud/services/ga2/ -run "TestUnitGa2GlobalAccelerator" -v -count=1 -gcflags="all=-l"

func TestUnitGa2GlobalAccelerator_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2GlobalAccelerator()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch CreateGlobalAcceleratorWithContext
	patches.ApplyMethodFunc(ga2Client, "CreateGlobalAcceleratorWithContext", func(_ interface{}, request *ga2v20250115.CreateGlobalAcceleratorRequest) (*ga2v20250115.CreateGlobalAcceleratorResponse, error) {
		assert.Equal(t, "tf-test-ga2", *request.Name)
		assert.Equal(t, "POSTPAID", *request.InstanceChargeType)
		assert.Equal(t, "test description", *request.Description)
		assert.Equal(t, "HighQuality", *request.CrossBorderType)
		assert.Equal(t, true, *request.CrossBorderPromiseFlag)

		resp := &ga2v20250115.CreateGlobalAcceleratorResponse{}
		resp.Response = &ga2v20250115.CreateGlobalAcceleratorResponseParams{
			GlobalAcceleratorId: ptrStringGa2("ga2-test001"),
			TaskId:              ptrStringGa2("task-create-001"),
			RequestId:           ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeTaskResultWithContext for WaitForGa2TaskFinish
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(_ interface{}, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := &ga2v20250115.DescribeTaskResultResponse{}
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeGlobalAcceleratorsWithContext for Read
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorsWithContext", func(_ interface{}, request *ga2v20250115.DescribeGlobalAcceleratorsRequest) (*ga2v20250115.DescribeGlobalAcceleratorsResponse, error) {
		resp := &ga2v20250115.DescribeGlobalAcceleratorsResponse{}
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorsResponseParams{
			GlobalAcceleratorSet: []*ga2v20250115.GlobalAcceleratorSet{
				{
					GlobalAcceleratorId: ptrStringGa2("ga2-test001"),
					Name:                ptrStringGa2("tf-test-ga2"),
					Description:         ptrStringGa2("test description"),
					InstanceChargeType:  ptrStringGa2("POSTPAID"),
					CrossBorderType:     ptrStringGa2("HighQuality"),
					CreateTime:          ptrStringGa2("2025-01-15 10:00:00"),
					State:               ptrStringGa2("RUNNING"),
					Status:              ptrStringGa2("NORMAL"),
					DdosId:              ptrStringGa2("ddos-001"),
					Cname:               ptrStringGa2("ga2-test001.example.com"),
					TagSet: []*ga2v20250115.Tag{
						{
							Key:   ptrStringGa2("createdBy"),
							Value: ptrStringGa2("terraform"),
						},
					},
				},
			},
			TotalCount: func() *uint64 { v := uint64(1); return &v }(),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2GlobalAccelerator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                      "tf-test-ga2",
		"instance_charge_type":      "POSTPAID",
		"description":               "test description",
		"cross_border_type":         "HighQuality",
		"cross_border_promise_flag": true,
		"tags": map[string]interface{}{
			"createdBy": "terraform",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ga2-test001", d.Id())
}

func TestUnitGa2GlobalAccelerator_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2GlobalAccelerator()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch DescribeGlobalAcceleratorsWithContext
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorsWithContext", func(_ interface{}, request *ga2v20250115.DescribeGlobalAcceleratorsRequest) (*ga2v20250115.DescribeGlobalAcceleratorsResponse, error) {
		resp := &ga2v20250115.DescribeGlobalAcceleratorsResponse{}
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorsResponseParams{
			GlobalAcceleratorSet: []*ga2v20250115.GlobalAcceleratorSet{
				{
					GlobalAcceleratorId: ptrStringGa2("ga2-test001"),
					Name:                ptrStringGa2("tf-test-ga2"),
					Description:         ptrStringGa2("test description"),
					InstanceChargeType:  ptrStringGa2("POSTPAID"),
					CrossBorderType:     ptrStringGa2("HighQuality"),
					CreateTime:          ptrStringGa2("2025-01-15 10:00:00"),
					State:               ptrStringGa2("RUNNING"),
					Status:              ptrStringGa2("NORMAL"),
					DdosId:              ptrStringGa2("ddos-001"),
					Cname:               ptrStringGa2("ga2-test001.example.com"),
					TagSet: []*ga2v20250115.Tag{
						{
							Key:   ptrStringGa2("createdBy"),
							Value: ptrStringGa2("terraform"),
						},
					},
				},
			},
			TotalCount: func() *uint64 { v := uint64(1); return &v }(),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2GlobalAccelerator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("ga2-test001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf-test-ga2", d.Get("name"))
	assert.Equal(t, "POSTPAID", d.Get("instance_charge_type"))
	assert.Equal(t, "test description", d.Get("description"))
	assert.Equal(t, "HighQuality", d.Get("cross_border_type"))
	assert.Equal(t, "2025-01-15 10:00:00", d.Get("create_time"))
	assert.Equal(t, "RUNNING", d.Get("state"))
	assert.Equal(t, "NORMAL", d.Get("status"))
	assert.Equal(t, "ddos-001", d.Get("ddos_id"))
	assert.Equal(t, "ga2-test001.example.com", d.Get("cname"))
}

func TestUnitGa2GlobalAccelerator_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2GlobalAccelerator()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch DescribeGlobalAcceleratorsWithContext - return empty
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorsWithContext", func(_ interface{}, request *ga2v20250115.DescribeGlobalAcceleratorsRequest) (*ga2v20250115.DescribeGlobalAcceleratorsResponse, error) {
		resp := &ga2v20250115.DescribeGlobalAcceleratorsResponse{}
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorsResponseParams{
			GlobalAcceleratorSet: []*ga2v20250115.GlobalAcceleratorSet{},
			TotalCount:           func() *uint64 { v := uint64(0); return &v }(),
			RequestId:            ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2GlobalAccelerator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("ga2-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestUnitGa2GlobalAccelerator_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2GlobalAccelerator()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch ModifyGlobalAcceleratorWithContext
	patches.ApplyMethodFunc(ga2Client, "ModifyGlobalAcceleratorWithContext", func(_ interface{}, request *ga2v20250115.ModifyGlobalAcceleratorRequest) (*ga2v20250115.ModifyGlobalAcceleratorResponse, error) {
		assert.Equal(t, "ga2-test001", *request.GlobalAcceleratorId)
		assert.Equal(t, "tf-test-ga2-updated", *request.Name)
		assert.Equal(t, "updated description", *request.Description)

		resp := &ga2v20250115.ModifyGlobalAcceleratorResponse{}
		resp.Response = &ga2v20250115.ModifyGlobalAcceleratorResponseParams{
			TaskId:    ptrStringGa2("task-modify-001"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeTaskResultWithContext for WaitForGa2TaskFinish
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(_ interface{}, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := &ga2v20250115.DescribeTaskResultResponse{}
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeGlobalAcceleratorsWithContext for Read after update
	patches.ApplyMethodFunc(ga2Client, "DescribeGlobalAcceleratorsWithContext", func(_ interface{}, request *ga2v20250115.DescribeGlobalAcceleratorsRequest) (*ga2v20250115.DescribeGlobalAcceleratorsResponse, error) {
		resp := &ga2v20250115.DescribeGlobalAcceleratorsResponse{}
		resp.Response = &ga2v20250115.DescribeGlobalAcceleratorsResponseParams{
			GlobalAcceleratorSet: []*ga2v20250115.GlobalAcceleratorSet{
				{
					GlobalAcceleratorId: ptrStringGa2("ga2-test001"),
					Name:                ptrStringGa2("tf-test-ga2-updated"),
					Description:         ptrStringGa2("updated description"),
					InstanceChargeType:  ptrStringGa2("POSTPAID"),
					CrossBorderType:     ptrStringGa2("HighQuality"),
					CreateTime:          ptrStringGa2("2025-01-15 10:00:00"),
					State:               ptrStringGa2("RUNNING"),
					Status:              ptrStringGa2("NORMAL"),
					DdosId:              ptrStringGa2("ddos-001"),
					Cname:               ptrStringGa2("ga2-test001.example.com"),
				},
			},
			TotalCount: func() *uint64 { v := uint64(1); return &v }(),
			RequestId:  ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2GlobalAccelerator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":                 "tf-test-ga2-updated",
		"instance_charge_type": "POSTPAID",
		"description":          "updated description",
		"cross_border_type":    "HighQuality",
	})
	d.SetId("ga2-test001")

	// Simulate HasChange by marking fields as changed
	d.MarkNewResource()

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf-test-ga2-updated", d.Get("name"))
	assert.Equal(t, "updated description", d.Get("description"))
}

func TestUnitGa2GlobalAccelerator_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaGa2GlobalAccelerator()
	ga2Client := &ga2v20250115.Client{}
	patches.ApplyMethodReturn(meta.client, "UseGa2V20250115Client", ga2Client)

	// Patch DeleteGlobalAcceleratorWithContext
	patches.ApplyMethodFunc(ga2Client, "DeleteGlobalAcceleratorWithContext", func(_ interface{}, request *ga2v20250115.DeleteGlobalAcceleratorRequest) (*ga2v20250115.DeleteGlobalAcceleratorResponse, error) {
		assert.Equal(t, "ga2-test001", *request.GlobalAcceleratorId)

		resp := &ga2v20250115.DeleteGlobalAcceleratorResponse{}
		resp.Response = &ga2v20250115.DeleteGlobalAcceleratorResponseParams{
			TaskId:    ptrStringGa2("task-delete-001"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeTaskResultWithContext for WaitForGa2TaskFinish
	patches.ApplyMethodFunc(ga2Client, "DescribeTaskResultWithContext", func(_ interface{}, request *ga2v20250115.DescribeTaskResultRequest) (*ga2v20250115.DescribeTaskResultResponse, error) {
		resp := &ga2v20250115.DescribeTaskResultResponse{}
		resp.Response = &ga2v20250115.DescribeTaskResultResponseParams{
			Status:    ptrStringGa2("SUCCESS"),
			RequestId: ptrStringGa2("fake-request-id"),
		}
		return resp, nil
	})

	res := ga2.ResourceTencentCloudGa2GlobalAccelerator()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("ga2-test001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
