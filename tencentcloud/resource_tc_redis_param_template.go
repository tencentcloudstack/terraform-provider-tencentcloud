/*
Provides a resource to create a redis param_template

# Example Usage

```hcl

	resource "tencentcloud_redis_param_template" "param_template" {
	  name = ""
	  description = ""
	  product_type = ""
	  template_id = ""
	  param_list {
				key = ""
				value = ""

	  }
	}

```
Import

redis param_template can be imported using the id, e.g.
```
$ terraform import tencentcloud_redis_param_template.param_template param_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisParamTemplate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudRedisParamTemplateRead,
		Create: resourceTencentCloudRedisParamTemplateCreate,
		Update: resourceTencentCloudRedisParamTemplateUpdate,
		Delete: resourceTencentCloudRedisParamTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "参数模板名称。.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "参数模板描述。.",
			},

			"product_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "产品类型：1 – Redis2.8内存版（集群架构），2 – Redis2.8内存版（标准架构），3 – CKV 3.2内存版(标准架构)，4 – CKV 3.2内存版(集群架构)，5 – Redis2.8内存版（单机），6 – Redis4.0内存版（标准架构），7 – Redis4.0内存版（集群架构），8 – Redis5.0内存版（标准架构），9 – Redis5.0内存版（集群架构）。创建模板时必填，从源模板复制则不需要传入该参数。.",
			},

			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "源参数模板 ID。.",
			},

			"param_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "参数列表。.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "设置参数的名称。例如timeout。当前支持自定义的参数，请参见&lt;a href=&quot;https://cloud.tencent.com/document/product/239/49925&quot;&gt;参数配置&lt;/a&gt;。.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "设置参数名称对应的运行值。例如timeout对应运行值可设置为120， 单位为秒（s）。指当客户端连接闲置时间达到120 s时，将关闭连接。更多参数取值信息，请参见&lt;a href=&quot;https://cloud.tencent.com/document/product/239/49925&quot;&gt;参数配置&lt;/a&gt;。.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudRedisParamTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_param_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := RedisService{client: client}

	var (
		request = redis.NewCreateParamTemplateRequest()
		id      string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_type"); ok {
		request.ProductType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("template_id"); ok {
		request.TemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			instanceParam := redis.InstanceParam{}
			if v, ok := dMap["key"]; ok {
				instanceParam.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				instanceParam.Value = helper.String(v.(string))
			}

			request.ParamList = append(request.ParamList, &instanceParam)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		resTemplateId, e := service.CreateParamTemplate(ctx, request)
		if e != nil {
			return retryError(e)
		}
		id = resTemplateId
		return nil
	})

	if err != nil {
		return err
	}

	if id == "" {
		return fmt.Errorf("cannot get redis template id")
	}

	d.SetId(id)
	return resourceTencentCloudRedisParamTemplateRead(d, meta)
}

func resourceTencentCloudRedisParamTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_param_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	paramTemplateId := d.Id()

	paramTemplate, err := service.DescribeParamTemplateInfo(ctx, paramTemplateId)

	if err != nil {
		return err
	}

	if paramTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `param_template` %s does not exist", paramTemplateId)
	}

	if paramTemplate.Name != nil {
		_ = d.Set("name", paramTemplate.Name)
	}

	if paramTemplate.Description != nil {
		_ = d.Set("description", paramTemplate.Description)
	}

	if paramTemplate.ProductType != nil {
		_ = d.Set("product_type", paramTemplate.ProductType)
	}

	if len(paramTemplate.Items) > 0 {
		result := make([]interface{}, 0)
		for i := range paramTemplate.Items {
			item := paramTemplate.Items[i]
			listMap := map[string]interface{}{}
			if item.Name != nil {
				listMap["key"] = item.Name
			}
			if item.CurrentValue != nil {
				listMap["value"] = item.CurrentValue
			}

			result = append(result, item)
		}
		_ = d.Set("param_list", result)
	}

	return nil
}

func resourceTencentCloudRedisParamTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_param_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := redis.NewModifyParamTemplateRequest()

	templateId := d.Id()

	request.TemplateId = &templateId

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

	}

	if d.HasChange("product_type") {
		return fmt.Errorf("`product_type` do not support change now.")

	}

	if d.HasChange("template_id") {
		if v, ok := d.GetOk("template_id"); ok {
			request.TemplateId = helper.String(v.(string))
		}

	}

	if d.HasChange("param_list") {
		if v, ok := d.GetOk("param_list"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				instanceParam := redis.InstanceParam{}
				if v, ok := dMap["key"]; ok {
					instanceParam.Key = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					instanceParam.Value = helper.String(v.(string))
				}

				request.ParamList = append(request.ParamList, &instanceParam)
			}
		}

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ModifyParamTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudRedisParamTemplateRead(d, meta)
}

func resourceTencentCloudRedisParamTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_param_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	paramTemplateId := d.Id()

	request := redis.NewDeleteParamTemplateRequest()
	request.TemplateId = &paramTemplateId
	if err := service.DeleteParamTemplate(ctx, request); err != nil {
		return err
	}

	return nil
}
