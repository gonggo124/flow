#ifndef MNC_Parser
#define MNC_Parser

#include "../Tokenizer/tokenizer.h"

typedef struct {
	TokenList *toks;
	TokenList stack;
	size_t stack_offset;
	TOK_line_t linenum;
} Parser;

void PAR_Parser_init(Parser *p, TokenList *toks);
int PAR_Parser_scan(Parser *p);

const char* PAR_get_error(int err_code);

#endif // MNC_Parser
