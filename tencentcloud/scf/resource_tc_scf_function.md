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

Using CFS config
```
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cfs_config {
	user_id	= "10000"
	user_group_id	= "10000"
	cfs_id	= "cfs-xxxxxxxx"
	mount_ins_id	= "cfs-xxxxxxxx"
	local_mount_dir	= "/mnt"
	remote_mount_dir	= "/"
  }
}
```

Import

SCF function can be imported, e.g.

-> **NOTE:** function id is `<function namespace>+<function name>`

```
$ terraform import tencentcloud_scf_function.test default+test
```