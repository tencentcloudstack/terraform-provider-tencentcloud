resource "tencentcloud_ags_sandbox_tool" "this" {
  tool_name       = var.tool_name
  tool_type       = var.tool_type
  description     = var.description
  default_timeout = var.default_timeout
  persistent      = var.persistent
  role_arn        = var.role_arn
  tags            = var.tags

  network_configuration {
    network_mode = var.network.network_mode
    dynamic "vpc_config" {
      for_each = var.network.vpc_config == null ? [] : [var.network.vpc_config]
      content {
        subnet_ids         = vpc_config.value.subnet_ids
        security_group_ids = vpc_config.value.security_group_ids
      }
    }
  }

  dynamic "custom_configuration" {
    for_each = var.custom_configuration == null ? [] : [var.custom_configuration]
    content {
      image               = custom_configuration.value.image
      image_registry_type = custom_configuration.value.image_registry_type
      command             = custom_configuration.value.command
      args                = custom_configuration.value.args
      env                 = custom_configuration.value.env
      resources {
        cpu    = custom_configuration.value.resources.cpu
        memory = custom_configuration.value.resources.memory
      }
      dynamic "ports" {
        for_each = custom_configuration.value.ports
        content {
          name     = ports.value.name
          port     = ports.value.port
          protocol = ports.value.protocol
        }
      }
      dynamic "probe" {
        for_each = custom_configuration.value.probe == null ? [] : [custom_configuration.value.probe]
        content {
          ready_timeout_ms  = probe.value.ready_timeout_ms
          probe_timeout_ms  = probe.value.probe_timeout_ms
          probe_period_ms   = probe.value.probe_period_ms
          success_threshold = probe.value.success_threshold
          failure_threshold = probe.value.failure_threshold
          http_get {
            path   = probe.value.http_get.path
            port   = probe.value.http_get.port
            scheme = probe.value.http_get.scheme
          }
        }
      }
    }
  }

  dynamic "storage_mounts" {
    for_each = var.storage_mounts
    content {
      name       = storage_mounts.value.name
      mount_path = storage_mounts.value.mount_path
      read_only  = storage_mounts.value.read_only
      dynamic "cos" {
        for_each = storage_mounts.value.cos == null ? [] : [storage_mounts.value.cos]
        content {
          endpoint    = cos.value.endpoint
          bucket_name = cos.value.bucket_name
          bucket_path = cos.value.bucket_path
        }
      }
      dynamic "cfs" {
        for_each = storage_mounts.value.cfs == null ? [] : [storage_mounts.value.cfs]
        content {
          file_system_id = cfs.value.file_system_id
          path           = cfs.value.path
        }
      }
    }
  }

  dynamic "log_configuration" {
    for_each = var.log_configuration == null ? [] : [var.log_configuration]
    content {
      topic_id = log_configuration.value.topic_id
      files    = log_configuration.value.files
    }
  }

  timeouts {
    create = var.timeouts.create
    update = var.timeouts.update
    delete = var.timeouts.delete
    read   = var.timeouts.read
  }
}
