/*
 * Kagome C Wrapper Implementation
 *
 * Wraps Go-exported functions to provide a stable API for FFI users.
 * Built with -fvisibility=hidden so only kagome_* functions are visible.
 */

#include "../bin/libkagome.h"
#include "kagome_wrapper.h"

__attribute__((visibility("default")))
void* kagome_init(void) {
    return KagomeInit();
}

__attribute__((visibility("default")))
void kagome_destroy(void* handle) {
    KagomeDestroy(handle);
}

__attribute__((visibility("default")))
TokenArray* kagome_tokenize(void* handle, const char* input) {
    return KagomeTokenizeStruct(handle, (char*)input);
}

__attribute__((visibility("default")))
void kagome_free_token_array(TokenArray* arr) {
    KagomeFreeTokenArray(arr);
}

__attribute__((visibility("default")))
char* kagome_echo(const char* input) {
    return Echo((char*)input);
}

__attribute__((visibility("default")))
void kagome_echo_free(char* str) {
    EchoFree(str);
}
