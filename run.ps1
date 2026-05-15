param (
    [Parameter(Mandatory=$true, Position=0)]
    [string]$Pattern
)

$targetDir = Get-ChildItem -Directory | Where-Object { 
    $_.Name -like "$Pattern*" -or $_.Name -like "*$Pattern*" 
} | Select-Object -First 1

if ($null -eq $targetDir) {
    Write-Host "Error: Pattern not found matching '$Pattern'" -ForegroundColor Red
    Write-Host "Please check your folder names and try again." -ForegroundColor Yellow
    Exit
}

$mainFile = Join-Path $targetDir.FullName "main.go"
if (-not (Test-Path $mainFile)) {
    Write-Host "Error: Found folder $($targetDir.Name) but main.go does not exist." -ForegroundColor Red
    Exit
}

Write-Host "Running pattern: " -NoNewline -ForegroundColor Cyan
Write-Host "$($targetDir.Name)" -ForegroundColor Green

go run $mainFile
