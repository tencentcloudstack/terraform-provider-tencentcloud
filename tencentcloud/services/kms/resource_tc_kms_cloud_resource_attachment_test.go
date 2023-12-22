package kms_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudKmsCloudResourceAttachmentResource_basic -v
func TestAccTencentCloudKmsCloudResourceAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsCloudResourceAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kms_cloud_resource_attachment.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kms_cloud_resource_attachment.example", "key_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kms_cloud_resource_attachment.example", "product_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kms_cloud_resource_attachment.example", "resource_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_kms_cloud_resource_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKmsCloudResourceAttachment = `
resource "tencentcloud_kms_cloud_resource_attachment" "example" {
  key_id      = "72688f39-1fe8-11ee-9f1a-525400cf25a4"
  product_id  = "mysql"
  resource_id = "cdb-fitq5t9h"
}
`
