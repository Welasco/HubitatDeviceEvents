
CREATE DATABASE IF NOT EXISTS hubitatdeviceevents;

-------------------------------------------------------------------------------------------------
-- {
--     "id": "768",
--     "name": "GE Enbrighten Z-Wave Smart Switch",
--     "label": "Kitchen Table",
--     "type": "GE Enbrighten Z-Wave Smart Switch",
--     "room": null
-- }

CREATE TABLE IF NOT EXISTS devices (
    id BIGINT NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    room VARCHAR(255) NULL
);

INSERT INTO devices (id, name, label, type, room) VALUES (1, 'device1', 'label1', 'type1', 'room1');

DELETE FROM devices WHERE id = ?;

UPDATE devices SET name = ?, displayname = ?, label = ? WHERE id = ?;

-------------------------------------------------------------------------------------------------
-- {
--     "content": {
--         "name": "switch",
--         "value": "on",
--         "displayName": "Office",
--         "deviceId": "187",
--         "descriptionText": "Office was turned on [digital]",
--         "unit": null,
--         "type": "digital",
--         "data": null
--     }
-- }

--{"content":{"name":"switch","value":"on","displayName":"Office","deviceId":"1","descriptionText":"Office was turned on [digital]","unit":null,"type":"digital","data":null}}

CREATE TABLE IF NOT EXISTS deviceevents (
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL,
    displayName VARCHAR(255) NOT NULL,
    deviceId BIGINT NOT NULL,
    descriptionText VARCHAR(255),
    unit VARCHAR(255) NULL,
    type VARCHAR(255) NULL,
    data VARCHAR(255) NULL,
    CONSTRAINT deviceId FOREIGN KEY (deviceId) REFERENCES devices (id)
);

INSERT INTO deviceevents (name, value, displayName, deviceId, descriptionText, unit, type, data) VALUES ('device1', 'value1', 'displayName1', 1, 'description1', NULL, 'type1', 'data1');
