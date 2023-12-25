package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverPublishSubscribeDataSource -v
func TestAccTencentCloudSqlserverPublishSubscribeDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSqlserverPublishSubscribeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudSqlServerPublishSubscribeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_sqlserver_publish_subscribes.example", "publish_subscribe_list.#"),
				),
			},
		},
	})
}

const testAccTencentCloudSqlServerPublishSubscribeDataSourceConfig = tcacctest.CommonPresetSQLServer + `
data "tencentcloud_sqlserver_publish_subscribes" "example" {
  instance_id = local.sqlserver_id
}
`
