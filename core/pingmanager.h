// Copyright (c) 2022 Husarnet sp. z o.o.
// Authors: listed in project_root/README.md
// License: specified in project_root/LICENSE.txt
#pragma once
#include <vector>
#include <chrono>
#include <mutex>
#include <condition_variable>
#include "husarnet/husarnet_manager.h"

class PingManager {
    HusarnetManager * husarnetManager;
    uint16_t currentSequence;
    std::mutex m;
    std::condition_variable v;

 public:
    explicit PingManager(HusarnetManager *husarnetManager) :husarnetManager(husarnetManager), currentSequence(0){};
    std::pair<bool,std::chrono::milliseconds> sendPing();
    void notify();
    uint16_t getCurrentSequence();
};