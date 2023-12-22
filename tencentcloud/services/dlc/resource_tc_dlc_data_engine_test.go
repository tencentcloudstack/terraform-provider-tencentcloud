package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcDataEngineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDataEngine,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_data_engine.data_engine", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "engine_type", "spark"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "data_engine_name", "testSpark"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "cluster_type", "spark_cu"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "auto_resume", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "size", "16"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "min_clusters", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "max_clusters", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "default_data_engine", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "cidr_block", "10.255.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "message", "test spark1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "time_span", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "time_unit", "h"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "auto_suspend", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "crontab_resume_suspend", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "engine_exec_type", "BATCH"),
				),
			}, {
				Config: testAccDlcDataEngineUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_data_engine.data_engine", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "engine_type", "spark"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "data_engine_name", "testSpark"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "cluster_type", "spark_cu"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "mode", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "auto_resume", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "size", "16"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "min_clusters", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "max_clusters", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "default_data_engine", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "cidr_block", "10.255.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "message", "test spark"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "time_span", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "time_unit", "h"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "auto_suspend", "false"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "crontab_resume_suspend", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_data_engine.data_engine", "engine_exec_type", "BATCH"),
				),
			},

			{
				ResourceName:            "tencentcloud_dlc_data_engine.data_engine",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pay_mode", "size", "time_span", "time_unit"},
			},
		},
	})
}

const testAccDlcDataEngine = `

resource "tencentcloud_dlc_data_engine" "data_engine" {
  engine_type = "spark"
  data_engine_name = "testSpark"
  cluster_type = "spark_cu"
  mode = 1
  auto_resume = false
  size = 16
  pay_mode = 0
  min_clusters = 1
  max_clusters = 1
  default_data_engine = false
  cidr_block = "10.255.0.0/16"
  message = "test spark1"
  time_span = 1
  time_unit = "h"
  auto_suspend = false
  crontab_resume_suspend = 0
  engine_exec_type = "BATCH"
}

`
const testAccDlcDataEngineUpdate = `

resource "tencentcloud_dlc_data_engine" "data_engine" {
  engine_type = "spark"
  data_engine_name = "testSpark"
  cluster_type = "spark_cu"
  mode = 1
  auto_resume = false
  size = 16
  pay_mode = 0
  min_clusters = 1
  max_clusters = 1
  default_data_engine = false
  cidr_block = "10.255.0.0/16"
  message = "test spark"
  time_span = 1
  time_unit = "h"
  auto_suspend = false
  crontab_resume_suspend = 0
  engine_exec_type = "BATCH"
}

`
