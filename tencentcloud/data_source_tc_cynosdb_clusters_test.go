package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbClustersDataSource_full(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterstDataSource_full(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_cynosdb_cluster.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.cluster_name", "tf-cynosdb"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.cluster_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.available_zone", "ap-guangzhou-4"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.db_type", "MYSQL"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.db_version", "5.7"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.cluster_status"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.port"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_clusters.cluster_full", "cluster_list.0.auto_renew_flag", "0"),
				),
			},
		},
	})
}

func testAccCynosdbClusterstDataSource_full() string {
	return testAccCynosdbCluster + `
data "tencentcloud_cynosdb_clusters" "cluster_full" {
  cluster_id   = tencentcloud_cynosdb_cluster.foo.id
  project_id   = 0
  db_type      = "MYSQL"
  cluster_name = "tf-cynosdb"
}`
}
