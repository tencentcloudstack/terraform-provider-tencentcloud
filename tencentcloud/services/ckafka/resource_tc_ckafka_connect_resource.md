Provides a resource to create a ckafka connect_resource

Example Usage

```hcl
resource "tencentcloud_ckafka_connect_resource" "connect_resource" {
  resource_name = "terraform-test"
  type          = "MYSQL"
  description   = "for terraform test"

  mysql_connect_param {
    port        = 3306
    user_name   = "root"
    password    = "xxxxxxxxx"
    resource    = "cdb-fitq5t9h"
    service_vip = "172.16.80.59"
    uniq_vpc_id = "vpc-4owdpnwr"
    self_built  = false
  }
}

```

Import

ckafka connect_resource can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_connect_resource.connect_resource connect_resource_id
```