#include "parser.h"
#include <string.h>
#include <stdio.h>

void PAR_Parser_init(Parser *p, TokenList *toks) {
	p->toks = toks;
	p->offset = 0;
}

typedef enum {
	STATE_NORMAL,
	STATE_MODULE, // starts with 'module'
	STATE_ERROR_UNEXPECTED_TOKEN = -1,
} State;

const char* PAR_get_error(int err_code) {
	switch (err_code) {
	case STATE_ERROR_UNEXPECTED_TOKEN:
		return "unexpected token";
	default:
		return "";
	}
}

static State do_state_normal(Token *tok, Token **stack, size_t *stack_offset) {
	(void)stack;
	(void)stack_offset;
	switch (tok->type) {
	case TOK_MODULE: return STATE_MODULE;
	default: return STATE_ERROR_UNEXPECTED_TOKEN;
	}
	return STATE_NORMAL;
}

static State do_state_module(Token *tok, Token **stack, size_t *stack_offset) {
	(void)stack;
	switch (tok->type) {
	case TOK_LITERAL: 
		memset(stack,0,(*stack_offset) * sizeof(char));
		*stack_offset = 0;
		return STATE_NORMAL;
	default: return STATE_ERROR_UNEXPECTED_TOKEN;
	}
	return STATE_NORMAL;
}



int PAR_Parser_scan(Parser *p) {
	Token *token_stack[64] = {0};
	size_t stack_offset = 0;
	State state = STATE_NORMAL;
	for (int i = 0; i < 169; i++) {
		Token *item = &p->toks->arr[i];
		if (item->type == 0) break;

		p->linenum = item->linenum+1;
		token_stack[stack_offset] = item;
		stack_offset++;

		printf("[");
		for (size_t j = 0; j < stack_offset; j++) {
			printf("%d,",token_stack[j]->type);
		}
		printf("]\n");


		if (state == STATE_NORMAL) state = do_state_normal(item,token_stack,&stack_offset);
		else if (state == STATE_MODULE) state = do_state_module(item,token_stack,&stack_offset);

		if (state < 0) return state;

//		printf("tok[%d] at %d: \"%s\"\n", item->type, item->linenum, item->value);
//		[1] + [10] = module statement
//		[3]|[4] + [2] + [5] + [params] + [6] = function statement
//		[7] + [expressions, something else...] + [8] = compound statement
//		params = [3]|[4] + [2] + [comma] + [params]
		
	}
	printf("==============================\n");
	return 0;
}
