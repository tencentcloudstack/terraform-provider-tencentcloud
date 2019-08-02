/*
Provide a resource to create a Mongodb instance.

Example Usage

```hcl
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name     = "mongodb"
  memory            = 4
  volume            = 100
  engine_version    = "MONGO_3_WT"
  machine_type      = "GIO"
  available_zone    = "ap-guangzhou-2"
  vpc_id            = "vpc-mz3efvbw"
  subnet_id         = "subnet-lk0svi3p"
  project_id        = 0
  password          = "mypassword"
}
```

Import

Mongodb instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_instance.mongodb cmgo-41s6jwy4
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20180408"
)

func resourceTencentCloudMongodbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceCreate,
		Read:   resourceTencentCloudMongodbInstanceRead,
		Update: resourceTencentCloudMongodbInstanceUpdate,
		Delete: resourceTencentCloudMongodbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
				Description:  "Name of the Mongodb instance.",
			},
			"memory": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerMin(4),
				Description:  "Memory size.",
			},
			"volume": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerMin(100),
				Description:  "Disk size.",
			},
			"engine_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(MONGODB_ENGINE_VERSION),
				Description:  "Version of the Mongodb, and available values include MONGO_3_WT, MONGO_3_ROCKS and MONGO_36_WT.",
			},
			"machine_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(MONGODB_MACHINE_TYPE),
				Description:  "Type of Mongodb instance, and available values include GIO and TGIO.",
			},
			"available_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone of the Mongodb.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Description: "ID of the VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the subnet within this VPC. The vaule is required if VpcId is set.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the project which the instance belongs.",
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
				Description: "ID of the security group.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password of this Mongodb account.",
			},

			// Computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2)",
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP of the Mongodb instance.",
			},
			"vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IP port of the Mongodb instance.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the Mongodb instance.",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	mongodbService := MongodbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := mongodb.NewCreateDBInstanceHourRequest()
	request.ReplicateSetNum = intToPointer(1)
	request.SecondaryNum = intToPointer(2)
	request.GoodsNum = intToPointer(1)
	request.InstanceRole = stringToPointer("MASTER")
	request.InstanceType = stringToPointer("REPLSET")
	request.Memory = intToPointer(d.Get("memory").(int))
	request.Volume = intToPointer(d.Get("volume").(int))
	request.EngineVersion = stringToPointer(d.Get("engine_version").(string))
	request.Machine = stringToPointer(d.Get("machine_type").(string))
	request.Zone = stringToPointer(d.Get("available_zone").(string))
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}
	if v, ok := d.GetOk("security_groups"); ok {
		securityGroups := v.(*schema.Set).List()
		request.SecurityGroup = make([]*string, 0, len(securityGroups))
		for _, v := range securityGroups {
			request.SecurityGroup = append(request.SecurityGroup, stringToPointer(v.(string)))
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseMongodbClient().CreateDBInstanceHour(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceIds) < 1 {
		return fmt.Errorf("mongodb instance id is nil")
	}
	instanceId := *response.Response.InstanceIds[0]

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		instance, e := mongodbService.DescribeInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if *instance.Status == MONGODB_INSTANCE_STATUS_INITIAL {
			return nil
		}
		if *instance.Status == MONGODB_INSTANCE_STATUS_PROCESSING {
			return resource.RetryableError(fmt.Errorf("mongodb instance status is processing"))
		}
		e = fmt.Errorf("mongodb instance status is %d, we won't wait for it finish.", *instance.Status)
		return resource.NonRetryableError(e)
	})
	if err != nil {
		log.Printf("[CRITAL]%s creating mongodb instance failed, reason:%s\n ", logId, err.Error())
		return err
	}

	// setting instance name
	instanceName := d.Get("instance_name").(string)
	err = mongodbService.ModifyInstanceName(ctx, instanceId, instanceName)
	if err != nil {
		return err
	}

	// init instance(setting password of mongouser)
	password := d.Get("password").(string)
	err = mongodbService.SetInstancePassword(ctx, instanceId, "mongouser", password)
	if err != nil {
		return err
	}

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		instance, e := mongodbService.DescribeInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if *instance.Status == MONGODB_INSTANCE_STATUS_RUNNING {
			return nil
		}
		if *instance.Status == MONGODB_INSTANCE_STATUS_PROCESSING {
			return resource.RetryableError(fmt.Errorf("mongodb instance status is processing"))
		}
		e = fmt.Errorf("mongodb instance status is %d, we won't wait for it finish.", *instance.Status)
		return resource.NonRetryableError(e)
	})
	if err != nil {
		log.Printf("[CRITAL]%s creating mongodb instance failed, reason:%s\n ", logId, err.Error())
		return err
	}
	d.SetId(instanceId)

	return resourceTencentCloudMongodbInstanceRead(d, meta)
}

func resourceTencentCloudMongodbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := d.Id()
	mongodbService := MongodbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instance, err := mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	d.Set("instance_name", instance.InstanceName)
	d.Set("memory", *instance.Memory/1024)
	d.Set("volume", *instance.Volume/1024)
	d.Set("engine_version", instance.MongoVersion)

	switch *instance.MachineType {
	case "HIO10G":
		d.Set("machine_type", MONGODB_MACHINE_TYPE_TGIO)

	case "HIO":
		d.Set("machine_type", MONGODB_MACHINE_TYPE_GIO)
	}

	d.Set("available_zone", instance.Zone)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("project_id", instance.ProjectId)
	d.Set("status", instance.Status)
	d.Set("vip", instance.Vip)
	d.Set("vport", instance.Vport)
	d.Set("create_time", instance.CreateTime)

	return nil
}

func resourceTencentCloudMongodbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := d.Id()
	mongodbService := MongodbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	d.Partial(true)

	if d.HasChange("memory") || d.HasChange("volume") {
		memory := d.Get("memory").(int)
		volume := d.Get("volume").(int)
		err := mongodbService.UpgradeInstance(ctx, instanceId, memory, volume)
		if err != nil {
			return err
		}

		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			instance, e := mongodbService.DescribeInstanceById(ctx, instanceId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if *instance.Status == MONGODB_INSTANCE_STATUS_RUNNING {
				return nil
			}
			if *instance.Status == MONGODB_INSTANCE_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("mongodb instance status is processing"))
			}
			e = fmt.Errorf("mongodb instance status is %d, we won't wait for it finish.", *instance.Status)
			return resource.NonRetryableError(e)
		})
		if err != nil {
			log.Printf("[CRITAL]%s upgrade mongodb instance failed, reason:%s\n ", logId, err.Error())
			return err
		}

		if d.HasChange("memory") {
			d.SetPartial("memory")
		}
		if d.HasChange("volume") {
			d.SetPartial("volume")
		}
	}

	if d.HasChange("instance_name") {
		instanceName := d.Get("instance_name").(string)
		err := mongodbService.ModifyInstanceName(ctx, instanceId, instanceName)
		if err != nil {
			return err
		}
		d.SetPartial("instance_name")
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := mongodbService.ModifyProjectId(ctx, instanceId, projectId)
		if err != nil {
			return err
		}
		d.SetPartial("project_id")
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		err := mongodbService.SetInstancePassword(ctx, instanceId, "mongouser", password)
		if err != nil {
			return err
		}

		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			instance, e := mongodbService.DescribeInstanceById(ctx, instanceId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if *instance.Status == MONGODB_INSTANCE_STATUS_RUNNING {
				return nil
			}
			if *instance.Status == MONGODB_INSTANCE_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("mongodb instance status is processing"))
			}
			e = fmt.Errorf("mongodb instance status is %d, we won't wait for it finish.", *instance.Status)
			return resource.NonRetryableError(e)
		})
		if err != nil {
			log.Printf("[CRITAL]%s setting mongodb instance password failed, reason:%s\n ", logId, err.Error())
			return err
		}

		d.SetPartial("password")
	}

	d.Partial(false)

	return nil
}

func resourceTencentCloudMongodbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := d.Id()
	mongodbService := MongodbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := mongodbService.DeleteInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	return nil
}
