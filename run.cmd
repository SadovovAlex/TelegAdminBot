SET CGO_ENABLED=1
SET WORKDIR=d:\PROGRAMMING\GO\botrecept

echo Change to the directory where the Go code is located
cd %WORKDIR%\cmd\bot

echo Check if the 'run' directory exists, if not, create it
IF NOT EXIST "%WORKDIR%\run" (
    mkdir "%WORKDIR%\run"
)
IF NOT EXIST "%WORKDIR%\run\locales" (
    mkdir "%WORKDIR%\run\locales"
)
IF NOT EXIST "%WORKDIR%\run\database" (
    mkdir "%WORKDIR%\run\database"
)
echo Copy .env
copy %WORKDIR%\.env %WORKDIR%\run\.env 
copy %WORKDIR%\internal\localization\locales %WORKDIR%\run\locales
copy %WORKDIR%\recept.json %WORKDIR%\run\recept.json

echo Build the Go application
go build -x -o %WORKDIR%\run\botrecept.exe main.go setup.go

echo Run the application with the environment variables file
cd %WORKDIR%\run
botrecept.exe

pause
