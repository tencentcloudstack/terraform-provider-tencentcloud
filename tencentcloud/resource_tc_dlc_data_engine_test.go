package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcDataEngineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDataEngine,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_data_engine.data_engine", "id")),
			},
			{
				ResourceName:      "tencentcloud_dlc_data_engine.data_engine",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcDataEngine = `

resource "tencentcloud_dlc_data_engine" "data_engine" {
  engine_type = "spark"
  data_engine_name = "testSpark"
  cluster_type = "spark_cu"
  mode = 2
  auto_resume = false
  min_clusters = 1
  max_clusters = 10
  default_data_engine = false
  cidr_block = "192.0.2.1/24"
  message = "test spark"
  pay_mode = 1
  time_span = 3600
  time_unit = "m"
  auto_renew = 0
  auto_suspend = false
  crontab_resume_suspend = 0
  crontab_resume_suspend_strategy {
		resume_time = "1000000-08:00:00"
		suspend_time = ""
		suspend_strategy = 

  }
  engine_exec_type = "SQL"
  max_concurrency = 5
  tolerable_queue_time = 0
  auto_suspend_time = 10
  resource_type = "Standard_CU"
  data_engine_config_pairs = 
  image_version_name = ""
  main_cluster_name = "testSpark"
  elastic_switch = false
  elastic_limit = 0
  session_resource_template {
		driver_size = "small"
		executor_size = "small"
		executor_nums = 1
		executor_max_numbers = 1

  }
}

`
