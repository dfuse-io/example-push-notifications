package main

const (
	IOS = iota
	ANDROID
)

type DeviceToken struct {
	deviceType int
	Token      string
}

type Database struct {
	devicesTokens map[string]*DeviceToken
	lastCursor    string
}

func NewDatabase() *Database {
	return &Database{
		devicesTokens: map[string]*DeviceToken{},
	}
}

func (d *Database) OptInDeviceToken(eosAccountName string, deviceToken string, deviceType int) {
	d.devicesTokens[eosAccountName] = &DeviceToken{deviceType: deviceType, Token: deviceToken}
}

func (d *Database) FindDeviceToken(eosAccountName string) *DeviceToken {
	return d.devicesTokens[eosAccountName]
}

func (d *Database) StoreCursor(cursor string) {
	d.lastCursor = cursor
}

func (d *Database) LoadCursor() string {
	return d.lastCursor
}
