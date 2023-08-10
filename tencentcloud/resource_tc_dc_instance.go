/*
Provides a resource to create a dc instance

Example Usage

```hcl
resource "tencentcloud_dc_instance" "instance" {
  access_point_id         = "ap-shenzhen-b-ft"
  bandwidth               = 10
  customer_contact_number = "0"
  direct_connect_name     = "terraform-for-test"
  line_operator           = "In-houseWiring"
  port_type               = "10GBase-LR"
  sign_law                = true
  vlan                    = -1
}
```

Import

dc instance can be imported using the id, e.g.

```
terraform import tencentcloud_dc_instance.instance dc_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcInstanceCreate,
		Read:   resourceTencentCloudDcInstanceRead,
		Update: resourceTencentCloudDcInstanceUpdate,
		Delete: resourceTencentCloudDcInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"direct_connect_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Connection name.",
			},

			"access_point_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Access point of connection.You can call `DescribeAccessPoints` to get the region ID. The selected access point must exist and be available.",
			},

			"line_operator": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ISP that provides connections. Valid values: ChinaTelecom (China Telecom), ChinaMobile (China Mobile), ChinaUnicom (China Unicom), In-houseWiring (in-house wiring), ChinaOther (other Chinese ISPs), InternationalOperator (international ISPs).",
			},

			"port_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Port type of connection. Valid values: 100Base-T (100-Megabit electrical Ethernet interface), 1000Base-T (1-Gigabit electrical Ethernet interface), 1000Base-LX (1-Gigabit single-module optical Ethernet interface; 10 KM), 10GBase-T (10-Gigabit electrical Ethernet interface), 10GBase-LR (10-Gigabit single-module optical Ethernet interface; 10 KM). Default value: 1000Base-LX.",
			},

			"circuit_code": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Circuit code of a connection, which is provided by the ISP or connection provider.",
			},

			"location": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Local IDC location.",
			},

			"bandwidth": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Connection port bandwidth in Mbps. Value range: [2,10240]. Default value: 1000.",
			},

			"redundant_direct_connect_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of redundant connection.",
			},

			"vlan": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "VLAN for connection debugging, which is enabled and automatically assigned by default.",
			},

			"tencent_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Tencent-side IP address for connection debugging, which is automatically assigned by default.",
			},

			"customer_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User-side IP address for connection debugging, which is automatically assigned by default.",
			},

			"customer_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Name of connection applicant, which is obtained from the account system by default.",
			},

			"customer_contact_mail": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Email address of connection applicant, which is obtained from the account system by default.",
			},

			"customer_contact_number": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Contact number of connection applicant, which is obtained from the account system by default.",
			},

			"fault_report_contact_person": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fault reporting contact person.",
			},

			"fault_report_contact_number": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fault reporting contact number.",
			},

			"sign_law": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether the connection applicant has signed the service agreement. Default value: true.",
			},
		},
	}
}

func resourceTencentCloudDcInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = dc.NewCreateDirectConnectRequest()
		response = dc.NewCreateDirectConnectResponse()
	)
	if v, ok := d.GetOk("direct_connect_name"); ok {
		request.DirectConnectName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("access_point_id"); ok {
		request.AccessPointId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("line_operator"); ok {
		request.LineOperator = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port_type"); ok {
		request.PortType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("circuit_code"); ok {
		request.CircuitCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("location"); ok {
		request.Location = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("bandwidth"); ok {
		request.Bandwidth = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("redundant_direct_connect_id"); ok {
		request.RedundantDirectConnectId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("vlan"); ok {
		request.Vlan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tencent_address"); ok {
		request.TencentAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("customer_address"); ok {
		request.CustomerAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("customer_name"); ok {
		request.CustomerName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("customer_contact_mail"); ok {
		request.CustomerContactMail = helper.String(v.(string))
	}

	if v, ok := d.GetOk("customer_contact_number"); ok {
		request.CustomerContactNumber = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fault_report_contact_person"); ok {
		request.FaultReportContactPerson = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fault_report_contact_number"); ok {
		request.FaultReportContactNumber = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sign_law"); ok {
		request.SignLaw = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().CreateDirectConnect(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dc instance failed, reason:%+v", logId, err)
		return err
	}

	dcSet := response.Response.DirectConnectIdSet
	if len(dcSet) < 1 {
		return fmt.Errorf("Create Direct Connect failed")
	}

	d.SetId(*dcSet[0])

	return resourceTencentCloudDcInstanceRead(d, meta)
}

func resourceTencentCloudDcInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}

	directConnectId := d.Id()

	instances, err := service.DescribeDirectConnects(ctx, directConnectId, "")
	if err != nil {
		return err
	}

	if len(instances) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	instance := instances[0]

	if instance.DirectConnectName != nil {
		_ = d.Set("direct_connect_name", instance.DirectConnectName)
	}

	if instance.AccessPointId != nil {
		_ = d.Set("access_point_id", instance.AccessPointId)
	}

	if instance.LineOperator != nil {
		_ = d.Set("line_operator", instance.LineOperator)
	}

	if instance.PortType != nil {
		_ = d.Set("port_type", instance.PortType)
	}

	if instance.CircuitCode != nil {
		_ = d.Set("circuit_code", instance.CircuitCode)
	}

	if instance.Location != nil {
		_ = d.Set("location", instance.Location)
	}

	if instance.Bandwidth != nil {
		_ = d.Set("bandwidth", instance.Bandwidth)
	}

	if instance.RedundantDirectConnectId != nil {
		_ = d.Set("redundant_direct_connect_id", instance.RedundantDirectConnectId)
	}

	if instance.Vlan != nil {
		_ = d.Set("vlan", instance.Vlan)
	}

	if instance.TencentAddress != nil {
		_ = d.Set("tencent_address", instance.TencentAddress)
	}

	if instance.CustomerAddress != nil {
		_ = d.Set("customer_address", instance.CustomerAddress)
	}

	if instance.CustomerName != nil {
		_ = d.Set("customer_name", instance.CustomerName)
	}

	if instance.CustomerContactMail != nil {
		_ = d.Set("customer_contact_mail", instance.CustomerContactMail)
	}

	if instance.CustomerContactNumber != nil {
		_ = d.Set("customer_contact_number", instance.CustomerContactNumber)
	}

	if instance.FaultReportContactPerson != nil {
		_ = d.Set("fault_report_contact_person", instance.FaultReportContactPerson)
	}

	if instance.FaultReportContactNumber != nil {
		_ = d.Set("fault_report_contact_number", instance.FaultReportContactNumber)
	}

	if instance.SignLaw != nil {
		_ = d.Set("sign_law", instance.SignLaw)
	}

	return nil
}

func resourceTencentCloudDcInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dc.NewModifyDirectConnectAttributeRequest()

	directConnectId := d.Id()

	request.DirectConnectId = &directConnectId
	needChange := false

	immutableArgs := []string{
		"access_point_id", "line_operator", "port_type",
		"bandwidth", "redundant_direct_connect_id",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	mutableArgs := []string{
		"direct_connect_name", "circuit_code", "location",
		"vlan", "tencent_address", "customer_address", "customer_name",
		"customer_contact_mail", "customer_contact_number", "fault_report_contact_person",
		"fault_report_contact_number", "sign_law",
	}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("direct_connect_name"); ok {
			request.DirectConnectName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("circuit_code"); ok {
			request.CircuitCode = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("bandwidth"); ok {
			request.Bandwidth = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("vlan"); ok {
			request.Vlan = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("tencent_address"); ok {
			request.TencentAddress = helper.String(v.(string))
		}

		if v, ok := d.GetOk("customer_address"); ok {
			request.CustomerAddress = helper.String(v.(string))
		}

		if v, ok := d.GetOk("customer_name"); ok {
			request.CustomerName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("customer_contact_mail"); ok {
			request.CustomerContactMail = helper.String(v.(string))
		}

		if v, ok := d.GetOk("customer_contact_number"); ok {
			request.CustomerContactNumber = helper.String(v.(string))
		}

		if v, ok := d.GetOk("fault_report_contact_person"); ok {
			request.FaultReportContactPerson = helper.String(v.(string))
		}

		if v, ok := d.GetOk("fault_report_contact_number"); ok {
			request.FaultReportContactNumber = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("sign_law"); ok {
			request.SignLaw = helper.Bool(v.(bool))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcClient().ModifyDirectConnectAttribute(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dc instance failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDcInstanceRead(d, meta)
}

func resourceTencentCloudDcInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dc_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcService{client: meta.(*TencentCloudClient).apiV3Conn}
	directConnectId := d.Id()

	if err := service.DeleteDcInstanceById(ctx, directConnectId); err != nil {
		return err
	}

	return nil
}
