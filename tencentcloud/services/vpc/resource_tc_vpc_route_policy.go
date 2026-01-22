package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcRoutePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcRoutePolicyCreate,
		Read:   resourceTencentCloudVpcRoutePolicyRead,
		Update: resourceTencentCloudVpcRoutePolicyUpdate,
		Delete: resourceTencentCloudVpcRoutePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"route_policy_description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Routing policy description.",
			},

			"route_policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the routing strategy name.",
			},

			// computed
			"route_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Route policy ID.",
			},
		},
	}
}

func resourceTencentCloudVpcRoutePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = vpcv20170312.NewCreateRoutePolicyRequest()
		response      = vpcv20170312.NewCreateRoutePolicyResponse()
		routePolicyId string
	)

	if v, ok := d.GetOk("route_policy_description"); ok {
		request.RoutePolicyDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_policy_name"); ok {
		request.RoutePolicyName = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateRoutePolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RoutePolicy == nil {
			return resource.NonRetryableError(fmt.Errorf("Create vpc route policy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vpc route policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.RoutePolicy.RoutePolicyId == nil {
		return fmt.Errorf("RoutePolicyId is nil.")
	}

	routePolicyId = *response.Response.RoutePolicy.RoutePolicyId
	d.SetId(routePolicyId)
	return resourceTencentCloudVpcRoutePolicyRead(d, meta)
}

func resourceTencentCloudVpcRoutePolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		routePolicyId = d.Id()
	)

	respData, err := service.DescribeVpcRoutePolicyById(ctx, routePolicyId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_route_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.RoutePolicyDescription != nil {
		_ = d.Set("route_policy_description", respData.RoutePolicyDescription)
	}

	if respData.RoutePolicyName != nil {
		_ = d.Set("route_policy_name", respData.RoutePolicyName)
	}

	_ = d.Set("route_policy_id", routePolicyId)

	return nil
}

func resourceTencentCloudVpcRoutePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = vpcv20170312.NewModifyRoutePolicyAttributeRequest()
		routePolicyId = d.Id()
	)

	if v, ok := d.GetOk("route_policy_description"); ok {
		request.RoutePolicyDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_policy_name"); ok {
		request.RoutePolicyName = helper.String(v.(string))
	}

	request.RoutePolicyId = &routePolicyId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyRoutePolicyAttributeWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify vpc route policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudVpcRoutePolicyRead(d, meta)
}

func resourceTencentCloudVpcRoutePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = vpcv20170312.NewDeleteRoutePolicyRequest()
		routePolicyId = d.Id()
	)

	request.RoutePolicyId = &routePolicyId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteRoutePolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vpc route policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
