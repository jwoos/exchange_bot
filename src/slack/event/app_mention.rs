use super::{Event, EventType};

#[derive(serde::Deserialize)]
pub struct AppMention {
    #[serde(alias = "type")]
    type_: EventType,
    event_ts: String,

    user: String,
    text: String,
    ts: String,
    channel: String,
}

impl Event for AppMention {
    fn as_any(&self) -> &dyn std::any::Any {
        self
    }

    fn get_type(&self) -> EventType {
        EventType::AppMention
    }

    fn get_event_ts<'a>(&'a self) -> &'a str {
        &self.event_ts
    }
}

impl AppMention {
    pub fn get_user<'a>(&'a self) -> &'a str {
        &self.user
    }

    pub fn get_text<'a>(&'a self) -> &'a str {
        &self.text
    }

    pub fn get_ts<'a>(&'a self) -> &'a str {
        &self.ts
    }

    pub fn get_channel<'a>(&'a self) -> &'a str {
        &self.channel
    }
}
