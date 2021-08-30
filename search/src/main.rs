extern crate lazy_static;
extern crate bincode;
extern crate rust_stemmers;

use lazy_static::lazy_static;
use std::collections::HashMap;
use std::io::{self, BufRead, BufReader};
use std::fs::File;
use regex::Regex;
use chrono::{NaiveDateTime};

struct Review<'a> {

    id:   u64,
    date: Option<NaiveDateTime>,
    text: &'a str,
    word_ids: Vec<usize>,
}

impl<'a> Review<'a> {

    pub fn new(review: (&'a str, &'a str)) -> &Review {

        &Review{
            id:   0,
            date: NaiveDateTime::parse_from_str(&review.1, "%Y-%m-%d %H:%M:%S").ok(),
            text: review.0,
            word_ids: Vec::new(),
        }
    }
}
 
pub fn update_word_ids<'b>(review: &'b mut Review<'b>) -> &'b mut Review<'b> {

    lazy_static! {
        static ref re_word: Regex = 
            Regex::new(r"(\w+)").unwrap();
    }

    use rust_stemmers::{Algorithm, Stemmer};
    let stemmer = Stemmer::create(Algorithm::English);

    fn get_word_id(word: &str) -> usize {

        let mut ID: &'static usize = &0;

        lazy_static! {
            static ref word_to_id: HashMap<&'static str, usize> = HashMap::new();
        }

        if word_to_id.contains_key(word) {
            *word_to_id.get(word).unwrap()
        } else {
            *ID += 1;
            word_to_id.insert(word.clone(), *ID);
            *ID
        }
    }

    let text = review.text.clone()
                .to_string()
                .to_lowercase()
                .retain( |c| c != '\'' );
    /*
    review.word_ids = re_word.captures_iter(text)
                    .map(|s| stemmer.stem(s))
                    .map(get_word_id)
                    .collect::<Vec<usize>>();
    */
    review
}

fn extract_data<R>(num_reviews: usize, reader: &mut R) -> io::Result<()> 
        where R: BufRead + std::fmt::Debug
{
    lazy_static! {
        static ref re_review: Regex = 
            Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap();
    }

    let reviews = reader.lines().take(num_reviews).collect::<io::Result<Vec<String>>>()?;
    let reviews = reviews.iter()
                    .map( |r| re_review.captures(r).unwrap() )
                    .map( |c| (c.get(1).map_or("", |m| m.as_str()),
                                 c.get(2).map_or("", |m| m.as_str())) )
                    .map( |r| Review::new(r) )
                    .map( |&r| update_word_ids(&mut r) );

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
