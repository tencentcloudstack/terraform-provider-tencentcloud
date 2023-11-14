/*
Provides a resource to create a apigateway service

Example Usage

```hcl
resource "tencentcloud_apigateway_service" "service" {
  service_name = ""
  protocol = ""
  service_desc = ""
  net_types =
  ip_version = ""
  set_server_name = ""
  app_id_type = ""
  instance_id = ""
  uniq_vpc_id = ""
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

apigateway service can be imported using the id, e.g.

```
terraform import tencentcloud_apigateway_service.service service_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudApigatewayService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApigatewayServiceCreate,
		Read:   resourceTencentCloudApigatewayServiceRead,
		Update: resourceTencentCloudApigatewayServiceUpdate,
		Delete: resourceTencentCloudApigatewayServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "User defined service name.",
			},

			"protocol": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The front-end request type of the service. Such as http, https, and http&amp;amp;amp;https.",
			},

			"service_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User defined service description.",
			},

			"net_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of network types used to specify the supported access types, INNER for internal network access and OUTER for external network access. The default is OUTER.",
			},

			"ip_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "IP version number, supports IPv4 and IPv6, defaults to IPv4.",
			},

			"set_server_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster name. Reserved field, used for tsf serverless type.",
			},

			"app_id_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User type. Reserved type for serverless users.",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Exclusive instance ID.",
			},

			"uniq_vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC Properties.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudApigatewayServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_service.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = apigateway.NewCreateServiceRequest()
		response  = apigateway.NewCreateServiceResponse()
		serviceId string
	)
	if v, ok := d.GetOk("service_name"); ok {
		request.ServiceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_desc"); ok {
		request.ServiceDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("net_types"); ok {
		netTypesSet := v.(*schema.Set).List()
		for i := range netTypesSet {
			netTypes := netTypesSet[i].(string)
			request.NetTypes = append(request.NetTypes, &netTypes)
		}
	}

	if v, ok := d.GetOk("ip_version"); ok {
		request.IpVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("set_server_name"); ok {
		request.SetServerName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("app_id_type"); ok {
		request.AppIdType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().CreateService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create apigateway service failed, reason:%+v", logId, err)
		return err
	}

	serviceId = *response.Response.ServiceId
	d.SetId(serviceId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::apigw:%s:uin/:serviceId/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayServiceRead(d, meta)
}

func resourceTencentCloudApigatewayServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_service.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}

	serviceId := d.Id()

	service, err := service.DescribeApigatewayServiceById(ctx, serviceId)
	if err != nil {
		return err
	}

	if service == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayService` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if service.ServiceName != nil {
		_ = d.Set("service_name", service.ServiceName)
	}

	if service.Protocol != nil {
		_ = d.Set("protocol", service.Protocol)
	}

	if service.ServiceDesc != nil {
		_ = d.Set("service_desc", service.ServiceDesc)
	}

	if service.NetTypes != nil {
		_ = d.Set("net_types", service.NetTypes)
	}

	if service.IpVersion != nil {
		_ = d.Set("ip_version", service.IpVersion)
	}

	if service.SetServerName != nil {
		_ = d.Set("set_server_name", service.SetServerName)
	}

	if service.AppIdType != nil {
		_ = d.Set("app_id_type", service.AppIdType)
	}

	if service.InstanceId != nil {
		_ = d.Set("instance_id", service.InstanceId)
	}

	if service.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", service.UniqVpcId)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "apigw", "serviceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudApigatewayServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_service.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		updateServiceRequest  = apigateway.NewUpdateServiceRequest()
		updateServiceResponse = apigateway.NewUpdateServiceResponse()
	)

	serviceId := d.Id()

	request.ServiceId = &serviceId

	immutableArgs := []string{"service_name", "protocol", "service_desc", "net_types", "ip_version", "set_server_name", "app_id_type", "instance_id", "uniq_vpc_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("service_name") {
		if v, ok := d.GetOk("service_name"); ok {
			request.ServiceName = helper.String(v.(string))
		}
	}

	if d.HasChange("protocol") {
		if v, ok := d.GetOk("protocol"); ok {
			request.Protocol = helper.String(v.(string))
		}
	}

	if d.HasChange("service_desc") {
		if v, ok := d.GetOk("service_desc"); ok {
			request.ServiceDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("net_types") {
		if v, ok := d.GetOk("net_types"); ok {
			netTypesSet := v.(*schema.Set).List()
			for i := range netTypesSet {
				netTypes := netTypesSet[i].(string)
				request.NetTypes = append(request.NetTypes, &netTypes)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseApigatewayClient().UpdateService(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update apigateway service failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("apigw", "serviceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudApigatewayServiceRead(d, meta)
}

func resourceTencentCloudApigatewayServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_apigateway_service.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ApigatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
	serviceId := d.Id()

	if err := service.DeleteApigatewayServiceById(ctx, serviceId); err != nil {
		return err
	}

	return nil
}
