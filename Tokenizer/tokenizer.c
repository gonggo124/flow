#include "tokenizer.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

enum {
	TOK_MODULE = 1, // 'module'
	TOK_IDENTIFIER, // 나머지
	TOK_MO, // 'mo'
	TOK_DI, // 'di'
	TOK_L_PAREN, // '('
	TOK_R_PAREN, // ')'
	TOK_L_BRACE, // '{'
	TOK_R_BRACE, // '}'
	TOK_SEMICOLON, // ';'
	TOK_LITERAL_STRING // "abc", 123 등..
};

typedef int (*TOK_act_c)(Tokenizer *tok, char chr);

enum {
	STATE_NORMAL = 0,
	STATE_STRING,
	STATE_STRING_SKIP,
	STATE_NUMBER
};

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
		.next      = {STATE_STRING},
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

enum {
	ERR_BUF_OVERFLOW,
	ERR_UNEXPECTED_TOKEN
};

void TOK_strerr(char *buf, int errno) {
	 switch (errno) {
		case ERR_BUF_OVERFLOW:
			strcpy(buf,"토큰의 길이가 초과되었습니다: 토큰의 최대 길이는 256 bytes 입니다");
		break;
	}
}

int TOK_Tokenizer_scan(Tokenizer *tok) {
	char chr;
	while ((chr = fgetc(tok->file)) != EOF) {
		TOK_State cstate = States[tok->state]; // current state
		if (chr == '\n') tok->linenum++;
		for (int i = 0; i < cstate.size+1; i++) {
			if (i < cstate.size && !strchr(cstate.condition[i],chr)) continue;
			int err = cstate.acts[i](tok,chr);
			if (err) return err;
			tok->state=cstate.next[i];
			break;
		}
	}
	return 0;
}

int TOK_Tokenizer_init(Tokenizer *tok, FILE *file) {
	tok->state = STATE_NORMAL;
	tok->file = file;
	memset(tok->buf,0,TOK_BUF_SIZE);
	TOK_TokenList_clear(&tok->toks);
	return 0;
}

static int Tokenize(Token *tok, char *buf) {
	if (buf[0]==0) return 0;
	// TODO: Do Tokenize Shit

	strcpy(tok->value,buf);

	size_t len = strlen(buf);

	// type specification
	if (len==1) {
		switch(buf[0]) {
			case '(': tok->type = TOK_L_PAREN; break;
			case ')': tok->type = TOK_R_PAREN; break;
			case '{': tok->type = TOK_L_BRACE; break;
			case '}': tok->type = TOK_R_BRACE; break;
			case ';': tok->type = TOK_SEMICOLON; break;
			default: return ERR_UNEXPECTED_TOKEN;
		}
	} else {
		if (strcmp(buf,"module")==0) tok->type = TOK_MODULE;
		else if (strcmp(buf,"mo")==0) tok->type = TOK_MO;
		else if (strcmp(buf,"di")==0) tok->type = TOK_DI;
		else tok->type = TOK_IDENTIFIER;
	}
	// type spec end

	return 0;
}

static int none(Tokenizer *tok, char chr) {
	(void)tok;
	(void)chr;
	return 0;
}
static int pushc(Tokenizer *tok, char chr) {
	if (tok->boffset < TOK_BUF_SIZE) {
		tok->buf[tok->boffset]=chr;
		tok->boffset++;
	} else {
		return ERR_BUF_OVERFLOW;
	}
	return 0;
}
static Token make_token(Tokenizer *tok) {
	Token new_token = {0};
	new_token.linenum = tok->linenum;
	return new_token;
}
static int tokc(Tokenizer *tok, char chr) {
	if (tok->buf[0]==0) return 0;
	(void)chr; // trash

	Token new_token = make_token(tok);
	Tokenize(&new_token,tok->buf);
	if (tok->state==STATE_STRING) new_token.type=TOK_LITERAL_STRING;

	TOK_TokenList_push(&(tok->toks),new_token);

	memset(tok->buf,0,TOK_BUF_SIZE); // empty buf
	tok->boffset = 0; // empty buf

	return 0;
}
static int toka(Tokenizer *tok, char chr) {
	tokc(tok,chr);
	pushc(tok,chr);
	tokc(tok,chr);
	return 0;
}

void TOK_TokenList_push(TokenList *toklist, Token tok) {
	if (toklist->offset < 169) {
		toklist->arr[toklist->offset]=tok;
		toklist->offset++;
	}
}

void TOK_TokenList_pop(TokenList *toklist) {
	(void)toklist;
}	
void TOK_TokenList_clear(TokenList *toklist) {
	memset(toklist,0,169*sizeof(Token));
	(void)toklist;
}
void TOK_TokenList_destroy(TokenList *toklist) {
	(void)toklist;
}
