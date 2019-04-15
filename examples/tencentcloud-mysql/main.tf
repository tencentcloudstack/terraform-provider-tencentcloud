resource "tencentcloud_mysql_instance" "main" {
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id = "${var.default_vpc_id}"
	subnet_id = "${var.default_subnet_id}"
	engine_version = "5.7"
	root_password = "test1234"
	availability_zone = "${var.availability_zone}"
	internet_service = 1
	slave_sync_mode =1
	intranet_port =3360
	tags = {
		purpose ="for test"
	}
	parameters = {
		max_connections = "1000"
	}
}

data "tencentcloud_mysql_parameter_list" "mysql" {
	mysql_id = "${tencentcloud_mysql_instance.main.id}"
}

resource "tencentcloud_mysql_account" "mysql_account" {
	mysql_id = "${tencentcloud_mysql_instance.main.id}"
	name = "test"
	password = "test1234"
	description = "for test"
}

resource "tencentcloud_mysql_account_privilege" "mysql_account_privilege"{
	mysql_id = "${tencentcloud_mysql_instance.main.id}"
	account_name = "${tencentcloud_mysql_account.mysql_account.name}"
	privileges = ["SELECT", "INSERT", "UPDATE", "DELETE"]
	database_names=["test"]
}

resource "tencentcloud_mysql_backup_policy" "mysql_backup_policy" {
	mysql_id = "${tencentcloud_mysql_instance.main.id}"
	retention_period = 56
	backup_model = "physical"
	backup_time = "10:00-14:00"
}