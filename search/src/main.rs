extern crate bincode;
extern crate rust_stemmers;

use chrono::NaiveDateTime;
use regex::{Captures, Regex};
use rust_stemmers::{Algorithm, Stemmer};
use std::collections::{HashMap, HashSet};
use std::fs::{self, File};
use std::io::{self, BufRead, BufReader};
use std::iter::FromIterator;

type WordId = usize;
type ReviewId = usize;

struct Context {
    max_review_id: ReviewId,
    id: WordId,
    word_to_id: HashMap<String, WordId>,
    re_word: Regex,
    stopwords: HashSet<String>,
    stemmer: Stemmer,
    re_review: Regex,
}

#[derive(Debug)]
struct Review {
    id: ReviewId,
    date: Option<NaiveDateTime>,
    text: String,
    word_ids: Vec<WordId>,
}

impl Review {
    pub fn new(context: &mut Context, review: (String, String)) -> Box<Review> {
        let review = Review {
            id: context.max_review_id,
            date: NaiveDateTime::parse_from_str(&review.1, "%Y-%m-%d %H:%M:%S").ok(),
            text: review.0.clone(),
            word_ids: vec![],
        };

        let b: Box<Review> = Box::new(review);

        context.max_review_id += 1;
        b
    }
}

fn get_word_id(c: &mut Context, word: String) -> WordId {
    if c.word_to_id.contains_key(&word) {
        *c.word_to_id.get(&word).unwrap()
    } else {
        c.id += 1;
        c.word_to_id.insert(word, c.id);
        c.id
    }
}

fn update_word_ids(context: &mut Context, mut review: Box<Review>) -> Box<Review> {
    let mut text = review.text.clone().to_string().to_lowercase();
    text.retain(|c| c != '\'');

    let mut words: Vec<String> = vec![];
    for caps in context.re_word.captures_iter(&text) {
        words.push(caps[0].to_string());
    }

    let stemmed_words = words
        .into_iter()
        .filter(|w| !context.stopwords.contains(w))
        .map(|w| context.stemmer.stem(&w).to_string())
        .collect::<Vec<String>>();
    review.word_ids = stemmed_words
        .iter()
        .map(|w| get_word_id(context, (&w).to_string()))
        .collect::<Vec<WordId>>();

    review
}

fn stopwords() -> Option<Vec<String>> {
    let data = fs::read_to_string("/Users/samirgadkari/work/datasets/yelp/english_stopwords.json")
        .expect("Unable to read file");

    let stopwords: serde_json::Value = serde_json::from_str(&data).expect("Unable to parse");

    Some(
        stopwords["english_stopwords"]
            .as_array()?
            .iter()
            .map(|v| v.as_str().unwrap().to_string())
            .collect::<Vec<String>>(),
    )
}

fn extract_data(num_reviews: usize) -> io::Result<()> {
    let f_in =
        File::open("/Users/samirgadkari/work/datasets/yelp/yelp_academic_dataset_review.json")?;
    let reader = BufReader::new(f_in);

    let mut context = Context {
        max_review_id: 0,
        id: 0,
        word_to_id: HashMap::new(),
        re_word: Regex::new(r#"(\w+)"#).unwrap(),
        stopwords: HashSet::from_iter((&stopwords()).as_ref().unwrap().to_vec()),
        stemmer: Stemmer::create(Algorithm::English),
        re_review: Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap(),
    };

    let reviews = reader
        .lines()
        .take(num_reviews)
        .collect::<io::Result<Vec<String>>>()?;
    let captures = reviews
        .iter()
        .map(|r| context.re_review.captures(r).unwrap())
        .collect::<Vec<Captures>>();
    let reviews = captures
        .iter()
        .map(|captures| {
            (
                captures.get(1).map_or("", |m| &m.as_str()),
                captures.get(2).map_or("", |m| &m.as_str()),
            )
        })
        .map(|r| (r.0.to_string(), r.1.to_string()))
        .map(|r| Review::new(&mut context, r))
        .collect::<Vec<Box<Review>>>();

    let reviews: Vec<Box<Review>> = reviews
        .into_iter()
        .map(|r| update_word_ids(&mut context, r))
        .collect::<Vec<Box<Review>>>();

    for r in reviews {
        println!("Id: {}\nDate: {:?}\nReview: {}\nWord IDs: {:?}\nWord ID len: {}\n---------------------------------\n",
            r.id, r.date, r.text, r.word_ids, r.word_ids.len());
    }

    // println!("{:?}", context.word_to_id);

    Ok(())
}

fn main() -> std::io::Result<()> {
    let _result = extract_data(20);

    Ok(())
}
