#include "search.h"
#include <immintrin.h>

bool search(const std::vector<int>& data, int value)
{
    // Your code goes here.
    int first = 0;
    int last = data.size();
    int length = last;
    constexpr int entries_per_256KB = 256 * 1024 / sizeof(data);
    if (length >= entries_per_256KB) {
        constexpr int num_per_cache_line = std::max(64 / int(sizeof(data)), 1);
        while (length >= 3 * num_per_cache_line) {
            int half = length / 2;
            __builtin_prefetch(&data[first + half / 2]);
            int first_half1 = first + (length - half);
            __builtin_prefetch(&data[first_half1 + half / 2]);
            first = data[first + half] < value ? first_half1 : first;
            length = half;
        }
    }
    while (length > 0) {
        int half = length / 2;
        int first_half1 = first + (length - half);
        first = data[first + half]< value ? first_half1 : first;
        length = half;
    }
    return data[first]==value;
}
