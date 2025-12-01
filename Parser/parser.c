#include "parser.h"

void PAR_Parser_init(Parser *p, TokenList *toks) {
	p->toks = toks;
	p->offset = 0;
}

int PAR_Parser_scan(Parser *p) {
	
}
