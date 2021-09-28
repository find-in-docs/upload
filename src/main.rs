mod doc;
mod errors;
mod utils;

use crate::errors::DocError;
use doc::extract_docs;
use std::result::Result;

fn main() -> Result<(), DocError> {
    let config = utils::load_config()?;
    let result = extract_docs(
        &config.in_docs_filename,
        &config.in_stopwords,
        &config.out_dirname,
        5,
    );

    println!("{:?}", result);
    Ok(())
}
