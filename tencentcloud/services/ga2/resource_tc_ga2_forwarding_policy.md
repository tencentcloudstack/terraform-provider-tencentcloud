Provides a resource to create a GA2 forwarding policy

Example Usage

```hcl
resource "tencentcloud_ga2_forwarding_policy" "example" {
  global_accelerator_id = "ga-fhhs8w84"
  listener_id           = "lsr-dyy8jhzp"
  host                  = "example.com"
}
```

Import

GA2 forwarding policy can be imported using the composite id `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>`, e.g.

```
terraform import tencentcloud_ga2_forwarding_policy.example ga-fhhs8w84#lsr-dyy8jhzp#dm-rjssxr8k
```