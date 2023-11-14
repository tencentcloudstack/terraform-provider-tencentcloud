/*
Provides a resource to create a tsf unit_namespace

Example Usage

```hcl
resource "tencentcloud_tsf_unit_namespace" "unit_namespace" {
  gateway_instance_id = ""
  unit_namespace_list {
		namespace_id = ""
		namespace_name = ""
		id = ""
		gateway_instance_id = ""
		created_time = ""
		updated_time = ""

  }
}
```

Import

tsf unit_namespace can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_unit_namespace.unit_namespace unit_namespace_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfUnitNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfUnitNamespaceCreate,
		Read:   resourceTencentCloudTsfUnitNamespaceRead,
		Delete: resourceTencentCloudTsfUnitNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Gateway instance Id.",
			},

			"unit_namespace_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Unit namespace list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Namespace id.",
						},
						"namespace_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Namespace name.",
						},
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unit namespace ID. Note: This field may return null, indicating that no valid value was found.",
						},
						"gateway_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Gateway instance id Note: This field may return null, indicating that no valid value was found.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Create time. Note: This field may return null, indicating that no valid value was found.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time. Note: This field may return null, indicating that no valid value was found.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfUnitNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_unit_namespace.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request           = tsf.NewCreateUnitNamespacesRequest()
		response          = tsf.NewCreateUnitNamespacesResponse()
		gatewayInstanceId string
	)
	if v, ok := d.GetOk("gateway_instance_id"); ok {
		gatewayInstanceId = v.(string)
		request.GatewayInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("unit_namespace_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			unitNamespace := tsf.UnitNamespace{}
			if v, ok := dMap["namespace_id"]; ok {
				unitNamespace.NamespaceId = helper.String(v.(string))
			}
			if v, ok := dMap["namespace_name"]; ok {
				unitNamespace.NamespaceName = helper.String(v.(string))
			}
			if v, ok := dMap["id"]; ok {
				unitNamespace.Id = helper.String(v.(string))
			}
			if v, ok := dMap["gateway_instance_id"]; ok {
				unitNamespace.GatewayInstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["created_time"]; ok {
				unitNamespace.CreatedTime = helper.String(v.(string))
			}
			if v, ok := dMap["updated_time"]; ok {
				unitNamespace.UpdatedTime = helper.String(v.(string))
			}
			request.UnitNamespaceList = append(request.UnitNamespaceList, &unitNamespace)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateUnitNamespaces(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf unitNamespace failed, reason:%+v", logId, err)
		return err
	}

	gatewayInstanceId = *response.Response.GatewayInstanceId
	d.SetId(gatewayInstanceId)

	return resourceTencentCloudTsfUnitNamespaceRead(d, meta)
}

func resourceTencentCloudTsfUnitNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_unit_namespace.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	unitNamespaceId := d.Id()

	unitNamespace, err := service.DescribeTsfUnitNamespaceById(ctx, gatewayInstanceId)
	if err != nil {
		return err
	}

	if unitNamespace == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfUnitNamespace` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if unitNamespace.GatewayInstanceId != nil {
		_ = d.Set("gateway_instance_id", unitNamespace.GatewayInstanceId)
	}

	if unitNamespace.UnitNamespaceList != nil {
		unitNamespaceListList := []interface{}{}
		for _, unitNamespaceList := range unitNamespace.UnitNamespaceList {
			unitNamespaceListMap := map[string]interface{}{}

			if unitNamespace.UnitNamespaceList.NamespaceId != nil {
				unitNamespaceListMap["namespace_id"] = unitNamespace.UnitNamespaceList.NamespaceId
			}

			if unitNamespace.UnitNamespaceList.NamespaceName != nil {
				unitNamespaceListMap["namespace_name"] = unitNamespace.UnitNamespaceList.NamespaceName
			}

			if unitNamespace.UnitNamespaceList.Id != nil {
				unitNamespaceListMap["id"] = unitNamespace.UnitNamespaceList.Id
			}

			if unitNamespace.UnitNamespaceList.GatewayInstanceId != nil {
				unitNamespaceListMap["gateway_instance_id"] = unitNamespace.UnitNamespaceList.GatewayInstanceId
			}

			if unitNamespace.UnitNamespaceList.CreatedTime != nil {
				unitNamespaceListMap["created_time"] = unitNamespace.UnitNamespaceList.CreatedTime
			}

			if unitNamespace.UnitNamespaceList.UpdatedTime != nil {
				unitNamespaceListMap["updated_time"] = unitNamespace.UnitNamespaceList.UpdatedTime
			}

			unitNamespaceListList = append(unitNamespaceListList, unitNamespaceListMap)
		}

		_ = d.Set("unit_namespace_list", unitNamespaceListList)

	}

	return nil
}

func resourceTencentCloudTsfUnitNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_unit_namespace.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	unitNamespaceId := d.Id()

	if err := service.DeleteTsfUnitNamespaceById(ctx, gatewayInstanceId); err != nil {
		return err
	}

	return nil
}
