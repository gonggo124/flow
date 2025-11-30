#include "tokenizer.h"
#include <stdio.h>
#include <string.h>

#define TOK_MODULE 0 // 'module'
#define TOK_IDENTIFIER 1 // 나머지
#define TOK_MO 2 // 'mo'
#define TOK_DI 3 // 'di'
#define TOK_L_PAREN 4 // '('
#define TOK_R_PAREN 5 // ')'
#define TOK_L_BRACE 6 // '{'
#define TOK_R_BRACE 7 // '}'
#define TOK_SEMICOLON 8 // ';'
#define TOK_LITERAL 9 // "abc", 123 등..

typedef int (*TOK_act_c)(Tokenizer *tok, char chr);

#define STATE_NORMAL 0
#define STATE_STRING 1
#define STATE_STRING_SKIP 2
#define STATE_NUMBER 3

#define CONDITION_LEN 3
#define NEXT_LEN CONDITION_LEN+1
#define ACTS_LEN CONDITION_LEN+1
typedef struct {
	char* condition[CONDITION_LEN];
	TOK_state_t next[NEXT_LEN];
	TOK_act_c acts[ACTS_LEN];
	int size;
} TOK_State;


static int none(Tokenizer *tok, char chr); // do nothing
static int pushc(Tokenizer *tok, char chr); // push char to buf
static int tokc(Tokenizer *tok, char chr); // tokenize current(only buf)
static int toka(Tokenizer *tok, char chr); // tokenize all(buf and chr)
// int TOK_default_state_default(Tokenizer *tok, char chr); // 보류

static const TOK_State States[] = {
	{ // normal state
		.condition = {"{}();"     ,"\""        ," \n\r\t"                },
		.next      = {STATE_NORMAL,STATE_STRING,STATE_NORMAL,STATE_NORMAL},
		.acts      = {toka        ,none        ,tokc        ,pushc       },
		.size      = 3
	},
	{ // string state
		.condition = {"\""        ,"\\"                          },
		.next      = {STATE_NORMAL,STATE_STRING_SKIP,STATE_STRING},
		.acts      = {tokc        ,none             ,pushc       },
		.size      = 2
	},
	{ // string-skip state
		.condition = {},
		.next      = {STATE_NORMAL},
		.acts      = {pushc},
		.size      = 0
	},
	{ // number state
		.condition = {},
		.next      = {},
		.acts      = {},
		.size      = 0
	}
};

#define ERR_BUF_OVERFLOW 0

void TOK_strerr(char *buf, int errno) {
	 switch (errno) {
		case ERR_BUF_OVERFLOW:
			strcpy(buf,"토큰의 길이가 초과되었습니다: 토큰의 최대 길이는 256 bytes 입니다");
		break;
	}
}

int TOK_Tokenizer_scan(Tokenizer *tok) {
	// TODO: tokenizer
	char chr;
	while ((chr = fgetc(tok->file)) != EOF) {
		TOK_State cstate = States[tok->state]; // current state
		//printf("current state: %d\n",tok->state);
		printf("chr: %c ",chr);
		for (int i = 0; i < cstate.size+1; i++) {
			if (i < cstate.size && !strchr(cstate.condition[i],chr)) continue;
			if (i < cstate.size) printf("i: %d ",i);
			int err = cstate.acts[i](tok,chr);
			if (err) return err;
			tok->state=cstate.next[i];
			break;
		}
		printf("\n");
	}
	return 0;
}

int TOK_Tokenizer_push(Tokenizer *tokenizer, Token tok) {
	if (tokenizer->offset < 169) {
		tokenizer->toks[tokenizer->offset]=tok;
		tokenizer->offset++;
	} else {
		printf("overflow shit\n");
		return 1;
	}
	return 0;
}

int TOK_Tokenizer_init(Tokenizer *tok, FILE *file) {
	tok->state = STATE_NORMAL;
	tok->file = file;
	memset(tok->buf,0,TOK_BUF_SIZE);
	return 0;
}

static Token Tokenize(char *buf) {
	// TODO: Do Tokenize Shit
	Token new_tok = {0};
	new_tok.type = 1;
	strcpy(new_tok.value,buf);
	return new_tok;
}

static int none(Tokenizer *tok, char chr) {
	(void)tok;
	(void)chr;
	printf("none ");
	return 0;
}
static int pushc(Tokenizer *tok, char chr) {
	printf("putchar ");
	if (tok->boffset < TOK_BUF_SIZE) {
		tok->buf[tok->boffset]=chr;
		tok->boffset++;
	} else {
		return ERR_BUF_OVERFLOW;
	}
	return 0;
}
static int tokc(Tokenizer *tok, char chr) {
	if (tok->buf[0]==0) return 0;
	printf("tokc ");
	(void)chr;
	TOK_Tokenizer_push(tok,Tokenize(tok->buf));
	memset(tok->buf,0,TOK_BUF_SIZE);
	tok->boffset = 0;
	return 0;
}
static int toka(Tokenizer *tok, char chr) {
	printf("toka ");
	tokc(tok,chr);
	pushc(tok,chr);
	tokc(tok,chr);
	return 0;
}
