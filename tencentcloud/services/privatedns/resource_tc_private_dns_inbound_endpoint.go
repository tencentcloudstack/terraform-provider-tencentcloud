package privatedns

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	privatednsv20201028 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPrivateDnsInboundEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPrivateDnsInboundEndpointCreate,
		Read:   resourceTencentCloudPrivateDnsInboundEndpointRead,
		Update: resourceTencentCloudPrivateDnsInboundEndpointUpdate,
		Delete: resourceTencentCloudPrivateDnsInboundEndpointDelete,
		Schema: map[string]*schema.Schema{
			"endpoint_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name.",
			},

			"endpoint_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region.",
			},

			"endpoint_vpc": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},

			"subnet_ip": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Subnet ID.",
						},
						"subnet_vip": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "IP address.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudPrivateDnsInboundEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_inbound_endpoint.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = privatednsv20201028.NewCreateInboundEndpointRequest()
		response   = privatednsv20201028.NewCreateInboundEndpointResponse()
		endpointId string
	)

	if v, ok := d.GetOk("endpoint_name"); ok {
		request.EndpointName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("endpoint_region"); ok {
		request.EndpointRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("endpoint_vpc"); ok {
		request.EndpointVpc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ip"); ok {
		for _, item := range v.([]interface{}) {
			subnetIpMap := item.(map[string]interface{})
			subnetIpInfo := privatednsv20201028.SubnetIpInfo{}
			if v, ok := subnetIpMap["subnet_id"].(string); ok && v != "" {
				subnetIpInfo.SubnetId = helper.String(v)
			}

			if v, ok := subnetIpMap["subnet_vip"].(string); ok && v != "" {
				subnetIpInfo.SubnetVip = helper.String(v)
			}

			request.SubnetIp = append(request.SubnetIp, &subnetIpInfo)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsV20201028Client().CreateInboundEndpointWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create private dns inbound endpoint failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create private dns inbound endpoint failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.EndpointId == nil {
		return fmt.Errorf("EndpointId is nil.")
	}

	endpointId = *response.Response.EndpointId
	d.SetId(endpointId)
	return resourceTencentCloudPrivateDnsInboundEndpointRead(d, meta)
}

func resourceTencentCloudPrivateDnsInboundEndpointRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_inbound_endpoint.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = PrivatednsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		endpointId = d.Id()
	)

	respData, err := service.DescribePrivateDnsInboundEndpointById(ctx, endpointId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_private_dns_inbound_endpoint` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.EndPointName != nil {
		_ = d.Set("endpoint_name", respData.EndPointName)
	}

	if respData.UniqVpcId != nil {
		_ = d.Set("endpoint_vpc", respData.UniqVpcId)
	}

	if respData.EndPointService != nil {
		tmpList := make([]map[string]interface{}, 0, len(respData.EndPointService))
		for _, item := range respData.EndPointService {
			dMap := make(map[string]interface{}, 0)
			if item.UniqSubnetId != nil {
				dMap["subnet_id"] = item.UniqSubnetId
			}

			if item.EndPointVip != nil {
				dMap["subnet_vip"] = item.EndPointVip
			}

			tmpList = append(tmpList, dMap)
		}

		_ = d.Set("subnet_ip", tmpList)
	}

	return nil
}

func resourceTencentCloudPrivateDnsInboundEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_inbound_endpoint.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		endpointId = d.Id()
	)

	if d.HasChange("endpoint_name") {
		request := privatednsv20201028.NewModifyInboundEndpointRequest()
		if v, ok := d.GetOk("endpoint_name"); ok {
			request.EndpointName = helper.String(v.(string))
		}

		request.EndpointId = &endpointId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsV20201028Client().ModifyInboundEndpointWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update private dns inbound endpoint failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudPrivateDnsInboundEndpointRead(d, meta)
}

func resourceTencentCloudPrivateDnsInboundEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_private_dns_inbound_endpoint.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = privatednsv20201028.NewDeleteInboundEndpointRequest()
		endpointId = d.Id()
	)

	request.EndpointId = &endpointId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePrivatednsV20201028Client().DeleteInboundEndpointWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete private dns inbound endpoint failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
