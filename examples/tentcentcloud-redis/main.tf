resource "tencentcloud_redis_instance" "redis_instance_test"{
	availability_zone="${var.availability_zone}"
	type="master_slave_redis"
	password="test12345789"
	mem_size=8192
	name="terrform_test"
	port=6379
}

resource "tencentcloud_redis_backup_config" "redis_backup_config" {
	redis_id="${tencentcloud_redis_instance.redis_instance_test.id}"
	backup_time="01:00-02:00"
	backup_period=["Saturday", "Sunday"]
}

