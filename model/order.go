package model


type Order struct {
	Ext1       string    `json:"ext1"`
	FinishTime int64     `json:"finishTime"`
	OrderEmt   int64     `json:"orderEmt"`
	OrderID    int64     `json:"orderId"`
	OrderTime  int64     `json:"orderTime"`
	ParentID   int64     `json:"parentId"`
	PayMonth   int64     `json:"payMonth"`
	Plus       int64     `json:"plus"`
	PopID      int64     `json:"popId"`
	SkuList    []SkuList `json:"skuList"`
	UnionID    int64     `json:"unionId"`
	ValidCode  int64     `json:"validCode"`
}

type SkuList struct {
	ActualCosPrice    float64 `json:"actualCosPrice"`
	ActualFee         float64 `json:"actualFee"`
	Cid1              int64   `json:"cid1"`
	Cid2              int64   `json:"cid2"`
	Cid3              int64   `json:"cid3"`
	CommissionRate    float64 `json:"commissionRate"`
	EstimateCosPrice  float64 `json:"estimateCosPrice"`
	EstimateFee       float64 `json:"estimateFee"`
	Ext1              string  `json:"ext1"`
	FinalRate         float64 `json:"finalRate"`
	FrozenSkuNum      int64   `json:"frozenSkuNum"`
	PayMonth          int64   `json:"payMonth"`
	PID               string  `json:"pid"`
	PopID             int64   `json:"popId"`
	PositionID        int64   `json:"positionId"`
	Price             float64 `json:"price"`
	SiteID            int64   `json:"siteId"`
	SkuID             int64   `json:"skuId"`
	SkuName           string  `json:"skuName"`
	SkuNum            int64   `json:"skuNum"`
	SkuReturnNum      int64   `json:"skuReturnNum"`
	SubSideRate       float64 `json:"subSideRate"`
	SubUnionID        string  `json:"subUnionId"`
	SubsidyRate       float64 `json:"subsidyRate"`
	TraceType         int64   `json:"traceType"`
	UnionAlias        string  `json:"unionAlias"`
	UnionTag          string  `json:"unionTag"`
	UnionTrafficGroup int64   `json:"unionTrafficGroup"`
	ValidCode         int64   `json:"validCode"`
}