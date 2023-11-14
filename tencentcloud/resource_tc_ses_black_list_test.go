package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesBlackListResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesBlackList,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_black_list.black_list", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_black_list.black_list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesBlackList = `

resource "tencentcloud_ses_black_list" "black_list" {
  email_address_list = 
}

`
