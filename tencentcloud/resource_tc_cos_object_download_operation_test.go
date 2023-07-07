package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosObjectDownloadOperationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosObjectDownloadOperation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_object_download_operation.object_download", "id"),
				),
			},
		},
	})
}

const testAccCosObjectDownloadOperation = `
resource "tencentcloud_cos_object_download_operation" "object_download" {
    bucket = "keep-test-1308919341"
    key = "download.txt"
    download_path = "./download.txt"
}
`
