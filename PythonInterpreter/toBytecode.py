import dis
import importlib
import types
import sys

source_py = sys.argv[1:]

source_py = source_py[0]

with open(source_py) as f_source:
    source_code = f_source.read()

byte_code = compile(source_code, source_py, "exec")
dis.dis(byte_code)

for x in byte_code.co_consts:
    if isinstance(x, types.CodeType):
        sub_byte_code = x
        dis.dis(sub_byte_code)