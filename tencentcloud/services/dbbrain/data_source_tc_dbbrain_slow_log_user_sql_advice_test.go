package dbbrain_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainSlowLogUserSqlAdviceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainSlowLogUserSqlAdviceDataSource, tcacctest.DefaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_slow_log_user_sql_advice.test"),
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
