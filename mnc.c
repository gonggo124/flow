#include <stdio.h>
#include <stdbool.h>
#include <string.h>

#include <stdlib.h>
#include <dirent.h>
#include <sys/stat.h>
#include <errno.h>
#include "Tokenizer/tokenizer.h"
#include "Parser/parser.h"

bool hasPrefix(const char* target, const char* prefix) {
	size_t t_len = strlen(target);
	size_t p_len = strlen(prefix);
	if (t_len < p_len) return false;
	return strncmp(target,prefix,p_len) == 0;
}

#define FLAG_NONE 0
#define FLAG_OUTPUT 1
typedef int flag_t;

#define PATH_SIZE 1025

void walkdir(const char* fn, void (*callback)(const char* path)) {
	DIR *dir;
	char path[PATH_SIZE];
	struct dirent *entry;
	struct stat info;

	dir = opendir(fn);
	if (dir==NULL) {
		perror("opendir() error");
	} else {
		while ((entry=readdir(dir))!=NULL) {
			if (entry->d_name[0]!='.') {
				strcpy(path,fn);
				strcat(path,"/");
				strcat(path,entry->d_name);
				if (stat(path,&info))
					fprintf(stderr, "stat() error on %s: %s\n", path, strerror(errno));
				else if (S_ISDIR(info.st_mode))
					walkdir(path,callback);
				else
					callback(path);
			}
		}
		closedir(dir);
	}
}

void wd_callback(const char* path) {
	FILE *file = fopen(path,"r");
	if (file==NULL)
		perror("fopen() error");
	else {
		char copied_path[PATH_SIZE]; strcpy(copied_path,path);
		char *extender = strrchr(copied_path,'.');
		if (extender) extender++;
		else { printf("couldn't find file extender"); return; }

		if (strcmp(extender,"mn")!=0) return;

		Tokenizer tokenizer = {0};
		TOK_Tokenizer_init(&tokenizer,file);
		TOK_Tokenizer_scan(&tokenizer);
		for (size_t i = 0; i < tokenizer.toks.size; i++) {
			Token *item = TOK_TokenList_getN(&tokenizer.toks,i);
			printf("tok[%d] at %d: \"%s\"\n", item->type, item->linenum+1, item->value);
		}
		printf("==============================\n");
		
		Parser parser = {0};
		PAR_Parser_init(&parser,&tokenizer.toks);
		int parse_err = PAR_Parser_scan(&parser);
		if (parse_err != 0) {
			printf("Parse Error: %s at %s:%d\n",PAR_get_error(parse_err),path,parser.linenum);
		}

		TOK_Tokenizer_destroy(&tokenizer); // Parser도 toks 써야해서 여기서 destroy

		fclose(file);
	}
}

int main(int argc, char** args) {
	flag_t flag = FLAG_NONE;
	char *input_dir;
	char *output_dir;
	for (int i = 1; i < argc; i++) {
		char *cur = args[i];
		if (cur[0]=='-') {
			if (cur[1]=='o') flag=FLAG_OUTPUT;
		} else {
			switch (flag) {
			case FLAG_NONE: {
				input_dir=cur;
			}; break;
			case FLAG_OUTPUT: {
				output_dir=cur;
			}; break;
			default:
				perror("Unknown Flag: %d");
			}
			flag=FLAG_NONE;
		}
	}
	printf("input: %s\noutput: %s\n",input_dir,output_dir);

	walkdir(input_dir,wd_callback);
}
