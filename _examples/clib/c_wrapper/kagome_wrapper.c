#include "../bin/libkagome.h"

__attribute__((visibility("default")))
void* kagome_init(void) {
    return KagomeInit();
}
