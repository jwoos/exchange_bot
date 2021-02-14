const API_BASE: &'static str = "https://www.alphavantage.co";

// TODO make this a trait and implement it per client
struct Client<'a> {
    client: &'a reqwest::Client,
    api_key: String,
}

impl <'a> Client<'a> {
    pub fn new(client: &'a reqwest::Client, token: String) -> Client<'a> {
        Client {
            client,
            token,
        }
    }

    pub fn get_base(&self) -> &'static str {
        //API_BASE
        API_SANDBOX_BASE
    }

    // TODO actually use errors properly
    pub async fn make_request<Req>(&self, builder: Req) -> Result<Req::Response, &'static str>
    where
        Req: RequestBuilder,
        Req::Response: for<'de> serde::Deserialize<'de>,
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

pub trait RequestBuilder {
    type Response;

    fn build(&self, base: &str) -> String;
}
