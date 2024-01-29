package types

type OBUDATA struct {
	OBUID int     `json:"obuid"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type Distance struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuid"`
	Unix  int64   `json:"unix"`
}
