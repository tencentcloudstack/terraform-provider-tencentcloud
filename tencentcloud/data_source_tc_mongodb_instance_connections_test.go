package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMongodbInstanceConnectionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceConnectionsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_connections.instance_connections")),
			},
		},
	})
}

const testAccMongodbInstanceConnectionsDataSource = `

data "tencentcloud_mongodb_instance_connections" "instance_connections" {
  instance_id = "cmgo-9d0p6umb"
  }

`
