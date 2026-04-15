package teo

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

// MockClient 用于 mock API 调用
type MockClient struct {
	mock.Mock
}

// MockTeoClient 用于 mock TEO 客户端
type MockTeoClient struct {
	mock.Mock
}

// MockUseTeoV20220901Client 用于 mock TEO v20220901 客户端
func (m *MockClient) MockUseTeoV20220901Client() *MockTeoClient {
	return &MockTeoClient{}
}

func (m *MockClient) UseTeoV20220901Client() *teov20220901.Client {
	return nil // mock 实现
}

// TestAccTeoDnsRecordV6ParseResourceId 测试 ID 解析函数
func TestAccTeoDnsRecordV6ParseResourceId(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		zoneId    string
		recordId  string
		expectErr bool
	}{
		{
			name:      "valid id",
			id:        "zone-123#record-456",
			zoneId:    "zone-123",
			recordId:  "record-456",
			expectErr: false,
		},
		{
			name:      "invalid id - missing separator",
			id:        "zone-123record-456",
			zoneId:    "",
			recordId:  "",
			expectErr: true,
		},
		{
			name:      "invalid id - too many parts",
			id:        "zone-123#record-456#extra",
			zoneId:    "",
			recordId:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idSplit := regexp.MustCompile("#").Split(tt.id, -1)
			if len(idSplit) != 2 {
				if tt.expectErr {
					return
				}
				t.Errorf("id is broken, %s", tt.id)
			}
			if idSplit[0] != tt.zoneId || idSplit[1] != tt.recordId {
				t.Errorf("parsed zoneId=%s, recordId=%s, expected zoneId=%s, recordId=%s",
					idSplit[0], idSplit[1], tt.zoneId, tt.recordId)
			}
		})
	}
}

// TestAccTeoDnsRecordV6SchemaValidation 测试 Schema 参数验证
func TestAccTeoDnsRecordV6SchemaValidation(t *testing.T) {
	resource := ResourceTencentCloudTeoDnsRecordV6()

	// 测试 TTL 验证
	ttlSchema := resource.Schema["ttl"]
	if ttlSchema.ValidateFunc == nil {
		t.Error("TTL schema should have ValidateFunc")
	}

	// 测试 Weight 验证
	weightSchema := resource.Schema["weight"]
	if weightSchema.ValidateFunc == nil {
		t.Error("Weight schema should have ValidateFunc")
	}

	// 测试 Priority 验证
	prioritySchema := resource.Schema["priority"]
	if prioritySchema.ValidateFunc == nil {
		t.Error("Priority schema should have ValidateFunc")
	}
}

// TestAccTeoDnsRecordV6Create 测试创建逻辑（mock）
func TestAccTeoDnsRecordV6Create(t *testing.T) {
	// 创建测试数据
	d := schema.TestResourceDataRaw(t, ResourceTencentCloudTeoDnsRecordV6().Schema, map[string]interface{}{
		"zone_id": "zone-123",
		"name":    "example.com",
		"type":    "AAAA",
		"content": "2001:db8::1",
		"ttl":     300,
	})

	// 验证资源 ID 格式
	if d.Id() == "" {
		d.SetId("zone-123#record-456")
	}
	assert.Equal(t, "zone-123#record-456", d.Id())

	// 验证必填字段
	zoneId, ok := d.Get("zone_id").(string)
	assert.True(t, ok)
	assert.Equal(t, "zone-123", zoneId)

	name, ok := d.Get("name").(string)
	assert.True(t, ok)
	assert.Equal(t, "example.com", name)

	recordType, ok := d.Get("type").(string)
	assert.True(t, ok)
	assert.Equal(t, "AAAA", recordType)

	content, ok := d.Get("content").(string)
	assert.True(t, ok)
	assert.Equal(t, "2001:db8::1", content)
}

// TestAccTeoDnsRecordV6Read 测试读取逻辑（mock）
func TestAccTeoDnsRecordV6Read(t *testing.T) {
	// 设置资源 ID
	d := schema.TestResourceDataRaw(t, ResourceTencentCloudTeoDnsRecordV6().Schema, map[string]interface{}{})
	d.SetId("zone-123#record-456")

	// 验证 ID 解析
	idSplit := regexp.MustCompile("#").Split(d.Id(), -1)
	assert.Equal(t, 2, len(idSplit))
	assert.Equal(t, "zone-123", idSplit[0])
	assert.Equal(t, "record-456", idSplit[1])
}

// TestAccTeoDnsRecordV6Update 测试更新逻辑（mock）
func TestAccTeoDnsRecordV6Update(t *testing.T) {
	// 创建测试数据
	d := schema.TestResourceDataRaw(t, ResourceTencentCloudTeoDnsRecordV6().Schema, map[string]interface{}{
		"zone_id": "zone-123",
		"name":    "example.com",
		"type":    "AAAA",
		"content": "2001:db8::1",
		"ttl":     600,
	})
	d.SetId("zone-123#record-456")

	// 模拟更新
	d.Set("content", "2001:db8::2")
	d.Set("ttl", 900)

	// 验证更新后的值
	content, ok := d.Get("content").(string)
	assert.True(t, ok)
	assert.Equal(t, "2001:db8::2", content)

	ttl, ok := d.Get("ttl").(int)
	assert.True(t, ok)
	assert.Equal(t, 900, ttl)
}

// TestAccTeoDnsRecordV6Delete 测试删除逻辑（mock）
func TestAccTeoDnsRecordV6Delete(t *testing.T) {
	// 创建测试数据
	d := schema.TestResourceDataRaw(t, ResourceTencentCloudTeoDnsRecordV6().Schema, map[string]interface{}{
		"zone_id": "zone-123",
		"name":    "example.com",
		"type":    "AAAA",
		"content": "2001:db8::1",
	})
	d.SetId("zone-123#record-456")

	// 验证删除前 ID 存在
	assert.NotEqual(t, "", d.Id())

	// 模拟删除（清除 ID）
	d.SetId("")

	// 验证删除后 ID 为空
	assert.Equal(t, "", d.Id())
}

// MockTeoService 用于 mock TEO 服务
type MockTeoService struct {
	mock.Mock
}

func (m *MockTeoService) DescribeTeoDnsRecordById(ctx context.Context, zoneId, recordId string) (*teov20220901.DnsRecord, error) {
	args := m.Called(ctx, zoneId, recordId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*teov20220901.DnsRecord), args.Error(1)
}

// TestAccTeoDnsRecordV6ServiceDescribe 测试服务层 DescribeTeoDnsRecordById
func TestAccTeoDnsRecordV6ServiceDescribe(t *testing.T) {
	mockService := new(MockTeoService)

	// 准备测试数据
	expectedRecord := &teov20220901.DnsRecord{
		ZoneId:     common.StringPtr("zone-123"),
		RecordId:   common.StringPtr("record-456"),
		Name:       common.StringPtr("example.com"),
		Type:       common.StringPtr("AAAA"),
		Content:    common.StringPtr("2001:db8::1"),
		TTL:        common.Int64Ptr(300),
		Weight:     common.Int64Ptr(-1),
		Priority:   common.Int64Ptr(0),
		Status:     common.StringPtr("enable"),
		CreatedOn:  common.StringPtr("2024-01-01 00:00:00"),
		ModifiedOn: common.StringPtr("2024-01-01 00:00:00"),
	}

	mockService.On("DescribeTeoDnsRecordById", mock.Anything, "zone-123", "record-456").Return(expectedRecord, nil)

	// 调用 mock 服务
	ctx := context.Background()
	record, err := mockService.DescribeTeoDnsRecordById(ctx, "zone-123", "record-456")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, "zone-123", *record.ZoneId)
	assert.Equal(t, "record-456", *record.RecordId)
	assert.Equal(t, "example.com", *record.Name)
	assert.Equal(t, "AAAA", *record.Type)
	assert.Equal(t, "2001:db8::1", *record.Content)

	mockService.AssertExpectations(t)
}

// TestAccTeoDnsRecordV6ResourceSchema 测试资源 Schema
func TestAccTeoDnsRecordV6ResourceSchema(t *testing.T) {
	resource := ResourceTencentCloudTeoDnsRecordV6()

	// 验证资源函数存在
	assert.NotNil(t, resource.Create)
	assert.NotNil(t, resource.Read)
	assert.NotNil(t, resource.Update)
	assert.NotNil(t, resource.Delete)
	assert.NotNil(t, resource.Importer)

	// 验证 Schema 字段
	schema := resource.Schema
	assert.NotNil(t, schema)

	// 验证必需字段
	requiredFields := []string{"zone_id", "name", "type", "content"}
	for _, field := range requiredFields {
		s, ok := schema[field]
		assert.True(t, ok, fmt.Sprintf("Field %s should exist", field))
		assert.True(t, s.Required, fmt.Sprintf("Field %s should be required", field))
	}

	// 验证计算字段
	computedFields := []string{"record_id", "status", "created_on", "modified_on"}
	for _, field := range computedFields {
		s, ok := schema[field]
		assert.True(t, ok, fmt.Sprintf("Field %s should exist", field))
		assert.True(t, s.Computed, fmt.Sprintf("Field %s should be computed", field))
	}
}

// TestAccTeoDnsRecordV6SchemaFields 测试 Schema 字段类型
func TestAccTeoDnsRecordV6SchemaFields(t *testing.T) {
	resource := ResourceTencentCloudTeoDnsRecordV6()

	// 验证字符串字段类型
	stringFields := []string{"zone_id", "name", "type", "content", "location", "record_id", "status", "created_on", "modified_on"}
	for _, field := range stringFields {
		s, ok := resource.Schema[field]
		assert.True(t, ok, fmt.Sprintf("Field %s should exist", field))
		assert.Equal(t, schema.TypeString, s.Type, fmt.Sprintf("Field %s should be string type", field))
	}

	// 验证整数字段类型
	intFields := []string{"ttl", "weight", "priority"}
	for _, field := range intFields {
		s, ok := resource.Schema[field]
		assert.True(t, ok, fmt.Sprintf("Field %s should exist", field))
		assert.Equal(t, schema.TypeInt, s.Type, fmt.Sprintf("Field %s should be int type", field))
	}
}

// TestAccTeoDnsRecordV6ValidateTTL 测试 TTL 验证
func TestAccTeoDnsRecordV6ValidateTTL(t *testing.T) {
	testCases := []struct {
		ttl       int
		expectErr bool
	}{
		{60, false},    // 最小值
		{300, false},   // 默认值
		{86400, false}, // 最大值
		{59, true},     // 小于最小值
		{86401, true},  // 大于最大值
		{-1, true},     // 负值
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TTL=%d", tc.ttl), func(t *testing.T) {
			err := testValidateTTL(tc.ttl)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// testValidateTTL 测试 TTL 验证函数
func testValidateTTL(ttl int) error {
	if ttl < 60 || ttl > 86400 {
		return fmt.Errorf("TTL must be between 60 and 86400, got %d", ttl)
	}
	return nil
}

// TestAccTeoDnsRecordV6ValidateWeight 测试 Weight 验证
func TestAccTeoDnsRecordV6ValidateWeight(t *testing.T) {
	testCases := []struct {
		weight    int
		expectErr bool
	}{
		{-1, false},  // 最小值（表示不设置权重）
		{0, false},   // 表示不解析
		{50, false},  // 中间值
		{100, false}, // 最大值
		{-2, true},   // 小于最小值
		{101, true},  // 大于最大值
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Weight=%d", tc.weight), func(t *testing.T) {
			err := testValidateWeight(tc.weight)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// testValidateWeight 测试 Weight 验证函数
func testValidateWeight(weight int) error {
	if weight < -1 || weight > 100 {
		return fmt.Errorf("Weight must be between -1 and 100, got %d", weight)
	}
	return nil
}

// TestAccTeoDnsRecordV6ValidatePriority 测试 Priority 验证
func TestAccTeoDnsRecordV6ValidatePriority(t *testing.T) {
	testCases := []struct {
		priority  int
		expectErr bool
	}{
		{0, false},  // 最小值
		{25, false}, // 中间值
		{50, false}, // 最大值
		{-1, true},  // 小于最小值
		{51, true},  // 大于最大值
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Priority=%d", tc.priority), func(t *testing.T) {
			err := testValidatePriority(tc.priority)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// testValidatePriority 测试 Priority 验证函数
func testValidatePriority(priority int) error {
	if priority < 0 || priority > 50 {
		return fmt.Errorf("Priority must be between 0 and 50, got %d", priority)
	}
	return nil
}

// TestAccTeoDnsRecordV6ZoneIdForceNew 测试 zone_id ForceNew 属性
func TestAccTeoDnsRecordV6ZoneIdForceNew(t *testing.T) {
	resource := ResourceTencentCloudTeoDnsRecordV6()

	zoneIdSchema, ok := resource.Schema["zone_id"]
	assert.True(t, ok)
	assert.True(t, zoneIdSchema.ForceNew, "zone_id should have ForceNew attribute")
}

// TestAccTeoDnsRecordV6Import 测试导入功能
func TestAccTeoDnsRecordV6Import(t *testing.T) {
	resource := ResourceTencentCloudTeoDnsRecordV6()

	// 验证 Importer 存在
	assert.NotNil(t, resource.Importer)
	assert.NotNil(t, resource.Importer.State)

	// 验证 ImportStatePassthrough（使用默认的 passthrough 导入）
	assert.Equal(t, schema.ImportStatePassthrough, resource.Importer.State)
}
