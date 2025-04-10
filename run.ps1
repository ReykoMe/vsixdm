param (
  [string]$cmd = "run",
  [Parameter(ValueFromRemainingArguments = $true)]
  [string[]]$args
)

switch ($cmd) {
  "dev" { .\scripts\win\dev.ps1 @args }
  "build" { .\scripts\win\build.ps1 }
  "init" { .\scripts\win\init.ps1 }
  default { Write-Host "Unknown command: $cmd" }
}