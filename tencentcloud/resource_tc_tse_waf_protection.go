package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
)

func resourceTencentCloudTseWafProtection() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_tse_waf_protection.create")()
	defer inconsistentCheck(d, meta)()

	var gatewayId string
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
	}

	d.SetId(gatewayId)

	return resourceTencentCloudTseWafProtectionUpdate(d, meta)
}

func resourceTencentCloudTseWafProtectionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_waf_protection.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

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
	defer logElapsed("resource.tencentcloud_tse_waf_protection.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().OpenWafProtection(request)
				if e != nil {
					return retryError(e)
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
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CloseWafProtection(request)
				if e != nil {
					return retryError(e)
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
	defer logElapsed("resource.tencentcloud_tse_waf_protection.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
