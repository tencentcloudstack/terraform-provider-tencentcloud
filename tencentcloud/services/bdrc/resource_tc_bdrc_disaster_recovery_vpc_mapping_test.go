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

type mockMetaBdrcVpcMapping struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBdrcVpcMapping) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBdrcVpcMapping{}

func newMockMetaBdrcVpcMapping() *mockMetaBdrcVpcMapping {
	return &mockMetaBdrcVpcMapping{client: &connectivity.TencentCloudClient{}}
}

func ptrUint64Bdrc(v uint64) *uint64 {
	return &v
}

func ptrInt64Bdrc(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/bdrc/ -run "TestBdrcDisasterRecoveryVpcMapping" -v -count=1 -gcflags="all=-l"

func TestBdrcDisasterRecoveryVpcMapping_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcVpcMapping().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "CreateDisasterRecoveryVpcMappingWithContext", func(ctx context.Context, request *bdrcv20260330.CreateDisasterRecoveryVpcMappingRequest) (*bdrcv20260330.CreateDisasterRecoveryVpcMappingResponse, error) {
		assert.Equal(t, "sp-001", *request.SitePairId)
		assert.Equal(t, "vpc-source", *request.SourceVpcId)
		assert.Equal(t, "subnet-source", *request.SourceSubnetId)
		assert.Equal(t, "vpc-target", *request.TargetVpcId)
		assert.Equal(t, "subnet-target", *request.TargetSubnetId)

		resp := bdrcv20260330.NewCreateDisasterRecoveryVpcMappingResponse()
		resp.Response = &bdrcv20260330.CreateDisasterRecoveryVpcMappingResponseParams{
			RequestId: ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(bdrcClient, "DescribeVpcMappingsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeVpcMappingsRequest) (*bdrcv20260330.DescribeVpcMappingsResponse, error) {
		assert.Equal(t, "sp-001", *request.SitePairId)
		assert.Equal(t, int64(100), *request.Limit)

		resp := bdrcv20260330.NewDescribeVpcMappingsResponse()
		resp.Response = &bdrcv20260330.DescribeVpcMappingsResponseParams{
			TotalCount: ptrInt64Bdrc(1),
			VpcMappingSet: []*bdrcv20260330.VpcMapping{
				{
					Id:           ptrUint64Bdrc(88),
					SitePairId:   ptrStringBdrc("sp-001"),
					SourceVpc:    ptrStringBdrc("vpc-source"),
					SourceSubnet: ptrStringBdrc("subnet-source"),
					TargetVpc:    ptrStringBdrc("vpc-target"),
					TargetSubnet: ptrStringBdrc("subnet-target"),
					Status:       ptrStringBdrc("ok"),
					LifeState:    ptrStringBdrc("NORMAL"),
				},
			},
			RequestId: ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcVpcMapping()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryVpcMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":     "sp-001",
		"source_vpc_id":    "vpc-source",
		"source_subnet_id": "subnet-source",
		"target_vpc_id":    "vpc-target",
		"target_subnet_id": "subnet-target",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sp-001#88", d.Id())
	assert.Equal(t, "vpc-source", d.Get("source_vpc").(string))
	assert.Equal(t, "ok", d.Get("status").(string))
	assert.Equal(t, "NORMAL", d.Get("life_state").(string))
}

func TestBdrcDisasterRecoveryVpcMapping_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcVpcMapping().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "DescribeVpcMappingsWithContext", func(ctx context.Context, request *bdrcv20260330.DescribeVpcMappingsRequest) (*bdrcv20260330.DescribeVpcMappingsResponse, error) {
		resp := bdrcv20260330.NewDescribeVpcMappingsResponse()
		resp.Response = &bdrcv20260330.DescribeVpcMappingsResponseParams{
			TotalCount:    ptrInt64Bdrc(0),
			VpcMappingSet: []*bdrcv20260330.VpcMapping{},
			RequestId:     ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcVpcMapping()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryVpcMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("sp-001#88")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestBdrcDisasterRecoveryVpcMapping_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcVpcMapping().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "DeleteDisasterRecoveryVpcMappingWithContext", func(ctx context.Context, request *bdrcv20260330.DeleteDisasterRecoveryVpcMappingRequest) (*bdrcv20260330.DeleteDisasterRecoveryVpcMappingResponse, error) {
		assert.Equal(t, 1, len(request.VpcMappingIds))
		assert.Equal(t, uint64(88), *request.VpcMappingIds[0])

		resp := bdrcv20260330.NewDeleteDisasterRecoveryVpcMappingResponse()
		resp.Response = &bdrcv20260330.DeleteDisasterRecoveryVpcMappingResponseParams{
			RequestId: ptrStringBdrc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcVpcMapping()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryVpcMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("sp-001#88")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestBdrcDisasterRecoveryVpcMapping_Update_Reject(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcVpcMapping().client, "UseBdrcV20260330Client", bdrcClient)

	meta := newMockMetaBdrcVpcMapping()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryVpcMapping()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":     "sp-001",
		"source_vpc_id":    "vpc-source",
		"source_subnet_id": "subnet-source",
		"target_vpc_id":    "vpc-target-changed",
		"target_subnet_id": "subnet-target",
	})
	d.SetId("sp-001#88")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CRD-only API")
}
