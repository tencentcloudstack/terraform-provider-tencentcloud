/*
Provides a resource to create a rum whitelist

Example Usage

```hcl
resource "tencentcloud_rum_whitelist" "whitelist" {
  instance_i_d = &lt;nil&gt;
  remark = &lt;nil&gt;
  whitelist_uin = &lt;nil&gt;
  aid = &lt;nil&gt;
        }
```

Import

rum whitelist can be imported using the id, e.g.

```
terraform import tencentcloud_rum_whitelist.whitelist whitelist_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudRumWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRumWhitelistCreate,
		Read:   resourceTencentCloudRumWhitelistRead,
		Update: resourceTencentCloudRumWhitelistUpdate,
		Delete: resourceTencentCloudRumWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_i_d": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, such as taw-123.",
			},

			"remark": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Remarks.",
			},

			"whitelist_uin": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Uin: business identifier.",
			},

			"aid": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Business identifier.",
			},

			"ttl": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"wid": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Auto-Increment allowlist ID.",
			},

			"create_user": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creator ID.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},
		},
	}
}

func resourceTencentCloudRumWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_whitelist.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = rum.NewCreateWhitelistRequest()
		response   = rum.NewCreateWhitelistResponse()
		instanceID string
		iD         string
	)
	if v, ok := d.GetOk("instance_i_d"); ok {
		instanceID = v.(string)
		request.InstanceID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("whitelist_uin"); ok {
		request.WhitelistUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("aid"); ok {
		request.Aid = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().CreateWhitelist(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create rum whitelist failed, reason:%+v", logId, err)
		return err
	}

	instanceID = *response.Response.InstanceID
	d.SetId(strings.Join([]string{instanceID, iD}, FILED_SP))

	return resourceTencentCloudRumWhitelistRead(d, meta)
}

func resourceTencentCloudRumWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_whitelist.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceID := idSplit[0]
	iD := idSplit[1]

	whitelist, err := service.DescribeRumWhitelistById(ctx, instanceID, iD)
	if err != nil {
		return err
	}

	if whitelist == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RumWhitelist` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if whitelist.InstanceID != nil {
		_ = d.Set("instance_i_d", whitelist.InstanceID)
	}

	if whitelist.Remark != nil {
		_ = d.Set("remark", whitelist.Remark)
	}

	if whitelist.WhitelistUin != nil {
		_ = d.Set("whitelist_uin", whitelist.WhitelistUin)
	}

	if whitelist.Aid != nil {
		_ = d.Set("aid", whitelist.Aid)
	}

	if whitelist.Ttl != nil {
		_ = d.Set("ttl", whitelist.Ttl)
	}

	if whitelist.Wid != nil {
		_ = d.Set("wid", whitelist.Wid)
	}

	if whitelist.CreateUser != nil {
		_ = d.Set("create_user", whitelist.CreateUser)
	}

	if whitelist.CreateTime != nil {
		_ = d.Set("create_time", whitelist.CreateTime)
	}

	return nil
}

func resourceTencentCloudRumWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_whitelist.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := rum.NewRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceID := idSplit[0]
	iD := idSplit[1]

	request.InstanceID = &instanceID
	request.ID = &iD

	immutableArgs := []string{"instance_i_d", "remark", "whitelist_uin", "aid", "ttl", "wid", "create_user", "create_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update rum whitelist failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRumWhitelistRead(d, meta)
}

func resourceTencentCloudRumWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_whitelist.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceID := idSplit[0]
	iD := idSplit[1]

	if err := service.DeleteRumWhitelistById(ctx, instanceID, iD); err != nil {
		return err
	}

	return nil
}
