package tke

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesAddonConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesAddonConfigCreate,
		Read:   resourceTencentCloudKubernetesAddonConfigRead,
		Update: resourceTencentCloudKubernetesAddonConfigUpdate,
		Delete: resourceTencentCloudKubernetesAddonConfigDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of cluster.",
			},

			"addon_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of addon.",
			},

			"addon_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Version of addon.",
			},

			"raw_values": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Params of addon, base64 encoded json format.",
			},

			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of addon.",
			},

			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason of addon failed.",
			},
		},
	}
}

func resourceTencentCloudKubernetesAddonConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_addon_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	clusterId := d.Get("cluster_id").(string)
	addonName := d.Get("addon_name").(string)

	d.SetId(clusterId + tccommon.FILED_SP + addonName)
	return resourceTencentCloudKubernetesAddonConfigUpdate(d, meta)
}

func resourceTencentCloudKubernetesAddonConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_addon_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	addonName := idSplit[1]

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("addon_name", addonName)

	respData, err := service.DescribeKubernetesAddonById(ctx, clusterId, addonName)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_addon_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.AddonVersion != nil {
		_ = d.Set("addon_version", respData.AddonVersion)
	}

	if respData.RawValues != nil {
		rawValues := respData.RawValues
		base64DecodeValues, _ := base64.StdEncoding.DecodeString(*rawValues)
		jsonValues := string(base64DecodeValues)
		_ = d.Set("raw_values", jsonValues)
	}

	if respData.Phase != nil {
		_ = d.Set("phase", respData.Phase)
	}

	if respData.Reason != nil {
		_ = d.Set("reason", respData.Reason)
	}

	_ = addonName
	return nil
}

func resourceTencentCloudKubernetesAddonConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_addon_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	addonName := idSplit[1]

	needChange := false
	mutableArgs := []string{"addon_version", "raw_values"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tke.NewUpdateAddonRequest()

		request.ClusterId = &clusterId
		request.AddonName = &addonName

		if v, ok := d.GetOk("addon_version"); ok {
			request.AddonVersion = helper.String(v.(string))
		}

		if v, ok := d.GetOk("raw_values"); ok {
			jsonValues := helper.String(v.(string))
			rawValues := base64.StdEncoding.EncodeToString([]byte(*jsonValues))
			request.RawValues = &rawValues
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().UpdateAddonWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes addon failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = addonName
	return resourceTencentCloudKubernetesAddonConfigRead(d, meta)
}

func resourceTencentCloudKubernetesAddonConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_addon_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
