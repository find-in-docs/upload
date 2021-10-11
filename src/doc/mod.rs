extern crate bincode;
extern crate rust_stemmers;
extern crate thiserror;

#[allow(unused_imports)]
pub use crate::errors::{DocError, ParserError};

mod doc_id;

use chrono::NaiveDateTime;
pub use doc_id::get_doc_id;
use regex::Regex;
// use rust_stemmers::{Algorithm, Stemmer};
// use rust_stemmers::Stemmer;
// use std::collections::{HashMap, HashSet};
use std::collections::HashSet;
use std::fs::File;

#[allow(unused_imports)]
use std::io::{BufRead, BufReader, ErrorKind, Lines};

// use std::iter::FromIterator;

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

pub struct Data<'a> {
    lines: Option<std::iter::Take<Lines<&'a mut BufReader<File>>>>,
    docs: Option<Vec<Doc<'a>>>,
    stopwords: Option<HashSet<String>>,
}

pub trait Processing {
    /*fn update_option<'a, T, Q>(
        opt: &'a Option<Q>,
        updater: &dyn FnMut(Q) -> Option<T>,
    ) -> Option<T>;*/
    fn string_to_doc<'a>(re_doc: &Regex, s: &'a str) -> Option<Doc<'a>>;
    fn extract_docs(&mut self, in_docs: &mut BufReader<File>, num_docs: usize, in_stopwords: &str);
    fn load_stopwords(&mut self, stopwords_fn: &str);
}

impl Data<'_> {
    pub fn new() -> Box<Data<'static>> {
        Box::new(Data {
            lines: None,
            docs: None,
            stopwords: None,
        })
    }
}

impl Processing for Data<'_> {
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

    fn extract_docs(
        &mut self,
        docs_reader: &mut BufReader<File>,
        num_docs: usize,
        in_stopwords: &str,
    ) {
        let re_doc: Regex =
            Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap();

        self.lines = Some(docs_reader.lines().take(num_docs));
        if let Some(lines) = self.lines {
            self.docs = lines.map(|x| {
                Some(
                    lines
                        .filter_map(Result::ok)
                        .filter_map(|s: &str| self::Processing::string_to_doc(&re_doc, s))
                        .collect::<Vec<Doc>>(),
                )
            });
        } else {
            self.docs = None;
        }
        /*self.docs =
        Processing::update_option::<Vec<Doc>, std::iter::Take<Lines<&mut BufReader<File>>>>(
            &self.lines,
            &(|lines: &std::iter::Take<Lines<&mut BufReader<File>>>| {
                Some(
                    (*lines)
                        .filter_map(Result::ok)
                        .filter_map(|s| self::Processing::string_to_doc(&re_doc, &s))
                        .collect::<Vec<Doc>>(),
                )
            }),
        );*/
        self.load_stopwords(in_stopwords);
    }

    /*fn update_option<'a, T, Q>(opt: &'a Q, updater: &dyn FnMut(Q) -> Option<T>) -> Option<T> {
        match opt {
            None => &None,
            Some(v) => updater(*v),
        }
    }*/

    fn load_stopwords(&mut self, in_stopwords: &str) {
        let data = std::fs::read_to_string(in_stopwords).expect("Unable to read file");
        let stopwords: serde_json::Value = serde_json::from_str(&data).expect("Unable to parse");

        self.stopwords = stopwords["english_stopwords"]
            .as_array()
            .map(|x| {
                Some(
                    x.iter()
                        .map(|v| v.as_str().unwrap().to_string())
                        .collect::<HashSet<String>>(),
                )
            })
            .unwrap_or(None);
        /*self.stopwords = Processing::update_option::<HashSet<String>, Vec<serde_json::value::Value>>(
            &stopwords["english_stopwords"].as_array(),
            &(|stopwords: Vec<serde_json::value::Value>| {
                Some(
                    stopwords
                        .iter()
                        .map(|v| v.as_str().unwrap().to_string())
                        .collect::<HashSet<String>>(),
                )
            }),
        );*/
    }

    /*
    fn extract_docs(
        &mut self,
        in_filename: &str,
        stopwords_fn: &str,
        _out_dirname: &str,
        num_docs: usize,
    ) {
        let f_in = File::open(in_filename).expect(&format!("Could not open file: {}", in_filename));

        let reader = BufReader::new(f_in);
        let re_doc: Regex =
            Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap();

        self.lines = reader.lines().take(num_docs);
        self.docs = self
            .lines
            .filter_map(Result::ok)
            .filter_map(|s| self.string_to_doc(&re_doc, &s))
            .collect::<Vec<Doc>>();

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
    */
}
