Provides a resource to manage a Demo widget. <!-- TEMPLATE: one-line description, no heading. -->

<!-- TEMPLATE: HARD RULE — the generated md contains ONLY the standard
     structure: one-line description, the type-specific version note
     (action / ephemeral_resource only), Example Usage, and the optional
     Import section. NEVER add any custom NOTE, API-call restriction,
     business-logic remark or caveat of your own — even when you
     understand such limits from the API docs. All such extra content is
     added manually by humans afterwards; mention findings in the final
     summary instead of the document. -->

<!-- TEMPLATE: type-specific version note, placed between the description and
     Example Usage — REQUIRED for these two types (delete for others):
     - action:              ~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.
     - ephemeral_resource:  ~> **NOTE:** Ephemeral Resources are supported in HashiCorp Terraform version 1.10 and later.
     (golden source: tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.md) -->

Example Usage

```hcl
resource "tencentcloud_demo_widget" "example" {
  name    = "tf-test-widget"
  zone_id = "zone-12345"
  remark  = "created by terraform"
}
```

<!-- TEMPLATE: for managed resources supporting import, keep this section.
     The import id wording MUST match the actual resource_id shape from the
     JSON input — never write a generic "id" for every resource:
     - single-id resources: "... using the id, e.g." + the single value;
     - composite-id resources (attachment, or general with a joined id):
       "... using the composite id, e.g." + document the segment names in
       parent-first order (e.g. instance_id#disk_id) and show a realistic
       joined value in the example, e.g. ins-12345#disk-67890. -->

Import

<!-- TEMPLATE: single-id form (delete this block for composite ids): -->

Demo widget can be imported using the id, e.g.

```
terraform import tencentcloud_demo_widget.example widget-12345
```

<!-- TEMPLATE: composite-id form (use this block instead when the id is a
     parent-first join of segment attributes):

Demo disk attachment can be imported using the composite id instance_id#disk_id, e.g.

```
terraform import tencentcloud_cbs_disk_attachment.example ins-12345#disk-67890
```

-->
