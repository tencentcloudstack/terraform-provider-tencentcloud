Provides a resource to create a CLS DataSight console.

Example Usage

If login_mode is set to 0

```hcl
resource "tencentcloud_cls_console" "example" {
  access_mode     = ["public", "internal"]
  login_mode      = 0
  domain_prefix   = "datasight"
  remarks         = "remarks."
  intranet_type   = 1
  intranet_region = "ap-chongqing"

  accounts {
    user_name  = "your_username1"
    password   = "your_password1"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
    email      = "demo@example.com"
  }

  accounts {
    user_name  = "your_username2"
    password   = "your_password2"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
    email      = "demo@example.com"
  }

  accounts {
    user_name  = "your_username3"
    password   = "your_password3"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
    email      = "demo@example.com"
  }

  menus = [
    "/cls/search",
    "/cls/dashboard",
    "/cls/alarm",
    "/cls/process",
  ]

  access_control_rules {
    access_mode = "public"
    action      = "DENY"
    cidr_blocks = [
      "1.1.1.1",
      "2.2.2.2",
      "3.3.3.3",
    ]
  }

  access_control_rules {
    access_mode = "internal"
    action      = "DENY"
    cidr_blocks = [
      "4.4.4.4",
      "5.5.5.5"
    ]
  }

  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

If login_mode is set to 1

```hcl
resource "tencentcloud_cls_console" "example" {
  access_mode     = ["internal"]
  login_mode      = 1
  domain_prefix   = "datasight"
  remarks         = "remarks."
  intranet_type   = 1
  intranet_region = "ap-chongqing"

  anonymous_login {
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
  }

  menus = [
    "/cls/search",
    "/cls/dashboard",
    "/cls/alarm",
  ]

  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

If login_mode is set to 2

```hcl
resource "tencentcloud_cls_console" "example2" {
  access_mode     = ["internal"]
  login_mode      = 2
  domain_prefix   = "datasight2"
  remarks         = "remarks."
  intranet_type   = 1
  intranet_region = "ap-chongqing"

  auth_roles {
    role_name  = "role1"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
  }

  auth_roles {
    role_name  = "role2"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
  }

  auth_roles {
    role_name  = "role3"
    secret_id  = "your_secret_id"
    secret_key = "your_secret_key"
  }

  menus = [
    "/cls/search",
    "/cls/dashboard",
    "/cls/alarm",
    "/cls/process",
  ]

  access_control_rules {
    access_mode = "internal"
    action      = "ACCEPT"
    cidr_blocks = [
      "1.1.1.1",
      "2.2.2.2"
    ]
  }

  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

Import

CLS DataSight console can be imported using the id, e.g.

```
terraform import tencentcloud_cls_console.example clsconsole-0d59e193
```
