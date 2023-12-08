import yaml
import json
import re



#
# def move(dictionary):
#     for file_name, content in dictionary.items():
#
#         if file_name=="tencentcloud/extension_billing.go":
#             continue
#         # 打开文件并读取内容
#         print(file_name)
#         with open("../"+file_name, "r") as file:
#             lines = file.readlines()
#
#         modified_lines = []
#         inside_code_block = False
#         replacement_start=""
#         replacement_end=""
#         # 遍历文件内容的每一行
#         for line in lines:
#
#             start_match = re.search(r"//internal version: replace (\w+) begin", line)
#             end_match = re.search(r"//internal version: replace (\w+) end", line)
#             if start_match:
#                 inside_code_block = True
#                 replacement_start = start_match.group(1)
#                 continue
#             elif end_match:
#                 inside_code_block = False
#                 replacement_end = end_match.group(1)
#                 continue
#             elif not inside_code_block:
#                 if replacement_start ==replacement_end:
#                     if replacement_start in content:
#                         line=content[replacement_start]
#                     else:
#                         line=""
#                 modified_lines.append(line)
#
#             # modified_lines.append(line)
#
#         # 将修改后的内容写回文件
#         with open("../"+file_name, "w") as file:
#             file.writelines(modified_lines)
#
#     print("success replace")
#



def replace_code(dictionary, code):

    matches = re.finditer(r"//internal version: replace (\w+) begin.*?//internal version: replace \w+ end", code, flags=re.DOTALL)

    for match in matches:
        key = match.group(1)
        print(key)
        if key in dictionary:
            # 使用字典中的代码替换
            replacement_code = dictionary[key]
        else:
            replacement_code=""
        # testKey=re.search(r"//internal version: replace \w+ begin.*?//internal version: replace \w+ end")
        # print(testKey,replacement_code)
        mark_str="//internal version: replace %s begin.*?//internal version: replace %s end"%(key,key)
        code = re.sub(r"%s"%mark_str, replacement_code, code, flags=re.DOTALL)
    return code

def replace(dictionary):
    for file_name, content in dictionary.items():
        if file_name in "tencentcloud/extension_billing.go" or file_name in "tencentcloud/service_tencentcloud_billing.go":
            with open("../" + file_name, "w") as file:
                file.write(content["all"])
            continue

        if file_name in "go.sum":
            with open("../" + file_name, "a") as file:
                file.write(content)
        # 打开文件并读取内容
        with open("../" + file_name, "r") as file:
            code = file.read()

        # 替换代码
        replaced_code = replace_code(content, code)

        # 将修改后的内容写回文件
        with open("../" + file_name, "w") as file:
            file.write(replaced_code)

    print("Success replace")


def run():
    # 读取YAML文件
    yaml_file = "yunti-code.yaml"

    # 将YAML文件转换为JSON
    with open(yaml_file, "r") as f:
        yaml_data = yaml.safe_load(f)
    json_data = json.dumps(yaml_data)

    # 将JSON转换为字典
    dictionary = json.loads(json_data)

    replace(dictionary)




run()
