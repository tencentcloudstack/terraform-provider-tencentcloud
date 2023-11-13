package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsSyncJobOperateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncJobOperate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job_operate.sync_job_operate", "id")),
			},
			{
				ResourceName:      "tencentcloud_dts_sync_job_operate.sync_job_operate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsSyncJobOperate = `

resource "tencentcloud_dts_sync_job_operate" "sync_job_operate" {
  job_id = "sync-werwfs23"
  new_instance_class = "large"
}

`
