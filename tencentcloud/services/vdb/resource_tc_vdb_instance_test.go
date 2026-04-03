package vdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	vdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vdb/v20230616"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// Mock VDB Client
// ============================================================

type MockVdbClient struct {
	mock.Mock
}

func (m *MockVdbClient) DescribeInstancesWithContext(ctx context.Context, request *vdb.DescribeInstancesRequest) (*vdb.DescribeInstancesResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.DescribeInstancesResponse), args.Error(1)
}

func (m *MockVdbClient) DescribeInstanceNodesWithContext(ctx context.Context, request *vdb.DescribeInstanceNodesRequest) (*vdb.DescribeInstanceNodesResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.DescribeInstanceNodesResponse), args.Error(1)
}

func (m *MockVdbClient) DescribeDBSecurityGroupsWithContext(ctx context.Context, request *vdb.DescribeDBSecurityGroupsRequest) (*vdb.DescribeDBSecurityGroupsResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.DescribeDBSecurityGroupsResponse), args.Error(1)
}

func (m *MockVdbClient) DescribeInstanceMaintenanceWindowWithContext(ctx context.Context, request *vdb.DescribeInstanceMaintenanceWindowRequest) (*vdb.DescribeInstanceMaintenanceWindowResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.DescribeInstanceMaintenanceWindowResponse), args.Error(1)
}

func (m *MockVdbClient) CreateInstanceWithContext(ctx context.Context, request *vdb.CreateInstanceRequest) (*vdb.CreateInstanceResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.CreateInstanceResponse), args.Error(1)
}

func (m *MockVdbClient) ScaleUpInstanceWithContext(ctx context.Context, request *vdb.ScaleUpInstanceRequest) (*vdb.ScaleUpInstanceResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.ScaleUpInstanceResponse), args.Error(1)
}

func (m *MockVdbClient) ScaleOutInstanceWithContext(ctx context.Context, request *vdb.ScaleOutInstanceRequest) (*vdb.ScaleOutInstanceResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.ScaleOutInstanceResponse), args.Error(1)
}

func (m *MockVdbClient) IsolateInstanceWithContext(ctx context.Context, request *vdb.IsolateInstanceRequest) (*vdb.IsolateInstanceResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.IsolateInstanceResponse), args.Error(1)
}

func (m *MockVdbClient) DestroyInstancesWithContext(ctx context.Context, request *vdb.DestroyInstancesRequest) (*vdb.DestroyInstancesResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.DestroyInstancesResponse), args.Error(1)
}

func (m *MockVdbClient) AssociateSecurityGroupsWithContext(ctx context.Context, request *vdb.AssociateSecurityGroupsRequest) (*vdb.AssociateSecurityGroupsResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.AssociateSecurityGroupsResponse), args.Error(1)
}

func (m *MockVdbClient) DisassociateSecurityGroupsWithContext(ctx context.Context, request *vdb.DisassociateSecurityGroupsRequest) (*vdb.DisassociateSecurityGroupsResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.DisassociateSecurityGroupsResponse), args.Error(1)
}

func (m *MockVdbClient) ModifyDBInstanceSecurityGroupsWithContext(ctx context.Context, request *vdb.ModifyDBInstanceSecurityGroupsRequest) (*vdb.ModifyDBInstanceSecurityGroupsResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vdb.ModifyDBInstanceSecurityGroupsResponse), args.Error(1)
}

// ============================================================
// Helper functions
// ============================================================

func newTestService(mockClient *MockVdbClient) *VdbService {
	return &VdbService{vdbClient: mockClient}
}

func strPtr(s string) *string       { return &s }
func uint64Ptr(v uint64) *uint64    { return &v }
func float64Ptr(v float64) *float64 { return &v }
func int64Ptr(v int64) *int64       { return &v }

func makeInstanceInfo(id, name, status string, cpu float64, memory float64, disk uint64, replicaNum uint64) *vdb.InstanceInfo {
	return &vdb.InstanceInfo{
		InstanceId: strPtr(id),
		Name:       strPtr(name),
		Status:     strPtr(status),
		Cpu:        float64Ptr(cpu),
		Memory:     float64Ptr(memory),
		Disk:       uint64Ptr(disk),
		ReplicaNum: uint64Ptr(replicaNum),
		Region:     strPtr("ap-guangzhou"),
		Zone:       strPtr("ap-guangzhou-3"),
		CreatedAt:  strPtr("2024-01-01 00:00:00"),
		PayMode:    int64Ptr(0),
	}
}

func makeDescribeResponse(items []*vdb.InstanceInfo) *vdb.DescribeInstancesResponse {
	return &vdb.DescribeInstancesResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeInstancesResponseParams{
			Items:      items,
			TotalCount: uint64Ptr(uint64(len(items))),
		},
	}
}

// ============================================================
// Tests: DescribeVdbInstanceById
// ============================================================

func TestDescribeVdbInstanceById_Found(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	instance := makeInstanceInfo("vdb-test123", "my-instance", "online", 4, 8, 100, 2)
	resp := makeDescribeResponse([]*vdb.InstanceInfo{instance})

	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	result, err := service.DescribeVdbInstanceById(context.Background(), "vdb-test123")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "vdb-test123", *result.InstanceId)
	assert.Equal(t, "my-instance", *result.Name)
	assert.Equal(t, "online", *result.Status)
	assert.Equal(t, float64(4), *result.Cpu)
	assert.Equal(t, float64(8), *result.Memory)
	assert.Equal(t, uint64(100), *result.Disk)
	assert.Equal(t, uint64(2), *result.ReplicaNum)
	mockClient.AssertExpectations(t)
}

func TestDescribeVdbInstanceById_NotFound(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	resp := makeDescribeResponse([]*vdb.InstanceInfo{})
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	result, err := service.DescribeVdbInstanceById(context.Background(), "vdb-nonexist")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestDescribeVdbInstanceById_NilItems(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	resp := &vdb.DescribeInstancesResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeInstancesResponseParams{
			Items: nil,
		},
	}
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	result, err := service.DescribeVdbInstanceById(context.Background(), "vdb-test")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestDescribeVdbInstanceById_ApiError(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(nil, fmt.Errorf("api error: InternalError"))

	result, err := service.DescribeVdbInstanceById(context.Background(), "vdb-test")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestDescribeVdbInstanceById_NilResponse(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	resp := &vdb.DescribeInstancesResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response:     nil,
	}
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	result, err := service.DescribeVdbInstanceById(context.Background(), "vdb-test")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "response is nil")
	mockClient.AssertExpectations(t)
}

// ============================================================
// Tests: DescribeVdbInstanceNodesById
// ============================================================

func TestDescribeVdbInstanceNodesById_Found(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	nodes := []*vdb.NodeInfo{
		{Name: strPtr("pod-1"), Status: strPtr("Running")},
		{Name: strPtr("pod-2"), Status: strPtr("Running")},
	}
	resp := &vdb.DescribeInstanceNodesResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeInstanceNodesResponseParams{
			Items:      nodes,
			TotalCount: uint64Ptr(2),
		},
	}
	mockClient.On("DescribeInstanceNodesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	result, err := service.DescribeVdbInstanceNodesById(context.Background(), "vdb-test123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "pod-1", *result[0].Name)
	assert.Equal(t, "pod-2", *result[1].Name)
	mockClient.AssertExpectations(t)
}

func TestDescribeVdbInstanceNodesById_Empty(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	resp := &vdb.DescribeInstanceNodesResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeInstanceNodesResponseParams{
			Items: nil,
		},
	}
	mockClient.On("DescribeInstanceNodesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	result, err := service.DescribeVdbInstanceNodesById(context.Background(), "vdb-test")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

// ============================================================
// Tests: DescribeDBSecurityGroupsByInstanceId
// ============================================================

func TestDescribeDBSecurityGroupsByInstanceId_Found(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	groups := []*vdb.SecurityGroup{
		{SecurityGroupId: strPtr("sg-111"), SecurityGroupName: strPtr("default")},
		{SecurityGroupId: strPtr("sg-222"), SecurityGroupName: strPtr("web")},
	}
	resp := &vdb.DescribeDBSecurityGroupsResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeDBSecurityGroupsResponseParams{
			Groups: groups,
		},
	}
	mockClient.On("DescribeDBSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	result, err := service.DescribeDBSecurityGroupsByInstanceId(context.Background(), "vdb-test")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "sg-111", *result[0].SecurityGroupId)
	assert.Equal(t, "sg-222", *result[1].SecurityGroupId)
	mockClient.AssertExpectations(t)
}

func TestDescribeDBSecurityGroupsByInstanceId_Empty(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	resp := &vdb.DescribeDBSecurityGroupsResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeDBSecurityGroupsResponseParams{
			Groups: nil,
		},
	}
	mockClient.On("DescribeDBSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	result, err := service.DescribeDBSecurityGroupsByInstanceId(context.Background(), "vdb-test")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

// ============================================================
// Tests: WaitForInstanceStatus
// ============================================================

func TestWaitForInstanceStatus_AlreadyOnline(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	instance := makeInstanceInfo("vdb-test", "test", "online", 2, 4, 50, 1)
	resp := makeDescribeResponse([]*vdb.InstanceInfo{instance})
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	err := service.WaitForInstanceStatus(context.Background(), "vdb-test", "online", 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForInstanceStatus_CaseInsensitive(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	instance := makeInstanceInfo("vdb-test", "test", "Online", 2, 4, 50, 1)
	resp := makeDescribeResponse([]*vdb.InstanceInfo{instance})
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	err := service.WaitForInstanceStatus(context.Background(), "vdb-test", "online", 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForInstanceStatus_IsolatedNotFound(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	resp := makeDescribeResponse([]*vdb.InstanceInfo{})
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	err := service.WaitForInstanceStatus(context.Background(), "vdb-test", "isolated", 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForInstanceStatus_DestroyedNotFound(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	resp := makeDescribeResponse([]*vdb.InstanceInfo{})
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	err := service.WaitForInstanceStatus(context.Background(), "vdb-test", "destroyed", 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// ============================================================
// Tests: WaitForInstanceScaleUp
// ============================================================

func TestWaitForInstanceScaleUp_AlreadyMatched(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	instance := makeInstanceInfo("vdb-test", "test", "online", 4, 8, 100, 2)
	resp := makeDescribeResponse([]*vdb.InstanceInfo{instance})
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	err := service.WaitForInstanceScaleUp(context.Background(), "vdb-test", 4, 8, 100, 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForInstanceScaleUp_ValuesConverge(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	// First call: old values
	oldInstance := makeInstanceInfo("vdb-test", "test", "online", 2, 4, 50, 1)
	oldResp := makeDescribeResponse([]*vdb.InstanceInfo{oldInstance})

	// Second call: new values
	newInstance := makeInstanceInfo("vdb-test", "test", "online", 4, 8, 100, 1)
	newResp := makeDescribeResponse([]*vdb.InstanceInfo{newInstance})

	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(oldResp, nil).Once()
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(newResp, nil)

	err := service.WaitForInstanceScaleUp(context.Background(), "vdb-test", 4, 8, 100, 30*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForInstanceScaleUp_NilCpuMemoryDisk(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	// Instance with nil cpu/memory/disk values
	instance := &vdb.InstanceInfo{
		InstanceId: strPtr("vdb-test"),
		Status:     strPtr("online"),
		Cpu:        nil,
		Memory:     nil,
		Disk:       nil,
	}
	resp := makeDescribeResponse([]*vdb.InstanceInfo{instance})

	// Second call with correct values
	newInstance := makeInstanceInfo("vdb-test", "test", "online", 0, 0, 0, 1)
	newResp := makeDescribeResponse([]*vdb.InstanceInfo{newInstance})

	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(resp, nil).Once()
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(newResp, nil)

	// target is 0,0,0 so nil values (treated as 0) should match on first call
	err := service.WaitForInstanceScaleUp(context.Background(), "vdb-test", 0, 0, 0, 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// ============================================================
// Tests: WaitForInstanceScaleOut
// ============================================================

func TestWaitForInstanceScaleOut_AlreadyMatched(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	instance := makeInstanceInfo("vdb-test", "test", "online", 2, 4, 50, 3)
	resp := makeDescribeResponse([]*vdb.InstanceInfo{instance})
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	err := service.WaitForInstanceScaleOut(context.Background(), "vdb-test", 3, 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForInstanceScaleOut_ValuesConverge(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	oldInstance := makeInstanceInfo("vdb-test", "test", "online", 2, 4, 50, 1)
	oldResp := makeDescribeResponse([]*vdb.InstanceInfo{oldInstance})

	newInstance := makeInstanceInfo("vdb-test", "test", "online", 2, 4, 50, 3)
	newResp := makeDescribeResponse([]*vdb.InstanceInfo{newInstance})

	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(oldResp, nil).Once()
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(newResp, nil)

	err := service.WaitForInstanceScaleOut(context.Background(), "vdb-test", 3, 30*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// ============================================================
// Tests: Resource Schema Validation
// ============================================================

func TestVdbInstanceResourceSchema(t *testing.T) {
	r := ResourceTencentCloudVdbInstance()

	// Verify all required fields exist
	requiredFields := []string{"vpc_id", "subnet_id", "pay_mode", "security_group_ids"}
	for _, field := range requiredFields {
		s, ok := r.Schema[field]
		assert.True(t, ok, "schema should contain field: %s", field)
		assert.True(t, s.Required, "field %s should be Required", field)
	}

	// Verify optional input fields
	optionalFields := []string{
		"instance_name", "pay_period", "auto_renew",
		"params", "resource_tags", "instance_type", "mode", "goods_num",
		"product_type", "node_type", "cpu", "memory", "disk_size",
		"worker_node_num", "force_delete",
	}
	for _, field := range optionalFields {
		s, ok := r.Schema[field]
		assert.True(t, ok, "schema should contain field: %s", field)
		assert.True(t, s.Optional, "field %s should be Optional", field)
	}

	// Verify computed fields
	computedFields := []string{
		"status", "region", "zone", "product", "shard_num", "api_version",
		"extend", "expired_at", "is_no_expired", "wan_address", "isolate_at",
		"task_status", "networks", "nodes", "created_at", "engine_name",
		"engine_version",
	}
	for _, field := range computedFields {
		s, ok := r.Schema[field]
		assert.True(t, ok, "schema should contain computed field: %s", field)
		assert.True(t, s.Computed, "field %s should be Computed", field)
	}

	// Verify ForceNew fields
	forceNewFields := []string{
		"vpc_id", "subnet_id", "instance_type",
		"mode", "product_type", "node_type",
	}
	for _, field := range forceNewFields {
		s, ok := r.Schema[field]
		assert.True(t, ok, "schema should contain field: %s", field)
		assert.True(t, s.ForceNew, "field %s should be ForceNew", field)
	}

	// Verify cpu, memory, disk_size, worker_node_num, security_group_ids are NOT ForceNew (they support Update)
	updateableFields := []string{"cpu", "memory", "disk_size", "worker_node_num", "security_group_ids"}
	for _, field := range updateableFields {
		s, ok := r.Schema[field]
		assert.True(t, ok, "schema should contain field: %s", field)
		assert.False(t, s.ForceNew, "field %s should NOT be ForceNew (supports Update)", field)
	}

	// Verify force_delete default
	fd := r.Schema["force_delete"]
	assert.Equal(t, false, fd.Default, "force_delete should default to false")

	// Verify CRUD functions
	assert.NotNil(t, r.Create)
	assert.NotNil(t, r.Read)
	assert.NotNil(t, r.Update)
	assert.NotNil(t, r.Delete)
	assert.NotNil(t, r.Importer)

	// Verify Timeouts
	assert.NotNil(t, r.Timeouts)
	assert.NotNil(t, r.Timeouts.Create)
	assert.NotNil(t, r.Timeouts.Update)
	assert.NotNil(t, r.Timeouts.Delete)
}

// ============================================================
// Tests: WaitForSecurityGroupsMatch
// ============================================================

func TestWaitForSecurityGroupsMatch_AlreadyMatched(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	groups := []*vdb.SecurityGroup{
		{SecurityGroupId: strPtr("sg-111")},
		{SecurityGroupId: strPtr("sg-222")},
	}
	resp := &vdb.DescribeDBSecurityGroupsResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeDBSecurityGroupsResponseParams{
			Groups: groups,
		},
	}
	mockClient.On("DescribeDBSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	err := service.WaitForSecurityGroupsMatch(context.Background(), "vdb-test", []string{"sg-111", "sg-222"}, 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForSecurityGroupsMatch_EventuallyMatch(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	// First call: old groups
	oldGroups := []*vdb.SecurityGroup{
		{SecurityGroupId: strPtr("sg-old")},
	}
	oldResp := &vdb.DescribeDBSecurityGroupsResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeDBSecurityGroupsResponseParams{
			Groups: oldGroups,
		},
	}

	// Second call: new groups
	newGroups := []*vdb.SecurityGroup{
		{SecurityGroupId: strPtr("sg-new1")},
		{SecurityGroupId: strPtr("sg-new2")},
	}
	newResp := &vdb.DescribeDBSecurityGroupsResponse{
		BaseResponse: &tchttp.BaseResponse{},
		Response: &vdb.DescribeDBSecurityGroupsResponseParams{
			Groups: newGroups,
		},
	}

	mockClient.On("DescribeDBSecurityGroupsWithContext", mock.Anything, mock.Anything).
		Return(oldResp, nil).Once()
	mockClient.On("DescribeDBSecurityGroupsWithContext", mock.Anything, mock.Anything).
		Return(newResp, nil)

	err := service.WaitForSecurityGroupsMatch(context.Background(), "vdb-test", []string{"sg-new1", "sg-new2"}, 30*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// ============================================================
// Tests: WaitForInstanceNotFound
// ============================================================

func TestWaitForInstanceNotFound_AlreadyGone(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	resp := makeDescribeResponse([]*vdb.InstanceInfo{})
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).Return(resp, nil)

	err := service.WaitForInstanceNotFound(context.Background(), "vdb-test", 10*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForInstanceNotFound_EventuallyGone(t *testing.T) {
	mockClient := new(MockVdbClient)
	service := newTestService(mockClient)

	// First call: still exists
	instance := makeInstanceInfo("vdb-test", "test", "isolated", 2, 4, 50, 1)
	existsResp := makeDescribeResponse([]*vdb.InstanceInfo{instance})

	// Second call: gone
	goneResp := makeDescribeResponse([]*vdb.InstanceInfo{})

	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(existsResp, nil).Once()
	mockClient.On("DescribeInstancesWithContext", mock.Anything, mock.Anything).
		Return(goneResp, nil)

	err := service.WaitForInstanceNotFound(context.Background(), "vdb-test", 30*time.Second)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// ============================================================
// Tests: Network sub-fields in schema
// ============================================================

func TestVdbInstanceNetworkSubFields(t *testing.T) {
	r := ResourceTencentCloudVdbInstance()

	networksSchema := r.Schema["networks"]
	assert.NotNil(t, networksSchema)

	elemResource := networksSchema.Elem.(*schema.Resource)
	expectedSubFields := []string{"vpc_id", "subnet_id", "vip", "port", "preserve_duration", "expire_time"}
	for _, field := range expectedSubFields {
		s, ok := elemResource.Schema[field]
		assert.True(t, ok, "networks should contain sub-field: %s", field)
		assert.True(t, s.Computed, "networks.%s should be Computed", field)
	}
}

func TestVdbInstanceNodesSubFields(t *testing.T) {
	r := ResourceTencentCloudVdbInstance()

	nodesSchema := r.Schema["nodes"]
	assert.NotNil(t, nodesSchema)

	elemResource := nodesSchema.Elem.(*schema.Resource)
	for _, field := range []string{"name", "status"} {
		s, ok := elemResource.Schema[field]
		assert.True(t, ok, "nodes should contain sub-field: %s", field)
		assert.True(t, s.Computed, "nodes.%s should be Computed", field)
	}
}
