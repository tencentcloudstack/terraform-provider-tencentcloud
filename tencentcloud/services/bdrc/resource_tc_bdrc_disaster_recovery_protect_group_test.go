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

type mockMetaBdrcProtectGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBdrcProtectGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBdrcProtectGroup{}

func newMockMetaBdrcProtectGroup() *mockMetaBdrcProtectGroup {
	return &mockMetaBdrcProtectGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringBdrcPg(s string) *string { return &s }
func ptrInt64BdrcPg(v int64) *int64    { return &v }
func ptrUint64BdrcPg(v uint64) *uint64 { return &v }

// buildProtectGroupResponse builds a canned DescribeDisasterRecoveryProtectGroups response
// containing a single ProtectGroup that matches the given id.
func buildProtectGroupResponse(id string) *bdrcv20260330.DescribeDisasterRecoveryProtectGroupsResponse {
	resp := bdrcv20260330.NewDescribeDisasterRecoveryProtectGroupsResponse()
	resp.Response = &bdrcv20260330.DescribeDisasterRecoveryProtectGroupsResponseParams{
		TotalCount: ptrInt64BdrcPg(1),
		ProtectGroupSet: []*bdrcv20260330.ProtectGroup{
			{
				ProtectGroupId:                   ptrStringBdrcPg(id),
				ProtectGroupName:                 ptrStringBdrcPg("tf-example-protect-group"),
				ProtectGroupType:                 ptrStringBdrcPg("DISK"),
				SitePairId:                       ptrStringBdrcPg("sitepair-xxxxxxxx"),
				SitePairName:                     ptrStringBdrcPg("tf-example-site-pair"),
				RecoveryPointObjective:           ptrInt64BdrcPg(900),
				SourceRegion:                     ptrStringBdrcPg("ap-guangzhou"),
				SourceZone:                       ptrStringBdrcPg("ap-guangzhou-3"),
				SourceVpc:                        ptrStringBdrcPg("vpc-source-xxx"),
				TargetRegion:                     ptrStringBdrcPg("ap-shanghai"),
				TargetZone:                       ptrStringBdrcPg("ap-shanghai-2"),
				TargetVpc:                        ptrStringBdrcPg("vpc-target-xxx"),
				CopyType:                         ptrStringBdrcPg("ASY"),
				DisasterRecoveryType:             ptrStringBdrcPg("CROSS_REGION"),
				DataDirection:                    ptrStringBdrcPg("POSITIVE"),
				PeerCloudName:                    ptrStringBdrcPg(""),
				CreateFrom:                       ptrStringBdrcPg("LOCAL"),
				LifeState:                        ptrStringBdrcPg("NORMAL"),
				AccountUin:                       ptrStringBdrcPg("100000000001"),
				SubAccountUin:                    ptrStringBdrcPg("200000000001"),
				CreateTime:                       ptrStringBdrcPg("2026-09-15T10:00:00Z"),
				ModifyTime:                       ptrStringBdrcPg("2026-09-15T11:00:00Z"),
				BindProtectedResourceCount:       ptrInt64BdrcPg(2),
				ErrorRecoveryPointObjectiveCount: ptrInt64BdrcPg(0),
				ProtectedResourceStatusSet: []*bdrcv20260330.ProtectedResourceStatus{
					{
						Status: ptrStringBdrcPg("AVAILABLE"),
						Count:  ptrUint64BdrcPg(2),
					},
				},
				AppId: ptrInt64BdrcPg(1250000001),
			},
		},
		RequestId: ptrStringBdrcPg("fake-request-id"),
	}
	return resp
}

// TestBdrcDisasterRecoveryProtectGroup_Create verifies the Create flow sets the id from the mock response.
func TestBdrcDisasterRecoveryProtectGroup_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcProtectGroup().client, "UseBdrcV20260330Client", bdrcClient)

	var capturedRequest *bdrcv20260330.CreateDisasterRecoveryProtectGroupRequest
	patches.ApplyMethodFunc(bdrcClient, "CreateDisasterRecoveryProtectGroupWithContext", func(_ context.Context, request *bdrcv20260330.CreateDisasterRecoveryProtectGroupRequest) (*bdrcv20260330.CreateDisasterRecoveryProtectGroupResponse, error) {
		capturedRequest = request
		assert.Equal(t, "sitepair-xxxxxxxx", *request.SitePairId)
		assert.Equal(t, "DISK", *request.ProtectGroupType)
		assert.Equal(t, int64(15), *request.RecoveryPointObjective)
		assert.Equal(t, "tf-example-protect-group", *request.ProtectGroupName)
		assert.Equal(t, "POSITIVE", *request.DataDirection)

		resp := bdrcv20260330.NewCreateDisasterRecoveryProtectGroupResponse()
		resp.Response = &bdrcv20260330.CreateDisasterRecoveryProtectGroupResponseParams{
			ProtectGroupId: ptrStringBdrcPg("pg-abc123"),
			RequestId:      ptrStringBdrcPg("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(bdrcClient, "DescribeDisasterRecoveryProtectGroupsWithContext", func(_ context.Context, request *bdrcv20260330.DescribeDisasterRecoveryProtectGroupsRequest) (*bdrcv20260330.DescribeDisasterRecoveryProtectGroupsResponse, error) {
		return buildProtectGroupResponse("pg-abc123"), nil
	})

	meta := newMockMetaBdrcProtectGroup()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sitepair-xxxxxxxx",
		"protect_group_type":       "DISK",
		"recovery_point_objective": 15,
		"protect_group_name":       "tf-example-protect-group",
		"data_direction":           "POSITIVE",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "pg-abc123#DISK", d.Id())
	assert.NotNil(t, capturedRequest)
}

// TestBdrcDisasterRecoveryProtectGroup_Read verifies Read populates computed fields from the response.
func TestBdrcDisasterRecoveryProtectGroup_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcProtectGroup().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "DescribeDisasterRecoveryProtectGroupsWithContext", func(_ context.Context, request *bdrcv20260330.DescribeDisasterRecoveryProtectGroupsRequest) (*bdrcv20260330.DescribeDisasterRecoveryProtectGroupsResponse, error) {
		assert.NotNil(t, request.ProtectGroupIds)
		assert.Equal(t, "pg-read123", *request.ProtectGroupIds[0])
		return buildProtectGroupResponse("pg-read123"), nil
	})

	meta := newMockMetaBdrcProtectGroup()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sitepair-xxxxxxxx",
		"protect_group_type":       "DISK",
		"recovery_point_objective": 15,
		"protect_group_name":       "tf-example-protect-group",
		"data_direction":           "POSITIVE",
	})
	d.SetId("pg-read123#DISK")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "pg-read123#DISK", d.Id())

	assert.Equal(t, "tf-example-protect-group", d.Get("protect_group_name"))
	assert.Equal(t, "DISK", d.Get("protect_group_type"))
	assert.Equal(t, "NORMAL", d.Get("life_state"))
	assert.Equal(t, "CROSS_REGION", d.Get("disaster_recovery_type"))
	assert.Equal(t, int(2), d.Get("bind_protected_resource_count"))

	statusSet := d.Get("protected_resource_status_set").([]interface{})
	assert.Len(t, statusSet, 1)
	statusMap := statusSet[0].(map[string]interface{})
	assert.Equal(t, "AVAILABLE", statusMap["status"])
	assert.Equal(t, int(2), statusMap["count"])
}

// TestBdrcDisasterRecoveryProtectGroup_Update_Rename verifies Update triggers ModifyProtectGroupAttribute with the new name.
func TestBdrcDisasterRecoveryProtectGroup_Update_Rename(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcProtectGroup().client, "UseBdrcV20260330Client", bdrcClient)

	var capturedModifyRequest *bdrcv20260330.ModifyProtectGroupAttributeRequest
	patches.ApplyMethodFunc(bdrcClient, "ModifyProtectGroupAttributeWithContext", func(_ context.Context, request *bdrcv20260330.ModifyProtectGroupAttributeRequest) (*bdrcv20260330.ModifyProtectGroupAttributeResponse, error) {
		capturedModifyRequest = request
		assert.Equal(t, "pg-update123", *request.ProtectGroupId)
		assert.Equal(t, "tf-example-protect-group-renamed", *request.ProtectGroupName)

		resp := bdrcv20260330.NewModifyProtectGroupAttributeResponse()
		resp.Response = &bdrcv20260330.ModifyProtectGroupAttributeResponseParams{
			RequestId: ptrStringBdrcPg("fake-request-id-modify"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(bdrcClient, "DescribeDisasterRecoveryProtectGroupsWithContext", func(_ context.Context, request *bdrcv20260330.DescribeDisasterRecoveryProtectGroupsRequest) (*bdrcv20260330.DescribeDisasterRecoveryProtectGroupsResponse, error) {
		return buildProtectGroupResponse("pg-update123"), nil
	})

	meta := newMockMetaBdrcProtectGroup()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sitepair-xxxxxxxx",
		"protect_group_type":       "DISK",
		"recovery_point_objective": 15,
		"protect_group_name":       "tf-example-protect-group-renamed",
		"data_direction":           "POSITIVE",
	})
	d.SetId("pg-update123#DISK")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "protect_group_name"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
}

// TestBdrcDisasterRecoveryProtectGroup_Update_ImmutableArg verifies Update rejects an immutable field change.
func TestBdrcDisasterRecoveryProtectGroup_Update_ImmutableArg(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcProtectGroup().client, "UseBdrcV20260330Client", bdrcClient)

	meta := newMockMetaBdrcProtectGroup()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sitepair-xxxxxxxx",
		"protect_group_type":       "INSTANCE",
		"recovery_point_objective": 15,
		"protect_group_name":       "tf-example-protect-group",
		"data_direction":           "POSITIVE",
	})
	d.SetId("pg-immutable123#INSTANCE")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "protect_group_type"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "protect_group_type")
	assert.Contains(t, err.Error(), "immutable")
}

// TestBdrcDisasterRecoveryProtectGroup_Delete verifies Delete issues DeleteDisasterRecoveryProtectGroups with the id.
func TestBdrcDisasterRecoveryProtectGroup_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcProtectGroup().client, "UseBdrcV20260330Client", bdrcClient)

	var capturedRequest *bdrcv20260330.DeleteDisasterRecoveryProtectGroupsRequest
	patches.ApplyMethodFunc(bdrcClient, "DeleteDisasterRecoveryProtectGroupsWithContext", func(_ context.Context, request *bdrcv20260330.DeleteDisasterRecoveryProtectGroupsRequest) (*bdrcv20260330.DeleteDisasterRecoveryProtectGroupsResponse, error) {
		capturedRequest = request
		assert.NotNil(t, request.ProtectGroups)
		assert.Equal(t, 1, len(request.ProtectGroups))
		assert.Equal(t, "pg-del123", *request.ProtectGroups[0])

		resp := bdrcv20260330.NewDeleteDisasterRecoveryProtectGroupsResponse()
		resp.Response = &bdrcv20260330.DeleteDisasterRecoveryProtectGroupsResponseParams{
			RequestId: ptrStringBdrcPg("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcProtectGroup()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sitepair-xxxxxxxx",
		"protect_group_type":       "DISK",
		"recovery_point_objective": 15,
		"protect_group_name":       "tf-example-protect-group",
		"data_direction":           "POSITIVE",
	})
	d.SetId("pg-del123#DISK")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
}

// TestBdrcDisasterRecoveryProtectGroup_Read_NotFound verifies Read clears state when the describe returns empty.
func TestBdrcDisasterRecoveryProtectGroup_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bdrcClient := &bdrcv20260330.Client{}
	patches.ApplyMethodReturn(newMockMetaBdrcProtectGroup().client, "UseBdrcV20260330Client", bdrcClient)

	patches.ApplyMethodFunc(bdrcClient, "DescribeDisasterRecoveryProtectGroupsWithContext", func(_ context.Context, request *bdrcv20260330.DescribeDisasterRecoveryProtectGroupsRequest) (*bdrcv20260330.DescribeDisasterRecoveryProtectGroupsResponse, error) {
		resp := bdrcv20260330.NewDescribeDisasterRecoveryProtectGroupsResponse()
		resp.Response = &bdrcv20260330.DescribeDisasterRecoveryProtectGroupsResponseParams{
			TotalCount:      ptrInt64BdrcPg(0),
			ProtectGroupSet: []*bdrcv20260330.ProtectGroup{},
			RequestId:       ptrStringBdrcPg("fake-request-id-empty"),
		}
		return resp, nil
	})

	meta := newMockMetaBdrcProtectGroup()
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"site_pair_id":             "sitepair-xxxxxxxx",
		"protect_group_type":       "DISK",
		"recovery_point_objective": 15,
		"protect_group_name":       "tf-example-protect-group",
		"data_direction":           "POSITIVE",
	})
	d.SetId("pg-missing#DISK")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestBdrcDisasterRecoveryProtectGroup_Schema validates the schema definition.
func TestBdrcDisasterRecoveryProtectGroup_Schema(t *testing.T) {
	res := bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "site_pair_id")
	assert.Contains(t, res.Schema, "protect_group_type")
	assert.Contains(t, res.Schema, "recovery_point_objective")
	assert.Contains(t, res.Schema, "protect_group_name")
	assert.Contains(t, res.Schema, "data_direction")

	sitePairId := res.Schema["site_pair_id"]
	assert.Equal(t, schema.TypeString, sitePairId.Type)
	assert.True(t, sitePairId.Required)
	assert.True(t, sitePairId.ForceNew)

	protectGroupType := res.Schema["protect_group_type"]
	assert.Equal(t, schema.TypeString, protectGroupType.Type)
	assert.True(t, protectGroupType.Required)
	assert.True(t, protectGroupType.ForceNew)

	recoveryPointObjective := res.Schema["recovery_point_objective"]
	assert.Equal(t, schema.TypeInt, recoveryPointObjective.Type)
	assert.True(t, recoveryPointObjective.Required)
	assert.True(t, recoveryPointObjective.ForceNew)

	protectGroupName := res.Schema["protect_group_name"]
	assert.Equal(t, schema.TypeString, protectGroupName.Type)
	assert.True(t, protectGroupName.Optional)
	assert.False(t, protectGroupName.ForceNew)

	dataDirection := res.Schema["data_direction"]
	assert.Equal(t, schema.TypeString, dataDirection.Type)
	assert.True(t, dataDirection.Optional)
	assert.True(t, dataDirection.ForceNew)

	statusSet := res.Schema["protected_resource_status_set"]
	assert.Equal(t, schema.TypeList, statusSet.Type)
	assert.True(t, statusSet.Computed)

	assert.NotContains(t, res.Schema, "protect_group_set")
	assert.NotContains(t, res.Schema, "protect_group_list")
}
