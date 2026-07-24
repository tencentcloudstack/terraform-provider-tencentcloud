Provides a resource to create a SCF (Serverless Cloud Function) trigger.

Example Usage

```hcl
resource "tencentcloud_scf_function" "function" {
  name    = "tf-example-function"
  runtime = "Nodejs16.13"
  handler = "index.main"
  memory_size = 128
  timeout     = 3
  cos_bucket_name = "tf-example-bucket-1300000000"
  cos_object_name = "function.zip"
}

resource "tencentcloud_scf_trigger" "timer" {
  function_name = tencentcloud_scf_function.function.name
  namespace     = "default"
  trigger_name  = "tf-example-trigger"
  type          = "timer"
  trigger_desc  = jsonencode({ cron = "*/5 * * * * * *" })
  enable        = "OPEN"
  description   = "tf example trigger"
  qualifier     = "$DEFAULT"
  custom_argument = "Information"
}
```

Import

SCF trigger can be imported using the composite id `function_name#namespace#trigger_name`, e.g.

```
terraform import tencentcloud_scf_trigger.timer tf-example-function#default#tf-example-trigger
```
