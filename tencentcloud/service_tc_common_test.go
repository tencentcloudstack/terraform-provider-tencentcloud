package tencentcloud

const MysqlInstanceCommonTestCase = `
resource "tencentcloud_mysql_instance" "default" {
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	engine_version = "5.7"
	root_password = "test1234"
	availability_zone = "ap-guangzhou-3"
}
`
