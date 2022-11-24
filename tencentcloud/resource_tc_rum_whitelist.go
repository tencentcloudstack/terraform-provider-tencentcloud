/*
Provides a resource to create a rum whitelist

Example Usage

```hcl
resource "tencentcloud_rum_whitelist" "whitelist" {
  instance_id = "rum-pasZKEI3RLgakj"
  remark = "white list remark"
  whitelist_uin = "20221122"
  # aid = ""
}

```
Import

rum whitelist can be imported using the id, e.g.
```
$ terraform import tencentcloud_rum_whitelist.whitelist whitelist_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRumWhitelist() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudRumWhitelistRead,
		Create: resourceTencentCloudRumWhitelistCreate,
		Update: resourceTencentCloudRumWhitelistUpdate,
		Delete: resourceTencentCloudRumWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID, such as taw-123.",
			},

			"remark": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remarks.",
			},

			"whitelist_uin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "uin: business identifier.",
			},

			"aid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Business identifier.",
			},

			"ttl": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "End time.",
			},

			"wid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auto-Increment allowlist ID.",
			},

			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creator ID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
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
		response   *rum.CreateWhitelistResponse
		instanceID string
		wid        string
	)

	if v, ok := d.GetOk("instance_id"); ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create rum whitelist failed, reason:%+v", logId, err)
		return err
	}

	wid = strconv.Itoa(int(*response.Response.ID))

	d.SetId(instanceID + FILED_SP + wid)
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
	wid := idSplit[1]

	whitelist, err := service.DescribeRumWhitelist(ctx, instanceID, wid)

	if err != nil {
		return err
	}

	if whitelist == nil {
		d.SetId("")
		return fmt.Errorf("resource `whitelist` %s does not exist", wid)
	}

	if whitelist.InstanceID != nil {
		_ = d.Set("instance_id", whitelist.InstanceID)
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

	if whitelist.ID != nil {
		_ = d.Set("wid", whitelist.ID)
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

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("remark") {
		return fmt.Errorf("`remark` do not support change now.")
	}

	if d.HasChange("whitelist_uin") {
		return fmt.Errorf("`whitelist_uin` do not support change now.")
	}

	if d.HasChange("aid") {
		return fmt.Errorf("`aid` do not support change now.")
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
	wid := idSplit[1]

	if err := service.DeleteRumWhitelistById(ctx, instanceID, wid); err != nil {
		return err
	}

	return nil
}
