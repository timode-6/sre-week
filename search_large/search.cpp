#include "search.h"

bool search(const std::vector<int>& data, int value)
{
    // Your code goes here.
    int base = 0;
    int len = data.size();
    while(len > 1){
        int half = len / 2;
        len -= half;
         __builtin_prefetch(&data[len / 2 - 1]);
         __builtin_prefetch(&data[half + len / 2 - 1]);
         base += (data[base + half - 1] < value) * half;
     } 
    return data[base] == value;
}
