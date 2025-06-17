import { checkAuth, logout } from "./utils/auth";

checkAuth().then((isLoggedIn) => {
  if (!isLoggedIn && window.location.pathname !== "/auth.html") {
    window.location.href = "/auth.html";
  }
});

const logoutBtn = document.getElementById("logout-btn") as HTMLButtonElement;

logoutBtn.addEventListener("click", async () => {
  logoutBtn.disabled = true;

  await logout();

  logoutBtn.disabled = false;

  window.location.href = "/auth.html";
});
