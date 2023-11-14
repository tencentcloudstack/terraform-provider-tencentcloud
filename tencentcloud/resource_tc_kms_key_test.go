package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKey,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kms_key.key", "id")),
			},
			{
				ResourceName:      "tencentcloud_kms_key.key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKmsKey = `

resource "tencentcloud_kms_key" "key" {
  alias = "test"
  description = "test"
  key_usage = "test"
  type = 11
  hsm_cluster_id = "ss"
  tags = {
    "createdBy" = "terraform"
  }
}

`
