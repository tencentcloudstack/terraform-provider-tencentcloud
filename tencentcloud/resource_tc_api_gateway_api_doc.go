/*
Provides a resource to create a APIGateway ApiDoc

Example Usage

```hcl
resource "tencentcloud_api_gateway_api_doc" "my_api_doc" {
  api_doc_name = "doc_test1"
  service_id   = "service_test1"
  environment  = "release"
  api_ids      = ["api-test1", "api-test2"]
}
```
*/
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

func resourceTencentCloudAPIGatewayAPIDoc() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayAPIDocCreate,
		Read:   resourceTencentCloudAPIGatewayAPIDocRead,
		Update: resourceTencentCloudAPIGatewayAPIDocUpdate,
		Delete: resourceTencentCloudAPIGatewayAPIDocDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_doc_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Api Document name.",
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
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of APIs for generating documents.",
			},
			"api_doc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Api Document ID.",
			},
			"api_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Api Document count.",
			},
			"view_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "API Document Viewing Times.",
			},
			"release_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of API document releases.",
			},
			"api_doc_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API Document Access URI.",
			},
			"share_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API Document Sharing Password.",
			},
			"api_doc_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API Document Build Status.",
			},
			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API Document update time.",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API Document service name.",
			},
			"api_names": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of names for generating documents.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayAPIDocCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_doc.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		request           = apiGateway.NewCreateAPIDocRequest()
		response          *apiGateway.CreateAPIDocResponse
		apiDocId          string
		err               error
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
		for _, item := range v.(*schema.Set).List() {
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

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		apiDocInfo, err := apiGatewayService.DescribeApiDoc(ctx, apiDocId)
		if err != nil {
			return retryError(err)
		}
		if apiDocInfo == nil {
			err = fmt.Errorf("api_doc_id %s not exists", apiDocId)
			return resource.NonRetryableError(err)
		}
		if *apiDocInfo.ApiDocStatus == API_GATEWAY_API_DOC_STATUS_PROCESSING {
			return resource.RetryableError(fmt.Errorf("create api_doc task status is %s", API_GATEWAY_API_DOC_STATUS_PROCESSING))
		}
		if *apiDocInfo.ApiDocStatus == API_GATEWAY_API_DOC_STATUS_COMPLETED {
			return nil
		}
		if *apiDocInfo.ApiDocStatus == API_GATEWAY_API_DOC_STATUS_FAIL {
			return resource.NonRetryableError(fmt.Errorf("create api_doc task status is %s", API_GATEWAY_API_DOC_STATUS_FAIL))
		}
		err = fmt.Errorf("create api_doc task status is %v, we won't wait for it finish", *apiDocInfo.ApiDocStatus)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create api_doc task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(apiDocId)

	return resourceTencentCloudAPIGatewayAPIDocRead(d, meta)
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
		return nil
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

func resourceTencentCloudAPIGatewayAPIDocUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_doc.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		request           = apiGateway.NewModifyAPIDocRequest()
		apiDocId          = d.Id()
		err               error
	)

	request.ApiDocId = &apiDocId

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
		for _, item := range v.(*schema.Set).List() {
			id := helper.String(item.(string))
			request.ApiIds = append(request.ApiIds, id)
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
		log.Printf("[CRITAL]%s update api_doc failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		apiDocInfo, err := apiGatewayService.DescribeApiDoc(ctx, apiDocId)
		if err != nil {
			return retryError(err)
		}
		if apiDocInfo == nil {
			err = fmt.Errorf("api_doc_id %s not exists", apiDocId)
			return resource.NonRetryableError(err)
		}
		if *apiDocInfo.ApiDocStatus == API_GATEWAY_API_DOC_STATUS_PROCESSING {
			return resource.RetryableError(fmt.Errorf("create api_doc task status is %s", API_GATEWAY_API_DOC_STATUS_PROCESSING))
		}
		if *apiDocInfo.ApiDocStatus == API_GATEWAY_API_DOC_STATUS_COMPLETED {
			return nil
		}
		if *apiDocInfo.ApiDocStatus == API_GATEWAY_API_DOC_STATUS_FAIL {
			return resource.NonRetryableError(fmt.Errorf("create api_doc task status is %s", API_GATEWAY_API_DOC_STATUS_PROCESSING))
		}
		err = fmt.Errorf("create api_doc task status is %v, we won't wait for it finish", *apiDocInfo.ApiDocStatus)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s update api_doc task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudAPIGatewayAPIDocRead(d, meta)
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
