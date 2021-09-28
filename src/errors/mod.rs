extern crate thiserror;

use std::io;
use thiserror::Error;

pub struct ParserError {
    message: String,
    line: usize,
}

#[derive(Error, Debug)]
pub enum DocError {
    #[error("Config error")]
    ConfigError(#[from] serde_yaml::Error),

    #[error("{message:} {line:}")]
    ParserError { message: String, line: usize },

    #[error("IO error")]
    IOError(#[from] io::Error),
}
