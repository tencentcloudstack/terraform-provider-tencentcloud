package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmRebootInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRebootInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_reboot_instance.reboot", "id"),
				),
			},
		},
	})
}

const testAccCvmRebootInstance = `
data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}
data "tencentcloud_instance_types" "default" {
  memory_size = 2
  exclude_sold_out = true
  filter {
    name = "instance-family"
    values = ["SA2","SA3","SA4","SA5","S2","S3"]
  }
  filter {
    name = "zone"
    values = ["ap-guangzhou-7"]
  }
  cpu_core_count = 2
}
resource "tencentcloud_instance" "test_cvm" {
  image_id = data.tencentcloud_images.default.images.0.image_id
  availability_zone = "ap-guangzhou-7"
  instance_type = "SA2.MEDIUM2"
  orderly_security_groups = ["sg-5275dorp"]
  instance_charge_type = "POSTPAID_BY_HOUR"
}

resource tencentcloud_cvm_reboot_instance reboot {
  instance_id = tencentcloud_instance.test_cvm.id
  stop_type = "SOFT"
}
`
