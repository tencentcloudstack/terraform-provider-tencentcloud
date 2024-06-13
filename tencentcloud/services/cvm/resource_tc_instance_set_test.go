package cvm_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixInstanceSetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceSetBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_instance_set.foo", "instance_ids.#", "2"),
				),
			},
		},
	})
}

const testAccInstanceSetBasic = `
data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
  image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {

  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }
  filter {
    name   = "zone"
    values = ["ap-guangzhou-7"]
  }
  cpu_core_count   = 2
  memory_size      = 2
  exclude_sold_out = true
}
resource "tencentcloud_instance_set" "foo" {
  timeouts {
    create = "10m"
    read   = "10m"
    delete = "10m"
  }

  instance_count    = 2
  instance_name     = "tf-ci-test"
  availability_zone = "ap-guangzhou-7"
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "user"
  project_id        = 0
}
`
