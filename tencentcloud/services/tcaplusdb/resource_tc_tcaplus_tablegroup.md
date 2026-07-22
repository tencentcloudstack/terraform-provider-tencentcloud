Use this resource to create TcaplusDB table group.

Example Usage

Create a tcaplusdb table group

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = "vpc-i5yyodl9"
  subnet_id                = "subnet-hhi88a58"
  password                 = "Password@2026"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "example" {
  cluster_id      = tencentcloud_tcaplus_cluster.example.id
  tablegroup_name = "tf_example_group_name"
  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

Create a tcaplusdb table group with user-specified table group id

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = "vpc-i5yyodl9"
  subnet_id                = "subnet-hhi88a58"
  password                 = "Password@2026"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "example" {
  cluster_id      = tencentcloud_tcaplus_cluster.example.id
  tablegroup_name = "tf_example_group_name"
  table_group_id  = "109"
  resource_tags {
    tag_key   = "CreatedBy"
    tag_value = "Terraform"
  }
}
```

Import

TcaplusDB table group can be imported using the clusterId:tableGroupId, e.g.

```
terraform import tencentcloud_tcaplus_tablegroup.example 5516511420:52
```
