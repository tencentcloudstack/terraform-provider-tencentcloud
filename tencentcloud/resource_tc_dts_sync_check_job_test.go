package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsSyncCheckJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncCheckJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_check_job.sync_check_job", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_sync_check_job.sync_check_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsSyncCheckJob = `

resource "tencentcloud_dts_sync_check_job" "sync_check_job" {
  job_id = &lt;nil&gt;
}

`
