package tencentcloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
)

func TestAccTencentCloudScfFunction_basic(t *testing.T) {
	var fnId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScfFunctionDestroy(&fnId),
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("first.zip", testAccScfFunctionBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestMatchResourceAttr("tencentcloud_scf_function.foo", "name", regexp.MustCompile(`ci-test-function`)),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "first.do_it_first"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "mem_size", "128"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python3.6"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vpc_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "subnet_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "namespace", "default"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "role", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_logset_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_topic_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "l5_enable", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "zip_file"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "triggers.#", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", "success"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_error", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "err_no"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "install_dependency", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status", SCF_FUNCTION_STATUS_ACTIVE),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status_desc", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eip_fixed", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "host", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vip", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "trigger_info.#", "0"),
				),
			},
			{
				Config: scfFunctionCodeEmbed("second.zip", testAccScfFunctionBasicUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", "test"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "second.do_it_second"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "mem_size", "1536"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "timeout", "300"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python3.6"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "l5_enable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eip_fixed", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "host", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vip", ""),
					resource.TestCheckNoResourceAttr("tencentcloud_scf_function.foo", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "tags.abc", "abc"),
				),
			},
			{
				ResourceName:      "tencentcloud_scf_function.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"zip_file",
				},
			},
		},
	})
}

func TestAccTencentCloudScfFunction_cos(t *testing.T) {
	var fnId string

	path := scfFunctionCodeFile("first.zip")

	// for unit test run on windows
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "\\\\", -1)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScfFunctionDestroy(&fnId),
		Steps: []resource.TestStep{
			{
				Config: testAccScfFunctionCosCode(path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestMatchResourceAttr("tencentcloud_scf_function.foo", "name", regexp.MustCompile(`ci-test-function`)),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "first.do_it_first"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "mem_size", "128"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python3.6"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vpc_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "subnet_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "namespace", "default"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "role", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_logset_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_topic_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "l5_enable", "false"),
					resource.TestCheckNoResourceAttr("tencentcloud_scf_function.foo", "zip_file"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "cos_bucket_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "cos_object_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "cos_bucket_region"),
					resource.TestCheckNoResourceAttr("tencentcloud_scf_function.foo", "demo_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", "success"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_error", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "err_no"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "install_dependency", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status", SCF_FUNCTION_STATUS_ACTIVE),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status_desc", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eip_fixed", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "host", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vip", ""),
				),
			},
			{
				Config: testAccScfFunctionCosUpdateCode(path),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python3.6"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", "success"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_error", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "err_no"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "install_dependency", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status", SCF_FUNCTION_STATUS_ACTIVE),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status_desc", ""),
				),
			},
		},
	})
}

func TestAccTencentCloudScfFunction_role(t *testing.T) {
	var fnId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScfFunctionDestroy(&fnId),
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("first.zip", testAccScfFunctionRole),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestMatchResourceAttr("tencentcloud_scf_function.foo", "name", regexp.MustCompile(`ci-test-function`)),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "first.do_it_first"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "mem_size", "128"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python3.6"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vpc_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "subnet_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "namespace", "default"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "role", "scf-role-test"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_logset_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_topic_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "l5_enable", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "zip_file"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", "success"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_error", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "err_no"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "install_dependency", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status", SCF_FUNCTION_STATUS_ACTIVE),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status_desc", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eip_fixed", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "host", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vip", ""),
				),
			},
			{
				Config: scfFunctionCodeEmbed("first.zip", testAccScfFunctionRoleUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "role", "scf-role-test-new"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "zip_file"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", "success"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_error", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "err_no"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status", SCF_FUNCTION_STATUS_ACTIVE),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status_desc", ""),
				),
			},
		},
	})
}

func TestAccTencentCloudScfFunction_trigger(t *testing.T) {
	var fnId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScfFunctionDestroy(&fnId),
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("first.zip", testAccScfFunctionTrigger),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestMatchResourceAttr("tencentcloud_scf_function.foo", "name", regexp.MustCompile(`ci-test-function`)),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "first.do_it_first"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "mem_size", "128"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python3.6"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vpc_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "subnet_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "namespace", "default"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "role", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_logset_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_topic_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "l5_enable", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "zip_file"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "triggers.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", "success"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_error", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "err_no"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "install_dependency", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status", SCF_FUNCTION_STATUS_ACTIVE),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status_desc", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eip_fixed", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "host", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vip", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "trigger_info.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "trigger_info.0.enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "trigger_info.0.create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "trigger_info.0.modify_time"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "trigger_info.1.enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "trigger_info.1.create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "trigger_info.1.modify_time"),
				),
			},
			{
				Config: scfFunctionCodeEmbed("first.zip", testAccScfFunctionTriggerUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "triggers.#", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status", SCF_FUNCTION_STATUS_ACTIVE),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status_desc", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "trigger_info.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "trigger_info.0.enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "trigger_info.0.create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "trigger_info.0.modify_time"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "trigger_info.1.enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "trigger_info.1.create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "trigger_info.1.modify_time"),
				),
			},
		},
	})
}

func TestAccTencentCloudScfFunction_customNamespace(t *testing.T) {
	var fnId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScfFunctionDestroy(&fnId),
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("first.zip", testAccScfFunctionCustomNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestMatchResourceAttr("tencentcloud_scf_function.foo", "name", regexp.MustCompile(`ci-test-function`)),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "first.do_it_first"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "mem_size", "128"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python3.6"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vpc_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "subnet_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "namespace", "ci-test-scf"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "role", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_logset_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_topic_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "l5_enable", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "zip_file"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "triggers.#", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", "success"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_error", ""),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "err_no"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "install_dependency", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status", SCF_FUNCTION_STATUS_ACTIVE),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "status_desc", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eip_fixed", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "host", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vip", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "trigger_info.#", "0"),
				),
			},
		},
	})
}

func testAccCheckScfFunctionExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return errors.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no scf function id is set")
		}

		split := strings.Split(rs.Primary.ID, "+")
		namespace, name := split[0], split[1]

		service := ScfService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		fn, err := service.DescribeFunction(context.TODO(), name, namespace)
		if err != nil {
			return err
		}

		if fn == nil {
			return errors.Errorf("scf function not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckScfFunctionDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		service := ScfService{client: client}

		split := strings.Split(*id, "+")
		if len(split) != 2 {
			return errors.Errorf("id is invalid; %s", *id)
		}
		namespace, name := split[0], split[1]

		fn, err := service.DescribeFunction(context.TODO(), name, namespace)
		if err != nil {
			return err
		}

		if fn != nil {
			return fmt.Errorf("scf function still exists")
		}

		return nil
	}
}

// base64 ecnode of zip file
var scfFunctionPy36Codes = map[string]string{
	"first.zip": `UEsDBBQAAAAIAN15nVCu+YHaVgAAAF4AAAAIABwAZmlyc3QucHlVVAkAA3EpqV5xKaledXgLAAEE
AAAAAAQAAAAAHctNCoAgEIbh/ZxicJWQHaDLuNCRJkhFv34gunvRu35e3Wpp4LWXTBQlcSxe4ZO2
jkEOyRg54LIz8VcT7C3/eor7Vvtwm3MpZmaD5II6SIf7Z/NYohdQSwECHgMUAAAACADdeZ1QrvmB
2lYAAABeAAAACAAYAAAAAAABAAAApIEAAAAAZmlyc3QucHlVVAUAA3EpqV51eAsAAQQAAAAABAAA
AABQSwUGAAAAAAEAAQBOAAAAmAAAAAAA`,
	"second.zip": `UEsDBBQAAAAIAOt5nVBo1OKHVgAAAGAAAAAJABwAc2Vjb25kLnB5VVQJAAOKKaleiimpXnV4CwAB
BAAAAAAEAAAAACXLQQqAIBBA0f2cYnClkB3Ay7jQiQx0RMcKorsX9dfvp1y5CW6dC0CkBSP7JL5T
4BI17VRkwiCncYBvjWS08vE5jly7vtSxsnKoZLEhWaEu9r/VbQAeUEsBAh4DFAAAAAgA63mdUGjU
4odWAAAAYAAAAAkAGAAAAAAAAQAAAKSBAAAAAHNlY29uZC5weVVUBQADiimpXnV4CwABBAAAAAAE
AAAAAFBLBQYAAAAAAQABAE8AAACZAAAAAAA=`,
}

func scfFunctionCodeFile(fileName string) string {
	fd, err := os.Create(fmt.Sprintf("%s/%s", os.TempDir(), fileName))
	if err != nil {
		panic(err)
	}
	data, _ := base64.StdEncoding.DecodeString(scfFunctionPy36Codes[fileName])
	_, _ = fd.Write(data)
	fd.Close()
	return fd.Name()
}

func scfFunctionCodeEmbed(fileName, cfg string) string {
	fileName = scfFunctionCodeFile(fileName)
	return fmt.Sprintf(cfg, scfFunctionRandomName(), fileName)
}

func scfFunctionRandomName() string {
	return fmt.Sprintf("ci-test-function-%d", rand.Intn(99999))
}

const testAccScfFunctionBasic = `
resource "tencentcloud_scf_function" "foo" {
  name      = "%s"
  handler   = "first.do_it_first"
  runtime   = "Python3.6"

  zip_file = "%s"

  tags = {
    "test" = "test"
  }
}
`

var testAccScfFunctionBasicUpdate = fmt.Sprintf(
	defaultVpcVariable+`
resource "tencentcloud_scf_function" "foo" {
  name      = "%s"
  handler   = "second.do_it_second"
  runtime   = "Python3.6"

  description = "test"
  mem_size    = 1536
  timeout     = 300

  environment = {
    "test" = "test"
  }

  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id

  tags = {
    "abc" = "abc"
  }

  zip_file = "%s"
}
`, "%s", "%s")

func testAccScfFunctionCosCode(codeSource string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos-%s"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket = tencentcloud_cos_bucket.foo.bucket
  key    = "/new_object_key.zip"
  source = "%s"
  acl    = "public-read"
}

resource "tencentcloud_scf_function" "foo" {
  name      = "ci-test-function-cos"
  handler   = "first.do_it_first"
  runtime   = "Python3.6"

  cos_bucket_name   = tencentcloud_cos_bucket.foo.id
  cos_object_name   = tencentcloud_cos_bucket_object.myobject.key
  cos_bucket_region = "ap-guangzhou"
}`, appid, codeSource)
}

func testAccScfFunctionCosUpdateCode(codeSource string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos-%s"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket = tencentcloud_cos_bucket.foo.bucket
  key    = "/new_object_key.zip"
  source = "%s"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket" "bar" {
  bucket = "scf-cos2-%s"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "bar" {
  bucket = tencentcloud_cos_bucket.bar.bucket
  key    = "/new_code.zip"
  source = "%s"
  acl    = "public-read"
}

resource "tencentcloud_scf_function" "foo" {
  name      = "%s"
  handler   = "first.do_it_first"
  runtime   = "Python3.6"

  cos_bucket_name   = tencentcloud_cos_bucket.bar.id
  cos_object_name   = tencentcloud_cos_bucket_object.bar.key
  cos_bucket_region = "ap-guangzhou"
}`, appid, codeSource, appid, codeSource, scfFunctionRandomName())
}

var testAccScfFunctionRole = fmt.Sprintf(`
variable "role_document" {
  default = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": ["name/sts:AssumeRole"],
      "effect": "allow",
      "principal": {"service":["scf.qcloud.com"]}
    }
  ]
}
EOF
}

resource "tencentcloud_cam_role" "foo" {
  name          = "scf-role-test"
  document      = var.role_document
  description   = "test"
  console_login = true
}

resource "tencentcloud_scf_function" "foo" {
  name    = "%s"
  handler = "first.do_it_first"
  runtime = "Python3.6"
  role    = tencentcloud_cam_role.foo.name

  zip_file = "%s"
}
`, "%s", "%s")

var testAccScfFunctionRoleUpdate = fmt.Sprintf(`
variable "role_document" {
  default = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": ["name/sts:AssumeRole"],
      "effect": "allow",
      "principal": {"service":["scf.qcloud.com"]}
    }
  ]
}
EOF
}

resource "tencentcloud_cam_role" "foo" {
  name          = "scf-role-test"
  document      = var.role_document
  description   = "test"
  console_login = true
}

resource "tencentcloud_cam_role" "bar" {
  name          = "scf-role-test-new"
  document      = var.role_document
  description   = "test"
  console_login = true
}

resource "tencentcloud_scf_function" "foo" {
  name    = "%s"
  handler = "first.do_it_first"
  runtime = "Python3.6"
  role    = tencentcloud_cam_role.bar.name

  zip_file = "%s"
}
`, "%s", "%s")

var testAccScfFunctionTrigger = fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos-%s"
  acl    = "public-read"
}

resource "tencentcloud_scf_function" "foo" {
  name      = "%s"
  handler   = "first.do_it_first"
  runtime   = "Python3.6"

  zip_file = "%s"

  triggers {
    name         = "ci-test-fn-api-gw"
    type         = "timer"
    trigger_desc = "*/5 * * * * * *"
  }

  triggers {
	name         = tencentcloud_cos_bucket.foo.id
	cos_region   = "ap-guangzhou"
    type         = "cos"
    trigger_desc = "{\"event\":\"cos:ObjectCreated:Put\",\"filter\":{\"Prefix\":\"\",\"Suffix\":\"\"}}"
  }
}
`, appid, "%s", "%s")

var testAccScfFunctionTriggerUpdate = fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos-%s"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket" "bar" {
  bucket = "scf-cos2-%s"
  acl    = "public-read"
}

resource "tencentcloud_scf_function" "foo" {
  name      = "%s"
  handler   = "first.do_it_first"
  runtime   = "Python3.6"

  zip_file = "%s"

  triggers {
    name         = "ci-test-fn-api-gw2"
    type         = "timer"
    trigger_desc = "*/5 * * * * * *"
  }

  triggers {
	name         = tencentcloud_cos_bucket.bar.id
	cos_region   = "ap-guangzhou"
    type         = "cos"
    trigger_desc = "{\"event\":\"cos:ObjectCreated:Put\",\"filter\":{\"Prefix\":\"\",\"Suffix\":\"\"}}"
  }
}
`, appid, appid, "%s", "%s")

const testAccScfFunctionCustomNamespace = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-scf"
  description = "test1"
}

resource "tencentcloud_scf_function" "foo" {
  namespace = tencentcloud_scf_namespace.foo.id
  name      = "%s"
  handler   = "first.do_it_first"
  runtime   = "Python3.6"

  zip_file = "%s"
}
`
