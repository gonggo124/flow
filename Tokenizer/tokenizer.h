#ifndef MNC_Tok
#define MNC_Tok
#include <stdio.h>
#include "dynamic_array.h"

enum {
        TOK_EOF = 0,
        TOK_MODULE, // 'module'
        TOK_FUNC, // 'func'
        TOK_MAC, // 'mac'
        TOK_IDENTIFIER, // 나머지
        TOK_MO, // 'mo'
        TOK_DI, // 'di'
        TOK_L_PAREN, // '('
        TOK_R_PAREN, // ')'
        TOK_L_BRACE, // '{'
        TOK_R_BRACE, // '}'
        TOK_SEMICOLON, // ';'
        TOK_COLON, // ':'
        TOK_DOT, // '.'
        TOK_LITERAL_STRING, // "abc", "Hello123" 등..
        TOK_LITERAL_TEXT, // 'abc', 'Hello123' 등...
        TOK_LITERAL_NUMBER, // 123, 6767 등..
};

typedef size_t TOK_size_t;
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
        Token* data;
        TOK_size_t size;
        TOK_size_t cap;
} TokenList;

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
void TOK_Tokenizer_destroy(Tokenizer *tok);

#endif // MNC_Tok
