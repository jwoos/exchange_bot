use std::boxed::Box;
use std::vec::Vec;

use super::event;

#[derive(serde::Deserialize)]
pub struct Authorization {
    enterprise_id: String,
    team_id: String,
    user_id: String,
    is_bot: bool,
}

#[derive(serde::Deserialize)]
pub struct EventWrapper {
    token: String,
    team_id: String,
    api_app_id: String,
    #[serde(deserialize_with = "event_deserialize")]
    event: Box<dyn event::Event>,
    #[serde(alias = "type")]
    type_: String,
    authed_users: Vec<String>,
    authed_teams: Vec<String>,
    authorizations: Option<Authorization>,
    event_context: Option<String>,
    event_id: Option<String>,
    event_time: Option<u64>,
}

/* deals with deserialization of the inner event, whose structure depends on
 * what type it is
 */
fn event_deserialize<'a, 'de, D>(deserializer: D) -> Result<Box<dyn event::Event>, D::Error>
where
    D: serde::de::Deserializer<'de>,
{
    struct EventVisitor;

    impl<'de> serde::de::Visitor<'de> for EventVisitor {
        type Value = Box<dyn event::Event>;

        fn expecting(&self, formatter: &mut std::fmt::Formatter) -> std::fmt::Result {
            formatter.write_str("a string containing event data")
        }

        fn visit_str<E>(self, v: &str) -> Result<Self::Value, E>
        where
            E: serde::de::Error,
        {
            let value: serde_json::Value =
                serde_json::from_str(v).map_err(serde::de::Error::custom)?;

            let map = value
                .as_object()
                .ok_or(serde::de::Error::custom("Value should have been an object"))?;
            let event_type = map
                .get("type")
                .ok_or(serde::de::Error::custom("Map doesn't have field `type`"))?
                .as_str()
                .ok_or(serde::de::Error::custom("Event type was not a string"))?;

            return match event_type {
                "app_mention" => {
                    let obj = serde_json::from_str::<event::app_mention::AppMention>(v)
                        .map_err(serde::de::Error::custom)?;
                    Ok(Box::new(obj))
                }
                _ => Err(serde::de::Error::custom("Unsupported event type")),
            };
        }
    }

    deserializer.deserialize_any(EventVisitor {})
}
