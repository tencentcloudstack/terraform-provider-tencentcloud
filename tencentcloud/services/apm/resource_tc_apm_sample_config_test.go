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
  instance_id    = tencentcloud_apm_instance.example.id
  sample_name    = "tf-example"
  sample_rate    = 90
  service_name   = "java-order-serive"
  operation_type = 0
  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
`

const testAccApmSampleConfigUpdate = `
resource "tencentcloud_apm_sample_config" "example" {
  instance_id    = tencentcloud_apm_instance.example.id
  sample_name    = "tf-example"
  sample_rate    = 95
  service_name   = "java-order-serive"
  operation_type = 1
  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
`
