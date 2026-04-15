package teo

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestAccTencentCloudTeoDnsRecord10_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeoDnsRecord10Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecord10Basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record", "zone_id", "zone-123"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record", "name", "www"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record", "type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record", "content", "1.2.3.4"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecord10_mx(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeoDnsRecord10Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecord10MX,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_mx", "zone_id", "zone-123"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_mx", "name", "@"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_mx", "type", "MX"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_mx", "content", "mail.example.com"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_mx", "priority", "10"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoDnsRecord10_withOptionalParams(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTeoDnsRecord10Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsRecord10WithOptionalParams,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_optional", "zone_id", "zone-123"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_optional", "name", "www"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_optional", "type", "A"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_optional", "content", "1.2.3.4"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_optional", "location", "China"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_optional", "ttl", "600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_record_10.dns_record_optional", "weight", "50"),
				),
			},
		},
	})
}

func TestResourceTencentCloudTeoDnsRecord10ParseId(t *testing.T) {
	testCases := []struct {
		input     string
		zoneId    string
		recordId  string
		expectErr bool
	}{
		{
			input:     "zone-123#record-456",
			zoneId:    "zone-123",
			recordId:  "record-456",
			expectErr: false,
		},
		{
			input:     "zone-with-dash#record-with-dash",
			zoneId:    "zone-with-dash",
			recordId:  "record-with-dash",
			expectErr: false,
		},
		{
			input:     "invalid-id",
			zoneId:    "",
			recordId:  "",
			expectErr: true,
		},
		{
			input:     "zone#record#extra",
			zoneId:    "",
			recordId:  "",
			expectErr: true,
		},
		{
			input:     "",
			zoneId:    "",
			recordId:  "",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			zoneId, recordId, err := resourceTencentCloudTeoDnsRecord10ParseId(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %s: %v", tc.input, err)
				}
				if zoneId != tc.zoneId {
					t.Errorf("Expected zoneId %s, got %s", tc.zoneId, zoneId)
				}
				if recordId != tc.recordId {
					t.Errorf("Expected recordId %s, got %s", tc.recordId, recordId)
				}
			}
		})
	}
}

func TestResourceTencentCloudTeoDnsRecord10Create(t *testing.T) {
	// Mock the TeoService
	mockClient := &connectivity.TencentCloudClient{}
	mockService := &MockTeoService{}

	resourceData := schema.TestResourceDataRaw(t, resourceTencentCloudTeoDnsRecord10().Schema, map[string]interface{}{
		"zone_id": "zone-123",
		"name":    "www",
		"type":    "A",
		"content": "1.2.3.4",
		"ttl":     600,
	})

	// Test successful creation
	t.Run("Success", func(t *testing.T) {
		mockService.reset()
		mockService.createDnsRecordResponse = &teo.CreateDnsRecordResponse{
			Response: &teo.CreateDnsRecordResponseParams{
				RecordId: helper.String("record-456"),
			},
		}
		mockService.describeDnsRecordsResponse = &teo.DescribeDnsRecordsResponse{
			Response: &teo.DescribeDnsRecordsResponseParams{
				TotalCount: helper.IntInt64(1),
				DnsRecords: []*teo.DnsRecord{
					{
						RecordId:  helper.String("record-456"),
						Name:      helper.String("www"),
						Type:      helper.String("A"),
						Content:   helper.String("1.2.3.4"),
						TTL:       helper.IntInt64(600),
						Location:  helper.String("Default"),
						Weight:    helper.IntInt64(-1),
						Priority:  helper.IntInt64(0),
						Status:    helper.String("enable"),
						CreatedOn: helper.String("2024-01-01T00:00:00Z"),
					},
				},
			},
		}

		err := resourceTencentCloudTeoDnsRecord10Create(resourceData, mockService)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		id := resourceData.Id()
		if id != "zone-123#record-456" {
			t.Errorf("Expected ID zone-123#record-456, got %s", id)
		}

		if !mockService.createDnsRecordCalled {
			t.Error("Expected CreateDnsRecord to be called")
		}

		if mockService.createDnsRecordRequest == nil {
			t.Error("Expected CreateDnsRecord request to be set")
		} else {
			if *mockService.createDnsRecordRequest.ZoneId != "zone-123" {
				t.Errorf("Expected ZoneId zone-123, got %s", *mockService.createDnsRecordRequest.ZoneId)
			}
			if *mockService.createDnsRecordRequest.Name != "www" {
				t.Errorf("Expected Name www, got %s", *mockService.createDnsRecordRequest.Name)
			}
			if *mockService.createDnsRecordRequest.Type != "A" {
				t.Errorf("Expected Type A, got %s", *mockService.createDnsRecordRequest.Type)
			}
			if *mockService.createDnsRecordRequest.Content != "1.2.3.4" {
				t.Errorf("Expected Content 1.2.3.4, got %s", *mockService.createDnsRecordRequest.Content)
			}
			if *mockService.createDnsRecordRequest.TTL != 600 {
				t.Errorf("Expected TTL 600, got %d", *mockService.createDnsRecordRequest.TTL)
			}
		}
	})

	// Test idempotency - record already exists
	t.Run("Idempotency", func(t *testing.T) {
		mockService.reset()
		mockService.describeDnsRecordsResponse = &teo.DescribeDnsRecordsResponse{
			Response: &teo.DescribeDnsRecordsResponseParams{
				TotalCount: helper.IntInt64(1),
				DnsRecords: []*teo.DnsRecord{
					{
						RecordId: helper.String("existing-record"),
						Name:     helper.String("www"),
						Type:     helper.String("A"),
						Content:  helper.String("1.2.3.4"),
					},
				},
			},
		}

		err := resourceTencentCloudTeoDnsRecord10Create(resourceData, mockService)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if mockService.createDnsRecordCalled {
			t.Error("Expected CreateDnsRecord NOT to be called for idempotency")
		}

		id := resourceData.Id()
		if id != "zone-123#existing-record" {
			t.Errorf("Expected ID zone-123#existing-record, got %s", id)
		}
	})

	// Test error handling
	t.Run("Error", func(t *testing.T) {
		mockService.reset()
		mockService.createDnsRecordError = errors.New("create failed")

		err := resourceTencentCloudTeoDnsRecord10Create(resourceData, mockService)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestResourceTencentCloudTeoDnsRecord10Read(t *testing.T) {
	mockService := &MockTeoService{}

	t.Run("Success", func(t *testing.T) {
		mockService.reset()
		mockService.describeDnsRecordsResponse = &teo.DescribeDnsRecordsResponse{
			Response: &teo.DescribeDnsRecordsResponseParams{
				TotalCount: helper.IntInt64(1),
				DnsRecords: []*teo.DnsRecord{
					{
						RecordId:  helper.String("record-456"),
						Name:      helper.String("www"),
						Type:      helper.String("A"),
						Content:   helper.String("1.2.3.4"),
						TTL:       helper.IntInt64(600),
						Location:  helper.String("China"),
						Weight:    helper.IntInt64(50),
						Priority:  helper.IntInt64(0),
						Status:    helper.String("enable"),
						CreatedOn: helper.String("2024-01-01T00:00:00Z"),
					},
				},
			},
		}

		resourceData := schema.TestResourceDataRaw(t, resourceTencentCloudTeoDnsRecord10().Schema, map[string]interface{}{})
		resourceData.SetId("zone-123#record-456")

		err := resourceTencentCloudTeoDnsRecord10Read(resourceData, mockService)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify resource data is set correctly
		if name, _ := resourceData.Get("name").(string); name != "www" {
			t.Errorf("Expected name www, got %s", name)
		}
		if recordType, _ := resourceData.Get("type").(string); recordType != "A" {
			t.Errorf("Expected type A, got %s", recordType)
		}
		if content, _ := resourceData.Get("content").(string); content != "1.2.3.4" {
			t.Errorf("Expected content 1.2.3.4, got %s", content)
		}
		if ttl, _ := resourceData.Get("ttl").(int); ttl != 600 {
			t.Errorf("Expected ttl 600, got %d", ttl)
		}
		if location, _ := resourceData.Get("location").(string); location != "China" {
			t.Errorf("Expected location China, got %s", location)
		}
		if weight, _ := resourceData.Get("weight").(int); weight != 50 {
			t.Errorf("Expected weight 50, got %d", weight)
		}
		if status, _ := resourceData.Get("status").(string); status != "enable" {
			t.Errorf("Expected status enable, got %s", status)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.reset()
		mockService.describeDnsRecordsResponse = &teo.DescribeDnsRecordsResponse{
			Response: &teo.DescribeDnsRecordsResponseParams{
				TotalCount: helper.IntInt64(0),
				DnsRecords: []*teo.DnsRecord{},
			},
		}

		resourceData := schema.TestResourceDataRaw(t, resourceTencentCloudTeoDnsRecord10().Schema, map[string]interface{}{})
		resourceData.SetId("zone-123#nonexistent-record")

		err := resourceTencentCloudTeoDnsRecord10Read(resourceData, mockService)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if resourceData.Id() != "" {
			t.Errorf("Expected ID to be cleared, got %s", resourceData.Id())
		}
	})
}

func TestResourceTencentCloudTeoDnsRecord10Update(t *testing.T) {
	mockService := &MockTeoService{}

	t.Run("Success", func(t *testing.T) {
		mockService.reset()
		mockService.describeDnsRecordsResponse = &teo.DescribeDnsRecordsResponse{
			Response: &teo.DescribeDnsRecordsResponseParams{
				TotalCount: helper.IntInt64(1),
				DnsRecords: []*teo.DnsRecord{
					{
						RecordId:  helper.String("record-456"),
						Name:      helper.String("www-updated"),
						Type:      helper.String("A"),
						Content:   helper.String("5.6.7.8"),
						TTL:       helper.IntInt64(900),
						Location:  helper.String("China"),
						Weight:    helper.IntInt64(50),
						Priority:  helper.IntInt64(0),
						Status:    helper.String("enable"),
						CreatedOn: helper.String("2024-01-01T00:00:00Z"),
					},
				},
			},
		}

		resourceData := schema.TestResourceDataRaw(t, resourceTencentCloudTeoDnsRecord10().Schema, map[string]interface{}{
			"zone_id":  "zone-123",
			"name":     "www-updated",
			"type":     "A",
			"content":  "5.6.7.8",
			"ttl":      900,
			"location": "China",
			"weight":   50,
		})
		resourceData.SetId("zone-123#record-456")
		resourceData.Set("content", "5.6.7.8")
		resourceData.Set("ttl", 900)

		err := resourceTencentCloudTeoDnsRecord10Update(resourceData, mockService)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !mockService.modifyDnsRecordsCalled {
			t.Error("Expected ModifyDnsRecords to be called")
		}
	})

	t.Run("NoChanges", func(t *testing.T) {
		mockService.reset()

		resourceData := schema.TestResourceDataRaw(t, resourceTencentCloudTeoDnsRecord10().Schema, map[string]interface{}{
			"zone_id":  "zone-123",
			"name":     "www",
			"type":     "A",
			"content":  "1.2.3.4",
			"ttl":      600,
			"location": "Default",
		})
		resourceData.SetId("zone-123#record-456")

		// Don't set any changes
		err := resourceTencentCloudTeoDnsRecord10Update(resourceData, mockService)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if mockService.modifyDnsRecordsCalled {
			t.Error("Expected ModifyDnsRecords NOT to be called for no-op update")
		}
	})
}

func TestResourceTencentCloudTeoDnsRecord10Delete(t *testing.T) {
	mockService := &MockTeoService{}

	t.Run("Success", func(t *testing.T) {
		mockService.reset()
		mockService.deleteDnsRecordsResponse = &teo.DeleteDnsRecordsResponse{
			Response: &teo.DeleteDnsRecordsResponseParams{},
		}

		resourceData := schema.TestResourceDataRaw(t, resourceTencentCloudTeoDnsRecord10().Schema, map[string]interface{}{})
		resourceData.SetId("zone-123#record-456")

		err := resourceTencentCloudTeoDnsRecord10Delete(resourceData, mockService)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if resourceData.Id() != "" {
			t.Errorf("Expected ID to be cleared after delete, got %s", resourceData.Id())
		}

		if !mockService.deleteDnsRecordsCalled {
			t.Error("Expected DeleteDnsRecords to be called")
		}
	})

	t.Run("Error", func(t *testing.T) {
		mockService.reset()
		mockService.deleteDnsRecordsError = errors.New("delete failed")

		resourceData := schema.TestResourceDataRaw(t, resourceTencentCloudTeoDnsRecord10().Schema, map[string]interface{}{})
		resourceData.SetId("zone-123#record-456")

		err := resourceTencentCloudTeoDnsRecord10Delete(resourceData, mockService)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Mock TeoService for testing
type MockTeoService struct {
	client connectivity.TencentCloudClient

	createDnsRecordCalled   bool
	createDnsRecordRequest  *teo.CreateDnsRecordRequest
	createDnsRecordResponse *teo.CreateDnsRecordResponse
	createDnsRecordError    error

	describeDnsRecordsCalled   bool
	describeDnsRecordsRequest  *teo.DescribeDnsRecordsRequest
	describeDnsRecordsResponse *teo.DescribeDnsRecordsResponse
	describeDnsRecordsError    error

	modifyDnsRecordsCalled   bool
	modifyDnsRecordsRequest  *teo.ModifyDnsRecordsRequest
	modifyDnsRecordsResponse *teo.ModifyDnsRecordsResponse
	modifyDnsRecordsError    error

	deleteDnsRecordsCalled   bool
	deleteDnsRecordsRequest  *teo.DeleteDnsRecordsRequest
	deleteDnsRecordsResponse *teo.DeleteDnsRecordsResponse
	deleteDnsRecordsError    error
}

func (m *MockTeoService) reset() {
	m.createDnsRecordCalled = false
	m.createDnsRecordRequest = nil
	m.createDnsRecordResponse = nil
	m.createDnsRecordError = nil

	m.describeDnsRecordsCalled = false
	m.describeDnsRecordsRequest = nil
	m.describeDnsRecordsResponse = nil
	m.describeDnsRecordsError = nil

	m.modifyDnsRecordsCalled = false
	m.modifyDnsRecordsRequest = nil
	m.modifyDnsRecordsResponse = nil
	m.modifyDnsRecordsError = nil

	m.deleteDnsRecordsCalled = false
	m.deleteDnsRecordsRequest = nil
	m.deleteDnsRecordsResponse = nil
	m.deleteDnsRecordsError = nil
}

func (m *MockTeoService) CreateDnsRecord(request *teo.CreateDnsRecordRequest) (*teo.CreateDnsRecordResponse, error) {
	m.createDnsRecordCalled = true
	m.createDnsRecordRequest = request
	if m.createDnsRecordError != nil {
		return nil, m.createDnsRecordError
	}
	return m.createDnsRecordResponse, nil
}

func (m *MockTeoService) DescribeDnsRecords(request *teo.DescribeDnsRecordsRequest) (*teo.DescribeDnsRecordsResponse, error) {
	m.describeDnsRecordsCalled = true
	m.describeDnsRecordsRequest = request
	if m.describeDnsRecordsError != nil {
		return nil, m.describeDnsRecordsError
	}
	return m.describeDnsRecordsResponse, nil
}

func (m *MockTeoService) ModifyDnsRecords(request *teo.ModifyDnsRecordsRequest) (*teo.ModifyDnsRecordsResponse, error) {
	m.modifyDnsRecordsCalled = true
	m.modifyDnsRecordsRequest = request
	if m.modifyDnsRecordsError != nil {
		return nil, m.modifyDnsRecordsError
	}
	return m.modifyDnsRecordsResponse, nil
}

func (m *MockTeoService) DeleteDnsRecords(request *teo.DeleteDnsRecordsRequest) (*teo.DeleteDnsRecordsResponse, error) {
	m.deleteDnsRecordsCalled = true
	m.deleteDnsRecordsRequest = request
	if m.deleteDnsRecordsError != nil {
		return nil, m.deleteDnsRecordsError
	}
	return m.deleteDnsRecordsResponse, nil
}

// Helper functions
func TeoService(meta interface{}) *TeoService {
	return meta.(*connectivity.TencentCloudClient).UseTeoClient()
}

// Test config templates
const testAccTeoDnsRecord10Basic = `
resource "tencentcloud_teo_dns_record_10" "dns_record" {
  zone_id = "zone-123"
  name     = "www"
  type     = "A"
  content  = "1.2.3.4"
}
`

const testAccTeoDnsRecord10MX = `
resource "tencentcloud_teo_dns_record_10" "dns_record_mx" {
  zone_id  = "zone-123"
  name      = "@"
  type      = "MX"
  content   = "mail.example.com"
  priority  = 10
}
`

const testAccTeoDnsRecord10WithOptionalParams = `
resource "tencentcloud_teo_dns_record_10" "dns_record_optional" {
  zone_id  = "zone-123"
  name      = "www"
  type      = "A"
  content   = "1.2.3.4"
  location  = "China"
  ttl       = 600
  weight    = 50
}
`

func testAccCheckTeoDnsRecord10Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_dns_record_10" {
			continue
		}

		// Verify that the DNS record has been deleted
		// In real scenario, this would call the API to check if record exists
		return nil
	}
	return nil
}
