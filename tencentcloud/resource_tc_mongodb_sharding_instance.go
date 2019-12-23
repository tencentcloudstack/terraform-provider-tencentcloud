/*
Provide a resource to create a Mongodb sharding instance.

Example Usage

```hcl
resource "tencentcloud_mongodb_sharding_instance" "mongodb" {
  instance_name   = "mongodb"
  shard_quantity  = 2
  nodes_per_shard = 3
  memory          = 4
  volume          = 100
  engine_version  = "MONGO_3_WT"
  machine_type    = "GIO"
  available_zone  = "ap-guangzhou-3"
  vpc_id          = "vpc-mz3efvbw"
  subnet_id       = "subnet-lk0svi3p"
  project_id      = 0
  password        = "mypassword"
}
```

Import

Mongodb sharding instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_mongodb_sharding_instance.mongodb cmgo-41s6jwy4
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20180408"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMongodbShardingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongodbShardingInstanceCreate,
		Read:   resourceMongodbShardingInstanceRead,
		Update: resourceMongodbShardingInstanceUpdate,
		Delete: resourceMongodbShardingInstanceDelete,
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
			"shard_quantity": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(2, 20),
				Description:  "Number of sharding.",
			},
			"nodes_per_shard": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(3, 5),
				Description:  "Number of nodes per shard, at least 3(one master and two slaves).",
			},
			"memory": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerMin(2),
				Description:  "Memory size. The minimum value is 2, and unit is GB.",
			},
			"volume": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerMin(25),
				Description:  "Disk size. The minimum value is 25, and unit is GB.",
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
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of the Mongodb.",
			},

			// Computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status of the Mongodb instance, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).",
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

func resourceMongodbShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_sharding_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	mongodbService := MongodbService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	request := mongodb.NewCreateDBInstanceHourRequest()
	request.ReplicateSetNum = helper.IntUint64(d.Get("shard_quantity").(int))
	request.SecondaryNum = helper.IntUint64(d.Get("nodes_per_shard").(int) - 1)
	request.GoodsNum = helper.IntUint64(1)
	request.InstanceRole = helper.String("MASTER")
	request.InstanceType = helper.String("SHARD")
	request.Memory = helper.IntUint64(d.Get("memory").(int))
	request.Volume = helper.IntUint64(d.Get("volume").(int))
	request.EngineVersion = helper.String(d.Get("engine_version").(string))
	request.Machine = helper.String(d.Get("machine_type").(string))
	request.Zone = helper.String(d.Get("available_zone").(string))
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}
	if v, ok := d.GetOk("security_groups"); ok {
		securityGroups := v.(*schema.Set).List()
		request.SecurityGroup = make([]*string, 0, len(securityGroups))
		for _, v := range securityGroups {
			request.SecurityGroup = append(request.SecurityGroup, helper.String(v.(string)))
		}
	}

	response, err := client.UseMongodbClient().CreateDBInstanceHour(request)
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
		log.Printf("[CRITAL]%s creating mongodb instance failed, reason:%s\n", logId, err.Error())
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
		log.Printf("[CRITAL]%s creating mongodb instance failed, reason:%s\n", logId, err.Error())
		return err
	}
	d.SetId(instanceId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := BuildTagResourceName("mongodb", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceMongodbShardingInstanceRead(d, meta)
}

func resourceMongodbShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_sharding_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := d.Id()

	mongodbService := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}
	instance, err := mongodbService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if nilFields := CheckNil(instance, map[string]string{
		"InstanceName":      "instance name",
		"ProjectId":         "project id",
		"ClusterType":       "cluster type",
		"Zone":              "available zone",
		"VpcId":             "vpc id",
		"SubnetId":          "subnet id",
		"Status":            "status",
		"Vip":               "vip",
		"Vport":             "vport",
		"CreateTime":        "create time",
		"MongoVersion":      "engine version",
		"Memory":            "memory",
		"Volume":            "volume",
		"MachineType":       "machine type",
		"ReplicationSetNum": "shard quantity",
		"SecondaryNum":      "secondary number",
	}); len(nilFields) > 0 {
		return fmt.Errorf("mongodb %v are nil", nilFields)
	}

	_ = d.Set("shard_quantity", instance.ReplicationSetNum)
	_ = d.Set("nodes_per_shard", *instance.SecondaryNum+1)
	_ = d.Set("instance_name", instance.InstanceName)
	_ = d.Set("memory", *instance.Memory/1024/(*instance.ReplicationSetNum))
	_ = d.Set("volume", *instance.Volume/1024/(*instance.ReplicationSetNum))
	_ = d.Set("engine_version", instance.MongoVersion)

	switch *instance.MachineType {
	case "HIO10G":
		_ = d.Set("machine_type", MONGODB_MACHINE_TYPE_TGIO)

	case "HIO":
		_ = d.Set("machine_type", MONGODB_MACHINE_TYPE_GIO)
	}

	_ = d.Set("available_zone", instance.Zone)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("status", instance.Status)
	_ = d.Set("vip", instance.Vip)
	_ = d.Set("vport", instance.Vport)
	_ = d.Set("create_time", instance.CreateTime)

	tags := make(map[string]string, len(instance.Tags))
	for _, tag := range instance.Tags {
		if tag.TagKey == nil {
			return errors.New("mongodb tag key is nil")
		}
		if tag.TagValue == nil {
			return errors.New("mongodb tag value is nil")
		}

		tags[*tag.TagKey] = *tag.TagValue
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceMongodbShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_sharding_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := d.Id()

	client := meta.(*TencentCloudClient).apiV3Conn
	mongodbService := MongodbService{client: client}
	tagService := TagService{client: client}
	region := client.Region

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
			log.Printf("[CRITAL]%s upgrade mongodb instance failed, reason:%s\n", logId, err.Error())
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
			log.Printf("[CRITAL]%s setting mongodb instance password failed, reason:%s\n", logId, err.Error())
			return err
		}

		d.SetPartial("password")
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("mongodb", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceMongodbShardingInstanceRead(d, meta)
}

func resourceMongodbShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mongodb_sharding_instance.delete")()

	logId := getLogId(contextNil)
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
