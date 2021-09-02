Project Architecture

1. Parsing reviews:
1a. Lowercase all characters
1b. Remove apostrophes (required? Maybe stemming will take care of this)
1c. Stem words using https://github.com/CurrySoftware/rust-stemmers
1d. Give an identifier to each word: word_id
1e. Give an identifier to each review: review_id
    review_id is a constantly-growing integer.
1f. Save this to binary file (pre_processed.reviews) for all reviews
    (review_id, review_date, review_len, [word_id]).
1g. Capture the location of each word in each review:
    (word, (review_id, review_date, location))
    As you get chunks of these words, write each chunk to N
    different files in a round-robin fashion. This way, the same word
    locations will not congregate in a single file. This helps when reading,
    since we can source the N files from different servers, thus
    decreasing latency.
    Call the files (review.locations.1, review.locations.2, etc.)
    We may find the same word has already been written to this file.
    That's ok. We will combine the words in the next step.
1h. Combine words from each review.locations.x file. 

