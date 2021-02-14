use super::super::RequestBuilder;

use std::collections::HashMap;

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
        std::format!("{}/query/{}", base, self.ticker)
    }
}

#[derive(serde::Deserialize)]
pub struct InnerStockQuoteResponse {
    #[serde(alias = "01. symbol")]
    pub symbol: String,

    #[serde(alias = "02. open", deserialize_with = "string_to_f64_deserialize")]
    pub open: f64,

    #[serde(alias = "03. high", deserialize_with = "string_to_i64_deserialize")]
    pub high: f64,

    #[serde(alias = "04. low", deserialize_with = "string_to_f64_deserialize")]
    pub low: f64,

    #[serde(alias = "05. price", deserialize_with = "string_to_f64_deserialize")]
    pub price: f64,

    #[serde(alias = "06. volume", deserialize_with = "string_to_i64_deserialize")]
    pub volume: u64,

    #[serde(alias = "07. latest trading day")]
    pub latest_trading_day: String,

    #[serde(alias = "08. previous close", deserialize_with = "string_to_f64_deserialize")]
    pub previous_close: f64,

    #[serde(alias = "09. change", deserialize_with = "string_to_f64_deserialize")]
    pub change: f64,

    #[serde(alias = "10. change percent")]
    pub change_percent: String,
}

#[derive(serde::Deserialize)]
pub struct StockQuoteResponse {
    #[serde(alias = "Global Quote")]
    pub global_quote: InnerStockQuoteResponse,
}

fn string_to_f64_deserialize<'a, 'de, D>(
    deserializer: D,
) -> Result<f64, D::Error>
where D: serde::de::Deserialize<'de>,
      {
          struct StringVisitor;

          impl<'de> serde::de::Visitor<'de> for StringVisitor {
              type Value = f64;

              fn expecting(&self, formatter: &mut std::fmt::Formatter) -> std::fmt::Result {
                  formatter.write_str("a map")
              }

              fn visit_str<E>(self, s: &str) -> Result<Self::Value, E>
                  where
                      E: serde::de::Error,
                  {
                      s.parse::<f64>().ok_or(serde::de::Error::custom(format!("Failed to parse {} into f64", s)))
                  }
          }

          deserializer.deserialize_any(StringVisitor{})
      }

fn string_to_i64_deserialize<'a, 'de, D>(
    deserializer: D,
) -> Result<f64, D::Error>
where D: serde::de::Deserialize<'de>,
      {
          struct StringVisitor;

          impl<'de> serde::de::Visitor<'de> for StringVisitor {
              type Value = i64;

              fn expecting(&self, formatter: &mut std::fmt::Formatter) -> std::fmt::Result {
                  formatter.write_str("a map")
              }

              fn visit_str<E>(self, s: &str) -> Result<Self::Value, E>
                  where
                      E: serde::de::Error,
                  {
                      s.parse::<i64>().ok_or(serde::de::Error::custom(format!("Failed to parse {} into i64", s)))
                  }
          }

          deserializer.deserialize_any(StringVisitor{})
      }


mod tests {
    use super::super::super::Client;
    use super::*;

    #[test]
    fn test_deserialize_resp() -> Result<(), serde_json::Error> {
        let raw_data = r#"{
            "Global Quote": {
                "01. symbol": "IBM",
                "02. open": "121.0000",
                "03. high": "121.3600",
                "04. low": "120.0900",
                "05. price": "120.8000",
                "06. volume": "3871195",
                "07. latest trading day": "2021-02-12",
                "08. previous close": "120.9100",
                "09. change": "-0.1100",
                "10. change percent": "-0.0910%"
            }
        }"#;

        let val: StockQuoteResponse = serde_json::from_str(raw_data);

        Ok(())
    }
}
