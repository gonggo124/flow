#include "parser.h"
#include <string.h>
#include <stdio.h>
#include <assert.h>

void PAR_Parser_init(Parser *p, TokenList *toks) {
	assert(false && "parser is not implemented");
	p->toks = toks;
	p->offset = 0;
//	p->stack = (Token**)malloc(sizeof(Token*)*200); // TODO: make dynamic token array
	p->stack_offset = 0;
}

//static int handle_module_token(Parser *p) {
//	(void)p;
//	return 0;
//}

const char* PAR_get_error(int err_code) {
	(void)err_code;
	return "this is error";
}

int PAR_Parser_scan(Parser *p) {
	Token *token_stack[64] = {0};
	size_t stack_offset = 0;
	for (int i = 0; i < 169; i++) {
		Token *item = &p->toks->data[i];
		if (item->type == 0) break;

		p->linenum = item->linenum+1;
		token_stack[stack_offset] = item;
		stack_offset++;

		printf("[");
		for (size_t j = 0; j < stack_offset; j++) {
			printf("%d,",token_stack[j]->type);
		}
		printf("]\n");
		
	}
	printf("==============================\n");
	return 0;
}

//		printf("tok[%d] at %d: \"%s\"\n", item->type, item->linenum, item->value);
//		[1] + [10] = module statement
//		[3]|[4] + [2] + [5] + [params] + [6] = function statement
//		[7] + [expressions, something else...] + [8] = compound statement
//		params = [3]|[4] + [2] + [comma] + [params]
