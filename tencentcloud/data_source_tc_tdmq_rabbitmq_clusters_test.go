package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqRabbitmqClustersDataSource -v
func TestAccTencentCloudTdmqRabbitmqClustersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdmqRabbitmqClusters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rabbitmq_clusters.rabbitmq_clusters"),
				),
			},
		},
	})
}

const testAccDataSourceTdmqRabbitmqClusters = `

data "tencentcloud_tdmq_rabbitmq_clusters" "rabbitmq_clusters" {
}

`
