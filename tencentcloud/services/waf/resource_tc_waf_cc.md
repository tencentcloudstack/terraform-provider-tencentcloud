Provides a resource to create a WAF cc

Example Usage

If advance is 0(IP model)

```hcl
resource "tencentcloud_waf_cc" "example" {
  domain      = "www.demo.com"
  name        = "tf-example"
  status      = 1
  advance     = "0"
  limit       = "60"
  interval    = "60"
  url         = "/cc_demo"
  match_func  = 0
  action_type = "22"
  priority    = 50
  valid_time  = 600
  edition     = "sparta-waf"
  type        = 1
  logical_op  = "and"
  options_arr = jsonencode(
    [
      {
        "key" : "URL",
        "args" : [
          "=cHJlZml4"
        ],
        "match" : "2",
        "encodeflag" : true
      },
      {
        "key" : "Method",
        "args" : [
          "=POST"  # if encodeflag is false, parameter value needs to be prefixed with an = sign.
        ],
        "match" : "0",
        "encodeflag" : false
      },
      {
        "key" : "Post",
        "args" : [
          "S2V5=VmFsdWU"
        ],
        "match" : "0",
        "encodeflag" : true
      },
      {
        "key" : "Referer",
        "args" : [
          "="
        ],
        "match" : "12",
        "encodeflag" : true
      },
      {
        "key" : "Cookie",
        "args" : [
          "S2V5=VmFsdWU"
        ],
        "match" : "3",
        "encodeflag" : true
      },
      {
        "key" : "IPLocation",
        "args" : [
          "=eyJMYW5nIjoiY24iLCJBcmVhcyI6W3siQ291bnRyeSI6IuWbveWkliJ9XX0"
        ],
        "match" : "13",
        "encodeflag" : true
      }
    ]
  )
}
```

If advance is 1(SESSION model)

```hcl
resource "tencentcloud_waf_cc" "example" {
  domain          = "news.bots.icu"
  name            = "tf-example"
  status          = 1
  advance         = "1"
  limit           = "60"
  interval        = "60"
  url             = "/cc_demo"
  match_func      = 0
  action_type     = "22"
  priority        = 50
  valid_time      = 600
  edition         = "sparta-waf"
  type            = 1
  session_applied = [0]
  limit_method    = "only_limit"
  logical_op      = "or"
  cel_rule        = "(has(request.url) && request.url.startsWith('/prefix')) && (has(request.method) && request.method == 'POST')"
}
```
