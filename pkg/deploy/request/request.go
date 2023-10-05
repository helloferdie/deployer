package request

// Deploy -
type Deploy struct {
	Repository string `query:"repository" loc:"repository" validate:"required"`
	Code       string `query:"code" loc:"code" validate:"required"`
	Async      bool   `query:"async" loc:"async"`
}
