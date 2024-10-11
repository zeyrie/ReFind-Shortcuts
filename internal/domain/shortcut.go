package domain

type Shortcut struct {
	ID          int64
	Value       string
	Description string
	Note        *string
	IconUrl     *string
	Category    int64
}
