package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlIsolateInstanceResource_basic -v
func TestAccTencentCloudMysqlIsolateInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlIsolateInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_isolate_instance.isolate_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_isolate_instance.isolate_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlIsolateInstance = `

resource "tencentcloud_mysql_instance" "this" {
	instance_name  = "test-nv"
	vpc_id         = "vpc-4owdpnwr"
	subnet_id      = "subnet-ahv6swf2"
	engine_version = "5.7"
	root_password  = "password123"
	availability_zone = "ap-guangzhou-3"
	mem_size       = 1000
	volume_size    = 25
	cpu            = 1
	intranet_port  = 3306
	security_groups   = ["sg-ngx2bo7j"]
  
	tags = {
	  createdBy = "terraform"
	}
  
	parameters = {
	  character_set_server = "gbk"
	  lower_case_table_names = "0"
	  max_connections      = "1000"
	}
}

resource "tencentcloud_mysql_isolate_instance" "isolate_instance" {
  instance_id = tencentcloud_mysql_instance.this.id
}

`
