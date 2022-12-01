// Copyright (c) 2022 Husarnet sp. z o.o.
// Authors: listed in project_root/README.md
// License: specified in project_root/LICENSE.txt
#include "husarnet/pingmanager.h"

std::pair<bool,std::chrono::milliseconds> PingManager::sendPing(){
    std::chrono::steady_clock::time_point begin = std::chrono::steady_clock::now();
    husarnetManager->getNGSocket()->sendPing(currentSequence);
    std::unique_lock lk(m);
    auto endTime = std::chrono::steady_clock::now() + std::chrono::seconds(5);
    auto result = v.wait_until(lk,endTime);
    if(result == std::cv_status::timeout){
        currentSequence++;
        return std::pair<bool,std::chrono::milliseconds>(false,std::chrono::steady_clock::duration::zero().count());

    }else{
        std::chrono::steady_clock::time_point end = std::chrono::steady_clock::now();
        currentSequence++;
        return std::pair<bool,std::chrono::milliseconds>(true,std::chrono::duration_cast<std::chrono::milliseconds>(end - begin).count());
    }
};

void PingManager::notify(){
    v.notify_all();
};

uint16_t PingManager::getCurrentSequence(){
    return currentSequence;
};