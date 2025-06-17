import { checkAuth, login, signup } from "./utils/auth";

checkAuth().then((isLoggedIn) => {
  if (isLoggedIn && window.location.pathname == "/auth.html") {
    window.location.href = "/";
  }
});

const loginTab = document.getElementById("login-tab") as HTMLButtonElement;
const signupTab = document.getElementById("signup-tab") as HTMLButtonElement;

const loginForm = document.getElementById("login-form") as HTMLFormElement;
const signupForm = document.getElementById("signup-form") as HTMLFormElement;

function setActiveTab(tab: "login" | "signup") {
  if (tab == "login") {
    loginTab.classList.add("active-tab");
    signupTab.classList.remove("active-tab");

    // Change the classes for the forms
    loginForm.classList.add("active-form");
    signupForm.classList.remove("active-form", "flex");
    signupForm.classList.add("hidden");
  } else {
    signupTab.classList.add("active-tab");
    loginTab.classList.remove("active-tab");

    signupForm.classList.add("active-form")
    loginForm.classList.remove("active-form", "flex");
    loginForm.classList.add("hidden");
  }
}

loginTab.addEventListener("click", () => {
  setActiveTab("login");
});

signupTab.addEventListener("click", () => {
  setActiveTab("signup");
});

loginForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const email = (document.getElementById("login-email") as HTMLInputElement).value;
  const password = (document.getElementById("login-password") as HTMLInputElement).value;
  const loginBtn = document.getElementById("signin-btn") as HTMLButtonElement;

  loginBtn.disabled = true;

  await login(email, password);

  loginBtn.disabled = false;

  window.location.href = "/";
})

signupForm.addEventListener("submit", async (event) => {
  event.preventDefault()

  const name = (document.getElementById("name") as HTMLInputElement).value;
  const email = (document.getElementById("signup-email") as HTMLInputElement).value;
  const password = (document.getElementById("signup-password") as HTMLInputElement).value;
  const signupBtn = document.getElementById("signup-btn") as HTMLButtonElement;

  signupBtn.disabled = true;

  await signup(name, email, password);

  signupBtn.disabled = false;

  window.location.href = "/";
})
