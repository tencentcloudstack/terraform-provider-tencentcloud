package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsWhiteBoxKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsWhiteBoxKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kms_white_box_key.white_box_key", "id")),
			},
			{
				ResourceName:      "tencentcloud_kms_white_box_key.white_box_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKmsWhiteBoxKey = `

resource "tencentcloud_kms_white_box_key" "white_box_key" {
  alias = "test_alias"
  description = "test_description"
  algorithm = "SM4"
  tags = {
    "createdBy" = "terraform"
  }
}

`
