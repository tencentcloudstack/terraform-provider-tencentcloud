package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainSlowLogUserSqlAdviceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSlowLogUserSqlAdviceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_slow_log_user_sql_advice.slow_log_user_sql_advice")),
			},
		},
	})
}

const testAccDbbrainSlowLogUserSqlAdviceDataSource = `

data "tencentcloud_dbbrain_slow_log_user_sql_advice" "slow_log_user_sql_advice" {
  instance_id = ""
  sql_text = ""
  schema = ""
  product = ""
          }

`
