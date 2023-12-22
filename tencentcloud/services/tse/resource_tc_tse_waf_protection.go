package tse

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
)

func ResourceTencentCloudTseWafProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseWafProtectionCreate,
		Read:   resourceTencentCloudTseWafProtectionRead,
		Update: resourceTencentCloudTseWafProtectionUpdate,
		Delete: resourceTencentCloudTseWafProtectionDelete,

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The type of protection resource. Reference value: `Global`: instance, `Service`: service, `Route`: route, `Object`: obejct (This interface does not currently support this type).",
			},

			"list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Means the list of services or routes when the resource type `Type` is `Service` or `Route`.",
			},

			"operate": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`open`: open the protection, `close`: close the protection.",
			},

			"global_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Global protection status.",
			},
		},
	}
}

func resourceTencentCloudTseWafProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_waf_protection.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var gatewayId string
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
	}

	d.SetId(gatewayId)

	return resourceTencentCloudTseWafProtectionUpdate(d, meta)
}

func resourceTencentCloudTseWafProtectionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_waf_protection.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	gatewayId := d.Id()

	wafProtection, err := service.DescribeTseWafProtectionById(ctx, gatewayId)
	if err != nil {
		return err
	}

	if wafProtection == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseWafProtection` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)

	if wafProtection.GlobalStatus != nil {
		_ = d.Set("global_status", wafProtection.GlobalStatus)
	}

	return nil
}

func resourceTencentCloudTseWafProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_waf_protection.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		resourceType string
		resourceList []*string
	)

	gatewayId := d.Id()
	if v, ok := d.GetOk("type"); ok {
		resourceType = v.(string)
	}

	if v, ok := d.GetOk("list"); ok {
		listSet := v.(*schema.Set).List()
		for i := range listSet {
			list := listSet[i].(string)
			resourceList = append(resourceList, &list)
		}
	}

	if v, ok := d.GetOk("operate"); ok {
		operate := v.(string)
		if operate == "open" {
			request := tse.NewOpenWafProtectionRequest()
			request.GatewayId = &gatewayId
			request.Type = &resourceType
			request.List = resourceList
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().OpenWafProtection(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s open tse wafProtection failed, reason:%+v", logId, err)
				return err
			}
		} else if operate == "close" {
			request := tse.NewCloseWafProtectionRequest()
			request.GatewayId = &gatewayId
			request.Type = &resourceType
			request.List = resourceList
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().CloseWafProtection(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s close tse wafProtection failed, reason:%+v", logId, err)
				return err
			}
		} else {
			return fmt.Errorf("")
		}
	}

	return resourceTencentCloudTseWafProtectionRead(d, meta)
}

func resourceTencentCloudTseWafProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_waf_protection.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
