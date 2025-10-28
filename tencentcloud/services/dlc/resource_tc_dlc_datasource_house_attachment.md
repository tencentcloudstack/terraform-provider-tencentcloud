Provides a resource to create a DLC datasource house attachment

Example Usage

```hcl
resource "tencentcloud_dlc_datasource_house_attachment" "example" {
  datasource_connection_name = "tf-example"
  datasource_connection_type = "Mysql"
  datasource_connection_config {
    mysql {
      location {
        vpc_id            = "vpc-khkyabcd"
        vpc_cidr_block    = "192.168.0.0/16"
        subnet_id         = "subnet-o7n9eg12"
        subnet_cidr_block = "192.168.0.0/24"
      }
    }
  }
  
  data_engine_names       = ["engine_demo"]
  network_connection_type = 4
  network_connection_desc = "remark."
}
```
