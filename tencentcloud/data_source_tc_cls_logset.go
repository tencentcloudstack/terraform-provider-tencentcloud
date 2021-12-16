package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func dataSourceTencentCloudClsLogset() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsLogsetRead,
		Schema: map[string]*schema.Schema{
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The offset of paging. The default value is 0",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "The limit number of single page paging. The default value is 20 and the maximum value is 100",
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    10,
				Description: "filters of cls logsets,The upper limit of filters per request is 10 ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields to be filtered, only supported logsetName,logsetId,tagKey,tag:tagKey",
						},
						"value": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							MaxItems:    5,
							Required:    true,
							Description: "Values to be filtered,the upper limit of filter.values is 5.",
						},
					},
				},
			},
			"logsets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "logset lists",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logset_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of logset",
						},
						"logset_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of logset",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"topic_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of log topics under the logset",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If assumeruin is not empty, it indicates the server role that created the log set",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Label of log set binding",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag Key",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value",
									},
								},
							},
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of pages",
			},
		},
	}
}

func dataSourceTencentCloudClsLogsetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cls_logset.read")()

	logId := getLogId(contextNil)

	config := meta.(*TencentCloudClient).apiV3Conn
	request := cls.NewDescribeLogsetsRequest()
	if v, ok := d.GetOk("limit"); ok {
		request.Limit = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("offset"); ok {
		request.Offset = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("filters"); ok {
		filters := v.([]interface{})
		request.Filters = make([]*cls.Filter, 0, len(filters))
		for i := range filters {
			var valueArray []string
			filter := filters[i].(map[string]interface{})
			if values, ok := filter["value"].([]interface{}); ok {
				valueArray = make([]string, 0, len(values))
				for _, value := range values {
					valueArray = append(valueArray, value.(string))
				}
				filterGet := cls.Filter{
					Key:    helper.String(filter["key"].(string)),
					Values: helper.Strings(valueArray),
				}
				request.Filters = append(request.Filters, &filterGet)
			}
		}
	}
	var logsets *cls.DescribeLogsetsResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := config.UseClsClient().DescribeLogsets(request)
		if e != nil {
			return retryError(e, InternalError)
		}
		logsets = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLS logsets failed, reason:%s\n", logId, err.Error())
		return err
	}
	logInfos := logsets.Response.Logsets
	logsetList := make([]map[string]interface{}, 0, len(logInfos))
	ids := make([]string, 0, len(logInfos))
	for _, logInfo := range logInfos {
		mapping := map[string]interface{}{
			"logset_id":   *logInfo.LogsetId,
			"logset_name": *logInfo.LogsetName,
			"create_time": *logInfo.CreateTime,
			"topic_count": *logInfo.TopicCount,
			"role_name":   *logInfo.RoleName,
			"tags":        flattenDataTagMappings(logInfo.Tags),
		}
		logsetList = append(logsetList, mapping)
		ids = append(ids, *logInfo.LogsetId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("logsets", logsetList); e != nil {
		log.Printf("[CRITAL]%s provider set logset list fail, reason:%s\n", logId, e.Error())
		return e
	}
	if e := d.Set("total_count", *logsets.Response.TotalCount); e != nil {
		log.Printf("[CRITAL]%s provider set logset list fail, reason:%s\n", logId, e.Error())
		return e
	}
	return nil
}
func flattenDataTagMappings(list []*cls.Tag) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		tag := map[string]interface{}{
			"key":   *v.Key,
			"value": *v.Value,
		}
		result = append(result, tag)
	}
	return result
}
