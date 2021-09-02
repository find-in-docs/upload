extern crate bincode;
extern crate rust_stemmers;

// use serde::Deserialize;
// use serde_json::{Result, Value};
use std::collections::HashMap;
use std::io::{self, BufRead, BufReader};
use std::fs::{self, File};
use regex::Regex;
use chrono::{NaiveDateTime};

struct Context {

    id: usize,
    word_to_id: HashMap<String, usize>,
    re_word: Regex,
    stopwords: Vec<String>,
    re_review: Regex,
}
    
#[derive(Debug)]
struct Review {

    id:   u64,
    date: Option<NaiveDateTime>,
    text: String,
    word_ids: Vec<usize>,
}


impl Review {

    pub fn new(review: (String, String)) -> Box<Review> {

        let review = Review{
            id:   0,
            date: NaiveDateTime::parse_from_str(&review.1, "%Y-%m-%d %H:%M:%S").ok(),
            text: review.0.clone(),
            word_ids: vec![],
        };

        let b: Box<Review> = Box::new(review);

        b
    }
}
 
fn get_word_id(c: &mut Context, word: String) -> usize {

    if c.word_to_id.contains_key(&word) {
        *c.word_to_id.get(&word).unwrap()
    } else {
        c.id += 1;
        c.word_to_id.insert(word, c.id);
        c.id
    }
}

fn update_word_ids(context: &mut Context, mut review: Box<Review>) 
    -> Box<Review> {

    use rust_stemmers::{Algorithm, Stemmer};
    let stemmer = Stemmer::create(Algorithm::English);

    let mut text = review.text.clone()
                .to_string()
                .to_lowercase();
    text.retain( |c| c != '\'' );

    let mut words: Vec<String> = vec![];
    for caps in context.re_word.captures_iter(&text) {
        words.push(caps[0].to_string());
    };
    
    review.word_ids = words.iter()
                    .map(|s| stemmer.stem(s))
                    .map(|w| get_word_id(context, (&w).to_string()))
                    .collect::<Vec<usize>>();
    
    review
}

fn stopwords() -> Option<Vec<String>> {
    let data = fs::read_to_string(
        "/Users/samirgadkari/work/datasets/yelp/english_stopwords.json")
        .expect("Unable to read file");

    let stopwords: serde_json::Value = 
        serde_json::from_str(&data)
        .expect("Unable to parse");

    Some(stopwords["english_stopwords"].as_array()?.iter()
        .map(|v| v.as_str().unwrap().to_string())
        .collect::<Vec<String>>())
}

fn extract_data(num_reviews: usize) -> io::Result<()> 
{
    let f_in = File::open("/Users/samirgadkari/work/datasets/yelp/yelp_academic_dataset_review.json")?;
    let reader = BufReader::new(f_in);
    
    let mut context = Context {
        id:         0,
        word_to_id: HashMap::new(),
        re_word:    Regex::new(r#"(\w+)"#).unwrap(),
        stopwords:  (&stopwords()).as_ref().unwrap().to_vec(),
        re_review:  
            Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap(),
    };

    let reviews = reader.lines().take(num_reviews).collect::<io::Result<Vec<String>>>()?;
    let reviews: Vec<Box<Review>> = reviews.iter()
        .map( |r| context.re_review.captures(r).unwrap() )
        .map( |captures| (captures.get(1).map_or("", |m| &m.as_str()),
                     captures.get(2).map_or("", |m| &m.as_str())) )
        .map( |r| (r.0.to_string(), r.1.to_string()) )
        .map( |r| Review::new(r) )
        .collect::<Vec<Box<Review>>>();

    let reviews: Vec<Box<Review>> =
            reviews.into_iter()
             .map( |r| update_word_ids(&mut context, r) )
             .collect::<Vec<Box<Review>>>();

    /*
    for r in reviews {
        println!("Id: {}\nDate: {:?}\nReview: {}\nWord IDs: {:?}\nWord ID len: {}\n---------------------------------\n",
            r.id, r.date, r.text, r.word_ids, r.word_ids.len());
    }
    */

    // println!("{:?}", context.word_to_id);

    Ok(())   
}

fn main() -> std::io::Result<()> {

    let _result = extract_data(20);

    Ok(())
}
