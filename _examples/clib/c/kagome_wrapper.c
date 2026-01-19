#include "../bin/libkagome.h"

__attribute__((visibility("default")))
uintptr_t kagome_init(const char* path) {
    return KagomeInit((char*)path);
}
