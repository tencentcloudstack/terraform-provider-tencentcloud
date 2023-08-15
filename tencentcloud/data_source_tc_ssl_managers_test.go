package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSslManagersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslManagersDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssl_managers.managers")),
			},
		},
	})
}

const testAccSslManagersDataSource = `

data "tencentcloud_ssl_managers" "managers" {
  company_id = 1
  manager_name = "leader"
  manager_mail = "xx@x.com"
  status = "none"
  search_key = "xxx"
}

`
