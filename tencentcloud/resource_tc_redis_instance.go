package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudRedisInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisInstanceCreate,
		Read:   resourceTencentCloudMysqlInstanceRead,
		Update: resourceTencentCloudMysqlInstanceUpdate,
		Delete: resourceTencentCloudMysqlInstanceDelete,

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  REDIS_NAMES[REDIS_VERSION_MASTER_SLAVE_REDIS],
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					for _, name := range REDIS_NAMES {
						if name == value {
							return
						}
					}
					errors = append(errors, fmt.Errorf("this redis type %s not support now.", value))
					return
				},
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateMysqlPassword,
			},
			"mem_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 100),
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(1, 100),
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return hashcode.String(v.(string))
				},
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  6379,
			},

			// Computed values
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTencentCloudRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("source.tencentcloud_redis_instance.create")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	availabilityZone := d.Get("availability_zone").(string)
	redisName := d.Get("name").(string)
	redisType := d.Get("type").(string)
	password := d.Get("password").(string)
	memSize := d.Get("mem_size").(int)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	securityGroups := d.Get("security_groups").(*schema.Set).List()
	projectId := d.Get("project_id").(int)
	port := d.Get("port").(int)

	requestSecurityGroup := make([]string, 0, len(securityGroups))

	for _, v := range securityGroups {
		requestSecurityGroup = append(requestSecurityGroup, v.(string))
	}

	dealId, err := service.CreateInstances(ctx,
		availabilityZone,
		redisType,
		password,
		vpcId,
		subnetId,
		int64(memSize),
		int64(projectId),
		int64(port),
		requestSecurityGroup)

	if err != nil {
		return err
	}

	if dealId == "" {
		return fmt.Errorf("redis api CreateInstances return  empty redis id")
	}
	var redisId = dealId
	err = resource.Retry(60*time.Minute, func() *resource.RetryError {
		has, online, _, err := service.CheckRedisCreateOk(ctx, dealId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("redis instance not exists."))
		}
		if online {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("create redis task  is processing"))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create redis  task fail, reason:%s\n ", logId, err.Error())
		return err
	}
	if redisName != "" {
		if err := service.ModifyInstanceName(ctx, redisId, redisName); err != nil {
			log.Printf("[CRITAL]%s  redis  set name error, reason:%s\n ", logId, err.Error())
		}
	}
	d.SetId(redisId)
	return nil
}
