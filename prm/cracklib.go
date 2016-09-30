package prm

// #cgo LDFLAGS: -lcrack
// #cgo CFLAGS:
// #include <stdlib.h>
// #include <string.h>
// #include <crack.h>
// char * password_check(char * pw ) {
//		char const * msg;
//		char * ret;
//		msg = FascistCheck(pw, GetDefaultCracklibDict());
//		if (msg){
//			ret = malloc(sizeof(char) * strlen(msg));
//			strcpy(ret,msg);
//		} else {
//			ret = malloc(sizeof(char) * 5);
//			strcpy(ret,"GOOD");
//		}
//		return ret;
// }
//
import "C"
import "unsafe"

func TestPassword(password string) string {
	var cchar *C.char = C.password_check(C.CString(password))
	defer C.free(unsafe.Pointer(cchar))
	return C.GoString(cchar)
}
