$name = "vscode-portable-win64-1.92.0-49"
$arcName = "$name.zip"

Expand-Archive ".\soft\$arcName" ".\soft\temp" && 
New-Item -Path ".\vscode\win" -ItemType Directory && 
Move-Item -Path ".\soft\temp\$name\*" -Destination ".\vscode\win"&& 
Remove-Item ".\soft\temp" -Recurse;
