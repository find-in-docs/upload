mod doc;
mod errors;
mod utils;

use crate::errors::DocError;
use doc::Processing;
use std::fs::File;
use std::io::BufReader;
use std::result::Result;

fn main() -> Result<(), DocError> {
    let config = utils::load_config()?;
    let mut data = doc::Data::new();

    let f_in = File::open(config.in_docs_filename)
        .expect(&format!("Could not open file: {}", config.in_docs_filename));
    let mut docs_reader = BufReader::new(f_in);

    data.extract_docs(&mut docs_reader, 5, &config.in_stopwords);

    // println!("{:?}", result);
    Ok(())
}
