package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMariadbCancelDcnJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbCancelDcnJob,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_cancel_dcn_job.cancel_dcn_job", "id"),
				),
			},
		},
	})
}

const testAccMariadbCancelDcnJob = `
resource "tencentcloud_mariadb_cancel_dcn_job" "cancel_dcn_job" {
  instance_id = "tdsql-9vqvls95"
}
`
