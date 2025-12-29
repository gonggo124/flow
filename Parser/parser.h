#ifndef MNC_Parser
#define MNC_Parser

#include "../Tokenizer/tokenizer.h"

typedef struct {
        TokenList *toks;
        TokenList stack;
        size_t stack_offset;
        TOK_line_t linenum;
        char cur_module[TOK_BUF_SIZE];
        FILE* output_file;
} Parser;

void PAR_Parser_init(Parser *p, TokenList *toks, FILE* output_file);
int PAR_Parser_scan(Parser *p);

const char* PAR_get_error(int err_code);

#endif // MNC_Parser
