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

typedef int (*TOK_act_c)(Tok *tok, char *buf);

#define STATE_DEFAULT 0
#define STATE_STRING  1
#define STATE_NUMBER  2

#define CONDITION_LEN 1
#define NEXT_LEN CONDITION_LEN
#define ACTS_LEN CONDITION_LEN+1
typedef struct {
	const char* condition[CONDITION_LEN];
	const TOK_statenum next[NEXT_LEN];
	const TOK_act_c acts[ACTS_LEN];
} TOK_State;


int TOK_none(Tok *tok, char *buf); // do nothing
int TOK_putc(Tok *tok, char *buf); // put char
int TOK_puts(Tok *tok, char *buf); // put string
int TOK_puttoken(Tok *tok, char *buf); // put token
// int TOK_default_state_default(Tok *tok, char *buf); // 보류

const TOK_State TOK_States = {
	{
		.condition = {"{}();","\""},
		.next = {STATE_DEFAULT,STATE_STRING},
		.acts = {TOK_puts,TOK_none,TOK_putc}
	}
};

#define BUF_SIZE 256

#define ERR_BUF_OVERFLOW 0

void TOK_strerr(char *buf, int errno) {
	 switch (errno) {
		case ERR_BUF_OVERFLOW:
			strcpy(buf,"토큰의 길이가 초과되었습니다: 토큰의 최대 길이는 256 bytes 입니다");
		break;
	}
}

int Tok_scan(Tok *tok) {
	// TODO: tokenizer
	char chr;
	char buf[BUF_SIZE] = {0};
	int boffset = 0;
	TOK_offset_t offset = 0;
	while ((chr = fgetc(tok->file)) != EOF) {
		if (boffset < BUF_SIZE) {
			buf[boffset]=chr;
			boffset+=1;
			offset+=1;
//			if (strcmp(buf,
		} else {
			return TOK_ERR_BUF_OVERFLOW;
		}
			
	}
	return 0;
}

int Tok_init(Tok *tok, FILE *file) {
	tok->file = file;
	return 0;
}

int TOK_putc(Tok *tok, char *buf) {
	return 0;
}

int TOK_puts(Tok *tok, char *buf) {
	return 0;
}

int TOK_none(Tok *tok, char *buf) {
	(void)tok;
	(void)buf;
	return 0;
}
