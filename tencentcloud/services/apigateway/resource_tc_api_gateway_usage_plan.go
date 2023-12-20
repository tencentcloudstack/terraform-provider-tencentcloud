package apigateway

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAPIGatewayUsagePlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayUsagePlanCreate,
		Read:   resourceTencentCloudAPIGatewayUsagePlanRead,
		Update: resourceTencentCloudAPIGatewayUsagePlanUpdate,
		Delete: resourceTencentCloudAPIGatewayUsagePlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"usage_plan_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom usage plan name.",
			},
			"usage_plan_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom usage plan description.",
			},
			"max_request_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
				ValidateFunc: func(i interface{}, s string) (strings []string, errors []error) {
					if value := int64(i.(int)); value == -1 {
						return
					}
					return tccommon.ValidateIntegerInRange(1, 99999999)(i, s)
				},
				Description: "Total number of requests allowed. Valid values: -1, [1,99999999]. The default value is -1, which indicates no limit.",
			},
			"max_request_num_pre_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
				ValidateFunc: func(i interface{}, s string) (strings []string, errors []error) {
					if value := int64(i.(int)); value == -1 {
						return
					}
					return tccommon.ValidateIntegerInRange(1, 2000)(i, s)
				},
				Description: "Limit of requests per second. Valid values: -1, [1,2000]. The default value is -1, which indicates no limit.",
			},
			// Computed values.
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
			},
			"attach_api_keys": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Attach API keys list.",
			},
			"attach_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attach service and API list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service ID.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service name.",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API ID, this value is empty if attach service.",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API name, this value is empty if attach service.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API path, this value is empty if attach service.",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API method, this value is empty if attach service.",
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment name.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in the format of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAPIGatewayUsagePlanCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_usage_plan.create")()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService   = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		usagePlanName       = d.Get("usage_plan_name").(string)
		maxRequestNum       = int64(d.Get("max_request_num").(int))
		maxRequestNumPreSec = int64(d.Get("max_request_num_pre_sec").(int))
		usagePlanDesc       *string
	)

	if v, has := d.GetOk("usage_plan_desc"); has {
		usagePlanDesc = helper.String(v.(string))
	}

	usagePlanId, err := apiGatewayService.CreateUsagePlan(ctx, usagePlanName, usagePlanDesc, maxRequestNum, maxRequestNumPreSec)
	if err != nil {
		return err
	}

	d.SetId(usagePlanId)

	//wait usage plan create ok
	if outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr := apiGatewayService.DescribeUsagePlan(ctx, usagePlanId)
		if inErr != nil {
			return tccommon.RetryError(inErr, tccommon.InternalError)
		}
		if has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" usage plan  %s not found on server", usagePlanId))

	}); outErr != nil {
		return outErr
	}

	return resourceTencentCloudAPIGatewayUsagePlanRead(d, meta)
}

func resourceTencentCloudAPIGatewayUsagePlanRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_usage_plan.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		usagePlanId       = d.Id()
		info              apigateway.UsagePlanInfo
		attachList        []*apigateway.UsagePlanEnvironment
		err               error
		has               bool
	)

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, err = apiGatewayService.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	//service attach and API
	for _, bindType := range API_GATEWAY_TYPES {
		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			list, inErr := apiGatewayService.DescribeUsagePlanEnvironments(ctx, usagePlanId, bindType)
			if inErr != nil {
				return tccommon.RetryError(inErr, tccommon.InternalError)
			}
			attachList = append(attachList, list...)
			return nil
		}); err != nil {
			return err
		}
	}

	infoAttachList := make([]map[string]interface{}, 0, len(attachList))
	for _, v := range attachList {
		infoAttachList = append(infoAttachList, map[string]interface{}{
			"service_id":   v.ServiceId,
			"service_name": v.ServiceName,
			"api_id":       v.ApiId,
			"api_name":     v.ApiName,
			"path":         v.Path,
			"method":       v.Method,
			"environment":  v.Environment,
			"modify_time":  v.ModifiedTime,
			"create_time":  v.CreatedTime,
		})
	}

	_ = d.Set("usage_plan_name", info.UsagePlanName)
	_ = d.Set("usage_plan_desc", info.UsagePlanDesc)
	_ = d.Set("max_request_num", info.MaxRequestNum)
	_ = d.Set("max_request_num_pre_sec", info.MaxRequestNumPreSec)
	_ = d.Set("modify_time", info.ModifiedTime)
	_ = d.Set("create_time", info.CreatedTime)
	_ = d.Set("attach_list", infoAttachList)

	attachApiKeys := make([]string, 0, len(info.BindSecretIds))
	for _, v := range info.BindSecretIds {
		attachApiKeys = append(attachApiKeys, *v)
	}
	_ = d.Set("attach_api_keys", attachApiKeys)

	return nil
}

func resourceTencentCloudAPIGatewayUsagePlanUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_usage_plan.update")()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		usagePlanId         = d.Id()
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService   = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		usagePlanName       = d.Get("usage_plan_name").(string)
		maxRequestNum       = int64(d.Get("max_request_num").(int))
		maxRequestNumPreSec = int64(d.Get("max_request_num_pre_sec").(int))
		err                 error
		usagePlanDesc       *string
	)

	if v, has := d.GetOk("usage_plan_desc"); has {
		usagePlanDesc = helper.String(v.(string))
	}

	if d.HasChange("usage_plan_name") || d.HasChange("usage_plan_desc") ||
		d.HasChange("max_request_num") || d.HasChange("max_request_num_pre_sec") {

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err = apiGatewayService.ModifyUsagePlan(ctx,
				usagePlanId,
				usagePlanName,
				usagePlanDesc,
				maxRequestNum,
				maxRequestNumPreSec)

			if nil != err {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayUsagePlanRead(d, meta)
}

func resourceTencentCloudAPIGatewayUsagePlanDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_usage_plan.delete")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		usagePlanId       = d.Id()
	)

	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr := apiGatewayService.DeleteUsagePlan(ctx, usagePlanId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
}
