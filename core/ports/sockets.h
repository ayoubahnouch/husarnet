// Copyright (c) 2022 Husarnet sp. z o.o.
// Authors: listed in project_root/README.md
// License: specified in project_root/LICENSE.txt

// Wrapper around OS sockets
#pragma once

#include <functional>
#include <memory>
#include <vector>

#ifndef _WIN32
#include <netinet/in.h>
#include <sys/socket.h>
#else
#include <winsock2.h>

#include <iphlpapi.h>
#include <ws2ipdef.h>
#include <ws2tcpip.h>
#endif

#include <condition_variable>
#include <mutex>
#include <thread>

#include "husarnet/ipaddress.h"
#include "husarnet/string_view.h"

using mutex_guard = std::lock_guard<std::recursive_mutex>;

namespace OsSocket {
  struct FramedTcpConnection;

  using PacketCallack = std::function<void(InetAddress, string_view)>;
  using TcpDataCallback = std::function<void(const std::string&)>;
  using TcpErrorCallback =
      std::function<void(std::shared_ptr<FramedTcpConnection>)>;

  // Helpers
  InetAddress ipFromSockaddr(struct sockaddr_storage st);

  // UDP
  bool udpListenUnicast(
      int port,
      PacketCallack callback,
      bool setAsDefaultOutgoing = true);
  void udpSend(InetAddress address, string_view data, int fd = -1);
  bool udpListenMulticast(InetAddress address, PacketCallack callback);
  void udpSendMulticast(InetAddress address, const std::string& data);
  int bindUdpSocket(InetAddress addr, bool reuse);

  // fd
  void bindCustomFd(int fd, std::function<void()> readyCallback);

  // Connect via TCP and return plain fd
  int connectTcpSocket(InetAddress addr);

  // TCP
  std::shared_ptr<FramedTcpConnection> tcpConnect(
      InetAddress address,
      TcpDataCallback dataCallback,
      TcpErrorCallback errorCallback);
  bool write(
      std::shared_ptr<FramedTcpConnection> conn,
      const std::string& data,
      bool queue);
  void close(std::shared_ptr<FramedTcpConnection> conn);

  // Execute pending callbacks. Remember to call it periodically
  void runOnce(int timeout);

}  // namespace OsSocket
