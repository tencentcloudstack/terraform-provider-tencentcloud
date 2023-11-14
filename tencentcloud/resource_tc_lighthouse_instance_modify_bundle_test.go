package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseInstanceModifyBundleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstanceModifyBundle,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_instance_modify_bundle.instance_modify_bundle", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_instance_modify_bundle.instance_modify_bundle",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseInstanceModifyBundle = `

resource "tencentcloud_lighthouse_instance_modify_bundle" "instance_modify_bundle" {
  instance_ids = 
  bundle_id = "bundle_gen_03"
  auto_voucher = true
}

`
