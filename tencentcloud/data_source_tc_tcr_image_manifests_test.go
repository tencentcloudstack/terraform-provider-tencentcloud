package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testTcrImageManifestsObjectName = "data.tencentcloud_tcr_image_manifests.image_manifests"

func TestAccTencentCloudTcrImageManifestsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrImageManifestsDataSource, defaultTCRInstanceId, defaultTCRNamespace, defaultTCRRepoName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testTcrImageManifestsObjectName),
					resource.TestCheckResourceAttrSet(testTcrImageManifestsObjectName, "id"),
					resource.TestCheckResourceAttr(testTcrImageManifestsObjectName, "registry_id", defaultTCRInstanceId),
					resource.TestCheckResourceAttr(testTcrImageManifestsObjectName, "namespace_name", defaultTCRNamespace),
					resource.TestCheckResourceAttr(testTcrImageManifestsObjectName, "repository_name", defaultTCRRepoName),
				),
			},
		},
	})
}

const testAccTcrImageManifestsDataSource = `

data "tencentcloud_tcr_image_manifests" "image_manifests" {
	registry_id = "%s"
	namespace_name = "%s" 
	repository_name = "%s"
	image_version = "v1"
}

`
