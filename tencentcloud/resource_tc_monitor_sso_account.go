/*
Provides a resource to create a monitor sso_account

Example Usage

```hcl
resource "tencentcloud_monitor_sso_account" "sso_account" {
  instance_id = &lt;nil&gt;
  user_id = &lt;nil&gt;
  notes = &lt;nil&gt;
  role {
		organization = &lt;nil&gt;
		role = &lt;nil&gt;

  }
}
```

Import

monitor sso_account can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_sso_account.sso_account sso_account_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMonitorSsoAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorSsoAccountCreate,
		Read:   resourceTencentCloudMonitorSsoAccountRead,
		Update: resourceTencentCloudMonitorSsoAccountUpdate,
		Delete: resourceTencentCloudMonitorSsoAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance id.",
			},

			"user_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sub account uin of specific user.",
			},

			"notes": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Account related description.",
			},

			"role": {
				Optional:    true,
				Computed:    true,
				Description: "Grafana role.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Grafana organization id string.",
						},
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Grafana role, one of {Admin,Editor,Viewer}.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMonitorSsoAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_sso_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreateSSOAccountRequest()
		response = monitor.NewCreateSSOAccountResponse()
		userId   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		request.UserId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notes"); ok {
		request.Notes = helper.String(v.(string))
	}

	if v, _ := d.GetOk("role"); v != nil {
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreateSSOAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor ssoAccount failed, reason:%+v", logId, err)
		return err
	}

	userId = *response.Response.UserId
	d.SetId(userId)

	return resourceTencentCloudMonitorSsoAccountRead(d, meta)
}

func resourceTencentCloudMonitorSsoAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_sso_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	ssoAccountId := d.Id()

	ssoAccount, err := service.DescribeMonitorSsoAccountById(ctx, userId)
	if err != nil {
		return err
	}

	if ssoAccount == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorSsoAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ssoAccount.InstanceId != nil {
		_ = d.Set("instance_id", ssoAccount.InstanceId)
	}

	if ssoAccount.UserId != nil {
		_ = d.Set("user_id", ssoAccount.UserId)
	}

	if ssoAccount.Notes != nil {
		_ = d.Set("notes", ssoAccount.Notes)
	}

	if ssoAccount.Role != nil {
	}

	return nil
}

func resourceTencentCloudMonitorSsoAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_sso_account.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateSSOAccountRequest()

	ssoAccountId := d.Id()

	request.UserId = &userId

	immutableArgs := []string{"instance_id", "user_id", "notes", "role"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateSSOAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor ssoAccount failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorSsoAccountRead(d, meta)
}

func resourceTencentCloudMonitorSsoAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_sso_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	ssoAccountId := d.Id()

	if err := service.DeleteMonitorSsoAccountById(ctx, userId); err != nil {
		return err
	}

	return nil
}
