package cdb_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	cdb_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
)

type mockMetaForCdbStartCpuExpand struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCdbStartCpuExpand) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCdbStartCpuExpand{}

func newMockMetaForCdbStartCpuExpand() *mockMetaForCdbStartCpuExpand {
	return &mockMetaForCdbStartCpuExpand{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCdbExpand(s string) *string { return &s }
func ptrInt64CdbExpand(v int64) *int64    { return &v }

// patchAsyncRequestSuccess mocks the async task polling so that Create/Delete
// finish without hitting the real cloud API.
func patchAsyncRequestSuccess(patches *gomonkey.Patches, cdbClient *cdb_sdk.Client) {
	patches.ApplyMethodFunc(cdbClient, "DescribeAsyncRequestInfo", func(request *cdb_sdk.DescribeAsyncRequestInfoRequest) (*cdb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := cdb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &cdb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status:    ptrStringCdbExpand("SUCCESS"),
			Info:      ptrStringCdbExpand("ok"),
			RequestId: ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})
}

func TestCdbStartCpuExpand_Create_Auto_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbStartCpuExpand().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "StartCpuExpandWithContext", func(_ context.Context, request *cdb_sdk.StartCpuExpandRequest) (*cdb_sdk.StartCpuExpandResponse, error) {
		resp := cdb_sdk.NewStartCpuExpandResponse()
		resp.Response = &cdb_sdk.StartCpuExpandResponseParams{
			AsyncRequestId: ptrStringCdbExpand("async-request-id-123"),
			RequestId:      ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cdbClient, "DescribeCPUExpandStrategyInfo", func(request *cdb_sdk.DescribeCPUExpandStrategyInfoRequest) (*cdb_sdk.DescribeCPUExpandStrategyInfoResponse, error) {
		resp := cdb_sdk.NewDescribeCPUExpandStrategyInfoResponse()
		resp.Response = &cdb_sdk.DescribeCPUExpandStrategyInfoResponseParams{
			Type: ptrStringCdbExpand("auto"),
			AutoStrategy: &cdb_sdk.AutoStrategy{
				ExpandThreshold:    ptrInt64CdbExpand(80),
				ShrinkThreshold:    ptrInt64CdbExpand(20),
				ExpandSecondPeriod: ptrInt64CdbExpand(300),
				ShrinkSecondPeriod: ptrInt64CdbExpand(600),
			},
			RequestId: ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	patchAsyncRequestSuccess(patches, cdbClient)

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpand()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-test1234",
		"type":        "auto",
		"auto_strategy": []interface{}{
			map[string]interface{}{
				"expand_threshold":     80,
				"shrink_threshold":     20,
				"expand_second_period": 300,
				"shrink_second_period": 600,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cdb-test1234", d.Id())
	assert.Equal(t, "auto", d.Get("type"))
}

func TestCdbStartCpuExpand_Create_Manual_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbStartCpuExpand().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "StartCpuExpandWithContext", func(_ context.Context, request *cdb_sdk.StartCpuExpandRequest) (*cdb_sdk.StartCpuExpandResponse, error) {
		resp := cdb_sdk.NewStartCpuExpandResponse()
		resp.Response = &cdb_sdk.StartCpuExpandResponseParams{
			AsyncRequestId: ptrStringCdbExpand("async-request-id-456"),
			RequestId:      ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cdbClient, "DescribeCPUExpandStrategyInfo", func(request *cdb_sdk.DescribeCPUExpandStrategyInfoRequest) (*cdb_sdk.DescribeCPUExpandStrategyInfoResponse, error) {
		resp := cdb_sdk.NewDescribeCPUExpandStrategyInfoResponse()
		resp.Response = &cdb_sdk.DescribeCPUExpandStrategyInfoResponseParams{
			Type:      ptrStringCdbExpand("manual"),
			ExpandCpu: ptrInt64CdbExpand(4),
			RequestId: ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	patchAsyncRequestSuccess(patches, cdbClient)

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpand()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-test5678",
		"type":        "manual",
		"expand_cpu":  4,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cdb-test5678", d.Id())
	assert.Equal(t, "manual", d.Get("type"))
	assert.Equal(t, 4, d.Get("expand_cpu"))
}

func TestCdbStartCpuExpand_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbStartCpuExpand().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "DescribeCPUExpandStrategyInfo", func(request *cdb_sdk.DescribeCPUExpandStrategyInfoRequest) (*cdb_sdk.DescribeCPUExpandStrategyInfoResponse, error) {
		resp := cdb_sdk.NewDescribeCPUExpandStrategyInfoResponse()
		resp.Response = &cdb_sdk.DescribeCPUExpandStrategyInfoResponseParams{
			Type: ptrStringCdbExpand("auto"),
			AutoStrategy: &cdb_sdk.AutoStrategy{
				ExpandThreshold:    ptrInt64CdbExpand(80),
				ShrinkThreshold:    ptrInt64CdbExpand(20),
				ExpandSecondPeriod: ptrInt64CdbExpand(300),
				ShrinkSecondPeriod: ptrInt64CdbExpand(600),
			},
			RequestId: ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpand()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-test1234",
		"type":        "auto",
	})
	d.SetId("cdb-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cdb-test1234", d.Id())
	assert.Equal(t, "auto", d.Get("type"))
	assert.Equal(t, "cdb-test1234", d.Get("instance_id"))

	autoStrategy := d.Get("auto_strategy").([]interface{})
	assert.NotEmpty(t, autoStrategy)
	autoStrategyMap := autoStrategy[0].(map[string]interface{})
	assert.Equal(t, 80, autoStrategyMap["expand_threshold"])
	assert.Equal(t, 20, autoStrategyMap["shrink_threshold"])
	assert.Equal(t, 300, autoStrategyMap["expand_second_period"])
	assert.Equal(t, 600, autoStrategyMap["shrink_second_period"])
}

func TestCdbStartCpuExpand_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbStartCpuExpand().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "DescribeCPUExpandStrategyInfo", func(request *cdb_sdk.DescribeCPUExpandStrategyInfoRequest) (*cdb_sdk.DescribeCPUExpandStrategyInfoResponse, error) {
		resp := cdb_sdk.NewDescribeCPUExpandStrategyInfoResponse()
		resp.Response = &cdb_sdk.DescribeCPUExpandStrategyInfoResponseParams{
			Type:      nil,
			RequestId: ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpand()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-test1234",
		"type":        "auto",
	})
	d.SetId("cdb-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestCdbStartCpuExpand_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbStartCpuExpand().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "StopCpuExpand", func(request *cdb_sdk.StopCpuExpandRequest) (*cdb_sdk.StopCpuExpandResponse, error) {
		resp := cdb_sdk.NewStopCpuExpandResponse()
		resp.Response = &cdb_sdk.StopCpuExpandResponseParams{
			AsyncRequestId: ptrStringCdbExpand("async-delete-request-id-789"),
			RequestId:      ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	patchAsyncRequestSuccess(patches, cdbClient)

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpand()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-test1234",
		"type":        "auto",
	})
	d.SetId("cdb-test1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestCdbStartCpuExpand_Update_ImmutableChange(t *testing.T) {
	res := cdb.ResourceTencentCloudCdbStartCpuExpand()

	// Build a prior state where type=auto, and a new config where type=manual,
	// so that d.HasChange("type") returns true inside Update and triggers the
	// immutable-argument error.
	state := &terraform.InstanceState{
		ID: "cdb-test1234",
		Attributes: map[string]string{
			"id":          "cdb-test1234",
			"instance_id": "cdb-test1234",
			"type":        "auto",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"instance_id": "cdb-test1234",
		"type":        "manual",
	})

	diff, err := res.Diff(nil, state, rawConfig, newMockMetaForCdbStartCpuExpand())
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	updateErr := res.Update(d, newMockMetaForCdbStartCpuExpand())
	assert.Error(t, updateErr)
	assert.Contains(t, updateErr.Error(), "immutable argument")
}
