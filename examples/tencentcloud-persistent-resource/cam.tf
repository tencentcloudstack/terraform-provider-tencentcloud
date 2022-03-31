resource "tencentcloud_cam_role" "role" {
  name          = var.cam_role_basic
  document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/${local.uin}:uin/${local.uin}\"]}}]}"
  description   = "test"
  console_login = true
}

resource "tencentcloud_cam_policy" "policy" {
  name        = var.cam_policy_basic
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
  description = "test"
}