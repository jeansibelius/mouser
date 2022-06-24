let statusEl;

const postData = async (data) => {
  try {
    const response = await fetch("/mouse", {
      method: "POST",
      body: data,
      headers: {
        "Content-Type": "application/text",
      },
    });

    if (!response.ok) {
      statusEl.classList.add("notConnected");
      throw new Error(`Request failed with status ${response.status}`);
    }
    statusEl.classList.remove("notConnected");
  } catch (e) {
    statusEl.classList.add("notConnected");
    console.error(e.message);
  }
};

const parseCoordinates = (x, y) => {
  return String(x) + "," + String(y);
};

const send = async (x, y) => {
  const coord = parseCoordinates(x, y);
  await postData(coord);
};

// Joystick
let canvas, ctx;
let width, height, radius, xOrig, yOrig;

const resizeCanvas = () => {
  width = window.innerWidth;
  radius = width / 4 > 200 ? 200 : width / 4;
  height = window.innerHeight;
  ctx.canvas.width = width;
  ctx.canvas.height = height;
  joystickBg();
  joystickCircle(width / 2, height / 3);
};

const joystickBg = () => {
  xOrig = width / 2;
  yOrig = height / 3;
  ctx.beginPath();
  ctx.arc(xOrig, yOrig, radius + radius * 0.13, 0, Math.PI * 2, true);
  ctx.shadowColor = "transparent";
  ctx.fillStyle = "#ECE9E9";
  ctx.fill();
};

const joystickCircle = (width, height) => {
  ctx.beginPath();
  ctx.arc(width, height, radius, 0, Math.PI * 2, true);
  ctx.shadowColor = "rgba(0, 0, 0, 0.4)";
  ctx.shadowBlur = 6;
  ctx.shadowOffsetX = 6;
  ctx.shadowOffsetY = 6;
  ctx.fillStyle = "rgb(33, 0, 60)";
  ctx.fill();
  ctx.strokeStyle = "rgb(182, 155, 204)";
  ctx.lineWidth = 8;
  ctx.stroke();
};

const drawJoystick = (newWidth = width / 2, newHeight = height / 3) => {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  joystickBg();
  joystickCircle(newWidth, newHeight);
};

// Cursor coordinates
const coord = { x: 0, y: 0 };

const getPosition = (event) => {
  const cursorX = event.clientX || event.touches[0].clientX;
  const cursorY = event.clientY || event.touches[0].clientY;
  coord.x = cursorX - canvas.offsetLeft;
  coord.y = cursorY - canvas.offsetTop;
};

const isInCircle = () => {
  const current_radius = Math.sqrt(Math.pow(coord.x - xOrig, 2) + Math.pow(coord.y - yOrig, 2));
  if (radius >= current_radius) return true;
  return false;
};

// State of painting
let isPainting = false;

const startDrawing = (event) => {
  isPainting = true;
  getPosition(event);
  if (isInCircle()) {
    Draw(event);
  }
};

const stopDrawing = () => {
  isPainting = false;
  drawJoystick();
};

const Draw = (event) => {
  if (isPainting) {
    const { x, y } = coord;
    const xDelta = Math.round(x - xOrig);
    const yDelta = Math.round(y - yOrig);
    // Constrain max "speed" of mouse to make it more controllable
    const x_relative = Math.abs(xDelta) > 8 ? Math.sign(xDelta) * 8 : xDelta;
    const y_relative = Math.abs(yDelta) > 8 ? Math.sign(yDelta) * 8 : yDelta;

    send(x_relative, y_relative);

    // Redraw joystick
    getPosition(event);
    drawJoystick(x, y);
  }
};

window.addEventListener("load", () => {
  statusEl = document.getElementById("status");
  canvas = document.getElementById("canvas");
  ctx = canvas.getContext("2d");
  resizeCanvas();

  // Mouse events
  document.addEventListener("mousedown", startDrawing);
  document.addEventListener("mouseup", stopDrawing);
  document.addEventListener("mousemove", Draw);

  // Touch events
  document.addEventListener("touchstart", startDrawing);
  document.addEventListener("touchend", stopDrawing);
  document.addEventListener("touchcancel", stopDrawing);
  document.addEventListener("touchmove", Draw);

  // Mouse buttons
  document.getElementById("left").onclick = () => {
    postData("left");
  };
  document.getElementById("right").onclick = () => {
    postData("right");
  };

  // Set canvas size
  window.addEventListener("resize", resizeCanvas);
});
