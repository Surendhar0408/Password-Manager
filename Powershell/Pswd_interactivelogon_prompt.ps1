$policy = Get-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" -Name "PasswordExpiryWarning"
$value = $policy.PasswordExpiryWarning

if ($value -eq 0) {
    Write-Output "The policy is not configured."
} else {
    Write-Output "$value"
}
