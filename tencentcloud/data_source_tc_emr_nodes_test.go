package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudEMRNodes(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEMRNodes,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_emr_nodes.my_emr_nodes"),
					resource.TestCheckResourceAttr("data.tencentcloud_emr_nodes.my_emr_nodes", "nodes.#", "1"),
				),
			},
		},
	})
}

const testAccEMRNodes = `

resource "tencentcloud_emr_cluster" "emrrrr" {
  product_id=4
  display_strategy="clusterList"
  vpc_settings={
    vpc_id="vpc-4owdpnwr"
    subnet_id:"subnet-ahv6swf2"
  }
  softwares=["zookeeper-3.6.1"]
  support_ha=0
  instance_name="emr-test"
  resource_spec {
    master_resource_spec {
  	mem_size=8192
  	cpu=4
  	disk_size=100
  	disk_type="CLOUD_PREMIUM"
  	spec="CVM.S2"
  	storage_type=5
    }
    core_resource_spec {
  	mem_size=8192
  	cpu=4
  	disk_size=100
  	disk_type="CLOUD_PREMIUM"
  	spec="CVM.S2"
  	storage_type=5
    }
    master_count=1
    core_count=2
  }
  login_settings={
    password="tencent@cloud123"
  }
  time_span=3600
  time_unit="s"
  pay_mode=0
  placement={
    zone="ap-guangzhou-3"
    project_id=0
  }
  sg_id="sg-qyy7jd2b"
}
data "tencentcloud_emr_nodes" "my_emr_nodes" {
  node_flag="master"
  instance_id=tencentcloud_emr_cluster.emrrrr.instance_id
}
`
