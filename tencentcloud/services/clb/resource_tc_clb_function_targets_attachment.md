Provides a resource to create a clb function_targets_attachment

Example Usage

```hcl
resource "tencentcloud_clb_function_targets_attachment" "function_targets" {
  domain           = "xxx.com"
  listener_id      = "lbl-nonkgvc2"
  load_balancer_id = "lb-5dnrkgry"
  url              = "/"

  function_targets {
    weight = 10

    function {
      function_name           = "keep-tf-test-1675954233"
      function_namespace      = "default"
      function_qualifier      = "$LATEST"
      function_qualifier_type = "VERSION"
    }
  }
}
```

Import

clb function_targets_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_clb_function_targets_attachment.function_targets loadBalancerId#listenerId#locationId or loadBalancerId#listenerId#domain#rule
```