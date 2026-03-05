data "tencentcloud_clb_instances" "example_internal" {
  network_type = "INTERNAL"
}

data "tencentcloud_clb_instances" "example_all" {
}

output "clb_instances_with_zones" {
  value = [
    for clb in data.tencentcloud_clb_instances.example_internal.clb_list : {
      clb_id       = clb.clb_id
      clb_name     = clb.clb_name
      network_type = clb.network_type
      zones        = clb.zones
      vpc_id       = clb.vpc_id
    }
  ]
  description = "Internal CLB instances with their zones information"
}

output "all_clb_instances" {
  value = [
    for clb in data.tencentcloud_clb_instances.example_all.clb_list : {
      clb_id       = clb.clb_id
      clb_name     = clb.clb_name
      network_type = clb.network_type
      zones        = try(clb.zones, null)
    }
  ]
  description = "All CLB instances with zones field (may be null for some instances)"
}
