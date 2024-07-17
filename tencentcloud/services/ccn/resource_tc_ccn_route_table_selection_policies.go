package ccn

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCcnRouteTableSelectionPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnRouteTableSelectionPoliciesCreate,
		Read:   resourceTencentCloudCcnRouteTableSelectionPoliciesRead,
		Update: resourceTencentCloudCcnRouteTableSelectionPoliciesUpdate,
		Delete: resourceTencentCloudCcnRouteTableSelectionPoliciesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CCN Instance ID.",
			},
			"selection_policies": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Select strategy information set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance Type: Private Network: VPC, Dedicated Gateway: DIRECTCONNECT, Blackstone Private Network: BMVPC, EDGE Device: EDGE, EDGE Tunnel: EDGE_TUNNEL, EDGE Gateway: EDGE_VPNGW, VPN Gateway: VPNGW.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance ID.",
						},
						"source_cidr_block": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source CIDR.",
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "route table ID.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "description.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCcnRouteTableSelectionPoliciesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_selection_policies.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vpc.NewModifyRouteTableSelectionPoliciesRequest()
		ccnId   string
	)

	if v, ok := d.GetOk("ccn_id"); ok {
		request.CcnId = helper.String(v.(string))
		ccnId = v.(string)
	}

	if v, ok := d.GetOk("selection_policies"); ok {
		for _, item := range v.([]interface{}) {
			selectionPolicy := item.(map[string]interface{})
			selectionPolicyMap := vpc.CcnRouteTableSelectPolicy{}
			if v, ok := selectionPolicy["instance_type"]; ok {
				selectionPolicyMap.InstanceType = helper.String(v.(string))
			}

			if v, ok := selectionPolicy["instance_id"]; ok {
				selectionPolicyMap.InstanceId = helper.String(v.(string))
			}

			if v, ok := selectionPolicy["source_cidr_block"]; ok {
				selectionPolicyMap.SourceCidrBlock = helper.String(v.(string))
			}

			if v, ok := selectionPolicy["route_table_id"]; ok {
				selectionPolicyMap.RouteTableId = helper.String(v.(string))
			}

			if v, ok := selectionPolicy["description"]; ok {
				selectionPolicyMap.Description = helper.String(v.(string))
			}

			request.SelectionPolicies = append(request.SelectionPolicies, &selectionPolicyMap)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyRouteTableSelectionPolicies(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create vpc ModifyRouteTableSelectionPolicies failed")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc ModifyRouteTableSelectionPolicies failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(ccnId)

	return resourceTencentCloudCcnRouteTableSelectionPoliciesRead(d, meta)
}

func resourceTencentCloudCcnRouteTableSelectionPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_selection_policies.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ccnId   = d.Id()
	)

	routeSelectionPolicySet, err := service.DescribeVpcReplaceCcnRouteTableSelectionPolicysById(ctx, ccnId)
	if err != nil {
		return err
	}

	if routeSelectionPolicySet == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RouteTableSelectionPolicies` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("ccn_id", ccnId)
	policyList := []interface{}{}
	for _, policy := range routeSelectionPolicySet {
		policyMap := map[string]interface{}{}
		if policy.InstanceType != nil {
			policyMap["instance_type"] = policy.InstanceType
		}

		if policy.InstanceId != nil {
			policyMap["instance_id"] = policy.InstanceId
		}

		if policy.SourceCidrBlock != nil {
			policyMap["source_cidr_block"] = policy.SourceCidrBlock
		}

		if policy.RouteTableId != nil {
			policyMap["route_table_id"] = policy.RouteTableId
		}

		if policy.Description != nil {
			policyMap["description"] = policy.Description
		}

		policyList = append(policyList, policyMap)
	}

	_ = d.Set("selection_policies", policyList)

	return nil
}

func resourceTencentCloudCcnRouteTableSelectionPoliciesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_selection_policies.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vpc.NewModifyRouteTableSelectionPoliciesRequest()
		ccnId   = d.Id()
	)

	if d.HasChange("selection_policies") {
		request.CcnId = &ccnId
		if v, ok := d.GetOk("selection_policies"); ok {
			for _, item := range v.([]interface{}) {
				selectionPolicy := item.(map[string]interface{})
				selectionPolicyMap := vpc.CcnRouteTableSelectPolicy{}
				if v, ok := selectionPolicy["instance_type"]; ok {
					selectionPolicyMap.InstanceType = helper.String(v.(string))
				}

				if v, ok := selectionPolicy["instance_id"]; ok {
					selectionPolicyMap.InstanceId = helper.String(v.(string))
				}

				if v, ok := selectionPolicy["source_cidr_block"]; ok {
					selectionPolicyMap.SourceCidrBlock = helper.String(v.(string))
				}

				if v, ok := selectionPolicy["route_table_id"]; ok {
					selectionPolicyMap.RouteTableId = helper.String(v.(string))
				}

				if v, ok := selectionPolicy["description"]; ok {
					selectionPolicyMap.Description = helper.String(v.(string))
				}

				request.SelectionPolicies = append(request.SelectionPolicies, &selectionPolicyMap)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyRouteTableSelectionPolicies(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("update vpc ModifyRouteTableSelectionPolicies failed")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update vpc ModifyRouteTableSelectionPolicies failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCcnRouteTableSelectionPoliciesRead(d, meta)
}

func resourceTencentCloudCcnRouteTableSelectionPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_selection_policies.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vpc.NewClearRouteTableSelectionPoliciesRequest()
		ccnId   = d.Id()
	)

	request.CcnId = &ccnId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ClearRouteTableSelectionPolicies(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete vpc ClearRouteTableSelectionPolicies failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
