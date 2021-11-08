variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "instance_type" {
  default = "SA1.SMALL1"
}

variable "min_size" {
  default = 0
}

variable "max_size" {
  default = 5
}

variable "desired_capacity" {
  default = 0
}

variable "cvm_product" {
  default = "cvm"
}
