Provides a resource to create a privatedns account

Example Usage

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account {
    uin = "100022770160"
  }
}
```

Or

```hcl
resource "tencentcloud_private_dns_account" "example" {
  account {
    uin      = "100022770160"
    account  = "example@tencent.com"
    nickname = "tf-example"
  }
}
```

Import

privatedns account can be imported using the id, e.g.

```
terraform import tencentcloud_private_dns_account.example 100022770160
```
