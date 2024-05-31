package cvm_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudReservedInstance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPlacementGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPlacementGroupExists("tencentcloud_reserved_instance.reserved_instance"),
					resource.TestCheckResourceAttr("tencentcloud_reserved_instance.reserved_instance", "config_id", "469043dd-28b9-4d89-b557-74f6a8326259"),
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
  config_id              = "469043dd-28b9-4d89-b557-74f6a8326259"
  instance_count         = 2
  reserved_instance_name = "test"
}
`
