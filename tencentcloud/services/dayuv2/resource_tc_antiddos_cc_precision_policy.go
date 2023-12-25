package dayuv2

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosCcPrecisionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosCcPrecisionPolicyCreate,
		Read:   resourceTencentCloudAntiddosCcPrecisionPolicyRead,
		Update: resourceTencentCloudAntiddosCcPrecisionPolicyUpdate,
		Delete: resourceTencentCloudAntiddosCcPrecisionPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Id.",
			},

			"ip": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Ip value.",
			},

			"protocol": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "protocol http or https.",
			},

			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "domain.",
			},

			"policy_action": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "policy type, alg or drop.",
			},

			"policy_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "policy list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "field type.",
						},
						"field_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration fields can take values of cgi, ua, cookie, referer, accept, srcip.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "value.",
						},
						"value_operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration item value comparison method, can take values of equal, not_ Equal, include.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAntiddosCcPrecisionPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_cc_precision_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request      = antiddos.NewCreateCCPrecisionPolicyRequest()
		instanceId   string
		ip           string
		domain       string
		protocol     string
		policyAction string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("ip"); ok {
		ip = v.(string)
		request.Ip = helper.String(ip)
	}

	if v, ok := d.GetOk("protocol"); ok {
		protocol = v.(string)
		request.Protocol = helper.String(protocol)
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(domain)
	}

	if v, ok := d.GetOk("policy_action"); ok {
		policyAction = v.(string)
		request.PolicyAction = helper.String(policyAction)
	}

	policyStrings := make([]string, 0)

	if v, ok := d.GetOk("policy_list"); ok {
		for _, item := range v.([]interface{}) {
			policyItems := make([]string, 0)
			dMap := item.(map[string]interface{})
			cCPrecisionPlyRecord := antiddos.CCPrecisionPlyRecord{}
			if v, ok := dMap["field_type"]; ok {
				cCPrecisionPlyRecord.FieldType = helper.String(v.(string))
				policyItems = append(policyItems, v.(string))
			}
			if v, ok := dMap["field_name"]; ok {
				cCPrecisionPlyRecord.FieldName = helper.String(v.(string))
				policyItems = append(policyItems, v.(string))
			}
			if v, ok := dMap["value"]; ok {
				cCPrecisionPlyRecord.Value = helper.String(v.(string))
				policyItems = append(policyItems, v.(string))
			}
			if v, ok := dMap["value_operator"]; ok {
				cCPrecisionPlyRecord.ValueOperator = helper.String(v.(string))
				policyItems = append(policyItems, v.(string))
			}
			request.PolicyList = append(request.PolicyList, &cCPrecisionPlyRecord)
			policyStrings = append(policyStrings, strings.Join(policyItems, tccommon.FILED_SP))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateCCPrecisionPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos ccPrecisionPolicy failed, reason:%+v", logId, err)
		return err
	}

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	ccPrecisionPolicys, err := service.DescribeAntiddosCcPrecisionPolicyById(ctx, "bgpip", instanceId, ip, domain, protocol)
	if err != nil {
		return err
	}
	var ccPrecisionPolicy *antiddos.CCPrecisionPolicy
	for _, item := range ccPrecisionPolicys {
		if *item.PolicyAction != policyAction {
			continue
		}
		tmpPolicyStrings := make([]string, 0)
		for _, policy := range item.PolicyList {
			policyItems := make([]string, 0)
			policyItems = append(policyItems, *policy.FieldType)
			policyItems = append(policyItems, *policy.FieldName)
			policyItems = append(policyItems, *policy.Value)
			policyItems = append(policyItems, *policy.ValueOperator)
			tmpPolicyStrings = append(tmpPolicyStrings, strings.Join(policyItems, tccommon.FILED_SP))
		}
		sort.Strings(policyStrings)
		sort.Strings(tmpPolicyStrings)
		if !reflect.DeepEqual(policyStrings, tmpPolicyStrings) {
			continue
		}
		ccPrecisionPolicy = item
	}
	if ccPrecisionPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosCcPrecisionPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	d.SetId(strings.Join([]string{instanceId, *ccPrecisionPolicy.PolicyId, ip, domain, protocol}, tccommon.FILED_SP))

	return resourceTencentCloudAntiddosCcPrecisionPolicyRead(d, meta)
}

func resourceTencentCloudAntiddosCcPrecisionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_cc_precision_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	policyId := idSplit[1]
	ip := idSplit[2]
	domain := idSplit[3]
	protocol := idSplit[4]

	ccPrecisionPolicys, err := service.DescribeAntiddosCcPrecisionPolicyById(ctx, "bgpip", instanceId, ip, domain, protocol)
	if err != nil {
		return err
	}

	var ccPrecisionPolicy *antiddos.CCPrecisionPolicy
	for _, item := range ccPrecisionPolicys {
		if *item.PolicyId == policyId {
			ccPrecisionPolicy = item
			break
		}
	}
	if ccPrecisionPolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosCcPrecisionPolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if ccPrecisionPolicy.InstanceId != nil {
		_ = d.Set("instance_id", ccPrecisionPolicy.InstanceId)
	}

	if ccPrecisionPolicy.Ip != nil {
		_ = d.Set("ip", ccPrecisionPolicy.Ip)
	}

	if ccPrecisionPolicy.Protocol != nil {
		_ = d.Set("protocol", ccPrecisionPolicy.Protocol)
	}

	if ccPrecisionPolicy.Domain != nil {
		_ = d.Set("domain", ccPrecisionPolicy.Domain)
	}

	if ccPrecisionPolicy.PolicyAction != nil {
		_ = d.Set("policy_action", ccPrecisionPolicy.PolicyAction)
	}

	if ccPrecisionPolicy.PolicyList != nil {
		policyListList := []interface{}{}
		for _, policy := range ccPrecisionPolicy.PolicyList {
			policyListMap := map[string]interface{}{}

			if policy.FieldType != nil {
				policyListMap["field_type"] = policy.FieldType
			}

			if policy.FieldName != nil {
				policyListMap["field_name"] = policy.FieldName
			}

			if policy.Value != nil {
				policyListMap["value"] = policy.Value
			}

			if policy.ValueOperator != nil {
				policyListMap["value_operator"] = policy.ValueOperator
			}

			policyListList = append(policyListList, policyListMap)
		}

		_ = d.Set("policy_list", policyListList)

	}

	return nil
}

func resourceTencentCloudAntiddosCcPrecisionPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_cc_precision_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := antiddos.NewModifyCCPrecisionPolicyRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	policyId := idSplit[1]

	request.InstanceId = &instanceId
	request.PolicyId = &policyId

	immutableArgs := []string{"instance_id", "ip", "protocol", "domain"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	if d.HasChange("policy_action") || d.HasChange("policy_list") {
		if v, ok := d.GetOk("policy_action"); ok {
			request.PolicyAction = helper.String(v.(string))
		}

		if v, ok := d.GetOk("policy_list"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				cCPrecisionPlyRecord := antiddos.CCPrecisionPlyRecord{}
				if v, ok := dMap["field_type"]; ok {
					cCPrecisionPlyRecord.FieldType = helper.String(v.(string))
				}
				if v, ok := dMap["field_name"]; ok {
					cCPrecisionPlyRecord.FieldName = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					cCPrecisionPlyRecord.Value = helper.String(v.(string))
				}
				if v, ok := dMap["value_operator"]; ok {
					cCPrecisionPlyRecord.ValueOperator = helper.String(v.(string))
				}
				request.PolicyList = append(request.PolicyList, &cCPrecisionPlyRecord)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().ModifyCCPrecisionPolicy(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update antiddos ccPrecisionPolicy failed, reason:%+v", logId, err)
			return err
		}
	}
	return resourceTencentCloudAntiddosCcPrecisionPolicyRead(d, meta)
}

func resourceTencentCloudAntiddosCcPrecisionPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_cc_precision_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	policyId := idSplit[1]

	if err := service.DeleteAntiddosCcPrecisionPolicyById(ctx, instanceId, policyId); err != nil {
		return err
	}

	return nil
}
