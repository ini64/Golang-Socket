@echo off

IF EXIST "golang" (
	rmdir /Q /S golang
)

IF EXIST "c#" (
	rmdir /Q /S c#
)


@echo on

flatc -n -o c# ../fbs/Client.fbs
flatc -g -o golang ../fbs/Client.fbs
flatc -g -o golang ../fbs/Server.fbs




