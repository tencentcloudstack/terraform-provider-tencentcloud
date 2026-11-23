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

data "tencentcloud_security_groups" "example_2" {
  name = "tf-example"
}

data "tencentcloud_security_groups" "example_3" {
  project_id = 0
}

data "tencentcloud_security_groups" "example_4" {
  tags = {
    createBy = "Terraform"
  }
}
```