@echo off
chcp 65001 >nul

REM ==============================
REM Скрипт для кроссплатформенной сборки PGen CLI
REM ==============================

set VERSION=1.2.0
set APP_NAME=pgen
set BUILD_DIR=build

REM Переходим в корень проекта
cd ..

REM Создаем директорию для сборки, если нет
if not exist %BUILD_DIR% mkdir %BUILD_DIR%

echo Сборка %APP_NAME% v%VERSION%...
echo.

REM ==============================
REM Windows
REM ==============================
echo --- Windows amd64 ---
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-windows-amd64.exe .
if errorlevel 1 ( echo Ошибка сборки Windows amd64 & exit /b 1 )

echo --- Windows i386 ---
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=386
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-windows-i386.exe .
if errorlevel 1 ( echo Ошибка сборки Windows i386 & exit /b 1 )

echo --- Windows arm64 ---
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=arm64
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-windows-arm64.exe .
if errorlevel 1 ( echo Ошибка сборки Windows arm64 & exit /b 1 )

REM ==============================
REM Linux
REM ==============================
echo --- Linux amd64 ---
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-linux-amd64 .
if errorlevel 1 ( echo Ошибка сборки Linux amd64 & exit /b 1 )

echo --- Linux i386 ---
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=386
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-linux-i386 .
if errorlevel 1 ( echo Ошибка сборки Linux i386 & exit /b 1 )

echo --- Linux arm ---
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-linux-arm .
if errorlevel 1 ( echo Ошибка сборки Linux arm & exit /b 1 )

echo --- Linux arm64 ---
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm64
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-linux-arm64 .
if errorlevel 1 ( echo Ошибка сборки Linux arm64 & exit /b 1 )

REM ==============================
REM macOS
REM ==============================
echo --- macOS amd64 ---
set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-darwin-amd64 .
if errorlevel 1 ( echo Ошибка сборки macOS amd64 & exit /b 1 )

echo --- macOS arm64 ---
set CGO_ENABLED=0
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-darwin-arm64 .
if errorlevel 1 ( echo Ошибка сборки macOS arm64 & exit /b 1 )

REM ==============================
REM FreeBSD
REM ==============================
echo --- FreeBSD amd64 ---
set CGO_ENABLED=0
set GOOS=freebsd
set GOARCH=amd64
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-freebsd-amd64 .
if errorlevel 1 ( echo Ошибка сборки FreeBSD amd64 & exit /b 1 )

echo --- FreeBSD arm64 ---
set CGO_ENABLED=0
set GOOS=freebsd
set GOARCH=arm64
go build -ldflags "-s -w -X \"main.Version=%VERSION%\"" -o %BUILD_DIR%\%APP_NAME%-%VERSION%-freebsd-arm64 .
if errorlevel 1 ( echo Ошибка сборки FreeBSD arm64 & exit /b 1 )

echo.
echo 🎉 Сборка завершена успешно!
echo Файлы находятся в директории %BUILD_DIR%\
dir %BUILD_DIR%\
exit /b 0
