Provides an available EIP for the user.

The EIP data source fetch proper EIP from user's EIP pool.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_eips.

Example Usage

```hcl
data "tencentcloud_eip" "my_eip" {
  filter {
    name   = "address-status"
    values = ["UNBIND"]
  }
}
```