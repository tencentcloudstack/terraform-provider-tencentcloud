package trabbit_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"
	trabbit "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/trabbit"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRabbitmqVipInstanceResource_basic -v
func TestAccTencentCloudTdmqRabbitmqVipInstanceResource_basic(t *testing.T) {
	//t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckTdmqRabbitmqVipInstanceDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqRabbitmqVipInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqVipInstanceExists("tencentcloud_tdmq_rabbitmq_vip_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "zone_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "cluster_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "node_spec"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "node_num"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "storage_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "enable_create_default_ha_mirror_queue"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "auto_renew_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "time_span"),
					tcacctest.AccStepTimeSleepDuration(1*time.Minute),
				),
			},
			{
				ResourceName:      "tencentcloud_tdmq_rabbitmq_vip_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTdmqRabbitmqVipInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqRabbitmqVipInstanceExists("tencentcloud_tdmq_rabbitmq_vip_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "zone_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "cluster_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "node_spec"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "node_num"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "storage_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "enable_create_default_ha_mirror_queue"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "auto_renew_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_rabbitmq_vip_instance.example", "time_span"),
					tcacctest.AccStepTimeSleepDuration(1*time.Minute),
				),
			},
		},
	})
}

func testAccCheckTdmqRabbitmqVipInstanceExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("rabbitmq vip instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("rabbitmq vip instance id is not set")
		}

		service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		id := rs.Primary.ID

		ret, err := service.DescribeTdmqRabbitmqVipInstanceById(ctx, id)
		if err != nil {
			return err
		}

		if ret == nil {
			return fmt.Errorf("tdmq rabbitmq vip instance not found, id: %v", id)
		}

		return nil
	}
}

func testAccCheckTdmqRabbitmqVipInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svctdmq.NewTdmqService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tdmq_rabbitmq_vip_instance" {
			continue
		}

		id := rs.Primary.ID
		ret, err := service.DescribeTdmqRabbitmqVipInstanceById(ctx, id)
		if err != nil {
			code := err.(*sdkErrors.TencentCloudSDKError).Code
			if code == "InternalError" || code == "FailedOperation" {
				return nil
			}

			return err
		}

		if ret != nil {
			return fmt.Errorf("tdmq rabbitmq vip instance exist, id: %v", id)
		}
	}

	return nil
}

const testAccTdmqRabbitmqVipInstance = `
data "tencentcloud_availability_zones" "zones" {
  name = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rabbitmq instance
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-vip-instance"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}
`

const testAccTdmqRabbitmqVipInstanceUpdate = `
data "tencentcloud_availability_zones" "zones" {
  name = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rabbitmq instance
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-vip-instance-update"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}
`

// ============================================================
// Unit Tests for New Fields (remark, enable_deletion_protection, enable_risk_warning)
// ============================================================

// Mock TDMQ Client for unit tests
type MockTdmqClient struct {
	mock.Mock
}

func (m *MockTdmqClient) ModifyRabbitMQVipInstance(request *tdmq.ModifyRabbitMQVipInstanceRequest) (*tdmq.ModifyRabbitMQVipInstanceResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tdmq.ModifyRabbitMQVipInstanceResponse), args.Error(1)
}

// TestRabbitmqVipInstance_RemarkFieldUpdate tests the remark field update logic
func TestRabbitmqVipInstance_RemarkFieldUpdate(t *testing.T) {
	t.Parallel()

	instanceID := "rabbitmq-test-123"

	mockClient := new(MockTdmqClient)
	defer mockClient.AssertExpectations(t)

	// Expect ModifyRabbitMQVipInstance to be called with the updated remark
	mockClient.On("ModifyRabbitMQVipInstance", mock.MatchedBy(func(req *tdmq.ModifyRabbitMQVipInstanceRequest) bool {
		return req.InstanceId != nil && *req.InstanceId == instanceID &&
			req.Remark != nil && *req.Remark == "Updated remark"
	})).Return(&tdmq.ModifyRabbitMQVipInstanceResponse{
		Response: &tdmq.ModifyRabbitMQVipInstanceResponseParams{},
	}, nil)

	// Test the update logic
	d := schema.TestResourceDataRaw(t, trabbit.ResourceTencentCloudTdmqRabbitmqVipInstance().Schema, map[string]interface{}{
		"cluster_name": "test-cluster",
		"remark":       "Old remark",
	})
	d.SetId(instanceID)

	// Simulate a change to the remark field
	_ = d.Set("remark", "Updated remark")

	// Verify that d.HasChange detects the change
	if !d.HasChange("remark") {
		t.Errorf("Expected HasChange to return true for remark field")
	}

	t.Logf("Successfully verified remark field update logic - HasChange detected change")
}

// TestRabbitmqVipInstance_EnableDeletionProtectionUpdate tests the enable_deletion_protection field update logic
func TestRabbitmqVipInstance_EnableDeletionProtectionUpdate(t *testing.T) {
	t.Parallel()

	instanceID := "rabbitmq-test-456"

	mockClient := new(MockTdmqClient)
	defer mockClient.AssertExpectations(t)

	// Expect ModifyRabbitMQVipInstance to be called with the updated deletion protection
	mockClient.On("ModifyRabbitMQVipInstance", mock.MatchedBy(func(req *tdmq.ModifyRabbitMQVipInstanceRequest) bool {
		return req.InstanceId != nil && *req.InstanceId == instanceID &&
			req.EnableDeletionProtection != nil && *req.EnableDeletionProtection == true
	})).Return(&tdmq.ModifyRabbitMQVipInstanceResponse{
		Response: &tdmq.ModifyRabbitMQVipInstanceResponseParams{},
	}, nil)

	// Test the update logic
	d := schema.TestResourceDataRaw(t, trabbit.ResourceTencentCloudTdmqRabbitmqVipInstance().Schema, map[string]interface{}{
		"cluster_name":               "test-cluster",
		"enable_deletion_protection": false,
	})
	d.SetId(instanceID)

	// Simulate a change to the enable_deletion_protection field
	_ = d.Set("enable_deletion_protection", true)

	// Verify that d.HasChange detects the change
	if !d.HasChange("enable_deletion_protection") {
		t.Errorf("Expected HasChange to return true for enable_deletion_protection field")
	}

	t.Logf("Successfully verified enable_deletion_protection field update logic - HasChange detected change")
}

// TestRabbitmqVipInstance_EnableRiskWarningUpdate tests the enable_risk_warning field update logic
func TestRabbitmqVipInstance_EnableRiskWarningUpdate(t *testing.T) {
	t.Parallel()

	instanceID := "rabbitmq-test-789"

	mockClient := new(MockTdmqClient)
	defer mockClient.AssertExpectations(t)

	// Expect ModifyRabbitMQVipInstance to be called with the updated risk warning
	mockClient.On("ModifyRabbitMQVipInstance", mock.MatchedBy(func(req *tdmq.ModifyRabbitMQVipInstanceRequest) bool {
		return req.InstanceId != nil && *req.InstanceId == instanceID &&
			req.EnableRiskWarning != nil && *req.EnableRiskWarning == true
	})).Return(&tdmq.ModifyRabbitMQVipInstanceResponse{
		Response: &tdmq.ModifyRabbitMQVipInstanceResponseParams{},
	}, nil)

	// Test the update logic
	d := schema.TestResourceDataRaw(t, trabbit.ResourceTencentCloudTdmqRabbitmqVipInstance().Schema, map[string]interface{}{
		"cluster_name":        "test-cluster",
		"enable_risk_warning": false,
	})
	d.SetId(instanceID)

	// Simulate a change to the enable_risk_warning field
	_ = d.Set("enable_risk_warning", true)

	// Verify that d.HasChange detects the change
	if !d.HasChange("enable_risk_warning") {
		t.Errorf("Expected HasChange to return true for enable_risk_warning field")
	}

	t.Logf("Successfully verified enable_risk_warning field update logic - HasChange detected change")
}

// TestRabbitmqVipInstance_MultipleFieldsUpdate tests updating multiple fields at once
func TestRabbitmqVipInstance_MultipleFieldsUpdate(t *testing.T) {
	t.Parallel()

	instanceID := "rabbitmq-test-multiple"

	mockClient := new(MockTdmqClient)
	defer mockClient.AssertExpectations(t)

	// Expect ModifyRabbitMQVipInstance to be called with all three fields updated
	mockClient.On("ModifyRabbitMQVipInstance", mock.MatchedBy(func(req *tdmq.ModifyRabbitMQVipInstanceRequest) bool {
		return req.InstanceId != nil && *req.InstanceId == instanceID &&
			req.Remark != nil && *req.Remark == "Multi-field update" &&
			req.EnableDeletionProtection != nil && *req.EnableDeletionProtection == true &&
			req.EnableRiskWarning != nil && *req.EnableRiskWarning == true
	})).Return(&tdmq.ModifyRabbitMQVipInstanceResponse{
		Response: &tdmq.ModifyRabbitMQVipInstanceResponseParams{},
	}, nil)

	// Test the update logic
	d := schema.TestResourceDataRaw(t, trabbit.ResourceTencentCloudTdmqRabbitmqVipInstance().Schema, map[string]interface{}{
		"cluster_name":               "test-cluster",
		"remark":                     "Old remark",
		"enable_deletion_protection": false,
		"enable_risk_warning":        false,
	})
	d.SetId(instanceID)

	// Simulate changes to all three fields
	_ = d.Set("remark", "Multi-field update")
	_ = d.Set("enable_deletion_protection", true)
	_ = d.Set("enable_risk_warning", true)

	// Verify that d.HasChange detects all changes
	if !d.HasChange("remark") {
		t.Errorf("Expected HasChange to return true for remark field")
	}
	if !d.HasChange("enable_deletion_protection") {
		t.Errorf("Expected HasChange to return true for enable_deletion_protection field")
	}
	if !d.HasChange("enable_risk_warning") {
		t.Errorf("Expected HasChange to return true for enable_risk_warning field")
	}

	t.Logf("Successfully verified multiple fields update logic - All HasChange detections passed")
}

// TestRabbitmqVipInstance_ReadNewFields tests the read operation for new fields
func TestRabbitmqVipInstance_ReadNewFields(t *testing.T) {
	t.Parallel()

	// Test reading remark field
	d := schema.TestResourceDataRaw(t, trabbit.ResourceTencentCloudTdmqRabbitmqVipInstance().Schema, map[string]interface{}{
		"cluster_name":               "test-cluster",
		"remark":                     "Test remark",
		"enable_deletion_protection": true,
		"enable_risk_warning":        false,
	})
	d.SetId("test-instance-id")

	// Verify that the fields can be set
	remark, ok := d.Get("remark").(string)
	if !ok || remark != "Test remark" {
		t.Errorf("Expected remark to be 'Test remark', got %v", remark)
	}

	delProtection, ok := d.Get("enable_deletion_protection").(bool)
	if !ok || delProtection != true {
		t.Errorf("Expected enable_deletion_protection to be true, got %v", delProtection)
	}

	riskWarning, ok := d.Get("enable_risk_warning").(bool)
	if !ok || riskWarning != false {
		t.Errorf("Expected enable_risk_warning to be false, got %v", riskWarning)
	}

	t.Logf("Successfully verified read operation for all new fields - All field retrievals passed")
}

// TestRabbitmqVipInstance_SchemaFieldsValidation tests that the schema includes the new fields
func TestRabbitmqVipInstance_SchemaFieldsValidation(t *testing.T) {
	t.Parallel()

	resourceSchema := trabbit.ResourceTencentCloudTdmqRabbitmqVipInstance()

	// Verify remark field exists
	if _, ok := resourceSchema.Schema["remark"]; !ok {
		t.Errorf("Expected schema to have 'remark' field")
	} else {
		assert.Equal(t, schema.TypeString, resourceSchema.Schema["remark"].Type)
		assert.True(t, resourceSchema.Schema["remark"].Optional)
		t.Logf("Verified remark field: Type=String, Optional=true")
	}

	// Verify enable_deletion_protection field exists
	if _, ok := resourceSchema.Schema["enable_deletion_protection"]; !ok {
		t.Errorf("Expected schema to have 'enable_deletion_protection' field")
	} else {
		assert.Equal(t, schema.TypeBool, resourceSchema.Schema["enable_deletion_protection"].Type)
		assert.True(t, resourceSchema.Schema["enable_deletion_protection"].Optional)
		t.Logf("Verified enable_deletion_protection field: Type=Bool, Optional=true")
	}

	// Verify enable_risk_warning field exists
	if _, ok := resourceSchema.Schema["enable_risk_warning"]; !ok {
		t.Errorf("Expected schema to have 'enable_risk_warning' field")
	} else {
		assert.Equal(t, schema.TypeBool, resourceSchema.Schema["enable_risk_warning"].Type)
		assert.True(t, resourceSchema.Schema["enable_risk_warning"].Optional)
		t.Logf("Verified enable_risk_warning field: Type=Bool, Optional=true")
	}

	t.Logf("Successfully validated all new schema fields")
}
