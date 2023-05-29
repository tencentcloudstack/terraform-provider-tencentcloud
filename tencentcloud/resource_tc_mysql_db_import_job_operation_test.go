package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlDbImportJobOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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

resource "tencentcloud_mysql_db_import_job_operation" "db_import_job" {
	instance_id = tencentcloud_mysql_instance.this.id
	user = "root"
	file_name = "mysql.sql"
	password = "password123"
	# db_name = "t_test"
	cos_url = "https://terraform-ci-1308919341.cos.ap-guangzhou.myqcloud.com/mysql/mysql.sql?q-sign-algorithm=sha1&q-ak=AKIDlchzcM5ppPlbSV7yhstd8narnfMtVzZRQsayiPEzHxN5pb5UDeOL2yrNqwn2Yztr&q-sign-time=1684998880;1685002480&q-key-time=1684998880;1685002480&q-header-list=host&q-url-param-list=&q-signature=b7e6165f971009fc3628eca370b9c784269cb948&x-cos-security-token=bF2h2255ZhdpmBoYFCFme0eH2h2wDJPa75147a141a51a01e7e7c49b25de5baa3r9JfJvPSQ-zrBNd5wWgKPO8slY1_PK34fw-6oxLB6EnUez5quPhZ2bPGjZ9Wmktp3st44c-0zipO4MFoQw5ZQYLuezMrpfgejRtgzcMA6xg9vAjfhYnDEmGLbAirarpWmjNotia7Xgo0sr6hVjz7pcOXhs327895IjQDQrnYMw4CcZYaekjm7sEv51XqK5V5"
  }

`
