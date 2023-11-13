package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbEncryptAttributesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbEncryptAttributes,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_encrypt_attributes.encrypt_attributes", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_encrypt_attributes.encrypt_attributes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbEncryptAttributes = `

resource "tencentcloud_dcdb_encrypt_attributes" "encrypt_attributes" {
  instance_id = ""
  encrypt_enabled = 
}

`
