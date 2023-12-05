package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAPIGatewayAPIApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayAPIAppCreate,
		Read:   resourceTencentCloudAPIGatewayAPIAppRead,
		Update: resourceTencentCloudAPIGatewayAPIAppUpdate,
		Delete: resourceTencentCloudAPIGatewayAPIAppDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_app_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Api app name.",
			},
			"api_app_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "App description.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"api_app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Api app ID.",
			},
			"api_app_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Api app key.",
			},
			"api_app_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Api app secret.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Api app created time.",
			},
			"modified_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Api app modified time.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayAPIAppCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_app.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		request  = apiGateway.NewCreateApiAppRequest()
		response *apiGateway.CreateApiAppResponse
		apiAppId string
		err      error
	)

	if v, ok := d.GetOk("api_app_name"); ok {
		request.ApiAppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_app_desc"); ok {
		request.ApiAppDesc = helper.String(v.(string))
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, err := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().CreateApiApp(request)
		if err != nil {
			return retryError(err)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create api_app failed, reason:%+v", logId, err)
		return err
	}

	apiAppId = *response.Response.Result.ApiAppId

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::apigateway:%s:uin/:apiAppId/%s", region, apiAppId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(apiAppId)
	return resourceTencentCloudAPIGatewayAPIAppRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIAppRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_app.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		apiAppId          = d.Id()
		apiAppInfo        *apiGateway.ApiAppInfos
		err               error
	)

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		apiAppInfo, err = apiGatewayService.DescribeApiApp(ctx, apiAppId)
		if err != nil {
			return retryError(err)
		}
		return nil
	})

	if apiAppInfo == nil {
		d.SetId("")
		log.Printf("resource `api_app` %s does not exist", apiAppId)
		return nil
	}

	apiAppData := apiAppInfo.ApiAppSet[0]
	if apiAppData.ApiAppId != nil {
		_ = d.Set("api_app_id", apiAppData.ApiAppId)
	}

	if apiAppData.ApiAppName != nil {
		_ = d.Set("api_app_name", apiAppData.ApiAppName)
	}

	if apiAppData.ApiAppKey != nil {
		_ = d.Set("api_app_key", apiAppData.ApiAppKey)
	}

	if apiAppData.ApiAppSecret != nil {
		_ = d.Set("api_app_secret", apiAppData.ApiAppSecret)
	}

	if apiAppData.CreatedTime != nil {
		_ = d.Set("created_time", apiAppData.CreatedTime)
	}

	if apiAppData.ModifiedTime != nil {
		_ = d.Set("modified_time", apiAppData.ModifiedTime)
	}

	if apiAppData.ApiAppDesc != nil {
		err = d.Set("api_app_desc", apiAppData.ApiAppDesc)
		if err != nil {
			return err
		}
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	tags, err := tagService.DescribeResourceTags(ctx, "apigateway", "apiAppId", tcClient.Region, apiAppId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudAPIGatewayAPIAppUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_app.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		request  = apiGateway.NewModifyApiAppRequest()
		apiAppId = d.Id()
		err      error
	)

	request.ApiAppId = &apiAppId
	if d.HasChange("api_app_name") {
		if v, ok := d.GetOk("api_app_name"); ok {
			request.ApiAppName = helper.String(v.(string))
		}
	}

	if d.HasChange("api_app_desc") {
		if v, ok := d.GetOk("api_app_desc"); ok {
			request.ApiAppDesc = helper.String(v.(string))
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().ModifyApiApp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update api_app failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("apigateway", "apiAppId", tcClient.Region, apiAppId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayAPIAppRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIAppDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_app.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		apiAppId          = d.Id()
		err               error
	)

	if err = apiGatewayService.DeleteAPIGatewayAPIAppById(ctx, apiAppId); err != nil {
		return err
	}

	return nil
}
