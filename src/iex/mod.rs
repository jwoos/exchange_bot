pub mod stock;

const API_BASE: &'static str = "https://cloud.iexapis.com/v1";
const API_SANDBOX_BASE: &'static str = "https://sandbox.iexapis.com/v1";

struct Client<'a> {
    client: &'a reqwest::Client,
    token: String,
}

impl<'a> Client<'a> {
    pub fn new(client: &'a reqwest::Client, token: String) -> Client<'a> {
        Client {
            client: client,
            token: token,
        }
    }

    pub fn get_base(&self) -> &'static str {
        //API_BASE
        API_SANDBOX_BASE
    }

    pub async fn make_request<R: RequestBuilder>(&self, builder: R) -> impl Response {
        DummyResponse {}
    }
}

pub enum Type {
    Dummy,
    Stock(stock::StockType),
}

pub trait RequestBuilder {
    fn build(&self, base: &str) -> String;
}

pub trait Response {
    fn get_type(&self) -> Type;
}

// only for testing/building out the actual fn
pub struct DummyResponse {}

impl Response for DummyResponse {
    fn get_type(&self) -> Type {
        Type::Dummy
    }
}
