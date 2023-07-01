package models

type TaxBucket string

const (
	TaxDeferred TaxBucket = "tax-deferred"
	Roth        TaxBucket = "roth"
	Taxable     TaxBucket = "taxable"
)
