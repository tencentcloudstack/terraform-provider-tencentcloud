/*
Provides a resource to create a tsf bind_api_group

Example Usage

```hcl
resource "tencentcloud_tsf_bind_api_group" "bind_api_group" {
  group_gateway_list {
		gateway_deploy_group_id = "group-vzd97zpy"
		group_id = "grp-qp0rj3zi"

  }
}
```

Import

tsf bind_api_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_bind_api_group.bind_api_group bind_api_group_id
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
	"strings"
)

func resourceTencentCloudTsfBindApiGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfBindApiGroupCreate,
		Read:   resourceTencentCloudTsfBindApiGroupRead,
		Delete: resourceTencentCloudTsfBindApiGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_gateway_list": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Api group bind with gateway Group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_deploy_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Gateway group id.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Group id.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfBindApiGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_bind_api_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request              = tsf.NewBindApiGroupRequest()
		response             = tsf.NewBindApiGroupResponse()
		groupId              string
		gatewayDeployGroupId string
	)
	if v, ok := d.GetOk("group_gateway_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			gatewayGroupIds := tsf.GatewayGroupIds{}
			if v, ok := dMap["gateway_deploy_group_id"]; ok {
				gatewayGroupIds.GatewayDeployGroupId = helper.String(v.(string))
			}
			if v, ok := dMap["group_id"]; ok {
				gatewayGroupIds.GroupId = helper.String(v.(string))
			}
			request.GroupGatewayList = append(request.GroupGatewayList, &gatewayGroupIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().BindApiGroup(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf bindApiGroup failed, reason:%+v", logId, err)
		return err
	}

	groupId = *response.Response.groupId
	d.SetId(strings.Join([]string{groupId, gatewayDeployGroupId}, FILED_SP))

	return resourceTencentCloudTsfBindApiGroupRead(d, meta)
}

func resourceTencentCloudTsfBindApiGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_bind_api_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	groupId := idSplit[0]
	gatewayDeployGroupId := idSplit[1]

	bindApiGroup, err := service.DescribeTsfBindApiGroupById(ctx, groupId, gatewayDeployGroupId)
	if err != nil {
		return err
	}

	if bindApiGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfBindApiGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if bindApiGroup.GroupGatewayList != nil {
		groupGatewayListList := []interface{}{}
		for _, groupGatewayList := range bindApiGroup.GroupGatewayList {
			groupGatewayListMap := map[string]interface{}{}

			if bindApiGroup.GroupGatewayList.GatewayDeployGroupId != nil {
				groupGatewayListMap["gateway_deploy_group_id"] = bindApiGroup.GroupGatewayList.GatewayDeployGroupId
			}

			if bindApiGroup.GroupGatewayList.GroupId != nil {
				groupGatewayListMap["group_id"] = bindApiGroup.GroupGatewayList.GroupId
			}

			groupGatewayListList = append(groupGatewayListList, groupGatewayListMap)
		}

		_ = d.Set("group_gateway_list", groupGatewayListList)

	}

	return nil
}

func resourceTencentCloudTsfBindApiGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_bind_api_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	groupId := idSplit[0]
	gatewayDeployGroupId := idSplit[1]

	if err := service.DeleteTsfBindApiGroupById(ctx, groupId, gatewayDeployGroupId); err != nil {
		return err
	}

	return nil
}
