use super::{Event, EventType};

#[derive(serde::Deserialize)]
pub struct AppMention {
    #[serde(alias = "type")]
    type_: EventType,
    user: String,
    text: String,
    ts: String,
    channel: String,
    event_ts: String,
}

impl Event for AppMention {
    fn get_type(&self) -> EventType {
        EventType::AppMention
    }

    fn get_event_ts<'a>(&'a self) -> &'a str {
        &self.event_ts
    }
}
