package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbEncryptAttributesResource_basic -v
func TestAccTencentCloudMariadbEncryptAttributesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbEncryptAttributes,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_encrypt_attributes.encrypt_attributes", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_encrypt_attributes.encrypt_attributes", "encrypt_enabled", "1"),
				),
			},
		},
	})
}

const testAccMariadbEncryptAttributes = testAccMariadbHourDbInstance + `

resource "tencentcloud_mariadb_encrypt_attributes" "encrypt_attributes" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
  encrypt_enabled = 1
}

`
