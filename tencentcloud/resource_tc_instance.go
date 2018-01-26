package tencentcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	tencentCloudApiInstanceChargeTypePrePaid        = "PREPAID"
	tencentCloudApiInstanceChargeTypePostPaidByHour = "POSTPAID_BY_HOUR"
)

const (
	tencentCloudApiInternetChargeTypeBandwithPrepaid         = "BANDWIDTH_PREPAID"
	tencentCloudApiInternetChargeTypeTrafficPostpaidByHour   = "TRAFFIC_POSTPAID_BY_HOUR"
	tencentCloudApiInternetChargeTypeBandwidthPostpaidByHour = "BANDWIDTH_POSTPAID_BY_HOUR"
	tencentCloudApiInternetChargeTypeBandwidthPackage        = "BANDWIDTH_PACKAGE"
)

const (
	tencentCloudApiInstanceChargeTypePrePaidRenewFlagNotifyAndAutoRenew          = "NOTIFY_AND_AUTO_RENEW"
	tencentCloudApiInstanceChargeTypePrePaidRenewFlagNotifyAndManualRenew        = "NOTIFY_AND_MANUAL_RENEW"
	tencentCloudApiInstanceChargeTypePrePaidRenewFlagDisableNotifyAndManualRenew = "DISABLE_NOTIFY_AND_MANUAL_RENEW"
)

const (
	tencentCloudApiDiskTypeLocalBaisc = "LOCAL_BASIC"
	tencentCloudApiDiskTypeLocalSSD   = "LOCAL_SSD"
	tencentCloudApiDiskTypeCloudBasic = "CLOUD_BASIC"
	tencentCloudApiDiskTypeCloudSSD   = "CLOUD_SSD"
)

var (
	availableInstanceChargeTypes = []string{
		tencentCloudApiInstanceChargeTypePrePaid,
		tencentCloudApiInstanceChargeTypePostPaidByHour,
	}
	availableInternetChargeTypes = []string{
		tencentCloudApiInternetChargeTypeBandwithPrepaid,
		tencentCloudApiInternetChargeTypeTrafficPostpaidByHour,
		tencentCloudApiInternetChargeTypeBandwidthPostpaidByHour,
		tencentCloudApiInternetChargeTypeBandwidthPackage,
	}
	availableDiskTypes = []string{
		tencentCloudApiDiskTypeLocalBaisc,
		tencentCloudApiDiskTypeLocalSSD,
		tencentCloudApiDiskTypeCloudBasic,
		tencentCloudApiDiskTypeCloudSSD,
	}
	availableInstanceChargeTypePrePaidPeriodValues    = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}
	availableInstanceChargeTypePrePaidRenewFlagValues = []string{
		tencentCloudApiInstanceChargeTypePrePaidRenewFlagNotifyAndAutoRenew,
		tencentCloudApiInstanceChargeTypePrePaidRenewFlagNotifyAndManualRenew,
		tencentCloudApiInstanceChargeTypePrePaidRenewFlagDisableNotifyAndManualRenew,
	}
)

var (
	// TODO remove me when related feature implemented
	unsupportedUpdateFields = []string{
		"instance_charge_type_prepaid_period",
		"instance_charge_type_prepaid_renew_flag",
		"internet_charge_type",
		"internet_max_bandwidth_out",
		"allocate_public_ip",
		"system_disk_size",
		"data_disks",
	}
)

func resourceTencentCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudInstanceCreate,
		Read:   resourceTencentCloudInstanceRead,
		Update: resourceTencentCloudInstanceUpdate,
		Delete: resourceTencentCloudInstanceDelete,

		Schema: map[string]*schema.Schema{
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "CVM-Instance",
				ValidateFunc: validateInstanceName,
			},
			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceType,
			},
			// payment
			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceChargeType,
			},
			"instance_charge_type_prepaid_period": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateInstanceChargeTypePrePaidPeriod,
			},
			"instance_charge_type_prepaid_renew_flag": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceChargeTypePrePaidRenewFlag,
			},
			// network
			"internet_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetChargeType,
			},
			"internet_max_bandwidth_out": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateInternetMaxBandwidthOut,
			},
			"allocate_public_ip": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// vpc
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			// security group
			"security_groups": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			// storage
			"system_disk_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskType,
			},
			// TODO finish me when integrating CBS
			//"system_disk_id": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
			"system_disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(50, 1000),
			},
			"data_disks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateDiskType,
						},
						// TODO finish me when integrating CBS
						//"data_disk_id": &schema.Schema{
						//	Type:     schema.TypeString,
						//	Optional: true,
						//},
						"data_disk_size": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateDiskSize,
						},
					},
				},
			},
			// enhance services
			"disable_security_service": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_monitor_service": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			// login
			"key_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			// Computed values.
			"instance_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudInstanceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*TencentCloudClient).commonConn

	params := map[string]string{
		"Version":        "2017-03-12",
		"Action":         "RunInstances",
		"Placement.Zone": d.Get("availability_zone").(string),
		"ImageId":        d.Get("image_id").(string),
	}

	if instanceType, ok := d.GetOk("instance_type"); ok {
		insType := instanceType.(string)
		if len(insType) > 0 {
			params["InstanceType"] = insType
		}
	}
	if instanceName, ok := d.GetOk("instance_name"); ok {
		insName := instanceName.(string)
		if len(insName) > 0 {
			params["InstanceName"] = insName
		}
	}

	if instanceChargeType, ok := d.GetOk("instance_charge_type"); ok {
		insChargeType := instanceChargeType.(string)
		if insChargeType == tencentCloudApiInstanceChargeTypePrePaid {
			period, ok := d.GetOk("instance_charge_type_prepaid_period")
			if !ok {
				return fmt.Errorf(
					"tencentcloud_instance instance_charge_type_prepaid_period is need when instance_charge_type is %v",
					tencentCloudApiInstanceChargeTypePrePaid,
				)
			}
			periodStr := fmt.Sprintf("%v", period.(int))
			params["InstanceChargePrepaid.Period"] = periodStr
			if renewFlag, ok := d.GetOk("instance_charge_type_prepaid_renew_flag"); ok {
				params["InstanceChargePrepaid.RenewFlag"] = renewFlag.(string)
			}
		}

		params["InstanceChargeType"] = insChargeType
	}

	// network releated
	if internetChargeType, ok := d.GetOk("internet_charge_type"); ok {
		netChargeType := internetChargeType.(string)
		params["InternetAccessible.InternetChargeType"] = netChargeType
	}
	if internetMaxBandwidthOut, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		maxBandwidthOut := internetMaxBandwidthOut.(int)
		params["InternetAccessible.InternetMaxBandwidthOut"] = fmt.Sprintf("%v", maxBandwidthOut)
	}
	if allocatePublicIP, ok := d.GetOk("allocate_public_ip"); ok {
		if allocatePublicIP.(bool) {
			params["InternetAccessible.PublicIpAssigned"] = "TRUE"
		} else {
			params["InternetAccessible.PublicIpAssigned"] = "FALSE"
		}
	}

	// security groups
	if v, ok := d.GetOk("security_groups"); ok {
		sgIds := v.(*schema.Set).List()
		if len(sgIds) > 0 {
			for i, sgId := range sgIds {
				paramKey := fmt.Sprintf("SecurityGroupIds.%v", i)
				params[paramKey] = sgId.(string)
			}
		}
	}

	// storage
	if systemDiskType, ok := d.GetOk("system_disk_type"); ok {
		params["SystemDisk.DiskType"] = systemDiskType.(string)
	}
	if systemDiskSize, ok := d.GetOk("system_disk_size"); ok {
		diskSize := systemDiskSize.(int)
		params["SystemDisk.DiskSize"] = fmt.Sprintf("%v", diskSize)
	}
	var dataDisksAttr []map[string]interface{}
	if dataDisks, ok := d.GetOk("data_disks"); ok {
		dataDiskList := dataDisks.([]interface{})
		if len(dataDiskList) > 1 {
			return fmt.Errorf("tencentcloud_instance currently only one data disk is supported during instance creation")
		}
		for i, dataDisk := range dataDiskList {
			dd := dataDisk.(map[string]interface{})
			if v, ok := dd["data_disk_type"].(string); ok && v != "" {
				paramKey := fmt.Sprintf("DataDisks.%v.DiskType", i)
				params[paramKey] = v
			}
			if v, ok := dd["data_disk_size"].(int); ok {
				paramKey := fmt.Sprintf("DataDisks.%v.DiskSize", i)
				paramValue := fmt.Sprintf("%v", v)
				params[paramKey] = paramValue
			}
			dataDisksAttr = append(dataDisksAttr, dd)
		}
	}

	// enhance services
	if v, ok := d.GetOk("disable_security_service"); ok {
		disable := v.(bool)
		if disable {
			params["EnhancedService.SecurityService.Enabled"] = "FALSE"
		}
	}
	if v, ok := d.GetOk("disable_monitor_service"); ok {
		disable := v.(bool)
		if disable {
			params["EnhancedService.MonitorService.Enabled"] = "FALSE"
		}
	}

	// login confidential
	if v, ok := d.GetOk("key_name"); ok {
		keyId := v.(string)
		params["LoginSettings.KeyIds.0"] = keyId
	}
	if v, ok := d.GetOk("password"); ok {
		passwd := v.(string)
		params["LoginSettings.Password"] = passwd
	}

	// vpc
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId := v.(string)
		params["VirtualPrivateCloud.VpcId"] = vpcId
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId := v.(string)
		params["VirtualPrivateCloud.SubnetId"] = subnetId
	}

	response, err := client.SendRequest("cvm", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:tag"Code"`
				Message string `json:tag"Message"`
			}
			InstanceIdSet []string
			RequestId     string
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Response.Error.Code != "" {
		return fmt.Errorf(
			"tencentcloud_instance got error, code:%v, message:%v",
			jsonresp.Response.Error.Code,
			jsonresp.Response.Error.Message,
		)
	}
	if len(jsonresp.Response.InstanceIdSet) == 0 {
		return fmt.Errorf("tencentcloud_instance no instance id returned")
	}

	var instanceStatusMap map[string]string
	instanceStatusMap, err = waitInstanceReachTargetStatus(client, jsonresp.Response.InstanceIdSet, "RUNNING")
	if err != nil {
		return err
	}

	instanceId := jsonresp.Response.InstanceIdSet[0]
	d.SetId(instanceId)
	d.Set("instance_status", instanceStatusMap[instanceId])
	d.Set("data_disks", dataDisksAttr)
	return resourceTencentCloudInstanceRead(d, m)
}

func resourceTencentCloudInstanceRead(d *schema.ResourceData, m interface{}) error {
	instanceId := d.Id()
	params := map[string]string{
		"Version": "2017-03-12",
		"Action":  "DescribeInstances",
	}
	params["InstanceIds.0"] = instanceId

	client := m.(*TencentCloudClient).commonConn
	response, err := client.SendRequest("cvm", params)
	if err != nil {
		return err
	}
	var jsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:tag"Code"`
				Message string `json:tag"Message"`
			}
			TotalCount  int
			InstanceSet []struct {
				Placement struct {
					Zone      string   `json:tag"Zone"`
					ProjectId int      `json:tag"ProjectId"`
					HostIds   []string `json:tag"HostIds"`
				} `json:tag"Placement"`

				InstanceId   string `json:tag"InstanceId"`
				InstanceType string `json:tag"InstanceType"`
				CPU          int    `json:tag"CPU"`
				Memory       int    `json:tag"Memory"`

				RestrictState      string `json:tag"RestrictState"`
				InstanceName       string `json:tag"InstanceName"`
				InstanceChargeType string `json:tag"InstanceChargeType"`

				SystemDisk struct {
					DiskType string `json:tag"DiskType"`
					DiskId   string `json:tag"DiskId"`
					DiskSize int    `json:tag"DiskSize"`
				} `json:tag"SystemDisk"`

				DataDisks []struct {
					DiskType string `json:tag"DiskType"`
					DiskId   string `json:tag"DiskId"`
					DiskSize int    `json:tag"DiskSize"`
				} `json:tag"DataDisks"`

				PrivateIpAddresses []string `json:tag"PrivateIpAddresses"`
				PublicIpAddresses  []string `json:tag"PublicIpAddresses"`

				InternetAccessible struct {
					InternetChargeType      string `json:tag"InternetChargeType"`
					InternetMaxBandwidthOut int    `json:tag"InternetMaxBandwidthOut"`
					PublicIpAssigned        bool   `json:tag"PublicIpAssigned"`
				} `json:tag"InternetAccessible"`

				VirtualPrivateCloud struct {
					VpcId              string   `json:tag"VpcId"`
					SubnetId           string   `json:tag"SubnetId"`
					AsVpcGateway       bool     `json:tag"AsVpcGateway"`
					PrivateIpAddresses []string `json:tag"PrivateIpAddresses"`
				}

				SecurityGroupIds []string `json:tag"SecurityGroupIds"`

				LoginSettings struct {
					KeyIds []string `json:tag"KeyIds"`
				} `json:tag"LoginSettings"`

				ImageId string `json:tag"ImageId"`

				RenewFlag   string    `json:tag"RenewFlag"`
				CreatedTime time.Time `json:tag"CreatedTime"`
				ExpiredTime time.Time `json:tag"ExpiredTime"`
			} `json:tag"InstanceSet"`
			RequestId string
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Response.Error.Code != "" {
		return fmt.Errorf(
			"tencentcloud_instance got error, code:%v, message:%v",
			jsonresp.Response.Error.Code,
			jsonresp.Response.Error.Message,
		)
	}
	if len(jsonresp.Response.InstanceSet) == 0 {
		d.SetId("")
		return nil
	}
	privateIPs := jsonresp.Response.InstanceSet[0].PrivateIpAddresses
	if len(privateIPs) > 0 {
		d.Set("private_ip", privateIPs[0])
	}
	publicIPs := jsonresp.Response.InstanceSet[0].PublicIpAddresses
	if len(publicIPs) > 0 {
		d.Set("public_ip", publicIPs[0])
	}
	systemDiskType := jsonresp.Response.InstanceSet[0].SystemDisk.DiskType
	d.Set("system_disk_type", systemDiskType)
	systemDiskSize := jsonresp.Response.InstanceSet[0].SystemDisk.DiskSize
	d.Set("system_disk_size", systemDiskSize)

	var dataDiskList []map[string]interface{}
	dataDisks := jsonresp.Response.InstanceSet[0].DataDisks
	for _, dataDisk := range dataDisks {
		m := make(map[string]interface{})
		diskType := dataDisk.DiskType
		diskSize := dataDisk.DiskSize
		m["data_disk_type"] = diskType
		m["data_disk_size"] = diskSize
		dataDiskList = append(dataDiskList, m)
	}
	//d.Set("data_disks", dataDiskList)

	securityGroupIds := jsonresp.Response.InstanceSet[0].SecurityGroupIds
	if len(securityGroupIds) > 0 {
		d.Set("security_groups", securityGroupIds)
	}

	loginSettings := jsonresp.Response.InstanceSet[0].LoginSettings
	keyIds := loginSettings.KeyIds
	if len(keyIds) > 0 {
		keyName := keyIds[0]
		d.Set("key_name", keyName)
	}
	if v, ok := d.GetOk("password"); ok {
		passwd := v.(string)
		d.Set("password", passwd)
	}

	virtualPrivateCloud := jsonresp.Response.InstanceSet[0].VirtualPrivateCloud
	vpcId := virtualPrivateCloud.VpcId
	if len(vpcId) > 0 {
		d.Set("vpc_id", vpcId)
	}
	subnetId := virtualPrivateCloud.SubnetId
	if len(subnetId) > 0 {
		d.Set("subnet_id", subnetId)
	}

	return nil
}

func resourceTencentCloudInstanceUpdate(d *schema.ResourceData, m interface{}) (err error) {
	client := m.(*TencentCloudClient).commonConn
	instanceId := d.Id()

	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_instance update on %v is not supported yet", field)
		}
	}

	d.Partial(true)

	var isUpdatingImage bool

	if d.HasChange("image_id") {
		isUpdatingImage = true
	}

	if d.HasChange("instance_name") {
		d.SetPartial("instance_name")
		oldInstanceName, newInstanceName := d.GetChange("instance_name")
		log.Printf("[DEBUG] tencentcloud_instance rename instance_name from %v to %v", oldInstanceName, newInstanceName)
		err = renameInstancesName(client, []string{instanceId}, newInstanceName.(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("key_name") {
		d.SetPartial("key_name")
		if isUpdatingImage {
			goto LABEL_REINSTALL
		}
		oldKey, newKey := d.GetChange("key_name")
		log.Printf("[DEBUG] tencentcloud_instance rebind key pair, old key: %v, new key: %v", oldKey, newKey)

		_, err := waitInstanceReachOneOfTargetStatusList(
			client,
			[]string{instanceId},
			[]string{
				"STOPPED",
				"RUNNING",
			},
		)
		if err != nil {
			return err
		}

		err = bindKeyPiar(client, instanceId, newKey.(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("password") {
		d.SetPartial("password")
		if isUpdatingImage {
			goto LABEL_REINSTALL
		}
		_, newValue := d.GetChange("password")
		log.Printf("[DEBUG] tencentcloud_instance reset password\n")
		err = resetInstancePassword(client, instanceId, newValue.(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("security_groups") {
		_, n := d.GetChange("security_groups")
		ns := n.(*schema.Set)

		sgIds := expandStringList(ns.List())
		if len(sgIds) == 0 {
			return fmt.Errorf("tencentcloud_instance security groups are not allow to be empty")
		}

		err = bindInstanceWithSgIds(client, d.Id(), sgIds)
		if err != nil {
			return err
		}
		d.SetPartial("security_groups")
	}

LABEL_REINSTALL:
	if d.HasChange("image_id") {
		d.SetPartial("image_id")
		oldValue, newValue := d.GetChange("image_id")
		log.Printf("[DEBUG] tencentcloud_instance reinstall image from %v to %v\n", oldValue, newValue)

		err = resetInstanceSystem(client, d, instanceId, newValue.(string))
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudInstanceRead(d, m)
}

func resourceTencentCloudInstanceDelete(d *schema.ResourceData, m interface{}) error {
	v, ok := d.GetOk("instance_charge_type")
	if ok {
		instanceChargeType := v.(string)
		if instanceChargeType == tencentCloudApiInstanceChargeTypePrePaid {
			return fmt.Errorf("prepaid instance %v is not allowed to be deleted in terraform", d.Id())
		}
	}

	client := m.(*TencentCloudClient).commonConn

	params := map[string]string{
		"Version":       "2017-03-12",
		"Action":        "TerminateInstances",
		"InstanceIds.0": d.Id(),
	}

	resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err := client.SendRequest("cvm", params)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		var jsonresp struct {
			Response struct {
				Error struct {
					Code    string `json:tag"Code"`
					Message string `json:tag"Message"`
				}
				RequestId string
			}
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if jsonresp.Response.Error.Code == "InternalError" {
			return resource.RetryableError(fmt.Errorf(jsonresp.Response.Error.Message))
		}
		return nil
	})
	d.SetId("")
	return nil
}
