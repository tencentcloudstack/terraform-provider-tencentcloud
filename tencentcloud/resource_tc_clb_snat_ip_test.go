package tencentcloud

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudClbSnatIp(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbSnatIpBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clb_snat_ip.snat_ips", "id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_snat_ip.snat_ips", "ips.#", "3"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 10)
				},
				Config: testAccClbSnatIpBasicUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clb_snat_ip.snat_ips", "id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_snat_ip.snat_ips", "ips.#", "3"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 10)
				},
				Config: testAccClbSnatIpBasicUpdate2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clb_snat_ip.snat_ips", "id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_snat_ip.snat_ips", "ips.#", "1"),
				),
			},
		},
	})
}

const testAccClbSnatIpBasic = `
data "tencentcloud_vpc_instances" "gz3vpc" {
  name = "Default-"
  is_default = true
}

data "tencentcloud_vpc_subnets" "gz3" {
  vpc_id = data.tencentcloud_vpc_instances.gz3vpc.instance_list.0.vpc_id
}

locals {
  keep_clb_subnets = [for subnet in data.tencentcloud_vpc_subnets.gz3.instance_list: lookup(subnet, "subnet_id") if lookup(subnet, "name") == "keep-clb-sub"]
  subnets = [for subnet in data.tencentcloud_vpc_subnets.gz3.instance_list: lookup(subnet, "subnet_id") ]
  subnet_for_clb_snat = concat(local.keep_clb_subnets, local.subnets)
}

resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = "tf-clb-snat-resource-test"
}

resource "tencentcloud_clb_snat_ip" "snat_ips" {
  clb_id = tencentcloud_clb_instance.foo.id
  ips {
    ip = "172.16.151.17"
	subnet_id = local.subnet_for_clb_snat.0
  }
  ips {
	ip = "172.16.151.15"
	subnet_id = local.subnet_for_clb_snat.0
  }
  ips {
	ip = "172.16.151.138"
	subnet_id = local.subnet_for_clb_snat.1
  }
}
`

const testAccClbSnatIpBasicUpdate = `
data "tencentcloud_vpc_instances" "gz3vpc" {
  name = "Default-"
  is_default = true
}

data "tencentcloud_vpc_subnets" "gz3" {
  vpc_id = data.tencentcloud_vpc_instances.gz3vpc.instance_list.0.vpc_id
}

locals {
  keep_clb_subnets = [for subnet in data.tencentcloud_vpc_subnets.gz3.instance_list: lookup(subnet, "subnet_id") if lookup(subnet, "name") == "keep-clb-sub"]
  subnets = [for subnet in data.tencentcloud_vpc_subnets.gz3.instance_list: lookup(subnet, "subnet_id") ]
  subnet_for_clb_snat = concat(local.keep_clb_subnets, local.subnets)
}

resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = "tf-clb-snat-resource-test"
}

resource "tencentcloud_clb_snat_ip" "snat_ips" {
  clb_id = tencentcloud_clb_instance.foo.id
  ips {
    ip = "172.16.151.17"
	subnet_id = local.subnet_for_clb_snat.0
  }
  ips {
	ip = "172.16.151.138"
	subnet_id = local.subnet_for_clb_snat.1
  }
  ips {
	ip = "172.16.151.139"
    subnet_id = local.subnet_for_clb_snat.1
  }
}
`

const testAccClbSnatIpBasicUpdate2 = `
data "tencentcloud_vpc_instances" "gz3vpc" {
  name = "Default-"
  is_default = true
}

data "tencentcloud_vpc_subnets" "gz3" {
  vpc_id = data.tencentcloud_vpc_instances.gz3vpc.instance_list.0.vpc_id
}

locals {
  keep_clb_subnets = [for subnet in data.tencentcloud_vpc_subnets.gz3.instance_list: lookup(subnet, "subnet_id") if lookup(subnet, "name") == "keep-clb-sub"]
  subnets = [for subnet in data.tencentcloud_vpc_subnets.gz3.instance_list: lookup(subnet, "subnet_id") ]
  subnet_for_clb_snat = concat(local.keep_clb_subnets, local.subnets)
}

resource "tencentcloud_clb_instance" "foo" {
  network_type = "OPEN"
  clb_name     = "tf-clb-snat-resource-test"
}

resource "tencentcloud_clb_snat_ip" "snat_ips" {
  clb_id = tencentcloud_clb_instance.foo.id
  ips {
    ip = "172.16.151.16"
	subnet_id = local.subnet_for_clb_snat.0
  }
}
`
