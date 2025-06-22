$username=$args[0]
Get-ADUser -identity $username  –Properties  "msDS-UserPasswordExpiryTimeComputed" |
Select-Object -Property @{Name="ExpiryDate";Expression={[datetime]::FromFileTime($_."msDS-UserPasswordExpiryTimeComputed")}}
 
 