package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixSsmRotateProductSecretResource_basic -v
func TestAccTencentCloudNeedFixSsmRotateProductSecretResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmRotateProductSecret,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ssm_rotate_product_secret.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssm_rotate_product_secret.example", "secret_name"),
				),
			},
		},
	})
}

const testAccSsmRotateProductSecret = `
resource "tencentcloud_ssm_rotate_product_secret" "example" {
  secret_name = "tf_example"
}
`
