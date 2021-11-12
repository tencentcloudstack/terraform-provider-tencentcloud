/*
Provides a mysql instance resource to create master database instances.

~> **NOTE:** If this mysql has readonly instance, the terminate operation of the mysql does NOT take effect immediately, maybe takes for several hours. so during that time, VPCs associated with that mysql instance can't be terminated also.

Example Usage

```hcl
resource "tencentcloud_mysql_instance" "default" {
  internet_service = 1
  engine_version   = "5.7"
  charge_type = "POSTPAID"
  root_password     = "********"
  slave_deploy_mode = 0
  first_slave_zone  = "ap-guangzhou-4"
  second_slave_zone = "ap-guangzhou-4"
  slave_sync_mode   = 1
  availability_zone = "ap-guangzhou-4"
  project_id        = 201901010001
  instance_name     = "myTestMysql"
  mem_size          = 128000
  volume_size       = 250
  vpc_id            = "vpc-12mt3l31"
  subnet_id         = "subnet-9uivyb1g"
  intranet_port     = 3306
  security_groups   = ["sg-ot8eclwz"]

  tags = {
    name = "test"
  }

  parameters = {
    max_connections = "1000"
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudEmrInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrInstanceCreate,
		Read:   resourceTencentCloudEmrInstanceRead,
		Update: resourceTencentCloudEmrInstanceUpdate,
		Delete: resourceTencentCloudEmrInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 30),
				Description:  "The product id of EMR instance.",
			},
			"vpc_settings": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "The private net config of EMR instance.",
			},
			"softwares": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The softwares of a EMR instance.",
			},
			"resource_spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_resource_spec": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec":         {Type: schema.TypeString, Optional: true},
									"storage_type": {Type: schema.TypeInt, Optional: true},
									"disk_type":    {Type: schema.TypeString, Optional: true},
									"mem_size":     {Type: schema.TypeInt, Optional: true},
									"cpu":          {Type: schema.TypeInt, Optional: true},
									"disk_size":    {Type: schema.TypeInt, Optional: true},
									"root_size":    {Type: schema.TypeInt, Optional: true},
									"tags": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"multi_disks": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeMap},
									},
									"instance_type":  {Type: schema.TypeString, Optional: true},
									"local_disk_num": {Type: schema.TypeInt, Optional: true},
									"disk_num":       {Type: schema.TypeInt, Optional: true},
								},
							},
						},
						// "core_resource_spec": {
						// 	Type: schema.TypeMap,
						// },
						// "task_resource_spec": {
						// 	Type: schema.TypeMap,
						// },
						// "master_count": {
						// 	Type: schema.TypeInt,
						// },
						// "core_count": {
						// 	Type: schema.TypeInt,
						// },
						// "task_count": {
						// 	Type: schema.TypeInt,
						// },
						// "common_resource_spec": {
						// 	Type: schema.TypeMap,
						// },
						// "common_count": {
						// 	Type: schema.TypeInt,
						// },
					},
				},
				Description: "Resource specification of EMR instance.",
			},
			// "support_ha": {
			// 	Type:         schema.TypeInt,
			// 	Required:     true,
			// 	ValidateFunc: validateIntegerInRange(0, 1),
			// 	Description:  "The flag whether the instance support high availability.(0=>not support, 1=>support)",
			// },
			// "instance_name": {
			// 	Type:         schema.TypeString,
			// 	Required:     true,
			// 	ValidateFunc: validateStringLengthInRange(6, 36),
			// 	Description:  "Name of the instance, which can contain 6 to 36 English letters, Chinese characters, digits, dashes(-), or underscores(_).",
			// },
			// "pay_mode": {
			// 	Type:         schema.TypeInt,
			// 	Required:     true,
			// 	ValidateFunc: validateIntegerInRange(0, 1),
			// 	Description:  "The pay mode of instance. 0 is pay on an annual basis, 1 is pay on a measure basis.",
			// },
			// "placement": {
			// 	Type:        schema.TypeMap,
			// 	Required:    true,
			// 	Description: "The location of the instance.",
			// },
			// "time_span": {
			// 	Type:        schema.TypeInt,
			// 	Required:    true,
			// 	Description: "The length of time the instance was purchased. Use with TimeUnit.When TimeUnit is s, the parameter can only be filled in at 3600, representing a metered instance.\nWhen TimeUnit is m, the number filled in by this parameter indicates the length of purchase of the monthly instance of the package year, such as 1 for one month of purchase.",
			// },
			// "time_unit": {
			// 	Type:        schema.TypeString,
			// 	Required:    true,
			// 	Description: "The unit of time in which the instance was purchased. When PayMode is 0, TimeUnit can only take values of s(second). When PayMode is 1, TimeUnit can only take the value m(month)",
			// },
			// "login_settings": {
			// 	Type:        schema.TypeMap,
			// 	Required:    true,
			// 	Description: "Instance login settings.",
			// },
			// "cos_settings": {
			// 	Type:        schema.TypeMap,
			// 	Optional:    true,
			// 	Description: "The parameters that need to be set to turn on COS access.",
			// },
			// "sg_id": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The ID of the security group to which the instance belongs",
			// },
		},
	}
}

func resourceTencentCloudEmrInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_emr_instance.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	emrService := EMRService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	if err := emrService.CreateInstance(ctx, d); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudEmrInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTencentCloudEmrInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceTencentCloudEmrInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
