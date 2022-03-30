locals {
  app_id = data.tencentcloud_user_info.info.app_id
  uin = data.tencentcloud_user_info.info.uin
  owner_uin = data.tencentcloud_user_info.info.owner_uin
}

data "tencentcloud_user_info" "info" {}