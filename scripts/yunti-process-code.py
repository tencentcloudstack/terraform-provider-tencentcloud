import yaml
import json
import re


# 处理每个文件
def replace(dictionary):
    for file_name, content in dictionary.items():
        # 读取文件内容
        if file_name=="tencentcloud/extension_billing.go":
            with open("../"+file_name, "w") as f:
                f.write(content["all"])
            continue

        with open("../"+file_name, "r") as f:
            file_content = f.readlines()

        # 修改文件内容
        new_file_content = []
        for line in file_content:
            match = re.search(r"//yunti mark (\w+)", line)
            if match:
                index = match.group(1)
                if index =="move":
                    continue
                if index in content:
                    line = content[index] +"\n"
            new_file_content.append(line)

        # 写入修改后的文件内容
        with open("../"+file_name, "w") as f:
            f.writelines(new_file_content)

def move(dictionary):
    start_marker = "//yunti mark move begin"
    end_marker = "//yunti mark move end"

    for file_name, _ in dictionary.items():
        if file_name=="tencentcloud/extension_billing.go":
            continue
        # 打开文件并读取内容
        with open("../"+file_name, "r") as file:
            lines = file.readlines()

        modified_lines = []
        inside_code_block = False
        # 遍历文件内容的每一行
        for line in lines:
            if start_marker in line:
                inside_code_block = True
            elif end_marker in line:
                inside_code_block = False
            elif not inside_code_block:
                modified_lines.append(line)

        # 将修改后的内容写回文件
        with open("../"+file_name, "w") as file:
            file.writelines(modified_lines)

    print("代码已成功删除。")

def run():
    # 读取YAML文件
    yaml_file = "yunti-code.yaml"

    # 将YAML文件转换为JSON
    with open(yaml_file, "r") as f:
        yaml_data = yaml.safe_load(f)
    json_data = json.dumps(yaml_data)

    # 将JSON转换为字典
    dictionary = json.loads(json_data)

    # 替换
    # replace(dictionary)

    # Move
    move(dictionary)

run()