#ifndef CLD_H
#define CLD_H

#ifdef __cplusplus 
extern "C" { 
#endif 

	void detect_language(char* utf8text, int length, char* res);
#ifdef __cplusplus 
} /* extern "C" */ 
#endif 
#endif
