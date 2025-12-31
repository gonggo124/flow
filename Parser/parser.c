#include "parser.h"
#include <string.h>
#include <stdio.h>
#include <assert.h>

#define notNullM(value) if (value==NULL) return ERR_UNEXPECTED_EOF
#define shouldBeM(tok,tok_type,err) if (tok->type!=tok_type) return err

enum {
        ERR_UNEXPECTED_EOF = 1,
        ERR_UNEXPECTED_TOKEN,
};

static Token* next(Parser* p) {
        if (p->toks->size <= p->stack_offset+1) return NULL;
        else return &p->toks->data[++p->stack_offset];
}

void PAR_Parser_init(Parser* p, TokenList* toks, FILE* output_file) {
        //assert(false && "parser is not implemented");
        p->toks = toks;
        p->stack_offset = 0;
        p->output_file = output_file;
}

static int start_module_statement(Parser* p) {
        Token* next_tok = next(p);
        notNullM(next_tok);
        shouldBeM(next_tok,TOK_LITERAL_STRING,ERR_UNEXPECTED_TOKEN);

        // TODO: do some module system
        strcpy(p->cur_module,next_tok->value);
        fprintf(p->output_file, "prefix %s\n", p->cur_module);
        return 0;
}

/* typedef struct { */
/*         Token** data; */
/*         size_t size; */
/*         size_t cap; */
/* } */

#define reset(x) do { buf_ptr = x; buf_range = 1; } while(0)
#define push() do {                                                     \
        if (buf_ptr) {                                                  \
                fprintf(p->output_file,"( ");                           \
                for (size_t i = 0; i < buf_range; i++)                  \
                        fprintf(p->output_file,"%s",buf_ptr[i].value);  \
                fprintf(p->output_file," ) ");                          \
                reset(NULL);                                            \
        }                                                               \
        } while(0)
static int start_block_statement(Parser* p) {
        // TODO: do some block statement thing
        // TODO: 괴랄한 코드 리팩터링
        fprintf(p->output_file,"[ block statement ]\n");
        Token* buf_ptr = NULL;
        size_t buf_range = 1;
        int brace_stack = 1;
        char print_prev = 1;
        while (brace_stack) {
                Token* cur_tok = next(p);
                notNullM(cur_tok);

                switch (cur_tok->type) {
                case TOK_L_BRACE: brace_stack++; break;
                case TOK_R_BRACE: brace_stack--; break;
                }

                switch (cur_tok->type) {
                case TOK_SEMICOLON:
                        push();
                        fprintf(p->output_file,"\n");
                        reset(NULL);
                        break;
                case TOK_DOT:
                        if (buf_ptr->type == TOK_IDENTIFIER) {
                                print_prev = 0;
                        }
                        break;
                default: if (print_prev) push(); else print_prev = 1;
                }

                switch (cur_tok->type) {
                case TOK_SEMICOLON: break;
                default: if (buf_ptr) buf_range++; else reset(cur_tok);
                }
        }
        fprintf(p->output_file,"[ block statement end ]\n");
        
        return 0;
}
#undef reset
#undef push

/* if (cur_tok->value[0]=='-') { */
/*         char *ptr = cur_tok->value; */
/*         do if (*ptr=='_') *ptr = ' '; while(*++ptr); */
/*         ptr = (cur_tok->value)+1; */
/*         fprintf(p->output_file,"%s ",ptr); */
/*  } else */
/*         fprintf(p->output_file,"%s ",cur_tok->value); */
/* break; */

static int start_param_statement(Parser* p) {
        // TODO: do some params thing
        Token* param_tok = next(p);
        notNullM(param_tok);
        shouldBeM(param_tok,TOK_R_PAREN,ERR_UNEXPECTED_TOKEN);
        fprintf(p->output_file, "param\n");
        return 0;
}

static int start_func_definition(Parser* p) {
        (void)p;
        Token* next_tok = next(p);
        notNullM(next_tok);
        shouldBeM(next_tok,TOK_IDENTIFIER,ERR_UNEXPECTED_TOKEN);
        fprintf(p->output_file, "make %s.mcfunction\n", next_tok->value);

        Token* param_tok = next(p);
        notNullM(param_tok);
        shouldBeM(param_tok,TOK_L_PAREN,ERR_UNEXPECTED_TOKEN);
        start_param_statement(p);

        Token* block_start_tok = next(p);
        notNullM(block_start_tok);
        shouldBeM(block_start_tok,TOK_L_BRACE,ERR_UNEXPECTED_TOKEN);
        start_block_statement(p);
        
        return 0;
}

const char* PAR_get_error(int err_code) {
        static char buf[128];
        switch (err_code) {
        case ERR_UNEXPECTED_EOF:
                snprintf(buf,sizeof(buf),"Unexpected EOF");
                break;
        case ERR_UNEXPECTED_TOKEN:
                snprintf(buf,sizeof(buf),"Unexpected Token");
                break;
        }
        return buf;
}

int PAR_Parser_scan(Parser *p) {
        fprintf(p->output_file,"namespace namespace\n");
        int err_code = 0;
        for (;p->stack_offset < p->toks->size; p->stack_offset++) {
                Token *cur_tok = &p->toks->data[p->stack_offset];
                p->linenum = cur_tok->linenum+1;

                printf("[%d] at %d: \"%s\"\n", cur_tok->type, cur_tok->linenum+1, cur_tok->value);
                if (cur_tok->type == TOK_MODULE) err_code = start_module_statement(p);
                if (cur_tok->type == TOK_FUNC) err_code = start_func_definition(p);
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
