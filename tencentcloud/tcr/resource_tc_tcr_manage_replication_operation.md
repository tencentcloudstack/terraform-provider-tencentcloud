Provides a resource to start a tcr instance replication operation

Example Usage

Sync source tcr instance to target instance

Synchronize an existing tcr instance to the destination instance. This operation is often used in the cross-multiple region scenario.
Assume you have had two TCR instances before this operation. This example shows how to sync a tcr instance from ap-guangzhou(gz) to ap-shanghai(sh).

```hcl
# tcr instance on ap-guangzhou
resource "tencentcloud_tcr_instance" "example_gz" {
  name          = "tf-example-tcr-gz"
  instance_type = "premium"
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example_gz" {
  instance_id    = tencentcloud_tcr_instance.example_gz.id
  name           = "tf_example_ns_gz"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

# tcr instance on ap-shanghai
resource "tencentcloud_tcr_instance" "example_sh" {
  name          = "tf-example-tcr-sh"
  instance_type = "premium"
  delete_bucket = true
}

resource "tencentcloud_tcr_namespace" "example_sh" {
  instance_id    = tencentcloud_tcr_instance.example_sh.id
  name           = "tf_example_ns_sh"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

```

```hcl
# Run this on region ap-guangzhou
locals {
  src_id  = tencentcloud_tcr_instance.example_gz.id
  dest_id = tencentcloud_tcr_instance.example_sh.id
  src_ns  = tencentcloud_tcr_namespace.example_gz.name
  dest_ns = tencentcloud_tcr_instance.example_sh.id
}

variable "tcr_region_map" {
  default = {
    "ap-guangzhou"     = 1
    "ap-shanghai"      = 4
    "ap-hongkong"      = 5
    "ap-beijing"       = 8
    "ap-singapore"     = 9
    "na-siliconvalley" = 15
    "ap-chengdu"       = 16
    "eu-frankfurt"     = 17
    "ap-seoul"         = 18
    "ap-chongqing"     = 19
    "ap-mumbai"        = 21
    "na-ashburn"       = 22
    "ap-bangkok"       = 23
    "eu-moscow"        = 24
    "ap-tokyo"         = 25
    "ap-nanjing"       = 33
    "ap-taipei"        = 39
    "ap-jakarta"       = 72
  }
}

resource "tencentcloud_tcr_manage_replication_operation" "example_sync" {
  source_registry_id      = local.src_id
  destination_registry_id = local.dest_id
  rule {
    name           = "tf_example_sync_gz_to_sh"
    dest_namespace = local.dest_ns
    override       = true
    filters {
      type  = "name"
      value = join("/", [local.src_ns, "**"])
    }
    filters {
      type  = "tag"
      value = ""
    }
    filters {
      type  = "resource"
      value = ""
    }
  }
  description           = "example for tcr sync operation"
  destination_region_id = var.tcr_region_map["ap-shanghai"] # 4
  peer_replication_option {
    peer_registry_uin       = ""
    peer_registry_token     = ""
    enable_peer_replication = false
  }
}
```