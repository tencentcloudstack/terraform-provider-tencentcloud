Provides a resource to create a tem gateway

Example Usage

```hcl
resource "tencentcloud_tem_gateway" "gateway" {
  ingress {
    ingress_name = "demo"
    environment_id = "en-853mggjm"
    address_ip_version = "IPV4"
    rewrite_type = "NONE"
    mixed = false
    rules {
      host = "test.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "demo"
            service_port = 80
          }
        }
      }
    }
    rules {
      host = "hello.com"
      protocol = "http"
      http {
        paths {
          path = "/"
          backend {
            service_name = "hello"
            service_port = 36000
          }
        }
      }
    }
  }
}

```
Import

tem gateway can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_gateway.gateway environmentId#gatewayName
```