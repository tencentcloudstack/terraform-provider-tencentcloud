package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosObjectAclResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosObjectAcl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_object_acl.object_acl", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cos_object_acl.object_acl", "x_cos_acl", "public-read"),
				),
			},
			{
				Config: testAccCosObjectAcl_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_object_acl.object_acl", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cos_object_acl.object_acl", "x_cos_acl", "private"),
				),
			},
		},
	})
}

const testAccCosObjectAcl = `
resource "tencentcloud_cos_object_acl" "object_acl" {
	bucket = "keep-test-1308919341"
	key = "acl.txt"
	x_cos_acl = "public-read"
}
`

const testAccCosObjectAcl_update = `
resource "tencentcloud_cos_object_acl" "object_acl" {
	bucket = "keep-test-1308919341"
	key = "acl.txt"
	x_cos_acl = "private"
}
`
