package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainVerifyUserAccountDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainVerifyUserAccountDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_verify_user_account.verify_user_account")),
			},
		},
	})
}

const testAccDbbrainVerifyUserAccountDataSource = `

data "tencentcloud_dbbrain_verify_user_account" "verify_user_account" {
  instance_id = ""
  user = ""
  password = ""
  product = ""
  }

`
