import os
import re
import json


lang_files = ["fr", "es", "it", "jp", "ko", "nl", "pt", "ro", "ru", "zh"]
folder_path = "./assets/lang"


def find_lang_translate():
    lang_dict = {}
    lang_pattern = re.compile('lang\.Translate\("([^"]*)"\)')
    go_files = 0
    for root, dirs, files in os.walk("."):
        for file_name in files:
            if file_name.endswith(".go"):
                go_files += 1
                file_path = os.path.join(root, file_name)
                with open(file_path, "r", encoding="utf-8") as file:
                    file_content = file.read()
                    matches = re.findall(lang_pattern, file_content)
                    for match in matches:
                        lang_dict[match] = ""

    print("total .go files checked {}".format(go_files))
    sorted_keys = sorted(lang_dict.keys())
    return {key: lang_dict[key] for key in sorted_keys}


def main():
    print("loading @lang.Translate() keys")
    new_lang_dict = find_lang_translate()
    print("total keys {}".format(len(new_lang_dict)))

    for lang in lang_files:
        print("updating {}".format(lang))
        file_path = "{}/{}.json".format(folder_path, lang)
        try:
            with open(file_path, encoding="utf-8") as file:
                lang_dict = json.load(file)
                current_keys = list(lang_dict.keys())
                new_keys = list(new_lang_dict.keys())

                # remove old keys
                for current_key in current_keys:
                    found = False
                    for new_key in new_keys:
                        if current_key == new_key:
                            found = True

                    if not found:
                        del lang_dict[current_key]

                # add new keys
                for new_key in new_keys:
                    found = False
                    for current_key in current_keys:
                        if current_key == new_key:
                            found = True

                    if not found:
                        lang_dict[new_key] = ""

        except FileNotFoundError:
            lang_dict = new_lang_dict

        with open(file_path, "w", encoding="utf-8") as file:
            json.dump(lang_dict, file, indent=2, ensure_ascii=False)

    print("done")


if __name__ == "__main__":
    main()
