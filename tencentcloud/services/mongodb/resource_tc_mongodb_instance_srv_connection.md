Provides a resource to manage MongoDB instance SRV connection URL configuration.

Example Usage

Enable SRV connection with default domain

```hcl
resource "tencentcloud_mongodb_instance_srv_connection" "example" {
  instance_id = "cmgo-p8vnipr5"
}

output "domain" {
  value = tencentcloud_mongodb_instance_srv_connection.example.domain
}
```

Enable SRV connection with custom domain

```hcl
resource "tencentcloud_mongodb_instance_srv_connection" "example" {
  instance_id = "cmgo-p8vnipr5"
  domain      = "exampleDomain"
}
```

Import

MongoDB instance SRV connection can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_instance_srv_connection.example cmgo-p8vnipr5
```
