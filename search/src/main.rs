extern crate bincode;
extern crate rust_stemmers;

use std::collections::HashMap;
use std::io::{self, BufRead, BufReader};
use std::fs::File;
use regex::Regex;
use chrono::{NaiveDateTime};

struct Context {

    id: usize,
    word_to_id: HashMap<String, usize>,
    re_word: Regex,
    re_review: Regex,
}
    
#[derive(Debug)]
struct Review {

    id:   u64,
    date: Option<NaiveDateTime>,
    text: String,
    word_ids: Option<Vec<usize>>,
}

impl Review {

    pub fn new(review: (String, String)) -> Box<Review> {

        let review = Review{
            id:   0,
            date: NaiveDateTime::parse_from_str(&review.1, "%Y-%m-%d %H:%M:%S").ok(),
            text: review.0.clone(),
            word_ids: None,
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

fn update_word_ids(_c: &mut Context, review: Box<Review>) 
    -> Box<Review> {

    /*
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
        
    review.word_ids = Some(words.iter()
                    .map(|s| stemmer.stem(s))
                    .map(|w| get_word_id(c, (&w).to_string()))
                    .collect::<Vec<usize>>());
    */
    
    review
}

fn extract_data<R>(num_reviews: usize, reader: &mut R) -> io::Result<()> 
        where R: BufRead + std::fmt::Debug
{
    let mut context = Context {
        id:         0,
        word_to_id: HashMap::new(),
        re_word:    Regex::new(r"(\w+)").unwrap(),
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

    println!("{:?}", reviews);
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
