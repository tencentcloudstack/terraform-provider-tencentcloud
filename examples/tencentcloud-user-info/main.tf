data "tencentcloud_user_info" "info" {}

output "appid" {
  value = data.tencentcloud_user_info.info.app_id
}