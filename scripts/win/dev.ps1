$env:VSIXDM_STAGE="DEV" 
Push-Location ./src
go run . @args
Pop-Location