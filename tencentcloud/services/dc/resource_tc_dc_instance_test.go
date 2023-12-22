package dc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dc_instance.instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_dc_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcInstance = `

resource "tencentcloud_dc_instance" "instance" {
  access_point_id         = "ap-shenzhen-b-ft"
  bandwidth               = 10
  customer_contact_number = "0"
  direct_connect_name     = "terraform-for-test"
  line_operator           = "In-houseWiring"
  port_type               = "10GBase-LR"
  sign_law                = true
  vlan                    = -1
}

`
