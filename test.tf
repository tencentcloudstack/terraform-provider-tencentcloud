provider "tencentcloud" {
  secret_id  = "AKIDwVqwp873h5Zt22AJ5sVv2va1CE5ipdx4"
  secret_key = "KQJ1VomW1b3GGYcBcpAMJ1823zeBb1DA"
  region     = "ap-guangzhou"
}

resource "tencentcloud_mysql_account_privilege" "default" {
  mysql_id = "cdb-q6muaqct"
  account_name= "tf_account"
  privileges = ["SELECT"]
  database_names = ["myTestMysql"]
}

resource "tencentcloud_mysql_account" "default" {
  mysql_id = "cdb-q6muaqct"
  name = "tf_account"
  password = "Wfxinjihao1998"
  description = "My test account"
}

