package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrImageManifestsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImageManifestsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_image_manifests.image_manifests")),
			},
		},
	})
}

const testAccTcrImageManifestsDataSource = `

data "tencentcloud_tcr_image_manifests" "image_manifests" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  repository_name = "repo"
  image_version = "v1"
    }

`
