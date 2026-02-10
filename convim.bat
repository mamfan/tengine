@echo off
setlocal

set "PROJECT_DIR=%~dp0"
set "PROJECT_DIR=%PROJECT_DIR:~0,-1%"

pushd "%PROJECT_DIR%\core"
@REM go run .\cmd\img2matrix -in "%PROJECT_DIR%\aaa.jpg" -out "%PROJECT_DIR%\assets\aaa.timg.json"
go run .\cmd\img2matbin -in "%PROJECT_DIR%\aaa.jpg" -out "%PROJECT_DIR%\assets\aaa.timg"
popd

endlocal
