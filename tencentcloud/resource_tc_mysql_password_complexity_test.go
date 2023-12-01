package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlPasswordComplexityResource_basic -v
func TestAccTencentCloudMysqlPasswordComplexityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlPasswordComplexity,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_password_complexity.password_complexity", "id"),
				),
			},
		},
	})
}

const testAccMysqlPasswordComplexity = `

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

resource "tencentcloud_mysql_password_complexity" "password_complexity" {
	instance_id = tencentcloud_mysql_instance.this.id
	param_list {
	  name = "validate_password_length"
	  current_value = "8"
	}
	param_list {
	  name = "validate_password_mixed_case_count"
	  current_value = "2"
	}
	param_list {
	  name = "validate_password_number_count"
	  current_value = "2"
	}
	param_list {
	  name = "validate_password_special_char_count"
	  current_value = "2"
	}
}

`
