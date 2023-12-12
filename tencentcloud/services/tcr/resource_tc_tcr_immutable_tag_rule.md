Provides a resource to create a tcr immutable tag rule.

Example Usage

Create a immutable tag rule with specified tags and exclude specified repositories

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "premium"
  delete_bucket = true
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_example_ns"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_immutable_tag_rule" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  rule {
    repository_pattern    = "deprecated_repo" # specify exclude repo
    tag_pattern           = "**"              # all tags
    repository_decoration = "repoExcludes"
    tag_decoration        = "matches"
    disabled              = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

With specified repositories and exclude specified version tag

```hcl
resource "tencentcloud_tcr_immutable_tag_rule" "example" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  rule {
    repository_pattern    = "**" # all repo
    tag_pattern           = "v1" # exlude v1 tags
    repository_decoration = "repoMatches"
    tag_decoration        = "excludes"
    disabled              = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Disabled the specified rule

```hcl
resource "tencentcloud_tcr_immutable_tag_rule" "example_rule_A" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  rule {
    repository_pattern    = "deprecated_repo" # specify exclude repo
    tag_pattern           = "**"              # all tags
    repository_decoration = "repoExcludes"
    tag_decoration        = "matches"
    disabled              = false
  }
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_immutable_tag_rule" "example_rule_B" {
  registry_id    = tencentcloud_tcr_instance.example.id
  namespace_name = tencentcloud_tcr_namespace.example.name
  rule {
    repository_pattern    = "**" # all repo
    tag_pattern           = "v1" # exlude v1 tags
    repository_decoration = "repoMatches"
    tag_decoration        = "excludes"
    disabled              = true # disable it
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tcr immutable_tag_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tcr_immutable_tag_rule.immutable_tag_rule immutable_tag_rule_id
```