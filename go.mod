module samplesort

go 1.12

require (
	github.com/bugra/kmeans v0.0.0-20140831011822-bf06fda928a7
	github.com/gonum/stat v0.0.0-20181125101827-41a0da705a5b // indirect
	github.com/jeremywohl/flatten v0.0.0-20180923035001-588fe0d4c603
	gonum.org/v1/gonum v0.0.0-20180816165407-929014505bf4
	gonum.org/v1/plot v0.0.0-20190515093506-e2840ee46a6b
)

replace samplesort/engine => ./engine

replace samplesort/parser => ./parser

replace samplesort/cache => ./cache

replace samplesort/analyze => ./analyze

replace samplesort/crypto => ./crypto

replace samplesort/collection => ./collection

replace samplesort/sample => ./sample
