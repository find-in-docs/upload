extern crate bincode;

use std::io::{self, BufRead, BufReader};
use std::fs::File;
// use std::error::Error;
use regex::Regex;


fn extract_data<R>(num_reviews: usize, reader: R) -> io::Result<()> 
        where R: BufRead + std::fmt::Debug
{
    let reviews = reader.lines().take(num_reviews);
    println!("{:?}", reviews);

    let re_review = Regex::new(r#"^.*?"text"\s*:\s*"(.+)",\s*"date"\s*:\s*"(.*)".*$"#).unwrap();

    for review in reviews {

        // println!("{:?}", &review); 
        let review: String = review?;
        let captures = re_review.captures(&review).unwrap();
        let review = captures.get(1).map_or("", |m| m.as_str());
        let date = captures.get(2).map_or("", |m| m.as_str());
     
        if review.len() == 0 {
     
            println!("{}", date);
            println!("Review not found\n");
        } else if date.len() == 0 {
     
            println!("Date not found");
            println!("{}\n", review);
        } else {

            println!("Date: {}\nReview: {}", date, review);
        }
        println!("-------------------------------------------------");
    }
 
    Ok(())   
}

fn main() -> std::io::Result<()> {

    let f_in = File::open("/Users/samirgadkari/work/datasets/yelp/yelp_academic_dataset_review.json")?;
    let reader = BufReader::new(f_in);

    let _result = extract_data(20, reader);

    Ok(())
}
