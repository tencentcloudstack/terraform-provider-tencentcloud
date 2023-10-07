package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsDescribeMediaMetadataOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsDescribeMediaMetadataOperation, defaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_describe_media_metadata_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_describe_media_metadata_operation.operation", "input_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_describe_media_metadata_operation.operation", "input_info.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_describe_media_metadata_operation.operation", "input_info.0.cos_input_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_describe_media_metadata_operation.operation", "input_info.0.cos_input_info.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_describe_media_metadata_operation.operation", "input_info.0.cos_input_info.0.region", defaultRegion),
					resource.TestCheckResourceAttr("tencentcloud_mps_describe_media_metadata_operation.operation", "input_info.0.cos_input_info.0.object", "/mps-test/test.mov"),
				),
			},
		},
	})
}

const testAccMpsDescribeMediaMetadataOperation = userInfoData + `
data "tencentcloud_cos_bucket_object" "object" {
	bucket = "keep-bucket-${local.app_id}"
	key    = "/mps-test/test.mov"
}

resource "tencentcloud_mps_describe_media_metadata_operation" "operation" {
  input_info {
		type = "COS"
		cos_input_info {
			bucket = data.tencentcloud_cos_bucket_object.object.bucket
			region = "%s"
			object = data.tencentcloud_cos_bucket_object.object.key
		}
  }
}

`
