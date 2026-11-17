variable "region" {
  type        = string
  nullable    = false
  description = "COS Vectors region enabled for your account."

  validation {
    condition     = length(trimspace(var.region)) > 0 && var.region == trimspace(var.region)
    error_message = "Provide an explicit non-empty region enabled for COS Vectors."
  }
}

variable "vector_bucket_name" {
  type        = string
  description = "Full vector bucket name including the real -APPID suffix. Import existing buckets instead of recreating them."
}

variable "index_name" {
  type    = string
  default = "example-index"
}

variable "dimension" {
  type    = number
  default = 3
}

variable "distance_metric" {
  type    = string
  default = "cosine"
}

variable "policy_principal_qcs" {
  type        = string
  default     = null
  description = "Optional authorized CAM user principal QCS. Null skips policy creation."

  validation {
    condition     = var.policy_principal_qcs == null ? true : can(regex("^qcs::cam::uin/[0-9]+:uin/[0-9]+$", var.policy_principal_qcs))
    error_message = "Use null or an explicitly authorized CAM user QCS: qcs::cam::uin/ACCOUNT_UIN:uin/USER_UIN. Wildcards are not allowed in this allow-policy example."
  }
}
