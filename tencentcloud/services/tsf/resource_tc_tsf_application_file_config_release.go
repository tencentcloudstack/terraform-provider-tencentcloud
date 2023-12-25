package tsf

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTsfApplicationFileConfigRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfApplicationFileConfigReleaseCreate,
		Read:   resourceTencentCloudTsfApplicationFileConfigReleaseRead,
		Delete: resourceTencentCloudTsfApplicationFileConfigReleaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "File config id.",
			},

			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Group Id.",
			},

			"release_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "release Description.",
			},
		},
	}
}

func resourceTencentCloudTsfApplicationFileConfigReleaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_application_file_config_release.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = tsf.NewReleaseFileConfigRequest()
		configId string
		groupId  string
	)
	if v, ok := d.GetOk("config_id"); ok {
		configId = v.(string)
		request.ConfigId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("release_desc"); ok {
		request.ReleaseDesc = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().ReleaseFileConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf applicationFileConfigRelease failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(configId + tccommon.FILED_SP + groupId)

	return resourceTencentCloudTsfApplicationFileConfigReleaseRead(d, meta)
}

func resourceTencentCloudTsfApplicationFileConfigReleaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_application_file_config_release.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	groupId := idSplit[1]

	applicationFileConfigRelease, err := service.DescribeTsfApplicationFileConfigReleaseById(ctx, configId, groupId)
	if err != nil {
		return err
	}

	if applicationFileConfigRelease == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfApplicationFileConfigRelease` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationFileConfigRelease.ConfigId != nil {
		_ = d.Set("config_id", applicationFileConfigRelease.ConfigId)
	}

	if applicationFileConfigRelease.GroupId != nil {
		_ = d.Set("group_id", applicationFileConfigRelease.GroupId)
	}

	if applicationFileConfigRelease.ReleaseDesc != nil {
		_ = d.Set("release_desc", applicationFileConfigRelease.ReleaseDesc)
	}

	return nil
}

func resourceTencentCloudTsfApplicationFileConfigReleaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_application_file_config_release.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	configId := idSplit[0]
	groupId := idSplit[1]

	if err := service.DeleteTsfApplicationFileConfigReleaseById(ctx, configId, groupId); err != nil {
		return err
	}

	return nil
}
