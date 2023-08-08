locals {
  uin = data.tencentcloud_user_info.info.uin
}

data "tencentcloud_user_info" "info" {}