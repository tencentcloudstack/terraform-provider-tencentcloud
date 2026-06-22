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

// go test ./tencentcloud/services/mqtt/ -run "TestMqttInstanceX509Mode" -v -count=1 -gcflags="all=-l"

// TestMqttInstanceX509Mode_Read tests that x509_mode is correctly read from DescribeInstance response
func TestMqttInstanceX509Mode_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response with X509Mode
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
			X509Mode:       ptrMqttInstanceString("mTLS"),
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
	assert.Equal(t, "mTLS", d.Get("x509_mode").(string))
}

// TestMqttInstanceX509Mode_ReadNil tests that x509_mode is not set when response field is nil
func TestMqttInstanceX509Mode_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response without X509Mode (nil)
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
			X509Mode:       nil,
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
	assert.Equal(t, "", d.Get("x509_mode").(string))
}

// TestMqttInstanceX509Mode_Schema tests the schema definition of x509_mode
func TestMqttInstanceX509Mode_Schema(t *testing.T) {
	res := svcmqtt.ResourceTencentCloudMqttInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "x509_mode")

	x509Mode := res.Schema["x509_mode"]
	assert.Equal(t, schema.TypeString, x509Mode.Type)
	assert.True(t, x509Mode.Optional)
	assert.True(t, x509Mode.Computed)
}

// go test ./tencentcloud/services/mqtt/ -run "TestMqttInstanceDeviceCertificateProvisionType" -v -count=1 -gcflags="all=-l"

// TestMqttInstanceDeviceCertificateProvisionType_Schema tests the schema definition of device_certificate_provision_type
func TestMqttInstanceDeviceCertificateProvisionType_Schema(t *testing.T) {
	res := svcmqtt.ResourceTencentCloudMqttInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "device_certificate_provision_type")

	field := res.Schema["device_certificate_provision_type"]
	assert.Equal(t, schema.TypeString, field.Type)
	assert.True(t, field.Optional)
	assert.True(t, field.Computed)
}

// TestMqttInstanceDeviceCertificateProvisionType_Read tests that device_certificate_provision_type is correctly read from DescribeInstance response
func TestMqttInstanceDeviceCertificateProvisionType_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response with DeviceCertificateProvisionType
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:                   ptrMqttInstanceString("PRO"),
			InstanceName:                   ptrMqttInstanceString("test-instance"),
			Remark:                         ptrMqttInstanceString("test remark"),
			SkuCode:                        ptrMqttInstanceString("pro_6k_1"),
			PayMode:                        ptrMqttInstanceString("POSTPAID"),
			RenewFlag:                      ptrMqttInstanceInt64(0),
			InstanceStatus:                 ptrMqttInstanceString("RUNNING"),
			DeviceCertificateProvisionType: ptrMqttInstanceString("JITP"),
			RequestId:                      ptrMqttInstanceString("fake-request-id"),
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
	assert.Equal(t, "JITP", d.Get("device_certificate_provision_type").(string))
}

// TestMqttInstanceDeviceCertificateProvisionType_ReadNil tests that device_certificate_provision_type is not set when response field is nil
func TestMqttInstanceDeviceCertificateProvisionType_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response without DeviceCertificateProvisionType (nil)
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:                   ptrMqttInstanceString("PRO"),
			InstanceName:                   ptrMqttInstanceString("test-instance"),
			Remark:                         ptrMqttInstanceString("test remark"),
			SkuCode:                        ptrMqttInstanceString("pro_6k_1"),
			PayMode:                        ptrMqttInstanceString("POSTPAID"),
			RenewFlag:                      ptrMqttInstanceInt64(0),
			InstanceStatus:                 ptrMqttInstanceString("RUNNING"),
			DeviceCertificateProvisionType: nil,
			RequestId:                      ptrMqttInstanceString("fake-request-id"),
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
	assert.Equal(t, "", d.Get("device_certificate_provision_type").(string))
}

// TestMqttInstanceDeviceCertificateProvisionType_Update tests that device_certificate_provision_type is sent in ModifyInstance request
func TestMqttInstanceDeviceCertificateProvisionType_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	var capturedRequest *mqttv20240516.ModifyInstanceRequest

	// Patch ModifyInstanceWithContext to capture the request
	patches.ApplyMethodFunc(mqttClient, "ModifyInstanceWithContext", func(_ context.Context, request *mqttv20240516.ModifyInstanceRequest) (*mqttv20240516.ModifyInstanceResponse, error) {
		capturedRequest = request
		resp := mqttv20240516.NewModifyInstanceResponse()
		resp.Response = &mqttv20240516.ModifyInstanceResponseParams{
			RequestId: ptrMqttInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeInstance for the Read call after Update
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:                   ptrMqttInstanceString("PRO"),
			InstanceName:                   ptrMqttInstanceString("test-instance"),
			Remark:                         ptrMqttInstanceString("test remark"),
			SkuCode:                        ptrMqttInstanceString("pro_6k_1"),
			PayMode:                        ptrMqttInstanceString("POSTPAID"),
			RenewFlag:                      ptrMqttInstanceInt64(0),
			InstanceStatus:                 ptrMqttInstanceString("RUNNING"),
			DeviceCertificateProvisionType: ptrMqttInstanceString("JITP"),
			RequestId:                      ptrMqttInstanceString("fake-request-id"),
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
		"instance_type":                     "PRO",
		"name":                              "test-instance",
		"sku_code":                          "pro_6k_1",
		"device_certificate_provision_type": "JITP",
	})
	d.SetId("mqtt-test-instance-id")

	// Patch HasChange to simulate a change in device_certificate_provision_type
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "device_certificate_provision_type"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.DeviceCertificateProvisionType)
	assert.Equal(t, "JITP", *capturedRequest.DeviceCertificateProvisionType)
}

// go test ./tencentcloud/services/mqtt/ -run "TestMqttInstanceMessageRate" -v -count=1 -gcflags="all=-l"

// TestMqttInstanceMessageRate_Schema tests the schema definition of message_rate
func TestMqttInstanceMessageRate_Schema(t *testing.T) {
	res := svcmqtt.ResourceTencentCloudMqttInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "message_rate")

	field := res.Schema["message_rate"]
	assert.Equal(t, schema.TypeInt, field.Type)
	assert.True(t, field.Optional)
	assert.True(t, field.Computed)
}

// TestMqttInstanceMessageRate_Read tests that message_rate is correctly read from DescribeInstance response
func TestMqttInstanceMessageRate_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response with MessageRate
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
			MessageRate:    ptrMqttInstanceInt64(100),
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
	assert.Equal(t, 100, d.Get("message_rate").(int))
}

// TestMqttInstanceMessageRate_ReadNil tests that message_rate is not set when response field is nil
func TestMqttInstanceMessageRate_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response without MessageRate (nil)
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
			MessageRate:    nil,
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
	assert.Equal(t, 0, d.Get("message_rate").(int))
}

// TestMqttInstanceMessageRate_Update tests that message_rate is sent in ModifyInstance request
func TestMqttInstanceMessageRate_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	var capturedRequest *mqttv20240516.ModifyInstanceRequest

	// Patch ModifyInstanceWithContext to capture the request
	patches.ApplyMethodFunc(mqttClient, "ModifyInstanceWithContext", func(_ context.Context, request *mqttv20240516.ModifyInstanceRequest) (*mqttv20240516.ModifyInstanceResponse, error) {
		capturedRequest = request
		resp := mqttv20240516.NewModifyInstanceResponse()
		resp.Response = &mqttv20240516.ModifyInstanceResponseParams{
			RequestId: ptrMqttInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeInstance for the Read call after Update
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
			MessageRate:    ptrMqttInstanceInt64(200),
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
		"message_rate":  200,
	})
	d.SetId("mqtt-test-instance-id")

	// Patch HasChange to simulate a change in message_rate
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "message_rate"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.MessageRate)
	assert.Equal(t, int64(200), *capturedRequest.MessageRate)
}

// go test ./tencentcloud/services/mqtt/ -run "TestMqttInstanceUseDefaultServerCert" -v -count=1 -gcflags="all=-l"

// TestMqttInstanceUseDefaultServerCert_Schema tests the schema definition of use_default_server_cert
func TestMqttInstanceUseDefaultServerCert_Schema(t *testing.T) {
	res := svcmqtt.ResourceTencentCloudMqttInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "use_default_server_cert")

	field := res.Schema["use_default_server_cert"]
	assert.Equal(t, schema.TypeBool, field.Type)
	assert.True(t, field.Optional)
	assert.True(t, field.Computed)
}

// TestMqttInstanceUseDefaultServerCert_Read tests that use_default_server_cert is correctly read from DescribeInstance response
func TestMqttInstanceUseDefaultServerCert_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response with UseDefaultServerCert
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:         ptrMqttInstanceString("PRO"),
			InstanceName:         ptrMqttInstanceString("test-instance"),
			Remark:               ptrMqttInstanceString("test remark"),
			SkuCode:              ptrMqttInstanceString("pro_6k_1"),
			PayMode:              ptrMqttInstanceString("POSTPAID"),
			RenewFlag:            ptrMqttInstanceInt64(0),
			InstanceStatus:       ptrMqttInstanceString("RUNNING"),
			UseDefaultServerCert: ptrMqttInstanceBool(true),
			RequestId:            ptrMqttInstanceString("fake-request-id"),
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
	assert.Equal(t, true, d.Get("use_default_server_cert").(bool))
}

// TestMqttInstanceUseDefaultServerCert_ReadNil tests that use_default_server_cert is not set when response field is nil
func TestMqttInstanceUseDefaultServerCert_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	// Patch DescribeInstance to return response without UseDefaultServerCert (nil)
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:         ptrMqttInstanceString("PRO"),
			InstanceName:         ptrMqttInstanceString("test-instance"),
			Remark:               ptrMqttInstanceString("test remark"),
			SkuCode:              ptrMqttInstanceString("pro_6k_1"),
			PayMode:              ptrMqttInstanceString("POSTPAID"),
			RenewFlag:            ptrMqttInstanceInt64(0),
			InstanceStatus:       ptrMqttInstanceString("RUNNING"),
			UseDefaultServerCert: nil,
			RequestId:            ptrMqttInstanceString("fake-request-id"),
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
	assert.Equal(t, false, d.Get("use_default_server_cert").(bool))
}

// TestMqttInstanceUseDefaultServerCert_Update tests that use_default_server_cert is sent in ModifyInstance request
func TestMqttInstanceUseDefaultServerCert_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mqttClient := &mqttv20240516.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseMqttV20240516Client", mqttClient)

	var capturedRequest *mqttv20240516.ModifyInstanceRequest

	// Patch ModifyInstanceWithContext to capture the request
	patches.ApplyMethodFunc(mqttClient, "ModifyInstanceWithContext", func(_ context.Context, request *mqttv20240516.ModifyInstanceRequest) (*mqttv20240516.ModifyInstanceResponse, error) {
		capturedRequest = request
		resp := mqttv20240516.NewModifyInstanceResponse()
		resp.Response = &mqttv20240516.ModifyInstanceResponseParams{
			RequestId: ptrMqttInstanceString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeInstance for the Read call after Update
	patches.ApplyMethodFunc(mqttClient, "DescribeInstance", func(request *mqttv20240516.DescribeInstanceRequest) (*mqttv20240516.DescribeInstanceResponse, error) {
		resp := mqttv20240516.NewDescribeInstanceResponse()
		resp.Response = &mqttv20240516.DescribeInstanceResponseParams{
			InstanceType:         ptrMqttInstanceString("PRO"),
			InstanceName:         ptrMqttInstanceString("test-instance"),
			Remark:               ptrMqttInstanceString("test remark"),
			SkuCode:              ptrMqttInstanceString("pro_6k_1"),
			PayMode:              ptrMqttInstanceString("POSTPAID"),
			RenewFlag:            ptrMqttInstanceInt64(0),
			InstanceStatus:       ptrMqttInstanceString("RUNNING"),
			UseDefaultServerCert: ptrMqttInstanceBool(true),
			RequestId:            ptrMqttInstanceString("fake-request-id"),
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
		"instance_type":           "PRO",
		"name":                    "test-instance",
		"sku_code":                "pro_6k_1",
		"use_default_server_cert": true,
	})
	d.SetId("mqtt-test-instance-id")

	// Patch HasChange to simulate a change in use_default_server_cert
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "use_default_server_cert"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.UseDefaultServerCert)
	assert.Equal(t, true, *capturedRequest.UseDefaultServerCert)
}
