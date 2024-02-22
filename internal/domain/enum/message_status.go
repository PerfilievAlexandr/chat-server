package messageStatus

import (
	"database/sql/driver"
	"fmt"
)

type MessageStatus string

const (
	NEW     MessageStatus = "NEW"
	EDITED  MessageStatus = "EDITED"
	DELETED MessageStatus = "DELETED"
)

func (m MessageStatus) string() string {
	return string(m)
}

func (m *MessageStatus) Scan(value interface{}) error {
	switch true {
	case value == NEW.string():
		*m = NEW
	case value == EDITED.string():
		*m = EDITED
	case value == DELETED.string():
		*m = DELETED
	default:
		return fmt.Errorf(`cannot parse:[%s] as message status`, value)
	}

	return nil
}

func (m MessageStatus) Value() (driver.Value, error) {
	return m.string(), nil
}
