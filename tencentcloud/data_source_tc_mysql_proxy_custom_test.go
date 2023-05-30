package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMysqlProxyCustomDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlProxyCustomDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_proxy_custom.proxy_custom")),
			},
		},
	})
}

const testAccMysqlProxyCustomDataSourceVar = `
variable "instance_id" {
	default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlProxyCustomDataSource = testAccMysqlProxyCustomDataSourceVar + `

data "tencentcloud_mysql_proxy_custom" "proxy_custom" {
	instance_id = var.instance_id
}

`
