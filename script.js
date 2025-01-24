const loginBtn = document.getElementById("login-btn");
const signupBtn = document.getElementById("signup-btn");
const loginForm = document.getElementById("login-form");
const signupForm = document.getElementById("signup-form");

loginBtn.addEventListener("click", () => {
  loginForm.classList.add("active");
  signupForm.classList.remove("active");
  loginBtn.classList.add("active");
  signupBtn.classList.remove("active");
});

signupBtn.addEventListener("click", () => {
  signupForm.classList.add("active");
  loginForm.classList.remove("active");
  signupBtn.classList.add("active");
  loginBtn.classList.remove("active");
});
