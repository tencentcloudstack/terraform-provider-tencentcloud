variable "role_document" {
  default = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}}]}"
}

variable "trigger_desc" {
  default = "{\"event\":\"cos:ObjectCreated:Put\",\"filter\":{\"Prefix\":\"\",\"Suffix\":\"\"}}"
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}
