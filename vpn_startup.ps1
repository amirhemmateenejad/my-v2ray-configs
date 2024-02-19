Add-Type -AssemblyName System.Windows.Forms
$FileBrowser = New-Object System.Windows.Forms.OpenFileDialog -Property @{ InitialDirectory = [Environment]::GetFolderPath('Desktop') }


$softwarePathKey = "HKEY_CURRENT_USER\SOFTWARE\ahp_vpn_app"
$filePath = ""
If (Test-Path -Path "Registry::$softwarePathKey") {
    $filePath = Get-ItemPropertyValue -Path "Registry::$softwarePathKey" -Name "filePath"
}
Else {
    $FileBrowser.ShowDialog() | Out-Null
    $filePath = $FileBrowser.FileName
    New-Item -Path "Registry::$softwarePathKey"
    New-ItemProperty -Path "Registry::$softwarePathKey" -Name "filePath" -Value "$filePath"
}



$hasInternet = Test-Connection -Quiet -ComputerName "viaq.ir"
Write-Host $hasInternet
if ($hasInternet -eq $false) {
    Write-Host "no internet"
    Exit
}

& "$filePath"

# $parentPath = (get-item $filePath).Directory

# Set-location $parentPath



git add .
git commit -m "new config released"
git push

ntfy publish ahp_vpn_download "update vpn please update"
