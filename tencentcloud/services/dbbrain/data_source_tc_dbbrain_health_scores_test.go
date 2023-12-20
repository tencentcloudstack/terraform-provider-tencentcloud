package dbbrain_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainHealthScoresDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	reportTime := time.Now().Add(-1 * time.Hour).In(loc).Format("2006-01-02 15:04:05")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainHealthScoresDataSource, reportTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_health_scores.health_scores"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_health_scores.health_scores", "time", reportTime),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_health_scores.health_scores", "product", "mysql"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "data.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "data.0.events_total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "data.0.health_score"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "data.0.health_level"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "data.0.issue_types.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "data.0.issue_types.0.issue_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "data.0.issue_types.0.events.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_health_scores.health_scores", "data.0.issue_types.0.total_count"),
				),
			},
		},
	})
}

const testAccDbbrainHealthScoresDataSource = tcacctest.CommonPresetMysql + `

data "tencentcloud_dbbrain_health_scores" "health_scores" {
  instance_id = local.mysql_id
  time = "%s"
  product = "mysql"
  }

`
