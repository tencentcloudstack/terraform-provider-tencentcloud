variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "image_id" {
  default = "img-2xnn7dex"
}

variable "instance_type" {
  default = ""
}

variable "instance_name" {
  default = "terraform-testing"
}

variable "os_name" {
  default = "centos"
}

variable "image_name_regex" {
  default = "^CentOS\\s+7\\.3\\s+64\\w*"
}
