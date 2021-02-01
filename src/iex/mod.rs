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

    // TODO actually use errors properly
    pub async fn make_request<Req, Res>(&self, builder: Req) -> Result<Res, &'static str>
    where
        Req: RequestBuilder,
        Res: for<'de> serde::Deserialize<'de>,
    {
        let url = builder.build(self.get_base());
        let resp = self
            .client
            .get(&url)
            .query(&[("token", &self.token)])
            .send()
            .await
            .or(Err("Error making request"))?;
        let json = resp.json().await.or(Err("Error getting json"))?;
        serde_json::from_value(json).or(Err("Unable to deserialize result"))
    }
}

pub enum Type {
    Dummy,
    Stock(stock::StockType),
}

pub trait RequestBuilder {
    fn build(&self, base: &str) -> String;
}
