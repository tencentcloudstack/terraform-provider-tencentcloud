package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudReservedInstance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists("tencentcloud_reserved_instance.reserved_instance"),
					resource.TestCheckResourceAttr("tencentcloud_reserved_instance.reserved_instance", "config_id", "ca8d39c3-4605-4e8d-9a21-05188cc9d1d4"),
					resource.TestCheckResourceAttr("tencentcloud_reserved_instance.reserved_instance", "instance_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_reserved_instance.reserved_instance", "reserved_instance_name", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserved_instance.reserved_instance", "start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserved_instance.reserved_instance", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserved_instance.reserved_instance", "status"),
				),
			},
			{
				ResourceName:      "tencentcloud_reserved_instance.reserved_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccReservedInstance = `
resource "tencentcloud_reserved_instance" "reserved_instance" {
  config_id              = "ca8d39c3-4605-4e8d-9a21-05188cc9d1d4"
  instance_count         = 2
  reserved_instance_name = "test"
}
`
