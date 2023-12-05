package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcdbDbParameters() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbDbParametersCreate,
		Read:   resourceTencentCloudDcdbDbParametersRead,
		Update: resourceTencentCloudDcdbDbParametersUpdate,
		Delete: resourceTencentCloudDcdbDbParametersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"params": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Parameter list, each element is a combination of Param and Value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of parameter.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of parameter.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDcdbDbParametersCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_parameters.create")()
	defer inconsistentCheck(d, meta)()

	return resourceTencentCloudDcdbDbParametersUpdate(d, meta)
}

func resourceTencentCloudDcdbDbParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_parameters.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	paramName := idSplit[1]

	dbParameters, err := service.DescribeDcdbDbParametersById(ctx, instanceId)
	if err != nil {
		return err
	}
	if dbParameters == nil {
		d.SetId("")
		return fmt.Errorf("resource `dbParameters` %s does not exist", d.Id())
	}

	if dbParameters.InstanceId != nil {
		_ = d.Set("instance_id", dbParameters.InstanceId)
	}

	paramsList := []interface{}{}
	for _, param := range dbParameters.Params {
		paramsMap := map[string]interface{}{}

		if param.Param != nil && *param.Param == paramName {
			paramsMap["param"] = param.Param
			if param.Value != nil {
				paramsMap["value"] = param.Value
			}
			paramsList = append(paramsList, paramsMap)
			break
		}
	}
	_ = d.Set("params", paramsList)

	return nil
}

func resourceTencentCloudDcdbDbParametersUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_parameters.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		instanceId string
		paramName  string
		request    = dcdb.NewModifyDBParametersRequest()
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if d.HasChange("params") {
		if dMap, ok := helper.InterfacesHeadMap(d, "params"); ok {
			dBParamValue := dcdb.DBParamValue{}
			if v, ok := dMap["param"]; ok {
				dBParamValue.Param = helper.String(v.(string))
				paramName = v.(string)
			}
			if v, ok := dMap["value"]; ok {
				dBParamValue.Value = helper.String(v.(string))
			}
			request.Params = append(request.Params, &dBParamValue)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBParameters(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb dbParameters failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
		dbParameters, err := service.DescribeDcdbDbParametersById(ctx, instanceId)
		if err != nil {
			return retryError(err)
		}

		if dbParameters == nil || dbParameters.Params == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeDcdbDbParametersById return result(dcdb db parameter) is nil!"))
		}

		for _, param := range dbParameters.Params {
			if param.Param != nil && *param.Param == paramName {
				if *param.Value == *param.SetValue {
					return nil
				} else {
					return resource.RetryableError(fmt.Errorf("db parameter initializing, retry..."))
				}
			}
		}
		return resource.NonRetryableError(fmt.Errorf("DescribeDcdbDbParametersById not found the target param:[%s], exit...", paramName))
	})
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{instanceId, paramName}, FILED_SP))

	return resourceTencentCloudDcdbDbParametersRead(d, meta)
}

func resourceTencentCloudDcdbDbParametersDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_parameters.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
