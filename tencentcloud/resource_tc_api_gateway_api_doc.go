/*
Provides a resource to create a APIGateway ApiDoc

Example Usage

```hcl
resource "tencentcloud_api_gateway_api_doc" "my_api_doc" {
  api_doc_name = "doc_test1"
  service_id   = "service_test1"
  environment  = "release"
  apilds       = ["api-test1", "api-test2"]
}
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudAPIGatewayAPIDoc() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudAPIGatewayAPIDocRead,
		Create: resourceTencentCloudAPIGatewayAPIDocCreate,
		Update: resourceTencentCloudAPIGatewayAPIDocUpdate,
		Delete: resourceTencentCloudAPIGatewayAPIDocDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_doc_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Api doc name.",
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Service name.",
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Env name.",
			},
			"api_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of APIs for generating documents.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayAPIDocRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_doc.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		apiDocId          = d.Id()
		apiDocInfo        *apiGateway.APIDocInfo
		err               error
	)

	apiDocInfo, err = apiGatewayService.DescribeApiDoc(ctx, apiDocId)
	if err != nil {
		return err
	}

	if apiDocInfo == nil {
		d.SetId("")
		return fmt.Errorf("resource `api_doc` %s does not exist", apiDocId)
	}

	if apiDocInfo.ApiDocId != nil {
		_ = d.Set("api_doc_id", apiDocInfo.ApiDocId)
	}

	if apiDocInfo.ApiDocName != nil {
		_ = d.Set("api_doc_name", apiDocInfo.ApiDocName)
	}

	if apiDocInfo.ApiCount != nil {
		_ = d.Set("api_count", apiDocInfo.ApiCount)
	}

	if apiDocInfo.ViewCount != nil {
		_ = d.Set("view_count", apiDocInfo.ViewCount)
	}

	if apiDocInfo.ReleaseCount != nil {
		_ = d.Set("release_count", apiDocInfo.ReleaseCount)
	}

	if apiDocInfo.ApiDocUri != nil {
		_ = d.Set("api_doc_uri", apiDocInfo.ApiDocUri)
	}

	if apiDocInfo.SharePassword != nil {
		_ = d.Set("share_password", apiDocInfo.SharePassword)
	}

	if apiDocInfo.ApiDocStatus != nil {
		_ = d.Set("api_doc_status", apiDocInfo.ApiDocStatus)
	}

	if apiDocInfo.UpdatedTime != nil {
		_ = d.Set("updated_time", apiDocInfo.UpdatedTime)
	}

	if apiDocInfo.ServiceId != nil {
		_ = d.Set("service_id", apiDocInfo.ServiceId)
	}

	if apiDocInfo.ServiceName != nil {
		_ = d.Set("service_name", apiDocInfo.ServiceName)
	}

	if apiDocInfo.ApiIds != nil {
		apiIdsList := []interface{}{}
		for _, id := range apiDocInfo.ApiIds {
			apiIdsList = append(apiIdsList, id)
		}
		_ = d.Set("api_ids", apiIdsList)
	}

	if apiDocInfo.ApiNames != nil {
		apiNameList := []interface{}{}
		for _, name := range apiDocInfo.ApiNames {
			apiNameList = append(apiNameList, name)
		}
		_ = d.Set("api_names", apiNameList)
	}

	if apiDocInfo.Environment != nil {
		_ = d.Set("environment", apiDocInfo.Environment)
	}

	return nil
}

func resourceTencentCloudAPIGatewayAPIDocCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_doc.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = apiGateway.NewCreateAPIDocRequest()
		response *apiGateway.CreateAPIDocResponse
		apiDocId string
		err      error
	)

	if v, ok := d.GetOk("api_doc_name"); ok {
		request.ApiDocName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("environment"); ok {
		request.Environment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("api_ids"); ok {
		for _, item := range v.([]interface{}) {
			id := helper.String(item.(string))
			request.ApiIds = append(request.ApiIds, id)
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, err := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().CreateAPIDoc(request)
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
		log.Printf("[CRITAL]%s create api_doc failed, reason:%+v", logId, err)
		return err
	}

	apiDocId = *response.Response.Result.ApiDocId
	d.SetId(apiDocId)
	return resourceTencentCloudAPIGatewayAPIDocRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIDocUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_doc.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = apiGateway.NewModifyAPIDocRequest()
		apiDocId = d.Id()
		err      error
	)

	request.ApiDocId = &apiDocId
	if d.HasChange("api_doc_name") {
		if v, ok := d.GetOk("api_doc_name"); ok {
			request.ApiDocName = helper.String(v.(string))
		}
	}

	if d.HasChange("service_id") {
		if v, ok := d.GetOk("service_id"); ok {
			request.ServiceId = helper.String(v.(string))
		}
	}

	if d.HasChange("environment") {
		if v, ok := d.GetOk("environment"); ok {
			request.Environment = helper.String(v.(string))
		}
	}

	if d.HasChange("api_ids") {
		if v, ok := d.GetOk("api_ids"); ok {
			for _, item := range v.([]interface{}) {
				id := helper.String(item.(string))
				request.ApiIds = append(request.ApiIds, id)
			}
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().ModifyAPIDoc(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tat command failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTatCommandRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIDocDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_doc.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		apiDocId          = d.Id()
		err               error
	)

	if err = apiGatewayService.DeleteAPIGatewayAPIDocById(ctx, apiDocId); err != nil {
		return err
	}

	return nil
}
