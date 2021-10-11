extern crate serde;
extern crate serde_yaml;

use crate::errors::DocError;
use ::serde::Deserialize;

#[derive(Deserialize, Debug)]
pub struct Config {
    pub in_docs_filename: String,
    pub out_dirname: String,
    pub in_stopwords: String,
}

pub fn load_config() -> std::result::Result<Config, DocError> {
    let yaml = include_str!("config.yaml");
    println!("{}", yaml);
    let config = serde_yaml::from_str(&yaml)?;
    println!("{:?}", config);
    Ok(config)
}

/*impl<T> std::iter::Iterator for Option {
    type Item = Option;

    fn next(&mut self) -> Self::Item {
        match self {
            None => None,
            Some(x) => x,
        }
    }
}*/
