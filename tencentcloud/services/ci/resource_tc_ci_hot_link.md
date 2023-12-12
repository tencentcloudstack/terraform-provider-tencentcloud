Provides a resource to create a ci hot_link

Example Usage

```hcl
resource "tencentcloud_ci_hot_link" "hot_link" {
	bucket = "terraform-ci-xxxxxx"
	url = ["10.0.0.1", "10.0.0.2"]
	type = "white"
}
```

Import

ci hot_link can be imported using the bucket, e.g.

```
terraform import tencentcloud_ci_hot_link.hot_link terraform-ci-xxxxxx
```