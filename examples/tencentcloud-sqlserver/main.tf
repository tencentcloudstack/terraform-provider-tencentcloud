data "tencentcloud_sqlserver_zone_config" "foo" {
}

resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = "mssql-3cdq7kx5"
  name = "db_brickzzhang_update"
  charset = "Chinese_PRC_BIN"
  remark = "test-remark-update"
}