package kms_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixKmsWhiteBoxKeyResource_basic -v
func TestAccTencentCloudNeedFixKmsWhiteBoxKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsWhiteBoxKey,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kms_white_box_key.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kms_white_box_key.example", "alias", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_kms_white_box_key.example", "description", "test desc."),
					resource.TestCheckResourceAttr("tencentcloud_kms_white_box_key.example", "algorithm", "SM4"),
					resource.TestCheckResourceAttr("tencentcloud_kms_white_box_key.example", "status", "Disabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_kms_white_box_key.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccKmsWhiteBoxKeyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kms_white_box_key.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kms_white_box_key.example", "alias", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_kms_white_box_key.example", "description", "test desc."),
					resource.TestCheckResourceAttr("tencentcloud_kms_white_box_key.example", "algorithm", "SM4"),
					resource.TestCheckResourceAttr("tencentcloud_kms_white_box_key.example", "status", "Enabled"),
				),
			},
		},
	})
}

const testAccKmsWhiteBoxKey = `
resource "tencentcloud_kms_white_box_key" "example" {
  alias       = "tf_example"
  description = "test desc."
  algorithm   = "SM4"
  status      = "Disabled"
  tags        = {
    "createdBy" = "terraform"
  }
}
`

const testAccKmsWhiteBoxKeyUpdate = `
resource "tencentcloud_kms_white_box_key" "example" {
  alias       = "tf_example"
  description = "test desc."
  algorithm   = "SM4"
  status      = "Enabled"
  tags        = {
    "createdByUpdate" = "terraformUpdate"
  }
}
`
