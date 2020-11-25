variable "bucket-name" {
  default = "bucket-test-1258798060"
}

variable "object-name" {
  default = "object-test"
}

variable "acl" {
  default = "public-read"
}

variable "object-content" {
  default = "terraform tencent cloud cos object"
}

variable "policy" {
  default = <<EOF
{
  "version": "2.0",
  "Statement": [
    {
      "Principal": {
        "qcs": [
          "qcs::cam::uin/100010835595:uin/100014918835"
        ]
      },
      "Action": [
        "name/cos:DeleteBucket",
        "name/cos:PutBucketACL"
      ],
      "Effect": "allow",
      "Resource": [
        "qcs::cos:ap-nanjing:uid/1259649581:hhermanwang-1259649581/*"
      ]
    }
  ]
}
EOF
}