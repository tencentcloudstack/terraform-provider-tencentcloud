package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudClsLogset_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsLogset_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cls_logset.logset_basic", "logset_name", "cls888"),
				),
			},
			{
				Config: testAccClsLogset_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cls_logset.logset_basic", "logset_name", "cls9999"),
				),
			},
		},
	})
}

const testAccClsLogset_basic = `
resource "tencentcloud_cls_logset" "logset_basic" {
  logset_name    = "cls888"
  tags{
		 key = "aaa"
		 value = "bbb"
      }
}
`
const testAccClsLogset_update = `
	resource "tencentcloud_cls_logset" "logset_basic" {
  	logset_name    = "cls9999"
	tags{
		 key = "ccc"
		 value = "ddd"
      }
}
`

//asdad
