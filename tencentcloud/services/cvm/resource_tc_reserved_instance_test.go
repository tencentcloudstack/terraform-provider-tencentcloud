package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudReservedInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_reserved_instance.reserved_instance", "config_id", "3ae825a8-bf6f-4c43-aa2c-c507c08fba92"),
					resource.TestCheckResourceAttr("tencentcloud_reserved_instance.reserved_instance", "instance_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_reserved_instance.reserved_instance", "reserved_instance_name", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserved_instance.reserved_instance", "start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserved_instance.reserved_instance", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_reserved_instance.reserved_instance", "status"),
				),
			},
			{
				ResourceName: "tencentcloud_reserved_instance.reserved_instance",
				ImportState:  true,
			},
		},
	})
}

const testAccReservedInstance = `
resource "tencentcloud_reserved_instance" "reserved_instance" {
  config_id              = "3ae825a8-bf6f-4c43-aa2c-c507c08fba92"
  instance_count         = 1
  reserved_instance_name = "test"
}
`
