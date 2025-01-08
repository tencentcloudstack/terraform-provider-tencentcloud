package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMongodbInstanceParams() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceParamsCreate,
		Read:   resourceTencentCloudMongodbInstanceParamsRead,
		Update: resourceTencentCloudMongodbInstanceParamsUpdate,
		Delete: resourceTencentCloudMongodbInstanceParamsDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance id.",
			},

			"instance_params": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Specify the parameter name and value to be modified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter names that need to be modified.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value corresponding to the parameter name to be modified.",
						},
					},
				},
			},

			"modify_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Operation types, including:\n" +
					"	- IMMEDIATELY: Adjust immediately;\n" +
					"	- DELAY: Delay adjustment;\n" +
					"Optional field. If this parameter is not configured, it defaults to immediate adjustment.",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceParamsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_params.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMongodbInstanceParamsUpdate(d, meta)
}

func resourceTencentCloudMongodbInstanceParamsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_params.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	_ = d.Set("instance_id", instanceId)

	instanceParams := d.Get("instance_params").(*schema.Set).List()
	paramNames := make([]string, 0)
	for _, instanceParam := range instanceParams {
		instanceParamMap := instanceParam.(map[string]interface{})
		paramNames = append(paramNames, instanceParamMap["key"].(string))
	}
	respData, err := service.DescribeMongodbInstanceParamValues(ctx, instanceId, paramNames)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `mongodb_instance_params` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	resInstanceParams := make([]interface{}, 0)
	for k, v := range respData {
		resInstanceParams = append(resInstanceParams, map[string]interface{}{
			"key":   k,
			"value": v,
		})
	}
	_ = d.Set("instance_params", resInstanceParams)

	return nil
}

func resourceTencentCloudMongodbInstanceParamsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_params.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	instanceId := d.Id()

	needChange := false
	mutableArgs := []string{"instance_params"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := mongodb.NewModifyInstanceParamsRequest()
		request.InstanceId = helper.String(d.Get("instance_id").(string))

		if v, ok := d.GetOk("instance_params"); ok {
			for _, item := range v.(*schema.Set).List() {
				instanceParamsMap := item.(map[string]interface{})
				modifyMongoDBParamType := mongodb.ModifyMongoDBParamType{}
				if v, ok := instanceParamsMap["key"]; ok {
					modifyMongoDBParamType.Key = helper.String(v.(string))
				}
				if v, ok := instanceParamsMap["value"]; ok {
					modifyMongoDBParamType.Value = helper.String(v.(string))
				}
				request.InstanceParams = append(request.InstanceParams, &modifyMongoDBParamType)
			}
		}

		if v, ok := d.GetOk("modify_type"); ok {
			request.ModifyType = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().ModifyInstanceParamsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mongodb instance params failed, reason:%+v", logId, err)
			return err
		}

		service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceParams := d.Get("instance_params").(*schema.Set).List()
		paramMap := make(map[string]string)
		paramNames := make([]string, 0)

		for _, instanceParam := range instanceParams {
			instanceParamMap := instanceParam.(map[string]interface{})
			key := instanceParamMap["key"].(string)
			value := instanceParamMap["value"].(string)
			paramMap[key] = value
			paramNames = append(paramNames, key)
		}

		err = resource.Retry(6*tccommon.WriteRetryTimeout, func() *resource.RetryError {
			respMap, e := service.DescribeMongodbInstanceParamValues(ctx, instanceId, paramNames)
			if e != nil {
				return tccommon.RetryError(e)
			}

			for k, v := range paramMap {
				if vv, ok := respMap[k]; ok {
					if v != vv {
						return resource.RetryableError(fmt.Errorf("param %s is still being updated", k))
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mongodb instance params failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = instanceId
	return resourceTencentCloudMongodbInstanceParamsRead(d, meta)
}

func resourceTencentCloudMongodbInstanceParamsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_params.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
