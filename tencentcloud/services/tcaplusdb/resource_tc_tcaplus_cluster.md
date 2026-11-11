Provides a resource to create a TcaplusDB cluster.

~> **NOTE:** TcaplusDB now only supports the following regions: `ap-shanghai`, `ap-hongkong`, `na-siliconvalley`, `ap-singapore`, `ap-seoul`, `ap-tokyo`, `eu-frankfurt`, `and na-ashburn`.

Example Usage

Create a tcaplus cluster instance with cluster_type is 1

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "MIX"
  cluster_name             = "tf_example"
  vpc_id                   = "vpc-jll1dzwr"
  subnet_id                = "subnet-ef14ogeu"
  password                 = "Password@2026"
  old_password_expire_last = 3600
  cluster_type             = 1
  resource_tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

Create a tcaplus cluster instance with cluster_type is 2

```hcl
resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "MIX"
  cluster_name             = "tf_example"
  vpc_id                   = "vpc-qtzga3pm"
  subnet_id                = "subnet-c063n9el"
  password                 = "Password@2026"
  old_password_expire_last = 3600
  cluster_type             = 2
  server_list {
    machine_type = "T1"
    machine_num  = 4
  }

  proxy_list {
    machine_type = "T1"
    machine_num  = 2
  }

  resource_tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

Import

TcaplusDB cluster can be imported using the id, e.g.

```
terraform import tencentcloud_tcaplus_cluster.example 35402666774
```
