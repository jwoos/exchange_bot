use std::collections::VecDeque;
use std::vec;

pub struct AssetMetadata {
    // date as YYYY-MM-DD
    date_purchased: String,
    price: i64,
}

pub enum AssetType {
    Stock,
    Cryptocurrency,
}

pub trait Asset {
    fn get_ticker(&self) -> &str;

    fn get_price_basis(&self) -> i64;

    fn get_prices(&self) -> vec::Vec<i64>;

    fn get_type(&self) -> AssetType;
}

// TODO implement Asset
pub struct Stock {
    ticker: String,
    metadata: VecDeque<AssetMetadata>,
    count: u64,
}
