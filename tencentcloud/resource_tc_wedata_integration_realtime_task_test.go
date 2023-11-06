package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataIntegrationRealtimeTaskResource_basic -v
func TestAccTencentCloudNeedFixWedataIntegrationRealtimeTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataIntegrationRealtimeTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "task_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "sync_type"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_integration_realtime_task.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataIntegrationRealtimeTaskUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "task_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_integration_realtime_task.example", "sync_type"),
				),
			},
		},
	})
}

const testAccWedataIntegrationRealtimeTask = `
resource "tencentcloud_wedata_integration_realtime_task" "example" {
  project_id  = "1455251608631480391"
  task_name   = "tf_example"
  description = "description."
  sync_type   = 1
  task_info {
    incharge    = "100028439226"
    executor_id = "20230704142425553913"
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "TaskManager"
      value = "1"
    }
    config {
      name  = "JobManager"
      value = "1"
    }
    config {
      name  = "TolerateDirtyData"
      value = "0"
    }
    config {
      name  = "CheckpointingInterval"
      value = "1"
    }
    config {
      name  = "CheckpointingIntervalUnit"
      value = "min"
    }
    config {
      name  = "RestartStrategyFixedDelayAttempts"
      value = "-1"
    }
    config {
      name  = "ResourceAllocationType"
      value = "0"
    }
    config {
      name  = "TaskAlarmRegularList"
      value = ""
    }
    mappings {
      source_id = "2"
      sink_id   = "1"
    }
    nodes {
      id               = "1"
      name             = "gf_poc"
      node_type        = "INPUT"
      data_source_type = "MYSQL"
      datasource_id    = "5737"
      config {
        name  = "StartupMode"
        value = "INIT"
      }
      config {
        name  = "Encode"
        value = "utf-8"
      }
      config {
        name  = "Database"
        value = "UNKNOW"
      }
      config {
        name  = "SourceRule"
        value = "all"
      }
      config {
        name  = "FilterOper"
        value = "update"
      }
      config {
        name  = "ServerTimeZone"
        value = "Asia/Shanghai"
      }
      config {
        name  = "GhostChange"
        value = "false"
      }
      config {
        name  = "TableNames"
        value = "gf_db.*,hx_db.*,information_schema.*,mysql.*,performance_schema.*,run_time.*,sys.*,test01.*"
      }
      config {
        name  = "FirstDataSource"
        value = "5737"
      }
      config {
        name  = "MultipleDataSources"
        value = "5737"
      }
      config {
        name  = "SiblingNodes"
        value = "[]"
      }
    }
  }
}
`

const testAccWedataIntegrationRealtimeTaskUpdate = `
resource "tencentcloud_wedata_integration_realtime_task" "example" {
  project_id  = "1455251608631480391"
  task_name   = "tf_example_update"
  description = "description update."
  sync_type   = 1
  task_info {
    incharge    = "100028439226"
    executor_id = "20230704142425553913"
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "TaskManager"
      value = "1"
    }
    config {
      name  = "JobManager"
      value = "1"
    }
    config {
      name  = "TolerateDirtyData"
      value = "0"
    }
    config {
      name  = "CheckpointingInterval"
      value = "1"
    }
    config {
      name  = "CheckpointingIntervalUnit"
      value = "min"
    }
    config {
      name  = "RestartStrategyFixedDelayAttempts"
      value = "-1"
    }
    config {
      name  = "ResourceAllocationType"
      value = "0"
    }
    config {
      name  = "TaskAlarmRegularList"
      value = ""
    }
    mappings {
      source_id = "2"
      sink_id   = "1"
    }
    nodes {
      id               = "1"
      name             = "gf_poc"
      node_type        = "INPUT"
      data_source_type = "MYSQL"
      datasource_id    = "5737"
      config {
        name  = "StartupMode"
        value = "INIT"
      }
      config {
        name  = "Encode"
        value = "utf-8"
      }
      config {
        name  = "Database"
        value = "UNKNOW"
      }
      config {
        name  = "SourceRule"
        value = "all"
      }
      config {
        name  = "FilterOper"
        value = "update"
      }
      config {
        name  = "ServerTimeZone"
        value = "Asia/Shanghai"
      }
      config {
        name  = "GhostChange"
        value = "false"
      }
      config {
        name  = "TableNames"
        value = "gf_db.*,hx_db.*,information_schema.*,mysql.*,performance_schema.*,run_time.*,sys.*,test01.*"
      }
      config {
        name  = "FirstDataSource"
        value = "5737"
      }
      config {
        name  = "MultipleDataSources"
        value = "5737"
      }
      config {
        name  = "SiblingNodes"
        value = "[]"
      }
    }
  }
}
`
