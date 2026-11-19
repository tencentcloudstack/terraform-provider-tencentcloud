Provides a resource to manage environment variables of a Makers (TEO makers) project.

~> **NOTE** If you need to adjust the order of the `env_vars` entries, you must delete and recreate them.

Example Usage

```hcl
resource "tencentcloud_makers_project_envs" "example" {
  zone_id    = "zone-3hzk0ndvnnyp"
  project_id = "makers-jw7jha6dktob"

  env_vars = [
    # Variables specific to the production environment
    {
      key     = "API_KEY"
      value   = "secret-value"
      comment = "production api key"
      env     = ["Production"]
    },

    # Variables specific to the preview environment (sets of environment variables with the same key cannot overlap)
    {
      key     = "API_KEY"
      value   = "preview-key"
      comment = "preview api key"
      env     = ["Preview"]
    },

    # All environment variables
    {
      key     = "DEBUG"
      value   = "1"
      comment = "debug switch"
      env     = ["Production", "Preview"]
    },
    {
      key     = "DEBUG1"
      value   = "1"
    }
  ]
}
```

Import

tencentcloud_makers_project_envs can be imported using the `<zone_id>#<project_id>`, e.g.

```
$ terraform import tencentcloud_makers_project_envs.example zone-xxxx#proj-yyyy
```
