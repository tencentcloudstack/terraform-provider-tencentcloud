package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcNativeSparkSessionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccDlcNativeSparkSessionsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dlc_native_spark_sessions.example"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_dlc_native_spark_sessions.example", "data_engine_id"),
				resource.TestCheckResourceAttrSet("data.tencentcloud_dlc_native_spark_sessions.example", "resource_group_id"),
			),
		}},
	})
}

const testAccDlcNativeSparkSessionsDataSource = `
data "tencentcloud_dlc_native_spark_sessions" "example" {
  data_engine_id    = "DataEngine-5plqp7q7"
  resource_group_id = "rg-j3zolzg77b"
}
`
