package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testTcrImageManifestsObjectName = "data.tencentcloud_tcr_image_manifests.image_manifests"

func TestAccTencentCloudTcrImageManifestsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImageManifestsDataSource,
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(testTcrImageManifestsObjectName),
					resource.TestCheckResourceAttrSet(testTcrImageManifestsObjectName, "id"),
					resource.TestCheckResourceAttrSet(testTcrImageManifestsObjectName, "registry_id"),
					resource.TestCheckResourceAttrSet(testTcrImageManifestsObjectName, "namespace_name"),
					resource.TestCheckResourceAttrSet(testTcrImageManifestsObjectName, "repository_name"),
				),
			},
		},
	})
}

const testAccTcrImageManifestsDataSource = TCRDataSource + `

data "tencentcloud_tcr_image_manifests" "image_manifests" {
	registry_id = local.tcr_id
	namespace_name = local.tcr_ns_name
	repository_name = local.tcr_repo
	image_version = "vv1"
}

`
