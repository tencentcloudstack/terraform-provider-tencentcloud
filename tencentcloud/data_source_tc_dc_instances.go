package tencentcloud

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func dataSourceTencentCloudDcInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcInstancesRead,

		Schema: map[string]*schema.Schema{
			"dc_id": {
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
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_point_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"circuit_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"redundant_dc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tencent_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fault_report_contact_person": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fault_report_contact_phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDcInstancesRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)

	defer LogElapsed(logId + "data_source.tencentcloud_dc_instances.read")()

	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		id   = ""
		name = ""
	)
	if temp, ok := d.GetOk("dc_id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			id = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeDirectConnects(ctx, id, name)

	if err != nil {
		return err
	}
	var instanceList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {

		var infoMap = make(map[string]interface{})
		infoMap["dc_id"] = *item.DirectConnectId
		infoMap["name"] = *item.DirectConnectName

		infoMap["state"] = strings.ToUpper(service.strPt2str(item.State))
		infoMap["access_point_id"] = service.strPt2str(item.AccessPointId)
		infoMap["line_operator"] = service.strPt2str(item.LineOperator)

		infoMap["location"] = service.strPt2str(item.Location)
		infoMap["bandwidth"] = service.int64Pt2int64(item.Bandwidth)
		infoMap["port_type"] = service.strPt2str(item.PortType)

		infoMap["circuit_code"] = service.strPt2str(item.CircuitCode)
		infoMap["redundant_dc_id"] = service.strPt2str(item.RedundantDirectConnectId)
		infoMap["tencent_address"] = service.strPt2str(item.TencentAddress)

		infoMap["customer_address"] = service.strPt2str(item.CustomerAddress)
		infoMap["customer_name"] = service.strPt2str(item.CustomerName)
		infoMap["customer_email"] = service.strPt2str(item.CustomerContactMail)

		infoMap["customer_phone"] = service.strPt2str(item.CustomerContactNumber)
		infoMap["fault_report_contact_person"] = service.strPt2str(item.FaultReportContactPerson)
		infoMap["fault_report_contact_phone"] = service.strPt2str(item.FaultReportContactNumber)

		infoMap["enabled_time"] = service.strPt2str(item.EnabledTime)
		infoMap["create_time"] = service.strPt2str(item.CreatedTime)
		infoMap["expired_time"] = service.strPt2str(item.ExpiredTime)

		instanceList = append(instanceList, infoMap)
	}

	if err := d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set  dc instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	m := md5.New()
	m.Write([]byte(id + "_" + name))
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
