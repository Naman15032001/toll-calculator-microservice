package types

type OBUDATA struct {
	OBUID int     `json:"obuid"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}