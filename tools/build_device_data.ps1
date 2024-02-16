function GetDeviceData {
    param (
        [string]$device_id
    )
    $deviceResponse = Invoke-WebRequest http://hubitat.localdomain/device/edit/$device_id
    $html_tbody = $deviceResponse.ParsedHtml.getElementsByTagName('tbody')
    foreach ($entry in $html_tbody){
        $device_data = $entry.IHTMLElement_innerText.Split("`n")
        $arr_device_data= @()
        $arr_device_data += " "
        $arr_device_data += "Device Name: " + $deviceResponse.ParsedHtml.title
        $arr_device_data += "Device_ID: " + $device_id
        foreach($i in $device_data){
            if([int][char]$i[0] -eq 13 ){
                break
            }
            $arr_device_data += $i
        }
    }
    $arr_device_data
}

$device_list = @("612","613","614","615","616","644","645","762","763","764","765","766","767")
foreach($device_id in $device_list){
    $device_data = GetDeviceData -device_id $device_id
    $device_data
}
