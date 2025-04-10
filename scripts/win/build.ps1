$name = "vscode-portable-win64-1.92.0-49"
$arcName = "$name.zip"
$env:VSIXDM_STAGE=""

Remove-Item -Path .\bin -Recurse
Push-Location .\src
go build -o ..\bin\vsixdm.exe .
Pop-Location
Expand-Archive -Path ".\soft\$arcName" -DestinationPath ".\bin\code" -Force
Move-Item -Path ".\bin\code\$name\*" -Destination ".\bin\code" -Force
Remove-Item -Path ".\bin\code\$name" -Recurse -Force

