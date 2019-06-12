package msignotify

const (
	IOS     = "ios"
	ANDROID = "android"
)

type DeviceToken struct {
	DeviceType string
	Token      string
}
