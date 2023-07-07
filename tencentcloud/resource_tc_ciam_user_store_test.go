package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
  user_pool_name = "tf_user_store"
  user_pool_desc = "for terraform test"
  user_pool_logo = "https://ciam-prd-1302490086.cos.ap-guangzhou.myqcloud.com/temporary/92630252a2c5422d9663db5feafd619b.png"
}

`
