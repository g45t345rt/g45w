import os
import re
import json


def main():
    lang_pattern = re.compile('lang\.Translate\("([^"]*)"\)')

    dict_lang = {}

    for root, dirs, files in os.walk("."):
        for file_name in files:
            if file_name.endswith(".go"):
                file_path = os.path.join(root, file_name)
                with open(file_path, 'r') as file:
                    print(file_name)
                    file_content = file.read()
                    matches = re.findall(lang_pattern, file_content)
                    for match in matches:
                        dict_lang[match] = match
                    print(len(matches), "keys")

    sorted_keys = sorted(dict_lang.keys())
    sorted_lang = {key: dict_lang[key] for key in sorted_keys}

    with open("assets/lang/gen_template.json", "w") as file:
        json.dump(sorted_lang, file, indent=2)


if __name__ == "__main__":
    main()
