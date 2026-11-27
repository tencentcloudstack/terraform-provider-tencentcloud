Provides a resource to manage a Demo widget. <!-- TEMPLATE: one-line description, no heading. Human-readable resource name ONLY: split the resource name on underscores into separate words (api_gateway_service_release -> "API gateway service release"); NEVER keep underscores in prose — they break markdown rendering and read like code. Product acronyms (API, SSL, VPC, ...) stay uppercase. Golden source: tencentcloud/services/apigateway/resource_tc_api_gateway_service_release.md. -->

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

<!-- TEMPLATE: a minimal, realistic HCL example with obviously-fake values
     (tf-test, example) — never real credentials. When the resource/action
     declares a `timeouts` block (async scenarios), SHOW it in the example
     with the default values (golden source: mongodb action md), e.g.

       timeouts {
         create = "15m"
       }

     HCL syntax is decided by the schema type of EACH field (conventions
     section 6), at every nesting depth:
       SingleNestedAttribute / ObjectAttribute -> attr = { k = v }
       List/SetNestedAttribute                 -> attr = [{ k = v }]
       MapNestedAttribute                      -> attr = { key = { k = v } }
       MapAttribute                            -> attr = { k = "v" }
       List/SetAttribute                       -> attr = ["a", "b"]
       Single/List/SetNestedBlock, `Blocks:`   -> name { ... }  (incl. timeouts)
     Rule of thumb: a type named *Attribute ALWAYS has `=`; a type named *Block
     NEVER has `=`. Real regression: encryption_configuration is a
     SingleNestedAttribute, so `encryption_configuration = { sse_type = "AES256" }`,
     never `encryption_configuration { ... }`. Enforced by check_md_hcl.py.

     TEMPLATE: HARD RULE — align the generated file line-by-line with THIS
     template (conventions section 6): a blank line after EVERY heading
     ("Example Usage", "Import") and before each code fence; the Example
     Usage fence is ```hcl, the terraform import fence is a BARE ``` (CLI
     command, not HCL). -->

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
       joined value in the example, e.g. ins-12345#disk-67890.
     The SUBJECT of the sentence uses the human-readable resource name
     (underscores split into words, same rule as the first line):
       "API gateway service release can be imported using the id, e.g."
     NEVER "API gateway_service_release can be imported ...".
     Underscores stay ONLY in code references: the terraform import command
     (tencentcloud_demo_widget.example) and the composite id segment names
     (instance_id#disk_id) — those are literal identifiers, not prose. -->

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
