package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_instance.instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrInstance = `

resource "tencentcloud_tcr_instance" "instance" {
  registry_id = "tcr-xxx"
  registry_charge_prepaid {
		period = 1
		renew_flag = 0

  }
  flag = 0
  tags = {
    "createdBy" = "terraform"
  }
}

`
