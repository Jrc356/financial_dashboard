package tax

type Bucket string

const (
	TaxDeferred Bucket = "tax-deferred"
	Roth        Bucket = "roth"
	Taxable     Bucket = "taxable"
)
