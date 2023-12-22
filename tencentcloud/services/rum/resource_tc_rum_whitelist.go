package rum

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRumWhitelist() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_rum_whitelist.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRumClient().CreateWhitelist(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	d.SetId(instanceID + tccommon.FILED_SP + wid)
	return resourceTencentCloudRumWhitelistRead(d, meta)
}

func resourceTencentCloudRumWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_rum_whitelist.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RumService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
	defer tccommon.LogElapsed("resource.tencentcloud_rum_whitelist.update")()
	defer tccommon.InconsistentCheck(d, meta)()

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
	defer tccommon.LogElapsed("resource.tencentcloud_rum_whitelist.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RumService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
