variable "vpc_cidr" {
  default = "10.1.0.0/21"
}

variable "bandwidth" {
  default = 100
}

variable "max_concurrent" {
  default = 1000000
}
