Use this data source to query detailed information of security groups.

Example Usage

Query all security groups

```hcl
data "tencentcloud_security_groups" "example" {}
```

Query security groups by filter

```hcl
data "tencentcloud_security_groups" "example" {
  security_group_id = "sg-e699atb7"
}

data "tencentcloud_security_groups" "example" {
  name = "tf-example"
}

data "tencentcloud_security_groups" "example" {
  project_id = 0
}

data "tencentcloud_security_groups" "example" {
  tags = {
    createBy = "Terraform"
  }
}
```