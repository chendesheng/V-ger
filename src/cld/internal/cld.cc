#include "../public/compact_lang_det.h"
#include "../public/cld.h"
#include <stdlib.h>

void detect_language(char* utf8text, int length, /*out*/char* res) {    
    bool is_reliable;
    CLD2::Language language3[3];
    int percent3[3];
    int text_bytes;

    CLD2::Language l = CLD2::DetectLanguageSummary(utf8text, length, true, language3, percent3, &text_bytes, &is_reliable);

    printf("languagecode:%s\n", LanguageCode(l));
    
    if (percent3[1] > 30) {
        sprintf(res, "%s,%s", LanguageCode(language3[0]), LanguageCode(language3[1]));
    } else {
        sprintf(res, "%s", LanguageCode(l));
    }
}