package skeleton

type DefaultRequest struct {
	Data string
}

func (r DefaultRequest) Valid() bool {
	return true
}

func (r DefaultRequest) GetString() string {
	return r.Data
}
