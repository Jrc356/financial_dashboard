package tax

type Bucket string

const (
	TaxDeferred Bucket = "tax-deferred"
	Roth        Bucket = "roth"
	Taxable     Bucket = "taxable"
)

func IsValidTaxBucket(tb Bucket) bool {
	switch tb {
	case TaxDeferred:
		return true
	case Roth:
		return true
	case Taxable:
		return true
	default:
		return false
	}
}
