use super::super::{RequestBuilder, Response, Type};
use super::StockType;

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
    fn build(&self, base: &str) -> String {
        std::format!("{}/stock/{}/quote", base, self.ticker)
    }
}

#[serde(rename_all = "snake_case")]
#[derive(serde::Deserialize)]
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
    pub is_us_market_open: bool,
}
