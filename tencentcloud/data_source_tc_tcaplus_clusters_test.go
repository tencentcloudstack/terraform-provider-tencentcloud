package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTcaplusClustersName = "data.tencentcloud_tcaplus_clusters.tcaplus"

func TestAccTencentCloudTcaplusClustersData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTcaplusClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusClustersBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testDataTcaplusClustersName, "cluster_name", defaultTcaPlusClusterName),
					resource.TestCheckResourceAttr(testDataTcaplusClustersName, "list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataTcaplusClustersName, "list.0.cluster_id"),
					resource.TestCheckResourceAttr(testDataTcaplusClustersName, "list.0.cluster_name", defaultTcaPlusClusterName),
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

const testAccTencentCloudDataTcaplusClustersBasic = defaultTcaPlusData
