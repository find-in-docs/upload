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


def extract_reviews(in_filename, out_filename):

    try:
        with open(in_filename) as fr:
            with open(out_filename, 'w') as fw: 
                for _ in range(0, NUM_CHUNKS):
                    chunk_len = config("CHUNK_LEN", cast=int)
                    for chunk in read_file_chunk(fr, chunk_len):
                        for element in chunk:
                            review = review_re.match(element).group(1)
                            print(review)

    except FileNotFoundError as e:
        print('Error occurred !!!\n\t\t{}', e)
        print("Could not find file: {} or {}", in_filename, out_filename)


if __name__ == '__main__':
    try:
        assert len(sys.argv) == 3
        extract_reviews(sys.argv[1], sys.argv[2])
    except AssertionError as e:
        print('Error occurred !!!\n\t\t{}', e)
        print('Usage: preprocess_text.py json_reviews_filepath')
