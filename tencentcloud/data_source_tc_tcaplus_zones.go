/*
Use this data source to query tcaplus zones

Example Usage

```hcl
data "tencentcloud_tcaplus_zones" "null" {
  app_id = "19162256624"
}
data "tencentcloud_tcaplus_zones" "id" {
  app_id  = "19162256624"
  zone_id = "19162256624:1"
}
data "tencentcloud_tcaplus_zones" "name" {
  app_id    = "19162256624"
  zone_name = "test"
}
data "tencentcloud_tcaplus_zones" "all" {
  app_id    = "19162256624"
  zone_id   = "19162256624:1"
  zone_name = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudTcaplusZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusZonesRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the tcapplus application to be query.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone id to be query.",
			},
			"zone_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone name to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of tcaplus zones. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the tcapplus zone.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the tcapplus zone.",
						},
						"table_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of tables.",
						},
						"total_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total storage(MB).",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the tcapplus zone.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusZonesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcaplus_zones.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcaplusService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	applicationId := d.Get("app_id").(string)
	zoneId := d.Get("zone_id").(string)
	zoneName := d.Get("zone_name").(string)

	apps, err := service.DescribeZones(ctx, applicationId, zoneId, zoneName)
	if err != nil {
		apps, err = service.DescribeZones(ctx, applicationId, zoneId, zoneName)
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(apps))

	for _, app := range apps {
		listItem := make(map[string]interface{})
		listItem["zone_name"] = *app.TableGroupName
		listItem["zone_id"] = fmt.Sprintf("%s:%s", applicationId, *app.TableGroupId)
		listItem["table_count"] = *app.TableCount
		listItem["total_size"] = *app.TotalSize
		listItem["create_time"] = *app.CreatedTime
		list = append(list, listItem)
	}

	d.SetId("zone." + applicationId + "." + zoneId + "." + zoneName)
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}
	return nil

}
