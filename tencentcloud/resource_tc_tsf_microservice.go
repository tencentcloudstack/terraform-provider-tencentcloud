/*
Provides a resource to create a tsf microservice

Example Usage

```hcl
resource "tencentcloud_tsf_microservice" "microservice" {
  namespace_id = "namespace-vjlkzkgy"
  microservice_name = "test-microservice"
  microservice_desc = "desc-microservice"
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

tsf microservice can be imported using the namespaceId#microserviceId, e.g.

```
terraform import tencentcloud_tsf_microservice.microservice namespace-vjlkzkgy#ms-vjeb43lw
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request          = tsf.NewCreateMicroserviceRequest()
		response         = tsf.NewCreateMicroserviceResponse()
		microserviceName string
		namespaceId      string
		microserviceId   string
	)
	if v, ok := d.GetOk("namespace_id"); ok {
		namespaceId = v.(string)
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("microservice_name"); ok {
		microserviceName = v.(string)
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

	if *response.Response.Result {
		service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
		microservice, err := service.DescribeTsfMicroserviceById(ctx, namespaceId, "", microserviceName)
		if err != nil {
			return err
		}

		microserviceId = *microservice.MicroserviceId
		d.SetId(namespaceId + FILED_SP + microserviceId)
	} else {
		return fmt.Errorf("[DEBUG]%s api[%s] Creation failed, and the return result of interface creation is false", logId, request.GetAction())
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:microservice/%s", region, microserviceId)
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

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespaceId := idSplit[0]
	microserviceId := idSplit[1]

	microservice, err := service.DescribeTsfMicroserviceById(ctx, namespaceId, microserviceId, "")
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
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "microservice", tcClient.Region, microserviceId)
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := tsf.NewModifyMicroserviceRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	// namespaceId := idSplit[0]
	microserviceId := idSplit[1]

	request.MicroserviceId = &microserviceId

	immutableArgs := []string{"namespace_id", "microservice_name"}

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
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tsf", "microservice", tcClient.Region, microserviceId)
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
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	// namespaceId := idSplit[0]
	microserviceId := idSplit[1]

	if err := service.DeleteTsfMicroserviceById(ctx, microserviceId); err != nil {
		return err
	}

	return nil
}
