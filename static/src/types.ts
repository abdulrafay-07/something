export interface ApiResponse {
  success: boolean;
  message: string;
  data?: any
};

export interface FlowerData {
  id: string;
  userId: number;
  thought: string;
  createdAt: string;
  visibility: "public" | "private";
};
