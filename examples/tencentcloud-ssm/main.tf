provider "tencentcloud" {
  region = "ap-guangzhou"
}

resource "tencentcloud_ssm_secret" "foo" {
  secret_name = "test"
  description = "test secret"
  recovery_window_in_days = 0
  is_enabled = true

  init_secret {
    version_id = "v1"
    secret_string = "123456"
  }

  tags = {
    test-tag = "test"
  }
}

resource "tencentcloud_ssm_secret_version" "v2" {
  secret_name = tencentcloud_ssm_secret.foo.secret_name
  version_id = "v2"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}

data "tencentcloud_ssm_secrets" "secret_list" {
  secret_name = tencentcloud_ssm_secret.foo.secret_name
  order_type = 1
  state = 1
}

data "tencentcloud_ssm_secret_versions" "secret_version_list" {
  secret_name = tencentcloud_ssm_secret_version.v2.secret_name
  version_id = tencentcloud_ssm_secret_version.v2.version_id
}