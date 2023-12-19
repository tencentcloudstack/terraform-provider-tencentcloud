package cdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlDbImportJobOperationResource_basic -v
func TestAccTencentCloudMysqlDbImportJobOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDbImportJobOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_db_import_job_operation.db_import_job_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_db_import_job_operation.db_import_job_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlDbImportJobOperation = `


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

resource "tencentcloud_cos_bucket_object" "object_content" {
	bucket       = "terraform-ci-1308919341"
	key          = "/mysql/mysql.sql"
	content      = "SELECT NOW(),SYSDATE();"
	content_type = "binary/octet-stream"
	acl 		 = "public-read"
  }

resource "tencentcloud_mysql_db_import_job_operation" "db_import_job" {
	instance_id = tencentcloud_mysql_instance.this.id
	user = "root"
	file_name = "mysql.sql"
	password = "password123"
	# db_name = "t_test"
	cos_url = "https://terraform-ci-1308919341.cos.ap-guangzhou.myqcloud.com${tencentcloud_cos_bucket_object.object_content.key}"
}

`
