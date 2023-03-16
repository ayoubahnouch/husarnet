// Copyright (c) 2022 Husarnet sp. z o.o.
// Authors: listed in project_root/README.md
// License: specified in project_root/LICENSE.txt
#pragma once
#include "husarnet/ngsocket_manager.h"

class HusarnetManager;

class MulticastDiscoveryServer {
 private:
  HusarnetManager* manager;
  NgSocketManager* ngsocket;

  void serverThread();
  void onPacket(InetAddress address, string_view data);

 public:
  MulticastDiscoveryServer(HusarnetManager* manager, NgSocketManager* ngsocket);

  void start();
};
