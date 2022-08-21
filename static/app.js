let connectionStatusIndicator;

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
      connectionStatusIndicator.classList.add("notConnected");
      throw new Error(`Request failed with status ${response.status}`);
    }
    connectionStatusIndicator.classList.remove("notConnected");
  } catch (e) {
    connectionStatusIndicator.classList.add("notConnected");
    console.error(e.message);
  }
};

let pointerIsDown = false;
let previousPos = { x: 0, y: 0 };

const moveMouse = (currentPos, previousPos) => {
  const newPosition = `${currentPos.x - previousPos.x},${currentPos.y - previousPos.y}`;
  postData(newPosition);
};

const getPointerPosition = (pointerEvent) => {
  return { x: pointerEvent.clientX, y: pointerEvent.clientY };
};

const pointerDown = (pointerEvent) => {
  pointerIsDown = true;
  previousPos = getPointerPosition(pointerEvent); // Initialise position so we can move relative to it
};

const pointerUp = () => {
  pointerIsDown = false;
};

const pointerMove = (pointerEvent) => {
  if (!pointerIsDown) return;
  const currentPos = getPointerPosition(pointerEvent);
  moveMouse(currentPos, previousPos);
  previousPos = currentPos;
};

window.addEventListener("load", () => {
  connectionStatusIndicator = document.getElementById("status");

  document.addEventListener("pointerdown", pointerDown);
  document.addEventListener("pointerup", pointerUp);
  document.addEventListener("pointermove", pointerMove);

  // Mouse buttons
  document.getElementById("left").onclick = () => {
    postData("left");
  };
  document.getElementById("right").onclick = () => {
    postData("right");
  };
});
