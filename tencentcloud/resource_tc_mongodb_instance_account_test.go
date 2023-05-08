package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMongodbInstanceAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_account.instance_account", "id")),
			},
		},
	})
}

const testAccMongodbInstanceAccount = `

resource "tencentcloud_mongodb_instance_account" "instance_account" {
  instance_id = "cmgo-lxaz2c9b"
  user_name = "test_account"
  password = "xxxxxxxx"
  mongo_user_password = "xxxxxxxxx"
  user_desc = "test account"
  auth_role {
    mask = 0
    namespace = "*"
  }
}

`
