Use this data source to query eip instances.

Example Usage

Query all eip instances

```hcl
data "tencentcloud_eips" "example" {}
```

Query eip instances by eip ID

```hcl
data "tencentcloud_eips" "example" {
  eip_id = "eip-ry9h95hg"
}
```

Query eip instances by eip name

```hcl
data "tencentcloud_eips" "example" {
  eip_name = "tf-example"
}
```

Query eip instances by public ip

```hcl
data "tencentcloud_eips" "example" {
  public_ip = "1.12.62.3"
}
```

Query eip instances by tags

```hcl
data "tencentcloud_eips" "example" {
  tags = {
    "test" = "test"
  }
}
```
