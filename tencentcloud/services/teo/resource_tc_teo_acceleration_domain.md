Provides a resource to create a TEO acceleration domain

~> **NOTE:** Before modifying resource content, you need to ensure that the `status` is `online`.

~> **NOTE:** Only `origin_type` is `IP_DOMAIN` can set `host_header`; And when `origin_type` is changed to `IP_DOMAIN`, `host_header` needs to be set to a legal value, such as a domain name string(like `domain_name`).

~> **NOTE:** If you use a third-party storage bucket configured for back-to-source, you need to ignore changes to `SecretAccessKey`.

~> **NOTE:** Before creating an accelerated domain, please verify the domain ownership of the site.

Example Usage

```hcl
resource "tencentcloud_teo_acceleration_domain" "example" {
  zone_id     = "zone-39quuimqg8r6"
  domain_name = "www.demo.com"

  origin_info {
    origin      = "150.109.8.1"
    origin_type = "IP_DOMAIN"
  }

  status            = "online"
  origin_protocol   = "FOLLOW"
  http_origin_port  = 80
  https_origin_port = 443
  ipv6_status       = "follow"
}
```

Back-to-source configuration using a third-party storage bucket. 
SecretAccessKey is sensitive data and can no longer be queried in plain text, so changes to SecretAccessKey need to be ignored.

```hcl
resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
    domain_name       = "cos.demo.cn"
    http_origin_port  = 80
    https_origin_port = 443
    ipv6_status       = "follow"
    origin_protocol   = "FOLLOW"
    status            = "online"
    zone_id           = "zone-39quuimqg8r6"

    origin_info {
        backup_origin    = null
        origin           = "example.s3.ap-northeast.amazonaws.com"
        origin_type      = "AWS_S3"
        private_access   = "on"
        vod_bucket_id    = null
        vod_origin_scope = null

        private_parameters {
            name  = "AccessKeyId"
            value = "aaaaaaa"
        }
        private_parameters {
            name  = "SecretAccessKey"
            value = "bbbbbbb"
        }
        private_parameters {
            name  = "SignatureVersion"
            value = "v4"
        }
        private_parameters {
            name  = "Region"
            value = "us-east1"
        }
    }
    lifecycle {
      ignore_changes = [ 
        origin_info[0].private_parameters[1].value,
       ]
    }
}
```

Import

TEO acceleration domain can be imported using the id, e.g.

```
terraform import tencentcloud_teo_acceleration_domain.example zone-39quuimqg8r6#www.demo.com
```
