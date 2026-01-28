struct MarketDataService {
	providerName string
}
func NewMarketDataService() *MarketDataService {
    return &MarketDataService{}
}

func main() {
    app := fx.New(
        fx.Provide(NewMarketDataService),
        fx.Invoke(RegisterGRPC),
    )
    app.Run()
}