extern crate bincode;
extern crate rust_stemmers;
extern crate thiserror;

#[allow(unused_imports)]
pub use crate::errors::{DocError, ParserError};

mod doc_id;

use chrono::NaiveDateTime;
pub use doc_id::get_doc_id;
use regex::Regex;
use rust_stemmers::{Algorithm, Stemmer};
use std::collections::{HashMap, HashSet};
use std::fs::{self, File};

#[allow(unused_imports)]
use std::io::{BufRead, BufReader, ErrorKind};

use std::iter::FromIterator;

pub type DocId = usize;
pub type WordId = usize;
pub type DocResult<'a> = Result<Vec<Doc<'a>>, DocError>;

#[derive(Debug)]
pub struct Doc<'a> {
    id: DocId,
    date: Option<NaiveDateTime>,
    text: &'a str,
}

impl Doc<'_> {
    pub fn new<'a>(doc: (&'a str, &str)) -> Box<Doc<'a>> {
        Box::new(Doc {
            id: get_doc_id(),
            date: NaiveDateTime::parse_from_str(&doc.1, "%Y-%m-%d %H:%M:%S").ok(),
            text: doc.0.clone(),
        })
    }
}

struct Context {
    word_to_id: HashMap<String, WordId>,
    stopwords: HashSet<String>,
    stemmer: Stemmer,
    re_doc: Regex,
    re_word: Regex,
}

fn stopwords(in_stopwords: &str) -> Option<Vec<String>> {
    let data = fs::read_to_string(in_stopwords).expect("Unable to read file");

    let stopwords: serde_json::Value = serde_json::from_str(&data).expect("Unable to parse");

    Some(
        stopwords["english_stopwords"]
            .as_array()?
            .iter()
            .map(|v| v.as_str().unwrap().to_string())
            .collect::<Vec<String>>(),
    )
}

fn string_to_doc<'a>(re_doc: &Regex, s: &'a str) -> Option<Doc<'a>> {
    let doc = match re_doc.captures(&s) {
        None => None,
        Some(c) => {
            let text = c.get(1);
            let date = c.get(2);
            if (text == None) || (date == None) {
                return None;
            }
            let text = text.map_or("", |m| m.as_str());
            let date = date.map_or("", |m| m.as_str());
            let document = Doc::new((text, date));
            println!("{:?}", document);
            Some(*document)
            // Some(Doc::new((text, date)))
        }
    };

    doc
}

pub fn extract_docs<'a>(
    in_filename: &'a str,
    in_stopwords: &str,
    _out_dirname: &str,
    num_docs: usize,
) -> DocResult<'a> {
    let f_in = File::open(in_filename)?;
    let reader = BufReader::new(f_in);

    let context = Context {
        id: 0,

        word_to_id: HashMap::new(),
        stopwords: HashSet::from_iter((&stopwords(in_stopwords)).as_ref().unwrap().to_vec()),
        stemmer: Stemmer::create(Algorithm::English),

        re_doc: Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap(),
        re_word: Regex::new(r#"(\w+)"#).unwrap(),
    };

    Ok(reader
        .lines() // each line is a document
        .take(num_docs)
        .filter_map(Result::ok)
        .filter_map(|s| string_to_doc(&context.re_doc, &s))
        .collect::<Vec<Doc>>())

    /*
    // let _docs = strings.filter_map(|s| string_to_doc(&mut context, &s));
    let mut docs: Vec<Doc> = vec![];
    for s in strings {
        match string_to_doc(&context.re_doc, &s) {
            None => {
                println!("None");
                ()
            }
            Some(d) => {
                println!("Some: d = {:?}", d);
                docs.push(d)
            }
        }
    }*/

    // Ok(docs)
}
