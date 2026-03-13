package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMonitorExternalClusterResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorExternalCluster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_external_cluster.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_external_cluster.example", "cluster_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_external_cluster.example", "cluster_name", "tf-external-cluster"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_external_cluster.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorExternalCluster = `
resource "tencentcloud_monitor_external_cluster" "example" {
  instance_id    = "prom-gzg3f1em"
  cluster_region = "ap-guangzhou"
  cluster_name   = "tf-external-cluster"

  external_labels {
    name  = "clusterName"
    value = "example"
  }

  external_labels {
    name  = "environment"
    value = "prod"
  }

  enable_external = false
}
`
