#include "tokenizer.h"
#include <stdlib.h>
#include <assert.h>
#include <string.h>

TokenList TOK_TokenList_make(size_t start_cap) {
	Token* new_data = malloc(sizeof(Token)*start_cap);
	if (new_data == NULL) {
		perror("malloc failed");
		exit(1);
	}
	TokenList new_list = {
		.data = new_data,
		.size = 0,
		.cap = start_cap,
	};
	return new_list;
}

void TOK_TokenList_extend(TokenList* list, size_t extend_amount) {
	assert(list->data!=NULL);
	size_t new_cap = list->cap+extend_amount;
	list->data = malloc(sizeof(Token)*new_cap);
	if (list->data == NULL) {
		perror("malloc failed");
		exit(1);
	}
	list->cap = new_cap;
}

void TOK_TokenList_push(TokenList* list, Token tok) {
	assert(list->data!=NULL);
	if (list->cap < list->size) {
		// 용량 초과
		TOK_TokenList_extend(list,8);
	}
	list->data[list->size] = tok;
	list->size++;
}

void TOK_TokenList_clear(TokenList *list) {
	memset(list,0,sizeof(Token)*list->size);
	list->size = 0;
}

void TOK_TokenList_destroy(TokenList *list) {
	if (list->data==NULL) return;
	free(list->data);
	list = NULL;
}

