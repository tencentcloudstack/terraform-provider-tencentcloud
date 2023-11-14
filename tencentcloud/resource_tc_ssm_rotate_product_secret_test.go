package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSsmRotateProductSecretResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmRotateProductSecret,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssm_rotate_product_secret.rotate_product_secret", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssm_rotate_product_secret.rotate_product_secret",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSsmRotateProductSecret = `

resource "tencentcloud_ssm_rotate_product_secret" "rotate_product_secret" {
  secret_name = ""
}

`
