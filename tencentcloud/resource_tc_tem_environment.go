/*
Provides a resource to create a tem environment

Example Usage

```hcl
resource "tencentcloud_tem_environment" "environment" {
  environment_name = "xxx"
  description = "xxx"
  vpc = "vpc-xxx"
  subnet_ids =
  tags {
		tag_key = "key"
		tag_value = "tag value"

  }
}
```

Import

tem environment can be imported using the id, e.g.

```
terraform import tencentcloud_tem_environment.environment environment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudTemEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemEnvironmentCreate,
		Read:   resourceTencentCloudTemEnvironmentRead,
		Update: resourceTencentCloudTemEnvironmentUpdate,
		Delete: resourceTencentCloudTemEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment name.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Environment description.",
			},

			"vpc": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Vpc ID.",
			},

			"subnet_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Subnet IDs.",
			},

			"tags": {
				Optional:    true,
				Description: "Environment tag list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_environment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewCreateEnvironmentRequest()
		response      = tem.NewCreateEnvironmentResponse()
		environmentId string
	)
	if v, ok := d.GetOk("environment_name"); ok {
		request.EnvironmentName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc"); ok {
		request.Vpc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIdsSet := v.(*schema.Set).List()
		for i := range subnetIdsSet {
			subnetIds := subnetIdsSet[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetIds)
		}
	}

	if v, _ := d.GetOk("tags"); v != nil {
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateEnvironment(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem environment failed, reason:%+v", logId, err)
		return err
	}

	environmentId = *response.Response.EnvironmentId
	d.SetId(environmentId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"NORMAL"}, 10*readRetryTimeout, time.Second, service.TemEnvironmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTemEnvironmentRead(d, meta)
}

func resourceTencentCloudTemEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_environment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	environmentId := d.Id()

	environment, err := service.DescribeTemEnvironmentById(ctx, environmentId)
	if err != nil {
		return err
	}

	if environment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TemEnvironment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if environment.EnvironmentName != nil {
		_ = d.Set("environment_name", environment.EnvironmentName)
	}

	if environment.Description != nil {
		_ = d.Set("description", environment.Description)
	}

	if environment.Vpc != nil {
		_ = d.Set("vpc", environment.Vpc)
	}

	if environment.SubnetIds != nil {
		_ = d.Set("subnet_ids", environment.SubnetIds)
	}

	if environment.tags != nil {
	}

	return nil
}

func resourceTencentCloudTemEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_environment.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyEnvironmentRequest()

	environmentId := d.Id()

	request.EnvironmentId = &environmentId

	immutableArgs := []string{"environment_name", "description", "vpc", "subnet_ids", "tags"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("environment_name") {
		if v, ok := d.GetOk("environment_name"); ok {
			request.EnvironmentName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("subnet_ids") {
		if v, ok := d.GetOk("subnet_ids"); ok {
			subnetIdsSet := v.(*schema.Set).List()
			for i := range subnetIdsSet {
				subnetIds := subnetIdsSet[i].(string)
				request.SubnetIds = append(request.SubnetIds, &subnetIds)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyEnvironment(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tem environment failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTemEnvironmentRead(d, meta)
}

func resourceTencentCloudTemEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_environment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	environmentId := d.Id()

	if err := service.DeleteTemEnvironmentById(ctx, environmentId); err != nil {
		return err
	}

	return nil
}
