Provide a resource to create a SCF function.

Example Usage

```hcl
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "scf-code-1234567890"
  cos_object_name   = "code.zip"
  cos_bucket_region = "ap-guangzhou"
}
```

Using Zip file

```hcl
resource "tencentcloud_scf_function" "foo" {
  name              = "ci-test-function"
  handler           = "first.do_it_first"
  runtime           = "Python3.6"
  enable_public_net = true
  dns_cache         = true
  intranet_config {
    ip_fixed = "ENABLE"
  }
  vpc_id    = "vpc-391sv4w3"
  subnet_id = "subnet-ljyn7h30"

  zip_file = "/scf/first.zip"

  tags = {
    "env" = "test"
  }
}

Using CFS config
```

Import

SCF function can be imported, e.g.

-> **NOTE:** function id is `<function namespace>+<function name>`

```
$ terraform import tencentcloud_scf_function.test default+test
```