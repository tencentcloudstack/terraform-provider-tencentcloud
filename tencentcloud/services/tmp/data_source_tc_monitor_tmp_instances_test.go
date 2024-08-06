package tmp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMonitorTmpInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_tmp_instances.tmp_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.instance_charge_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.region_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.data_retention_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.instance_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.created_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.enable_grafana"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.ipv4_address"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.tag_specification.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.expire_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.charge_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.auto_renew_flag"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.is_near_expire"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.auth_token"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.remote_write"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.api_root_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.proxy_address"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.grafana_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.grant.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.alert_rule_limit"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.recording_rule_limit"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.migration_type"),
					resource.TestCheckResourceAttr("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.grafana_instance_id", ""),
					resource.TestCheckResourceAttr("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.grafana_url", ""),
					resource.TestCheckResourceAttr("data.tencentcloud_monitor_tmp_instances.tmp_instances", "instance_set.0.grafana_ip_white_list", ""),
				),
			},
		},
	})
}

const testAccMonitorTmpInstancesDataSource = `
variable "availability_zone" {
	default = "ap-guangzhou-4"
}
  
resource "tencentcloud_vpc" "vpc" {
	cidr_block = "10.0.0.0/16"
	name       = "tf_monitor_tmp_vpc"
}
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	availability_zone = var.availability_zone
	name              = "tf_monitor_tmp_subnet"
	cidr_block        = "10.0.1.0/24"
}
  
resource "tencentcloud_monitor_tmp_instance" "example" {
	instance_name       = "tf-monitor-tmp-instance"
	vpc_id              = tencentcloud_vpc.vpc.id
	subnet_id           = tencentcloud_subnet.subnet.id
	data_retention_time = 30
	zone                = var.availability_zone
}
  
data "tencentcloud_monitor_tmp_instances" "tmp_instances" {
	instance_ids = [tencentcloud_monitor_tmp_instance.example.id]
}
`
