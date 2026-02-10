@echo off
setlocal

REM Launch the Go project in a new WezTerm window.
REM Uses the local WezTerm binary in the wezterm folder.

set "PROJECT_DIR=%~dp0"
set "PROJECT_DIR=%PROJECT_DIR:~0,-1%"
set "WEZTERM_EXE=%PROJECT_DIR%\wezterm\wezterm.exe"

if not exist "%WEZTERM_EXE%" (
  echo WezTerm not found at "%WEZTERM_EXE%".
  exit /b 1
)

copy /Y "c:\Users\YaALi-pc\goProjs\tengine\.wezterm.lua" "%USERPROFILE%\.wezterm.lua"

pushd "%PROJECT_DIR%\core"
"%WEZTERM_EXE%" start --cwd "%PROJECT_DIR%\core" -- cmd.exe /k "go run ."
popd

endlocal
