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

const listenMouseMove = (onImmediate, onFrame) => {
  let isDown = false;
  let current = { x: 0, y: 0 };
  let previous = { x: 0, y: 0 };

  const setCurrent = (event) => {
    current = {
      x: event.clientX || event.touches[0].clientX,
      y: event.clientY || event.touches[0].clientY,
    };
  };

  const checkRedraw = () => {
    if (current !== previous) {
      const prev = previous;
      previous = current;
      onFrame(current, prev);
    }
    if (isDown) {
      window.requestAnimationFrame(checkRedraw);
    }
  };

  const onStart = (event) => {
    isDown = true;
    setCurrent(event);
    onMove(event);
    checkRedraw();
  };

  const onEnd = () => {
    isDown = false;
  };

  const onMove = (event) => {
    if (isDown) {
      previous = current;
      setCurrent(event);
      onImmediate(current, previous);
    }
  };

  document.addEventListener("mousedown", onStart);
  document.addEventListener("mouseup", onEnd);
  document.addEventListener("mousemove", onMove);

  // Touch events
  document.addEventListener("touchstart", onStart);
  document.addEventListener("touchend", onEnd);
  document.addEventListener("touchcancel", onEnd);
  document.addEventListener("touchmove", onMove);
  return () => {
    document.removeEventListener("mousedown", onStart);
    document.removeEventListener("mouseup", onEnd);
    document.removeEventListener("mousemove", onMove);

    // Touch events
    document.removeEventListener("touchstart", onStart);
    document.removeEventListener("touchend", onEnd);
    document.removeEventListener("touchcancel", onEnd);
    document.removeEventListener("touchmove", onMove);
  };
};

const unlisten = listenMouseMove(
  (current, previous) => {
    postData(`${current.x - previous.x},${current.y - previous.y}`);
  },
  (current, previous) => {
    postData(`${current.x - previous.x},${current.y - previous.y}`);
  }
);

window.addEventListener("load", () => {
  statusEl = document.getElementById("status");

  // Mouse buttons
  document.getElementById("left").onclick = () => {
    postData("left");
  };
  document.getElementById("right").onclick = () => {
    postData("right");
  };
});
