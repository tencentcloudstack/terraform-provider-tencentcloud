package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbEncryptAttributesOperationResource_basic -v
func TestAccTencentCloudMariadbEncryptAttributesOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbEncryptAttributesOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_encrypt_attributes_operation.encrypt_attributes_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_encrypt_attributes_operation.encrypt_attributes_operation", "encrypt_enabled", "1"),
				),
			},
		},
	})
}

const testAccMariadbEncryptAttributesOperation = testAccMariadbHourDbInstance + `

resource "tencentcloud_mariadb_encrypt_attributes_operation" "encrypt_attributes_operation" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
  encrypt_enabled = 1
}

`
