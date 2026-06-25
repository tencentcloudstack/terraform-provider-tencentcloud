package cdb_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
func ptrBoolCdbExpand(b bool) *bool       { return &b }

func TestCdbStartCpuExpandAttachment_Create_Auto_Success(t *testing.T) {
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

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpandAttachment()
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
	assert.Equal(t, "async-request-id-123", d.Get("async_request_id"))
	assert.Equal(t, "auto", d.Get("type"))
}

func TestCdbStartCpuExpandAttachment_Create_Manual_Success(t *testing.T) {
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

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpandAttachment()
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

func TestCdbStartCpuExpandAttachment_Read_Success(t *testing.T) {
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
	res := cdb.ResourceTencentCloudCdbStartCpuExpandAttachment()
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
	if len(autoStrategy) > 0 {
		autoStrategyMap := autoStrategy[0].(map[string]interface{})
		assert.Equal(t, 80, autoStrategyMap["expand_threshold"])
		assert.Equal(t, 20, autoStrategyMap["shrink_threshold"])
	}
}

func TestCdbStartCpuExpandAttachment_Read_NotFound(t *testing.T) {
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
	res := cdb.ResourceTencentCloudCdbStartCpuExpandAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-test1234",
		"type":        "auto",
	})
	d.SetId("cdb-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestCdbStartCpuExpandAttachment_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbStartCpuExpand().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "StopCpuExpandWithContext", func(_ context.Context, request *cdb_sdk.StopCpuExpandRequest) (*cdb_sdk.StopCpuExpandResponse, error) {
		resp := cdb_sdk.NewStopCpuExpandResponse()
		resp.Response = &cdb_sdk.StopCpuExpandResponseParams{
			AsyncRequestId: ptrStringCdbExpand("async-delete-request-id-789"),
			RequestId:      ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cdbClient, "DescribeCPUExpandStrategyInfo", func(request *cdb_sdk.DescribeCPUExpandStrategyInfoRequest) (*cdb_sdk.DescribeCPUExpandStrategyInfoResponse, error) {
		resp := cdb_sdk.NewDescribeCPUExpandStrategyInfoResponse()
		resp.Response = &cdb_sdk.DescribeCPUExpandStrategyInfoResponseParams{
			Type:      nil,
			RequestId: ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpandAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-test1234",
		"type":        "auto",
	})
	d.SetId("cdb-test1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestCdbStartCpuExpandAttachment_Update_ImmutableChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cdbClient := &cdb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCdbStartCpuExpand().client, "UseMysqlClient", cdbClient)

	patches.ApplyMethodFunc(cdbClient, "DescribeCPUExpandStrategyInfo", func(request *cdb_sdk.DescribeCPUExpandStrategyInfoRequest) (*cdb_sdk.DescribeCPUExpandStrategyInfoResponse, error) {
		resp := cdb_sdk.NewDescribeCPUExpandStrategyInfoResponse()
		resp.Response = &cdb_sdk.DescribeCPUExpandStrategyInfoResponseParams{
			Type: ptrStringCdbExpand("auto"),
			AutoStrategy: &cdb_sdk.AutoStrategy{
				ExpandThreshold: ptrInt64CdbExpand(80),
				ShrinkThreshold: ptrInt64CdbExpand(20),
			},
			RequestId: ptrStringCdbExpand("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCdbStartCpuExpand()
	res := cdb.ResourceTencentCloudCdbStartCpuExpandAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cdb-test1234",
		"type":        "auto",
		"auto_strategy": []interface{}{
			map[string]interface{}{
				"expand_threshold": 80,
				"shrink_threshold": 20,
			},
		},
	})
	d.SetId("cdb-test1234")

	_ = d.Set("type", "manual")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "immutable argument")
}
