resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos1-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket" "bar" {
  bucket = "scf-cos2-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket = "${tencentcloud_cos_bucket.foo.bucket}"
  key    = "/new_object_key.zip"
  source = "code.zip"
  acl    = "public-read"
}

resource "tencentcloud_cam_role" "foo" {
  name          = "ci-scf-role"
  document      = "${var.role_document}"
  description   = "ci-scf-role"
  console_login = true
}

resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-scf"
  description = "test1"
}

resource "tencentcloud_scf_function" "foo" {
  name        = "ci-test-function"
  description = "test"
  handler     = "main.do_it"
  runtime     = "Python3.6"
  namespace   = "${tencentcloud_scf_namespace.foo.id}"
  role        = "${tencentcloud_cam_role.foo.id}"

  cos_bucket_name   = "${tencentcloud_cos_bucket.foo.id}"
  cos_object_name   = "${tencentcloud_cos_bucket_object.myobject.key}"
  cos_bucket_region = "ap-guangzhou"

  triggers {
    name         = "ci-test-fn-api-gw"
    type         = "timer"
    trigger_desc = "*/5 * * * * * *"
  }

  triggers {
    name         = "${tencentcloud_cos_bucket.bar.id}"
    type         = "cos"
    trigger_desc = "{\"event\":\"cos:ObjectCreated:Put\",\"filter\":{\"Prefix\":\"\",\"Suffix\":\"\"}}"
  }

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_scf_functions" "foo" {
  name        = "${tencentcloud_scf_function.foo.name}"
  description = "${tencentcloud_scf_function.foo.description}"
  namespace   = "${tencentcloud_scf_function.foo.namespace}"
  tags        = "${tencentcloud_scf_function.foo.tags}"
}

data "tencentcloud_scf_namespaces" "foo" {
  namespace   = "${tencentcloud_scf_namespace.foo.id}"
  description = "${tencentcloud_scf_namespace.foo.description}"
}

data "tencentcloud_scf_logs" "foo" {
  function_name = "${tencentcloud_scf_function.foo.name}"
  namespace     = "${tencentcloud_scf_function.foo.namespace}"
  offset        = 0
  limit         = 100
  order         = "desc"
  order_by      = "duration"
  ret_code      = "UserCodeException"
  start_time    = "2017-05-16 20:00:00"
  end_time      = "2017-05-17 20:00:00"
}
