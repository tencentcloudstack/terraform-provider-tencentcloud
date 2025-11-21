package apm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudApmSampleConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApmSampleConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "sample_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "sample_rate"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "service_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "operation_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "operation_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "sample_config_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "tags"),
				),
			},
			{
				Config: testAccApmSampleConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "sample_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "sample_rate"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "service_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "operation_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "operation_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "sample_config_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_sample_config.example", "tags"),
				),
			},
			{
				ResourceName:      "tencentcloud_apm_sample_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApmSampleConfig = `
resource "tencentcloud_apm_sample_config" "example" {
  instance_id          = ""
  sample_name          = ""
  sample_rate          = 10
  service_name         = ""
  operation_name       = ""
  operation_type       = ""
  sample_config_status = 1
  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
`

const testAccApmSampleConfigUpdate = `
resource "tencentcloud_apm_sample_config" "example" {
  instance_id          = ""
  sample_name          = ""
  sample_rate          = 20
  service_name         = ""
  operation_name       = ""
  operation_type       = ""
  sample_config_status = 0
  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
`
