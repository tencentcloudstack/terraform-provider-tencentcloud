package mqtt_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcmqtt "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mqtt"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

func TestAccTencentCloudNeedFixMqttInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqtt,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "instance_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "sku_code"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "pay_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "device_certificate_provision_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "automatic_activation"),
				),
			},
			{
				Config: testAccMqttUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "instance_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "sku_code"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "pay_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "device_certificate_provision_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "automatic_activation"),
				),
			},
		},
	})
}

const testAccMqtt = `
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  is_multicast      = false
}

// create mqtt instance
resource "tencentcloud_mqtt_instance" "example" {
  instance_type = "BASIC"
  name          = "tf-example"
  sku_code      = "basic_2k"
  remark        = "remarks."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode                          = 0
  device_certificate_provision_type = "API"
  automatic_activation              = false
  tags = {
    createBy = "Terraform"
  }
}
`

const testAccMqttUpdate = `
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  is_multicast      = false
}

// create mqtt instance
resource "tencentcloud_mqtt_instance" "example" {
  instance_type = "BASIC"
  name          = "tf-example-update"
  sku_code      = "basic_2k"
  remark        = "remarks update."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode                          = 0
  device_certificate_provision_type = "JITP"
  automatic_activation              = true
  tags = {
    createBy = "Terraform"
  }
}
`

// mockMetaMqttInstance implements tccommon.ProviderMeta
type mockMetaMqttInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaMqttInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaMqttInstance{}

func ptrMqttInstanceInt64(i int64) *int64 {
	return &i
}

func ptrMqttInstanceString(s string) *string {
	return &s
}

func ptrMqttInstanceBool(b bool) *bool {
	return &b
}

// go test ./tencentcloud/services/mqtt/ -run "TestMqttInstanceBlockRuleLimit" -v -count=1 -gcflags="all=-l"

// TestMqttInstanceBlockRuleLimit_Read tests that block_rule_limit is correctly read from DescribeInstance response
func TestMqttInstanceBlockRuleLimit_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response with BlockRuleLimit
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:   ptrMqttInstanceString("PRO"),
			InstanceName:   ptrMqttInstanceString("test-instance"),
			Remark:         ptrMqttInstanceString("test remark"),
			SkuCode:        ptrMqttInstanceString("pro_6k_1"),
			PayMode:        ptrMqttInstanceString("POSTPAID"),
			RenewFlag:      ptrMqttInstanceInt64(0),
			InstanceStatus: ptrMqttInstanceString("RUNNING"),
			BlockRuleLimit: ptrMqttInstanceInt64(100),
			RequestId:      ptrMqttInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch TagService.DescribeResourceTags
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := &mockMetaMqttInstance{client: mockClient}
	res := svcmqtt.ResourceTencentCloudMqttInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "PRO",
		"name":          "test-instance",
		"sku_code":      "pro_6k_1",
	})
	d.SetId("mqtt-test-instance-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 100, d.Get("block_rule_limit").(int))
}

// TestMqttInstanceBlockRuleLimit_ReadNil tests that block_rule_limit is not set when response field is nil
func TestMqttInstanceBlockRuleLimit_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response without BlockRuleLimit (nil)
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:   ptrMqttInstanceString("PRO"),
			InstanceName:   ptrMqttInstanceString("test-instance"),
			Remark:         ptrMqttInstanceString("test remark"),
			SkuCode:        ptrMqttInstanceString("pro_6k_1"),
			PayMode:        ptrMqttInstanceString("POSTPAID"),
			RenewFlag:      ptrMqttInstanceInt64(0),
			InstanceStatus: ptrMqttInstanceString("RUNNING"),
			BlockRuleLimit: nil,
			RequestId:      ptrMqttInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch TagService.DescribeResourceTags
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := &mockMetaMqttInstance{client: mockClient}
	res := svcmqtt.ResourceTencentCloudMqttInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "PRO",
		"name":          "test-instance",
		"sku_code":      "pro_6k_1",
	})
	d.SetId("mqtt-test-instance-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, d.Get("block_rule_limit").(int))
}

// TestMqttInstanceBlockRuleLimit_Schema tests the schema definition of block_rule_limit
func TestMqttInstanceBlockRuleLimit_Schema(t *testing.T) {
	res := svcmqtt.ResourceTencentCloudMqttInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "block_rule_limit")

	blockRuleLimit := res.Schema["block_rule_limit"]
	assert.Equal(t, schema.TypeInt, blockRuleLimit.Type)
	assert.True(t, blockRuleLimit.Computed)
}
