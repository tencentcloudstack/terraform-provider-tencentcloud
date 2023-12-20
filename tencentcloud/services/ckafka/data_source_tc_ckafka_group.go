package ckafka

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaGroupRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search for the keyword.",
			},

			"group_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "GroupList.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "groupId.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol used by this group.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCkafkaGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["search_word"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var groups []*ckafka.DescribeGroup

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaGroupByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		groups = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groups))
	groupMapList := []interface{}{}
	for _, group := range groups {
		groupMap := map[string]interface{}{}

		if group.Group != nil {
			groupMap["group"] = group.Group
			ids = append(ids, *group.Group)
		}

		if group.Protocol != nil {
			groupMap["protocol"] = group.Protocol
		}

		groupMapList = append(groupMapList, groupMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("group_list", groupMapList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupMapList); e != nil {
			return e
		}
	}
	return nil
}
