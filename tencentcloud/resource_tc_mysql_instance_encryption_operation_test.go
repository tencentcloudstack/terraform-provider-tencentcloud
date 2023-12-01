package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlInstanceEncryptionOperationResource_basic -v
func TestAccTencentCloudMysqlInstanceEncryptionOperationResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlInstanceEncryptionOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance_encryption_operation.instance_encryption_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance_encryption_operation.instance_encryption_operation", "key_id", "KMS-CDB"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance_encryption_operation.instance_encryption_operation", "key_region", "ap-guangzhou"),
				),
			},
		},
	})
}

const testAccMysqlInstanceEncryptionOperationVar = `
resource "tencentcloud_mysql_instance" "mysql8" {
	internet_service = 1
	engine_version   = "8.0"
	charge_type = "POSTPAID"
	root_password     = "password123"
	slave_deploy_mode = 0
	first_slave_zone  = "ap-guangzhou-4"
	second_slave_zone = "ap-guangzhou-4"
	slave_sync_mode   = 1
	availability_zone = "ap-guangzhou-4"
	project_id        = 0
	instance_name     = "myTestMysql"
	mem_size          = 1000
	volume_size       = 25
	intranet_port     = 3306
	security_groups   = ["sg-ngx2bo7j"]
  
	tags = {
	  createdBy = "terraform"
	}
  
	parameters = {
	  character_set_server = "utf8"
	  lower_case_table_names = 0
	  max_connections = "1000"
	}
}`

const testAccMysqlInstanceEncryptionOperation = testAccMysqlInstanceEncryptionOperationVar + `

resource "tencentcloud_mysql_instance_encryption_operation" "instance_encryption_operation" {
  instance_id = tencentcloud_mysql_instance.mysql8.id
  key_id = "KMS-CDB"
  key_region = "ap-guangzhou"
}

`
