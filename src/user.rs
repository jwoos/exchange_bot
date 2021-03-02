use crate::portfolio;

pub struct User {
    id: String,
    // name: String,
    // money in CENTS to avoid weird floating point issues
    // also signed, for future margins/options
    money: i64,
    portfolio: portfolio::Portfolio,
}

impl User {
    fn new(id: String) -> User {
        User {
            id: id,
            money: 10000 * 100,
            portfolio: portfolio::Portfolio::new(),
        }
    }
}
