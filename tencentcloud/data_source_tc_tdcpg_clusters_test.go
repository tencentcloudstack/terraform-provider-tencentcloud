package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTdcpgClustersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTdcpgClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgClusters_id, defaultTdcpgClusterId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_clusters.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.id", "list.0.cluster_id", defaultTdcpgClusterId),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.id", "list.0.cluster_name", defaultTdcpgClusterName),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.id", "list.0.region", defaultRegion),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.id", "list.0.zone", defaultTdcpgZone),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.db_version"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.id", "list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.status_desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.storage_used"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.storage_limit"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.id", "list.0.pay_mode", defaultTdcpgPayMode),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.auto_renew_flag"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.db_charset"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.instance_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.0.endpoint_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.0.endpoint_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.0.endpoint_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.0.private_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.endpoint_set.0.private_port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.db_major_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.id", "list.0.db_kernel_version"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.id", "list.0.pay_mode", defaultTdcpgPayMode),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgClusters_name, defaultTdcpgClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_clusters.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.name", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.name", "list.0.cluster_id", defaultTdcpgClusterId),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.name", "list.0.cluster_name", defaultTdcpgClusterName),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.name", "list.0.region", defaultRegion),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.name", "list.0.zone", defaultTdcpgZone),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.name", "list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.name", "list.0.instance_count"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.name", "list.0.pay_mode", defaultTdcpgPayMode),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgClusters_status, "running"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_clusters.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.status", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.status", "list.0.status", "running"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgClusters_paymode, defaultTdcpgPayMode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_clusters.paymode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.paymode", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.paymode", "list.0.pay_mode", defaultTdcpgPayMode),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceTdcpgClusters_project, "0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdcpg_clusters.project"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tdcpg_clusters.project", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdcpg_clusters.project", "list.0.project_id", "0"),
				),
			},
		},
	})
}

const testAccDataSourceTdcpgClusters_id = `

data "tencentcloud_tdcpg_clusters" "id" {
  cluster_id = "%s"
  }

`

const testAccDataSourceTdcpgClusters_name = `

data "tencentcloud_tdcpg_clusters" "name" {
  cluster_name = "%s"
  }

`

const testAccDataSourceTdcpgClusters_status = `

data "tencentcloud_tdcpg_clusters" "status" {
  status = "%s"
  }

`

const testAccDataSourceTdcpgClusters_paymode = `

data "tencentcloud_tdcpg_clusters" "paymode" {
  pay_mode = "%s"
  }

`

const testAccDataSourceTdcpgClusters_project = `

data "tencentcloud_tdcpg_clusters" "project" {
  project_id = "%s"
  }

`
