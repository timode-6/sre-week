#include "search.h"

bool search(const std::vector<int>& data, int value)
{
    // Your code goes here.
    // int left = 0;
    // int right = data.size()-1;
    // int mid1 = 0;
    // int mid2 = 0;
    // while(left <= right){
    //     mid1 = left + (right - left) / 3;
    //     mid2 = right - (right - left) / 3;
    //     if(data[mid1] == value)
    //         return true;
    //     if(data[mid2] == value)
    //         return true;
    //     if(value < data[mid1]){
    //         right = mid1 - 1;
    //     }else if(value > data[mid2]){
    //         left = mid2 + 1;
    //     }else{
    //         left = mid1 + 1;
    //         right = mid2 - 1;
    //     }
    // }
    // return false;


    // int base = 0;
    // int len = data.size();
    // while(len > 1){
    //     int half = len / 2;
    //     base += (data[base + half - 1] < value) * half;
    //     len -= half;
    // } 
    // return data[base] == value;

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
