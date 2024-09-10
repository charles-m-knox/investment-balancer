module investment-balancer-v3

go 1.23.0

require (
	github.com/charles-m-knox/investment-balancer/pkg/balancer v0.0.0-00010101000000-000000000000
	github.com/shopspring/decimal v1.4.0
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/charles-m-knox/investment-balancer/pkg/balancer => ./pkg/balancer
