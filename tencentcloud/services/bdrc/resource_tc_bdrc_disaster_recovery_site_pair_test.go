package bdrc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"
)

type mockMetaBdrcSitePair struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBdrcSitePair) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBdrcSitePair{}

func newMockMetaBdrcSitePair() *mockMetaBdrcSitePair {
	return &mockMetaBdrcSitePair{client: &connectivity.TencentCloudClient{}}
}

func ptrStringSitePair(s string) *string { return &s }
func ptrInt64SitePair(v int64) *int64    { return &v }

// go test ./tencentcloud/services/bdrc/ -run "TestBdrcDisasterRecoverySitePair" -v -count=1 -gcflags="all=-l"

// TestBdrcDisasterRecoverySitePair_Create tests the Create operation
func TestBdrcDisasterRecoverySitePair_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSitePair().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "CreateDisasterRecoverySitePairWithContext", func(ctx context.Context, request *bdrcv20260330.CreateDisasterRecoverySitePairRequest) (*bdrcv20260330.CreateDisasterRecoverySitePairResponse, error) {
		assert.Equal(t, "CROSS_REGION", *request.DisasterRecoveryType)
		assert.Equal(t, "ap-guangzhou", *request.SourceRegion)
		assert.Equal(t, "ap-guangzhou-3", *request.SourceZone)
		assert.Equal(t, "ap-shanghai", *request.TargetRegion)
		assert.Equal(t, "ap-shanghai-2", *request.TargetZone)
		assert.Equal(t, "vpc-source-xxx", *request.SourceVpc)
		assert.Equal(t, "vpc-target-yyy", *request.TargetVpc)
		assert.Equal(t, "DISK", *request.SitePairProductType)
		assert.Equal(t, "tf-example-site-pair", *request.SitePairName)
		assert.Equal(t, "ASY", *request.CopyType)

		resp := bdrcv20260330.NewCreateDisasterRecoverySitePairResponse()
		resp.Response = &bdrcv20260330.CreateDisasterRecoverySitePairResponseParams{
			SitePairId: ptrStringSitePair("site-pair-create-001"),
			RequestId:  ptrStringSitePair("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(bdrcClient, "DescribeDisasterRecoverySitePairsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeDisasterRecoverySitePairsRequest) (*bdrcv20260330.DescribeDisasterRecoverySitePairsResponse, error) {
		resp := bdrcv20260330.NewDescribeDisasterRecoverySitePairsResponse()
		resp.Response = &bdrcv20260330.DescribeDisasterRecoverySitePairsResponseParams{
			TotalCount: ptrInt64SitePair(1),
			SitePairSet: []*bdrcv20260330.SitePair{
				{
					SitePairId:            ptrStringSitePair("site-pair-create-001"),
					SitePairName:          ptrStringSitePair("tf-example-site-pair"),
					SitePairType:          ptrStringSitePair("DISK"),
					SitePairState:         ptrStringSitePair("RUNNING"),
					DisasterRecoveryType:  ptrStringSitePair("CROSS_REGION"),
					SourceRegion:          ptrStringSitePair("ap-guangzhou"),
					SourceZone:            ptrStringSitePair("ap-guangzhou-3"),
					TargetRegion:          ptrStringSitePair("ap-shanghai"),
					TargetZone:            ptrStringSitePair("ap-shanghai-2"),
					SourceVpc:             ptrStringSitePair("vpc-source-xxx"),
					TargetVpc:             ptrStringSitePair("vpc-target-yyy"),
					SitePairProductType:   ptrStringSitePair("DISK"),
					CopyType:              ptrStringSitePair("ASY"),
					CreateFrom:            ptrStringSitePair("LOCAL"),
					AccountUin:            ptrStringSitePair("100000000001"),
					SubAccountUin:         ptrStringSitePair("200000000002"),
					CreateTime:            ptrStringSitePair("2026-09-15 10:00:00"),
					BindProtectGroupCount: ptrInt64SitePair(0),
				},
			},
			RequestId: ptrStringSitePair("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcSitePair()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoverySitePair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"disaster_recovery_type": "CROSS_REGION",
		"source_region":          "ap-guangzhou",
		"source_zone":            "ap-guangzhou-3",
		"target_region":          "ap-shanghai",
		"target_zone":            "ap-shanghai-2",
		"source_vpc":             "vpc-source-xxx",
		"target_vpc":             "vpc-target-yyy",
		"site_pair_product_type": "DISK",
		"site_pair_name":         "tf-example-site-pair",
		"copy_type":              "ASY",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "site-pair-create-001", d.Id())
	assert.Equal(t, "tf-example-site-pair", d.Get("site_pair_name").(string))
	assert.Equal(t, "RUNNING", d.Get("site_pair_state").(string))
	assert.Equal(t, "DISK", d.Get("site_pair_type").(string))
}

// TestBdrcDisasterRecoverySitePair_Create_EmptyId tests Create when SitePairId is empty
func TestBdrcDisasterRecoverySitePair_Create_EmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSitePair().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "CreateDisasterRecoverySitePairWithContext", func(ctx context.Context, request *bdrcv20260330.CreateDisasterRecoverySitePairRequest) (*bdrcv20260330.CreateDisasterRecoverySitePairResponse, error) {
		resp := bdrcv20260330.NewCreateDisasterRecoverySitePairResponse()
		resp.Response = &bdrcv20260330.CreateDisasterRecoverySitePairResponseParams{
			SitePairId: ptrStringSitePair(""),
			RequestId:  ptrStringSitePair("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcSitePair()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoverySitePair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"disaster_recovery_type": "CROSS_REGION",
		"source_region":          "ap-guangzhou",
		"source_zone":            "ap-guangzhou-3",
		"target_region":          "ap-shanghai",
		"target_zone":            "ap-shanghai-2",
		"source_vpc":             "vpc-source-xxx",
		"target_vpc":             "vpc-target-yyy",
		"site_pair_product_type": "DISK",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestBdrcDisasterRecoverySitePair_Read tests the Read operation
func TestBdrcDisasterRecoverySitePair_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSitePair().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "DescribeDisasterRecoverySitePairsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeDisasterRecoverySitePairsRequest) (*bdrcv20260330.DescribeDisasterRecoverySitePairsResponse, error) {
		resp := bdrcv20260330.NewDescribeDisasterRecoverySitePairsResponse()
		resp.Response = &bdrcv20260330.DescribeDisasterRecoverySitePairsResponseParams{
			TotalCount: ptrInt64SitePair(1),
			SitePairSet: []*bdrcv20260330.SitePair{
				{
					SitePairId:            ptrStringSitePair("site-pair-read-001"),
					SitePairName:          ptrStringSitePair("tf-read-site-pair"),
					SitePairType:          ptrStringSitePair("DISK"),
					SitePairState:         ptrStringSitePair("RUNNING"),
					DisasterRecoveryType:  ptrStringSitePair("CROSS_REGION"),
					SourceRegion:          ptrStringSitePair("ap-guangzhou"),
					SourceZone:            ptrStringSitePair("ap-guangzhou-3"),
					TargetRegion:          ptrStringSitePair("ap-shanghai"),
					TargetZone:            ptrStringSitePair("ap-shanghai-2"),
					SourceVpc:             ptrStringSitePair("vpc-source-xxx"),
					TargetVpc:             ptrStringSitePair("vpc-target-yyy"),
					CopyType:              ptrStringSitePair("ASY"),
					CreateFrom:            ptrStringSitePair("LOCAL"),
					AccountUin:            ptrStringSitePair("100000000001"),
					SubAccountUin:         ptrStringSitePair("200000000002"),
					CreateTime:            ptrStringSitePair("2026-09-15 10:00:00"),
					BindProtectGroupCount: ptrInt64SitePair(2),
				},
			},
			RequestId: ptrStringSitePair("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcSitePair()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoverySitePair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_product_type": "DISK",
	})
	d.SetId("site-pair-read-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "site-pair-read-001", d.Id())
	assert.Equal(t, "tf-read-site-pair", d.Get("site_pair_name").(string))
	assert.Equal(t, "RUNNING", d.Get("site_pair_state").(string))
	assert.Equal(t, "DISK", d.Get("site_pair_type").(string))
	assert.Equal(t, "CROSS_REGION", d.Get("disaster_recovery_type").(string))
	assert.Equal(t, 2, d.Get("bind_protect_group_count").(int))
}

// TestBdrcDisasterRecoverySitePair_Read_NotFound tests Read when the site pair is not found
func TestBdrcDisasterRecoverySitePair_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSitePair().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "DescribeDisasterRecoverySitePairsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeDisasterRecoverySitePairsRequest) (*bdrcv20260330.DescribeDisasterRecoverySitePairsResponse, error) {
		resp := bdrcv20260330.NewDescribeDisasterRecoverySitePairsResponse()
		resp.Response = &bdrcv20260330.DescribeDisasterRecoverySitePairsResponseParams{
			TotalCount:  ptrInt64SitePair(0),
			SitePairSet: []*bdrcv20260330.SitePair{},
			RequestId:   ptrStringSitePair("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcSitePair()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoverySitePair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_product_type": "DISK",
	})
	d.SetId("site-pair-not-found")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestBdrcDisasterRecoverySitePair_Update tests the Update operation for site_pair_name
func TestBdrcDisasterRecoverySitePair_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSitePair().client, "UseBdrcV20260330Client", bdrcClient)

	modifyCalled := false
	patches.ApplyMethodFunc(bdrcClient, "ModifySitePairAttributeWithContext", func(ctx context.Context, request *bdrcv20260330.ModifySitePairAttributeRequest) (*bdrcv20260330.ModifySitePairAttributeResponse, error) {
		modifyCalled = true
		assert.Equal(t, "site-pair-update-001", *request.SitePairId)
		assert.Equal(t, "tf-updated-site-pair", *request.SitePairName)

		resp := bdrcv20260330.NewModifySitePairAttributeResponse()
		resp.Response = &bdrcv20260330.ModifySitePairAttributeResponseParams{
			RequestId: ptrStringSitePair("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(bdrcClient, "DescribeDisasterRecoverySitePairsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeDisasterRecoverySitePairsRequest) (*bdrcv20260330.DescribeDisasterRecoverySitePairsResponse, error) {
		resp := bdrcv20260330.NewDescribeDisasterRecoverySitePairsResponse()
		resp.Response = &bdrcv20260330.DescribeDisasterRecoverySitePairsResponseParams{
			TotalCount: ptrInt64SitePair(1),
			SitePairSet: []*bdrcv20260330.SitePair{
				{
					SitePairId:            ptrStringSitePair("site-pair-update-001"),
					SitePairName:          ptrStringSitePair("tf-updated-site-pair"),
					SitePairType:          ptrStringSitePair("DISK"),
					SitePairState:         ptrStringSitePair("RUNNING"),
					DisasterRecoveryType:  ptrStringSitePair("CROSS_REGION"),
					SourceRegion:          ptrStringSitePair("ap-guangzhou"),
					SourceZone:            ptrStringSitePair("ap-guangzhou-3"),
					TargetRegion:          ptrStringSitePair("ap-shanghai"),
					TargetZone:            ptrStringSitePair("ap-shanghai-2"),
					SourceVpc:             ptrStringSitePair("vpc-source-xxx"),
					TargetVpc:             ptrStringSitePair("vpc-target-yyy"),
					CopyType:              ptrStringSitePair("ASY"),
					BindProtectGroupCount: ptrInt64SitePair(1),
				},
			},
			RequestId: ptrStringSitePair("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcSitePair()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoverySitePair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"disaster_recovery_type": "CROSS_REGION",
		"source_region":          "ap-guangzhou",
		"source_zone":            "ap-guangzhou-3",
		"target_region":          "ap-shanghai",
		"target_zone":            "ap-shanghai-2",
		"source_vpc":             "vpc-source-xxx",
		"target_vpc":             "vpc-target-yyy",
		"site_pair_product_type": "DISK",
		"site_pair_name":         "tf-updated-site-pair",
		"copy_type":              "ASY",
	})
	d.SetId("site-pair-update-001")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "tf-updated-site-pair", d.Get("site_pair_name").(string))
}

// TestBdrcDisasterRecoverySitePair_Delete tests the Delete operation
func TestBdrcDisasterRecoverySitePair_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcSitePair().client, "UseBdrcV20260330Client", bdrcClient)

	deleteCalled := false
	patches.ApplyMethodFunc(bdrcClient, "DeleteDisasterRecoverySitePairsWithContext", func(ctx context.Context, request *bdrcv20260330.DeleteDisasterRecoverySitePairsRequest) (*bdrcv20260330.DeleteDisasterRecoverySitePairsResponse, error) {
		deleteCalled = true
		assert.Equal(t, 1, len(request.SitePairIds))
		assert.Equal(t, "site-pair-delete-001", *request.SitePairIds[0])

		resp := bdrcv20260330.NewDeleteDisasterRecoverySitePairsResponse()
		resp.Response = &bdrcv20260330.DeleteDisasterRecoverySitePairsResponseParams{
			RequestId: ptrStringSitePair("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcSitePair()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoverySitePair()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("site-pair-delete-001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}
