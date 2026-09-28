package dbdc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	dbdc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

type mockMetaDbdc struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDbdc) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDbdc{}

func newMockMetaDbdc() *mockMetaDbdc {
	return &mockMetaDbdc{client: &connectivity.TencentCloudClient{}}
}

func ptrStringDbdc(s string) *string {
	return &s
}

func ptrUint64Dbdc(u uint64) *uint64 {
	return &u
}

// TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_create tests the
// full Create -> Read flow with gomonkey mocks.
func TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdc().client, "UseDbdcV20201029Client", dbdcClient)

	// Mock CreateDBCustomDisasterRecoverGroupWithContext
	patches.ApplyMethodFunc(dbdcClient, "CreateDBCustomDisasterRecoverGroupWithContext", func(_ context.Context, request *dbdcv20201029.CreateDBCustomDisasterRecoverGroupRequest) (*dbdcv20201029.CreateDBCustomDisasterRecoverGroupResponse, error) {
		assert.NotNil(t, request.Name)
		assert.Equal(t, "tf-example", *request.Name)
		groupId := "dcg-test-123"
		resp := dbdcv20201029.NewCreateDBCustomDisasterRecoverGroupResponse()
		resp.Response = &dbdcv20201029.CreateDBCustomDisasterRecoverGroupResponseParams{
			DisasterRecoverGroupId: &groupId,
			RequestId:              ptrStringDbdc("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DbdcService.DescribeDBCustomDisasterRecoverGroupById for create polling + read
	patches.ApplyMethodFunc(&dbdc.DbdcService{}, "DescribeDBCustomDisasterRecoverGroupById", func(_ context.Context, disasterRecoverGroupId string) (*dbdcv20201029.DisasterRecoverGroup, error) {
		assert.Equal(t, "dcg-test-123", disasterRecoverGroupId)
		return &dbdcv20201029.DisasterRecoverGroup{
			DisasterRecoverGroupId: ptrStringDbdc("dcg-test-123"),
			Name:                   ptrStringDbdc("tf-example"),
			Type:                   ptrStringDbdc("HOST"),
			Status:                 ptrStringDbdc("Available"),
			Strategy:               ptrStringDbdc("SPREAD"),
			Affinity:               ptrUint64Dbdc(1),
			NodeQuotaTotal:         ptrUint64Dbdc(10),
			CurrentNum:             ptrUint64Dbdc(0),
			CreatedTime:            ptrStringDbdc("2024-01-01 00:00:00"),
			Tags: []*dbdcv20201029.Tag{
				{
					Key:   ptrStringDbdc("createBy"),
					Value: ptrStringDbdc("Terraform"),
				},
			},
		}, nil
	})

	meta := newMockMetaDbdc()
	res := dbdc.ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":     "tf-example",
		"type":     "HOST",
		"strategy": "SPREAD",
		"affinity": 1,
		"tags": []interface{}{
			map[string]interface{}{
				"key":   "createBy",
				"value": "Terraform",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "dcg-test-123", d.Id())
	assert.Equal(t, "tf-example", d.Get("name"))
	assert.Equal(t, "HOST", d.Get("type"))
	assert.Equal(t, "Available", d.Get("status"))
	assert.Equal(t, "SPREAD", d.Get("strategy"))
	assert.Equal(t, 1, d.Get("affinity"))
}

// TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_read tests the
// Read flow when the group no longer exists (should clear ID).
func TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DbdcService.DescribeDBCustomDisasterRecoverGroupById to return nil (not found)
	patches.ApplyMethodFunc(&dbdc.DbdcService{}, "DescribeDBCustomDisasterRecoverGroupById", func(_ context.Context, disasterRecoverGroupId string) (*dbdcv20201029.DisasterRecoverGroup, error) {
		return nil, nil
	})

	meta := newMockMetaDbdc()
	res := dbdc.ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "tf-example",
	})
	d.SetId("dcg-not-found")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_update tests the
// Update flow with name change.
func TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdc().client, "UseDbdcV20201029Client", dbdcClient)

	// Mock ModifyDBCustomDisasterRecoverGroupAttributeWithContext
	patches.ApplyMethodFunc(dbdcClient, "ModifyDBCustomDisasterRecoverGroupAttributeWithContext", func(_ context.Context, request *dbdcv20201029.ModifyDBCustomDisasterRecoverGroupAttributeRequest) (*dbdcv20201029.ModifyDBCustomDisasterRecoverGroupAttributeResponse, error) {
		assert.NotNil(t, request.DisasterRecoverGroupId)
		assert.Equal(t, "dcg-test-123", *request.DisasterRecoverGroupId)
		assert.NotNil(t, request.Name)
		assert.Equal(t, "tf-example-updated", *request.Name)
		resp := dbdcv20201029.NewModifyDBCustomDisasterRecoverGroupAttributeResponse()
		resp.Response = &dbdcv20201029.ModifyDBCustomDisasterRecoverGroupAttributeResponseParams{
			RequestId: ptrStringDbdc("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DbdcService.DescribeDBCustomDisasterRecoverGroupById for Read after update
	patches.ApplyMethodFunc(&dbdc.DbdcService{}, "DescribeDBCustomDisasterRecoverGroupById", func(_ context.Context, disasterRecoverGroupId string) (*dbdcv20201029.DisasterRecoverGroup, error) {
		return &dbdcv20201029.DisasterRecoverGroup{
			DisasterRecoverGroupId: ptrStringDbdc("dcg-test-123"),
			Name:                   ptrStringDbdc("tf-example-updated"),
			Type:                   ptrStringDbdc("HOST"),
			Status:                 ptrStringDbdc("Available"),
			Strategy:               ptrStringDbdc("SPREAD"),
			Affinity:               ptrUint64Dbdc(1),
		}, nil
	})

	meta := newMockMetaDbdc()
	res := dbdc.ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":     "tf-example-updated",
		"type":     "HOST",
		"strategy": "SPREAD",
		"affinity": 1,
	})
	d.SetId("dcg-test-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "tf-example-updated", d.Get("name"))
}

// TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_updateImmutable
// tests that changing an immutable field returns an error.
func TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_updateImmutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDbdc()
	res := dbdc.ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":     "tf-example",
		"type":     "HOST",
		"strategy": "SPREAD",
		"affinity": 1,
	})
	d.SetId("dcg-test-123")

	// Mock HasChange to simulate that type is being changed
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "type"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "argument `type` cannot be changed")
}

// TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_delete tests
// the Delete flow with async task polling.
func TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdc().client, "UseDbdcV20201029Client", dbdcClient)

	// Mock DeleteDBCustomDisasterRecoverGroupsWithContext
	patches.ApplyMethodFunc(dbdcClient, "DeleteDBCustomDisasterRecoverGroupsWithContext", func(_ context.Context, request *dbdcv20201029.DeleteDBCustomDisasterRecoverGroupsRequest) (*dbdcv20201029.DeleteDBCustomDisasterRecoverGroupsResponse, error) {
		assert.NotNil(t, request.DisasterRecoverGroupIds)
		assert.Len(t, request.DisasterRecoverGroupIds, 1)
		assert.Equal(t, "dcg-test-123", *request.DisasterRecoverGroupIds[0])
		taskId := uint64(1001)
		resp := dbdcv20201029.NewDeleteDBCustomDisasterRecoverGroupsResponse()
		resp.Response = &dbdcv20201029.DeleteDBCustomDisasterRecoverGroupsResponseParams{
			TaskId:    &taskId,
			RequestId: ptrStringDbdc("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeDBCustomTaskStatusWithContext for waitDBCustomTaskSucceeded
	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomTaskStatusWithContext", func(_ context.Context, request *dbdcv20201029.DescribeDBCustomTaskStatusRequest) (*dbdcv20201029.DescribeDBCustomTaskStatusResponse, error) {
		status := "Succeeded"
		resp := dbdcv20201029.NewDescribeDBCustomTaskStatusResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomTaskStatusResponseParams{
			Status:    &status,
			RequestId: ptrStringDbdc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaDbdc()
	res := dbdc.ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "tf-example",
	})
	d.SetId("dcg-test-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_createFailed
// tests that Create returns an error when the group status is CreateFailed.
func TestAccTencentCloudDbdcDbCustomDisasterRecoverGroupResource_createFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdc().client, "UseDbdcV20201029Client", dbdcClient)

	// Mock CreateDBCustomDisasterRecoverGroupWithContext
	patches.ApplyMethodFunc(dbdcClient, "CreateDBCustomDisasterRecoverGroupWithContext", func(_ context.Context, request *dbdcv20201029.CreateDBCustomDisasterRecoverGroupRequest) (*dbdcv20201029.CreateDBCustomDisasterRecoverGroupResponse, error) {
		groupId := "dcg-test-failed"
		resp := dbdcv20201029.NewCreateDBCustomDisasterRecoverGroupResponse()
		resp.Response = &dbdcv20201029.CreateDBCustomDisasterRecoverGroupResponseParams{
			DisasterRecoverGroupId: &groupId,
			RequestId:              ptrStringDbdc("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DbdcService.DescribeDBCustomDisasterRecoverGroupById - return CreateFailed status
	patches.ApplyMethodFunc(&dbdc.DbdcService{}, "DescribeDBCustomDisasterRecoverGroupById", func(_ context.Context, disasterRecoverGroupId string) (*dbdcv20201029.DisasterRecoverGroup, error) {
		return &dbdcv20201029.DisasterRecoverGroup{
			DisasterRecoverGroupId: ptrStringDbdc("dcg-test-failed"),
			Name:                   ptrStringDbdc("tf-example"),
			Status:                 ptrStringDbdc("CreateFailed"),
		}, nil
	})

	meta := newMockMetaDbdc()
	res := dbdc.ResourceTencentCloudDbdcDbCustomDisasterRecoverGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "tf-example",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CreateFailed")
}
