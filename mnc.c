#include <stdio.h>
#include <stdbool.h>
#include <string.h>
#include <stdlib.h>
#include <dirent.h>
#include <sys/stat.h>
#include <errno.h>
#include "Tokenizer/tokenizer.h"

// dtpk/src/main.mn
// -o
// dtpk/data

bool hasPrefix(const char* target, const char* prefix) {
	size_t t_len = strlen(target);
	size_t p_len = strlen(prefix);
	if (t_len < p_len) return false;
	return strncmp(target,prefix,p_len) == 0;
}

#define FLAG_NONE 0
#define FLAG_OUTPUT 1
typedef int flag_t;

void walkdir(const char* fn, void (*callback)(const char* path)) {
	DIR *dir;
	char path[1025];
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
		Tok tokenizer = {0};
		Tok_setFile(&tokenizer,file);
		Tok_scan(&tokenizer);
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
