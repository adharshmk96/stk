import os
from dataclasses import dataclass
from pathlib import Path
from typing import List, Tuple, Callable

BASE_PATH = Path('./')
WORKING_DIR = '../pkg/project/tpl'
OUT_PROJECT_PATH = Path(WORKING_DIR+'/project.go')
OUT_MODULE_PATH = Path(WORKING_DIR+'/modules.go')

IGNORE_DIRS = ['.git', 'mocks']
IGNORE_FILES = [
    'template.go', 
    "go.mod", 
    "go.sum",
    "database.db",
    "generate.py",
]

PLACEHOLDERS = { 
    "app": "{{ .AppName }}",
    "pkg": "{{ .PkgName }}",
    "mod": "{{ .ModName }}",
    "exported": "{{ .ExportedName }}",
}

def get_pkg_and_app_name_from_go_mod():
    with open('go.mod', 'r') as f:
        lines = f.readlines()
        for line in lines:
            if 'module' in line:
                pkg_name = line.split(' ')[1].strip('\n')
                app_name = pkg_name.split('/')[-1]
                return pkg_name, app_name
    return None, None

pkg_name, app_name = get_pkg_and_app_name_from_go_mod()
mod_name = "ping"
exported_mod_name = "Ping"

@dataclass
class TemplateMap:
    absolute_path: Path
    relative_path: str
    var_name: str
    var_suffix: str = ""

def replace_strings(content: str, replacements: dict) -> str:
    for old, new in replacements.items():
        content = content.replace(old, new)
    return content

def ignore_dir(dir: Path) -> bool:
    return any(ignore in str(dir) for ignore in IGNORE_DIRS)


def generate_var_name(relative_path: str) -> str:
    var_name = (
        relative_path
        .replace('.', '')
        .replace('/', '_')
        .replace('-', '')
        .strip('_')
        .upper()
    )
    return var_name

def find_template_map(base_path: Path, filename_condition: Callable[[str], bool], template_var_suffix: str) -> List[TemplateMap]:
    template_map = []
    for path in base_path.rglob('*'):
        # Check if it's a file, satisfies the filename condition, not in ignore dirs/files
        if (
            path.is_file() 
            and filename_condition(path.name) 
            and not ignore_dir(path.parent) 
            and path.name not in IGNORE_FILES
        ):
            relative_path = path.relative_to(base_path).as_posix()
            var_name = generate_var_name(relative_path)
            var_name += template_var_suffix
            template_map.append(TemplateMap(path, relative_path, var_name))  # Here, path is a Path object                
    return template_map


def generate_template(template_map: TemplateMap, replacements: dict) -> str:    
    template_content = "package tpl\n\n"
    for template in template_map:
        with template.absolute_path.open('r', encoding='utf-8') as input_file:
            content = replace_strings(input_file.read(), replacements)
            structure = f'var {template.var_name} = Template{{\n\tFilePath: "{template.relative_path}",\n\tContent: `{content}`,\n}}\n\n'
            template_content += structure
    return template_content

def write_template_to_file(target_path: Path, file_content: str, template_map: List[TemplateMap], var_name: str) -> None:
    with target_path.open('w', encoding='utf-8') as output_file:
        output_file.write(file_content)
        output_file.write(f"var {var_name} = []Template{{\n")
        for template in template_map:
            output_file.write(f'\t{template.var_name},\n')
        output_file.write("}\n")

def create_template(
    base_directory: Path, 
    target_file_path: Path, 
    filename_condition, 
    replacements: dict, 
    var_name: str,
    template_var_suffix: str
) -> None:
    template_map = find_template_map(base_directory, filename_condition, template_var_suffix)
    template_content = generate_template(template_map, replacements)
    write_template_to_file(target_file_path, template_content, template_map, var_name)


def remove_files(files):
    for file in files:
        try:
            os.remove(file)
        except:
            pass

if __name__ == '__main__':

    remove_files([OUT_PROJECT_PATH, OUT_MODULE_PATH])
    replacements = {
        # order is important, pkg should be before app_name
        pkg_name: PLACEHOLDERS['pkg'],        
        app_name: PLACEHOLDERS['app'],
        mod_name: PLACEHOLDERS['mod'],
        exported_mod_name: PLACEHOLDERS['exported']
    }

    print("Generating Project Templates...\n")

    create_template(
        BASE_PATH, 
        OUT_PROJECT_PATH, 
        filename_condition=lambda f: "ping" not in f,
        replacements=replacements,
        var_name="ProjectTemplates",
        template_var_suffix="_TPL"
    )

    print("Generating Module Templates...\n")
    
    create_template(
        BASE_PATH,
        OUT_MODULE_PATH,
        filename_condition=lambda f: "ping" in f,
        replacements=replacements,
        var_name="ModuleTemplates",
        template_var_suffix="_MOD_TPL"
    )

    print(f"All files have been written to {OUT_PROJECT_PATH}, {OUT_MODULE_PATH}.")
