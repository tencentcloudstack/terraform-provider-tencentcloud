variable "password" {
  default = "Gail@1234"
}

variable "phone_num" {
  default = "13631555963"
}

variable "country_code" {
  default = "86"
}

variable "email" {
  default = "1234@qq.com"
}

variable "policy_document" {
  default = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]},{\"action\":[\"name/cos:PutObject\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
}

variable "role_document" {
  default = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}}]}"
}
