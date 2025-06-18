import type { FlowerData } from "./types";

import { getDateText, getRandomFlower, changeVisibility } from "./utils/utils";
import { checkAuth, logout } from "./utils/auth";
import { getUserThoughts, updateVisibility } from "./utils/thought";

let initLeft = 200;

const h1Greet = document.getElementById("greet") as HTMLHeadingElement;
checkAuth().then((data) => {
  const isLoggedIn = data.success;

  if (!isLoggedIn && window.location.pathname !== "/auth.html") {
    window.location.href = "/auth.html";
    return;
  };

  h1Greet.textContent = `Howdy, ${data.data.name}`;
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
      h1Greet.style.opacity = "100%";
      
      // Set background opacity to low
      containerBackground.style.opacity = "50%";
      h1Greet.style.opacity = "50%";

      // Set opacity of all flowers to 50%, except the clicked one
      const allFlowers = container.getElementsByClassName("flower");
      for (let i = 0; i < allFlowers.length; i++) {
        const flowerEl = allFlowers[i] as HTMLDivElement;
        if (flowerEl === el) {
          flowerEl.style.opacity = "100%";
        } else {
          flowerEl.style.opacity = "50%";
        }
      }

      // Set opacity of all dates to 50%
      const allDates = container.getElementsByClassName("date")
      for (let i = 0; i < allDates.length; i++) {
        const date = allDates[i] as HTMLParagraphElement;
        date.style.opacity = "50%";
      }

      // Popup content
      flowerPopup.textContent = flower.thought;
      flowerPopup.style.display = "flex";
      flowerPopup.style.textTransform = "lowercase";
      flowerPopup.style.left = (parseInt(el.style.left, 10) + 140) + "px";
      
      const hrEl = document.createElement("hr");
      hrEl.style.marginTop = "4px";
      flowerPopup.appendChild(hrEl);

      // Create another div inside popup to show current
      const divEl = document.createElement("div");
      divEl.className = "popup-container";
      flowerPopup.appendChild(divEl);

      const buttonEl = document.createElement("button");
      buttonEl.textContent = `change to ${changeVisibility(flower.visibility)}`;
      buttonEl.style.cursor = "pointer";
      buttonEl.addEventListener("mouseenter", () => {
        buttonEl.style.textDecoration = "underline";
      });
      buttonEl.addEventListener("mouseleave", () => {
        buttonEl.style.textDecoration = "none";
      });
      buttonEl.addEventListener("click", async () => {
        await updateVisibility(flower.id, changeVisibility(flower.visibility));
        
        window.location.reload();
      });
      divEl.appendChild(buttonEl);
      
      const visibilityEl = document.createElement("p");
      visibilityEl.textContent = flower.visibility;
      divEl.appendChild(visibilityEl);
    });

    container.appendChild(el);

    const dateEl = document.createElement("p");
    dateEl.className = "date";
    dateEl.style.left = initLeft + 250 + "px";

    dateEl.textContent = getDateText(flower.createdAt);

    initLeft += 600;

    container.appendChild(dateEl);
  });
});

// Reset styles upon click
document.addEventListener("click", () => {
  flowerPopup.style.display = "none";
  containerBackground.style.opacity = "100%";
  h1Greet.style.opacity = "100%";

  const allFlowers = container.getElementsByClassName("flower");
  for (let i = 0; i < allFlowers.length; i++) {
    const flowerEl = allFlowers[i] as HTMLDivElement;
    flowerEl.style.opacity = "100%";
  }
  const allDates = container.getElementsByClassName("date")
  for (let i = 0; i < allDates.length; i++) {
    const date = allDates[i] as HTMLParagraphElement;
    date.style.opacity = "100%";
  }
});

const logoutBtn = document.getElementById("logout-btn") as HTMLButtonElement;
logoutBtn.addEventListener("click", async () => {
  logoutBtn.disabled = true;

  await logout();

  logoutBtn.disabled = false;

  window.location.href = "/auth.html";
});
