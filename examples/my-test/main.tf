terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

provider "tencentcloud" {
  region = "ap-guangzhou"
}

resource "tencentcloud_sqlserver_general_backup" "test_sqlserver_backup" {
  instance_id = "mssql-qelbzgwf"
  backup_name = "update_backup"
}