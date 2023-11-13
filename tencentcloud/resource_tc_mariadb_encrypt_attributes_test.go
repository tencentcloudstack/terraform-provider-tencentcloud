package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_encrypt_attributes.encrypt_attributes", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_encrypt_attributes.encrypt_attributes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbEncryptAttributes = `

resource "tencentcloud_mariadb_encrypt_attributes" "encrypt_attributes" {
  instance_id = "tdsql-e9tklsgz"
  encrypt_enabled = 
}

`
