Manage Guetzli compression functionality

Example Usage

```hcl

resource "tencentcloud_ci_guetzli" "foo" {
	bucket = "examplebucket-1250000000"
	status = "on"
}

```

Import

Resource guetzli can be imported using the id, e.g.

```
$ terraform import tencentcloud_ci_guetzli.example examplebucket-1250000000
```