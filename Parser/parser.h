#ifndef MNC_Parser
#define MNC_Parser

#include "../Tokenizer/tokenizer.h"

typedef struct {
	TokenList *toks;
	TOK_size_t offset;
} Parser;

void PAR_Parser_init(Parser *p, TokenList *toks);
int PAR_Parser_scan(Parser *p);

#endif // MNC_Parser
