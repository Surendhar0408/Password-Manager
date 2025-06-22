Import-Module ActiveDirectory
$username = $args[0]
#Get-ADUser -LDAPFilter "(pwdLastSet=0)" | Select SamAccountName|
#Export-CSV "C:\Crayonte\Dontremove\ChangePasswordAtNextLogon.csv" -NoTypeInformation -Encoding UTF8
get-aduser -identity $username -properties * | select pwdlastset | format-list #accountexpirationdate, accountexpires, accountlockouttime, badlogoncount, padpwdcount, lastbadpasswordattempt, lastlogondate, lockedout, passwordexpired, passwordlastset, pwdlastset | format-list
