use super::{Event, EventType};

#[derive(Debug, serde::Deserialize)]
pub struct UrlVerification {
    token: String,
    challenge: String,
    #[serde(alias = "type")]
    type_: String,
}

impl Event for UrlVerification {
    fn as_any(&self) -> &dyn std::any::Any {
        self
    }

    fn get_type(&self) -> EventType {
        EventType::UrlVerification
    }

    fn get_event_ts<'a>(&'a self) -> &'a str {
        &""
    }
}

impl UrlVerification {
    pub fn get_token<'a>(&'a self) -> &'a str {
        &self.token
    }

    pub fn get_challenge<'a>(&'a self) -> &'a str {
        &self.challenge
    }
}
