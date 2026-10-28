Provide a resource to create a SnatIp of CLB instance.

~> **NOTE:** Target CLB instance must enable `snat_pro` before creating snat ips.

~> **NOTE:** Dynamic allocate IP doesn't support for now.

Example Usage

```hcl
resource "tencentcloud_clb_snat_ip" "example" {
  clb_id = "lb-jnx618r2"
  ips {
    subnet_id = "subnet-hhi88a58"
    ip        = "10.0.30.10"
  }

  ips {
    subnet_id = "subnet-d4umunpy"
  }
}
```

Import

Clb instance snat ip can be imported by clb instance id, e.g.

```
terraform import tencentcloud_clb_snat_ip.example lb-jnx618r2
```