Provide a resource to configure kubernetes cluster app addons.

Example Usage

Install tcr addon

```hcl
resource "tencentcloud_kubernetes_addon" "example" {
  cluster_id = "cls-k2o1ws9g"
  addon_name = "tcr"
  raw_values = jsonencode({
    global = {
      imagePullSecretsCrs = [
        {
          name            = "tcr-h3ff76s9"
          namespaces      = "*"
          serviceAccounts = "*"
          type            = "docker"
          dockerUsername  = "100038911322"
          dockerPassword  = "eyJhbGciOiJSUzI1NiIsImtpZCI6************"
          dockerServer    = "testcd.tencentcloudcr.com"
        }
      ]
    }
  })
}
```

Import

kubernetes cluster app addons can be imported using the id, e.g.
```
$ terraform import tencentcloud_kubernetes_addon.example cls-k2o1ws9g#tcr
```