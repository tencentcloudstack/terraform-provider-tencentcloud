package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseInstanceLoginKeyPairResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstanceLoginKeyPair,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_instance_login_key_pair.instance_login_key_pair", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_instance_login_key_pair.instance_login_key_pair",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseInstanceLoginKeyPair = `

resource "tencentcloud_lighthouse_instance_login_key_pair" "instance_login_key_pair" {
  instance_ids = 
  permit_login = ""
}

`
