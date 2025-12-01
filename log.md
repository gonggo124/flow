2025.11.23
- 입력 폴더와 출력 폴더 입력 구현.
- 입력 폴더 탐색 구현.

2025.11.26
- Lexer -> Tokenizer로 이름 변경.
- '\n' 단위 파일 탐색 -> char 단위 탐색으로 변경.

2025.11.27
- 토큰 추가.
- 에러 반환 시스템 추가.
- Tok_scan이 buf에 fgetc 결과를 밀어넣게 함.
- TOK_offset_t, TOK_act_c 등 타입 규격 추가.
- 상태 머신 방식으로 개발 계획 중 -> states, acts 추가.
- Tok_State 추가. header로 옮겨야 할듯? 그리고 Tok에 state field도 추가해야함.

2025.11.29
- 상태머신 일부 구현.
- Tok에 상태 추가.
- 상태 표 일부 정의.
- 변수 이름 일부 수정. Tok_State -> State(typedef는 tok.c 내부에서만 인식되기 때문에 필요 없을 듯?) 등.
- header로 옮기지 않았음. 필요없어서.

2025.11.30
- 상태표 구현. (number 제외)
- 이름 개편. (Tok->Tokenizer)
- Token 구조체 추가.
- Tokenizer에 buf, boffset, offset, toks(임시), state 필드 추가.
- putc, none 등 static으로 변경.
- 이외에도 일부 수정.

2025.11.30 2차 개발
- Tokenizer 상당 부분 개발
- 상태 머신 구현 완료
- 이름 개편2222(putc->pushc => 이름 겹침;;, puts->toka, putt->tokc => 명확성을 위해.)
- 이제 샘플코드 토큰으로 변환 가능
- none, pushc, tokc, toka 상당 부분 개발 완료
- 필요 함수 추가. (Tokenize, TOK_Tokenizer_push 등)

2025.12.1
- Dynamic Array 준비 (TokenList 만듦. 아직 고정 길이 배열)
- Parser 준비.
- TODO: Implement Tokenizing
- TODO: Dynamic Array
