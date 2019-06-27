package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudCnnInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCnnInstancesRead,

		Schema: map[string]*schema.Schema{
			"cnn_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			// Computed values
			"instance_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cnn_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qos": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCnnInstancesRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "data_source.tencentcloud_cnn_instances.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		cnnId string = ""
		name  string = ""
	)

	if temp, ok := d.GetOk("cnn_id"); ok {
		if tempStr := temp.(string); tempStr != "" {
			cnnId = tempStr
		}
	}

	if temp, ok := d.GetOk("name"); ok {
		if tempStr := temp.(string); tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeCcns(ctx, cnnId, name)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["cnn_id"] = item.cnnId
		infoMap["name"] = item.name
		infoMap["description"] = item.description
		infoMap["qos"] = item.qos
		infoMap["state"] = strings.ToUpper(item.state)
		infoMap["instance_count"] = item.instanceCount
		infoMap["create_time"] = item.createTime
		infoList = append(infoList, infoMap)
	}
	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  cnn instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId("cnn_instances" + cnnId + "_" + name)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
