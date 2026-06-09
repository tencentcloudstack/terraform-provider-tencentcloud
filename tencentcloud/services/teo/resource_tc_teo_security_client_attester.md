Provides a resource to create a TEO security client attester.

Example Usage

```hcl
resource "tencentcloud_teo_security_client_attester" "example" {
  zone_id = "zone-3fkff38fyw8s"

  client_attesters {
    name              = "tf-example"
    attester_source   = "TC-RCE"
    attester_duration = "300s"

    tc_rce_option {
      channel = "12399223"
      region  = "ap-beijing"
    }
  }
}
```

Import

TEO security client attester can be imported using the zoneId#clientAttesterId, e.g.

```
terraform import tencentcloud_teo_security_client_attester.example zone-3fkff38fyw8s#attest-0000361666
```
