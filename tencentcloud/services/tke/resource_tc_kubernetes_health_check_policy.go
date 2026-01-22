package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20220501 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20220501"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesHealthCheckPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesHealthCheckPolicyCreate,
		Read:   resourceTencentCloudKubernetesHealthCheckPolicyRead,
		Update: resourceTencentCloudKubernetesHealthCheckPolicyUpdate,
		Delete: resourceTencentCloudKubernetesHealthCheckPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the cluster.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Health Check Policy Name.",
			},

			"rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Health check policy rule list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_repair_enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable repair or not.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable detection of this project or not.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Health check rule details.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudKubernetesHealthCheckPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_health_check_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = tkev20220501.NewCreateHealthCheckPolicyRequest()
		response  = tkev20220501.NewCreateHealthCheckPolicyResponse()
		clusterId string
		name      string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	healthCheckPolicy := tkev20220501.HealthCheckPolicy{}
	if v, ok := d.GetOk("name"); ok {
		healthCheckPolicy.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			rulesMap := item.(map[string]interface{})
			healthCheckPolicyRule := tkev20220501.HealthCheckPolicyRule{}
			if v, ok := rulesMap["auto_repair_enabled"]; ok {
				healthCheckPolicyRule.AutoRepairEnabled = helper.Bool(v.(bool))
			}

			if v, ok := rulesMap["enabled"]; ok {
				healthCheckPolicyRule.Enabled = helper.Bool(v.(bool))
			}

			if v, ok := rulesMap["name"]; ok {
				healthCheckPolicyRule.Name = helper.String(v.(string))
			}

			healthCheckPolicy.Rules = append(healthCheckPolicy.Rules, &healthCheckPolicyRule)
		}
	}

	request.HealthCheckPolicy = &healthCheckPolicy
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20220501Client().CreateHealthCheckPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create kubernetes health check policy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes health check policy failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.HealthCheckPolicyName == nil {
		return fmt.Errorf("HealthCheckPolicyName is nil.")
	}

	name = *response.Response.HealthCheckPolicyName
	d.SetId(strings.Join([]string{clusterId, name}, tccommon.FILED_SP))
	return resourceTencentCloudKubernetesHealthCheckPolicyRead(d, meta)
}

func resourceTencentCloudKubernetesHealthCheckPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_health_check_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	name := idSplit[1]

	respData, err := service.DescribeKubernetesHealthCheckPolicyById(ctx, clusterId, name)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_health_check_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Rules != nil {
		rulesList := make([]map[string]interface{}, 0, len(respData.Rules))
		for _, rules := range respData.Rules {
			rulesMap := map[string]interface{}{}
			if rules.AutoRepairEnabled != nil {
				rulesMap["auto_repair_enabled"] = rules.AutoRepairEnabled
			}

			if rules.Enabled != nil {
				rulesMap["enabled"] = rules.Enabled
			}

			if rules.Name != nil {
				rulesMap["name"] = rules.Name
			}

			rulesList = append(rulesList, rulesMap)
		}

		_ = d.Set("rules", rulesList)
	}

	return nil
}

func resourceTencentCloudKubernetesHealthCheckPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_health_check_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	name := idSplit[1]

	if d.HasChange("rules") {
		request := tkev20220501.NewModifyHealthCheckPolicyRequest()
		healthCheckPolicy := tkev20220501.HealthCheckPolicy{}
		healthCheckPolicy.Name = helper.String(name)
		if v, ok := d.GetOk("rules"); ok {
			for _, item := range v.([]interface{}) {
				rulesMap := item.(map[string]interface{})
				healthCheckPolicyRule := tkev20220501.HealthCheckPolicyRule{}
				if v, ok := rulesMap["auto_repair_enabled"]; ok {
					healthCheckPolicyRule.AutoRepairEnabled = helper.Bool(v.(bool))
				}

				if v, ok := rulesMap["enabled"]; ok {
					healthCheckPolicyRule.Enabled = helper.Bool(v.(bool))
				}

				if v, ok := rulesMap["name"]; ok {
					healthCheckPolicyRule.Name = helper.String(v.(string))
				}

				healthCheckPolicy.Rules = append(healthCheckPolicy.Rules, &healthCheckPolicyRule)
			}
		}

		request.HealthCheckPolicy = &healthCheckPolicy
		request.ClusterId = helper.String(clusterId)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20220501Client().ModifyHealthCheckPolicyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes health check policy failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudKubernetesHealthCheckPolicyRead(d, meta)
}

func resourceTencentCloudKubernetesHealthCheckPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_health_check_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tkev20220501.NewDeleteHealthCheckPolicyRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	name := idSplit[1]

	request.ClusterId = helper.String(clusterId)
	request.HealthCheckPolicyName = helper.String(name)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20220501Client().DeleteHealthCheckPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes health check policy failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
