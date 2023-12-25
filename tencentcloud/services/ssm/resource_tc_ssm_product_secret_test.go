package ssm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSsmProductSecretResource_basic -v
func TestAccTencentCloudSsmProductSecretResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmProductSecret,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_product_secret.product_secret", "description", "for ssm product test"),
					resource.TestCheckResourceAttr("tencentcloud_ssm_product_secret.product_secret", "status", "Disabled"),
				),
			},
			{
				Config: testAccSsmProductSecretUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_ssm_product_secret.product_secret", "description", "for ssm product"),
					resource.TestCheckResourceAttr("tencentcloud_ssm_product_secret.product_secret", "status", "Enabled"),
				),
			},
		},
	})
}

const testAccSsmProductSecret = `
data "tencentcloud_kms_keys" "kms" {
  key_state = 1
}

data "tencentcloud_mysql_instance" "mysql" {
  mysql_id = "cdb-fitq5t9h"
}

resource "tencentcloud_ssm_product_secret" "product_secret" {
  secret_name      = "tf-product-ssm-test"
  user_name_prefix = "test"
  product_name     = "Mysql"
  instance_id      = data.tencentcloud_mysql_instance.mysql.instance_list.0.mysql_id
  domains          = ["10.0.0.0"]
  privileges_list {
    privilege_name = "GlobalPrivileges"
    privileges     = ["ALTER ROUTINE"]
  }
  description = "for ssm product test"
  kms_key_id  = data.tencentcloud_kms_keys.kms.key_list.0.key_id
  status      = "Disabled"

  tags = {
    "createdBy" = "terraform"
  }
}
`

const testAccSsmProductSecretUpdate = `
data "tencentcloud_kms_keys" "kms" {
  key_state = 1
}

data "tencentcloud_mysql_instance" "mysql" {
  mysql_id = "cdb-fitq5t9h"
}

resource "tencentcloud_ssm_product_secret" "product_secret" {
  secret_name      = "tf-product-ssm-test"
  user_name_prefix = "test"
  product_name     = "Mysql"
  instance_id      = data.tencentcloud_mysql_instance.mysql.instance_list.0.mysql_id
  domains          = ["10.0.0.0"]
  privileges_list {
    privilege_name = "GlobalPrivileges"
    privileges     = ["ALTER ROUTINE"]
  }
  description = "for ssm product"
  kms_key_id  = data.tencentcloud_kms_keys.kms.key_list.0.key_id
  status      = "Enabled"

  tags = {
    "createdBy" = "terraform"
  }
}
`
