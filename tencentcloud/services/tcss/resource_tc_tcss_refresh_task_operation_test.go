package tcss_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTcssRefreshTaskOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcssRefreshTaskOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcss_refresh_task_operation.example", "id"),
				),
			},
		},
	})
}

const testAccTcssRefreshTaskOperation = `
resource "tencentcloud_tcss_refresh_task_operation" "example" {
  cluster_ids = [
    "cls-fdy7hm1q"
  ]
  is_sync_list_only = false
}
`
