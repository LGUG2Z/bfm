package brew

const (
	RemoveAll = iota
	RemovePackageOnly
	RemovePackageAndRequired
)

const (
	AddAll = iota
	AddPackageOnly
	AddPackageAndRequired
)
