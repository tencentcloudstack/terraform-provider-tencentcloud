package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEsDiagnoseResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEsDiagnose,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_diagnose.diagnose", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_diagnose.diagnose", "diagnose_job_metas.#"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_diagnose.diagnose", "cron_time", "15:00:00"),
				),
			},
			{
				Config: testAccEsDiagnoseUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_diagnose.diagnose", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_elasticsearch_diagnose.diagnose", "diagnose_job_metas.#"),
					resource.TestCheckResourceAttr("tencentcloud_elasticsearch_diagnose.diagnose", "cron_time", "16:00:00"),
				),
			},
			{
				ResourceName:      "tencentcloud_elasticsearch_diagnose.diagnose",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccEsDiagnose = `
resource "tencentcloud_elasticsearch_diagnose" "diagnose" {
	instance_id = "es-nni6pm4s"
	cron_time = "15:00:00"
}
`

const testAccEsDiagnoseUpdate = `
resource "tencentcloud_elasticsearch_diagnose" "diagnose" {
	instance_id = "es-nni6pm4s"
	cron_time = "16:00:00"
}
`
