package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Config: fmt.Sprintf(testAccDbbrainSlowLogUserSqlAdviceDataSource, defaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_slow_log_user_sql_advice.test"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_user_sql_advice.test", "sql_text", "select * from tf_test_ci where 1=1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_user_sql_advice.test", "product", "mysql"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_user_sql_advice.test", "advices"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_user_sql_advice.test", "comments"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_user_sql_advice.test", "tables"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_user_sql_advice.test", "sql_plan"),
				),
			},
		},
	})
}

const testAccDbbrainSlowLogUserSqlAdviceDataSource = `

data "tencentcloud_dbbrain_slow_log_user_sql_advice" "test" {
  instance_id = "%s"
  sql_text = "select * from tf_test_ci where 1=1"
  product = "mysql"
}

`
