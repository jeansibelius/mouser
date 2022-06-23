# 🐭 Mouser

Mouser (mouse remotely): use a smartphone (or any device with a browser) to control your computer mouse.

### 🙈 Why?

Sometimes the computer is far away and clicking something would require getting up from the sofa. This simplifies that task.

## 🚴 How to run

- First build `go build`
- Then run `./mouser` (or install in gopath/path)
- Optionally provide flag `-addr [host:port]` to run at custom location (default 192.168.1.5:8080)
- Access the address above to control your computer mouse.\*

\*_Tested and used only on Ubuntu 22.04._

### 📎 Dependencies

- [Go](https://go.dev/)
- [Gorilla](https://github.com/gorilla/websocket)
- [libevdev](https://gitlab.freedesktop.org/libevdev/libevdev/-/tree/master/libevdev)

## 🤸 Notes/ideas

- Use smartphone movement (gyro/accleration) to simulate mouse movement?

## 📑 References

- https://gitlab.freedesktop.org/libevdev/libevdev/-/tree/master/
- https://github.com/gorilla/websocket/tree/master/examples/echo
- https://github.com/joshuar/gokbd
- https://www.instructables.com/Making-a-Joystick-With-HTML-pure-JavaScript/
