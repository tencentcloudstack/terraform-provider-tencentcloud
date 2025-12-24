package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcRoutePolicyAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcRoutePolicyAssociationCreate,
		Read:   resourceTencentCloudVpcRoutePolicyAssociationRead,
		Update: resourceTencentCloudVpcRoutePolicyAssociationUpdate,
		Delete: resourceTencentCloudVpcRoutePolicyAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique route table ID.",
			},
			"route_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the unique ID of the route reception policy.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority.",
			},
		},
	}
}

func resourceTencentCloudVpcRoutePolicyAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_association.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = vpcv20170312.NewCreateRoutePolicyAssociationsRequest()
		routePolicyId string
		routeTableId  string
	)

	routePolicyAssociation := vpcv20170312.RoutePolicyAssociation{}
	if v, ok := d.GetOk("route_table_id"); ok {
		routePolicyAssociation.RouteTableId = helper.String(v.(string))
		routeTableId = v.(string)
	}

	if v, ok := d.GetOk("route_policy_id"); ok {
		routePolicyAssociation.RoutePolicyId = helper.String(v.(string))
		routePolicyId = v.(string)
	}

	if v, ok := d.GetOkExists("priority"); ok {
		routePolicyAssociation.Priority = helper.IntUint64(v.(int))
	}

	request.RoutePolicyAssociationSet = append(request.RoutePolicyAssociationSet, &routePolicyAssociation)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateRoutePolicyAssociationsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vpc route policy association failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{routePolicyId, routeTableId}, tccommon.FILED_SP))
	return resourceTencentCloudVpcRoutePolicyAssociationRead(d, meta)
}

func resourceTencentCloudVpcRoutePolicyAssociationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_association.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	routePolicyId := idSplit[0]
	routeTableId := idSplit[1]

	respData, err := service.DescribeVpcRoutePolicyAssociationById(ctx, routePolicyId, routeTableId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_route_policy_association` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.RouteTableId != nil {
		_ = d.Set("route_table_id", respData.RouteTableId)
	}

	if respData.RoutePolicyId != nil {
		_ = d.Set("route_policy_id", respData.RoutePolicyId)
	}

	if respData.Priority != nil {
		_ = d.Set("priority", respData.Priority)
	}

	return nil
}

func resourceTencentCloudVpcRoutePolicyAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_association.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = vpcv20170312.NewResetRoutePolicyAssociationsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	routePolicyId := idSplit[0]
	routeTableId := idSplit[1]

	routePolicyAssociation := vpcv20170312.RoutePolicyAssociation{}
	routePolicyAssociation.RouteTableId = &routeTableId
	routePolicyAssociation.RoutePolicyId = &routePolicyId

	if v, ok := d.GetOkExists("priority"); ok {
		routePolicyAssociation.Priority = helper.IntUint64(v.(int))
	}

	request.RouteTableId = &routeTableId
	request.RoutePolicyAssociationSet = append(request.RoutePolicyAssociationSet, &routePolicyAssociation)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ResetRoutePolicyAssociationsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update vpc route policy association failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudVpcRoutePolicyAssociationRead(d, meta)
}

func resourceTencentCloudVpcRoutePolicyAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_association.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = vpcv20170312.NewDeleteRoutePolicyAssociationsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	routePolicyId := idSplit[0]
	routeTableId := idSplit[1]

	routePolicyAssociation := vpcv20170312.RoutePolicyAssociation{}
	routePolicyAssociation.RouteTableId = &routeTableId
	routePolicyAssociation.RoutePolicyId = &routePolicyId
	request.RoutePolicyAssociationSet = append(request.RoutePolicyAssociationSet, &routePolicyAssociation)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteRoutePolicyAssociationsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vpc route policy association failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
