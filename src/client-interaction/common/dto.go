package common

type SetKeyRequest struct {
	Key   string
	Value any
}

type SetKeyResponse struct {
	Code     string
	LeaderId string
}

type GetKeyRequest struct {
	Key string
}

type GetKeyResponse struct {
	Value    any
	Code     string
	LeaderId string
}

type DeleteKeyRequest struct {
	Key string
}

type DeleteKeyResponse struct {
	Code     string
	LeaderId string
}
