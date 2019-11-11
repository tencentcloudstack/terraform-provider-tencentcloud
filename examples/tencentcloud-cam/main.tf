resource "tencentcloud_cam_group" "example" {
  name   = "example"
  remark = "example"
}

resource "tencentcloud_cam_user" "example" {
  name                = "example"
  remark              = "example"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "${var.password}"
  phone_num           = "${var.phone_num}"
  country_code        = "${var.country_code}"
  email               = "${var.email}"
}

resource "tencentcloud_cam_policy" "example" {
  name     = "example"
  document = "${var.policy_document}"
}

resource "tencentcloud_cam_role" "example" {
  name          = "example"
  document      = "${var.role_document}"
  description   = "test"
  console_login = true
}

resource "tencentcloud_cam_group_membership" "example" {
  group_id = "${tencentcloud_cam_group.example.id}"
  user_ids = ["${tencentcloud_cam_user.example.id}"]
}

resource "tencentcloud_cam_role_policy_attachment" "example" {
  role_id   = "${tencentcloud_cam_role.example.id}"
  policy_id = "${tencentcloud_cam_policy.example.id}"
}

resource "tencentcloud_cam_user_policy_attachment" "example" {
  user_id   = "${tencentcloud_cam_user.example.id}"
  policy_id = "${tencentcloud_cam_policy.example.id}"
}

resource "tencentcloud_cam_group_policy_attachment" "example" {
  group_id  = "${tencentcloud_cam_group.example.id}"
  policy_id = "${tencentcloud_cam_policy.example.id}"
}

resource "tencentcloud_cam_saml_provider" "example" {
  name        = "example"
  meta_data   = "${var.meta_data}"
  description = "test"
}

data "tencentcloud_cam_users" "users" {
  name = "${tencentcloud_cam_user.example.id}"
}

data "tencentcloud_cam_roles" "roles" {
  role_id = "${tencentcloud_cam_role.example.id}"
}

data "tencentcloud_cam_policies" "policies" {
  policy_id = "${tencentcloud_cam_policy.example.id}"
}

data "tencentcloud_cam_groups" "groups" {
  group_id = "${tencentcloud_cam_group.example.id}"
}

data "tencentcloud_cam_group_memberships" "memberships" {
  group_id = "${tencentcloud_cam_group_membership.example.id}"
}

data "tencentcloud_cam_user_policy_attachments" "user_policy_attachments" {
  user_id = "${tencentcloud_cam_user_policy_attachment.example.user_id}"
}

data "tencentcloud_cam_role_policy_attachments" "role_policy_attachments" {
  role_id = "${tencentcloud_cam_role_policy_attachment.example.role_id}"
}

data "tencentcloud_cam_group_policy_attachments" "group_policy_attachments" {
  group_id = "${tencentcloud_cam_group_policy_attachment.example.group_id}"
}

data "tencentcloud_cam_saml_providers" "saml_providers" {
  name = "${tencentcloud_cam_saml_provider.example.id}"
}