package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsCloudResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsCloudResource,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kms_cloud_resource.cloud_resource", "id")),
			},
			{
				ResourceName:      "tencentcloud_kms_cloud_resource.cloud_resource",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKmsCloudResource = `

resource "tencentcloud_kms_cloud_resource" "cloud_resource" {
  key_id = "23e80852-1e38-11e9-b129-5cb9019b4b01"
  product_id = "ssm"
  resource_id = "ins-123456"
}

`
