use std::collections::HashMap;

use crate::asset::Stock;

pub struct Portfolio {
    stock: HashMap<String, Stock>,
}

impl Portfolio {
    pub fn new() -> Portfolio {
        Portfolio {
            stock: HashMap::new(),
        }
    }
}
