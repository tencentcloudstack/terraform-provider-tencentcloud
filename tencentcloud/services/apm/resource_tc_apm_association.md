Provides a resource to create a APM association

Example Usage

```hcl
resource "tencentcloud_apm_association" "example" {
  instance_id  = tencentcloud_apm_instance.example.id
  product_name = "Prometheus"
  status       = 1
  peer_id      = "prom-kx3eqdby"
}
```

Import

APM association can be imported using the id, e.g.

```
terraform import tencentcloud_apm_association.example apm-jPr5iQL77#Prometheus
```
