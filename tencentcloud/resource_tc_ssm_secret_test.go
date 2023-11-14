package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSsmSecretResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmSecret,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssm_secret.secret", "id")),
			},
			{
				ResourceName:      "tencentcloud_ssm_secret.secret",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSsmSecret = `

resource "tencentcloud_ssm_secret" "secret" {
  secret_name = "test_name"
  user_name_prefix = "test_prefix"
  enable_rotation = False
  rotation_begin_time = "2006-01-02 15:04:05"
  rotation_frequency = 1
  instance_i_d = "cdb-xxxxxxxx"
  description = ""
  kms_key_id = ""
  product_name = "Mysql"
  domains = 
  privileges_list {
		privilege_name = "GlobalPrivileges"
		privileges = 
		database = ""
		table_name = ""
		column_name = ""

  }
  project_id = 
  s_s_h_key_name = ""
  version_id = "v1.0"
  secret_type = 0
  secret_binary = ""
  secret_string = "test"
  additional_config = "{}"
  tags = {
    "createdBy" = "terraform"
  }
}

`
