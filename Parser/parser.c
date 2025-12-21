#include "parser.h"
#include <string.h>
#include <stdio.h>
#include <assert.h>

#define optionM(value,err) if (value==NULL) return err

enum {
	ERR_UNEXPECTED_EOF,
	ERR_UNEXPECTED_TOKEN,
};

static Token* next(Parser* p) {
	p->stack_offset++;
	return TOK_TokenList_getN(p->toks,p->stack_offset);
}

void PAR_Parser_init(Parser* p, TokenList* toks) {
//	assert(false && "parser is not implemented");
	p->toks = toks;
	p->stack = TOK_TokenList_make(32);
	p->stack_offset = 0;
}

static int start_module_statement(Parser* p, Token* cur_tok) {
	(void)cur_tok;
	Token* next_tok = next(p); optionM(next_tok,ERR_UNEXPECTED_EOF);
	if (next_tok->type != TOK_IDENTIFIER) return ERR_UNEXPECTED_TOKEN;
	return 0;
}

const char* PAR_get_error(int err_code) {
	(void)err_code;
	return "something went wrong";
}

int PAR_Parser_scan(Parser *p) {
	Token *cur_tok = NULL;
	int err_code = 0;
	while ((cur_tok = TOK_TokenList_getN(p->toks,p->stack_offset))!=NULL) {
		p->linenum = cur_tok->linenum+1;

		printf("[%d] at %d: \"%s\"\n", cur_tok->type, cur_tok->linenum+1, cur_tok->value);
		if (cur_tok->type == TOK_MODULE) err_code = start_module_statement(p,cur_tok);

		p->stack_offset++;
	}


	if (err_code < 0) return err_code;
	
	printf("==============================\n");
	return 0;
}

//		printf("tok[%d] at %d: \"%s\"\n", item->type, item->linenum, item->value);
//		[1] + [10] = module statement
//		[3]|[4] + [2] + [5] + [params] + [6] = function statement
//		[7] + [expressions, something else...] + [8] = compound statement
//		params = [3]|[4] + [2] + [comma] + [params]
