@echo off

setlocal

set GOOS=linux

if "%1"=="" (
    echo Please specify a build target:
    echo - mocktifactory
    goto end
)

if "%1"=="mocktifactory" (
    go build -o bin/linux/mocktifactory ./cmd/server
    if %errorlevel% neq 0 goto end
    docker build -t 172.30.1.1:5000/myproject/mocktifactory:latest -f ./build/Dockerfile ./bin/linux
    if %errorlevel% neq 0 goto end
    docker push 172.30.1.1:5000/myproject/mocktifactory:latest
    if %errorlevel% neq 0 goto end
    oc delete deployment mocktifactory
    oc apply -f deployments/minishift/mocktifactory.yaml
    if %errorlevel% neq 0 goto end
) else (
    echo Unknown build target %1
)

:end
endlocal
exit /b 0

