package tse

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionCreate,
		Read:   resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionRead,
		Update: resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionUpdate,
		Delete: resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"source_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Resource type bound to the IP restriction plugin: route|service.",
			},

			"source_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Route or service ID.",
			},

			"enabled": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the plugin.",
			},

			"restriction_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "IP restriction type: whiteList|blackList.",
			},

			"address_list": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IP/CIDR address list.",
			},
		},
	}
}

func resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cloud_native_api_gateway_ip_restriction.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		gatewayId  string
		sourceType string
		sourceId   string
	)

	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
	}

	if v, ok := d.GetOk("source_type"); ok {
		sourceType = v.(string)
	}

	if v, ok := d.GetOk("source_id"); ok {
		sourceId = v.(string)
	}

	d.SetId(strings.Join([]string{gatewayId, sourceType, sourceId}, tccommon.FILED_SP))

	return resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionUpdate(d, meta)
}

func resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cloud_native_api_gateway_ip_restriction.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	sourceType := idSplit[1]
	sourceId := idSplit[2]

	request := tse.NewDescribeCloudNativeAPIGatewayIPRestrictionRequest()
	request.GatewayId = helper.String(gatewayId)
	request.SourceType = helper.String(sourceType)
	request.SourceId = helper.String(sourceId)

	var (
		result *tse.DescribeCloudNativeAPIGatewayIPRestrictionResponse
		reqErr error
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, reqErr = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().DescribeCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)
		if reqErr != nil {
			return tccommon.RetryError(reqErr)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("api returned no response"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read tse_cloud_native_api_gateway_ip_restriction failed, reason:%+v", logId, err)
		return err
	}

	if result.Response.Result == nil {
		log.Printf("[WARN] read tse_cloud_native_api_gateway_ip_restriction id=%s empty result", d.Id())
		if d.IsNewResource() {
			return fmt.Errorf("tse_cloud_native_api_gateway_ip_restriction [%s] not found after creation", d.Id())
		}
		d.SetId("")
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("source_type", sourceType)
	_ = d.Set("source_id", sourceId)

	ipRestriction := result.Response.Result
	if ipRestriction.Enabled != nil {
		_ = d.Set("enabled", ipRestriction.Enabled)
	}
	if ipRestriction.RestrictionType != nil {
		_ = d.Set("restriction_type", ipRestriction.RestrictionType)
	}
	if ipRestriction.AddressList != nil {
		addressList := make([]interface{}, 0, len(ipRestriction.AddressList))
		for _, v := range ipRestriction.AddressList {
			if v != nil {
				addressList = append(addressList, *v)
			}
		}
		_ = d.Set("address_list", addressList)
	}

	return nil
}

func resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cloud_native_api_gateway_ip_restriction.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := tse.NewCreateOrModifyCloudNativeAPIGatewayIPRestrictionRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	request.GatewayId = helper.String(idSplit[0])
	request.SourceType = helper.String(idSplit[1])
	request.SourceId = helper.String(idSplit[2])

	if v, ok := d.GetOkExists("enabled"); ok {
		request.Enabled = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("restriction_type"); ok {
		request.RestrictionType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_list"); ok {
		addressList := v.(*schema.Set).List()
		for i := range addressList {
			address := addressList[i].(string)
			request.AddressList = append(request.AddressList, &address)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().CreateOrModifyCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("api returned no response"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse_cloud_native_api_gateway_ip_restriction failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionRead(d, meta)
}

func resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cloud_native_api_gateway_ip_restriction.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	sourceType := idSplit[1]
	sourceId := idSplit[2]

	request := tse.NewDeleteCloudNativeAPIGatewayIPRestrictionRequest()
	request.GatewayId = helper.String(gatewayId)
	request.SourceType = helper.String(sourceType)
	request.SourceId = helper.String(sourceId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().DeleteCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("api returned no response"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete tse_cloud_native_api_gateway_ip_restriction failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
