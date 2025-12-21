#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "tokenizer.h"

#define must(expr) do {\
		int _e = (expr); \
		if (_e) return _e; \
	} while(0)

typedef int (*TOK_act_c)(Tokenizer *tok, char chr);

enum {
	STATE_NORMAL = 0,
	STATE_WORD,
	STATE_STRING,
	STATE_STRING_SKIP,
	STATE_NUMBER
};

#define CONDITION_LEN 7
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
static int tokc(Tokenizer *tok, char chr); // tokenize current(only buf, ignore chr)
static int toka(Tokenizer *tok, char chr); // tokenize all(buf and chr)
static int tokcw(Tokenizer *tok, char chr); // tokenize current end wind back(offset-=1)

#define IDCHARS "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
#define NUMBERCHARS "0123456789"

// MARK: #001
static const TOK_State States[] = {
	{ // normal state
		.condition = {IDCHARS   ,NUMBERCHARS ," \n\r\t"   ,"{}();"     ,"\""         /*DEFAULT*/ },
		.next      = {STATE_WORD,STATE_NUMBER,STATE_NORMAL,STATE_NORMAL,STATE_STRING,STATE_NORMAL},
		.acts      = {pushc     ,pushc       ,tokc        ,toka        ,none        ,pushc       },
		.size      = 5 // TODO: it's not DRY. FIX IT!!! I WAS ABOUT TO SPEND 10 HOURS TO FIX THE BUG BECAUSE OF IT!!
	},
	{ // word state
		.condition = {IDCHARS NUMBERCHARS /*DEFAULT*/ },
		.next      = {STATE_WORD         ,STATE_NORMAL},
		.acts      = {pushc              ,tokcw       },
		.size      = 1
	},
	{ // string state
		.condition = {"\""        ,"\\"              /*DEFAULT*/ },
		.next      = {STATE_NORMAL,STATE_STRING_SKIP,STATE_STRING},
		.acts      = {tokc        ,none             ,pushc       },
		.size      = 2
	},
	{ // string-skip state
		.condition = {/*DEFAULT*/ },
		.next      = {STATE_STRING},
		.acts      = {pushc       },
		.size      = 0
	},
	{ // number state
		.condition = {NUMBERCHARS ,/*DEFAULT*/ },
		.next      = {STATE_NUMBER,STATE_NORMAL},
		.acts      = {pushc       ,tokcw       },
		.size      = 1
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

// MARK: #003
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

int TOK_Tokenizer_wind(Tokenizer *tok, long amount) {
	return fseek(tok->file,-amount,SEEK_CUR);
}

int TOK_Tokenizer_init(Tokenizer *tok, FILE *file) {
	tok->state = STATE_NORMAL;
	tok->file = file;
	memset(tok->buf,0,TOK_BUF_SIZE);
	tok->toks = TOK_TokenList_make(64);
	return 0;
}

void TOK_Tokenizer_destroy(Tokenizer *tok) {
	TOK_TokenList_destroy(&tok->toks);
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
		else if (strcmp(buf,"func")==0) tok->type = TOK_FUNC;
		else if (strcmp(buf,"mac")==0) tok->type = TOK_MAC;
		else if (strcmp(buf,"mo")==0) tok->type = TOK_MO;
		else if (strcmp(buf,"di")==0) tok->type = TOK_DI;
		else tok->type = TOK_IDENTIFIER;
	}
	// type spec end

	return 0;
}


// MARK: #002
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
	if (tok->state==STATE_STRING || tok->state==STATE_NUMBER) new_token.type=TOK_LITERAL;

	TOK_TokenList_push(&(tok->toks),new_token);

	memset(tok->buf,0,TOK_BUF_SIZE); // empty buf
	tok->boffset = 0; // empty buf

	return 0;
}
static int tokcw(Tokenizer *tok, char chr) {
	must(tokc(tok,chr));
	TOK_Tokenizer_wind(tok,1);
	return 0;
}
static int toka(Tokenizer *tok, char chr) {
	must(tokc(tok,chr));
	must(pushc(tok,chr));
	must(tokc(tok,chr));
	return 0;
}
