package tencentcloud

/*
resource tencentcloud_ccn main{
	name ="ci-temp-test-ccn"
	description="ci-temp-test-ccn-des"
	qos ="AG"
}
data tencentcloud_ccn_instances test{
	ccn_id = "${tencentcloud_ccn.main.id}"
}
*/
import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudCcnInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnInstancesRead,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
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
						"ccn_id": {
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
						"attachment_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"attached_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr_block": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
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

func dataSourceTencentCloudCcnInstancesRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "data_source.tencentcloud_ccn_instances.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		ccnId string = ""
		name  string = ""
	)

	if temp, ok := d.GetOk("ccn_id"); ok {
		if tempStr := temp.(string); tempStr != "" {
			ccnId = tempStr
		}
	}

	if temp, ok := d.GetOk("name"); ok {
		if tempStr := temp.(string); tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeCcns(ctx, ccnId, name)
	if err != nil {
		return err
	}

	var infoList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {
		var infoMap = make(map[string]interface{})
		infoMap["ccn_id"] = item.ccnId
		infoMap["name"] = item.name
		infoMap["description"] = item.description
		infoMap["qos"] = item.qos
		infoMap["state"] = strings.ToUpper(item.state)
		infoMap["create_time"] = item.createTime
		infoList = append(infoList, infoMap)

		instances, err := service.DescribeCcnAttachedInstances(ctx, item.ccnId)
		if err != nil {
			return err
		}
		attachmentList := make([]interface{}, 0, len(instances))

		for _, instance := range instances {

			instanceMap := map[string]interface{}{
				"instance_type":   instance.instanceType,
				"instance_region": instance.instanceRegion,
				"instance_id":     instance.instanceId,
				"state":           strings.ToUpper(instance.state),
				"attached_time":   instance.attachedTime,
				"cidr_block":      instance.cidrBlock,
			}
			attachmentList = append(attachmentList, instanceMap)

		}

		infoMap["attachment_list"] = attachmentList

	}
	if err := d.Set("instance_list", infoList); err != nil {
		log.Printf("[CRITAL]%s provider set  ccn instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	m := md5.New()
	m.Write([]byte("ccn_instances" + ccnId + "_" + name))
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), infoList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
