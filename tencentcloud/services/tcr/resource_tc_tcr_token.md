Use this resource to create tcr long term token.

Example Usage

Create a token for tcr instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_token" "example" {
  instance_id		= tencentcloud_tcr_instance.example.id
  description		= "example for the tcr token"
}
```

Import

tcr token can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_token.example instance_id#token_id
```