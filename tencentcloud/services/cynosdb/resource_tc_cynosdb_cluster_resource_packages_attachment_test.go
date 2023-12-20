package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCynosdbClusterResourcePackagesAttachmentResource_basic -v
func TestAccTencentCloudNeedFixCynosdbClusterResourcePackagesAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterResourcePackagesAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_resource_packages_attachment.cluster_resource_packages_attachment", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_resource_packages_attachment.cluster_resource_packages_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbClusterResourcePackagesAttachment = `
resource "tencentcloud_cynosdb_cluster_resource_packages_attachment" "cluster_resource_packages_attachment" {
  cluster_id  = "cynosdbmysql-q1d8151n"
  package_ids = ["package-hy4d2ppl"]
}
`
