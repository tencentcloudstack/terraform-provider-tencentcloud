
package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDtsMigrateJob_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateJob,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job.migrate_job", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dts_migrate_job.migrateJob",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDtsMigrateJob = `

resource "tencentcloud_dts_migrate_job" "migrate_job" {
  src_database_type = ""
  dst_database_type = ""
  src_region = ""
  dst_region = ""
  instance_class = ""
  count = ""
  job_name = ""
  tags {
			tag_key = ""
			tag_value = ""

  }
            }

`
