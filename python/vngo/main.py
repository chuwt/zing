# coding:utf-8
"""
@Time :    2021/2/7 下午1:41
@Author:  chuwt
"""

import os
import importlib
from go_template import CtaTemplate

classes = dict()


def load_strategy(path):
    load_strategy_class_from_folder(path, "strategies")


def load_strategy_class_from_folder(path, module_name=""):
    for dirpath, dirnames, filenames in os.walk(str(path)):
        for filename in filenames:
            if filename.split(".")[-1] in ("py", "pyd", "so"):
                strategy_module_name = ".".join([module_name, filename.split(".")[0]])
                load_strategy_class_from_module(strategy_module_name)


def load_strategy_class_from_module(module_name):
    try:
        module = importlib.import_module(module_name)
        for name in dir(module):
            value = getattr(module, name)
            if isinstance(value, type) and \
                    value.__base__.__name__ == CtaTemplate.__name__ and \
                    value is not CtaTemplate:
                classes[value.__name__] = value
    except Exception as e:  # noqa
        print("not found", str(e))
        return "not found"


def get_strategy_instance(path, strategy_class_name, strategy_name, vt_symbol, setting):
    load_strategy(path)
    strategy_class = classes.get(strategy_class_name, None)
    if not strategy_class:
        print("none")
        return None
    return strategy_class(None, strategy_name, vt_symbol, setting)


if __name__ == "__main__":
    load_strategy("/Volumes/hdd1000gb/workspace/src/vngo/python/vngo/strategies")
    print(classes)
