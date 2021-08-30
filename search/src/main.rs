extern crate bincode;
extern crate rust_stemmers;

use std::collections::HashMap;
use std::io::{self, BufRead, BufReader};
use std::fs::File;
use regex::Regex;
use chrono::{NaiveDateTime};

struct Context {

    ID: usize,
    word_to_id: HashMap<&'static str, usize>,
    re_word: Regex,
    re_review: Regex,
}
    
struct Review<'a> {

    id:   u64,
    date: Option<NaiveDateTime>,
    text: &'a str,
    word_ids: Vec<usize>,
}

impl<'a> Review<'a> {

    pub fn new(review: (&'a str, &'a str)) -> Box<Review> {

        Box::new(Review{
            id:   0,
            date: NaiveDateTime::parse_from_str(&review.1, "%Y-%m-%d %H:%M:%S").ok(),
            text: review.0,
            word_ids: Vec::new(),
        })
    }
}
 
fn get_word_id(c: &mut Context, word: &str) -> usize {

    if c.word_to_id.contains_key(word) {
        *c.word_to_id.get(word).unwrap()
    } else {
        c.ID += 1;
        c.word_to_id.insert(&word, c.ID);
        c.ID
    }
}

pub fn update_word_ids<'b>(c: &mut Context, review: &'b mut Review<'b>) -> &'b mut Review<'b> {

    use rust_stemmers::{Algorithm, Stemmer};
    let stemmer = Stemmer::create(Algorithm::English);

    let mut text = review.text.clone()
                .to_string()
                .to_lowercase();
    text.retain( |c| c != '\'' );
    
    let captures = c.re_word.captures(&text);
    let mut words: Vec<String> = vec![];
    for (i, capture) in captures.iter().enumerate() {

        if i == 0 {   // ignore full match
            continue;
        }

        words.push(capture.get(i).map_or("".to_string(), |m| m.as_str().to_string()))
    }
        
    review.word_ids = words.iter()
                    .map(|s| stemmer.stem(s))
                    .map(|w| get_word_id(&mut c, &w))
                    .collect::<Vec<usize>>();
    
    review
}

fn extract_data<R>(c: &mut Context, num_reviews: usize, reader: &mut R) -> io::Result<()> 
        where R: BufRead + std::fmt::Debug
{
    let reviews = reader.lines().take(num_reviews).collect::<io::Result<Vec<String>>>()?;
    let reviews = reviews.iter()
                    .map( |r| c.re_review.captures(r).unwrap() )
                    .map( |captures| (captures.get(1).map_or("", |m| m.as_str()),
                                 captures.get(2).map_or("", |m| m.as_str())) )
                    .map( |r| Review::new(r) )
                    .map( |r| update_word_ids(c, &mut r) );

    for r in reviews {
        println!("Id: {}\nDate: {:?}\nReview: {}\n---------------------------------\n",
            r.id, r.date, r.text);
    }
 
    Ok(())   
}

fn main() -> std::io::Result<()> {

    let f_in = File::open("/Users/samirgadkari/work/datasets/yelp/yelp_academic_dataset_review.json")?;
    let mut reader = BufReader::new(f_in);

    let mut context = Context {
        ID:         0,
        word_to_id: HashMap::new(),
        re_word:    Regex::new(r"(\w+)").unwrap(),
        re_review:  
            Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap(),
    };
    let _result = extract_data(&mut context, 20, &mut reader);

    Ok(())
}
