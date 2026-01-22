package cvm_test

import (
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_instance", &resource.Sweeper{
		Name: "tencentcloud_instance",
		F:    testSweepCvmInstance,
	})
}

func testSweepCvmInstance(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	cvmService := svccvm.NewCvmService(client)

	instances, err := cvmService.DescribeInstanceByFilter(ctx, nil, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range instances {
		if !tccommon.CheckResourcePersist(*v.InstanceName, *v.CreatedTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.InstanceId,
				Name: *v.InstanceName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:         *v.InstanceId,
			Name:       *v.InstanceName,
			CreateTime: *v.CreatedTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "RunInstances")

	for _, v := range instances {
		instanceId := *v.InstanceId
		//instanceName := *v.InstanceName
		now := time.Now()
		createTime := tccommon.StringToTime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()

		//if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
		//	continue
		//}

		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = cvmService.DeleteInstance(ctx, instanceId, false); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudInstanceResourceBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_basic"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_basic", "tags.hostname", "tci"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_basic", "instance_status", "RUNNING"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_basic", "private_ip"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_basic", "vpc_id"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_basic", "subnet_id"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_basic", "project_id")),
			},
			{
				Config: testAccCvmInstanceResource_BasicChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_basic"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_basic", "instance_type"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_basic", "instance_status", "RUNNING")),
			},
			{
				Config: testAccCvmInstanceResource_BasicChangeCamRoleName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_basic"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_basic", "cam_role_name", "CVM_QcsRole"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_basic", "instance_status", "RUNNING"),
				),
			},
			{
				ResourceName:            "tencentcloud_instance.cvm_basic",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"disable_monitor_service", "disable_security_service", "disable_automation_service", "hostname", "password", "force_delete"},
			},
		},
	})
}
func TestAccTencentCloudInstanceResource_UserDataRaw(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_UserDataRaw,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_user_data_raw"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_user_data_raw", "user_data_raw", "test"),
				),
			},
			{
				Config: testAccCvmInstanceResource_UserDataRawUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_user_data_raw"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_user_data_raw", "user_data_raw", "test-update"),
				),
			},
			{
				ResourceName:            "tencentcloud_instance.cvm_user_data_raw",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"disable_monitor_service", "disable_security_service", "disable_automation_service", "hostname", "password", "force_delete"},
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_UserData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_UserData,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_user_data"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_user_data", "user_data", "dGVzdA=="),
				),
			},
			{
				Config: testAccCvmInstanceResource_UserDataUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_user_data"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_user_data", "user_data", "dGVzdC11cGRhdGU="),
				),
			},
			{
				ResourceName:            "tencentcloud_instance.cvm_user_data",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"disable_monitor_service", "disable_security_service", "disable_automation_service", "hostname", "password", "force_delete"},
			},
		},
	})
}

func testAccCheckCvmInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := common.GetLogId(common.ContextNil)
		ctx := context.WithValue(context.TODO(), common.LogIdKey, logId)
		service := cvm.NewCvmService(acctest.AccProvider.Meta().(common.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource `%s` is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource `%s` is not found", n)
		}
		id := rs.Primary.ID

		result, err := service.DescribeInstanceById(ctx, id)
		if err != nil {
			return err
		}
		if result == nil {
			return fmt.Errorf("resource `%s` create failed", id)
		}
		return nil
	}
}
func testAccCheckCvmInstanceDestroy(s *terraform.State) error {
	logId := common.GetLogId(common.ContextNil)
	ctx := context.WithValue(context.TODO(), common.LogIdKey, logId)
	service := cvm.NewCvmService(acctest.AccProvider.Meta().(common.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		id := rs.Primary.ID
		if rs.Type != "tencentcloud_cvm_instance" {
			continue
		}
		result, err := service.DescribeInstanceById(ctx, id)
		if err != nil {
			return err
		}
		if result != nil {
			return fmt.Errorf("resource `%s` still exist", id)
		}
	}
	return nil
}

const testAccCvmInstanceResource_UserData = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-basic-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_user_data" {
    instance_name = "tf-test-user-data"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    tags = {
        hostname = "tci"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
    user_data = "dGVzdA=="
}

`
const testAccCvmInstanceResource_UserDataUpdate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-basic-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_user_data" {
    instance_name = "tf-test-user-data"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    tags = {
        hostname = "tci"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
    user_data = "dGVzdC11cGRhdGU="
}

`

const testAccCvmInstanceResource_UserDataRaw = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-basic-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_user_data_raw" {
    instance_name = "tf-test-user-data-raw"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    tags = {
        hostname = "tci"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
    user_data_raw = "test"
}

`
const testAccCvmInstanceResource_UserDataRawUpdate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-basic-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_user_data_raw" {
    instance_name = "tf-test-user-data-raw"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    tags = {
        hostname = "tci"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
    user_data_raw = "test-update"
}

`

const testAccCvmInstanceResource_BasicCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-basic-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_basic" {
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    tags = {
        hostname = "tci"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
}

`
const testAccCvmInstanceResource_BasicChange1 = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        values = ["ap-guangzhou-7"]
        name = "zone"
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-basic-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_basic" {
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    tags = {
        hostname = "tci"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
}

`

const testAccCvmInstanceResource_BasicChangeCamRoleName = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        values = ["ap-guangzhou-7"]
        name = "zone"
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-basic-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_basic" {
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    tags = {
        hostname = "tci"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
    cam_role_name = "CVM_QcsRole"
}

`

func TestAccTencentCloudInstance_SystemDiskResizeOnline(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_SystemDiskResizeOnline,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_system_disk_resize_online"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_system_disk_resize_online", "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_system_disk_resize_online", "system_disk_size", "100"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_system_disk_resize_online", "instance_status", "RUNNING"),
				),
			},
			{
				Config: testAccCvmInstanceResource_SystemDiskResizeOnlineUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_system_disk_resize_online"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_system_disk_resize_online", "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_system_disk_resize_online", "system_disk_size", "200"),
					resource.TestCheckResourceAttr("tencentcloud_instance.cvm_system_disk_resize_online", "instance_status", "RUNNING"),
				),
			},
		},
	})
}

const testAccCvmInstanceResource_SystemDiskResizeOnline = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-resize-online-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-resize-online-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_system_disk_resize_online" {
    instance_name = "tf-system-disk-resize-online"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    system_disk_size = 100
    project_id = 0
}
`

const testAccCvmInstanceResource_SystemDiskResizeOnlineUpdate = `
data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-resize-online-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-resize-online-subnet"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_instance" "cvm_system_disk_resize_online" {
    instance_name = "tf-system-disk-resize-online"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    system_disk_size = 200
    system_disk_resize_online = true
    project_id = 0
}
`

func TestAccTencentCloudInstanceResourcePrepaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_PrepaidCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_prepaid_basic"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_prepaid_basic", "vpc_id"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_prepaid_basic", "subnet_id"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_prepaid_basic", "project_id"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_prepaid_basic", "tags.hostname", "tci"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_prepaid_basic", "instance_status", "RUNNING"), resource.TestCheckResourceAttrSet("tencentcloud_instance.cvm_prepaid_basic", "private_ip")),
			},
		},
	})
}

const testAccCvmInstanceResource_PrepaidCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-prepaid-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-prepaid-basic-subnet"
    cidr_block = "10.0.0.0/16"
    availability_zone = "ap-guangzhou-7"
}
resource "tencentcloud_instance" "cvm_prepaid_basic" {
    force_delete = true
    availability_zone = "ap-guangzhou-7"
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
    instance_charge_type = "PREPAID"
    instance_charge_type_prepaid_period = 1
    instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
    
    tags = {
        hostname = "tci"
    }
    instance_name = "tf-ci-test"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    vpc_id = tencentcloud_vpc.vpc.id
}

`

func TestAccTencentCloudInstanceResourceWithDataDisk(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithDataDiskCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.data_disk_size", "100"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.data_disk_snapshot_id", ""), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.1.data_disk_type", "CLOUD_PREMIUM"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.1.data_disk_size", "100"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "system_disk_size", "100"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "system_disk_type", "CLOUD_PREMIUM"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.data_disk_type", "CLOUD_PREMIUM")),
			},
			{
				Config: testAccCvmInstanceResource_WithDataDiskChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.data_disk_type", "CLOUD_PREMIUM"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.data_disk_size", "150"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.delete_with_instance", "true"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.1.data_disk_type", "CLOUD_PREMIUM"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.1.data_disk_size", "150"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.1.delete_with_instance", "true")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithDataDiskCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_instance" "foo" {
    system_disk_type = "CLOUD_PREMIUM"
    disable_monitor_service = true
    instance_name = "tf-ci-test"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_size = 100
    system_disk_name = "sys_cbs_test1"
    
    data_disks {
        delete_with_instance = true
        data_disk_type = "CLOUD_PREMIUM"
        data_disk_size = 100
        data_disk_name = "data_cbs_test1"
    }
    data_disks {
        data_disk_type = "CLOUD_PREMIUM"
        data_disk_size = 100
        delete_with_instance = true
        data_disk_name = "data_cbs_test2"
    }
    disable_security_service = true
    disable_automation_service = true
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    availability_zone = "ap-guangzhou-7"
}

`
const testAccCvmInstanceResource_WithDataDiskChange1 = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_instance" "foo" {
    system_disk_type = "CLOUD_PREMIUM"
    disable_monitor_service = true
    instance_name = "tf-ci-test"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_size = 100
    system_disk_name = "sys_cbs_test1"
    
    data_disks {
        data_disk_size = 150
        delete_with_instance = true
        data_disk_type = "CLOUD_PREMIUM"
        data_disk_name = "data_cbs_test1_update"
    }
    data_disks {
        data_disk_size = 150
        delete_with_instance = true
        data_disk_type = "CLOUD_PREMIUM"
        data_disk_name = "data_cbs_test2_update"
    }
    disable_security_service = true
    disable_automation_service = true
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    availability_zone = "ap-guangzhou-7"
}

`

func TestAccTencentCloudInstanceResourceWithNetwork(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithNetworkCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttrSet("tencentcloud_instance.foo", "public_ip")),
			},
			{
				Config: testAccCvmInstanceResource_WithNetworkChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "internet_max_bandwidth_out", "5"), resource.TestCheckResourceAttrSet("tencentcloud_instance.foo", "public_ip"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithNetworkCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_instance" "foo" {
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    internet_max_bandwidth_out = 5
    allocate_public_ip = true
    system_disk_type = "CLOUD_PREMIUM"
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
}

`
const testAccCvmInstanceResource_WithNetworkChange1 = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_instance" "foo" {
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    internet_max_bandwidth_out = "5"
    allocate_public_ip = true
    system_disk_type = "CLOUD_PREMIUM"
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
}

`

func TestAccTencentCloudInstanceResourceWithPrivateIp(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithPrivateIpCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithPrivateIpCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-with-privateip-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-with-privateip-subnet"
    cidr_block = "10.0.0.0/16"
    availability_zone = "ap-guangzhou-7"
}
resource "tencentcloud_instance" "foo" {
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    vpc_id = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
    private_ip = "10.0.0.123"
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
}

`

// TODO generate
func TestAccTencentCloudInstanceResource_WithKeyPairs(t *testing.T) {
	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithKeyPair_withoutKeyPair,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithKeyPair(
					"[tencentcloud_key_pair.key_pair_0.id, tencentcloud_key_pair.key_pair_1.id]",
				),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "key_ids.#", "2"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 5)
				},
				Config: testAccTencentCloudInstanceWithKeyPair("[tencentcloud_key_pair.key_pair_2.id]"),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "key_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResourceWithPassword(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithPasswordCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttrSet("tencentcloud_instance.foo", "instance_status"), resource.TestCheckResourceAttrSet("tencentcloud_instance.foo", "password")),
			},
			{
				Config: testAccCvmInstanceResource_WithPasswordChange1,
				PreConfig: func() {
					time.Sleep(time.Second * 5)
				},
				Check: resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttrSet("tencentcloud_instance.foo", "instance_status"), resource.TestCheckResourceAttrSet("tencentcloud_instance.foo", "password")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithPasswordCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_instance" "foo" {
    password = "TF_test_123"
    system_disk_type = "CLOUD_PREMIUM"
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
}

`
const testAccCvmInstanceResource_WithPasswordChange1 = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        values = ["S1","S2","S3","S4","S5"]
        name = "instance-family"
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_instance" "foo" {
    password = "TF_test_123"
    system_disk_type = "CLOUD_PREMIUM"
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
}

`

func TestAccTencentCloudInstanceResourceWithImageLogin(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithImageLoginCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "keep_image_login", "true"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "disable_api_termination", "false")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithImageLoginCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_images" "zoo" {
    image_type = ["PRIVATE_IMAGE"]
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        values = ["S1","S2","S3","S4","S5"]
        name = "instance-family"
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_instance" "foo" {
    system_disk_type = "CLOUD_PREMIUM"
    disable_api_termination = false
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.zoo.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    keep_image_login = true
}

`

func TestAccTencentCloudInstanceResourceWithName(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithNameCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_name", "tf-ci-test")),
			},
			{
				Config: testAccCvmInstanceResource_WithNameChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_name", "tf-ci-test-update")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithNameCreate = `

data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
data "tencentcloud_availability_zones" "default" {
}
resource "tencentcloud_instance" "foo" {
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    
    lifecycle {
        ignore_changes = [instance_type]
    }
}

`
const testAccCvmInstanceResource_WithNameChange1 = `

data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
data "tencentcloud_availability_zones" "default" {
}
resource "tencentcloud_instance" "foo" {
    instance_name = "tf-ci-test-update"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    
    lifecycle {
        ignore_changes = [instance_type]
    }
}

`

func TestAccTencentCloudInstanceResourceWithHostname(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithHostnameCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "hostname", "tf-ci-test")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithHostnameCreate = `

data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
resource "tencentcloud_instance" "foo" {
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    hostname = "tf-ci-test"
    system_disk_type = "CLOUD_PREMIUM"
    instance_name = "tf-ci-test"
}

`

func TestAccTencentCloudInstanceResourceWithSecurityGroup(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithSecurityGroupCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttrSet("tencentcloud_instance.foo", "instance_status"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "security_groups.#", "1")),
			},
			{
				Config: testAccCvmInstanceResource_WithSecurityGroupChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "security_groups.#", "2")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithSecurityGroupCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_name_regex = "Final"
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_instance" "foo" {
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    security_groups = ["sg-cm7fbbf3"]
    
    lifecycle {
        ignore_changes = [instance_type]
    }
}

`
const testAccCvmInstanceResource_WithSecurityGroupChange1 = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_name_regex = "Final"
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_instance" "foo" {
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    security_groups = ["sg-cm7fbbf3","sg-kensue7b"]
    
    lifecycle {
        ignore_changes = [instance_type]
    }
}

`

func TestAccTencentCloudInstanceResourceWithOrderSecurityGroup(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithOrderSecurityGroupCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_with_orderly_sg"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_with_orderly_sg", "orderly_security_groups.0", "sg-cm7fbbf3"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_with_orderly_sg", "orderly_security_groups.1", "sg-kensue7b"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_with_orderly_sg", "orderly_security_groups.2", "sg-05f7wnhn")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithOrderSecurityGroupCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_instance" "cvm_with_orderly_sg" {
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    orderly_security_groups = ["sg-cm7fbbf3","sg-kensue7b","sg-05f7wnhn"]
    instance_name = "test-orderly-sg-cvm"
    availability_zone = "ap-guangzhou-7"
}

`

func TestAccTencentCloudInstanceResourceWithTags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithTagsCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "tags.hello", "world"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "tags.happy", "hour")),
			},
			{
				Config: testAccCvmInstanceResource_WithTagsChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "tags.hello", "hello")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithTagsCreate = `

data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
data "tencentcloud_availability_zones" "default" {
}
resource "tencentcloud_instance" "foo" {
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    
    data_disks {
        data_disk_size = 150
        delete_with_instance = true
        data_disk_type = "CLOUD_PREMIUM"
    }
    
    tags = {
        hello = "world"
        happy = "hour"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_name = "tf-ci-test"
}

`
const testAccCvmInstanceResource_WithTagsChange1 = `

data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        values = ["S1","S2","S3","S4","S5"]
        name = "instance-family"
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
data "tencentcloud_availability_zones" "default" {
}
resource "tencentcloud_instance" "foo" {
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    
    data_disks {
        data_disk_type = "CLOUD_PREMIUM"
        data_disk_size = 150
        delete_with_instance = true
    }
    
    tags = {
        hello = "hello"
    }
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    instance_name = "tf-ci-test"
}

`

func TestAccTencentCloudInstanceResourceWithPlacementGroup(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithPlacementGroupCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING"), resource.TestCheckResourceAttrSet("tencentcloud_instance.foo", "placement_group_id")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithPlacementGroupCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_name_regex = "Final"
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
}
resource "tencentcloud_instance" "foo" {
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    placement_group_id = "ps-1y147q03"
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
}

`

func TestAccTencentCloudInstanceResourceWithSpotpaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithSpotpaidCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithSpotpaidCreate = `

data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
}
data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
resource "tencentcloud_instance" "foo" {
    spot_instance_type = "ONE-TIME"
    spot_max_price = 0.5
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    instance_charge_type = "SPOTPAID"
}

`

func TestAccTencentCloudInstanceResourceWithDataDiskOrder(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_WithDataDiskOrderCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.0.data_disk_size", "70"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.1.data_disk_size", "50"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "data_disks.2.data_disk_size", "100")),
			},
		},
	})
}

const testAccCvmInstanceResource_WithDataDiskOrderCreate = `

data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_name_regex = "Final"
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-with-cbs-order-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-with-cbs-order-subnet"
    cidr_block = "10.0.0.0/16"
    availability_zone = "ap-guangzhou-7"
}
resource "tencentcloud_instance" "foo" {
    vpc_id = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
    
    data_disks {
        delete_with_instance = true
        data_disk_size = 70
        data_disk_type = "CLOUD_PREMIUM"
    }
    data_disks {
        data_disk_type = "CLOUD_PREMIUM"
        delete_with_instance = true
        data_disk_size = 50
    }
    data_disks {
        data_disk_size = 100
        data_disk_type = "CLOUD_PREMIUM"
        delete_with_instance = true
    }
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
}

`

func TestAccTencentCloudInstanceResourceDataDiskByCbs(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_DataDiskByCbsCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.cvm_add_data_disk_by_cbs"), resource.TestCheckResourceAttr("tencentcloud_instance.cvm_add_data_disk_by_cbs", "instance_status", "RUNNING")),
			},
		},
	})
}

const testAccCvmInstanceResource_DataDiskByCbsCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["SA2"]
    }
    filter {
        values = ["ap-guangzhou-7"]
        name = "zone"
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-attach-cbs-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    vpc_id = tencentcloud_vpc.vpc.id
    name = "cvm-attach-cbs-subnet"
    cidr_block = "10.0.0.0/16"
    availability_zone = "ap-guangzhou-7"
}
resource "tencentcloud_instance" "cvm_add_data_disk_by_cbs" {
    subnet_id = tencentcloud_subnet.subnet.id
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
    instance_name = "cvm-add-data-disk-by-cbs"
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    vpc_id = tencentcloud_vpc.vpc.id
}
resource "tencentcloud_cbs_storage" "cbs_disk2" {
    storage_size = 100
    availability_zone = "ap-guangzhou-7"
    project_id = 0
    encrypt = true
    storage_name = "cbs_disk2"
    storage_type = "CLOUD_SSD"
}
resource "tencentcloud_cbs_storage" "cbs_disk1" {
    availability_zone = "ap-guangzhou-7"
    project_id = 0
    encrypt = true
    storage_name = "cbs_disk1"
    storage_type = "CLOUD_SSD"
    storage_size = 200
}
resource "tencentcloud_cbs_storage_attachment" "attachment_cbs_disk1" {
    storage_id = tencentcloud_cbs_storage.cbs_disk1.id
    instance_id = tencentcloud_instance.cvm_add_data_disk_by_cbs.id
}
resource "tencentcloud_cbs_storage_attachment" "attachment_cbs_disk2" {
    storage_id = tencentcloud_cbs_storage.cbs_disk2.id
    instance_id = tencentcloud_instance.cvm_add_data_disk_by_cbs.id
}

`

func TestAccTencentCloudInstanceResourceWithLocalDisk(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_LocalDiskCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.local_disk"),
					resource.TestCheckResourceAttr("tencentcloud_instance.local_disk", "instance_status", "RUNNING")),
			},
		},
	})
}

const testAccCvmInstanceResource_LocalDiskCreate = `
data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet"
  cidr_block        = "10.0.0.0/16"
  availability_zone = "ap-guangzhou-6"
}

resource "tencentcloud_instance" "local_disk" {
  instance_name     = "tf-example"
  availability_zone = "ap-guangzhou-6"
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = "IT5.4XLARGE64"
  system_disk_type  = "LOCAL_BASIC"
  system_disk_size  = 50
  hostname          = "user"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
    data_disk_name = "tf-test1"
  }

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 60
    encrypt        = false
    data_disk_name = "tf-test2"
  }

  tags = {
    tagKey = "tagValue"
  }
}
`

func TestAccTencentCloudNeedFixInstancePostpaidToPrepaid(t *testing.T) {

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstancePostPaid,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
			{
				Config: testAccTencentCloudInstanceBasicToPrepaid,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_period", "1"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
		},
	})
}

func TestAccTencentCloudNeedFixCvmInstanceResource_PrepaidToPostpaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_PrepaidToPostpaidCreate,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_charge_type_prepaid_period", "1"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_charge_type_prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_charge_type", "PREPAID")),
			},
			{
				Config: testAccCvmInstanceResource_PrepaidToPostpaidChange1,
				Check:  resource.ComposeTestCheckFunc(testAccCheckCvmInstanceExists("tencentcloud_instance.foo"), resource.TestCheckResourceAttr("tencentcloud_instance.foo", "instance_status", "RUNNING")),
			},
		},
	})
}

const testAccCvmInstanceResource_PrepaidToPostpaidCreate = `

data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1"]
    }
    cpu_core_count = 2
    memory_size = 2
}
resource "tencentcloud_instance" "foo" {
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-3"
    instance_charge_type = "PREPAID"
    instance_charge_type_prepaid_period = 1
    instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    image_id = "img-2lr9q49h"
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    force_delete = true
}

`
const testAccCvmInstanceResource_PrepaidToPostpaidChange1 = `

data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1"]
    }
    cpu_core_count = 2
    memory_size = 2
}
resource "tencentcloud_instance" "foo" {
    instance_name = "tf-ci-test"
    availability_zone = "ap-guangzhou-3"
    
    lifecycle {
        ignore_changes = [instance_type]
    }
    image_id = "img-2lr9q49h"
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    force_delete = true
}

`

func testAccCheckSecurityGroupExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no security group ID is set")
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		sg, err := service.DescribeSecurityGroup(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if sg == nil {
			return fmt.Errorf("security group not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckTencentCloudInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cvm instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cvm instance id is not set")
		}

		cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, err = cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("cvm instance id is not found")
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_instance" {
			continue
		}

		instance, err := cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, err = cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance != nil && *instance.InstanceState != svccvm.CVM_STATUS_SHUTDOWN && *instance.InstanceState != svccvm.CVM_STATUS_TERMINATING {
			return fmt.Errorf("cvm instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccTencentCloudInstanceBasic = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-basic-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-basic-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0

  tags = {
    hostname = "tci"
  }
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstancePrepaidBasic = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-prepaid-basic-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-prepaid-basic-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "cvm_prepaid_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
  instance_charge_type                    = "PREPAID"
  instance_charge_type_prepaid_period     = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  force_delete = true
  tags = {
    hostname = "tci"
  }
}
`

const testAccTencentCloudInstanceWithDataDiskOrder = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-with-cbs-order-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-with-cbs-order-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0

  data_disks {
    data_disk_size         = 100
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
  data_disks {
    data_disk_size         = 50
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
  data_disks {
    data_disk_size         = 70
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
}
`

const testAccTencentCloudInstanceAddDataDiskByCbs = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-attach-cbs-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-attach-cbs-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "cvm_add_data_disk_by_cbs" {
  instance_name     = "cvm-add-data-disk-by-cbs"
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}

resource "tencentcloud_cbs_storage" "cbs_disk1" {
	storage_name = "cbs_disk1"
	storage_type = "CLOUD_SSD"
	storage_size = 200
	availability_zone = var.availability_cvm_zone
	project_id = 0
	encrypt = false
}
resource "tencentcloud_cbs_storage" "cbs_disk2" {
	storage_name = "cbs_disk2"
	storage_type = "CLOUD_SSD"
	storage_size = 100
	availability_zone = var.availability_cvm_zone
	project_id = 0
	encrypt = false
}
resource "tencentcloud_cbs_storage_attachment" "attachment_cbs_disk1" {
	storage_id = tencentcloud_cbs_storage.cbs_disk1.id
	instance_id = tencentcloud_instance.cvm_add_data_disk_by_cbs.id
}
resource "tencentcloud_cbs_storage_attachment" "attachment_cbs_disk2" {
	storage_id = tencentcloud_cbs_storage.cbs_disk2.id
	instance_id = tencentcloud_instance.cvm_add_data_disk_by_cbs.id
}
`

const testAccTencentCloudInstancePostPaid = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 2
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "` + tcacctest.DefaultInsName + `"
  availability_zone = "` + tcacctest.DefaultAZone + `"
  image_id          = "` + tcacctest.DefaultTkeOSImageId + `"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  force_delete = true
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstanceBasicToPrepaid = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 2
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "` + tcacctest.DefaultInsName + `"
  availability_zone = "` + tcacctest.DefaultAZone + `"
  image_id          = "` + tcacctest.DefaultTkeOSImageId + `"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  instance_charge_type       = "PREPAID"
  instance_charge_type_prepaid_period = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  force_delete = true
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstanceModifyInstanceType = tcacctest.DefaultInstanceVariable + `
data "tencentcloud_instance_types" "new_type" {
	availability_zone = var.availability_cvm_zone
  
	cpu_core_count = 2
	memory_size    = 2
  }

resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-basic-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-basic-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.new_type.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0

  tags = {
    hostname = "tci"
  }
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstanceWithDataDisk = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    delete_with_instance  = true
	// encrypt = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    # data_disk_snapshot_id = "snap-nvzu3dmh"
    delete_with_instance  = true
  }

  disable_security_service = true
  disable_monitor_service  = true
  disable_automation_service = true
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstanceWithDataDiskUpdate = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }



  disable_security_service = true
  disable_monitor_service  = true
  disable_automation_service = true
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

func testAccTencentCloudInstanceWithNetworkFalse(hasPublicIp string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		hasPublicIp,
	)
}

func testAccTencentCloudInstanceWithNetwork(hasPublicIp string, maxBandWidthOut int64) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  internet_max_bandwidth_out = %d
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		maxBandWidthOut, hasPublicIp,
	)
}

const testAccTencentCloudInstanceWithPrivateIP = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-with-privateip-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-with-privateip-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  private_ip        = "10.0.0.123"
}
`

const testAccTencentCloudInstanceWithKeyPair_withoutKeyPair = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
	instance_name     = var.instance_name
	availability_zone = var.availability_cvm_zone
	image_id          = data.tencentcloud_images.default.images.0.image_id
	instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
	system_disk_type  = "CLOUD_PREMIUM"
	lifecycle {
		ignore_changes = [instance_type]
	}
}
`

func testAccTencentCloudInstanceWithKeyPair(keyIds string) string {

	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_key_pair" "key_pair_0" {
  key_name = "key_pair_0"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

resource "tencentcloud_key_pair" "key_pair_1" {
  key_name = "key_pair_1"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzwYE6KI8uULEvSNA2k1tlsLtMDe+x1Saw6yL3V1mk9NFws0K2BshYqsnP/BlYiGZv/Nld5xmGoA9LupOcUpyyGGSHZdBrMx1Dz9ajewe7kGowRWwwMAHTlzh9+iqeg/v6P5vW6EwK4hpGWgv06vGs3a8CzfbHu1YRbZAO/ysp3ymdL+vGvw/vzC0T+YwPMisn9wFD5FTlJ+Em6s9PzxqR/41t4YssmCwUV78ZoYL8CyB0emuB8wALvcXbdUVxMxpBEHd5U6ZP5+HPxU2WFbWqiFCuErLIZRuxFw8L/Ot+JOyNnadN1XU4crYDX5cML1i/ExXKVIDoBaLtgAJOpyeP"
}

resource "tencentcloud_key_pair" "key_pair_2" {
  key_name = "key_pair_2"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDJ1zyoM55pKxJptZBKceZSEypPN7BOunqBR1Qj3Tz5uImJ+dwfKzggu8PGcbHtuN8D2n1BH/GDkiGFaz/sIYUJWWZudcdut+ra32MqUvk953Sztf12rsFC1+lZ1CYEgon8Lt6ehxn+61tsS31yfUmpL1mq2vuca7J0NLdPMpxIYkGlifyAMISMmxi/m7gPYpbdZTmhQQS2aOhuLm+B4MwtTvT58jqNzIaFU0h5sqAvGQfzI5pcxwYvFTeQeXjJZfaYapDHN0MAg0b/vIWWNrDLv7dlv//OKBIaL0LIzIGQS8XXhF3HlyqfDuf3bjLBIKzYGSV/DRqlEsGBgzinJZXvJZug5oq1n2njDFsdXEvL6fYsP4WLvBLiQlceQ7oXi7m5nfrwFTaX+mpo7dUOR9AcyQ1AAgCcM67orB4E33ycaArGHtpjnCnWUjqQ+yCj4EXsD4yOL77wGsmhkbboVNnYAD9MJWsFP03hZE7p/RHY0C5NfLPT3mL45oZxBpC5mis="
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  key_ids           = %s
  system_disk_type  = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		keyIds,
	)
}

func testAccTencentCloudInstanceWithPassword(password string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  password                   = "%s"
  system_disk_type           = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		password,
	)
}

const testAccTencentCloudInstanceWithImageLogin = tcacctest.DefaultInstanceVariable + `
data "tencentcloud_images" "zoo" {
  image_type = ["PRIVATE_IMAGE"]
}
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.zoo.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  keep_image_login 			 = true
  system_disk_type           = "CLOUD_PREMIUM"
  disable_api_termination    = false
}
`

func testAccTencentCloudInstanceWithName(instanceName string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = "%s"
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		instanceName,
	)
}

const testAccTencentCloudInstanceWithHostname = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  hostname          = var.instance_name
  system_disk_type  = "CLOUD_PREMIUM"
}
`

func testAccTencentCloudInstanceWithSecurityGroup(ids string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  security_groups            = %s
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		ids,
	)
}

func testAccTencentCloudInstanceWithTags(tags string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  lifecycle {
	ignore_changes = [instance_type]
  }
  tags = %s
}
`,
		tags,
	)
}

const testAccTencentCloudInstanceWithPlacementGroup = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name      = var.instance_name
  availability_zone  = var.availability_cvm_zone
  image_id           = data.tencentcloud_images.default.images.0.image_id
  instance_type      = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type   = "CLOUD_PREMIUM"
  placement_group_id = "ps-1y147q03"
}
`

const testAccTencentCloudInstanceWithSpotpaid = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name        = var.instance_name
  availability_zone    = var.availability_cvm_zone
  image_id             = data.tencentcloud_images.default.images.0.image_id
  instance_type        = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type     = "CLOUD_PREMIUM"
  instance_charge_type = "SPOTPAID"
  spot_instance_type   = "ONE-TIME"
  spot_max_price       = "0.5"
}
`

const testAccTencentCloudInstanceOrderlySecurityGroups = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "cvm_with_orderly_sg" {
	instance_name              = "test-orderly-sg-cvm"
	availability_zone          = var.availability_cvm_zone
	image_id                   = data.tencentcloud_images.default.images.0.image_id
	instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
	system_disk_type           = "CLOUD_PREMIUM"
	orderly_security_groups    = ["sg-cm7fbbf3", "sg-kensue7b", "sg-05f7wnhn"]
}
`

func TestAccTencentCloudInstanceResourceHPC(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceHPC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.hpc"),
					resource.TestCheckResourceAttr("tencentcloud_instance.hpc", "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet("tencentcloud_instance.hpc", "hpc_cluster_id"),
				),
			},
		},
	})
}

const testAccTencentCloudInstanceHPC = `
locals {
  availability_zone = "ap-shanghai-5"
  instance_type     = "HCCS5.24XLARGE384"
}

resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
  zone   = local.availability_zone
  name   = "terraform-test"
  remark = "create for test"
}

data "tencentcloud_images" "default" {
  instance_type = local.instance_type
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = local.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_instance" "hpc" {
  instance_name     = "tf-example"
  availability_zone = local.availability_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = local.instance_type
  system_disk_type  = "LOCAL_BASIC"
  system_disk_size  = 440
  hostname          = "user"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  hpc_cluster_id    = tencentcloud_cvm_hpc_cluster.hpc_cluster.id
}
`

func TestAccTencentCloudNeedFixCvmInstanceResource_IPv6AddressType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_IPv6AddressType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "ipv6_address_type", "EIPv6"),
				),
			},
		},
	})
}

const testAccCvmInstanceResource_IPv6AddressType = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }
  cpu_core_count = 2
  memory_size    = 2
}
resource "tencentcloud_instance" "foo" {
  instance_name                           = "tf-ci-test"
  availability_zone                       = "ap-guangzhou-3"
  vpc_id                                  = "vpc-mvhjjprd"
  subnet_id                               = "subnet-2qfyfvv8"
  lifecycle {
    ignore_changes = [instance_type]
  }
  image_id          = "img-2lr9q49h"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  force_delete      = true
  orderly_security_groups = ["sg-05f7wnhn"]
  ipv6_address_type = "EIPv6"
}
`

func TestAccTencentCloudNeedFixCvmInstanceResource_ReleaseAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
			tcacctest.AccStepSetRegion(t, "ap-hongkong")
		},
		Providers:    acctest.AccProviders,
		CheckDestroy: testAccCheckCvmInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceResource_ReleaseAddress,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCvmInstanceExists("tencentcloud_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "ipv4_address_type", "HighQualityEIP"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "ipv6_address_type", "HighQualityEIPv6"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "ipv6_address_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "release_address", "true"),
					resource.TestCheckResourceAttr("tencentcloud_instance.foo", "ipv6_addresses.#", "1"),
				),
			},
		},
	})
}

const testAccCvmInstanceResource_ReleaseAddress = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S5"]
  }
  cpu_core_count = 2
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "tf-ci-test"
  availability_zone = "ap-hongkong-3"
  vpc_id            = "vpc-0rlhss58"
  subnet_id         = "subnet-73ckwxez"
  lifecycle {
    ignore_changes = [instance_type]
  }
  image_id                   = "img-l8og963d"
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  force_delete               = true
  internet_charge_type       = "BANDWIDTH_PACKAGE"
  bandwidth_package_id       = "bwp-rcql8f4c"
  ipv4_address_type          = "HighQualityEIP"
  ipv6_address_type          = "HighQualityEIPv6"
  ipv6_address_count         = 1
  allocate_public_ip         = true
  internet_max_bandwidth_out = 1
  release_address            = true
}
`
