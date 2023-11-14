package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsCancelKeyArchiveResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsCancelKeyArchive,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kms_cancel_key_archive.cancel_key_archive", "id")),
			},
			{
				ResourceName:      "tencentcloud_kms_cancel_key_archive.cancel_key_archive",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKmsCancelKeyArchive = `

resource "tencentcloud_kms_cancel_key_archive" "cancel_key_archive" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
}

`
