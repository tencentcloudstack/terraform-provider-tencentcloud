Provides an ephemeral resource that returns a synthetic short-lived
credential bundle. The returned value lives only inside a single
Terraform graph walk and is never persisted to the state file. This
reference is suitable for handing transient secrets to downstream
provider configuration blocks.

Example Usage

```hcl
ephemeral "tencentcloud_temp_credential" "this" {
  ttl_seconds = 900
}

provider "vault" {
  address = "https://vault.example.com"
  token   = ephemeral.tencentcloud_temp_credential.this.session_token
}
```
