import type { FlowerData } from "./types";

import { getDateText, getRandomFlower } from "./utils/utils";
import { checkAuth, logout } from "./utils/auth";
import { getUserThoughts } from "./utils/thought";

let initLeft = 200;

checkAuth().then((isLoggedIn) => {
  if (!isLoggedIn && window.location.pathname !== "/auth.html") {
    window.location.href = "/auth.html";
  };
});

const container = document.getElementById("main-container") as HTMLDivElement;
const containerBackground = document.getElementsByClassName("container-background")[0] as HTMLDivElement;
const flowerPopup = document.getElementById("flower-popup") as HTMLDivElement;

getUserThoughts().then((data) => {
  data.thoughts.forEach((flower: FlowerData) => {
    const el = document.createElement("div");
    el.className = `flower ${getRandomFlower()}`;
    el.style.left = initLeft + "px";

    el.addEventListener("click", (e) => {
      e.stopPropagation();

      // Set previous popup display to none
      flowerPopup.style.display = "none";
      containerBackground.style.opacity = "100%";
      
      // Popup content
      flowerPopup.textContent = flower.thought;
      flowerPopup.style.display = "block";
      flowerPopup.style.textTransform = "lowercase";
      flowerPopup.style.left = (parseInt(el.style.left, 10) + 140) + "px";

      // Set background opacity to low
      containerBackground.style.opacity = "50%";
    });

    container.appendChild(el);

    const dateEl = document.createElement("p");
    dateEl.className = "date";
    dateEl.style.left = initLeft + 170 + "px";

    dateEl.textContent = getDateText(flower.createdAt);

    initLeft += 600;

    container.appendChild(dateEl);
  });
});

document.addEventListener("click", () => {
  flowerPopup.style.display = "none";
  containerBackground.style.opacity = "100%";
});

const logoutBtn = document.getElementById("logout-btn") as HTMLButtonElement;

logoutBtn.addEventListener("click", async () => {
  logoutBtn.disabled = true;

  await logout();

  logoutBtn.disabled = false;

  window.location.href = "/auth.html";
});
