Use this resource to create tcr namespace.

Example Usage

Create a tcr namespace instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id		= tencentcloud_tcr_instance.example.id
  name          	= "example"
  is_public		 	= true
  is_auto_scan		= true
  is_prevent_vul	= true
  severity			= "medium"
  cve_whitelist_items	{
    cve_id = "cve-xxxxx"
  }
}
```

Import

tcr namespace can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_namespace.example tcr_instance_id#namespace_name
```