package cbs_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cbssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	cbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
)

type mockMetaCbsCopySnapshot struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCbsCopySnapshot) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCbsCopySnapshot{}

func newMockMetaCbsCopySnapshot() *mockMetaCbsCopySnapshot {
	return &mockMetaCbsCopySnapshot{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCSCR(s string) *string {
	return &s
}

func TestCbsCopySnapshotCrossRegion_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbssdk.Client{}
	patches.ApplyMethodReturn(newMockMetaCbsCopySnapshot().client, "UseCbsClient", cbsClient)

	// Mock CopySnapshotCrossRegionsWithContext
	patches.ApplyMethodFunc(cbsClient, "CopySnapshotCrossRegionsWithContext", func(_ context.Context, request *cbssdk.CopySnapshotCrossRegionsRequest) (*cbssdk.CopySnapshotCrossRegionsResponse, error) {
		assert.NotNil(t, request.SnapshotId)
		assert.Equal(t, "snap-source-xxx", *request.SnapshotId)
		assert.NotNil(t, request.DestinationRegions)
		assert.Equal(t, 2, len(request.DestinationRegions))
		assert.Equal(t, "ap-shanghai", *request.DestinationRegions[0])
		assert.Equal(t, "ap-chengdu", *request.DestinationRegions[1])

		resp := cbssdk.NewCopySnapshotCrossRegionsResponse()
		resp.Response = &cbssdk.CopySnapshotCrossRegionsResponseParams{
			SnapshotCopyResultSet: []*cbssdk.SnapshotCopyResult{
				{
					SnapshotId:        ptrStringCSCR("snap-shanghai-yyy"),
					Code:              ptrStringCSCR("Success"),
					Message:           ptrStringCSCR(""),
					DestinationRegion: ptrStringCSCR("ap-shanghai"),
				},
				{
					SnapshotId:        ptrStringCSCR("snap-chengdu-zzz"),
					Code:              ptrStringCSCR("Success"),
					Message:           ptrStringCSCR(""),
					DestinationRegion: ptrStringCSCR("ap-chengdu"),
				},
			},
			RequestId: ptrStringCSCR("fake-request-id"),
		}
		return resp, nil
	})

	// Mock CbsService.DescribeSnapshotById for async polling
	patches.ApplyMethodFunc(&cbs.CbsService{}, "DescribeSnapshotById", func(_ context.Context, snapshotId string) (*cbssdk.Snapshot, error) {
		return &cbssdk.Snapshot{
			SnapshotId:    ptrStringCSCR(snapshotId),
			SnapshotState: ptrStringCSCR("NORMAL"),
			SnapshotName:  ptrStringCSCR("my-copied-snapshot"),
		}, nil
	})

	meta := newMockMetaCbsCopySnapshot()
	res := cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"snapshot_id":         "snap-source-xxx",
		"destination_regions": []interface{}{"ap-shanghai", "ap-chengdu"},
		"snapshot_name":       "my-copied-snapshot",
		"delete_bind_images":  false,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "snap-source-xxx#snap-shanghai-yyy", d.Id())
	assert.Equal(t, "snap-source-xxx", d.Get("snapshot_id"))

	// Verify snapshot_copy_result_set
	results := d.Get("snapshot_copy_result_set").([]interface{})
	assert.Equal(t, 2, len(results))

	result0 := results[0].(map[string]interface{})
	assert.Equal(t, "snap-shanghai-yyy", result0["snapshot_id"])
	assert.Equal(t, "Success", result0["code"])
	assert.Equal(t, "ap-shanghai", result0["destination_region"])

	result1 := results[1].(map[string]interface{})
	assert.Equal(t, "snap-chengdu-zzz", result1["snapshot_id"])
	assert.Equal(t, "Success", result1["code"])
	assert.Equal(t, "ap-chengdu", result1["destination_region"])
}

func TestCbsCopySnapshotCrossRegion_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbssdk.Client{}
	patches.ApplyMethodReturn(newMockMetaCbsCopySnapshot().client, "UseCbsClient", cbsClient)

	// Mock CbsService.DescribeSnapshotById for Read
	patches.ApplyMethodFunc(&cbs.CbsService{}, "DescribeSnapshotById", func(_ context.Context, snapshotId string) (*cbssdk.Snapshot, error) {
		assert.Equal(t, "snap-shanghai-yyy", snapshotId)
		return &cbssdk.Snapshot{
			SnapshotId:    ptrStringCSCR("snap-shanghai-yyy"),
			SnapshotState: ptrStringCSCR("NORMAL"),
			SnapshotName:  ptrStringCSCR("my-copied-snapshot"),
		}, nil
	})

	meta := newMockMetaCbsCopySnapshot()
	res := cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"snapshot_id":         "snap-source-xxx",
		"destination_regions": []interface{}{"ap-shanghai"},
		"delete_bind_images":  false,
	})
	d.SetId("snap-source-xxx#snap-shanghai-yyy")

	// Set snapshot_copy_result_set in state before Read
	snapshotCopyResultSet := []map[string]interface{}{
		{
			"snapshot_id":        "snap-shanghai-yyy",
			"code":               "Success",
			"message":            "",
			"destination_region": "ap-shanghai",
		},
	}
	_ = d.Set("snapshot_copy_result_set", snapshotCopyResultSet)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "snap-source-xxx#snap-shanghai-yyy", d.Id())
	assert.Equal(t, "snap-source-xxx", d.Get("snapshot_id"))
	assert.Equal(t, "my-copied-snapshot", d.Get("snapshot_name"))

	// Verify snapshot_copy_result_set is preserved (not overwritten since it was already set)
	results := d.Get("snapshot_copy_result_set").([]interface{})
	assert.Equal(t, 1, len(results))
	result0 := results[0].(map[string]interface{})
	assert.Equal(t, "snap-shanghai-yyy", result0["snapshot_id"])
}

func TestCbsCopySnapshotCrossRegion_ReadImportScenario(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbssdk.Client{}
	patches.ApplyMethodReturn(newMockMetaCbsCopySnapshot().client, "UseCbsClient", cbsClient)

	// Mock CbsService.DescribeSnapshotById for Read after import (no snapshot_copy_result_set in state)
	patches.ApplyMethodFunc(&cbs.CbsService{}, "DescribeSnapshotById", func(_ context.Context, snapshotId string) (*cbssdk.Snapshot, error) {
		assert.Equal(t, "snap-shanghai-yyy", snapshotId)
		return &cbssdk.Snapshot{
			SnapshotId:    ptrStringCSCR("snap-shanghai-yyy"),
			SnapshotState: ptrStringCSCR("NORMAL"),
			SnapshotName:  ptrStringCSCR("my-copied-snapshot"),
			Placement: &cbssdk.Placement{
				Zone: ptrStringCSCR("ap-shanghai-2"),
			},
		}, nil
	})

	meta := newMockMetaCbsCopySnapshot()
	res := cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"snapshot_id":         "snap-source-xxx",
		"destination_regions": []interface{}{"ap-shanghai"},
		"delete_bind_images":  false,
	})
	d.SetId("snap-source-xxx#snap-shanghai-yyy")
	// Note: snapshot_copy_result_set is NOT set, simulating import scenario

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "snap-source-xxx#snap-shanghai-yyy", d.Id())
	assert.Equal(t, "snap-source-xxx", d.Get("snapshot_id"))
	assert.Equal(t, "my-copied-snapshot", d.Get("snapshot_name"))

	// Verify snapshot_copy_result_set was reconstructed from the composite ID and API response
	results := d.Get("snapshot_copy_result_set").([]interface{})
	assert.Equal(t, 1, len(results))
	result0 := results[0].(map[string]interface{})
	assert.Equal(t, "snap-shanghai-yyy", result0["snapshot_id"])
	assert.Equal(t, "Success", result0["code"])
}

func TestCbsCopySnapshotCrossRegion_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbssdk.Client{}
	patches.ApplyMethodReturn(newMockMetaCbsCopySnapshot().client, "UseCbsClient", cbsClient)

	// Mock CbsService.DescribeSnapshotById returning nil (snapshot not found)
	patches.ApplyMethodFunc(&cbs.CbsService{}, "DescribeSnapshotById", func(_ context.Context, snapshotId string) (*cbssdk.Snapshot, error) {
		return nil, nil
	})

	meta := newMockMetaCbsCopySnapshot()
	res := cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"snapshot_id":         "snap-source-xxx",
		"destination_regions": []interface{}{"ap-shanghai"},
		"delete_bind_images":  false,
	})
	d.SetId("snap-source-xxx#snap-shanghai-yyy")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestCbsCopySnapshotCrossRegion_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbssdk.Client{}
	patches.ApplyMethodReturn(newMockMetaCbsCopySnapshot().client, "UseCbsClient", cbsClient)

	// Mock DeleteSnapshotsWithContext
	patches.ApplyMethodFunc(cbsClient, "DeleteSnapshotsWithContext", func(_ context.Context, request *cbssdk.DeleteSnapshotsRequest) (*cbssdk.DeleteSnapshotsResponse, error) {
		assert.NotNil(t, request.SnapshotIds)
		assert.Equal(t, 2, len(request.SnapshotIds))
		assert.Equal(t, "snap-shanghai-yyy", *request.SnapshotIds[0])
		assert.Equal(t, "snap-chengdu-zzz", *request.SnapshotIds[1])
		assert.NotNil(t, request.DeleteBindImages)
		assert.Equal(t, false, *request.DeleteBindImages)

		resp := cbssdk.NewDeleteSnapshotsResponse()
		resp.Response = &cbssdk.DeleteSnapshotsResponseParams{
			RequestId: ptrStringCSCR("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaCbsCopySnapshot()
	res := cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"snapshot_id":         "snap-source-xxx",
		"destination_regions": []interface{}{"ap-shanghai", "ap-chengdu"},
		"delete_bind_images":  false,
	})
	d.SetId("snap-source-xxx#snap-shanghai-yyy")

	// Set snapshot_copy_result_set in state before Delete
	snapshotCopyResultSet := []map[string]interface{}{
		{
			"snapshot_id":        "snap-shanghai-yyy",
			"code":               "Success",
			"message":            "",
			"destination_region": "ap-shanghai",
		},
		{
			"snapshot_id":        "snap-chengdu-zzz",
			"code":               "Success",
			"message":            "",
			"destination_region": "ap-chengdu",
		},
	}
	_ = d.Set("snapshot_copy_result_set", snapshotCopyResultSet)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestCbsCopySnapshotCrossRegion_DeleteImportScenario(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbssdk.Client{}
	patches.ApplyMethodReturn(newMockMetaCbsCopySnapshot().client, "UseCbsClient", cbsClient)

	// Mock DeleteSnapshotsWithContext - should use composite ID fallback
	patches.ApplyMethodFunc(cbsClient, "DeleteSnapshotsWithContext", func(_ context.Context, request *cbssdk.DeleteSnapshotsRequest) (*cbssdk.DeleteSnapshotsResponse, error) {
		assert.NotNil(t, request.SnapshotIds)
		assert.Equal(t, 1, len(request.SnapshotIds))
		assert.Equal(t, "snap-shanghai-yyy", *request.SnapshotIds[0])
		assert.NotNil(t, request.DeleteBindImages)
		assert.Equal(t, false, *request.DeleteBindImages)

		resp := cbssdk.NewDeleteSnapshotsResponse()
		resp.Response = &cbssdk.DeleteSnapshotsResponseParams{
			RequestId: ptrStringCSCR("fake-request-id-delete-import"),
		}
		return resp, nil
	})

	meta := newMockMetaCbsCopySnapshot()
	res := cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"snapshot_id":         "snap-source-xxx",
		"destination_regions": []interface{}{"ap-shanghai"},
		"delete_bind_images":  false,
	})
	d.SetId("snap-source-xxx#snap-shanghai-yyy")
	// Note: snapshot_copy_result_set is NOT set, simulating import scenario

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestCbsCopySnapshotCrossRegion_DeleteWithBindImages(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbssdk.Client{}
	patches.ApplyMethodReturn(newMockMetaCbsCopySnapshot().client, "UseCbsClient", cbsClient)

	// Mock DeleteSnapshotsWithContext with delete_bind_images=true
	patches.ApplyMethodFunc(cbsClient, "DeleteSnapshotsWithContext", func(_ context.Context, request *cbssdk.DeleteSnapshotsRequest) (*cbssdk.DeleteSnapshotsResponse, error) {
		assert.NotNil(t, request.SnapshotIds)
		assert.Equal(t, 1, len(request.SnapshotIds))
		assert.Equal(t, "snap-shanghai-yyy", *request.SnapshotIds[0])
		assert.NotNil(t, request.DeleteBindImages)
		assert.Equal(t, true, *request.DeleteBindImages)

		resp := cbssdk.NewDeleteSnapshotsResponse()
		resp.Response = &cbssdk.DeleteSnapshotsResponseParams{
			RequestId: ptrStringCSCR("fake-request-id-delete-bind-images"),
		}
		return resp, nil
	})

	meta := newMockMetaCbsCopySnapshot()
	res := cbs.ResourceTencentCloudCbsCopySnapshotCrossRegionAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"snapshot_id":         "snap-source-xxx",
		"destination_regions": []interface{}{"ap-shanghai"},
		"delete_bind_images":  true,
	})
	d.SetId("snap-source-xxx#snap-shanghai-yyy")

	// Set snapshot_copy_result_set in state before Delete
	snapshotCopyResultSet := []map[string]interface{}{
		{
			"snapshot_id":        "snap-shanghai-yyy",
			"code":               "Success",
			"message":            "",
			"destination_region": "ap-shanghai",
		},
	}
	_ = d.Set("snapshot_copy_result_set", snapshotCopyResultSet)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
