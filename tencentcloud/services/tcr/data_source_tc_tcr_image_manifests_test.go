package tcr_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testTcrImageManifestsObjectName = "data.tencentcloud_tcr_image_manifests.image_manifests"

func TestAccTencentCloudTcrImageManifestsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrImageManifestsDataSource, tcacctest.DefaultTCRInstanceId, tcacctest.DefaultTCRNamespace, tcacctest.DefaultTCRRepoName),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(testTcrImageManifestsObjectName),
					resource.TestCheckResourceAttrSet(testTcrImageManifestsObjectName, "id"),
					resource.TestCheckResourceAttr(testTcrImageManifestsObjectName, "registry_id", tcacctest.DefaultTCRInstanceId),
					resource.TestCheckResourceAttr(testTcrImageManifestsObjectName, "namespace_name", tcacctest.DefaultTCRNamespace),
					resource.TestCheckResourceAttr(testTcrImageManifestsObjectName, "repository_name", tcacctest.DefaultTCRRepoName),
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
	image_version = "vv1"
}

`
