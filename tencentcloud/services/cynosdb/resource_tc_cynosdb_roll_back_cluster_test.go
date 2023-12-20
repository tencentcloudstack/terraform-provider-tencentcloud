package cynosdb_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbRollBackClusterResource_basic -v
func TestAccTencentCloudCynosdbRollBackClusterResource_basic(t *testing.T) {

	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -1).In(loc).Format("2006-01-02 15:04:05")
	timeUnix := time.Now().AddDate(0, 0, -1).Unix()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCynosdbRollBackCluster, startTime, timeUnix, timeUnix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_roll_back_cluster.roll_back_cluster", "id"),
				),
			},
		},
	})
}

const testAccCynosdbRollBackCluster = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_roll_back_cluster" "roll_back_cluster" {
	cluster_id        =  var.cynosdb_cluster_id
	rollback_strategy = "snapRollback"
	rollback_id       = 732725
	expect_time = "%v"
	expect_time_thresh = 0
	rollback_databases {
	  old_database = "users"
	  new_database = "users_bak_%v"
	}
	rollback_tables {
	  database = "tf_ci_test"
	  tables {
		old_table = "test"
		new_table = "test_bak_%v"
	  }
  
	}
	rollback_mode = "full"
  }

`
