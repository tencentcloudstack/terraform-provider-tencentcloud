/*
Provide a resource to create a SCF function.

Example Usage

```hcl
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "scf-code-1234567890"
  cos_object_name   = "code.zip"
  cos_bucket_region = "ap-guangzhou"
}
```

Using Zip file
```
resource "tencentcloud_scf_function" "foo" {
  name              = "ci-test-function"
  handler           = "first.do_it_first"
  runtime           = "Python3.6"
  enable_public_net = true
  dns_cache         = true
  intranet_config {
    ip_fixed = "ENABLE"
  }
  vpc_id    = "vpc-391sv4w3"
  subnet_id = "subnet-ljyn7h30"

  zip_file = "/scf/first.zip"

  tags = {
    "env" = "test"
  }
}

Using CFS config
```
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cfs_config {
	user_id	= "10000"
	user_group_id	= "10000"
	cfs_id	= "cfs-xxxxxxxx"
	mount_ins_id	= "cfs-xxxxxxxx"
	local_mount_dir	= "/mnt"
	remote_mount_dir	= "/"
  }
}
```

Import

SCF function can be imported, e.g.

-> **NOTE:** function id is `<function namespace>+<function name>`

```
$ terraform import tencentcloud_scf_function.test default+test
```
*/
package tencentcloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func scfFunctionValidate(allowDot bool) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (wss []string, errs []error) {
		runes := []rune(v.(string))

		if !unicode.IsLetter(runes[0]) {
			errs = append(errs, errors.Errorf("%s should start with letter", k))
			return
		}

		switch runes[len(runes)-1] {
		case '-', '_':
			errs = append(errs, errors.Errorf(`%s can't end with "-" or "_"`, k))
			return
		}

		for _, r := range runes {
			switch {
			case unicode.IsLetter(r), unicode.IsNumber(r), r == '-', r == '_', r == ':', r == '.' && allowDot:
			default:
				if !allowDot {
					errs = append(errs, errors.Errorf(`invalid %s, %s only can contain a-Z, 0-9, "-" , "_" and ":"`, k, k))
				} else {
					errs = append(errs, errors.Errorf(`invalid %s, %s only can contain a-Z, 0-9, "-", ".", "_" and ":"`, k, k))
				}
				return
			}
		}
		return
	}
}

func resourceTencentCloudScfFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfFunctionCreate,
		Read:   resourceTencentCloudScfFunctionRead,
		Update: resourceTencentCloudScfFunctionUpdate,
		Delete: resourceTencentCloudScfFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: helper.ComposeValidateFunc(
					validateStringLengthInRange(2, 60),
					scfFunctionValidate(false),
				),
				Description: "Name of the SCF function. Name supports 26 English letters, numbers, connectors, and underscores, it should start with a letter. The last character cannot be `-` or `_`. Available length is 2-60.",
			},
			"handler": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: helper.ComposeValidateFunc(
					validateStringLengthInRange(2, 60),
					scfFunctionValidate(true),
				),
				Description: "Handler of the SCF function. The format of name is `<filename>.<method_name>`, and it supports 26 English letters, numbers, connectors, and underscores, it should start with a letter. The last character cannot be `-` or `_`. Available length is 2-60.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validateStringLengthInRange(0, 1000),
				Description:  "Description of the SCF function. Description supports English letters, numbers, spaces, commas, newlines, periods and Chinese, the maximum length is 1000.",
			},
			"mem_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  128,
				ValidateFunc: helper.ComposeValidateFunc(
					//validateIntegerInRange(128, 1536),
					func(v interface{}, k string) (wss []string, errs []error) {
						if v.(int)%128 != 0 {
							errs = append(errs, errors.Errorf("%s should be with 128M as the ladder", k))
						}
						return
					},
				),
				Description: "Memory size of the SCF function, unit is MB. The default is `128`MB. The ladder is 128M.",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validateIntegerInRange(1, 900),
				Description:  "Timeout of the SCF function, unit is second. Default `3`. Available value is 1-900.",
			},
			"environment": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Environment of the SCF function.",
			},
			"runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Runtime of the SCF function, only supports `Python2.7`, `Python3.6`, `Nodejs6.10`, `Nodejs8.9`, `Nodejs10.15`, `PHP5`, `PHP7`, `Golang1`, and `Java8`.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC ID of the SCF function.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet ID of the SCF function.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				ForceNew:    true,
				Description: "Namespace of the SCF function, default is `default`.",
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role of the SCF function.",
			},
			"cls_logset_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "cls logset id of the SCF function.",
			},
			"cls_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "cls topic id of the SCF function.",
			},
			"func_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Event",
				Description: "Function type. The default value is Event. Enter Event if you need to create a trigger function. Enter HTTP if you need to create an HTTP function service.",
			},
			"l5_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable L5 for SCF function, default is `false`.",
			},
			"layers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of association layers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"layer_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of Layer.",
						},
						"layer_version": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The version of layer.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the SCF function.",
			},
			"async_run_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{SCF_FUNCTION_OPEN, SCF_FUNCTION_CLOSE}),
				Description:  "Whether SCF function asynchronous attribute is enabled. `TRUE` is open, `FALSE` is close.",
			},
			"enable_public_net": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether public net config enabled. Default `false`. NOTE: only `vpc_id` specified can disable public net config.",
			},
			"enable_eip_config": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether EIP config set to `ENABLE` when `enable_public_net` was true. Default `false`.",
			},
			// cos code
			"cos_bucket_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"zip_file", "image_config"},
				Description:   "Cos bucket name of the SCF function, such as `cos-1234567890`, conflict with `zip_file`.",
			},
			"cos_object_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"zip_file", "image_config"},
				ValidateFunc:  validateStringSuffix(".zip", ".jar"),
				Description:   "Cos object name of the SCF function, should have suffix `.zip` or `.jar`, conflict with `zip_file`.",
			},
			"cos_bucket_region": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"zip_file", "image_config"},
				Description:   "Cos bucket region of the SCF function, conflict with `zip_file`.",
			},

			// zip upload
			"zip_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"cos_bucket_name", "cos_object_name", "cos_bucket_region", "image_config"},
				Description:   "Zip file of the SCF function, conflict with `cos_bucket_name`, `cos_object_name`, `cos_bucket_region`.",
			},

			// image
			"image_config": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"cos_bucket_name", "cos_object_name", "cos_bucket_region", "zip_file"},
				Description:   "Image of the SCF function, conflict with `cos_bucket_name`, `cos_object_name`, `cos_bucket_region`, `zip_file`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue([]string{"personal", "enterprise"}),
							Description:  "The image type. personal or enterprise.",
						},
						"image_uri": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The uri of image.",
						},
						"registry_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The registry id of TCR. When image type is enterprise, it must be set.",
						},
						"entry_point": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The entrypoint of app.",
						},
						"command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The command of entrypoint.",
						},
						"args": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "the parameters of command.",
						},
						"container_image_accelerate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Image accelerate switch.",
						},
						"image_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      9000,
							ValidateFunc: validateIntegerInRange(-1, 65535),
							Description:  "Image function port setting. Default is `9000`, -1 indicates no port mirroring function. Other value ranges 0 ~ 65535.",
						},
					},
				},
			},

			"triggers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Trigger list of the SCF function, note that if you modify the trigger list, all existing triggers will be deleted, and then create triggers in the new list. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateStringLengthInRange(1, 100),
							Description:  "Name of the SCF function trigger, if `type` is `ckafka`, the format of name must be `<ckafkaInstanceId>-<topicId>`; if `type` is `cos`, the name is cos bucket id, other In any case, it can be combined arbitrarily. It can only contain English letters, numbers, connectors and underscores. The maximum length is 100.",
						},
						"cos_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region of cos bucket. if `type` is `cos`, `cos_region` is required.",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue(SCF_TRIGGER_TYPES),
							Description:  "Type of the SCF function trigger, support `cos`, `cmq`, `timer`, `ckafka`, `apigw`.",
						},
						"trigger_desc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "TriggerDesc of the SCF function trigger, parameter format of `timer` is linux cron expression; parameter of `cos` type is json string `{\"bucketUrl\":\"<name-appid>.cos.<region>.myqcloud.com\",\"event\":\"cos:ObjectCreated:*\",\"filter\":{\"Prefix\":\"\",\"Suffix\":\"\"}}`, where `bucketUrl` is cos bucket (optional), `event` is the cos event trigger, `Prefix` is the corresponding file prefix filter condition, `Suffix` is the suffix filter condition, if not need filter condition can not pass; `cmq` type does not pass this parameter; `ckafka` type parameter format is json string `{\"maxMsgNum\":\"1\",\"offset\":\"latest\"}`; `apigw` type parameter format is json string `{\"api\":{\"authRequired\":\"FALSE\",\"requestConfig\":{\"method\":\"ANY\"},\"isIntegratedResponse\":\"FALSE\"},\"service\":{\"serviceId\":\"service-dqzh68sg\"},\"release\":{\"environmentName\":\"test\"}}`.",
						},
					},
				},
			},

			//cfs config
			"cfs_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of CFS configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of user.",
						},
						"user_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of user group.",
						},
						"cfs_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "File system instance ID.",
						},
						"mount_ins_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "File system mount instance ID.",
						},
						"local_mount_dir": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Local mount directory.",
						},
						"remote_mount_dir": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Remote mount directory.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(Readonly) File system ip address.",
						},
						"mount_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(Readonly) File system virtual private network ID.",
						},
						"mount_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "(Readonly) File system subnet ID.",
						},
					},
				},
			},

			// computed
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function last modified time.",
			},
			"code_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "SCF function code size, unit is M.",
			},
			"code_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function code is correct.",
			},
			"code_error": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function code error message.",
			},
			"err_no": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "SCF function code error code.",
			},
			"install_dependency": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to automatically install dependencies.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function status.",
			},
			"status_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF status description.",
			},
			"eip_fixed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether EIP is a fixed IP.",
			},
			"eips": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "SCF function EIP list.",
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function domain name.",
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF function vip.",
			},
			"trigger_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "SCF trigger details list. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of SCF function trigger.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of SCF function trigger.",
						},
						"trigger_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TriggerDesc of SCF function trigger.",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether SCF function trigger is enable.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of SCF function trigger.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modify time of SCF function trigger.",
						},
						"custom_argument": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User-defined parameters of SCF function trigger.",
						},
					},
				},
			},
			"dns_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable Dns caching capability, only the EVENT function is supported. Default is false.",
			},
			"intranet_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Intranet access configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_fixed": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable fixed intranet IP, ENABLE is enabled, DISABLE is disabled.",
						},
						"ip_address": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "If fixed intranet IP is enabled, this field returns the IP list used.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudScfFunctionCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := m.(*TencentCloudClient).apiV3Conn
	scfService := ScfService{client: client}

	var functionInfo scfFunctionInfo

	functionInfo.name = d.Get("name").(string)
	functionInfo.desc = helper.String(d.Get("description").(string))
	functionInfo.memSize = helper.Int(d.Get("mem_size").(int))
	functionInfo.timeout = helper.Int(d.Get("timeout").(int))
	functionInfo.environment = helper.GetTags(d, "environment")
	functionInfo.namespace = helper.String(d.Get("namespace").(string))

	if raw, ok := d.GetOk("handler"); ok {
		functionInfo.handler = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("runtime"); ok {
		functionInfo.runtime = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("vpc_id"); ok {
		functionInfo.vpcId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("subnet_id"); ok {
		functionInfo.subnetId = helper.String(raw.(string))
	}
	if err := helper.CheckIfSetTogether(d, "vpc_id", "subnet_id"); err != nil {
		return err
	}

	if raw, ok := d.GetOk("role"); ok {
		functionInfo.role = helper.String(raw.(string))
	}

	if raw, ok := d.GetOk("cls_logset_id"); ok {
		functionInfo.clsLogsetId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("cls_topic_id"); ok {
		functionInfo.clsTopicId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("func_type"); ok {
		functionInfo.funcType = helper.String(raw.(string))
	}
	if err := helper.CheckIfSetTogether(d, "cls_logset_id", "cls_topic_id"); err != nil {
		return err
	}

	type scfFunctionCodeType int
	const (
		scfFunctionCosCode scfFunctionCodeType = iota + 1 // start at 1 so we can check if codeType set or not
		scfFunctionZipFileCode
		scfFunctionImageCode
	)

	var codeType scfFunctionCodeType

	if raw, ok := d.GetOk("cos_bucket_name"); ok {
		codeType = scfFunctionCosCode
		functionInfo.cosBucketName = helper.String(raw.(string))
		// to remove string like -1234567890 from bucket id
		split := strings.Split(*functionInfo.cosBucketName, "-")
		if len(split) > 1 {
			functionInfo.cosBucketName = helper.String(strings.Join(split[:len(split)-1], "-"))
		}
	}
	if raw, ok := d.GetOk("cos_object_name"); ok {
		codeType = scfFunctionCosCode
		functionInfo.cosObjectName = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("cos_bucket_region"); ok {
		codeType = scfFunctionCosCode
		functionInfo.cosBucketRegion = helper.String(raw.(string))
	}

	if v, ok := d.GetOk("layers"); ok {
		layers := make([]*scf.LayerVersionSimple, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			layer := scf.LayerVersionSimple{}
			layer.LayerName = helper.String(m["layer_name"].(string))
			layer.LayerVersion = helper.IntInt64(m["layer_version"].(int))
			layers = append(layers, &layer)
		}
		functionInfo.layers = layers
	}

	enablePublicNet, _ := d.GetOk("enable_public_net")
	enableEipConfig, _ := d.GetOk("enable_eip_config")

	if enablePublicNet != nil {
		enable := enablePublicNet.(bool)
		publicNetStatus := helper.String("ENABLE")
		if !enable {
			publicNetStatus = helper.String("DISABLE")
		}
		functionInfo.publicNetConfig = &scf.PublicNetConfigIn{
			PublicNetStatus: publicNetStatus,
			EipConfig: &scf.EipConfigIn{
				EipStatus: helper.String("DISABLE"),
			},
		}
	}

	if enableEipConfig != nil {
		enableEip := enableEipConfig.(bool)
		eipStatus := "DISABLE"
		if enableEip {
			if !enablePublicNet.(bool) {
				return fmt.Errorf("cannot set enable_eip_config to true if enable_public_net was disable")
			}
			eipStatus = "ENABLE"
		}
		functionInfo.publicNetConfig.EipConfig = &scf.EipConfigIn{
			EipStatus: helper.String(eipStatus),
		}
	}

	if raw, ok := d.GetOk("zip_file"); ok {
		path, err := homedir.Expand(raw.(string))
		if err != nil {
			return fmt.Errorf("zip file (%s) homedir expand error: %s", raw.(string), err.Error())
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("zip file (%s) open error: %s", path, err.Error())
		}
		defer file.Close()
		body, err := ioutil.ReadAll(file)
		if err != nil {
			return fmt.Errorf("zip file (%s) read error: %s", path, err.Error())
		}

		codeType = scfFunctionZipFileCode
		content := base64.StdEncoding.EncodeToString(body)
		functionInfo.zipFile = &content
	}

	var imageConfigs = make([]*scf.ImageConfig, 0)

	if raw, ok := d.GetOk("image_config"); ok {
		configs := raw.([]interface{})
		for _, v := range configs {
			value := v.(map[string]interface{})
			imageType := value["image_type"].(string)
			imageUri := value["image_uri"].(string)
			registryId := value["registry_id"].(string)
			entryPoint := value["entry_point"].(string)
			command := value["command"].(string)
			args := value["args"].(string)
			containerImageAccelerate := value["container_image_accelerate"].(bool)
			imagePort := int64(value["image_port"].(int))

			config := &scf.ImageConfig{
				ImageType:                &imageType,
				ImageUri:                 &imageUri,
				RegistryId:               &registryId,
				EntryPoint:               &entryPoint,
				Command:                  &command,
				Args:                     &args,
				ContainerImageAccelerate: &containerImageAccelerate,
				ImagePort:                &imagePort,
			}
			imageConfigs = append(imageConfigs, config)
		}
		codeType = scfFunctionImageCode
		functionInfo.imageConfig = imageConfigs[0]
	}

	switch codeType {
	case scfFunctionCosCode:
		if err := helper.CheckIfSetTogether(d, "cos_bucket_name", "cos_object_name", "cos_bucket_region"); err != nil {
			return err
		}

	case scfFunctionZipFileCode:
	case scfFunctionImageCode:
	default:
		return errors.New("no function code set")
	}

	if v, ok := d.GetOk("cfs_config"); ok {
		configs := v.([]interface{})
		cfsList := make([]*scf.CfsInsInfo, 0, len(configs))

		for i := range configs {
			var (
				item        = configs[i].(map[string]interface{})
				userId      = item["user_id"].(string)
				userGroupId = item["user_group_id"].(string)
				cfsId       = item["cfs_id"].(string)
				mountInsId  = item["mount_ins_id"].(string)
				localMount  = item["local_mount_dir"].(string)
				remoteMount = item["remote_mount_dir"].(string)
			)

			cfsInfo := &scf.CfsInsInfo{
				UserId:         &userId,
				UserGroupId:    &userGroupId,
				CfsId:          &cfsId,
				MountInsId:     &mountInsId,
				LocalMountDir:  &localMount,
				RemoteMountDir: &remoteMount,
			}
			cfsList = append(cfsList, cfsInfo)
		}

		functionInfo.cfsConfig = &scf.CfsConfig{
			CfsInsList: cfsList,
		}
	}

	// Pass tag as creation param instead of modify and time.Sleep
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		functionInfo.tags = tags

		tagService := TagService{client: m.(*TencentCloudClient).apiV3Conn}
		region := m.(*TencentCloudClient).apiV3Conn.Region
		functionId := fmt.Sprintf("%s/function/%s", *functionInfo.namespace, functionInfo.name)
		resourceName := BuildTagResourceName(SCF_SERVICE, SCF_FUNCTION_RESOURCE_PREFIX, region, functionId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			log.Printf("[CRITAL]%s create function tags failed: %+v", logId, err)
			return err
		}
		// wait for tags add successfully
		time.Sleep(time.Second)
	}

	if v, ok := d.GetOk("async_run_enable"); ok && v != nil {
		enableStr := v.(string)
		functionInfo.asyncRunEnable = helper.String(enableStr)
	}

	funcType := *functionInfo.funcType
	if funcType == SCF_FUNCTION_TYPE_EVENT {
		if v, ok := d.GetOk("dns_cache"); ok {
			dnsCache := v.(bool)
			if dnsCache {
				functionInfo.dnsCache = helper.String("TRUE")
			} else {
				functionInfo.dnsCache = helper.String("FALSE")
			}
		}
	}
	if raw, ok := d.GetOk("intranet_config"); ok {
		configs := raw.([]interface{})
		var intranetConfigs = make([]*scf.IntranetConfigIn, 0)

		for _, v := range configs {
			value := v.(map[string]interface{})
			ipFixed := value["ip_fixed"].(string)

			config := &scf.IntranetConfigIn{
				IpFixed: &ipFixed,
			}
			intranetConfigs = append(intranetConfigs, config)
		}
		functionInfo.intranetConfig = intranetConfigs[0]
	}

	if err := scfService.CreateFunction(ctx, functionInfo); err != nil {
		log.Printf("[CRITAL]%s create function failed: %+v", logId, err)
		return err
	}

	// id format is [namespace]+[function name], so that we can support import with enough info
	d.SetId(fmt.Sprintf("%s+%s", *functionInfo.namespace, functionInfo.name))

	err := waitScfFunctionReady(ctx, functionInfo.name, *functionInfo.namespace, client.UseScfClient())
	if err != nil {
		return err
	}

	if d.Get("l5_enable").(bool) {
		if err := scfService.ModifyFunctionConfig(ctx, scfFunctionInfo{
			name:      functionInfo.name,
			namespace: functionInfo.namespace,
			l5Enable:  helper.Bool(true),
		}); err != nil {
			log.Printf("[CRITAL]%s enable function L5 failed: %+v", logId, err)
			return err
		}
	}

	if raw, ok := d.GetOk("triggers"); ok {
		set := raw.(*schema.Set)
		triggers := make([]scfTrigger, 0, set.Len())
		for _, rawTrigger := range set.List() {
			tg := rawTrigger.(map[string]interface{})

			switch tg["type"].(string) {
			case SCF_TRIGGER_TYPE_COS:
				if tg["cos_region"].(string) == "" {
					return fmt.Errorf("type if cos, cos_region is required")
				}
				// scf cos trigger name format is xxx-1234567890.cos.ap-guangzhou.myqcloud.com
				tg["name"] = fmt.Sprintf("%s.cos.%s.myqcloud.com", tg["name"].(string), tg["cos_region"].(string))
			}

			triggers = append(triggers, scfTrigger{
				name:        tg["name"].(string),
				triggerType: tg["type"].(string),
				triggerDesc: tg["trigger_desc"].(string),
			})
		}

		if err := scfService.CreateTriggers(ctx, functionInfo.name, *functionInfo.namespace, triggers); err != nil {
			log.Printf("[CRITAL]%s create triggers failed: %+v", logId, err)
			return err
		}
	}

	_, err = scfService.DescribeFunction(ctx, functionInfo.name, *functionInfo.namespace)
	if err != nil {
		log.Printf("[CRITAL]%s get function id failed: %+v", logId, err)
		return err
	}

	return resourceTencentCloudScfFunctionRead(d, m)
}

func resourceTencentCloudScfFunctionRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

	id := d.Id()
	split := strings.Split(id, "+")
	if len(split) != 2 {
		log.Printf("[CRITAL]%s function id is invalid", logId)
		d.SetId("")
		return nil
	}
	namespace, name := split[0], split[1]

	response, err := service.DescribeFunction(ctx, name, namespace)
	if err != nil {
		log.Printf("[CRITAL]%s read function failed: %+v", logId, err)
	}

	if response == nil {
		d.SetId("")
		return nil
	}

	resp := response.Response

	_ = d.Set("name", resp.FunctionName)
	_ = d.Set("handler", resp.Handler)
	_ = d.Set("description", resp.Description)
	_ = d.Set("mem_size", resp.MemorySize)
	_ = d.Set("timeout", resp.Timeout)

	environment := make(map[string]string, len(resp.Environment.Variables))
	for _, v := range resp.Environment.Variables {
		environment[*v.Key] = *v.Value
	}
	_ = d.Set("environment", environment)

	_ = d.Set("runtime", resp.Runtime)
	_ = d.Set("vpc_id", resp.VpcConfig.VpcId)
	_ = d.Set("subnet_id", resp.VpcConfig.SubnetId)
	_ = d.Set("namespace", resp.Namespace)
	_ = d.Set("role", resp.Role)
	_ = d.Set("cls_logset_id", resp.ClsLogsetId)
	_ = d.Set("cls_topic_id", resp.ClsTopicId)
	_ = d.Set("func_type", resp.Type)
	_ = d.Set("l5_enable", *resp.L5Enable == "TRUE")

	tags := make(map[string]string, len(resp.Tags))
	for _, tag := range resp.Tags {
		tags[*tag.Key] = *tag.Value
	}
	_ = d.Set("tags", tags)
	_ = d.Set("async_run_enable", resp.AsyncRunEnable)

	_ = d.Set("modify_time", resp.ModTime)
	_ = d.Set("code_size", resp.CodeSize)
	_ = d.Set("code_result", resp.CodeResult)
	_ = d.Set("code_error", resp.CodeError)
	_ = d.Set("err_no", resp.ErrNo)
	_ = d.Set("install_dependency", *resp.InstallDependency == "TRUE")
	_ = d.Set("status", resp.Status)
	_ = d.Set("status_desc", resp.StatusDesc)
	_ = d.Set("eip_fixed", *resp.EipConfig.EipFixed == "TRUE")
	_ = d.Set("eips", resp.EipConfig.Eips)
	_ = d.Set("host", resp.AccessInfo.Host)
	_ = d.Set("vip", resp.AccessInfo.Vip)
	if resp.PublicNetConfig != nil {
		_ = d.Set("enable_public_net", *resp.PublicNetConfig.PublicNetStatus == "ENABLE")
		_ = d.Set("enable_eip_config", *resp.PublicNetConfig.EipConfig.EipStatus == "ENABLE")
	}

	if resp.CfsConfig != nil {
		cfsList := make([]interface{}, 0, len(resp.CfsConfig.CfsInsList))
		for i := range resp.CfsConfig.CfsInsList {
			item := resp.CfsConfig.CfsInsList[i]
			cfs := map[string]interface{}{
				"user_id":          item.UserId,
				"user_group_id":    item.UserGroupId,
				"cfs_id":           item.CfsId,
				"mount_ins_id":     item.MountInsId,
				"local_mount_dir":  item.LocalMountDir,
				"remote_mount_dir": item.RemoteMountDir,
			}

			if item.IpAddress != nil {
				cfs["ip_address"] = item.IpAddress
			}
			if item.MountVpcId != nil {
				cfs["mount_vpc_id"] = item.MountVpcId
			}
			if item.MountSubnetId != nil {
				cfs["mount_subnet_id"] = item.MountSubnetId
			}
			cfsList = append(cfsList, cfs)
		}
		if err := d.Set("cfs_config", cfsList); err != nil {
			return err
		}
	}

	triggers := make([]map[string]interface{}, 0, len(resp.Triggers))
	for _, trigger := range resp.Triggers {
		switch *trigger.Type {
		case SCF_TRIGGER_TYPE_TIMER:
			data := struct {
				Cron string `json:"cron"`
			}{}
			if err := json.Unmarshal([]byte(*trigger.TriggerDesc), &data); err != nil {
				log.Printf("[WARN]%s unmarshal timer trigger trigger_desc failed: %+v", logId, errors.WithStack(err))
				continue
			}
			*trigger.TriggerDesc = data.Cron
		}

		triggers = append(triggers, map[string]interface{}{
			"name":            trigger.TriggerName,
			"type":            trigger.Type,
			"trigger_desc":    trigger.TriggerDesc,
			"enable":          *trigger.Enable == 1,
			"create_time":     trigger.AddTime,
			"modify_time":     trigger.ModTime,
			"custom_argument": trigger.CustomArgument,
		})
	}
	_ = d.Set("trigger_info", triggers)

	if resp.ImageConfig != nil {
		imageConfigs := make([]map[string]interface{}, 0, 1)
		imageConfigResp := resp.ImageConfig

		imageConfig := map[string]interface{}{
			"image_type": imageConfigResp.ImageType,
			"image_uri":  imageConfigResp.ImageUri,
		}
		if imageConfigResp.RegistryId != nil {
			imageConfig["registry_id"] = imageConfigResp.RegistryId
		}
		if imageConfigResp.EntryPoint != nil {
			imageConfig["entry_point"] = imageConfigResp.EntryPoint
		}
		if imageConfigResp.Command != nil {
			imageConfig["command"] = imageConfigResp.Command
		}
		if imageConfigResp.Args != nil {
			imageConfig["args"] = imageConfigResp.Args
		}
		if imageConfigResp.ContainerImageAccelerate != nil {
			imageConfig["container_image_accelerate"] = imageConfigResp.ContainerImageAccelerate
		}
		if imageConfigResp.ImagePort != nil {
			imageConfig["image_port"] = imageConfigResp.ImagePort
		}
		imageConfigs = append(imageConfigs, imageConfig)

		if err = d.Set("image_config", imageConfigs); err != nil {
			return err
		}
	}

	if resp.DnsCache != nil {
		_ = d.Set("dns_cache", *resp.DnsCache == "TRUE")
	}
	if resp.IntranetConfig != nil {
		intranetConfigs := make([]map[string]interface{}, 0, 1)
		intranetConfigResp := resp.IntranetConfig

		intranetConfig := map[string]interface{}{
			"ip_fixed": intranetConfigResp.IpFixed,
		}
		if intranetConfigResp.IpAddress != nil {
			intranetConfig["ip_address"] = intranetConfigResp.IpAddress
		}
		intranetConfigs = append(intranetConfigs, intranetConfig)

		if err = d.Set("intranet_config", intranetConfigs); err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudScfFunctionUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := m.(*TencentCloudClient).apiV3Conn
	scfService := ScfService{client: client}
	tagService := TagService{client: client}
	region := client.Region

	id := d.Id()
	split := strings.Split(id, "+")
	if len(split) != 2 {
		log.Printf("[CRITAL]%s function id is invalid", logId)
		d.SetId("")
		return nil
	}
	namespace, name := split[0], split[1]

	if d.HasChange("func_type") {
		return fmt.Errorf("`func_type` do not support change now.")
	}

	d.Partial(true)

	functionInfo := scfFunctionInfo{
		name:      name,
		namespace: helper.String(namespace),
	}

	var updateAttrs []string

	if d.HasChange("handler") {
		updateAttrs = append(updateAttrs, "handler")
	}

	if d.HasChange("cos_bucket_name") {
		updateAttrs = append(updateAttrs, "cos_bucket_name")
	}
	if d.HasChange("cos_object_name") {
		updateAttrs = append(updateAttrs, "cos_object_name")
	}
	if d.HasChange("cos_bucket_region") {
		updateAttrs = append(updateAttrs, "cos_bucket_region")
	}
	if raw, ok := d.GetOk("cos_bucket_name"); ok {
		functionInfo.cosBucketName = helper.String(raw.(string))
		// to remove string like -1234567890 from bucket id
		split := strings.Split(*functionInfo.cosBucketName, "-")
		if len(split) > 1 {
			functionInfo.cosBucketName = helper.String(strings.Join(split[:len(split)-1], "-"))
		}
	}
	if raw, ok := d.GetOk("cos_object_name"); ok {
		functionInfo.cosObjectName = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("cos_bucket_region"); ok {
		functionInfo.cosBucketRegion = helper.String(raw.(string))
	}
	if err := helper.CheckIfSetTogether(d, "cos_bucket_name", "cos_object_name", "cos_bucket_region"); err != nil {
		return err
	}

	if d.HasChange("zip_file") {
		updateAttrs = append(updateAttrs, "zip_file")
	}
	if raw, ok := d.GetOk("zip_file"); ok {
		path, err := homedir.Expand(raw.(string))
		if err != nil {
			return fmt.Errorf("zip file (%s) homedir expand error: %s", raw.(string), err.Error())
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("zip file (%s) open error: %s", path, err.Error())
		}
		defer file.Close()
		body, err := ioutil.ReadAll(file)
		if err != nil {
			return fmt.Errorf("zip file (%s) read error: %s", path, err.Error())
		}

		content := base64.StdEncoding.EncodeToString(body)
		functionInfo.zipFile = &content
	}

	if d.HasChange("image_config") {
		updateAttrs = append(updateAttrs, "image_config")
		if raw, ok := d.GetOk("image_config"); ok {
			var imageConfigs = make([]*scf.ImageConfig, 0)
			configs := raw.([]interface{})
			for _, v := range configs {
				value := v.(map[string]interface{})
				imageType := value["image_type"].(string)
				imageUri := value["image_uri"].(string)
				registryId := value["registry_id"].(string)
				entryPoint := value["entry_point"].(string)
				command := value["command"].(string)
				args := value["args"].(string)
				containerImageAccelerate := value["container_image_accelerate"].(bool)
				imagePort := int64(value["image_port"].(int))

				config := &scf.ImageConfig{
					ImageType:                &imageType,
					ImageUri:                 &imageUri,
					RegistryId:               &registryId,
					EntryPoint:               &entryPoint,
					Command:                  &command,
					Args:                     &args,
					ContainerImageAccelerate: &containerImageAccelerate,
					ImagePort:                &imagePort,
				}
				imageConfigs = append(imageConfigs, config)
			}
			functionInfo.imageConfig = imageConfigs[0]
		}
	}

	if d.HasChange("cfs_config") {
		updateAttrs = append(updateAttrs, "cfs_config")
		if v, ok := d.GetOk("cfs_config"); ok {
			configs := v.([]interface{})
			cfsList := make([]*scf.CfsInsInfo, 0, len(configs))

			for i := range configs {
				var (
					item        = configs[i].(map[string]interface{})
					userId      = item["user_id"].(string)
					userGroupId = item["user_group_id"].(string)
					cfsId       = item["cfs_id"].(string)
					mountInsId  = item["mount_ins_id"].(string)
					localMount  = item["local_mount_dir"].(string)
					remoteMount = item["remote_mount_dir"].(string)
				)

				cfsInfo := &scf.CfsInsInfo{
					UserId:         &userId,
					UserGroupId:    &userGroupId,
					CfsId:          &cfsId,
					MountInsId:     &mountInsId,
					LocalMountDir:  &localMount,
					RemoteMountDir: &remoteMount,
				}
				cfsList = append(cfsList, cfsInfo)
			}

			functionInfo.cfsConfig = &scf.CfsConfig{
				CfsInsList: cfsList,
			}
		} else {
			functionInfo.cfsConfig = &scf.CfsConfig{
				CfsInsList: []*scf.CfsInsInfo{},
			}
		}
	}

	// update function code
	if len(updateAttrs) > 0 {
		if len(updateAttrs) == 0 && updateAttrs[0] == "handler" {
			return errors.New("can't only change handler")
		}

		functionInfo.handler = helper.String(d.Get("handler").(string))

		if err := scfService.ModifyFunctionCode(ctx, functionInfo); err != nil {
			log.Printf("[CRITAL]%s update function code failed: %+v", logId, err)
			return err
		}

	}

	updateAttrs = updateAttrs[:0]
	functionInfo = scfFunctionInfo{
		name:      name,
		namespace: helper.String(namespace),
	}

	if d.HasChange("description") {
		updateAttrs = append(updateAttrs, "description")
		functionInfo.desc = helper.String(d.Get("description").(string))
	}
	if d.HasChange("mem_size") {
		updateAttrs = append(updateAttrs, "mem_size")
		functionInfo.memSize = helper.Int(d.Get("mem_size").(int))
	}
	if d.HasChange("timeout") {
		updateAttrs = append(updateAttrs, "timeout")
		functionInfo.timeout = helper.Int(d.Get("timeout").(int))
	}

	if d.HasChange("environment") {
		updateAttrs = append(updateAttrs, "environment")
	}
	functionInfo.environment = helper.GetTags(d, "environment")

	if d.HasChange("runtime") {
		updateAttrs = append(updateAttrs, "runtime")
		functionInfo.runtime = helper.String(d.Get("runtime").(string))
	}

	if d.HasChange("vpc_id") {
		updateAttrs = append(updateAttrs, "vpc_id")
	}
	if d.HasChange("subnet_id") {
		updateAttrs = append(updateAttrs, "subnet_id")
	}
	if raw, ok := d.GetOk("vpc_id"); ok {
		functionInfo.vpcId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("subnet_id"); ok {
		functionInfo.subnetId = helper.String(raw.(string))
	}
	if err := helper.CheckIfSetTogether(d, "vpc_id", "subnet_id"); err != nil {
		return err
	}

	if d.HasChange("role") {
		updateAttrs = append(updateAttrs, "role")
		functionInfo.role = helper.String(d.Get("role").(string))
	}

	if d.HasChange("cls_logset_id") {
		updateAttrs = append(updateAttrs, "cls_logset_id")
	}
	if d.HasChange("cls_topic_id") {
		updateAttrs = append(updateAttrs, "cls_topic_id")
	}
	if raw, ok := d.GetOk("cls_logset_id"); ok {
		functionInfo.clsLogsetId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("cls_topic_id"); ok {
		functionInfo.clsTopicId = helper.String(raw.(string))
	}
	if err := helper.CheckIfSetTogether(d, "cls_logset_id", "cls_topic_id"); err != nil {
		return err
	}

	if d.HasChange("l5_enable") {
		updateAttrs = append(updateAttrs, "l5_enable")
		functionInfo.l5Enable = helper.Bool(d.Get("l5_enable").(bool))
	}

	if d.HasChange("enable_public_net") {
		updateAttrs = append(updateAttrs, "enable_public_net")
	}

	if d.HasChange("enable_eip_config") {
		updateAttrs = append(updateAttrs, "enable_eip_config")
	}

	if raw, _ := d.GetOk("enable_public_net"); raw != nil {
		enablePublicNet := raw.(bool)
		publicNetStatus := helper.String("ENABLE")
		if !enablePublicNet {
			publicNetStatus = helper.String("DISABLE")
		}
		functionInfo.publicNetConfig = &scf.PublicNetConfigIn{
			PublicNetStatus: publicNetStatus,
			EipConfig: &scf.EipConfigIn{
				EipStatus: helper.String("DISABLE"),
			},
		}
	}

	if raw, _ := d.GetOk("enable_eip_config"); raw != nil {
		status := "DISABLE"
		enablePublicNet := d.Get("enable_public_net").(bool)
		if raw.(bool) {
			if !enablePublicNet {
				return fmt.Errorf("cannot set enable_eip_config to true if enable_public_net was disable")
			}
			status = "ENABLE"
		}
		functionInfo.publicNetConfig.EipConfig = &scf.EipConfigIn{
			EipStatus: helper.String(status),
		}
	}

	var funcType string
	if raw, ok := d.GetOk("func_type"); ok {
		funcType = raw.(string)
	}
	if funcType == SCF_FUNCTION_TYPE_EVENT && d.HasChange("dns_cache") {
		updateAttrs = append(updateAttrs, "dns_cache")
		if v, ok := d.GetOkExists("dns_cache"); ok {
			dnsCache := v.(bool)
			if dnsCache {
				functionInfo.dnsCache = helper.String("TRUE")
			} else {
				functionInfo.dnsCache = helper.String("FALSE")
			}
		}
	}
	if d.HasChange("intranet_config") {
		updateAttrs = append(updateAttrs, "intranet_config")
		if raw, ok := d.GetOk("intranet_config"); ok {
			configs := raw.([]interface{})
			var intranetConfigs = make([]*scf.IntranetConfigIn, 0)

			for _, v := range configs {
				value := v.(map[string]interface{})
				ipFixed := value["ip_fixed"].(string)

				config := &scf.IntranetConfigIn{
					IpFixed: &ipFixed,
				}
				intranetConfigs = append(intranetConfigs, config)
			}
			functionInfo.intranetConfig = intranetConfigs[0]
		}
	}

	// update function configuration
	if len(updateAttrs) > 0 {
		if err := scfService.ModifyFunctionConfig(ctx, functionInfo); err != nil {
			log.Printf("[CRITAL]%s update function configuration failed: %+v", logId, err)
			return err
		}
	}

	if d.HasChange("triggers") {
		oldRaw, newRaw := d.GetChange("triggers")
		oldSet := oldRaw.(*schema.Set)
		newSet := newRaw.(*schema.Set)

		oldTriggers := make([]scfTrigger, 0, oldSet.Len())
		for _, trigger := range oldSet.List() {
			tg := trigger.(map[string]interface{})

			switch tg["type"].(string) {
			case SCF_TRIGGER_TYPE_COS:
				tg["name"] = fmt.Sprintf("%s.cos.%s.myqcloud.com", tg["name"].(string), tg["cos_region"].(string))
			}

			oldTriggers = append(oldTriggers, scfTrigger{
				name:        tg["name"].(string),
				triggerType: tg["type"].(string),
				triggerDesc: tg["trigger_desc"].(string),
			})
		}
		if err := scfService.DeleteTriggers(ctx, name, namespace, oldTriggers); err != nil {
			log.Printf("[CRITAL]%s delete old triggers failed: %+v", logId, err)
			return err
		}

		newTriggers := make([]scfTrigger, 0, newSet.Len())
		for _, trigger := range newSet.List() {
			tg := trigger.(map[string]interface{})

			switch tg["type"].(string) {
			case SCF_TRIGGER_TYPE_COS:
				if tg["cos_region"].(string) == "" {
					return fmt.Errorf("type if cos, cos_region is required")
				}
				tg["name"] = fmt.Sprintf("%s.cos.%s.myqcloud.com", tg["name"].(string), tg["cos_region"].(string))
			}

			newTriggers = append(newTriggers, scfTrigger{
				name:        tg["name"].(string),
				triggerType: tg["type"].(string),
				triggerDesc: tg["trigger_desc"].(string),
			})
		}
		if err := scfService.CreateTriggers(ctx, name, namespace, newTriggers); err != nil {
			log.Printf("[CRITAL]%s create new triggers failed: %+v", logId, err)
			return err
		}

	}

	if d.HasChange("tags") {
		resp, err := scfService.DescribeFunction(ctx, functionInfo.name, *functionInfo.namespace)
		if err != nil {
			log.Printf("[CRITAL]%s get function id failed: %+v", logId, err)
			return err
		}
		fnName := *resp.Response.FunctionName
		fnNamespace := *resp.Response.Namespace
		functionId := fmt.Sprintf("%s/function/%s", fnNamespace, fnName)

		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName(SCF_SERVICE, SCF_FUNCTION_RESOURCE_PREFIX, region, functionId)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			log.Printf("[CRITAL]%s update function tags failed: %+v", logId, err)
			return err
		}

		// wait for tags add successfully
		time.Sleep(time.Second)
	}

	d.Partial(false)

	return resourceTencentCloudScfFunctionRead(d, m)
}

func resourceTencentCloudScfFunctionDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_function.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

	id := d.Id()
	split := strings.Split(id, "+")
	if len(split) != 2 {
		log.Printf("[CRITAL]%s function id is invalid", logId)
		return nil
	}
	namespace, name := split[0], split[1]

	return service.DeleteFunction(ctx, name, namespace)
}
