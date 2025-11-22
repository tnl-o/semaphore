param(
    [string]$RepoRoot = (Resolve-Path "$PSScriptRoot\..").Path
)

Write-Host "=== Cursor terminal fix ===" -ForegroundColor Cyan

# 1. Set execution policy for current user (allows running local scripts)
Write-Host "Setting execution policy to RemoteSigned for CurrentUser..."
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force

# 2. Ensure PowerShell profile exists with UTF-8 settings
if (-not (Test-Path -LiteralPath $PROFILE)) {
    Write-Host "Creating PowerShell profile: $PROFILE"
    New-Item -ItemType File -Path $PROFILE -Force | Out-Null
} else {
    Write-Host "PowerShell profile already exists, updating..."
}

$profileSnippet = @'
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$PSDefaultParameterValues['*:Encoding'] = 'utf8'
if ((Get-Culture).TextInfo.OEMCodePage -ne 65001) {
    chcp 65001 | Out-Null
}
'@

if (-not (Get-Content -LiteralPath $PROFILE -Raw | Select-String -SimpleMatch '[Console]::OutputEncoding')) {
    Add-Content -LiteralPath $PROFILE -Value "`r`n# Cursor fix`r`n$profileSnippet"
    Write-Host "Added UTF-8 settings to PowerShell profile."
} else {
    Write-Host "UTF-8 settings already present in PowerShell profile."
}

# 3. Update .vscode/settings.json to force Command Prompt
$settingsDir = Join-Path $RepoRoot ".vscode"
if (-not (Test-Path -LiteralPath $settingsDir)) {
    New-Item -ItemType Directory -Path $settingsDir -Force | Out-Null
}
$settingsPath = Join-Path $settingsDir "settings.json"

function ConvertFrom-JsonWithComments {
    param([string]$JsonText)
    $clean = $JsonText -split "`n" | ForEach-Object {
        $line = $_.Trim()
        if ($line -like "//*" -or $line -like "//*") {
            return ""
        }
        if ($line -match "^\s*//") {
            return ""
        }
        return $_
    } | Where-Object { $_ -ne "" } | Out-String
    if ([string]::IsNullOrWhiteSpace($clean)) {
        return [ordered]@{}
    }
    return $clean | ConvertFrom-Json
}

if (Test-Path -LiteralPath $settingsPath) {
    Write-Host "Parsing existing settings with comment stripping..."
    $rawSettings = Get-Content -LiteralPath $settingsPath -Raw
    $settingsJson = ConvertFrom-JsonWithComments -JsonText $rawSettings
} else {
    $settingsJson = [ordered]@{}
}

if (-not ($settingsJson.PSObject.Properties.Name -contains "terminal.integrated.defaultProfile.windows")) {
    $settingsJson | Add-Member -NotePropertyName "terminal.integrated.defaultProfile.windows" -NotePropertyValue "Command Prompt"
} else {
    $settingsJson."terminal.integrated.defaultProfile.windows" = "Command Prompt"
}

$settingsJson | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $settingsPath -Encoding UTF8
Write-Host "Updated $settingsPath to enforce Command Prompt terminal (comments removed during rewrite)."

Write-Host "Done. Restart Cursor to apply the changes." -ForegroundColor Green

