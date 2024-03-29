package tmp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInternationalMonitorResource_tmpInstance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMonInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testInternationalInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("tencentcloud_monitor_tmp_instance.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.basic", "instance_name", "demo-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.basic", "data_retention_time", "30"),
				),
			},
			{
				Config: testInternationalInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("tencentcloud_monitor_tmp_instance.update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.update", "instance_name", "demo-test-update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.update", "data_retention_time", "30"),
				),
			},
		},
	})
}

const testInternationalInstanceVar = `
variable "vpc_id" {
  default = "` + tcacctest.DefaultInternationalGrafanaVpcId + `"
}
variable "subnet_id" {
  default = "` + tcacctest.DefaultInternationalGrafanaSubnetId + `"
}
`
const testInternationalInstance_basic = testInternationalInstanceVar + `
resource "tencentcloud_monitor_tmp_instance" "basic" {
 instance_name 		= "demo-test"
 vpc_id 				= var.vpc_id
 subnet_id				= var.subnet_id
 data_retention_time	= 30
 zone 					= "ap-guangzhou-4"
 tags = {
   "createdBy" = "terraform"
 }
}`

const testInternationalInstance_update = testInternationalInstanceVar + `
resource "tencentcloud_monitor_tmp_instance" "update" {
 instance_name 		= "demo-test-update"
 vpc_id 				= var.vpc_id
 subnet_id				= var.subnet_id
 data_retention_time	= 30
 zone 					= "ap-guangzhou-4"
 tags = {
   "createdBy" = "terraform"
 }
}`
