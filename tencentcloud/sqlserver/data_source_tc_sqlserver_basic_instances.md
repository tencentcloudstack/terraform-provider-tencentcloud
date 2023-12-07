Use this data source to query SQL Server basic instances

Example Usage

Filter instance by Id

```hcl
data "tencentcloud_sqlserver_basic_instances" "example_id" {
  id = "mssql-3l3fgqn7"
}
```

Filter instance by project Id

```hcl
data "tencentcloud_sqlserver_basic_instances" "example_project" {
  project_id = 0
}
```

Filter instance by VPC/Subnet

```hcl
data "tencentcloud_sqlserver_basic_instances" "example_vpc" {
  vpc_id    = "vpc-409mvdvv"
  subnet_id = "subnet-nf9n81ps"
}
```