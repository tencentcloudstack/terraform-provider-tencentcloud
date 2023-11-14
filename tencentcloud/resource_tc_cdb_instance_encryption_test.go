package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbInstanceEncryptionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbInstanceEncryption,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_instance_encryption.instance_encryption", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_instance_encryption.instance_encryption",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbInstanceEncryption = `

resource "tencentcloud_cdb_instance_encryption" "instance_encryption" {
  instance_id = ""
  key_id = ""
  key_region = ""
}

`
