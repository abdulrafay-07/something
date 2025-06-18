import axios from "axios";

import type { ApiResponse } from "../types";

const baseUrl = "http://localhost:8080/api";

export async function getUserThoughts() {
  const response = await axios.get<ApiResponse>(`${baseUrl}/thought`, {
    withCredentials: true,
  });

  if (response.data.success) {
    return response.data.data;
  }
}

export async function updateVisibility(id: string, visibility: "public" | "private") {
  const response = await axios.patch<ApiResponse>(`${baseUrl}/thought/${id}`, {
    visibility,
  }, {
    withCredentials: true,
  });

  if (response.data.success) {
    return response.data
  }
}
