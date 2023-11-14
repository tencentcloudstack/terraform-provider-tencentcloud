package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresBackupPlanResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresBackupPlan,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_backup_plan.backup_plan", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_backup_plan.backup_plan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresBackupPlan = `

resource "tencentcloud_postgres_backup_plan" "backup_plan" {
  d_b_instance_id = "postgres-xxxxx"
  min_backup_start_time = "01:00:00"
  max_backup_start_time = "02:00:00"
  base_backup_retention_period = 7
  backup_period = 
}

`
