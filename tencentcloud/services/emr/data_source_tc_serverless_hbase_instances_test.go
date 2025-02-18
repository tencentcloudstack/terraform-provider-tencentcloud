package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudServerlessHbaseInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessHbaseInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_serverless_hbase_instances.serverless_hbase_instances"),
					resource.TestCheckResourceAttr("data.tencentcloud_serverless_hbase_instances.serverless_hbase_instances", "instance_list.#", "1"),
				),
			},
		},
	})
}

const testAccServerlessHbaseInstancesDataSource = `
resource "tencentcloud_serverless_hbase_instance" "serverless_hbase_instance" {
  instance_name = "tf-test-datasource"
  pay_mode = 0
  disk_type = "CLOUD_HSSD"
  disk_size = 100
  node_type = "4C16G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id = "vpc-muytmxhk"
      subnet_id = "subnet-9ye3xm5v"
    }
    node_num = 3
  }
}

data "tencentcloud_serverless_hbase_instances" "serverless_hbase_instances" {
  display_strategy = "clusterList"
  filters {
	name = "ClusterId"
	values = [tencentcloud_serverless_hbase_instance.serverless_hbase_instance.id]
  }
}
`
