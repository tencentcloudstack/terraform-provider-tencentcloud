/*
Use this data source to query detailed information of DC instances.

Example Usage

```hcl
data "tencentcloud_dc_instances" "name_select" {
  name = "t"
}

data "tencentcloud_dc_instances" "id" {
  dcx_id = "dc-kax48sg7"
}
```
*/
package tencentcloud

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTencentCloudDcInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcInstancesRead,

		Schema: map[string]*schema.Schema{
			"dc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the DC to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the DC to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the DC.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the DC.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the DC.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the DC, and available values include `REJECTED`, `TOPAY`, `PAID`, `ALLOCATED`, `AVAILABLE`, `DELETING` and `DELETED`.",
						},
						"access_point_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access point ID of tne DC.",
						},
						"line_operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operator of the DC, and available values include `ChinaTelecom`, `ChinaMobile`, `ChinaUnicom`, `In-houseWiring`, `ChinaOther` and `InternationalOperator`.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The DC location where the connection is located.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bandwidth of the DC.",
						},
						"port_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port type of the DC in client, and available values include `100Base-T`, `1000Base-T`, `1000Base-LX`, `10GBase-T` and `10GBase-LR`. The default value is `1000Base-LX`.",
						},
						"circuit_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The circuit code provided by the operator for the DC.",
						},
						"redundant_dc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the redundant DC.",
						},
						"tencent_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect IP of the DC within Tencent. Note: This field may return null, indicating that no valid values are taken.",
						},
						"customer_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect IP of the DC within client. Note: This field may return null, indicating that no valid values are taken.",
						},
						"customer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Applicant name of the DC, the default is obtained from the account. Note: This field may return null, indicating that no valid values are taken.",
						},
						"customer_email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Applicant email of the DC, the default is obtained from the account. Note: This field may return null, indicating that no valid values are taken.",
						},
						"customer_phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Applicant phone number of the DC, the default is obtained from the account. Note: This field may return null, indicating that no valid values are taken.",
						},
						"fault_report_contact_person": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contact of reporting a faulty. Note: This field may return null, indicating that no valid values are taken.",
						},
						"fault_report_contact_phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Phone number of reporting a faulty. Note: This field may return null, indicating that no valid values are taken.",
						},
						"enabled_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enable time of resource.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of resource.",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expire date of resource.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDcInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dc_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

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
	_, err = m.Write([]byte(id + "_" + name))
	if err != nil {
		return err
	}
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
