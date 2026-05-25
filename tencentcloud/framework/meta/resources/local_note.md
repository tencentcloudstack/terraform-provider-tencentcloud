Provides a local-only resource that stores a free-form note inside the
Terraform state. This reference resource does not call any cloud API; it
exists to demonstrate the framework resource lifecycle (Create / Read /
Update / Delete) and serves as a template for new framework resources.

Example Usage

```hcl
resource "tencentcloud_local_note" "example" {
  content = "hello, framework"
}
```

Import

A local note can be imported using its id (the SHA-256 of `content`):

```
terraform import tencentcloud_local_note.example 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
```
