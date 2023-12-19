package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlVerifyRootAccountResource_basic -v
func TestAccTencentCloudMysqlVerifyRootAccountResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlVerifyRootAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_verify_root_account.verify_root_account", "id"),
				),
			},
		},
	})
}

const testAccMysqlVerifyRootAccount = testAccMysqlInstanceEncryptionOperationVar + `

resource "tencentcloud_mysql_verify_root_account" "verify_root_account" {
  instance_id = tencentcloud_mysql_instance.mysql8.id
  password = "password123"
}

`
