package device

type Device struct {
	ID     int    // Уникальный идентификатор
	UserID int    // Владелец устройства
	Name   string // Название устройства (например, "Кухонный ESP")
	Type   string // Тип (ESP8266, NeoPixel, Arduino и т.д.)
	Status string // online / offline / error
}
