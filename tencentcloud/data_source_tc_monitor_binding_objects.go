/*
Use this data source to query policy group binding objects.

Example Usage

```hcl
data "tencentcloud_monitor_policy_groups" "name" {
  name = "test"
}

data "tencentcloud_monitor_binding_objects" "objects" {
  group_id = data.tencentcloud_monitor_policy_groups.name.list[0].group_id
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
)

func dataSourceTencentMonitorBindingObjects() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentMonitorBindingObjectRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy group ID for query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list objects. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Object unique ID.",
						},
						"dimensions_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Represents a collection of dimensions of an object instance, json format.",
						},
						"is_shielded": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the object is shielded or not, `0` means unshielded and `1` means shielded.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region where the object is located.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorBindingObjectRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_binding_objects.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		objects        []*monitor.DescribeBindingPolicyObjectListInstance
		err            error

		list = make([]interface{}, 0, len(objects))
	)

	id := int64(d.Get("group_id").(int))

	objects, err = monitorService.DescribeBindingPolicyObjectList(ctx, id)
	if err != nil {
		return err
	}

	for _, event := range objects {
		var listItem = map[string]interface{}{}
		listItem["region"] = event.Region
		listItem["unique_id"] = event.UniqueId
		listItem["dimensions_json"] = event.Dimensions
		listItem["is_shielded"] = event.IsShielded
		listItem["region"] = event.Region
		list = append(list, listItem)
	}
	if err = d.Set("list", list); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", id))
	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), list)
	}
	return nil
}
