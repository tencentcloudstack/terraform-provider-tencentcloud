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
	"crypto/md5"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func dataSourceTencentMonitorBindingObjects() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentMonitorBindingObjectRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy group id for query",
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
							Description: "Object unique id.",
						},
						"dimensions_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Represents a collection of dimensions of an object instance, json format.",
						},
						"is_shielded": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the object is shielded or not, 0 means unshielded and 1 means shielded.",
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

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDescribeBindingPolicyObjectListRequest()
		response       *monitor.DescribeBindingPolicyObjectListResponse
		objects        []*monitor.DescribeBindingPolicyObjectListInstance
		offset         int64 = 0
		limit          int64 = 20
		err            error
		finish         bool
		list           = make([]interface{}, 0, len(objects))
	)

	request.GroupId = helper.IntInt64(d.Get("group_id").(int))
	request.Module = helper.String("monitor")
	request.Offset = &offset
	request.Limit = &limit

	for {
		if finish {
			break
		}
		if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if response, err = monitorService.client.UseMonitorClient().DescribeBindingPolicyObjectList(request); err != nil {
				return retryError(err, InternalError)
			}
			objects = append(objects, response.Response.List...)
			if len(response.Response.List) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			return err
		}

		offset = offset + limit
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

	md := md5.New()
	md.Write([]byte(request.ToJsonString()))
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)

	if output, ok := d.GetOk("result_output_file"); ok {
		return writeToFile(output.(string), list)
	}
	return nil
}
