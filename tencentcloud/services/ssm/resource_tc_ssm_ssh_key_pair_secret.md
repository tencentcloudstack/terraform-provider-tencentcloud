Provides a resource to create a ssm ssh key pair secret

Example Usage

```hcl
resource "tencentcloud_kms_key" "example" {
  alias                = "tf-example-kms-key"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = true

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_ssh_key_pair_secret" "example" {
  secret_name   = "tf-example"
  project_id    = 0
  description   = "desc."
  kms_key_id    = tencentcloud_kms_key.example.id
  ssh_key_name  = "tf_example_ssh"
  status        = "Enabled"
  clean_ssh_key = true

  tags = {
    createdBy = "terraform"
  }
}
```

Import

ssm ssh_key_pair_secret can be imported using the id, e.g.

```
terraform import tencentcloud_ssm_ssh_key_pair_secret.ssh_key_pair_secret ssh_key_pair_secret_name
```