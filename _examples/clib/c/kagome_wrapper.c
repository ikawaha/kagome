#include "../bin/libkagome.h"

__attribute__((visibility("default")))
void* kagome_init(const char* path) {
    return KagomeInit((char*)path);
}
