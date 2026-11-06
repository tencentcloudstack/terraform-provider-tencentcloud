variable "tool_name" {
  type        = string
  description = "Unique Tool name in the selected account."
}

variable "tool_type" {
  type        = string
  description = "browser, code-interpreter or custom."
  validation {
    condition     = contains(["browser", "code-interpreter", "custom"], var.tool_type)
    error_message = "Only the three V1 Tool types are supported."
  }
}

variable "description" {
  type        = string
  default     = null
  description = "Omitted accepts remote defaults; empty requests explicit clearing."
}

variable "default_timeout" {
  type        = string
  default     = null
  description = "Default instance lifetime; changes replace the Tool."
}

variable "persistent" {
  type        = bool
  default     = null
  description = "Future instance persistence; does not create an instance."
}

variable "role_arn" {
  type        = string
  default     = null
  description = "Existing CAM data-access role. This module never creates or deletes it."
}

variable "tags" {
  type        = map(string)
  default     = null
  description = "User tags; {} requests removal of all user tags, excluding immutable system tags."
}

variable "network" {
  type = object({
    network_mode = string
    vpc_config = optional(object({
      subnet_ids         = set(string)
      security_group_ids = set(string)
    }))
  })
  description = "Network mode and existing VPC bindings. No dependencies are created implicitly."
}

variable "custom_configuration" {
  type = object({
    image               = string
    image_registry_type = string
    command             = optional(list(string))
    args                = optional(list(string))
    env                 = optional(map(string))
    resources = object({
      cpu    = string
      memory = string
    })
    ports = optional(set(object({
      name     = string
      port     = number
      protocol = string
    })), [])
    probe = optional(object({
      http_get = object({
        path   = string
        port   = number
        scheme = string
      })
      ready_timeout_ms  = optional(number)
      probe_timeout_ms  = optional(number)
      probe_period_ms   = optional(number)
      success_threshold = optional(number)
      failure_threshold = optional(number)
    }))
  })
  default     = null
  description = "Required for custom Tools. Uses a prebuilt, AGS-compatible image and non-sensitive environment values."
}

variable "storage_mounts" {
  type = set(object({
    name       = string
    mount_path = string
    read_only  = optional(bool, false)
    cos = optional(object({
      endpoint    = string
      bucket_name = string
      bucket_path = string
    }))
    cfs = optional(object({
      file_system_id = string
      path           = string
    }))
  }))
  default     = []
  description = "Existing COS/CFS mounts. Storage lifecycle and retention remain with their owners."
}

variable "log_configuration" {
  type = object({
    topic_id = string
    files    = set(string)
  })
  default     = null
  description = "Existing CLS Topic and approved log files. No Topic or role is created."
}

variable "timeouts" {
  type = object({
    create = optional(string, "30m")
    update = optional(string, "30m")
    delete = optional(string, "30m")
    read   = optional(string, "5m")
  })
  default     = {}
  description = "Terraform waiting budgets, separate from instance lifetime."
}
