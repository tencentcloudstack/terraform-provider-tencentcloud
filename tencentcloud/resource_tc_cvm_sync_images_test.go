package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmSyncImagesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmSyncImages,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_sync_images.sync_images", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_sync_images.sync_images",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmSyncImages = `

resource "tencentcloud_cvm_sync_images" "sync_images" {
  image_ids = 
  destination_regions = 
  dry_run = false
  image_name = "img-evhmf3fy"
  image_set_required = false
}

`
