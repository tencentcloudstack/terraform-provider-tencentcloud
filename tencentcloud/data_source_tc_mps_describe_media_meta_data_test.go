package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsDescribeMediaMetaDataDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsDescribeMediaMetaDataDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_describe_media_meta_data.describe_media_meta_data")),
			},
		},
	})
}

const testAccMpsDescribeMediaMetaDataDataSource = `

data "tencentcloud_mps_describe_media_meta_data" "describe_media_meta_data" {
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
