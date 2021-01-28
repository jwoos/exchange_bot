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
            formatter.write_str("a map")
        }

        fn visit_map<M>(self, mut v: M) -> Result<Self::Value, M::Error>
        where
            M: serde::de::MapAccess<'de>,
        {
            let mut map = serde_json::map::Map::new();

            while let Some((key, value)) = v.next_entry()? {
                map.insert(key, value);
            }

            let event_type = map
                .get("type")
                .ok_or(serde::de::Error::custom("Map doesn't have field `type`"))?
                .as_str()
                .ok_or(serde::de::Error::custom("Event type was not a string"))?;

            return match event_type {
                "app_mention" => {
                    let obj = serde_json::from_value::<event::app_mention::AppMention>(
                        serde_json::value::Value::Object(map),
                    )
                    .map_err(serde::de::Error::custom)?;
                    Ok(Box::new(obj))
                }
                _ => Err(serde::de::Error::custom("Unsupported event type")),
            };
        }
    }

    deserializer.deserialize_any(EventVisitor {})
}

#[cfg(test)]
mod tests {
    use super::event::app_mention::AppMention;
    use super::event::{Event, EventType};
    use super::EventWrapper;

    #[test]
    fn test_app_mention() {
        let test_str = r#"{
                "token": "test_token",
                "team_id": "test_team_id",
                "api_app_id": "test_api_app_id",
                "type": "test_type",
                "authed_users": ["user_1", "user2"],
                "authed_teams": ["team_1", "team2"],
                "event": {
                    "type": "app_mention",
                    "user": "test_user",
                    "text": "text",
                    "ts": "100000000.00000",
                    "channel": "channel",
                    "event_ts": "10000000000000"
                }
            }"#;

        let event_wrapper: EventWrapper = serde_json::from_str(test_str).unwrap();
        let event: &AppMention = event_wrapper
            .event
            .as_any()
            .downcast_ref::<AppMention>()
            .unwrap();

        assert_eq!(event.get_type(), EventType::AppMention);
        assert_eq!(event.get_event_ts(), "10000000000000");
        assert_eq!(event.get_user(), "test_user");
        assert_eq!(event.get_text(), "text");
        assert_eq!(event.get_channel(), "channel");
    }
}
