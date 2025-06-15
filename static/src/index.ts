import { checkAuth } from "./utils/auth";

checkAuth().then((isLoggedIn) => {
  if (!isLoggedIn && window.location.pathname !== "/auth.html") {
    window.location.href = "/auth.html";
  }
});
