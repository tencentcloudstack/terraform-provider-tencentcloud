data "tencentcloud_audit_cos_regions" "example" {
}

data "tencentcloud_audits" "example" {
}

data "tencentcloud_audit_key_alias" "example" {
  region = "ap-hongkong"
}

resource "tencentcloud_audit" "example_kms" {
  name                 = "example_kms"
  cos_bucket           = "test"
  cos_region           = "ap-hongkong"
  enable_kms_encry     = true
  log_file_prefix      = "exampleprefix"
  key_id               = data.tencentcloud_audit_key_alias.example.audit_key_alias_list.0.key_id
  audit_switch         = true
  read_write_attribute = 3
}

data "tencentcloud_audits" "name_example" {
  name = "example_kms"
}
