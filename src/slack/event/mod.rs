pub mod app_mention;

pub trait Event {
    // for downcasting to concrete type from trait object
    fn as_any(&self) -> &dyn std::any::Any;

    fn get_type(&self) -> EventType;

    fn get_event_ts<'a>(&'a self) -> &'a str;
}

#[derive(Debug, PartialEq, serde::Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum EventType {
    AppMention,
    //Message(&'a str),
    UrlVerification,
    //Member(&'a str),
}
