package clients

type DataResult struct {
	CurrentValue float64
	TargetValue  float64
	ShouldNotify bool
}

// Common client interface that will be implemented by the concrete types of clients
type Client interface {
	FetchAndExtractData(trackerCode string) (*DataResult, error)
}
