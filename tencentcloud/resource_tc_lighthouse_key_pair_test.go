package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseKeyPairResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseKeyPair,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_key_pair.key_pair", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_key_pair.key_pair",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseKeyPair = `

resource "tencentcloud_lighthouse_key_pair" "key_pair" {
  key_name = "key_name_test"
  public_key = "public key content"
}

`
