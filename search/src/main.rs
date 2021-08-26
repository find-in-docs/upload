extern crate bincode;

use std::io::{self, BufRead, BufReader};
use std::fs::File;
use regex::Regex;
use chrono::{NaiveDateTime};

struct Review {

    id:   u64,
    date: Option<NaiveDateTime>,
    text: String
}
 
// Lowercase all letters.
// Remove apostrophes.
// Split review into words.
// Remove punctuation.
// Remove empty words.
// Remove stopwords.
fn build_review(review: (&str, &str)) -> Review {

    Review{
        id:   0,
        date: NaiveDateTime::parse_from_str(&review.1, "%Y-%m-%d %H:%M:%S").ok(),
        text: review.0.to_string(),
    }
}

fn extract_data<R>(num_reviews: usize, reader: &mut R) -> io::Result<()> 
        where R: BufRead + std::fmt::Debug
{
    let re_review = Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap();

    let reviews = reader.lines().take(num_reviews).collect::<io::Result<Vec<String>>>()?;
    let reviews = reviews.iter()
                    .map( |r| re_review.captures(r).unwrap() )
                    .map( |c| (c.get(1).map_or("", |m| m.as_str()),
                                 c.get(2).map_or("", |m| m.as_str())))
                    .map( |r| build_review(r) );

    for r in reviews {
        println!("Id: {}\nDate: {:?}\nReview: {}\n---------------------------------\n",
            r.id, r.date, r.text);
    }
 
    Ok(())   
}

fn main() -> std::io::Result<()> {

    let f_in = File::open("/Users/samirgadkari/work/datasets/yelp/yelp_academic_dataset_review.json")?;
    let mut reader = BufReader::new(f_in);

    let _result = extract_data(20, &mut reader);

    Ok(())
}
