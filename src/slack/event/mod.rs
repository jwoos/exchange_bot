pub mod app_mention;

pub trait Event {
    fn get_type(&self) -> EventType;

    fn get_event_ts<'a>(&'a self) -> &'a str;
}

#[derive(serde::Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum EventType {
    AppMention,
    //Message(&'a str),
    UrlVerification,
    //Member(&'a str),
}
