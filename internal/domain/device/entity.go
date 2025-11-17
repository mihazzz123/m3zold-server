package device

type Device struct {
	ID     string // Уникальный идентификатор
	UserID string // Владелец устройства
	Name   string // Название устройства (например, "Кухонный ESP")
	Type   string // Тип (ESP8266, NeoPixel, Arduino и т.д.)
	Status string // online / offline / error
}
