package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cam_user
	resource.AddTestSweepers("tencentcloud_cam_user", &resource.Sweeper{
		Name: "tencentcloud_cam_user",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := CamService{client: client}

			users, err := service.DescribeUsersByFilter(ctx, nil)
			if err != nil {
				return err
			}

			for _, v := range users {
				if persistResource.MatchString(*v.Name) {
					continue
				}
				if !strings.Contains(*v.Name, "test") {
					continue
				}
				request := cam.NewDeleteUserRequest()
				request.Name = v.Name
				request.Force = helper.IntUint64(1)
				_, err := client.UseCamClient().DeleteUser(request)
				if err != nil {
					log.Printf("[%s] error, request: %s \nreason: %s ", request.GetAction(), request.ToJsonString(), err.Error())
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudCamUser_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUser_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamUserExists("tencentcloud_cam_user.user_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "name", "cam-user-test0"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "remark", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "console_login", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "need_reset_password", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "use_api", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "phone_num", "12345678910"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "country_code", "86"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "password", "Gail@1234"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "email", "1234@qq.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_basic", "uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_basic", "uid"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_basic", "secret_key"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_basic", "secret_id"),
				),
			},
			{
				ResourceName:            "tencentcloud_cam_user.user_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key", "secret_id", "password"},
			},
			{
				Config: testAccCamUser_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamUserExists("tencentcloud_cam_user.user_basic"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "name", "cam-user-test0"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "remark", "test1235"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "console_login", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "need_reset_password", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "use_api", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "phone_num", "13670093505"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "country_code", "72"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "password", "Gail@12346"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_basic", "email", "141515@qq.com"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_basic", "uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_basic", "uid"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_basic", "secret_key"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_basic", "secret_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudCamUser_nilPassword(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUser_nilPassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamUserExists("tencentcloud_cam_user.user_nil_password"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_nil_password", "name", "cam-user-testnil"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_nil_password", "remark", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_nil_password", "console_login", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_nil_password", "need_reset_password", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_nil_password", "use_api", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_nil_password", "phone_num", "12345678910"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_nil_password", "country_code", "86"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_nil_password", "uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_nil_password", "uid"),
				),
			},
		},
	})
}

func TestAccTencentCloudCamUser_withoutKey(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUser_withoutKey,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamUserExists("tencentcloud_cam_user.user_without_key"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_without_key", "name", "cam-user-testkey"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_without_key", "remark", "test"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_without_key", "console_login", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_without_key", "need_reset_password", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_without_key", "use_api", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_without_key", "phone_num", "12345678910"),
					resource.TestCheckResourceAttr("tencentcloud_cam_user.user_without_key", "country_code", "86"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_without_key", "uin"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_user.user_without_key", "uid"),
				),
			},
		},
	})
}
func testAccCheckCamUserDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_user" {
			continue
		}

		instance, err := camService.DescribeUserById(ctx, rs.Primary.ID)
		if err == nil && (instance != nil && instance.Response != nil && instance.Response.Uid != nil) {
			return fmt.Errorf("[CHECK][CAM user][Destroy] check: CAM user still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM user][Exists] check: CAM user %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM user][Exists] check: CAM user id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		instance, err := camService.DescribeUserById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil || instance.Response == nil || instance.Response.Uid == nil {
			return fmt.Errorf("[CHECK][CAM user][Exists] check: CAM user %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCamUser_basic = `
resource "tencentcloud_cam_user" "user_basic" {
  name                = "cam-user-test0"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "12345678910"
  country_code        = "86"
  email               = "1234@qq.com"
  force_delete		  = false
}
`

const testAccCamUser_update = `
resource "tencentcloud_cam_user" "user_basic" {
  name                = "cam-user-test0"
  remark              = "test1235"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@12346"
  phone_num           = "13670093505"
  country_code        = "72"
  email               = "141515@qq.com"
  force_delete		  = true
}
`
const testAccCamUser_nilPassword = `
resource "tencentcloud_cam_user" "user_nil_password" {
	name                = "cam-user-testnil"
	remark              = "test"
	console_login       = true
	use_api             = true
	need_reset_password = true
	phone_num           = "12345678910"
	country_code        = "86"
	email               = "141515@qq.com"
    force_delete		= true
}
`
const testAccCamUser_withoutKey = `
resource "tencentcloud_cam_user" "user_without_key" {
  name                = "cam-user-testkey"
  remark              = "test"
  console_login       = false
  use_api             = false
  need_reset_password = true
  phone_num           = "12345678910"
  country_code        = "86"
  email               = "141515@qq.com"
  force_delete		  = true
}
`
