#include "parser.h"
#include <string.h>
#include <stdio.h>
#include <assert.h>

#define notNullM(value,err) if (value==NULL) return err
#define shouldBeM(tok,tok_type,err) if (tok->type!=tok_type) return err

enum {
        ERR_UNEXPECTED_EOF = 1,
        ERR_UNEXPECTED_TOKEN,
};

static Token* next(Parser* p) {
        if (p->toks->size <= p->stack_offset+1) return NULL;
        else return &p->toks->data[++p->stack_offset];
}

void PAR_Parser_init(Parser* p, TokenList* toks) {
        p->toks = toks;
        p->stack_offset = 0;
}

static int start_module_statement(Parser* p, Token* cur_tok) {
        (void)cur_tok;
        Token* next_tok = next(p);
        notNullM(next_tok,ERR_UNEXPECTED_EOF);
        shouldBeM(next_tok,TOK_LITERAL,ERR_UNEXPECTED_TOKEN);

        // TODO: do some module system
        printf("set module\n");
        return 0;
}

static int start_func_definition(Parser* p, Token* cur_tok) {
        // TODO: do some func definition
        (void)cur_tok;
        Token* next_tok = next(p);
        notNullM(next_tok,ERR_UNEXPECTED_EOF);
        shouldBeM(next_tok,TOK_IDENTIFIER,ERR_UNEXPECTED_TOKEN);
        return 0;
}

const char* PAR_get_error(int err_code) {
        static char buf[128];
        snprintf(buf,sizeof(buf),"something went wrong with %d", err_code);
        return buf;
}

int PAR_Parser_scan(Parser *p) {
        int err_code = 0;
        for (;p->stack_offset < p->toks->size; p->stack_offset++) {
                Token *cur_tok = &p->toks->data[p->stack_offset];
                p->linenum = cur_tok->linenum+1;

                printf("[%d] at %d: \"%s\"\n", cur_tok->type, cur_tok->linenum+1, cur_tok->value);
                if (cur_tok->type == TOK_MODULE) err_code = start_module_statement(p,cur_tok);
                if (cur_tok->type == TOK_FUNC) err_code = start_func_definition(p,cur_tok);
                if (err_code != 0) return err_code;
        }

        printf("==============================\n");
        return 0;
}

//		printf("tok[%d] at %d: \"%s\"\n", item->type, item->linenum, item->value);
//		[1] + [10] = module statement
//		[3]|[4] + [2] + [5] + [params] + [6] = function statement
//		[7] + [expressions, something else...] + [8] = compound statement
//		params = [3]|[4] + [2] + [comma] + [params]
