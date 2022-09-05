/*
Provides a resource to create a monitor ssoAccount

Example Usage

```hcl
resource "tencentcloud_monitor_sso_account" "ssoAccount" {
  instance_id = ""
  user_id = ""
  notes = ""
  role {
			organization = ""
			role = ""

  }
}

```
Import

monitor ssoAccount can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_sso_account.ssoAccount ssoAccount_id
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorSsoAccount() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorSsoAccountRead,
		Create: resourceTencentCloudMonitorSsoAccountCreate,
		Update: resourceTencentCloudMonitorSsoAccountUpdate,
		Delete: resourceTencentCloudMonitorSsoAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "grafana instance id.",
			},

			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "sub account uin of specific user.",
			},

			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "account related description.",
			},

			"role": {
				Optional:    true,
				Computed:    true,
				Description: "grafana role.",
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
		request    = monitor.NewCreateSSOAccountRequest()
		response   *monitor.CreateSSOAccountResponse
		instanceId string
		userId     string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_id"); ok {
		request.UserId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notes"); ok {
		request.Notes = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreateSSOAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor ssoAccount failed, reason:%+v", logId, err)
		return err
	}

	//userId := *response.Response.UserId

	d.SetId(strings.Join([]string{instanceId, userId}, FILED_SP))
	return resourceTencentCloudMonitorSsoAccountRead(d, meta)
}

func resourceTencentCloudMonitorSsoAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_sso_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userId := idSplit[1]

	ssoAccount, err := service.DescribeMonitorSsoAccount(ctx, instanceId, userId)

	if err != nil {
		return err
	}

	if ssoAccount == nil {
		d.SetId("")
		return fmt.Errorf("resource `ssoAccount` %s does not exist", userId)
	}

	_ = d.Set("instance_id", instanceId)

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

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userId := idSplit[1]

	request.InstanceId = &instanceId
	request.UserId = &userId

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("user_id") {
		return fmt.Errorf("`user_id` do not support change now.")
	}

	if d.HasChange("notes") {
		return fmt.Errorf("`notes` do not support change now.")
	}

	if d.HasChange("role") {
		return fmt.Errorf("`role` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateSSOAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
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

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userId := idSplit[1]

	if err := service.DeleteMonitorSsoAccountById(ctx, instanceId, userId); err != nil {
		return err
	}

	return nil
}
