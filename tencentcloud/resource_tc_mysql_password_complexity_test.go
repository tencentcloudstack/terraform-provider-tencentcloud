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

const testAccMysqlPasswordComplexityVar = `
variable "instance_id" {
  default = "` + defaultDbBrainInstanceId + `"
}
`

const testAccMysqlPasswordComplexity = testAccMysqlPasswordComplexityVar + `

resource "tencentcloud_mysql_password_complexity" "password_complexity" {
	instance_id = var.instance_id
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
