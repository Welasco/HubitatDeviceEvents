package model

// {
//     "content": {
//         "name": "switch",
//         "value": "off",
//         "displayName": "Office",
//         "deviceId": "187",
//         "descriptionText": "Office was turned off [digital]",
//         "unit": null,
//         "type": "digital",
//         "data": null
//     }
// }

type Content struct {
	DeviceEvent DeviceEvent `json:"content"`
}

type DeviceEvent struct {
	TimeStamp       string
	Name            string `json:"name"`
	Value           string `json:"value"`
	DisplayName     string `json:"displayName"`
	DeviceId        string `json:"deviceId"`
	DescriptionText string `json:"descriptionText"`
	Unit            string `json:"unit"`
	Type            string `json:"type"`
	Data            string `json:"data"`
}
