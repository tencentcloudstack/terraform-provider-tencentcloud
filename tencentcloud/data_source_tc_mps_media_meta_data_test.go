package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsMediaMetaDataDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsMediaMetaDataDataSource, defaultRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_media_meta_data.metadata"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "input_info.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_media_meta_data.metadata", "input_info.0.type", "COS"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "input_info.0.cos_input_info.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "input_info.0.cos_input_info.0.bucket"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_media_meta_data.metadata", "input_info.0.cos_input_info.0.region", defaultRegion),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_media_meta_data.metadata", "input_info.0.cos_input_info.0.object", "/mps-test/test.mov"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "meta_data.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "meta_data.0.size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "meta_data.0.container"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "meta_data.0.bitrate"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "meta_data.0.height"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "meta_data.0.width"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_media_meta_data.metadata", "meta_data.0.duration"),
				),
			},
		},
	})
}

const testAccMpsMediaMetaDataDataSource = userInfoData + `
data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

data "tencentcloud_mps_media_meta_data" "metadata" {
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
