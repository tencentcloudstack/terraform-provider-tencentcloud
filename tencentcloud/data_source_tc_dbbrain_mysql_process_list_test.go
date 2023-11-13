package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainMysqlProcessListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainMysqlProcessListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_mysql_process_list.mysql_process_list")),
			},
		},
	})
}

const testAccDbbrainMysqlProcessListDataSource = `

data "tencentcloud_dbbrain_mysql_process_list" "mysql_process_list" {
  instance_id = ""
  i_d = 
  user = ""
  host = ""
  d_b = ""
  state = ""
  command = ""
  time = 
  info = ""
  product = ""
  }

`
