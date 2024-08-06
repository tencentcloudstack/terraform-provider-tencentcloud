package ccn

import (
	"context"
	"fmt"
	"log"
	"strings"

	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func ResourceTencentCloudCcnRouteTableInputPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCcnRouteTableInputPoliciesCreate,
		Read:   resourceTencentCloudCcnRouteTableInputPoliciesRead,
		Update: resourceTencentCloudCcnRouteTableInputPoliciesUpdate,
		Delete: resourceTencentCloudCcnRouteTableInputPoliciesDelete,
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
			"route_table_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CCN Route table ID.",
			},
			"policies": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Routing reception strategy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Routing behavior, `accept` allows, `drop` rejects.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Policy description.",
						},
						"route_conditions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Routing conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "condition type.",
									},
									"values": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "List of conditional values.",
									},
									"match_pattern": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Matching mode, `1` precise matching, `0` fuzzy matching.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCcnRouteTableInputPoliciesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_input_policies.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = vpc.NewReplaceCcnRouteTableInputPolicysRequest()
		ccnId        string
		routeTableId string
	)

	if v, ok := d.GetOk("ccn_id"); ok {
		request.CcnId = helper.String(v.(string))
		ccnId = v.(string)
	}

	if v, ok := d.GetOk("route_table_id"); ok {
		request.RouteTableId = helper.String(v.(string))
		routeTableId = v.(string)
	}

	if v, ok := d.GetOk("policies"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			ccnRouteTableInputPolicy := vpc.CcnRouteTableInputPolicy{}
			if v, ok := dMap["route_conditions"]; ok {
				for _, item := range v.([]interface{}) {
					routeConditionsMap := item.(map[string]interface{})
					ccnRouteBroadcastPolicyRouteCondition := vpc.CcnRouteBroadcastPolicyRouteCondition{}
					if v, ok := routeConditionsMap["name"]; ok {
						ccnRouteBroadcastPolicyRouteCondition.Name = helper.String(v.(string))
					}

					if v, ok := routeConditionsMap["values"]; ok {
						valuesSet := v.([]interface{})
						for i := range valuesSet {
							if valuesSet[i] != nil {
								values := valuesSet[i].(string)
								ccnRouteBroadcastPolicyRouteCondition.Values = append(ccnRouteBroadcastPolicyRouteCondition.Values, &values)
							}
						}
					}

					if v, ok := routeConditionsMap["match_pattern"]; ok {
						ccnRouteBroadcastPolicyRouteCondition.MatchPattern = helper.IntUint64(v.(int))
					}

					ccnRouteTableInputPolicy.RouteConditions = append(ccnRouteTableInputPolicy.RouteConditions, &ccnRouteBroadcastPolicyRouteCondition)
				}
			}

			if v, ok := dMap["action"]; ok {
				ccnRouteTableInputPolicy.Action = helper.String(v.(string))
			}

			if v, ok := dMap["description"]; ok {
				ccnRouteTableInputPolicy.Description = helper.String(v.(string))
			}

			request.Policys = append(request.Policys, &ccnRouteTableInputPolicy)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReplaceCcnRouteTableInputPolicys(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create vpc ReplaceCcnRouteTableInputPolicys failed")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc ReplaceCcnRouteTableInputPolicys failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{ccnId, routeTableId}, tccommon.FILED_SP))

	return resourceTencentCloudCcnRouteTableInputPoliciesRead(d, meta)
}

func resourceTencentCloudCcnRouteTableInputPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_input_policies.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	ccnId := items[0]
	routeTableId := items[1]

	policySet, err := service.DescribeVpcReplaceCcnRouteTableInputPolicysById(ctx, ccnId, routeTableId)
	if err != nil {
		return err
	}

	if policySet == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcCcnRouteTable` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("ccn_id", ccnId)
	_ = d.Set("route_table_id", routeTableId)

	if policySet.Policys != nil {
		policysList := []interface{}{}
		for _, policys := range policySet.Policys {
			policysMap := map[string]interface{}{}

			if policys.RouteConditions != nil {
				routeConditionsList := []interface{}{}
				for _, routeConditions := range policys.RouteConditions {
					routeConditionsMap := map[string]interface{}{}

					if routeConditions.Name != nil {
						routeConditionsMap["name"] = routeConditions.Name
					}

					if routeConditions.Values != nil {
						routeConditionsMap["values"] = routeConditions.Values
					}

					if routeConditions.MatchPattern != nil {
						routeConditionsMap["match_pattern"] = routeConditions.MatchPattern
					}

					routeConditionsList = append(routeConditionsList, routeConditionsMap)
				}

				policysMap["route_conditions"] = routeConditionsList
			}

			if policys.Action != nil {
				policysMap["action"] = policys.Action
			}

			if policys.Description != nil {
				policysMap["description"] = policys.Description
			}

			policysList = append(policysList, policysMap)
		}

		_ = d.Set("policies", policysList)
	}

	return nil
}

func resourceTencentCloudCcnRouteTableInputPoliciesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_input_policies.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vpc.NewReplaceCcnRouteTableInputPolicysRequest()
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	ccnId := items[0]
	routeTableId := items[1]

	if d.HasChange("policies") {
		request.CcnId = &ccnId
		request.RouteTableId = &routeTableId

		if v, ok := d.GetOk("policies"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				ccnRouteTableInputPolicy := vpc.CcnRouteTableInputPolicy{}
				if v, ok := dMap["route_conditions"]; ok {
					for _, item := range v.([]interface{}) {
						routeConditionsMap := item.(map[string]interface{})
						ccnRouteBroadcastPolicyRouteCondition := vpc.CcnRouteBroadcastPolicyRouteCondition{}
						if v, ok := routeConditionsMap["name"]; ok {
							ccnRouteBroadcastPolicyRouteCondition.Name = helper.String(v.(string))
						}

						if v, ok := routeConditionsMap["values"]; ok {
							valuesSet := v.([]interface{})
							for i := range valuesSet {
								if valuesSet[i] != nil {
									values := valuesSet[i].(string)
									ccnRouteBroadcastPolicyRouteCondition.Values = append(ccnRouteBroadcastPolicyRouteCondition.Values, &values)
								}
							}
						}

						if v, ok := routeConditionsMap["match_pattern"]; ok {
							ccnRouteBroadcastPolicyRouteCondition.MatchPattern = helper.IntUint64(v.(int))
						}

						ccnRouteTableInputPolicy.RouteConditions = append(ccnRouteTableInputPolicy.RouteConditions, &ccnRouteBroadcastPolicyRouteCondition)
					}
				}

				if v, ok := dMap["action"]; ok {
					ccnRouteTableInputPolicy.Action = helper.String(v.(string))
				}

				if v, ok := dMap["description"]; ok {
					ccnRouteTableInputPolicy.Description = helper.String(v.(string))
				}

				request.Policys = append(request.Policys, &ccnRouteTableInputPolicy)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReplaceCcnRouteTableInputPolicys(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("update vpc ReplaceCcnRouteTableInputPolicys failed")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update vpc ReplaceCcnRouteTableInputPolicys failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCcnRouteTableInputPoliciesRead(d, meta)
}

func resourceTencentCloudCcnRouteTableInputPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ccn_route_table_input_policies.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	ccnId := items[0]
	routeTableId := items[1]

	body := map[string]interface{}{
		"CcnId":        ccnId,
		"RouteTableId": routeTableId,
		"Policys":      []interface{}{},
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("vpc")
	request := tchttp.NewCommonRequest("vpc", "2017-03-12", "ReplaceCcnRouteTableInputPolicys")
	err := request.SetActionParameters(body)
	if err != nil {
		return err
	}

	response := tchttp.NewCommonResponse()
	err = client.Send(request, response)
	if err != nil {
		fmt.Printf("delete vpc ReplaceCcnRouteTableInputPolicys failed: %v \n", err)
		return err
	}

	return nil
}
