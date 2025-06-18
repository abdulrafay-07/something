import axios from "axios";
import type { ApiResponse } from "../types";

const baseUrl = "http://localhost:8080/api";

export async function checkAuth() {
  const response = await axios.get<ApiResponse>(`${baseUrl}/me`, {
    withCredentials: true,
  });

  return response.data;
}

export async function login(email: string, password: string) {
  const response = await axios.post<ApiResponse>(`${baseUrl}/signin`, {
    email,
    password,
  }, {
    withCredentials: true,
  });

  return response.data;
};

export async function signup(name: string, email: string, password: string) {
  const response = await axios.post<ApiResponse>(`${baseUrl}/signup`, {
    name,
    email,
    password,
  });

  return response.data;
}

export async function logout() {
  const response = await axios.post<ApiResponse>(`${baseUrl}/logout`, {}, { withCredentials: true });

  return response.data;
}
