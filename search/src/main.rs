extern crate bincode;

use std::io::{self, BufRead, BufReader};
use std::fs::File;
use regex::Regex;


fn extract_data<R>(num_reviews: usize, reader: &mut R) -> io::Result<()> 
        where R: BufRead + std::fmt::Debug
{
    let re_review = Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap();

    let reviews = reader.lines().take(num_reviews).collect::<io::Result<Vec<String>>>()?;
    let reviews = reviews.iter()
                    .map( |r| re_review.captures(r).unwrap() )
                    .map( |c| (c.get(1).map_or("", |m| m.as_str()),
                                 c.get(2).map_or("", |m| m.as_str())));

    for r in reviews {
        println!("Date: {}\nReview: {}", r.1, r.0);
    }
 
    Ok(())   
}

fn main() -> std::io::Result<()> {

    let f_in = File::open("/Users/samirgadkari/work/datasets/yelp/yelp_academic_dataset_review.json")?;
    let mut reader = BufReader::new(f_in);

    let _result = extract_data(20, &mut reader);

    Ok(())
}
