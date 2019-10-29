package tencentcloud

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
				Config: scfFunctionCodeEmbed("main.py", testAccScfFunctionBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "name", "ci-test-function"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "main.do_it"),
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
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", ""),
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

func TestAccTencentCloudScfFunction_update(t *testing.T) {
	var fnId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScfFunctionDestroy(&fnId),
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("main.py", testAccScfFunctionBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "name", "ci-test-function"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "main.do_it"),
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
					resource.TestCheckNoResourceAttr("tencentcloud_scf_function.foo", "cos_bucket_name"),
					resource.TestCheckNoResourceAttr("tencentcloud_scf_function.foo", "cos_object_name"),
					resource.TestCheckNoResourceAttr("tencentcloud_scf_function.foo", "cos_bucket_region"),
					resource.TestCheckNoResourceAttr("tencentcloud_scf_function.foo", "demo_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", ""),
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
				Config: scfFunctionCodeEmbed("second.py", testAccScfFunctionUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", "test"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "second.do_it"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "mem_size", "1536"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "timeout", "300"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python2.7"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "l5_enable", "true"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eip_fixed", "false"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "eips.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "host", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vip", ""),
					resource.TestCheckNoResourceAttr("tencentcloud_scf_function.foo", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "tags.abc", "abc"),
				),
			},
		},
	})
}

func TestAccTencentCloudScfFunction_cos(t *testing.T) {
	var fnId string

	f, err := ioutil.TempFile(os.TempDir(), "scf")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		name := f.Name()
		f.Close()
		os.Remove(name)
	}()

	t.Log("file name:", f.Name())

	writer := zip.NewWriter(f)
	file, err := writer.Create("main.py")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(file, scfFunctionPy36Code); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	path := f.Name()
	// for unit test run on windows
	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "\\", "\\\\")
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
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "name", "ci-test-function"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "main.do_it"),
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
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", ""),
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
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", ""),
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
				Config: scfFunctionCodeEmbed("main.py", testAccScfFunctionRole),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "name", "ci-test-function"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "main.do_it"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "mem_size", "128"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "timeout", "3"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "runtime", "Python3.6"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "vpc_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "subnet_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "namespace", "default"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "role"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_logset_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "cls_topic_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "l5_enable", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "zip_file"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", ""),
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
				Config: scfFunctionCodeEmbed("main.py", testAccScfFunctionRoleUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "role"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "zip_file"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "modify_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_function.foo", "code_size"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", ""),
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
				Config: scfFunctionCodeEmbed("main.py", testAccScfFunctionTrigger),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "name", "ci-test-function"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "main.do_it"),
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
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", ""),
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
				Config: scfFunctionCodeEmbed("main.py", testAccScfFunctionTriggerUpdate),
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
				Config: scfFunctionCodeEmbed("main.py", testAccScfFunctionCustomNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScfFunctionExists("tencentcloud_scf_function.foo", &fnId),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "name", "ci-test-function"),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "handler", "main.do_it"),
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
					resource.TestCheckResourceAttr("tencentcloud_scf_function.foo", "code_result", ""),
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

const scfFunctionPy36Code = `
import json

def do_it(event,ctx):
    return "ci test"
`

func scfFunctionCodeEmbed(fileName, cfg string) string {
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)
	fileWriter, _ := writer.Create(fileName)
	_, _ = io.WriteString(fileWriter, scfFunctionPy36Code)
	writer.Close()

	return fmt.Sprintf(cfg, base64.StdEncoding.EncodeToString(buf.Bytes()))
}

const testAccScfFunctionBasic = `
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  zip_file = "%s"

  tags = {
    "test" = "test"
  }
}
`

var testAccScfFunctionUpdate = fmt.Sprintf(`
resource "tencentcloud_scf_function" "foo" {
  name        = "ci-test-function"
  handler     = "second.do_it"
  description = "test"
  mem_size    = 1536
  timeout     = 300
  
  environment = {
    "test" = "test"
  }

  runtime   = "Python2.7"
  vpc_id    = "%s"
  subnet_id = "%s"
  l5_enable = true

  tags = {
    "abc" = "abc"
  }

  zip_file = "%s"
}
`, DefaultVpcId, DefaultSubnetId, "%s")

func testAccScfFunctionCosCode(codeSource string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos1-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket = "${tencentcloud_cos_bucket.foo.bucket}"
  key    = "/new_object_key.zip"
  source = "%s"
  acl    = "public-read"
}

resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "${tencentcloud_cos_bucket.foo.id}"
  cos_object_name   = "${tencentcloud_cos_bucket_object.myobject.key}"
  cos_bucket_region = "ap-guangzhou"
}`, codeSource)
}

func testAccScfFunctionCosUpdateCode(codeSource string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos1-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket = "${tencentcloud_cos_bucket.foo.bucket}"
  key    = "/new_object_key.zip"
  source = "%s"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket" "bar" {
  bucket = "scf-cos2-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "bar" {
  bucket = "${tencentcloud_cos_bucket.bar.bucket}"
  key    = "/new_code.zip"
  source = "%s"
  acl    = "public-read"
}

resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "${tencentcloud_cos_bucket.bar.id}"
  cos_object_name   = "${tencentcloud_cos_bucket_object.bar.key}"
  cos_bucket_region = "ap-guangzhou"
}`, codeSource, codeSource)
}

const testAccScfFunctionRole = `
variable "role_document" {
  default = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}}]}"
}

resource "tencentcloud_cam_role" "foo" {
  name          = "ci-scf-role"
  document      = "${var.role_document}"
  description   = "ci-scf-role"
  console_login = true
}

resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"
  role    = "${tencentcloud_cam_role.foo.id}"

  zip_file = "%s"
}
`

const testAccScfFunctionRoleUpdate = `
variable "role_document" {
  default = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/100009461222:uin/100009461222\"]}}]}"
}

resource "tencentcloud_cam_role" "foo" {
  name          = "ci-scf-role"
  document      = "${var.role_document}"
  description   = "ci-scf-role"
  console_login = true
}

resource "tencentcloud_cam_role" "bar" {
  name          = "ci-scf-role-new"
  document      = "${var.role_document}"
  description   = "ci-scf-role"
  console_login = true
}

resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"
  role    = "${tencentcloud_cam_role.bar.id}"

  zip_file = "%s"
}
`

const testAccScfFunctionTrigger = `
resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos1-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  zip_file = "%s"

  triggers {
    name         = "ci-test-fn-api-gw"
    type         = "timer"
    trigger_desc = "*/5 * * * * * *"
  }

  triggers {
    name         = "${tencentcloud_cos_bucket.foo.id}"
    type         = "cos"
    trigger_desc = "{\"event\":\"cos:ObjectCreated:Put\",\"filter\":{\"Prefix\":\"\",\"Suffix\":\"\"}}"
  }
}
`

const testAccScfFunctionTriggerUpdate = `
resource "tencentcloud_cos_bucket" "foo" {
  bucket = "scf-cos1-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket" "bar" {
  bucket = "scf-cos2-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  zip_file = "%s"

  triggers {
    name         = "ci-test-fn-api-gw2"
    type         = "timer"
    trigger_desc = "*/5 * * * * * *"
  }

  triggers {
    name         = "${tencentcloud_cos_bucket.bar.id}"
    type         = "cos"
    trigger_desc = "{\"event\":\"cos:ObjectCreated:Put\",\"filter\":{\"Prefix\":\"\",\"Suffix\":\"\"}}"
  }
}
`

const testAccScfFunctionCustomNamespace = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace   = "ci-test-scf"
  description = "test1"
}

resource "tencentcloud_scf_function" "foo" {
  name      = "ci-test-function"
  namespace = "${tencentcloud_scf_namespace.foo.id}"
  handler   = "main.do_it"
  runtime   = "Python3.6"

  zip_file = "%s"

  tags = {
    "test" = "test"
  }
}
`
