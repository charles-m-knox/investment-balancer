module go-fltk-investment-balancer

go 1.23.4

require (
	github.com/adrg/xdg v0.5.3
	github.com/charles-m-knox/investment-balancer v0.7.0
	github.com/pwiecz/go-fltk v0.0.0-20241102195420-49be1fb41870
	github.com/shopspring/decimal v1.4.0
)

require golang.org/x/sys v0.29.0 // indirect

replace github.com/charles-m-knox/investment-balancer => ../../
