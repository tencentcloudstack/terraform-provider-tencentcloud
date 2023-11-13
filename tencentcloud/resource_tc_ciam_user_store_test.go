package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCiamUserStoreResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiamUserStore,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ciam_user_store.user_store", "id")),
			},
			{
				ResourceName:      "tencentcloud_ciam_user_store.user_store",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiamUserStore = `

resource "tencentcloud_ciam_user_store" "user_store" {
  user_pool_name = ""
  user_pool_desc = ""
  user_pool_logo = ""
}

`
