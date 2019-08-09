package ratelimit

//default cgi limit

const (
	DefaultLimit int64 = 20
)

func init() {

	//old  (filename . key)
	limitConfig["resource_tc_instance"] = 50
	limitConfig["resource_tc_instance.create"] = 10
	limitConfig["resource_tc_instance.update"] = 10
	limitConfig["resource_tc_instance.delete"] = 10


	//new(filename . action)
	limitConfig["service_tencentcloud_mysql"] = 50
	limitConfig["service_tencentcloud_mysql.CreateDBInstanceHour"] = 20
	limitConfig["service_tencentcloud_mysql.OfflineIsolatedInstances"] = 20
	limitConfig["service_tencentcloud_mysql.CreateBackup"] = 5
	limitConfig["service_tencentcloud_mysql.ModifyInstanceParam"] = 20

	limitConfig["service_tencentcloud_cos"] = 10


}
