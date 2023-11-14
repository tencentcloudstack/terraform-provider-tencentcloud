/*
Provides a resource to create a tsf microservice

Example Usage

```hcl
resource "tencentcloud_tsf_microservice" "microservice" {
  namespace_id = ""
  microservice_name = ""
  microservice_desc = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf microservice can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_microservice.microservice microservice_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfMicroservice() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfMicroserviceCreate,
		Read:   resourceTencentCloudTsfMicroserviceRead,
		Update: resourceTencentCloudTsfMicroserviceUpdate,
		Delete: resourceTencentCloudTsfMicroserviceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace ID.",
			},

			"microservice_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Microservice name.",
			},

			"microservice_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Microservice description information.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfMicroserviceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_microservice.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = tsf.NewCreateMicroserviceRequest()
		response       = tsf.NewCreateMicroserviceResponse()
		microserviceId string
	)
	if v, ok := d.GetOk("namespace_id"); ok {
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("microservice_name"); ok {
		request.MicroserviceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("microservice_desc"); ok {
		request.MicroserviceDesc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateMicroservice(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf microservice failed, reason:%+v", logId, err)
		return err
	}

	microserviceId = *response.Response.MicroserviceId
	d.SetId(microserviceId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:microserviceId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfMicroserviceRead(d, meta)
}

func resourceTencentCloudTsfMicroserviceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_microservice.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	microserviceId := d.Id()

	microservice, err := service.DescribeTsfMicroserviceById(ctx, microserviceId)
	if err != nil {
		return err
	}

	if microservice == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfMicroservice` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if microservice.NamespaceId != nil {
		_ = d.Set("namespace_id", microservice.NamespaceId)
	}

	if microservice.MicroserviceName != nil {
		_ = d.Set("microservice_name", microservice.MicroserviceName)
	}

	if microservice.MicroserviceDesc != nil {
		_ = d.Set("microservice_desc", microservice.MicroserviceDesc)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "microserviceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfMicroserviceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_microservice.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyMicroserviceRequest()

	microserviceId := d.Id()

	request.MicroserviceId = &microserviceId

	immutableArgs := []string{"namespace_id", "microservice_name", "microservice_desc"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("microservice_desc") {
		if v, ok := d.GetOk("microservice_desc"); ok {
			request.MicroserviceDesc = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyMicroservice(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf microservice failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tsf", "microserviceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfMicroserviceRead(d, meta)
}

func resourceTencentCloudTsfMicroserviceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_microservice.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	microserviceId := d.Id()

	if err := service.DeleteTsfMicroserviceById(ctx, microserviceId); err != nil {
		return err
	}

	return nil
}
