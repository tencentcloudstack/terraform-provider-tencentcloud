package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbCancelDcnJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbCancelDcnJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_cancel_dcn_job.cancel_dcn_job", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_cancel_dcn_job.cancel_dcn_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbCancelDcnJob = `

resource "tencentcloud_dcdb_cancel_dcn_job" "cancel_dcn_job" {
  instance_id = ""
}

`
