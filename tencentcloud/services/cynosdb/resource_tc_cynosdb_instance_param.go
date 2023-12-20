package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbInstanceParam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbInstanceParamCreate,
		Read:   resourceTencentCloudCynosdbInstanceParamRead,
		Update: resourceTencentCloudCynosdbInstanceParamUpdate,
		Delete: resourceTencentCloudCynosdbInstanceParamDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"instance_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"instance_param_list": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "Instance parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter Name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Current value of parameter.",
						},
					},
				},
			},

			"is_in_maintain_period": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Yes: modify within the operation and maintenance time window, no: execute immediately (default value).",
			},
		},
	}
}

func resourceTencentCloudCynosdbInstanceParamCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_instance_param.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(clusterId + tccommon.FILED_SP + instanceId)

	return resourceTencentCloudCynosdbInstanceParamUpdate(d, meta)
}

func resourceTencentCloudCynosdbInstanceParamRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_instance_param.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]

	instanceParam, err := service.DescribeCynosdbInstanceParamById(ctx, clusterId, instanceId)
	if err != nil {
		return err
	}

	if instanceParam == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbInstanceParam` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("instance_id", instanceId)

	if instanceParam.ParamsItems != nil {
		checkFlag := true
		paramItem := make(map[string]string)
		if v, ok := d.GetOk("instance_param_list"); ok {
			for _, v := range v.(*schema.Set).List() {
				dMap := v.(map[string]interface{})
				key := dMap["param_name"].(string)
				value := dMap["current_value"].(string)
				paramItem[key] = value
			}
		} else {
			checkFlag = false
		}

		paramInfoSetList := []interface{}{}
		for _, paramInfoSet := range instanceParam.ParamsItems {
			if _, ok := paramItem[*paramInfoSet.ParamName]; !ok && checkFlag {
				continue
			}
			paramInfoSetList = append(paramInfoSetList, map[string]interface{}{
				"param_name":    *paramInfoSet.ParamName,
				"current_value": *paramInfoSet.CurrentValue,
			})
		}
		_ = d.Set("instance_param_list", paramInfoSetList)
	}

	return nil
}

func resourceTencentCloudCynosdbInstanceParamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_instance_param.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	request := cynosdb.NewModifyInstanceParamRequest()
	response := cynosdb.NewModifyInstanceParamResponse()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	instanceId := idSplit[1]
	request.ClusterId = &clusterId
	request.InstanceIds = []*string{&instanceId}

	if d.HasChange("instance_param_list") {
		oldParam, _ := d.GetChange("instance_param_list")
		oldItem := oldParam.(*schema.Set).List()
		oldParamItem := make(map[string]string)
		for _, v := range oldItem {
			dMap := v.(map[string]interface{})
			key := dMap["param_name"].(string)
			value := dMap["current_value"].(string)
			oldParamItem[key] = value
		}

		if v, ok := d.GetOk("instance_param_list"); ok {
			for _, item := range v.(*schema.Set).List() {
				dMap := item.(map[string]interface{})
				paramItem := cynosdb.ModifyParamItem{}
				if v, ok := dMap["param_name"]; ok {
					paramItem.ParamName = helper.String(v.(string))
				}
				if v, ok := dMap["current_value"]; ok {
					paramItem.CurrentValue = helper.String(v.(string))
				}
				if oldParamItem[*paramItem.ParamName] != "" {
					paramItem.OldValue = helper.String(oldParamItem[*paramItem.ParamName])
				}
				request.InstanceParamList = append(request.InstanceParamList, &paramItem)
			}
		}
	}

	if v, ok := d.GetOk("is_in_maintain_period"); ok {
		request.IsInMaintainPeriod = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyInstanceParam(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb instanceParam failed, reason:%+v", logId, err)
		return err
	}

	flowId := *response.Response.FlowId
	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("update cynosdb instanceParam is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb instanceParam fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbInstanceParamRead(d, meta)
}

func resourceTencentCloudCynosdbInstanceParamDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_instance_param.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
