package tencentcloud

const MysqlInstanceCommonTestCase = `
resource "tencentcloud_mysql_instance" "default" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id = "vpc-fzdzrsir"
	subnet_id = "subnet-he8ldxx6"
	engine_version = "5.7"
	root_password = "test1234"
	availability_zone = "ap-guangzhou-4"
}
`
