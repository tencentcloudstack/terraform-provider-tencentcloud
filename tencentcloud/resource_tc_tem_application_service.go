/*
Provides a resource to create a tem application_service

Example Usage

```hcl
resource "tencentcloud_tem_application_service" "application_service" {
  environment_id = "en-xxx"
  application_id = "xxx"
  service {
		type = "CLUSTER"
		service_name = "consumer"
		i_p = ""
		subnet_id = "subnet-xxxx"
		port_mapping_item_list {
			port = 80
			target_port = 80
			protocol = "TCP"
		}

  }
}
```

Import

tem application_service can be imported using the id, e.g.

```
terraform import tencentcloud_tem_application_service.application_service application_service_id
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
	"strings"
)

func resourceTencentCloudTemApplicationService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemApplicationServiceCreate,
		Read:   resourceTencentCloudTemApplicationServiceRead,
		Update: resourceTencentCloudTemApplicationServiceUpdate,
		Delete: resourceTencentCloudTemApplicationServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment ID.",
			},

			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application ID.",
			},

			"service": {
				Optional:    true,
				Description: "Service detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Application service type: EXTERNAL | VPC | CLUSTER.",
						},
						"service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Application service name.",
						},
						"i_p": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ip address of application service.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Subnet id of instance for type VPC.",
						},
						"port_mapping_item_list": {
							Optional:    true,
							Description: "Port mapping item list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Container port.",
									},
									"target_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Application listen port.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UDP or TCP.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemApplicationServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application_service.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewCreateApplicationServiceRequest()
		response      = tem.NewCreateApplicationServiceResponse()
		environmentId string
		applicationId string
	)
	if v, ok := d.GetOk("environment_id"); ok {
		environmentId = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		applicationId = v.(string)
		request.ApplicationId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("service"); v != nil {
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateApplicationService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem applicationService failed, reason:%+v", logId, err)
		return err
	}

	environmentId = *response.Response.EnvironmentId
	d.SetId(strings.Join([]string{environmentId, applicationId}, FILED_SP))

	return resourceTencentCloudTemApplicationServiceRead(d, meta)
}

func resourceTencentCloudTemApplicationServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application_service.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]

	applicationService, err := service.DescribeTemApplicationServiceById(ctx, environmentId, applicationId)
	if err != nil {
		return err
	}

	if applicationService == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TemApplicationService` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if applicationService.EnvironmentId != nil {
		_ = d.Set("environment_id", applicationService.EnvironmentId)
	}

	if applicationService.ApplicationId != nil {
		_ = d.Set("application_id", applicationService.ApplicationId)
	}

	if applicationService.Service != nil {
	}

	return nil
}

func resourceTencentCloudTemApplicationServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application_service.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyApplicationServiceRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]

	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId

	immutableArgs := []string{"environment_id", "application_id", "service"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, _ := d.GetOk("service"); v != nil {
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyApplicationService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tem applicationService failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTemApplicationServiceRead(d, meta)
}

func resourceTencentCloudTemApplicationServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_application_service.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]

	if err := service.DeleteTemApplicationServiceById(ctx, environmentId, applicationId); err != nil {
		return err
	}

	return nil
}
