Use this data source to query the list of SQL Server accounts.

Example Usage

Pull instance account list

```hcl
data "tencentcloud_sqlserver_accounts" "example" {
  instance_id = "mssql-3cdq7kx5"
}
```

Pull instance account list Filter by name

```hcl
data "tencentcloud_sqlserver_accounts" "example" {
  instance_id = "mssql-3cdq7kx5"
  name        = "myaccount"
}
```