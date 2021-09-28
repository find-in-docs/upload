use std::sync::atomic::{AtomicUsize, Ordering};

use super::DocId;

static DOC_ID: AtomicUsize = AtomicUsize::new(0 as DocId);

pub fn get_doc_id() -> DocId {
    println!("DOC_ID: {}", DOC_ID.fetch_add(0, Ordering::SeqCst));
    DOC_ID.fetch_add(1, Ordering::SeqCst)
}
