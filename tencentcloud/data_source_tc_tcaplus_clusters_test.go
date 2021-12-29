package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTcaplusClustersName = "data.tencentcloud_tcaplus_clusters.id_test"

func TestAccTencentCloudDataTcaplusClusters(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusClustersBaic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTcaplusClusterExists("tencentcloud_tcaplus_cluster.test_cluster"),
					resource.TestCheckResourceAttrSet(testDataTcaplusClustersName, "cluster_id"),
					resource.TestCheckResourceAttr(testDataTcaplusClustersName, "list.#", "1"),
					resource.TestCheckResourceAttr(testDataTcaplusClustersName, "list.0.cluster_name", "tf_tcaplus_data_guagua"),
					resource.TestCheckResourceAttr(testDataTcaplusClustersName, "list.0.idl_type", "PROTO"),
					resource.TestCheckResourceAttrSet(testDataTcaplusClustersName, "list.0.network_type"),
					resource.TestCheckResourceAttrSet(testDataTcaplusClustersName, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataTcaplusClustersName, "list.0.password_status"),
					resource.TestCheckResourceAttrSet(testDataTcaplusClustersName, "list.0.api_access_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusClustersName, "list.0.api_access_ip"),
					resource.TestCheckResourceAttrSet(testDataTcaplusClustersName, "list.0.api_access_port"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTcaplusClustersBaic = `
variable "availability_zone" {
 default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "test_cluster" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_tcaplus_data_guagua"
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id                = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}
data "tencentcloud_tcaplus_clusters" "id_test" {
  cluster_id = tencentcloud_tcaplus_cluster.test_cluster.id
}
`
