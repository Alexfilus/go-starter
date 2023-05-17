package event

const (
	SubjectStockUpdates = "stockupdates"
)

type Stock struct {
	Symbol string
	Price  int
}
