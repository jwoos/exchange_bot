pub mod chat;

pub struct SlackClient {
    token: String,
}

impl SlackClient {
    pub fn new(token: String) -> SlackClient {
        SlackClient{token}
    }
}
