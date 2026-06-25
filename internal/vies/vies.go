package vies

type VIESResult struct {
	IsValid     bool
	CountryCode string
	TradeName   string
	Address     string
	CleanedVAT  string
}

func Check(countryCode, VATNumber string) (*VIESResult, error) {
	res := &VIESResult{IsValid: true, CountryCode: countryCode, TradeName: "AXON CABLE", Address: "Višķu iela 21C, Daugavpils, LV-5410", CleanedVAT: VATNumber[2:]}

	return res, nil
}
