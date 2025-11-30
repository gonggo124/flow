#ifndef MNC_Tok
#define MNC_Tok
#include <stdio.h>

typedef int TOK_offset_t;
typedef int TOK_state_t;
typedef int TOK_line;
typedef int TOK_Type;

#define TOK_BUF_SIZE 256

typedef struct {
	TOK_Type type;
	char     value[TOK_BUF_SIZE];
	TOK_line line;
} Token;

typedef struct {
	FILE* file;  // 대상 파일.
	TOK_offset_t errpos; // 문제 발생지.
	TOK_state_t state; // 상태.
	char buf[TOK_BUF_SIZE];
	TOK_offset_t boffset;
	Token toks[169]; // TODO: Dynamic Array로 구현해야 함.
	TOK_offset_t offset;
} Tokenizer;


// TOK==Module Prefix, Tok==struct name

int TOK_Tokenizer_scan(Tokenizer *tok); // 0 == 정상, 나머지 == 파싱 에러 발생
int TOK_Tokenizer_init(Tokenizer *tok, FILE* file); // 0 == 정상, 나머지 == 에러 발생
void TOK_strerr(char *buf, int errno); // 파싱 에러 코드를 문자열로 변환

#endif // MNC_Tok
