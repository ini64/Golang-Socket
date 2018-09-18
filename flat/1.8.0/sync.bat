@echo off

ROBOCOPY golang\packet ..\..\golang\src\packet /MIR /XA:H /W:0 /R:1 /REG

