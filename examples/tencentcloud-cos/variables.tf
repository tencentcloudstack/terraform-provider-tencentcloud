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
        "*"
      ]
    }
  ]
}
EOF
}

variable "acl_body" {
  default = <<EOF
<AccessControlPolicy>
    <Owner>
        <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
    </Owner>
    <AccessControlList>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
                <URI>http://cam.qcloud.com/groups/global/AllUsers</URI>
            </Grantee>
            <Permission>READ</Permission>
        </Grant>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
                <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            </Grantee>
            <Permission>WRITE</Permission>
        </Grant>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
                <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            </Grantee>
            <Permission>READ_ACP</Permission>
        </Grant>
    </AccessControlList>
</AccessControlPolicy>
EOF
}