Provides a resource to create a tcr webhook trigger

Example Usage

Create a tcr webhook trigger instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
	instance_id 	 = tencentcloud_tcr_instance.example.id
	name			 = "tf_example_ns_retention"
	is_public		 = true
	is_auto_scan	 = true
	is_prevent_vul = true
	severity		 = "medium"
	cve_whitelist_items	{
	  cve_id = "cve-xxxxx"
	}
  }

data "tencentcloud_tcr_namespaces" "example" {
	instance_id = tencentcloud_tcr_namespace.example.instance_id
  }

locals {
    ns_id = data.tencentcloud_tcr_namespaces.example.namespace_list.0.id
  }

resource "tencentcloud_tcr_webhook_trigger" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  namespace = tencentcloud_tcr_namespace.example.name
  trigger {
		name = "trigger-example"
		targets {
			address = "http://example.org/post"
			headers {
				key = "X-Custom-Header"
				values = ["a"]
			}
		}
		event_types = ["pushImage"]
		condition = ".*"
		enabled = true
		description = "example for trigger description"
		namespace_id = local.ns_id

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr webhook_trigger can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_webhook_trigger.example webhook_trigger_id
```