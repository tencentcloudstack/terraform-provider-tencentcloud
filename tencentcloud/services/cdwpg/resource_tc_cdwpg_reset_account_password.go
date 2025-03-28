package cdwpg

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwpgv20201230 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCdwpgResetAccountPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdwpgResetAccountPasswordCreate,
		Read:   resourceTencentCloudCdwpgResetAccountPasswordRead,
		Update: resourceTencentCloudCdwpgResetAccountPasswordUpdate,
		Delete: resourceTencentCloudCdwpgResetAccountPasswordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id.",
			},

			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Username.",
			},

			"new_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "New password.",
			},
		},
	}
}

func resourceTencentCloudCdwpgResetAccountPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_reset_account_password.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)
	userName := d.Get("user_name").(string)

	d.SetId(strings.Join([]string{instanceId, userName}, tccommon.FILED_SP))

	return resourceTencentCloudCdwpgResetAccountPasswordUpdate(d, meta)
}

func resourceTencentCloudCdwpgResetAccountPasswordRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_reset_account_password.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CdwpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	respData, err := service.DescribeCdwpgAccountById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `cdwpg_account` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.UserName != nil {
		_ = d.Set("user_name", respData.UserName)
	}

	_ = userName
	return nil
}

func resourceTencentCloudCdwpgResetAccountPasswordUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_reset_account_password.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	needChange := false
	mutableArgs := []string{"new_password"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cdwpgv20201230.NewResetAccountPasswordRequest()

		request.InstanceId = helper.String(instanceId)

		request.UserName = helper.String(userName)

		if v, ok := d.GetOk("new_password"); ok {
			request.NewPassword = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwpgV20201230Client().ResetAccountPasswordWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update cdwpg account failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudCdwpgResetAccountPasswordRead(d, meta)
}

func resourceTencentCloudCdwpgResetAccountPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdwpg_reset_account_password.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
