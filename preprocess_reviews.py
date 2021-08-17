"""
 Filename: preprocess_text.py
   Take a list of reviews and apply these transformations
   to each one:
   Lowercase characters
   Remove punctuation
   Remove stopwords
   Lemmatize
"""

import re
import sys
from itertools import islice
from decouple import config
from nltk.corpus import stopwords
from nltk.tokenize import RegexpTokenizer

CHUNK_LEN_DEFAULT = 2
NUM_CHUNKS = 3

review_re = re.compile(r'^.*\"text\"\:\"(.+?)\".*$')


def read_file_chunk(file_handle, chunk_len=CHUNK_LEN_DEFAULT):

    num_chunks = 0
    chunk = []
    for line in file_handle:
        chunk.append(line)

        if chunk_len == len(chunk):
            yield chunk
            chunk = []
            num_chunks += 1

        if num_chunks > NUM_CHUNKS:
            break

    if len(chunk) > 0:
        yield chunk


def extract_reviews(in_filename):

    try:
        with open(in_filename) as fr:
            for _ in range(0, NUM_CHUNKS):
                chunk_len = config("CHUNK_LEN", cast=int)
                for chunk in read_file_chunk(fr, chunk_len):
                    for element in chunk:
                        review = review_re.match(element).group(1)
                        yield review

    except FileNotFoundError as e:
        print('Error occurred !!!\n\t\t{}', e)
        print("Could not find file: {} or {}", in_filename)


def normalize_reviews(review):

    # Remove apostrophe from the whole text before tokenizing.
    # Tokenizing splits "Don't" into "Don" and "t".
    # We would rather have "Dont"
    # Also remove all number characters.
    remove_chars = set(["'", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"])
    review = ''.join([c for c in review if c not in remove_chars])

    words = RegexpTokenizer(r'\w+').tokenize(review)   # split text on word boundaries.
    words = [word for word in words if len(word) != 0] # remove empty words.
    words = [word.lower() for word in words]        

    stop_words = set(stopwords.words('english'))
    words = [w for w in words if w not in stop_words]

    return ' '.join(words)


if __name__ == '__main__':
    try:
        assert len(sys.argv) == 3
        in_filename = sys.argv[1]
        out_filename = sys.argv[2]

        with open(out_filename, "w") as fw:
            for review in extract_reviews(in_filename):
                normalized_review = normalize_reviews(review) + '\n'
                fw.write(normalized_review)

    except FileNotFoundError as e:
        print('Error occurred !!!\n\t\t{}', e)
        print("Could not find file: {}", out_filename)

    except AssertionError as e:
        print('Error occurred !!!\n\t\t{}', e)
        print('Usage: preprocess_text.py json_reviews_filepath output_filepath')
