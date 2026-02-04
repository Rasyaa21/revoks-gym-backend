param(
  [string]$BaseUrl = "http://localhost:3000",
  [string]$CreateUser = "true",
  [string]$RenewMembership = "true",
  [int]$RenewMonths = 1,
  [string]$Plan = "standard",
  [string]$Email = "",
  [string]$Password = "password123",
  [string]$Name = "Test User",
  [string]$Phone = "08123456789",
  [string]$Address = "Jakarta"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function To-Bool([string]$Value, [string]$ParamName) {
  if ($null -eq $Value) { return $false }
  $v = $Value.Trim().ToLowerInvariant()
  if ($v -match '^\$?(true|1|yes|y)$') { return $true }
  if ($v -match '^\$?(false|0|no|n)$') { return $false }
  throw "Invalid value for -${ParamName}: '$Value'. Use true/false or 1/0."
}

function Write-Section([string]$Title) {
  Write-Host ""
  Write-Host "=== $Title ===" -ForegroundColor Cyan
}

function Get-ErrorBody($err) {
  try {
    $resp = $err.Exception.Response
    if ($null -eq $resp) { return $null }
    $stream = $resp.GetResponseStream()
    if ($null -eq $stream) { return $null }
    $reader = New-Object System.IO.StreamReader($stream)
    return $reader.ReadToEnd()
  } catch {
    return $null
  }
}

function Invoke-Json(
  [Parameter(Mandatory=$true)][ValidateSet('GET','POST','PUT','DELETE','PATCH')][string]$Method,
  [Parameter(Mandatory=$true)][string]$Path,
  [hashtable]$Headers = @{},
  $Body = $null
) {
  $uri = "$BaseUrl$Path"
  try {
    if ($null -ne $Body) {
      $json = $Body | ConvertTo-Json -Depth 10
      return Invoke-RestMethod -Method $Method -Uri $uri -Headers $Headers -ContentType "application/json" -Body $json
    }
    return Invoke-RestMethod -Method $Method -Uri $uri -Headers $Headers
  } catch {
    $bodyText = Get-ErrorBody $_
    Write-Host "Request failed: $Method $Path" -ForegroundColor Red
    if ($bodyText) {
      Write-Host $bodyText -ForegroundColor DarkRed
    } else {
      Write-Host $_ -ForegroundColor DarkRed
    }
    throw
  }
}

$doCreateUser = To-Bool -Value $CreateUser -ParamName "CreateUser"
$doRenewMembership = To-Bool -Value $RenewMembership -ParamName "RenewMembership"

# Generate a unique email only when we are creating a user
if ($doCreateUser -and [string]::IsNullOrWhiteSpace($Email)) {
  $stamp = (Get-Date).ToString("yyyyMMddHHmmss")
  $Email = "test+$stamp@example.com"
}

if (-not $doCreateUser -and [string]::IsNullOrWhiteSpace($Email)) {
  throw "-Email is required when -CreateUser is false (so the script can login)."
}

Write-Section "Health"
$health = Invoke-RestMethod -Method Get -Uri "$BaseUrl/health"
$health | ConvertTo-Json -Depth 5

if ($doCreateUser) {
  Write-Section "Create User"
  $createBody = @{ name = $Name; email = $Email; password = $Password; phone = $Phone; address = $Address }
  try {
    $created = Invoke-Json -Method POST -Path "/api/v1/users" -Body $createBody
    $created | ConvertTo-Json -Depth 10
  } catch {
    Write-Host "Create user failed (maybe already exists). Continuing..." -ForegroundColor Yellow
  }
}

Write-Section "Login"
$login = Invoke-Json -Method POST -Path "/api/v1/auth/login" -Body @{ email = $Email; password = $Password }
$token = $login.data.access_token
if ([string]::IsNullOrWhiteSpace($token)) {
  throw "Login succeeded but access_token is missing."
}
Write-Host "Token prefix: $($token.Substring(0,20))..." -ForegroundColor Green
$authHeaders = @{ Authorization = "Bearer $token" }

Write-Section "Me"
$me = Invoke-Json -Method GET -Path "/api/v1/me" -Headers $authHeaders
$me.data | ConvertTo-Json -Depth 10

Write-Section "Membership (before renew)"
$membershipBefore = Invoke-Json -Method GET -Path "/api/v1/membership" -Headers $authHeaders
$membershipBefore.data | ConvertTo-Json -Depth 10

if ($doRenewMembership) {
  Write-Section "Membership Renew"
  $renew = Invoke-Json -Method POST -Path "/api/v1/membership/renew" -Headers $authHeaders -Body @{ months = $RenewMonths; plan = $Plan }
  $renew.data | ConvertTo-Json -Depth 10
}

Write-Section "QR Code"
$qr = $null
try {
  $qr = Invoke-Json -Method GET -Path "/api/v1/qr/code" -Headers $authHeaders
  $qr.data | ConvertTo-Json -Depth 10
} catch {
  Write-Host "QR code request failed. This is expected if membership is not ACTIVE." -ForegroundColor Yellow
}

if ($null -ne $qr -and $null -ne $qr.data -and -not [string]::IsNullOrWhiteSpace($qr.data.token)) {
  Write-Section "QR Scan (gate)"
  try {
    $scan = Invoke-Json -Method POST -Path "/api/v1/qr/scan" -Body @{ token = $qr.data.token }
    $scan.data | ConvertTo-Json -Depth 10
  } catch {
    Write-Host "QR scan failed. This may be expected if membership is expired or token is invalid." -ForegroundColor Yellow
  }
} else {
  Write-Host "Skipping QR scan because no QR token was returned." -ForegroundColor Yellow
}

Write-Section "Attendance History"
$att = Invoke-Json -Method GET -Path "/api/v1/attendance/history?limit=10" -Headers $authHeaders
$att.data | ConvertTo-Json -Depth 10

Write-Section "Workout Progress (create + list)"
$workoutCreate = Invoke-Json -Method POST -Path "/api/v1/workouts/progress" -Headers $authHeaders -Body @{ title = "Push day"; notes = "Bench + triceps" }
$workoutCreate.data | ConvertTo-Json -Depth 10
$workoutList = Invoke-Json -Method GET -Path "/api/v1/workouts/progress?limit=5" -Headers $authHeaders
$workoutList.data | ConvertTo-Json -Depth 10

Write-Section "Settings (get + update)"
$settings = Invoke-Json -Method GET -Path "/api/v1/settings" -Headers $authHeaders
$settings.data | ConvertTo-Json -Depth 10
$settingsUpdated = Invoke-Json -Method PUT -Path "/api/v1/settings" -Headers $authHeaders -Body @{ push_enabled = $false }
$settingsUpdated.data | ConvertTo-Json -Depth 10

Write-Section "Notifications"
$notif = Invoke-Json -Method GET -Path "/api/v1/notifications?limit=5" -Headers $authHeaders
$notif.data | ConvertTo-Json -Depth 10

Write-Section "Program/PT (may be empty without seed)"
$templatesFollowed = Invoke-Json -Method GET -Path "/api/v1/templates/followed" -Headers $authHeaders
$templatesFollowed.data | ConvertTo-Json -Depth 10
$targets = Invoke-Json -Method GET -Path "/api/v1/targets?period=weekly" -Headers $authHeaders
$targets.data | ConvertTo-Json -Depth 10
$pt = Invoke-Json -Method GET -Path "/api/v1/pt" -Headers $authHeaders
$pt.data | ConvertTo-Json -Depth 10

Write-Host "" 
Write-Host "Smoke test completed." -ForegroundColor Green
Write-Host "User email: $Email" -ForegroundColor Green
