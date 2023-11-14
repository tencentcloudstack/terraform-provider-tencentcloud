package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsArchiveKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsArchiveKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kms_archive_key.archive_key", "id")),
			},
			{
				ResourceName:      "tencentcloud_kms_archive_key.archive_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKmsArchiveKey = `

resource "tencentcloud_kms_archive_key" "archive_key" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
}

`
