use super::super::RequestBuilder;

pub struct StockQuoteRequestBuilder {
    ticker: String,
}

impl StockQuoteRequestBuilder {
    pub fn new(ticker: &str) -> StockQuoteRequestBuilder {
        StockQuoteRequestBuilder {
            ticker: ticker.to_owned(),
        }
    }
}

impl RequestBuilder for StockQuoteRequestBuilder {
    type Response = StockQuoteResponse;

    fn build(&self, base: &str) -> String {
        std::format!("{}/stock/{}/quote", base, self.ticker)
    }
}

#[derive(serde::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct StockQuoteResponse {
    pub symbol: String,
    pub company_name: String,
    pub primary_exchange: String,
    pub calculation_price: String,
    pub open: f64,
    pub open_time: i64,
    pub open_source: String,
    pub close: f64,
    pub close_time: i64,
    pub close_source: String,
    pub high: f64,
    pub high_time: i64,
    pub high_source: String,
    pub low: f64,
    pub low_time: i64,
    pub low_source: String,
    pub latest_price: f64,
    pub latest_source: String,
    pub latest_time: String,
    pub latest_update: i64,
    pub latest_volume: i64,
    pub iex_realtime_price: Option<f64>,
    pub iex_realtime_size: Option<f64>,
    pub iex_last_updated: Option<f64>,
    pub delayed_price: f64,
    pub delayed_price_time: i64,
    pub odd_lot_delayed_price: f64,
    pub odd_lot_delayed_price_time: i64,
    pub extended_price: f64,
    pub extended_change: f64,
    pub extended_change_percent: f64,
    pub extended_price_time: i64,
    pub previous_close: f64,
    pub previous_volume: i64,
    pub change: f64,
    pub change_percent: f64,
    pub volume: i64,
    pub iex_market_percent: Option<f64>,
    pub iex_volume: Option<f64>,
    pub avg_total_volume: i64,
    pub iex_bid_price: Option<f64>,
    pub iex_bid_size: Option<i64>,
    pub iex_ask_price: Option<f64>,
    pub iex_ask_size: Option<f64>,
    pub iex_open: f64,
    pub iex_open_time: i64,
    pub iex_close: f64,
    pub iex_close_time: i64,
    pub market_cap: i64,
    pub pe_ratio: f64,
    pub week52_high: f64,
    pub week52_low: f64,
    pub ytd_change: f64,
    pub last_trade_time: i64,
    #[serde(rename(deserialize = "isUSMarketOpen"))]
    pub is_us_market_open: bool,
}

mod tests {
    use super::super::super::Client;
    use super::*;

    #[test]
    fn test_deserialize_resp() -> Result<(), serde_json::Error> {
        let raw_data = r#"{
            "avgTotalVolume":119072197,
            "calculationPrice":"close",
            "change":-5.33,
            "changePercent":-0.03903,
            "close":134.51,
            "closeSource":"fafiiclo",
            "closeTime":1618852107068,
            "companyName":"Apple Inc",
            "delayedPrice":137.94,
            "delayedPriceTime":1675582700093,
            "extendedChange":-0.18,
            "extendedChangePercent":-0.00142,
            "extendedPrice":134.52,
            "extendedPriceTime":1673900572692,
            "high":143.08,
            "highSource":"l  amnue tep1riyid5dece",
            "highTime":1652934553133,
            "iexAskPrice":null,
            "iexAskSize":null,
            "iexBidPrice":null,
            "iexBidSize":null,
            "iexClose":133.12,
            "iexCloseTime":1661074865704,
            "iexLastUpdated":null,
            "iexMarketPercent":null,
            "iexOpen":136.57,
            "iexOpenTime":1629195416523,
            "iexRealtimePrice":null,
            "iexRealtimeSize":null,
            "iexVolume":null,
            "isUSMarketOpen":false,
            "lastTradeTime":1666364457060,
            "latestPrice":137.93,
            "latestSource":"Close",
            "latestTime":"January 29, 2021",
            "latestUpdate":1617261572432,
            "latestVolume":178713003,
            "low":136.04,
            "lowSource":"naep1iut m li5eecd yder",
            "lowTime":1633883662832,
            "marketCap":2272704860911,
            "oddLotDelayedPrice":134.3,
            "oddLotDelayedPriceTime":1665274496269,
            "open":141.3,
            "openSource":"aifoiclf",
            "openTime":1649340563111,
            "peRatio":37.37,
            "previousClose":138.65,
            "previousVolume":143335271,
            "primaryExchange":"NKL)ETBOA/SR ETQGMLACES AN(ALD SG",
            "symbol":"AAPL",
            "volume":184630595,
            "week52High":147.68,
            "week52Low":57.36,
            "ytdChange":-0.04294152218978952
        }"#;
        let _val: StockQuoteResponse = serde_json::from_str(raw_data)?;

        Ok(())
    }

    #[test]
    #[ignore]
    fn test_request() -> Result<(), &'static str> {
        let token = std::env::var("EXCHANGE_IEX_TOKEN")
            .or(Err("Could not find token in environment variable"))?;
        let inner_client = reqwest::Client::new();
        let client = Client::new(&inner_client, token);

        let request = StockQuoteRequestBuilder::new("aapl");

        let _resp = tokio_test::block_on(client.make_request(request))?;

        Ok(())
    }
}
