package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Config: testAccMpsDescribeMediaMetadataOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_describe_media_metadata_operation.describe_media_metadata_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_describe_media_metadata_operation.describe_media_metadata_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsDescribeMediaMetadataOperation = `

resource "tencentcloud_mps_describe_media_metadata_operation" "describe_media_metadata_operation" {
  input_info {
		type = ""
		cos_input_info {
			bucket = ""
			region = ""
			object = ""
		}
		url_input_info {
			url = ""
		}
		s3_input_info {
			s3_bucket = ""
			s3_region = ""
			s3_object = ""
			s3_secret_id = ""
			s3_secret_key = ""
		}

  }
}

`
