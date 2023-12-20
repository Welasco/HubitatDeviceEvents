package model

// {
//     "deviceTypeId": 104,
//     "displayName": "Kitchen",
//     "groupId": null,
//     "label": "Kitchen",
//     "lanId": null,
//     "version": 2,
//     "driverType": "sys",
//     "zigbeeId": null,
//     "currentStates": {},
//     "parentDeviceId": null,
//     "name": "GE Enbrighten Z-Wave Smart Switch",
//     "id": 252,
//     "deviceNetworkId": "3A"
// }

type Device struct {
	//gorm.Model
	Id          int    `json:"Id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayname"`
	Label       int    `json:"label"`
}