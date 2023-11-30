Provides a resource to create a mps ai_analysis_template

Example Usage

```hcl
resource "tencentcloud_mps_ai_analysis_template" "ai_analysis_template" {
  name = "terraform-test"

  classification_configure {
    switch = "OFF"
  }

  cover_configure {
    switch = "ON"
  }

  frame_tag_configure {
    switch = "ON"
  }

  tag_configure {
    switch = "ON"
  }
}

```

Import

mps ai_analysis_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_ai_analysis_template.ai_analysis_template ai_analysis_template_id
```