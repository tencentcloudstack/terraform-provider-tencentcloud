/*
Provides a resource to create a waf web_shell

Example Usage

```hcl
resource "tencentcloud_waf_web_shell" "example" {
  domain = "demo.waf.com"
  status = 0
}
```

Import

waf web_shell can be imported using the id, e.g.

```
terraform import tencentcloud_waf_web_shell.example demo.waf.com
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWafWebShell() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafWebShellCreate,
		Read:   resourceTencentCloudWafWebShellRead,
		Update: resourceTencentCloudWafWebShellUpdate,
		Delete: resourceTencentCloudWafWebShellDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain.",
			},
			"status": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
				Description:  "Webshell status, 1: open; 0: closed; 2: log.",
			},
		},
	}
}

func resourceTencentCloudWafWebShellCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_web_shell.create")()
	defer inconsistentCheck(d, meta)()

	var domain string

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	d.SetId(domain)

	return resourceTencentCloudWafWebShellUpdate(d, meta)
}

func resourceTencentCloudWafWebShellRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_web_shell.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
		domain  = d.Id()
	)

	webShell, err := service.DescribeWafWebShellById(ctx, domain)
	if err != nil {
		return err
	}

	if webShell == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafWebShell` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if webShell.Domain != nil {
		_ = d.Set("domain", webShell.Domain)
	}

	if webShell.Status != nil {
		_ = d.Set("status", webShell.Status)
	}

	return nil
}

func resourceTencentCloudWafWebShellUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_web_shell.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		request = waf.NewModifyWebshellStatusRequest()
		domain  = d.Id()
	)

	webShellStatus := waf.WebshellStatus{}
	webShellStatus.Domain = helper.String(domain)

	if v, ok := d.GetOkExists("status"); ok {
		webShellStatus.Status = helper.IntUint64(v.(int))
	}

	request.Webshell = &webShellStatus

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyWebshellStatus(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update waf webShell failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafWebShellRead(d, meta)
}

func resourceTencentCloudWafWebShellDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_web_shell.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
