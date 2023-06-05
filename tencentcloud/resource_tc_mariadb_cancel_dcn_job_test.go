package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMariadbCancelDcnJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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
