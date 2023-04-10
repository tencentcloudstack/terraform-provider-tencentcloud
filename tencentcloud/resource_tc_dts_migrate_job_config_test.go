package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDtsMigrateJobConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsMigrateJobConfig_pause(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "job_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "action"),
				),
			},
			{
				Config: testAccDtsMigrateJobConfig_continue(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "job_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job_config.config", "action", "continue"),
				),
			},
			// {
			// 	Config: testAccDtsMigrateJobConfig_complete(),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "id"),
			// 		resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "job_id"),
			// 		resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job_config.config", "action", "complete"),
			// 		resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job_config.config", "action", "immediately"),
			// 	),
			// },
			{
				Config: testAccDtsMigrateJobConfig_stop(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "job_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job_config.config", "action", "stop"),
				),
			},
			{
				Config: testAccDtsMigrateJobConfig_isolate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "job_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job_config.config", "action", "isolate"),
				),
			},
			{
				Config: testAccDtsMigrateJobConfig_recover(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_migrate_job_config.config", "job_id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_migrate_job_config.config", "action", "recover"),
				),
			},
		},
	})
}

func testAccDtsMigrateJobConfig_basic() string {
	curSec := fmt.Sprint(time.Now().Unix())
	return fmt.Sprintf(testAccDtsMigrateJob_basic, "migrate_job_config", curSec)
}

func testAccDtsMigrateJobConfig_pause() string {
	ret := fmt.Sprintf(testAccDtsMigrateJobConfig_basic() + `
	
resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "pause"
}
	
`)
	return ret
}

func testAccDtsMigrateJobConfig_continue() string {
	ret := fmt.Sprintf(testAccDtsMigrateJobConfig_basic() + `

resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "continue"
}

`)
	return ret
}

func testAccDtsMigrateJobConfig_stop() string {
	ret := fmt.Sprintf(testAccDtsMigrateJobConfig_basic() + `

resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "stop"
  complete_mode = "immediately"
}

`)
	return ret
}

func testAccDtsMigrateJobConfig_isolate() string {
	ret := fmt.Sprintf(testAccDtsMigrateJobConfig_basic() + `

resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "isolate"
  complete_mode = "immediately"
}

`)
	return ret
}

func testAccDtsMigrateJobConfig_recover() string {
	ret := fmt.Sprintf(testAccDtsMigrateJobConfig_basic() + `

resource "tencentcloud_dts_migrate_job_config" "config" {
  job_id = tencentcloud_dts_migrate_job_start_operation.start.id
  action = "recover"
  complete_mode = "immediately"
}

`)
	return ret
}
