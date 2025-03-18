package restapi

type SetKeyRequest struct {
	Value any `binding:"required"`
}

type SetKeyResponse struct {
	Code     string
	LeaderId string
}

type GetKeyResponse struct {
	Value    any
	Code     string
	LeaderId string
}

type DeleteKeyResponse struct {
	Code     string
	LeaderId string
}

type ClusterInfoResponse struct {
	Code     string
	LeaderId string

	Info struct {
		currentTerm int
		// TODO: add fields
	}
}
