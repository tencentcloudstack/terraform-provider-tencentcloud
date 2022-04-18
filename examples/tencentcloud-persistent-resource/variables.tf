variable "default_az" {
  default = "ap-guangzhou-3"
}

variable "cvm_name" {
  default = "keep-cvm"
}

# this persist exist and barely removed
variable "image_id" {
  default = "img-2lr9q49h"
}

variable "vpn_gw" {
  default = "kepp-vpn-gw"
}

variable "vpn_conn" {
  default = "keep-vpn-conn"
}

variable "vpn_cgw" {
  default = "keep-vpn-cgw"
}

variable "cam_role_basic" {
  default = "keep-cam-role"
}

variable "cam_policy_basic" {
  default = "keep-cam-policy"
}

variable "cam_group_basic" {
  default = "keep-cam-group"
}