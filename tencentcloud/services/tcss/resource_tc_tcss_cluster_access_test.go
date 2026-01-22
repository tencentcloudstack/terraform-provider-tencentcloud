package tcss_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTcssClusterAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcssClusterAccess,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_cluster_access.example", "id"),
				),
			},
			{
				Config: testAccTcssClusterAccessUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_cluster_access.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcss_cluster_access.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcssClusterAccess = `
resource "tencentcloud_tcss_cluster_access" "example" {
  cluster_id = "cls-fdy7hm1q"
  switch_on  = true

  timeouts {
    create = "20m"
    update = "20m"
    delete = "20m"
  }
}
`

const testAccTcssClusterAccessUpdate = `
resource "tencentcloud_tcss_cluster_access" "example" {
  cluster_id = "cls-fdy7hm1q"
  switch_on  = false

  timeouts {
    create = "20m"
    update = "20m"
    delete = "20m"
  }
}
`
