#ifndef MNC_Tok
#define MNC_Tok
#include <stdio.h>

typedef int TOK_size_t;
typedef int TOK_state_t;
typedef int TOK_line_t;
typedef int TOK_Type;

#define TOK_BUF_SIZE 256

typedef struct {
	TOK_Type type;
	char     value[TOK_BUF_SIZE];
	TOK_line_t linenum;
} Token;

typedef struct {
	Token arr[169]; // TODO: dynamic array로 교체
	TOK_size_t size;
	TOK_size_t cap;
	TOK_size_t offset;
} TokenList;

void TOK_TokenList_push(TokenList *toklist, Token tok);
void TOK_TokenList_pop(TokenList *toklist);
void TOK_TokenList_clear(TokenList *toklist);
void TOK_TokenList_destroy(TokenList *toklist);

typedef struct {
	FILE* file;
	TOK_size_t errpos;

	TOK_state_t state;

	char buf[TOK_BUF_SIZE];
	TOK_size_t boffset;

	TokenList toks;
	TOK_size_t offset;

	TOK_line_t linenum;
} Tokenizer;


// TOK==Module Prefix, Tok==struct name

int TOK_Tokenizer_scan(Tokenizer *tok); // 0 == 정상, 나머지 == 파싱 에러 발생
int TOK_Tokenizer_init(Tokenizer *tok, FILE* file); // 0 == 정상, 나머지 == 에러 발생
void TOK_strerr(char *buf, int errno); // 파싱 에러 코드를 문자열로 변환

#endif // MNC_Tok
