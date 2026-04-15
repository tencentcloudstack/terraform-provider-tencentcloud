package teo_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// MockClient is a mock for the TEO API client
type MockTeoV20220901Client struct {
	mock.Mock
}

// CreateCLSIndex mocks the CreateCLSIndex API call
func (m *MockTeoV20220901Client) CreateCLSIndex(request *teo_v20220901.CreateCLSIndexRequest) (response *teo_v20220901.CreateCLSIndexResponse, err error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*teo_v20220901.CreateCLSIndexResponse), args.Error(1)
}

// MockProviderMeta is a mock for the provider metadata
type MockProviderMeta struct {
	mock.Mock
}

// GetAPIV3Conn returns a mock API connection
func (m *MockProviderMeta) GetAPIV3Conn() interface{} {
	args := m.Called()
	return args.Get(0)
}

func TestResourceTencentCloudTeoCreateCLSIndexOperationCreate_Success(t *testing.T) {
	// Create mock client and provider meta
	mockClient := new(MockTeoV20220901Client)
	mockProviderMeta := new(MockProviderMeta)

	// Setup expectations
	mockProviderMeta.On("GetAPIV3Conn").Return(mockClient)
	mockClient.On("CreateCLSIndex", mock.Anything).Return(&teo_v20220901.CreateCLSIndexResponse{}, nil)

	// Create resource data schema
	resourceSchema := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	resourceData := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"task_id": "task-87654321",
	})

	// Call Create function
	err := resourceData.Create(mockProviderMeta, common.ProviderMeta{})

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", resourceData.Id())
	mockClient.AssertExpectations(t)
	mockProviderMeta.AssertExpectations(t)
}

func TestResourceTencentCloudTeoCreateCLSIndexOperationCreate_APIError(t *testing.T) {
	// Create mock client and provider meta
	mockClient := new(MockTeoV20220901Client)
	mockProviderMeta := new(MockProviderMeta)

	// Setup expectations
	mockProviderMeta.On("GetAPIV3Conn").Return(mockClient)
	mockClient.On("CreateCLSIndex", mock.Anything).Return(nil, errors.New("API error: Invalid zone_id"))

	// Create resource data schema
	resourceSchema := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	resourceData := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"task_id": "task-87654321",
	})

	// Call Create function
	err := resourceData.Create(mockProviderMeta, common.ProviderMeta{})

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
	mockClient.AssertExpectations(t)
	mockProviderMeta.AssertExpectations(t)
}

func TestResourceTencentCloudTeoCreateCLSIndexOperationCreate_MissingZoneId(t *testing.T) {
	// Create resource data schema without zone_id
	resourceSchema := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	resourceData := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"task_id": "task-87654321",
	})

	// Call Create function
	err := resourceData.Create(nil, common.ProviderMeta{})

	// Assertions
	// The error should indicate that zone_id is required
	assert.Error(t, err)
}

func TestResourceTencentCloudTeoCreateCLSIndexOperationCreate_MissingTaskId(t *testing.T) {
	// Create resource data schema without task_id
	resourceSchema := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	resourceData := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})

	// Call Create function
	err := resourceData.Create(nil, common.ProviderMeta{})

	// Assertions
	// The error should indicate that task_id is required
	assert.Error(t, err)
}

func TestResourceTencentCloudTeoCreateCLSIndexOperationRead(t *testing.T) {
	// Create resource data schema
	resourceSchema := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	resourceData := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"task_id": "task-87654321",
	})
	resourceData.SetId("zone-12345678")

	// Call Read function (should return nil for operation resource)
	err := resourceData.Read(nil, common.ProviderMeta{})

	// Assertions
	assert.NoError(t, err)
}

func TestResourceTencentCloudTeoCreateCLSIndexOperationDelete(t *testing.T) {
	// Create resource data schema
	resourceSchema := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()
	resourceData := schema.TestResourceDataRaw(t, resourceSchema.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"task_id": "task-87654321",
	})
	resourceData.SetId("zone-12345678")

	// Call Delete function (should return nil for operation resource)
	err := resourceData.Delete(nil, common.ProviderMeta{})

	// Assertions
	assert.NoError(t, err)
}

func TestResourceTencentCloudTeoCreateCLSIndexOperationSchema(t *testing.T) {
	// Get the resource schema
	resourceSchema := teo.ResourceTencentCloudTeoCreateCLSIndexOperation()

	// Assertions
	assert.NotNil(t, resourceSchema)
	assert.NotNil(t, resourceSchema.Create)
	assert.NotNil(t, resourceSchema.Read)
	assert.Nil(t, resourceSchema.Update) // Operation resources don't have Update
	assert.NotNil(t, resourceSchema.Delete)

	// Check schema fields
	assert.NotNil(t, resourceSchema.Schema)
	assert.Contains(t, resourceSchema.Schema, "zone_id")
	assert.Contains(t, resourceSchema.Schema, "task_id")

	// Check zone_id
	zoneIdSchema := resourceSchema.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneIdSchema.Type)
	assert.True(t, zoneIdSchema.Required)
	assert.True(t, zoneIdSchema.ForceNew)

	// Check task_id
	taskIdSchema := resourceSchema.Schema["task_id"]
	assert.Equal(t, schema.TypeString, taskIdSchema.Type)
	assert.True(t, taskIdSchema.Required)
	assert.True(t, taskIdSchema.ForceNew)
}
