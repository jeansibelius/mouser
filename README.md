# 🐭 Mouser

Mouser (mouse remotely): use a smartphone to control your computer mouse.

### 🙈 Why?

Sometimes the computer is far away and clicking something would require getting up from the sofa. This simplifies that task.

## 🚴 How to run

- First build `go build`
- Then run `./mouser` (or install in gopath/path)
- Optionally provide flag `-addr [host:port]` to run at custom location (default 192.168.1.5:8080)

## 🤸 Notes/ideas

- Use smartphone movement (gyro/accleration) to simulate mouse movement?

## 📑 References

- https://gitlab.freedesktop.org/libevdev/libevdev/-/tree/master/
- https://github.com/gorilla/websocket/tree/master/examples/echo
- https://github.com/joshuar/gokbd
